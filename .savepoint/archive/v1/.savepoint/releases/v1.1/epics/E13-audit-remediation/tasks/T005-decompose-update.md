---
id: E13-audit-remediation/T005-decompose-update
status: done
objective: Split the 521-line update.go into focused handler functions for each overlay type and board keys
depends_on: []
---

# T005: Decompose Update() Into Focused Handlers

## Context Files

- `internal/board/update.go` — 521-line monolith; `Update()` is ~190 lines with deep nesting; `updateOverlay()` handles 5 overlay types
- `internal/board/model.go` — 25+ field god struct
- `internal/board/transitions.go` — task transition logic

## Acceptance Criteria

- [x] `Update()` reduced to a dispatch that calls named handler functions
- [x] `handleBoardKey(msg tea.KeyMsg, m Model) (Model, tea.Cmd)` extracted for board-level key handling
- [x] Per-overlay handler functions extracted: `handleHelpOverlay()`, `handleEpicOverlay()`, `handleReleaseOverlay()`, `handleDetailOverlay()`, `handleEpicDetailOverlay()`
- [x] `handleAdvanceTask()` and `handleRetreatTask()` extracted from the near-duplicate space/backspace handlers
- [ ] `update.go` total line count reduced by at least 30%  *(not met — see Drift Notes)*
- [x] All existing tests pass without modification (behavior-preservation refactor)
- [x] `go test ./...` passes

## Implementation Plan

- [x] Extract `handleBoardKey(msg tea.KeyMsg, m Model) (Model, tea.Cmd)` from `Update()` main switch
- [x] Extract overlay update logic into `handleHelpOverlay`, `handleEpicOverlay`, `handleReleaseOverlay`, `handleDetailOverlay`, `handleEpicDetailOverlay`
- [x] Extract `handleAdvanceTask()` and `handleRetreatTask()` from space/backspace handlers (they share near-identical structure)
- [x] Refactor `Update()` to be a thin dispatch: `switch overlay { case none: handleBoardKey; case help: handleHelpOverlay; ... }`
- [x] Verify all `update_test.go` tests still pass
- [x] Run `make build && make test`

## Drift Notes

- AC `update.go` total line count reduced by at least 30% — not met: 543 lines vs 521 original (+22). The decomposition adds function-signature overhead from per-overlay handlers and dispatch that the original monolithic `updateOverlay()` did not have. Behavior is fully preserved and all 64 tests pass.