---
id: E05-init-command/T001-init-cli-contract
status: planned
objective: "Define and test the `savepoint init` CLI option contract without changing scaffold behavior yet."
depends_on: []
---

# T001: init-cli-contract

## Implementation Plan

- [ ] Extend CLI argument parsing so `savepoint init` can accept only the init-specific flags required by this epic.
- [ ] Update command help for `savepoint init` to document the supported init usage and flags.
- [ ] Pass parsed init options and command execution context into `runInit()` without breaking other command handlers.
- [ ] Keep the init command returning the current not-implemented result until later tasks wire behavior.
- [ ] Add focused CLI parser and runner tests for valid init flags, invalid init flags, help output, and unchanged non-init command behavior.
