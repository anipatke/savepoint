---
id: E04-board-components/T006-help-overlay
status: done
objective: "Render help overlay with keyboard shortcuts"
depends_on: [E04-board-components/T005-release-dropdown]
---

# T006: Help Overlay

## Acceptance Criteria

- `?` opens help overlay.
- Shows all keyboard shortcuts.
- `Esc` or `q` closes help.
- Overlay renders centered.

## Implementation Plan

- [x] Create `internal/board/help.go`.
- [x] Implement help content.
- [x] Wire `?` key into Update loop.
- [x] Write tests.
- [x] Run `go test`.

## Context Log

Files read: `.savepoint/router.md`, `.savepoint/releases/v1/epics/E04-board-components/E04-Detail.md`, this task file, `agent-skills/ink-tui-design/SKILL.md`, `.savepoint/visual-identity.md`, `internal/board/model.go`, `internal/board/update.go`, `internal/board/view.go`, `internal/board/detail.go`, `internal/board/epic_panel.go`, `internal/board/release.go`, `internal/board/update_test.go`, `internal/board/view_test.go`, `internal/board/detail_test.go`, `internal/board/release_test.go`, `internal/styles/styles.go`, `go.mod`, `Makefile`
Estimated input tokens: ~8500
Notes: Added help overlay renderer, `?` key handling, centered overlay rendering over dimmed board, and focused render/update/view tests. Focused test `go test ./internal/board/...` initially failed in the sandbox because Go could not access the external build cache; rerun with approved cache access passed. Quality gates: `go build ./...` pass, `go vet ./...` pass, `go test ./...` pass.

## Drift Notes

- Drift: `internal/board/help.go` added, not yet in Codebase Map.
- Drift: `internal/board/help_test.go` added, not yet in Codebase Map.
