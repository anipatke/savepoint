# Agent State Machine

This file routes the agent. The active skill is the canonical workflow source for each state; this router records state and next action only.

## Read order

1. This file
2. The matching skill for `state`
3. `.savepoint/Idea.md`, when it exists, and `.savepoint/Design.md`
4. The active Objective, when `objective` is set
5. Files explicitly listed by the active skill or Task

## Current state

```yaml
state: idea
objective: none
task: none
next_action: "Turn a rough idea — a single sentence is a valid starting point, no prepared requirements document needed — into .savepoint/Idea.md through a short back-and-forth with the owner. If the directory already holds a codebase, the same conversation grounds intent while Design.md's Current Technical State is reconstructed from targeted reads, per AGENTS.md's Existing Codebase Adoption section: one route, two branches, not two workflows."
```

## Skill Activation

| State | Skill |
|-------|-------|
| idea | savepoint-idea |
| design | savepoint-design |
| task | savepoint-task |
| check | savepoint-check |

`REPLAN REQUIRED`, returned by `savepoint-task` on a materially invalid plan, routes back into `savepoint-design`. It is not a fifth state; `state` stays `design` while the planner resolves what broke.

Use the `skill` tool when the listed skill is available. If the agent says the skill is not found, read `agent-skills/{skill}/SKILL.md` directly and follow it as the active skill.

## State Meanings

- `idea`: intent and boundary are being defined in `.savepoint/Idea.md`.
- `design`: `Design.md`, `Guardrails.md`, and the current Objective's Tasks are being kept ready; see the Readiness Gate in `savepoint-design`.
- `task`: the active Task is being built by `savepoint-task`, within the boundaries the planner already set.
- `check`: the active Task or Objective is being independently verified by a fresh `savepoint-check` session.

## Terminology

- Router `state` is the current state: `idea`, `design`, `task`, or `check`.
- Task `status` is only `planned`, `in_progress`, or `done`.
- Task `stage` is required when task `status` is `in_progress`: `build` → `test` → `audit`.
- An Issue is a durable follow-up record entered from any of the four skills' own workflow, not a fifth state; see `agent-skills/references/issue-capture.md`.
