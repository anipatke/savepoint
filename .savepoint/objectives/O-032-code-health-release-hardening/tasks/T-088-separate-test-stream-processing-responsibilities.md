---
id: T-088
title: Separate test-stream processing responsibilities
objective: O-032
status: done
complexity_tier: high
complexity_reason: Streaming subprocess failures and report completion have multiple error and output ownership boundaries.
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: planning-o032-20261002-owner-confirmed}
check_waiver:
    task: T-088
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-02T04:04:59Z"
---

# Separate test-stream processing responsibilities

## Outcome

The build tool preserves test timing, failure diagnostics and atomic test reports with focused helpers instead of a large streaming orchestration function.

## User Check

Run a temporary module with failing tests and an interrupted run; verify diagnostics and previous complete reports survive.

## Done When

- Separate event decoding/aggregation and failed-package output from runGoTestStream orchestration, preserving subprocess lifecycle and output order.
- Malformed/oversized event, scanner failure, writer failure, process nonzero exit and interruption retain clear primary failure handling; no partial output becomes a completed health test report.
- Existing complete-report preservation, temporary-file cleanup, focused-no-report, timing and package failure regressions remain valid; add exact missing failure evidence rather than implementation-mirroring tests.
- Record before/after Lizard complexity on the same scope; target touched functions <=20 CCN, aiming for <=10, without deleting failure handling.

## Context Files

`internal/buildtool/main.go`; `internal/buildtool/main_test.go`.

## Design References

O-032 Confirmed Release-Hardening Design and Owner-Requested Health Analysis; Design sections 1, 6, 7, 9, 12 and 13 as applicable.

## Guardrails

FS-01, FS-04, FS-05, FS-06, ARCH-01, ARCH-02, ARCH-03, CFG-01, CFG-02, TEST-01, TEST-02, TEST-03, TEST-04, TEST-05, TEST-08, TEST-09, STYLE-01, STYLE-02, STYLE-07.

## Implementation Plan

1. Trace streaming process/read/wait/output order and report completion tests.
2. Extract coherent accumulator and output helpers while keeping lifecycle ownership in orchestration.
3. Probe malformed stream and output/process failures with helper subprocesses.
4. Run buildtool focused tests and handoff gate.

## Boundaries

No gate definition changes, timing format redesign, provider behavior or package/distribution changes.

## Technical Verification

For before/after complexity evidence, the owner-confirmed scope authorizes a direct project-owned Lizard invocation over the scoped source paths with temporary output; run no Savepoint health command. If the executable is unavailable, record the measurement gap for the owner instead of installing it.

`make build && make test-fast` for ordinary handoff; focused `make test-focused TEST=... PKGS=...` only during iteration. Record named per-criterion evidence and limitations. The later independent Check applies `agent-skills/references/check-method.md`; the mandatory Full Objective Check runs its fresh full gate and official Code Health collection. This executor never writes a Check or runs savepoint health setup/report/check.

## Technical Evidence

Executed 2026-10-02 by executor; toolchain go1.26.2 linux/amd64, lizard 1.24.0.

Per criterion:
1. Separation: `runGoTestStream` now owns only pipe/start/wait/output ordering. Event aggregation moved to `testStreamSummary` (`add`, `addOutput`, `addPackageResult`, `addTestResult`, `writeFailedPackages`); stream reading/tee moved to `consumeTestStream`. Output order unchanged: malformed lines live, then failed-package output (sorted, deduped), stderr, timing summary. Evidence: `TestRunGoTestStreamPrintsMalformedLinesAndFailedPackageOutputInOrder` plus existing `TestRunGoTestCommandReportsTimingAndPropagatesFailure`.
2. Failure handling: malformed line passthrough (new test), tee/report write failure reported after a successful child with timing still printed (`TestRunGoTestStreamReportsReportWriteFailureAfterChildSucceeds`), oversized line reported as `read go test output` (`TestRunGoTestStreamReportsOversizedEventLine`), nonzero exit propagates as the primary error (existing), interruption/early exit keep prior reports (existing `...WhenInterrupted`, `...WhenChildExitsEarly`). No partial stream becomes a report: `goTestRunCompleted`/`goTestStreamComplete` untouched.
3. Existing regressions (complete-report preservation, temp cleanup, focused-no-report, timing, package failure) pass unchanged. Four tests added, each for evidence that was missing.
4. Lizard (`lizard internal/buildtool/main.go`, same file): before `runGoTestStream` CCN 31 / 87 NLOC; after CCN 8 / 30 NLOC. New helpers max CCN 8 (`consumeTestStream`), `writeFailedPackages` 5, others <=5. Untouched `runGoTestWithReports` and `goTestStreamComplete` remain at CCN 10.

Behavior fix found while testing: a single event line over the 10 MB scanner limit stopped reading and left the child blocked on a full pipe, so `cmd.Wait` hung forever. `consumeTestStream` now drains the pipe after a scan error, preserving the read error as the failure. This is a small necessary change within the stream-processing scope.

Commands: `make test-focused TEST='TestRunGoTest|TestGoTest' PKGS=./internal/buildtool` passed; `make build && make test-fast` exit 0.

Files read: internal/buildtool/main.go, internal/buildtool/main_test.go (Context Files). Extra reads: AGENTS.md, agent-skills/savepoint-task/SKILL.md, router.md, O-032 Objective.md (workflow routing and Objective boundaries). Files changed: internal/buildtool/main.go, internal/buildtool/main_test.go, this Task.

Limitations: no real-module interrupted-run manual User Check was performed; interruption is covered by the existing helper-subprocess tests. `make test-full` not run (ordinary handoff). The 11 MB oversized test emits ~11 MB through a pipe (about 40ms).

## Drift Notes

Record implemented responsibility or interface changes for planner reconciliation before the mandatory Full Objective Check.

- Added `testStreamSummary` and `consumeTestStream` in `internal/buildtool/main.go`; `runGoTestStream` signature unchanged. Oversized stream lines now drain the pipe instead of deadlocking Wait. If evidence requires a material unknown repair, return REPLAN REQUIRED rather than inventing another scope.
