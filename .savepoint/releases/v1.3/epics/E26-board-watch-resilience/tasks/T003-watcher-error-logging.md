---
id: E26-board-watch-resilience/T003-watcher-error-logging
status: planned
objective: Surface fsnotify watcher errors through the board debug log instead of discarding them.
depends_on: []
complexity_tier: low
complexity_reason: One-branch logging change in the watch loop with a small test.
---

# T003: Watcher Error Logging

## Problem

The watch loop drains `w.Errors` but discards the values — not even `debugf` sees them (`internal/board/watch.go:85-88`). fsnotify failures such as watch-buffer overflow on Windows vanish, making stale-board reports undiagnosable. Audit finding L6.

## Context Files

- `internal/board/watch.go`
- `internal/board/debug.go`

## Acceptance Criteria

- [ ] Every value received from the watcher error channel is logged via `debugf` with a `watcher:` prefix.
- [ ] The watch loop continues running after logging an error (behavior otherwise unchanged).
- [ ] `go test ./internal/board` passes.

## Implementation Plan

- [ ] Capture the error value in the `w.Errors` case and call `debugf("watcher: error %v", err)` before continuing the loop.
- [ ] If practical with the existing debug hooks, add a test asserting the log call; otherwise verify by inspection and note it in the context log.
- [ ] Run `go test ./internal/board`.

## Context Log

Pending.
