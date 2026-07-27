package init

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/opencode/savepoint/internal/testutil"
)

// splitAuditTemplates is the post-split package payload: both audit skills plus
// the shared, non-triggerable method reference.
func splitAuditTemplates() fstest.MapFS {
	return fstest.MapFS{
		"agent-skills/savepoint-audit-task/SKILL.md": &fstest.MapFile{Data: []byte("# Audit Task")},
		"agent-skills/savepoint-audit-epic/SKILL.md": &fstest.MapFile{Data: []byte("# Audit Epic")},
		"agent-skills/references/audit-method.md":    &fstest.MapFile{Data: []byte("# Shared Method")},
		"agent-skills/savepoint-build-task/SKILL.md": &fstest.MapFile{Data: []byte("# Build Task")},
	}
}

func newLegacyProject(t *testing.T, legacyContent string) string {
	t.Helper()

	target := t.TempDir()
	if err := os.MkdirAll(filepath.Join(target, ".savepoint"), 0755); err != nil {
		t.Fatal(err)
	}
	if legacyContent != "" {
		dir := filepath.Join(target, "agent-skills", "savepoint-audit")
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		testutil.WriteFile(t, filepath.Join(dir, "SKILL.md"), legacyContent)
	}
	return target
}

func archivePath(target string, index int) string {
	name := "savepoint-audit-SKILL.md"
	if index > 0 {
		name = "savepoint-audit-SKILL." + string(rune('0'+index)) + ".md"
	}
	return filepath.Join(target, ".savepoint", "migrations", name)
}

