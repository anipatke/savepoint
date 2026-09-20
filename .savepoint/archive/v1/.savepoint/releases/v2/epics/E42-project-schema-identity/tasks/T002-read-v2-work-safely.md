---
id: E42-project-schema-identity/T002-read-v2-work-safely
title: Read V2 work safely
status: done
objective: Decode strict V2 Objective and Task records without leaking unsafe V1 defaults into the new schema.
depends_on:
    - E42-project-schema-identity/T001-recognize-v2-projects
complexity_tier: high
complexity_reason: Establishes new record contracts across parsing, lifecycle validation, and V1 compatibility boundaries.
---

# T002: Read V2 work safely

## Problem

V2 needs explicit Objective and Task identity, ownership, and dependencies, while current V1 parsing heals values and derives meaning that must not carry into V2.

## Context Files

- `internal/data/parser.go`
- `internal/data/parser_test.go`
- `internal/data/task.go`
- `internal/data/task_test.go`
- `internal/data/lifecycle.go`
- `internal/data/lifecycle_test.go`
- `internal/data/objective_v2.go`
- `internal/data/objective_v2_test.go`
- `internal/data/task_v2.go`
- `internal/data/task_v2_test.go`
- `internal/data/project.go`
- `internal/data/project_test.go`

## Acceptance Criteria

- [x] V2 Objectives require a valid global `O` ID, title, supported lifecycle value, and typed Objective dependencies.
- [x] V2 Tasks require a valid global `T` ID, a non-empty human-facing `title`, exactly one `objective: O###` owner, supported lifecycle fields, and typed Task dependencies.
- [x] Task `title` is a distinct display label; the detailed objective or Outcome is not silently substituted when the title is missing.
- [x] Task dependency records accept `requires: clear|accepted`, default only an omitted requirement to `clear`, and reject every other value.
- [x] Optional release metadata is decoded as filtering metadata and never used to establish identity or ownership.
- [x] Missing titles or other required fields, malformed IDs, unknown lifecycle values, and invalid stage/status combinations return named V2 diagnostics without completion-capable healing.
- [x] Existing V1 Task parsing and compatibility aliases remain behaviorally unchanged and type-isolated from V2 records.

## Implementation Plan

- [x] Add V2 Objective, Task, dependency, source-metadata, and validation types without broadening the active V1 Task type.
- [x] Extend the frontmatter parser boundary so schema-specific decoders can retain their source document while projecting typed fields.
- [x] Implement strict Objective and Task decoders with required titles, anchored global-ID validation, and explicit ownership.
- [x] Apply V2 lifecycle validation separately from the V1 defaulting and alias path.
- [x] Test that missing V2 titles are rejected and detailed objectives are never used as an implicit display-title fallback.
- [x] Add positive, malformed, other missing-field, unsafe-default, and V1-regression tests; run the focused `internal/data` suite.

## Context Log

**Files read:** `internal/data/parser.go`, `internal/data/parser_test.go`, `internal/data/task.go`, `internal/data/task_test.go`, `internal/data/lifecycle.go`, `internal/data/lifecycle_test.go`, `internal/data/project.go`, `internal/data/project_test.go`, `internal/data/config.go`, `internal/data/errors.go`, `internal/data/write.go`, `internal/data/dependency.go` (V1 reference-matching shape, not reused), `AGENTS.md`, `.savepoint/Guardrails.md`, `.savepoint/router.md`, `E42-Detail.md`, `T001-recognize-v2-projects.md`, `.savepoint/releases/v2/v2-Design.md` (sections 1-6, identity/lifecycle/dependency contracts).

**Files edited:**
- `internal/data/errors.go` — added `ErrV2Malformed`, `ErrV2MissingField`, `ErrV2InvalidID`, `ErrV2InvalidOwnership`, `ErrV2InvalidLifecycle`, `ErrV2InvalidDependency` named sentinels.
- `internal/data/parser.go` — added `V2SourceDocument` (path + raw frontmatter `yaml.Node` + body) and `ParseV2Document(path, content)`, the raw-document/typed-record boundary strict V2 decoders use, built on the existing `SplitFrontmatterBody`.
- `internal/data/objective_v2.go` (new) — `ObjectiveV2` type and `DecodeObjectiveV2(path, content)`: anchored `O###` ID validation, required non-empty title, status validated via `IsCanonicalTaskStatus` (no alias healing), typed `O###` dependency list, inert `Release` field.
- `internal/data/objective_v2_test.go` (new) — valid/minimal, malformed-ID table, missing-title, unknown/legacy-alias-status-not-healed table, invalid dependency, malformed YAML, no-frontmatter, release-is-inert-metadata.
- `internal/data/task_v2.go` (new) — `TaskV2`, `TaskDependencyV2`, `TaskDependencyRequirement` (`clear`/`accepted`) types and `DecodeTaskV2(path, content)`: anchored `T###` ID validation, required title with no objective-as-title fallback, single `O###` owner validated via `objectiveIDPattern` (shared with objective_v2.go), lifecycle validated by reusing the existing strict `ValidateTaskLifecycleStateForWrite` (not the V1 healing `ParseTaskLifecycle` path), dependency records defaulting only an omitted `requires` to `clear` and rejecting all other values.
- `internal/data/task_v2_test.go` (new) — valid/minimal, malformed-ID table, title-not-backfilled-from-objective, missing/malformed-objective table, non-scalar-objective-rejected (proves single-owner structurally), lifecycle-not-healed table (missing status, legacy todo/complete aliases, missing/unknown/legacy-implementation stage, stale stage, `phase` field not treated as a stage alias), dependency requires-defaults-only-when-omitted, requires-rejects-unknown-values, malformed dependency task ID, malformed YAML, no-frontmatter, release-never-establishes-ownership, type-isolation from V1 `Task`, and a V1 regression test (`TestParseTaskFile_v1BehaviorUnaffectedByV2Types`) proving `ParseTaskFile`'s legacy `todo` healing and objective-as-title fallback are unchanged.

**Named cases and results:**
- `TestDecodeObjectiveV2_*` (9 tests, several table-driven subtests) — PASS
- `TestDecodeTaskV2_*` (14 tests, several table-driven subtests) — PASS
- `TestParseTaskFile_v1BehaviorUnaffectedByV2Types` — PASS
- Full `go test ./internal/data/...` — PASS (no regressions)
- `go vet ./...` — clean

**Quality gates:** `make build && make test` — PASS across all packages (`.`, `cmd`, `internal/board`, `internal/buildtool`, `internal/data`, `internal/doctor`, `internal/init`, `internal/styles`).

No `.savepoint/Health-Check.md` in this project; Quick check step skipped per `savepoint-build-task`.

No Drift Notes: `objective_v2.go`, `task_v2.go`, and their tests are exactly the files E42-Detail.md's Components table already names as new; `parser.go` and `errors.go` are extended, not repurposed. No unplanned files or architecture changes. `project.go` and `dependency.go` were read for context only and left unchanged — record discovery/index wiring is later E42 task scope, not this task's decode-only boundary.
