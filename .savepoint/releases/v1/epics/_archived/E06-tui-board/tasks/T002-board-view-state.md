---
id: E06-tui-board/T002-board-view-state
status: done
objective: "Create pure board view-state reducers for column grouping, selected-task movement, refresh replacement, and keyboard-intent actions."
depends_on:
  - E06-tui-board/T001-board-command-data
---

# T002: board-view-state

## Implementation Plan

- [x] Add `src/tui/state/` modules that derive board columns and task indexes from the T001 board data model without storing duplicate workflow truth.
- [x] Model selection state so it survives same-task refreshes, falls back predictably when tasks disappear, and handles empty columns.
- [x] Implement reducer actions for left/right/up/down movement plus vim-style `h`, `j`, `k`, and `l` intents.
- [x] Add refresh replacement behavior that swaps in newly read task data while preserving selection and visible status columns.
- [x] Keep transition and filesystem write behavior out of this reducer task so the reducer remains pure.
- [x] Add focused reducer tests for initial selection, horizontal and vertical navigation, empty columns, and refresh after task reorder/removal.

## Context Log

- Files read: `.savepoint/router.md`; `.savepoint/releases/v1/epics/E06-tui-board/Design.md`; `.savepoint/releases/v1/epics/E06-tui-board/tasks/T002-board-view-state.md`; `.savepoint/releases/v1/epics/E06-tui-board/tasks/T001-board-command-data.md`; `src/tui/board-data.ts`; `src/tui/render/plain-table.ts`; `src/domain/status.ts`; `src/domain/task.ts`; `src/commands/board.ts`; `test/domain/status.test.ts`; `package.json`; `tsconfig.json`; `.savepoint/metrics/context-bench.md`
- Estimated input tokens: ~6,925
- Notes: Rehydrated from router, E06 Design, task file, and directly touched source/test files. Implemented pure reducer with deterministic selection fallbacks and refresh behavior.
