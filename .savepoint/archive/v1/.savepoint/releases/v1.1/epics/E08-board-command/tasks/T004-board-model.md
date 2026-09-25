---
id: E08-board-command/T004-board-model
status: done
objective: "Implement TUI board model with 3-column Kanban"
depends_on: ["E08-board-command/T003-tui-app-shell"]
---

# T004: Board Model

## Acceptance Criteria

- Three columns: planned, in_progress, done
- Column headers with task count
- Task cards showing ID and objective
- Keyboard navigation: arrows, vim-style j/k/h/l
- Selection highlighting with accent color
- Filter by release and epic

## Implementation Plan

- [x] Add `internal/board/model.go` (extend existing or create new)
- [x] Implement column structure (planned, in_progress, done)
- [x] Implement card rendering (ID + objective truncated)
- [x] Implement selection state (current column, current card)
- [x] Handle key events: Up/Down/Left/Right, j/k/h/l
- [x] Apply theme colors for focus/accent
- [x] Implement release/epic filtering from args
- [x] Test board rendering in TUI
- [x] Run `make build && make test`