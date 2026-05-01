---
id: E06-atari-noir-layout/T005-restore-nav-hints
status: done
objective: "Restore persistent navigation keybinding hints to the main board footer"
depends_on: []
---

# T005: Restore Navigation Hints in Footer

## Acceptance Criteria

- The footer shows a second line below the phase labels with keybinding hints: `←/→:nav  E:epic  R:release  ?:help  q:quit`
- Hints are styled in a subdued color to not overwhelm the phase labels
- The footer still fits within a single terminal width without wrapping
- All existing footer behavior (centered phase labels, alignment) remains intact
- Hints refresh correctly on terminal resize

## Implementation Plan

- [x] Edit `internal/board/view.go` — update `renderFooter()` to append a second line of keybinding hints below the phase labels.
- [x] Add a new `FooterHints` style in `internal/styles/styles.go` using `clrDim` (dim/subdued appearance).
- [x] Run `make build && make test` to verify no regressions.

## Context Log

Files read:
- `internal/board/view.go`
- `internal/board/update.go`
- `internal/styles/styles.go`
- `internal/board/view_test.go`
- `internal/board/layout.go`
- `internal/board/model.go`

Estimated input tokens: 1400

Notes:
- `make build` could not run because `make` is not installed on PATH in this environment.
- Verified the task with `go build ./...`, `go test ./...`, and focused `go test ./internal/board`.
- Footer render now uses a width-guarded centered phase line, a blank spacer line, and a second subdued hint line.
