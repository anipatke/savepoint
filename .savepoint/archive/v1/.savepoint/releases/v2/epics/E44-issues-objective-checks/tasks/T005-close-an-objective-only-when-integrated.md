---
id: E44-issues-objective-checks/T005-close-an-objective-only-when-integrated
title: Close an objective only when it is integrated
status: done
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

- [x] `ResolveObjectiveCompletion` returns a `GateDecision` using the existing blocker vocabulary; no parallel decision or blocker type is introduced.
- [x] Completion requires every Task in `index.ObjectiveTasks` for that Objective to be `done`; membership is read from Task ownership and never from a list maintained on the Objective.
- [x] An Objective with an unfinished Task is blocked with a reason naming which Task, and every unfinished Task is named, not only the first.
- [x] Integration clearance is resolved by the existing `ResolveClearance` over the Objective's `scope.kind: objective` Checks; missing, `needs_work`, stale, and unknown each block with the distinct existing blocker kind.
- [x] An Objective with all Tasks done and current integration clearance is allowed under checker authority.
- [x] An Objective declaring `owner_validation.required` additionally needs owner acceptance of that same current Check; acceptance bound to a superseded Check does not close it.
- [x] A recorded exception naming the Objective's latest Check grants completion by exception under owner authority, reported as allowed-by-exception with the exception attached and never as a CLEAR result or current clearance; an exception naming any other Check does not apply.
- [x] An Objective with no owned Tasks and no Check resolves to blocked with missing clearance rather than to allowed.
- [x] A Task's own completion decision is unchanged: no Objective Check closes a Task, and `ResolveTaskCompletion` returns identical decisions before and after this task.

## Implementation Plan

- [x] Add `objective_gate_v2.go` with `ResolveObjectiveCompletion`, reusing `ResolveClearance`, `GateDecision`, `GateBlocker`, and the existing clearance blocker kinds.
- [x] Reuse the E43 owner-acceptance and exception helpers rather than restating their rules; export or relocate them within `internal/data` only if reuse requires it.
- [x] Add a blocker reason for an Objective whose owned Tasks are not all done, naming each one.
- [x] Test every unfinished-Task, clearance-state, acceptance, superseded-acceptance, and exception combination, plus the empty-Objective case.
- [x] Add a regression case proving `ResolveTaskCompletion` decisions are unchanged.
- [x] Run the focused `internal/data` suite.

## Context Log

**Files read:** `objective_gate_v2.go` (new, none existed), `objective_gate_v2_test.go` (new), `gate_v2.go`, `gate_v2_test.go`, `objective_v2.go`, `project.go`, `check_v2.go`, `evidence_v2.go`, `dependency.go`, `dependency_test.go`.

**Files edited:**
- `internal/data/objective_gate_v2.go` (new) — `ResolveObjectiveCompletion(index, objectiveID) GateDecision`. Walks `index.ObjectiveTasks[objectiveID]`; any Task not `done` produces a `GateBlockInvalidState` blocker naming that Task and returns immediately (not excusable by exception, matching how an invalid Task state short-circuits `ResolveTaskCompletion` before its own exception check). Otherwise resolves integration clearance via the existing `ResolveClearance`, mirrors `ResolveTaskCompletion`'s clearance-state switch (missing/needs_work/stale/unknown/current+owner-acceptance) using the existing `GateBlockClearance*`, `GateBlockCheckerAuthority`, and `GateBlockOwnerAcceptance` kinds, and falls back to `applicableException` when blocked. No new `GateBlockKind` value or parallel decision type was introduced; `GateBlockInvalidState`'s doc comment was widened by one clause to cover this second use.
- `internal/data/gate_v2.go` — refactored `ownerValidationRequired`, `ownerAcceptedCheck`, and `applicableException` to take `*Evidence` instead of `*TaskV2`, so `ResolveObjectiveCompletion` reuses them unchanged rather than restating the rules. Updated the two `ResolveTaskCompletion` call sites accordingly and widened the `GateBlockInvalidState` doc comment.
- `internal/data/dependency.go` — updated `ResolveTaskDependencyV2`'s call to `ownerAcceptedCheck(target, ...)` to `ownerAcceptedCheck(target.Evidence, ...)` for the same signature change.

**Quality gates:**
- `go build ./...` — clean.
- `go test ./internal/data/... -run 'ResolveObjectiveCompletion|ResolveTaskCompletion|ResolveClearance|ResolveTaskStart|ResolveTaskDependencyV2' -v` — all pass, including the new `TestResolveObjectiveCompletion_*` cases and the `TestResolveTaskCompletion_unchangedByObjectiveGate` regression case.
- `make build && make test` — full suite green across all packages.

No `.savepoint/Health-Check.md` exists in this project, so the Quick check step was skipped per the skill's instruction that its absence is not a finding.

No drift: no new module/file outside the epic's documented `objective_gate_v2.go` target, and no architectural change beyond what E44's design already specifies.
