package doctor

import (
	"errors"
	"fmt"
	"strings"

	"github.com/opencode/savepoint/internal/data"
)

// typedRepairRules is checked first, in order, with errors.Is so a wrapped
// sentinel wins over any message wording.
var typedRepairRules = []struct {
	target error
	repair string
}{
	{data.ErrConfigNotFound, "Run `savepoint init` to scaffold a new project"},
	{data.ErrInvalidStatus, "Set router state to a recognized workflow state (see router.md State → action section)"},
	{data.ErrMissingFrontmatter, "Fix the YAML frontmatter between the --- delimiters"},
	{data.ErrStructureProblem, "Review the file and fix the reported issue"},
}

// messageRepairRule matches a message when every term of any one alternative is
// a substring of it.
type messageRepairRule struct {
	alternatives [][]string
	repair       string
}

func (r messageRepairRule) matches(msg string) bool {
	for _, terms := range r.alternatives {
		if containsAll(msg, terms) {
			return true
		}
	}
	return false
}

func containsAll(msg string, terms []string) bool {
	for _, term := range terms {
		if !strings.Contains(msg, term) {
			return false
		}
	}
	return true
}

// messageRepairRules is ordered: the first matching rule wins, so a broad
// term such as "release" deliberately shadows every later, narrower rule.
// Legacy V1 wording stays here for compatibility with older projects.
var messageRepairRules = []messageRepairRule{
	{[][]string{{"config.yml not found"}}, "Run `savepoint init` to scaffold a new project"},
	{[][]string{{"config.yml missing required field"}}, "Add the missing field to config.yml — see the project template for reference"},
	{[][]string{{"invalid YAML"}}, "Fix the YAML syntax error at the indicated line"},
	{[][]string{{"router.md not found"}}, "Run `savepoint init` to scaffold a new project"},
	{[][]string{{"unknown state"}}, "Set router state to a recognized workflow state (see router.md State → action section)"},
	{[][]string{{"release PRD file not found"}}, "Create a {release}-PRD.md file with frontmatter for the release"},
	{[][]string{{"defect uses non-canonical status"}}, "Replace the defect status with the canonical value named in the problem"},
	{[][]string{{"defect status invalid"}}, "Set the defect status to open, in_progress, or resolved"},
	{[][]string{{"defect stage is required"}}, "Add stage: build, stage: test, or stage: audit while the defect status is in_progress"},
	{[][]string{{"defect stage invalid"}}, "Set the defect stage to build, test, or audit"},
	{[][]string{{"defect stage", "only valid"}}, "Remove stage unless the defect status is in_progress"},
	{[][]string{{"epic uses non-canonical status"}, {"epic status invalid"}}, "Set the epic status to planned, in_progress, done, or audited"},
	{[][]string{{"release"}}, "Create the release directory at releases/<release-id>/"},
	{[][]string{{"epic", "directory not found"}}, "Create the epic directory at releases/<release>/epics/<epic-id>/"},
	{[][]string{{"epic detail file not found"}}, "Create an E##-Detail.md with frontmatter for the epic"},
	{[][]string{{"invalid frontmatter"}}, "Fix the YAML frontmatter between the --- delimiters"},
	{[][]string{{"task missing required frontmatter field"}}, "Add the missing field to the task frontmatter"},
	{[][]string{{"task uses non-canonical status"}}, "Replace the task status with the canonical value named in the problem"},
	{[][]string{{"task uses legacy frontmatter field phase"}}, "Use stage: build, stage: test, or stage: audit only while status is in_progress; otherwise remove phase"},
	{[][]string{{"task stage is required"}}, "Add stage: build, stage: test, or stage: audit while status is in_progress"},
	{[][]string{{"task stage invalid"}}, "Set stage to build, test, or audit"},
	{[][]string{{"task stage field", "only valid"}}, "Remove stage unless the task status is in_progress"},
	{[][]string{{"task complexity invalid", "complexity_reason"}}, "Shorten complexity_reason to the configured word limit and keep complexity_tier set"},
	{[][]string{{"missing ## Acceptance Criteria"}}, "Add an ## Acceptance Criteria section with checkable items"},
	{[][]string{{"depends_on must be a list"}}, "Change depends_on to a YAML list format"},
	{[][]string{{"references non-existent"}}, "Create the referenced task or remove the dependency"},
	{[][]string{{"duplicate task ID"}}, "Rename one of the tasks to have a unique ID"},
	{[][]string{{"dependency cycle"}}, "Break the circular dependency chain between tasks"},
	{[][]string{{"audit proposal exists"}}, "Set router state to audit-pending for the matching epic, or remove stale audit files"},
	{[][]string{{"orphaned"}}, "Move the task directory to the correct epic or create the referenced epic"},
	{[][]string{{"defect parse error"}}, "Fix the defect frontmatter — required fields: id, severity, title (or objective); status must be open, in_progress, or resolved"},
	{[][]string{{"defect missing required frontmatter field"}}, "Add the missing field to the defect frontmatter"},
	{[][]string{{"defect reference", "empty"}}, "Fix the reference field — format must be 'EPIC-slug/TASK-slug' with non-empty parts"},
	{[][]string{{"defect reference", "does not match"}}, "Update the reference field to match an existing task ID, or remove it if the task was deleted"},
	{[][]string{{"quality gate"}}, "Fix the issue reported by the quality gate tool"},
}

