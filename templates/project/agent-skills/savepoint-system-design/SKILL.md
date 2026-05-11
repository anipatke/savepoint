---
name: savepoint-system-design
description: Produces Savepoint system or epic design documents from the PRD when the router is in epic-design or the user asks for architectural design.
---

# Savepoint Skill: System Design

## Purpose

Turn approved product scope into a concise technical design. For project design, update `.savepoint/Design.md`. For epic design, update the active `.savepoint/releases/{release}/epics/{E##-slug}/E##-Detail.md`.

## Trigger

Use this skill when router `state` is `epic-design`, or when the user explicitly asks for architecture/design work.

## Read

- `.savepoint/router.md`
- Active PRD or release PRD only when the design depends on product scope
- Active epic detail file when designing an epic
- `.savepoint/Design.md` only for architecture context

## Workflow

1. Read the router and the design inputs needed for the active scope.
2. Describe what the scope adds, which components it touches, and the architectural delta.
3. Record firm implementation boundaries and out-of-scope items.
4. Keep task breakdown and build steps out of the design.
5. When the design is ready, update the router toward `epic-task-breakdown` for the active epic.

## Artifact Template

Write `.savepoint/releases/{release}/epics/{E##-slug}/E##-Detail.md` with this structure:

```markdown
---
type: epic-design
status: planned
---

# E##: Epic Title

## Purpose

Short statement of why this epic exists.

## What this epic adds

- User-facing or workflow capability this epic adds

## Components and files

| Module | Purpose |
|--------|---------|
| `path/or/module` | Responsibility in this epic |

## Architectural delta

What changes from the current system design.

## Boundaries

**In scope:**
- Included work

**Out of scope:**
- Excluded work

## Quality gates

- Required verification outcome

## Open decisions

None.
```

## Rules

- Do not write product code.
- Do not create detailed task implementation plans.
- Keep design documents short enough to be read every session.
- Use `state` only for router phase, task `status` only for task lifecycle, and `stage` only when an item is `in_progress`.
