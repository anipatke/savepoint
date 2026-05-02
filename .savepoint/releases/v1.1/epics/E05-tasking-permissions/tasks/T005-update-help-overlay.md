---
id: E05-tasking-permissions/T005-update-help-overlay
status: planned
objective: "Add `m` key binding to the TUI help overlay"
depends_on: ["E05-tasking-permissions/T004-implement-m-hotkey"]
---

# T005: Update Help Overlay

## Acceptance Criteria

- Help overlay (`?` key) shows `m` key binding: "m: update router"
- Key binding appears with consistent styling in the shortcuts list
- No other help entries modified

## Implementation Plan

- [ ] Add `helpRow("m", "update router")` to `RenderHelp` in `help.go`
- [ ] Run `make build && make test` to verify

## Context Log

Files read:
- .savepoint/releases/v1.1/epics/E05-tasking-permissions/E05-Detail.md
- .savepoint/releases/v1.1/epics/E05-tasking-permissions/tasks/T005-update-help-overlay.md
- internal/board/help.go

Estimated input tokens: 2000

Notes:
