---
 type: epic-design
 status: planned
 ---

 # Epic E03: CLI Simplify

 ## Purpose

 Strip the CLI down to the board command only. Remove init, audit, doctor commands and their supporting infrastructure. Update CLI tests to reflect the board-only interface.

 ## What this epic adds

 - Board-only argument parsing (`args.ts`).
 - Board-only help text (`help.ts`).
 - Board-only command dispatch (`run.ts`).
 - Deletion of init, audit, doctor command files and their supporting directories.

 ## Definition of Done

 - `savepoint --help` shows only `board` command.
 - `savepoint init/audit/doctor` returns unknown-command error.
 - `savepoint board` still launches from TTY.
 - No dead command code remains in `src/commands/`, `src/init/`, `src/templates/`.
 - All CLI tests pass with board-only interface.
 - Build and typecheck pass.

 ## Components and files

 | Path | Purpose |
 |------|---------|
 | `src/cli/args.ts` | Board-only argument parsing |
 | `src/cli/help.ts` | Board-only help text |
 | `src/cli/run.ts` | Board-only command dispatch |
 | `test/cli/args.test.ts` | Updated args tests |
 | `test/cli/help.test.ts` | Updated help tests |
 | `test/cli/run.test.ts` | Updated runner tests |
 | (deleted) `src/commands/init.ts` | Removed |
 | (deleted) `src/commands/audit.ts` | Removed |
 | (deleted) `src/commands/doctor.ts` | Removed |
 | (deleted) `src/commands/result.ts` | Removed |
 | (deleted) `src/init/` | Removed (5 files) |
 | (deleted) `src/templates/` | Removed (5 files) |
