---
id: E03-board-tui-core/T002-update-loop
status: done
objective: "Implement Update(msg) with key and window size handling"
depends_on: [E03-board-tui-core/T001-model]
---

# T002: Update Loop

## Acceptance Criteria

- `Update(msg)` handles `tea.KeyMsg` for `q`, `ctrl+c`.
- `Update(msg)` handles `tea.WindowSizeMsg` for terminal resize.
- Returns `tea.Quit` on `q`.
- Tests verify update behavior.

## Implementation Plan

- [x] Create `internal/board/update.go`.
- [x] Implement `Update(msg tea.Msg)` with type switch.
- [x] Add key handling.
- [x] Add window size handling.
- [x] Write tests.
- [x] Run `go test`.

## Context Log

Files read: `internal/board/model.go`, `T001-model.md`, `E03-board-tui-core/Design.md`
Estimated input tokens: ~600
Notes: `go test ./internal/board/ -v` — 8 pass. `go build ./...` clean. Update returns `(Model, tea.Cmd)` not `(tea.Model, tea.Cmd)` — keeps type concrete for T003 View composition.
