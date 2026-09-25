---
id: E12-validation-fix/T005-tests
status: done
objective: Verify all validation scenarios work correctly
depends_on:
    - E12-validation-fix/T001-default-phase
    - E12-validation-fix/T002-default-status
    - E12-validation-fix/T003-better-errors
    - E12-validation-fix/T004-validate-on-write
---

# T005: Tests

## Context Files

- `internal/data/parser_test.go`
- `internal/data/lifecycle_test.go`
- `internal/data/write_test.go`

## Acceptance Criteria

- [x] All new defaults work in tests
- [x] Error messages tested
- [x] go test ./... passes

## Implementation Plan

- [x] Run `go test ./...` - all pass
- [x] Updated lifecycle tests for new behavior:
  - TestValidateTaskLifecycle_allowsPlannedWithoutPhase ✓
  - TestValidateTaskLifecycle_allowsInProgressWithPhase ✓
  - TestValidateTaskLifecycle_rejectsUnknownStatus ✓
  - TestValidateTaskLifecycle_rejectsPhaseOutsideInProgress ✓
- [x] Updated parser tests:
  - TestParseTaskFile_includesDefaultBuildForInProgress ✓
  - TestParseTaskFile_allowsPhaseOutsideInProgress ✓
- [x] All other existing tests pass
