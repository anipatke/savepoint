---
id: E06-tui-board/T006-board-integration-audit-entry
status: done
objective: "Integrate the rich board command end to end and add the initial audit-review entry point when audit proposals are present."
depends_on:
  - E06-tui-board/T005-ink-board-ui
---

# T006: board-integration-audit-entry

## Implementation Plan

- [x] Update `runBoard()` so TTY sessions launch the Ink app and non-TTY sessions use the plain table fallback from T001.
- [x] Ensure command execution handles missing project files, task graph errors, config errors, refresh failures, write conflicts, and user exit with clear stdout/stderr behavior.
- [x] Add audit proposal discovery for the active epic and expose an initial board entry point when `.savepoint/audit/{epic}/proposals.md` exists.
- [x] Keep audit-review behavior limited to entry signaling and navigation handoff; leave full proposal diff review to `E07-audit-pipeline`.
- [x] Add integration tests for TTY vs non-TTY dispatch, fallback rendering, audit-entry detection, and unchanged help/unknown-flag behavior.
- [x] Run the focused TUI/CLI tests first, then run typecheck, lint, format check, build, and the smoke suite before moving the task to review.

## Context Log

- Files read: `.savepoint/router.md`; `.savepoint/releases/v1/epics/E06-tui-board/Design.md`; `.savepoint/releases/v1/epics/E06-tui-board/tasks/T006-board-integration-audit-entry.md`; `.savepoint/metrics/context-bench.md`; `package.json`; `vitest.config.js`; `src/commands/board.ts`; `src/cli/run.ts`; `src/cli/environment.ts`; `src/cli/help.ts`; `src/commands/audit.ts`; `src/tui/App.tsx`; `src/tui/board-data.ts`; `src/tui/Board.tsx`; `src/tui/DetailPane.tsx`; `src/tui/io/gates.ts`; `src/tui/io/write-status.ts`; `src/tui/render/plain-table.ts`; `src/tui/state/app-reducer.ts`; `src/tui/state/reducer.ts`; `src/tui/state/view-state.ts`; `src/tui/theme/capability.ts`; `src/tui/theme/palette.ts`; `src/fs/project.ts`; `test/cli/run.test.ts`; `test/commands/board.test.ts`; `test/commands/board-tty.test.ts`; `test/tui/state/app-reducer.test.ts`; `test/tui/render/plain-table.test.ts`; `test/tui/state/reducer.test.ts`
- Estimated input tokens: ~28,568
- Notes: Implemented audit proposal discovery, plain-table signal, Ink banner and key handler, and handoff message. Added board-tty integration tests with mocked ink. All quality gates passed.
