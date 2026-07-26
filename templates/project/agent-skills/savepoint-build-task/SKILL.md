---
name: savepoint-build-task
description: Executes Savepoint task-building work when .savepoint/router.md state is task-building, including implementing one active task, checking acceptance criteria, running quality gates, and stopping for user review.
---

# Savepoint Skill: Build Task

## Purpose

Implement one active task exactly as planned, prove each acceptance criterion, run the quality gates, and hand control back to the user.

## Trigger

Use this skill when router `state` is `task-building`, or `defect-building` for a release defect repair.

Do not use this skill to review your own work. If the user explicitly asks to audit or re-audit the task, hand off to `savepoint-audit-task`; if the epic is complete or the router is `audit-pending`, hand off to `savepoint-audit-epic`. Both audits keep their own read scope and result contract.

## Read

- `.savepoint/router.md`
- Active epic detail file or active defect file
- Active task file when in `task-building`
- Only the files listed in the task `## Context Files`, unless the task itself requires a targeted verification read

## Workflow

1. Read the task acceptance criteria and implementation plan.
2. Before editing code, verify task frontmatter uses canonical lifecycle fields: `status`, and `stage` only when `status: in_progress`.
3. Set the task frontmatter to `status: in_progress` and `stage: build` when starting.
4. Mark the focused task as router priority in the TUI when available.
5. Implement the checklist in order and tick completed items.
6. Verify every acceptance criterion with a concrete outcome.
7. Run `make build && make test`.
8. Apply the `.savepoint/Health-Check.md` Quick check at task handoff and record its evidence block in the task `## Context Log`. Skip this step when the project has no `.savepoint/Health-Check.md`.
9. Fill the task `## Context Log` with files read/edited and quality-gate results.
10. Add `## Drift Notes` only if files/modules or architecture changed beyond the documented map/design.
11. Stop for user review; only the user may mark the task `done`.

## Rules

- Stay within the active task scope.
- Do not audit the epic you just built. Defer explicit audit and re-audit requests to `savepoint-audit-task` for one in-progress task, or `savepoint-audit-epic` for a completed epic.
- Do not use Savepoint CLI commands; edit files directly.
- Use `state` only for router phase, task `status` only for task lifecycle, and `stage` only when an item is `in_progress`.
- Never write `stage: implementation`; implementation work starts at `stage: build`.
- Treat legacy task `phase` as parser compatibility only; new task guidance and writes must use `stage`.
