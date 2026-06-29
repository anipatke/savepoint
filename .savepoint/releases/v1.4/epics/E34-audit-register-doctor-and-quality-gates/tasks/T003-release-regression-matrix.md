---
id: E34-audit-register-doctor-and-quality-gates/T003-release-regression-matrix
status: planned
objective: Prove v1.4 behavior with a focused release regression matrix.
depends_on:
  - E30-audit-register-templates/T004-upgrade-assets-and-template-regression
  - E32-audit-register-tui-review/T004-help-footer-and-render-regression
  - E34-audit-register-doctor-and-quality-gates/T002-doctor-report-and-repair-guidance
complexity_tier: medium
complexity_reason: The matrix coordinates package-level coverage across the release surface.
---

# T003: Release Regression Matrix

## Problem

The release touches scaffold, data parsing, board review, workflow guidance, and doctor diagnostics; it needs a final matrix proving the old workflow still works.

## Context Files

- `internal/init/integration_test.go`
- `internal/data/fuzz_test.go`
- `internal/board/integration_test.go`
- `internal/doctor/checks_test.go`
- `README.md`
- `.savepoint/releases/v1.4/v1.4-PRD.md`

## Acceptance Criteria

- [ ] New-project scaffold includes audit-register assets and existing core assets.
- [ ] Existing project upgrade preserves edited audit-register files.
- [ ] Data parsing tolerates absent audit registers and reports malformed present files.
- [ ] Board TUI opens and closes the Audit Register overlay without changing task, defect, release-doc, or epic-detail behavior.
- [ ] Doctor reports audit-register issues only when audit-register files are present and malformed.
- [ ] `make build && make test` passes.

## Implementation Plan

- [ ] Run targeted package tests for init, data, board, and doctor.
- [ ] Add or adjust integration coverage for the release-level acceptance paths.
- [ ] Run the full quality gate.
- [ ] Record quality-gate results in the task context log.

## Context Log

Pending.
