---
id: T-121
title: Honour still-valid owner decisions when deciding if work can close
objective: O-044
status: done
depends_on: [{task: T-120, requires: clear}]
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: planning-o044-20261010}
check_waiver:
    task: T-121
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-10T06:13:56Z"
---

# Honour still-valid owner decisions when deciding if work can close

## Outcome

The completion gates treat an owner acceptance or exception as applying at the latest Check when it was recorded against that Check, a checker assessed it as still applying there, or the owner renewed it there. A decision with no assessment at the latest Check is reported as needing assessment; one assessed as changed is reported as needing renewal for that decision only, naming the change. Clearance stays NEEDS WORK while evidence is waived, and an exception only grants completion when it covers every unmet requirement the latest Check lists.

## User Check

Run the TheShed O-002 regression test and read its fixture: the Objective with C-005 → C-006 and the decisions carried forward resolves as completion allowed by exception, while its clearance still reads NEEDS WORK.

## Done When

1. One resolver in `internal/data` returns, for a decision at a record's latest Check: applies (origin, checker `applies: true`, or owner renewal), unassessed (names the latest Check), or materially changed (names the change). An owner renewal at a Check outranks a checker `applies: false` there.
2. `ResolveTaskCompletion`, `ResolveObjectiveCompletion`, `ResolveTaskDependencyV2` (`requires: accepted`), `ResolveObjectiveDependency`, `taskDoneByOwnerDecision`, `InspectTaskConsistency` and `InspectObjectiveConsistency` use it in place of exact-Check equality; no other copy of the rule remains (STYLE-07).
3. New typed blockers: decision unassessed, decision materially changed (carrying the change and which decision), and exception scope (carrying the uncovered requirement IDs). `GateDecision` exposes the deciding carry entry when completion is allowed by exception.
4. Regression (TEST-05), TheShed O-002 shape: done Tasks; owner acceptance and an evidence exception recorded against C-005; C-005 NEEDS WORK with `unmet` covering Design drift and the waived evidence; C-006 supersedes it, NEEDS WORK with `unmet` the waived evidence only, and checker entries carry both decisions forward. Completion is allowed by exception with owner authority; clearance is NEEDS WORK, never current. The same fixture without the carry entries reports unassessed at C-006, and before this change reported a plain NEEDS WORK block (recorded as the failing scenario).
5. A material change to accepted behavior blocks renewal of that acceptance only while another decision in the same index still applies; a material change to an exception's scope (checker `applies: false`, or an `unmet` ID outside `requirements`) blocks that exception only and names the change or IDs.
6. Legacy records with no new fields resolve exactly as before when the decision names the latest Check; prior Check files are never read for writing. No owner decision is created and no status changes.
7. `make build && make test-fast` passes (TEST-08); per-criterion evidence follows AGENTS.md's Verification Policy, including an explicit owner waiver if the optional Task Check is skipped.

## Context Files

`internal/data/gate_v2.go`, `internal/data/gate_v2_test.go`, `internal/data/objective_gate_v2.go`, `internal/data/objective_gate_v2_test.go`, `internal/data/dependency.go`, `internal/data/dependency_test.go`, `internal/data/evidence_v2.go`, `internal/data/check_v2.go`, `internal/testutil/`, `.savepoint/objectives/O-044-keep-owner-decisions-across-rechecks/Objective.md`, `.savepoint/issues/I-141-recheck-invalidates-owner-acceptance-and-exceptions.md`.

## Design References

O-044 Confirmed Design (all bullets). I-141 Evidence (TheShed C-003 to C-006) and Root cause 1–2.

## Guardrails

DATA-02, STYLE-04, STYLE-07, TEST-01, TEST-02, TEST-04, TEST-05, TEST-08.

## Implementation Plan

1. Confirm T-120's fields exist as specified; return REPLAN REQUIRED if not.
2. Add the applicability resolver beside the existing gate helpers and replace `ownerAcceptedCheck` and `applicableException` call sites.
3. Add the three blocker kinds and the exception-scope check against the latest Check's `unmet`.
4. Build the TheShed-shaped fixture in a temporary directory (TEST-04), not from the real TheShed files.
5. Add the regression, material-change, renewal and legacy tests.

## Boundaries

Gates only: no Next kinds, resume or board wording, or skill changes; later O-044 Tasks own those. Do not weaken the mandatory Objective Check, treat waived evidence as CLEAR, or let an exception excuse an unfinished Task or an open material Issue.

## Technical Verification

Focused `make test-focused TEST=...` during iteration; `make build && make test-fast` for handoff. The Full Objective Check needs fresh `make test-full` under `agent-skills/references/check-method.md`.

