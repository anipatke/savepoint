---
id: E43-task-check-gates/T001-record-check-results
title: Record independent check results
status: done
objective: Load immutable Check records into the project index with validated scope, identity, and supersedes chains.
depends_on: []
complexity_tier: high
complexity_reason: Adds a third record family across decoding, confined discovery, and index validation.
---

# T001: Record independent check results

## Problem

A completed evaluation has nowhere to live. Without a Check record family in the index, clearance has no evidence to read and a rerun has no way to say which earlier result it replaces.

## Context Files

- `internal/data/check_v2.go`
- `internal/data/check_v2_test.go`
- `internal/data/errors.go`
- `internal/data/parser.go`
- `internal/data/objective_v2.go`
- `internal/data/task_v2.go`
- `internal/data/discover.go`
- `internal/data/discover_test.go`
- `internal/data/project.go`
- `internal/data/project_test.go`

## Acceptance Criteria

- [x] A Check record requires a valid global `C` ID, a `scope` of `{kind: task|objective, id: T###|O###}`, a `result` of `CLEAR` or `NEEDS WORK`, `checked_by` actor provenance, and a parseable `checked_at` timestamp.
- [x] `reviewed` decodes its base/head commit, files, and dependency entries as recorded basis; an absent or empty entry stays absent rather than becoming a passing default.
- [x] `issues` decodes as `I###` identity references only, with no Issue record loading, lookup, or validation.
- [x] `supersedes` resolves to an existing Check with the same scope; a missing target, a scope mismatch, a cycle, or two Checks superseding the same Check each return a distinct named diagnostic.
- [x] Checks are discovered only from `.savepoint/checks/`, through the same path confinement as V2 records, rejecting traversal, symlink escape, case aliasing, duplicate IDs, and a filename that does not start with its declared ID.
- [x] `V2Index` exposes Checks keyed by C ID and, for each Task or Objective target, its Checks in recorded order with the latest identifiable.
- [x] A Check whose scope names a record that does not exist in the index returns a named diagnostic.
- [x] A project with no `checks/` directory loads with no Checks and no error.
- [x] Malformed Check records fail the load closed, in deterministic ID order, exactly as E42's structural diagnostics do.

## Implementation Plan

- [x] Add `CheckV2`, its result/scope/actor types, and `DecodeCheckV2` in a new `check_v2.go`, reusing `ParseV2Document` and the existing anchored ID patterns.
- [x] Add the named Check sentinels to `errors.go`: malformed record, missing scope target, missing reference, supersedes conflict.
- [x] Add confined `.savepoint/checks/` discovery in `discover.go` reusing `v2PathConfiner` and the existing path/ID mismatch rule.
- [x] Extend `V2Index` and `LoadV2Index` in `project.go` with the Checks map, the per-target scope index, and scope-target resolution.
- [x] Validate supersedes chains during index construction: target exists, same scope, no cycle, no fork.
- [x] Test valid records, every malformed field, duplicate IDs, path mismatch, unsafe paths, missing scope targets, valid chains, forked and cyclic chains, absent `checks/`, and deterministic diagnostic ordering.
- [x] Run the focused `internal/data` suite.

## Context Log

