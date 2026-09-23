package init

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// checkMethodSections are the named V2 method sections the shared check
// reference must carry forward from the V1 audit method, in V2 vocabulary.
var checkMethodSections = []string{
	"## Task Check And Objective Check Depth",
	"## Quick And Full Evidence Modes",
	"## Freeze The Check Scope",
	"## Turn Acceptance Into Invariants",
	"## Build The Mandatory Coverage Matrix",
	"## Workflow And Side-Effect Check Lock",
	"## Perform The Adversarial Pass",
	"## Re-check After Remediation",
	"## Summarize Materiality",
}

// registerVocabulary is V1 audit-register bookkeeping that must never leak
// into the V2 check method, which replaces the register with Check records
// and derived Issues.
var registerVocabulary = []string{
	".savepoint/audit/",
	"register",
	"F###",
	"audit run",
}

func TestSharedCheckMethodIsNonTriggerableReference(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "references", "check-method.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		if got := frontmatterField(content, "type"); got != "check-method-reference" {
			t.Errorf("%s: %s type = %q, want check-method-reference", tree, path, got)
		}
		if got := frontmatterField(content, "triggerable"); got != "false" {
			t.Errorf("%s: %s triggerable = %q, want false", tree, path, got)
		}
		if frontmatterField(content, "name") != "" {
			t.Errorf("%s: %s carries a skill name and would be discoverable as a skill", tree, path)
		}
	}
}

func TestSharedCheckMethodCarriesRequiredSections(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "references", "check-method.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		for _, heading := range checkMethodSections {
			body, found := sectionBody(content, heading)
			if !found {
				t.Errorf("%s: %s missing %s", tree, path, heading)
				continue
			}
			if body == "" {
				t.Errorf("%s: %s has an empty %s section", tree, path, heading)
			}
		}
	}
}

func TestSharedCheckMethodHasNoRegisterVocabulary(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "references", "check-method.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		for _, stale := range registerVocabulary {
			if strings.Contains(content, stale) {
				t.Errorf("%s: %s contains register vocabulary %q", tree, path, stale)
			}
		}
	}
}

func TestSharedCheckMethodLiveAndTemplateMatch(t *testing.T) {
	root := filepath.Join("..", "..")
	assertFileMatches(t, root,
		filepath.Join("agent-skills", "references", "check-method.md"),
		filepath.Join("templates", "project-v2", "agent-skills", "references", "check-method.md"),
	)
}

// ideaSkillForbiddenOutputs are the outputs the idea planner must name as
// out of bounds, so an idea session can't quietly start designing.
var ideaSkillForbiddenOutputs = []string{
	"Design",
	"Guardrails",
	"Objectives",
	"Tasks",
	"Checks",
	"Issues",
	"production code",
}

// ideaTemplateSections are the required headings of the Idea artifact
// template the planner writes into .savepoint/Idea.md.
var ideaTemplateSections = []string{
	"Intent",
	"User",
	"Core Experience",
	"Scope",
	"Out of Scope",
	"Success Criteria",
}

func TestSavepointIdeaSkillReadBoundary(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-idea", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}

		body, found := sectionBody(string(data), "## Read")
		if !found || body == "" {
			t.Errorf("%s: %s missing or empty ## Read section", tree, path)
			continue
		}
		for _, phrase := range []string{"router.md", ".savepoint/Idea.md", "stated intent", "existing-project evidence"} {
			if !strings.Contains(body, phrase) {
				t.Errorf("%s: %s ## Read section missing %q", tree, path, phrase)
			}
		}
	}
}

func TestSavepointIdeaSkillWriteBoundary(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-idea", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		if !strings.Contains(content, ".savepoint/Idea.md") {
			t.Errorf("%s: %s does not name .savepoint/Idea.md as a write target", tree, path)
		}
		if !strings.Contains(content, "routing handoff") {
			t.Errorf("%s: %s does not name the routing handoff as a write target", tree, path)
		}
		for _, forbidden := range ideaSkillForbiddenOutputs {
			if !strings.Contains(content, forbidden) {
				t.Errorf("%s: %s does not name %q as forbidden output", tree, path, forbidden)
			}
		}
		if !strings.Contains(content, "savepoint-design") {
			t.Errorf("%s: %s does not name savepoint-design as the handoff target", tree, path)
		}
		if !strings.Contains(content, "does not detail any Objective") {
			t.Errorf("%s: %s does not state that it leaves Objective detailing to savepoint-design", tree, path)
		}
	}
}

func TestSavepointIdeaSkillHasArtifactTemplate(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-idea", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		for _, heading := range ideaTemplateSections {
			if !strings.Contains(content, "## "+heading) {
				t.Errorf("%s: %s missing Idea template section %q", tree, path, heading)
			}
		}
	}
}

func TestSavepointIdeaSkillAcceptsRoughInputAndEscalates(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-idea", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		if !strings.Contains(content, "rough sentence") {
			t.Errorf("%s: %s does not state a rough sentence is a valid starting input", tree, path)
		}
		if !strings.Contains(content, "material product uncertainty") {
			t.Errorf("%s: %s does not state the material-product-uncertainty escalation rule", tree, path)
		}
	}
}

func TestSavepointIdeaSkillLiveAndTemplateMatch(t *testing.T) {
	root := filepath.Join("..", "..")
	assertFileMatches(t, root,
		filepath.Join("agent-skills", "savepoint-idea", "SKILL.md"),
		filepath.Join("templates", "project-v2", "agent-skills", "savepoint-idea", "SKILL.md"),
	)
}

// designStructureHeadings are the required non-empty headings for
// savepoint-design, mirroring the split audit skill structure check.
var designStructureHeadings = []string{"## Purpose", "## Trigger", "## Read", "## Workflow", "## Rules"}

func TestSavepointDesignSkillPassesStructureValidation(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-design", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		if got := frontmatterField(content, "name"); got != "savepoint-design" {
			t.Errorf("%s: %s frontmatter name = %q, want savepoint-design", tree, path, got)
		}
		for _, heading := range designStructureHeadings {
			body, found := sectionBody(content, heading)
			if !found {
				t.Errorf("%s: %s missing %s", tree, path, heading)
				continue
			}
			if body == "" {
				t.Errorf("%s: %s has an empty %s section", tree, path, heading)
			}
		}
	}
}

