---
id: I014
title: Done column and task card focus accent should be green not orange
type: defect
status: resolved
source:
  kind: report
  actor:
    role: owner
    session: user-review
  at: '2026-09-20T21:36:36Z'
severity: low
resolution:
  disposition: accepted
  actor:
    role: owner
    session: user-request
  at: '2026-09-22T09:18:44Z'
  reason: >-
    Owner waived an independent Check and accepted I014 as resolved based on
    the repair evidence recorded by the implementing task; this is not a
    technical CLEAR verdict.
history:
  - at: '2026-09-20T21:36:36Z'
    actor:
      role: owner
      session: user-review
    kind: observed
    note: User requested the Done column and focused task card border in Done use green accent rather than orange
  - at: '2026-09-20T21:46:00Z'
    actor:
      role: executor
      session: user-request
    kind: repair_attempted
    note: Added ColumnFocusedDone, ColumnTitleFocusedDone, CardBoxFocusedDone, and TaskItemFocusedDone with green accent; updated column and card rendering and tests
  - at: '2026-09-22T09:18:44Z'
    actor:
      role: owner
      session: user-request
    kind: owner_decision
    note: >-
      Owner waived independent verification and accepted I014 as resolved based
      on the repair evidence from the implementing task.
---

# I014: Done column and task card focus accent should be green not orange

## Summary

In `internal/board/v2/column.go` and `internal/board/v2/card.go`, focusing the `DONE` column or a task card in the `DONE` column currently uses the default orange focus accent (`styles.ColumnFocused`, `styles.ColumnTitleFocused`, `styles.CardBoxFocused`).

The Done column and focused task cards in the Done column should instead use green (`clrGreen` / `NPPGreen` `#A4C639`) to align with the visual semantics of completed work.

## Evidence

- `internal/board/v2/column.go:columnStyle`:
  Applies `styles.ColumnFocused` (orange border) when focused for non-planned columns including `DONE`.
- `internal/board/v2/column.go:renderColumn`:
  Applies `styles.ColumnTitleFocused` (orange title) when focused for non-planned columns including `DONE`.
- `internal/board/v2/card.go:cardStyle`:
  Applies `styles.CardBoxFocused` (orange border) when a done card is focused.
- `internal/styles/styles.go`:
  Defines `TagDone = lipgloss.NewStyle().Foreground(clrGreen)` and `BadgeClear = lipgloss.NewStyle().Foreground(clrGreen)`, establishing green as the canonical completion accent in Atari-Noir palette.

## Proof Needed

1. The Done column when focused renders with a green border and header (`ColumnFocusedDone`, `ColumnTitleFocusedDone`).
2. Focused task cards with `status == data.ColumnDone` render with a green border (`CardBoxFocusedDone`) and title (`TaskItemFocusedDone`).
3. Unit tests in `internal/board/v2` assert green focus styling for the Done column and its cards.

## Repair Evidence

1. `internal/styles/styles.go`:
   - Defined `ColumnFocusedDone`, `ColumnTitleFocusedDone`, `CardBoxFocusedDone`, and `TaskItemFocusedDone` using `clrGreen`.
2. `internal/board/v2/column.go`:
   - Added `isDoneColumn` helper and wired `ColumnFocusedDone` / `ColumnTitleFocusedDone` when `isDoneColumn` is focused.
3. `internal/board/v2/card.go`:
   - Updated `cardStyle` to return `CardBoxFocusedDone` when `status == data.ColumnDone && focused`.
   - Updated `titleStyle` to return `TaskItemFocusedDone` when `status == data.ColumnDone && focused`.
4. Tests in `internal/board/v2/column_test.go` (`TestRenderColumnPlannedFocusUsesPlannedAccent`) and `internal/board/v2/card_test.go` (`TestRenderCardDoneFocusUsesGreenStyling`) verify the green ANSI sequences (`163;198;56`) and absence of orange.
