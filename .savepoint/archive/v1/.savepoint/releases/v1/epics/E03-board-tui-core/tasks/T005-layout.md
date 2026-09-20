---
id: E03-board-tui-core/T005-layout
status: done
objective: "Implement responsive layout with 3 breakpoints"
depends_on: [E03-board-tui-core/T004-styles]
---

# T005: Layout

## Acceptance Criteria

- >=120 cols: epic panel (28w) + 3 columns. ✓ (CalculateLayout, renderEpicPanel)
- 80-119 cols: 3 columns only. ✓
- <80 cols: 1 column, left/right switches columns. ✓ (renderColumns + update.go left/right/h/l keys)
- Column widths calculated from available width. ✓
- No hardcoded widths that break at small sizes. ✓ (minColWidth floor in CalculateLayout)

## Implementation Plan

- [x] Create `internal/board/layout.go`.
- [x] Implement `CalculateLayout(width, height)`.
- [x] Return column count, widths, epic panel visibility.
- [x] Wire into `View()`.
- [x] Test at multiple widths.

## Context Log

Files read: internal/board/model.go, view.go, update.go, board.go, view_test.go, update_test.go, model_test.go, internal/styles/styles.go, go.mod, .savepoint/releases/v1/epics/E03-board-tui-core/E03-Detail.md

Estimated input tokens: ~6 000

Notes:
- Moved colOverhead/minColWidth constants from view.go to layout.go (logical home).
- Audit closeout removed dead columnContentWidth() from view.go; layout width coverage lives in layout_test.go.
- Added left/right (and h/l) column navigation to update.go to satisfy the <80-col acceptance criterion.
- Added EpicPanel lipgloss style to internal/styles/styles.go (purple border, rounded).
- Quality gates: go build ./... ✓ · go vet ./... ✓ · go test ./... ✓ (all 3 packages)
- Audit closeout gates: go test ./internal/board ✓ · go test ./... ✓ · go build ./... ✓

## Drift Notes

- Drift: `internal/board/layout.go` added, not yet in Codebase Map.
