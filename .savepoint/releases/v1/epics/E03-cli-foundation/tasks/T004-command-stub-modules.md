---
id: E03-cli-foundation/T004-command-stub-modules
status: done
objective: "Create deterministic stub handlers for the fixed command surface without project file side effects."
depends_on:
  - E03-cli-foundation/T001-argument-parser-contract
---

# T004: Command Stub Modules

## Scope

Create command modules for the v1 surface so dispatch can call stable handler boundaries before the real command behavior exists.

## Implementation Plan

- [x] Read `.savepoint/router.md`, this epic `Design.md`, this task file, and the directly touched command/test files.
- [x] Add a small command result type if needed by the stub handlers and runner.
- [x] Add stub modules under `src/commands/` for `init`, `board`, `audit`, and `doctor`.
- [x] Make each stub return deterministic "not implemented yet" output and the shared not-implemented exit code.
- [x] Ensure stubs do not read, write, or discover project files.
- [x] Add focused tests for each command stub result.
- [x] Run the focused command-stub test file.

## Acceptance Criteria

- Each fixed command has a handler module with a testable function.
- Stub outputs are clear and deterministic.
- Command stubs have no project filesystem side effects.
