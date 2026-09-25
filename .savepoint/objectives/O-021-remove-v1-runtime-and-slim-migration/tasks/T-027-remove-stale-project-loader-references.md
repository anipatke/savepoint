---
id: T-027
title: Remove stale project loader references
objective: O-021
planned_by: {role: planner, session: o021-remediation-20260924}
status: done
complexity_tier: low
complexity_reason: "Remove an unused schema-dispatch wrapper, preserve the active schema/index APIs and migration readers, and update their existing test callers and one Design paragraph."
depends_on: [{task: T-022, requires: clear}, {task: T-026, requires: clear}]
owner_validation: {required: false, accepted_check: ""}
check_waiver:
  task: T-027
  reason: Owner explicitly waived the optional Task Check in chat and directed the evidence to the mandatory Full Objective Check.
  actor: {role: owner, session: owner-chat-20260924}
  recorded_at: '2026-09-23T23:01:49Z'
---

# T-027: Remove stale project loader references

## Outcome

O-021's data dead-code scan is clean and Design describes the live board load path. The runtime schema gate and V2 index remain the only current project-loading path, while `internal/migrate` retains the V1 readers it directly uses.

## User Check

No interactive behavior changes. The independent Full Objective Check verifies the data dead-code result, existing behavior, and Design reconciliation.

## Done When

- `data.LoadProject`, `loadProjectV1`, `loadProjectV2`, and the now-unused `Project` wrapper are removed. `CheckRuntimeSchema`, `LoadV2Index`, and migration's directly used V1 readers remain.
- Test callers no longer depend on the retired wrapper; existing assertions use the live schema/index APIs where applicable, and V1 dispatch-only assertions are removed without weakening tests of migration conversion.
- A fresh `deadcode .` report contains no `internal/data` entries.
- `.savepoint/Design.md` describes the board's actual schema check, V2 index load, router decoding, and `ResolveNext` flow; no active guidance describes `LoadProject` as a runtime path.
- Updated Go production/test line counts are recorded in this Task's evidence.
- `git diff --check` and a fresh `make test-full` pass.

## Context Files

`.savepoint/Design.md`, `.savepoint/checks/C-914-o021-objective-check.md`, `.savepoint/issues/I-040-data-loadproject-unreachable-v1-dispatch.md`, `.savepoint/issues/I-041-design-board-refresh-names-loadproject.md`, `internal/data/project.go`, `internal/data/runtime.go`, `internal/data/project_test.go`, `internal/data/migration_history_test.go`, `internal/data/migration_source_test.go`, `main_test.go`, `main_resume_matrix_test.go`, `internal/migrate/apply_test.go`, `internal/init/v2_scaffold_test.go`.

## Design References

Design sections 1, 2, and 10.

## Guardrails

ARCH-04, DATA-01, TPL-02, TEST-06, TEST-08.

## Implementation Plan

1. Confirm the recorded I-040/I-041 evidence and the active schema/index interfaces in the listed Context Files.
2. Remove the unused `Project` type and `LoadProject` V1/V2 dispatch helpers; keep APIs and V1 readers used by current commands or `internal/migrate`.
3. Update or remove only tests whose assertions require the retired wrapper, keeping live schema-gate, V2-index, and migration conversion coverage.
4. Correct the Board persistence and refresh paragraph in Design to name the real load path.
5. Record the fresh `deadcode .` result, updated line counts, `git diff --check`, and `make test-full` result.

## Boundaries

No runtime behavior changes, no new V1 runtime support, no migration-reader changes, and no edits to immutable Check evidence. Keep I-040 and I-041 open for the independent checker to resolve after a CLEAR recheck.

## Technical Verification

`deadcode .`, `git diff --check`, and a fresh `make test-full`. This is required to replace the full-gate evidence invalidated by the code change and to satisfy I-040's proof.

## Technical Evidence

Executor session `o021-remediation-20260924`; toolchain go1.26.2 linux/amd64.

### Per-criterion outcome

1. **Met.** Removed `Project`, `LoadProject`, `loadProjectV1`, and `loadProjectV2` from `internal/data/project.go`. `CheckRuntimeSchema`, `LoadV2Index`, and the direct V1 readers used by `internal/migrate` remain.
2. **Met.** Replaced wrapper callers in the current schema/index tests and migration apply tests. V1 fixture reader cases now exercise `NewDiscover` directly. Runtime loader cases in `project_test.go` now exercise `CheckRuntimeSchema`; existing `LoadV2Index` tests remain.
3. **Met.** Fresh `/home/user/go/bin/deadcode .` reports no `internal/data` entries. The only remaining entries are `internal/migrate.FindProjectRoot`, `UnmarshalManifest`, and `WriteManifestCreateOnly`, outside O-021's dead-code packages.
4. **Met.** Design's board refresh paragraph now names `CheckRuntimeSchema`, `LoadV2Index`, router decoding, and `ResolveNext`. Two doctor comments naming the removed API were corrected. `LoadProject` remains only in boundary tests that forbid that API and in unrelated `internal/board/v2` loader test names.
5. **Met.** Updated per-package newline counts are recorded below against T-026's post-edit counts.
6. **Met.** `git diff --check` and a fresh `make build && make test-full` passed.

