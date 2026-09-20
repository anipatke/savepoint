---
id: E06-tui-board/T004-terminal-theme
status: done
objective: "Load board theme tokens from config and map them into deterministic terminal color palettes for rich and fallback environments."
depends_on:
  - E06-tui-board/T001-board-command-data
---

# T004: terminal-theme

## Implementation Plan

- [x] Add `src/tui/theme/` modules that consume `SavepointConfig["theme"]` and expose board-specific semantic tokens for statuses, borders, text, focus, and blocked states.
- [x] Implement terminal color capability mapping for 24-bit, 256-color, 16-color, `NO_COLOR=1`, `FORCE_COLOR`, `TERM=dumb`, and non-TTY plain output.
- [x] Keep color fallback data deterministic and separate from Ink components so render tests can avoid full-screen snapshots.
- [x] Respect configured border style values while preserving readable defaults when config values are missing or colors are disabled.
- [x] Add tests for config defaults, custom accents, no-color mode, dumb terminals, and reduced palettes.

## Context Log

- Files read: `.savepoint/router.md`; `.savepoint/releases/v1/epics/E06-tui-board/Design.md`; `.savepoint/releases/v1/epics/E06-tui-board/tasks/T004-terminal-theme.md`; `.savepoint/releases/v1/epics/E06-tui-board/tasks/T003-transition-gates-and-writes.md`; `package.json`; `tsconfig.json`; `src/domain/config.ts`; `src/cli/environment.ts`; `src/commands/board.ts`; `src/domain/status.ts`; `src/tui/board-data.ts`; `src/tui/render/plain-table.ts`; `src/tui/state/view-state.ts`; `src/tui/state/reducer.ts`; `src/tui/io/gates.ts`; `src/tui/io/write-status.ts`; `test/tui/state/reducer.test.ts`; `test/tui/render/plain-table.test.ts`
- Estimated input tokens: ~14,752
- Notes: Implemented theme capability detection, 256/16 color approximation, and palette builder with semantic tokens. Rehydrated from router, E06 Design, T004 task, T003 task for dependency state, and directly touched source/test files.
