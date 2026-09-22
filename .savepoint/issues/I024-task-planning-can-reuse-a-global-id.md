---
id: I024
title: Task planning can reuse an existing global ID
type: defect
status: open
source:
  kind: report
  actor: {role: planner, session: task-id-collision-review-20260922}
  at: '2026-09-22T09:56:33Z'
severity: high
history:
  - at: '2026-09-22T09:56:33Z'
    actor: {role: planner, session: task-id-collision-review-20260922}
    kind: observed
    note: I018 remediation allocated T006 by scanning only O012, colliding with O013's existing T006 and making strict V2 reload fail until the Task was renamed T009.
---

# I024: Task planning can reuse an existing global ID

## Summary

The design workflow can create a Task ID that already belongs to another
Objective. During I018 remediation, the planner saw O012's local sequence
T003-T005 and selected T006 without scanning the rest of the active project.
O013 already owned T006, so the watcher reload refused the entire new index and
showed a duplicate-record diagnostic until the new Task was renamed T009.

This was not only an agent typo. The canonical design skill explicitly defines
global first-unused allocation for Releases but gives no equivalent instruction
for Tasks, and it does not require strict V2 loading immediately after creating
an identity-bearing record. The loader correctly fails closed, but only after
an invalid file has reached the live project.

## Evidence

- `agent-skills/savepoint-design/SKILL.md:65` explicitly requires global
  first-unused allocation for `R###`, but its Task artifact/workflow sections
  provide no global allocation procedure for `T###`.
- The scaffold copy has the same omission.
- `internal/data/discover.go:440` correctly rejects duplicate Task identities
  across Objectives and names both source paths.
- The I018 subagent reported choosing T006 from O012's local sequence and not
  running strict loading immediately after record creation.
- A global active-record scan showed T001-T008 already occupied; renaming the
  remediation Task to T009 restored strict loading.

## Proof Needed

- Define one Task-ID allocation procedure that scans the complete active V2
  index and chooses the first unused global identity; never infer the next ID
  from the current Objective's directory alone.
- Require planners to run strict V2 loading immediately after creating or
  renaming any Task, Objective, Release, Check, or Issue identity-bearing
  record, before doing further work or allowing the watcher to surface a
  prolonged invalid state.
- Reconcile the canonical design/check/Issue-capture guidance and byte-identical
  V2 scaffold copies without duplicating allocation policy in conflicting
  places.
- Add agent/workflow scenarios that create Tasks across multiple Objectives
  and fail if a planner reuses an occupied global ID.
- Keep the loader's fail-closed duplicate diagnostic and prove that both paths
  remain named; prevention complements rather than weakens strict validation.
- Coordinate with I022 so the same allocation guarantee applies after IDs move
  to the hyphenated `T-###` form.
