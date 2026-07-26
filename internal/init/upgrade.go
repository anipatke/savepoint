package init

import (
	"bytes"
	"errors"
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
	// ActionConflict reports work deliberately not done: the file on disk was
	// changed by the user, so it is kept and the incoming version is written
	// beside it. It is the one action that needs the user to act.
	ActionConflict UpgradeAction = "conflict"
	// ActionFailed reports work that was attempted and did not complete. It
	// exists so a write failure still names its path and any sidecar already
	// written for it, rather than vanishing behind the returned error.
	ActionFailed UpgradeAction = "failed"
)

// Sidecar suffixes for content upgrade must not silently destroy: the incoming
// version of a conflicted file, and the previous content of a replaced one.
const (
	incomingSuffix = ".new"
	backupSuffix   = ".bak"
)

// Per-entry notes explaining a sidecar file, so the report names the recovery
// path rather than leaving the user to find it.
const (
	noteConflict = "kept your version; incoming written to " + incomingSuffix
	noteBackup   = "previous content saved to " + backupSuffix
)

type UpgradeReport struct {
	Actions []UpgradeEntry
}

type UpgradeEntry struct {
	Path   string
	Action UpgradeAction
	// Note names a sidecar file written alongside the target, empty when the
	// action needed none.
	Note string
}

func (r *UpgradeReport) Format() string {
	if len(r.Actions) == 0 {
		return "No assets to upgrade."
	}

	var updated, installed, merged, migrated, unchanged, skipped, conflicts, failed int
	for _, e := range r.Actions {
		switch e.Action {
		case ActionFailed:
			failed++
		case ActionConflict:
			conflicts++
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
	// Failures and conflicts lead the summary: they are the counts that need
	// the user to do something before the upgrade is complete.
	if failed > 0 {
		fmt.Fprintf(&b, "  Failed: %d\n", failed)
	}
	if conflicts > 0 {
		fmt.Fprintf(&b, "  Conflicts: %d\n", conflicts)
	}
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
		if e.Note != "" {
			fmt.Fprintf(&b, "  %s  %s  (%s)\n", e.Action, e.Path, e.Note)
			continue
		}
		fmt.Fprintf(&b, "  %s  %s\n", e.Action, e.Path)
	}

	return b.String()
}

// assetWriter performs one upgrade write. Production passes AtomicWrite; tests
// substitute a writer that fails at a chosen call, which is the only way to
// prove what a project looks like after a failure part-way through an upgrade.
type assetWriter func(path string, content []byte) error

func UpgradeProjectAssets(templates fs.FS, targetDir string, dryRun, force bool) (*UpgradeReport, error) {
	return upgradeProjectAssets(templates, targetDir, dryRun, force, AtomicWrite)
}

func upgradeProjectAssets(templates fs.FS, targetDir string, dryRun, force bool, write assetWriter) (*UpgradeReport, error) {
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

	manifest, err := LoadManifest(absTarget)
	if err != nil {
		return nil, err
	}

	// The manifest is the last thing an upgrade persists. Proving it can be
	// written turns an unwritable .savepoint/ into a clean refusal with nothing
	// changed, instead of an error arriving after skills and the agent guide
	// have already been rewritten. The check runs immediately before the first
	// write and only then: an upgrade that changes nothing must touch nothing,
	// and a probe file would be an observable change to .savepoint/ itself.
	write = beforeFirstWrite(write, func() error { return ensureManifestWritable(absTarget) })

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
			action, err := installMissingAsset(absTarget, templates, path, installAction, dryRun, write)
			if err != nil {
				report.Actions = append(report.Actions, UpgradeEntry{Path: path, Action: ActionFailed})
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
			report.Actions = append(report.Actions, UpgradeEntry{Path: path, Action: ActionFailed})
			return fmt.Errorf("read template %s: %w", path, err)
		}

		// A failed entry is appended too, and carries whatever note its helper
		// set: a backup written before the replacement that failed has to reach
		// the user, or the command creates a .bak without saying so.
		if isSkill {
			entry, err := upgradeSkillAsset(targetPath, path, content, manifest, dryRun, force, write)
			report.Actions = append(report.Actions, entry)
			if err != nil {
				return err
			}
			return nil
		}

		entry, err := upgradeAgentGuide(absTarget, path, content, dryRun, force, write)
		report.Actions = append(report.Actions, entry)
		if err != nil {
			return err
		}
		return nil
	})
	sort.Slice(report.Actions, func(i, j int) bool {
		return report.Actions[i].Path < report.Actions[j].Path
	})

	// A failure part-way through leaves the assets already written on disk.
	// Record their provenance anyway and hand the partial report back with the
	// error, so the user sees what was applied and the next upgrade still knows
	// which files are ours. The write failure stays the primary error; a
	// manifest failure behind it is secondary context.
	if err != nil {
		if dryRun {
			return &report, err
		}
		return &report, errors.Join(err, manifest.save(absTarget, write))
	}

	// A dry run reports only; the manifest is state, so it stays untouched.
	if !dryRun {
		if err := manifest.save(absTarget, write); err != nil {
			return &report, err
		}
	}

	return &report, nil
}

