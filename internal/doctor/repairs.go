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
		return "Fix the tasks entry to name an existing task (T### or epic/task id), or remove the entry"
	case data.AuditDefectRefMissing:
		return "Fix the defects entry to name an existing defect id, or remove the entry"
	case data.AuditDuplicateRefMissing:
		return "Point the duplicate_of field at an existing finding other than this one"
	default:
		return "Review the finding file and fix the reported field"
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
