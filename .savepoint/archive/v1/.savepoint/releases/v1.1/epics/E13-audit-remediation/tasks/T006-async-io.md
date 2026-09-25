---
id: E13-audit-remediation/T006-async-io
status: done
objective: Extract synchronous filesystem I/O from Update() into tea.Cmd async pattern following Bubble Tea conventions
depends_on:
  - E13-audit-remediation/T005-decompose-update
---

# T006: Extract I/O From Update() Into tea.Cmd Async Pattern

## Context Files

- `internal/board/update.go` — `writeTaskStatus()`, `writeRouterTask()`, `writeRouterReleaseEpic()`, `readEpicDetailFile()`, `selectEpicPanelEpic()` perform synchronous I/O inside `Update()`
- `internal/board/model.go` — `writeRouterReleaseEpic()` and `writeRouterTask()` methods on Model
- `internal/board/watch.go` — `reloadTasks` silently swallows errors (returns `nil` on error)

## Acceptance Criteria

- [x] All filesystem reads in `Update()` dispatch `tea.Cmd` functions that return result messages
- [x] All filesystem writes in `Update()` dispatch `tea.Cmd` functions that return result messages
- [x] `writeRouterTask` extracted to a `tea.Cmd` that returns a `routerWriteMsg` on success or `errorMsg` on failure
- [x] `writeRouterReleaseEpic` extracted to a `tea.Cmd` that returns a `routerWriteMsg` on success or `errorMsg` on failure
- [x] `writeTaskStatus` for task transitions extracted to a `tea.Cmd` that returns a `taskWriteMsg` or `errorMsg`
- [x] `readEpicDetailFile` extracted to a `tea.Cmd` that returns the file content as a message
- [x] `reloadTasks` in `watch.go` returns an `errorMsg` on failure instead of `nil`
- [x] Status messages for write errors are displayed to the user (existing error display mechanism)
- [x] No synchronous `os.ReadFile`, `os.WriteFile`, or `os.Stat` calls remain in `Update()` or `handleBoardKey()`
- [x] `go test ./...` passes with no regressions

## Implementation Plan

- [x] Define message types: `routerWriteMsg`, `taskWriteMsg`, `epicDetailMsg`, `auditContentMsg`, `errorMsg`
- [x] Create `writeRouterTaskCmd` and `writeRouterReleaseEpicCmd` that write router state and return `routerWriteMsg` or `errorMsg`
- [x] Create `writeTaskStatusCmd` for task transitions returning `taskWriteMsg` or `errorMsg`
- [x] Create `readEpicDetailCmd` and `readEpicAuditCmd` for file reads
- [x] Extract `selectEpicPanelEpic` I/O into a `tea.Cmd`
- [x] Update `Update()` to dispatch commands and handle result messages in new branches
- [x] Fix `reloadTasks` in `watch.go` to return `errorMsg` on error
- [x] Update tests for async message pattern
- [x] Run `make build && make test`