func TestSavepointDesignSkillTrigger(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-design", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		body, found := sectionBody(string(data), "## Trigger")
		if !found {
			t.Errorf("%s: %s missing ## Trigger", tree, path)
			continue
		}
		if !strings.Contains(body, "state` is `design`") {
			t.Errorf("%s: %s trigger does not name router state design", tree, path)
		}
		if !strings.Contains(body, "REPLAN REQUIRED") {
			t.Errorf("%s: %s trigger does not name REPLAN REQUIRED routing back into this skill", tree, path)
		}
	}
}

// designWriteTargets and designForbiddenOutputs pin the write boundary named
// in the Objective/Design contract: only these outputs, never production
// code or backlog beyond the next Objective.
var designWriteTargets = []string{"Design", "Guardrails", "Objective", "routing"}
var designForbiddenOutputs = []string{"production code"}

func TestSavepointDesignSkillWriteBoundary(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-design", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		for _, target := range designWriteTargets {
			if !strings.Contains(content, target) {
				t.Errorf("%s: %s does not name %q as a write target", tree, path, target)
			}
		}
		for _, forbidden := range designForbiddenOutputs {
			if !strings.Contains(content, forbidden) {
				t.Errorf("%s: %s does not name %q as forbidden output", tree, path, forbidden)
			}
		}
		if !strings.Contains(content, "next Objective") {
			t.Errorf("%s: %s does not scope detailed Tasks to the next Objective", tree, path)
		}
		if !strings.Contains(content, "beyond the next Objective") {
			t.Errorf("%s: %s does not forbid detailed backlog beyond the next Objective", tree, path)
		}
	}
}

var designObjectiveFrontmatterFields = []string{"id", "title", "status", "depends_on", "release", "last_check", "freshness"}
var designObjectiveBodySections = []string{"Outcome", "Why", "Success Conditions", "Architectural Considerations", "Boundaries"}

func TestSavepointDesignSkillObjectiveTemplate(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-design", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		for _, field := range designObjectiveFrontmatterFields {
			if !strings.Contains(content, field+":") {
				t.Errorf("%s: %s Objective template missing frontmatter field %q", tree, path, field)
			}
		}
		if !strings.Contains(content, "status: planned|in_progress|done") {
			t.Errorf("%s: %s Objective template does not constrain status to planned|in_progress|done", tree, path)
		}
		if !strings.Contains(content, "depends_on: [O-###]") {
			t.Errorf("%s: %s Objective template does not show depends_on: [O-###]", tree, path)
		}
		for _, heading := range designObjectiveBodySections {
			if !strings.Contains(content, "## "+heading) {
				t.Errorf("%s: %s Objective template missing body section %q", tree, path, heading)
			}
		}
		if !strings.Contains(content, "derived from which Tasks name this Objective") {
			t.Errorf("%s: %s does not state Task membership is derived from Task ownership", tree, path)
		}
		if !strings.Contains(content, "second, manually kept list") {
			t.Errorf("%s: %s does not forbid a second manually maintained membership list", tree, path)
		}
	}
}

var designTaskFrontmatterFields = []string{"id", "title", "objective", "status", "depends_on", "owner_validation", "planned_by"}
var designTaskBodySections = []string{"Outcome", "User Check", "Done When", "Context Files", "Design References", "Guardrails", "Implementation Plan", "Boundaries", "Technical Verification", "Technical Evidence", "Drift Notes"}

func TestSavepointDesignSkillTaskTemplate(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-design", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		for _, field := range designTaskFrontmatterFields {
			if !strings.Contains(content, field+":") {
				t.Errorf("%s: %s Task template missing frontmatter field %q", tree, path, field)
			}
		}
		if !strings.Contains(content, "depends_on: [{task: T-013, requires: clear}]") {
			t.Errorf("%s: %s Task template does not show depends_on in {task: T-###, requires: clear} form", tree, path)
		}
		if !strings.Contains(content, "owner_validation: {required: true}") {
			t.Errorf("%s: %s Task template does not show owner_validation: {required: true|false}", tree, path)
		}
		for _, heading := range designTaskBodySections {
			if !strings.Contains(content, "## "+heading) {
				t.Errorf("%s: %s Task template missing body section %q", tree, path, heading)
			}
		}
		if !strings.Contains(content, "no globs, no directory-only entries") {
			t.Errorf("%s: %s Task template does not state Context Files must name exact paths", tree, path)
		}
	}
}

func TestSavepointDesignSkillTaskTitleNoReuseRule(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-design", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		if !strings.Contains(content, "must never be the Outcome text, a truncation of it, or a restatement of the technical objective") {
			t.Errorf("%s: %s does not state the title no-reuse rule", tree, path)
		}
		if !strings.Contains(content, "always an `O-###` reference to the owning Objective, never free text") {
			t.Errorf("%s: %s does not document objective as an O-### reference, never free text", tree, path)
		}
		if !strings.Contains(content, "every Task belongs to exactly one Objective") {
			t.Errorf("%s: %s does not state one-Objective Task membership", tree, path)
		}
		if !strings.Contains(content, "E50") {
			t.Errorf("%s: %s does not note title readability is evaluated by E50 agent scenarios", tree, path)
		}
	}
}

var designSectionHeadings = []string{"Architecture", "Components/Codebase Map", "Interfaces and Data Flow", "Boundaries", "Decisions", "Current Technical State"}

func TestSavepointDesignSkillDesignTemplate(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-design", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		for _, heading := range designSectionHeadings {
			if !strings.Contains(content, heading) {
				t.Errorf("%s: %s Design template missing section %q", tree, path, heading)
			}
		}
		if !strings.Contains(content, "Design describes implemented reality") {
			t.Errorf("%s: %s does not state Design describes implemented reality", tree, path)
		}
		if !strings.Contains(content, "Objective's deltas") || !strings.Contains(content, "reconciled") {
			t.Errorf("%s: %s does not state Objective deltas hold planned change until reconciliation", tree, path)
		}
	}
}

func TestSavepointDesignSkillGuardrailsRule(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-design", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		for _, phrase := range []string{"durable project constraints only", "10–20 substantive rules", "stable category IDs", "owns severity and exception authority", "must not duplicate guardrail prose"} {
			if !strings.Contains(content, phrase) {
				t.Errorf("%s: %s does not state guardrails rule phrase %q", tree, path, phrase)
			}
		}
	}
}

