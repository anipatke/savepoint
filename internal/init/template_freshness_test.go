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
	auditSkill := readTemplate(t, root, "templates", "project", "agent-skills", "savepoint-audit-epic", "SKILL.md")

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
	var templateSkillCount int
	for _, entry := range templateEntries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), "savepoint-") {
			templateSkillCount++
		}
	}
	if templateSkillCount != len(skillNames) {
		t.Fatalf("scaffolded skill count = %d, want %d", templateSkillCount, len(skillNames))
	}

	// The shared audit method is a non-triggerable reference, so it is mirrored
	// outside the savepoint-* skill folders.
	assertFileMatches(t, root,
		filepath.Join("agent-skills", "references", "audit-method.md"),
		filepath.Join("templates", "project", "agent-skills", "references", "audit-method.md"),
	)
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

func TestProjectDesignTemplateListsGuardrailsAndHealthCheck(t *testing.T) {
	root := filepath.Join("..", "..")
	design := readTemplate(t, root, "templates", "project", ".savepoint", "Design.md")

	// The Design.md directory listing is what tells an agent these two files
	// exist, so a scaffolded file with no listing row is invisible in practice.
	assertContains(t, design, "Guardrails.md")
	assertContains(t, design, "Health-Check.md")
	assertContains(t, design, "engineering policy, severity model, rule index")
	assertContains(t, design, "Quick/Full/Deep evidence gates")

	// The listing must name the split audit skills, never the retired generic one.
	assertContains(t, design, "savepoint-audit-task")
	assertContains(t, design, "savepoint-audit-epic")
}

func TestProjectGuardrailsAndHealthCheckTemplatesExist(t *testing.T) {
	root := filepath.Join("..", "..")

	guardrails := readTemplate(t, root, "templates", "project", ".savepoint", "Guardrails.md")
	assertContains(t, guardrails, "type: guardrails")
	assertContains(t, guardrails, "## Severity Model")
	assertContains(t, guardrails, "| Blocker |")
	assertContains(t, guardrails, "## Rule Index")

	health := readTemplate(t, root, "templates", "project", ".savepoint", "Health-Check.md")
	assertContains(t, health, "type: health-check")
	assertContains(t, health, "## Modes")
	assertContains(t, health, "| Quick | Before task handoff |")
	assertContains(t, health, "| Full | Before epic audit closeout |")
	assertContains(t, health, "## Quick Check")
}

