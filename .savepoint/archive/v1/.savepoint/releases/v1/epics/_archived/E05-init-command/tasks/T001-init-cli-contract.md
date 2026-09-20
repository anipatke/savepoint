---
id: E05-init-command/T001-init-cli-contract
status: done
objective: "Define and test the `savepoint init` CLI option contract without changing scaffold behavior yet."
depends_on: []
---

# T001: init-cli-contract

## Implementation Plan

- [x] Extend CLI argument parsing so `savepoint init` can accept only the init-specific flags required by this epic.
- [x] Update command help for `savepoint init` to document the supported init usage and flags.
- [x] Pass parsed init options and command execution context into `runInit()` without breaking other command handlers.
- [x] Keep the init command returning the current not-implemented result until later tasks wire behavior.
- [x] Add focused CLI parser and runner tests for valid init flags, invalid init flags, help output, and unchanged non-init command behavior.

## Context Log

- Files read: not captured before context tracking was added to task files.
- Estimated input tokens: n/a
- Notes: Use `.savepoint/metrics/context-bench.md` for measured rows going forward.
