---
id: E50-release-validation-cutover/T001-enforce-one-fail-closed-cutover-preflight
status: done
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
- `internal/migrate/cutover.go`
- `internal/migrate/cutover_test.go`
- `main.go`
- `main_test.go`

## Acceptance Criteria

- [x] One typed result states whether cutover may proceed and names every blocking operational condition without rewriting Release completion rules.
- [x] Ambiguous decisions, source conflicts, an incomplete operation, an unrecoverable plan, invalid staged V2 state, and unsupported schema all fail closed with actionable diagnostics.
- [x] Every declared Release is evaluated through `data.ResolveReleaseCutover`; no E50 code rechecks Check freshness, material Issues, or owner acceptance itself.
- [x] Release-free V2 projects remain valid when every operational prerequisite passes.
- [x] Preflight and dry-run create no files, directories, probes, backups, manifests, or mtime changes.
- [x] Table-driven tests cover every refusal, a clean Release-free candidate, and a fully accepted multi-Release candidate in temporary roots.

## Implementation Plan

- [x] Define the cutover-preflight result and stable refusal categories at the migration/data boundary.
- [x] Compose migration applicability, operation recovery state, candidate V2 loading, and `ResolveReleaseCutover` in execution order.
- [x] Keep per-Release reasoning inside the canonical data resolver and translate only its typed result.
- [x] Add temporary-project tests for every refusal and both successful project shapes.
- [x] Prove the entire preflight path is read-only and preserves bytes and mtimes.
- [x] Run focused data/migrate/root tests and the required full gates; record evidence.

## Context Log

- `go test ./internal/data ./internal/migrate . -run 'TestPreflightCutover|TestResolveReleaseCutover' -count=1`: PASS.
- `make build`: PASS.
- `make test`: PASS (`go test ./...`; `internal/migrate` 136.768s).
- Final loader adjustment verification: `make build`: PASS; `make test`: PASS (`go test ./...`; `internal/migrate` 134.441s).
- `git diff --check`: PASS.
- Preflight tests use temporary project roots and prove file bytes, directory entries, and mtimes are unchanged across V1 planning, valid V2 loading, and pending-operation inspection.
- No `.savepoint/Health-Check.md` is present, so the Quick health check is skipped per `AGENTS.md`.