**Files read:** `internal/data/check_v2.go` (n/a, new), `internal/data/errors.go`, `internal/data/parser.go`, `internal/data/objective_v2.go`, `internal/data/task_v2.go`, `internal/data/discover.go`, `internal/data/discover_test.go`, `internal/data/project.go`, `internal/data/project_test.go`, `internal/data/dependency.go`, `internal/data/lifecycle.go`, `internal/doctor/checks.go` (targeted read to confirm doctor wiring is out of scope for this task), `.savepoint/releases/v2/v2-Design.md` sections 4-5 (Check record schema fields, targeted read since the task's Context Files did not fully spell out `checked_by`/`reviewed` field names).

**Files edited:**
- `internal/data/check_v2.go` (new) — `CheckV2`, `CheckResult` (`CLEAR`/`NEEDS WORK`), `CheckScopeKind`/`CheckScope`, `ActorRole`/`Actor`, `ReviewedBasis`, and `DecodeCheckV2`, reusing `ParseV2Document`, `taskIDPatternV2`, and `objectiveIDPattern`.
- `internal/data/errors.go` — added `ErrV2CheckMalformed`, `ErrV2CheckMissingScopeTarget`, `ErrV2CheckMissingReference`, `ErrV2CheckSupersedesConflict`.
- `internal/data/discover.go` — added `v2ChecksDirName` and `DiscoverV2Checks(root)`, a flat confined `.savepoint/checks/` walk reusing `v2PathConfiner`.
- `internal/data/project.go` — extended `V2Index` with `Checks`, `ScopeChecks`, `LatestCheck`; `LoadV2Index` now calls `DiscoverV2Checks` and `indexChecks`; added `indexChecks`, `checkScopeTargetExists`, `validateCheckSupersedesChains` (reuses `findV2Cycle`/`describeV2Cycle` from `dependency.go`).
- `internal/data/check_v2_test.go` (new) — decode coverage: valid, minimal valid, reviewed-empty-stays-absent, malformed ID, scope (missing/invalid kind, missing id, kind/id mismatch), result (missing/invalid), checked_by (missing role/session, invalid role), checked_at (missing/unparseable), invalid issue reference, invalid supersedes reference, malformed YAML, no frontmatter.
- `internal/data/discover_test.go` — added `writeV2CheckFixture` and `TestDiscoverV2Checks_{valid,absentChecksDir,fileNameMismatch,duplicateID,rejectsSymlinkEscape,rejectsCaseCollision}`.
- `internal/data/project_test.go` — added `TestLoadV2Index_checks{Valid,AbsentDirLoadsEmpty,MissingScopeTarget,SupersedesChain,SupersedesMissingTarget,SupersedesScopeMismatch,SupersedesFork,SupersedesCycle,DeterministicDiagnosticOrder}`.

**Named cases and results:**
- `TestDecodeCheckV2_*` (13 top-level, several table-driven) — PASS
- `TestDiscoverV2Checks_*` (6) — PASS
- `TestLoadV2Index_check*` (9) — PASS
- Full `go test ./internal/data/...` — PASS (no regressions; existing E42 suites unchanged)

**Quality gates:** `make build && make test` — PASS across all packages (`.`, `cmd`, `internal/board`, `internal/buildtool`, `internal/data`, `internal/doctor`, `internal/init`, `internal/styles`). `go vet ./...` — clean.

No `.savepoint/Health-Check.md` in this project; Quick check step skipped per `savepoint-build-task`.

**Design decisions not fully spelled out in the task/epic:**
- `checked_by.role` and `scope.kind`/`result` invalid-enum cases share one sentinel, `ErrV2CheckMalformed`, per the Implementation Plan's four-sentinel budget; distinctness for `errors.Is` callers is at missing-scope-target / missing-reference / supersedes-conflict granularity, with message text distinguishing sub-cases within `ErrV2CheckSupersedesConflict` (scope mismatch, fork, cycle) — mirrors `ErrV2InvalidDependency`'s existing dual use in `task_v2.go`.
- `checked_at` is parsed strictly as RFC 3339 (matches the design doc's `'2026-09-14T00:00:00Z'` example).
- "Recorded order" for `ScopeChecks` and "latest identifiable" for `LatestCheck` are both defined as ascending `C###` ID order, since IDs are allocated monotonically and this avoids requiring a fully linear supersedes chain per scope to answer "latest."
- `checked_by.role` is restricted to the four canonical actor roles named in E43-Detail.md (`planner`, `executor`, `checker`, `owner`) rather than left as a free string.

No Drift Notes: all edited/added files (`check_v2.go`, `errors.go`, `discover.go`, `project.go`, and their tests) are the exact files E43-Detail.md's Components table already names for this schema.
