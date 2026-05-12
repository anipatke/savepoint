---
id: E17-defect-workflow-tui/T002-defect-router-priority
status: done
objective: Extend router state and board priority formatting so a defect can become the active repair item
depends_on: [E17-defect-workflow-tui/T001-defect-data-model]
complexity_tier: medium
complexity_reason: "Extends router with new defect state and wires distinct Next Activity rendering in the board view"
---

# T002: Defect Router Priority

## Problem

The board can currently highlight task-building and audit activity, but defects need their own active repair state. Without router support, a defect can be documented but cannot become the visible next action in the methodology.

## Context Files

- `internal/data/router.go` - router state parsing and serialization
- `internal/data/router_test.go` - router parsing and lifecycle tests
- `internal/board/view.go` - Next Activity phase rendering and compact activity formatting
- `internal/board/card.go` - router-priority matching behavior for task cards
- `internal/board/view_test.go` - Next Activity rendering tests
- `.savepoint/router.md` - live router state format example

## Acceptance Criteria

- [x] Router state supports `state: defect-building`
- [x] Router state supports a `defect` field for the active defect id
- [x] Existing task, epic, and audit router behavior remains unchanged
- [x] Next Activity renders defect work as `DEFECT: ...`
- [x] Compact activity formatting supports defect priority
- [x] Router parsing tolerates older router files with no defect field
- [x] Tests cover defect router parsing, formatting, and unchanged existing states
- [x] `make build && make test` passes

## Implementation Plan

- [x] Extend the router state model with an optional defect id field
- [x] Add `defect-building` to router state parsing and formatting paths
- [x] Render defect priority in the Next Activity line with a distinct warning style
- [x] Keep task-card priority matching scoped to task router states only
- [x] Add router and board view tests for defect activity
- [x] Run `make build && make test`

## Context Log

- Files read: router.go, router_test.go, view.go, card.go, view_test.go, styles.go, palette.go, router.md
- Files edited: router.go, styles.go, view.go, card.go, router_test.go, view_test.go
- Quality gates: make build && make test — all pass