// beforeFirstWrite returns a writer that runs guard once, immediately before
// the first write it is asked to perform. A run with no writes never runs the
// guard, so a no-op upgrade stays a no-op on disk.
func beforeFirstWrite(write assetWriter, guard func() error) assetWriter {
	guarded := false
	return func(path string, content []byte) error {
		if !guarded {
			if err := guard(); err != nil {
				return err
			}
			guarded = true
		}
		return write(path, content)
	}
}

// ensureManifestWritable reports whether the manifest can be committed, without
// leaving anything behind.
func ensureManifestWritable(absTarget string) error {
	dir := filepath.Dir(filepath.Join(absTarget, manifestRelPath))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create parent dir for the upgrade manifest: %w", err)
	}

	probe, err := os.CreateTemp(dir, ".tmp-*.manifest-probe")
	if err != nil {
		return fmt.Errorf("cannot write the upgrade manifest in %s: %w", dir, err)
	}
	return discardProbe(probe)
}

// discardProbe closes and removes the writability probe. Both results are
// checked: a probe left behind is litter inside the user's project, and a close
// failure is how a full or failing filesystem shows up on a file whose write
// otherwise looked fine.
func discardProbe(probe *os.File) error {
	name := probe.Name()
	if err := errors.Join(probe.Close(), os.Remove(name)); err != nil {
		return fmt.Errorf("clean up manifest probe %s: %w", name, err)
	}
	return nil
}

// upgradeAgentGuide applies the AGENTS.md upgrade policy: refresh the managed
// block in place when the marker pair is present, and otherwise keep the file
// byte-identical and offer the merged result beside it.
//
// Appending to an unmarked guide would leave two sets of workflow instructions
// in the one file agents read — the user's above, never refreshed again, and
// the managed block below — drifting apart with every release. Which half is
// Savepoint-owned cannot be inferred without pattern-matching prose, and a
// wrong guess wraps the user's own text in markers the next upgrade overwrites.
// So the unmarked file conflicts, and adopting it stays an explicit --force.
//
// A dry run takes exactly the same decisions and writes nothing.
func upgradeAgentGuide(absTarget, path string, content []byte, dryRun, force bool, write assetWriter) (UpgradeEntry, error) {
	entry := UpgradeEntry{Path: path}

	// The guide keeps whatever casing it has on disk, so every write and every
	// sidecar targets the found name rather than the canonical one.
	dest := FindAgentGuide(absTarget)
	if dest == "" {
		dest = filepath.Join(absTarget, path)
	}

	block := managedBegin + "\n" + strings.TrimSpace(string(content)) + "\n" + managedEnd

	existing, err := os.ReadFile(dest)
	if os.IsNotExist(err) {
		entry.Action = ActionUpdated
		if dryRun {
			return entry, nil
		}
		if err := write(dest, []byte(block+"\n")); err != nil {
			return failedEntry(entry, ""), fmt.Errorf("write %s: %w", path, err)
		}
		return entry, nil
	}
	if err != nil {
		return failedEntry(entry, ""), fmt.Errorf("read existing %s: %w", path, err)
	}

	merged, marked := replaceManagedBlock(string(existing), block)

	switch {
	case marked && merged == string(existing):
		entry.Action = ActionUnchanged
		return entry, nil
	case marked:
		entry.Action = ActionMerged
	case force:
		entry.Action = ActionMerged
		entry.Note = noteBackup
	default:
		entry.Action = ActionConflict
		entry.Note = noteConflict
		if err := writeSidecar(dest, path, incomingSuffix, []byte(merged), dryRun, write); err != nil {
			return failedEntry(entry, ""), err
		}
		return entry, nil
	}

	if entry.Note == noteBackup {
		if err := writeSidecar(dest, path, backupSuffix, existing, dryRun, write); err != nil {
			return failedEntry(entry, ""), err
		}
	}
	if dryRun {
		return entry, nil
	}
	if err := write(dest, []byte(merged)); err != nil {
		// The backup, if this path wrote one, is on disk and stays named.
		return failedEntry(entry, entry.Note), fmt.Errorf("write %s: %w", path, err)
	}
	return entry, nil
}

