---
id: E32-audit-register-tui-review/T001-board-audit-data-loading
status: planned
objective: Load audit-register data into the board model.
depends_on:
  - E31-audit-register-data-model/T003-audit-register-discovery
complexity_tier: medium
complexity_reason: Loading follows existing async patterns but introduces a new board data set.
---

# T001: Board Audit Data Loading

## Problem

The board cannot render audit-register state until it can load prompt, findings, register summary, and runs through the existing message flow.

## Context Files

- `internal/board/board.go`
- `internal/board/io.go`
- `internal/board/watch.go`
- `internal/board/interfaces.go`
- `internal/board/io_test.go`
- `internal/data/audit.go`

## Acceptance Criteria

- [ ] Board startup can load audit-register data when `.savepoint/audit/` exists.
- [ ] Opening the Audit Register overlay triggers a fresh audit-register load.
- [ ] Reload/watch behavior refreshes audit-register data after audit files change.
- [ ] Missing audit-register files produce empty data state, not board startup failure.
- [ ] Tests cover successful, missing, and error load paths.

## Implementation Plan

- [ ] Add audit-register fields to board data state.
- [ ] Add an injected data-loading dependency for tests.
- [ ] Add Bubble Tea messages and commands for audit-register loads.
- [ ] Refresh audit data through existing reload/watch paths.
- [ ] Add focused board IO tests.

## Context Log

Pending.