func TestSavepointDesignSkillReadinessGate(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-design", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		for _, phrase := range []string{"Interfaces and data ownership", "settled", "Constraints on the work are scoped", "outcomes of any Task or Objective this one depends on are known", "verification approach for the work"} {
			if !strings.Contains(content, phrase) {
				t.Errorf("%s: %s does not state readiness gate phrase %q", tree, path, phrase)
			}
		}
	}
}

func TestSavepointDesignSkillOneObjectiveResearchAndSplitRules(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-design", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		if !strings.Contains(content, "Keep exactly one Objective active") {
			t.Errorf("%s: %s does not state the one-Objective rule", tree, path)
		}
		if !strings.Contains(content, "bounded research Task with a named decision deliverable") {
			t.Errorf("%s: %s does not state the research-Task rule", tree, path)
		}
		if !strings.Contains(content, "Split a Task that carries multiple unrelated outcomes or an unresolved architectural decision") {
			t.Errorf("%s: %s does not state the split rule", tree, path)
		}
		if !strings.Contains(content, "Route product choices to the owner") || !strings.Contains(content, "does not require the owner to review code") {
			t.Errorf("%s: %s does not state the product-choice/technical-readiness rule", tree, path)
		}
	}
}

func TestSavepointDesignSkillLiveAndTemplateMatch(t *testing.T) {
	root := filepath.Join("..", "..")
	assertFileMatches(t, root,
		filepath.Join("agent-skills", "savepoint-design", "SKILL.md"),
		filepath.Join("templates", "project-v2", "agent-skills", "savepoint-design", "SKILL.md"),
	)
}

// taskStructureHeadings are the required non-empty headings for
// savepoint-task, mirroring the split audit skill structure check.
var taskStructureHeadings = []string{"## Purpose", "## Trigger", "## Read", "## Workflow", "## Rules"}

func TestSavepointTaskSkillPassesStructureValidation(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-task", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		if got := frontmatterField(content, "name"); got != "savepoint-task" {
			t.Errorf("%s: %s frontmatter name = %q, want savepoint-task", tree, path, got)
		}
		for _, heading := range taskStructureHeadings {
			body, found := sectionBody(content, heading)
			if !found {
				t.Errorf("%s: %s missing %s", tree, path, heading)
				continue
			}
			if body == "" {
				t.Errorf("%s: %s has an empty %s section", tree, path, heading)
			}
		}
	}
}

func TestSavepointTaskSkillTrigger(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-task", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		body, found := sectionBody(string(data), "## Trigger")
		if !found {
			t.Errorf("%s: %s missing ## Trigger", tree, path)
			continue
		}
		if !strings.Contains(body, "state` is `task`") {
			t.Errorf("%s: %s trigger does not name router state task", tree, path)
		}
	}
}

func TestSavepointTaskSkillReadBudget(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-task", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		body, found := sectionBody(string(data), "## Read")
		if !found || body == "" {
			t.Errorf("%s: %s missing or empty ## Read section", tree, path)
			continue
		}
		for _, phrase := range []string{"router.md", "active Task", "owning Objective's boundaries", "Context Files", "Applicable policy", "the read budget"} {
			if !strings.Contains(body, phrase) {
				t.Errorf("%s: %s ## Read section missing %q", tree, path, phrase)
			}
		}
	}
}

var taskWriteTargets = []string{"scoped implementation", "recorded evidence", "lifecycle progress", "replan handoff"}
var taskForbiddenClaims = []string{"edit the Task's acceptance criteria", "write a Check record", "close an Issue", "claim clearance or owner acceptance"}

func TestSavepointTaskSkillWriteBoundaryAndForbiddenClaims(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-task", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		for _, target := range taskWriteTargets {
			if !strings.Contains(content, target) {
				t.Errorf("%s: %s does not name %q as a write target", tree, path, target)
			}
		}
		for _, forbidden := range taskForbiddenClaims {
			if !strings.Contains(content, forbidden) {
				t.Errorf("%s: %s does not name %q as forbidden", tree, path, forbidden)
			}
		}
	}
}

func TestSavepointTaskSkillLifecycleStages(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-task", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		for _, phrase := range []string{
			"in_progress` with `stage: build`",
			"`build` → `test` → `audit`",
			"ready for a Check",
			"does not mean it passed",
		} {
			if !strings.Contains(content, phrase) {
				t.Errorf("%s: %s does not state lifecycle phrase %q", tree, path, phrase)
			}
		}
	}
}

func TestSavepointTaskSkillBlockedStartReported(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-task", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		for _, phrase := range []string{"dependencies are satisfied", "owning Objective is ready", "blocked start is reported"} {
			if !strings.Contains(content, phrase) {
				t.Errorf("%s: %s does not state %q", tree, path, phrase)
			}
		}
	}
}

func TestSavepointTaskSkillExtraReadsLogged(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-task", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		if !strings.Contains(content, "logged with what was read and why") {
			t.Errorf("%s: %s does not require extra reads to be logged with what and why", tree, path)
		}
	}
}

func TestSavepointTaskSkillReplanRequiredContract(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-task", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		for _, phrase := range []string{
			"REPLAN REQUIRED",
			"current `status` and `stage` unchanged",
			"preserves whatever partial work",
			"hands control to the planner",
			"never silently redesigned",
		} {
			if !strings.Contains(content, phrase) {
				t.Errorf("%s: %s does not state replan phrase %q", tree, path, phrase)
			}
		}
	}
}

func TestSavepointTaskSkillEvidenceRequirement(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-task", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		for _, phrase := range []string{
			"per-criterion outcome",
			"make build && make test",
			"files read and the files changed",
			"stated limitations",
		} {
			if !strings.Contains(content, phrase) {
				t.Errorf("%s: %s does not state evidence phrase %q", tree, path, phrase)
			}
		}
	}
}

func TestSavepointTaskSkillFreshCheckHandoff(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-task", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		if !strings.Contains(content, "fresh `savepoint-check` session") {
			t.Errorf("%s: %s does not name a fresh savepoint-check session as the handoff target", tree, path)
		}
		if !strings.Contains(content, "can never be that Check") {
			t.Errorf("%s: %s does not state the executor's own session can never be the Check", tree, path)
		}
	}
}

