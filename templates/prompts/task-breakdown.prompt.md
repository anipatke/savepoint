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

## Acceptance Criteria

- [ ] Criterion one (observable outcome)
- [ ] Criterion two (observable outcome)

## Implementation Plan

- [ ] Step one
- [ ] Step two

## Context Log

- Files read:
- Estimated input tokens:
- Notes:
```

## Rules

1. **Each task must be independently buildable.** A task should compile, pass its own tests, and leave the repo in a valid state when done.
2. **Objective is one sentence.** If you need more words, split the task.
3. **depends_on must be explicit.** List full task IDs. Avoid circular dependencies.
4. **Acceptance criteria come first.** Every task needs `## Acceptance Criteria` (observable outcomes that define "done") before `## Implementation Plan` (the build checklist that satisfies them).
5. **Checkboxes are the plan.** Use inline `- [ ]` items under `## Implementation Plan` before any code is written.
6. **Use only canonical task status.** New task frontmatter must use `status: planned`. Do not write `todo`, `doing`, `blocked`, `review`, `audit`, or phase names into `status`.
7. **Use phase separately.** Only implementation flow may add `phase: build`, `phase: test`, or `phase: audit`, and only when `status: in_progress`.
8. **Order tasks by dependency.** The first task should have no `depends_on` (or only depend on prior epic tasks).
9. **Update the router.** After all tasks are written, set `.savepoint/router.md` `state: task-building` and point `task` to the first unblocked task.
10. **Do not implement.** Stop at planned tasks with unchecked boxes.
