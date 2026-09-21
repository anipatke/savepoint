---
id: T005
title: Refresh the O900 visual test spread
objective: O012
planned_by: {role: planner, session: task-card-simplification-20260921}
status: done
depends_on: [{task: T003, requires: clear}]
owner_validation:
    required: true
    accepted_check: ""
check_waiver:
    task: T005
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-21T11:35:50Z"
---

# T005: Refresh the O900 visual test spread

## Outcome

The disposable O900 Objective gives the owner one place to inspect every
retained Task-card review outcome and each independent blocker after the card
simplification lands.

## User Check

Open O900 in the V2 board and confirm its cards visibly cover `CHECK`, checked
`CHECK`, `NEEDS WORK`, `REVIEW`, `WAIVED`, and `OWNER ACCEPTED`, alongside
build/test/audit stages and wait/replan/owner blockers. No retired stale,
unverified, exception, or completion badge should appear. O900 must remain
explicitly disposable and required to be removed before R006 closes.

## Done When

- O900's Objective text identifies it as the visual verification spread for
  the simplified Task-card outcomes.
- Its twelve Tasks preserve their existing IDs, O900 ownership, and useful
  planned/in-progress/done column spread.
- The fixture provides an unambiguous example for each outcome: T903 has no
  Check, T907 has a current clear Check, T905 has a current needs-work Check,
  T909 needs review, T910 has an owner waiver, and T908 is owner accepted.
- T901 still demonstrates a wait blocker, T904 a replan blocker, and T906 an
  owner blocker, independently from the outcome badge.
- Fixture titles and descriptions use the retained vocabulary rather than
  presenting stale, exception, or completion terminology as card labels.
- The refreshed records pass strict V2 loading and render through both board
  card and non-TTY summary paths.

## Context Files

- `.savepoint/objectives/O012-simplify-task-card-outcomes/tasks/T003-give-every-task-card-one-clear-review-outcome.md`
- `.savepoint/objectives/O900-v2-live-ui-test-data/Objective.md`
- `.savepoint/objectives/O900-v2-live-ui-test-data/tasks/T900-planned-ready.md`
- `.savepoint/objectives/O900-v2-live-ui-test-data/tasks/T901-planned-waiting.md`
- `.savepoint/objectives/O900-v2-live-ui-test-data/tasks/T902-planned-independent.md`
- `.savepoint/objectives/O900-v2-live-ui-test-data/tasks/T903-building.md`
- `.savepoint/objectives/O900-v2-live-ui-test-data/tasks/T904-testing-replan.md`
- `.savepoint/objectives/O900-v2-live-ui-test-data/tasks/T905-audit-needs-work.md`
- `.savepoint/objectives/O900-v2-live-ui-test-data/tasks/T906-audit-clear-owner.md`
- `.savepoint/objectives/O900-v2-live-ui-test-data/tasks/T907-done-clear.md`
- `.savepoint/objectives/O900-v2-live-ui-test-data/tasks/T908-done-exception.md`
- `.savepoint/objectives/O900-v2-live-ui-test-data/tasks/T909-done-stale.md`
- `.savepoint/objectives/O900-v2-live-ui-test-data/tasks/T910-done-unverified.md`
- `.savepoint/objectives/O900-v2-live-ui-test-data/tasks/T911-planned-spare.md`
- `.savepoint/checks/C900-o900-needs-work.md`
- `.savepoint/checks/C901-o900-clear-owner.md`
- `.savepoint/checks/C902-o900-done-clear.md`
- `.savepoint/checks/C903-o900-owner-accepted.md`
- `.savepoint/checks/C904-o900-stale.md`
- `internal/board/v2/card.go`
- `internal/board/v2/badges.go`
- `internal/board/v2/plain.go`

## Design References

Design sections 4, 7, 8, and 13.

## Guardrails

DATA-02, DATA-03, ARCH-02, TEST-01, TEST-02, TEST-04, TEST-08, STYLE-05,
STYLE-07, and STYLE-10.

## Implementation Plan

1. Read T003's implemented outcome mapping and assign each O900 Task exactly
   one intended visual role without changing its lifecycle column.
2. Update O900's body and fixture titles/descriptions to explain the retained
   outcome and blocker coverage in current vocabulary.
3. Add a valid owner `check_waiver` to T910, retain T908's owner-acceptance
   exception evidence, and do not edit C900-C904.
4. Preserve valid ownership, dependencies, stages, Check links, and evidence
   while removing fixture wording that advertises retired card badges.
5. Load and render O900 through the existing strict data and board paths,
   checking interactive and non-TTY output against the User Check.
6. Run the configured quality gates and record per-criterion evidence.

## Boundaries

- This Task changes only disposable O900 fixture records and tests directly
  needed to keep that fixture valid; production card logic belongs to T003.
- Check records C900-C904 are immutable and must not be edited.
- O900 remains a deliberate R006 blocker until it is separately deleted; this
  Task does not remove it or weaken that release guard.

## Technical Verification

Focused V2 data/board tests covering O900 loading and rendered labels,
`git diff --check`, `make build`, and `make test`.

## Technical Evidence

