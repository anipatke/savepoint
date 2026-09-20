---
id: E42-project-schema-identity/T001-recognize-v2-projects
title: Recognize V2 projects
status: done
objective: Load projects through one schema-aware boundary while preserving transitional V1 behavior and rejecting invalid explicit versions.
depends_on: []
complexity_tier: medium
complexity_reason: Introduces a shared loader boundary while preserving existing V1 discovery behavior.
---

# T001: Recognize V2 projects

## Problem

Project consumers currently enter V1-specific discovery directly, so they cannot distinguish a V2 project from legacy input or fail safely on unsupported schema declarations.

## Context Files

- `internal/data/config.go`
- `internal/data/config_test.go`
- `internal/data/discover.go`
- `internal/data/discover_test.go`
- `internal/data/project.go`
- `internal/data/project_test.go`

## Acceptance Criteria

- [x] A single data-layer project loader detects schema from `.savepoint/config.yml` before record discovery.
- [x] Absent `schema_version` selects transitional V1 behavior; numeric version `2` selects the V2 loader boundary.
- [x] Malformed and unsupported explicit versions fail with stable named diagnostics that identify the config path and supplied value.
- [x] Schema selection cannot be affected by package versions, release metadata, or `.upgrade-manifest.yml` content.
- [x] Existing V1 discovery behavior remains covered and unchanged behind the transitional dispatch path.

## Implementation Plan

- [x] Add the schema-version type, parser, and named errors to the config boundary without changing unrelated config defaults.
- [x] Introduce the project loader API and isolate V1 discovery behind its schema-specific adapter.
- [x] Add a V2 dispatch seam that later tasks can populate without silently falling back to V1 on V2 errors.
- [x] Cover absent, valid, malformed, unsupported, and unrelated-version inputs with table-driven temporary-project tests.
- [x] Run the focused `internal/data` tests and record the named cases and results.

## Context Log

**Files read:** `internal/data/config.go`, `internal/data/config_test.go`, `internal/data/discover.go`, `internal/data/discover_test.go`, `internal/data/errors.go`, `internal/data/task.go`, `internal/data/dependency.go`, `internal/data/migration_source_test.go`, `internal/testutil/fixture.go`, `internal/testutil/fs.go`, `.savepoint/releases/v2/v2-Design.md` (sections 1-3), `.savepoint/releases/v2/epics/E42-project-schema-identity/tasks/T002-read-v2-work-safely.md`, `T003-keep-tasks-connected.md` (for downstream seam shape).

**Files edited:**
- `internal/data/errors.go` — added `ErrMalformedSchemaVersion`, `ErrUnsupportedSchemaVersion`, `ErrV2ProjectNotImplemented` named sentinels.
- `internal/data/config.go` — added `SchemaVersion` type/constants and `ReadSchemaVersion(path string)`, decoding only the `schema_version` scalar via `yaml.Node` so theme/quality-gate/launcher defaults are untouched.
- `internal/data/config_test.go` — added `TestReadSchemaVersion` table covering absent file, absent field, valid `2`, malformed non-integer, unsupported `3`, unsupported `1`, and unrelated version-shaped fields.
- `internal/data/project.go` (new) — added `Project`, `LoadProject(root)` schema dispatch, `loadProjectV1` (wraps `Discover` unchanged) and `loadProjectV2` (fails closed via `ErrV2ProjectNotImplemented`, no V1 fallback).
- `internal/data/project_test.go` (new) — `TestLoadProject` table (absent config, absent field, valid V2 seam, malformed, unsupported, unrelated-version isolation via a competing `.upgrade-manifest.yml`), `TestLoadProjectV1DiscoveryUnchanged` (diffs `project.V1.ListEpics` against `NewDiscover().ListEpics` directly), `TestLoadProjectRejectsDirectoryReadFailure` (non-schema read error is not misreported as a schema diagnostic).

**Named cases and results:**
- `TestReadSchemaVersion` (7 subtests) — PASS
- `TestLoadProject` (6 subtests) — PASS
- `TestLoadProjectV1DiscoveryUnchanged` — PASS
- `TestLoadProjectRejectsDirectoryReadFailure` — PASS
- Full `go test ./internal/data/...` — PASS (no regressions)

**Quality gates:** `make build && make test` — PASS across all packages (`.`, `cmd`, `internal/board`, `internal/buildtool`, `internal/data`, `internal/doctor`, `internal/init`, `internal/styles`).

No `.savepoint/Health-Check.md` in this project; Quick check step skipped per `savepoint-build-task`.

No Drift Notes: `project.go` is the file already named in E42-Detail.md's Components table; no unplanned files or architecture changes.
