---
id: T-044
title: List finished Objectives under a Done heading
objective: O-024
status: done
complexity_tier: small
complexity_reason: Changes one data ordering function and the duplicate-rank filter, adds one heading rule to two renderers, and filters done rows from the reorder helper.
depends_on: []
owner_validation:
    required: true
    accepted_check: C-929
    accepted_by: {role: owner, session: owner-chat-20260925}
planned_by: {role: planner, session: o024-plan-20260925}
check_waiver:
    task: T-044
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-25T08:19:52Z"
---

# List finished Objectives under a Done heading

## Outcome

Open Objectives stay at the top of the sidebar under their priority headings.
Finished Objectives move below them under `DONE` in ID order, in both the TUI
and the piped board output, and the reorder keys ignore them.

## User Check

Open the board on this repository. Confirm O-015, O-017, and the other open
Objectives appear first under their priority headings. Confirm every finished
Objective is listed under `DONE` at the bottom in ID order, with no Planned /
In Progress / Done row line. Press `2`, `K`, and `J` on a done row: nothing
changes and no file is written. Press `K`/`J` on an open row: it moves only
among open rows. Run `savepoint board | cat` and confirm the list matches.

## Done When

- `OrderedObjectiveIDsForGoal` returns open Objectives in priority, rank, then
  ID order, followed by done Objectives in ascending ID order.
- `duplicateObjectiveRankFacts` skips done Objectives.
- The sidebar and `renderPlain` draw a `DONE` heading before the first done
  row and priority headings for open rows only. Headings appear only for
  non-empty groups. Sidebar windowing still counts every heading line.
- `sidebarObjectiveOrderChange` returns no change for a done row, and
  `sidebarObjectiveIDsInPriority` excludes done rows. No reorder writes a done
  Objective's file.
- Tests cover the Success Conditions test list in O-024.

## Context Files

`internal/data/project.go`, `internal/data/project_test.go`,
`internal/board/v2/objectives.go`, `internal/board/v2/objectives_test.go`,
`internal/board/v2/plain.go`, `internal/board/v2/update.go`,
`internal/board/v2/releases_test.go`,
`internal/doctor/v2_runtime_test.go`,
`.savepoint/objectives/O-024-move-done-objectives-below-open-work/Objective.md`.

## Design References

Design sections 1 (Objective planning order) and 8 (TUI layout, keybindings).
O-024 Confirmed Design Decisions.

## Guardrails

ARCH-02, DATA-02, FS-01, TEST-01..04, TEST-08, STYLE-07.

## Implementation Plan

1. In `OrderedObjectiveIDsForGoal`, sort done Objectives after open ones,
   and sort done Objectives by ID only.
2. Skip done Objectives in `duplicateObjectiveRankFacts`.
3. Add one board helper that returns a row's heading (`DONE` for a done
   Objective, otherwise its priority heading). Use it in `renderSidebar`,
   `visibleObjectiveWindow`, and `renderPlain` wherever they now compare
   priorities.
4. In `update.go`, return no change for a done focused row, and exclude done
   rows from `sidebarObjectiveIDsInPriority`.
5. Update the existing ordering tests that assumed done rows stay in rank
   position, and add the listed cases.

## Boundaries

No new keys, fields, or statuses. No collapsing of `DONE`. No change to Next,
the router, gates, or Task columns.

## Technical Verification

Focused tests while iterating; `make build && make test-fast` at handoff.

## Technical Evidence

