---
id: E50-release-validation-cutover/T003-confine-v1-compatibility-to-migration-and-history
status: planned
objective: Keep legacy parsing only for explicit migration and frozen history while preserving every archive and reference contract.
depends_on:
    - E50-release-validation-cutover/T002-make-live-command-routing-v2-only
complexity_tier: high
complexity_reason: Retires broad legacy reachability while preserving exact migration and archive behavior.
---

# T003: Confine V1 compatibility to migration and history

## Problem

Making commands V2-only is insufficient if legacy parsers remain general runtime dependencies. The retained V1 surface must be deliberately isolated behind migration inventory/conversion and immutable test fixtures.

## Context Files

- `.savepoint/releases/v2/epics/E50-release-validation-cutover/E50-Detail.md`
- `.savepoint/releases/v2/v2-Design.md`
- `internal/data/discover.go`
- `internal/data/parser.go`
- `internal/data/router.go`
- `internal/data/task.go`
- `internal/data/migration_source_test.go`
- `internal/data/migration_history_test.go`
- `internal/migrate/inventory.go`
- `internal/migrate/classify.go`
- `internal/migrate/convert.go`
- `internal/migrate/convert_docs.go`
- `internal/migrate/convert_issues.go`
- `internal/migrate/convert_releases.go`
- `internal/migrate/manifest.go`
- `internal/migrate/fixture_test.go`
- `internal/migrate/end_to_end_test.go`
- `internal/data/testdata/migration/README.md`

## Acceptance Criteria

- [ ] Production references to V1 record parsing originate only inside explicit migration code; ordinary data, board, doctor, resume, and init paths cannot reach it.
- [ ] Frozen V1 fixtures remain readable for migration regression tests without becoming accepted live-project input.
- [ ] Every source file, Release PRD, authored reference, audit artifact, and unknown preserved byte still has a manifest-backed live or archive destination.
- [ ] Preview, apply, every publish-boundary interruption, backup verification, conflict retry, and unchanged second apply retain their E51 guarantees.
- [ ] Static or package-boundary tests fail if a live consumer reintroduces a V1 dependency.
- [ ] Obsolete transitional adapters are removed rather than retained behind an undocumented compatibility flag.

## Implementation Plan

- [ ] Identify the minimal legacy types and readers required by inventory and conversion.
- [ ] Move or narrow those contracts to the migration boundary without changing preserved bytes or mappings.
- [ ] Remove obsolete live adapters and add an enforceable dependency-boundary test.
- [ ] Re-run fixture, archive, manifest, interruption, conflict, and no-op matrices.
- [ ] Confirm unknown historical content remains inspectable but excluded from live diagnostics.
- [ ] Run focused data/migrate tests and the required full gates; record evidence.

## Context Log

Pending.
