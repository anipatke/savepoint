---
type: epic-design
status: audited
---

# E08: Board Command

## Purpose

Implement `savepoint board`, a Bubble Tea terminal UI for viewing and moving tasks through the Savepoint workflow, including audit-review mode entry points.

## Interface

```bash
savepoint board              # launch TUI
savepoint                    # alias for board
savepoint board --release v1 # filter by release
savepoint board --epic E03   # filter by epic
```

## What this epic adds

- Ink/Bubble Tea app shell launched by `savepoint board`
- Three-column Kanban view for planned, in_progress, done
- Detail pane for selected task
- Keyboard navigation (arrows, vim-style j/k/h/l)
- Status transition actions with gate enforcement
- Manual refresh with `r` key
- Theme tokens loaded from config.yml
- Render fallbacks: 24-bit → 256-color → 16-color → NO_COLOR → non-TTY plain table
- Optimistic mtime conflict detection before writes
- Audit entry point: signal if proposals exist, exit toward audit

## Components

| Module | Purpose |
|--------|---------|
| `cmd/board.go` | CLI registration, arg parsing |
| `internal/board/cmd.go` | Non-TTY fallback, dispatch |
| `internal/board/tui.go` | Bubble Tea program setup |
| `internal/board/model.go` | Board state, columns, cards |
| `internal/board/detail.go` | Task detail pane |
| `internal/board/transitions.go` | Status transition gate enforcement |
| `internal/board/theme.go` | Theme loading, color fallbacks |
| `internal/board/plain.go` | Non-TTY plain text table |

## Implemented As

- `main.go` dispatches `board` and uses the no-arg command path as the board default.
- `cmd/board.go` owns board-specific flag parsing for `--release`, `--epic`, and `--help`.
- `internal/board/board.go`, `tui.go`, `model.go`, `plain.go`, `transitions.go`, and related render/update files own project discovery, non-TTY fallback, TUI startup, filtering, detail overlays, keyboard flow, and task status writes.
- Audit entry-point work is represented as the non-TTY audit proposal signal and the epic Detail/Audit tab flow added around E06.

## Boundaries

**In scope:**
- Board navigation and status transitions
- Gate enforcement for dependencies and status rules
- Plain table fallback when not in a TTY
- Manual `r` refresh
- Initial audit-review mode entry point if proposals exist

**Out of scope:**
- File watching
- Search
- Mouse interaction
- Drag-and-drop
- Full audit proposal diff implementation
