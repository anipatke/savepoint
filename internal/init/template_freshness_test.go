package init

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestProjectTemplatesUseCurrentWorkflow(t *testing.T) {
	root := filepath.Join("..", "..")
	agents := readTemplate(t, root, "templates", "project", "AGENTS.md")
	router := readTemplate(t, root, "templates", "project", ".savepoint", "router.md")
	auditSkill := readTemplate(t, root, "templates", "project", "agent-skills", "savepoint-audit", "SKILL.md")

	assertNotContains(t, agents, "`phase` (build/test/audit)")
	assertNotContains(t, agents, "npm run build && npm run test")
	assertContains(t, agents, "`stage` (build/test/audit): **required** when `status: in_progress`")
	assertContains(t, agents, "make build && make test")

	assertNotContains(t, router, ".savepoint/audit/{E##-epic}/snapshot.md")
	assertNotContains(t, router, ".savepoint/audit/{release}/{E##-epic}/proposals.md")
	assertNotContains(t, router, ".savepoint/audit/{E##-epic}/proposals.md")
	assertContains(t, router, ".savepoint/releases/{release}/epics/{E##-epic}/E##-Audit.md")
	assertContains(t, router, "`## Proposed Changes` — admin/apply metadata")
	assertContains(t, agents, "During audit apply/close, update the same `E##-Audit.md` visible sections")
	assertContains(t, auditSkill, "Update `E##-Audit.md` visible sections")
	assertContains(t, auditSkill, "Updated audit findings")
}

func TestProjectGuidanceTemplatesMirrorLiveGuidance(t *testing.T) {
	root := filepath.Join("..", "..")

	liveAgents := readTemplate(t, root, "AGENTS.md")
	templateAgents := readTemplate(t, root, "templates", "project", "AGENTS.md")
	for _, canonical := range []string{
		"The phase skill is the canonical workflow source.",
		"Task `stage` (build/test/audit): **required** when `status: in_progress`",
		"Task lifecycle rules are owned by `internal/data`; legacy `phase` is parse compatibility only and must not be used in new task guidance.",
		"Never write `stage: implementation`; use `stage: build` when starting implementation work.",
		"Only the user may set a task to `status: done`",
		"Never run `savepoint` commands.",
		"make build && make test",
	} {
		assertContains(t, liveAgents, canonical)
		assertContains(t, templateAgents, canonical)
	}

	liveSkillRoot := filepath.Join(root, "agent-skills")
	templateSkillRoot := filepath.Join(root, "templates", "project", "agent-skills")
	entries, err := os.ReadDir(liveSkillRoot)
	if err != nil {
		t.Fatalf("read live skill root: %v", err)
	}

	var skillNames []string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "savepoint-") {
			skillNames = append(skillNames, entry.Name())
		}
	}
	sort.Strings(skillNames)

	if len(skillNames) == 0 {
		t.Fatal("no live skills found")
	}

	for _, name := range skillNames {
		assertFileMatches(t, root,
			filepath.Join("agent-skills", name, "SKILL.md"),
			filepath.Join("templates", "project", "agent-skills", name, "SKILL.md"),
		)
		if _, err := os.Stat(filepath.Join(templateSkillRoot, name, "SKILL.md")); err != nil {
			t.Fatalf("missing scaffolded skill %s: %v", name, err)
		}
	}

	templateEntries, err := os.ReadDir(templateSkillRoot)
	if err != nil {
		t.Fatalf("read scaffold skill root: %v", err)
	}
	if len(templateEntries) != len(skillNames) {
		t.Fatalf("scaffolded skill count = %d, want %d", len(templateEntries), len(skillNames))
	}
}

