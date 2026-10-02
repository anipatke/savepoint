---
id: T-097
title: Decide how to bound the history a dashboard load reads
objective: O-032
status: done
depends_on: [{task: T-094, requires: clear}]
complexity_tier: spike
complexity_reason: Finding the newest snapshots without reading every file touches storage identity and migration, so the approach must be decided before any build.
owner_validation:
    required: true
    accepted_check: ""
planned_by: {role: planner, session: planning-o032-followup-t094}
check_waiver:
    task: T-097
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-02T04:51:56Z"
---

# Decide how to bound the history a dashboard load reads

## Outcome

A named decision says how opening Code Health can read a bounded window of saved snapshots instead of all of them, with the exact behaviour beyond the window, and proposes the bounded implementation Task or records that no change is warranted.

## User Check

Read the decision and its measured trade-offs; approve or reject the chosen approach and any change to wording or storage separately. Nothing is built by this Task.

## Done When

- T-094's owner-approved budgets are the target; the decision restates them and says whether T-094's measurements still hold on the current code.
- At least two ways to find the newest snapshots without opening every file are compared on correctness, storage identity, migration of existing projects, failure behaviour and measured cost, using the T-094 benchmarks. A cache is excluded unless the decision shows no other way meets the budgets.
- The decision states which earlier snapshots each signal needs (trend window 5, baseline 3, sparkline 10, per series) and whether the "N earlier official results not compared" and "manual results shown, not counted" wording may change when the window is bounded.
- The decision states whether sign-off wording (`hasOfficial`) and `SnapshotLabels` stay full-history.
- The decision states what happens to a damaged file outside the window: still fails the load, or is reported without blocking.
- A bounded implementation Task is drafted through design with measurable Done When items, or a justified no-change decision is recorded. No production code, retention change or storage identity change lands in this Task.

## Context Files

`.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-094-measure-the-cost-of-health-history.md` (measurements, budgets and hot paths); `internal/codehealth/storage.go`; `internal/codehealth/dashboard.go`; `internal/codehealth/history.go`; `internal/codehealth/snapshot.go`; `internal/codehealth/history_bench_test.go`.

## Design References

O-032 Confirmed Release-Hardening Design; Design sections on Code Health storage and history.

## Guardrails

FS-01, FS-04, FS-05, ARCH-01, ARCH-02, CFG-01, TEST-01, TEST-02, TEST-03, STYLE-01, STYLE-02.

## Implementation Plan

1. Re-run the T-094 benchmarks to confirm the baseline still matches.
2. Sketch each candidate lookup against the snapshot file layout and measure it with throwaway benchmarks that are not committed as product code.
3. Check each candidate against existing projects, damaged-history handling and snapshot identity.
4. Write the decision and, if warranted, the ID-free draft of the implementation Task for the planner.

## Boundaries

Research only; no speculative cache, retention change, storage identity change, provider timing promise or machine-dependent failing test.

## Technical Verification

`make build && make test-fast` for ordinary handoff; focused `make test-focused TEST=... PKGS=...` only during iteration. Record per-criterion evidence and limitations. The later independent Check applies `agent-skills/references/check-method.md`. Research completion requires the named decision deliverable; it does not claim any recommended repair has landed.

## Technical Evidence

Toolchain/platform: go1.26.2 linux/amd64 (WSL2), AMD Ryzen 7 7800X3D, 16 threads, warm filesystem cache. Candidate lookups were measured with a throwaway benchmark file (`internal/codehealth/zz_t097_throwaway_test.go`) built on the committed T-094 fixtures, run `-count=3`, then deleted. It is not product code and not in the tree; numbers below are medians, spread under 3%. Candidate A below is a sketch (substring scan of a 512-byte head), not the proposed implementation.

### Decision: bound-the-dashboard-history-window

**Target (T-094 budgets, restated).** Popover render constant (<= 700 allocs, 0.2 ms); load no worse than linear (<= 100 allocs/result, 15 allocs/KB); `single` at n=1,000 under 250 ms, with a 2x regression reopening the decision. T-094's measurements still hold on current code (re-run today): `LoadDashboard` single n=100 8.1 ms, n=1,000 85 ms, 62.7 MB, 377k allocs; heavy n=1,000 1.25 s, 1.21 GB, 4.15M allocs. All three budgets are met today.

**Finding not in T-094: the cost is paid more than once per board open.** `LoadChip` (`board/v2/load.go:165`, part of the index load) calls `LoadDashboard`; the Health screen's own load (`healthLoadCmd`, async) calls it again; `loadHealthLabels` (`load.go:173`) calls `SnapshotLabels`, a third full read, whenever some Check names a snapshot. `RefreshReport` and the collector also load. So 1,000 normal snapshots cost about 85 ms x 2-3 per open, and the chip load is in the project-load path, not an async command. Manual snapshots are never pruned automatically (`Prune` has no non-test caller), so history only grows.

**Candidates for finding the newest snapshots without decoding every file** (single / heavy at n=1,000; full `LoadSnapshots` is 76 / 1,150 ms):

