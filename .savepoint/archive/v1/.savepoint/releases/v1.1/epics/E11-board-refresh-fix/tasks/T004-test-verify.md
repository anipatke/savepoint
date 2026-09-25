---
id: E11-board-refresh-fix/T004-test-verify
status: done
objective: Test and verify board refreshes on file change
depends_on:
    - E11-board-refresh-fix/T001-debug-logging
    - E11-board-refresh-fix/T002-increase-debounce
    - E11-board-refresh-fix/T003-error-handling
---

# T004: Test & Verify

## Acceptance Criteria

- [ ] Board starts with debug logs visible
- [ ] Create new task file while board running
- [ ] Board auto-refreshes within 1 second
- [ ] New task appears in Planned column

## Implementation Plan

- [ ] Start board from terminal
- [ ] Verify debug log output shows watchFiles() started
- [ ] Create test task file: `echo "---" > E06/tasks/T999-test.md`
- [ ] Wait 2 seconds for debounce + reload
- [ ] Verify task appears in board
- [ ] Or: debug log shows fileChangeMsg received
- [ ] Clean up test file
- [ ] Run `make build && make test`