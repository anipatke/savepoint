---
id: E06-tui-board/T003-transition-gates-and-writes
status: done
objective: "Add task transition gate evaluation and optimistic mtime-safe status writes for board actions."
depends_on:
  - E06-tui-board/T001-board-command-data
  - E06-tui-board/T002-board-view-state
---

# T003: transition-gates-and-writes

## Implementation Plan

- [x] Add pure gate evaluation that combines `validateTransition()` with dependency rules so tasks cannot advance when required dependencies are not `done`.
- [x] Return structured gate results that the TUI can display without duplicating workflow logic in rendering components.
- [x] Add `src/tui/io/` helpers that read a task file's mtime before action and verify it is unchanged immediately before writing.
- [x] Update only the task frontmatter `status` field while preserving the task body and unrelated frontmatter fields.
- [x] Surface stale mtime conflicts as non-destructive action results that require a manual refresh before retry.
- [x] Add tests for allowed transitions, blocked dependencies, disallowed status jumps, successful status write preservation, and stale mtime conflict handling without relying on real-time delays.

## Context Log

- Files read: `.savepoint/router.md`; `.savepoint/releases/v1/epics/E06-tui-board/Design.md`; `.savepoint/releases/v1/epics/E06-tui-board/tasks/T003-transition-gates-and-writes.md`; `.savepoint/releases/v1/epics/E06-tui-board/tasks/T001-board-command-data.md`; `.savepoint/releases/v1/epics/E06-tui-board/tasks/T002-board-view-state.md`; `.savepoint/metrics/context-bench.md`; `package.json`; `tsconfig.json`; `vitest.config.js`; `src/domain/status.ts`; `src/domain/task.ts`; `src/domain/ids.ts`; `src/validation/dependencies.ts`; `src/fs/markdown.ts`; `src/readers/tasks.ts`; `src/commands/board.ts`; `src/tui/board-data.ts`; `src/tui/state/reducer.ts`; `src/tui/state/view-state.ts`; `src/tui/render/plain-table.ts`; `test/tui/state/reducer.test.ts`; `test/tui/render/plain-table.test.ts`; `test/fs/markdown.test.ts`
- Estimated input tokens: ~19,630
- Notes: Implementation rehydrated from router, E06 Design, task file, and directly touched source/test files. Extended `MarkdownDoc` with `rawFrontmatter` to preserve extra frontmatter fields during status writes. Gate evaluator allows idempotent same-status transitions and blocks only forward moves with undone dependencies.
