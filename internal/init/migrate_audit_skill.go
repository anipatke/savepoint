package init

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

// The generic savepoint-audit skill was replaced by savepoint-audit-task and
// savepoint-audit-epic. Upgrade must retire it without discarding local edits,
// so its content moves to a Savepoint-owned archive that no agent loads before
// the triggerable copy is deleted.
const (
	legacyAuditSkillDir     = "agent-skills/savepoint-audit"
	legacyAuditSkillFile    = "agent-skills/savepoint-audit/SKILL.md"
	migrationsDir           = ".savepoint/migrations"
	legacyAuditArchiveStem  = "savepoint-audit-SKILL"
	migrationsReadmeName    = "README.md"
	maxArchiveConflictTries = 100
)

const migrationsReadme = `# Savepoint Migrations

Archived copies of package-owned assets that Savepoint retired during
` + "`savepoint upgrade-assets`" + `. Files here are records, not instructions:
nothing in this directory is loaded by an agent or triggerable as a skill.

- ` + "`savepoint-audit-SKILL.md`" + ` — the generic audit skill retired when audit
  split into ` + "`savepoint-audit-task`" + ` (read-only review of one in-progress
  task) and ` + "`savepoint-audit-epic`" + ` (audit-pending closeout). Kept so local
  edits are recoverable. A numbered suffix means a differing copy was archived by
  a later upgrade.

- The nine V1 skills and the shared audit method reference, retired when a
  project upgrades on the V2 workflow: ` + "`savepoint-draft-prd`" + `,
  ` + "`savepoint-create-plan`" + `, ` + "`savepoint-system-design`" + `,
  ` + "`savepoint-create-task`" + `, ` + "`savepoint-build-task`" + `,
  ` + "`savepoint-audit-task`" + `, ` + "`savepoint-audit-epic`" + `,
  ` + "`savepoint-audit-register`" + `, ` + "`savepoint-create-defect`" + `, and
  ` + "`references/audit-method.md`" + `. Every one is archived, edited or not,
  because deciding which local edits were worth keeping is not this command's
  call to make. A numbered suffix means a differing copy was archived by a later
  upgrade.

Delete anything here once you have salvaged what you need.
`

// migrateLegacyAuditSkill retires agent-skills/savepoint-audit by archiving its
// content under .savepoint/migrations/ and deleting the triggerable copy. It
// returns nil when the project never had the legacy skill, so projects that are
// already on the split skills gain no archive.
func migrateLegacyAuditSkill(absTarget string, dryRun bool, write assetWriter) (*UpgradeEntry, error) {
	legacyPath := filepath.Join(absTarget, filepath.FromSlash(legacyAuditSkillFile))

	content, err := os.ReadFile(legacyPath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read legacy audit skill: %w", err)
	}

	archivePath, err := resolveArchivePath(absTarget, legacyAuditArchiveStem, content)
	if err != nil {
		return nil, err
	}

	entry := &UpgradeEntry{Path: legacyAuditSkillFile, Action: ActionMigrated}
	if dryRun {
		return entry, nil
	}

	// Preserve first, delete second: the triggerable copy is removed only after
	// its content is safely on disk somewhere else.
	if archivePath != "" {
		if err := writeArchive(absTarget, archivePath, content, write); err != nil {
			return nil, err
		}
	}

	if err := os.Remove(legacyPath); err != nil {
		return nil, fmt.Errorf("remove legacy audit skill: %w", err)
	}
	removeDirIfEmpty(filepath.Join(absTarget, filepath.FromSlash(legacyAuditSkillDir)))

	return entry, nil
}

// resolveArchivePath applies the archive conflict policy: reuse nothing, never
// overwrite. It returns an empty path when an identical archive already exists
// under stem, which makes repeated upgrades idempotent, and a numbered sibling
// path when a differing archive is already present.
func resolveArchivePath(absTarget, stem string, content []byte) (string, error) {
	dir := filepath.Join(absTarget, filepath.FromSlash(migrationsDir))

	for i := 0; i < maxArchiveConflictTries; i++ {
		name := stem + ".md"
		if i > 0 {
			name = fmt.Sprintf("%s.%d.md", stem, i)
		}
		candidate := filepath.Join(dir, name)

		existing, err := os.ReadFile(candidate)
		if os.IsNotExist(err) {
			return candidate, nil
		}
		if err != nil {
			return "", fmt.Errorf("read migration archive %s: %w", name, err)
		}
		if bytes.Equal(existing, content) {
			return "", nil
		}
	}

	return "", fmt.Errorf("cannot archive %s: %s already holds %d differing copies", stem, migrationsDir, maxArchiveConflictTries)
}

func writeArchive(absTarget, archivePath string, content []byte, write assetWriter) error {
	dir := filepath.Join(absTarget, filepath.FromSlash(migrationsDir))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create %s: %w", migrationsDir, err)
	}

	readmePath := filepath.Join(dir, migrationsReadmeName)
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		if err := write(readmePath, []byte(migrationsReadme)); err != nil {
			return fmt.Errorf("write %s/%s: %w", migrationsDir, migrationsReadmeName, err)
		}
	} else if err != nil {
		return fmt.Errorf("stat %s/%s: %w", migrationsDir, migrationsReadmeName, err)
	}

	if err := write(archivePath, content); err != nil {
		return fmt.Errorf("write migration archive: %w", err)
	}
	return nil
}

// removeDirIfEmpty drops the now-skill-less legacy folder. A non-empty folder is
// user content, so a failure here is deliberately not an error.
func removeDirIfEmpty(dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) > 0 {
		return
	}
	_ = os.Remove(dir)
}
