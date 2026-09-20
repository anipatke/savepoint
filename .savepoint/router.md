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
release: R006
objective: O001
task: T002
next_action: Run an independent Check for T002; then obtain T001 clearance and owner completion evidence before reevaluating R006 cutover.
```

## State → action

| State | Skill | Next action |
| --- | --- | --- |
| `idea` | `savepoint-idea` | Capture intent and boundaries in `.savepoint/Idea.md`. |
| `design` | `savepoint-design` | Reconcile architecture, guardrails, and Objective plan. |
| `task` | `savepoint-task` | Execute the active Task within its Context Files and hand it to a fresh Check. |
| `check` | `savepoint-check` | Independently verify the Task or Objective and record immutable evidence. |

`REPLAN REQUIRED` is not a fifth state. It routes the current plan back to
`design` while preserving partial work and the executor's current lifecycle.

## V2 lifecycle

- Task `status` is `planned`, `in_progress`, or `done`.
- Task `stage` is required only for `in_progress`: `build` → `test` → `audit`.
- A Task's `done` transition remains owner-authorized; a Check does not silently
  close it.
- Issues are durable follow-up records, not a fourth task column or router
  state. Only `savepoint-check` closes an Issue.
- Releases are optional delivery boundaries. Their membership and completion
  are derived by `internal/data`; they do not publish, deploy, tag, or create
  changelogs.

## Migration boundary

The V1 source hierarchy and retired skills remain under `.savepoint/archive/v1/`
or the V1 scaffold. They are not active routing inputs for this schema-2
repository.
