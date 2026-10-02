---
id: I-120
title: Test-stream writer failures are ignored or replace the child failure
type: defect
status: resolved
source:
  kind: check
  check: C-957
  actor: {role: checker, session: check-o032-20261002-independent}
  at: '2026-10-02T05:23:00Z'
tasks: [T-088]
checks: [C-957, C-958, C-960]
guardrail_ids: [TEST-01, TEST-02]
severity: medium
resolution:
  disposition: verified
  check: C-960
  actor: {role: checker, session: recheck-o032-20261003-independent}
  at: '2026-10-02T19:20:00Z'
  reason: "Code unchanged since C-958 proof; broken-sink stream tests pass in a fresh make ci."
history:
  - at: '2026-10-02T05:23:00Z'
    actor: {role: checker, session: check-o032-20261002-independent}
    kind: observed
    check: C-957
    note: Initial independent Full O-032 Check; frozen scope and reproducible harness in C-957.
  - at: '2026-10-02T05:27:50Z'
    actor: {role: executor, session: repair-o032-20261002}
    kind: repair_attempted
    note: 'runGoTestStream wraps output in a first-error recorder, runs every output step, and joins the child ExitError with read and output failures. Added broken-sink tests for passing child, failing child and malformed passthrough.'
  - at: '2026-10-02T06:19:18Z'
    actor: {role: checker, session: recheck-o032-20261002-independent}
    kind: rechecked
    check: C-958
    note: "Original broken output sink reproductions pass, child ExitError remains discoverable; fresh full stream tests pass. Technically proven, open pending CLEAR Check proof."
  - at: '2026-10-02T19:20:00Z'
    actor: {role: checker, session: recheck-o032-20261003-independent}
    kind: rechecked
    check: C-960
    note: "CLEAR. Code unchanged since C-958 proof; broken-sink stream tests pass in a fresh make ci."
---

# I-120: Test-stream writer failures are ignored or replace the child failure

## Summary

T-088 Done When 2 promises clear primary handling of writer and process failures. The refactored stream silently succeeds on timing-output failure, and a diagnostic-output failure replaces the actual child ExitError.

## Evidence

C-957 `TestO032StreamOutputFailureMatrix`: `runGoTestStream(streamHelperCommand("ok"), failingWriter{}, nil)` returns nil even though every Write returns `disk full`. `writeTestTimingSummary` (`internal/buildtool/main.go:318-324`) ignores all print errors; malformed passthrough at :270-272 also ignores its sink failure. With helper `failed-package`, the child exits 1, but the same failing output writer returns only `write failed package output: disk full`; `errors.As(err, *exec.ExitError)` is false. `runGoTestStream` has already obtained `waitErr` at :299 but returns the output error at :300-305 before preserving it. Existing new tests cover a failing report tee and an oversized event with a good diagnostic sink, and therefore miss these two cells. Normal event aggregation/tee ordering passed an independent expected-values probe. Some behavior predates the refactor; it is admitted because T-088 explicitly promises this failure contract and touched these paths.

## Proof Needed

Propagate diagnostic/timing output failures and preserve the primary read/tee/child error when secondary output also fails (joining or otherwise retaining it). Add successful-child/broken-output and failed-child/broken-output cases, plus malformed passthrough with failing sink. Preserve report completion, cleanup and process drain behavior; verify no partial report publication.

Repair directly under this Issue by default; preserve all owner-completed Task statuses. This record grants no clearance or owner acceptance.
