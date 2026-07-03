package init

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type UpgradeAction string

const (
	ActionUpdated   UpgradeAction = "updated"
	ActionMerged    UpgradeAction = "merged"
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

	var updated, merged, unchanged, skipped int
	for _, e := range r.Actions {
		switch e.Action {
		case ActionUpdated:
			updated++
		case ActionMerged:
			merged++
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
	if merged > 0 {
		fmt.Fprintf(&b, "  Merged: %d\n", merged)
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

		if isAuditAsset(path) {
			action, err := upgradeAuditAsset(absTarget, templates, path, dryRun)
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

		isSkill := strings.HasPrefix(path, "agent-skills/") && strings.HasSuffix(path, "/SKILL.md")

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

// auditAssetPrefix marks the project documentation area whose files are
// user-maintained audit state. Upgrade may add a missing scaffold file here,
// but must never overwrite an existing prompt, register, finding, or run file.
const auditAssetPrefix = ".savepoint/audit/"

func isAuditAsset(path string) bool {
	return strings.HasPrefix(path, auditAssetPrefix)
}

// upgradeAuditAsset adds a missing audit-register scaffold file and reports
// ActionUpdated. An existing file is left untouched and reported ActionUnchanged
// so user-edited audit state is never overwritten.
func upgradeAuditAsset(absTarget string, templates fs.FS, path string, dryRun bool) (UpgradeAction, error) {
	targetPath := filepath.Join(absTarget, path)

	if _, err := os.Stat(targetPath); err == nil {
		return ActionUnchanged, nil
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("stat %s: %w", path, err)
	}

	if dryRun {
		return ActionUpdated, nil
	}

	content, err := fs.ReadFile(templates, path)
	if err != nil {
		return "", fmt.Errorf("read template %s: %w", path, err)
	}
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return "", fmt.Errorf("create parent dir for %s: %w", path, err)
	}
	if err := AtomicWrite(targetPath, content); err != nil {
		return "", fmt.Errorf("write %s: %w", path, err)
	}
	return ActionUpdated, nil
}
