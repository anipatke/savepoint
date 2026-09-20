---
type: audit-findings
audited: 2026-07-03
---

# Audit Findings: E34 Audit Register Doctor and Quality Gates

## Main Findings

E34 is largely complete against its task acceptance criteria. Doctor now reports audit-register field, lifecycle, and link diagnostics; plain report output includes a dedicated Audit Register Check section with typed repair suggestions; absent or partially adopted audit-register trees remain clean; and release regression coverage exists across doctor, data, board, init scaffold, and upgrade behavior.

Open finding: verified-proof repair guidance is misleading. `internal/data/audit_validate.go` only clears `AuditVerifiedMissingProof` when `verified_proof` is non-empty, but both the validation message and doctor repair suggestion tell the user they can alternatively record proof in `## Proof`. A user following that body-section instruction would still see `savepoint doctor` report the same problem. Proposed changes below align the messages with the implemented rule and add a regression assertion.

Verification run during audit:

- `go test ./internal/doctor ./internal/data ./internal/board ./internal/init -count=1` passed.
- `go test ./cmd -count=1` passed.
- `make build` passed.
- `make test` passed.

Process note: router state is still `task-building` for E34 and the E34 detail frontmatter still says `status: planned`, even though all E34 tasks are marked `done`. This audit was still performed because the user explicitly requested it. Apply/close should update the normal audit handoff metadata after the proposed changes are approved.

## Code Style Review

- [x] One job per file
- [x] One job per function
- [ ] Test branches - verified-proof repair wording is not tested tightly enough to catch the unsupported `## Proof` alternative.
- [x] Types document intent
- [x] Build only what is needed
- [x] Handle errors at boundaries
- [ ] One source of truth - verified-proof guidance currently diverges from validator behavior.
- [x] Comments explain WHY
- [x] Content in data files
- [x] Small diffs

## Proposed Changes

### Target File
internal/data/audit_validate.go

### Replace
```go
				add(f, AuditVerifiedMissingProof,
					"verified finding has no named proof; set verified_proof or record the passing regression test in ## Proof")
```

### With
```go
				add(f, AuditVerifiedMissingProof,
					"verified finding has no named proof; set verified_proof to the passing regression test")
```

### Target File
internal/doctor/repairs.go

### Replace
```go
	case data.AuditVerifiedMissingProof:
		return "A verified finding requires named proof — set the verified_proof field to the passing regression test, or record it in the ## Proof section"
```

### With
```go
	case data.AuditVerifiedMissingProof:
		return "A verified finding requires named proof — set the verified_proof field to the passing regression test"
```

### Target File
internal/doctor/repairs_test.go

### Replace
```go
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
}
```

### With
```go
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
```
