package doctor

import (
	"errors"
	"fmt"
	"strings"

	"github.com/opencode/savepoint/internal/data"
)

func SuggestRepair(err error) string {
	switch {
	case errors.Is(err, data.ErrConfigNotFound):
		return "Run `savepoint init` to scaffold a new project"
	case errors.Is(err, data.ErrInvalidStatus):
		return "Set router state to a recognized workflow state (see router.md State → action section)"
	case errors.Is(err, data.ErrMissingFrontmatter):
		return "Fix the YAML frontmatter between the --- delimiters"
	case errors.Is(err, data.ErrStructureProblem):
		return "Review the file and fix the reported issue"
	}

	msg := err.Error()
	switch {
	case strings.Contains(msg, "config.yml not found"):
		return "Run `savepoint init` to scaffold a new project"
	case strings.Contains(msg, "config.yml missing required field"):
		return "Add the missing field to config.yml — see the project template for reference"
	case strings.Contains(msg, "invalid YAML"):
		return "Fix the YAML syntax error at the indicated line"
	case strings.Contains(msg, "router.md not found"):
		return "Run `savepoint init` to scaffold a new project"
	case strings.Contains(msg, "unknown state"):
		return "Set router state to a recognized workflow state (see router.md State → action section)"
	case strings.Contains(msg, "release PRD file not found"):
		return "Create a {release}-PRD.md file with frontmatter for the release"
	case strings.Contains(msg, "defect uses non-canonical status"):
		return "Replace the defect status with the canonical value named in the problem"
	case strings.Contains(msg, "defect status invalid"):
		return "Set the defect status to open, in_progress, or resolved"
	case strings.Contains(msg, "defect stage is required"):
		return "Add stage: build, stage: test, or stage: audit while the defect status is in_progress"
	case strings.Contains(msg, "defect stage invalid"):
		return "Set the defect stage to build, test, or audit"
	case strings.Contains(msg, "defect stage") && strings.Contains(msg, "only valid"):
		return "Remove stage unless the defect status is in_progress"
	case strings.Contains(msg, "epic uses non-canonical status"), strings.Contains(msg, "epic status invalid"):
		return "Set the epic status to planned, in_progress, done, or audited"
	case strings.Contains(msg, "release"):
		return "Create the release directory at releases/<release-id>/"
	case strings.Contains(msg, "epic") && strings.Contains(msg, "directory not found"):
		return "Create the epic directory at releases/<release>/epics/<epic-id>/"
	case strings.Contains(msg, "epic detail file not found"):
		return "Create an E##-Detail.md with frontmatter for the epic"
	case strings.Contains(msg, "invalid frontmatter"):
		return "Fix the YAML frontmatter between the --- delimiters"
	case strings.Contains(msg, "task missing required frontmatter field"):
		return "Add the missing field to the task frontmatter"
	case strings.Contains(msg, "task uses non-canonical status"):
		return "Replace the task status with the canonical value named in the problem"
	case strings.Contains(msg, "task uses legacy frontmatter field phase"):
		return "Use stage: build, stage: test, or stage: audit only while status is in_progress; otherwise remove phase"
	case strings.Contains(msg, "task stage is required"):
		return "Add stage: build, stage: test, or stage: audit while status is in_progress"
	case strings.Contains(msg, "task stage invalid"):
		return "Set stage to build, test, or audit"
	case strings.Contains(msg, "task stage field") && strings.Contains(msg, "only valid"):
		return "Remove stage unless the task status is in_progress"
	case strings.Contains(msg, "task complexity invalid") && strings.Contains(msg, "complexity_reason"):
		return "Shorten complexity_reason to the configured word limit and keep complexity_tier set"
	case strings.Contains(msg, "missing ## Acceptance Criteria"):
		return "Add an ## Acceptance Criteria section with checkable items"
	case strings.Contains(msg, "depends_on must be a list"):
		return "Change depends_on to a YAML list format"
	case strings.Contains(msg, "references non-existent"):
		return "Create the referenced task or remove the dependency"
	case strings.Contains(msg, "duplicate task ID"):
		return "Rename one of the tasks to have a unique ID"
	case strings.Contains(msg, "dependency cycle"):
		return "Break the circular dependency chain between tasks"
	case strings.Contains(msg, "audit proposal exists"):
		return "Set router state to audit-pending for the matching epic, or remove stale audit files"
	case strings.Contains(msg, "orphaned"):
		return "Move the task directory to the correct epic or create the referenced epic"
	case strings.Contains(msg, "defect parse error"):
		return "Fix the defect frontmatter — required fields: id, severity, title (or objective); status must be open, in_progress, or resolved"
	case strings.Contains(msg, "defect missing required frontmatter field"):
		return "Add the missing field to the defect frontmatter"
	case strings.Contains(msg, "defect reference") && strings.Contains(msg, "empty"):
		return "Fix the reference field — format must be 'EPIC-slug/TASK-slug' with non-empty parts"
	case strings.Contains(msg, "defect reference") && strings.Contains(msg, "does not match"):
		return "Update the reference field to match an existing task ID, or remove it if the task was deleted"
	case strings.Contains(msg, "quality gate"):
		return "Fix the issue reported by the quality gate tool"
	default:
		return "Review the file and fix the reported issue"
	}
}

