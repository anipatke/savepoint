---
id: E05-tasking-permissions/T004-implement-m-hotkey
status: done
objective: "Implement TUI `m` hotkey that explicitly sets the router's task from the currently focused task"
depends_on: ["E05-tasking-permissions/T002-update-router-md", "E05-tasking-permissions/T005-add-help-entry"]
---

# T004: Implement TUI `m` hotkey

## Acceptance Criteria

- Pressing `m` in the TUI board (when no overlay is open) reads the currently focused task and updates `router.md` with:
  - `release`: from focused task
  - `epic`: from focused task
  - `task`: full task ID
  - `state`: `task-building` (or `audit-pending` if last uncompleted task in epic)
- Pressing `m` when no task is focused does nothing
- Pressing `m` when focused on a `done` task does nothing (or shows status message)
- Router `task` field in `writeRouterReleaseEpic` is now set explicitly (currently it writes release+epic only)
- The `m` key does NOT work when any overlay is open (consistent with other board keys)
- The `m` key displays a status message confirming the router update: "Router set to {release} {epic}/{task}" or "Audit pending for {epic}" if last task

## Implementation Plan

- [x] Add `writeRouterTask(task data.Task)` method to `model.go`: reads current router, sets release/epic/task/state fields, writes via `data.WriteRouterState`
- [x] Add `m` key handler in `update.go` main switch (before epic panel routing, like `q`/`e`/`r`/`?`)
- [x] Handler: calls `focusedTask()` to get current task, calls `m.writeRouterTask()`, sets `m.StatusMessage` on success/error
- [x] Detect "last uncompleted task in epic" logic: scan `m.AllTasks` for the same epic, check if any are not done
- [x] Update `writeRouterReleaseEpic()` is unchanged (still only writes release+epic, used by navigation)
- [x] Run `make build && make test` to verify

## Strategic Visual Rework Plan

This task is reopened because the `m` hotkey work exposed router-priority and board rendering paths where terminal color fallback behavior still makes the TUI look patchy between black and gray. Do not continue with incremental one-widget fixes. Apply one bottom-up rendering policy.

- [x] Force Lipgloss to a deterministic 256-color profile before board rendering starts.
- [x] Change black 256-color fallbacks for Background, Surface, and Surface 2 to the same actual ANSI256 black value.
- [x] Centralize the background policy in `internal/styles/styles.go`: TUI surfaces emit no explicit background escapes; terminal default background remains the single source of truth.
- [x] Remove or avoid explicit background fills from nested panels/cards where they create padded gray bars.
- [x] Normalize focus borders across task columns, task cards, overlays, and epic sidebar. Prefer single-line box borders if rounded borders continue rendering as dash bars in Warp.
- [x] Add regression tests that fail if rendered board/card/router-priority output emits explicit background escapes.
- [x] Verify representative board states in both Windows Terminal and Warp after code tests pass.

## Context Log

Files read:
- .savepoint/router.md
- .savepoint/releases/v1.1/epics/E05-tasking-permissions/E05-Detail.md
- .savepoint/releases/v1.1/epics/E05-tasking-permissions/tasks/T004-implement-m-hotkey.md
- internal/board/model.go
- internal/board/update.go
- internal/board/update_test.go
- internal/board/board.go
- internal/board/view.go
- internal/board/column_test.go
- internal/board/view_test.go
- internal/board/render_policy_test.go
- internal/styles/palette.go
- internal/styles/styles.go
- internal/data/router.go
- internal/data/write.go

Estimated input tokens: 12000

Notes:
- Focused task comes from `m.Tasks[m.FocusedColumn][m.FocusedTask]`
- Router state derivation: if all other tasks in the same epic+release are `done`, set `state: audit-pending` and clear `task`; otherwise `state: task-building` with the task ID
- Focused test: `go test ./internal/board` passed
- Quality gate: `make build` could not run because `make` is unavailable in this shell
- Quality gate equivalent: `go run ./internal/buildtool build` passed
- Quality gate equivalent: `go test ./...` passed
- Strategic color rework: `lipgloss.SetColorProfile(termenv.ANSI256)` is set before board rendering.
- Strategic color rework: Background, Surface, and Surface 2 ANSI256 fallbacks now all use `16`.
- Strategic color rework: nested task/card/panel text no longer paints gray background fills; root-level full-width lines also avoid explicit background escapes.
- Strategic border rework: rounded borders replaced with single-line borders to avoid Warp dash-bar artifacts.
- Regression scan: no `RoundedBorder`, rounded border glyphs, or ANSI256 gray background fallback codes (`232`/`233`) remain under `internal/`.
- Remaining manual check: visually compare representative board states in Windows Terminal and Warp.
- Follow-up color consistency pass: removed explicit background painting from `Divider`, `HeaderFrame`, `BoardFrame`, and `RootLine`; this avoids mixing terminal-default cells with `48;5;16` cells after nested style resets.
- Regression tightened: `internal/board/render_policy_test.go` now fails on any `48;` background escape or basic black `40m` escape across board, card, detail, release/epic dropdown, and help render paths.
- Focused test: `go test ./internal/board` passed after the no-background policy change.
- Quality gate: `make build` still cannot run because `make` is unavailable in this shell.
- Quality gate equivalent: `go run ./internal/buildtool build` passed after the no-background policy change.
- Quality gate equivalent: `go test ./...` passed after the no-background policy change.
- Manual verification: user confirmed closeout after the PowerShell/Warp color consistency pass.
- Closeout focused test: `go test ./internal/board` passed.
- Closeout quality gate: `make build` failed because `make` is unavailable in this shell.
- Closeout quality gate equivalent: `go run ./internal/buildtool build` passed.
- Closeout quality gate equivalent: `go test ./...` passed.
