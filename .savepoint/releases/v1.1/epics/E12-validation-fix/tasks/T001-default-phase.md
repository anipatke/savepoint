---
id: E12-validation-fix/T001-default-phase
status: done
objective: Default phase to build when status=in_progress but phase missing
depends_on: []
---

# T001: Default Phase to Build

## Context Files

- `internal/data/parser.go`
- `internal/data/lifecycle.go`
- `internal/data/parser_test.go`

## Acceptance Criteria

- [x] firstStage() returns "build" when both stage and phase are empty
- [x] Tasks without phase field parse successfully
- [x] Make build passes
- [x] go test ./... passes

## Implementation Plan

- [x] Set default at parse time after validation:
  - In parser.go: ParseTaskFile() line 60-63
  - After validation: `if task.Column == in_progress && task.Stage == "" { task.Stage = StageBuild }`
- [x] Run `make build` to verify change
- [x] Run `go test ./...` all pass

Notes:
- Default set at parsing time, not validation time
- Works with both planned and done tasks (phase ignored)
- Better error messages in lifecycle.go already implemented
- Tests added for new behavior