func TestUpgradeMigratesLegacyAuditSkillFromRealTemplates(t *testing.T) {
	root := filepath.Join("..", "..")
	templates := os.DirFS(filepath.Join(root, "templates", "project"))

	target := t.TempDir()
	if err := os.MkdirAll(filepath.Join(target, ".savepoint"), 0755); err != nil {
		t.Fatal(err)
	}
	legacyDir := filepath.Join(target, "agent-skills", "savepoint-audit")
	if err := os.MkdirAll(legacyDir, 0755); err != nil {
		t.Fatal(err)
	}
	legacy := "# Locally Edited Generic Audit Skill\n"
	if err := os.WriteFile(filepath.Join(legacyDir, "SKILL.md"), []byte(legacy), 0644); err != nil {
		t.Fatal(err)
	}

	report, err := UpgradeProjectAssets(templates, target, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	if got := upgradeActionFor(t, report, "agent-skills/savepoint-audit/SKILL.md"); got != ActionMigrated {
		t.Errorf("legacy skill action = %v, want migrated", got)
	}
	for _, path := range []string{
		"agent-skills/savepoint-audit-task/SKILL.md",
		"agent-skills/savepoint-audit-epic/SKILL.md",
		"agent-skills/references/audit-method.md",
	} {
		if got := upgradeActionFor(t, report, path); got != ActionUpdated {
			t.Errorf("%s action = %v, want updated", path, got)
		}
		if _, err := os.Stat(filepath.Join(target, filepath.FromSlash(path))); err != nil {
			t.Errorf("%s not installed by upgrade: %v", path, err)
		}
	}

	// The legacy content survives in the non-triggerable archive, and no
	// triggerable copy remains.
	archived, err := os.ReadFile(filepath.Join(target, ".savepoint", "migrations", "savepoint-audit-SKILL.md"))
	if err != nil {
		t.Fatalf("legacy content not archived: %v", err)
	}
	if string(archived) != legacy {
		t.Errorf("archived = %q, want verbatim %q", string(archived), legacy)
	}
	if _, err := os.Stat(filepath.Join(target, "agent-skills", "savepoint-audit", "SKILL.md")); !os.IsNotExist(err) {
		t.Errorf("legacy skill still triggerable after upgrade, stat err = %v", err)
	}
}

func TestProjectAuditRegisterTemplatesExist(t *testing.T) {
	root := filepath.Join("..", "..")

	prompt := readTemplate(t, root, "templates", "project", ".savepoint", "audit", "prompt.md")
	register := readTemplate(t, root, "templates", "project", ".savepoint", "audit", "register.md")
	findings := readTemplate(t, root, "templates", "project", ".savepoint", "audit", "findings", "README.md")
	runs := readTemplate(t, root, "templates", "project", ".savepoint", "audit", "runs", "README.md")

	// Prompt and register/findings guidance are seeded by T001/T002; verify they survive.
	assertContains(t, prompt, "register-backed Savepoint audit")
	assertContains(t, register, "Convergence summary")
	assertContains(t, register, "first real finding as `F001`")
	assertNotContains(t, register, "Example finding")
	assertContains(t, findings, "F###-slug.md")

	// Run history guidance: naming convention.
	assertContains(t, runs, "YYYY-MM-DD-label.md")

	// Run records require date, auditor/model, prompt version, commit SHA, mode,
	// coverage, source audits, and headline counts.
	for _, field := range []string{
		"date:",
		"auditor:",
		"model:",
		"prompt_version:",
		"commit:",
		"mode:",
		"coverage:",
		"source_audits:",
		"net_new:",
		"reopened:",
		"verified:",
		"deferred:",
		"coverage_gaps:",
	} {
		assertContains(t, runs, field)
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

func TestUpgradeAddsAuditRegisterTemplatesFromRealTemplates(t *testing.T) {
	root := filepath.Join("..", "..")
	templates := os.DirFS(filepath.Join(root, "templates", "project"))

	target := t.TempDir()
	if err := os.MkdirAll(filepath.Join(target, ".savepoint"), 0755); err != nil {
		t.Fatal(err)
	}

	auditAssets := []string{
		".savepoint/audit/prompt.md",
		".savepoint/audit/register.md",
		".savepoint/audit/findings/README.md",
		".savepoint/audit/runs/README.md",
	}

	report, err := UpgradeProjectAssets(templates, target, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	for _, path := range auditAssets {
		if got := upgradeActionFor(t, report, path); got != ActionUpdated {
			t.Errorf("audit asset %s action = %v, want updated", path, got)
		}
		if _, err := os.Stat(filepath.Join(target, filepath.FromSlash(path))); err != nil {
			t.Errorf("audit asset %s not written: %v", path, err)
		}
	}

	// A second upgrade must leave the now-present audit assets untouched.
	rerun, err := UpgradeProjectAssets(templates, target, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() rerun error = %v", err)
	}
	for _, path := range auditAssets {
		if got := upgradeActionFor(t, rerun, path); got != ActionUnchanged {
			t.Errorf("audit asset %s rerun action = %v, want unchanged", path, got)
		}
	}
}

func TestUpgradeDeliversPolicyAssetsFromRealTemplates(t *testing.T) {
	root := filepath.Join("..", "..")
	templates := os.DirFS(filepath.Join(root, "templates", "project"))

	target := t.TempDir()
	if err := os.MkdirAll(filepath.Join(target, ".savepoint"), 0755); err != nil {
		t.Fatal(err)
	}

	report, err := UpgradeProjectAssets(templates, target, false, false)
	if err != nil {
		t.Fatalf("UpgradeProjectAssets() error = %v", err)
	}

	for _, path := range []string{".savepoint/Guardrails.md", ".savepoint/Health-Check.md"} {
		if got := upgradeActionFor(t, report, path); got != ActionInstalled {
			t.Errorf("policy asset %s action = %v, want installed", path, got)
		}
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
