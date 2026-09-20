---
id: E50-release-validation-cutover/T003-confine-v1-compatibility-to-migration-and-history
status: done
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

- [x] Production references to V1 record parsing originate only inside explicit migration code; ordinary data, board, doctor, resume, and init paths cannot reach it.
- [x] Frozen V1 fixtures remain readable for migration regression tests without becoming accepted live-project input.
- [x] Every source file, Release PRD, authored reference, audit artifact, and unknown preserved byte still has a manifest-backed live or archive destination.
- [x] Preview, apply, every publish-boundary interruption, backup verification, conflict retry, and unchanged second apply retain their E51 guarantees.
- [x] Static or package-boundary tests fail if a live consumer reintroduces a V1 dependency.
- [x] Obsolete transitional adapters are removed rather than retained behind an undocumented compatibility flag.

## Implementation Plan

- [x] Identify the minimal legacy types and readers required by inventory and conversion.
- [x] Move or narrow those contracts to the migration boundary without changing preserved bytes or mappings.
- [x] Remove obsolete live adapters and add an enforceable dependency-boundary test.
- [x] Re-run fixture, archive, manifest, interruption, conflict, and no-op matrices.
- [x] Confirm unknown historical content remains inspectable but excluded from live diagnostics.
- [x] Run focused data/migrate tests and the required full gates; record evidence.

## Context Log

- `go test ./internal/data ./internal/doctor ./internal/board/v2 ./internal/migrate -run 'TestLiveSourcesDoNotReachV1Readers|TestFindProjectRoot|TestRunV2Checks|TestMigrationSourceBasicInterpretation|TestMigrationHistoryScopedReferences' -count=1`: PASS.
- `make build && make test`: PASS (`go test ./...`; `internal/migrate` 133.321s).
- `git diff --check`: PASS.
- Live root resolution now uses `migrate.FindProjectRoot`; V2 board load/write, resume, and the live doctor use `data.LoadV2Index`/`ReadStateV2` rather than schema-dispatched V1 readers.
- `doctor.RunV2Checks` is the live V2 health boundary; the documented `RunAllChecks` adapter remains only for frozen migration/history coverage. Malformed legacy audit content is ignored by the live V2 report, while legacy router fields are rejected by the strict reader.
- The obsolete `runV1Board`/plain-output dispatch adapter was removed. Remaining V1 board helpers are isolated in `internal/board/legacy.go` for retained historical tests; migration production readers remain under `internal/migrate`'s explicit plan/conversion path.
- `internal/migrate/live_boundary_test.go` scans all live command/consumer sources and fails on V1 reader calls, V1 record types, or schema-dispatch loading.
- Existing migration fixture, archive, manifest, interruption, backup, conflict-retry, and unchanged-second-apply matrices remain covered by `internal/data/migration_*_test.go` and `internal/migrate/*_test.go` and passed in the full gate.
- No `.savepoint/Health-Check.md` is present, so the Quick health check is skipped per `AGENTS.md`.
