---
id: E44-issues-objective-checks/T006-wait-for-objectives-that-are-not-ready
title: Wait for objectives that are not ready
status: planned
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

- [ ] An Objective dependency is satisfied only when the dependency Objective is `done` and its integration clearance is `current`.
- [ ] A dependency Objective whose Tasks are all individually cleared but whose Objective-scoped Check is absent, stale, unknown, or `needs_work` does not satisfy the dependency, and each state is reported with its own reason.
- [ ] A dependency Objective that is cleared by exception rather than by clearance is reported distinctly and does not silently read as current clearance.
- [ ] `ResolveTaskStart` blocks when the owning Objective has an unsatisfied dependency, using a distinct blocker kind that identifies the waiting Objective and the dependency Objective, so a consumer can explain that the wait is at the Objective level.
- [ ] Every unsatisfied Objective dependency is named, not only the first.
- [ ] A Task whose owning Objective declares no dependencies, or whose dependency Objectives are all satisfied, starts exactly as it did before this task.
- [ ] A Task in an Objective unrelated to a blocked Objective is unaffected: readiness consults only declared dependencies, never project-wide state.
- [ ] Objective dependency evaluation re-resolves clearance from the index at decision time rather than reading a cached or stored judgement.

## Implementation Plan

- [ ] Add `ResolveObjectiveDependency` to `objective_gate_v2.go`, returning a typed decision consistent with `ResolveTaskDependencyV2`'s shape.
- [ ] Add the Objective-dependency blocker kind to `gate_v2.go` and extend `ResolveTaskStart` to evaluate the owning Objective's dependencies after its Task dependencies.
- [ ] Keep `ResolveTaskAdvance` and `ResolveTaskCompletion` unchanged; only start consults Objective readiness.
- [ ] Test each unsatisfied clearance state, the exception case, multiple unsatisfied dependencies, the satisfied path, the no-dependency path, and an unrelated Objective proceeding while another is blocked.
- [ ] Add a regression case proving advance and completion decisions are unchanged.
- [ ] Run the focused `internal/data` suite.

## Context Log

Pending.
