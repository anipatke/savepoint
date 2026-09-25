---
id: E27-parallel-planning-contracts/T001-advanced-parallel-config
status: planned
objective: Add disabled-by-default advanced parallel-worktree launcher configuration.
depends_on:
  - E22-agent-launcher-foundation/T001-launcher-config-contract
complexity_tier: medium
complexity_reason: Nested opt-in defaults and limits must preserve every existing launcher configuration path.
---

# T001: Advanced Parallel Config

## Problem

The launcher has no independent opt-in or capacity limit for parallel worktree execution.

## Context Files

- `internal/data/config.go`
- `internal/data/config_test.go`
- `.savepoint/config.yml`

## Acceptance Criteria

- [ ] `agent_launcher.advanced.parallel_worktrees.enabled` defaults to `false` when absent.
- [ ] `max_agents` has a conservative default and rejects values outside the documented range.
- [ ] Advanced mode cannot be enabled unless the base launcher and builder profile are valid.
- [ ] Disabling advanced mode does not disable or alter sequential launcher actions.
- [ ] Existing launcher and pre-v1.3 config fixtures parse with unchanged behavior.

## Implementation Plan

- [ ] Add typed advanced and parallel-worktree config structs under the existing launcher config.
- [ ] Define defaults and validation in the canonical config boundary.
- [ ] Return field-specific errors for invalid enablement and capacity values.
- [ ] Add absent, disabled, enabled, partial, and malformed table tests.
- [ ] Add a disabled repository-local example without enabling the feature.

## Context Log

Pending.
