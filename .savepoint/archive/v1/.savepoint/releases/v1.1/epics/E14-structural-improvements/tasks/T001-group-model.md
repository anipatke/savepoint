---
id: E14-structural-improvements/T001-group-model
status: done
objective: Group the 27-field Model struct into focused sub-structs; delete newProgramModel and loadAllTasks
depends_on: []
---

# T001: Group Model Fields Into Sub-structs and Remove Dead Code

## Context Files

- `internal/board/model.go` — Model struct with 27 flat fields
- `internal/board/board.go:32-34` — newProgramModel with hardcoded slug
- `internal/board/board.go:125-128` — loadAllTasks dead wrapper
- `internal/board/board_test.go:13` — references newProgramModel

## Acceptance Criteria

- [x] Model fields grouped into sub-structs (e.g. ViewConfig, NavigationState, DataState)
- [x] newProgramModel() removed from board.go
- [x] loadAllTasks() removed from board.go
- [x] board_test.go updated to not reference removed functions
- [x] `go test ./...` passes with no regressions

## Implementation Plan

- [x] Identify logical field groups in Model struct
- [x] Extract sub-struct type definitions
- [x] Update all field references across the board package
- [x] Delete newProgramModel() and loadAllTasks()
- [x] Update board_test.go to construct Model directly
- [x] Run `make build && make test`

## Context Log

- Files read: `.savepoint/router.md`, `.savepoint/releases/v1.1/epics/E14-structural-improvements/E14-Detail.md`, `.savepoint/releases/v1.1/epics/E14-structural-improvements/tasks/T001-group-model.md`, `agent-skills/savepoint-build-task/SKILL.md`, `internal/board/model.go`, `internal/board/board.go`, `internal/board/board_test.go`
- Files edited: `.savepoint/releases/v1.1/epics/E14-structural-improvements/tasks/T001-group-model.md`, `internal/board/model.go`, `internal/board/board.go`, `internal/board/board_test.go`
- Tokens used: approximately 16k
- Quality gates: `make build` passed; `make test` passed. Initial sandboxed runs hit Go build-cache access errors, then passed with approved escalation.
