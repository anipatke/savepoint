---
id: T003
title: Give every Task card one clear review outcome
objective: O012
planned_by: {role: planner, session: task-card-simplification-20260921}
status: done
owner_validation:
    required: true
    accepted_check: ""
check_waiver:
    task: T003
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-21T11:00:33Z"
---

# T003: Give every Task card one clear review outcome

## Outcome

The V2 board derives one compact review outcome for each non-planned Task card,
removes redundant completion badges, and preserves actionable blockers across
interactive and plain rendering without changing any data-layer decision.

## User Check

After the focused and full gates pass, open a representative board as the
owner and confirm that current Checks, owner waivers, owner-accepted risk,
failed Checks, and review-needed evidence are visually distinct; Done cards do
not repeat their column status; and the layout remains readable at a narrow
terminal width. Record acceptance against the current Task Check if one is
requested. If the optional Task Check is waived, record the explicit Task
waiver and rely on O012's mandatory Full Objective Check for integration.

## Done When

- Card outcome precedence is `OWNER ACCEPTED` → `WAIVED` → resolved clearance,
  so one card never prints competing review outcomes.
- Missing, current, needs-work, stale/unknown, waiver, and exception-approved
  paths render exactly as `[ ] CHECK`, `[✓] CHECK`, `[!] NEEDS WORK`,
  `[!] REVIEW`, `[✓] WAIVED`, and `[✓] OWNER ACCEPTED` respectively.
- `WAIVED` and `OWNER ACCEPTED` use the existing green clear style without
  changing their underlying clearance, dependency, or gate semantics.
- `✓ DONE`, `⚠ DONE`, `BY WAIVER`, `BY EXCEPTION`, `Check (stale)`, and
  `Check (unverified)` are absent from Task cards and plain Task-card output.
- Replan, Task wait, Objective wait, and unresolved owner-action badges remain;
  checker-authority is represented by `[!] REVIEW`, and an accepted exception
  does not also show its overridden owner blocker.
- Planned cards still omit review outcomes; stage labels and card/column focus
  treatments remain unchanged.
- Full-colour, colour-disabled, narrow-width, interactive, and non-TTY tests
  prove the new mapping and fixed badge order.

## Context Files

`internal/board/v2/card.go`, `internal/board/v2/badges.go`, `internal/board/v2/plain.go`, `internal/board/v2/card_test.go`, `internal/board/v2/badges_test.go`, `internal/board/v2/columns_view_test.go`, `internal/board/v2/run_test.go`, `internal/board/v2/fixture_test.go`, `internal/data/gate_v2.go`, `internal/data/evidence_v2.go`, `internal/styles/styles.go`, `.savepoint/visual-identity.md`.

## Design References

Design sections 1, 4, 7, 8, and 13.

## Guardrails

DATA-02, ARCH-02, TEST-01, TEST-02, TEST-08, STYLE-03, STYLE-07, STYLE-08, STYLE-09, STYLE-10.

## Implementation Plan

1. Preserve `newTaskCard`'s existing resolver inputs and waiver/exception flags;
   introduce no card-local gate decision.
2. Replace the separate completion-plus-Check composition with one review
   outcome mapping using the precedence and exact labels above.
3. Remove or retire the Task-only completion mapping after confirming it has no
   remaining consumer; leave Objective badge behavior untouched.
4. Keep actionable blockers, fold checker-authority into generic review, and
   suppress an owner blocker only when the recorded exception has already
   produced `OWNER ACCEPTED`.
5. Update card, badge, full-board, plain-output, no-colour, and width tests to
   cover every outcome and reject the retired strings.
6. Run focused V2 board tests, `git diff --check`, `make build`, and `make test`;
   record exact results and any limitations in Technical Evidence.

## Boundaries

No changes to `internal/data` types or resolvers, Objective/sidebar badges,
detail/Next/resume/doctor wording, Task status/stage vocabulary, palette values,
or project records used as live UI fixture data.

## Technical Verification

Focused `go test ./internal/board/v2`, named card-state and non-TTY assertions,
colour-disabled distinctness, narrow-width rendering, `git diff --check`,
`make build`, and `make test`. A requested Task Check uses Quick evidence from
`agent-skills/references/check-method.md`; O012 closure requires the mandatory
Full Objective Check.

## Technical Evidence

**Files changed:**

