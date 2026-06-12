---
id: E24-launcher-diagnostics-and-guidance/T001-doctor-launcher-diagnostics
status: planned
objective: Diagnose invalid enabled launcher configuration without affecting disabled projects.
depends_on:
  - E22-agent-launcher-foundation/T004-launcher-service
complexity_tier: medium
complexity_reason: Diagnostics must distinguish optional absence from actionable enabled-state configuration errors.
---

# T001: Doctor Launcher Diagnostics

## Problem

Users who opt in need clear diagnostics for missing commands, invalid placeholders, and unsupported terminal settings before attempting a board launch.

## Context Files

- `internal/doctor/checks.go`
- `internal/doctor/checks_test.go`
- `internal/doctor/report.go`
- `internal/doctor/report_test.go`
- `internal/doctor/repairs.go`
- `internal/doctor/repairs_test.go`
- `internal/doctor/interfaces.go`
- `internal/data/config.go`
- `internal/launcher/launcher.go`

## Acceptance Criteria

- [ ] Absent or disabled launcher configuration produces no doctor problem.
- [ ] Enabled configuration reports missing builder executable, invalid argument placeholders, and invalid terminal settings.
- [ ] Missing optional auditor configuration is not reported as an error.
- [ ] Diagnostics distinguish malformed configuration from an executable not found on the current host.
- [ ] Repair suggestions are descriptive and never install software or edit config automatically.
- [ ] Doctor remains deterministic and read-only.

## Implementation Plan

- [ ] Add a focused launcher configuration check using shared data validation.
- [ ] Add injectable executable lookup where host availability is diagnosed.
- [ ] Map typed validation failures to concise report messages and repair suggestions.
- [ ] Add tests for absent, disabled, valid, malformed, missing-command, and optional-auditor cases.
- [ ] Keep existing required config checks and exit semantics unchanged.

## Context Log

Pending.