// failedEntry marks an attempted action as failed, keeping only the note that
// is still true: a sidecar already written stays named, and a note describing
// work that never happened is dropped.
func failedEntry(entry UpgradeEntry, note string) UpgradeEntry {
	entry.Action = ActionFailed
	entry.Note = note
	return entry
}

// upgradeSkillAsset applies the skill upgrade policy for one template path:
// replace a file only when it is provably unmodified, and otherwise keep the
// user's version and offer the incoming one beside it. The manifest supplies
// the provenance that makes "the user edited this" distinguishable from "this
// is the previous version's copy"; template comparison alone cannot.
//
// A dry run takes exactly the same decisions and writes nothing.
func upgradeSkillAsset(targetPath, path string, content []byte, manifest *Manifest, dryRun, force bool, write assetWriter) (UpgradeEntry, error) {
	entry := UpgradeEntry{Path: path}

	existing, err := os.ReadFile(targetPath)
	if os.IsNotExist(err) {
		entry.Action = ActionUpdated
		if err := writeSkillAsset(targetPath, path, content, manifest, dryRun, write); err != nil {
			return failedEntry(entry, ""), err
		}
		return entry, nil
	}
	if err != nil {
		return failedEntry(entry, ""), fmt.Errorf("read existing %s: %w", path, err)
	}

	// Cheapest test first: identical content needs no hashing and no write,
	// only the provenance record a pre-manifest project still lacks.
	if bytes.Equal(existing, content) {
		if !dryRun {
			manifest.Record(path, existing)
		}
		entry.Action = ActionUnchanged
		return entry, nil
	}

	// Shared references under agent-skills/references/ are package-owned and
	// outside the manifest, so they refresh directly.
	if !isManifestPath(path) {
		entry.Action = ActionUpdated
		if err := writeSkillAsset(targetPath, path, content, manifest, dryRun, write); err != nil {
			return failedEntry(entry, ""), err
		}
		return entry, nil
	}

	recorded, tracked := manifest.Hash(path)
	switch {
	case !tracked:
		// Pre-manifest project: provenance is unknown for every skill, and
		// conflicting on all of them would bury the common untouched-but-
		// outdated case in noise. Back the old copy up instead — one time,
		// recoverable, and every later upgrade has exact provenance.
		entry.Action = ActionUpdated
		entry.Note = noteBackup
	case recorded == hashContent(existing):
		// Ours, merely outdated.
		entry.Action = ActionUpdated
	case force:
		entry.Action = ActionUpdated
		entry.Note = noteBackup
	default:
		entry.Action = ActionConflict
		entry.Note = noteConflict
	}

	if entry.Action == ActionConflict {
		// Writing to a fixed sidecar path replaces any stale .new from an
		// earlier upgrade rather than accumulating variants.
		if err := writeSidecar(targetPath, path, incomingSuffix, content, dryRun, write); err != nil {
			return failedEntry(entry, ""), err
		}
		return entry, nil
	}

	// The backup goes first, so a failure in the replacement that follows can
	// only ever leave the user's content still recoverable beside it.
	if entry.Note == noteBackup {
		if err := writeSidecar(targetPath, path, backupSuffix, existing, dryRun, write); err != nil {
			return failedEntry(entry, ""), err
		}
	}
	if err := writeSkillAsset(targetPath, path, content, manifest, dryRun, write); err != nil {
		// The backup, if this path wrote one, is on disk and stays named: the
		// user must be told where their content went, not just that a write
		// failed.
		return failedEntry(entry, entry.Note), err
	}
	return entry, nil
}

// writeSkillAsset writes a skill file and records its provenance, or does
// nothing on a dry run.
func writeSkillAsset(targetPath, path string, content []byte, manifest *Manifest, dryRun bool, write assetWriter) error {
	if dryRun {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return fmt.Errorf("create parent dir for %s: %w", path, err)
	}
	if err := write(targetPath, content); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	manifest.Record(path, content)
	return nil
}

// writeSidecar writes content next to targetPath under suffix, or does nothing
// on a dry run. Sidecars carry no provenance: they are never the live file.
func writeSidecar(targetPath, path, suffix string, content []byte, dryRun bool, write assetWriter) error {
	if dryRun {
		return nil
	}
	if err := write(targetPath+suffix, content); err != nil {
		return fmt.Errorf("write %s%s: %w", path, suffix, err)
	}
	return nil
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
func installMissingAsset(absTarget string, templates fs.FS, path string, installAction UpgradeAction, dryRun bool, write assetWriter) (UpgradeAction, error) {
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
	if err := write(targetPath, []byte(rendered)); err != nil {
		return "", fmt.Errorf("write %s: %w", path, err)
	}
	return installAction, nil
}