func TestUpgrade_installsSplitSkillsAndSharedReference(t *testing.T) {
	target := newLegacyProject(t, "")
	templates := splitAuditTemplates()

	report, err := UpgradeProjectAssets(templates, target, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	for path, file := range templates {
		action, found := actionFor(report, path)
		if !found {
			t.Fatalf("%s not in report", path)
		}
		if action != ActionUpdated {
			t.Errorf("%s action = %v, want updated", path, action)
		}
		data, err := os.ReadFile(filepath.Join(target, filepath.FromSlash(path)))
		if err != nil {
			t.Errorf("%s not installed: %v", path, err)
			continue
		}
		if string(data) != string(file.Data) {
			t.Errorf("%s = %q, want %q", path, string(data), string(file.Data))
		}
	}
}

func TestUpgrade_sharedReferenceRefreshesLikeASkill(t *testing.T) {
	target := newLegacyProject(t, "")
	refDir := filepath.Join(target, "agent-skills", "references")
	if err := os.MkdirAll(refDir, 0755); err != nil {
		t.Fatal(err)
	}
	testutil.WriteFile(t, filepath.Join(refDir, "audit-method.md"), "# Old Method")

	report, err := UpgradeProjectAssets(splitAuditTemplates(), target, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	action, found := actionFor(report, "agent-skills/references/audit-method.md")
	if !found {
		t.Fatal("shared reference not in report")
	}
	if action != ActionUpdated {
		t.Errorf("shared reference action = %v, want updated", action)
	}

	data, err := os.ReadFile(filepath.Join(refDir, "audit-method.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "# Shared Method" {
		t.Errorf("shared reference = %q, want refreshed content", string(data))
	}
}

func TestUpgrade_migratesStockLegacyAuditSkill(t *testing.T) {
	legacy := "# Old Generic Audit Skill"
	target := newLegacyProject(t, legacy)

	report, err := UpgradeProjectAssets(splitAuditTemplates(), target, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	action, found := actionFor(report, "agent-skills/savepoint-audit/SKILL.md")
	if !found {
		t.Fatal("legacy skill not in report")
	}
	if action != ActionMigrated {
		t.Errorf("legacy action = %v, want migrated", action)
	}

	archived, err := os.ReadFile(archivePath(target, 0))
	if err != nil {
		t.Fatalf("legacy content not archived: %v", err)
	}
	if string(archived) != legacy {
		t.Errorf("archived = %q, want %q", string(archived), legacy)
	}

	if _, err := os.Stat(filepath.Join(target, "agent-skills", "savepoint-audit", "SKILL.md")); !os.IsNotExist(err) {
		t.Errorf("legacy skill still triggerable, stat err = %v", err)
	}
	if _, err := os.Stat(filepath.Join(target, "agent-skills", "savepoint-audit")); !os.IsNotExist(err) {
		t.Errorf("empty legacy folder not removed, stat err = %v", err)
	}

	readme, err := os.ReadFile(filepath.Join(target, ".savepoint", "migrations", "README.md"))
	if err != nil {
		t.Fatalf("migrations README not written: %v", err)
	}
	if !strings.Contains(string(readme), "nothing in this directory is loaded by an agent") {
		t.Errorf("migrations README does not document non-triggerable status: %q", string(readme))
	}
}

func TestUpgrade_preservesUserModifiedLegacySkill(t *testing.T) {
	legacy := "# My Locally Edited Audit Skill\n\nCustom project rules."
	target := newLegacyProject(t, legacy)

	if _, err := UpgradeProjectAssets(splitAuditTemplates(), target, false, false); err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	archived, err := os.ReadFile(archivePath(target, 0))
	if err != nil {
		t.Fatalf("edited legacy content not archived: %v", err)
	}
	if string(archived) != legacy {
		t.Errorf("archived = %q, want verbatim %q", string(archived), legacy)
	}
}

func TestUpgrade_withoutLegacySkillCreatesNoArchive(t *testing.T) {
	target := newLegacyProject(t, "")

	report, err := UpgradeProjectAssets(splitAuditTemplates(), target, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	if _, found := actionFor(report, "agent-skills/savepoint-audit/SKILL.md"); found {
		t.Error("legacy skill reported for a project that never had it")
	}
	if _, err := os.Stat(filepath.Join(target, ".savepoint", "migrations")); !os.IsNotExist(err) {
		t.Errorf("unnecessary migrations archive created, stat err = %v", err)
	}
}

func TestUpgrade_migrationIsIdempotent(t *testing.T) {
	legacy := "# Old Generic Audit Skill"
	target := newLegacyProject(t, legacy)

	if _, err := UpgradeProjectAssets(splitAuditTemplates(), target, false, false); err != nil {
		t.Fatalf("first UpgradeProjectAssets() error = %v", err)
	}
	report, err := UpgradeProjectAssets(splitAuditTemplates(), target, false, false)
	if err != nil {
		t.Fatalf("second UpgradeProjectAssets() error = %v", err)
	}

	if _, found := actionFor(report, "agent-skills/savepoint-audit/SKILL.md"); found {
		t.Error("legacy skill reported again after migration")
	}

	entries, err := os.ReadDir(filepath.Join(target, ".savepoint", "migrations"))
	if err != nil {
		t.Fatal(err)
	}
	// README.md plus exactly one archived copy.
	if len(entries) != 2 {
		t.Errorf("migrations dir has %d entries, want 2", len(entries))
	}
}

func TestUpgrade_reArchivesDifferingLegacyContentWithoutOverwriting(t *testing.T) {
	first := "# First Legacy Copy"
	target := newLegacyProject(t, first)

	if _, err := UpgradeProjectAssets(splitAuditTemplates(), target, false, false); err != nil {
		t.Fatalf("first UpgradeProjectAssets() error = %v", err)
	}

	// A second project state reintroduces the legacy skill with different content.
	second := "# Second Legacy Copy"
	dir := filepath.Join(target, "agent-skills", "savepoint-audit")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	testutil.WriteFile(t, filepath.Join(dir, "SKILL.md"), second)

	if _, err := UpgradeProjectAssets(splitAuditTemplates(), target, false, false); err != nil {
		t.Fatalf("second UpgradeProjectAssets() error = %v", err)
	}

	original, err := os.ReadFile(archivePath(target, 0))
	if err != nil {
		t.Fatal(err)
	}
	if string(original) != first {
		t.Errorf("existing archive overwritten: got %q, want %q", string(original), first)
	}

	conflict, err := os.ReadFile(archivePath(target, 1))
	if err != nil {
		t.Fatalf("differing copy not archived to a numbered sibling: %v", err)
	}
	if string(conflict) != second {
		t.Errorf("conflict archive = %q, want %q", string(conflict), second)
	}
}

func TestUpgrade_dryRunReportsMigrationWithoutWriting(t *testing.T) {
	legacy := "# Old Generic Audit Skill"
	target := newLegacyProject(t, legacy)

	report, err := UpgradeProjectAssets(splitAuditTemplates(), target, true, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() dry-run error = %v", err)
	}

	action, found := actionFor(report, "agent-skills/savepoint-audit/SKILL.md")
	if !found {
		t.Fatal("legacy skill not in dry-run report")
	}
	if action != ActionMigrated {
		t.Errorf("dry-run legacy action = %v, want migrated", action)
	}

	data, err := os.ReadFile(filepath.Join(target, "agent-skills", "savepoint-audit", "SKILL.md"))
	if err != nil {
		t.Fatalf("dry-run removed the legacy skill: %v", err)
	}
	if string(data) != legacy {
		t.Errorf("dry-run altered legacy content: %q", string(data))
	}
	if _, err := os.Stat(filepath.Join(target, ".savepoint", "migrations")); !os.IsNotExist(err) {
		t.Errorf("dry-run wrote the archive, stat err = %v", err)
	}
}

func TestUpgradeReport_formatCountsMigrated(t *testing.T) {
	r := &UpgradeReport{
		Actions: []UpgradeEntry{
			{Path: "agent-skills/savepoint-audit/SKILL.md", Action: ActionMigrated},
			{Path: "agent-skills/savepoint-audit-epic/SKILL.md", Action: ActionUpdated},
		},
	}

	output := r.Format()
	if !strings.Contains(output, "Migrated: 1") {
		t.Errorf("missing migrated count: %q", output)
	}
	if !strings.Contains(output, "migrated  agent-skills/savepoint-audit/SKILL.md") {
		t.Errorf("missing migrated entry line: %q", output)
	}
}
