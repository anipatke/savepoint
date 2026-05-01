---
id: E03-cli-foundation/T005-cli-runner-dispatch
status: done
objective: "Wire the parser, help text, environment helpers, and command stubs into a testable CLI runner."
depends_on:
  - E03-cli-foundation/T001-argument-parser-contract
  - E03-cli-foundation/T002-help-text-generation
  - E03-cli-foundation/T003-terminal-environment-detection
  - E03-cli-foundation/T004-command-stub-modules
---

# T005: CLI Runner Dispatch

## Scope

Add the `runCli()` boundary that later command implementations will use. This task owns command dispatch behavior, injected streams/environment, and result-to-output handling, but not the process entrypoint.

## Implementation Plan

- [x] Read `.savepoint/router.md`, this epic `Design.md`, this task file, and the directly touched CLI/command/test files.
- [x] Add `src/cli/run.ts` with an injectable runner input for `argv`, `stdin`, `stdout`, `stderr`, environment, and platform data.
- [x] Route bare invocation and global help flags to top-level help output with a success exit code.
- [x] Route global version flags to the existing package version source with a success exit code.
- [x] Route command-level help requests to command help output with a success exit code.
- [x] Route known commands to their stub handlers.
- [x] Route unknown commands and unknown flags to stderr with the usage-error exit code.
- [x] Add focused runner tests for bare invocation, help, version, each command dispatch path, command-level help, unknown command, and unknown flag.
- [x] Run the focused runner test file.

## Acceptance Criteria

- `runCli()` is testable without spawning Node or reading process globals.
- Output routing between stdout and stderr is deterministic.
- Every command dispatch branch required by the epic quality gates has focused coverage.