- `internal/board/v2/badges.go` — added `taskReviewOutcomeBadge` (the single
  review-outcome mapping, precedence: exception → waiver → resolved
  clearance); removed `completionBadge` and `exceptionBadge` (no remaining
  consumer after the change, confirmed by grep before removal); folded
  `GateBlockCheckerAuthority` into the suppressed-blocker set in
  `blockerBadge` (already-current REVIEW outcome states it); removed the
  now-unused `glyphClear` and `glyphUnknown` constants. `taskCheckBadge` and
  `objectiveCheckBadge` are untouched — `taskCheckBadge` is still read by
  `detail_view.go`'s Task detail overlay (out of this Task's Context Files
  and Boundaries), and `objectiveCheckBadge` is the sidebar's own vocabulary.
- `internal/board/v2/card.go` — `TaskCard.badges()` now emits stage, one
  `taskReviewOutcomeBadge` for every non-planned card (Done included, so
  completion is never double-stated), then remaining blockers. No resolver
  input or `newTaskCard` change.
- `internal/board/v2/plain.go` — unchanged; it already renders through
  `card.badges()`, so the new vocabulary reaches non-TTY output for free.
- `internal/board/v2/badges_test.go`, `card_test.go`, `columns_view_test.go`,
  `run_test.go` — updated to the new labels and added
  `TestTaskReviewOutcomeBadgeRendersEveryPathAtExactText` (all 6 Done-When
  strings, plus the 2 precedence-overrides-clearance cases),
  `TestTaskReviewOutcomeBadgeStaysDistinctAcrossThePrecedenceLadder`,
  `TestRenderCardDoneCardsShowOneReviewOutcomeAndNoCompletionBadge` (replaces
  the old completion-badge test; asserts the retired strings are absent), and
  `TestRenderCardOpenExceptionOmitsTheOwnerBlockerItResolved`.

**Per-criterion evidence (Done When):**

- Precedence `OWNER ACCEPTED` → `WAIVED` → resolved clearance:
  `taskReviewOutcomeBadge`'s switch order, proven by the "owner accepted
  outranks waived" case in `taskReviewOutcomeCases`.
- The six exact renderings: `TestTaskReviewOutcomeBadgeRendersEveryPathAtExactText`
  checks `Text()` equals `"[ ] CHECK"`, `"[✓] CHECK"`, `"[!] NEEDS WORK"`,
  `"[!] REVIEW"` (both stale and unknown), `"[✓] WAIVED"`, `"[✓] OWNER ACCEPTED"`
  byte-for-byte.
- `WAIVED`/`OWNER ACCEPTED` keep the green clear style: both branches use
  `styles.BadgeClear`, same as `ClearanceCurrent`; no clearance, dependency,
  or gate semantics touched (`internal/data` untouched).
- Retired strings absent from cards and plain output:
  `TestRenderCardDoneCardsShowOneReviewOutcomeAndNoCompletionBadge` and
  `TestBoardShowsTasksInTheColumnTheirStatusNames` assert `"✓ DONE"`,
  `"⚠ DONE"`, `"BY EXCEPTION"`, `"BY WAIVER"`, `"Check (stale)"`,
  `"Check (unverified)"` do not appear.
- Replan/Task-wait/Objective-wait/owner-action badges remain; checker
  authority folds into `[!] REVIEW`: `blockerBadge` keeps Replan, Dependency,
  ObjectiveDependency, and OwnerAcceptance stated (`TestBlockerBadgeCoversEveryKind`),
  now excludes CheckerAuthority. An owner-accepted override showing no
  overridden blocker is proven directly by
  `TestRenderCardOpenExceptionOmitsTheOwnerBlockerItResolved`, and is
  structurally guaranteed upstream: `data.ResolveTaskCompletion` returns an
  empty `Blockers` slice whenever `AllowedByException` is true
  (`internal/data/gate_v2.go`), which this Task does not modify.
- Planned cards omit review outcomes, stage labels/focus unchanged:
  `TestRenderCardPlannedOmitsCheckBadge`, `TestRenderCardStageAbsentOffAnInProgressTask`,
  `TestRenderCardPlannedFocusUsesMutedStyling`, `TestRenderCardFocusChangesColorNotGeometry`
  pass unmodified — the planned/stage/focus code paths were not touched.
- Full-colour / colour-disabled / narrow-width / interactive / non-TTY:
  `TestBadgesCarryGlyphAndLabelUnderColor` and `TestBadgesStayDistinctWithColorDisabled`
  cover `allBadges()` including every `taskReviewOutcomeCases` entry;
  `TestRenderCardNeverExceedsItsWidth` covers narrow widths (20/28/40/72) with
  a `ClearanceUnknown` (now `[!] REVIEW`) card; `TestBoardShowsTasksInTheColumnTheirStatusNames`
  / `TestBoardLinesNeverExceedTheTerminalWidth` cover the interactive board;
  `TestRunWithoutTTYIncludesTitlesBadgesAndIssueSummary` covers plain output.

**Commands run:**

- `go build ./...` — clean, after the badges.go/card.go edit and again after
  the final test-file edit.
- `go test ./internal/board/v2/...` (plain and `-v`) — all pass, including
  after each incremental test-file edit.
- `git diff --check` — no whitespace errors.
- `make build && make test` — full run, exit 0. All packages `ok`, including
  `internal/migrate` (121s) and the changed `internal/board/v2` (1.7s).

**Limitations:**

- `internal/board/v2/detail_view.go`'s comment on `clearanceLines` says the
  Task detail badge is "the badge a card ... would show, so the overlay
  agrees with them at a glance." After this change that is no longer
  literally true — the card now uses `taskReviewOutcomeBadge`'s wording
  (`"[!] REVIEW"`) where the overlay still uses `taskCheckBadge`'s
  (`"Check (stale)"`). Both state the same underlying clearance fact in
  different words; this is presentation drift in a source comment, not a
  functional or policy defect, and `detail_view.go` is outside this Task's
  Context Files and explicitly out of scope (Boundaries: no changes to
  detail wording). Flagging here rather than editing it or opening an Issue,
  since T004 owns active-documentation reconciliation for this Objective.
- No visual/owner review has been performed yet (User Check is pending); this
  evidence covers only the technical verification this skill can perform.
- The Go test fixture `writeBadgeProject` (fixture_test.go) was not extended
  with a waiver-state Task; the waiver path is exercised directly at the
  function level (`taskReviewOutcomeCases`, `allBadges()`) rather than through
  a full project fixture, consistent with the prior scope of the board-level
  tests (which also never exercised a `BY WAIVER` case at that layer).

## Drift Notes

Any required data-layer or gate-policy change is REPLAN REQUIRED. Presentation
drift discovered in active documentation is handed to T004 rather than folded
into this implementation Task.

## Addendum: build/test-stage owner refinement (post-waiver)

After this Task's owner waiver, the owner reviewed the rendered board and
asked to narrow the review-outcome badge further: an in-progress Task at
`build` or `test` stage with `ClearanceMissing` (nothing recorded yet) no
longer shows the pending `[ ] CHECK` badge — it is redundant with the stage
badge, since a build/test-stage Task cannot yet have been reviewed for the
work it is doing now. The badge still shows in that same stage when there is
a real recorded outcome (a NEEDS WORK Check that sent the Task back to build
for repair, or clearance gone stale/unknown), and always shows at `audit`
stage and on Done cards — those are exactly the cases carrying information
the stage badge does not.

This revises this Task's own Done-When line ("Missing... paths render
exactly as `[ ] CHECK`... respectively") to add the build/test-with-nothing-
recorded exception; every other path (current, needs-work, stale/unknown,
waived, exception) is unchanged. No `internal/data` types, resolvers, or
gate semantics changed — this is presentation composition in `card.go`
(`TaskCard.showsReviewOutcome`), calling the unchanged `taskReviewOutcomeBadge`
conditionally rather than changing what it renders.

**Files changed:** `internal/board/v2/card.go` (added `showsReviewOutcome`,
called from `badges()` in place of the unconditional
`c.Task.Status != data.ColumnPlanned` check); `internal/board/v2/card_test.go`
(replaced the in-progress half of `TestRenderCardPlannedOmitsCheckBadge` —
which had asserted the now-retired always-shown behavior — with
`TestRenderCardBuildOrTestOmitsCheckBadgeWhenNothingIsRecorded`,
`TestRenderCardBuildOrTestStillShowsARealOutcome`, and
`TestRenderCardAuditAlwaysShowsCheckBadgeEvenWhenMissing`);
`internal/board/v2/columns_view_test.go` (`TestBoardShowsTasksInTheColumnTheirStatusNames`
no longer asserts `"[ ] CHECK"` present — this fixture's build/test-stage
Tasks now correctly suppress it — and asserts it does not reappear there).

