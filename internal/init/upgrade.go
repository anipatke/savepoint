package init

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
)

type UpgradeAction string

const (
	ActionUpdated   UpgradeAction = "updated"
	ActionInstalled UpgradeAction = "installed"
	ActionMerged    UpgradeAction = "merged"
	ActionMigrated  UpgradeAction = "migrated"
	ActionUnchanged UpgradeAction = "unchanged"
	ActionSkipped   UpgradeAction = "skipped"
)

type UpgradeReport struct {
	Actions []UpgradeEntry
}

type UpgradeEntry struct {
	Path   string
	Action UpgradeAction
}

func (r *UpgradeReport) Format() string {
	if len(r.Actions) == 0 {
		return "No assets to upgrade."
	}

	var updated, installed, merged, migrated, unchanged, skipped int
	for _, e := range r.Actions {
		switch e.Action {
		case ActionUpdated:
			updated++
		case ActionInstalled:
			installed++
		case ActionMerged:
			merged++
		case ActionMigrated:
			migrated++
		case ActionUnchanged:
			unchanged++
		case ActionSkipped:
			skipped++
		}
	}

	var b strings.Builder
	b.WriteString("Upgrade Report:\n")
	if updated > 0 {
		fmt.Fprintf(&b, "  Updated: %d\n", updated)
	}
	if installed > 0 {
		fmt.Fprintf(&b, "  Installed: %d\n", installed)
	}
	if merged > 0 {
		fmt.Fprintf(&b, "  Merged: %d\n", merged)
	}
	if migrated > 0 {
		fmt.Fprintf(&b, "  Migrated: %d\n", migrated)
	}
	if unchanged > 0 {
		fmt.Fprintf(&b, "  Unchanged: %d\n", unchanged)
	}
	if skipped > 0 {
		fmt.Fprintf(&b, "  Skipped: %d\n", skipped)
	}

	for _, e := range r.Actions {
		fmt.Fprintf(&b, "  %s  %s\n", e.Action, e.Path)
	}

	return b.String()
}

func UpgradeProjectAssets(templates fs.FS, targetDir string, dryRun, force bool) (*UpgradeReport, error) {
	absTarget, err := filepath.Abs(targetDir)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve target directory: %w", err)
	}

	info, err := os.Stat(absTarget)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("target directory %q does not exist", targetDir)
		}
		return nil, fmt.Errorf("cannot access %q: %w", targetDir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("target %q is not a directory", targetDir)
	}

	savepointDir := filepath.Join(absTarget, ".savepoint")
	if _, err := os.Stat(savepointDir); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("target %q is not a Savepoint project: no .savepoint directory", targetDir)
		}
		return nil, fmt.Errorf("cannot check .savepoint directory: %w", err)
	}

	var report UpgradeReport

	// Retire the legacy generic audit skill before installing the split skills,
	// so an interrupted upgrade never leaves the old alias triggerable next to
	// its replacements.
	migration, err := migrateLegacyAuditSkill(absTarget, dryRun)
	if err != nil {
		return nil, err
	}
	if migration != nil {
		report.Actions = append(report.Actions, *migration)
	}

	err = fs.WalkDir(templates, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walk error at %s: %w", path, err)
		}
		if d.IsDir() {
			return nil
		}
		if path == "." {
			return nil
		}

		if installAction, ok := installMissingAction(path); ok {
			action, err := installMissingAsset(absTarget, templates, path, installAction, dryRun)
			if err != nil {
				return err
			}
			report.Actions = append(report.Actions, UpgradeEntry{Path: path, Action: action})
			return nil
		}

		if strings.HasPrefix(path, ".savepoint") {
			report.Actions = append(report.Actions, UpgradeEntry{Path: path, Action: ActionSkipped})
			return nil
		}

		targetPath := filepath.Join(absTarget, path)

		isSkill := isPackageSkillAsset(path)

		if !isSkill && path != "AGENTS.md" {
			report.Actions = append(report.Actions, UpgradeEntry{Path: path, Action: ActionSkipped})
			return nil
		}

		content, err := fs.ReadFile(templates, path)
		if err != nil {
			return fmt.Errorf("read template %s: %w", path, err)
		}

		if dryRun {
			if path == "AGENTS.md" {
				dest := FindAgentGuide(absTarget)
				if dest == "" {
					report.Actions = append(report.Actions, UpgradeEntry{Path: path, Action: ActionUpdated})
					return nil
				}

				existingContent, err := os.ReadFile(dest)
				if err != nil {
					return fmt.Errorf("read existing %s: %w", path, err)
				}
				block := managedBegin + "\n" + strings.TrimSpace(string(content)) + "\n" + managedEnd
				merged := replaceManagedBlock(string(existingContent), block)
				if merged == string(existingContent) {
					report.Actions = append(report.Actions, UpgradeEntry{Path: path, Action: ActionUnchanged})
				} else {
					report.Actions = append(report.Actions, UpgradeEntry{Path: path, Action: ActionMerged})
				}
				return nil
			}

			if _, err := os.Stat(targetPath); err != nil {
				if os.IsNotExist(err) {
					report.Actions = append(report.Actions, UpgradeEntry{Path: path, Action: ActionUpdated})
				} else {
					report.Actions = append(report.Actions, UpgradeEntry{Path: path, Action: ActionSkipped})
				}
			} else {
				existingContent, err := os.ReadFile(targetPath)
				if err != nil {
					return fmt.Errorf("read existing %s: %w", path, err)
				}
				if string(existingContent) == string(content) {
					report.Actions = append(report.Actions, UpgradeEntry{Path: path, Action: ActionUnchanged})
				} else {
					report.Actions = append(report.Actions, UpgradeEntry{Path: path, Action: ActionUpdated})
				}
			}
			return nil
		}

		if isSkill {
			existingContent, err := os.ReadFile(targetPath)
			if os.IsNotExist(err) {
				if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
					return fmt.Errorf("create parent dir for %s: %w", path, err)
				}
				if err := AtomicWrite(targetPath, content); err != nil {
					return fmt.Errorf("write %s: %w", path, err)
				}
				report.Actions = append(report.Actions, UpgradeEntry{Path: path, Action: ActionUpdated})
				return nil
			}
			if err != nil {
				return fmt.Errorf("read existing %s: %w", path, err)
			}

			if string(existingContent) == string(content) {
				report.Actions = append(report.Actions, UpgradeEntry{Path: path, Action: ActionUnchanged})
				return nil
			}

			if err := AtomicWrite(targetPath, content); err != nil {
				return fmt.Errorf("write %s: %w", path, err)
			}
			report.Actions = append(report.Actions, UpgradeEntry{Path: path, Action: ActionUpdated})
			return nil
		}

		if path == "AGENTS.md" {
			dest := FindAgentGuide(absTarget)
			if dest == "" {
				dest = targetPath
			}

			existingContent, readErr := os.ReadFile(dest)
			hasExisting := readErr == nil

			if hasExisting {
				hasBlock := strings.Contains(string(existingContent), managedBegin)
				if hasBlock {
					rendered := string(content)
					merged := replaceManagedBlock(string(existingContent), managedBegin+"\n"+strings.TrimSpace(rendered)+"\n"+managedEnd)
					if merged == string(existingContent) {
						report.Actions = append(report.Actions, UpgradeEntry{Path: path, Action: ActionUnchanged})
						return nil
					}
				}
			}

			if err := MergeAgentGuide(dest, string(content)); err != nil {
				return fmt.Errorf("merge agent guide %s: %w", path, err)
			}
			report.Actions = append(report.Actions, UpgradeEntry{Path: path, Action: ActionMerged})
			return nil
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(report.Actions, func(i, j int) bool {
		return report.Actions[i].Path < report.Actions[j].Path
	})

	return &report, nil
}

