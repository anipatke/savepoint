---
id: E15-hardening/T001-benchmarks
status: done
objective: Add benchmark tests for board render functions
depends_on: []
---

# T001: Add Benchmark Tests for Render Functions

## Context Files

- `internal/board/view.go` — renderColumn, renderColumnHeaders
- `internal/board/card.go` — RenderCard
- `internal/board/column.go` — RenderColumn
- `internal/board/layout.go` — layout calculations

## Acceptance Criteria

- [x] Benchmark added for RenderCard with varied content widths
- [x] Benchmark added for renderColumn with multiple tasks
- [x] Benchmark added for layout calculations at different widths
- [x] Benchmarks are repeatable with consistent results
- [x] Benchmarks don't modify package state
- [x] `go test -bench=. ./internal/board/` runs without errors

## Implementation Plan

- [x] Add benchmark functions in view_test.go, card_test.go, column_test.go
- [x] Create test data with representative task mixtures
- [x] Run benchmarks and document baseline
- [x] Run `make build && make test`
