---
id: I-055
title: Design does not describe the DONE Objective section or done-row order keys
type: drift
status: resolved
source:
  kind: check
  check: C-928
  actor: {role: checker, session: o024-objective-check-20260925}
  at: '2026-09-25T08:22:37Z'
tasks: [T-044]
checks: [C-928, C-929]
severity: low
resolution:
  disposition: verified
  check: C-929
  actor: {role: checker, session: o024-objective-recheck-20260925}
  at: '2026-09-25T08:45:28Z'
  reason: Design §1 and §8 now document the DONE section, done-last ID order, open-only duplicate-rank detection, write-set exclusion, and done-row key behavior; the current sources agree and the Full Objective gate passed.
history:
  - at: '2026-09-25T08:22:37Z'
    actor: {role: checker, session: o024-objective-check-20260925}
    kind: observed
    check: C-928
    note: Design §1 and §8 still describe the O-020 order, in which done Objectives keep their ranked position; O-024 moved them under a DONE heading and made order keys no-ops on them.
  - at: '2026-09-25T08:42:09Z'
    actor: {role: planner, session: o024-design-reconciliation-20260925}
    kind: repair_attempted
    note: Updated Design §§1 and 8 to document finished Objective ordering, the DONE heading, duplicate-rank behavior, and reorder-key behavior; no code changed. Awaiting independent recheck.
  - at: '2026-09-25T08:45:28Z'
    actor: {role: checker, session: o024-objective-recheck-20260925}
    kind: rechecked
    check: C-929
    note: Design §1 and §8 correction matched the implementation and tests; the Full Objective recheck passed.
---

# I-055: Design does not describe the DONE Objective section or done-row order keys

## Summary

O-024 changed where finished Objectives appear and how the order keys treat
them, but `.savepoint/Design.md` still describes the O-020 behavior. O-024's
Boundaries list "the Design §8 wording at reconciliation" as in scope, and
T-044's Drift Notes say §1 and §8 would be updated at Objective
reconciliation. That did not happen before the Full Objective Check, which
covers reconciliation against Design.

## Evidence

- `.savepoint/Design.md:29` (§1 **Objective planning order**) describes the
  canonical order without saying done Objectives come last in ID order, and
  says duplicate ranks are reported without saying done Objectives are
  excluded. Code: `internal/data/project.go` `OrderedObjectiveIDsForGoal`
  (done rows sort after open rows, by ID) and `duplicateObjectiveRankFacts`
  (skips `done`).
- `.savepoint/Design.md:179` (§8 **Layout**) says the sidebar is grouped under
  `CRITICAL`/`HIGH`/`MEDIUM`/`LOW` and "Sidebar rows follow priority, rank,
  then stable ID order". It does not mention the `DONE` heading, which
  `internal/board/v2/objectives.go` `objectiveRowHeading` and
  `internal/board/v2/plain.go` now render for every done Objective after the
  open ones.
- `.savepoint/Design.md:210` (§8 **Keybindings**) says `1`–`4` and
  `K`/`J`/`shift+↑`/`shift+↓` act on "the focused Objective". It does not say
  they do nothing on a done row or that moves only swap among open rows.
  Code: `internal/board/v2/update.go` `sidebarObjectiveOrderChange` returns no
  change for a done row, and `sidebarObjectiveIDsInPriority` excludes done
  rows.
- T-044 Drift Notes: "Design §8 Layout and Keybindings gain the `DONE` section
  and the done-row key rule at Objective reconciliation. Design §1's
  'Objective planning order' line gains 'done Objectives last, by ID'."

## Proof Needed

Planner updates Design §1 (Objective planning order) and §8 (Layout,
Keybindings) to describe the implemented behavior. The update changes no
code. A recheck confirms the text matches `project.go`, `objectives.go`,
`plain.go`, and `update.go`.