func TestSavepointTaskSkillObjectiveCheckNeedsWorkDoesNotRetreatDoneTasks(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-task", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		if !strings.Contains(content, "never retreats a Task that is already `done`") {
			t.Errorf("%s: %s does not state an Objective/Release Check's NEEDS WORK never retreats a done Task", tree, path)
		}
	}
}

func TestSavepointTaskSkillStyleAdvisoryRule(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-task", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		if !strings.Contains(content, "STYLE` guardrail rules as advisory") {
			t.Errorf("%s: %s does not treat STYLE guardrail rules as advisory", tree, path)
		}
		if !strings.Contains(content, "guardrail rule IDs rather than restating rule prose") {
			t.Errorf("%s: %s does not reference guardrail rule IDs instead of restating rule prose", tree, path)
		}
	}
}

func TestSavepointTaskSkillLiveAndTemplateMatch(t *testing.T) {
	root := filepath.Join("..", "..")
	assertFileMatches(t, root,
		filepath.Join("agent-skills", "savepoint-task", "SKILL.md"),
		filepath.Join("templates", "project-v2", "agent-skills", "savepoint-task", "SKILL.md"),
	)
}

func TestSavepointBuildTaskSkillPresentInBothTrees(t *testing.T) {
	for tree, root := range skillRoots() {
		path := filepath.Join(root, "savepoint-build-task", "SKILL.md")
		if _, err := os.Stat(path); err != nil {
			t.Errorf("%s: savepoint-build-task/SKILL.md missing: %v", tree, err)
		}
	}
}

// checkStructureHeadings are the required non-empty headings for
// savepoint-check, mirroring the split audit skill structure check.
var checkStructureHeadings = []string{"## Purpose", "## Trigger", "## Read", "## Workflow", "## Rules"}

func TestSavepointCheckSkillPassesStructureValidation(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-check", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		if got := frontmatterField(content, "name"); got != "savepoint-check" {
			t.Errorf("%s: %s frontmatter name = %q, want savepoint-check", tree, path, got)
		}
		for _, heading := range checkStructureHeadings {
			body, found := sectionBody(content, heading)
			if !found {
				t.Errorf("%s: %s missing %s", tree, path, heading)
				continue
			}
			if body == "" {
				t.Errorf("%s: %s has an empty %s section", tree, path, heading)
			}
		}
	}
}

func TestSavepointCheckSkillFreshSessionRule(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-check", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		body, found := sectionBody(string(data), "## Trigger")
		if !found {
			t.Errorf("%s: %s missing ## Trigger", tree, path)
			continue
		}
		for _, phrase := range []string{"state` is `check`", "independent from the executor's conversation", "same model is allowed", "same session is not", "Model names are optional"} {
			if !strings.Contains(body, phrase) {
				t.Errorf("%s: %s trigger missing fresh-session phrase %q", tree, path, phrase)
			}
		}
	}
}

func TestSavepointCheckSkillLoadsMethodWithoutRestating(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-check", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		if !strings.Contains(content, "references/check-method.md` in full") {
			t.Errorf("%s: %s does not state it loads check-method.md in full", tree, path)
		}
		if !strings.Contains(content, "does not restate that method") {
			t.Errorf("%s: %s does not state it avoids restating the shared method", tree, path)
		}
	}
}

var checkWriteTargets = []string{"the Check record", "Issues", "evaluation metadata", "authorized closure"}
var checkForbiddenActions = []string{"repair implementation", "edit acceptance criteria", "update Design as a form of remediation"}

func TestSavepointCheckSkillWriteBoundaryAndForbiddenActions(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-check", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		for _, target := range checkWriteTargets {
			if !strings.Contains(content, target) {
				t.Errorf("%s: %s does not name %q as a write target", tree, path, target)
			}
		}
		for _, forbidden := range checkForbiddenActions {
			if !strings.Contains(content, forbidden) {
				t.Errorf("%s: %s does not name %q as forbidden", tree, path, forbidden)
			}
		}
	}
}

var checkTemplateFields = []string{"id: C-###", "scope: {kind: task|objective|release, id:", "result: CLEAR|NEEDS WORK", "checked_by:", "executed_session:", "checked_at:", "reviewed:", "files:", "dependencies:", "issues:", "supersedes:"}

func TestSavepointCheckSkillArtifactTemplate(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-check", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		for _, field := range checkTemplateFields {
			if !strings.Contains(content, field) {
				t.Errorf("%s: %s Check template missing field %q", tree, path, field)
			}
		}
		if !strings.Contains(content, "immutable Check record") {
			t.Errorf("%s: %s does not state each run is a new immutable Check record", tree, path)
		}
		if !strings.Contains(content, "never edits the superseded record") {
			t.Errorf("%s: %s does not state a recheck never edits the superseded record", tree, path)
		}
	}
}

func TestSavepointCheckSkillScopes(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-check", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		if !strings.Contains(content, "a Task Check evaluates one Task's outcome and evidence") {
			t.Errorf("%s: %s does not scope a Task Check to one Task's outcome and evidence", tree, path)
		}
		if !strings.Contains(content, "integration across the Objective's owned Tasks and reconciliation against Design") {
			t.Errorf("%s: %s does not scope an Objective Check to integration and Design reconciliation", tree, path)
		}
		if !strings.Contains(content, "a Release Check uses `scope.kind: release` to evaluate integration across all member Objectives") {
			t.Errorf("%s: %s does not scope a Release Check to cross-Objective integration", tree, path)
		}
	}
}

func TestSavepointCheckSkillClosureRules(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-check", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		for _, phrase := range []string{
			"complete a technical Task",
			"no unexcepted material blocker",
			"owner_validation.required` additionally needs the owner's recorded acceptance naming this same current Check",
			"Release `done` requires at least one member Objective",
			"The checker never supplies that acceptance",
			"cannot support completion",
			"Stale or unknown freshness blocks normal completion",
			"never waived through by re-asserting",
		} {
			if !strings.Contains(content, phrase) {
				t.Errorf("%s: %s does not state closure phrase %q", tree, path, phrase)
			}
		}
	}
}

