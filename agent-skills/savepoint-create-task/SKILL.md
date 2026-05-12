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

## Complexity Rubric

Assign one tier per task using the highest tier that applies to any step:

| Tier | Meaning |
|------|---------|
| `low` | Single-file or well-bounded change. No cross-module ripple. No design decisions. |
| `medium` | Multi-file change with some design choices. Moderate cross-module impact. |
| `high` | Cross-module redesign, significant interface changes, or many coordinated steps with regression risk. |
| `spike` | Scope or approach is unknown or exploratory. Use when the problem is not well-understood regardless of estimated size. |

**Highest-tier rule:** when steps span multiple tiers, record the highest.  
**Spike-override rule:** if any step has uncertain scope, always use `spike`.

The `complexity_reason` must be a single sentence (≤120 characters) stating why the tier was chosen.

## Workflow

1. Read the active epic and any existing task plans for that epic.
2. Create task files at `.savepoint/releases/{release}/epics/{E##-slug}/tasks/TNNN-slug.md`.
3. Use frontmatter with `id`, `status: planned`, `objective`, `depends_on`, `complexity_tier`, and `complexity_reason`.
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
complexity_tier: low|medium|high|spike
complexity_reason: One sentence ≤120 chars explaining the tier choice.
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
- Always include `complexity_tier` and `complexity_reason` in every new task file.
