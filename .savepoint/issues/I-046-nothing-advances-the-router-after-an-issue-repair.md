---
id: I-046
title: Nothing moves the router off an Issue once its repair is done
type: drift
status: open
source:
  kind: report
  actor: {role: owner, session: owner-chat-20260924}
  at: '2026-09-24T11:10:00Z'
tasks: [T-031, T-034]
checks: [C-916]
severity: low
history:
  - at: '2026-09-24T11:10:00Z'
    actor: {role: owner, session: owner-chat-20260924}
    kind: observed
    note: >-
      After I-044 and I-045 were repaired, the board header still read
      "Fix I-044 — …" until the router was edited by hand to select O-014.
  - at: '2026-09-24T09:37:40Z'
    actor: {role: executor, session: repair-I-046-20260924}
    kind: repair_attempted
    note: >-
      Added Issue-only post-repair router advancement to live and scaffold
      AGENTS.md, savepoint-task, and issue-capture guidance. The rule selects
      the one Objective identified by linked Tasks and Objective Checks, or
      clears only the Issue when there is no unique Objective. Preserves the
      Goal selection. make build and make test-fast passed; git diff --check
      passed. The router already selected O-014 before this repair.
---

# I-046: Nothing moves the router off an Issue once its repair is done

## Summary

AGENTS.md's Router Selection section names who advances the router for every
case O-014 considered: the board after the owner closes a Task,
`savepoint-design` for the next Objective, and `savepoint-task` for the Task
it starts. An Issue-only selection (T-031) has no such step. After a direct
repair, `savepoint-task` records `repair_attempted` and leaves the Issue open
for a checker, so the Issue is not resolved and the stale-selection warning
does not fire. Next keeps saying `Fix I-### — …` after the fix has landed, and
the actual next step (the Objective's re-check) is not shown.

## Violated requirement

O-014 Success Conditions: guidance must "name who advances the router after an
owner closure", and Next must state the real next step from the router
selection alone.

## Reproduction

Router `objective: none`, `task: none`, `issue: I-044`. Repair I-044 and
record `repair_attempted`. `savepoint resume` still prints
`Fix I-044 — The Next line's Objective word never reads In Progress`.

## Repair

Proposed, guidance only: when `savepoint-task` finishes a direct Issue repair
and records `repair_attempted`, it points the router at the Issue's Objective,
using the same rule as the board's Issue-only filter (the one Objective the
Issue's linked Tasks and Objective-scoped Checks belong to), with `task` and
`issue` cleared. That usually makes Next read `… O-### · Check`. If the Issue
has no single Objective, it clears `issue` only, and Next reads
`Nothing selected`. `release:` is never changed. The rule goes into
`savepoint-task`, `issue-capture.md` (Out-Of-Scope Repair) and AGENTS.md's
Router Selection section, all in live and scaffold copies.

## Repair evidence

- Read I-046, the task skill, issue-capture guidance, O-014's boundaries,
  router state, the board's existing Issue-to-Objective rule, and applicable
  `STYLE-*` guardrails to keep the guidance aligned with current behavior.
- Changed the live and scaffold copies of AGENTS.md, savepoint-task, and
  issue-capture guidance. The new rule records `repair_attempted` before
  changing an Issue-only router selection, preserves `release:`, and leaves
  the Issue open for independent verification.
- `make build && make test-fast`: passed on 2026-09-24. `git diff --check`:
  passed. No runtime code, tests, fixtures, dependencies, or gate definitions
  changed in this repair.
- The router selected O-014 before this work, so this session could not
  observe an Issue-only transition in the live project. The checker should
  verify that transition and the ambiguous-link case in a fresh session.