func TestV2SkillsTeachOptionalReleaseWorkflow(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		idea := string(readSkillFile(t, root, "savepoint-idea"))
		for _, phrase := range []string{
			"Release is optional planning context",
			"navigable delivery/package promise",
			"Ask whether the owner needs a navigable delivery/package promise across multiple Objectives",
			"otherwise continue with Objective → Task",
		} {
			if !strings.Contains(idea, phrase) {
				t.Errorf("%s: savepoint-idea missing optional-Release phrase %q", tree, phrase)
			}
		}

		design := string(readSkillFile(t, root, "savepoint-design"))
		for _, phrase := range []string{
			"## Optional Release Boundary",
			"stable global `R-###` identity from the first unused number",
			"Release sections `Outcome`, `Why`, `Success Conditions`, and `Boundaries`",
			"one optional `release: R-###` field",
			"do not maintain a second membership list",
			"continue through Idea → Design → Task → Check with no missing-record error or extra phase",
			"current CLEAR integration evidence exists",
			"not whether it has been published or deployed",
		} {
			if !strings.Contains(design, phrase) {
				t.Errorf("%s: savepoint-design missing Release design phrase %q", tree, phrase)
			}
		}
		if strings.Contains(design, "release: optional-release-name") {
			t.Errorf("%s: savepoint-design retains the pre-E51 optional release-string placeholder", tree)
		}

		check := string(readSkillFile(t, root, "savepoint-check"))
		for _, phrase := range []string{
			"scope: {kind: task|objective|release, id: T-###, O-###, or R-###}",
			"cross-Objective integration",
			"creates or reuses ordinary Issues",
			"never records owner acceptance on the owner's behalf",
			"Release `done` does not mean published or deployed",
		} {
			if !strings.Contains(check, phrase) {
				t.Errorf("%s: savepoint-check missing Release Check phrase %q", tree, phrase)
			}
		}
	}
}

func readSkillFile(t *testing.T, root, skill string) []byte {
	t.Helper()

	data, err := os.ReadFile(filepath.Join(root, skill, "SKILL.md"))
	if err != nil {
		t.Fatalf("read %s/%s/SKILL.md: %v", root, skill, err)
	}
	return data
}

func TestSavepointCheckSkillNeedsWorkPath(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-check", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		if !strings.Contains(content, "hand remediation back to the executor or planner") {
			t.Errorf("%s: %s does not hand NEEDS WORK remediation back to executor or planner", tree, path)
		}
		if !strings.Contains(content, "A Task Check's `NEEDS WORK` resumes the executor at `stage: build` inside that same Task") {
			t.Errorf("%s: %s does not state a Task Check's NEEDS WORK resumes the executor at stage: build inside that same Task", tree, path)
		}
		if !strings.Contains(content, "must not retreat a Task that is already `done`") && !strings.Contains(content, "must never retreat a Task that is already `done`") {
			t.Errorf("%s: %s does not state an Objective/Release Check's NEEDS WORK must not retreat a done Task", tree, path)
		}
	}
}

func TestSavepointCheckSkillAdvisoryObservationsNonBlocking(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-check", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		if !strings.Contains(content, "STYLE` guardrail rules, as non-blocking") {
			t.Errorf("%s: %s does not treat STYLE guardrail rules as non-blocking", tree, path)
		}
	}
}

func TestSavepointCheckSkillLiveAndTemplateMatch(t *testing.T) {
	root := filepath.Join("..", "..")
	assertFileMatches(t, root,
		filepath.Join("agent-skills", "savepoint-check", "SKILL.md"),
		filepath.Join("templates", "project-v2", "agent-skills", "savepoint-check", "SKILL.md"),
	)
}

// v1AuditSkillNames are the pre-existing V1 audit skills that must survive
// the V2 savepoint-check addition untouched in both trees.
var v1AuditSkillNames = []string{"savepoint-audit-task", "savepoint-audit-epic"}

func TestV1AuditSkillsAndSharedMethodPresentInBothTrees(t *testing.T) {
	for tree, root := range skillRoots() {
		for _, name := range v1AuditSkillNames {
			path := filepath.Join(root, name, "SKILL.md")
			if _, err := os.Stat(path); err != nil {
				t.Errorf("%s: %s missing: %v", tree, name, err)
			}
		}
		path := filepath.Join(root, "references", "audit-method.md")
		if _, err := os.Stat(path); err != nil {
			t.Errorf("%s: %s missing: %v", tree, path, err)
		}
	}
}

func issueCapturePath(root string) string {
	return filepath.Join(root, "references", "issue-capture.md")
}

func TestSharedIssueCaptureIsNonTriggerableReference(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := issueCapturePath(root)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		if got := frontmatterField(content, "type"); got != "issue-capture-reference" {
			t.Errorf("%s: %s type = %q, want issue-capture-reference", tree, path, got)
		}
		if got := frontmatterField(content, "triggerable"); got != "false" {
			t.Errorf("%s: %s triggerable = %q, want false", tree, path, got)
		}
		if frontmatterField(content, "name") != "" {
			t.Errorf("%s: %s carries a skill name and would be discoverable as a skill", tree, path)
		}
	}
}

// issueArtifactFields are the IssueV2 frontmatter fields the reference's
// template must name, matching internal/data/issue_v2.go exactly.
var issueArtifactFields = []string{
	"id:", "title:", "type: defect|drift|guardrail|verification|other",
	"status: open|in_progress|resolved", "source:", "tasks:", "checks:",
	"guardrail_ids:", "severity:", "resolution:", "duplicate_of:", "history:",
}

var issueArtifactBodySections = []string{"Summary", "Evidence", "Proof Needed"}

func TestSharedIssueCaptureArtifactTemplate(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := issueCapturePath(root)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		for _, field := range issueArtifactFields {
			if !strings.Contains(content, field) {
				t.Errorf("%s: %s Issue template missing field %q", tree, path, field)
			}
		}
		for _, heading := range issueArtifactBodySections {
			if !strings.Contains(content, "## "+heading) {
				t.Errorf("%s: %s Issue template missing body section %q", tree, path, heading)
			}
		}
	}
}

func TestSharedIssueCaptureTypeNeverBlocksAlone(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := issueCapturePath(root)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		if !strings.Contains(string(data), "never sufficient on its own to block a Task") {
			t.Errorf("%s: %s does not state type is never sufficient on its own to block a Task", tree, path)
		}
	}
}

func TestSharedIssueCaptureSearchFirstRule(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := issueCapturePath(root)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		for _, phrase := range []string{
			"same symptom, the same location, the same violated requirement, or the same linked work",
			"No automatic deduplication is assumed",
		} {
			if !strings.Contains(content, phrase) {
				t.Errorf("%s: %s does not state search-first phrase %q", tree, path, phrase)
			}
		}
	}
}

func TestSharedIssueCaptureDispositionsAndHistory(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := issueCapturePath(root)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		for _, phrase := range []string{
			"proven by a Check that recorded `CLEAR`",
			"owner decision, not a `CLEAR` Check or independent proof",
			"names that Issue and proves nothing itself",
			"reuses the same `I-###` with new, dated evidence",
			"{at, actor, kind, note, check}",
			"shortens, reorders, or edits a recorded entry is refused",
			"not a fourth lifecycle state",
		} {
			if !strings.Contains(content, phrase) {
				t.Errorf("%s: %s does not state disposition/history phrase %q", tree, path, phrase)
			}
		}
	}
}