| | Lookup alone | Correctness | Storage identity / migration | Damaged file |
|---|---|---|---|---|
| A. Read a 512-byte head per file for `created_at` and `origin` (both sit within the first 200 bytes of every canonical file: version, id, origin, retention, created_at), then fully decode only the window | 6.5 / 6.4 ms; with window decode 7.7 / 25 ms (3-4 ms and 19 ms of that is decoding about 15 files) | Same ordering as `LoadSnapshots` (stored time, then identity); file times never consulted. Cost is independent of snapshot size. A head that does not parse falls back to a full decode of that file | None: file names, bytes and retention unchanged; no new file | Unreadable head still fails; a bad body outside the window is not seen |
| B. Whole-file read plus partial JSON decode (what `LatestSnapshot` does today) | 32 / 540 ms | Correct | None | Same as A, but 540 ms heavy: no better than today for large files |
| C. File mtime | 2.0 ms | Wrong: git checkout, copy and restore reset mtimes; `LatestSnapshot` already states file times are never consulted | None | n/a |
| D. Time-ordered file names | 0.42 ms (names only) | Correct | Changes the "only name a snapshot may have is its identity digest" rule; every existing project needs a rename migration, and old and new binaries disagree | Names fail cheaply |
| E. Index or manifest file | about 0 | Correct only while it matches the directory | New file, rebuild and staleness rules: a cache, which T-094 excludes | n/a |

A is the only candidate that is exact, needs no migration and does not depend on snapshot size. C is rejected on correctness, D and E on identity, migration and the no-cache boundary. A still lists and opens every file, so load stays linear in file count at about 6.4 us per file (about 6.4 ms per 1,000 files, 0.64 s per 100,000), but no longer in bytes. The bytes were the dominant cost.

**Which earlier snapshots each signal needs** (`history.go`, `spark.go`, `dashboard.go`): per series (comparable official, measured, same `SeriesID`; manual and other-series entries only add to counts): trend window 5 values including current; decline baseline median of the last 3 earlier; sparkline 10 values including current. With the newest snapshot manual, current adds no value, so up to 10 earlier official snapshots are needed. The window is counted in snapshots, not series, because the series of a result is unknown until it is decoded.

**Chosen window:** the 10 newest official snapshots (`MaxSparkPoints`), every snapshot newer than the oldest of them, and the newest snapshot. This keeps every manual snapshot between officials, so manual counts inside the window stay exact. The newest 10 snapshots of either origin (the `History` list, `MaxDashboardHistory`) are always inside it, so that list is unchanged. Rows are identical to today's for any histories with at most 10 official snapshots, and for any signal present in every window snapshot beyond that. A signal absent from some window snapshots (an instance added or removed recently) can show fewer earlier values than full history would, and its trend can be thinner.

**Wording.** `basisWords` prints `len(s.values)`, which today is the whole history ("Compared with 47 earlier comparable official checks"), even though trends read at most 5 and sparklines 10. Beyond the window that number is no longer known, so the wording must change, only beyond the window. Proposed: when official history was cut, say "Compared with the 9 most recent comparable official checks" and append "Older history was not read." The "N earlier official results not compared" and "manual results shown, not counted" counts become counts within the window and say so. Histories inside the window keep today's text byte for byte. The wording change needs the owner's separate approval (see User Check).

**Sign-off and labels.**
- `hasOfficial`: stays exact full-history, at no extra cost, because candidate A reads `origin` from every file's head, not just the window.
- `SnapshotLabels`: stays full-history in the proposed Task and is not changed by it. It decodes every file to label every ID, yet a Check names a specific snapshot and the file name is that ID, so a lookup by ID is exact and needs no window. It is a separate, smaller follow-up; flagged here, not proposed as part of the window Task.
- The chip needs no change: it calls `LoadDashboard` and gets the bounded load.

**Damaged file outside the window.** The dashboard load fails on a name that is not a snapshot name (cheap, from the directory listing), on a head that cannot be read for time and origin, and on any damage in a window file (decoded, validated and checked against its name as today). It does not decode a body outside the window, so damage there is not shown on the dashboard. That damage is still caught where history is used in full: the collector (fails before any tool runs, `collect.go` `history()`), `Prune`, and doctor when a Check names a snapshot. Nothing is repaired or rewritten. A warning for it would need a new dashboard field and wording, so it is not proposed. This is a weakening of today's all-or-nothing dashboard and is the owner's call.

**Not proposed.** A cache or index, a storage-identity change, a retention change, mtime ordering, and any change to the collector, `Prune`, doctor or report text. No change at all is also defensible: all three T-094 budgets are met and 1,000 official or manual snapshots at ordinary size load in 85 ms. The case for change is the heavy end (1.25 s, 1.2 GB per load, two or three loads per open) and unbounded growth of manual history.

**Recommendation:** the bounded Task below, since the lookup is exact, needs no migration and removes the size-dependent cost; the owner may equally choose no change and keep T-094's budgets under watch.

