---
id: E34-audit-register-doctor-and-quality-gates/T002-doctor-report-and-repair-guidance
status: planned
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

- [ ] Audit-register diagnostics appear in plain doctor output with stable wording.
- [ ] Repair suggestions name the exact file and field to edit.
- [ ] Missing proof suggestions explain that `verified` requires named proof.
- [ ] Broken link suggestions distinguish task, defect, and duplicate finding references.
- [ ] Doctor does not offer destructive or automatic repairs for audit-register files.

## Implementation Plan

- [ ] Add report formatting for audit-register check results.
- [ ] Add typed repair suggestions for each diagnostic category.
- [ ] Add command-level output tests.
- [ ] Verify existing doctor output remains stable when no audit diagnostics exist.

## Context Log

Pending.
