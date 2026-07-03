package doctor

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/opencode/savepoint/internal/data"
)

func TestSuggestRepair(t *testing.T) {
	tests := []struct {
		err      error
		contains string
	}{
		{Problem{Message: "config.yml not found"}, "savepoint init"},
		{Problem{Message: "config.yml missing required field: theme"}, "Add the missing field"},
		{Problem{Message: "config.yml invalid YAML: yaml: line 3: could not find expected"}, "Fix the YAML syntax"},
		{Problem{Message: "router.md not found"}, "savepoint init"},
		{Problem{Message: "router.md unknown state \"bogus\""}, "Set router state to a recognized"},
		{Problem{Message: "router.md release \"v99\" directory not found"}, "Create the release directory"},
		{Problem{Message: "router.md epic \"E99-foo\" directory not found"}, "Create the epic directory"},
		{Problem{Message: "release PRD file not found"}, "Create a {release}-PRD.md"},
		{Problem{Message: "epic detail file not found"}, "Create an E##-Detail.md"},
		{Problem{Message: `epic uses non-canonical status "epic-design"; use planned, in_progress, done, or audited (loads as planned)`}, "Set the epic status to planned, in_progress, done, or audited"},
		{Problem{Message: `epic status invalid "wip"; use planned, in_progress, done, or audited (loads as planned)`}, "Set the epic status to planned, in_progress, done, or audited"},
		{Problem{Message: "invalid frontmatter: yaml: line 5:"}, "Fix the YAML frontmatter"},
		{Problem{Message: "task missing required frontmatter field: status"}, "Add the missing field"},
		{Problem{Message: "task uses non-canonical status \"complete\"; replace with \"done\""}, "canonical value"},
		{Problem{Message: "task uses legacy frontmatter field phase"}, "Use stage: build"},
		{Problem{Message: "task stage is required when status is in_progress"}, "Add stage: build"},
		{Problem{Message: "task stage invalid: invalid stage \"done\""}, "Set stage to build"},
		{Problem{Message: "task stage field \"build\" is only valid when status is in_progress"}, "Remove stage"},
		{Problem{Message: "task complexity invalid: complexity_reason has 21 words; maximum is 20 words"}, "Shorten complexity_reason"},
		{Problem{Message: "task missing ## Acceptance Criteria section"}, "Add an ## Acceptance Criteria section"},
		{Problem{Message: "task frontmatter field depends_on must be a list"}, "Change depends_on to a YAML list"},
		{Problem{Message: "depends_on references non-existent task \"E99/T999\""}, "Create the referenced task"},
		{Problem{Message: "duplicate task ID \"E01-foo/T001-task\" (first seen in"}, "Rename one of the tasks"},
		{Problem{Message: "dependency cycle detected:"}, "Break the circular dependency chain"},
		{Problem{Message: "audit proposal exists but router state is"}, "Set router state to audit-pending"},
		{Problem{Message: "orphaned task: epic \"E99-ghost\" does not exist"}, "Move the task directory"},
		{Problem{Message: "quality gate \"lint\" failed"}, "Fix the issue reported by the quality gate"},
		{Problem{Message: "some random unknown problem"}, "Review the file and fix"},
	}
	for _, tt := range tests {
		got := SuggestRepair(tt.err)
		if !strings.Contains(got, tt.contains) {
			t.Errorf("SuggestRepair(%q) = %q, want containing %q", tt.err.Error(), got, tt.contains)
		}
	}
}

func TestGateSuggestion(t *testing.T) {
	tests := []struct {
		name     string
		contains string
	}{
		{"lint", "make lint"},
		{"typecheck", "make typecheck"},
		{"test", "make test"},
		{"custom", "Run \"custom\" locally"},
	}
	for _, tt := range tests {
		got := GateSuggestion(tt.name)
		if !strings.Contains(got, tt.contains) {
			t.Errorf("GateSuggestion(%q) = %q, want containing %q", tt.name, got, tt.contains)
		}
	}
}

