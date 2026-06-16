---
type: epic-design
status: planned
---

# E29: Parallel Mode Guidance and Hardening

## Purpose

Make advanced parallel worktrees diagnosable, safely scaffolded, documented, and proven backward compatible before release.

## What this epic adds

- Doctor diagnostics for configuration, planning metadata, manifests, and live Git mismatches.
- Disabled scaffold examples and ignore rules for local runtime/worktree state.
- User and agent guidance for planning, launching, integrating, toggling, and cleaning up parallel work.
- A release regression matrix covering disabled, enabled, restart, partial-failure, and cross-platform behavior.

## Components and files

| Module | Purpose |
|--------|---------|
| `internal/doctor` | Diagnose advanced config, metadata, and runtime/Git drift |
| `templates/project` | Scaffold disabled configuration, ignore rules, and agent workflow guidance |
| `README.md` | Explain advanced-mode operation and safety boundaries |
| `internal/board` and `internal/worktree` tests | Prove compatibility and failure behavior |

## Architectural delta

Doctor gains read-only awareness of local runtime coordination while remaining non-destructive. Generated projects advertise the capability without enabling it, and the release contract explicitly separates isolation/launch from integration and judgment.

## Boundaries

**In scope:**
- Read-only diagnostics and repair suggestions
- Disabled scaffold and upgrade-safe managed guidance
- User workflow documentation
- Automated and manual release verification

**Out of scope:**
- Doctor cleanup or repair execution
- Automatic branch integration
- Paid-agent end-to-end automation
- Additional agent vendors or terminal emulators

## Quality gates

- Doctor tests cover metadata and runtime drift without mutating Git state.
- Init and upgrade tests prove advanced mode remains disabled.
- Disabled-mode board and config behavior remains unchanged.
- Cross-platform builds pass and manual worktree operations are documented.
- `make build && make test` passes.

## Open decisions

None.
