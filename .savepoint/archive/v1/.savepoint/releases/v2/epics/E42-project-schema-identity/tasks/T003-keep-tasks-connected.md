---
id: E42-project-schema-identity/T003-keep-tasks-connected
title: Keep tasks connected when files move
status: done
objective: Discover V2 records safely and expose one identity-keyed project index with validated ownership and reference graphs.
depends_on:
    - E42-project-schema-identity/T002-read-v2-work-safely
complexity_tier: high
complexity_reason: Coordinates confined discovery, global indexing, and cycle-safe graph validation across project records.
---

# T003: Keep tasks connected when files move

## Problem

Later V2 consumers need one authoritative view of Objectives, Tasks, ownership, and references; path-derived lookup would make moved records unstable and duplicate validation inconsistent.

## Context Files

- `internal/data/project.go`
- `internal/data/project_test.go`
- `internal/data/objective_v2.go`
- `internal/data/objective_v2_test.go`
- `internal/data/task_v2.go`
- `internal/data/task_v2_test.go`
- `internal/data/discover.go`
- `internal/data/discover_test.go`
- `internal/data/dependency.go`
- `internal/data/dependency_test.go`
- `internal/data/errors.go`

## Acceptance Criteria

- [x] V2 discovery loads Objective and Task records into deterministic indexes keyed by their explicit global IDs.
- [x] Moving a Task without changing its record ID retains its identity and explicit Objective ownership while a path/record naming mismatch is diagnosed.
- [x] Duplicate IDs, missing Objective owners, missing dependency targets, self-dependencies, Task cycles, and Objective cycles produce stable diagnostics with paths and IDs.
- [x] Stored record paths are normalized project-relative paths, and discovery rejects traversal, symlink escape, and filesystem-aliasing case collisions.
- [x] Index construction validates structure only; it does not infer Check clearance, owner acceptance, or dependency satisfaction.
- [x] V1 discovery remains isolated and its release/epic-scoped behavior is unchanged.

## Implementation Plan

- [x] Implement confined V2 Objective and Task discovery using the record decoders from T002.
- [x] Build deterministic identity maps and explicit Objective-to-Task membership from each Task's owner reference.
- [x] Add path/name consistency, duplicate identity, missing target, and path-confinement diagnostics.
- [x] Add separate Task and Objective graph validation with self-edge and cycle-path reporting.
- [x] Cover moved records, ordering, invalid graphs, traversal, symlink, case-collision, and V1-regression scenarios in temporary-project tests.
- [x] Run the focused `internal/data` tests and record the diagnostic names and outcomes.

## Context Log

**Files read:** `internal/data/project.go`, `internal/data/project_test.go`, `internal/data/objective_v2.go`, `internal/data/objective_v2_test.go`, `internal/data/task_v2.go`, `internal/data/task_v2_test.go`, `internal/data/discover.go`, `internal/data/discover_test.go`, `internal/data/dependency.go`, `internal/data/errors.go`, `internal/data/parser.go`, `internal/data/config.go`, `internal/testutil/fs.go`, `internal/testutil/fixture.go`, `AGENTS.md`, `.savepoint/router.md`, `E42-Detail.md`, `T001-recognize-v2-projects.md`, `T002-read-v2-work-safely.md`, `T004-keep-project-notes-intact.md`, `T005-show-clear-project-errors.md`, `.savepoint/releases/v2/v2-Design.md` (sections 2-3, file model and identity/reference contracts).

