---
id: E34-audit-register-doctor-and-quality-gates/T001-doctor-audit-diagnostics
status: done
objective: Add doctor diagnostics for audit-register consistency.
depends_on:
    - E31-audit-register-data-model/T004-audit-register-validation
complexity_tier: medium
complexity_reason: Doctor checks consume new validation results while preserving absent-register behavior.
---

# T001: Doctor Audit Diagnostics

## Problem

Malformed audit-register state should be visible through doctor instead of being discovered only when the board or an agent reads the files.

## Context Files

- `internal/doctor/checks.go`
- `internal/doctor/checks_test.go`
- `internal/doctor/interfaces.go`
- `internal/data/audit.go`
- `internal/data/audit_validate.go`

## Acceptance Criteria

- [x] Doctor reports invalid audit finding IDs, statuses, severity, and confidence.
- [x] Doctor reports verified findings with missing proof.
- [x] Doctor reports duplicate findings whose target does not exist.
- [x] Doctor reports broken task and defect links.
- [x] Projects without `.savepoint/audit/` produce no audit-register diagnostic failures.

## Implementation Plan

- [x] Add audit-register checks using the data validation API.
- [x] Map validation results to doctor severities and messages.
- [x] Add fixtures for clean, absent, malformed, and partially adopted audit registers.
- [x] Add tests proving existing doctor checks still run.

## Context Log

- Read: `internal/doctor/checks.go`, `internal/doctor/interfaces.go`, `internal/data/audit.go`, `internal/data/audit_validate.go`, `internal/data/audit_finding.go`, `internal/data/audit_run.go`.
- Edited: `internal/doctor/checks.go` — added `CheckAuditRegister` (structural load errors for findings/runs, per-finding `DiagnoseFinding` on raw frontmatter, `ValidateAuditFindings` against discovered work items), `checkAuditFindingFields`, `collectAuditWorkItems`, and a shared `readTaskID` helper reused by `collectTaskIDs`.
- Edited: `internal/data/audit_finding.go` — split `ParseFindingFile` into `ParseRawFindingFile` + `NormalizeFindingForLoad` so doctor can diagnose original (un-healed) field values; behavior of `ParseFindingFile` unchanged.
- Edited: `internal/doctor/interfaces.go` (+ test fake in `interfaces_test.go`) — `taskParser` gains `ParseRawFindingFile`.
- Tests: `internal/doctor/checks_test.go` — absent audit dir, partial adoption (register.md + READMEs only), clean finding, healed fields (id mismatch, status, severity, confidence), missing required fields, verified without proof, duplicate without/with unknown target, broken release/epic/task/defect links, malformed finding and run files.
- Quality gates: `go test ./internal/doctor ./internal/data` pass; `make build && make test` pass (see T003 log).
