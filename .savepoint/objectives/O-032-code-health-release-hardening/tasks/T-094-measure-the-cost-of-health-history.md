---
id: T-094
title: Measure the cost of health history
objective: O-032
status: planned
depends_on: [{task: T-090, requires: clear}, {task: T-089, requires: clear}]
complexity_tier: spike
complexity_reason: Performance budgets and any necessary optimization must follow reproducible measurements.
owner_validation: {required: true}
planned_by: {role: planner, session: planning-o032-20261002-owner-confirmed}
---

# Measure the cost of health history

## Outcome

Reproducible benchmark evidence and a performance decision identify whether history loading or rendering needs a bounded release fix.

## User Check

Review measured history sizes and timings; approve any proposed budget or optimization separately.

## Done When

- Create deterministic benchmark fixtures at 10, 100 and 1,000 snapshots with named result/instance/evidence sizes; distinguish filesystem dashboard load from pure report/popover render and Git freshness.
- Record wall time, allocations, fixture size, toolchain/platform and repeatability; no live provider, network or developer-project dependency.
- Measure fixed-frame normal/history/multi-instance output with bounded snapshots; state asymptotic growth and observed limits without inventing a universal machine-speed gate.
- Deliver a named performance-budget-and-remediation decision with proposed measurable release budgets and exact hot paths. If an optimization is warranted, propose a bounded implementation Task through design; do not prematurely add a cache or change storage identity.

## Context Files

`.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-089-simplify-board-reload-state-restoration.md` (dependency evidence); `.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-090-protect-saved-health-evidence-under-failures.md` (dependency evidence); `internal/codehealth/dashboard.go`; `internal/codehealth/dashboard_test.go`; `internal/codehealth/dashboard_freshness.go`; `internal/codehealth/spark.go`; `internal/codehealth/report.go`; `internal/codehealth/report_test.go`; `internal/board/v2/health_view.go`; `internal/board/v2/health_realcopy_test.go`; `internal/board/v2/health_popover_test.go`; `internal/board/v2/fixture_test.go`.

## Design References

O-032 Confirmed Release-Hardening Design and Owner-Requested Health Analysis; Design sections 1, 6, 7, 9, 12 and 13 as applicable.

## Guardrails

FS-01, FS-04, FS-05, FS-06, ARCH-01, ARCH-02, ARCH-03, CFG-01, CFG-02, TEST-01, TEST-02, TEST-03, TEST-04, TEST-05, TEST-08, TEST-09, STYLE-01, STYLE-02, STYLE-07.

## Implementation Plan

1. Establish deterministic history and dashboard fixture builders using existing test helpers.
2. Benchmark load and pure rendering separately across recorded sizes.
3. Compare repeated measurements and inspect dominant costs.
4. Record proposed budget, justified no-change decision or exact optimization follow-up.

## Boundaries

Research/benchmarks only; no speculative cache, provider timing promise, retention change or machine-dependent failing unit test.

## Technical Verification

`make build && make test-fast` for ordinary handoff; focused `make test-focused TEST=... PKGS=...` only during iteration. Record named per-criterion evidence and limitations. The later independent Check applies `agent-skills/references/check-method.md`; the mandatory Full Objective Check runs its fresh full gate and official Code Health collection. This executor never writes a Check or runs savepoint health setup/report/check. Research completion requires the named decision deliverable; it does not claim any recommended repair has landed.

## Technical Evidence

Pending execution: named per-criterion cases, command/time/toolchain/results, files read/changed, decision deliverable where applicable and limitations.

## Drift Notes

Record implemented responsibility or interface changes for planner reconciliation before the mandatory Full Objective Check. If evidence requires a material unknown repair, return REPLAN REQUIRED rather than inventing another scope.
