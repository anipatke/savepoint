---
id: E29-parallel-mode-guidance-and-hardening/T004-parallel-regression-matrix
status: planned
objective: Prove advanced parallel worktrees are optional, bounded, recoverable, and cross-platform before release audit.
depends_on:
  - E24-launcher-diagnostics-and-guidance/T004-launcher-regression-matrix
  - E29-parallel-mode-guidance-and-hardening/T003-parallel-workflow-documentation
complexity_tier: high
complexity_reason: Release verification spans config, task data, Git worktrees, launcher, board, doctor, init, and multiple platforms.
---

# T004: Parallel Regression Matrix

## Problem

Package tests must be combined into a release-level proof that optional parallel behavior cannot damage sequential workflows or Git state.

## Context Files

- `internal/data/config_test.go`
- `internal/data/task_test.go`
- `internal/data/parallel_test.go`
- `internal/worktree/service_test.go`
- `internal/worktree/manifest_test.go`
- `internal/board/integration_test.go`
- `internal/board/update_test.go`
- `internal/doctor/checks_test.go`
- `internal/init/integration_test.go`
- `Makefile`
- `README.md`

## Acceptance Criteria

- [ ] Tests prove absent and disabled advanced configuration preserve sequential parsing, board actions, lifecycle writes, and plain output.
- [ ] Enabled tests cover valid launch, overlap rejection, dirty checkout, branch/path collision, duplicate dispatch, restart reconciliation, and partial launch failure.
- [ ] Cleanup tests prove dirty or mismatched worktrees and all task branches are preserved.
- [ ] Cross-platform builds cover Git and launcher process boundaries without shell interpolation.
- [ ] A documented manual matrix covers create, launch, toggle off/on, integrate, and clean removal in a disposable repository.
- [ ] `make build && make test` passes before handoff to a fresh audit agent.

## Implementation Plan

- [ ] Add end-to-end fixtures for disabled, eligible, active, stale, and conflicting parallel groups.
- [ ] Exercise board commands with fake launchers and temporary real Git repositories.
- [ ] Add regression assertions for existing sequential launcher, navigation, lifecycle, and doctor behavior.
- [ ] Run host tests and all supported cross-platform builds.
- [ ] Record the disposable-repository manual matrix and outcomes in the Context Log.
- [ ] Run full quality gates and stop for fresh-session audit handoff.

## Context Log

Pending.
