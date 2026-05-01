---
id: E06-tui-board/T001-board-command-data
status: done
objective: "Replace the board stub with a data-backed command path that can read the active Savepoint project and render deterministic plain output."
depends_on: []
---

# T001: board-command-data

## Implementation Plan

- [x] Add the command context needed by `runBoard()` to know the current working directory, stdout TTY capability, environment, and platform without breaking existing command handlers.
- [x] Resolve the Savepoint project root from the command context and surface boundary errors through existing CLI exit-code conventions.
- [x] Load router state, config, and the active epic task set through existing readers instead of parsing workflow files inside the command.
- [x] Introduce a small board data model that groups tasks by the five canonical statuses and records the active release, epic, task, and reader errors.
- [x] Implement `src/tui/render/plain-table.ts` as the deterministic non-TTY fallback for columns, selected task marker, and empty states.
- [x] Update board command and runner tests to prove `savepoint board` no longer returns the not-implemented stub in an initialized project and still reports readable errors outside a Savepoint project.

## Context Log

- Files read: `.savepoint/router.md`; `.savepoint/releases/v1/epics/E06-tui-board/Design.md`; `.savepoint/releases/v1/epics/E06-tui-board/tasks/T001-board-command-data.md`; `.savepoint/metrics/context-bench.md`; `src/commands/board.ts`; `src/cli/run.ts`; `src/cli/environment.ts`; `src/domain/config.ts`; `src/domain/status.ts`; `src/domain/task.ts`; `src/fs/markdown.ts`; `src/fs/project.ts`; `src/readers/config.ts`; `src/readers/router.ts`; `src/readers/tasks.ts`; `src/validation/dependencies.ts`; `test/cli/run.test.ts`; `test/domain/status.test.ts`; `package.json`; `src/cli.ts`; `src/commands/audit.ts`; `src/commands/doctor.ts`; `src/commands/init.ts`; `test/init/validate-target.test.ts`; `test/readers/router.test.ts`; `test/domain/router.test.ts`; `test/cli/commands.test.ts`; `src/domain/router.ts`; `src/domain/ids.ts`; `src/cli/exit-codes.ts`; `src/commands/result.ts`
- Estimated input tokens: ~19,660
- Notes: Rehydrated from router, E06 Design, task file, and directly touched source/test files. Extended router schema with optional `task` field to match live `router.md` format. Board command now async with rich result shape matching `runInit` pattern.
