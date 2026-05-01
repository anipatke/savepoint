---
type: audit-snapshot
epic: E06-tui-board
source: manual
created: 2026-04-27
---

# E06-tui-board Audit Snapshot

Manual snapshot created because the audit CLI is not available yet.

## Epic Scope

Implement `savepoint board`, an Ink-based terminal UI for viewing and moving tasks through the Savepoint workflow, with deterministic non-TTY output, terminal theme fallbacks, mtime-safe task status writes, and an initial audit-review entry signal.

## Completed Tasks

- E06-tui-board/T001-board-command-data
- E06-tui-board/T002-board-view-state
- E06-tui-board/T003-transition-gates-and-writes
- E06-tui-board/T004-terminal-theme
- E06-tui-board/T005-ink-board-ui
- E06-tui-board/T006-board-integration-audit-entry

## Changed Files To Audit

- `package.json`
- `package-lock.json`
- `tsconfig.json`
- `vitest.config.js`
- `src/cli/environment.ts`
- `src/cli/help.ts`
- `src/cli/run.ts`
- `src/commands/board.ts`
- `src/domain/router.ts`
- `src/fs/markdown.ts`
- `src/tui/board-data.ts`
- `src/tui/App.tsx`
- `src/tui/Board.tsx`
- `src/tui/DetailPane.tsx`
- `src/tui/io/gates.ts`
- `src/tui/io/write-status.ts`
- `src/tui/render/plain-table.ts`
- `src/tui/state/app-reducer.ts`
- `src/tui/state/reducer.ts`
- `src/tui/state/view-state.ts`
- `src/tui/theme/capability.ts`
- `src/tui/theme/index.ts`
- `src/tui/theme/palette.ts`
- `test/commands/board.test.ts`
- `test/commands/board-tty.test.ts`
- `test/tui/components/Board.test.tsx`
- `test/tui/components/DetailPane.test.tsx`
- `test/tui/io/gates.test.ts`
- `test/tui/io/write-status.test.ts`
- `test/tui/render/plain-table.test.ts`
- `test/tui/state/app-reducer.test.ts`
- `test/tui/state/reducer.test.ts`
- `test/tui/theme/capability.test.ts`
- `test/tui/theme/palette.test.ts`

## Known Non-E06 Working Tree Noise

The working tree also contains prior E05/init, template, README, AGENTS, and router changes. They are not audit inputs for this E06 proposal except where the router explicitly asks for proposal targets.
