---
id: E17-defect-workflow-tui/T008-footer-defects-shortcut
status: done
objective: Add the defects overlay shortcut to the board footer helper line
depends_on: [E17-defect-workflow-tui/T004-defects-overlay]
---

# T008: Footer Defects Shortcut

## Problem

The board supports `d` for the Defects overlay, but the footer helper line does not advertise the shortcut. Users can miss the defect overlay unless they open help or already know the keybinding.

## Context Files

- `internal/board/view.go` - footer helper line rendering
- `internal/board/help.go` - shortcut wording reference
- `internal/board/view_test.go` - footer rendering test patterns

## Acceptance Criteria

- [x] The board footer helper line includes `d: Defects`
- [x] Existing footer shortcuts remain visible and reasonably compact
- [x] Footer truncation behavior remains stable on narrow terminals
- [x] Tests cover the updated footer helper text
- [x] `make build && make test` passes

## Implementation Plan

- [x] Update the footer helper string in `internal/board/view.go`
- [x] Add or update a focused view/footer test for `d: Defects`
- [x] Run `make build && make test`

## Context Log

- Files read: `internal/board/view.go`, `internal/board/view_test.go`, `internal/board/help.go`, `E17-Detail.md`
- Files edited: `internal/board/view.go`, `internal/board/view_test.go`
- Token estimate: ~7
- Quality gates: `make build && make test` — pass
