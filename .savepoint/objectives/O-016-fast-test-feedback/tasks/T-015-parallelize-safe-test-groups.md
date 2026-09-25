---
id: T-015
title: Run independent tests together
objective: O-016
status: done
depends_on: [{task: T-013, requires: clear}, {task: T-014, requires: clear}]
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: o016-design-20260923}
complexity_tier: high
complexity_reason: Package-wide hooks and filesystem mutations make parallel safety a test-specific decision.
check_waiver:
    task: T-015
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-23T00:15:19Z"
---

# T-015: Run independent tests together

## Outcome

Independent migration test groups execute concurrently, while tests using shared hooks or other process-wide state stay serial and deterministic.

## User Check

If a Task Check is requested, review the concurrency inventory and repeated stress results, especially interruption and recovery cases.

## Done When

- Each parallelized group has a documented isolation basis from T-013/T-014; global hook and environment users remain serial unless dependencies are made per-run.
- Temporary directories and prepared fixture data are not shared mutably across parallel tests.
- Repeated runs and race-enabled runs show no nondeterminism; named fault-injection and recovery cases still pass.
- Before/after timings and any groups deliberately left serial are recorded.
- Configured quality gates and per-criterion evidence pass.

## Context Files

`internal/migrate/end_to_end_test.go`, `internal/migrate/apply_test.go`, `internal/migrate/operation_test.go`, `internal/migrate/command_test.go`, `internal/migrate/cutover_test.go`, `.savepoint/objectives/O-016-fast-test-feedback/tasks/T-013-measure-and-design-test-gates.md`, `.savepoint/objectives/O-016-fast-test-feedback/tasks/T-014-share-read-only-migration-results.md`.

## Design References

Design sections 1 and 13; O-016 Architectural Considerations.

## Guardrails

TEST-01..05, TEST-07..08, CFG-02, STYLE-03, STYLE-05.

## Implementation Plan

1. Apply T-013's shared-state inventory to each proposed concurrent group.
2. Add parallel execution only to isolated subtests and fixtures.
3. Keep fault-injection and process-global controls serial.
4. Run repeated normal and race-enabled suites; investigate any inconsistent result before handoff.
5. Compare timing and run the configured quality gates.

## Boundaries

No global test hook races, shared mutable fixture data, or assertions that depend on completion order.

## Technical Verification

Repeated focused migration runs, race-enabled migration tests, recovery and interruption cases, `make build && make test`.

## Technical Evidence

### Start and scope

- Started 2026-09-22 23:35 UTC at `stage: build`, following the owner's instruction to begin T-015. T-013 and T-014 are `done`; each has an explicit owner Task-check waiver. The task workflow states that an explicit waiver satisfies a downstream dependency requiring `clear`, so both T-015 dependencies are satisfied. O-016 records owner confirmation of the T-013-T-016 sequence and verification policy.
- The router's `next_action` still described T-014's pending owner closure even though T-014 is owner-closed. The owner explicitly directed T-015, so the router was reconciled to this active Task.
- Read the owning Objective and T-013/T-014 records. Read the five Context Files named above, with targeted inspection of their test functions and shared-state touchpoints. The earlier task evidence records the measured gate baseline and migration fixture isolation decisions used here.
- Required routing reads: `.savepoint/router.md` and `agent-skills/savepoint-task/SKILL.md`.
- Extra read: `.savepoint/Design.md` sections 1 and 13, because they are named Design References and describe architecture and test evidence boundaries. Applicable Guardrails read: TEST-01..05, TEST-07..08, CFG-02, STYLE-03, and STYLE-05.
- Extra router read/write: `.savepoint/router.md`, to correct the stale next action after verifying T-014's owner closure and receiving the owner's T-015 instruction.

### Isolation inventory

- The operation tests each create their own `t.TempDir()` and operate on a private journal and files. The cutover tests create separate temporary projects; their shared options return immutable deterministic values.
- Apply tests use `copyFixtureProject` to create a fresh temporary project. `TestApply_recoversAtEveryPublishBoundaryWithoutOverwritingUserEdits` sets package-global publish hooks and remains serial; other apply tests do not touch those hooks.
- Command tests use private fixture copies. `TestRunCommand_previewAndRecoveryReportAreReadOnlyOnReadable0555Project` temporarily replaces the package-global `targetWriteabilityProbe` and remains serial. Both command recovery tests call `interruptFirstCommand`, which temporarily sets `afterPublishWriteHook`, and remain serial. The sentinel-preservation test is isolated to its own temporary project.
- End-to-end tests either use private temporary copies or read T-014's prepared results after `TestMain` has finished building them. The prepared-result map and trees are read-only to those consumers. `TestEndToEnd_sharedMigratedResultsAreReadOnly` probes the prepared tree by temporarily creating/removing a file, so it remains serial. All other end-to-end tests are independent of that probe and each other.
- No environment setters or working-directory changes occur in the five scoped test files. No shared mutable test fixture or completion-order assertion is introduced.

