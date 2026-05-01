---
id: E04-board-components/T005-release-dropdown
status: done
objective: "Render release selector dropdown"
depends_on: [E04-board-components/T004-detail-overlay]
---

# T005: Release Dropdown

## Acceptance Criteria

- `r` opens release dropdown.
- Shows list of releases with current one marked.
- Selecting a release updates board.
- `Esc` closes without selecting.

## Implementation Plan

- [x] Add `OverlayRelease` to `model.go` overlay type constants.
- [x] Add `Releases []string` and `ReleaseCursor int` fields to `Model`.
- [x] Create `internal/board/release.go` with `RenderReleaseDropdown` and `releaseIndex`.
- [x] Wire `r` key in `update.go` to open `OverlayRelease` and set cursor.
- [x] Handle `up`/`down`/`enter`/`esc` for release overlay in `updateOverlay`.
- [x] Render `OverlayRelease` in `view.go`.
- [x] Write `release_test.go` covering render, nav, select, and esc.
- [x] Run `go test ./internal/board/...` — all pass.
- [x] Run `go build ./... && go vet ./... && go test ./...` — all pass.

## Context Log

Files read: model.go, update.go, view.go, epic_panel.go, board.go, update_test.go, model_test.go
Estimated input tokens: ~2500
Notes: Mirrored epic dropdown pattern exactly. `r` always opens regardless of screen width (unlike `e` which is narrow-only). Quality gates: build OK, vet OK, all tests pass.
