---
id: E05-tasking-permissions/T006-tests-and-quality-gates
status: planned
objective: "Write unit tests for the `m` hotkey router update logic and verify quality gates pass"
depends_on: ["E05-tasking-permissions/T004-implement-m-hotkey", "E05-tasking-permissions/T005-update-help-overlay"]
---

# T006: Tests & Quality Gates

## Acceptance Criteria

- Unit tests for `writeRouterTask()` covering:
  - Focused task → router updated with release/epic/task/state
  - Last uncompleted task → router shows `audit-pending`
  - Focused `done` task → no-op
- Unit test for `m` key handler in update_test.go:
  - Pressing `m` updates the model's StatusMessage
- Unit test for help overlay:
  - `RenderHelp` output contains `m: update router`
- `go build -o savepoint main.go` passes
- `go test ./...` passes

## Implementation Plan

- [ ] Add `writeRouterTask_test.go` (or add to model_test.go) for router state derivation
- [ ] Add `m` key test to update_test.go
- [ ] Add help overlay test to help_test.go (if none exists, add basic test)
- [ ] Run `go build -o savepoint main.go`
- [ ] Run `go test ./...`
- [ ] Fix any failures

## Context Log

Files read:
- .savepoint/releases/v1.1/epics/E05-tasking-permissions/E05-Detail.md
- .savepoint/releases/v1.1/epics/E05-tasking-permissions/tasks/T006-tests-and-quality-gates.md
- internal/board/model.go
- internal/board/update.go
- internal/board/help.go
- internal/board/update_test.go
- internal/board/model_test.go
- internal/board/help_test.go

Estimated input tokens: 4000

Notes:
