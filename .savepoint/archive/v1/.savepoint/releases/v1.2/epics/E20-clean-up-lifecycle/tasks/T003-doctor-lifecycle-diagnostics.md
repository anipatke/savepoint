---
id: E20-clean-up-lifecycle/T003-doctor-lifecycle-diagnostics
status: done
objective: Align doctor task lifecycle diagnostics and repair suggestions with the shared contract.
depends_on: [T002-parser-writer-lifecycle]
complexity_tier: medium
complexity_reason: Doctor diagnostics need shared policy adoption plus repair and test updates.
---

# T003: Doctor Lifecycle Diagnostics

## Problem

Doctor currently re-implements task lifecycle checks directly from raw frontmatter. It should remain stricter than board loading where appropriate, but it should use the same lifecycle policy so diagnostics do not drift.

## Context Files

- `internal/doctor/checks.go`
- `internal/doctor/checks_test.go`
- `internal/doctor/repairs.go`
- `internal/doctor/repairs_test.go`
- `internal/data/lifecycle.go`
- `internal/data/lifecycle_test.go`

## Acceptance Criteria

- [x] Doctor reports lifecycle issues using shared lifecycle diagnostics rather than hand-rolled status/stage branching.
- [x] Doctor flags missing status, missing in-progress stage, invalid in-progress stage, legacy `phase`, and stale non-in-progress `stage`.
- [x] Repair suggestions remain concrete and canonical.
- [x] Doctor stays read-only and does not auto-repair lifecycle metadata.

## Implementation Plan

- [x] Add or reuse lifecycle diagnostic output suitable for doctor.
- [x] Update `checkTaskLifecycle` to consume lifecycle diagnostics.
- [x] Update repair suggestion matching only where message text changes.
- [x] Add tests for compatibility cases that board can load but doctor should still report.

## Context Log

- Read: `internal/doctor/checks.go`, `internal/doctor/checks_test.go`, `internal/doctor/repairs.go`, `internal/doctor/repairs_test.go`, `internal/data/lifecycle.go`, `internal/data/lifecycle_test.go`.
- Edited: `internal/data/lifecycle.go`, `internal/data/lifecycle_test.go`, `internal/doctor/checks.go`, `internal/doctor/checks_test.go`, `internal/doctor/repairs.go`, `internal/doctor/repairs_test.go`.
- Added shared task lifecycle diagnostics in `internal/data` for missing status, invalid status, legacy `phase`, missing/invalid in-progress stage, and stale non-`in_progress` `stage`.
- Updated doctor structure checks to consume lifecycle diagnostics instead of duplicating status/stage branching.
- Updated legacy `phase` repair guidance so canonical fixes remove `phase` outside `in_progress` instead of blindly renaming it.
- Added tests for board-load-compatible metadata that doctor still reports, and a read-only assertion that diagnostics do not rewrite task files.
- Quality gates: `go test ./internal/data ./internal/doctor`, `make build`, and `make test` passed.
