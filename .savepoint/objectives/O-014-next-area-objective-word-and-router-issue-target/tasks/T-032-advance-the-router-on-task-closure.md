---
id: T-032
title: Move the router to the next Task when the owner closes one
objective: O-014
planned_by: {role: planner, session: o014-design-20260924}
status: done
complexity_tier: medium
complexity_reason: "Fixes a live bug where the p key writes release: none, and adds a router selection update to the board's completion writes with conflict handling."
depends_on: [{task: T-031, requires: clear}]
owner_validation:
    required: false
    accepted_check: ""
check_waiver:
    task: T-032
    reason: Waived - deferred to objective check
    actor:
        role: owner
        session: owner-chat
    recorded_at: "2026-09-24T08:02:25Z"
---

# T-032: Move the router to the next Task when the owner closes one

## Outcome

No board action ever blanks or changes the router's `release:` except the
Goal selector itself. Closing the selected Task on the board moves the
router to the Objective's next unfinished Task, so the Next line is already
right for the next session without an agent editing the router.

## User Check

Press `p` on a Task: router.md still says `release: R-006`. Close that Task
with Space: router.md now selects the Objective's next unfinished Task, keeps
`release: R-006`, and the board's Next line shows it. Close the last Task:
`task:` becomes `none` and Next reads `In Progress O-### · Check`.

## Done When

- Regression fixed: `selectionForTarget` (`actions.go`) currently builds
  `RouterSelectionV2{Objective, Task}` with an empty Release, which the
  writer records as `release: none`. Every board selection write now carries
  the router's current Release (and Issue) forward.
- Closing a Task on the board (Space to done, including the waiver path and
  exception completion) also sets router `task` to the lowest-ID Task in the
  same Objective that is not done, or `none` when all are done, when the
  router named the closed Task. Closing an Objective clears router
  `objective` and `task` when it named that Objective. A closure of a Task
  the router did not select leaves the router alone. `release:` and `issue:` are unchanged byte-for-byte.
- The router write happens after the record write; a router mtime conflict or
  write failure leaves the completed record in place and reports a clear
  message that the router still selects it (the T-030 warning then shows).
- Tests: `p` on Task and Objective preserves Release; each closure path
  advances or clears only the matching selection; a blocked next Task is
  still selected; a non-matching selection is left
  alone; router conflict path; `release:` bytes identical in every case.
- `git diff --check`, `make build && make test-fast` pass.

## Context Files

`internal/board/v2/actions.go`, `internal/board/v2/actions_test.go`,
`internal/board/v2/io.go`, `internal/board/v2/update.go`,
`internal/data/write.go`, `internal/data/write_test.go`.

## Design References

Design sections 4, 8, and 9.

## Guardrails

FS-01, FS-04, DATA-01, DATA-05, ARCH-02, TEST-03, TEST-05, TEST-08.

## Implementation Plan

1. Write a failing test for `p` blanking the Release; fix by reading the
   current router selection into the written selection.
2. Add a small helper that, given the fresh router and a closed record,
   returns the advanced/cleared selection or "no change"; the selection
   rule is a `data` function so it is not re-derived in the board.
3. Call it after each board completion write; handle conflicts as above.
4. Add tests.

## Boundaries

No automatic selection of the next Objective; that stays with
`savepoint-design`.
No change to who may close records.

## Technical Verification

Focused `./internal/board/v2/...` and data write tests during iteration;
`make build && make test-fast` at handoff.

## Technical Evidence

Extra reads before implementation:

- `.savepoint/objectives/O-014-next-area-objective-word-and-router-issue-target/tasks/T-031-let-the-router-name-an-issue.md` — checked the recorded state and owner waiver for T-032's `requires: clear` predecessor.
- `.savepoint/Guardrails.md` — read the Task's named guardrails and applicable style rules before code changes.
- Targeted `ResolveTaskDependencyV2` symbol search in `internal/data` — located the runtime gate needed to verify the predecessor without deciding from Task prose.
- `.savepoint/Design.md` Sections 4, 8, and 9 — checked the lifecycle ownership, board write boundaries, and mtime conflict approach named by the Task.
- `go doc` for `LoadV2Index`, `ResolveTaskStart`, `ResolveTaskDependencyV2`, `ResolveObjectiveDependency`, and their index/gate record types — identified the runtime APIs needed to verify the Task and Objective gates.
- A temporary `go run` helper loaded `.savepoint` through `LoadV2Index` and invoked the runtime start and dependency resolvers; it was removed after the preflight.
- `git status --short` — identified existing shared-worktree changes from T-030/T-031 and the router selection; these are pre-existing inputs to preserve while editing overlapping router-writer files.
- Targeted search outside the Context Files for the `writeValidProject`/`writeEvidenceProject` board-test fixtures — needed to build valid Release and Objective selection scenarios for the scoped action tests.
- `internal/board/v2/fixture_test.go` Release and Objective fixture helpers — needed to attach a declared R-006 to a selected Objective so the `p` action tests pass the normal target validation.
- `go doc` for `ReleaseV2` — checked the minimum record shape needed to make those R-006 action fixtures valid under `LoadV2Index`.
- The live R-006 `Release.md` record — used only as a schema reference for the temporary Release in action tests, not as a test fixture dependency.
- `internal/board/v2/releases_test.go` `TestSelectionWriteRejectsCrossReleaseObjective` — inspected after the focused gate exposed a regression caused by merging Release context inside `writeSelectionCmd`; the fix must retain its explicit cross-Release refusal while preserving an omitted Release.
- Targeted board-state/Next projection references — needed to assert that a router update changes the loaded board's Next line, as specified in T-032's User Check.
- Final `git diff` review for the scoped board/data files and `git status --short` — checked the implementation diff and separated T-032 changes from the existing T-030/T-031 worktree changes.

