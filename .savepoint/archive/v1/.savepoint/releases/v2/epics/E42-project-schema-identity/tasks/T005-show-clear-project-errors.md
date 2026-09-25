---
id: E42-project-schema-identity/T005-show-clear-project-errors
title: Show clear project errors
status: done
objective: Surface V2 loader diagnostics through doctor and prove E42 behavior against frozen V1 fixtures and full regression gates.
depends_on:
    - E42-project-schema-identity/T003-keep-tasks-connected
    - E42-project-schema-identity/T004-keep-project-notes-intact
complexity_tier: medium
complexity_reason: Integrates new data diagnostics with doctor and validates cross-task behavior against frozen fixtures.
---

# T005: Show clear project errors

## Problem

The new schema and index contracts are not useful until doctor reports their diagnostics consistently and the combined behavior is proven without regressing frozen V1 inputs.

## Context Files

- `internal/data/project.go`
- `internal/data/project_test.go`
- `internal/data/errors.go`
- `internal/data/migration_source_test.go`
- `internal/data/migration_history_test.go`
- `internal/data/testdata/migration/v1-basic/manifest.yml`
- `internal/data/testdata/migration/v1-history/manifest.yml`
- `internal/doctor/checks.go`
- `internal/doctor/checks_test.go`
- `internal/doctor/interfaces.go`
- `internal/doctor/interfaces_test.go`
- `.savepoint/releases/v2/epics/E42-project-schema-identity/E42-Detail.md`

## Acceptance Criteria

- [x] Doctor obtains schema, record, identity, path, and graph findings from the shared data-layer project loader rather than duplicating those rules.
- [x] Doctor output gives each V2 problem a stable diagnostic name plus actionable record path and identity context.
- [x] Doctor reports a missing V2 Task title as an actionable schema error and does not accept the detailed objective as its replacement.
- [x] Doctor remains read-only for valid and invalid V1/V2 projects, including when configured quality gates are present.
- [x] Frozen E41 basic and history fixtures still load through transitional V1 dispatch with their source manifests unchanged.
- [x] One integrated V2 fixture covers distinct display titles and objectives, valid Objectives/Tasks, a moved Task, preserved unknown content, and deterministic indexing.
- [x] Focused data/doctor tests and `make build && make test` pass, with named outcome evidence recorded before audit handoff.

## Implementation Plan

- [x] Adapt doctor checks and its data interface to consume the schema-aware project result and named diagnostics.
- [x] Map data diagnostics to concise doctor reports without adding repair writes or a second validation implementation.
- [x] Add read-only assertions and representative malformed-project cases, including a Task with an objective but no title, to doctor tests.
- [x] Extend E41 fixture characterization only where needed to prove transitional dispatch and unchanged source manifests.
- [x] Add the integrated V2 project scenario and assert identity, ownership, reference, and preservation outcomes across the completed E42 boundary.
- [x] Run focused data and doctor packages, then `make build && make test`; record commands, named cases, results, and limitations.

## Context Log

