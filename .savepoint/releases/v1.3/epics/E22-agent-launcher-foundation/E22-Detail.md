---
type: epic-design
status: planned
---

# E22: Agent Launcher Foundation

## Purpose

Create the opt-in configuration, scoped request model, prompt generation, and cross-platform terminal runtime needed to launch an interactive agent without coupling Savepoint to a vendor.

## What this epic adds

- Disabled-by-default launcher configuration with builder and optional auditor profiles.
- Structured command and argument handling without shell command-string interpolation.
- Typed launch requests for task build, task audit, defect build, defect audit, and epic audit.
- Deterministic prompts that identify the selected item, workflow action, and scope limits.
- Platform adapters that open a separate interactive terminal in the project root.
- A launcher service with testable validation and error boundaries.

## Components and files

| Module | Purpose |
|--------|---------|
| `internal/data/config.go` | Parse launcher opt-in and agent profile configuration |
| `internal/launcher` | Own launch requests, prompts, validation, and terminal process startup |
| `.savepoint/config.yml` | Hold local launcher preferences |

## Architectural delta

Savepoint gains an outbound process boundary behind a small `internal/launcher` package. The package accepts structured arguments, works only when explicitly enabled, and returns launch results to callers without owning task lifecycle or TUI state.

## Boundaries

**In scope:**
- Builder and auditor role configuration
- Structured prompt and process launch contracts
- Windows, macOS, and Linux terminal adapters
- Unit-test seams for executable lookup and process startup

**Out of scope:**
- Board keybindings and lifecycle writes
- Doctor diagnostics and scaffold documentation
- Agent output capture or process supervision
- Vendor-specific behavior

## Quality gates

- Launcher package tests cover disabled, malformed, missing-role, prompt, quoting, and platform-selection behavior.
- `go test ./internal/data ./internal/launcher` passes.
- `make build` passes for the host platform.

## Open decisions

- Confirm the smallest reliable automatic terminal set for macOS and Linux during T003; retain an explicit terminal override for unsupported environments.