## Technical Evidence

Executor session, 2026-10-10. Stage advanced build → test → audit; ready for a Check, not passed. T-120's dependency (`requires: clear`) was met by its recorded owner Task-check waiver; `savepoint resume` reported "Ready" before start. T-120's fields were confirmed present and used as specified, so no REPLAN REQUIRED.

**Per-criterion outcomes**

1. Met. `assessDecision` (`internal/data/decision_gate_v2.go`) is the single rule: applies (origin == latest, checker `applies: true`, or owner renewal), unassessed (names the latest Check), changed (carries the change). An owner entry is checked first so it outranks a checker `applies: false` at the same Check. `TestAssessDecision` (8 cases).
2. Met. `ResolveTaskCompletion`, `ResolveObjectiveCompletion`, `ResolveTaskDependencyV2` (`requires: accepted`), `ResolveObjectiveDependency`, `taskDoneByOwnerDecision`, `InspectTaskConsistency` and `InspectObjectiveConsistency` all go through it. `ownerAcceptedCheck` and `applicableException` are deleted; grep finds no remaining use.
3. Met. New blockers `GateBlockDecisionUnassessed`, `GateBlockDecisionChanged` (with `Decision`, `Change`) and `GateBlockExceptionScope` (with `Requirements`); `GateDecision.ExceptionCarry` is set when completion is allowed by exception and the deciding entry exists. Tests: `TestResolveObjectiveCompletion_unassessedExceptionNamesLatestCheck`, `_changedExceptionBlocksOnlyTheException`, `_exceptionMustCoverEveryUnmetRequirement`, `_exceptionCarriedForwardGrantsCompletion`.
4. Met. `TestTheShedO002_decisionsCarriedForwardAllowCompletionByException` loads a temporary project (not TheShed's files): done Tasks, acceptance and exception against C-005, C-005 NEEDS WORK `unmet [DESIGN-01, TEST-08]`, C-006 supersedes with `unmet [TEST-08]` and checker carries. Completion is allowed by exception, owner authority, carry C-006, clearance NEEDS WORK on C-006. `TestTheShedO002_withoutCarryEntriesReportsUnassessed` shows the same records without carries report unassessed at C-006. The failing-before scenario (a plain NEEDS WORK block, no way to carry a decision) is shown by the existing tests that had to change; I did not run the new tests against the old code.
5. Met. A checker `applies: false` on the exception blocks only the exception and names the change, while the acceptance still applies; an `unmet` ID outside `requirements` gives an exception_scope blocker naming the IDs; an owner renewal restores a changed acceptance (`TestResolveTaskCompletion_changedAcceptanceNamesTheChange`). One reading to confirm: I took "blocks renewal ... only while another decision still applies" to mean each decision is judged on its own.
6. Met. Decisions naming the latest Check resolve as before. No file was written other than test fixtures, no owner decision created, no status changed by code.
7. Met. `make build` exit 0, `make test-fast` exit 0 (go1.26.2 linux/amd64, 2026-10-10T06:13Z). No Task Check requested and no owner waiver recorded for T-121.

**Behaviour change to existing tests (five):** a decision recorded against a superseded Check used to be ignored, giving a plain `owner_acceptance_required` or `clearance_needs_work` block. It now gives the new unassessed blocker (and for exceptions it is added after the clearance blocker). Updated: `TestE43_EpicScenario` step, `TestResolveTaskCompletion_acceptanceOfSupersededCheckDoesNotClose`, `TestResolveTaskCompletion_exceptionDoesNotCarryToASupersedingCheck`, `TestResolveObjectiveCompletion_acceptanceOfSupersededCheckDoesNotClose`, `TestResolveObjectiveCompletion_exceptionNamingOtherCheckDoesNotApply`. All still assert not allowed.

**Files read:** `gate_v2.go`, `objective_gate_v2.go`, `dependency.go` (parts), plus test helpers in `gate_v2_test.go`, `objective_gate_v2_test.go`, `dependency_test.go`, `e2e_v2_test.go` (to update the five tests). Extra reads: none beyond the Context Files and the failing-test excerpts above.

**Files changed:** `internal/data/decision_gate_v2.go` (new), `internal/data/decision_gate_v2_test.go` (new), `gate_v2.go`, `objective_gate_v2.go`, `dependency.go`, and the four existing test files named above; this Task's status.

**Limitations:** `make test-full` not run. Resume, board and Next still treat these blockers as generic; wording and the Assess step belong to the next Task. `ResolveTaskCompletion` for a Task with a recorded exception but no Check at all is unchanged (no blocker added). Release-level evidence is not covered by this rule.

## Drift Notes

None yet.