func TestSuggestRepair_typedErrors(t *testing.T) {
	tests := []struct {
		err      error
		contains string
	}{
		{fmt.Errorf("config.yml not found: %w", data.ErrConfigNotFound), "savepoint init"},
		{fmt.Errorf("router.md not found: %w", data.ErrConfigNotFound), "savepoint init"},
	}
	for _, tt := range tests {
		if !errors.Is(tt.err, data.ErrConfigNotFound) {
			t.Fatalf("test error should wrap %v", data.ErrConfigNotFound)
		}
		got := SuggestRepair(tt.err)
		if !strings.Contains(got, tt.contains) {
			t.Errorf("SuggestRepair(%q) = %q, want containing %q", tt.err.Error(), got, tt.contains)
		}
	}
}

func TestAuditFindingRepair(t *testing.T) {
	tests := []struct {
		code     data.FindingDiagnosticCode
		contains string
	}{
		{data.FindingMissingFieldCode, "Add the named field"},
		{data.FindingInvalidIDCode, "id field"},
		{data.FindingIDMismatchCode, "match the filename's F### id"},
		{data.FindingInvalidStatusCode, "status field"},
		{data.FindingInvalidSeverityCode, "severity field to critical, high, medium, or low"},
		{data.FindingInvalidConfidenceCode, "confidence field to high, medium, or low"},
		{data.FindingDiagnosticCode("unknown"), "Review the finding file"},
	}
	for _, tt := range tests {
		got := AuditFindingRepair(tt.code)
		if !strings.Contains(got, tt.contains) {
			t.Errorf("AuditFindingRepair(%q) = %q, want containing %q", tt.code, got, tt.contains)
		}
	}
}

func TestAuditValidationRepair(t *testing.T) {
	tests := []struct {
		code     data.AuditValidationCode
		contains string
	}{
		{data.AuditVerifiedMissingProof, "verified finding requires named proof"},
		{data.AuditVerifiedMissingProof, "verified_proof field"},
		{data.AuditDuplicateMissingTarget, "duplicate_of field"},
		{data.AuditDeferredMissingRationale, "deferral_reason field"},
		{data.AuditOwnerDecisionMissingRationale, "deferral_reason or waiver_reason"},
		{data.AuditWaivedMissingRationale, "waiver_reason field"},
		{data.AuditReleaseRefMissing, "releases entry"},
		{data.AuditEpicRefMissing, "epics entry"},
		{data.AuditTaskRefMissing, "tasks entry"},
		{data.AuditDefectRefMissing, "defects entry"},
		{data.AuditDuplicateRefMissing, "duplicate_of field"},
		{data.AuditValidationCode("unknown"), "Review the finding file"},
	}
	for _, tt := range tests {
		got := AuditValidationRepair(tt.code)
		if !strings.Contains(got, tt.contains) {
			t.Errorf("AuditValidationRepair(%q) = %q, want containing %q", tt.code, got, tt.contains)
		}
	}

	verifiedRepair := AuditValidationRepair(data.AuditVerifiedMissingProof)
	if strings.Contains(verifiedRepair, "## Proof") {
		t.Errorf("AuditValidationRepair(%q) = %q, should require verified_proof only", data.AuditVerifiedMissingProof, verifiedRepair)
	}
}

// TestAuditRepairsDistinguishReferenceKinds proves the broken-link suggestions
// name different frontmatter lists for task, defect, and duplicate references.
func TestAuditRepairsDistinguishReferenceKinds(t *testing.T) {
	task := AuditValidationRepair(data.AuditTaskRefMissing)
	defect := AuditValidationRepair(data.AuditDefectRefMissing)
	duplicate := AuditValidationRepair(data.AuditDuplicateRefMissing)
	if task == defect || task == duplicate || defect == duplicate {
		t.Fatalf("repair suggestions must distinguish reference kinds: task=%q defect=%q duplicate=%q", task, defect, duplicate)
	}
}
