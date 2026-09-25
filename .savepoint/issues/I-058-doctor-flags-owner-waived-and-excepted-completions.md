---
id: I-058
title: Doctor flags Tasks and Objectives the owner completed by waiver or exception
type: defect
status: resolved
source:
  kind: report
  actor: {role: executor, session: doctor-rca-20260925}
  at: '2026-09-25T11:01:14Z'
severity: medium
history:
  - at: '2026-09-25T11:01:14Z'
    actor: {role: executor, session: doctor-rca-20260925}
    kind: observed
    note: >-
      Live savepoint doctor exits 1 with 108 errors on this repository. 98 are
      false positives: the doctor's consistency checks reject completions that
      the completion gates allow by owner waiver or exception.
  - at: '2026-09-25T11:02:54Z'
    actor: {role: executor, session: doctor-rca-20260925}
    kind: repair_attempted
    note: >-
      InspectTaskConsistency and InspectObjectiveConsistency now accept the
      owner waiver and exception that the completion gates accept, and the
      acceptance-superseded check skips Tasks with no Task Check. New
      regression tests fail without the fix. make build && make test-fast
      passed. Live doctor now reports only the I-050..I-053 proof and R-006
      Goal Check items. Issue remains open for independent verification.
  - at: "2026-09-25T21:14:52Z"
    actor:
      role: owner
      session: board-owner
    kind: owner_decision
    note: Moved from open to in_progress by the owner from the board.
  - at: "2026-09-25T21:14:53Z"
    actor:
      role: owner
      session: board-owner
    kind: owner_decision
    note: Moved from in_progress to open by the owner from the board.
  - at: "2026-09-25T21:18:12Z"
    actor:
      role: owner
      session: board-owner
    kind: owner_decision
    note: Moved from open to in_progress by the owner from the board.
  - at: "2026-09-25T21:18:13Z"
    actor:
      role: owner
      session: board-owner
    kind: owner_decision
    note: Resolved by the owner from the board.
resolution:
  disposition: accepted
  actor:
    role: owner
    session: board-owner
  at: "2026-09-25T21:18:13Z"
  reason: Resolved by the owner from the board.
---

# I-058: Doctor flags Tasks and Objectives the owner completed by waiver or exception

## Summary

`savepoint doctor` checks done work with a stricter rule than the
completion gates. `ResolveTaskCompletion` lets a Task complete with no Task
Check when the owner recorded a `check_waiver`, or by an `exception` naming
its latest Check. `ResolveObjectiveCompletion` lets an Objective complete by
an `exception`. The doctor's `InspectTaskConsistency` and
`InspectObjectiveConsistency` accept only current clearance, so every
Task the owner closes on the board without a Task Check becomes a doctor
error.

A second check, `v2-acceptance-superseded`, fires when a waived Task's
`owner_validation.accepted_check` names the Objective Check. The Task has no
Task Check, so nothing superseded the acceptance, and the message ends with a
blank Check ID.

## Evidence

- `savepoint doctor` at a1b2b9f: exit 1. Diagnostic counts (each is listed
  once under Project Check and once in the Health Summary):
  - `v2-done-without-clearance`: 44 Tasks, all with clearance `missing` and
    a `check_waiver`.
  - `v2-objective-done-without-clearance`: O-001 (C-905) and O-013 (C-911),
    each with an owner `exception` naming that NEEDS WORK Check.
  - `v2-acceptance-superseded`: T-020, T-044, T-046, for example
    `owner accepted C-919 but the latest check is ` (C-919, C-929, and C-932
    are Objective Checks).
- `internal/data/gate_v2.go` `InspectTaskConsistency` and
  `internal/data/objective_gate_v2.go` `InspectObjectiveConsistency` compare
  only `ResolveClearance(...).State` against `ClearanceCurrent`.
- The doctor check dates from 032095a (2026-09-15). Waivers joined the
  completion gate on 2026-09-20 (a4e3db6, 3cdf717); the doctor was never
  updated to match.
- Out of scope and correct under the current rules: `v2-issue-verified-proof-superseded`
  for I-050..I-053 (C-925 superseded their proof C-924), and
  `v2-release-clearance-missing` for R-006 (Goal Check still pending).

## Repair Attempt

- `internal/data/gate_v2.go`: new `taskDoneByOwnerDecision` applies
  `ResolveTaskCompletion`'s owner rules (a waiver while no Check exists, or an
  exception naming the latest Check). `InspectTaskConsistency` uses it before
  reporting `ConsistencyDoneWithoutClearance`. The acceptance check now
  fires only when the Task has a latest Task Check that differs from the
  accepted one.
- `internal/data/objective_gate_v2.go`: `InspectObjectiveConsistency` skips
  `ObjectiveConsistencyDoneWithoutClearance` when `applicableException`
  names the latest Objective Check.
- Tests: `TestInspectTaskConsistency_ownerDecisionCompletionsAreNotReported`,
  `TestInspectTaskConsistency_waiverDoesNotExcuseARecordedCheck`, and
  `TestInspectObjectiveConsistency_exceptionCompletionIsNotReported` (which
  also checks that an exception naming a superseded Check is still flagged).
  Both new "not reported" tests fail against the unfixed code.
- `make build && make test-fast` passed (exit 0, go1.26.2 linux/amd64).
- Live `savepoint doctor` after the fix reports only 8
  `v2-issue-verified-proof-superseded` (I-050..I-053, listed twice) and 2
  `v2-release-clearance-missing` (R-006, listed twice).
- Limitation: the doctor still does not check a done Task's
  `owner_validation` against a current Task Check; it never did, and this
  repair does not add that.

## Proof Needed

Regression tests showing that the doctor reports nothing for: a done Task with
an applicable waiver and no Check; a done Task or Objective with an exception
naming its latest Check; and a waived Task whose `owner_validation` names a
non-Task Check. Existing tests for real contradictions must still pass. Live
`savepoint doctor` on this repository must report none of the three
false-positive families.
