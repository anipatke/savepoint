package data

import (
	"errors"
	"testing"
	"time"
)

func TestDecodeEvidenceV2_absentEvidenceIsNil(t *testing.T) {
	evidence, err := decodeEvidenceV2("test.md", "task", "T001", evidenceV2Frontmatter{})
	if err != nil {
		t.Fatalf("decodeEvidenceV2() error = %v", err)
	}
	if evidence != nil {
		t.Fatalf("decodeEvidenceV2() = %+v, want nil for a record carrying no evidence", evidence)
	}
}

func TestDecodeEvidenceV2_lastCheckOnly(t *testing.T) {
	evidence, err := decodeEvidenceV2("test.md", "task", "T001", evidenceV2Frontmatter{LastCheck: "C001"})
	if err != nil {
		t.Fatalf("decodeEvidenceV2() error = %v", err)
	}
	if evidence == nil || evidence.LastCheck != "C001" {
		t.Fatalf("decodeEvidenceV2() = %+v, want LastCheck C001", evidence)
	}
	if evidence.Freshness != nil || evidence.OwnerValidation != nil || evidence.Exception != nil || evidence.Replan != nil {
		t.Fatalf("decodeEvidenceV2() = %+v, want every other sub-block nil", evidence)
	}
}

func TestDecodeEvidenceV2_lastCheckMalformed(t *testing.T) {
	_, err := decodeEvidenceV2("test.md", "task", "T001", evidenceV2Frontmatter{LastCheck: "C1"})
	if !errors.Is(err, ErrV2InvalidID) {
		t.Fatalf("decodeEvidenceV2() error = %v, want ErrV2InvalidID", err)
	}
}

func validFreshnessFrontmatter() freshnessV2Frontmatter {
	return freshnessV2Frontmatter{
		State:      "current",
		Check:      "C001",
		AssessedBy: evidenceActorFrontmatter{Role: "checker", Session: "sess-1"},
		AssessedAt: "2026-09-15T00:00:00Z",
		Basis:      "reviewed the diff",
	}
}

func TestDecodeEvidenceV2_freshnessValid(t *testing.T) {
	fresh := validFreshnessFrontmatter()
	evidence, err := decodeEvidenceV2("test.md", "task", "T001", evidenceV2Frontmatter{Freshness: &fresh})
	if err != nil {
		t.Fatalf("decodeEvidenceV2() error = %v", err)
	}
	if evidence == nil || evidence.Freshness == nil {
		t.Fatal("decodeEvidenceV2() Freshness = nil, want a decoded Freshness")
	}
	f := evidence.Freshness
	if f.State != FreshnessCurrent || f.Check != "C001" || f.Basis != "reviewed the diff" {
		t.Errorf("Freshness = %+v, want state current, check C001, basis set", f)
	}
	if f.AssessedBy.Role != ActorRoleChecker || f.AssessedBy.Session != "sess-1" {
		t.Errorf("Freshness.AssessedBy = %+v, want checker/sess-1", f.AssessedBy)
	}
	wantTime, _ := time.Parse(time.RFC3339, "2026-09-15T00:00:00Z")
	if !f.AssessedAt.Equal(wantTime) {
		t.Errorf("Freshness.AssessedAt = %v, want %v", f.AssessedAt, wantTime)
	}
}

func TestDecodeEvidenceV2_freshnessEachStateValue(t *testing.T) {
	for _, state := range []FreshnessState{FreshnessCurrent, FreshnessStale, FreshnessUnknown} {
		t.Run(string(state), func(t *testing.T) {
			fresh := validFreshnessFrontmatter()
			fresh.State = string(state)
			evidence, err := decodeEvidenceV2("test.md", "task", "T001", evidenceV2Frontmatter{Freshness: &fresh})
			if err != nil {
				t.Fatalf("decodeEvidenceV2() error = %v", err)
			}
			if evidence.Freshness.State != state {
				t.Errorf("Freshness.State = %q, want %q", evidence.Freshness.State, state)
			}
		})
	}
}

func TestDecodeEvidenceV2_freshnessUnknownStateValueRejected(t *testing.T) {
	fresh := validFreshnessFrontmatter()
	fresh.State = "expired"
	_, err := decodeEvidenceV2("test.md", "task", "T001", evidenceV2Frontmatter{Freshness: &fresh})
	if !errors.Is(err, ErrV2EvidenceMalformed) {
		t.Fatalf("decodeEvidenceV2() error = %v, want ErrV2EvidenceMalformed", err)
	}
}

func TestDecodeEvidenceV2_currentFreshnessRequiresCheckerProvenance(t *testing.T) {
	fresh := validFreshnessFrontmatter()
	fresh.AssessedBy = evidenceActorFrontmatter{Role: "executor", Session: "executor-1"}
	_, err := decodeEvidenceV2("test.md", "task", "T001", evidenceV2Frontmatter{Freshness: &fresh})
	if !errors.Is(err, ErrV2EvidenceMalformed) {
		t.Fatalf("decodeEvidenceV2() error = %v, want ErrV2EvidenceMalformed", err)
	}
}

// TestDecodeEvidenceV2_freshnessPartiallyFilledRejected proves every
// required freshness field is enforced: clearing any one of them, with the
// rest left valid, must fail rather than decode a partial block.
func TestDecodeEvidenceV2_freshnessPartiallyFilledRejected(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*freshnessV2Frontmatter)
		wantErr error
	}{
		{"missing state", func(f *freshnessV2Frontmatter) { f.State = "" }, ErrV2MissingField},
		{"missing check", func(f *freshnessV2Frontmatter) { f.Check = "" }, ErrV2MissingField},
		{"malformed check", func(f *freshnessV2Frontmatter) { f.Check = "C1" }, ErrV2InvalidID},
		{"missing assessed_by role", func(f *freshnessV2Frontmatter) { f.AssessedBy.Role = "" }, ErrV2MissingField},
		{"missing assessed_by session", func(f *freshnessV2Frontmatter) { f.AssessedBy.Session = "" }, ErrV2MissingField},
		{"unknown assessed_by role", func(f *freshnessV2Frontmatter) { f.AssessedBy.Role = "auditor" }, ErrV2EvidenceMalformed},
		{"missing assessed_at", func(f *freshnessV2Frontmatter) { f.AssessedAt = "" }, ErrV2MissingField},
		{"unparseable assessed_at", func(f *freshnessV2Frontmatter) { f.AssessedAt = "yesterday" }, ErrV2EvidenceMalformed},
		{"missing basis", func(f *freshnessV2Frontmatter) { f.Basis = "" }, ErrV2MissingField},
		{"whitespace-only basis", func(f *freshnessV2Frontmatter) { f.Basis = "   " }, ErrV2MissingField},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fresh := validFreshnessFrontmatter()
			tt.mutate(&fresh)
			_, err := decodeEvidenceV2("test.md", "task", "T001", evidenceV2Frontmatter{Freshness: &fresh})
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("decodeEvidenceV2() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecodeEvidenceV2_ownerValidationAbsentMeansNotRequired(t *testing.T) {
	evidence, err := decodeEvidenceV2("test.md", "task", "T001", evidenceV2Frontmatter{})
	if err != nil {
		t.Fatalf("decodeEvidenceV2() error = %v", err)
	}
	if evidence != nil {
		t.Fatalf("decodeEvidenceV2() = %+v, want nil evidence entirely when nothing is set", evidence)
	}
}

func TestDecodeEvidenceV2_ownerValidationRequiredOnly(t *testing.T) {
	ov := ownerValidationV2Frontmatter{Required: true}
	evidence, err := decodeEvidenceV2("test.md", "task", "T001", evidenceV2Frontmatter{OwnerValidation: &ov})
	if err != nil {
		t.Fatalf("decodeEvidenceV2() error = %v", err)
	}
	if evidence == nil || evidence.OwnerValidation == nil {
		t.Fatal("decodeEvidenceV2() OwnerValidation = nil, want decoded block")
	}
	if !evidence.OwnerValidation.Required || evidence.OwnerValidation.AcceptedCheck != "" {
		t.Errorf("OwnerValidation = %+v, want Required true, AcceptedCheck empty", evidence.OwnerValidation)
	}
}

func TestDecodeEvidenceV2_ownerValidationWithAcceptedCheck(t *testing.T) {
	ov := ownerValidationV2Frontmatter{Required: true, AcceptedCheck: "C002", AcceptedBy: &evidenceActorFrontmatter{Role: "owner", Session: "owner-1"}}
	evidence, err := decodeEvidenceV2("test.md", "task", "T001", evidenceV2Frontmatter{OwnerValidation: &ov})
	if err != nil {
		t.Fatalf("decodeEvidenceV2() error = %v", err)
	}
	if evidence.OwnerValidation.AcceptedCheck != "C002" {
		t.Errorf("OwnerValidation.AcceptedCheck = %q, want C002", evidence.OwnerValidation.AcceptedCheck)
	}
}

func TestDecodeEvidenceV2_ownerAcceptanceRequiresProvenance(t *testing.T) {
	ov := ownerValidationV2Frontmatter{Required: true, AcceptedCheck: "C002"}
	_, err := decodeEvidenceV2("test.md", "task", "T001", evidenceV2Frontmatter{OwnerValidation: &ov})
	if !errors.Is(err, ErrV2MissingField) {
		t.Fatalf("decodeEvidenceV2() error = %v, want ErrV2MissingField", err)
	}
}

func TestDecodeEvidenceV2_ownerAcceptanceRequiresOwnerRole(t *testing.T) {
	ov := ownerValidationV2Frontmatter{
		Required:      true,
		AcceptedCheck: "C002",
		AcceptedBy:    &evidenceActorFrontmatter{Role: "executor", Session: "executor-1"},
	}
	_, err := decodeEvidenceV2("test.md", "task", "T001", evidenceV2Frontmatter{OwnerValidation: &ov})
	if !errors.Is(err, ErrV2EvidenceMalformed) {
		t.Fatalf("decodeEvidenceV2() error = %v, want ErrV2EvidenceMalformed", err)
	}
}

func TestDecodeEvidenceV2_ownerValidationMalformedAcceptedCheck(t *testing.T) {
	ov := ownerValidationV2Frontmatter{AcceptedCheck: "C1"}
	_, err := decodeEvidenceV2("test.md", "task", "T001", evidenceV2Frontmatter{OwnerValidation: &ov})
	if !errors.Is(err, ErrV2InvalidID) {
		t.Fatalf("decodeEvidenceV2() error = %v, want ErrV2InvalidID", err)
	}
}

func validExceptionFrontmatter() exceptionV2Frontmatter {
	return exceptionV2Frontmatter{
		Requirements: []string{"TEST-08", "STYLE-01"},
		Reason:       "owner accepted the tradeoff",
		Owner:        "ani",
		RecordedAt:   "2026-09-15T00:00:00Z",
		Check:        "C003",
	}
}

func TestDecodeEvidenceV2_exceptionValid(t *testing.T) {
	exception := validExceptionFrontmatter()
	evidence, err := decodeEvidenceV2("test.md", "task", "T001", evidenceV2Frontmatter{Exception: &exception})
	if err != nil {
		t.Fatalf("decodeEvidenceV2() error = %v", err)
	}
	if evidence == nil || evidence.Exception == nil {
		t.Fatal("decodeEvidenceV2() Exception = nil, want decoded block")
	}
	e := evidence.Exception
	if len(e.Requirements) != 2 || e.Requirements[0] != "TEST-08" || e.Requirements[1] != "STYLE-01" {
		t.Errorf("Exception.Requirements = %v, want [TEST-08 STYLE-01]", e.Requirements)
	}
	if e.Reason != "owner accepted the tradeoff" || e.Owner != "ani" || e.Check != "C003" {
		t.Errorf("Exception = %+v, want reason/owner/check set", e)
	}
}

// TestDecodeEvidenceV2_exceptionPartiallyFilledRejected proves every
// required exception field is enforced individually.
func TestDecodeEvidenceV2_exceptionPartiallyFilledRejected(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*exceptionV2Frontmatter)
		wantErr error
	}{
		{"missing requirements", func(e *exceptionV2Frontmatter) { e.Requirements = nil }, ErrV2MissingField},
		{"empty requirements list", func(e *exceptionV2Frontmatter) { e.Requirements = []string{} }, ErrV2MissingField},
		{"requirements contains empty entry", func(e *exceptionV2Frontmatter) { e.Requirements = []string{"TEST-08", "  "} }, ErrV2EvidenceMalformed},
		{"missing reason", func(e *exceptionV2Frontmatter) { e.Reason = "" }, ErrV2MissingField},
		{"missing owner", func(e *exceptionV2Frontmatter) { e.Owner = "" }, ErrV2MissingField},
		{"missing recorded_at", func(e *exceptionV2Frontmatter) { e.RecordedAt = "" }, ErrV2MissingField},
		{"unparseable recorded_at", func(e *exceptionV2Frontmatter) { e.RecordedAt = "soon" }, ErrV2EvidenceMalformed},
		{"missing check", func(e *exceptionV2Frontmatter) { e.Check = "" }, ErrV2MissingField},
		{"malformed check", func(e *exceptionV2Frontmatter) { e.Check = "C1" }, ErrV2InvalidID},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exception := validExceptionFrontmatter()
			tt.mutate(&exception)
			_, err := decodeEvidenceV2("test.md", "task", "T001", evidenceV2Frontmatter{Exception: &exception})
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("decodeEvidenceV2() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func validReplanFrontmatter() replanV2Frontmatter {
	return replanV2Frontmatter{
		Reason:     "scope changed mid-build",
		RecordedBy: evidenceActorFrontmatter{Role: "executor", Session: "sess-4"},
		RecordedAt: "2026-09-15T00:00:00Z",
	}
}

func TestDecodeEvidenceV2_replanValid(t *testing.T) {
	replan := validReplanFrontmatter()
	evidence, err := decodeEvidenceV2("test.md", "task", "T001", evidenceV2Frontmatter{Replan: &replan})
	if err != nil {
		t.Fatalf("decodeEvidenceV2() error = %v", err)
	}
	if evidence == nil || evidence.Replan == nil {
		t.Fatal("decodeEvidenceV2() Replan = nil, want decoded block")
	}
	if evidence.Replan.Reason != "scope changed mid-build" {
		t.Errorf("Replan.Reason = %q, want set reason", evidence.Replan.Reason)
	}
	if evidence.Replan.RecordedBy.Role != ActorRoleExecutor || evidence.Replan.RecordedBy.Session != "sess-4" {
		t.Errorf("Replan.RecordedBy = %+v, want executor/sess-4", evidence.Replan.RecordedBy)
	}
}

func TestDecodeEvidenceV2_replanPartiallyFilledRejected(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*replanV2Frontmatter)
		wantErr error
	}{
		{"missing reason", func(r *replanV2Frontmatter) { r.Reason = "" }, ErrV2MissingField},
		{"missing recorded_by role", func(r *replanV2Frontmatter) { r.RecordedBy.Role = "" }, ErrV2MissingField},
		{"unknown recorded_by role", func(r *replanV2Frontmatter) { r.RecordedBy.Role = "auditor" }, ErrV2EvidenceMalformed},
		{"missing recorded_by session", func(r *replanV2Frontmatter) { r.RecordedBy.Session = "" }, ErrV2MissingField},
		{"missing recorded_at", func(r *replanV2Frontmatter) { r.RecordedAt = "" }, ErrV2MissingField},
		{"unparseable recorded_at", func(r *replanV2Frontmatter) { r.RecordedAt = "later" }, ErrV2EvidenceMalformed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			replan := validReplanFrontmatter()
			tt.mutate(&replan)
			_, err := decodeEvidenceV2("test.md", "task", "T001", evidenceV2Frontmatter{Replan: &replan})
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("decodeEvidenceV2() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

// TestDecodeEvidenceV2_fullBlockAllSubBlocksTogether proves last_check,
// freshness, owner_validation, exception, and replan all decode together on
// one record without interfering with each other.
func TestDecodeEvidenceV2_fullBlockAllSubBlocksTogether(t *testing.T) {
	fresh := validFreshnessFrontmatter()
	ov := ownerValidationV2Frontmatter{Required: true, AcceptedCheck: "C001", AcceptedBy: &evidenceActorFrontmatter{Role: "owner", Session: "owner-1"}}
	exception := validExceptionFrontmatter()
	replan := validReplanFrontmatter()

	evidence, err := decodeEvidenceV2("test.md", "objective", "O001", evidenceV2Frontmatter{
		LastCheck:       "C001",
		Freshness:       &fresh,
		OwnerValidation: &ov,
		Exception:       &exception,
		Replan:          &replan,
	})
	if err != nil {
		t.Fatalf("decodeEvidenceV2() error = %v", err)
	}
	if evidence.LastCheck != "C001" || evidence.Freshness == nil || evidence.OwnerValidation == nil || evidence.Exception == nil || evidence.Replan == nil {
		t.Fatalf("decodeEvidenceV2() = %+v, want every sub-block populated", evidence)
	}
}
