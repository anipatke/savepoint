---
id: T-009
title: Restore the missing O-900 Check example
objective: O-012
status: done
depends_on: []
owner_validation: {required: false}
planned_by: {role: planner, session: i018-remediation-20260922}
check_waiver:
  task: T-009
  reason: Owner requested O-012's mandatory Full Objective Check after reviewing the audit-ready remediation evidence.
  actor: {role: owner, session: user}
  recorded_at: '2026-09-22T09:58:47Z'
---

# T-009: Restore the missing O-900 Check example

## Outcome

The disposable O-900 visual spread genuinely renders every retained Task-card
review outcome, including pending `[ ] CHECK`, through both interactive card
and non-TTY paths while production Task-card behavior remains unchanged.

## User Check

No separate owner behavior check is required. The mandatory fresh Full O-012
Check will verify the live O-900 spread and I-018 repair after implementation.

## Done When

- One valid O-900 Task renders `[ ] CHECK` under the final card rules without
  changing production Task-card behavior or losing the 4/4/4 column spread.
- O-900 still renders all other retained outcomes plus wait, replan, and owner
  blockers through both interactive card and non-TTY paths.
- Durable fixture-level regression coverage proves the complete O-900 outcome
  and blocker matrix through both rendering paths.
- The stale audit-stage test comment and T-005's inaccurate technical evidence
  are corrected while T-003, T-004, and T-005 remain `done`.
- I-018 carries repair evidence for a fresh independent Check; I-019 remains
  open and is not implemented by this Task.

## Context Files

- `.savepoint/issues/I-018-o900-missing-pending-check-outcome.md`
- `.savepoint/objectives/O-012-simplify-task-card-outcomes/tasks/T-005-refresh-the-o900-visual-test-spread.md`
- `.savepoint/objectives/O-900-test-objective-for-ui-checks/Objective.md`
- `.savepoint/objectives/O-900-test-objective-for-ui-checks/tasks/T-900.md`
- `.savepoint/objectives/O-900-test-objective-for-ui-checks/tasks/T-901.md`
- `.savepoint/objectives/O-900-test-objective-for-ui-checks/tasks/T-902.md`
- `.savepoint/objectives/O-900-test-objective-for-ui-checks/tasks/T-903.md`
- `.savepoint/objectives/O-900-test-objective-for-ui-checks/tasks/T-904.md`
- `.savepoint/objectives/O-900-test-objective-for-ui-checks/tasks/T-905.md`
- `.savepoint/objectives/O-900-test-objective-for-ui-checks/tasks/T-906.md`
- `.savepoint/objectives/O-900-test-objective-for-ui-checks/tasks/T-907.md`
- `.savepoint/objectives/O-900-test-objective-for-ui-checks/tasks/T-908.md`
- `.savepoint/objectives/O-900-test-objective-for-ui-checks/tasks/T-909.md`
- `.savepoint/objectives/O-900-test-objective-for-ui-checks/tasks/T-910.md`
- `.savepoint/objectives/O-900-test-objective-for-ui-checks/tasks/T-911.md`
- `internal/board/v2/badges.go`
- `internal/board/v2/card.go`
- `internal/board/v2/card_test.go`
- `internal/board/v2/fixture_test.go`
- `internal/board/v2/plain.go`
- `internal/board/v2/run_test.go`

## Design References

Design sections 4, 8, and 13.

## Guardrails

DATA-02, DATA-03, ARCH-02, TEST-01, TEST-02, TEST-04, TEST-05, TEST-08,
STYLE-05, STYLE-07, and STYLE-10.

## Implementation Plan

1. Confirm the existing final card rule and identify a valid Done-column O-900
   record that can carry missing clearance without removing another outcome.
2. Adjust only disposable O-900 fixture records needed to restore the pending
   Check example while preserving all retained outcomes and blockers.
3. Add fixture-level regression coverage that loads an O-900-equivalent record
   set and asserts the full matrix through interactive card and non-TTY output.
4. Correct the stale audit-stage comment and T-005 evidence to match the final
   supported fixture behavior.
5. Record repair evidence on I-018 without closing it; leave I-019 unchanged.
6. Run focused board/data tests, `git diff --check`, `make build`, and
   `make test`, then hand O-012 to a fresh Full Objective Check.

