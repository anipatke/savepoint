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

func TestV2ProblemRepair_checkAndEvidenceNames(t *testing.T) {
	tests := []struct {
		name     string
		contains string
	}{
		{"v2-check-malformed", "result must be CLEAR or NEEDS WORK"},
		{"v2-check-missing-scope-target", "Create the Task or Objective"},
		{"v2-check-missing-release-scope-target", "Create the Release"},
		{"v2-release-invalid-id", "valid R### identity"},
		{"v2-invalid-release-reference", "existing R### Release"},
		{"v2-missing-release", "referenced R### Release"},
		{"v2-release-missing-section", "required Release body section"},
		{"v2-release-legacy-malformed", "legacy_completion"},
		{"v2-check-missing-reference", "Create the Check named in supersedes"},
		{"v2-check-supersedes-conflict", "supersedes chain"},
		{"v2-evidence-malformed", "evidence field"},
		{"v2-evidence-missing-reference", "name a Check that exists"},
		{"v2-check-immutable", "immutable once written"},
		{"v2-issue-malformed", "type, status, source, resolution, and history"},
		{"v2-issue-missing-duplicate-target", "duplicate_of to an existing I### Issue"},
		{"v2-issue-self-duplicate", "own id from its duplicate_of"},
		{"v2-issue-duplicate-cycle", "circular duplicate_of chain"},
		{"v2-issue-missing-link-target", "tasks, checks, or issues reference"},
		{"v2-issue-unpaired-check-link", "checks field"},
		{"v2-issue-resolution-required", "resolution block naming a disposition"},
		{"v2-issue-resolution-not-allowed", "Remove the resolution block"},
		{"v2-issue-resolution-missing-proof", "CLEAR C### Check"},
		{"v2-issue-resolution-unusable-proof", "recorded CLEAR"},
		{"v2-issue-resolution-field-mismatch", "accepted needs an owner actor"},
		{"v2-issue-already-exists", "cannot be overwritten by create"},
		{"v2-issue-history-not-append-only", "Append new history entries"},
		{"unknown-name", "Review the V2 project diagnostic"},
	}
	for _, tt := range tests {
		got := V2ProblemRepair(tt.name)
		if !strings.Contains(got, tt.contains) {
			t.Errorf("V2ProblemRepair(%q) = %q, want containing %q", tt.name, got, tt.contains)
		}
	}
}

func TestV2ProblemRepair_invalidIDNamesEveryV2RecordFamily(t *testing.T) {
	got := V2ProblemRepair("v2-invalid-id")
	for _, identity := range []string{"O### (Objective)", "T### (Task)", "C### (Check)", "I### (Issue reference)"} {
		if !strings.Contains(got, identity) {
			t.Errorf("V2ProblemRepair(v2-invalid-id) = %q, want %q guidance", got, identity)
		}
	}
}

func TestV2ConsistencyRepair(t *testing.T) {
	tests := []struct {
		name     string
		contains string
	}{
		{"v2-done-without-clearance", "Record a Check and a current freshness assessment"},
		{"v2-acceptance-superseded", "owner accept the Task's actual latest Check"},
		{"v2-evidence-contradicts-status", "Advance the Task's status to done"},
		{"v2-objective-done-without-clearance", "Record an Objective-scoped Check"},
		{"v2-objective-done-with-incomplete-task", "Finish the named owned Task"},
		{"v2-issue-verified-proof-superseded", "Record a fresh proof Check"},
		{"unknown-name", "Review the Task's recorded evidence"},
	}
	for _, tt := range tests {
		got := V2ConsistencyRepair(tt.name)
		if !strings.Contains(got, tt.contains) {
			t.Errorf("V2ConsistencyRepair(%q) = %q, want containing %q", tt.name, got, tt.contains)
		}
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