// auditFrontmatterRepair is the suggestion for an audit record whose frontmatter
// fails to parse structurally.
const auditFrontmatterRepair = "Fix the YAML frontmatter between the --- delimiters in the named audit file"

// AuditFindingRepair maps a load-time finding diagnostic to a manual repair
// suggestion naming the frontmatter field to edit in the reported finding file.
// Doctor never repairs audit-register files itself.
func AuditFindingRepair(code data.FindingDiagnosticCode) string {
	switch code {
	case data.FindingMissingFieldCode:
		return "Add the named field to the finding file's frontmatter"
	case data.FindingInvalidIDCode:
		return "Set the id field to the finding's stable F### id (at least three digits)"
	case data.FindingIDMismatchCode:
		return "Set the id field to match the filename's F### id, or rename the file to match the id"
	case data.FindingInvalidStatusCode:
		return "Set the status field to one of the canonical finding statuses named in the problem"
	case data.FindingInvalidSeverityCode:
		return "Set the severity field to critical, high, medium, or low"
	case data.FindingInvalidConfidenceCode:
		return "Set the confidence field to high, medium, or low"
	default:
		return "Review the finding file and fix the reported field"
	}
}

// AuditValidationRepair maps a cross-record audit validation code to a manual
// repair suggestion naming the frontmatter field to edit in the reported finding
// file. Broken links distinguish release, epic, task, defect, and duplicate
// finding references so the author knows which list to fix.
func AuditValidationRepair(code data.AuditValidationCode) string {
	switch code {
	case data.AuditVerifiedMissingProof:
		return "A verified finding requires named proof — set the verified_proof field to the passing regression test"
	case data.AuditDuplicateMissingTarget:
		return "Set the duplicate_of field to the canonical F### finding this record duplicates"
	case data.AuditDeferredMissingRationale:
		return "Set the deferral_reason field to explain why the finding is deferred"
	case data.AuditOwnerDecisionMissingRationale:
		return "Set the deferral_reason or waiver_reason field to record the owner's decision"
	case data.AuditWaivedMissingRationale:
		return "Set the waiver_reason field to explain why the finding is waived"
	case data.AuditReleaseRefMissing:
		return "Fix the releases entry to name an existing release id, or remove the entry"
	case data.AuditEpicRefMissing:
		return "Fix the epics entry to name an existing epic (E## or full epic id), or remove the entry"
	case data.AuditTaskRefMissing:
		return "Fix the tasks entry to name an existing task (T-### or epic/task id), or remove the entry"
	case data.AuditDefectRefMissing:
		return "Fix the defects entry to name an existing defect id, or remove the entry"
	case data.AuditDuplicateRefMissing:
		return "Point the duplicate_of field at an existing finding other than this one"
	default:
		return "Review the finding file and fix the reported field"
	}
}