**Files edited:**
- `internal/data/errors.go` — added `ErrV2DuplicateID`, `ErrV2PathMismatch`, `ErrV2UnsafePath`, `ErrV2MissingOwner`, `ErrV2MissingDependencyTarget`, `ErrV2SelfDependency`, `ErrV2DependencyCycle` named sentinels.
- `internal/data/discover.go` — added `DiscoverV2Records(root)`, confined to `root/objectives/O###-slug/Objective.md` and `.../tasks/T###-slug.md` (the V2 file model layout). Decodes each record via T002's `DecodeObjectiveV2`/`DecodeTaskV2`, then checks directory/file-name-vs-ID consistency and duplicate IDs. Added `v2PathConfiner`, whose `stat` method is the only way discovery inspects a path: it resolves symlinks and rejects any path resolving outside the project root, and rejects a project-relative path that only differs by case from one already visited (case-insensitive-filesystem aliasing). A missing `objectives/` directory is a valid empty V2 project, not an error; a missing `Objective.md` or `tasks/` directory is likewise not fatal there — an orphaned Task still surfaces through the missing-owner diagnostic once ownership is validated.
- `internal/data/project.go` — added `V2Index` (identity-keyed `Objectives`/`Tasks` maps plus `ObjectiveTasks`, the explicit ownership membership built from each Task's own `objective` field, never from its containing directory) and `LoadV2Index(root)`, which calls `DiscoverV2Records`, builds ownership in sorted-ID order (missing-owner diagnostic here), and calls `ValidateV2ReferenceGraphs`. Left `loadProjectV2`/`ErrV2ProjectNotImplemented` untouched — wiring `LoadProject`'s V2 dispatch to this index is T005's scope (doctor integration), and changing it here would break T001's already-approved `TestLoadProject` "dispatches to the V2 seam" case.
- `internal/data/dependency.go` — added `ValidateV2ReferenceGraphs(index)`, which validates the Task graph and the Objective graph independently (self-dependency, missing target, then cycle), each walked in sorted-ID order for repeatable diagnostics. Cycle detection is a small iterative-DFS `findV2Cycle` shared by both graphs, colored white/gray/black, returning the first cycle as an ID sequence; `describeV2Cycle` renders it with each record's source path.
- `internal/data/discover_test.go` — added `TestDiscoverV2Records_*`: valid two-objective/three-task tree, empty project (no `objectives/` dir), moved-task ownership independence, objective-directory and task-file naming mismatches, duplicate objective/task IDs, symlink escape (skipped on windows), and case-collision aliasing. Added shared `writeV2ObjectiveFixture`/`writeV2TaskFixture` helpers reused by `project_test.go`.
- `internal/data/project_test.go` — added `TestLoadV2Index_*`: valid index assembly with sorted `ObjectiveTasks`, empty project, moved-task ownership through the full assembled index, missing-owner diagnostic, and a disk-round-trip Task cycle (`depends_on` written as raw YAML, proving the discover-then-validate path, not just the in-memory graph algorithm).
- `internal/data/dependency_test.go` (new) — `TestValidateV2ReferenceGraphs_*` built directly against hand-constructed `V2Index` values: valid acyclic graph, Task self-dependency, Task missing target, three-node Task cycle, Objective self-dependency, Objective missing target, two-node Objective cycle, and a case proving Task and Objective cycles are detected independently rather than sharing one combined graph.

**Named cases and results:**
- `TestDiscoverV2Records_*` (9 tests) — PASS
- `TestLoadV2Index_*` (5 tests) — PASS
- `TestValidateV2ReferenceGraphs_*` (8 tests) — PASS
- Existing V1 discovery tests (`TestFindSavepointRoot`, `TestListReleases`, `TestListRootDirs*`, `TestListEpics`, `TestListTasks`) and existing `TestLoadProject*` — PASS, unchanged
- Full `go test ./internal/data/...` — PASS (no regressions)
- `go vet ./...` — clean

**Quality gates:** `make build && make test` — PASS across all packages (`.`, `cmd`, `internal/board`, `internal/buildtool`, `internal/data`, `internal/doctor`, `internal/init`, `internal/styles`).

No `.savepoint/Health-Check.md` in this project; Quick check step skipped per `savepoint-build-task`.

No Drift Notes: `discover.go`, `project.go`, `dependency.go`, and `errors.go` are exactly the files E42-Detail.md's Components table names for this responsibility split (confined V2 discovery, the identity-keyed project/index API, and reference-graph validation, respectively). No unplanned files or architecture changes. `loadProjectV2`'s dispatch stub was deliberately left alone; see the `project.go` edit note above.
