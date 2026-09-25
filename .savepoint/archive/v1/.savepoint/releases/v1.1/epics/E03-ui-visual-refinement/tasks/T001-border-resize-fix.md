---
id: E03-ui-visual-refinement/T001-border-resize-fix
status: done
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
- At every layout breakpoint (120+, 80–119, 40–79), total rendered column width + borders + padding equals the available terminal width with zero underflow
- Column widths are calculated as `floor(innerWidth / colCount)` and the remainder is distributed as 1 extra column to the first `remainder` columns
- The board frame overhead (left border + inter-column separators + right border) is exactly accounted for in every breakpoint
- Reducing terminal width by 1 character removes exactly 1 column of content (no more, no less)
- All existing layout tests pass

## Implementation Plan

- [x] Verified — no artifacts observed. Existing layout.go + tests satisfy acceptance criteria.

## Context Log

Files read:
- `internal/board/layout.go`
- `internal/board/view.go`
- `internal/board/update.go`
- `internal/board/layout_test.go`

Estimated input tokens: 800

Notes:
- Consolidated from former T001 + T004 (duplicate width arithmetic scopes).
- Former T004 at `.savepoint/releases/v1.1/epics/E03-ui-visual-refinement/tasks/T004-width-arithmetic-audit.md` has been deleted.
- Moved from `E06-atari-noir-layout/T006-border-resize-fix` (release v1) to `E03-ui-visual-refinement` (release v1.1).
