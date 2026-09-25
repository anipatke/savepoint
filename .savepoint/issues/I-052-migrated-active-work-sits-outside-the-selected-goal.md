---
id: I-052
title: Migration selects an empty continuation Goal and hides active work
type: defect
status: resolved
source:
  kind: check
  check: C-923
  actor: {role: checker, session: o022-objective-check-20260925}
  at: '2026-09-25T03:27:22Z'
tasks: [T-039, T-037]
checks: [C-923, C-924]
resolution:
  disposition: verified
  check: C-924
  actor: {role: checker, session: o022-independent-recheck-20260925}
  at: '2026-09-25T04:53:45Z'
severity: medium
history:
  - at: '2026-09-25T03:27:22Z'
    actor: {role: checker, session: o022-objective-check-20260925}
    kind: observed
    check: C-923
    note: After a fallback migration the router selects an empty R-002 while in-progress O-001 stays in R-001; the board is empty and, in the archived case, the router selection is not honored.
  - at: '2026-09-25T04:29:53Z'
    actor: {role: executor, session: o022-repair-followup-20260925}
    kind: repair_attempted
    note: >-
      Following the owner's choice, fallback now selects an existing live Goal
      containing the active work when possible. If the selected work belongs
      only to a historical Goal, migration creates a live continuation and
      assigns those active Objectives to it. A pending release lifecycle
      decision is left for the preview instead of selecting a Goal prematurely.
      End-to-end tests verify strict loading, resolved Next, and the selected
      Goal's Objective and Task set for missing, unresolvable, and archived
      router cases. The focused migration tests and golden reproducibility
      checks passed; git diff --check passed. Issue remains open for independent
      verification.
  - at: '2026-09-25T04:48:34Z'
    actor: {role: executor, session: o022-repair-followup-20260925}
    kind: repair_attempted
    note: >-
      Updated the stale ConvertRouter fallback expectation to match the owner's
      choice: missing and unresolvable selections reuse R-001, while the
      archived case creates R-002. The focused test passed. Reconciled O-022
      Success Condition 5, T-039 criteria and evidence, Design §10, and the
      live/scaffold Goal guidance with that policy. The independent Objective
      Check is rerunning the full gate and review on this contract update.
  - at: '2026-09-25T04:53:45Z'
    actor: {role: checker, session: o022-independent-recheck-20260925}
    kind: rechecked
    check: C-924
    note: >-
      C-924 verified the owner-selected fallback: reuse the live Goal that
      contains the router's active Objective, or the sole live Goal when the
      choice is clear; when selected active work belongs only to a historical
      Goal, create a continuation and move that Objective into it. End-to-end
      checks resolve Next and expose the active work in the selected Goal.
---

# I-052: Migration selects an empty continuation Goal and hides active work

## Summary

O-022's Outcome says migrated projects start valid. With T-037 the board shows
only the selected Goal's work. T-039's continuation Goal gets no Objectives:
every converted Objective keeps its V1 release's Goal. After a fallback
migration, the router selects an empty Goal. The project's in-progress work
is not on the board. In the archived case, the migrated router selection is
also rejected. Before O-022 the same migration resumed on the active Task.

## Evidence

Materialize the `internal/migrate/testdata/golden` outputs, set
`schema_version: 2`, and run the binary:

- `v1-router-missing.yml`: V1 release `v1` is live (R-001, `in_progress`) and owns
  in-progress O-001/T-001. The router selects the generated R-002. `savepoint
  board` shows `Selected: all Objectives in Goal R-002` and 0 cards in every
  column. `savepoint resume` reads `Nothing selected`. R-001 was a live Goal that
  could have been selected.
- `v1-router-archived.yml`: the router selects R-002, O-001, and T-001, but O-001
  names the historical R-001. `savepoint resume` reads `Nothing selected` with
  `The router selects Objective O-001 inside Goal R-002, but its record names
  Goal R-001, so this selection is not honored.` The board shows no work. R-001
  is archived, so `g` cannot select it, and O-001 is never shown on the board.
- `internal/migrate/plan.go` `planRouterGoalSelection` creates the continuation
  Goal whenever the router's release is absent or not live. It does not check
  whether another converted Goal is live, or where the active Objectives are.
  `planEpic` always assigns the Objective to its V1 release's Goal.
- `TestEndToEnd_routerFallbackCasesMigrateToLiveGoal` checks strict load and
  the router's Goal ID. It never resolves Next or the board's Goal view on the
  result.

## Proof Needed

- After migration, the router's selected Goal contains the converted active
  work. Next resolves the migrated selection without a mismatch diagnostic.
  The planner decides how: select an existing single live Goal, or move active
  Objectives from a historical release into the continuation Goal.
- End-to-end tests for the missing, unresolvable, and archived cases resolve
  Next and the Goal-scoped Task set on the migrated project.