**Files read:** `internal/data/project.go`, `internal/data/project_test.go`, `internal/data/errors.go`, `internal/data/task_v2.go`, `internal/data/objective_v2.go`, `internal/data/discover.go`, `internal/data/dependency.go`, `internal/data/config.go`, `internal/data/parser.go` (V2SourceDocument section), `internal/data/write.go` (`mappingFieldValue`, `patchV2Mapping` — reused, not duplicated), `internal/data/migration_source_test.go`, `internal/data/migration_history_test.go`, `internal/data/testdata/migration/v1-basic/manifest.yml`, `internal/data/testdata/migration/v1-history/manifest.yml`, `internal/doctor/checks.go`, `internal/doctor/checks_test.go`, `internal/doctor/interfaces.go`, `internal/doctor/interfaces_test.go`, `internal/doctor/report.go`, `internal/doctor/report_test.go`, `internal/doctor/repairs.go`, `AGENTS.md`, `.savepoint/Guardrails.md`, `.savepoint/router.md`, `E42-Detail.md`, `T001-recognize-v2-projects.md`, `T002-read-v2-work-safely.md`, `T003-keep-tasks-connected.md` (Context Log — explicitly deferred wiring `LoadProject`'s V2 dispatch to this task), `T004-keep-project-notes-intact.md`.

**Files edited:**
- `internal/data/project.go` — wired `loadProjectV2` to call `LoadV2Index(root)` instead of failing closed with a stub error; added `Project.V2 *V2Index`. This is the dispatch-wiring T003's Context Log explicitly deferred to T005.
- `internal/data/errors.go` — removed `ErrV2ProjectNotImplemented`: dead now that `loadProjectV2` returns real results/errors from `LoadV2Index` instead of the stub sentinel.
- `internal/data/project_test.go` — updated `TestLoadProject`'s "explicit version 2" case from expecting the stub error to expecting a real (empty) V2 index, and split the post-assertion between `project.V1`/`project.V2` by schema version. Added `TestLoadV2Index_integratedProjectScenario` (the epic's required integrated fixture): two Objectives with distinct human titles, three Tasks including one with unknown top-level and nested-mapping frontmatter fields (proving the raw `V2SourceDocument` survives the full discover path, not just direct-decode) and one moved Task filed under a different Objective's `tasks/` directory (proving ownership comes from its own field), plus a same-input re-run asserting `ObjectiveTasks` is byte-identical (deterministic indexing). Added `findMappingChild` test helper to walk into a nested mapping node (`mappingFieldValue` in write.go only reaches scalar values).
- `internal/data/migration_source_test.go` — added `TestMigrationSourceBasicLoadProjectDispatch`: loads the frozen v1-basic fixture through `LoadProject` (not `NewDiscover()` directly), asserts V1 dispatch and its `ListEpics` result, then re-asserts the fixture's manifest-recorded hashes are unchanged.
- `internal/data/migration_history_test.go` — added `TestMigrationHistoryLoadProjectDispatch`, the same proof for the multi-release v1-history fixture (`ListReleases` returns `[v1 v1.1]`), plus the unchanged-manifest assertion.
- `internal/doctor/checks.go` — added `CheckProject(root) []Problem`, the sole consumer of `data.LoadProject` for V2 diagnostics (no re-derived schema/record/graph rules), and `v2DiagnosticName(err) string`, an `errors.Is` switch giving every V2 structural sentinel in `internal/data/errors.go` (schema-version malformed/unsupported plus every `ErrV2*`) a stable name doctor reports it under.
- `internal/doctor/repairs.go` — added `V2ProblemRepair(name string) string`, mirroring the existing `AuditFindingRepair`/`AuditValidationRepair` typed-repair pattern, mapping each `v2DiagnosticName` value to a manual fix suggestion.
- `internal/doctor/report.go` — added `DiagnosticReport.Project []Problem`, wired `CheckProject` into `RunAllChecks`, `HasProblems`, and `Format` (new "Project Check" section, ordered right after Router Check).
- `internal/doctor/report_test.go` — added "Project Check" to `TestDiagnosticReport_FormatContainsSections`'s expected section list.
- `internal/doctor/checks_test.go` — added `writeV2Objective`/`writeV2Task` doctor-side fixture helpers (doctor cannot reach the data package's unexported V2 fixture helpers) and a `CheckProject` test suite: V1-project no-problems, valid-V2 no-problems, malformed/unsupported schema_version naming, the required missing-Task-title-not-backfilled-from-objective case, missing-owner, a Task dependency cycle, and `TestCheckProject_ReadOnly` (three subtests: valid V1 with quality gates configured via `RunAllChecks`, valid V2, invalid V2) each byte-comparing the inspected file before/after.

**Named cases and results:**
- `TestLoadProject` (7 subtests, including the updated V2-dispatch case) — PASS
- `TestLoadV2Index_integratedProjectScenario` — PASS
- `TestMigrationSourceBasicLoadProjectDispatch` — PASS
- `TestMigrationHistoryLoadProjectDispatch` — PASS
- `TestCheckProject_v1ProjectNoProblems`, `TestCheckProject_v2ValidNoProblems`, `TestCheckProject_SchemaVersionMalformed`, `TestCheckProject_SchemaVersionUnsupported`, `TestCheckProject_MissingTaskTitleNotBackfilledFromObjective`, `TestCheckProject_MissingOwner`, `TestCheckProject_DependencyCycle` — PASS
- `TestCheckProject_ReadOnly` (3 subtests) — PASS
- `TestDiagnosticReport_FormatContainsSections` (updated) — PASS
- Full `go test ./internal/data/...` — PASS (no regressions)
- Full `go test ./internal/doctor/...` — PASS (no regressions)
- `go vet ./...` — clean

**Quality gates:** `make build && make test` — PASS across all packages (`.`, `cmd`, `internal/board`, `internal/buildtool`, `internal/data`, `internal/doctor`, `internal/init`, `internal/styles`).

No `.savepoint/Health-Check.md` in this project; Quick check step skipped per `savepoint-build-task`.

No Drift Notes: all edits are to files E42-Detail.md's Components table already names (`project.go`, `errors.go`, `internal/doctor/checks.go`) or their existing `_test.go` counterparts, plus `internal/doctor/report.go`/`report_test.go`/`repairs.go`, which are the existing doctor report/repair-suggestion machinery every other doctor check already integrates with — not a new module or architecture change. `interfaces.go` was read but not edited: `data.LoadProject` is a direct, deterministic data-layer call (no filesystem enumeration doctor needs to inject/stub), so it doesn't fit the existing `taskDiscoverer`/`taskParser` injection seam, and forcing it through that seam would be an unscoped abstraction beyond what this task's ACs ask for.
