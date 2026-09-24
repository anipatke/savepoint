---
id: T-032
title: Move the router to the next Task when the owner closes one
objective: O-014
planned_by: {role: planner, session: o014-design-20260924}
status: planned
complexity_tier: medium
complexity_reason: "Fixes a live bug where the p key writes release: none, and adds a router selection update to the board's completion writes with conflict handling."
depends_on: [{task: T-031, requires: clear}]
owner_validation: {required: false}
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

Pending execution.

## Drift Notes

Design section 4's "press `p`" wording is reconciled in T-034.
