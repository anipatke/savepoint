---
id: E11-board-refresh-fix/T001-debug-logging
status: done
objective: Add debug logging to watch.go to trace where refresh fails
depends_on: []
---

# T001: Debug Logging

## Acceptance Criteria

- [ ] Print when watchFiles() starts watching
- [ ] Print when event received from watcher
- [ ] Print when fileChangeMsg received in update.go
- [ ] Print any errors from watcher Errors channel
- [ ] Print when reloadTasks() is called

## Implementation Plan

- [ ] In `watch.go` - add fmt.Printf at start of watchFiles()
- [ ] In `watch.go` - add fmt.Printf for each event received
- [ ] In `watch.go` - add fmt.Printf for watcher errors
- [ ] In `update.go` - add fmt.Printf when fileChangeMsg received
- [ ] In `watch.go` - add fmt.Printf when reloadTasks() called
- [ ] Run board, make file change, observe logs