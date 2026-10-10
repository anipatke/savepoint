---
id: I-149
title: Lost-run detection watches the agent instead of the runner
type: defect
status: resolved
source:
  kind: check
  check: C-978
  actor: {role: checker, session: check-o045-20261011}
  at: '2026-10-10T21:10:26Z'
tasks: [T-129, T-130]
checks: [C-978, C-979]
guardrail_ids: [TEST-01]
resolution:
  disposition: verified
  check: C-979
  actor: {role: checker, session: check-o045-recheck-20261011}
  at: '2026-10-10T21:44:08Z'
  reason: 'Actual helper runner exits with live agent; saved runner_pid now makes RunnerLost and LostRuns report lost.'
history:
  - at: '2026-10-10T21:10:26Z'
    actor: {role: checker, session: check-o045-20261011}
    kind: observed
    check: C-978
    note: Initial Full Objective Check of O-045.
  - at: '2026-10-10T21:44:08Z'
    actor: {role: checker, session: check-o045-recheck-20261011}
    kind: rechecked
    check: C-979
    note: 'Actual helper runner exits with live agent; saved runner_pid now makes RunnerLost and LostRuns report lost.'
---

# I-149: Lost-run detection watches the agent instead of the runner

## Summary

A dead heartbeat writer with a stale record must be reported as runner_lost even while its agent survives. The recorded pid is the agent PID; doctor therefore omits precisely this lost-run case.

Violated requirements: O-045 SC6; T-129 DW3; T-130 DW6.

## Evidence

- Location: `internal/run/run.go:158 and internal/run/stop.go:93-101`.
- Smallest reproduction: Start a fake agent from a helper runner, let the helper exit without waiting, and retain its stale running record. The real subprocess probe observed runnerPID=145 alive=false, agentPID=151 alive=true, RunnerLost=false and LostRuns=[].
- Independent harness and full output are embedded in C-978. Probe PASS means the observation ran successfully, not that the requirement passed.

## Proof Needed

Store and test the runner identity separately from the agent identity. Repeat dead runner/live agent, live runner/live agent, dead both, and fresh/exact30s/over30s cases. Existing tests only check one PID, not independent runner and agent lifetimes.
