package init

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/opencode/savepoint/internal/testutil"
)

// Tests in this file cover retirement of the nine V1 skills and the shared
// audit method reference from a V2 project: T006's job. Version dispatch
// itself — which tree an upgrade walks — is covered in upgrade_schema_test.go;
// these tests only need retirement to fire, so their template trees are
// minimal and distinct from the V1 skills being retired.

// retirementTemplates returns a tiny, easily told-apart V1/V2 template pair.
// Neither tree includes any of the retired paths: retirement acts on what the
// project already has on disk, not on what either tree ships.
func retirementTemplates() (fstest.MapFS, fstest.MapFS) {
	v1 := fstest.MapFS{
		"agent-skills/savepoint-check/SKILL.md": &fstest.MapFile{Data: []byte("# not a real v1 skill, just tree filler")},
	}
	v2 := fstest.MapFS{
		"agent-skills/savepoint-idea/SKILL.md": &fstest.MapFile{Data: []byte("# Idea")},
	}
	return v1, v2
}

// v2ProjectWithLegacySkills builds a V2 project (schema_version: 2) that still
// carries every retired V1 asset, the state migrate leaves behind and T005
// deliberately does not touch. contents overrides the stock body for any path
// present as a key; every other retired path gets a distinct stock body.
func v2ProjectWithLegacySkills(t *testing.T, contents map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	testutil.MkdirAll(t, filepath.Join(dir, ".savepoint"))
	testutil.WriteFile(t, filepath.Join(dir, ".savepoint", "config.yml"), "schema_version: 2\n")
	for _, path := range retiredV1AssetPaths() {
		body, ok := contents[path]
		if !ok {
			body = "# stock " + path
		}
		testutil.WriteFile(t, filepath.Join(dir, filepath.FromSlash(path)), body)
	}
	return dir
}

// v1Project builds a project with no schema_version at all: the state a
// project that has never run migrate looks like.
func v1Project(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	testutil.MkdirAll(t, filepath.Join(dir, ".savepoint"))
	return dir
}

func retiredArchivePath(dir, path string, index int) string {
	name := archiveStem(path) + ".md"
	if index > 0 {
		name = fmt.Sprintf("%s.%d.md", archiveStem(path), index)
	}
	return filepath.Join(dir, ".savepoint", "migrations", name)
}

