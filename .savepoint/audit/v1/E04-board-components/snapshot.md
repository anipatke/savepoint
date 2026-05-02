# Manual Audit Snapshot: E04-board-components

Created: 2026-05-01

Reason: `.savepoint/router.md` is in `audit-pending` for `E04-board-components`, but `.savepoint/audit/E04-board-components/snapshot.md` was missing. Per router instructions, this is a one-time manual snapshot from the known epic scope.

## Router State

```yaml
state: audit-pending
release: v1
epic: E04-board-components
task: E04-board-components/T006-help-overlay
next_action: "Epic E04-board-components complete. Start a new agent session for audit."
```

## Epic Scope

Epic design: `.savepoint/releases/v1/epics/E04-board-components/E04-Detail.md`

Tasks:

- `.savepoint/releases/v1/epics/E04-board-components/tasks/T001-column.md`
- `.savepoint/releases/v1/epics/E04-board-components/tasks/T002-card.md`
- `.savepoint/releases/v1/epics/E04-board-components/tasks/T003-epic-panel.md`
- `.savepoint/releases/v1/epics/E04-board-components/tasks/T004-detail-overlay.md`
- `.savepoint/releases/v1/epics/E04-board-components/tasks/T005-release-dropdown.md`
- `.savepoint/releases/v1/epics/E04-board-components/tasks/T006-help-overlay.md`

## Files Listed By Epic Design

- `internal/board/column.go`
- `internal/board/card.go`
- `internal/board/epic_panel.go`
- `internal/board/detail.go`
- `internal/board/dropdown.go` (planned; no live file found)
- `internal/board/help.go`
- `internal/board/card_test.go`
- `internal/board/column_test.go`

## Additional Files Touched By Task Logs

- `internal/board/model.go`
- `internal/board/update.go`
- `internal/board/view.go`
- `internal/board/layout.go`
- `internal/board/detail_test.go`
- `internal/board/epic_panel_test.go`
- `internal/board/help_test.go`
- `internal/board/release.go`
- `internal/board/release_test.go`
- `internal/board/update_test.go`
- `internal/board/view_test.go`
- `internal/styles/styles.go`
- `internal/styles/palette.go`

## Verification Re-run During Audit

- `go test ./internal/board/...` passed.
- `go test ./...` passed.

## Scoped Review Notes

- `RenderCard` is implemented and tested, but `RenderColumn` still renders plain task labels through `taskLabel`, and `rg RenderCard internal/board` shows no production call site outside `card.go`.
- Epic and release selection currently mutate `SelectedEpic` / `SelectedRelease` only. There is no data reload or regrouping path in `internal/board/update.go`, so the task-list update acceptance criteria are not yet satisfied.
- Detail/help overlays are composited over dimmed base board. Epic/release dropdowns use `lipgloss.Place` over a blank canvas, so the board is not visible behind those overlays.

## Post-Audit Implementation

Updated after user request to implement all audit findings:

- `RenderColumn` now renders `RenderCard` for each task.
- `Model` now stores `AllTasks` and regroups visible tasks when epic/release selection changes.
- Base board `up`/`k` and `down`/`j` move focused task selection.
- Epic/release dropdowns now composite over the dimmed board.
- Detail overlay text now wraps long rows, description, and acceptance criteria.
- Verification after implementation: `go build ./...`, `go vet ./...`, and `go test ./...` passed.
