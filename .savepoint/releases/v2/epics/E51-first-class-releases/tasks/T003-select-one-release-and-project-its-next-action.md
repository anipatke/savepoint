---
id: E51-first-class-releases/T003-select-one-release-and-project-its-next-action
status: planned
objective: Persist optional Release context and make the shared Next projection interpret it without hiding unrelated global work.
depends_on:
  - E51-first-class-releases/T002-require-release-integration-evidence-and-owner-acceptance
complexity_tier: high
complexity_reason: Changes router persistence and the shared precedence projection consumed by every V2 user surface.
---

# T003: Select one Release and project its next action

## Problem

Board and resume cannot agree on Release context if selection exists only inside the TUI. The router and `data.Next` need one optional Release selection with strict mismatch handling, scoped progression inside a selected Release, and unchanged global behavior when no Release is selected.

## Context Files

- `.savepoint/releases/v2/epics/E51-first-class-releases/E51-Detail.md`
- `internal/data/release_v2.go`
- `internal/data/release_gate_v2.go`
- `internal/data/project.go`
- `internal/data/router_v2.go`
- `internal/data/router_v2_test.go`
- `internal/data/write.go`
- `internal/data/write_test.go`
- `internal/data/next.go`
- `internal/data/next_test.go`

## Acceptance Criteria

- [ ] V2 router state accepts optional `release: R###` and preserves it with unknown fields and the non-YAML body through managed selection writes.
- [ ] Selecting a Release writes Release context without changing router phase, Task/Objective lifecycle, evidence, or `next_action` prose.
- [ ] A selected Objective must belong to the selected Release; missing, archived, unassigned, and mismatched selections produce distinct named diagnostics.
- [ ] No missing or mismatched selection silently substitutes another Release or a same-numbered/same-titled Objective.
- [ ] With a valid Release selected, `ResolveNext` considers eligible Tasks/Objectives in that Release, then Release Check, Release owner acceptance, and Release-ready outcomes.
- [ ] With no Release selected, the pre-E51 global precedence and selected Objective/Task interpretation remain unchanged.
- [ ] Unassigned Objectives and Objectives in other Releases remain independently actionable globally and are not treated as dependencies of the selected Release.
- [ ] Pending migration and replan continue to outrank Release-scoped execution; existing Task and Objective dependency semantics retain their relative precedence.
- [ ] Repeat selection writes are byte- and mtime-idempotent, and stale mtime or pending migration refuses without partial changes.
- [ ] The projection exposes enough typed Release identity, evidence, blocker, and action data for board and resume to render without consulting the index again.

## Implementation Plan

- [ ] Extend typed V2 router parsing with optional Release selection and cross-selection validation.
- [ ] Extend the canonical router selection writer to preserve or change Release context safely.
- [ ] Add Release identity, readiness, evidence, and selection-diagnostic fields/rungs to `data.Next`.
- [ ] Integrate Release selection into the precedence ladder without changing the no-selection path.
- [ ] Make Release Check and owner-acceptance actions occur only after member Objective work is complete.
- [ ] Add table-driven coverage for valid, missing, archived, unassigned, and mismatched selection combinations.
- [ ] Freeze the pre-E51 no-Release matrix as regression evidence and verify write idempotence/conflict behavior.

## Context Log

Pending.
