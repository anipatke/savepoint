---
id: E34-audit-register-doctor-and-quality-gates/T003-release-regression-matrix
status: done
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

- [x] New-project scaffold includes audit-register assets and existing core assets.
- [x] Existing project upgrade preserves edited audit-register files.
- [x] Data parsing tolerates absent audit registers and reports malformed present files.
- [x] Board TUI opens and closes the Audit Register overlay without changing task, defect, release-doc, or epic-detail behavior.
- [x] Doctor reports audit-register issues only when audit-register files are present and malformed.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Run targeted package tests for init, data, board, and doctor.
- [x] Add or adjust integration coverage for the release-level acceptance paths.
- [x] Run the full quality gate.
- [x] Record quality-gate results in the task context log.

## Context Log

Regression matrix (2026-07-03):

| Surface | Coverage | Result |
|---------|----------|--------|
| Scaffold (init) | `TestIntegration_EmptyDirectory` asserts audit assets alongside core assets | pass |
| Upgrade (init) | `TestIntegration_AuditAssetsDoNotAlterExistingFiles` preserves edited audit files | pass |
| Data parsing | Existing audit loader tests + new `FuzzParseFindingFile` (seeds + 15s fuzz, 3.1M execs, no failures): healed enums stay canonical, filename ID wins | pass |
| Board TUI | New `TestAuditOverlay_openCloseLeavesBoardBehaviorUnchanged` in `internal/board/integration_test.go`: real on-disk register, `A` open → load → render → esc close, then task detail (enter), defect overlay (`d`), release-docs overlay (`D`) all behave as before | pass |
| Doctor | T001/T002 suites: absent register clean, partial adoption clean, malformed present files reported | pass |

- Targeted packages: `go test ./internal/doctor ./internal/data ./internal/board ./internal/init -count=1` — all pass.
- Full quality gate: `make build && make test` — all packages pass (board 0.404s, data 0.100s, doctor 0.920s, init 2.353s, cmd re-run uncached with `-count=1`).
