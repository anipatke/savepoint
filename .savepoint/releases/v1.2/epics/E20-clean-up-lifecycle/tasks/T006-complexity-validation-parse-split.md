---
id: E20-clean-up-lifecycle/T006-complexity-validation-parse-split
status: in_progress
stage: build
objective: Remove complexity_reason word-count validation from the parser so board rendering is never blocked by a policy violation.
depends_on: [T002-parser-writer-lifecycle]
complexity_tier: low
complexity_reason: Single call-site removal; validation remains in write path and doctor.
---

# T006: Split Complexity Validation from Parse

## Problem

`ParseTaskFile` calls `ValidateComplexity`, which returns a hard error when `complexity_reason` exceeds 20 words. This turns a policy violation into a parse failure, preventing the board from rendering any task that trips the limit — even read-only.

## Context Files

- `internal/data/parser.go`
- `internal/data/lifecycle.go`
- `internal/data/parser_test.go`
- `internal/doctor/checks.go`

## Acceptance Criteria

- [x] `ParseTaskFile` does not call `ValidateComplexity`; tasks with long `complexity_reason` values load without error.
- [x] `ValidateTaskLifecycle` (write path) still calls `ValidateComplexity` and rejects over-limit reasons on write.
- [x] `doctor/checks.go` still surfaces long `complexity_reason` as a diagnostic problem.
- [x] Existing parser tests for valid complexity fields continue to pass.
- [x] A new parser test confirms a task with a 21-word reason loads successfully.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Delete the `ValidateComplexity` call and surrounding error return from `ParseTaskFile` in `internal/data/parser.go`.
- [x] Add a parser test: task file with a `complexity_reason` of 21+ words parses without error.
- [x] Confirm `ValidateTaskLifecycle` and `doctor/checks.go` still call `ValidateComplexity` (no changes needed there).
- [x] Run `go test ./internal/data ./internal/doctor` then `make build && make test`.
