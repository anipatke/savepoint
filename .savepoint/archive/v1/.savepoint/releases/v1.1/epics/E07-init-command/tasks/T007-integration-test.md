---
id: E07-init-command/T007-integration-test
status: done
objective: "End-to-end integration test of init flow"
depends_on: ["E07-init-command/T006-clipboard"]
---

# T007: Integration Test

## Acceptance Criteria

- Full init pipeline runs end-to-end in temp directory
- Validates: empty dir, existing project, already-initialized, --force, --install
- Tests pass on Windows, Linux, macOS
- Output includes magic prompt to stdout

## Implementation Plan

- [x] Add `internal/init/integration_test.go`
- [x] Test: init in empty directory
- [x] Test: init in compatible project
- [x] Test: init fails without --force on existing .savepoint
- [x] Test: init with --force overwrites
- [x] Test: init with --install triggers npm install
- [x] Verify magic prompt in stdout
- [x] Run full test suite: `make build && make test`