### Per-criterion outcome
- T-094 budgets restated, measurements still hold: met (re-run `LoadDashboard` single n=100/1,000 and heavy n=1,000 within 2% of T-094; budgets all met today).
- At least two lookups compared on correctness, storage identity, migration, failure behaviour and measured cost, cache excluded unless needed: met (table: five candidates, A-E, measured with the throwaway benchmark; E excluded as a cache because A meets the budgets).
- Which earlier snapshots each signal needs, and whether wording may change: met (trend 5, baseline 3, sparkline 10, per series; `basisWords` full-history count must change beyond the window, inside it unchanged).
- `hasOfficial` and `SnapshotLabels`: met (`hasOfficial` exact at no cost; `SnapshotLabels` full-history, lookup by ID flagged as separate).
- Damaged file outside the window: met (not decoded, so not shown; header, name and window damage still fail; collector, Prune and doctor still validate in full).
- Bounded Task drafted or justified no-change recorded; no production code, retention or identity change: met (draft below; the throwaway benchmark is deleted, `git status` shows only T-094's two benchmark files and Savepoint records).

### Commands
- `go test ./internal/codehealth -run '^$' -bench 'HealthHistoryLoadDashboard/...' -benchmem -count=3` (baseline) and `-bench T097 -benchmem -count=3` (candidates, throwaway), all exit 0.
- `make build` exit 0; `make test-fast` exit 0 (no production code changed).

### Files
- Read: Context Files (`T-094` Task, `storage.go`, `dashboard.go`, `history.go`, `snapshot.go`, `history_bench_test.go` lines 1-400). Extra reads, to find who else loads history and to see the existing head-reading precedent: `AGENTS.md`, `agent-skills/savepoint-task/SKILL.md`, `.savepoint/router.md`, `internal/codehealth/{spark,chip,collect}.go`, `internal/board/v2/{load,io}.go`, `internal/doctor/health_snapshot_refs.go`, and the O-032 Objective header (status line only).
- Changed: this Task file only (status/stage, evidence). The throwaway benchmark was created and deleted in this session.

### Limitations
- One host, warm cache; cold disk, Windows/macOS and network drives unmeasured. Candidate A's per-file open cost is the part most likely to differ there.
- Candidate A was measured as a sketch with a substring head scan; a strict implementation (and its fallback for a head that does not parse) will cost somewhat more. The 512-byte head assumes canonical layout, so a hand-edited file with a late `created_at` must take the fallback.
- Whether the chip load runs before first paint was read from the code (`load.go` builds `HealthChip` in the index load), not timed in a running board.
- The window rule was reasoned from the code, not built; equivalence on windowed histories is claimed for the proposed Task to prove with tests.
- Owner approval of the decision, the wording change and the damaged-file behaviour is separate. No owner Task-check waiver recorded; no Check written.

## Drift Notes

Record any responsibility or interface change for planner reconciliation before the mandatory Full Objective Check. If evidence requires a material unknown repair, return REPLAN REQUIRED rather than inventing another scope.

No production interface changed; no benchmark file added by this Task (the throwaway one was deleted).

### Draft follow-up Task for the planner (proposal only; no Task record created, no ID assigned)

Title: Read a bounded window of health history for the dashboard

Outcome: Opening Code Health and the board chip read the newest 10 official snapshots plus the manual ones between them, instead of every saved snapshot, so load cost no longer grows with snapshot size, with unchanged output for histories inside the window.

Done when:
- A store function lists snapshot files, reads a bounded head of each for `created_at` and `origin` (falling back to a full decode when the head does not parse), orders them as `LoadSnapshots` does, and returns the window: the 10 newest official snapshots, every snapshot newer than the oldest of them, and the newest snapshot. Only those files are fully decoded and validated.
- `LoadDashboard` uses it; rows, headline, sign-off, `History` and `Chip` are byte-identical to today's for histories with at most 10 official snapshots, covered by tests built on the T-094 fixtures (compare against the full-load result).
- Beyond the window, specified and tested: trend, baseline and sparkline equal the full-history result for signals present in every window snapshot; `basisWords` says the window was applied ("most recent", "Older history was not read") only when history was cut; `hasOfficial` still reads every head.
- A name that is not a snapshot name, an unreadable head, or any damaged window file still fails the load with today's errors; a damaged body outside the window does not fail it; nothing is repaired or rewritten; `LoadSnapshots`, `Collect`, `Prune`, doctor and `SnapshotLabels` are unchanged and still read everything.
- `HealthHistoryLoadDashboard` single n=1,000 loads in under 250 ms and heavy n=1,000 allocates no more than heavy n=100 plus a stated constant, recorded the way T-094 recorded them.
- No cache, no index, no retention change, no storage identity change, no file written by a read.

Boundaries: `internal/codehealth` load path only; board and report layout do not change. `SnapshotLabels` lookup by ID is a separate follow-up.
Dependencies: this decision, owner-approved (window, wording, damaged-file behaviour).
Complexity: high (persistence semantics; a wrong window silently changes trends).

