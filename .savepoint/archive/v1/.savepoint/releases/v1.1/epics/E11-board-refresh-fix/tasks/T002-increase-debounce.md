---
id: E11-board-refresh-fix/T002-increase-debounce
status: done
objective: Increase debounce timer from 100ms to 500ms for reliable detection
depends_on:
    - E11-board-refresh-fix/T001-debug-logging
---

# T002: Increase Debounce

## Acceptance Criteria

- [ ] Change timer from 100ms to 500ms
- [ ] Verify board still works

## Implementation Plan

- [ ] In `internal/board/watch.go` line 31:
  - FROM: `timer := time.NewTimer(100 * time.Millisecond)`
  - TO: `timer := time.NewTimer(500 * time.Millisecond)`
- [ ] Run `make build` to verify