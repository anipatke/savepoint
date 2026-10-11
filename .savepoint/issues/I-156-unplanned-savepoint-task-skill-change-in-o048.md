---
id: I-156
title: O-048 working tree changes the savepoint-task skill, which the Objective puts out of scope
type: drift
status: resolved
source:
  kind: check
  check: C-980
  actor: {role: checker, session: check-o048-20261011}
  at: '2026-10-11T05:45:00Z'
tasks: [T-137, T-138]
checks: [C-980]
guardrail_ids: [TPL-02]
severity: low
resolution:
  disposition: accepted
  actor: {role: owner, session: owner-20261011}
  at: '2026-10-11T06:30:00Z'
  reason: Owner keeps the stage-writing step; O-048 Boundaries now name it as in scope.
history:
  - at: '2026-10-11T05:45:00Z'
    actor: {role: checker, session: check-o048-20261011}
    kind: observed
    check: C-980
    note: Found by the O-048 Full Objective Check.
  - at: '2026-10-11T06:30:00Z'
    actor: {role: owner, session: owner-20261011}
    kind: owner_decision
    note: Accepted by owner decision; Objective boundary updated.
---

# I-156: O-048 working tree changes the savepoint-task skill, which the Objective puts out of scope

## Summary

`agent-skills/savepoint-task/SKILL.md` and its shipped copy `templates/project-v2/agent-skills/savepoint-task/SKILL.md` gain a new Workflow step 5 (write `stage: test` before the gates, `stage: audit` after, `stage: build` again on a failed gate) and renumber steps 5-8. O-048's Boundaries list "changes to ... the workflow skills" as out of scope, T-137 and T-138 do not plan it, and neither Task's evidence records it. Design §8 does not mention it.

Violated requirement: O-048 Boundaries (out of scope: workflow skills); Task evidence completeness (unrecorded file change).

## Evidence

- `git diff HEAD -- agent-skills/savepoint-task/SKILL.md`: +5 / -4 lines, the two copies are byte-identical (so TPL-02 sync holds).
- Not listed under "Files changed" in T-137 or T-138 Technical Evidence.
- The change is plausible support for the live tile (agents must write the stage for the board to show it), but it changes every executor's workflow and was never planned or owner-approved in the records.

## Proof Needed

Either the owner accepts the change as part of O-048 (recorded as an Issue resolution or by the planner extending the Objective boundary and Design §8), or it is reverted from this Objective and planned separately.
