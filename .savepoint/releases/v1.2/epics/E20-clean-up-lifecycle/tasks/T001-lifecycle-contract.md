---
id: E20-clean-up-lifecycle/T001-lifecycle-contract
status: done
objective: Define a single task lifecycle contract API in internal/data.
depends_on: []
complexity_tier: high
complexity_reason: Cross-module lifecycle rules need a shared API before parser, doctor, and board changes.
---

# T001: Lifecycle Contract

## Problem

Task status and stage rules are currently repeated in parser validation, write validation, doctor checks, and board transitions. That makes compatibility fixes easy to implement inconsistently.

## Context Files

- `internal/data/task.go`
- `internal/data/lifecycle.go`
- `internal/data/lifecycle_test.go`
- `internal/data/parser.go`
- `internal/data/parser_test.go`
- `internal/data/write.go`
- `internal/data/write_test.go`

## Acceptance Criteria

- [x] `internal/data` exposes a small task lifecycle contract for canonical statuses, canonical stages, legacy aliases, parse normalization, write validation, and transition validation.
- [x] The contract distinguishes load-time compatibility from write-time canonical validation.
- [x] Missing task status, legacy `todo`, legacy `phase`, stale non-in-progress `stage`, and invalid in-progress stage cases are covered by focused tests.
- [x] Existing public behavior remains unchanged except where explicitly covered by compatibility tests.

## Implementation Plan

- [x] Add focused lifecycle types or helper functions in `internal/data` without changing board or doctor call sites yet.
- [x] Move canonical status/stage and legacy alias decisions behind those helpers.
- [x] Add table-driven tests for canonical, legacy, malformed, and stale metadata cases.
- [x] Keep error messages canonical around `status` and `stage`.

## Context Log

- Read: `internal/data/task.go`, `internal/data/lifecycle.go`, `internal/data/lifecycle_test.go`, `internal/data/parser.go`, `internal/data/parser_test.go`, `internal/data/write.go`, `internal/data/write_test.go`.
- Edited: `internal/data/lifecycle.go`, `internal/data/parser.go`, `internal/data/lifecycle_test.go`.
- Quality gates: `go test ./internal/data` passed; `make build` passed; `make test` passed.
- TUI priority note: no active TUI session was available for pressing `p`; router already pointed at this task.