func TestSharedIssueCaptureRoleBoundariesAndRepairRouting(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := issueCapturePath(root)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		for _, phrase := range []string{
			"executor** reports repair evidence without independently closing the Issue",
			"checker** verifies Check proof and closes the Issue as `verified`",
			"owner** may close the Issue as `accepted`",
			"planner** (`savepoint-design`) closes an Issue as `escalated`",
			"becomes a new, bounded Task",
			"within an existing Objective is not an escalation",
		} {
			if !strings.Contains(content, phrase) {
				t.Errorf("%s: %s does not state role/routing phrase %q", tree, path, phrase)
			}
		}
	}
}

func TestSharedIssueCaptureDefectWordMapping(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := issueCapturePath(root)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		for _, phrase := range []string{
			"\"Defect\" stays a word the user says",
			"maps to `type: defect`",
			"does not resurrect a separate defect record, status, or public phase in V2",
		} {
			if !strings.Contains(content, phrase) {
				t.Errorf("%s: %s does not state defect word-mapping phrase %q", tree, path, phrase)
			}
		}
	}
}

func TestSharedIssueCaptureLiveAndTemplateMatch(t *testing.T) {
	root := filepath.Join("..", "..")
	assertFileMatches(t, root,
		filepath.Join("agent-skills", "references", "issue-capture.md"),
		filepath.Join("templates", "project-v2", "agent-skills", "references", "issue-capture.md"),
	)
}

// issueCaptureEntrySkills maps each working skill's file to the phrase that
// must appear in its Issue capture entry point, distinct per role so the
// entry describes that skill's own write boundary rather than the shared
// reference's identical role-boundary prose.
var issueCaptureEntrySkills = map[string]string{
	"savepoint-design": "may read and reference an Issue; it\ndoes not close one",
	"savepoint-task":   "may add evidence to an Issue; it may record an explicit owner\n`accepted` closure",
	"savepoint-check":  "may close an Issue as `verified` after\nverifying Check proof",
}

func TestWorkingSkillsNameIssueCaptureEntry(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		for skill, rolePhrase := range issueCaptureEntrySkills {
			path := filepath.Join(root, skill, "SKILL.md")
			data, err := os.ReadFile(path)
			if err != nil {
				t.Errorf("%s: read %s: %v", tree, path, err)
				continue
			}
			content := string(data)

			if !strings.Contains(content, "## Issue Capture") {
				t.Errorf("%s: %s missing ## Issue Capture entry point", tree, path)
			}
			if !strings.Contains(content, "Enter Issue capture as an entry from this workflow") {
				t.Errorf("%s: %s does not name Issue capture as an entry from its own workflow", tree, path)
			}
			if !strings.Contains(content, "agent-skills/references/issue-capture.md") {
				t.Errorf("%s: %s does not point at the shared issue-capture reference", tree, path)
			}
			if !strings.Contains(content, rolePhrase) {
				t.Errorf("%s: %s does not state its own role phrase %q", tree, path, rolePhrase)
			}
		}
	}
}

// issueCaptureRuleProse is rule-defining language that must live only in the
// shared reference. If any of it leaks into a working skill file verbatim,
// the rule has been duplicated instead of referenced.
var issueCaptureRuleProse = []string{
	"never sufficient on its own to block a Task",
	"No automatic deduplication is assumed",
	"shortens, reorders, or edits a recorded entry is refused",
	"is proven by a Check that recorded",
	"does not resurrect a separate defect record, status, or public phase in V2",
}

func TestIssueCaptureRuleProseNotDuplicatedInWorkingSkills(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		for skill := range issueCaptureEntrySkills {
			path := filepath.Join(root, skill, "SKILL.md")
			data, err := os.ReadFile(path)
			if err != nil {
				t.Errorf("%s: read %s: %v", tree, path, err)
				continue
			}
			content := string(data)

			for _, prose := range issueCaptureRuleProse {
				if strings.Contains(content, prose) {
					t.Errorf("%s: %s duplicates Issue capture rule prose %q instead of referencing it", tree, path, prose)
				}
			}
		}
	}
}

func commandsAndProceduresPath(root string) string {
	return filepath.Join(root, "references", "commands-and-procedures.md")
}

func TestSharedCommandsAndProceduresIsNonTriggerableReference(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := commandsAndProceduresPath(root)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		if got := frontmatterField(content, "type"); got != "commands-reference" {
			t.Errorf("%s: %s type = %q, want commands-reference", tree, path, got)
		}
		if got := frontmatterField(content, "triggerable"); got != "false" {
			t.Errorf("%s: %s triggerable = %q, want false", tree, path, got)
		}
		if frontmatterField(content, "name") != "" {
			t.Errorf("%s: %s carries a skill name and would be discoverable as a skill", tree, path)
		}
	}
}

// commandsAndProceduresConfigKeys are the QualityGates field names
// internal/data/config.go decodes; the reference must name each one and
// must never introduce a parallel checks.technical surface.
var commandsAndProceduresConfigKeys = []string{"quality_gates", "lint", "typecheck", "build", "test", "block_on_failure", "gate_timeout"}

func TestSharedCommandsAndProceduresConfigContract(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := commandsAndProceduresPath(root)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		for _, key := range commandsAndProceduresConfigKeys {
			if !strings.Contains(content, key) {
				t.Errorf("%s: %s does not name config key %q", tree, path, key)
			}
		}
		if strings.Contains(content, "checks.technical") == false {
			t.Errorf("%s: %s does not name checks.technical as the rejected parallel surface", tree, path)
		}
		if !strings.Contains(content, "No parallel `checks.technical` surface") {
			t.Errorf("%s: %s does not state that no parallel checks.technical surface is introduced", tree, path)
		}
		if !strings.Contains(content, "explicit, optional gate command") {
			t.Errorf("%s: %s does not state build is an explicit, optional gate command", tree, path)
		}
		if !strings.Contains(content, "runs it in the project root") {
			t.Errorf("%s: %s does not describe where the build command runs", tree, path)
		}
	}
}

