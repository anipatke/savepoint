---
name: savepoint-create-plan
description: Creates Savepoint release epics from the PRD and Design documents when the router is in pre-implementation planning.
---

# Savepoint Skill: Create Plan

## Purpose

Convert approved product and architecture scope into an ordered release plan of independently valuable epics.

## Trigger

Use this skill when router `state` is `pre-implementation` and the PRD/design are approved but release epics are not defined.

## Read

- `.savepoint/router.md`
- `.savepoint/PRD.md`
- `.savepoint/Design.md`
- `.savepoint/releases/{release}/{release}-PRD.md` when release scope already exists

## Workflow

1. Read the approved scope and current release context.
2. Define ordered epics with clear deliverable value and dependencies.
3. Create or update epic detail stubs under `.savepoint/releases/{release}/epics/{E##-slug}/`.
4. Keep task files high-level unless the router has moved to `epic-task-breakdown`.
5. Update the router to the first epic that needs `epic-design` or `epic-task-breakdown`.

## Rules

- Do not write product code.
- Do not write detailed task acceptance criteria or implementation plans.
- Keep epics small enough for task-level build/audit cycles.
- Use `state` only for router phase, task `status` only for task lifecycle, and `stage` only when an item is `in_progress`.
