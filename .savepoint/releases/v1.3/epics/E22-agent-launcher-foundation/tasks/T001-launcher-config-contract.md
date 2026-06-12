---
id: E22-agent-launcher-foundation/T001-launcher-config-contract
status: planned
objective: Add a disabled-by-default, vendor-neutral launcher configuration contract.
depends_on: []
complexity_tier: medium
complexity_reason: The config adds nested defaults and validation while preserving all existing config behavior.
---

# T001: Launcher Config Contract

## Problem

Savepoint has no configuration for enabling agent launches, assigning builder and auditor roles, or selecting terminal behavior.

## Context Files

- `internal/data/config.go`
- `internal/data/config_test.go`
- `.savepoint/config.yml`

## Acceptance Criteria

- [ ] `agent_launcher.enabled` defaults to `false` when absent.
- [ ] Builder and auditor profiles support an executable plus an ordered argument list with documented placeholders.
- [ ] Auditor configuration is optional; builder configuration is required only when the launcher is enabled.
- [ ] Terminal mode supports `auto` plus a structured override without requiring shell command strings.
- [ ] Existing config files parse with unchanged theme and quality-gate behavior.
- [ ] Validation errors name the exact missing or invalid launcher field.

## Implementation Plan

- [ ] Add typed launcher, agent-profile, and terminal configuration fields to `internal/data`.
- [ ] Define defaults and validation helpers separately from YAML parsing.
- [ ] Keep command arguments as structured values and define supported placeholders centrally.
- [ ] Add table-driven tests for absent, disabled, enabled, partial, and malformed configurations.
- [ ] Add a disabled example block to the repository-local config for development.

## Context Log

Pending.
