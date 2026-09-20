---
id: E20-clean-up-lifecycle/T004-board-lifecycle-transitions
status: done
objective: Move board task transitions onto shared lifecycle operations without changing the TUI contract.
depends_on: [T003-doctor-lifecycle-diagnostics]
complexity_tier: high
complexity_reason: Board transitions must preserve dependency gates, mtime writes, and user-owned completion.
---

# T004: Board Lifecycle Transitions

## Problem

The board owns transition code for planned, in-progress stages, and done. This keeps the TUI usable, but it duplicates task lifecycle knowledge that should belong to `internal/data`.

## Context Files

- `internal/board/transitions.go`
- `internal/board/transitions_test.go`
- `internal/board/update.go`
- `internal/board/update_test.go`
- `internal/board/io.go`
- `internal/board/integration_test.go`
- `internal/data/lifecycle.go`
- `internal/data/lifecycle_test.go`

## Acceptance Criteria

- [x] Board advance and retreat logic delegates lifecycle state changes to `internal/data` helpers.
- [x] Dependency and epic audit gates remain enforced before completing an audit-stage task.
- [x] User-owned completion behavior remains unchanged: agents still stop before marking tasks done outside explicit user action.
- [x] Mtime conflict retry behavior still compares the correct lifecycle state.
- [x] Existing board transition tests pass with added coverage for invalid or legacy-loaded task states.

## Implementation Plan

- [x] Add lifecycle transition helpers or result types in `internal/data`.
- [x] Replace board-local state mutation in `Advance` and `Retreat` with the shared helpers.
- [x] Keep board-specific dependency and status-message formatting in `internal/board`.
- [x] Add tests around audit-to-done gating, done-to-audit retreat, and stale metadata loaded from disk.

## Context Log

- Read: `internal/board/transitions.go`, `internal/board/transitions_test.go`, `internal/board/update.go`, `internal/board/update_test.go`, `internal/board/io.go`, `internal/board/integration_test.go`, `internal/data/lifecycle.go`, `internal/data/lifecycle_test.go`.
- Edited: `internal/data/lifecycle.go`, `internal/data/lifecycle_test.go`, `internal/board/transitions.go`, `internal/board/transitions_test.go`, `internal/board/update.go`, `internal/board/update_test.go`, `internal/board/io.go`.
- Verification: board `Advance`/`Retreat` now call `internal/data` lifecycle helpers; `CanAdvance` uses the shared advance target to preserve audit-to-done dependency and epic gates; mtime retry compares normalized lifecycle state through `data.SameTaskLifecycleForTransition`.
- TUI priority: not available from this non-interactive session; router already points at this active task.
- Quality gates: `go test ./internal/data ./internal/board` passed; `make build` passed; `make test` passed.
