---
id: T-044
title: List finished Objectives under a Done heading
objective: O-024
status: planned
complexity_tier: small
complexity_reason: Changes one data ordering function and the duplicate-rank filter, adds one heading rule to two renderers, and filters done rows from the reorder helper.
depends_on: []
owner_validation: {required: true}
planned_by: {role: planner, session: o024-plan-20260925}
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

Pending execution.

## Drift Notes

Design §8 Layout and Keybindings gain the `DONE` section and the done-row key
rule at Objective reconciliation. Design §1's "Objective planning order" line
gains "done Objectives last, by ID".