- Per-criterion outcomes:
  - `OrderedObjectiveIDsForGoal` keeps open Objectives in priority/rank/ID
    order, then done Objectives in ID order: covered by
    `TestOrderedObjectiveIDsForGoalUsesPriorityRankAndID`.
  - Duplicate-rank facts and doctor warnings ignore done Objectives: covered by
    `TestLoadV2IndexIgnoresDuplicateRanksForDoneObjectives` and
    `TestRunV2ChecksIgnoresDuplicateRanksOnDoneObjectives`.
  - TUI and plain output share priority headings for open rows and one `DONE`
    heading for finished rows; empty and all-done groups and heading-aware
    windowing are covered by `TestSidebarAndPlainOutputUsePriorityAndRankOrder`,
    `TestDoneHeadingAppearsOnlyForNonEmptyDoneGroup`,
    `TestOnlyDoneObjectivesUseOneIDOrderedHeading`, and
    `TestSidebarCursorSkipsPriorityHeadingsAndWindowKeepsCurrentGroup`.
  - Done-row order keys are no-ops, reorder sets exclude done rows, and board
    writes leave done files untouched: covered by
    `TestSidebarReorderGroupsExcludeDoneObjectives`,
    `TestSidebarDoneObjectiveOrderKeysDoNotWriteProjectFiles`,
    `TestSidebarOpenReorderDoesNotWriteDoneObjectiveFile`, and
    `TestSidebarOrderHealsAfterInjectedMidWriteFailure`.
  - Reopening a done Objective restores its recorded rank position; the
    collision is detected and healed by a reorder: covered by
    `TestReopenedDoneObjectiveReturnsToRankedGroupAndHealsCollision`.
- Focused verification passed: `make test-focused` with the Objective ordering,
  heading, reorder, navigation, doctor, and reopen cases listed above; the same
  pattern passed with direct `go test ./... -run`, and the affected packages
  passed with direct `go test`.
- The first `make test-focused` attempt returned a package-setup failure without
  a diagnostic; subsequent direct package and full-repository runs passed,
  followed by a successful retry of `make test-focused`.
- `make build` passed after implementation and again after the final test fixes.
- The first `make test-fast` run exposed stale expectations in
  `TestSidebarOrderHealsAfterInjectedMidWriteFailure` and
  `TestSidebarSeparatesFinishedTasksFromAFinishedObjective`. The partial-write
  test now moves an open row and verifies the done row remains untouched; the
  sidebar row test helper now stops at actual `O-###` row starts. Both targeted
  cases passed after repair.
- Fresh `make test-fast` passed. `git diff --check` passed.
- Context files read: this Task, the O-024 Objective, `internal/data/project.go`,
  `internal/data/project_test.go`, `internal/board/v2/objectives.go`,
  `internal/board/v2/objectives_test.go`, `internal/board/v2/plain.go`,
  `internal/board/v2/update.go`, `internal/board/v2/releases_test.go`, and
  `internal/doctor/v2_runtime_test.go`.
- Extra reads and reasons: `agent-skills/savepoint-task/SKILL.md` for the active
  workflow; `.savepoint/router.md` to confirm the selected records;
  `.savepoint/Guardrails.md` for ARCH-02, DATA-02, FS-01, TEST-01..04,
  TEST-08, and STYLE-07;
  `.savepoint/objectives/O-020-rank-release-objectives/Objective.md` to confirm
  O-024's dependency was done; `internal/board/v2/fixture_test.go` and
  `internal/board/v2/columns_view_test.go` to inspect the existing sidebar
  fixture and key helpers; `internal/board/v2/actions_test.go` to inspect and
  update the partial-write regression after `make test-fast` exposed its stale
  done-row assumption.
- Files changed: O-024 `Objective.md`; this T-044 record; `internal/data/project.go`
  and `project_test.go`; `internal/board/v2/objectives.go`, `objectives_test.go`,
  `plain.go`, `update.go`, and `actions_test.go`; and
  `internal/doctor/v2_runtime_test.go`.
- Limitation: the required owner User Check (interactive keys, visual sidebar,
  and `savepoint board | cat`) remains for the owner. No Task Check was
  requested, and no owner waiver has been recorded.

## Drift Notes

Design §8 Layout and Keybindings gain the `DONE` section and the done-row key
rule at Objective reconciliation. Design §1's "Objective planning order" line
gains "done Objectives last, by ID".
