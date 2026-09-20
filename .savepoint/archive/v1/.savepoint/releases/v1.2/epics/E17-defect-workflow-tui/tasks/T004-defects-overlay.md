---
id: E17-defect-workflow-tui/T004-defects-overlay
status: done
objective: Add a keyboard-driven Defects overlay for browsing release defects by status
depends_on: [E17-defect-workflow-tui/T001-defect-data-model, E17-defect-workflow-tui/T003-board-defect-summary]
complexity_tier: high
complexity_reason: "Full new keyboard-driven overlay with grouping, navigation, empty state, and help text integration"
---

# T004: Defects Overlay

## Problem

Users need to inspect release defects without leaving the board or navigating the filesystem. The overlay should behave like existing release and epic selectors so the TUI stays consistent.

## Context Files

- `internal/board/model.go` - overlay and navigation state definitions
- `internal/board/update.go` - keyboard handling and overlay transitions
- `internal/board/view.go` - overlay composition
- `internal/board/release.go` - release selector overlay patterns
- `internal/board/epic_panel.go` - epic selector and sidebar list rendering patterns
- `internal/board/help.go` - keyboard shortcut documentation
- `internal/board/update_test.go` - keyboard and overlay tests
- `internal/board/help_test.go` - help overlay tests

## Acceptance Criteria

- [ ] Pressing `d` opens a Defects overlay
- [ ] The overlay lists defects grouped or filtered by open, in-progress, and resolved state
- [ ] Defect rows show severity, defect id, title, and linked epic/task when available
- [ ] Arrow and vim navigation move through defects consistently with existing overlays
- [ ] `esc` closes the overlay without changing board selection
- [ ] Help text documents the `d` shortcut
- [ ] Empty defect lists render a clear empty state
- [ ] Tests cover shortcut handling, navigation, empty state, and overlay rendering
- [ ] `make build && make test` passes

## Implementation Plan

- [ ] Add a defect overlay type and defect cursor state
- [ ] Implement defect list rendering using existing selector visual patterns
- [ ] Wire `d`, navigation, enter, and escape behavior into update handling
- [ ] Update help overlay shortcut text
- [ ] Add focused rendering and update tests
- [ ] Run `make build && make test`

## Context Log

- Files read:
- Files edited:
- Token estimate:
- Quality gates:

