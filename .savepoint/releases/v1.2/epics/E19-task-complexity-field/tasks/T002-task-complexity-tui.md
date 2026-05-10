---
id: E19-task-complexity-field/T002-task-complexity-tui
status: planned
objective: Render task complexity on cards and in task detail
depends_on: [E19-task-complexity-field/T001-task-complexity-model]
---

# T002: Task Complexity TUI

## Problem

The board currently shows task identity, status, and phase, but not implementation complexity. Users need the tier visible on the cards they scan and the full reason visible in the task detail overlay.

## Context Files

- `internal/board/card.go`
- `internal/board/detail.go`
- `internal/board/view.go`
- `internal/board/card_test.go`
- `internal/board/detail_test.go`
- `internal/board/render_policy_test.go`

## Acceptance Criteria

- [ ] Task cards show the complexity tier without breaking the existing ID/title/phase layout
- [ ] The task detail overlay shows a `Complexity:` row with the tier and reason
- [ ] Long complexity reasons wrap cleanly and remain readable on narrow widths
- [ ] Existing card and detail fields continue to render as before
- [ ] Tests cover the card, detail, and narrow-width rendering cases
- [ ] `make build && make test` passes for the updated board package

## Implementation Plan

- [ ] Update card rendering to include a compact complexity signal
- [ ] Update task detail rendering to show the full complexity reason
- [ ] Preserve existing layout and wrapping behavior on narrow terminals
- [ ] Add focused board rendering tests for the new field
- [ ] Run `make build && make test`

## Context Log

- Files read:
- Estimated input tokens:
- Notes:
