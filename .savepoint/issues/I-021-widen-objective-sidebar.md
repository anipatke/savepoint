---
id: I-021
title: Widen the Objective sidebar by about twenty percent
type: other
status: resolved
source:
  kind: report
  actor: {role: owner, session: user}
  at: '2026-09-22T09:47:13Z'
severity: low
resolution:
  disposition: accepted
  actor: {role: owner, session: user}
  at: '2026-09-22T11:16:02Z'
  reason: Owner explicitly accepted the remaining risk and requested this Issue be marked resolved.
history:
  - at: '2026-09-22T09:47:13Z'
    actor: {role: owner, session: user}
    kind: observed
    note: The Objective column feels too narrow; widen it by approximately twenty percent while preserving the board's responsive layout.
  - at: '2026-09-22T10:45:07Z'
    actor: {role: executor, session: codex}
    kind: repair_attempted
    note: Widened the sidebar to 34 cells, moved its content-width breakpoint to 124 so Task columns remain 30 cells wide at the boundary, added responsive width coverage, and passed the focused V2 tests plus git diff --check, make build, and make test. Independent Check remains required.
  - at: '2026-09-22T11:16:02Z'
    actor: {role: owner, session: user}
    kind: owner_decision
    note: Owner accepted the remaining risk and requested this Issue be marked resolved; no technical CLEAR is implied.
---

# I-021: Widen the Objective sidebar by about twenty percent

## Summary

The V2 board's Objective sidebar is currently fixed at 28 cells. The owner
wants roughly 20% more room so Objective titles and status content are easier
to scan. A target near 34 cells preserves the existing visual proportion while
giving wrapped titles materially more space.

The change must remain responsive: widening the sidebar must not squeeze the
three Task columns below their supported geometry, create overflow, or expose
the sidebar at terminal widths that can no longer carry the full layout.

## Evidence

- `internal/board/v2/view.go:22` fixes `sidebarWidth` at 28 cells.
- `internal/board/v2/view.go:25` uses a 120-cell breakpoint for showing the
  sidebar beside all three Task columns.
- `internal/board/v2/view.go:343-344` subtracts the fixed sidebar width before
  allocating the Task columns.
- The owner reported that the current Objective column is visually too narrow
  and requested an increase of about 20%.

## Proof Needed

- Increase the Objective sidebar from 28 cells to approximately 34 cells, or
  document an adjacent value justified by the board's cell geometry.
- Re-evaluate the sidebar breakpoint so every supported visible-sidebar width
  still fits three readable Task columns without overflow.
- Confirm Objective titles, badges, selection/focus, scrolling, and two-line
  wrapping remain correct with and without colour.
- Confirm narrower terminals still collapse the sidebar cleanly and all board
  lines remain within terminal width.
- Focused Objective/width tests, `git diff --check`, `make build`, and
  `make test` pass.
