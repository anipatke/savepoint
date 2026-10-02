---
id: T-098
title: Read a bounded window of health history for the dashboard
objective: O-032
status: done
depends_on: [{task: T-097, requires: clear}]
complexity_tier: high
complexity_reason: Persistence semantics; a wrong window silently changes trends.
owner_validation:
    required: true
    accepted_check: ""
planned_by: {role: planner, session: planning-o032-t097-window}
check_waiver:
    task: T-098
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-02T04:56:32Z"
---

# Read a bounded window of health history for the dashboard

## Outcome

Opening Code Health and the board chip read the newest 10 official snapshots plus the manual ones between them, instead of every saved snapshot, so load cost no longer grows with snapshot size, with unchanged output for histories inside the window.

## User Check

Open the board with a small and a large saved history; the Health screen and chip look the same, and a long history says "Older history was not read" in the basis line.

## Done When

- A store function lists snapshot files, reads a bounded head (512 bytes) of each for `created_at` and `origin` (falling back to a full decode of that file when the head does not parse), orders them as `LoadSnapshots` does (stored time, then identity), and returns the window: the 10 newest official snapshots, every snapshot newer than the oldest of them, and the newest snapshot. Only those files are fully decoded, validated and checked against their name.
- `LoadDashboard` uses it. Rows, headline, sign-off, `History` and `Chip` are byte-identical to the full-load result for histories with at most 10 official snapshots, proven by tests built on the T-094 fixtures.
- Beyond the window, specified and tested: trend, baseline and sparkline equal the full-history result for signals present in every window snapshot; `basisWords` says "most recent" and "Older history was not read" only when history was cut, and its manual and not-compared counts are within the window; `hasOfficial` still uses every head.
- A name that is not a snapshot name, an unreadable head, or any damaged window file still fails the load with today's errors; a damaged body outside the window does not fail it; nothing is repaired or rewritten.
- `LoadSnapshots`, `Collect`, `Prune`, doctor and `SnapshotLabels` are unchanged and still read everything.
- `HealthHistoryLoadDashboard` single n=1,000 loads in under 250 ms and heavy n=1,000 allocates no more than heavy n=100 plus a stated constant, recorded the way T-094 recorded them.
- No cache, no index, no retention change, no storage identity change, no file written by a read.

## Context Files

`.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-097-decide-how-to-bound-the-history-a-dashboard-load-reads.md` (decision, window rule, wording, damaged-file behaviour); `internal/codehealth/storage.go`; `internal/codehealth/dashboard.go`; `internal/codehealth/history.go`; `internal/codehealth/spark.go`; `internal/codehealth/chip.go`; `internal/codehealth/history_bench_test.go`.

## Design References

O-032 Confirmed Release-Hardening Design; T-097 decision `bound-the-dashboard-history-window`, owner-approved.

## Guardrails

FS-01, FS-04, FS-05, ARCH-01, ARCH-02, CFG-01, TEST-01, TEST-02, TEST-03, STYLE-01, STYLE-02.

## Implementation Plan

1. Add the head reader and window selection to the store, with tests for ordering, ties, fallback and damaged names and heads.
2. Switch `LoadDashboard` to the window; keep `hasOfficial` on the heads.
3. Change `basisWords` only when history was cut; add tests for inside and beyond the window.
4. Re-run the T-094 benchmarks and record the numbers.

## Boundaries

`internal/codehealth` load path only; board and report layout do not change. `SnapshotLabels` lookup by ID is a separate follow-up.

## Technical Verification

`make build && make test-fast` for ordinary handoff; focused `make test-focused TEST=... PKGS=...` only during iteration. Record per-criterion evidence and limitations. The later independent Check applies `agent-skills/references/check-method.md`.

## Technical Evidence

