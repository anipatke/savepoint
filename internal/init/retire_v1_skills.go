package init

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// migrationsDir and migrationsReadmeName locate the non-triggerable archive
// retirement writes content to; maxArchiveConflictTries bounds the numbered-
// sibling search resolveArchivePath uses when an archive path already holds
// differing content.
const (
	migrationsDir           = ".savepoint/migrations"
	migrationsReadmeName    = "README.md"
	maxArchiveConflictTries = 100
)

const migrationsReadme = `# Savepoint Migrations

Archived copies of package-owned assets that Savepoint retired during
` + "`savepoint upgrade-assets`" + `. Files here are records, not instructions:
nothing in this directory is loaded by an agent or triggerable as a skill.

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

// legacyMigrationsReadme is the exact README the V1 upgrade path shipped
// before V1 skill retirement existed, when only the legacy generic audit
// skill was migrated. A migrations directory can still carry this exact text
// from that era; retirement upgrades it to migrationsReadme like any other
// stock (unedited) copy, keeping it byte-for-byte stable so this comparison
// never misses a real one.
const legacyMigrationsReadme = `# Savepoint Migrations

Archived copies of package-owned assets that Savepoint retired during
` + "`savepoint upgrade-assets`" + `. Files here are records, not instructions:
nothing in this directory is loaded by an agent or triggerable as a skill.

- ` + "`savepoint-audit-SKILL.md`" + ` — the generic audit skill retired when audit
  split into ` + "`savepoint-audit-task`" + ` (read-only review of one in-progress
  task) and ` + "`savepoint-audit-epic`" + ` (audit-pending closeout). Kept so local
  edits are recoverable. A numbered suffix means a differing copy was archived by
  a later upgrade.

Delete anything here once you have salvaged what you need.
`

// resolveArchivePath applies the archive conflict policy: reuse an identical
// archive, never overwrite. It returns the exact selected path and whether
// that path needs a write. A differing archive gets a numbered sibling path.
func resolveArchivePath(absTarget, stem string, content []byte) (string, bool, error) {
	dir := filepath.Join(absTarget, filepath.FromSlash(migrationsDir))

	for i := 0; i < maxArchiveConflictTries; i++ {
		name := stem + ".md"
		if i > 0 {
			name = fmt.Sprintf("%s.%d.md", stem, i)
		}
		candidate := filepath.Join(dir, name)

		existing, err := os.ReadFile(candidate)
		if os.IsNotExist(err) {
			return candidate, true, nil
		}
		if err != nil {
			return "", false, fmt.Errorf("read migration archive %s: %w", name, err)
		}
		if bytes.Equal(existing, content) {
			return candidate, false, nil
		}
	}

	return "", false, fmt.Errorf("cannot archive %s: %s already holds %d differing copies", stem, migrationsDir, maxArchiveConflictTries)
}

// writeArchive writes archivePath's content and ensures the migrations README
// is current: it creates the README from readme when absent, and upgrades it
// in place when it still holds the exact pre-retirement legacyMigrationsReadme
// text. Any other existing README is left alone — it is either already
// current or the user's own edit. archiveWrite is false when an identical
// archive already exists, so only the README write is attempted.
func writeArchive(absTarget, archivePath string, content []byte, readme string, archiveWrite bool, write assetWriter) error {
	dir := filepath.Join(absTarget, filepath.FromSlash(migrationsDir))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create %s: %w", migrationsDir, err)
	}

	readmePath := filepath.Join(dir, migrationsReadmeName)
	existingReadme, err := os.ReadFile(readmePath)
	switch {
	case os.IsNotExist(err):
		if err := write(readmePath, []byte(readme)); err != nil {
			return fmt.Errorf("write %s/%s: %w", migrationsDir, migrationsReadmeName, err)
		}
	case err != nil:
		return fmt.Errorf("stat %s/%s: %w", migrationsDir, migrationsReadmeName, err)
	case bytes.Equal(existingReadme, []byte(legacyMigrationsReadme)) && !bytes.Equal(existingReadme, []byte(readme)):
		if err := write(readmePath, []byte(readme)); err != nil {
			return fmt.Errorf("write %s/%s: %w", migrationsDir, migrationsReadmeName, err)
		}
	}

	if archiveWrite {
		if err := write(archivePath, content); err != nil {
			return fmt.Errorf("write migration archive: %w", err)
		}
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
// reference from a V2 project: archive the content under
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
		if !dryRun {
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

	archivePath, archiveWrite, err := resolveArchivePath(absTarget, archiveStem(path), content)
	if err != nil {
		return &UpgradeEntry{Path: path, Action: ActionFailed}, err
	}

	entry := &UpgradeEntry{
		Path:   path,
		Action: ActionRetired,
		Note:   archiveRecoveryNote(absTarget, archivePath),
	}
	if dryRun {
		return entry, nil
	}

	// Preserve first, delete second: the triggerable copy is removed only after
	// its content is safely on disk somewhere else.
	if err := writeArchive(absTarget, archivePath, content, migrationsReadme, archiveWrite, write); err != nil {
		return &UpgradeEntry{Path: path, Action: ActionFailed, Note: entry.Note}, fmt.Errorf("archive retired asset %s: %w", path, err)
	}

	if err := os.Remove(targetPath); err != nil {
		return &UpgradeEntry{Path: path, Action: ActionFailed, Note: entry.Note}, fmt.Errorf("remove retired asset %s: %w", path, err)
	}
	removeDirIfEmpty(filepath.Dir(targetPath))

	return entry, nil
}

func archiveRecoveryNote(absTarget, archivePath string) string {
	rel, err := filepath.Rel(absTarget, archivePath)
	if err != nil {
		return "archived copy saved to " + filepath.ToSlash(archivePath)
	}
	return "archived copy saved to " + filepath.ToSlash(rel)
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
