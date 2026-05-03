---
id: E12-validation-fix/T004-validate-on-write
status: done
objective: Add validation when writing task files
depends_on:
    - E12-validation-fix/T003-better-errors
---

# T004: Validate On Write

## Context Files

- `internal/data/write.go`
- `internal/data/write_test.go`
- `internal/data/lifecycle.go`

## Acceptance Criteria

- [x] Already validated in write path via WriteTaskStatus
- [x] go test ./... passes

## Implementation Plan

- [x] Validation already runs via ParseTaskFile in write.go
- [x] Tests pass