func TestSharedCommandsAndProceduresHealthCheckMapping(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := commandsAndProceduresPath(root)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		for _, phrase := range []string{
			"quality_gates",
			"check-method.md",
			"one preserved optional procedure file",
			"never installed as a mandatory",
			"no project is expected to have one",
			"needs no",
			"generated substitute",
			"not a finding",
		} {
			if !strings.Contains(content, phrase) {
				t.Errorf("%s: %s missing mapping phrase %q", tree, path, phrase)
			}
		}

		if !strings.Contains(content, "## Technical Verification") {
			t.Errorf("%s: %s does not name Task's own Technical Verification section", tree, path)
		}
		if !strings.Contains(content, "not something migration performs automatically") {
			t.Errorf("%s: %s does not state migration never auto-creates the procedure file or writes config.yml", tree, path)
		}
	}
}

func TestSharedCommandsAndProceduresLiveAndTemplateMatch(t *testing.T) {
	root := filepath.Join("..", "..")
	assertFileMatches(t, root,
		filepath.Join("agent-skills", "references", "commands-and-procedures.md"),
		filepath.Join("templates", "project-v2", "agent-skills", "references", "commands-and-procedures.md"),
	)
}

func TestSavepointDesignSkillReferencesCommandsAndProcedures(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		path := filepath.Join(root, "savepoint-design", "SKILL.md")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("%s: read %s: %v", tree, path, err)
			continue
		}
		content := string(data)

		if !strings.Contains(content, "agent-skills/references/commands-and-procedures.md") {
			t.Errorf("%s: %s does not point at the shared commands-and-procedures reference", tree, path)
		}
	}
}

func TestSavepointCreateDefectSkillPresentInBothTrees(t *testing.T) {
	for tree, root := range skillRoots() {
		path := filepath.Join(root, "savepoint-create-defect", "SKILL.md")
		if _, err := os.Stat(path); err != nil {
			t.Errorf("%s: savepoint-create-defect/SKILL.md missing: %v", tree, err)
		}
	}
}

// v2Skills and v2References are the complete V2 skill/reference set the
// routing and role-matrix tests below treat as one partitioned whole.
var v2Skills = []string{"savepoint-idea", "savepoint-design", "savepoint-task", "savepoint-check"}
var v2References = []string{"check-method.md", "issue-capture.md", "commands-and-procedures.md"}

// v1SkillNames are the nine pre-existing V1 skills that must survive the V2
// routing documentation untouched, in both trees.
var v1SkillNames = []string{
	"savepoint-draft-prd", "savepoint-system-design", "savepoint-create-task",
	"savepoint-build-task", "savepoint-audit-epic", "savepoint-audit-task",
	"savepoint-audit-register", "savepoint-create-defect", "savepoint-create-plan",
}

func TestV1SkillSetUnchangedAlongsideV2Routing(t *testing.T) {
	if len(v1SkillNames) != 9 {
		t.Fatalf("v1SkillNames has %d entries, want 9", len(v1SkillNames))
	}
	for tree, root := range skillRoots() {
		for _, name := range v1SkillNames {
			path := filepath.Join(root, name, "SKILL.md")
			if _, err := os.Stat(path); err != nil {
				t.Errorf("%s: V1 skill %s missing: %v", tree, name, err)
			}
		}
	}
}

func TestV2SkillSetIsCompleteWithByteParity(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		for _, name := range v2Skills {
			path := filepath.Join(root, name, "SKILL.md")
			if _, err := os.Stat(path); err != nil {
				t.Errorf("%s: V2 skill %s missing: %v", tree, name, err)
			}
		}
		for _, name := range v2References {
			path := filepath.Join(root, "references", name)
			if _, err := os.Stat(path); err != nil {
				t.Errorf("%s: V2 reference %s missing: %v", tree, name, err)
			}
		}
	}

	root := filepath.Join("..", "..")
	for _, name := range v2Skills {
		assertFileMatches(t, root,
			filepath.Join("agent-skills", name, "SKILL.md"),
			filepath.Join("templates", "project-v2", "agent-skills", name, "SKILL.md"))
	}
	for _, name := range v2References {
		assertFileMatches(t, root,
			filepath.Join("agent-skills", "references", name),
			filepath.Join("templates", "project-v2", "agent-skills", "references", name))
	}
}

// v2ObsoleteVocabulary is V1-only vocabulary that must never appear in the V2
// skill set: V2 renamed Epic/PRD to Objective/Idea, replaced the audit
// register with Check records, and folded defect handling into Issue
// capture, so any of these terms resurfacing signals V1 vocabulary leaking
// into V2 guidance.
var v2ObsoleteVocabulary = []string{
	"epic", "PRD", "audit register", "defect-building",
	"phase: build", "phase: test", "phase: audit", "phase: implementation",
}

func TestV2SkillSetHasNoObsoleteVocabulary(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		files := map[string]string{}
		for _, name := range v2Skills {
			files[name] = filepath.Join(root, name, "SKILL.md")
		}
		for _, name := range v2References {
			files["references/"+name] = filepath.Join(root, "references", name)
		}

		for label, path := range files {
			data, err := os.ReadFile(path)
			if err != nil {
				t.Errorf("%s: read %s: %v", tree, path, err)
				continue
			}
			content := string(data)
			for _, stale := range v2ObsoleteVocabulary {
				if strings.Contains(content, stale) {
					t.Errorf("%s: %s (%s) contains obsolete vocabulary %q", tree, label, path, stale)
				}
			}
		}
	}
}

// v2SkillAuthority pins the single owner of one sensitive write action across
// the four V2 skills, plus which of the remaining skills must explicitly
// disclaim or never discuss it. A second skill silently gaining permission
// for one of these actions — a checker that can edit acceptance criteria, an
// executor that can close an Issue — is the failure that makes the four-role
// split decorative, and it is only visible when all four are read as a set.
type v2SkillAuthority struct {
	name        string
	ownerSkill  string
	ownerPhrase string
	disclaims   map[string]string
	absentFrom  []string
}