// V2ProblemRepair maps a CheckProject diagnostic name (v2DiagnosticName in
// checks.go) to a manual repair suggestion. Doctor never repairs V2 project
// files itself; the reported Problem's File already names the project root,
// and the diagnostic message carries the specific record path and identity.
func V2ProblemRepair(name string) string {
	switch name {
	case "schema-version-malformed":
		return "Set config.yml's schema_version to the integer 2, or remove the field for a V1 project"
	case "schema-version-unsupported":
		return "Set config.yml's schema_version to 2, the only supported explicit value, or remove the field for a V1 project"
	case "v2-missing-field":
		return "Add the missing required field named in the diagnostic to the record's frontmatter"
	case "v2-invalid-id":
		return "Set the record's id or reference to a valid family identity: O-### (Objective), T-### (Task), C-### (Check), or I-### (Issue reference), each with at least three digits"
	case "v2-release-invalid-id":
		return "Set the Release's id to a valid R-### identity with at least three digits"
	case "v2-invalid-ownership":
		return "Set the Task's objective field to exactly one existing O-### Objective id"
	case "v2-invalid-lifecycle":
		return "Set status (and stage while status is in_progress) to a supported V2 lifecycle value"
	case "v2-invalid-dependency":
		return "Fix the depends_on entry: task must be a T-### id and requires must be clear or accepted"
	case "v2-invalid-release-reference":
		return "Set the Objective's release field to an existing R-### Release, or remove it for an unassigned Objective"
	case "v2-duplicate-id":
		return "Rename one of the two records reporting the same id so each global id is declared once"
	case "v2-path-mismatch":
		return "Rename the directory or file so its name starts with the id the record declares"
	case "v2-unsafe-path":
		return "Remove the symlink or case-aliasing path reported in the diagnostic; V2 records must resolve inside the project root"
	case "v2-missing-owner":
		return "Create the referenced O-### Objective, or fix the Task's objective field to reference one that exists"
	case "v2-missing-release":
		return "Create the referenced R-### Release, or remove or correct the Objective's release field"
	case "v2-release-missing-section":
		return "Add the missing required Release body section: Outcome, Why, Success Conditions, or Boundaries"
	case "v2-release-legacy-malformed":
		return "Fix legacy_completion.source_path, archive_path, and sha256 on the named historical Release"
	case "v2-missing-dependency-target":
		return "Create the referenced dependency record, or remove it from depends_on"
	case "v2-self-dependency":
		return "Remove the record's own id from its depends_on list"
	case "v2-dependency-cycle":
		return "Break the circular dependency chain named in the diagnostic"
	case "v2-record-malformed":
		return "Fix the YAML frontmatter between the --- delimiters in the named record file"
	case "v2-check-malformed":
		return "Fix the named Check field in the record's frontmatter — result must be CLEAR or NEEDS WORK, checked_by.role a supported actor role, and checked_at a parseable RFC 3339 timestamp"
	case "v2-check-missing-scope-target":
		return "Create the Task or Objective the Check's scope names, or fix scope.id to reference one that exists"
	case "v2-check-missing-release-scope-target":
		return "Create the Release the Check's scope names, or fix scope.id to reference one that exists"
	case "v2-check-missing-reference":
		return "Create the Check named in supersedes, or remove the supersedes field"
	case "v2-check-supersedes-conflict":
		return "Fix the supersedes chain: it must name a Check with the same scope, no two Checks may supersede the same target, and the chain must not cycle"
	case "v2-evidence-malformed":
		return "Fix the named evidence field in the record's frontmatter — state and role values must be one of the supported options and timestamps must be parseable RFC 3339"
	case "v2-evidence-missing-reference":
		return "Fix the evidence field to name a Check that exists, or remove the reference"
	case "v2-check-immutable":
		return "A Check record is immutable once written — record a new Check with supersedes naming this one instead of editing it"
	case "v2-issue-malformed":
		return "Fix the named Issue field in the record's frontmatter — type, status, source, resolution, and history each have a fixed vocabulary"
	case "v2-issue-missing-duplicate-target":
		return "Set duplicate_of to an existing I-### Issue id, or remove the field"
	case "v2-issue-missing-escalation-target":
		return "Set escalated_to to an existing O-### Objective id, or remove the field"
	case "v2-issue-self-duplicate":
		return "Remove the Issue's own id from its duplicate_of field, or point it at a different canonical Issue"
	case "v2-issue-duplicate-cycle":
		return "Break the circular duplicate_of chain named in the diagnostic so one Issue in the ring is canonical"
	case "v2-issue-missing-link-target":
		return "Fix the named tasks, checks, or issues reference to name a record that exists, or remove it"
	case "v2-issue-unpaired-check-link":
		return "Add the Check's id to the named Issue's checks field, or remove the Issue from the Check's issues field"
	case "v2-issue-resolution-required":
		return "Add a resolution block naming a disposition, or set the Issue's status back to open or in_progress"
	case "v2-issue-resolution-not-allowed":
		return "Remove the resolution block, or set the Issue's status to resolved"
	case "v2-issue-resolution-missing-proof":
		return "Set resolution.check to an existing, CLEAR C-### Check that appears in the Issue's checks field"
	case "v2-issue-resolution-unusable-proof":
		return "Point resolution.check at a Check that recorded CLEAR and is listed in the Issue's checks field"
	case "v2-issue-resolution-field-mismatch":
		return "Fix the resolution field for its disposition: accepted needs an owner actor and reason and no proof check; duplicate needs duplicate_of and no proof check; escalated needs a planner actor, escalated_to, and no proof check"
	case "v2-issue-already-exists":
		return "Choose a different Issue id or file path — an existing Issue record cannot be overwritten by create"
	case "v2-issue-history-not-append-only":
		return "Append new history entries after the recorded ones rather than editing, reordering, or removing any existing entry"
	default:
		return "Review the V2 project diagnostic and fix the reported record"
	}
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
		return "Record a Check and a current freshness assessment before treating the Task as done, or set status back to reflect its real progress"
	case "v2-acceptance-superseded":
		return "Have the owner accept the Task's actual latest Check; an acceptance naming a superseded Check no longer applies"
	case "v2-evidence-contradicts-status":
		return "Advance the Task's status to done to match its recorded evidence, or correct the evidence if completion was not actually reached"
	case "v2-objective-done-without-clearance":
		return "Record an Objective-scoped Check and a current freshness assessment before treating the Objective as done, or set status back to reflect its real progress"
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
