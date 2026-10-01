---
id: I-106
title: Interrupted test subprocess overwrites the complete test report
type: defect
status: open
source:
  kind: check
  check: C-948
  actor: {role: checker, session: check-o031-20261002-independent}
  at: '2026-10-01T20:47:19Z'
tasks: [T-074]
checks: [C-948]
guardrail_ids: [TEST-01, TEST-02]
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