## Boundaries

- Do not change production Task-card behavior or verification policy.
- Keep T-003, T-004, and T-005 `done`; this Task is the new remediation history.
- Do not edit immutable Check records or implement any part of I-019.
- O-900 remains disposable and remains a Release blocker until separately
  removed under its existing policy.

## Technical Verification

Focused `internal/board/v2` and `internal/data` tests, `git diff --check`,
`make build`, and `make test`; fixture-level assertions cover both card paths.

## Technical Evidence

**Resolved replan:** `internal/board/v2/plain_test.go`, originally named by this Task as
the non-TTY regression-test context, does not exist. The executor confirmed
all other declared Context Files exist before implementation and made no
fixture or production changes. Planning must identify the existing non-TTY
test seam or explicitly authorize a new test file, then return this Task to
execution. T-009 remains `status: in_progress`, `stage: build` as required for
a replan handoff.

Design reconciliation identified `internal/board/v2/fixture_test.go` as the
fixture-backed card seam and `internal/board/v2/run_test.go` as the existing
non-TTY seam. The Context Files now name those real files and execution may
resume without widening I-018.

**Per-criterion evidence:**

- Live O-900 now keeps a 4/4/4 spread with T-903 as the supported Done/no-Check
  `[ ] CHECK` example and T-909 as an in-progress build card retaining stale
  `[!] REVIEW`; T-907/T-905/T-910/T-908 retain CHECK/NEEDS WORK/WAIVED/OWNER
  ACCEPTED, and T-901/T-904/T-906 retain wait/replan/owner blockers.
- `TestO900OutcomeSpreadRendersEveryOutcomeAndBlockerOnCards` loads a complete
  temporary O-900-equivalent project and asserts all six outcomes, all three
  blockers, and the 4/4/4 spread through the interactive card renderer.
- `TestRunWithoutTTYRendersTheCompleteO900OutcomeSpread` asserts the same
  matrix and counts through the non-TTY `Run` path.
- Production files `card.go`, `badges.go`, and `plain.go` are unchanged.
- T-003-T-005 remain `done`; T-005's false T-903 evidence is corrected, and stale
  comments now distinguish stored `audit` from its displayed `CHECK` label.
- I-018 links T-009 and records repair evidence for independent recheck. I-019 is
  unchanged and remains open.

**Files changed:**

- `.savepoint/router.md`
- `.savepoint/issues/I-018-o900-missing-pending-check-outcome.md`
- `.savepoint/objectives/O-012-simplify-task-card-outcomes/tasks/T-005-refresh-the-o900-visual-test-spread.md`
- `.savepoint/objectives/O-012-simplify-task-card-outcomes/tasks/T-009-restore-the-o900-pending-check-example.md`
- `.savepoint/objectives/O-900-test-objective-for-ui-checks/tasks/T-903.md`
- `.savepoint/objectives/O-900-test-objective-for-ui-checks/tasks/T-909.md`
- `internal/board/v2/card_test.go`
- `internal/board/v2/fixture_test.go`
- `internal/board/v2/run_test.go`

**Extra reads:** global active Task-ID/frontmatter scan to correct the initial
duplicate T-006 allocation; `internal/data/project.go` to invoke the real strict
loader; the `internal/board/v2` test-file list to resolve the nonexistent
`plain_test.go` plan; and final `git status`/targeted diff to separate this
Task's changes from pre-existing workspace edits.

**Commands run:**

- Real `data.LoadProject(".savepoint")` probe after T-009 renaming — pass.
- Real `data.LoadProject(".savepoint")` probe after fixture repair — pass.
- `go test ./internal/board/v2 -run 'TestO900OutcomeSpread|TestRunWithoutTTYRendersTheCompleteO900OutcomeSpread'` — pass.
- `go test ./internal/board/v2/... ./internal/data/...` — pass.
- `git diff --check` — pass.
- `make build` — pass.
- `make test` — pass, including `internal/migrate` in 118.614s.

**Limitations:** No optional Task Check has been requested or waived, T-009 is
not owner-closed, and I-018 remains open until an independent Check verifies
the repair. I-019 is intentionally outside scope and remains open.

## Drift Notes

The originally planned non-TTY test file did not exist. Design reconciliation
replaced it with the existing fixture and run test seams before implementation.
