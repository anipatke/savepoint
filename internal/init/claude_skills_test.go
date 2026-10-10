package init

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var workflowSkills = []string{"savepoint-idea", "savepoint-design", "savepoint-task", "savepoint-check"}

func pointerPath(skill string) string {
	return ".claude/skills/" + skill + "/SKILL.md"
}

func readTemplateFS(t *testing.T, templates fs.FS, path string) string {
	t.Helper()
	data, err := fs.ReadFile(templates, path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestClaudeSkillPointers_matchCanonicalNameAndDescription(t *testing.T) {
	templates := shippedTemplates(t)
	for _, skill := range workflowSkills {
		pointer := readTemplateFS(t, templates, pointerPath(skill))
		canonical := readTemplateFS(t, templates, "agent-skills/"+skill+"/SKILL.md")
		for _, key := range []string{"name", "description"} {
			want := frontmatterField(canonical, key)
			if want == "" {
				t.Fatalf("%s: canonical %s is empty", skill, key)
			}
			if got := frontmatterField(pointer, key); got != want {
				t.Errorf("%s: pointer %s = %q, want canonical %q", skill, key, got, want)
			}
		}
	}
}

func TestClaudeSkillPointers_bodyOnlyPointsAtCanonicalSkill(t *testing.T) {
	templates := shippedTemplates(t)
	for _, skill := range workflowSkills {
		pointer := strings.ReplaceAll(readTemplateFS(t, templates, pointerPath(skill)), "\r\n", "\n")
		body := strings.TrimSpace(pointer[strings.Index(pointer[4:], "\n---")+8:])
		if !strings.Contains(body, "agent-skills/"+skill+"/SKILL.md") {
			t.Errorf("%s: body does not point at the canonical skill: %q", skill, body)
		}
		if strings.Contains(body, "\n\n") || strings.Count(body, ". ") > 1 {
			t.Errorf("%s: body should be one or two sentences: %q", skill, body)
		}
	}
}

func TestClaudeSkillPointers_onlyWorkflowSkillsHavePointers(t *testing.T) {
	entries, err := os.ReadDir(filepath.Join("..", "..", "templates", "project-v2", ".claude", "skills"))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != len(workflowSkills) {
		t.Errorf("pointer directories = %d, want %d", len(entries), len(workflowSkills))
	}
}

func TestClaudeSkillPointers_repoCopiesAreByteIdenticalToTemplates(t *testing.T) {
	for _, skill := range workflowSkills {
		want := mustReadFile(t, filepath.Join("..", "..", "templates", "project-v2", pointerPath(skill)))
		got := mustReadFile(t, filepath.Join("..", "..", pointerPath(skill)))
		if string(got) != string(want) {
			t.Errorf("%s differs from its template", pointerPath(skill))
		}
	}
}

func TestScaffold_installsTrackedClaudeSkillPointers(t *testing.T) {
	dir := t.TempDir()
	if err := Scaffold(shippedTemplates(t), dir, "p", false); err != nil {
		t.Fatal(err)
	}
	manifest, err := LoadManifest(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, skill := range workflowSkills {
		path := pointerPath(skill)
		content := readProjectFile(t, dir, path)
		if hash, ok := manifest.Hash(path); !ok || hash != hashContent([]byte(content)) {
			t.Errorf("manifest entry for %s missing or wrong", path)
		}
	}
}

func upgradedPointerEntry(t *testing.T, dir string, dryRun bool, path string) UpgradeEntry {
	t.Helper()
	report, err := upgradeProjectAssets(shippedTemplates(t), dir, dryRun, false, AtomicWrite, false)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range report.Actions {
		if e.Path == path {
			return e
		}
	}
	t.Fatalf("no report entry for %s", path)
	return UpgradeEntry{}
}

func TestUpgrade_claudeSkillPointers(t *testing.T) {
	path := pointerPath("savepoint-task")
	setup := func(t *testing.T) string {
		dir := t.TempDir()
		if err := Scaffold(shippedTemplates(t), dir, "p", false); err != nil {
			t.Fatal(err)
		}
		return dir
	}

	t.Run("missing pointer is installed", func(t *testing.T) {
		dir := setup(t)
		if err := os.Remove(filepath.Join(dir, path)); err != nil {
			t.Fatal(err)
		}
		if e := upgradedPointerEntry(t, dir, false, path); e.Action != ActionUpdated {
			t.Errorf("action = %s, want updated", e.Action)
		}
		if _, err := os.Stat(filepath.Join(dir, path)); err != nil {
			t.Errorf("pointer not restored: %v", err)
		}
	})

	t.Run("unedited outdated pointer is refreshed", func(t *testing.T) {
		dir := setup(t)
		old := []byte("old pointer\n")
		if err := os.WriteFile(filepath.Join(dir, path), old, 0644); err != nil {
			t.Fatal(err)
		}
		manifest, _ := LoadManifest(dir)
		manifest.Record(path, old)
		if err := manifest.Save(dir); err != nil {
			t.Fatal(err)
		}
		if e := upgradedPointerEntry(t, dir, false, path); e.Action != ActionUpdated || e.Note != "" {
			t.Errorf("entry = %+v, want updated without backup", e)
		}
		want := mustReadFile(t, filepath.Join("..", "..", "templates", "project-v2", path))
		if got := mustReadFile(t, filepath.Join(dir, path)); string(got) != string(want) {
			t.Error("pointer was not refreshed")
		}
	})

	t.Run("edited pointer is kept and reported", func(t *testing.T) {
		dir := setup(t)
		mine := []byte("my own pointer\n")
		if err := os.WriteFile(filepath.Join(dir, path), mine, 0644); err != nil {
			t.Fatal(err)
		}
		if e := upgradedPointerEntry(t, dir, false, path); e.Action != ActionConflict {
			t.Errorf("action = %s, want conflict", e.Action)
		}
		if got := mustReadFile(t, filepath.Join(dir, path)); string(got) != string(mine) {
			t.Error("edited pointer was overwritten")
		}
		if _, err := os.Stat(filepath.Join(dir, path+incomingSuffix)); err != nil {
			t.Errorf("incoming pointer not written beside it: %v", err)
		}
	})

	t.Run("dry run writes nothing", func(t *testing.T) {
		dir := setup(t)
		if err := os.Remove(filepath.Join(dir, path)); err != nil {
			t.Fatal(err)
		}
		if e := upgradedPointerEntry(t, dir, true, path); e.Action != ActionUpdated {
			t.Errorf("action = %s, want updated", e.Action)
		}
		if _, err := os.Stat(filepath.Join(dir, path)); !os.IsNotExist(err) {
			t.Error("dry run wrote the pointer")
		}
	})
}
