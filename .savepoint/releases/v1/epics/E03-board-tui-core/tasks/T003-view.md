---
id: E03-board-tui-core/T003-view
status: done
objective: "Render a bordered 3-column grid via View()"
depends_on: [E03-board-tui-core/T002-update-loop]
---

# T003: View Rendering

## Acceptance Criteria

- `View()` renders a header, 3-column board, and status bar.
- Columns have borders using Lip Gloss.
- Board fills terminal width.
- No wrapping or layout breaks.

## Implementation Plan

- [x] Create `internal/board/view.go`.
- [x] Implement `View()` composing header, board, footer.
- [x] Use `lipgloss.JoinHorizontal` and `JoinVertical`.
- [x] Verify rendering at multiple terminal sizes.

## Context Log

Files read: `internal/board/model.go`, `internal/board/update.go`, `.savepoint/visual-identity.md`, `internal/data/task.go`, `E03-board-tui-core/Design.md`, `T003-view.md`
Estimated input tokens: ~2000
Notes: `go test ./internal/board/ ./internal/styles/ -v` — 17 pass. `go build ./...` clean. Created `internal/styles/palette.go` (Atari-Noir hex constants) and `internal/styles/styles.go` (Lip Gloss style vars). `view.go` uses `lipgloss.JoinHorizontal` for 3 columns, `JoinVertical` for header+board+status. Column content width = (termW - 12) / 3; fallback default=80 when Width=0. Focused column gets orange border; focused task gets orange text.

## Drift Notes

- Drift: `internal/styles/palette.go` added, not yet in Codebase Map.
- Drift: `internal/styles/styles.go` added, not yet in Codebase Map.
- Drift: `internal/board/view.go` added, not yet in Codebase Map.
