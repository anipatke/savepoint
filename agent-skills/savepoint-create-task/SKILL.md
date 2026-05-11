---
name: savepoint-create-task
description: Plans Savepoint task files during epic-task-breakdown by writing acceptance criteria, implementation checklists, dependencies, and context-log shells.
---

# Savepoint Skill: Create Task

## Purpose

Turn an epic design into independently buildable task files with explicit acceptance criteria, scoped context files, dependencies, and implementation checklists.

## Trigger

Use this skill when router `state` is `epic-task-breakdown`.

## Read

- `.savepoint/router.md`
- Active epic detail file
- Existing task files for the active epic only when needed to order dependencies

## Workflow

1. Read the active epic and any existing task plans for that epic.
2. Create task files at `.savepoint/releases/{release}/epics/{E##-slug}/tasks/TNNN-slug.md`.
3. Use frontmatter with `id`, `status: planned`, `objective`, and `depends_on`.
4. Add exact `## Context Files`; no globs or directory-only entries.
5. Add observable `## Acceptance Criteria` before `## Implementation Plan`.
6. Add a `## Context Log` shell.
7. Update the router to `task-building` for the first unblocked planned task and stop for approval.

## Artifact Template

Write `.savepoint/releases/{release}/epics/{E##-slug}/tasks/T###-slug.md` with this structure:

```markdown
---
id: E##-slug/T###-slug
status: planned
objective: One-sentence build outcome
depends_on: []
---

# T###: Task Title

## Problem

The concrete gap this task closes.

## Context Files

- `path/to/file.ext`

## Acceptance Criteria

- [ ] Observable outcome

## Implementation Plan

- [ ] Implementation step

## Context Log

Pending.
```

## Rules

- Do not write product code.
- Do not set task `status` to `in_progress` during planning.
- Keep each task isolated and buildable.
- Use `state` only for router phase, task `status` only for task lifecycle, and `stage` only when an item is `in_progress`.