**Commands run:** `go build ./...` clean; `go test ./internal/board/v2/...`
all pass (including `TestBadgeVocabularyLivesInOneFile`, which required
routing the missing-clearance check through the existing `clearanceIsMissing`
helper in `badges.go` rather than referencing `data.ClearanceMissing` from
`card.go` directly); `git diff --check` clean; `make build && make test`
full run, exit 0, all packages `ok`. Re-verified the live O900 fixture
(T005) renders correctly under this change: T903 (build) and T904 (test),
both previously `[ ] CHECK`, now show no review-outcome badge; T905's
`[!] NEEDS WORK` at audit is unaffected.

**Limitations:** No fresh User Check or owner waiver has been recorded for
this addendum specifically; it rides on T003's existing waiver and the same
mandatory O012 Full Objective Check that was always going to review this
file.

## Addendum 2: Check outcome is Done-column vocabulary; OWNER blocker reworded

The owner reviewed the first addendum's result live (T906: `[◆ CHECK]
[✓] CHECK] [! OWNER]`) and judged it still overcooked — a green "checked and
clear" badge sitting next to an orange blocker reads as a contradiction, even
though the two facts (an independent Check cleared it; the owner still has
to sign off) are both true. Rather than reword around that one collision, the
rule was generalized and the build/test-only suppression from Addendum 1 was
replaced outright:

- `TaskCard.showsReviewOutcome` (`card.go`) now applies to every open stage,
  not just build/test: an open Task's card states its review outcome only
  when the outcome is itself actionable (`NEEDS WORK`, or the stale/unknown
  `REVIEW` fold) or the Task carries an owner's own waiver/exception.
  `ClearanceCurrent` and `ClearanceMissing` — "checked and clear" and "not
  checked yet" — are completion-outcome vocabulary now reserved for the Done
  column, where a Task's review history is actually reported. This
  supersedes Addendum 1's narrower, stage-scoped rule (audit stage no longer
  gets special-cased to always show it).
- The new boolean lives in `badges.go` as `reviewOutcomeIsActionable(state
  data.ClearanceState)` — `TestBadgeVocabularyLivesInOneFile` requires any
  `data.Clearance*` comparison to live there, not in `card.go`.
- `blockerBadge`'s `GateBlockOwnerAcceptance` label changed from `"OWNER"` to
  `"AWAITS OWNER"` (`badges.go`) — with the Check outcome now silent whenever
  a blocker is the whole story, `AWAITS OWNER` has to carry its meaning on
  its own rather than sit next to a `CHECK` badge for contrast.

T906's card is now exactly `[◆ CHECK] [! AWAITS OWNER]` — one clean signal:
at audit, waiting on the owner. Re-verified against the live O900 fixture
the same way as Addendum 1.

**Files changed (beyond Addendum 1's):** `internal/board/v2/card.go`
(`showsReviewOutcome` rewritten); `internal/board/v2/badges.go`
(`reviewOutcomeIsActionable` added; `GateBlockOwnerAcceptance` label
renamed); `internal/board/v2/card_test.go` (replaced the Addendum-1
build/test-only tests with `TestRenderCardInProgressOmitsCheckBadgeWhenNotActionable`
— build/test/audit × missing/current — `TestRenderCardInProgressStillShowsAnActionableOutcome`,
and the concrete regression case
`TestRenderCardCurrentCheckAwaitingOwnerShowsOnlyTheOwnerBlocker`);
`internal/board/v2/columns_view_test.go` (`"! OWNER"` → `"! AWAITS OWNER"`
in the board-level must-appear list); `internal/board/v2/run_test.go`
(`TestRunWithoutTTYIncludesTitlesBadgesAndIssueSummary`'s plain-output
assertion updated from the old `"[◆ CHECK  [✓] CHECK"` to the current
`"[◆ CHECK]"` for T002 and `"[◆ CHECK  ! AWAITS OWNER]"` for T006, proving
the same behavior end to end through the non-TTY path).

**Commands run:** same as Addendum 1 — `go build ./...` clean; `go test
./internal/board/v2/...` all pass (including `TestBadgeVocabularyLivesInOneFile`);
`git diff --check` clean; `make build && make test` full run, exit 0, all
packages `ok`; live O900 fixture re-verified via the same throwaway
`groupTaskCardsFor`/`badges()` check (not committed).

**Limitations:** Same as Addendum 1 — this rides on T003's existing waiver
and the mandatory O012 Full Objective Check. `taskReviewOutcomeBadge` itself
and its six exact strings are unchanged; only when `badges()` calls it, and
the unrelated `AWAITS OWNER` blocker label, changed.

