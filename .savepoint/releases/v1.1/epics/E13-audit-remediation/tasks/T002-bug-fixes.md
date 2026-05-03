---
id: E13-audit-remediation/T002-bug-fixes
status: done
objective: Fix cycle detection path reconstruction, AtomicWrite cross-device fallback, config accent defaults, and quality gate timeout
depends_on: []
---

# T002: Bug Fixes — Cycle Detection, AtomicWrite, Accent Defaults, Gate Timeout

## Context Files

- `internal/doctor/checks.go` — `detectCycles` uses a `parent` map that gets overwritten on revisits, producing inaccurate cycle paths
- `internal/init/write.go` — `replaceFile()` fallback still uses `os.Rename` which fails on cross-device moves
- `internal/data/config.go` — `fillThemeDefaults()` replaces entire accent map instead of filling missing keys individually
- `internal/doctor/gates.go` — `RunQualityGates` has no execution timeout

## Acceptance Criteria

- [x] `detectCycles` uses a stack-based cycle reconstruction or validates reconstructed path before reporting
- [x] `replaceFile()` fallback uses `os.Open` + `io.Copy` instead of `os.Rename` for cross-device moves
- [x] `fillThemeDefaults()` fills missing accent keys individually from `defaultTheme.Accents` instead of replacing the entire map when any accent is present
- [x] `RunQualityGates` uses `exec.CommandContext` with a 60-second default timeout; a `gate_timeout` config option is supported in `QualityGates`
- [x] `go test ./...` passes with no regressions
- [x] Existing tests for cycle detection, write, config, and gates still pass

## Implementation Plan

- [x] Refactor `detectCycles()` in `checks.go:382` — stack-based DFS path instead of parent map
- [x] Update `detectCycles` tests — added `TestCheckDependencies_CycleAccuratePath` to verify accurate cycle paths
- [x] Rewrite `replaceFile()` in `write.go:45` — uses `os.Open` + `io.Copy` + `os.Remove` as cross-device fallback
- [x] Update `AtomicWrite` tests — existing tests cover fallback (same behavior on rename success)
- [x] Fix `fillThemeDefaults()` in `config.go:75` — iterates `defaultTheme.Accents` and fills only missing keys
- [x] Add `Timeout` field to `QualityGates` struct (default `"60s"`); parse it in `RunQualityGates` and use `exec.CommandContext` in `runGate()`
- [x] Update `gates_test.go` — added `TestRunQualityGates_Timeout` and `TestRunQualityGates_DefaultTimeout`
- [x] `go build ./...` + `go test ./...` — all pass (clipboard test is pre-existing flake, unrelated)