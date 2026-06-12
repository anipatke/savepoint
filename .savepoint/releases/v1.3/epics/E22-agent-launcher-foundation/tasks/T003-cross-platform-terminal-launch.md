---
id: E22-agent-launcher-foundation/T003-cross-platform-terminal-launch
status: planned
objective: Launch configured agents in new interactive terminals on supported desktop platforms.
depends_on:
  - E22-agent-launcher-foundation/T001-launcher-config-contract
  - E22-agent-launcher-foundation/T002-launch-request-and-prompts
complexity_tier: high
complexity_reason: Platform-specific detached terminal startup and argument quoting carry cross-platform regression risk.
---

# T003: Cross-Platform Terminal Launch

## Problem

Starting an interactive CLI in a separate terminal differs across Windows, macOS, and Linux, and unsafe shell interpolation would corrupt prompts or create command-injection risk.

## Context Files

- `internal/launcher/terminal.go`
- `internal/launcher/terminal_test.go`
- `internal/launcher/terminal_windows.go`
- `internal/launcher/terminal_windows_test.go`
- `internal/launcher/terminal_unix.go`
- `internal/launcher/terminal_unix_test.go`
- `go.mod`
- `go.sum`

## Acceptance Criteria

- [ ] Launches use structured executable and argument values rather than concatenated shell commands.
- [ ] The agent process starts in the project root inside a separate interactive terminal.
- [ ] Windows uses a new console or supported terminal host without reusing the board console.
- [ ] macOS and Linux use documented automatic terminal candidates with a configured override fallback.
- [ ] Missing agent executables, missing terminal hosts, unsupported modes, and startup failures return actionable errors.
- [ ] Platform choice and process startup are injectable for unit tests.

## Implementation Plan

- [ ] Define a terminal-launch interface and normalized launch specification.
- [ ] Implement executable lookup and placeholder expansion before process startup.
- [ ] Add Windows-specific detached console behavior behind build tags.
- [ ] Add macOS/Linux terminal candidate selection and explicit override behavior behind build tags.
- [ ] Add tests for paths and prompts containing spaces, quotes, newlines, and non-ASCII text.
- [ ] Document any unsupported terminal environment discovered during implementation in epic Drift Notes.

## Context Log

Pending.