func TestRetireV1Skills_fullRetirement(t *testing.T) {
	dir := v2ProjectWithLegacySkills(t, nil)
	v1, v2 := retirementTemplates()

	// Seed the manifest as a project that upgraded on V1 before migrating would
	// have it: every retired skill already carries a recorded hash. The shared
	// audit-method reference is never manifest-tracked, matching every other
	// agent-skills/references/ asset. Retirement must drop each recorded skill,
	// not merely find nothing to drop.
	seed := NewManifest()
	for _, skill := range retiredV1SkillDirs {
		path := "agent-skills/" + skill + "/SKILL.md"
		seed.Record(path, []byte("# stock "+path))
	}
	if err := seed.Save(dir); err != nil {
		t.Fatal(err)
	}
	for _, skill := range retiredV1SkillDirs {
		path := "agent-skills/" + skill + "/SKILL.md"
		if _, tracked := seed.Hash(path); !tracked {
			t.Fatalf("test setup: %s not recorded in the seeded manifest", path)
		}
	}

	report, err := UpgradeProjectAssets(v1, v2, dir, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	manifest, err := LoadManifest(dir)
	if err != nil {
		t.Fatal(err)
	}

	for _, path := range retiredV1AssetPaths() {
		entry, found := entryFor(report, path)
		if !found || entry.Action != ActionRetired {
			t.Errorf("%s entry = %+v (found %v), want retired", path, entry, found)
		}
		wantNote := "archived copy saved to " + filepath.ToSlash(filepath.Join(migrationsDir, archiveStem(path)+".md"))
		if found && entry.Note != wantNote {
			t.Errorf("%s note = %q, want %q", path, entry.Note, wantNote)
		}

		if _, err := os.Stat(filepath.Join(dir, filepath.FromSlash(path))); !os.IsNotExist(err) {
			t.Errorf("%s still on disk, stat err = %v", path, err)
		}

		archived, err := os.ReadFile(retiredArchivePath(dir, path, 0))
		if err != nil {
			t.Fatalf("%s not archived: %v", path, err)
		}
		if string(archived) != "# stock "+path {
			t.Errorf("archived %s = %q, want stock content", path, string(archived))
		}

		if _, tracked := manifest.Hash(path); tracked {
			t.Errorf("manifest still records retired %s", path)
		}
	}

	for _, skill := range retiredV1SkillDirs {
		if _, err := os.Stat(filepath.Join(dir, "agent-skills", skill)); !os.IsNotExist(err) {
			t.Errorf("emptied skill dir %s not removed, stat err = %v", skill, err)
		}
	}
}

func TestRetireV1Skills_missingFileForgetsStaleManifestEntry(t *testing.T) {
	path := "agent-skills/savepoint-draft-prd/SKILL.md"
	dir := v2ProjectWithLegacySkills(t, nil)
	if err := os.Remove(filepath.Join(dir, filepath.FromSlash(path))); err != nil {
		t.Fatal(err)
	}

	manifest := NewManifest()
	manifest.Record(path, []byte("# stale provenance"))
	if err := manifest.Save(dir); err != nil {
		t.Fatal(err)
	}

	v1, v2 := retirementTemplates()
	report, err := UpgradeProjectAssets(v1, v2, dir, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}
	if _, found := actionFor(report, path); found {
		t.Errorf("missing retired path reported an action: %+v", report.Actions)
	}

	loaded, err := LoadManifest(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, tracked := loaded.Hash(path); tracked {
		t.Errorf("manifest retained stale entry for missing %s", path)
	}
}

func TestRetireV1Skills_preservesUserEditsAndArchivesUnmodifiedAlike(t *testing.T) {
	edited := "agent-skills/savepoint-build-task/SKILL.md"
	editedBody := "# My locally edited build-task skill\n\nCustom project rules."
	dir := v2ProjectWithLegacySkills(t, map[string]string{edited: editedBody})
	v1, v2 := retirementTemplates()

	if _, err := UpgradeProjectAssets(v1, v2, dir, false, false); err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	archived, err := os.ReadFile(retiredArchivePath(dir, edited, 0))
	if err != nil {
		t.Fatalf("edited skill not archived: %v", err)
	}
	if string(archived) != editedBody {
		t.Errorf("archived edited skill = %q, want verbatim %q", string(archived), editedBody)
	}

	unmodified := "agent-skills/savepoint-create-task/SKILL.md"
	archivedStock, err := os.ReadFile(retiredArchivePath(dir, unmodified, 0))
	if err != nil {
		t.Fatalf("unmodified skill not archived: %v", err)
	}
	if string(archivedStock) != "# stock "+unmodified {
		t.Errorf("archived unmodified skill = %q, want stock content", string(archivedStock))
	}
}

func TestRetireV1Skills_reArchivesDifferingContentWithoutOverwriting(t *testing.T) {
	path := "agent-skills/savepoint-draft-prd/SKILL.md"
	first := "# First retirement pass"
	dir := v2ProjectWithLegacySkills(t, map[string]string{path: first})
	v1, v2 := retirementTemplates()

	if _, err := UpgradeProjectAssets(v1, v2, dir, false, false); err != nil {
		t.Fatalf("first UpgradeProjectAssets() error = %v", err)
	}

	// A later state reintroduces the same retired path with different content.
	second := "# Second retirement pass"
	testutil.WriteFile(t, filepath.Join(dir, filepath.FromSlash(path)), second)

	report, err := UpgradeProjectAssets(v1, v2, dir, false, false)
	if err != nil {
		t.Fatalf("second UpgradeProjectAssets() error = %v", err)
	}
	entry, found := entryFor(report, path)
	if !found || entry.Action != ActionRetired {
		t.Fatalf("second retirement entry = %+v (found %v), want retired", entry, found)
	}
	wantNote := "archived copy saved to " + filepath.ToSlash(filepath.Join(migrationsDir, archiveStem(path)+".1.md"))
	if entry.Note != wantNote {
		t.Errorf("second retirement note = %q, want %q", entry.Note, wantNote)
	}
	if !strings.Contains(report.Format(), wantNote) {
		t.Errorf("formatted second retirement report omits numbered archive: %q", report.Format())
	}

	original, err := os.ReadFile(retiredArchivePath(dir, path, 0))
	if err != nil {
		t.Fatal(err)
	}
	if string(original) != first {
		t.Errorf("existing archive overwritten: got %q, want %q", string(original), first)
	}

	conflict, err := os.ReadFile(retiredArchivePath(dir, path, 1))
	if err != nil {
		t.Fatalf("differing copy not archived to a numbered sibling: %v", err)
	}
	if string(conflict) != second {
		t.Errorf("conflict archive = %q, want %q", string(conflict), second)
	}
}

func TestRetireV1Skills_v1ProjectRetainsEverything(t *testing.T) {
	dir := v1Project(t)
	for _, path := range retiredV1AssetPaths() {
		testutil.WriteFile(t, filepath.Join(dir, filepath.FromSlash(path)), "# stock "+path)
	}
	v1, v2 := retirementTemplates()

	report, err := UpgradeProjectAssets(v1, v2, dir, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	for _, path := range retiredV1AssetPaths() {
		if _, found := actionFor(report, path); found {
			t.Errorf("%s reported on a V1 project", path)
		}
		data, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(path)))
		if err != nil {
			t.Fatalf("%s removed from a V1 project: %v", path, err)
		}
		if string(data) != "# stock "+path {
			t.Errorf("%s changed on a V1 project: %q", path, string(data))
		}
	}

	if _, err := os.Stat(filepath.Join(dir, ".savepoint", "migrations")); !os.IsNotExist(err) {
		t.Errorf("V1 project upgrade created a migrations archive, stat err = %v", err)
	}
}

func TestRetireV1Skills_alreadyRetiredProjectIsANoOp(t *testing.T) {
	dir := v2ProjectWithLegacySkills(t, nil)
	v1, v2 := retirementTemplates()

	if _, err := UpgradeProjectAssets(v1, v2, dir, false, false); err != nil {
		t.Fatalf("first UpgradeProjectAssets() error = %v", err)
	}

	before := dirSnapshot(t, dir)
	beforeTimes := mtimeSnapshot(t, dir)

	report, err := UpgradeProjectAssets(v1, v2, dir, false, false)
	if err != nil {
		t.Fatalf("second UpgradeProjectAssets() error = %v", err)
	}

	for _, path := range retiredV1AssetPaths() {
		if _, found := actionFor(report, path); found {
			t.Errorf("%s reported again on an already-retired project", path)
		}
	}

	assertNoChange(t, dir, before)
	after := mtimeSnapshot(t, dir)
	for path, mt := range beforeTimes {
		if after[path] != mt {
			t.Errorf("%s mtime changed on a no-op retirement rerun", path)
		}
	}
}

func TestRetireV1Skills_dryRunReportsAndWritesNothing(t *testing.T) {
	dir := v2ProjectWithLegacySkills(t, nil)
	v1, v2 := retirementTemplates()
	before := dirSnapshot(t, dir)
	beforeTimes := mtimeSnapshot(t, dir)

	report, err := UpgradeProjectAssets(v1, v2, dir, true, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() dry-run error = %v", err)
	}

	for _, path := range retiredV1AssetPaths() {
		action, found := actionFor(report, path)
		if !found || action != ActionRetired {
			t.Errorf("%s dry-run action = %v (found %v), want retired", path, action, found)
		}
	}

	assertNoChange(t, dir, before)
	after := mtimeSnapshot(t, dir)
	for path, mt := range beforeTimes {
		if after[path] != mt {
			t.Errorf("%s mtime changed by a dry run", path)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, ".savepoint", "migrations")); !os.IsNotExist(err) {
		t.Errorf("dry-run wrote a migrations archive, stat err = %v", err)
	}
}

func TestRetireV1Skills_writeFailureLeavesOriginalInPlace(t *testing.T) {
	path := "agent-skills/savepoint-draft-prd/SKILL.md"
	dir := v2ProjectWithLegacySkills(t, nil)
	_, v2 := retirementTemplates()

	w := &countingWriter{failAt: 1}
	report, err := upgradeProjectAssets(v2, dir, false, false, w.write, true)
	if !strings.Contains(err.Error(), "injected write failure") {
		t.Fatalf("error = %v, want the injected failure to surface", err)
	}

	original, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(path)))
	if err != nil {
		t.Fatalf("%s removed despite a failed archive write: %v", path, err)
	}
	if string(original) != "# stock "+path {
		t.Errorf("%s content changed despite a failed archive write: %q", path, string(original))
	}

	if _, err := os.Stat(filepath.Join(dir, ".savepoint", "migrations", archiveStem(path)+".md")); !os.IsNotExist(err) {
		t.Errorf("archive written despite reporting a write failure, stat err = %v", err)
	}

	action, found := actionFor(report, path)
	if !found || action != ActionFailed {
		t.Errorf("report action for %s = %v (found %v), want failed", path, action, found)
	}
}

