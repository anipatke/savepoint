---
id: I-152
title: Runner reports success after heartbeat or final record writes fail
type: defect
status: open
source:
  kind: check
  check: C-978
  actor: {role: checker, session: check-o045-20261011}
  at: '2026-10-10T21:10:26Z'
tasks: [T-127, T-128]
checks: [C-978]
guardrail_ids: [FS-06, TEST-01, TEST-02]
history:
  - at: '2026-10-10T21:10:26Z'
    actor: {role: checker, session: check-o045-20261011}
    kind: observed
    check: C-978
    note: Initial Full Objective Check of O-045.
---

# I-152: Runner reports success after heartbeat or final record writes fail

## Summary

A record failure must be surfaced and must not present successful saved completion. The heartbeat error, summary error and final record error are discarded; Wait promises a final saved record that does not exist.

Violated requirements: O-045 SC1 and required run artifacts; T-127 DW3; T-128 DW3.

## Evidence

- Location: `internal/run/run.go:185,195-196,221-224`.
- Smallest reproduction: Start a slow fake provider, replace run.yml with a directory after Start, and wait. Wait returns finished while run.yml is a directory. With a 20ms heartbeat it continues with the agent alive after the failed heartbeat and still returns finished. Precreate summary.md as a directory and it also returns finished.
- Independent harness and full output are embedded in C-978. Probe PASS means the observation ran successfully, not that the requirement passed.

## Proof Needed

Exercise initial/PID/heartbeat/summary/final writes, including failure after the child has started and cleanup failure; preserve the primary diagnostic, avoid claiming persisted success, and prove the chosen stop behavior and temp-file cleanup. Existing tests assert only successful persistence.
