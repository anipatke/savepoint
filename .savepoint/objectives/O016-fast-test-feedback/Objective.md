---
id: O016
title: Make test feedback fast without weakening verification
status: planned
---

# O016: Make test feedback fast without weakening verification

## Outcome

Developers and agents receive fast, trustworthy feedback for ordinary changes,
while migration recovery, cross-platform behavior, and full integration remain
covered by an explicit complete gate at the boundaries where that depth is
needed.

## Why

The current `make test` path is dominated by repeated full migration scenarios
and regularly takes about two minutes. Requiring that same cycle after every
small code or metadata correction creates long idle periods and encourages
people to skip validation altogether. I023 records the observed problem and
the evidence behind this Objective.

## Success Conditions

- Focused tests for a changed package normally complete within two seconds on
  a warm development machine.
- A documented fast repository gate completes within fifteen seconds and is
  sufficient for ordinary implementation feedback and metadata-only changes.
- A documented full gate runs every required migration interruption, recovery,
  integration, and platform-sensitive scenario, with a target duration below
  forty-five seconds on the same machine.
- Migration tests stop repeating equivalent fixture copies and full applies
  when several read-only assertions can share one independently prepared
  result without weakening isolation or failure diagnosis.
- Safe test groups run concurrently where they do not share mutable hooks,
  filesystem state, environment state, or other process-wide controls.
- Verification policy states exactly when focused, fast, and full evidence is
  required, including when an unchanged full-gate result may remain current
  after a metadata-only correction.
- Make targets, CI, Guardrails, active skills, scaffold copies, and repository
  guidance agree on the same gate names and responsibilities.
- CI reports package/test timing so a regression in the fast or full budget is
  visible and attributable rather than discovered through agent wait time.

## Architectural Considerations

- Coverage ownership remains with the existing package tests; splitting gates
  changes scheduling and evidence requirements, not the behavior promised by
  migration or other packages.
- Tests using package-global fault-injection hooks remain serial until those
  hooks become injected per-run dependencies. Parallelism must never introduce
  nondeterminism into recovery evidence.
- Shared fixture preparation is allowed only within a test owner that can keep
  every mutating scenario isolated. No test may depend on another test's order
  or previously mutated directory.
- Guardrail TEST-08 and the Check/task skills must be reconciled atomically with
  any new fast/full distinction; no renamed target may silently drop the slow
  suite from CI, Objective Checks, or Release Checks.
- Timing budgets are measured guidance with recorded environments, not flaky
  assertions embedded in ordinary unit tests.

## Boundaries

**In scope:**

- Measurement and attribution of current test duration.
- Fast/full Make targets and CI routing.
- Migration fixture reuse, duplicated end-to-end scenario consolidation, and
  safe test parallelism.
- Agent evidence-reuse rules for metadata-only changes.
- Guardrail, skill, template, and active guidance reconciliation.
- Timing reporting and regression visibility.

**Out of scope:**

- Removing migration interruption, recovery, idempotency, filesystem-safety,
  Windows, or cross-platform coverage merely to meet a timing target.
- Replacing deterministic tests with production mocks that bypass the behavior
  under test.
- General runtime performance work unrelated to test execution.
- Treating Go's test cache as proof that an unexecuted required gate passed on
  changed inputs.

## Originating Issue

I023 — Shorten the repository test feedback cycle.
