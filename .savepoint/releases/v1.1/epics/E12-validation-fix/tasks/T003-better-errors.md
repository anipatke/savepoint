---
id: E12-validation-fix/T003-better-errors
status: done
objective: Improve error messages with hints for common issues
depends_on:
    - E12-validation-fix/T001-default-phase
    - E12-validation-fix/T002-default-status
---

# T003: Better Error Messages

## Context Files

- `internal/data/lifecycle.go`
- `internal/data/lifecycle_test.go`

## Acceptance Criteria

- [x] Error messages include suggested fix
- [x] Show which field is problematic with context
- [x] go test ./... passes

## Implementation Plan

- [x] Improved errors in lifecycle.go (lines 7, 16, 22)
- Examples of better messages:
  - `invalid status %q: use planned, in_progress, or done. Add 'status: planned'...`
  - `invalid phase %q: use build, test, or audit. Add 'phase: build'...`
- [x] Tests pass
