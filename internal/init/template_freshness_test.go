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

func TestIdeaGuidanceFillsTheFreshProjectsGoal(t *testing.T) {
	root := filepath.Join("..", "..")
	for _, path := range []string{
		filepath.Join(root, "agent-skills", "savepoint-idea", "SKILL.md"),
		filepath.Join(root, "templates", "project-v2", "agent-skills", "savepoint-idea", "SKILL.md"),
	} {
		content := readTemplate(t, path)
		for _, phrase := range []string{
			"Every Savepoint project has at least one live Goal selected by the router",
			"G-001, titled after the project",
			"Outcome, Why, Success Conditions, and Boundaries",
			"do not create another Goal for the same initial outcome",
			"only fill the fresh scaffold placeholder from owner-provided answers",
		} {
			assertContains(t, content, phrase)
		}
		assertNotContains(t, content, "Goal is optional planning context")
	}
}

func TestProjectGuidanceRequiresGoalContext(t *testing.T) {
	root := filepath.Join("..", "..")
	guidance := [][]string{
		{"AGENTS.md"},
		{"templates", "project-v2", "AGENTS.md"},
		{".savepoint", "Design.md"},
		{".savepoint", "router.md"},
		{"templates", "project-v2", ".savepoint", "router.md"},
		{"README.md"},
	}
	for _, skill := range []string{"savepoint-idea", "savepoint-design", "savepoint-task", "savepoint-check"} {
		guidance = append(guidance,
			[]string{"agent-skills", skill, "SKILL.md"},
			[]string{"templates", "project-v2", "agent-skills", skill, "SKILL.md"},
		)
	}
	stalePhrases := []string{
		"goal is optional",
		"goal is not required",
		"goal is not mandatory",
		"goal not required",
		"goal may be omitted",
		"goal can be omitted",
		"goal may be left unassigned",
		"goals are optional",
		"goals are not required",
		"optional goal",
		"goal context is optional",
		"goal selection is optional",
		"release is optional",
		"release may be omitted",
		"release can be omitted",
		"goal check is mandatory whenever a goal exists",
		"goal check is mandatory when a goal exists",
		"when a goal exists, a mandatory full goal check",
		"whenever a goal exists",
		"without a goal, continue",
		"without a goal",
		"a v2 project may have no goal",
		"goal may be omitted",
		"goal is not required",
		"omit `release`",
		"optional `release:",
		"intentionally unassigned",
	}
	for _, parts := range guidance {
		content := readTemplate(t, root, parts...)
		lower := strings.ToLower(content)
		for _, phrase := range stalePhrases {
			if strings.Contains(lower, phrase) {
				t.Errorf("%s still contains optional-Goal wording %q", filepath.Join(parts...), phrase)
			}
		}
	}

	goalGuides := [][]string{
		{"AGENTS.md"},
		{"templates", "project-v2", "AGENTS.md"},
		{"agent-skills", "savepoint-idea", "SKILL.md"},
		{"templates", "project-v2", "agent-skills", "savepoint-idea", "SKILL.md"},
		{"agent-skills", "savepoint-design", "SKILL.md"},
		{"templates", "project-v2", "agent-skills", "savepoint-design", "SKILL.md"},
		{"agent-skills", "savepoint-task", "SKILL.md"},
		{"templates", "project-v2", "agent-skills", "savepoint-task", "SKILL.md"},
		{"agent-skills", "savepoint-check", "SKILL.md"},
		{"templates", "project-v2", "agent-skills", "savepoint-check", "SKILL.md"},
	}
	for _, parts := range goalGuides {
		content := readTemplate(t, root, parts...)
		for _, phrase := range []string{"Choose a Goal", "savepoint doctor", "savepoint init", "savepoint migrate"} {
			assertContains(t, content, phrase)
		}
		assertContains(t, content, "G-001")
	}
	for _, parts := range [][]string{{"AGENTS.md"}, {"templates", "project-v2", "AGENTS.md"}} {
		content := readTemplate(t, root, parts...)
		assertContains(t, content, "Every Savepoint project must have a live Goal selected by the router")
		assertContains(t, content, "every live Objective must name exactly one Goal")
	}
	for _, parts := range [][]string{
		{"agent-skills", "savepoint-idea", "SKILL.md"},
		{"templates", "project-v2", "agent-skills", "savepoint-idea", "SKILL.md"},
	} {
		content := readTemplate(t, root, parts...)
		assertContains(t, content, "at least one live Goal selected by the router")
		assertContains(t, content, "every live Objective names exactly one Goal through `release:`")
	}
	for _, parts := range [][]string{
		{"agent-skills", "savepoint-design", "SKILL.md"},
		{"templates", "project-v2", "agent-skills", "savepoint-design", "SKILL.md"},
	} {
		content := readTemplate(t, root, parts...)
		assertContains(t, content, "at least one live Goal selected by the router")
		assertContains(t, content, "every live Objective must name exactly one live Goal")
	}
	for _, parts := range [][]string{
		{"agent-skills", "savepoint-task", "SKILL.md"},
		{"templates", "project-v2", "agent-skills", "savepoint-task", "SKILL.md"},
		{"agent-skills", "savepoint-check", "SKILL.md"},
		{"templates", "project-v2", "agent-skills", "savepoint-check", "SKILL.md"},
	} {
		content := readTemplate(t, root, parts...)
		assertContains(t, content, "Every Savepoint project has at least one live Goal selected by the router")
		assertContains(t, content, "every live Objective names exactly one Goal through `release:`")
	}

	liveAgents := readTemplate(t, root, "AGENTS.md")
	templateAgents := readTemplate(t, root, "templates", "project-v2", "AGENTS.md")
	section := func(content string) string {
		t.Helper()
		const heading = "## Required Goal Context\n"
		start := strings.Index(content, heading)
		if start < 0 {
			t.Fatalf("missing %q section", strings.TrimSpace(heading))
		}
		body := content[start+len(heading):]
		if end := strings.Index(body, "\n## "); end >= 0 {
			body = body[:end]
		}
		return body
	}
	if live, scaffold := section(liveAgents), section(templateAgents); live != scaffold {
		t.Error("live and scaffold Required Goal Context sections differ")
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
		assertContains(t, content, "Exception: agents may run `savepoint create-task --objective O-### --draft <path> [dir]` only to create a new Task from an ID-free draft.")
		assertContains(t, content, "No other `savepoint` command is for agents except the narrow Task creation operation below.")
		assertContains(t, content, "After creating or renaming any other identity-bearing V2 record, run `savepoint resume` to require strict loading of the full V2 index.")
	}

	for _, content := range []string{liveAgents, templateAgents} {
		assertNotContains(t, content, "phase:")
	}
}

