---
id: E15-hardening/T001-benchmarks
status: planned
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

- [ ] Benchmark added for RenderCard with varied content widths
- [ ] Benchmark added for renderColumn with multiple tasks
- [ ] Benchmark added for layout calculations at different widths
- [ ] Benchmarks are repeatable with consistent results
- [ ] Benchmarks don't modify package state
- [ ] `go test -bench=. ./internal/board/` runs without errors

## Implementation Plan

- [ ] Add benchmark functions in view_test.go, card_test.go, column_test.go
- [ ] Create test data with representative task mixtures
- [ ] Run benchmarks and document baseline
- [ ] Run `make build && make test`
