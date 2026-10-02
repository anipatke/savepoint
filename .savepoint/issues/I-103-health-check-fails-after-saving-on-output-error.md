---
id: I-103
title: Health check fails after saving a snapshot when output fails
type: defect
status: resolved
source:
  kind: check
  check: C-946
  actor: {role: checker, session: check-o030-20261001-independent}
  at: '2026-10-01T20:06:17Z'
tasks: [T-068]
checks: [C-946, C-947]
guardrail_ids: [TEST-01, TEST-02]
resolution:
  disposition: verified
  check: C-947
  actor: {role: checker, session: recheck-o030-20261001-independent}
  at: '2026-10-01T20:15:49Z'
  reason: Original reproduction and frozen-scope recheck passed; C-947 records CLEAR.
history:
  - at: '2026-10-01T20:06:17Z'
    actor: {role: checker, session: check-o030-20261001-independent}
    kind: observed
    check: C-946
    note: Reproduced during the initial Full Objective Check of O-030.
  - at: '2026-10-01T20:15:49Z'
    actor: {role: checker, session: recheck-o030-20261001-independent}
    kind: rechecked
    check: C-947
    note: Repair independently verified against C-946's frozen scope and a fresh full gate.
---

# I-103: Health check fails after saving a snapshot when output fails

## Summary

T-068 Done When 4 and the confirmed O-030 command contract require exit 0 whenever a snapshot was saved, with non-zero exit only when no snapshot was saved.

## Evidence

`internal/healthcheck/healthcheck.go:98-99` returns the stdout write error after successful immutable persistence and evaluation; `main.go:93-100` maps that error to exit 1.

Smallest reproduction: temporary Git V2 project with O-001, passing Go test report, configured fake Lizard runner, and an `io.Writer` whose Write returns `(0, errors.New("output sink failed"))`. Call `healthcheck.Run`. Actual: `error=output sink failed`, while `Store.LoadSnapshots()` returns exactly one permanent official snapshot. Expected under the explicit contract: saved collection stays successful; report the secondary output problem without claiming collection failed. A redirected failing sink or closed pipe is a supported output failure path. Because no snapshot identity is printed, the checker cannot infer collection status from this error.

Independent command: `go test ./internal/healthcheck -run TestO030WorkflowMatrix -v -count=1`, output-failure-after-persistence cell fails. Harness is embedded in C-946. Existing tests cover pre-save collection failure and blocking verdict success, but have no failing output writer after save.

Scope: post-save output failure is the reproduced case. Load/evaluation post-save errors at :75-91 also contradict the unconditional comment; no concurrent corruption scenario was reproduced and it is not an additional finding.

## Proof Needed

Prove exit/success semantics for a saved snapshot followed by an output error, preserving its identity through a secondary diagnostic where possible. Keep genuine no-snapshot failures non-zero, blocking verdict success zero, and avoid deleting or rewriting the saved official snapshot. Reconcile any proposed contract change with the owner rather than silently changing acceptance.
