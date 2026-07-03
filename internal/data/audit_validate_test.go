package data

import "testing"

// sampleWorkItems is the discovered project the reference tests resolve against:
// one release, one epic (referenced full and short), one task (referenced full
// and short), and one defect.
var sampleWorkItems = AuditWorkItems{
	Releases: []string{"v1.4"},
	Epics:    []string{"E31-audit-register-data-model"},
	Tasks:    []string{"E31-audit-register-data-model/T004-audit-register-validation"},
	Defects:  []string{"D001-token-leak"},
}

func TestValidateAuditFindings_Clean(t *testing.T) {
	findings := []AuditFinding{
		{
			ID:            "F001",
			Status:        FindingVerified,
			Releases:      []string{"v1.4"},
			Epics:         []string{"E31-audit-register-data-model"},
			Tasks:         []string{"E31-audit-register-data-model/T004-audit-register-validation"},
			Defects:       []string{"D001-token-leak"},
			VerifiedProof: "regression test TestRedactTokens",
		},
		{ID: "F002", Status: FindingDuplicate, DuplicateOf: "F001"},
	}
	if got := ValidateAuditFindings(findings, sampleWorkItems); len(got) != 0 {
		t.Fatalf("ValidateAuditFindings() = %v, want none for a consistent register", got)
	}
}

func TestValidateAuditFindings_Branches(t *testing.T) {
	// Every case is a single finding exercising exactly one validation branch.
	// findingSet supplies the duplicate-of resolution target where needed.
	tests := []struct {
		name     string
		finding  AuditFinding
		wantCode AuditValidationCode
		wantNone bool
	}{
		{
			name:     "verified missing proof",
			finding:  AuditFinding{ID: "F001", Status: FindingVerified},
			wantCode: AuditVerifiedMissingProof,
		},
		{
			name:     "verified with proof",
			finding:  AuditFinding{ID: "F001", Status: FindingVerified, VerifiedProof: "manual note"},
			wantNone: true,
		},
		{
			name:     "duplicate missing target",
			finding:  AuditFinding{ID: "F002", Status: FindingDuplicate},
			wantCode: AuditDuplicateMissingTarget,
		},
		{
			name:     "deferred missing rationale",
			finding:  AuditFinding{ID: "F003", Status: FindingDeferred},
			wantCode: AuditDeferredMissingRationale,
		},
		{
			name:     "deferred with rationale",
			finding:  AuditFinding{ID: "F003", Status: FindingDeferred, DeferralReason: "blocked on upstream"},
			wantNone: true,
		},
		{
			name:     "owner_decision missing rationale",
			finding:  AuditFinding{ID: "F004", Status: FindingOwnerDecision},
			wantCode: AuditOwnerDecisionMissingRationale,
		},
		{
			name:     "owner_decision with waiver reason",
			finding:  AuditFinding{ID: "F004", Status: FindingOwnerDecision, WaiverReason: "owner accepted risk"},
			wantNone: true,
		},
		{
			name:     "waived missing rationale",
			finding:  AuditFinding{ID: "F005", Status: FindingWaived},
			wantCode: AuditWaivedMissingRationale,
		},
		{
			name:     "waived with deferral reason (template compat)",
			finding:  AuditFinding{ID: "F005", Status: FindingWaived, DeferralReason: "documented exception"},
			wantNone: true,
		},
		{
			name:     "unknown release",
			finding:  AuditFinding{ID: "F006", Status: FindingOpen, Releases: []string{"v9.9"}},
			wantCode: AuditReleaseRefMissing,
		},
		{
			name:     "known release",
			finding:  AuditFinding{ID: "F006", Status: FindingOpen, Releases: []string{"v1.4"}},
			wantNone: true,
		},
		{
			name:     "unknown epic",
			finding:  AuditFinding{ID: "F007", Status: FindingOpen, Epics: []string{"E99-missing"}},
			wantCode: AuditEpicRefMissing,
		},
		{
			name:     "wrong epic slug does not resolve by short id",
			finding:  AuditFinding{ID: "F007", Status: FindingOpen, Epics: []string{"E31-wrong-slug"}},
			wantCode: AuditEpicRefMissing,
		},
		{
			name:     "epic by short id",
			finding:  AuditFinding{ID: "F007", Status: FindingOpen, Epics: []string{"E31"}},
			wantNone: true,
		},
		{
			name:     "unknown task",
			finding:  AuditFinding{ID: "F008", Status: FindingOpen, Tasks: []string{"E31-audit-register-data-model/T099-ghost"}},
			wantCode: AuditTaskRefMissing,
		},
		{
			name:     "wrong full task path does not resolve by short id",
			finding:  AuditFinding{ID: "F008", Status: FindingOpen, Tasks: []string{"E99-missing/T004-audit-register-validation"}},
			wantCode: AuditTaskRefMissing,
		},
		{
			name:     "task by short id",
			finding:  AuditFinding{ID: "F008", Status: FindingOpen, Tasks: []string{"T004"}},
			wantNone: true,
		},
		{
			name:     "unknown defect",
			finding:  AuditFinding{ID: "F009", Status: FindingOpen, Defects: []string{"D404-nope"}},
			wantCode: AuditDefectRefMissing,
		},
		{
			name:     "duplicate_of unknown finding",
			finding:  AuditFinding{ID: "F010", Status: FindingDuplicate, DuplicateOf: "F404"},
			wantCode: AuditDuplicateRefMissing,
		},
		{
			name:     "duplicate_of self reference",
			finding:  AuditFinding{ID: "F011", Status: FindingDuplicate, DuplicateOf: "F011"},
			wantCode: AuditDuplicateRefMissing,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateAuditFindings([]AuditFinding{tt.finding}, sampleWorkItems)
			if tt.wantNone {
				if len(got) != 0 {
					t.Fatalf("ValidateAuditFindings() = %v, want none", got)
				}
				return
			}
			if len(got) != 1 {
				t.Fatalf("ValidateAuditFindings() = %v, want exactly one result", got)
			}
			if got[0].Code != tt.wantCode {
				t.Errorf("Code = %q, want %q", got[0].Code, tt.wantCode)
			}
			if got[0].FindingID != tt.finding.ID {
				t.Errorf("FindingID = %q, want %q", got[0].FindingID, tt.finding.ID)
			}
			if got[0].Message == "" {
				t.Error("Message = empty, want actionable text")
			}
		})
	}
}

