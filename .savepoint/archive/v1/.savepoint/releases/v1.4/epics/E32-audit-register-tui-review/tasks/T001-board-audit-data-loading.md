---
id: E32-audit-register-tui-review/T001-board-audit-data-loading
status: done
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

- [x] Board startup can load audit-register data when `.savepoint/audit/` exists.
- [x] Opening the Audit Register overlay triggers a fresh audit-register load. *(Mechanism delivered: `loadAuditRegisterCmd` performs the fresh load and reload/watch keep the model current; the `A`-key overlay entry point is wired in T002, which owns the overlay UI.)*
- [x] Reload/watch behavior refreshes audit-register data after audit files change.
- [x] Missing audit-register files produce empty data state, not board startup failure.
- [x] Tests cover successful, missing, and error load paths.

## Implementation Plan

- [x] Add audit-register fields to board data state.
- [x] Add an injected data-loading dependency for tests.
- [x] Add Bubble Tea messages and commands for audit-register loads.
- [x] Refresh audit data through existing reload/watch paths.
- [x] Add focused board IO tests.

## Context Log

Read: `router.md`, `E32-Detail.md`, this task, `internal/data/audit.go`, `internal/data/audit_register.go`,
`internal/data/audit_run.go`, and board `io.go`/`watch.go`/`interfaces.go`/`model.go`/`board.go`/`update.go`/`view.go`
plus `io_test.go`/`board_test.go` for existing async-load and reload patterns.

Edited:
- `internal/board/model.go` — added `AuditState{ Audit data.AuditRegisterSet }` and embedded it in `Model`.
- `internal/board/interfaces.go` — added the injectable `auditLoader` interface, the production `dataAuditLoader`
  adapter over `data.LoadAuditRegisterSet`, and threaded `AuditLoader` through `ModelDependencies`/defaults/overrides.
- `internal/board/io.go` — added `loadAuditRegisterCmd` (strict fresh load → `auditRegisterMsg`/`errorMsg`) and
  `loadAuditBestEffort` (used by startup + aggregate reload; degrades a malformed set to empty like router reads).
- `internal/board/watch.go` — added `auditRegisterMsg`, the `audit` field on `reloadMsg`, the best-effort audit load in
  `reloadTasksWithMessage`, and `newWatcher` now watches the `audit/` subtree so audit edits drive the reload path.
- `internal/board/board.go` — startup loads audit best-effort into `model.Audit` (missing tree → empty, never fatal).
- `internal/board/update.go` — `reloadMsg` sets `m.Audit`; new `auditRegisterMsg` case sets `m.Audit`.
- `internal/board/io_test.go` — added command tests (populated / missing / malformed), reload refresh + malformed
  tolerance tests, and a startup-load test.

Scope note: the audit overlay UI (the `A` key, tabs, and view rendering) belongs to T002; T001 delivers only the
loading foundation and the reusable load command T002 will trigger on open.

Quality gates: `make build && make test` pass; `go vet ./internal/board` clean.
