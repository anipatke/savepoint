---
id: T-093
title: Exercise all health readers together
objective: O-032
status: done
depends_on: [{task: T-090, requires: clear}, {task: T-091, requires: clear}, {task: T-086, requires: clear}]
complexity_tier: high
complexity_reason: Nine readers and partial polyglot states need cross-instance integration evidence.
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: planning-o032-20261002-owner-confirmed}
check_waiver:
    task: T-093
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-02T04:26:58Z"
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

Run 2026-10-02, go1.26.2, linux/WSL2. `make build && make test-fast` passed (codehealth 10.2s). Focused: `go test ./internal/codehealth -run Polyglot -count=3` passed. No provider, network, official collection, or `savepoint health` command ran; the only new file is `internal/codehealth/readers_polyglot_test.go` (temporary projects, fake tool runner, existing fixtures).

### Per-criterion evidence

1. **Nine readers, Go/JS/Python mix.** `TestPolyglotReadersKeepEachStackIndependent` runs `DefaultReaders()` over one config with an instance per reader. Outcomes: go-test measured (failing), vitest-junit malformed=failed, pytest-junit measured (failing), go-cover measured but stale, vitest-v8 absent (required), coveragepy malformed=failed, lizard go measured / js unavailable / py partial, jscpd js and py measured, osv scanner-error=partial. Every unmeasured result has no value, details or evidence, every partial is never `good`, and every non-measured summary is `unknown`. Each instance run alone in a fresh project matches its result beside the others (outcome, value, details).
2. **Scope/config/provider-version, rename/sibling, manual/official, trend and staleness.** `TestPolyglotHistoryComparabilityAcrossChanges`: three official checks draw a trend with the early note; a rename keeps the series; a scope change restarts only that instance (sibling keeps its drawing); a jscpd provider-version change restarts only jscpd; a newest manual snapshot leaves the trend and shows "manual refresh does not affect sign-off"; a newer official snapshot restores sign-off. Staleness: the Go coverage profile is older than the Go sources and is judged `stale` (reported, optional).
3. **Strict references, official-only Check, required vs optional, report warning.** Same test: `Evaluate` refuses the manual snapshot (`ErrManualSnapshot`); the absent vitest-v8 report blocks while required and is only reported with `Required` off (same snapshot, config digest unchanged). `TestPolyglotDashboardAgreesWithTheVerdict`: dashboard sign-off and every row's blocks/advisory state agree with `Evaluate`, and neither text ever says "healthy". Existing and not duplicated: strict snapshot references `internal/doctor/health_snapshot_refs_test.go` (missing, manual, unreadable, read-only); report-failure warning `TestCollectReportWriteFailureKeepsTheSnapshot` and `TestRun_reportWriteFailureIsOnlyAWarning`. No classification policy changed.
4. **Duplication denominator.** The scoped jscpd instances use the jscpd 5 report (no per-file counts); the test asserts each keeps only the report's own totals (3.33% of 600 lines) and asserts no scoped figure.
5. **Coverage matrix.** In the header comment of `readers_polyglot_test.go`; each reader's good/partial/malformed/boundary cases stay in `reader_*_test.go`, and `TestCollectWithRealReadersKeepsEveryMeasureTruthful` stays as the all-measured case.

### Files

Read: the Context Files plus `collect.go`, `gate.go`, `dashboard.go`, `spark.go`, `identity.go`, `reader_lizard.go`, `reader_osv.go`, `reader_jscpd.go`, `dashboard_test.go`, `repository_test.go`, `reader_repair_test.go`, `collect_test.go` (extra reads: to learn outcome rules, series identity and helpers). Changed: `internal/codehealth/readers_polyglot_test.go` (new), this Task file.

### Limitations

Test-only: no production behavior changed. Reports are made stale/fresh by file mtime, so the stale case depends on filesystem timestamps (set explicitly, hour offsets). The OSV scanner-error fixture reads as partial with a value, not failed; I asserted that existing behavior rather than changing it.

## Drift Notes

Record implemented responsibility or interface changes for planner reconciliation before the mandatory Full Objective Check.

None from this Task. Observation only: the `JscpdReader.Read` doc comment says a report without per-file counts "gives a partial reading with no value", but the code (and `TestJscpdReaderWithoutPerFileCountsKeepsTheReportTotal`) keeps the report total, available. Left unchanged as out of scope. If evidence requires a material unknown repair, return REPLAN REQUIRED rather than inventing another scope.
