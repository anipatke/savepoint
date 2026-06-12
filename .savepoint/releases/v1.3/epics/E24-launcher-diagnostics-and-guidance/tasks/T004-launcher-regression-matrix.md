---
id: E24-launcher-diagnostics-and-guidance/T004-launcher-regression-matrix
status: planned
objective: Verify the complete launcher action matrix and unchanged disabled behavior before release audit.
depends_on:
  - E23-board-launch-actions/T004-epic-audit-action-and-help
  - E24-launcher-diagnostics-and-guidance/T001-doctor-launcher-diagnostics
  - E24-launcher-diagnostics-and-guidance/T003-launcher-workflow-documentation
complexity_tier: high
complexity_reason: Release verification spans data, launcher, board, doctor, templates, and three operating systems.
---

# T004: Launcher Regression Matrix

## Problem

The feature crosses configuration, lifecycle writes, TUI dispatch, detached process startup, and platform behavior, so package tests alone cannot prove the release contract.

## Context Files

- `internal/board/integration_test.go`
- `internal/board/update_test.go`
- `internal/launcher/launcher_test.go`
- `internal/launcher/terminal_test.go`
- `internal/data/config_test.go`
- `internal/doctor/checks_test.go`
- `cmd/board_test.go`
- `Makefile`
- `README.md`

## Acceptance Criteria

- [ ] Integration tests cover the task, defect, and epic action matrix with fake launchers and no real agent process.
- [ ] Disabled and absent configuration produce byte-for-byte-equivalent action-free board behavior where practical.
- [ ] Tests cover lifecycle write failure, router write failure, missing role, startup failure, and duplicate dispatch.
- [ ] Cross-compilation covers Windows, macOS, and Linux launcher build-tag files.
- [ ] A documented manual matrix verifies one interactive terminal launch per supported OS without paid API calls.
- [ ] `make build && make test` passes before handoff to a fresh audit agent.

## Implementation Plan

- [ ] Add board-level integration fixtures for configured and disabled launcher states.
- [ ] Exercise launch messages with fake terminal and executable dependencies.
- [ ] Add regression assertions for existing navigation, lifecycle, and plain-output behavior.
- [ ] Run host tests plus cross-platform compilation for all launcher adapters.
- [ ] Record manual terminal verification commands and outcomes in the Context Log.
- [ ] Run the full repository quality gates and stop for user handoff.

## Context Log

Pending.
