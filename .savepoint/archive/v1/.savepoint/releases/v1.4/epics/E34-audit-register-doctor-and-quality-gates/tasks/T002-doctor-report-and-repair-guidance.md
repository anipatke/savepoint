---
id: E34-audit-register-doctor-and-quality-gates/T002-doctor-report-and-repair-guidance
status: done
objective: Add report output and repair guidance for audit-register diagnostics.
depends_on:
    - E34-audit-register-doctor-and-quality-gates/T001-doctor-audit-diagnostics
complexity_tier: medium
complexity_reason: Report and repair wording must stay actionable without introducing auto-repair.
---

# T002: Doctor Report and Repair Guidance

## Problem

Doctor findings need concise messages and repair suggestions so users and agents can fix malformed audit-register files manually.

## Context Files

- `internal/doctor/report.go`
- `internal/doctor/report_test.go`
- `internal/doctor/repairs.go`
- `internal/doctor/repairs_test.go`
- `cmd/doctor.go`
- `cmd/doctor_test.go`

## Acceptance Criteria

- [x] Audit-register diagnostics appear in plain doctor output with stable wording.
- [x] Repair suggestions name the exact file and field to edit.
- [x] Missing proof suggestions explain that `verified` requires named proof.
- [x] Broken link suggestions distinguish task, defect, and duplicate finding references.
- [x] Doctor does not offer destructive or automatic repairs for audit-register files.

## Implementation Plan

- [x] Add report formatting for audit-register check results.
- [x] Add typed repair suggestions for each diagnostic category.
- [x] Add command-level output tests.
- [x] Verify existing doctor output remains stable when no audit diagnostics exist.

## Context Log

- Read: `internal/doctor/report.go`, `internal/doctor/repairs.go`, `cmd/doctor.go`, `cmd/doctor_test.go`, `main.go` (doctor wiring).
- Edited: `internal/doctor/report.go` — `DiagnosticReport` gains `AuditRegister`; `RunAllChecks` runs `CheckAuditRegister`; new "Audit Register Check" section (category `audit-register`) prints before Quality Gates; `HasProblems` includes it; `problemRepair` prefers a problem's typed `Repair` over `SuggestRepair` message matching.
- Edited: `internal/doctor/checks.go` — `Problem` gains an optional typed `Repair` field.
- Edited: `internal/doctor/repairs.go` — typed `AuditFindingRepair` (per `FindingDiagnosticCode`) and `AuditValidationRepair` (per `AuditValidationCode`); every suggestion names the frontmatter field; broken-link suggestions distinguish releases/epics/tasks/defects entries and duplicate_of; all suggestions are manual edits only.
- Tests: `internal/doctor/repairs_test.go` covers every code plus reference-kind distinction; `internal/doctor/report_test.go` asserts plain-output wording (`✗ audit-register: <file>`, message, repair line) and that output stays ALL CLEAN with no audit tree. Doctor output is what `savepoint doctor` prints verbatim (`main.go` pipes `report.Format()`); `go test ./cmd -count=1` re-run to prove CLI arg/exit behavior unchanged.
- Quality gates: `make build && make test` pass (see T003 log).
