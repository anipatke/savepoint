---
id: E11-board-refresh-fix/T003-error-handling
status: done
objective: Add error handling instead of silent returns in watcher
depends_on:
    - E11-board-refresh-fix/T001-debug-logging
---

# T003: Error Handling

## Acceptance Criteria

- [ ] Log instead of silent return on channel close
- [ ] Log watcher errors that are currently ignored
- [ ] Log any reload failures

## Implementation Plan

- [ ] In `watch.go` - replace `if !ok { return nil }` with error log
- [ ] In `watch.go` - add logging for `w.Errors` channel events
- [ ] In `reloadTasks()` - add error handling for loadBoardData failure
- [ ] Run tests