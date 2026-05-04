---
id: E14-structural-improvements/T006-unify-enums
status: planned
objective: Consolidate ColumnType and TaskStatus into a single status type and remove syncTaskStatus
depends_on: []
---

# T006: Unify ColumnType and TaskStatus Enumerations

## Context Files

- `internal/data/task.go:13-36` — ColumnType and TaskStatus parallel definitions
- `internal/board/transitions.go:57-59` — syncTaskStatus manually syncs both enums
- Multiple consumer files reference Task.Status or Task.Column

## Acceptance Criteria

- [ ] ColumnType and TaskStatus consolidated into a single type
- [ ] syncTaskStatus removed from transitions.go
- [ ] All consumer references updated
- [ ] `go test ./...` passes with no regressions

## Implementation Plan

- [ ] Choose the surviving type (ColumnType has wider usage)
- [ ] Merge TaskStatus constants into ColumnType
- [ ] Remove syncTaskStatus and inline any callers
- [ ] Update all references across the codebase
- [ ] Run `make build && make test`