var v2SkillAuthorities = []v2SkillAuthority{
	{
		name:        "Check record writing",
		ownerSkill:  "savepoint-check",
		ownerPhrase: "may write: the Check record",
		disclaims: map[string]string{
			"savepoint-task": "write a Check record",
		},
		absentFrom: []string{"savepoint-idea", "savepoint-design"},
	},
	{
		name:        "Issue closure",
		ownerSkill:  "savepoint-check",
		ownerPhrase: "may close an Issue",
		disclaims: map[string]string{
			"savepoint-task":   "close an Issue",
			"savepoint-design": "does not close one",
		},
		absentFrom: []string{"savepoint-idea"},
	},
	{
		name:        "Acceptance-criteria edits",
		ownerSkill:  "savepoint-design",
		ownerPhrase: "## Done When",
		disclaims: map[string]string{
			"savepoint-task":  "edit the Task's acceptance criteria",
			"savepoint-check": "edit acceptance criteria to match a result",
		},
		absentFrom: []string{"savepoint-idea"},
	},
	{
		name:       "Owner acceptance claims",
		ownerSkill: "",
		disclaims: map[string]string{
			"savepoint-task": "claim clearance or owner acceptance",
		},
	},
}

func TestV2SkillRoleMatrixPartitionsSensitiveWrites(t *testing.T) {
	for tree, root := range v2SkillRoots() {
		content := map[string]string{}
		for _, name := range v2Skills {
			data, err := os.ReadFile(filepath.Join(root, name, "SKILL.md"))
			if err != nil {
				t.Fatalf("%s: read %s: %v", tree, name, err)
			}
			content[name] = string(data)
		}

		for _, authority := range v2SkillAuthorities {
			if authority.ownerSkill != "" && !strings.Contains(content[authority.ownerSkill], authority.ownerPhrase) {
				t.Errorf("%s: %s does not carry owner phrase %q for %s", tree, authority.ownerSkill, authority.ownerPhrase, authority.name)
			}
			for skill, phrase := range authority.disclaims {
				if !strings.Contains(content[skill], phrase) {
					t.Errorf("%s: %s does not disclaim %q for %s", tree, skill, phrase, authority.name)
				}
			}
			for _, skill := range authority.absentFrom {
				if authority.ownerPhrase != "" && strings.Contains(content[skill], authority.ownerPhrase) {
					t.Errorf("%s: %s unexpectedly carries owner phrase %q for %s", tree, skill, authority.ownerPhrase, authority.name)
				}
			}
		}
	}
}

// v2RoutingRequiredPhrases are the load-bearing statements the V2 routing
// section must carry in both AGENTS.md guides: the four-state mapping, the
// active-for-V2-but-not-this-repository statement, the REPLAN REQUIRED
// routing rule, and the shared-reference ownership list.
var v2RoutingRequiredPhrases = []string{
	"| idea | savepoint-idea |",
	"| design | savepoint-design |",
	"| task | savepoint-task |",
	"| check | savepoint-check |",
	"is not active",
	"E47",
	"REPLAN REQUIRED",
	"not a fifth router state",
	"agent-skills/references/check-method.md",
	"agent-skills/references/issue-capture.md",
	"agent-skills/references/commands-and-procedures.md",
	"non-triggerable",
}

func TestAgentsGuidesCarryV2RoutingSection(t *testing.T) {
	root := filepath.Join("..", "..")
	liveAgents := readTemplate(t, root, "AGENTS.md")
	templateAgents := readTemplate(t, root, "templates", "project", "AGENTS.md")

	for _, phrase := range v2RoutingRequiredPhrases {
		assertContains(t, liveAgents, phrase)
		assertContains(t, templateAgents, phrase)
	}

	liveSection, found := sectionBody(liveAgents, "## V2 Routing")
	if !found || liveSection == "" {
		t.Fatal("AGENTS.md missing or empty V2 routing section")
	}
	templateSection, found := sectionBody(templateAgents, "## V2 Routing")
	if !found || templateSection == "" {
		t.Fatal("templates/project/AGENTS.md missing or empty V2 routing section")
	}
	if liveSection != templateSection {
		t.Fatal("V2 routing section differs between live and template AGENTS.md")
	}

	// The V1 activation table must still govern current routing, unchanged.
	assertContains(t, liveAgents, "| task-building | savepoint-build-task |")
	assertContains(t, templateAgents, "| task-building | savepoint-build-task |")
}

// The design skill's outgoing handoff must use the same V2 state/skill pair
// as the routing table. Checking both contracts together prevents a design
// handoff from silently falling back to the still-active V1 state.
func TestV2DesignHandoffMatchesRoutingContract(t *testing.T) {
	root := filepath.Join("..", "..")
	const routingHeading = "## V2 Routing"
	const routeToTask = "| task | savepoint-task |"
	const designHandoff = "set router `state: task` for the first unblocked planned Task and update `next_action` to execute it with `savepoint-task`."

	pairs := []struct {
		name       string
		agentsPath []string
		designPath []string
	}{
		{
			name:       "live",
			agentsPath: []string{"AGENTS.md"},
			designPath: []string{"agent-skills", "savepoint-design", "SKILL.md"},
		},
		{
			name:       "template",
			agentsPath: []string{"templates", "project", "AGENTS.md"},
			designPath: []string{"templates", "project-v2", "agent-skills", "savepoint-design", "SKILL.md"},
		},
	}

	for _, pair := range pairs {
		agents := readTemplate(t, root, pair.agentsPath...)
		routingSection, found := sectionBody(agents, routingHeading)
		if !found {
			t.Fatalf("%s AGENTS.md missing %s", pair.name, routingHeading)
		}
		if !strings.Contains(routingSection, routeToTask) {
			t.Errorf("%s routing section missing V2 task route %q", pair.name, routeToTask)
		}

		design := readTemplate(t, root, pair.designPath...)
		workflow, found := sectionBody(design, "## Workflow")
		if !found {
			t.Fatalf("%s design skill missing ## Workflow", pair.name)
		}
		if !strings.Contains(workflow, designHandoff) {
			t.Errorf("%s design workflow missing V2 handoff %q", pair.name, designHandoff)
		}
		if strings.Contains(workflow, "task-building") {
			t.Errorf("%s design workflow still names the V1 task-building state", pair.name)
		}
	}
}
