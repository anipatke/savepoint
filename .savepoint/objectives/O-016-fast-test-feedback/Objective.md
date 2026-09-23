---
id: O-016
title: Make test feedback fast without weakening verification
status: done
release: R-006
last_check: C-909
freshness:
  state: current
  check: C-909
  assessed_by: {role: checker, session: o016-full-recheck-20260923}
  assessed_at: '2026-09-23T04:20:00Z'
  basis: C-909 ran a fresh passing Linux make test-full with native npm and the focused Windows migration check on the current working tree (HEAD e7e72bb plus diff sha256 a1e3a570); later edits are Check metadata only.
---

# O-016: Make test feedback fast without weakening verification

## Outcome

Developers and agents receive fast, trustworthy feedback for ordinary changes,
while migration recovery, cross-platform behavior, and full integration remain
covered by an explicit complete gate at the boundaries where that depth is
needed.

## Why

The current `make test` path is dominated by repeated full migration scenarios
and regularly takes about two minutes. Requiring that same cycle after every
small code or metadata correction creates long idle periods and encourages
people to skip validation altogether. I-023 records the observed problem and
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

## Confirmed Verification Policy

- Ordinary Task handoff requires a successful build and fast repository gate.
- Migration and platform-sensitive changes also require a fresh full gate at
  Task handoff. CI and every Full Objective or Release Check require the full
  gate. The full gate includes host-platform tests and cross-build checks; CI
  additionally runs a focused Windows migration setup check on Windows.
- Focused package tests support iteration but do not replace the fast handoff
  gate.
- The two-second focused, fifteen-second fast, and forty-five-second full
  targets are measured goals on a recorded warm machine. Timing regressions
  are reported with attribution, not encoded as flaky test timeouts.
- After a metadata-only correction, an unchanged full-gate result may remain
  current only if code, tests, fixtures, and gate definitions are unchanged.
  Evidence names the original run, the correction, and the unchanged inputs.
  Any relevant input change requires a fresh full run.

## Design Confirmation

The owner confirmed the verification policy, Windows coverage, and T-013-T-016
implementation sequence on 2026-09-22 at 21:51 UTC. The confirmed plan begins
with measurement and gate selection; implementation decisions that measurement
cannot settle return to design rather than silently changing these boundaries.

The owner approved the C-908 remediation design on 2026-09-23 at 02:11 UTC:
keep T-013-T-016 done and add separate O-016 Tasks for the migration package
TestMain conflict, Windows watcher exclusion, and the guide-casing regression
assertion. Preserve the existing Windows coverage and verification policy;
record the full Windows suite after all three repairs. The casing report will
be verified through directory-entry spelling because Windows Stat cannot
distinguish paths that differ only by case.

The owner narrowed the Windows verification scope on 2026-09-23 at 03:05 UTC:
keep a focused Windows migration setup check and the existing cross-builds;
the full Windows runtime suite and unrelated Windows platform repairs are not
required for O-016. The Linux `make test-full` gate remains required. This
supersedes the 02:11 instruction to record a full Windows suite after the
repairs.

The owner confirmed on 2026-09-23 at 03:55 UTC that T-018 and T-019 are
not needed without the full Windows runtime suite. Remove both from
O-016; keep the focused T-017 Windows check and the Linux full gate.

The owner closed O-016 on 2026-09-23 at 04:24 UTC on the current CLEAR Full
Objective Check C-909, after T-013-T-017 were done. I-030 stays open and deferred.

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

- Removing migration interruption, recovery, idempotency, or filesystem-safety
  coverage merely to meet a timing target. Keep the focused Windows migration
  setup check and cross-builds; the broad Windows runtime suite is deferred.
- Replacing deterministic tests with production mocks that bypass the behavior
  under test.
- General runtime performance work unrelated to test execution.
- Treating Go's test cache as proof that an unexecuted required gate passed on
  changed inputs.

## Originating Issue

I-023 — Shorten the repository test feedback cycle.