func TestRetireV1Skills_nonEmptySkillDirIsNotRemoved(t *testing.T) {
	dir := v2ProjectWithLegacySkills(t, nil)
	extra := filepath.Join(dir, "agent-skills", "savepoint-build-task", "notes.md")
	testutil.WriteFile(t, extra, "# user notes left inside the skill folder")
	v1, v2 := retirementTemplates()

	if _, err := UpgradeProjectAssets(v1, v2, dir, false, false); err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "agent-skills", "savepoint-build-task")); err != nil {
		t.Errorf("non-empty skill dir removed: %v", err)
	}
	if _, err := os.Stat(extra); err != nil {
		t.Errorf("user file inside a retired skill dir lost: %v", err)
	}
}

func TestRetireV1Skills_migrationsReadmeNamesRetiredSkillsAndNonTriggerableStatus(t *testing.T) {
	dir := v2ProjectWithLegacySkills(t, nil)
	v1, v2 := retirementTemplates()

	if _, err := UpgradeProjectAssets(v1, v2, dir, false, false); err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	readme, err := os.ReadFile(filepath.Join(dir, ".savepoint", "migrations", "README.md"))
	if err != nil {
		t.Fatalf("migrations README not written: %v", err)
	}
	content := string(readme)
	if !strings.Contains(content, "nothing in this directory is loaded by an agent") {
		t.Errorf("migrations README does not document non-triggerable status: %q", content)
	}
	for _, skill := range retiredV1SkillDirs {
		if !strings.Contains(content, skill) {
			t.Errorf("migrations README does not name retired skill %s", skill)
		}
	}
}

