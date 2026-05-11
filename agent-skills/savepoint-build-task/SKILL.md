---
name: savepoint-build-task
description: Executes Savepoint task-building work when .savepoint/router.md state is task-building, including implementing one active task, checking acceptance criteria, running quality gates, and stopping for user review.
---

# Savepoint Skill: Build Task

## Purpose

Implement one active task exactly as planned, prove each acceptance criterion, run the quality gates, and hand control back to the user.

## Trigger

Use this skill when router `state` is `task-building`, or `defect-building` for a release defect repair.

## Read

- `.savepoint/router.md`
- Active epic detail file or active defect file
- Active task file when in `task-building`
- Only the files listed in the task `## Context Files`, unless the task itself requires a targeted verification read

## Workflow

1. Read the task acceptance criteria and implementation plan.
2. Set the task frontmatter to `status: in_progress` and `stage: build` when starting.
3. Mark the focused task as router priority in the TUI when available.
4. Implement the checklist in order and tick completed items.
5. Verify every acceptance criterion with a concrete outcome.
6. Run `make build && make test`.
7. Fill the task `## Context Log` with files read/edited and quality-gate results.
8. Add `## Drift Notes` only if files/modules or architecture changed beyond the documented map/design.
9. Stop for user review; only the user may mark the task `done`.

## Rules

- Stay within the active task scope.
- Do not audit the epic you just built.
- Do not use Savepoint CLI commands; edit files directly.
- Use `state` only for router phase, task `status` only for task lifecycle, and `stage` only when an item is `in_progress`.
