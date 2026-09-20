---
id: E24-launcher-diagnostics-and-guidance/T002-scaffold-launcher-config
status: planned
objective: Scaffold a clear disabled launcher configuration without changing existing project upgrades.
depends_on:
  - E22-agent-launcher-foundation/T001-launcher-config-contract
complexity_tier: low
complexity_reason: The change is limited to scaffold data and focused init assertions.
---

# T002: Scaffold Launcher Config

## Problem

New projects need a discoverable opt-in configuration example, while existing projects must remain untouched and disabled by default.

## Context Files

- `templates/project/.savepoint/config.yml`
- `internal/init/scaffold_test.go`
- `internal/init/write_test.go`
- `cmd/init_test.go`
- `cmd/upgrade-assets_test.go`

## Acceptance Criteria

- [ ] Newly initialized projects contain an `agent_launcher` block with `enabled: false`.
- [ ] The example shows builder, optional auditor, and terminal settings without assuming a vendor.
- [ ] Existing projects without the block continue to load with the launcher disabled.
- [ ] `upgrade-assets` does not overwrite project-owned `config.yml`.
- [ ] Init and upgrade regression tests cover the ownership boundary.

## Implementation Plan

- [ ] Add a concise commented launcher block to the scaffolded config template.
- [ ] Use placeholder example values that cannot accidentally launch a process while disabled.
- [ ] Extend init tests to assert disabled defaults and expected YAML structure.
- [ ] Confirm upgrade-assets exclusion tests continue to protect config ownership.

## Context Log

Pending.