### Test evidence

The full gate ran all packages and the linux, darwin, and windows builds. Changed test cases exercised by that gate:

- `internal/data/project_test.go`: `TestCheckRuntimeSchema` (absent config; absent schema version; schema version 2; malformed and unsupported versions; unrelated version fields), `TestCheckRuntimeSchemaRejectsDirectoryReadFailure`, and the existing `TestLoadV2Index_*` cases.
- `internal/data/migration_history_test.go`: `TestMigrationHistoryDiscoverReleases`.
- `internal/data/migration_source_test.go`: `TestMigrationSourceBasicDiscoverEpics`.
- `main_test.go`: `TestMainInitScaffoldsV2ProjectWithEmptyValidIndex`.
- `main_resume_matrix_test.go`: `TestResumeMatrix_everyRungReachedExactlyOnce`, `TestResumeMatrix_runResumeThroughTheRealCommandAlsoWritesNothingTwice`.
- `internal/migrate/apply_test.go`: `TestApply_v1BasicEndToEnd`, `TestApply_v1HistoryKeepsReleaseScopedIDs`.
- `internal/init/v2_scaffold_test.go`: `TestV2ScaffoldLoadsCleanThroughRuntimeSchemaAndIndex`.

### Line counts

Counts are newline counts for non-recursive Go files in each package directory, matching T-026's method. Unchanged package rows remain as recorded in T-026.

| Package | After T-026 (prod / test) | After T-027 (prod / test) | Change |
| --- | ---: | ---: | ---: |
| `.` | 293 / 2,004 | 293 / 1,998 | 0 / -6 |
| `internal/data` | 7,272 / 12,416 | 7,224 / 12,341 | -48 / -75 |
| `internal/doctor` | 1,334 / 1,722 | 1,332 / 1,722 | -2 / 0 |
| `internal/init` | 1,541 / 7,672 | 1,541 / 7,667 | 0 / -5 |
| `internal/migrate` | 4,248 / 5,936 | 4,248 / 5,939 | 0 / +3 |
| **Total** | **23,204 / 37,829** | **23,154 / 37,746** | **-50 / -83** |

### Commands and results

- `gofmt -w` on the changed Go files: exit 0.
- `git diff --check`: exit 0.
- `make build && make test-full`: the sandbox attempt could not read `/home/user/.cache/go-build`; after approval to use that existing cache, the command exited 0. Full tests and linux/darwin/windows builds passed.
- `/home/user/go/bin/deadcode .`: the sandbox attempt could not read the Go build cache; after approval, exit 0 with only the three `internal/migrate` findings named above and no `internal/data` findings.
- Targeted Go symbol search found no `LoadProject`, `loadProjectV1`, `loadProjectV2`, or `data.Project` API references. The remaining strings are boundary guards and unrelated V2 board loader test names.

### Read and change record

Task context files read: `.savepoint/Design.md`, `.savepoint/checks/C-914-o021-objective-check.md`, `.savepoint/issues/I-040-data-loadproject-unreachable-v1-dispatch.md`, `.savepoint/issues/I-041-design-board-refresh-names-loadproject.md`, `internal/data/project.go`, `internal/data/runtime.go`, `internal/data/project_test.go`, `internal/data/migration_history_test.go`, `internal/data/migration_source_test.go`, `main_test.go`, `main_resume_matrix_test.go`, `internal/migrate/apply_test.go`, and `internal/init/v2_scaffold_test.go`.

Workflow reads: `AGENTS.md`, `.savepoint/router.md`, `agent-skills/savepoint-task/SKILL.md`, `agent-skills/savepoint-design/SKILL.md`, `.savepoint/Idea.md`, `.savepoint/Guardrails.md`, `agent-skills/references/issue-capture.md`, the O-021 Objective, and T-022. Extra reads and reasons: `internal/data/config.go` to confirm missing-config schema semantics while moving tests to `CheckRuntimeSchema`; `internal/doctor/checks.go` because the removed API search found two active comments naming it; searched for direct `CheckRuntimeSchema` test coverage, found no existing test file, and moved the former wrapper cases to `project_test.go`.

Files changed: `.savepoint/Design.md`, `.savepoint/router.md`, both Issue records I-040 and I-041 (task links and repair evidence), this Task record, `internal/data/project.go`, `internal/data/project_test.go`, `internal/data/migration_history_test.go`, `internal/data/migration_source_test.go`, `internal/doctor/checks.go`, `internal/init/v2_scaffold_test.go`, `internal/migrate/apply_test.go`, `main_test.go`, and `main_resume_matrix_test.go`.

At task start, `git status --short` showed the existing T-026 edits to `AGENTS.md`, `README.md`, `.savepoint/Design.md`, and T-026 itself, plus the unrelated existing I-031 edit. Those changes were preserved.

### Limitations and handoff

The owner explicitly waived T-027's optional Task Check and directed the evidence to the mandatory Full Objective Check. I-040 and I-041 remain open for that independent Check to verify and resolve. This executor has not granted clearance or marked the Task done. `stage: audit` means ready for that handoff, not passed.

## Drift Notes

This Task closes the two related findings from C-914 without changing O-021's confirmed outcome or success conditions.
