# Agent State Machine

This file routes the agent. The active skill is the canonical workflow source for each phase; this router records state and next action only.

## Read order

1. This file
2. The matching skill for `state`
3. Active epic detail, when `epic` is set
4. Active task or defect file, when `task` or `defect` is set
5. Files explicitly listed by the active task or skill

Read `.savepoint/PRD.md` only for vision changes. Read `.savepoint/Design.md` only for architecture or audit work. Read `.savepoint/releases/{release}/{release}-PRD.md` only for release planning or epic-order questions.

## Current state

```yaml
state: pre-implementation
release: v{{RELEASE_NUMBER}}
epic: none
task: none
next_action: "The project has its PRD and Design locked but no epics defined yet. Help the user define the epics list and confirm priority."
```

## Skill Activation

| State | Skill |
|-------|-------|
| pre-implementation | savepoint-draft-prd or savepoint-create-plan |
| epic-design | savepoint-system-design |
| epic-task-breakdown | savepoint-create-task |
| task-building | savepoint-build-task |
| audit-pending | savepoint-audit-epic |
| defect-building | savepoint-build-task |

`savepoint-audit-task` has no row: it is selected by `task-building` plus an explicit task audit or re-audit request. See `## Manual Overrides`.

## State Meanings

- `pre-implementation`: PRD/design or release epic order needs definition.
- `epic-design`: active epic detail needs design.
- `epic-task-breakdown`: active epic needs planned task files.
- `task-building`: active task is being implemented.
- `defect-building`: active release defect is being repaired.
- `audit-pending`: completed epic needs fresh-session audit.

## Manual Overrides

- If the user explicitly asks to audit or re-audit a completed epic, use `savepoint-audit-epic` even if the router is not `audit-pending`.
- If the user explicitly asks to audit or re-audit one in-progress task, use `savepoint-audit-task` and leave `state: task-building` unchanged. This is a request-qualified skill override, not a new state; the review is read-only and writes no audit file.
- If the user reports a concrete bug, regression, or broken expectation and asks to record it, use `savepoint-create-defect`.
- Audit writes exactly one `.savepoint/releases/{release}/epics/{E##-epic}/E##-Audit.md` before applying any proposals.
- Audit files keep `## Proposed Changes` — admin/apply metadata outside the visible findings sections.
- During audit apply, Update `E##-Audit.md` visible sections so they describe the applied outcome.

## Terminology

- Router `state` is the phase.
- Task `status` is only `planned`, `in_progress`, or `done`.
- Task `stage` is required when task `status` is `in_progress`.
- Defect `status` is only `open`, `in_progress`, or `resolved`.
