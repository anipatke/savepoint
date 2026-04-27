# Prompt: Task Breakdown

<!-- AGENT: You have an epic Design. Break it into independently buildable tasks with depends_on and checkbox plans. Stop when every task file is written and the router is updated to task-building for the first unblocked task. -->

You are breaking an approved epic into a set of tasks that agents can build one at a time.

## Input

- `.savepoint/releases/v{N}/epics/{E##-slug}/Design.md` — the epic to break down.
- `.savepoint/Design.md` — project architecture (for boundary checks).

## Output

One task file per task at `.savepoint/releases/v{N}/epics/{E##-slug}/tasks/TNNN-slug.md`:

```
---
id: {E##-slug}/TNNN-slug
status: planned
objective: "<one sentence>"
depends_on: []
---

# TNNN: slug

## Implementation Plan

- [ ] Step one
- [ ] Step two
```

## Rules

1. **Each task must be independently buildable.** A task should compile, pass its own tests, and leave the repo in a valid state when done.
2. **Objective is one sentence.** If you need more words, split the task.
3. **depends_on must be explicit.** List full task IDs. Avoid circular dependencies.
4. **Checkboxes are the plan.** Every task needs an `## Implementation Plan` with `- [ ]` items before any code is written.
5. **Order tasks by dependency.** The first task should have no `depends_on` (or only depend on prior epic tasks).
6. **Update the router.** After all tasks are written, set `.savepoint/router.md` `state: task-building` and point `task` to the first unblocked task.
7. **Do not implement.** Stop at planned tasks with unchecked boxes.
