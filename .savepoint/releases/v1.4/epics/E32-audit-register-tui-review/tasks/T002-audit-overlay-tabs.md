---
id: E32-audit-register-tui-review/T002-audit-overlay-tabs
status: planned
objective: Add the read-only Audit Register overlay with Prompt, Findings, and Runs tabs.
depends_on:
  - E32-audit-register-tui-review/T001-board-audit-data-loading
complexity_tier: medium
complexity_reason: The overlay reuses existing patterns but adds a new top-level interaction path.
---

# T002: Audit Overlay Tabs

## Problem

Users need a dedicated TUI section for audit prompt, current findings, and run history instead of hunting through markdown files.

## Context Files

- `internal/board/model.go`
- `internal/board/update.go`
- `internal/board/view.go`
- `internal/board/audit_overlay.go`
- `internal/board/audit_overlay_test.go`
- `internal/board/epic_panel.go`
- `internal/board/release_docs_overlay_test.go`

## Acceptance Criteria

- [ ] Pressing `A` opens the Audit Register overlay.
- [ ] Overlay tabs are Prompt, Findings, and Runs.
- [ ] `[`/`]`, left/right, and `h`/`l` switch tabs consistently with existing document overlays.
- [ ] `up`/`down`, `k`/`j`, `pgup`, and `pgdown` scroll the selected tab body.
- [ ] `esc` and `q` close the overlay.
- [ ] Missing prompt, findings, or runs render useful empty states.

## Implementation Plan

- [ ] Add audit overlay constants and navigation state.
- [ ] Add top-level `A` key handling.
- [ ] Render the tab strip and selected body using existing overlay styles.
- [ ] Add scroll offset handling per audit tab.
- [ ] Add tests for tab switching, scrolling, empty states, and close behavior.

## Context Log

Pending.
