package init

import (
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// assertPhrase matches prose across line wraps: the skills and shared method are
// hard-wrapped markdown, so a contract phrase may be split by a newline and
// indentation without any change in meaning.
func assertPhrase(t *testing.T, content, want string) {
	t.Helper()

	if !strings.Contains(collapseSpace(content), collapseSpace(want)) {
		t.Errorf("missing contract phrase %q", want)
	}
}

func assertNoPhrase(t *testing.T, content, stale string) {
	t.Helper()

	if strings.Contains(collapseSpace(content), collapseSpace(stale)) {
		t.Errorf("contains retired contract phrase %q", stale)
	}
}

var spaceRun = regexp.MustCompile(`\s+`)

func collapseSpace(s string) string {
	return strings.TrimSpace(spaceRun.ReplaceAllString(s, " "))
}

// materialityHeader is the one table both audit skills must reproduce verbatim,
// so a rename of any column is caught in both contracts.
const materialityHeader = "| Finding | Likelihood | Impact | Materiality | Recommendation |"

// auditSkillPair returns the live and scaffolded copies of one skill. Contracts
// are asserted against both so a generated project never drifts from the
// behavior this repository tests.
func auditSkillPair(t *testing.T, name string) []string {
	t.Helper()

	root := filepath.Join("..", "..")
	return []string{
		readTemplate(t, root, "agent-skills", name, "SKILL.md"),
		readTemplate(t, root, "templates", "project", "agent-skills", name, "SKILL.md"),
	}
}

func sharedAuditMethods(t *testing.T) []string {
	t.Helper()

	root := filepath.Join("..", "..")
	return []string{
		readTemplate(t, root, "agent-skills", "references", "audit-method.md"),
		readTemplate(t, root, "templates", "project", "agent-skills", "references", "audit-method.md"),
	}
}

// Routing must select the epic audit from the audit-pending state and the task
// audit from an explicit request, without inventing a router state for the
// latter.
func TestRoutingSelectsSplitAuditsWithoutANewState(t *testing.T) {
	root := filepath.Join("..", "..")
	router := readTemplate(t, root, "templates", "project", ".savepoint", "router.md")

	assertPhrase(t, router, "| audit-pending | savepoint-audit-epic |")
	assertPhrase(t, router, "`savepoint-audit-task` has no row")
	assertPhrase(t, router, "use `savepoint-audit-task` and leave `state: task-building` unchanged")
	assertPhrase(t, router, "not a new state")
	assertNoPhrase(t, router, "| task-audit |")
	assertNoPhrase(t, router, "state: task-audit")

	for _, agents := range []string{
		readTemplate(t, root, "AGENTS.md"),
		readTemplate(t, root, "templates", "project", "AGENTS.md"),
	} {
		assertPhrase(t, agents, "| audit-pending | savepoint-audit-epic |")
		assertPhrase(t, agents, "uses `savepoint-audit-task` while `state` stays `task-building`")
		assertPhrase(t, agents, "not a router state")
	}

	// The builder must hand explicit audit requests to the split skills rather
	// than reviewing its own work.
	for _, build := range auditSkillPair(t, "savepoint-build-task") {
		assertPhrase(t, build, "savepoint-audit-task")
		assertPhrase(t, build, "savepoint-audit-epic")
		assertPhrase(t, build, "Do not audit the epic you just built.")
	}
}

func TestTaskAuditSkillContract(t *testing.T) {
	for _, skill := range auditSkillPair(t, "savepoint-audit-task") {
		// No epic artifact, no lifecycle writes.
		assertPhrase(t, skill, "Do not create an epic audit artifact.")
		assertPhrase(t, skill, "`E##-Audit.md` belongs to")
		assertPhrase(t, skill, "Do not write any file.")
		assertPhrase(t, skill, "Do not change task `status`, task `stage`, or router `state`.")

		// Quick health check only, and one of exactly two result values.
		assertPhrase(t, skill, "Apply the Quick health check.")
		assertPhrase(t, skill, "Do not use a health-check mode other than Quick.")
		assertPhrase(t, skill, "Return exactly one result value: `CLEAR` or `NEEDS WORK`.")
		assertPhrase(t, skill, "Do not return a result value other than `CLEAR` or `NEEDS WORK`.")

		// Enriched rigor: file reality, materiality, sizing, handoff format.
		assertPhrase(t, skill, "Verify file reality")
		assertPhrase(t, skill, "phantom file as a finding")
		assertPhrase(t, skill, materialityHeader)
		assertPhrase(t, skill, "350–600 words")
		assertPhrase(t, skill, "## Final Response Output")
		assertPhrase(t, skill, "State the gate result.")
		assertPhrase(t, skill, "Link to the task file under audit.")

		// Shared method is loaded, not restated.
		assertPhrase(t, skill, "agent-skills/references/audit-method.md")
	}
}

func TestEpicAuditSkillContract(t *testing.T) {
	for _, skill := range auditSkillPair(t, "savepoint-audit-epic") {
		// Full health check, whole epic, independent session.
		assertPhrase(t, skill, "Apply the Full health check")
		assertPhrase(t, skill, "Do not use a health-check mode other than Full.")
		assertPhrase(t, skill, "every completed task in the epic")
		assertPhrase(t, skill, "requires a session independent from the builder")
		assertPhrase(t, skill, "Do not audit an epic you built in this session.")

		// Exactly one artifact, approval before apply, explicit closeout authority.
		assertPhrase(t, skill, ".savepoint/releases/{release}/epics/{E##-slug}/E##-Audit.md")
		assertPhrase(t, skill, "Do not create more than one audit file for an epic.")
		assertPhrase(t, skill, "Do not apply proposals before approval.")
		assertPhrase(t, skill, "## Apply And Close")
		assertPhrase(t, skill, "Mark the epic audited and advance the router.")

		// Enriched rigor: materiality, guardrails, repository handoff, sizing.
		assertPhrase(t, skill, materialityHeader)
		assertPhrase(t, skill, "### Guardrails Verification")
		assertPhrase(t, skill, "CLEAR TO COMMIT/PUSH")
		assertPhrase(t, skill, "NOT READY TO COMMIT/PUSH")
		assertPhrase(t, skill, "500–900 words")

		// File reality evidence and the user-facing handoff format.
		assertPhrase(t, skill, "Verify file reality")
		assertPhrase(t, skill, "File reality evidence:")
		assertPhrase(t, skill, "## Final Response Output")
		assertPhrase(t, skill, "State the gate result.")
		assertPhrase(t, skill, "Link to the audit file.")

		assertPhrase(t, skill, "agent-skills/references/audit-method.md")
	}
}

func TestSharedAuditMethodRequiresScenarioClassesAndMatrices(t *testing.T) {
	for _, method := range sharedAuditMethods(t) {
		// Frozen scope lock with its named fields.
		assertPhrase(t, method, "## Freeze The Audit Scope")
		assertPhrase(t, method, "the acceptance criteria, guardrails, and release gates being tested;")
		assertPhrase(t, method, "the changed files and supported public entry points;")
		assertPhrase(t, method, "the selected matrix axes and cells, including explicit not-applicable cells;")
		assertPhrase(t, method, "the supported-path and materiality boundary used to admit findings.")
		assertPhrase(t, method, "the lock is immutable for every re-audit")

		// Mandatory coverage matrix with its named axes.
		assertPhrase(t, method, "## Build The Mandatory Coverage Matrix")
		assertPhrase(t, method, "A prose checklist is not a matrix.")
		for _, axis := range []string{
			"**Public surfaces:**",
			"**Input shape:**",
			"**State:**",
			"**Environment/output:**",
			"**Boundaries:**",
			"**Sequences:**",
			"**Representations:**",
			"**Text classes:**",
		} {
			assertPhrase(t, method, axis)
		}
		assertPhrase(t, method, "### Finite External-Boundary Matrix")
		assertPhrase(t, method, "### Matrix Completion Lock")

		// Workflow and side-effect lock.
		assertPhrase(t, method, "## Workflow And Side-Effect Audit Lock")
		assertPhrase(t, method, "| Order | Real operation | Side effect / state change | Failure timing | Failure owner and final state | Cleanup / secondary failure | Independent oracle |")
		assertPhrase(t, method, "Use an independent oracle.")

		// Full re-audit, admission ledger, convergence limit, credible blocker.
		assertPhrase(t, method, "## Re-audit After Remediation")
		assertPhrase(t, method, "every original matrix cell;")
		assertPhrase(t, method, "every original reproduction;")
		assertPhrase(t, method, "| Re-audit check | Prior finding or remediation claim | Exact frozen matrix cell | Allowed result |")
		assertPhrase(t, method, "Default convergence limit:")
		assertPhrase(t, method, "credible blocker")

		// Evidence classification and consolidated finding fields.
		assertPhrase(t, method, "passed, finding, unverified, or not-applicable")
		assertPhrase(t, method, "**Proven:**")
		assertPhrase(t, method, "**Finding:**")
		assertPhrase(t, method, "**Unverified:**")
		assertPhrase(t, method, "Each finding must include:")
		for _, field := range []string{
			"the violated acceptance criterion or guardrail rule;",
			"the smallest reproducible scenario;",
			"expected and actual behavior;",
			"exact file and line evidence; and",
			"the missing or inadequate test evidence.",
		} {
			assertPhrase(t, method, field)
		}

		// File reality and materiality live in the shared method, not per-skill.
		assertPhrase(t, method, "## Verify File Reality")
		assertPhrase(t, method, "## Summarize Materiality")
		assertPhrase(t, method, materialityHeader)

		// The method never triggers on its own.
		assertPhrase(t, method, "triggerable: false")
		assertPhrase(t, method, "This reference is not a skill and never triggers on its own.")
	}
}
