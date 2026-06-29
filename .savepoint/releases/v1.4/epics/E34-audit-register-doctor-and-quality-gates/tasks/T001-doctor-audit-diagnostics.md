---
id: E34-audit-register-doctor-and-quality-gates/T001-doctor-audit-diagnostics
status: planned
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

- [ ] Doctor reports invalid audit finding IDs, statuses, severity, and confidence.
- [ ] Doctor reports verified findings with missing proof.
- [ ] Doctor reports duplicate findings whose target does not exist.
- [ ] Doctor reports broken task and defect links.
- [ ] Projects without `.savepoint/audit/` produce no audit-register diagnostic failures.

## Implementation Plan

- [ ] Add audit-register checks using the data validation API.
- [ ] Map validation results to doctor severities and messages.
- [ ] Add fixtures for clean, absent, malformed, and partially adopted audit registers.
- [ ] Add tests proving existing doctor checks still run.

## Context Log

Pending.