func TestProjectTemplatesRejectStaleWorkflowTerms(t *testing.T) {
	root := filepath.Join("..", "..")
	liveAgents := readTemplate(t, root, "AGENTS.md")
	agents := readTemplate(t, root, "templates", "project", "AGENTS.md")
	router := readTemplate(t, root, "templates", "project", ".savepoint", "router.md")
	liveBuildSkill := readTemplate(t, root, "agent-skills", "savepoint-build-task", "SKILL.md")
	buildSkill := readTemplate(t, root, "templates", "project", "agent-skills", "savepoint-build-task", "SKILL.md")

	for _, content := range []string{liveAgents, agents, router, liveBuildSkill, buildSkill} {
		assertNotContains(t, content, "status: todo")
		assertNotContains(t, content, "status: doing")
		assertNotContains(t, content, "status: blocked")
		assertNotContains(t, content, "status: review")
		assertNotContains(t, content, "status: audit")
		assertNotContains(t, content, "phase: build")
		assertNotContains(t, content, "phase: test")
		assertNotContains(t, content, "phase: audit")
		assertNotContains(t, content, "phase: implementation")
		assertNotContains(t, content, "`phase` (build/test/audit)")
		assertNotContains(t, content, "prompt-based phase")
	}

	assertContains(t, agents, "Task `stage` (build/test/audit): **required** when `status: in_progress`")
	assertContains(t, buildSkill, "Set the task frontmatter to `status: in_progress` and `stage: build`")
	assertContains(t, buildSkill, "Never write `stage: implementation`; implementation work starts at `stage: build`.")
	assertContains(t, buildSkill, "legacy task `phase` as parser compatibility only")
}

func TestProjectConceptTemplateExists(t *testing.T) {
	root := filepath.Join("..", "..")
	concept := readTemplate(t, root, "templates", "project", ".savepoint", "Concept.md")

	assertContains(t, concept, "type: project-concept")
	assertContains(t, concept, "## When to use this")
	assertContains(t, concept, "## Core impulse")
	assertContains(t, concept, "## Target feeling")
	assertContains(t, concept, "## The problem in one sentence")
	assertContains(t, concept, "## Who this is NOT for")
	assertContains(t, concept, "## Open questions")

	for _, stale := range []string{
		"status: todo",
		"status: doing",
		"status: blocked",
		"status: review",
		"status: audit",
		"phase: build",
		"phase: test",
		"phase: audit",
		"phase: implementation",
		"`phase` (build/test/audit)",
	} {
		assertNotContains(t, concept, stale)
	}
}

func TestProjectDocumentTemplatesHaveTypeFrontmatter(t *testing.T) {
	root := filepath.Join("..", "..")
	cases := []struct {
		path []string
		want string
	}{
		{[]string{"templates", "project", ".savepoint", "PRD.md"}, "type: project-prd"},
		{[]string{"templates", "project", ".savepoint", "Design.md"}, "type: project-design"},
		{[]string{"templates", "project", ".savepoint", "Concept.md"}, "type: project-concept"},
	}
	for _, c := range cases {
		content := readTemplate(t, root, c.path...)
		assertContains(t, content, c.want)
	}
}

func TestProjectAgentsGuidesLifecycleTerminologyConsistency(t *testing.T) {
	root := filepath.Join("..", "..")
	liveAgents := readTemplate(t, root, "AGENTS.md")
	templateAgents := readTemplate(t, root, "templates", "project", "AGENTS.md")

	canonicalStatuses := []string{"planned", "in_progress", "done"}
	for _, status := range canonicalStatuses {
		assertContains(t, liveAgents, status)
		assertContains(t, templateAgents, status)
	}

	assertContains(t, liveAgents, "Task `status`: only `planned`, `in_progress`, or `done`")
	assertContains(t, templateAgents, "Task `status`: only `planned`, `in_progress`, or `done`")

	for _, content := range []string{liveAgents, templateAgents} {
		assertNotContains(t, content, "phase:")
	}
}

func readTemplate(t *testing.T, root string, parts ...string) string {
	t.Helper()

	path := filepath.Join(append([]string{root}, parts...)...)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read template %s: %v", path, err)
	}
	return string(data)
}

func assertFileMatches(t *testing.T, root, livePath, templatePath string) {
	t.Helper()

	live := readTemplate(t, root, livePath)
	template := readTemplate(t, root, templatePath)
	if live != template {
		t.Fatalf("%s does not match %s", templatePath, livePath)
	}
}

func assertContains(t *testing.T, content, want string) {
	t.Helper()

	if !strings.Contains(content, want) {
		t.Fatalf("template missing %q", want)
	}
}

func assertNotContains(t *testing.T, content, stale string) {
	t.Helper()

	if strings.Contains(content, stale) {
		t.Fatalf("template contains stale text %q", stale)
	}
}
