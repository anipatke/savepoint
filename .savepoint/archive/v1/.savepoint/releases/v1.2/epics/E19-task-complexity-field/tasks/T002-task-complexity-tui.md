---
id: E19-task-complexity-field/T002-task-complexity-tui
status: done
objective: Render task complexity on cards and in task detail
depends_on: [E19-task-complexity-field/T001-task-complexity-model]
complexity_tier: medium
complexity_reason: "Updates card and detail rendering with layout-sensitive tests"
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

- [x] Task cards show the complexity tier without breaking the existing ID/title/phase layout
- [x] The task detail overlay shows a `Complexity:` row with the tier and reason
- [x] Long complexity reasons wrap cleanly and remain readable on narrow widths
- [x] Existing card and detail fields continue to render as before
- [x] Tests cover the card, detail, and narrow-width rendering cases
- [x] `make build && make test` passes for the updated board package

## Implementation Plan

- [x] Update card rendering to include a compact complexity signal
- [x] Update task detail rendering to show the full complexity reason
- [x] Preserve existing layout and wrapping behavior on narrow terminals
- [x] Add focused board rendering tests for the new field
- [x] Run `make build && make test`

## Context Log

- Files read: card.go, detail.go, card_test.go, detail_test.go, render_policy_test.go, task.go, styles.go
- Notes: Card gets `· low/med/high/spike` 3rd line (CardMeta style). Detail gets `Complexity: tier [— reason]` row via existing detailRow (handles wrapping). Pre-existing root test failure (savepoint-audit SKILL.md missing frontmatter) unrelated to board changes.
