---
id: E51-first-class-releases/T007-switch-releases-with-the-existing-r-key
status: in_progress
objective: Preserve the existing r-key Release selector while filtering V2 Objectives and Tasks through canonical Release links.
depends_on:
    - E51-first-class-releases/T003-select-one-release-and-project-its-next-action
complexity_tier: high
complexity_reason: Adds a stateful overlay, filtered navigation, persistence, reload behavior, and compatibility coverage to the V2 board.
stage: build
---

# T007: Switch Releases with the existing r key

## Problem

The current board lets a user press `r`, see Releases, move from the current selection, cancel, or press Enter to switch the visible work. E51 must preserve that muscle memory and behavior in V2. Replacing it with a passive label or a different navigation model would regress existing Release functionality.

## Context Files

- `.savepoint/releases/v2/epics/E51-first-class-releases/E51-Detail.md`
- `internal/board/release.go`
- `internal/board/release_test.go`
- `internal/board/update.go`
- `internal/board/help.go`
- `internal/board/v2/model.go`
- `internal/board/v2/load.go`
- `internal/board/v2/load_test.go`
- `internal/board/v2/objectives.go`
- `internal/board/v2/objectives_test.go`
- `internal/board/v2/update.go`
- `internal/board/v2/view.go`
- `internal/board/v2/view_test.go`
- `internal/board/v2/io.go`
- `internal/board/v2/watch.go`
- `internal/board/v2/watch_test.go`
- `internal/board/v2/releases.go`
- `internal/board/v2/releases_test.go`

## Acceptance Criteria

- [ ] On the V2 board, pressing `r` opens a Release selector overlay while leaving the board visible behind it.
- [ ] The selector cursor starts on the currently selected Release, or the first Release when the prior selection is absent.
- [ ] Arrow keys and `j`/`k` move within bounds; Esc or `q` closes without changing Release, router state, focus, bytes, or mtimes.
- [ ] Enter selects the focused Release, closes the overlay, refreshes visible Objectives and Tasks from `ReleaseObjectives`, and resets only invalid cursors/focus.
- [ ] Selection persists optional router `release` context through the canonical async writer and changes no lifecycle/evidence record or `next_action` prose.
- [ ] A pending migration or mtime conflict refuses persistence, reports the reason, and reloads/retains a truthful prior selection without partial state.
- [ ] Selecting a Release never nests Objective or Task ownership and never derives membership from titles, paths, or IDs outside the index.
- [ ] A project with no Releases keeps the shortcut safe and presents the established `(none)` selector state or an equally explicit inert result.
- [ ] Reload after Release/Objective edits preserves selection and focus when records still exist; a removed selection produces the canonical diagnostic.
- [ ] Help text continues to advertise `r` as the Release selector, and the V1 selector tests remain unchanged and passing.

## Implementation Plan

- [x] Model the Release overlay, ordered Release list, cursor, selected context, and return focus in the V2 board state.
- [x] Port the established `r`, navigation, Enter, cancel, and behind-board rendering contract to V2.
- [x] Filter Objective membership and Task cards only through indexed Release links.
- [x] Persist selection with a Bubble Tea command using the canonical router writer and migration/conflict guards.
- [x] Add `releases/` to the watch set and preserve/diagnose selection through canonical reload.
- [x] Update V2 help and narrow-width rendering without changing geometry on focus.
- [x] Add compatibility tests mirroring the V1 selector cases plus conflict, removal, no-Release, and reload cases.
- [x] Run the unchanged V1 Release selector suite as a regression gate.

## Context Log

Implemented the V2 Release selector and canonical Release-context switching.

Read: the active router, E51 detail, Guardrails, all listed V1/V2 selector,
load, navigation, rendering, IO, and watch files, plus the targeted V2 data
index, Release, router, gate, and Next APIs required to consume
`ReleaseObjectives` and persist `release: R###` safely.

Edited: `internal/board/v2/model.go`, `releases.go`, `releases_test.go`,
`objectives.go`, `card.go`, `update.go`, `io.go`, `load.go`, `view.go`,
`watch.go`, `help.go`, `next_panel.go`, `plain.go`, `run.go`, the V2 watch
and boundary compatibility tests, and this task checklist/log.

Verification:

- `go test ./internal/board/v2` — passed, including selector, filtering,
  reload, removal, conflict, pending-migration, no-Release, and help cases.
- `make build && make test` — passed for every package.
- The unchanged V1 Release selector suite in `internal/board/release_test.go`
  passed as part of `make test`.
- Quick Health-Check evidence was skipped because `.savepoint/Health-Check.md`
  is absent.
