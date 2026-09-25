---
id: E12-validation-fix/T002-default-status
status: done
objective: Default status to planned when both status and column missing
depends_on:
    - E12-validation-fix/T001-default-phase
---

# T002: Default Status to Planned

## Context Files

- `internal/data/parser.go`
- `internal/data/parser_test.go`

## Acceptance Criteria

- [x] Already implemented in normalizeColumn() - returns ColumnPlanned for empty
- [x] go test ./... passes