Toolchain/platform: go1.26.2 linux/amd64 (WSL2), Ryzen 7 7800X3D, warm cache. Benchmark: `go test ./internal/codehealth -run '^$' -bench HealthHistoryLoadDashboard -benchmem -count=3` (T-094's fixtures, medians, spread under 3%).

Built: `internal/codehealth/window.go` (`Store.LoadWindow`, 512-byte head read with whole-file fallback, `windowStart`); `LoadDashboard` now calls it through a new `buildDashboard`; `series.cut` and two wording constants make `basisWords` name the window only when history was cut. `LoadSnapshots`, `Collect`, `Prune`, doctor and `SnapshotLabels` are unchanged.

### Per-criterion outcome
- Store function lists files, reads a bounded head for time and origin, falls back to a full decode, orders like `LoadSnapshots`, returns the window, decodes only the window: met. `TestLoadWindowKeepsTheNewestOfficialAndEverythingNewer` (sizes 0, 1, 5, 15, 16, 40 against the full load), `TestLoadWindowBreaksTimeTiesByIdentityLikeLoadSnapshots`, `TestLoadWindowDecodesAFileWhoseHeadIsInAnotherOrder`, `TestLoadWindowWithoutHealthStorageIsEmpty`.
- Dashboard byte-identical within the window: met. `TestLoadDashboardIsUnchangedWithinTheWindow` compares the whole `Dashboard` with `reflect.DeepEqual` against the full-history build for sizes 1, 2, 3, 8 and 15 (15 = exactly 10 official plus manual). The chip is derived from the dashboard, so it follows.
- Beyond the window: met. `TestLoadDashboardBeyondTheWindowKeepsTrendsAndSaysOlderWasNotRead` (sizes 16, 30, 100): every row equals full history except `Basis`, which ends "Older history was not read." and says "most recent"; headline, sign-off and the history list equal full history. `TestBasisWordsNamesTheWindowOnlyWhenHistoryWasCut` fixes the exact wording both ways. `hasOfficial` runs on the window, which always holds an official snapshot when one exists; `TestLoadDashboardWithOnlyManualHistoryStillSaysNoOfficial` keeps the no-official text. This is how it stays equal to full history without reading every head separately.
- Failures: met. `TestLoadWindowStillFailsOnNamesAndUnreadableHeads` (a non-snapshot name, a head that is not JSON, a bad time, a directory in place of a file, each outside the window) fails the load; `TestLoadWindowDamageInsideAndOutsideTheWindow` shows a damaged body outside the window loads, the same damage inside it fails, and `LoadSnapshots` still rejects the file; `TestLoadDashboardReadsNoSnapshotBodyOutsideTheWindow`. Nothing is written or repaired.
- Other readers unchanged and still reading everything: met by diff (no edit to `LoadSnapshots`, `Collect`, `Prune`, doctor, `SnapshotLabels`); existing tests for them pass.
- Budget numbers: met, recorded below.
- No cache, index, retention or identity change, no write on read: met by diff.

| LoadDashboard | before (T-094, re-run) | now |
|---|---|---|
| single n=100 / n=1,000 | 8.1 / 85 ms | 2.4 / 13.2 ms |
| single n=1,000 memory, allocs | 62.7 MB, 377k | 4.9 MB, 91k |
| multi n=1,000 | 368 ms | 18 ms |
| heavy n=100 / n=1,000 | 122 ms / 1.25 s | 22.8 / 34.2 ms |
| heavy n=1,000 memory, allocs | 1.21 GB, 4.15M | 23.7 MB, 149k |

Single n=1,000 is 13 ms against the 250 ms budget. Heavy n=1,000 allocates 149k against 72k at n=100: the stated constant is about 85 allocations and 3.5 KB per file beyond the window (head read, token decode, sort), measured as the identical 76.5k difference for single and heavy; it does not depend on snapshot size. Load is still linear in file count (about 13 us per file), no longer in bytes.

### Commands
- `go vet ./internal/codehealth`; `go test ./internal/codehealth -run 'LoadWindow|LoadDashboard|BasisWords'`: pass. Benchmark above: exit 0.
- `make build && make test-fast`: exit 0.

### Files
- Read: Context Files, plus `internal/codehealth/dashboard_test.go` (first 120 lines, fixtures), `storage_test.go` helpers (grep), `primitives.go` (timestamp check, grep), `AGENTS.md`, savepoint-task skill, router.
- Changed: `internal/codehealth/window.go` (new), `window_test.go` (new), `dashboard.go`, `history.go` (one field), `history_bench_test.go` (one call argument), `AGENTS.md` (codebase map sentence), this Task file, router selection.

### Limitations
- One host, warm cache; cold disk, Windows, macOS and network drives unmeasured; per-file open cost is what would change there.
- Head parsing trusts the head: a file with a valid head and damaged body outside the window is not seen by the dashboard (owner-approved). A file with a duplicate key in its head is read as its first value.
- Per-file cost is higher than the T-094 sketch (13 ms vs 7.7 ms at n=1,000) because it uses a strict token decode; not tuned.
- `SnapshotLabels` still reads every file; a lookup by ID is the separate follow-up. No board screen was run; the board tests pass unchanged.
- No Check written; no waiver recorded.

## Drift Notes

Record any responsibility or interface change for planner reconciliation before the mandatory Full Objective Check.

`LoadDashboard` now calls the new `Store.LoadWindow` and a private `buildDashboard`; the dashboard's `Basis` text names the window when history was cut. `Design.md` sections on Code Health history may want one sentence: the dashboard reads a bounded window; collection, prune and doctor read everything.
