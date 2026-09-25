---
type: epic-design
status: audited
---

# E11: Board Auto-Refresh Fix

## Purpose

Fix the TUI board so it auto-refreshes when task files change on disk. The watch mechanism exists but isn't triggering reloads on file changes.

## What this epic adds

- Debug logging to trace where the refresh fails
- Increased debounce timer for reliable event detection
- Error handling for silent watcher failures
- Verified auto-refresh on file changes

## Components

| Module | Purpose |
|--------|---------|
| `internal/board/watch.go` | Debug logs, increased debounce, error handling |
| `internal/board/update.go` | Debug logs for reload flow |

## Boundaries

**In scope:**
- Debug logging to identify failure point
- Fix debounce timing
- Fix error handling
- Verify refresh works

**Out of scope:**
- No new UI features
- No new core functionality