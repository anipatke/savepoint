# Agent State Machine

This file routes the active V2 agent workflow. The matching skill is the
canonical workflow source; this router records state and next action only.

## Read order

1. This file
2. The matching skill for `state`
3. `.savepoint/Idea.md`, when it exists, and `.savepoint/Design.md`
4. The active Objective, when `objective` is set
5. Files explicitly listed by the active Task

## Current state

```yaml
state: task
release: R-006
objective: O-021
task: T-026
next_action: "T-025 is done (owner accepted C-913). Next: execute T-026 (reconcile guidance with the smaller code base), then the mandatory Full Objective Check for O-021."
```

## State → action

| State | Skill | Next action |
| --- | --- | --- |
| `idea` | `savepoint-idea` | Capture intent and boundaries in `.savepoint/Idea.md`. |
| `design` | `savepoint-design` | Reconcile architecture, guardrails, and Objective plan. |
| `task` | `savepoint-task` | Execute the active Task within its Context Files; choose an optional Task Check or route evidence to the mandatory Full Objective Check. |
| `check` | `savepoint-check` | Independently verify a requested Task or the mandatory Objective/Goal scope and record immutable evidence. |

`REPLAN REQUIRED` is not a fifth state. It routes the current plan back to
`design` while preserving partial work and the executor's current lifecycle.

## V2 lifecycle

- Task `status` is `planned`, `in_progress`, or `done`.
- Task `stage` is required only for `in_progress`: `build` → `test` → `audit`.
- A Task's `done` transition remains owner-authorized; a Check does not silently
  close it.
- A Task Check is optional and may be skipped only with an explicit owner
  waiver recorded in the Task evidence; that waiver is not technical `CLEAR`.
- A Full Objective Check is mandatory before Objective completion, and a
  Goal Check is mandatory whenever a Goal exists.
- Issues are durable follow-up records, not a fourth task column or router
  state. A checker closes a proven Issue as `verified`; the owner may explicitly close one as `accepted`, and the planner may close a promoted Issue as `escalated`.
- Goals are optional delivery contexts, stored as `R-###` Release records.
  Their membership and completion are derived by `internal/data`; they do not
  publish, deploy, tag, or create changelogs.

## Migration boundary

The V1 source hierarchy and retired skills remain under `.savepoint/archive/v1/`
or the V1 scaffold. They are not active routing inputs for this schema-2
repository.