func TestReadmeDocumentsTaskCreationCommand(t *testing.T) {
	readme := strings.Join(strings.Fields(readTemplate(t, filepath.Join("..", ".."), "README.md")), " ")
	for _, phrase := range []string{
		"`savepoint create-task --objective O-### --draft <path> [dir]`",
		"complete Task draft without an `id`",
		"assigns the next project-wide ID and path",
		"strict-loads the full V2 index before it reports success",
	} {
		if !strings.Contains(readme, phrase) {
			t.Errorf("README.md does not document Task creation behavior %q", phrase)
		}
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
	if _, err := os.Stat(filepath.Join(target, ".savepoint", "releases")); !os.IsNotExist(err) {
		t.Fatalf("upgrade-assets added a Goal to an existing project, releases stat err = %v", err)
	}

	if got := upgradeActionFor(t, report, ".savepoint/Guardrails.md"); got != ActionInstalled {
		t.Errorf("policy asset .savepoint/Guardrails.md action = %v, want installed", got)
	}

	// End to end: the AGENTS.md code-style pointer must resolve in an upgraded
	// project, all the way to the STYLE rules in the installed Guardrails.md.
	agents := readTemplate(t, target, "AGENTS.md")
	assertContains(t, agents, "## Code Style")
	assertContains(t, agents, "`STYLE` rules in `.savepoint/Guardrails.md`")
	assertContains(t, agents, "Exception: agents may run `savepoint create-task --objective O-### --draft <path> [dir]` only to create a new Task from an ID-free draft.")

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

// TestGuidancePlansAndBoundsWorktreeLanes locks I-074: the planner shapes and
// names parallel lanes, and a lane never writes the router or new identities.
func TestGuidancePlansAndBoundsWorktreeLanes(t *testing.T) {
	root := filepath.Join("..", "..")
	for _, parts := range [][]string{{"AGENTS.md"}, {"templates", "project-v2", "AGENTS.md"}} {
		content := readTemplate(t, root, parts...)
		for _, phrase := range []string{
			"## Worktree Lanes",
			"git rev-parse --git-common-dir",
			"Do not edit `.savepoint/router.md`",
			"Do not create Tasks, Checks, or Issues",
			"Full Objective Check run on the main branch after the lane merges",
		} {
			assertContains(t, content, phrase)
		}
	}
	for _, parts := range [][]string{
		{"agent-skills", "savepoint-design", "SKILL.md"},
		{"templates", "project-v2", "agent-skills", "savepoint-design", "SKILL.md"},
	} {
		content := readTemplate(t, root, parts...)
		assertContains(t, content, "no `depends_on` path between them and no overlapping Context Files")
		assertContains(t, content, "name the parallel lanes in plain words")
	}
	for _, parts := range [][]string{
		{"agent-skills", "savepoint-task", "SKILL.md"},
		{"templates", "project-v2", "agent-skills", "savepoint-task", "SKILL.md"},
	} {
		assertContains(t, readTemplate(t, root, parts...), "follow AGENTS.md's Worktree Lanes section")
	}
}
