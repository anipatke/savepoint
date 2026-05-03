---
id: E08-board-command/T008-integration-test
status: done
objective: "End-to-end integration test of board flow"
depends_on: ["E06-audit-command/T007-apply-close", "E07-init-command/T007-integration-test"]
---

# T008: Integration Test

## Acceptance Criteria

- Full board pipeline runs end-to-end
- TUI navigation, transitions, detail pane work
- Non-TTY fallback outputs correct table
- Status writes preserve task content
- Tests pass on Windows, Linux, macOS

## Implementation Plan

- [x] Add `internal/board/integration_test.go`
- [x] Test: board launches in TTY mode
- [x] Test: board falls back to plain table in non-TTY
- [x] Test: keyboard navigation (arrows, vim keys)
- [x] Test: status transitions preserve content
- [x] Test: mtime conflict detection works
- [x] Test: release/epic filters work
- [x] Run full test suite: `make build && make test`