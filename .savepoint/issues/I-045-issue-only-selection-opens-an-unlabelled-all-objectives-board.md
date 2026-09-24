---
id: I-045
title: An Issue-only router selection opens the board on an unlabelled all-Objectives view
type: defect
status: open
source:
  kind: report
  actor: {role: owner, session: owner-chat-20260924}
  at: '2026-09-24T09:40:00Z'
tasks: [T-031]
checks: [C-916]
severity: low
history:
  - at: '2026-09-24T09:40:00Z'
    actor: {role: owner, session: owner-chat-20260924}
    kind: observed
    note: >-
      The owner saw O-019's T-010..T-012 and O-023's T-035 in the Planned
      column after the router moved to `issue: I-044`; the Tasks vanished once
      another Objective was picked in the sidebar.
  - at: '2026-09-24T09:55:00Z'
    actor: {role: executor, session: o014-i045-repair-20260924}
    kind: repair_attempted
    note: >-
      Board startup and reload read routerObjective, which resolves an
      Issue-only selection to the Issue's one owning Objective
      (internal/board/v2/model.go). The terminal board labels the unfiltered
      view ALL OBJECTIVES and the plain output prints "Selected: all
      Objectives". Tests: TestIssueOnlySelectionOpensOnTheIssuesObjective,
      TestIssueOnlySelectionWithUnknownIssueOpensUnfiltered,
      TestFilteredViewHasNoAllObjectivesLabel. make build && make test-fast
      pass; git diff --check clean. Design section 8 updated. Left open for
      the O-014 re-check.
---

# I-045: An Issue-only router selection opens the board on an unlabelled all-Objectives view

## Summary

The board's column filter starts from the router's `objective:`
(`internal/board/v2/model.go` `selectedObjectiveUnscoped`). T-031 made an
Issue-only selection legitimate (`objective: none`, `issue: I-###`), so the
board now opens with no Objective filter and shows every Task in the Goal. That
view existed before only behind `esc:clear`. Nothing on screen says it is
unfiltered: no sidebar row is marked and the columns look like one Objective's
Tasks, so it reads as a broken filter.

## Reproduction

Router `release: R-006`, `objective: none`, `task: none`, `issue: I-044`.
`savepoint board` lists `T-010`, `T-011`, `T-012` (O-019) and `T-035` (O-023)
under PLANNED, and 30 Tasks under DONE.

## Repair

Owner decision, 2026-09-24: both parts, repaired directly under this Issue.

1. When the router selects an Issue alone, the board opens on the one
   Objective that the Issue's linked Tasks and Objective-scoped Checks belong
   to (I-044 → O-014). An Issue linked to no Objective, or to several, keeps
   the unfiltered view.
2. When no Objective filter is in effect, the board says so: the terminal
   board's Goal line adds `ALL OBJECTIVES`, and the plain output prints
   `Selected: all Objectives`.
