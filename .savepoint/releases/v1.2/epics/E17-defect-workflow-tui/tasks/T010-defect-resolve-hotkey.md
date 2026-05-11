---
type: task
id: E17-defect-workflow-tui/T010-defect-resolve-hotkey
status: done
objective: Allow a user to resolve an open defect directly from the Defects overlay by pressing `space` on the focused defect row.
depends_on: [T004-defects-overlay, T005-defect-detail-and-markers]
---

# T010: Defect Resolve Hotkey

## Objective

Allow a user to resolve an open defect directly from the Defects overlay by pressing `space` on the focused defect row.

## Acceptance Criteria

- [x] In the Defects overlay, pressing `space` on a selected defect with `status: open` updates that defect file to `status: resolved`.
- [x] The status update preserves the defect body and unrelated frontmatter fields.
- [x] The resolved defect no longer appears in open-defect counts after the board reloads or refreshes its defect state.
- [x] Pressing `space` on a selected defect that is already `resolved` does not rewrite the defect file or produce an error state.
- [x] Pressing `space` on a selected defect with `status: in_progress` does not silently skip required lifecycle handling; it must either leave the file unchanged with a clear no-op path or implement a valid transition that removes `stage`.
- [x] Keyboard help or footer copy documents the `space` resolve shortcut where defect overlay shortcuts are shown.
- [x] Focus and overlay mode remain stable after the update; the user stays in the Defects overlay with a valid selected row.
- [x] `make build` passes.
- [x] `make test` passes.

## Implementation Plan

- [x] Add focused-defect status transition handling to the board update path for the Defects overlay `space` key.
- [x] Use existing data write/frontmatter helpers rather than hand-editing markdown strings in board code.
- [x] Refresh in-memory defect state after a successful write so counts, rows, and related markers reflect the resolved status.
- [x] Add or update board tests covering open-to-resolved, resolved no-op, in-progress lifecycle behavior, and focus stability.
- [x] Add or update data write tests if existing helpers do not already cover preserving defect markdown while changing status.
- [x] Update defect overlay shortcut copy in help/footer rendering and its tests.
- [x] Run `make build` and `make test`.

## Context Files

- `internal/board/update.go`
- `internal/board/model.go`
- `internal/board/io.go`
- `internal/board/defect_overlay.go`
- `internal/board/defect_overlay_test.go`
- `internal/board/help.go`
- `internal/board/help_test.go`
- `internal/data/defect.go`
- `internal/data/write.go`
- `internal/data/write_test.go`

## Context Log

Files read:
- `.savepoint/router.md`
- `agent-skills/savepoint-create-defect/SKILL.md`
- `agent-skills/savepoint-build-task/SKILL.md`
- `.savepoint/releases/v1.2/epics/E17-defect-workflow-tui/tasks/T010-defect-resolve-hotkey.md`
- `internal/board/update.go`
- `internal/board/model.go`
- `internal/board/io.go`
- `internal/board/defect_overlay.go`
- `internal/board/defect_overlay_test.go`
- `internal/board/help.go`
- `internal/board/help_test.go`
- `internal/data/defect.go`
- `internal/data/write.go`
- `internal/data/write_test.go`
- `.savepoint/releases/v1.2/defects/D002-closed-e17-placeholder.md`
- `internal/data/lifecycle.go`
- `internal/data/parser.go`
- `internal/board/watch.go`
- `internal/board/board.go`
- `internal/board/interfaces.go`

Estimated input tokens: ~47k

Notes:
- Implemented defect resolve as canonical `open` -> `resolved` while preserving the UI label "RESOLVED".
- Normalized v1.2 defect fixtures to the defect-specific `open`/`resolved` lifecycle so board parsing no longer panics.
- Verification: `go test ./internal/data ./internal/board`, `go run main.go`, `make build`, and `make test` passed.
