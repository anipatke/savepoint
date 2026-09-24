package init

import (
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"testing"
)

func TestProjectGuidanceTemplatesMirrorLiveGuidance(t *testing.T) {
	root := filepath.Join("..", "..")

	// Byte parity and set-completeness for the V2 skills and their shared
	// references (see .savepoint/Guardrails.md TPL-01), in both directions:
	// a skill added to the live tree and forgotten in the shipped tree still
	// fails, and so does a skill shipped to the wrong tree.
	assertSkillTreeParity(t, root, "agent-skills", "templates/project-v2/agent-skills", v2Skills, v2References)

	// Every live skill belongs to the declared V2 set; a skill added live and
	// forgotten there would otherwise pass silently.
	liveSkills := savepointSkillDirs(t, filepath.Join(root, "agent-skills"))
	declared := append([]string(nil), v2Skills...)
	sort.Strings(declared)
	if !slices.Equal(liveSkills, declared) {
		t.Errorf("live agent-skills/ = %v, want exactly v2Skills = %v", liveSkills, declared)
	}
}

func TestRouterFilesOmitRetiredNextActionKey(t *testing.T) {
	root := filepath.Join("..", "..")
	for _, path := range [][]string{
		{".savepoint", "router.md"},
		{"templates", "project-v2", ".savepoint", "router.md"},
	} {
		content := readTemplate(t, root, path...)
		assertNotContains(t, content, "next_action:")
	}
}

// assertSkillTreeParity asserts, for one shipped tree, that every named
// canonical skill and reference is byte-identical in that tree, and that the
// tree carries exactly that set — no fewer, and nothing belonging to another
// tree.
func assertSkillTreeParity(t *testing.T, root, liveRel, shippedRel string, skillNames, referenceNames []string) {
	t.Helper()

	for _, name := range skillNames {
		assertFileMatches(t, root,
			filepath.Join(liveRel, name, "SKILL.md"),
			filepath.Join(shippedRel, name, "SKILL.md"),
		)
	}
	for _, name := range referenceNames {
		assertFileMatches(t, root,
			filepath.Join(liveRel, "references", name),
			filepath.Join(shippedRel, "references", name),
		)
	}

	gotSkills := savepointSkillDirs(t, filepath.Join(root, filepath.FromSlash(shippedRel)))
	wantSkills := append([]string(nil), skillNames...)
	sort.Strings(wantSkills)
	if !slices.Equal(gotSkills, wantSkills) {
		t.Errorf("%s skill set = %v, want %v", shippedRel, gotSkills, wantSkills)
	}

	gotReferences := referenceFileNames(t, filepath.Join(root, filepath.FromSlash(shippedRel)))
	wantReferences := append([]string(nil), referenceNames...)
	sort.Strings(wantReferences)
	if !slices.Equal(gotReferences, wantReferences) {
		t.Errorf("%s reference set = %v, want %v", shippedRel, gotReferences, wantReferences)
	}
}

// referenceFileNames lists the non-triggerable reference files shipped
// alongside a skill root's savepoint-* directories.
func referenceFileNames(t *testing.T, skillRoot string) []string {
	t.Helper()

	entries, err := os.ReadDir(filepath.Join(skillRoot, "references"))
	if err != nil {
		t.Fatalf("read references dir under %s: %v", skillRoot, err)
	}

	var names []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	return names
}

func TestProjectTemplatesRejectStaleWorkflowTerms(t *testing.T) {
	root := filepath.Join("..", "..")
	liveAgents := readTemplate(t, root, "AGENTS.md")

	for _, content := range []string{liveAgents} {
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
}

func TestV2WorkflowAssetsHaveNoActiveV1Routing(t *testing.T) {
	root := filepath.Join("..", "..")
	paths := []string{
		filepath.Join(root, "agent-skills", "savepoint-idea", "SKILL.md"),
		filepath.Join(root, "agent-skills", "savepoint-design", "SKILL.md"),
		filepath.Join(root, "agent-skills", "savepoint-task", "SKILL.md"),
		filepath.Join(root, "agent-skills", "savepoint-check", "SKILL.md"),
		filepath.Join(root, "agent-skills", "references", "check-method.md"),
		filepath.Join(root, "agent-skills", "references", "issue-capture.md"),
		filepath.Join(root, "agent-skills", "references", "commands-and-procedures.md"),
		filepath.Join(root, "templates", "project-v2", "agent-skills", "savepoint-idea", "SKILL.md"),
		filepath.Join(root, "templates", "project-v2", "agent-skills", "savepoint-design", "SKILL.md"),
		filepath.Join(root, "templates", "project-v2", "agent-skills", "savepoint-task", "SKILL.md"),
		filepath.Join(root, "templates", "project-v2", "agent-skills", "savepoint-check", "SKILL.md"),
		filepath.Join(root, "templates", "project-v2", "agent-skills", "references", "check-method.md"),
		filepath.Join(root, "templates", "project-v2", "agent-skills", "references", "issue-capture.md"),
		filepath.Join(root, "templates", "project-v2", "agent-skills", "references", "commands-and-procedures.md"),
	}
	forbidden := []string{
		"pre-implementation",
		"epic-design",
		"epic-task-breakdown",
		"task-building",
		"audit-pending",
		"defect-building",
		"savepoint-draft-prd",
		"savepoint-system-design",
		"savepoint-create-task",
		"savepoint-build-task",
		"savepoint-audit-epic",
		"savepoint-audit-task",
		"savepoint-audit-register",
		"savepoint-create-defect",
	}
	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read V2 workflow asset %s: %v", path, err)
		}
		for _, stale := range forbidden {
			if strings.Contains(string(content), stale) {
				t.Errorf("V2 workflow asset %s contains active V1 route %q", path, stale)
			}
		}
	}
}

func TestProjectGuardrailsTemplateExists(t *testing.T) {
	root := filepath.Join("..", "..")

	guardrails := readTemplate(t, root, "templates", "project-v2", ".savepoint", "Guardrails.md")
	assertContains(t, guardrails, "type: guardrails")
	assertContains(t, guardrails, "## Severity Model")
	assertContains(t, guardrails, "| Blocker |")
	assertContains(t, guardrails, "## Rule Index")
}

func TestProjectDocumentTemplatesHaveTypeFrontmatter(t *testing.T) {
	root := filepath.Join("..", "..")
	cases := []struct {
		path []string
		want string
	}{
		{[]string{"templates", "project-v2", ".savepoint", "Design.md"}, "type: project-design"},
		{[]string{"templates", "project-v2", ".savepoint", "Idea.md"}, "type: idea"},
		{[]string{"templates", "project-v2", ".savepoint", "Guardrails.md"}, "type: guardrails"},
	}
	for _, c := range cases {
		content := readTemplate(t, root, c.path...)
		assertContains(t, content, c.want)
	}
}

func TestProjectAgentsGuidesLifecycleTerminologyConsistency(t *testing.T) {
	root := filepath.Join("..", "..")
	liveAgents := readTemplate(t, root, "AGENTS.md")
	templateAgents := readTemplate(t, root, "templates", "project-v2", "AGENTS.md")

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

func TestUpgradeDeliversPolicyAssetsFromRealTemplates(t *testing.T) {
	root := filepath.Join("..", "..")
	templates := os.DirFS(filepath.Join(root, "templates", "project-v2"))

	target := t.TempDir()
	if err := os.MkdirAll(filepath.Join(target, ".savepoint"), 0755); err != nil {
		t.Fatal(err)
	}

	report, err := upgradeAssetsFromTree(templates, target, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	if got := upgradeActionFor(t, report, ".savepoint/Guardrails.md"); got != ActionInstalled {
		t.Errorf("policy asset .savepoint/Guardrails.md action = %v, want installed", got)
	}

	// End to end: the AGENTS.md code-style pointer must resolve in an upgraded
	// project, all the way to the STYLE rules in the installed Guardrails.md.
	agents := readTemplate(t, target, "AGENTS.md")
	assertContains(t, agents, "## Code Style")
	assertContains(t, agents, "`STYLE` rules in `.savepoint/Guardrails.md`")

	guardrails := readTemplate(t, target, ".savepoint", "Guardrails.md")
	assertContains(t, guardrails, "STYLE-01")
	assertContains(t, guardrails, "STYLE-10")
	assertNotContains(t, guardrails, "{{PROJECT_NAME}}")
}

func upgradeActionFor(t *testing.T, report *UpgradeReport, path string) UpgradeAction {
	t.Helper()

	for _, e := range report.Actions {
		if e.Path == path {
			return e.Action
		}
	}
	t.Fatalf("path %s not in upgrade report", path)
	return ""
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
