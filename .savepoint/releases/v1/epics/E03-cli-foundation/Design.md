---
type: epic-design
status: active
---

# Epic E03: cli-foundation

## Context Budget

For implementation tasks in this epic, read only:

- `.savepoint/router.md`
- this epic `Design.md`
- the active task file
- directly touched source/test files

Read `.savepoint/Design.md` only if the task changes architecture. Read `.savepoint/releases/v1/PRD.md`, prior epic docs, audit proposals, or `.savepoint/visual-identity.md` only when the active task explicitly requires them.

## Purpose

Introduce the real `savepoint` command shell: argument parsing, command dispatch, help/version output, exit codes, and terminal capability detection. This epic wires a command framework around stub command handlers so later epics can fill behavior without changing the CLI contract.

## What this epic adds

- Command dispatch for `init`, `board`, `audit`, and `doctor`.
- Top-level `--help` and `--version`.
- Command-level help text.
- Unknown command and unknown flag handling.
- Exit code conventions.
- TTY and color capability detection helpers.
- A testable CLI runner separated from process globals.

## Implementation strategy

- Use a small local parser instead of a CLI framework unless implementation proves it would reduce code.
- Treat the allowed commands as a closed set: `init`, `board`, `audit`, `doctor`.
- Treat allowed global flags as a closed set: `--help`, `-h`, `--version`, `-v`.
- Keep command stubs shallow. They should return deterministic "not implemented yet" results and exit codes without reading or writing project files.
- Keep process globals in `src/cli.ts`; tests should exercise `runCli()` with injected argv, streams, and environment.

## Components and files

Expected files introduced or extended by this epic:

| Path                     | Purpose                                            |
| ------------------------ | -------------------------------------------------- |
| `src/cli.ts`             | Process entrypoint that invokes the CLI runner.    |
| `src/cli/run.ts`         | Testable command runner.                           |
| `src/cli/args.ts`        | Argument parsing and normalization.                |
| `src/cli/help.ts`        | Help text generation.                              |
| `src/cli/exit-codes.ts`  | Shared exit code constants.                        |
| `src/cli/environment.ts` | TTY, color, and platform detection.                |
| `src/commands/*.ts`      | Stub command modules for the five-command surface. |
| `test/cli/**/*.test.ts`  | CLI parser and dispatch tests.                     |

## Architectural delta

Before this epic, the binary is a placeholder. After this epic, `savepoint` has a stable command boundary with stubbed behaviors behind it.

The command layer should not introduce data-model behavior. Commands that need project data stay as deterministic stubs until a later epic owns the data model.

## Boundaries

In scope:

- Preserve the exact v1 CLI surface from the project Design.
- Make bare `savepoint` print help.
- Make command stubs return clear "not implemented yet" output where necessary.
- Keep command functions testable without spawning Node processes.

Out of scope:

- Implementing `init`, `board`, `audit`, or `doctor` behavior.
- Adding extra commands such as `task new`, `plan`, or `status`.
- TUI rendering.
- Project file writes.

## Quality gates

- Tests must cover bare invocation, help, version, unknown commands, unknown flags, and each command dispatch path.
- CLI behavior should be deterministic with injected `argv`, `stdin`, `stdout`, `stderr`, and environment data.
- Every branch in argument parsing and command dispatch must have a focused test.

## Design constraints

- Keep process I/O at the entrypoint boundary.
- Avoid a heavy CLI framework unless it materially reduces complexity.
- Command names are fixed: `init`, `board`, `audit`, `doctor`.

## Close Criteria

- All `E03-cli-foundation` tasks are `done`.
- `npm run build`, `npm run typecheck`, `npm run lint`, `npm run format:check`, and `npm test` pass.
- CLI branch tests cover bare invocation, help, version, unknown commands, unknown flags, and each command dispatch path.
- Audit snapshot exists at `.savepoint/audit/E03-cli-foundation/snapshot.md`.
- Audit proposals are accepted, rejected, or explicitly carried forward.
- This epic `Design.md` has `status: audited`.
- Router points to the next epic state.
