---
id: E44-issues-objective-checks/T005-close-an-objective-only-when-integrated
title: Close an objective only when it is integrated
status: planned
objective: Decide Objective completion from owned Task completion plus current integration clearance, acceptance, or exception.
depends_on: []
complexity_tier: high
complexity_reason: Concentrates integration authority, acceptance, and exceptions into one new decision contract.
---

# T005: Close an objective only when it is integrated

## Problem

Nothing decides whether an Objective is finished. Its Tasks each carry their own clearance, but no control asks whether the assembled outcome was independently checked, so an Objective could be treated as complete from Task-only evidence and its integration Check would mean nothing.

## Context Files

- `internal/data/objective_gate_v2.go`
- `internal/data/objective_gate_v2_test.go`
- `internal/data/gate_v2.go`
- `internal/data/gate_v2_test.go`
- `internal/data/objective_v2.go`
- `internal/data/project.go`
- `internal/data/check_v2.go`
- `internal/data/evidence_v2.go`

## Acceptance Criteria

- [ ] `ResolveObjectiveCompletion` returns a `GateDecision` using the existing blocker vocabulary; no parallel decision or blocker type is introduced.
- [ ] Completion requires every Task in `index.ObjectiveTasks` for that Objective to be `done`; membership is read from Task ownership and never from a list maintained on the Objective.
- [ ] An Objective with an unfinished Task is blocked with a reason naming which Task, and every unfinished Task is named, not only the first.
- [ ] Integration clearance is resolved by the existing `ResolveClearance` over the Objective's `scope.kind: objective` Checks; missing, `needs_work`, stale, and unknown each block with the distinct existing blocker kind.
- [ ] An Objective with all Tasks done and current integration clearance is allowed under checker authority.
- [ ] An Objective declaring `owner_validation.required` additionally needs owner acceptance of that same current Check; acceptance bound to a superseded Check does not close it.
- [ ] A recorded exception naming the Objective's latest Check grants completion by exception under owner authority, reported as allowed-by-exception with the exception attached and never as a CLEAR result or current clearance; an exception naming any other Check does not apply.
- [ ] An Objective with no owned Tasks and no Check resolves to blocked with missing clearance rather than to allowed.
- [ ] A Task's own completion decision is unchanged: no Objective Check closes a Task, and `ResolveTaskCompletion` returns identical decisions before and after this task.

## Implementation Plan

- [ ] Add `objective_gate_v2.go` with `ResolveObjectiveCompletion`, reusing `ResolveClearance`, `GateDecision`, `GateBlocker`, and the existing clearance blocker kinds.
- [ ] Reuse the E43 owner-acceptance and exception helpers rather than restating their rules; export or relocate them within `internal/data` only if reuse requires it.
- [ ] Add a blocker reason for an Objective whose owned Tasks are not all done, naming each one.
- [ ] Test every unfinished-Task, clearance-state, acceptance, superseded-acceptance, and exception combination, plus the empty-Objective case.
- [ ] Add a regression case proving `ResolveTaskCompletion` decisions are unchanged.
- [ ] Run the focused `internal/data` suite.

## Context Log

Pending.
