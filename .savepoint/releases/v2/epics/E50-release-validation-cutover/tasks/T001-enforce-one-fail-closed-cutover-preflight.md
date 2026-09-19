---
id: E50-release-validation-cutover/T001-enforce-one-fail-closed-cutover-preflight
status: planned
objective: Expose one fail-closed cutover preflight over migration safety, V2 validity, and canonical Release readiness.
depends_on: []
complexity_tier: high
complexity_reason: Combines migration, loading, recovery, and Release gates without duplicating policy.
---

# T001: Enforce one fail-closed cutover preflight

## Problem

E51 proved the individual migration and Release decisions, but E50 still needs one operational answer that refuses cutover unless the complete candidate is safe. That composition must not become a second Release-readiness policy.

## Context Files

- `.savepoint/releases/v2/epics/E50-release-validation-cutover/E50-Detail.md`
- `.savepoint/releases/v2/epics/E51-first-class-releases/tasks/T009-prove-release-migration-and-gate-the-cutover.md`
- `.savepoint/releases/v2/v2-Design.md`
- `internal/data/project.go`
- `internal/data/release_cutover.go`
- `internal/data/release_gate_v2.go`
- `internal/data/release_gate_v2_test.go`
- `internal/migrate/plan.go`
- `internal/migrate/preview.go`
- `internal/migrate/apply.go`
- `internal/migrate/operation.go`
- `internal/migrate/end_to_end_test.go`
- `main.go`
- `main_test.go`

## Acceptance Criteria

- [ ] One typed result states whether cutover may proceed and names every blocking operational condition without rewriting Release completion rules.
- [ ] Ambiguous decisions, source conflicts, an incomplete operation, an unrecoverable plan, invalid staged V2 state, and unsupported schema all fail closed with actionable diagnostics.
- [ ] Every declared Release is evaluated through `data.ResolveReleaseCutover`; no E50 code rechecks Check freshness, material Issues, or owner acceptance itself.
- [ ] Release-free V2 projects remain valid when every operational prerequisite passes.
- [ ] Preflight and dry-run create no files, directories, probes, backups, manifests, or mtime changes.
- [ ] Table-driven tests cover every refusal, a clean Release-free candidate, and a fully accepted multi-Release candidate in temporary roots.

## Implementation Plan

- [ ] Define the cutover-preflight result and stable refusal categories at the migration/data boundary.
- [ ] Compose migration applicability, operation recovery state, candidate V2 loading, and `ResolveReleaseCutover` in execution order.
- [ ] Keep per-Release reasoning inside the canonical data resolver and translate only its typed result.
- [ ] Add temporary-project tests for every refusal and both successful project shapes.
- [ ] Prove the entire preflight path is read-only and preserves bytes and mtimes.
- [ ] Run focused data/migrate/root tests and the required full gates; record evidence.

## Context Log

Pending.
