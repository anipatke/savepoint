---
id: E06-tui-board/T005-ink-board-ui
status: done
objective: "Build the Ink board UI shell with keyboard navigation, detail pane, manual refresh, and transition actions."
depends_on:
  - E06-tui-board/T002-board-view-state
  - E06-tui-board/T003-transition-gates-and-writes
  - E06-tui-board/T004-terminal-theme
---

# T005: ink-board-ui

## Implementation Plan

- [x] Add the minimum Ink and React dependencies and TypeScript configuration needed for `src/tui/**/*.tsx` while keeping the CLI build ESM-compatible on Node 20.
- [x] Implement `src/tui/App.tsx` as the stateful app root that receives initial board data, theme, refresh/write callbacks, and terminal capability data.
- [x] Implement `src/tui/Board.tsx` with five compact status columns, selected-task focus treatment, empty-column states, and keyboard hints that do not dominate the screen.
- [x] Implement `src/tui/DetailPane.tsx` for selected task objective, status, dependencies, gate state, and last action/conflict messages.
- [x] Wire arrow keys, vim-style navigation, manual `r` refresh, and transition actions through the pure reducer and IO helpers.
- [x] Add focused tests around reducer-connected interaction behavior and component output where deterministic, avoiding brittle full-screen TUI snapshots.

## Context Log

- Files read: `.savepoint/router.md`; `.savepoint/releases/v1/epics/E06-tui-board/Design.md`; `.savepoint/releases/v1/epics/E06-tui-board/tasks/T005-ink-board-ui.md`; `.savepoint/visual-identity.md`; `.savepoint/metrics/context-bench.md`; `package.json`; `tsconfig.json`; `vitest.config.js`; `src/commands/board.ts`; `src/cli/run.ts`; `src/cli/environment.ts`; `src/domain/config.ts`; `src/domain/ids.ts`; `src/domain/status.ts`; `src/domain/task.ts`; `src/fs/project.ts`; `src/readers/tasks.ts`; `src/tui/board-data.ts`; `src/tui/state/reducer.ts`; `src/tui/state/view-state.ts`; `src/tui/theme/index.ts`; `src/tui/theme/palette.ts`; `src/tui/theme/capability.ts`; `src/tui/io/gates.ts`; `src/tui/io/write-status.ts`; `src/tui/render/plain-table.ts`; `test/tui/state/reducer.test.ts`; `test/tui/io/gates.test.ts`; `test/tui/io/write-status.test.ts`; `test/tui/render/plain-table.test.ts`; `test/tui/theme/capability.test.ts`; `test/tui/theme/palette.test.ts`; `.savepoint/releases/v1/epics/E06-tui-board/tasks/T002-board-view-state.md`; `.savepoint/releases/v1/epics/E06-tui-board/tasks/T003-transition-gates-and-writes.md`; `.savepoint/releases/v1/epics/E06-tui-board/tasks/T004-terminal-theme.md`
- Estimated input tokens: ~18,000
- Notes: Rehydrated from router, E06 Design, T005 task, prerequisite task files (T002-T004), visual-identity, and directly touched source/test files. Added Ink 7 + React 19 with JSX support, extended EnvMap with COLORTERM, and built App/Board/DetailPane components with keyboard-driven navigation and transition actions.
