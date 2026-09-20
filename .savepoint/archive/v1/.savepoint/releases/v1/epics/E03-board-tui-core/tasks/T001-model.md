---
id: E03-board-tui-core/T001-model
status: done
objective: "Define Model struct with all state fields"
depends_on: [E02-data-readers/T005-discovery]
---

# T001: Model Struct

## Acceptance Criteria

- `internal/board/model.go` defines Model with all state.
- Fields: Tasks, FocusedColumn, FocusedTask, SelectedEpic, SelectedRelease, Overlay, Width, Height.
- `Init()` returns `tea.Batch()` with initial commands.
- `NewModel()` accepts loaded data and returns initialized Model.

## Implementation Plan

- [x] Create `internal/board/model.go`.
- [x] Define Model struct and constructor.
- [x] Implement `Init()` method.
- [x] Write tests.
- [x] Run `go test`.

## Context Log

Files read: `internal/data/task.go`, `internal/board/board.go`, `go.mod`, `E03-board-tui-core/Design.md`
Estimated input tokens: ~900
Notes: `go test ./internal/board/ -v` — 4 tests pass. `go build ./...` clean. Model groups tasks by column on construction; empty column defaults to planned.
