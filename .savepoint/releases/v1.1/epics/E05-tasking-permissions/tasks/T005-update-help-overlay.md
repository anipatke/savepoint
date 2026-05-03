---
id: E05-tasking-permissions/T005-update-help-overlay
status: done
objective: "Add `p` priority key binding to the TUI help overlay"
depends_on: ["E05-tasking-permissions/T004-implement-m-hotkey"]
---

# T005: Update Help Overlay

## Acceptance Criteria

- Help overlay (`?` key) shows `p` key binding for router priority
- Key binding appears with consistent styling in the shortcuts list
- No unrelated help entries modified

## Implementation Plan

- [x] Add the `p` priority shortcut to `RenderHelp` in `help.go`
- [x] Verify with `go build -o savepoint.exe main.go` and `go test ./...`

## Context Log

Files read:
- .savepoint/releases/v1.1/epics/E05-tasking-permissions/E05-Detail.md
- .savepoint/releases/v1.1/epics/E05-tasking-permissions/tasks/T005-update-help-overlay.md
- internal/board/help.go

Estimated input tokens: 2000

Notes:
- Implementation uses `p` rather than the originally planned `m` key.
- `go build -o savepoint.exe main.go` passed during E05 audit.
- `go test ./...` passed during E05 audit.