func TestValidateAuditFindings_DuplicateResolvesAgainstRegister(t *testing.T) {
	// A duplicate whose target exists in the same set is clean; only its absence
	// is a problem, so the finding set itself is the duplicate-of lookup.
	findings := []AuditFinding{
		{ID: "F001", Status: FindingOpen},
		{ID: "F002", Status: FindingDuplicate, DuplicateOf: "F001"},
	}
	if got := ValidateAuditFindings(findings, AuditWorkItems{}); len(got) != 0 {
		t.Fatalf("ValidateAuditFindings() = %v, want none when duplicate_of resolves", got)
	}
}

func TestValidateAuditFindings_Ordering(t *testing.T) {
	// One finding tripping several branches: lifecycle problems come before
	// reference problems, references in release → epic → task → defect order.
	findings := []AuditFinding{{
		ID:       "F001",
		Status:   FindingVerified,
		Releases: []string{"v9.9"},
		Epics:    []string{"E99-missing"},
		Tasks:    []string{"T999"},
		Defects:  []string{"D999"},
	}}

	got := ValidateAuditFindings(findings, sampleWorkItems)
	want := []AuditValidationCode{
		AuditVerifiedMissingProof,
		AuditReleaseRefMissing,
		AuditEpicRefMissing,
		AuditTaskRefMissing,
		AuditDefectRefMissing,
	}
	if len(got) != len(want) {
		t.Fatalf("got %d results %v, want %d", len(got), got, len(want))
	}
	for i, code := range want {
		if got[i].Code != code {
			t.Errorf("result[%d].Code = %q, want %q", i, got[i].Code, code)
		}
		if got[i].FindingID != "F001" {
			t.Errorf("result[%d].FindingID = %q, want F001", i, got[i].FindingID)
		}
	}
}
