---
id: I-106
title: Interrupted test subprocess overwrites the complete test report
type: defect
status: resolved
source:
  kind: check
  check: C-948
  actor: {role: checker, session: check-o031-20261002-independent}
  at: '2026-10-01T20:47:19Z'
tasks: [T-074]
checks: [C-948, C-950, C-951]
guardrail_ids: [TEST-01, TEST-02]
resolution:
  disposition: verified
  check: C-951
  actor: {role: checker, session: recheck-o031-final-repair-20261002-independent}
  at: '2026-10-01T21:56:44Z'
  reason: Original frozen-scope repair proof passed; C-951 records current independent CLEAR.
history:
  - at: '2026-10-01T20:47:19Z'
    actor: {role: checker, session: check-o031-20261002-independent}
    kind: observed
    check: C-948
    note: Reproduced inside the frozen initial Full Objective Check scope of O-031.
  - at: '2026-10-02T07:30:00Z'
    actor: {role: executor, session: user-request}
    kind: repair_attempted
    note: 'runGoTestWithReports keeps the previous go-test.json and coverage.out unless go test ran to its own exit (goTestRunCompleted: success or non-signal exit). Added TestRunGoTestWithReportsKeepsCompleteReportsWhenInterrupted with a fake go that SIGTERMs itself. Native Windows evidence is still needed.'
  - at: '2026-10-01T21:45:21Z'
    actor: {role: checker, session: native-o031-recheck-20261002}
    kind: rechecked
    check: C-950
    note: Native Windows TerminateProcess reproduction still replaces the previous complete go-test.json with a start-only stream; ordinary full Windows suite passes but skips its interruption test.
  - at: '2026-10-01T21:56:44Z'
    actor: {role: checker, session: recheck-o031-final-repair-20261002-independent}
    kind: rechecked
    check: C-951
    note: Exact native Windows TerminateProcess repro now preserves both reports; Linux interruption, portable native early-exit and completed pass/fail report tests pass.
---

# I-106: Interrupted test subprocess overwrites the complete test report

## Summary

Interrupted test subprocess overwrites the complete test report. Repair directly under this Issue; keep the completed Task status unchanged.

## Evidence

T-074 Done When 2 requires an interrupted run never to replace a complete report with a half-written report.

`internal/buildtool/main.go:107-113` calls runGoTestStream, closes the temporary JSON file, and renames it over go-test.json regardless of whether the child completed or was killed. Atomic rename protects bytes from torn writes but does not establish a complete test event stream.

Independent `TestO031IndependentInterruptedReports` (embedded in C-948) starts with a sentinel complete go-test.json and places a controlled executable named go on PATH in a temporary directory. The executable emits only a package start event and terminates itself with SIGTERM. runGoTestWithReports returns signal: terminated, yet replaces the complete file with `{"Action":"start","Package":"fixture"}` plus newline. Expected: preserve the old complete report on process interruption, while still publishing completed passing/failing test runs. This simulates the finite subprocess interruption boundary without timing races or killing the review process.

Existing report tests cover ordinary test pass, ordinary failure, focused no-report, and temporary cleanup; none verifies interrupted child preservation. Loss is contained to a generated report, but it defeats the explicit evidence preservation guarantee and removes the last complete test evidence.

## Proof Needed

Prove completed pass/fail runs still replace both reports, while an interrupted or incomplete test subprocess keeps previous complete reports byte-identical. Verify temporary cleanup and nonzero exit. Native Windows equivalent and repeat replacement evidence remain needed for the platform-sensitive gate.


## Native Windows recheck — C-950

The full native suite passed on go1.26.2 windows/amd64 in an isolated Windows-local temporary copy of the reviewed source. The current regression TestRunGoTestWithReportsKeepsCompleteReportsWhenInterrupted skips Windows because its fixture is a shell script; this does not establish the promised Windows behavior.

A Windows equivalent of C-948 B2 was compiled with a Go overlay and executed natively through PowerShell. The fake go.exe emitted only a package start event, then called Windows TerminateProcess on itself with exit code 1. Both destination reports started with a complete sentinel. Expected: nonzero error and both old reports byte-identical. Actual: returned exit status 1 and replaced go-test.json with `{"Action":"start","Package":"fixture"}` plus newline. coverage.out stayed unchanged in this fixture because the child wrote no profile.

`internal/buildtool/main.go:131-136` treats any exec.ExitError with ExitCode >= 0 as a completed run; on Windows a terminated process returns a nonnegative code, unlike Unix signal termination. The resulting rename at lines 113-117 publishes the incomplete JSON stream. This remains the same I-106 preservation requirement (T-074 Done When 2), in the originally frozen native Windows subprocess boundary; it is not a new Issue or a widened scope.

Independent command: native o031-buildtool-windows.exe `-test.run=^TestO031NativeWindowsInterruptedReports$ -test.v -test.count=1`. FAIL, 0.30 seconds. Exact harness and output are retained in C-950. A platform-independent completeness check needs to preserve abnormal/incomplete reports while continuing to publish completed passing and failing runs. Keep T-074 done; repair directly under I-106, then obtain new independent CLEAR proof.
