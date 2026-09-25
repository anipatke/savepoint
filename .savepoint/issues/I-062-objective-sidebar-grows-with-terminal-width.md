---
id: I-062
title: Let the Objective sidebar grow with terminal width
type: other
status: open
source:
  kind: report
  actor: {role: owner, session: user}
  at: '2026-09-25T21:34:02Z'
severity: low
history:
  - at: '2026-09-25T21:34:02Z'
    actor: {role: owner, session: user}
    kind: observed
    note: >-
      The Objectives sidebar still feels cramped after I-021. Owner chose to
      let the sidebar take a share of the terminal width on wide terminals
      instead of staying fixed at 34 cells.
  - at: '2026-09-25T21:38:22Z'
    actor: {role: executor, session: i062-repair-20260926}
    kind: repair_attempted
    note: >-
      internal/board/v2/view.go replaces the fixed sidebarWidth constant with
      sidebarWidth(termW): 30% of the content width, capped so the Task
      columns keep their 30-cell breakpoint width, clamped to 34..52 cells.
      renderColumns and columnWidth share it. Widths: 124 -> 34, 126 -> 36,
      160 -> 48, 200+ -> 52; below 124 the sidebar collapses as before.
      TestSidebarGrowsWithTerminalWidth covers breakpoint, column-room limit,
      intermediate, cap, and past-cap widths, each asserting columns >= 30
      and the rendered board fits. Sidebar tests now measure against the
      model's derived width. Design.md states no sidebar geometry, so it is
      unchanged. go test ./internal/board/..., gofmt, git diff --check, and
      make build && make test-fast passed. Not verified by eye in a live
      terminal or with NO_COLOR at the new widths. Issue remains open for
      independent verification.
---

# I-062: Let the Objective sidebar grow with terminal width

## Summary

The Objective sidebar is fixed at 34 cells at every terminal width. After the
frame (4 cells), the `▸●` marker column (3), and the `O-### ` prefix (6), a
title has about 21 cells per line and truncates after two wrapped lines, so
anything past roughly 42 characters is cut. Most live Objective titles are
45–69 characters. On wide terminals the spare width goes entirely to the three
Task columns, which need only 30 cells each.

I-021 widened the sidebar from 28 to 34 cells as a fixed value. This Issue
changes the rule instead: above the sidebar breakpoint, the sidebar takes about
30% of the terminal width, clamped between today's 34 cells and about 52 cells,
and the Task columns split the rest.

| Terminal | Sidebar now | Sidebar after | Title cells/line |
|---|---|---|---|
| 124 (breakpoint) | 34 | 34 | 21 → 21 |
| 160 | 34 | ~48 | 21 → ~35 |
| 200+ | 34 | ~52 | 21 → ~39 |

## Evidence

- `internal/board/v2/view.go:22` fixes `sidebarWidth = 34`.
- `internal/board/v2/view.go:325` passes that constant to `renderSidebar`, and
  `internal/board/v2/view.go:347-360` (`columnWidth`) subtracts it before
  splitting the remaining width three ways.
- `internal/board/v2/objectives.go:329-343` wraps `O-### Title` to at most two
  lines of `width - columnChrome - rowMarkerCells`.
- Owner report, 2026-09-26: "Objectives feels cramped". Owner selected the
  width-share option over the in-place row tweaks or a toggle key.

## Proof Needed

- The sidebar width is derived from the terminal width: exactly 34 cells at the
  124-cell breakpoint, growing with width, capped near 52 cells.
- `columnWidth` and `renderSidebar` use the same derived width, so the sidebar
  plus three columns never exceed the terminal width and no column drops below
  its existing minimum at any visible-sidebar width.
- Below the breakpoint the sidebar still collapses exactly as today, and the
  compact single-column layout is unchanged.
- Titles, badges, selection/focus markers, priority headings, and scrolling
  render correctly at the breakpoint, at an intermediate width, and at and
  above the cap, with and without colour.
- `.savepoint/Design.md` describes the width rule if it states sidebar
  geometry.
- Width tests cover the breakpoint, an intermediate width, and the cap; `git
  diff --check`, `make build`, and `make test-fast` pass.
