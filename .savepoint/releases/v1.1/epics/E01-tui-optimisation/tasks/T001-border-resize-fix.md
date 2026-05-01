---
id: E01-tui-optimisation/T001-border-resize-fix
status: planned
objective: "Fix right-border clipping and ensure clean rendering on terminal resize"
depends_on: []
---

# T001: Fix Right-Border Clipping and Resize Robustness

## Acceptance Criteria

- The board frame right border (`│` or rounded corner) is always visible at any terminal width ≥ 40 chars
- Reducing terminal width below a breakpoint (120→119, 80→79) does not leave stray pixels or broken border artifacts
- Reducing to very narrow widths (< 50) degrades gracefully (no visual corruption)
- Expanding terminal width back renders cleanly with no leftover characters from previous dimensions
- All existing layout breakpoints (120/80) still function correctly

## Implementation Plan

- [ ] Edit `internal/board/layout.go` — audit width arithmetic to ensure total content width exactly fills `width - boardFrameOverhead` at every breakpoint; use `totalWidth := cw*colCount + colOverhead*colCount` and pad any remainder so it matches `inner` exactly.
- [ ] Edit `internal/board/view.go` — add minimum terminal width clamping (e.g., floor to 40) in `View()` to prevent degenerate states.
- [ ] Edit `internal/board/layout.go` — add a `minBoardWidth` const and guard `CalculateLayout` inputs below it.
- [ ] Run `make build && make test` to verify no regressions.

## Context Log

Files read:
- `internal/board/layout.go`
- `internal/board/view.go`
- `internal/board/update.go`

Estimated input tokens: 600

Notes:
- Moved from `E06-atari-noir-layout/T006-border-resize-fix` (release v1) to `E01-tui-optimisation` (release v1.1).
