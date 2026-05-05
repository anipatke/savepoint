---
id: E14-structural-improvements/T006-unify-enums
status: done
objective: Consolidate ColumnType and TaskStatus into a single status type and remove syncTaskStatus
depends_on: []
---

# T006: Unify ColumnType and TaskStatus Enumerations

## Context Files

- `internal/data/task.go:13-36` — ColumnType and TaskStatus parallel definitions
- `internal/board/transitions.go:57-59` — syncTaskStatus manually syncs both enums
- Multiple consumer files reference Task.Status or Task.Column

## Acceptance Criteria

- [x] ColumnType and TaskStatus consolidated into a single type
- [x] syncTaskStatus removed from transitions.go
- [x] All consumer references updated
- [x] `go test ./...` passes with no regressions

## Implementation Plan

- [x] Choose the surviving type (ColumnType has wider usage)
- [x] Merge TaskStatus constants into ColumnType
- [x] Remove syncTaskStatus and inline any callers
- [x] Update all references across the codebase
- [x] Run `make build && make test`
