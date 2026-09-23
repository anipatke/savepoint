---
id: T014
title: Reuse migration fixture results
objective: O016
status: done
depends_on: [{task: T013, requires: clear}]
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: o016-design-20260923}
complexity_tier: high
complexity_reason: Migration assertions share expensive setup while mutation and recovery evidence must stay isolated.
check_waiver:
    task: T014
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-22T23:31:39Z"
---

# T014: Reuse migration fixture results

## Outcome

Read-only end-to-end migration assertions reuse independently prepared converted fixture results where safe, reducing repeated copies and full applies without weakening their assertions.

## User Check

If a Task Check is requested, compare the before/after test inventory and a representative converted fixture result; confirm mutating recovery scenarios still use fresh directories.

## Done When

- The T013 decision identifies every assertion moved to shared read-only preparation.
- Shared preparation cannot be mutated by a consumer, and each mutating, interruption, recovery, no-op, and user-file-safety scenario retains isolated input.
- Failure messages still identify the fixture and assertion; test ordering cannot change results.
- Before/after timing and the retained test inventory are recorded.
- Named happy and failure paths, configured quality gates, and per-criterion evidence pass.

## Context Files

`internal/migrate/end_to_end_test.go`, `internal/migrate/apply_test.go`, `internal/migrate/fixture_test.go`, `internal/migrate/operation_test.go`, `.savepoint/objectives/O016-fast-test-feedback/tasks/T013-measure-and-design-test-gates.md`.

## Design References

Design sections 1 and 13; O016 Architectural Considerations.

## Guardrails

FS-01, TEST-01..05, TEST-07..08, STYLE-03, STYLE-07.

## Implementation Plan

1. Use T013's assertion map to group read-only checks by fixture and one prepared result.
2. Keep every write or fault-injection scenario on an independent temporary project.
3. Preserve named subtests and fixture-specific diagnostics.
4. Compare test inventory and timing before and after; rerun recovery and file-preservation cases.
5. Run the currently configured quality gates before handoff.

## Boundaries

No production migration behavior change, skipped assertion, shared mutable fixture, or test-order dependency.

## Technical Verification

Focused migration test cases, full migration package, filesystem preservation and recovery cases, `make build && make test`.

## Technical Evidence

### Start and shared preparation

- T013 is owner-closed with an explicit optional Task Check waiver. Under the active savepoint-design verification contract, that waiver satisfies T014's `requires: clear` dependency as the owner's completion decision; it does not claim technical CLEAR.
- `TestMain` now copies each frozen V1 fixture once, applies migration through the real command path with the fixed clock and operation ID, and makes the converted result read-only before tests start. It removes the temporary result tree after `m.Run`.
- On Unix, read-only assertions use the shared prepared result directly. `TestEndToEnd_sharedMigratedResultsAreReadOnly` checks every tree entry's write bits and probes directory and file write access; if the test process is privileged, it cleans up its probe and relies on the verified read-only modes. On Windows, consumers get a private copy of the prepared converted result because directory mode bits do not reliably prevent new files there. This Windows path is implemented but was not exercised in this WSL run.

### Assertion mapping

These named read-only assertions now consume the prepared result rather than copying and applying the frozen fixture again: `TestEndToEnd_migratedProjectLoadsCleanThroughLoadV2Index`, `TestEndToEnd_everyReferenceResolves`, `TestEndToEnd_noConvertedRecordCarriesClearance`, `TestEndToEnd_convertedActiveTaskIsNotCompletable`, `TestEndToEnd_everyFixtureFileHasAnAccountableDestination`, `TestEndToEnd_archivedFilesMatchTheFixtureManifestHashes`, `TestEndToEnd_recurringSharedTaskIsRecordedPerRelease`, `TestEndToEnd_legacyFactsAreResolvableThroughTheManifest`, `TestEndToEnd_legacyPrerequisiteIsTyped`, and `TestEndToEnd_goldenConvertedOutput`.

The accountability and golden assertions compute the deterministic plan from the frozen, write-free V1 input, then compare it with the prepared converted output. Archive hashes are checked against the independently authored fixture manifest and the migration manifest's archive paths. The active-Task completion assertion changes only the newly loaded in-memory record; it does not write to the shared result.

### Isolation and retained scenarios