**Context File naming drift (extra reads, logged):** The Objective and Task
paths listed above under `O900-v2-live-ui-test-data` and the descriptive
`T900-planned-ready.md`-style filenames do not exist on disk. The live O900
fixture is at `.savepoint/objectives/O900-test-objective-for-ui-checks/`, with
plain `T900.md`…`T911.md` task filenames and plain `C900.md`…`C904.md` check
filenames (no descriptive suffixes). Confirmed by directory listing before
editing; content and roles otherwise match this Task's description exactly
(same 12 Task IDs, same C900-C904 scope/result values), so this is read
naming drift in the Task doc, not a materially invalid plan — resolved by
locating the actual paths rather than a REPLAN.

**Files changed:**

- `.savepoint/objectives/O900-test-objective-for-ui-checks/Objective.md` —
  retitled and reworded to identify O900 as the visual spread for O012's
  simplified Task-card vocabulary (`CHECK`, `NEEDS WORK`, `REVIEW`, `WAIVED`,
  `OWNER ACCEPTED`, plus wait/replan/owner blockers), keeping the existing
  disposable/must-delete-before-R006 language unchanged.
- `tasks/T905.md` — title reworded from "check found problems" to "Check
  needs work" to name the retained outcome directly.
- `tasks/T908.md` — title reworded from "done by owner exception" (echoed the
  retired `BY EXCEPTION` card label) to "done, owner accepted" (the retained
  `OWNER ACCEPTED` outcome). No frontmatter change; its `exception` block was
  already valid evidence.
- `tasks/T909.md` — title reworded from "done, but clearance has gone stale"
  to "done, needs review" (stale collapses to `REVIEW` on the card per O012).
  No frontmatter change; its stale `freshness` block was already valid.
- `tasks/T910.md` — title reworded to "done, owner waived the Check"; added a
  valid `check_waiver` block (`task: T910`, owner actor, reason, RFC 3339
  `recorded_at`) so it demonstrates the `WAIVED` outcome, per the Done When
  requirement that T910 carry an owner waiver.
- Not edited: `T900`-`T904`, `T906`, `T907`, `T911`, and `C900`-`C904` (all
  already matched their intended role and used no retired vocabulary;
  C900-C904 are immutable per Boundaries and were not touched).

**Per-criterion evidence (Done When):**

- O900's Objective text identifies it as the visual verification spread for
  the simplified Task-card outcomes: see the reworded `Objective.md` body
  above, naming every retained outcome by exact label.
- Twelve Tasks preserve IDs/ownership/column spread: verified by loading the
  real project (`data.LoadProject(".savepoint")`) and grouping O900's cards —
  4 planned (T900-T902, T911), 4 in progress (T903-T906), 4 done
  (T907-T910), all owned by O900, all original IDs unchanged.
- Unambiguous example per outcome — confirmed by rendering actual card
  badges through `groupTaskCardsFor` + `TaskCard.badges()` (the real
  production path, not a re-implementation):
  - T903 → `[ ] CHECK` (no Check)
  - T907 → `[✓] CHECK` (current clear Check)
  - T905 → `[!] NEEDS WORK`
  - T909 → `[!] REVIEW` (stale clearance collapsed)
  - T910 → `[✓] WAIVED`
  - T908 → `[✓] OWNER ACCEPTED`
- T901 still shows `→ WAITS T900` (wait blocker), T904 still shows
  `⚠ REPLAN`, and T906 shows `[✓] CHECK` plus an independent `! OWNER`
  blocker — all confirmed in the same render, each shown alongside, not in
  place of, its outcome badge.
- No retired stale/exception/completion vocabulary presented as a card label
  in any reviewed title or description (targeted read of all 12 Task files
  and the Objective file); the two prior titles that echoed retired labels
  (T908, T909) were reworded above.
- The refreshed records pass strict V2 loading and render through the real
  card and non-TTY-backing paths: `data.LoadProject` loaded the live
  `.savepoint` project with no decode error, and `groupTaskCardsFor` /
  `TaskCard.badges()` (the same functions `renderCard` and `renderPlain`
  call) rendered every card above with no panic or missing badge.

**Commands run:**

- Temporary verification: a throwaway `_test.go` in `internal/board/v2`
  called `data.LoadProject("/home/user/code/savepoint/.savepoint")` and
  printed every O900 card's resolved badges through `groupTaskCardsFor` /
  `TaskCard.badges()`; output matched every Done-When outcome above, then the
  file was deleted (not part of this diff — verification only, per Boundaries
  restricting production/test changes to what O900 fixture validity needs,
  and none was needed here).
- `go build ./...` — clean.
- `go test ./internal/board/v2/... ./internal/data/...` — all pass.
- `git diff --check` — no whitespace errors.
- `make build && make test` — full run, exit 0, all packages `ok`.

**Limitations:**

- No colour-disabled or narrow-terminal screenshot was taken; verification
  read resolved `Badge.Text()` values (the same glyph+label pairs colour and
  no-colour rendering both draw from), not a rendered terminal frame. The
  User Check (opening O900 in the actual interactive board) is still
  pending — that is the owner's / a Task Check's job, not this evidence.
- T900, T902, T911 remain visually identical ("planned, nothing recorded"),
  unchanged from before this Task; T005's Done-When does not require they
  differ, since planned cards show no outcome badge to distinguish.

## Drift Notes

None at planning time.
