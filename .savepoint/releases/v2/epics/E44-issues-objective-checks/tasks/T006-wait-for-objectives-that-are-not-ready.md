---
id: E44-issues-objective-checks/T006-wait-for-objectives-that-are-not-ready
title: Wait for objectives that are not ready
status: done
objective: Satisfy Objective dependencies only from integration clearance and block Task start when the owning Objective waits.
depends_on:
    - E44-issues-objective-checks/T005-close-an-objective-only-when-integrated
complexity_tier: medium
complexity_reason: Adds one readiness rule and extends an existing start decision with a new blocker kind.
---

# T006: Wait for objectives that are not ready

## Problem

Objective `depends_on` is validated for structure but never answered, and `ResolveTaskStart` only consults Task dependencies. Work in a dependent Objective could therefore start from Task-only clearance in the Objective it depends on, which is exactly what an integration Check exists to prevent — while an unrelated Objective must stay free to continue.

## Context Files

- `internal/data/objective_gate_v2.go`
- `internal/data/objective_gate_v2_test.go`
- `internal/data/gate_v2.go`
- `internal/data/gate_v2_test.go`
- `internal/data/dependency.go`
- `internal/data/dependency_test.go`
- `internal/data/project.go`
- `internal/data/objective_v2.go`

## Acceptance Criteria

- [x] An Objective dependency is satisfied only when the dependency Objective is `done` and its integration clearance is `current`.
- [x] A dependency Objective whose Tasks are all individually cleared but whose Objective-scoped Check is absent, stale, unknown, or `needs_work` does not satisfy the dependency, and each state is reported with its own reason.
- [x] A dependency Objective that is cleared by exception rather than by clearance is reported distinctly and does not silently read as current clearance.
- [x] `ResolveTaskStart` blocks when the owning Objective has an unsatisfied dependency, using a distinct blocker kind that identifies the waiting Objective and the dependency Objective, so a consumer can explain that the wait is at the Objective level.
- [x] Every unsatisfied Objective dependency is named, not only the first.
- [x] A Task whose owning Objective declares no dependencies, or whose dependency Objectives are all satisfied, starts exactly as it did before this task.
- [x] A Task in an Objective unrelated to a blocked Objective is unaffected: readiness consults only declared dependencies, never project-wide state.
- [x] Objective dependency evaluation re-resolves clearance from the index at decision time rather than reading a cached or stored judgement.

## Implementation Plan

- [x] Add `ResolveObjectiveDependency` to `objective_gate_v2.go`, returning a typed decision consistent with `ResolveTaskDependencyV2`'s shape.
- [x] Add the Objective-dependency blocker kind to `gate_v2.go` and extend `ResolveTaskStart` to evaluate the owning Objective's dependencies after its Task dependencies.
- [x] Keep `ResolveTaskAdvance` and `ResolveTaskCompletion` unchanged; only start consults Objective readiness.
- [x] Test each unsatisfied clearance state, the exception case, multiple unsatisfied dependencies, the satisfied path, the no-dependency path, and an unrelated Objective proceeding while another is blocked.
- [x] Add a regression case proving advance and completion decisions are unchanged.
- [x] Run the focused `internal/data` suite.

## Context Log

**Files read:** `objective_gate_v2.go`, `objective_gate_v2_test.go`, `gate_v2.go`, `gate_v2_test.go`, `dependency.go`, `dependency_test.go`, `project.go`, `objective_v2.go`.

**Files edited:**
- `internal/data/objective_gate_v2.go` — added `ObjectiveDependencyBlockKind` (`not_done`, `not_cleared`, `cleared_by_exception`), `ObjectiveDependencyBlock`, `ObjectiveDependencyDecision`, and `ResolveObjectiveDependency`, mirroring `ResolveTaskDependencyV2`'s shape and reusing `ResolveClearance` and `applicableException` rather than restating either.
- `internal/data/gate_v2.go` — added `GateBlockObjectiveDependency` and a `GateBlocker.ObjectiveDependency` field; extended `ResolveTaskStart` to evaluate the owning Objective's `DependsOn` after Task dependencies clear, naming every unsatisfied one.
- `internal/data/objective_gate_v2_test.go` — added `TestResolveObjectiveDependency_*` covering not-done, missing objective, satisfied, each clearance state (missing/needs_work/unknown/stale), cleared-by-exception, and an exception naming a superseded check.
- `internal/data/gate_v2_test.go` — added `TestResolveTaskStart_*` covering the blocked/satisfied/no-dependency/unrelated-objective paths and multiple unsatisfied dependencies, plus `TestResolveTaskAdvanceAndCompletion_unaffectedByObjectiveDependencyGate` as the required regression case.

**Quality gates:**
- `go build ./...` — clean.
- `go test ./internal/data/... -run 'ObjectiveDependency|ResolveTaskStart|ResolveTaskAdvance|ResolveObjectiveCompletion|ResolveTaskCompletion' -v` — all pass.
- `make build && make test` — all packages pass.
- No `.savepoint/Health-Check.md` in this project; Quick check step skipped per AGENTS.md (absence is not a finding).
- Guardrails STYLE (advisory): no duplicated logic — `ResolveObjectiveDependency` reuses `ResolveClearance`/`applicableException` rather than reimplementing clearance or exception rules (STYLE-07); typed block/decision shape mirrors the existing Task-dependency pattern (STYLE-04); every branch (not-done, each clearance state, exception, satisfied, no-dependency, unrelated-objective, multi-block) is covered by a test (STYLE-03).

## Drift Notes

None. No new files or modules outside the epic's documented `objective_gate_v2.go`/`gate_v2.go` scope, and no architectural change beyond what `E44-Detail.md` describes.
