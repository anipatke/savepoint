## Target File

`.savepoint/Design.md`

## Replace

## 6. CLI surface (5 commands, no more)

| Command                | Purpose                                                                           |
| ---------------------- | --------------------------------------------------------------------------------- |
| `savepoint init`       | Scaffold `.savepoint/`, print magic prompt to stdout + clipboard                  |
| `savepoint board`      | Launch TUI; auto-falls-back to plain table on non-TTY                             |
| `savepoint audit`      | Run audit pipeline (`--skip --reason`, `--epic`)                                  |
| `savepoint doctor`     | Integrity check + ad-hoc quality-gate run + Layer-2 prompt for AI semantic review |
| `--version` / `--help` | Standard                                                                          |

- Bare `savepoint` prints help.
- **Explicitly rejected:** `task new`, `epic new`, `release new`, `plan`, `next`, `status`, `task done`. All are file edits or TUI actions.
- **Agents must not run `savepoint` commands.** Stated in AGENTS.md.

**Names:** npm package `savepoint`; binary `savepoint`. No `vk` alias.

## With

## 6. CLI surface (4 commands, no extras)

| Command                | Purpose                                                                           |
| ---------------------- | --------------------------------------------------------------------------------- |
| `savepoint init`       | Scaffold `.savepoint/`, print magic prompt to stdout + clipboard                  |
| `savepoint board`      | Launch TUI; auto-falls-back to plain table on non-TTY                             |
| `savepoint audit`      | Run audit pipeline (`--skip --reason`, `--epic`)                                  |
| `savepoint doctor`     | Integrity check + ad-hoc quality-gate run + Layer-2 prompt for AI semantic review |
| `--version` / `--help` | Standard global flags                                                             |

- Bare `savepoint` prints help.
- As of E03, the binary has a real command shell around stub handlers:
  - `src/cli.ts` owns process globals and delegates to `runCli()`.
  - `src/cli/args.ts` parses the fixed command and global flag surface.
  - `src/cli/run.ts` dispatches help, version, command stubs, unknown commands, and unknown top-level flags.
  - `src/cli/help.ts` generates deterministic top-level and command-level help.
  - `src/cli/environment.ts` detects TTY, color, and platform capability from injected inputs.
  - `src/commands/*.ts` contains deterministic not-yet-implemented stubs.
- Implemented global flags: `--help`, `-h`, `--version`, `-v`.
- Implemented command help flags: `--help`, `-h` after `init`, `board`, `audit`, or `doctor`.
- Command behavior remains future work. Until later epics fill it in, command stubs return `EXIT_NOT_IMPLEMENTED`.
- **Explicitly rejected:** `task new`, `epic new`, `release new`, `plan`, `next`, `status`, `task done`. All are file edits or TUI actions.
- **Agents must not run `savepoint` commands.** Stated in AGENTS.md.

**Names:** npm package `savepoint`; binary `savepoint`. No `vk` alias.

## Target File

`.savepoint/Design.md`

## Replace

`- **Read-only data model:** established in epic `E02-data-model`(2026-04-27).`js-yaml` is the only runtime dependency added for structured YAML/frontmatter parsing.`

## With

- **Read-only data model:** established in epic `E02-data-model` (2026-04-27). `js-yaml` is the only runtime dependency added for structured YAML/frontmatter parsing.
- **CLI foundation:** established in epic `E03-cli-foundation` (2026-04-27). The binary now has a testable parser, help/version handling, command dispatch, terminal capability detection, and deterministic stubs for `init`, `board`, `audit`, and `doctor`.

## Target File

`.savepoint/Design.md`

## Replace

`E02 adds focused unit coverage for ID parsing, status validation, task/release/epic/router/config frontmatter validation, markdown read-boundary failures, project root discovery, dependency graph errors, and epic task-set reading. As of the E02 audit, `npm test` reports 16 passing test files and 176 passing tests when run outside the sandbox.`

## With

E02 adds focused unit coverage for ID parsing, status validation, task/release/epic/router/config frontmatter validation, markdown read-boundary failures, project root discovery, dependency graph errors, and epic task-set reading. As of the E02 audit, `npm test` reports 16 passing test files and 176 passing tests when run outside the sandbox.

E03 adds focused CLI coverage for bare invocation, global help/version flags, command-level help, command dispatch, unknown top-level commands and flags, environment detection, and command stubs. As of the E03 audit snapshot, `npm test` reports 21 passing test files and 256 passing tests.
