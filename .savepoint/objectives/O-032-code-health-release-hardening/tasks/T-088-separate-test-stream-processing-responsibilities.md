---
id: T-088
title: Separate test-stream processing responsibilities
objective: O-032
status: planned
complexity_tier: high
complexity_reason: Streaming subprocess failures and report completion have multiple error and output ownership boundaries.
owner_validation: {required: false}
planned_by: {role: planner, session: planning-o032-20261002-owner-confirmed}
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

Pending execution: named per-criterion cases, command/time/toolchain/results, files read/changed, decision deliverable where applicable and limitations.

## Drift Notes

Record implemented responsibility or interface changes for planner reconciliation before the mandatory Full Objective Check. If evidence requires a material unknown repair, return REPLAN REQUIRED rather than inventing another scope.
