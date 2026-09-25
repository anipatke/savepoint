---
id: I-015
title: Allow up to two lines for task card and objective card titles
type: defect
status: resolved
source:
  kind: report
  actor:
    role: owner
    session: user-review
  at: '2026-09-20T21:38:00Z'
severity: low
resolution:
  disposition: accepted
  actor:
    role: owner
    session: user-request
  at: '2026-09-22T09:18:44Z'
  reason: >-
    Owner waived an independent Check and accepted I-015 as resolved based on
    the repair evidence recorded by the implementing task; this is not a
    technical CLEAR verdict.
history:
  - at: '2026-09-20T21:38:00Z'
    actor:
      role: owner
      session: user-review
    kind: observed
    note: User requested allowing a line break for task title (up to two lines) and objective card title so full titles can be read
  - at: '2026-09-20T21:49:00Z'
    actor:
      role: executor
      session: user-request
    kind: repair_attempted
    note: Implemented wrapTitleLines, updated renderCard and renderObjectiveRow to wrap titles up to two lines before truncating, and added tests
  - at: '2026-09-22T09:18:44Z'
    actor:
      role: owner
      session: user-request
    kind: owner_decision
    note: >-
      Owner waived independent verification and accepted I-015 as resolved based
      on the repair evidence from the implementing task.
---

# I-015: Allow up to two lines for task card and objective card titles

## Summary

In Savepoint V2's board, task titles on task cards (`internal/board/v2/card.go`) and objective titles in the objective sidebar (`internal/board/v2/objectives.go`) are currently truncated strictly to a single line using `truncateCells`. Long titles are cut off with an ellipsis (`…`), making them difficult to read without opening detail overlays.

Task card titles and objective card titles should support wrapping across up to two lines before truncating.

## Evidence

- `internal/board/v2/card.go:renderCard`:
  ```go
  title := truncateCells(card.Task.Title, textW)
  ```
  Only a single line is rendered for the task title.
- `internal/board/v2/objectives.go:renderObjectiveRow`:
  ```go
  titleW := width - rowMarkerCells - 1
  ...
  title := truncateCells(titleContent, titleW)
  ```
  The objective title row is strictly truncated to a single line.

## Proof Needed

1. Long task card titles wrap naturally up to two lines within `textW` width, truncating with ellipsis (`…`) on the second line only if exceeding two lines.
2. Long objective card titles wrap naturally up to two lines within the title area, properly aligning/indenting the second line, and truncating with ellipsis on the second line only if exceeding two lines.
3. Short titles that fit in one line remain single-line (no extra blank lines added).
4. Unit tests in `internal/board/v2` verify two-line wrapping and truncation for both task cards and objective rows.

## Repair Evidence

1. `internal/board/v2/width.go`:
   - Added `wrapTitleLines(text string, width int, maxLines int) []string` to wrap titles across up to `maxLines` lines at word boundaries, truncating the final line with an ellipsis (`…`) only when exceeding `maxLines`.
2. `internal/board/v2/card.go`:
   - Updated `renderCard` to format `card.Task.Title` with `wrapTitleLines(card.Task.Title, textW, 2)`, rendering each title line with the active `titleStyle`.
3. `internal/board/v2/objectives.go`:
   - Updated `renderObjectiveRow` to format the objective title with `wrapTitleLines(row.ID()+" "+row.Objective.Title, textW, 2)`, rendering line 1 with row markers and line 2 indented past the marker cells.
4. Unit tests:
   - `internal/board/v2/width_test.go`: `TestWrapTitleLines` verifies single-line fitting, two-line clean wrapping, two-line truncation with ellipsis, and oversized words.
   - `internal/board/v2/card_test.go`: `TestRenderCardTitleWrapsUpToTwoLines` verifies 1-line, 2-line wrapped, and 2-line truncated card behavior and dynamic height adjustments.
   - `internal/board/v2/objectives_test.go`: `TestRenderObjectiveRow_WrapsUpToTwoLinesAndTruncates` verifies short, two-line wrapped, and two-line truncated objective row behavior.
