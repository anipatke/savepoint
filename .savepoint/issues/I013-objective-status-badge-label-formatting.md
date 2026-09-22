---
id: I013
title: Objective status labels should be human-readable (Planned, In Progress, Done)
type: defect
status: resolved
source:
  kind: report
  actor:
    role: owner
    session: user-review
  at: '2026-09-20T21:24:55Z'
severity: low
resolution:
  disposition: accepted
  actor:
    role: owner
    session: user-request
  at: '2026-09-22T09:17:35Z'
  reason: >-
    Owner waived an independent Check and accepted I013 as resolved based on
    the repair evidence recorded by the implementing task; this is not a
    technical CLEAR verdict.
history:
  - at: '2026-09-20T21:24:55Z'
    actor:
      role: owner
      session: user-review
    kind: observed
    note: User requested Objective status badges/labels be formatted as Planned, In Progress, and Done instead of raw machine values
  - at: '2026-09-21T19:15:00Z'
    actor:
      role: executor
      session: user-request
    kind: repair_attempted
    note: >-
      Added objectiveStatusLabel mapping planned/in_progress/done to Planned/In
      Progress/Done for the sidebar row's status line; frontmatter and resolvers
      unchanged. Bundled in the same pass as I012's badge consolidation since both
      touched renderObjectiveRow.
  - at: '2026-09-22T09:17:35Z'
    actor:
      role: owner
      session: user-request
    kind: owner_decision
    note: >-
      Owner waived independent verification and accepted I013 as resolved based
      on the repair evidence from the implementing task.
---

# I013: Objective status labels should be human-readable (Planned, In Progress, Done)

## Summary

The Objective row in the sidebar (`internal/board/v2/objectives.go:renderObjectiveRow`) and related Objective presentation surfaces display raw internal lifecycle status values (`planned`, `in_progress`, `done`) rather than human-readable title-cased labels (`Planned`, `In Progress`, `Done`).

This is a presentation/label change only: internal canonical constants and YAML frontmatter values remain `planned`, `in_progress`, and `done`, while the UI renders human-friendly labels.

## Evidence

- `internal/board/v2/objectives.go:renderObjectiveRow` (line 286):
  Directly stringifies the record status: `styles.CardMeta.Render(indent(xansi.Truncate(string(row.Objective.Status), textW, "…")))`, outputting raw lowercase identifiers with underscores like `in_progress`.
- In contrast, Task card stages and completion states use human-facing labels in `internal/board/v2/badges.go` (e.g., `BUILD`, `TEST`, `CHECK`, `DONE`) and `internal/board/v2/next_panel.go` (`Planned`, `Build`, `Test`, `Check`, `Done`).

## Proof Needed

1. Add a label mapping for Objective status presentation:
   - `planned` -> `Planned`
   - `in_progress` -> `In Progress`
   - `done` -> `Done`
2. The Objective sidebar row (and any associated badge/display functions) renders these formatted labels rather than raw identifiers.
3. Frontmatter serialization and internal storage remain canonical `planned`, `in_progress`, and `done`.
4. Board tests in `internal/board/v2` assert the updated human-readable labels.

## Repair Evidence

1. `internal/board/v2/objectives.go`:
   - Added `objectiveStatusLabel(status data.ColumnType) string`, mapping `ColumnPlanned`/`ColumnInProgress`/`ColumnDone` to `Planned`/`In Progress`/`Done`; any other value reports itself rather than guessing.
   - `renderObjectiveRow` now renders `objectiveStatusLabel(row.Objective.Status)` instead of stringifying the raw record field. No frontmatter, resolver, or storage value changed.
2. Tests: `internal/board/v2/objectives_test.go` (`TestSidebarListsEveryObjectiveInOrder` updated to assert `Done`/`In Progress`/`Planned`).

The implementing task did not close this Issue. It was subsequently resolved by
owner acceptance with an independent Check waived; this is not a technical
`CLEAR` verdict.
