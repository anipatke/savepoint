---
id: E17-defect-workflow-tui/T003-board-defect-summary
status: done
objective: Load defects into board state and show a compact open-defect signal without changing the three-column board layout
depends_on: [E17-defect-workflow-tui/T001-defect-data-model]
---

# T003: Board Defect Summary Signal

## Problem

Defects need to be visible from the primary board, but a fourth column would crowd the layout and break the existing terminal-width model. The board should surface defect pressure while preserving the planned, in-progress, and done columns.

## Context Files

- `internal/board/model.go` - board state grouping and selected release/epic fields
- `internal/board/board.go` - board startup data loading
- `internal/board/interfaces.go` - board data-access boundary
- `internal/board/view.go` - header, Next Activity, footer, and board rendering
- `internal/board/layout.go` - current three-column layout breakpoints
- `internal/board/board_test.go` - board startup and data-loading tests
- `internal/board/layout_test.go` - layout breakpoint tests
- `internal/styles/styles.go` - shared visual styles and semantic colors

## Acceptance Criteria

- [ ] Board data loading includes release-scoped defects
- [ ] The board header or Next Activity area shows open defect count when defects exist
- [ ] Zero open defects do not add noisy visual chrome
- [ ] The main board remains three task columns at existing breakpoints
- [ ] Defect summary respects selected release
- [ ] Existing epic and task filtering behavior remains unchanged
- [ ] Tests cover zero defects, open defects, and unchanged three-column layout
- [ ] `make build && make test` passes

## Implementation Plan

- [ ] Add defect storage to board model state
- [ ] Extend board data loading through the existing dependency boundary
- [ ] Add a compact defect summary renderer
- [ ] Place the summary in existing header/Next Activity chrome without reducing column width
- [ ] Add board rendering and layout regression tests
- [ ] Run `make build && make test`

## Context Log

- Files read:
- Files edited:
- Token estimate:
- Quality gates:

