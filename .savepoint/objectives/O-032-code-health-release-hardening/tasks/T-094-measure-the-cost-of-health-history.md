---
id: T-094
title: Measure the cost of health history
objective: O-032
status: done
depends_on: [{task: T-090, requires: clear}, {task: T-089, requires: clear}]
complexity_tier: spike
complexity_reason: Performance budgets and any necessary optimization must follow reproducible measurements.
owner_validation:
    required: true
    accepted_check: ""
planned_by: {role: planner, session: planning-o032-20261002-owner-confirmed}
check_waiver:
    task: T-094
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-02T04:43:49Z"
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

Toolchain/platform: go1.26.2 linux/amd64 (WSL2), AMD Ryzen 7 7800X3D, 16 threads. Benchmarks: `go test ./internal/codehealth -run '^$' -bench HealthHistory -benchmem -count=5` and `go test ./internal/board/v2 -run '^$' -bench HealthRender -benchmem -count=5`. Medians below; run-to-run spread 1-6% (heavy n=1000 load 6%), so repeatable on one host. Not measured: other machines, cold disk cache (fixtures are in the OS cache), Windows/macOS.

Fixtures (deterministic, built in the benchmark file; no provider, network or developer project): sizes 10/100/1,000 snapshots, every third manual, values drifting by index. Profiles: `single` 5 results x 1 evidence (5.3 KB/snapshot); `multi` 4 instances per signal = 20 results x 3 evidence (27.9 KB); `heavy` 20 results x 20 evidence (the MaxEvidence cap) with 200-byte notes (139 KB).

| Measurement | single | multi | heavy |
|---|---|---|---|
| LoadSnapshots n=10 / 100 / 1000 (ms) | 0.73 / 7.4 / 77.8 | 3.3 / 34.4 / 316.7 | 13.6 / 122 / 1,147 |
| LoadDashboard n=1000 (ms) | 83.8 | 368 | 1,198 |
| LoadDashboard n=1000 allocation (MB, allocs) | 62.7, 377k | 328, 2.05M | 1,214, 4.15M |
| In-memory projection n=1000 (ms) | 5.9 | 57.7 | 54.7 |
| RenderReport, any n (ms) | 0.007 | 0.043 | 0.20 |

