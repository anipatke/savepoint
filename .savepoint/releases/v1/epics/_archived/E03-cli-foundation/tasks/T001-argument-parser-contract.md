---
id: E03-cli-foundation/T001-argument-parser-contract
status: done
objective: "Define the closed CLI argument parser and exit-code contract for the v1 command surface."
depends_on: []
---

# T001: Argument Parser Contract

## Scope

Create the argument parsing foundation for `savepoint` without wiring command behavior yet. This task owns the closed set of global flags, command names, parser result shapes, and shared exit-code constants.

## Implementation Plan

- [x] Read `.savepoint/router.md`, this epic `Design.md`, this task file, and the directly touched CLI/test files.
- [x] Add `src/cli/exit-codes.ts` with named exit-code constants for success, usage errors, and not-yet-implemented command stubs.
- [x] Add `src/cli/args.ts` with typed parser results for bare invocation, global help/version flags, command invocation, unknown commands, and unknown flags.
- [x] Keep allowed commands closed to `init`, `board`, `audit`, and `doctor`.
- [x] Treat `--help`, `-h`, `--version`, and `-v` as the only allowed global flags.
- [x] Add focused parser tests covering bare invocation, each global flag alias, each known command, unknown command, and unknown flag branches.
- [x] Run the focused parser test file.

## Acceptance Criteria

- The parser is deterministic from an injected `argv` array.
- Unknown commands and unknown flags are represented as typed parser results, not thrown exceptions.
- Every parser branch introduced by this task has a focused test.
