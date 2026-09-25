---
id: E22-agent-launcher-foundation/T004-launcher-service
status: planned
objective: Expose a small launcher service that validates, prepares, and starts configured agent sessions.
depends_on:
  - E22-agent-launcher-foundation/T002-launch-request-and-prompts
  - E22-agent-launcher-foundation/T003-cross-platform-terminal-launch
complexity_tier: medium
complexity_reason: The service coordinates config, prompts, and terminal startup across clear package boundaries.
---

# T004: Launcher Service

## Problem

Board code needs one testable entrypoint that resolves the correct role, expands prompt placeholders, and starts a terminal without knowing platform details.

## Context Files

- `internal/launcher/launcher.go`
- `internal/launcher/launcher_test.go`
- `internal/launcher/request.go`
- `internal/launcher/prompt.go`
- `internal/launcher/terminal.go`
- `internal/data/config.go`

## Acceptance Criteria

- [ ] The service rejects launches when disabled before any process action occurs.
- [ ] Build actions resolve only the builder profile and audit actions resolve only the auditor profile.
- [ ] Missing optional auditor configuration returns an unavailable result rather than a process error.
- [ ] Prompt and project placeholders expand as single argument values.
- [ ] Successful launch results identify the action and selected item for board status messages.
- [ ] The package has no dependency on Bubble Tea or board lifecycle writes.

## Implementation Plan

- [ ] Add a launcher service with injected terminal and executable dependencies.
- [ ] Separate action availability checks from launch execution.
- [ ] Resolve role profiles and expand only the documented placeholders.
- [ ] Return typed unavailable, validation, and startup errors.
- [ ] Add focused service tests covering every role/action combination and failure boundary.

## Context Log

Pending.
