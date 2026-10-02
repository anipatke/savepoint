---
id: T-093
title: Exercise all health readers together
objective: O-032
status: planned
depends_on: [{task: T-090, requires: clear}, {task: T-091, requires: clear}, {task: T-086, requires: clear}]
complexity_tier: high
complexity_reason: Nine readers and partial polyglot states need cross-instance integration evidence.
owner_validation: {required: false}
planned_by: {role: planner, session: planning-o032-20261002-owner-confirmed}
---

# Exercise all health readers together

## Outcome

Local cross-stack fixtures prove all approved readers, history and Check integration preserve measured, incomplete and unavailable states independently.

## User Check

Read a mixed-instance fixture result: one failed tool does not erase another signal or claim healthy evidence.

## Done When

- Cover all nine registered report readers with named Go/JS/Python mixed-instance scenarios; good, partial, failed, missing, malformed and unsupported evidence never collapses into a healthy zero.
- Cross scope/config/provider-version changes, renamed/sibling instances and newest manual/official snapshots with history/trend comparability and staleness expectations.
- Verify strict snapshot references and official-only Full Check support, required versus optional instance blockers and report failure warning; no new classification policy.
- Apply duplication research denominator findings to fixtures; do not assert scoped totals from a global report that lacks the required counts.
- Record an explicit provider/scenario coverage matrix with exact existing or new tests; normal gate is offline and temporary-project based.

## Context Files

`.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-086-measure-duplication-in-maintained-code.md` (dependency evidence); `.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-090-protect-saved-health-evidence-under-failures.md` (dependency evidence); `.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-091-verify-provider-processes-stop-safely.md` (dependency evidence); `internal/codehealth/readers.go`; `internal/codehealth/readers_integration_test.go`; `internal/codehealth/reader_repair_test.go`; `internal/codehealth/instances_test.go`; `internal/codehealth/classification_test.go`; `internal/codehealth/history.go`; `internal/codehealth/gate_test.go`; `internal/codehealth/dashboard_signoff_test.go`; `internal/healthcheck/healthcheck.go`; `internal/healthcheck/healthcheck_test.go`; `internal/doctor/health_snapshot_refs_test.go`; `internal/codehealth/testdata/readers/duplication/jscpd5.json`; `internal/codehealth/testdata/readers/vulnerabilities/scanner-error.json`.

## Design References

O-032 Confirmed Release-Hardening Design and Owner-Requested Health Analysis; Design sections 1, 6, 7, 9, 12 and 13 as applicable.

## Guardrails

FS-01, FS-04, FS-05, FS-06, ARCH-01, ARCH-02, ARCH-03, CFG-01, CFG-02, TEST-01, TEST-02, TEST-03, TEST-04, TEST-05, TEST-08, TEST-09, STYLE-01, STYLE-02, STYLE-07.

## Implementation Plan

1. Expand the existing DefaultReaders/Collect integration matrix only where outcome evidence is missing.
2. Add controlled fixtures and fake tools for polyglot/partial and history transitions.
3. Verify command/doctor/gate interpretation agrees on official evidence and independent failures.
4. Run focused readers/integration and handoff gate.

## Boundaries

No new provider, real provider/network requirement in tests, official collection by executor, or implicit clearance from labels.

## Technical Verification

`make build && make test-fast` for ordinary handoff; focused `make test-focused TEST=... PKGS=...` only during iteration. Record named per-criterion evidence and limitations. The later independent Check applies `agent-skills/references/check-method.md`; the mandatory Full Objective Check runs its fresh full gate and official Code Health collection. This executor never writes a Check or runs savepoint health setup/report/check.

## Technical Evidence

Pending execution: named per-criterion cases, command/time/toolchain/results, files read/changed, decision deliverable where applicable and limitations.

## Drift Notes

Record implemented responsibility or interface changes for planner reconciliation before the mandatory Full Objective Check. If evidence requires a material unknown repair, return REPLAN REQUIRED rather than inventing another scope.
