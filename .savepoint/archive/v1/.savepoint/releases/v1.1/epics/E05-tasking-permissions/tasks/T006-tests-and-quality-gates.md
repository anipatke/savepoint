---
id: E05-tasking-permissions/T006-tests-and-quality-gates
status: done
objective: "Write unit tests for the `p` priority router update logic and verify quality gates pass"
depends_on: ["E05-tasking-permissions/T004-implement-m-hotkey", "E05-tasking-permissions/T005-update-help-overlay"]
---

# T006: Tests & Quality Gates

## Acceptance Criteria

- Unit tests for `writeRouterTask()` covering:
  - Focused task → router updated with release/epic/task/state
  - Last uncompleted task → router still stays in `task-building` for the focused task
  - Focused `done` task → no-op
- Unit test for `p` priority key handler in update_test.go:
  - Pressing `p` updates the model's StatusMessage
- Unit test for help overlay:
  - `RenderHelp` output contains `p` priority shortcut
- `go build -o savepoint main.go` passes
- `go test ./...` passes

## Implementation Plan

- [x] Add/update router priority tests for focused task state derivation
- [x] Add `p` key tests to update_test.go
- [x] Add help overlay test coverage for the `p` shortcut
- [x] Run build quality gate equivalent
- [x] Run `go test ./...`
- [x] Fix priority key audit-pending regression

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
- internal/board/view.go
- internal/board/view_test.go
- internal/board/card.go
- internal/board/card_test.go
- internal/board/transitions.go
- internal/board/transitions_test.go
- .savepoint/router.md

Estimated input tokens: 11000

Notes:
- Priority key was renamed from `m` to `p`; footer now shows `p: Priority`, and help overlay includes the `p` shortcut.
- Fixed priority action regression: pressing `p` on the last incomplete task no longer writes `audit-pending`; it always writes `task-building` for the focused non-done task.
- Router restored from accidental `audit-pending` to `E05-tasking-permissions/T006-tests-and-quality-gates`.
- Added/update regression coverage in `internal/board/update_test.go`.
- Focused test: `go test ./internal/board` passed.
- Quality gate equivalent: `go test ./...` passed.
- Quality gate equivalent: `go run ./internal/buildtool build` passed.
- Explicit build gate: `go build -o savepoint main.go` passed.
- Repo wrapper gates: `make build` and `make test` could not run because `make` is not installed in this Windows shell.
