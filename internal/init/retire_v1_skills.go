package init

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// retiredV1SkillDirs names the nine V1 skills a project loses when it upgrades
// on the V2 workflow: the four V2 skills replace what these did, and shipping
// both leaves two routing vocabularies an agent could read first.
var retiredV1SkillDirs = []string{
	"savepoint-draft-prd",
	"savepoint-create-plan",
	"savepoint-system-design",
	"savepoint-create-task",
	"savepoint-build-task",
	"savepoint-audit-task",
	"savepoint-audit-epic",
	"savepoint-audit-register",
	"savepoint-create-defect",
}

// retiredV1AuditMethod is the shared, non-triggerable reference the retired V1
// audit skills loaded. It retires alongside them: nothing in the V2 tree reads
// it, and keeping it would be a dangling reference no skill points at anymore.
const retiredV1AuditMethod = "agent-skills/references/audit-method.md"

// retiredV1AssetPaths lists every path retirement considers, project-relative
// with forward slashes: one SKILL.md per retired skill, plus the shared audit
// method reference.
func retiredV1AssetPaths() []string {
	paths := make([]string, 0, len(retiredV1SkillDirs)+1)
	for _, skill := range retiredV1SkillDirs {
		paths = append(paths, "agent-skills/"+skill+"/SKILL.md")
	}
	return append(paths, retiredV1AuditMethod)
}

// retireV1Skills retires every V1-era skill and the shared audit method
// reference from a V2 project, reusing the archive-then-delete order
// migrateLegacyAuditSkill established: archive the content under
// .savepoint/migrations/, verify it landed, then remove the triggerable copy,
// and only then drop its manifest entry. A path the project does not have —
// already retired, or a project that never carried it — is silently skipped,
// so a repeat run reports nothing and writes nothing.
//
// It returns every entry retirement produced before a failure, including the
// failed entry itself, so a partial run still names what happened in the
// report even though the error stops it from going further.
func retireV1Skills(absTarget string, manifest *Manifest, dryRun bool, write assetWriter) ([]UpgradeEntry, error) {
	var entries []UpgradeEntry

	for _, path := range retiredV1AssetPaths() {
		entry, err := retireV1Asset(absTarget, path, dryRun, write)
		if entry != nil {
			entries = append(entries, *entry)
		}
		if err != nil {
			return entries, err
		}
		if entry != nil && !dryRun {
			manifest.Forget(path)
		}
	}

	return entries, nil
}

// retireV1Asset retires one path and returns nil, nil when the project never
// had it — already retired, or a fresh V2 project that never carried V1
// assets — so callers report nothing for a path with nothing to retire.
func retireV1Asset(absTarget, path string, dryRun bool, write assetWriter) (*UpgradeEntry, error) {
	targetPath := filepath.Join(absTarget, filepath.FromSlash(path))

	content, err := os.ReadFile(targetPath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read retired asset %s: %w", path, err)
	}

	archivePath, err := resolveArchivePath(absTarget, archiveStem(path), content)
	if err != nil {
		return &UpgradeEntry{Path: path, Action: ActionFailed}, err
	}

	entry := &UpgradeEntry{Path: path, Action: ActionRetired}
	if dryRun {
		return entry, nil
	}

	// Preserve first, delete second: the triggerable copy is removed only after
	// its content is safely on disk somewhere else.
	if archivePath != "" {
		if err := writeArchive(absTarget, archivePath, content, write); err != nil {
			return &UpgradeEntry{Path: path, Action: ActionFailed}, fmt.Errorf("archive retired asset %s: %w", path, err)
		}
	}

	if err := os.Remove(targetPath); err != nil {
		return &UpgradeEntry{Path: path, Action: ActionFailed}, fmt.Errorf("remove retired asset %s: %w", path, err)
	}
	removeDirIfEmpty(filepath.Dir(targetPath))

	return entry, nil
}

// archiveStem derives a migrations-archive file stem from a template-relative
// path: agent-skills/savepoint-build-task/SKILL.md becomes
// savepoint-build-task-SKILL, and agent-skills/references/audit-method.md
// becomes references-audit-method. The scheme matches legacyAuditArchiveStem,
// which is this function's rule applied by hand to one fixed path.
func archiveStem(path string) string {
	rel := strings.TrimSuffix(strings.TrimPrefix(path, "agent-skills/"), ".md")
	return strings.ReplaceAll(rel, "/", "-")
}
