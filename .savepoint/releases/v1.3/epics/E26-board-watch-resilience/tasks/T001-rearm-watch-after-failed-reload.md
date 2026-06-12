---
id: E26-board-watch-resilience/T001-rearm-watch-after-failed-reload
status: planned
objective: Guarantee file watching is re-armed after a failed reload so live updates never silently stop.
depends_on: []
complexity_tier: medium
complexity_reason: Requires a new message type, an Update branch change, and a regression test of the watch cycle.
---

# T001: Re-arm Watch After Failed Reload

## Problem

`watchFiles` delivers exactly one message and terminates; it is re-armed only in `Init()` and the `reloadMsg` branch of `Update()`. When `reloadTasks` fails (any malformed frontmatter aborts `loadBoardData`), the command returns a generic `errorMsg`, whose `Update` branch only sets `StatusMessage` — the watcher is never re-armed and live updates silently stop until a manual `ctrl+r`. Audit finding H2.

## Context Files

- `internal/board/watch.go`
- `internal/board/update.go`
- `internal/board/update_test.go`
- `internal/board/model.go`

## Acceptance Criteria

- [ ] A failed watch-triggered reload still results in the next file change being observed (watch cycle is never broken).
- [ ] The reload failure message is still shown to the user in the status line.
- [ ] A regression test injects a failing loader, delivers `fileChangeMsg`, executes the returned command chain, and asserts a watch command is re-armed alongside the error status.
- [ ] Manual `ctrl+r` behavior is unchanged.
- [ ] `go test ./internal/board` passes.

## Implementation Plan

- [ ] Add a dedicated `reloadFailedMsg{message string}` in `internal/board/watch.go` and return it from `reloadTasksWithMessage` instead of generic `errorMsg`.
- [ ] Handle `reloadFailedMsg` in `Update()`: set `StatusMessage` and, when `m.Watcher != nil`, return `watchFiles(m.Watcher)`.
- [ ] Keep `errorMsg` semantics unchanged for non-reload errors (task/router/defect write failures).
- [ ] Add the regression test to `internal/board/update_test.go`.
- [ ] Run `go test ./internal/board` and `make build`.

## Context Log

Pending.