- Command-path, preview, recovery, no-op, interruption, user-file-safety, and deliberately mutating tests retain their own fresh input directories. This includes `TestEndToEnd_bothFixturesMigrateThroughTheCommandPath`, `TestEndToEnd_frozenFixturesAreNeverMutated`, `TestEndToEnd_releaseRecordsMappingsAndCutoverGateAgree`, `TestEndToEnd_temporaryRepositoryCopyMigratesWithReleaseAccountability`, `TestEndToEnd_recurringSharedTaskAllocatesDistinctGlobalIDs`, `TestEndToEnd_goldenIsReproducible`, `TestEndToEnd_secondFullRunChangesNothing`, `TestEndToEnd_previewWritesNothing`, `TestEndToEnd_absentOptionalArtifactsAreNotFindings`, and the Apply and Operation test scenarios in their existing files.
- Each fixture remains a named subtest, and failure messages retain the fixture or assertion label. Preparation finishes before `m.Run`; consumers only read the result, and the scoped migration files contain no `t.Parallel` calls, so test order cannot mutate or race the shared preparation.
- `go test -list . ./internal/migrate` reports 170 top-level entries after the change. T013 recorded 169 before; this adds only `TestEndToEnd_sharedMigratedResultsAreReadOnly`. No existing test name was removed or renamed. Full coverage still runs all tests; the proposed fast selection omits only the same three T013 full-only test parents.

### Before and after timing

All package samples used Go 1.26.2 under WSL2 Ubuntu 24.04.4 on the same host, `-count=1`, and warm Go build cache. The OS page cache was not flushed.

| Measure | T013 warm baseline | T014 after change |
|---|---:|---:|
| `internal/migrate` package test time | 120.869 s, 116.895 s (mean 118.882 s) | 111.750 s, 113.239 s (mean 112.495 s) in two standalone samples; final full gate: 114.873 s |
| Proposed fast selection wall time | 10.98 s, 10.91 s | 7.43 s on final code |
| Proposed fast selection `internal/migrate` time | 10.342 s, 10.267 s | 6.778 s on final code |

The two standalone package samples averaged 6.388 seconds faster (about 5.4%) than baseline. They preceded a small robustness edit that lets privileged test processes clean up their own write probe; the final full gate, after that edit, reported 114.873 seconds, 4.009 seconds (about 3.4%) below the baseline warm mean. The final fast selection wall time fell by 3.515 seconds (about 32.1%).

### Acceptance evidence and commands

1. Assertion map: met. The ten read-only assertions above use one independently prepared converted result per fixture; their original test names and fixture subtests remain.
2. Isolation: met. Shared output is protected before test execution and has an explicit write-rejection test on Unix; Windows consumers receive private copies. Mutating, interruption, recovery, no-op, and user-file-safety cases remain isolated.
3. Diagnostics and order: met by code inspection and focused tests. Fixture subtest names and assertion-specific errors remain; preparation is complete before consumers run; no scoped test uses `t.Parallel`.
4. Timing and inventory: met. Two uncached package samples and the fast candidate were measured after the change; the retained inventory is 170 migration tests versus 169 before, with one guard test added and none removed.
5. Verification: passed. The read-only assertion selection passed on final code (`internal/migrate` 0.528 s). The fresh-directory, no-op, preview, recovery, interruption, and edited-file safety selection passed (`internal/migrate` 16.065 s). Two standalone full-package samples passed before the small privileged-user guard adjustment (111.750 s and 113.239 s); the final code then passed `make build && make test` across all packages, with `internal/migrate` at 114.873 s. The exact proposed fast selection passed on final code in 7.43 s wall time.

The full gate ran under WSL with a temporary `/tmp/t014-tools/npm` wrapper so Windows npm could use WSL's temporary test directory through `pushd`; the isolated `TestIntegration_InstallDependencies` check passed with that wrapper. An initial wrapper attempt failed before the migration package completed; the corrected wrapper and final full gate both passed. No wrapper or generated test output was added to the repository.

### Files read and changed

- Read T014 Context Files: `internal/migrate/end_to_end_test.go`, `internal/migrate/apply_test.go`, `internal/migrate/fixture_test.go`, `internal/migrate/operation_test.go`, and the T013 measurement record.
- Changed for T014: `internal/migrate/end_to_end_test.go` and this T014 record. `.savepoint/router.md` points to T014's audit handoff; the T013 record also retains the owner's closure and waiver. No production source, frozen fixture, Apply/Operation test, or workflow changed.
