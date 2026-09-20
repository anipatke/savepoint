---
id: E43-task-check-gates/T005-decide-when-tasks-close
title: Decide when a task can close
status: done
objective: Return one canonical start and completion decision per Task with its actor authority and blocking reasons.
depends_on:
    - E43-task-check-gates/T004-tell-fresh-from-stale
complexity_tier: high
complexity_reason: Concentrates completion authority, acceptance, exceptions, and replan into one decision contract.
---

# T005: Decide when a task can close

## Problem

Completion eligibility has no single owner, so every future consumer would re-derive it from status fields. Nothing distinguishes a technical Task that a checker may close from one waiting on the owner, and nothing records that a close happened by exception.

## Context Files

- `internal/data/gate_v2.go`
- `internal/data/gate_v2_test.go`
- `internal/data/dependency.go`
- `internal/data/evidence_v2.go`
- `internal/data/task_v2.go`
- `internal/data/lifecycle.go`
- `internal/data/project.go`

## Acceptance Criteria

- [x] Start and completion each return one decision: whether it is allowed, which actor may act, and typed blockers naming every unmet requirement.
- [x] Starting requires no replan flag and all dependencies satisfied; each blocker names the specific dependency or flag.
- [x] A technical Task — no declared owner validation — may complete under checker authority when clearance is `current`.
- [x] A Task declaring `owner_validation.required` may not complete on clearance alone; it needs owner acceptance of the current Check, and the blocker says so.
- [x] Acceptance of a Check that a later Check supersedes does not close the Task until the owner accepts again.
- [x] Missing, stale, unknown, and needs-work clearance each block completion with a distinct reason.
- [x] A decision allowed only by a recorded exception is returned as allowed-by-exception with the exception attached, and never reports a CLEAR result or current clearance.
- [x] An exception applies only to the Check it names and does not carry to a superseding Check.
- [x] A replan flag preserves the Task's status, stage, and partial work, blocks start and advance, and introduces no fourth status value.
- [x] Stage advance through `build → test → audit` reuses the canonical lifecycle vocabulary in `internal/data`; no parallel status or stage vocabulary is added.
- [x] Inconsistent external edits are reported as named diagnostics without rewriting the record: a `done` Task without current clearance, an acceptance naming a superseded Check, and evidence contradicting the recorded status.
- [x] Inconsistency reporting returns every problem found, not only the first, and claims no authentication or prevention of external writes.
- [x] No V1 consumer behavior changes: `internal/board/transitions.go`, the router, and V1 task handling are untouched by this task.

## Implementation Plan

- [x] Define the decision and blocker types in `gate_v2.go` with the actor authority the decision grants.
- [x] Implement the start decision over replan state and T004's dependency answers.
- [x] Implement the completion decision over clearance, owner validation, acceptance binding, and exception handling.
- [x] Implement stage advance validation by reusing the existing canonical lifecycle validators rather than restating them.
- [x] Add evidence-inconsistency inspection returning all named diagnostics for a loaded project.
- [x] Test technical versus owner-validated completion, each blocking clearance state, superseded acceptance, allowed-by-exception, exception not carrying forward, replan blocking with preserved work, and each inconsistent external edit.
- [x] Add a test proving V1 board transition behavior is unchanged, then run the focused `internal/data` suite.

## Context Log

**Files read:** `internal/data/gate_v2.go`, `internal/data/gate_v2_test.go`, `internal/data/dependency.go`, `internal/data/dependency_test.go`, `internal/data/evidence_v2.go`, `internal/data/task_v2.go`, `internal/data/objective_v2.go`, `internal/data/lifecycle.go`, `internal/data/project.go`, `internal/data/check_v2.go`, `internal/data/errors.go`, `.savepoint/Guardrails.md`, `.savepoint/releases/v2/epics/E43-task-check-gates/E43-Detail.md`.

**Files edited:** `internal/data/gate_v2.go` (added `GateBlockKind`/`GateBlocker`/`GateDecision`, `ResolveTaskStart`, `ResolveTaskAdvance`, `ResolveTaskCompletion`, `ConsistencyDiagnosticKind`/`ConsistencyDiagnostic`/`InspectTaskConsistency`), `internal/data/gate_v2_test.go` (18 new test functions, one table-driven).

**Design notes:**
- `ResolveTaskAdvance` defers to `ResolveTaskCompletion` when a Task's stage is already `audit`, rather than restating completion inside the stage machine, and otherwise calls `AdvanceTaskLifecycleState` (lifecycle.go) to confirm the next stage move is canonical — one state machine, not two.
- An exception grants completion only when `Allowed` would otherwise be false and `exception.Check` equals the target's current latest Check (`index.LatestCheck[taskID]`), so it never overrides an already-CLEAR decision and never survives a superseding Check.
- `InspectTaskConsistency` reuses `ResolveClearance` and `ResolveTaskCompletion` rather than re-deriving clearance/completion rules, so the diagnostics and the gate decisions can never drift apart.

**Quality gates:**
- `go build ./...` — clean.
- `go vet ./...` — clean.
- `make build && make test` — all packages pass (`internal/data` 0.344s, `internal/doctor` 1.121s, full suite green).
- New tests: `TestResolveTaskStart_allowedWhenNoReplanAndDependenciesSatisfied`, `TestResolveTaskStart_blockedByReplan`, `TestResolveTaskStart_namesEveryUnsatisfiedDependency`, `TestResolveTaskAdvance_blockedByReplan`, `TestResolveTaskAdvance_movesThroughBuildAndTest`, `TestResolveTaskAdvance_fromAuditDefersToCompletionDecision`, `TestResolveTaskCompletion_technicalTaskAllowedUnderCheckerAuthorityWhenCurrent`, `TestResolveTaskCompletion_ownerValidationRequiredBlocksOnClearanceAlone`, `TestResolveTaskCompletion_ownerValidationSatisfiedWhenAcceptedCurrentCheck`, `TestResolveTaskCompletion_acceptanceOfSupersededCheckDoesNotClose`, `TestResolveTaskCompletion_eachClearanceStateBlocksWithADistinctReason` (missing/needs_work/unknown/stale subtests), `TestResolveTaskCompletion_allowedByExceptionWhenOtherwiseBlocked`, `TestResolveTaskCompletion_exceptionDoesNotCarryToASupersedingCheck`, `TestInspectTaskConsistency_reportsEveryProblemNotOnlyTheFirst`, `TestInspectTaskConsistency_reportsNothingForAConsistentProject`, `TestV1BoardTransitions_unaffectedByV2GateAdditions` (proves `AdvanceTaskLifecycle`/`RetreatTaskLifecycle`, the exact functions `internal/board/transitions.go` calls, are untouched — `internal/data`, `internal/board`, and `internal/doctor` were not edited).
- No `.savepoint/Health-Check.md` in this project; Quick check step skipped per skill instructions.

**Drift:** None — no new files/modules beyond the ones `E43-Detail.md`'s component table already names, and no architecture change beyond what the epic describes.