// isPackageSkillAsset reports whether a template path is a package-owned skill
// asset that upgrade refreshes in place: a skill entrypoint, or a shared
// reference such as agent-skills/references/audit-method.md that the split
// audit skills load but that never triggers on its own.
func isPackageSkillAsset(path string) bool {
	if !strings.HasPrefix(path, "agent-skills/") {
		return false
	}
	return strings.HasSuffix(path, "/SKILL.md") || strings.HasPrefix(path, "agent-skills/references/")
}

// auditAssetPrefix marks the project documentation area whose files are
// user-maintained audit state. Upgrade may add a missing scaffold file here,
// but must never overwrite an existing prompt, register, finding, or run file.
const auditAssetPrefix = ".savepoint/audit/"

func isAuditAsset(path string) bool {
	return strings.HasPrefix(path, auditAssetPrefix)
}

// policyAssets are the project-owned policy files upgrade delivers to an
// existing project so guidance that references them resolves after an upgrade.
// The allowlist is exact: every other `.savepoint/` path stays skipped.
var policyAssets = []string{
	".savepoint/Guardrails.md",
	".savepoint/Health-Check.md",
}

func isPolicyAsset(path string) bool {
	return slices.Contains(policyAssets, path)
}

// installMissingAction reports whether a template path is install-if-missing,
// and the action to record when the file is added. Audit scaffold files keep
// reporting ActionUpdated; policy files report ActionInstalled so the upgrade
// report distinguishes a newly delivered policy asset from a refreshed one.
func installMissingAction(path string) (UpgradeAction, bool) {
	switch {
	case isAuditAsset(path):
		return ActionUpdated, true
	case isPolicyAsset(path):
		return ActionInstalled, true
	default:
		return "", false
	}
}

// installMissingAsset adds a missing install-if-missing asset and reports
// installAction. An existing file is left untouched and reported
// ActionUnchanged so user-edited content is never overwritten.
func installMissingAsset(absTarget string, templates fs.FS, path string, installAction UpgradeAction, dryRun bool) (UpgradeAction, error) {
	targetPath := filepath.Join(absTarget, path)

	if _, err := os.Stat(targetPath); err == nil {
		return ActionUnchanged, nil
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("stat %s: %w", path, err)
	}

	if dryRun {
		return installAction, nil
	}

	content, err := fs.ReadFile(templates, path)
	if err != nil {
		return "", fmt.Errorf("read template %s: %w", path, err)
	}
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return "", fmt.Errorf("create parent dir for %s: %w", path, err)
	}
	rendered := interpolate(string(content), ProjectNameFromDir(absTarget))
	if err := AtomicWrite(targetPath, []byte(rendered)); err != nil {
		return "", fmt.Errorf("write %s: %w", path, err)
	}
	return installAction, nil
}