const defaultRepair = "Review the file and fix the reported issue"

func SuggestRepair(err error) string {
	for _, rule := range typedRepairRules {
		if errors.Is(err, rule.target) {
			return rule.repair
		}
	}

	msg := err.Error()
	for _, rule := range messageRepairRules {
		if rule.matches(msg) {
			return rule.repair
		}
	}
	return defaultRepair
}

var v2ProblemRepairs = map[string]string{
	"schema-version-malformed":              "Set config.yml's schema_version to the integer 2, or remove the field for a V1 project",
	"schema-version-unsupported":            "Set config.yml's schema_version to 2, the only supported explicit value, or remove the field for a V1 project",
	"health-snapshot-missing":               "Set health_snapshot to the ID of an official snapshot under .savepoint/health/snapshots, or remove the field from the Check",
	"health-snapshot-manual":                "Point health_snapshot at an official snapshot recorded by a Full Objective Check; a manual refresh cannot back a Check",
	"health-snapshot-unreadable":            "Fix or remove the unreadable file named in the message under .savepoint/health/snapshots, then run doctor again",
	"v2-missing-field":                      "Add the missing required field named in the diagnostic to the record's frontmatter",
	"v2-invalid-id":                         "Set the record's id or reference to a valid family identity: G-### (Goal), O-### (Objective), T-### (Task), C-### (Check), or I-### (Issue reference), each with at least three digits",
	"v2-release-invalid-id":                 "Set the Goal's id to a valid R-### or G-### identity with at least three digits",
	"v2-invalid-ownership":                  "Set the Task's objective field to exactly one existing O-### Objective id",
	"v2-invalid-lifecycle":                  "Set status (and stage while status is in_progress) to a supported V2 lifecycle value",
	"v2-invalid-dependency":                 "Fix the depends_on entry: task must be a T-### id and requires must be clear or accepted",
	"v2-invalid-release-reference":          "Set the Objective's release field to an existing R-### or G-### Goal ID, or remove it for an unassigned Objective",
	"v2-duplicate-id":                       "Rename one of the two records reporting the same id so each global id is declared once",
	"v2-path-mismatch":                      "Rename the directory or file so its name starts with the id the record declares",
	"v2-unsafe-path":                        "Remove the symlink or case-aliasing path reported in the diagnostic; V2 records must resolve inside the project root",
	"v2-missing-owner":                      "Create the referenced O-### Objective, or fix the Task's objective field to reference one that exists",
	"v2-missing-release":                    "Create the referenced R-### or G-### Goal, or remove or correct the Objective's release field",
	"v2-release-missing-section":            "Add the missing required Release body section: Outcome, Why, Success Conditions, or Boundaries",
	"v2-release-legacy-malformed":           "Fix legacy_completion.source_path, archive_path, and sha256 on the named historical Release",
	"v2-missing-dependency-target":          "Create the referenced dependency record, or remove it from depends_on",
	"v2-self-dependency":                    "Remove the record's own id from its depends_on list",
	"v2-dependency-cycle":                   "Break the circular dependency chain named in the diagnostic",
	"v2-record-malformed":                   "Fix the YAML frontmatter between the --- delimiters in the named record file",
	"v2-check-malformed":                    "Fix the named Check field in the record's frontmatter — result must be CLEAR or NEEDS WORK, checked_by.role a supported actor role, and checked_at a parseable RFC 3339 timestamp",
	"v2-check-missing-scope-target":         "Create the Task or Objective the Check's scope names, or fix scope.id to reference one that exists",
	"v2-check-missing-release-scope-target": "Create the Release the Check's scope names, or fix scope.id to reference one that exists",
	"v2-check-missing-reference":            "Create the Check named in supersedes, or remove the supersedes field",
	"v2-check-supersedes-conflict":          "Fix the supersedes chain: it must name a Check with the same scope, no two Checks may supersede the same target, and the chain must not cycle",
	"v2-evidence-malformed":                 "Fix the named evidence field in the record's frontmatter — state and role values must be one of the supported options and timestamps must be parseable RFC 3339",
	"v2-evidence-missing-reference":         "Fix the evidence field to name a Check that exists, or remove the reference",
	"v2-check-immutable":                    "A Check record is immutable once written — record a new Check with supersedes naming this one instead of editing it",
	"v2-issue-malformed":                    "Fix the named Issue field in the record's frontmatter — type, status, source, resolution, and history each have a fixed vocabulary",
	"v2-issue-missing-duplicate-target":     "Set duplicate_of to an existing I-### Issue id, or remove the field",
	"v2-issue-missing-escalation-target":    "Set escalated_to to an existing O-### Objective id, or remove the field",
	"v2-issue-self-duplicate":               "Remove the Issue's own id from its duplicate_of field, or point it at a different canonical Issue",
	"v2-issue-duplicate-cycle":              "Break the circular duplicate_of chain named in the diagnostic so one Issue in the ring is canonical",
	"v2-issue-missing-link-target":          "Fix the named tasks, checks, or issues reference to name a record that exists, or remove it",
	"v2-issue-unpaired-check-link":          "Add the Check's id to the named Issue's checks field, or remove the Issue from the Check's issues field",
	"v2-issue-resolution-required":          "Add a resolution block naming a disposition, or set the Issue's status back to open or in_progress",
	"v2-issue-resolution-not-allowed":       "Remove the resolution block, or set the Issue's status to resolved",
	"v2-issue-resolution-missing-proof":     "Set resolution.check to an existing, CLEAR C-### Check that appears in the Issue's checks field",
	"v2-issue-resolution-unusable-proof":    "Point resolution.check at a Check that recorded CLEAR and is listed in the Issue's checks field",
	"v2-issue-resolution-field-mismatch":    "Fix the resolution field for its disposition: accepted needs an owner actor and reason and no proof check; duplicate needs duplicate_of and no proof check; escalated needs a planner actor, escalated_to, and no proof check",
	"v2-issue-already-exists":               "Choose a different Issue id or file path — an existing Issue record cannot be overwritten by create",
	"v2-issue-history-not-append-only":      "Append new history entries after the recorded ones rather than editing, reordering, or removing any existing entry",
}