func TestRetireV1Skills_upgradesStockLegacyMigrationsReadmeButPreservesEdits(t *testing.T) {
	t.Run("stock README is upgraded", func(t *testing.T) {
		dir := v2ProjectWithLegacySkills(t, nil)
		testutil.WriteFile(t, filepath.Join(dir, filepath.FromSlash(legacyAuditSkillFile)), "# Old Generic Audit Skill")
		testutil.WriteFile(t, filepath.Join(dir, migrationsDir, migrationsReadmeName), legacyFixture(t, "migrations-readme-pre-e47.md"))

		v1, v2 := retirementTemplates()
		if _, err := UpgradeProjectAssets(v1, v2, dir, false, false); err != nil {
			t.Fatalf("UpgradeProjectAssets() error = %v", err)
		}

		readme, err := os.ReadFile(filepath.Join(dir, migrationsDir, migrationsReadmeName))
		if err != nil {
			t.Fatal(err)
		}
		if string(readme) != migrationsReadme {
			t.Error("V2 retirement did not replace the exact stock pre-E47 README")
		}
	})

	t.Run("edited README is preserved", func(t *testing.T) {
		dir := v2ProjectWithLegacySkills(t, nil)
		custom := "# Local migration notes\n\nKeep this recovery advice.\n"
		testutil.WriteFile(t, filepath.Join(dir, migrationsDir, migrationsReadmeName), custom)

		v1, v2 := retirementTemplates()
		if _, err := UpgradeProjectAssets(v1, v2, dir, false, false); err != nil {
			t.Fatalf("UpgradeProjectAssets() error = %v", err)
		}

		readme, err := os.ReadFile(filepath.Join(dir, migrationsDir, migrationsReadmeName))
		if err != nil {
			t.Fatal(err)
		}
		if string(readme) != custom {
			t.Errorf("edited migrations README changed: %q", string(readme))
		}
	})
}

func TestUpgradeReport_formatDistinguishesRetiredFromInstallsAndUpdates(t *testing.T) {
	r := &UpgradeReport{
		Actions: []UpgradeEntry{
			{Path: "agent-skills/savepoint-draft-prd/SKILL.md", Action: ActionRetired},
			{Path: "agent-skills/savepoint-idea/SKILL.md", Action: ActionInstalled},
			{Path: "AGENTS.md", Action: ActionUpdated},
		},
	}

	output := r.Format()
	if !strings.Contains(output, "Retired: 1") {
		t.Errorf("missing retired count: %q", output)
	}
	if !strings.Contains(output, "retired   agent-skills/savepoint-draft-prd/SKILL.md") &&
		!strings.Contains(output, "retired  agent-skills/savepoint-draft-prd/SKILL.md") {
		t.Errorf("missing retired entry line: %q", output)
	}
	if strings.Contains(output, "Installed: 2") || strings.Contains(output, "Updated: 2") {
		t.Errorf("retired entry counted as install/update: %q", output)
	}
}