- Load is the only cost that grows with history: linear in snapshot count and bytes (n x10 -> time x10 and allocs x10; about 78 us per 5 KB snapshot, 8 us per KB). Allocation is about 12x the bytes on disk. `LoadSnapshots` (read, JSON decode, validate every file) is 90-94% of `LoadDashboard`; the projection (`dashboardHistory` + `dashboardRows`) is 5-7%.
- Pure render does not grow with history: popover 0.085-0.096 ms (1, 4, 12 instances), history view 0.052 ms with 10 or 1,000 entries (identical allocs), overlay 0.13 ms at 80x24 and 0.18 ms at 200x50. Frame is always 17 lines (`TestHealthBenchFixturesFillTheFixedFrame`). `RenderReport` is flat across n and scales only with result/evidence size.
- Git freshness is independent of snapshot count (reads only the newest snapshot's identity): 2 us with scripted Git (60 allocs), 6.9 ms with real git against a 20-file throwaway repository (process launch dominates).

### Decision: performance-budget-and-remediation

Decision: no optimization for this release; record budgets and watch them. Rendering and freshness are already bounded. Load is linear, and ordinary projects stay far under any perceptible delay: 1,000 normal snapshots load in about 84 ms. Only an implausible history (1,000 snapshots of 20 results at the evidence cap, 139 MB on disk) costs about 1.2 s and 1.2 GB of allocation.

Proposed measurable budgets (benchmark expectations on the reference host above, reviewed at release; not unit-test failures, no universal machine-speed gate):
1. Popover render is constant: at most 700 allocs/op and 0.2 ms for every history size and up to 12 instances per signal (observed 542 allocs, 0.096 ms).
2. Load grows no worse than linearly: at most 100 allocs per result and 15 allocs per KB of snapshot on disk, summed over LoadDashboard (observed single 75/result, heavy 207/result; the heavy figure is evidence-driven, so judge by bytes: 12 allocated bytes per stored byte, observed 11.8-12.0).
3. Reference profile `single` at n=1000 loads in under 250 ms (observed 84 ms); a regression past 2x on repeated runs reopens this decision.

Exact hot paths: `Store.LoadSnapshots` (`internal/codehealth/storage.go`: `readRecord` + `DecodeSnapshot` + validation of every file, all retained in memory) and `dashboardHistory` (`dashboard.go`: groups every earlier result although trends read at most the newest 10 per series). `LoadDashboard` needs only the newest snapshot plus enough official history for a trend window.

Remediation, only if the owner wants it: propose a bounded implementation Task through design that loads the newest snapshot plus a bounded window of earlier ones. Snapshot file names are identity digests, so ordering needs `created_at`; any such Task must decide how to find the newest files without reading all of them, and must not change storage identity, add a cache, or change retention. Not proposed now.

### Per-criterion outcome
- Fixtures at 10/100/1,000 with named result/instance/evidence sizes; load vs render vs Git separated: met (`history_bench_test.go`, `health_bench_test.go`; table above).
- Wall time, allocations, fixture size, toolchain/platform, repeatability; no live dependency: met (5 runs per case, spread stated, `fixture-bytes` metric reported).
- Fixed-frame normal/history/multi-instance output with bounded snapshots, asymptotic growth and observed limits: met (render constant, load linear, limits above). No machine-speed gate added.
- Named decision with proposed budgets and exact hot paths; no cache or identity change: met (above). The bounded optimization Task is proposed, not created.

### Commands
- Benchmarks as above, all exit 0. The first codehealth 5-run and a board 5-run overlapped in time once; the board numbers above come from a clean re-run afterwards.
- `make build && make test-fast`: exit 0.

### Files
- Read: both dependency Task files' status/waiver headers only, `internal/codehealth/{dashboard,dashboard_freshness,spark,report}.go`, `dashboard_test.go`, `dashboard_freshness_test.go`, `internal/board/v2/health_view.go`, `health_realcopy_test.go`, `health_popover_test.go`. Extra reads (outside Context Files, for fixture building and exact hot paths): `internal/codehealth/{history,storage,snapshot,classification,model,repository}.go`, `storage_test.go`, `model_test.go`, `internal/board/v2/{health,health_test,model,releases}.go`, `Makefile`, `AGENTS.md`, `agent-skills/savepoint-task/SKILL.md`, router, Objective.
- Changed: `internal/codehealth/history_bench_test.go` (new), `internal/board/v2/health_bench_test.go` (new), this Task file. `fixture_test.go` was not needed.

### Limitations
- One host, warm filesystem cache; Windows/macOS and slow disks unmeasured. Whether the board runs the load off the UI loop was not verified here.
- Real-Git freshness used a 20-file repository; large repositories cost more in Git itself and were not measured.
- Budgets are proposals; owner approval is separate. No owner Task-check waiver recorded; no Check written.

## Drift Notes

Record implemented responsibility or interface changes for planner reconciliation before the mandatory Full Objective Check. If evidence requires a material unknown repair, return REPLAN REQUIRED rather than inventing another scope.

No production interface changed; only two benchmark test files were added.

### Draft follow-up Task for the planner (proposal only; no Task record created, no ID assigned)

Title: Bound the history a dashboard load reads

Outcome: Opening Code Health reads a bounded window of saved snapshots instead of all of them, so load cost stops growing with history, with identical dashboard output.

Why: T-094 measured `LoadDashboard` as linear in snapshot count and bytes (84 ms at 1,000 normal snapshots, 1.2 s and 1.2 GB at 1,000 maximum-size ones), 90-94% of it in `Store.LoadSnapshots`. Official snapshots are never pruned, so history only grows. Render and Git freshness are already constant.

Open design questions (the planner decides before any build):
- Snapshot file names are identity digests, so the newest files cannot be found without reading them. Options: read only the small header fields needed to order files first, or add an ordering to names/an index. Either touches storage identity and migration, which T-094 forbade, so this needs an explicit design decision. A cache is not proposed.
- Which earlier snapshots a row needs: the newest N official and manual results per series (trend window is 5, baseline 3, sparkline 10), plus the counts of manual and incompatible results that `basisWords` states ("N earlier official results not compared"). Bounding the window would change those counts unless they are kept some other way. Decide whether that wording may change.
- `hasOfficial` (sign-off wording) and `SnapshotLabels` also read every snapshot; decide whether they stay full-history.

Done when (sketch):
- Dashboard output is byte-identical to today's for histories within the window, covered by tests built on the existing fixtures.
- Behaviour beyond the window is specified and tested, including the basis counts and sign-off text.
- `HealthHistoryLoadDashboard` at n=1000 (`single`) loads in under 250 ms and allocates no more than at n=100 plus a stated constant; the numbers are recorded the way T-094 recorded them.
- Damaged-history behaviour is unchanged: an invalid file still fails the load, nothing is repaired or rewritten (state whether files outside the window are still validated).
- No retention change, no storage identity change, no cache.

Boundaries: `internal/codehealth` load path only; the board and report text do not change.
Dependencies: T-094 owner-approved budgets.
Complexity: high (persistence semantics and migration safety).

Not a release blocker on T-094's evidence; schedule it only if the owner wants the headroom.

