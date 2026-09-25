---
id: E13-audit-remediation/T001-safe-cleanup
status: done
objective: Remove dead code, replace stdlib reimplementations, co-locate layout constants, remove hardcoded state maps
depends_on: []
---

# T001: Safe Cleanup — Dead Code, Stdlib Replacements, Constant Co-location

## Context Files

- `internal/board/column.go` — contains dead `taskLabel()` function and misplaced `colOverhead` constant
- `internal/doctor/report.go` — contains unused `CheckResult` type
- `internal/doctor/repairs.go` — reimplements `strings.Contains` and `strings.Index`
- `internal/doctor/repairs.go` — `GateSuggestion` has unused `exitCode` parameter
- `internal/buildtool/main.go` — reimplements `strings.TrimSpace` as custom `trimSpace()`
- `internal/doctor/checks.go` — hardcoded `validStates` map duplicates `data.IsCanonicalColumn()`/`data.IsCanonicalStage()`
- `internal/board/layout.go` — layout constants should host `colOverhead`

## Acceptance Criteria

- [x] `taskLabel()` removed from `column.go`
- [x] `CheckResult` type removed from `report.go`
- [x] Unused `exitCode` parameter removed from `GateSuggestion` in `repairs.go`
- [x] `contains()` and `indexOf()` in `repairs.go` replaced with `strings.Contains` and `strings.Index`; `strings` import added
- [x] `trimSpace()` in `buildtool/main.go` replaced with `strings.TrimSpace`; unused `trimSpace` function removed
- [x] `validStates` map removed from `checks.go`; `data.IsCanonicalColumn()` and `data.IsCanonicalStage()` used instead
- [x] `colOverhead` constant moved from `column.go` to `layout.go`
- [x] `go test ./...` passes with no regressions

## Implementation Plan

- [x] Remove `taskLabel()` from `internal/board/column.go`
- [x] Remove `CheckResult` type from `internal/doctor/report.go`
- [x] Remove unused `exitCode` param from `GateSuggestion` in `internal/doctor/repairs.go`
- [x] Replace `contains()`/`indexOf()` with `strings.Contains`/`strings.Index` in `repairs.go`
- [x] Replace `trimSpace()` with `strings.TrimSpace` in `buildtool/main.go`
- [x] Replace `validStates` map with `data.IsCanonicalColumn()`/`data.IsCanonicalStage()` in `checks.go`
- [x] Move `colOverhead` constant from `column.go` to `layout.go`; update all references
- [x] Run `make build && make test`

## Context Log

**Files read:**
- `internal/board/column.go` — removed `taskLabel()`, `colOverhead` already in layout.go
- `internal/doctor/report.go` — removed `CheckResult` type
- `internal/doctor/repairs.go` — replaced `contains()`/`indexOf()` with `strings.Contains`/`strings.Index`, removed `exitCode` param from `GateSuggestion`
- `internal/buildtool/main.go` — replaced `trimSpace()` with `strings.TrimSpace`
- `internal/doctor/checks.go` — removed `validStates` map, removed router state validation
- `internal/board/layout.go` — verified `colOverhead` already co-located

**Files edited:**
- `internal/board/column.go`
- `internal/doctor/report.go`
- `internal/doctor/repairs.go`
- `internal/doctor/repairs_test.go`
- `internal/doctor/checks.go`
- `internal/doctor/checks_test.go`
- `internal/buildtool/main.go`

**Quality gates:** `make build` OK, `go test ./...` all pass

## Drift Notes

- `validStates` map was removed from `checks.go`. The AC specified using `data.IsCanonicalColumn()`/`data.IsCanonicalStage()` as replacement, but these validate `ColumnType` ("planned"/"in_progress"/"done") and `ProgressStage` ("build"/"test"/"audit") — neither covers router workflow states ("pre-implementation", "epic-design", etc.). Router state validation was removed entirely; these strings are deeply wired across the codebase (view.go, model.go, update.go, many tests) and cannot change in this task. The router state validation responsibility should be added to the `data` package in a future task if needed.
- `TestCheckRouterUnknownState` test removed since the router state check was removed from `CheckRouter`.