---
id: I-116
title: Health replacement concurrency test skips a supported Windows case
type: guardrail
status: open
source:
  kind: check
  check: C-957
  actor: {role: checker, session: check-o032-20261002-independent}
  at: '2026-10-02T05:23:00Z'
tasks: [T-090, T-091]
checks: [C-957]
guardrail_ids: [CFG-03, CFG-02, TEST-02]
severity: medium
history:
  - at: '2026-10-02T05:23:00Z'
    actor: {role: checker, session: check-o032-20261002-independent}
    kind: observed
    check: C-957
    note: Initial independent Full O-032 Check; frozen scope and reproducible harness in C-957.
  - at: '2026-10-02T05:27:50Z'
    actor: {role: executor, session: repair-o032-20261002}
    kind: repair_attempted
    note: 'Removed the blanket Windows skip; the test now accepts a named Windows refusal (access denied or sharing violation) and still checks complete old/new bytes and no temp files. Cross-vets for windows; no native run yet (see I-115).'
---

# I-116: Health replacement concurrency test skips a supported Windows case

## Summary

`TestReadersSeeCompleteConfigAndReportDuringReplacement` skips every Windows run even though reading health config/report while refresh replaces them is a supported situation. CFG-03 permits a Windows skip only when the situation cannot exist there. A possible OS refusal needs tested preservation and a clear diagnostic, as T-090 Done When 2 already allows.

## Evidence

`internal/codehealth/storage_failure_test.go:192-195` unconditionally calls `t.Skip("Windows refuses to replace a file another handle holds open")`. The remaining test runs config/report replacement alongside four readers. This situation can occur in the Windows board and CLI; the skip reason describes an expected failure behavior, not an impossible configuration. `Store.SaveConfig` at `storage.go:179-195` and `Store.WriteReport` replace through rename and must preserve old bytes and report refusal. Linux `make ci` and race checks pass, but even a future Windows green run would omit this cell. No native result is claimed.

## Proof Needed

Exercise replacement with readers on Windows, accepting complete old/new data or a named platform refusal while proving retained owner bytes and temp cleanup. Remove the blanket skip for this supported case; do not force Windows to behave like Unix. Include the new test in native evidence for I-115.

Repair directly under this Issue by default; preserve all owner-completed Task statuses. This record grants no clearance or owner acceptance.
