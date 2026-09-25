---
id: T-042
title: Show Objectives in priority groups
objective: O-020
status: planned
complexity_tier: medium
complexity_reason: Reworks sidebar row layout and windowing around group headings, removes the status line, and adds a matching Objective list to plain output.
depends_on: [{task: T-041, requires: clear}]
owner_validation: {required: true}
planned_by: {role: planner, session: o020-rank-objectives-20260925}
---

# Show Objectives in priority groups

## Outcome

The board's Objective sidebar lists the selected Goal's Objectives under
`CRITICAL`, `HIGH`, `MEDIUM`, and `LOW` headings in the canonical order, without
the separate Planned / In Progress / Done line. Piped board output lists the
same Objectives in the same order under the same headings.

## User Check

Open the board on this repository. Confirm the sidebar shows only the
headings that have Objectives, rows appear in rank order under each, and no
row carries a Planned / In Progress / Done line. Confirm the Check badge,
wait badges, the `●` selection and `▸` cursor glyphs are still there, also
with `NO_COLOR=1` and in a narrow terminal. Run `savepoint board | cat` and
confirm the Objective list matches the sidebar.

## Done When

- `objectiveRowsForRelease` uses the ordered function from `internal/data`;
  the sidebar sorts nothing itself.
- Headings are drawn only for non-empty groups. They are not rows: the cursor
  skips them, and up/down still move between Objectives. Windowing keeps the
  cursor row visible and counts heading lines in its budget.
- The status line and `objectiveStatusLabel` are removed from the sidebar.
  Objective detail and Next still report the recorded status and evidence.
- Plain output adds a grouped Objective list after the `Selected:` line, one
  line per Objective with its ID, title, and badge text, in the same order and
  headings as the sidebar. Two runs produce identical bytes.
- Completed Objectives stay in their ranked place; blocked ones stay in place
  with their wait badge.
- Tests cover ordering across all four groups, empty groups hidden, cursor
  skipping headings, scroll windowing with headings, narrow widths, no-colour
  legibility, completed and blocked rows in place, no status line, and TTY and
  non-TTY agreement.

## Context Files

`internal/board/v2/objectives.go`, `internal/board/v2/objectives_test.go`,
`internal/board/v2/plain.go`, `internal/board/v2/releases_test.go`,
`internal/board/v2/view.go`, `internal/board/v2/view_test.go`,
`internal/board/v2/width_test.go`, `internal/data/project.go`,
`.savepoint/visual-identity.md`,
`.savepoint/objectives/O-020-rank-release-objectives/Objective.md`.

## Design References

Design section 8 (TUI layout and render fallbacks). O-020 Confirmed Design
Decisions.

## Guardrails

ARCH-02, DATA-02, TEST-01..04, TEST-08, STYLE-07, STYLE-09.

## Implementation Plan

1. Switch row construction to the ordered projection.
2. Render group headings in the sidebar and adjust windowing so heading lines
   are budgeted and the cursor index still refers to Objective rows.
3. Drop the status line from `renderObjectiveRow`.
4. Add the grouped Objective list to `renderPlain` using the same heading
   labels and badge text.
5. Update and add the tests listed above.

## Boundaries

No keys, writes, or data parsing. No change to Next, the Task columns, or
Objective detail.

## Technical Verification

Focused tests while iterating; `make build && make test-fast` at handoff.

## Technical Evidence

Pending execution.

## Drift Notes

Design section 8 gains the grouped sidebar and plain Objective list at
Objective reconciliation.