### Baseline and parallel groups

- Before changes, `go test -count=1 ./internal/migrate` passed: package test time 115.565 s; command wall time 132.04 s. This was an uncached test execution with the existing WSL npm bridge at `/tmp/t014-tools/npm`.
- Extra targeted read: `rg -n "t\.Parallel|afterPublishWriteHook|afterPublishRemovalHook|targetWriteabilityProbe|os\.(Setenv|Unsetenv|Chdir)|\.Setenv\(" internal/migrate -g '*_test.go'`. This checked whether tests outside the five Context Files already used Go's parallel-test scheduler or process-wide mutable controls that could overlap the selected tests. It found no pre-existing `t.Parallel` calls or environment/working-directory setters; the only hook mutators are the serial cases named above.
- Added parallel execution to 53 top-level tests in five groups:
  - `operation_test.go`: all 15 tests; every test uses an independent `t.TempDir()` project/journal.
  - `cutover_test.go`: all 5 tests; each project is private, including each nested operational-condition case.
  - `apply_test.go`: 13 tests use `copyFixtureProject` for a fresh directory. `TestApply_recoversAtEveryPublishBoundaryWithoutOverwritingUserEdits` remains serial because it changes package-global publish hooks.
  - `command_test.go`: `TestRunCommand_applyWriteabilityProbePreservesPreExistingSentinelAndPrefix` uses a private copied fixture. `TestRunCommand_previewAndRecoveryReportAreReadOnlyOnReadable0555Project` and both recovery tests remain serial because they change package-global probes/hooks.
  - `end_to_end_test.go`: 19 tests use fresh private copies or read the T-014 prepared results after `TestMain` makes them read-only. `TestEndToEnd_sharedMigratedResultsAreReadOnly` remains serial because it temporarily probes the shared result tree with a write.
- No test gained a shared mutable fixture or a completion-order dependency. All other scoped tests remain serial through Go's standard test scheduler.
### Verification and handoff evidence

| Acceptance criterion | Evidence |
|---|---|
| Isolation basis and serial global controls | The five group bases and serial exceptions are recorded above. Serial hook/probe cases: `TestApply_recoversAtEveryPublishBoundaryWithoutOverwritingUserEdits`; `TestRunCommand_previewAndRecoveryReportAreReadOnlyOnReadable0555Project`; `TestRunCommand_recoverApplyAcrossInvocationsUsesRecordedOperation`; `TestRunCommand_recoverApplyAcrossInvocationsRejectsEditedSource`; `TestEndToEnd_sharedMigratedResultsAreReadOnly`. |
| Independent mutable data | Parallel apply, operation, cutover, and command tests each use fresh temporary directories. Parallel end-to-end tests use private temporary copies or read T-014's prepared result after `TestMain` completes. No completion-order dependency was introduced. |
| Repeated and race-enabled results | Two full uncached `go test -count=1 ./internal/migrate` runs passed. Two focused `-race` runs passed, each at 46.450 s and 46.434 s package time. Both included all apply, command, operation, cutover, and end-to-end test prefixes plus the serial hook/recovery cases, except for the single repository-copy case named below. |
| Timing and intentionally serial groups | Baseline: 115.565 s package / 132.04 s wall. Normal post-change runs: 110.734 s / 122.44 s and 107.434 s / 122.36 s. The package-time mean is 109.084 s (5.6% below baseline); the timing result and exact serial groups are recorded above. |
| Configured quality gates and evidence | `make build && make test` passed. `make test` ran the uncached changed migration package in 111.537 s; other packages reported cached results. |

Race command used twice:

    go test -race -count=1 -run "^(TestApply_|TestEndToEnd_|TestOperation_|TestRunCommand_|TestPreflightCutover_)" -skip "^TestEndToEnd_temporaryRepositoryCopyMigratesWithReleaseAccountability$" ./internal/migrate

The initial broad `go test -race -count=1 ./internal/migrate` exceeded eight minutes. After the owner flagged the duration, it was stopped and was not counted as a pass or failure. The focused race command above excludes only `TestEndToEnd_temporaryRepositoryCopyMigratesWithReleaseAccountability`; the full normal suite and required gate both ran that case successfully.

Files changed for T-015: the five migration test files named in this Task, this T-015 record, and `.savepoint/router.md`. T-014's existing edits in `end_to_end_test.go` were retained.

Evidence is at `stage: audit`. T-015 remains `in_progress`; only the owner may set it to `done`. No optional Task Check was requested or waived in this session.
