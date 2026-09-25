---
id: E03-cli-foundation/T003-terminal-environment-detection
status: done
objective: "Provide injectable terminal capability detection for TTY, color, and platform behavior."
depends_on: []
---

# T003: Terminal Environment Detection

## Scope

Add environment capability helpers that can be used by later CLI and TUI work. This task stays independent from command dispatch and does not read project files.

## Implementation Plan

- [x] Read `.savepoint/router.md`, this epic `Design.md`, this task file, and the directly touched CLI/test files.
- [x] Add `src/cli/environment.ts` with typed inputs for stdout/stderr TTY state, environment variables, and platform.
- [x] Detect whether stdout and stderr are TTY-capable from injected stream-like values.
- [x] Detect color support from injected environment/platform data using conservative rules for `NO_COLOR`, `FORCE_COLOR`, and common CI/non-TTY cases.
- [x] Add focused tests for TTY true/false branches, color disabled/enabled branches, and platform passthrough.
- [x] Run the focused environment test file.

## Acceptance Criteria

- Capability detection is deterministic with injected streams, environment, and platform.
- No helper reads directly from `process`.
- Branching behavior for TTY and color detection is covered by focused tests.
