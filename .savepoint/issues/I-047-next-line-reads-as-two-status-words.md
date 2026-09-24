---
id: I-047
title: The Next line strings two status words together instead of saying what to do
type: defect
status: resolved
source:
  kind: report
  actor: {role: owner, session: owner-chat-20260924}
  at: '2026-09-24T11:40:00Z'
tasks: [T-028]
checks: [C-916, C-917]
severity: low
resolution:
  disposition: verified
  check: C-917
  actor: {role: checker, session: o014-objective-recheck-20260924}
  at: '2026-09-24T09:49:09Z'
history:
  - at: '2026-09-24T11:40:00Z'
    actor: {role: owner, session: owner-chat-20260924}
    kind: observed
    note: >-
      "In Progress O-014 · Check — Give the Next area an Objective word and
      let the router target Issues" does not read naturally.
  - at: '2026-09-24T11:45:00Z'
    actor: {role: owner, session: owner-chat-20260924}
    kind: owner_decision
    note: >-
      Replace O-014's confirmed line format with action-first lines (verb,
      record, title), repaired directly under this Issue.
  - at: '2026-09-24T12:05:00Z'
    actor: {role: executor, session: o014-i047-repair-20260924}
    kind: repair_attempted
    note: >-
      resume.NextLine now prints `<verb> <ID> — <title>`, with `(O-###)` after
      a Task, from the new resume.NextVerb; ObjectiveWord and TaskStageWord
      are removed. Goal completion rungs print Check/Accept/Close R-### instead
      of "Nothing selected". The board accents only the verb (verbStyle).
      AGENTS.md (live and scaffold) maps verbs to skills; Design section 8 and
      O-014's Confirmed Design Decisions record the format. Tests: NextLine
      table extended with Blocked, Replan, Accept on a Task and an Objective,
      the not-only-owner-acceptance fallback, a done Objective, a Task with no
      Objective record, and the three Goal verbs;
      TestRenderNextAccentsOnlyTheVerb; board, parity and main expectations
      moved to the new lines, including internal/board's V2 dispatch test,
      which the first fast-gate run caught. make build && make test-fast then
      passed; git diff --check clean. Left open for the O-014 re-check.
  - at: '2026-09-24T09:49:09Z'
    actor: {role: checker, session: o014-objective-recheck-20260924}
    kind: rechecked
    check: C-917
    note: "Next lines lead with action verbs, and current Checks use distinct Accept and Close owner steps."
---

# I-047: The Next line strings two status words together instead of saying what to do

## Summary

O-014's confirmed format `<Objective word> O-### · <Task word> T-### — <title>`
puts two status words in a row and then a title that reads as the object of the
second one. `· Check` also means both "run the Check" and "the Check passed,
record done" (C-916 observation 2). The Issue line `Fix I-044 — …` already
reads naturally because it leads with what to do.

## Repair

Owner decision, 2026-09-24: every Next line is `<verb> <ID> — <title>`, with a
Task line ending in its Objective, `(O-###)`.

| Selection | Verb |
| --- | --- |
| Task planned, may start / blocked | `Start` / `Blocked` |
| Task at build / test | `Build` / `Test` |
| Task at audit: needs Check / awaits owner acceptance / may close | `Check` / `Accept` / `Close` |
| Task carrying a replan | `Replan` |
| Objective, all Tasks done: needs Check / Check current, awaits owner / may close | `Check` / `Accept` / `Close` |
| Objective with no Tasks / Tasks left, none selected | `Plan` / `Pick a Task in` |
| Objective or Issue already finished | `Done` / `Resolved` (plus the stale warning) |
| Issue open or in progress | `Fix` |
| Goal (no Objective): needs Check / awaits owner / may close | `Check` / `Accept` / `Close` |

This also fixes Goal completion rungs, which printed `Nothing selected`.
Update `resume.NextLine`, the board's word accent, tests, Design section 8,
O-014's Confirmed Design Decisions, and AGENTS.md's verb-to-skill mapping
(live and scaffold).
