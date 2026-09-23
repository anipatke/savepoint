---
id: I-012
title: Objective sidebar Needs Integration badge is misleading and duplicative with Check badge
type: defect
status: resolved
source:
  kind: report
  actor:
    role: owner
    session: user-review
  at: '2026-09-20T21:23:03Z'
severity: medium
resolution:
  disposition: accepted
  actor:
    role: owner
    session: user-request
  at: '2026-09-21T19:30:00Z'
  reason: >-
    Owner reviewed the repair evidence directly and accepted it as resolved
    without routing through an independent savepoint-check session.
history:
  - at: '2026-09-20T21:23:03Z'
    actor:
      role: owner
      session: user-review
    kind: observed
    note: User noted that Needs Integration badge on Objective sidebar is misleading and duplicative with the Check badge
  - at: '2026-09-21T19:15:00Z'
    actor:
      role: executor
      session: user-request
    kind: repair_attempted
    note: >-
      Removed objectiveIntegrationBadge entirely rather than reconciling it with
      objectiveCheckBadge; the sidebar row now carries exactly one completion-state
      badge. Owner also directed collapsing the Check badge's needs_work/stale/unverified
      states into one "Check (needs work)" wording (I-013's status-label fix bundled in
      the same pass) and folding a recorded owner exception into the plain current/"[✓]
      Check" badge rather than a fourth "BY EXCEPTION" notch, since the audience does not
      need that granularity on the compact row — the detail overlay's EXCEPTION section
      still names the owner, Check, and reason in full.
  - at: '2026-09-21T19:30:00Z'
    actor:
      role: owner
      session: user-request
    kind: owner_decision
    note: Owner reviewed the repair evidence directly and accepted it as resolved, without routing through an independent savepoint-check session.
---

# I-012: Objective sidebar Needs Integration badge is misleading and duplicative with Check badge

## Summary

In `internal/board/v2/objectives.go` and `internal/board/v2/badges.go`, the board renders both `objectiveIntegrationBadge` ("NEEDS INTEGRATION" / "INTEGRATED") and `objectiveCheckBadge` ("[ ] Check" / "[✓] Check") on the Objective sidebar row.

This is misleading and duplicative:
1. **Misleading vocabulary:** "Needs Integration" implies that development work (merging, glue code, wiring) remains incomplete, when in reality 100% of owned tasks are already marked `done` and the only remaining action is an independent verification check (`savepoint-check`).
2. **Duplicative state:** Both badges read the exact same underlying `clearance.State`. When all tasks are done and the check is pending, the row renders `⚠ NEEDS INTEGRATION   [ ] Check`. When verified CLEAR, the row renders `✓ INTEGRATED   [✓] Check`, stating the same clearance fact twice side-by-side.

## Evidence

- `internal/board/v2/badges.go:objectiveIntegrationBadge` (lines 175-184):
  Returns `⚠ NEEDS INTEGRATION` when `tasksComplete` is true and `clearance != ClearanceCurrent`, and `✓ INTEGRATED` when `clearance == ClearanceCurrent`.
- `internal/board/v2/badges.go:objectiveCheckBadge` (lines 135-140):
  Returns `[ ] Check` when `clearance != ClearanceCurrent`, and `[✓] Check` when `clearance == ClearanceCurrent`.
- `internal/board/v2/objectives.go:ObjectiveRow.badges` (lines 180-190):
  Unconditionally appends both `objectiveIntegrationBadge` and `objectiveCheckBadge`, placing them next to each other on the same row.
- In `internal/board/v2/badges.go:blockerBadge` (lines 193-198), Savepoint's design explicitly forbids showing duplicate clearance badges because *"showing both would print the same fact twice in two wordings"*, but `ObjectiveRow` violates this principle.

## Proof Needed

1. The Objective sidebar badge representation is consolidated or renamed to eliminate redundancy with the Check badge and remove misleading "integration" terminology.
2. An Objective where all tasks are complete clearly indicates that the mandatory Full Objective Check is needed without implying missing integration code.
3. Unit and view tests in `internal/board/v2` reflect the unified/clarified badge presentation.

## Repair Evidence

1. `internal/board/v2/badges.go`:
   - Deleted `objectiveIntegrationBadge` ("INTEGRATED" / "NEEDS INTEGRATION" / "BY EXCEPTION") outright.
   - `objectiveCheckBadge` is now the row's sole completion-state badge: `func objectiveCheckBadge(clearance data.ClearanceState, byException bool) Badge`. `byException` folds to the same green `[✓] Check` a current Check gets; otherwise `ClearanceCurrent` → `[✓] Check`, `ClearanceMissing` → grey `[ ] Check`, and `needs_work`/`stale`/`unknown` collapse to one amber `[!] Check (needs work)` (see I-013's bundled scope note there — this three-notch simplification was an explicit owner call in the same conversation, not a separate defect).
2. `internal/board/v2/objectives.go`:
   - `ObjectiveRow.badges()` now appends exactly one completion badge (`objectiveCheckBadge(r.Clearance.State, r.ByException)`), no longer gated on `TasksComplete`, plus wait badges. No second badge duplicates the same clearance fact.
3. `internal/board/v2/detail.go` / `detail_view.go`:
   - `RecordDetail` carries a new `ByException` field, resolved the same way `ObjectiveRow.ByException` is (`ResolveObjectiveCompletion(...).AllowedByException`), so the Objective detail overlay's CLEARANCE badge agrees with the sidebar row for the same record. The overlay's EXCEPTION section (unchanged) still states the owner, Check, and reason in full regardless of what the compact badge shows.
4. Tests: `internal/board/v2/badges_test.go` (`TestObjectiveCheckBadgeIsThreeNotchOnly`, `TestObjectiveCheckBadgeFoldsExceptionIntoCurrent`), `internal/board/v2/objectives_test.go` (`TestSidebarShowsEachObjectivesOwnCheckBadge`, `TestSidebarSeparatesFinishedTasksFromAFinishedObjective` rewritten, `TestObjectiveRowBadgesFoldExceptionIntoCheck` new), `internal/board/v2/detail_test.go` updated for the new `clearanceLines` signature.

## Resolution

Closed by owner acceptance (`resolution.disposition: accepted`), not by an independent `savepoint-check` verification — the owner reviewed the repair evidence above directly and chose to accept it rather than route it through a separate Check session. This is a deliberate exception to this project's usual pattern (see I-014/I-015/I-017), recorded here rather than silently done.
