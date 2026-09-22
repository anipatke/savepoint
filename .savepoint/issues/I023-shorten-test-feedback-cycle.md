---
id: I023
title: Shorten the repository test feedback cycle
type: verification
status: resolved
source:
  kind: report
  actor: {role: owner, session: user}
  at: '2026-09-22T09:54:20Z'
severity: high
resolution:
  disposition: escalated
  actor: {role: planner, session: i023-escalation-20260922}
  at: '2026-09-22T00:00:00Z'
  reason: Repair promoted into Objective O016; under the new escalation rule the Issue retires immediately rather than staying open for a later Check.
escalated_to: O016
history:
  - at: '2026-09-22T09:54:20Z'
    actor: {role: owner, session: user}
    kind: observed
    note: Local and agent test cycles are taking far too long; the repeated full-suite wait is disrupting ordinary implementation and review work.
  - at: '2026-09-22T10:09:15Z'
    actor: {role: planner, session: o016-fast-test-feedback}
    kind: observed
    note: Planned O016 as the bounded Objective for measurement, fast/full gates, migration-test consolidation, safe parallelism, evidence reuse, and timing visibility.
  - at: '2026-09-22T10:11:38Z'
    actor: {role: owner, session: user}
    kind: deferred
    note: Defer implementation and verified closure to O016; keep I023 open until that Objective supplies repair evidence and a Check verifies it.
  - at: '2026-09-22T00:00:00Z'
    actor: {role: owner, session: user}
    kind: owner_decision
    note: New rule adopted — when an Issue's repair is promoted into an Objective, the Issue retires immediately (disposition escalated) instead of staying open for a later Check. Supersedes the prior deferred entry's plan to keep I023 open until O016 completes.
  - at: '2026-09-22T00:00:00Z'
    actor: {role: planner, session: i023-escalation-20260922}
    kind: escalated
    note: Retired under the new rule; repair is tracked as O016 going forward. O016's own mandatory Full Objective Check and owner acceptance are the proof, not a recheck of this Issue.
---

# I023: Shorten the repository test feedback cycle

Planned remediation Objective: O016 — Make test feedback fast without
weakening verification.

## Summary

The repository's routine verification loop is too slow for normal development.
`make test` is dominated by `internal/migrate`, which has repeatedly taken
about 121 seconds even when the scoped change is unrelated to migration.
Agents then repeat the full suite after small corrections, turning a short
fix-and-review cycle into several minutes of idle waiting.

Savepoint still needs a reliable full integration suite. The repair should
make focused and ordinary handoff feedback fast while retaining an explicit
full/slow gate at the appropriate CI, Objective Check, or release boundary.
It must not silently weaken TEST-08 or allow slow tests to disappear from the
required verification contract.

## Evidence

- C906 recorded `internal/migrate` at 121.009 seconds while the changed behavior
  was confined to board presentation and fixture data.
- The I018 subagent required repeated full verification after correcting one
  planning-record ID collision, producing multiple multi-minute waits.
- `.savepoint/Guardrails.md:90` currently requires `make build && make test`
  before every Task handoff, so the slowest package sits on every change's
  critical path regardless of scope.
- The owner reported that these test cycles are taking "WAY too long" and
  requested durable follow-up.

## Proof Needed

- Measure package and individual-test durations from a cold and warm cache;
  identify the real `internal/migrate` bottlenecks rather than guessing from
  package totals.
- Remove unnecessary sleeps, duplicated matrices, repeated compilation, or
  serial fixture setup where behavior can remain deterministic and race-safe.
- Define a fast local/focused gate with a concrete duration budget and an
  explicit full integration gate that still runs every slow test at the
  correct boundary. Reconcile TEST-08, skills, Make targets, CI, and guidance
  together if the gate names or timing change.
- Prevent agents from rerunning the unchanged full suite after metadata-only
  corrections when current evidence can be safely reused; document exactly
  when a fresh full run is mandatory.
- Add timing/selection tests or CI reporting that makes regressions visible,
  and demonstrate a materially shorter median feedback cycle without reducing
  coverage of migration interruption, recovery, and cross-platform behavior.
- Pass the revised focused and full gates and record their measured durations.