Preflight:

- `go run tmp_task032_preflight.go` — `ResolveTaskStart(T-032)` allowed; T-031's `requires: clear` dependency satisfied; O-014 is `planned` and its O-012 Objective dependency is satisfied. The temporary helper has been removed.

Focused iteration:

- `go test ./internal/board/v2 ./internal/data -count=1` — final run passed. Earlier iterations exposed test-fixture errors and a cross-Release selection regression; those were corrected before the passing run.

Acceptance evidence:

- The `p` action preserves the router's Release and Issue for both Task and Objective selection; `TestRecordSelectionPreservesReleaseAndIssueForTaskAndObjective` compares the entire router document with only the intended Task selection changed. `TestGoalSelectionPreservesIssueContext` covers the Goal selector's permitted Release change while keeping Issue bytes unchanged.
- Space completion advances a selected Task to the lowest-ID unfinished Task, even when that Task is blocked, and the loaded board Next line names it: `TestTaskCompletionMovesRouterToLowestUnfinishedBlockedTask`. The no-Check waiver path advances and preserves Release/Issue in `TestTaskWaiverCompletionAdvancesRouterAndKeepsReleaseAndIssue`.
- Exception completion advances only when the router selected that Task: `TestExceptionCompletionAdvancesOnlyWhenRouterSelectedClosedTask`. The last-Task and Objective closure cases clear only their matching selection fields and preserve the other context: `TestLastTaskCompletionClearsRouterTaskAndShowsObjectiveCheck` and `TestObjectiveExceptionCompletionClearsOnlyObjectiveAndTaskSelection`.
- The data selection rule selects the first unfinished Task without consulting its gate, clears `task` when all owned Tasks are done, clears a matching Objective and Task on Objective closure, and leaves nonmatching selections unchanged: `TestRouterSelectionAfterClosureV2_advancesToLowestUnfinishedEvenWhenBlocked` and `TestRouterSelectionAfterClosureV2_clearsMatchingSelectionsAndLeavesOthersAlone`.
- A router mtime conflict leaves the already-completed Task in place, retains the stale selection, requests a reload, and reports that the router still selects the Task: `TestRouterConflictAfterCompletionKeepsRecordAndReportsStaleSelection`.

Handoff gates:

- `make build && make test-fast` — passed on 2026-09-24 at 07:43:32 UTC with `go1.26.2 linux/amd64`.
- `git diff --check` — passed after the final code changes.
- The current runtime preflight was `go run tmp_task032_preflight.go`: T-032's start decision was allowed, T-031's `requires: clear` dependency was satisfied, and O-014's O-012 dependency was satisfied. The temporary helper was removed immediately afterward.
- No Full Objective Check was run; O-014 still requires current `make test-full` evidence before it can close. The owner waived the optional Task Check and deferred review to that Objective Check; this waiver is not technical CLEAR.

Files read:

- Router, active `savepoint-task` skill, O-014 Objective, and this Task record.
- Task Context Files: `internal/board/v2/actions.go`, `actions_test.go`, `io.go`, `update.go`, `internal/data/write.go`, and `write_test.go`.
- Extra reads recorded above: T-031's dependency record; `.savepoint/Guardrails.md`; `.savepoint/Design.md` Sections 4, 8, and 9; targeted `ResolveTaskDependencyV2` search and runtime API documentation; `internal/board/v2/fixture_test.go`; the live R-006 Release record; `internal/board/v2/releases_test.go`; board state/Next projection references in `load.go`, `next_panel.go`, and `next_panel_test.go`; and scoped `git status`/diff output.

Files changed for T-032:

- This Task record, `internal/board/v2/actions.go`, `internal/board/v2/actions_test.go`, `internal/board/v2/io.go`, `internal/board/v2/update.go`, `internal/data/write.go`, and `internal/data/write_test.go`.
- The shared worktree already contained T-030/T-031 changes. Existing Issue-selection changes in `internal/data/write.go` and `write_test.go`, and other T-030/T-031 files, were left intact. The temporary preflight helper is not present in the final worktree.

Limitations and handoff:

- The owner marked T-032 done and waived the optional Task Check, deferring review to O-014's mandatory Full Objective Check. The waiver records the owner's decision; it is not technical CLEAR, and no Task Check was written. The Full Objective Check remains required before Objective closure.

## Drift Notes

Design section 4's "press `p`" wording is reconciled in T-034.
