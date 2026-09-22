---
id: T009
title: Restore the missing O900 Check example
objective: O012
status: done
depends_on: []
owner_validation: {required: false}
planned_by: {role: planner, session: i018-remediation-20260922}
check_waiver:
  task: T009
  reason: Owner requested O012's mandatory Full Objective Check after reviewing the audit-ready remediation evidence.
  actor: {role: owner, session: user}
  recorded_at: '2026-09-22T09:58:47Z'
---

# T009: Restore the missing O900 Check example

## Outcome

The disposable O900 visual spread genuinely renders every retained Task-card
review outcome, including pending `[ ] CHECK`, through both interactive card
and non-TTY paths while production Task-card behavior remains unchanged.

## User Check

No separate owner behavior check is required. The mandatory fresh Full O012
Check will verify the live O900 spread and I018 repair after implementation.

## Done When

- One valid O900 Task renders `[ ] CHECK` under the final card rules without
  changing production Task-card behavior or losing the 4/4/4 column spread.
- O900 still renders all other retained outcomes plus wait, replan, and owner
  blockers through both interactive card and non-TTY paths.
- Durable fixture-level regression coverage proves the complete O900 outcome
  and blocker matrix through both rendering paths.
- The stale audit-stage test comment and T005's inaccurate technical evidence
  are corrected while T003, T004, and T005 remain `done`.
- I018 carries repair evidence for a fresh independent Check; I019 remains
  open and is not implemented by this Task.

## Context Files

- `.savepoint/issues/I018-o900-missing-pending-check-outcome.md`
- `.savepoint/objectives/O012-simplify-task-card-outcomes/tasks/T005-refresh-the-o900-visual-test-spread.md`
- `.savepoint/objectives/O900-test-objective-for-ui-checks/Objective.md`
- `.savepoint/objectives/O900-test-objective-for-ui-checks/tasks/T900.md`
- `.savepoint/objectives/O900-test-objective-for-ui-checks/tasks/T901.md`
- `.savepoint/objectives/O900-test-objective-for-ui-checks/tasks/T902.md`
- `.savepoint/objectives/O900-test-objective-for-ui-checks/tasks/T903.md`
- `.savepoint/objectives/O900-test-objective-for-ui-checks/tasks/T904.md`
- `.savepoint/objectives/O900-test-objective-for-ui-checks/tasks/T905.md`
- `.savepoint/objectives/O900-test-objective-for-ui-checks/tasks/T906.md`
- `.savepoint/objectives/O900-test-objective-for-ui-checks/tasks/T907.md`
- `.savepoint/objectives/O900-test-objective-for-ui-checks/tasks/T908.md`
- `.savepoint/objectives/O900-test-objective-for-ui-checks/tasks/T909.md`
- `.savepoint/objectives/O900-test-objective-for-ui-checks/tasks/T910.md`
- `.savepoint/objectives/O900-test-objective-for-ui-checks/tasks/T911.md`
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

1. Confirm the existing final card rule and identify a valid Done-column O900
   record that can carry missing clearance without removing another outcome.
2. Adjust only disposable O900 fixture records needed to restore the pending
   Check example while preserving all retained outcomes and blockers.
3. Add fixture-level regression coverage that loads an O900-equivalent record
   set and asserts the full matrix through interactive card and non-TTY output.
4. Correct the stale audit-stage comment and T005 evidence to match the final
   supported fixture behavior.
5. Record repair evidence on I018 without closing it; leave I019 unchanged.
6. Run focused board/data tests, `git diff --check`, `make build`, and
   `make test`, then hand O012 to a fresh Full Objective Check.

## Boundaries

- Do not change production Task-card behavior or verification policy.
- Keep T003, T004, and T005 `done`; this Task is the new remediation history.
- Do not edit immutable Check records or implement any part of I019.
- O900 remains disposable and remains a Release blocker until separately
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
execution. T009 remains `status: in_progress`, `stage: build` as required for
a replan handoff.

Design reconciliation identified `internal/board/v2/fixture_test.go` as the
fixture-backed card seam and `internal/board/v2/run_test.go` as the existing
non-TTY seam. The Context Files now name those real files and execution may
resume without widening I018.

**Per-criterion evidence:**

- Live O900 now keeps a 4/4/4 spread with T903 as the supported Done/no-Check
  `[ ] CHECK` example and T909 as an in-progress build card retaining stale
  `[!] REVIEW`; T907/T905/T910/T908 retain CHECK/NEEDS WORK/WAIVED/OWNER
  ACCEPTED, and T901/T904/T906 retain wait/replan/owner blockers.
- `TestO900OutcomeSpreadRendersEveryOutcomeAndBlockerOnCards` loads a complete
  temporary O900-equivalent project and asserts all six outcomes, all three
  blockers, and the 4/4/4 spread through the interactive card renderer.
- `TestRunWithoutTTYRendersTheCompleteO900OutcomeSpread` asserts the same
  matrix and counts through the non-TTY `Run` path.
- Production files `card.go`, `badges.go`, and `plain.go` are unchanged.
- T003-T005 remain `done`; T005's false T903 evidence is corrected, and stale
  comments now distinguish stored `audit` from its displayed `CHECK` label.
- I018 links T009 and records repair evidence for independent recheck. I019 is
  unchanged and remains open.

**Files changed:**

- `.savepoint/router.md`
- `.savepoint/issues/I018-o900-missing-pending-check-outcome.md`
- `.savepoint/objectives/O012-simplify-task-card-outcomes/tasks/T005-refresh-the-o900-visual-test-spread.md`
- `.savepoint/objectives/O012-simplify-task-card-outcomes/tasks/T009-restore-the-o900-pending-check-example.md`
- `.savepoint/objectives/O900-test-objective-for-ui-checks/tasks/T903.md`
- `.savepoint/objectives/O900-test-objective-for-ui-checks/tasks/T909.md`
- `internal/board/v2/card_test.go`
- `internal/board/v2/fixture_test.go`
- `internal/board/v2/run_test.go`

**Extra reads:** global active Task-ID/frontmatter scan to correct the initial
duplicate T006 allocation; `internal/data/project.go` to invoke the real strict
loader; the `internal/board/v2` test-file list to resolve the nonexistent
`plain_test.go` plan; and final `git status`/targeted diff to separate this
Task's changes from pre-existing workspace edits.

**Commands run:**

- Real `data.LoadProject(".savepoint")` probe after T009 renaming — pass.
- Real `data.LoadProject(".savepoint")` probe after fixture repair — pass.
- `go test ./internal/board/v2 -run 'TestO900OutcomeSpread|TestRunWithoutTTYRendersTheCompleteO900OutcomeSpread'` — pass.
- `go test ./internal/board/v2/... ./internal/data/...` — pass.
- `git diff --check` — pass.
- `make build` — pass.
- `make test` — pass, including `internal/migrate` in 118.614s.

**Limitations:** No optional Task Check has been requested or waived, T009 is
not owner-closed, and I018 remains open until an independent Check verifies
the repair. I019 is intentionally outside scope and remains open.

## Drift Notes

The originally planned non-TTY test file did not exist. Design reconciliation
replaced it with the existing fixture and run test seams before implementation.