// V2ProblemRepair maps a RunV2Checks diagnostic name (v2DiagnosticName in
// checks.go) to a manual repair suggestion. Doctor never repairs V2 project
// files itself; the reported Problem's File already names the project root,
// and the diagnostic message carries the specific record path and identity.
func V2ProblemRepair(name string) string {
	if repair, ok := v2ProblemRepairs[name]; ok {
		return repair
	}
	return "Review the V2 project diagnostic and fix the reported record"
}

// V2ConsistencyRepair maps a doctor consistency diagnostic name
// (v2ConsistencyDiagnosticName in checks.go, derived from
// data.ConsistencyDiagnosticKind) to a manual repair suggestion. Doctor never
// rewrites a Check, an evidence field, or a record status itself; every
// suggestion here names the record the diagnostic already carries the path
// to.
func V2ConsistencyRepair(name string) string {
	switch name {
	case "v2-done-without-clearance":
		return "Record a CLEAR Check before treating the Task as done, or set status back to reflect its real progress"
	case "v2-acceptance-superseded":
		return "Have the owner accept the Task's actual latest Check; an acceptance naming a superseded Check no longer applies"
	case "v2-evidence-contradicts-status":
		return "Advance the Task's status to done to match its recorded evidence, or correct the evidence if completion was not actually reached"
	case "v2-objective-done-without-clearance":
		return "Record a CLEAR Objective-scoped Check before treating the Objective as done, or set status back to reflect its real progress"
	case "v2-objective-planned-with-started-task":
		return "Set the Objective's status to in_progress; its work has started"
	case "v2-objective-done-with-incomplete-task":
		return "Finish the named owned Task before closing the Objective, or set the Objective's status back to reflect its real progress"
	case "v2-issue-verified-proof-superseded":
		return "Record a fresh proof Check against the Issue's verified resolution, or update resolution.check to name the current latest Check for that scope"
	default:
		return "Review the Task's recorded evidence against its status and fix the mismatch"
	}
}

// GateSuggestion returns a command-specific repair hint.
func GateSuggestion(name string) string {
	switch name {
	case "lint":
		return "Run `make lint` locally and fix reported issues"
	case "typecheck":
		return "Run `make typecheck` locally and fix type errors"
	case "test":
		return "Run `make test` locally and fix failing tests"
	default:
		return fmt.Sprintf("Run %q locally and fix reported issues", name)
	}
}
