---
id: E20-clean-up-lifecycle/T002-parser-writer-lifecycle
status: done
objective: Route task parser and writer lifecycle behavior through the shared contract.
depends_on: [T001-lifecycle-contract]
complexity_tier: medium
complexity_reason: Parser and writer changes span several files but follow the T001 contract.
---

# T002: Parser and Writer Lifecycle

## Problem

The parser currently has load-time compatibility behavior that is not expressed through the same code path as write-time validation. This creates room for future legacy metadata defects.

## Context Files

- `internal/data/parser.go`
- `internal/data/parser_test.go`
- `internal/data/write.go`
- `internal/data/write_test.go`
- `internal/data/lifecycle.go`
- `internal/data/lifecycle_test.go`

## Acceptance Criteria

- [x] `ParseTaskFile` uses the shared lifecycle contract for status normalization and stage compatibility.
- [x] `WriteTaskStatus` uses the shared lifecycle contract for canonical write validation.
- [x] Legacy `stage: implementation` and legacy `phase: implementation` on non-in-progress tasks load without blocking board operations and are removed on canonical writes.
- [x] Invalid `status: in_progress` stages still fail with clear canonical errors.

## Implementation Plan

- [x] Replace parser-local lifecycle decisions with calls into the shared contract.
- [x] Replace writer lifecycle validation with the same write-mode contract.
- [x] Add regression tests for both `stage: implementation` and `phase: implementation` outside `in_progress`.
- [x] Verify canonical writes still remove stale `phase` and non-active `stage` fields.

## Context Log

- Read: `internal/data/parser.go`, `internal/data/parser_test.go`, `internal/data/write.go`, `internal/data/write_test.go`, `internal/data/lifecycle.go`, `internal/data/lifecycle_test.go`.
- Edited: `internal/data/lifecycle.go`, `internal/data/lifecycle_test.go`, `internal/data/parser_test.go`, `internal/data/write_test.go`.
- Lifecycle contract now treats `implementation` as a load-only legacy stage alias for stale non-`in_progress` metadata while preserving canonical write validation.
- Regression coverage verifies stale `stage: implementation` and `phase: implementation` load/cleanup behavior, plus canonical rejection for `in_progress` `implementation` stages.
- Quality gates: `go test ./internal/data`, `make build`, and `make test` passed.
