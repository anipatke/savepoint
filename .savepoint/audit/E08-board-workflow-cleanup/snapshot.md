---
type: audit-snapshot
epic: E08-board-workflow-cleanup
source: manual
created: 2026-04-29
---

# E08-board-workflow-cleanup Audit Snapshot

Manual snapshot created because the audit CLI is not the active workflow entry for this epic closeout.

## Epic Scope

Extend the board from an active-epic-only view into a release workflow surface with acceptance-criteria-aware task documents, release-level task loading, epic browsing, alternate-screen TTY rendering, a task detail modal, and updated workflow templates/prompts.

## Completed Tasks

- E08-board-workflow-cleanup/T001-acceptance-criteria-model
- E08-board-workflow-cleanup/T002-release-task-set-reader
- E08-board-workflow-cleanup/T003-board-data-and-plain-output
- E08-board-workflow-cleanup/T004-board-selection-state
- E08-board-workflow-cleanup/T005-ink-board-layout-cleanup
- E08-board-workflow-cleanup/T006-task-detail-popup
- E08-board-workflow-cleanup/T007-templates-acceptance-criteria
- E08-board-workflow-cleanup/T008-board-workflow-integration

## Changed Files To Audit

- `.savepoint/router.md`
- `AGENTS.md`
- `src/commands/board.ts`
- `src/domain/task.ts`
- `src/readers/tasks.ts`
- `src/tui/App.tsx`
- `src/tui/Board.tsx`
- `src/tui/DetailPane.tsx`
- `src/tui/board-data.ts`
- `src/tui/io/gates.ts`
- `src/tui/render/plain-table.ts`
- `src/tui/state/app-reducer.ts`
- `src/tui/state/reducer.ts`
- `src/tui/state/view-state.ts`
- `templates/project/.savepoint/router.md`
- `templates/project/AGENTS.md`
- `templates/prompts/task-breakdown.prompt.md`
- `templates/prompts/task-planning.prompt.md`
- `templates/prompts/task-building.prompt.md`
- `test/commands/board.test.ts`
- `test/commands/board-tty.test.ts`
- `test/domain/task.test.ts`
- `test/readers/tasks.test.ts`
- `test/templates/project-templates.test.ts`
- `test/templates/prompt-templates.test.ts`
- `test/templates/router-template.test.ts`
- `test/tui/board-data.test.ts`
- `test/tui/components/App.test.tsx`
- `test/tui/components/Board.test.tsx`
- `test/tui/components/DetailPane.test.tsx`
- `test/tui/render/plain-table.test.ts`
- `test/tui/state/app-reducer.test.ts`
- `test/tui/state/reducer.test.ts`

## Known Non-E08 Working Tree Noise

The working tree also contains prior and parallel changes in E05, E06, E07, release-planning metadata, README/docs copies, and other scaffold/test areas. Those files are not audit inputs for this Epic08 bundle except where proposal targets explicitly update shared architecture or codebase-map documentation.
