---
id: C-960
scope: {kind: objective, id: O-032}
result: CLEAR
checked_by: {role: checker, session: recheck-o032-20261003-independent}
executed_session: repair-o032-20261002b
checked_at: '2026-10-02T19:20:00Z'
health_snapshot: sha256:b8e4f20a27cc62d79498d3647002804246bc2e89d6e762e0d8c60d631f43111a
reviewed:
  base_commit: d00543d5252fc5557aeba067995427236da2823f
  head_commit: 37dffd3dbb72adcda32376b0ce89a916c96238e6
  files:
    - .savepoint/Design.md
    - .savepoint/issues/I-118-snapshot-header-shortcut-can-hide-the-newest-valid-snapshot.md
    - .savepoint/objectives/O-032-code-health-release-hardening/tasks/T-097-decide-how-to-bound-the-history-a-dashboard-load-reads.md
    - .savepoint/objectives/O-032-code-health-release-hardening/tasks/T-098-read-a-bounded-window-of-health-history-for-the-dashboard.md
    - internal/codehealth/dashboard.go
    - internal/codehealth/runner.go
    - internal/codehealth/window.go
    - internal/codehealth/window_test.go
  dependencies: []
issues: []
supersedes: C-959
---

# C-960: O-032 Final Recheck After the Owner's I-118 Decision

## Result and authority

CLEAR. The code reviewed by C-959 is unchanged. The owner chose C-959's option 1: keep the whole-file read. The planner then amended T-098 (Outcome, Done When 1, 3, 4 and 6, and User Check), the T-097 decision and Design.md line 30 to describe it (commit 37dffd3; I-118 history `owner_decision` at 2026-10-02T07:05:00Z). Measured against those amended criteria, I-118 is proven. I-115, I-116, I-117, I-119, I-120 and I-121 keep their proven assessments, and every cell was reproduced again in this run.

This conversation started after `/clear` and performed no implementation, so it is independent from repair-o032-20261002b. The owner asked for this recheck after making the decision that C-959's convergence handoff required. It changes no implementation, acceptance criterion, Design, router or Task status. C-957, C-958 and C-959 are immutable.

## Closure map

| Issue | Assessment | Disposition |
|---|---|---|
| I-115 | Current-head native evidence: run https://github.com/anipatke/savepoint/actions/runs/36975366853 on exact head 37dffd3 has `ci` (110737921697) and `windows-tests` (110737921938) both successful. | Closed `verified` |
| I-116 | Its code has not changed since C-958's proof. The native Windows suite is green on the current head. | Closed `verified` |
| I-117 | The 36-cell official/manual count matrix (official 0–11 × leading/trailing/interleaved) passes 36 of 36. | Closed `verified` |
| I-118 | Proven against the amended criteria (see below). | Closed `verified` |
| I-119 | README and CHANGELOG are unchanged since C-958 proved them. Design.md now agrees with the code and README ("not checked", whole-file read, measured cost). | Closed `verified` |
| I-120 | Its code has not changed since C-958's proof. The stream tests pass in the fresh full gate. | Closed `verified` |
| I-121 | 130 of 130 cells pass through the real ExecRunner and Collect, and the saved snapshot is checked too. The sizes are 64KiB −1, exact and +12, every overflow that cuts the URI, and 2× and 3× the cap. Each runs with newline and non-space filler. No credential fragment remains, the trailing failure text is kept, the reason is valid UTF-8 and it is bounded. | Closed `verified` |

## Frozen scope and admission ledger

C-957's scope lock and the C-958/C-959 ledgers apply unchanged. Since C-959's head d2791fa, `git diff --name-only d2791fa HEAD` touches only `.savepoint/` records: Design.md, C-959, one snapshot, I-118, I-121, T-097 and T-098. No code, test, fixture, dependency or gate definition changed. The only new input is the owner-amended wording of the criteria C-959 measured against. It is not a new axis.

| Item | Prior claim | Frozen cell | Result |
|---|---|---|---|
| Amended DW1: each file read for `created_at`/`origin` with the full decoder's last-wins rule | I-118 | W valid duplicate headers, newest/oldest/boundary | Pass: 8 valid representation cells (repeat before) and 6 same-value repeats after |
| Different-value repeat after the normal fields | I-118 adjacency | W damage | Same 4 cells as C-959. The full decoder rejects them, so they are damage, not representation. The window tolerates them only when the changed time moves the file outside the window, which amended DW4 allows |
| Amended DW4: outside the window, non-JSON or missing/invalid time or origin fails; other readable damage is tolerated; inside the window, all damage fails; nothing is written | I-118 (a) | W damage and no-write | Pass, 17 of 17. Truncated, trailing-garbage, bad-time, missing-time and bad-origin fail outside and inside. Renamed-field, unknown-field and bad-results-value pass outside and fail inside. A non-snapshot name fails. `LoadWindow` and `LoadDashboard` agree, and the directory is byte-identical after every load |
| Representations: compact, reordered, 512-byte prefix, wrong type, trailing, symlink | I-118 adjacency | W representation | Pass |
| Amended DW6: single n=1,000 under 250 ms; heavy cases measured and recorded | I-118 (b) | F benchmarks | Pass (see below) |
| DW3 and User Check wording "Older history was not checked" | I-118 | W basis | Pass: `dashboard.go:164`, `window_test.go:237-240` |
| I-121 credential matrix | I-121 | E secret-safe stderr overflow | Pass, 130 of 130 |
| Fresh full gate, checksums, current-head native CI, official health after gate | Release gate | R/P | Pass |

The temporary harness `internal/codehealth/o032_c960_probe_test.go` was C-959's appendix harness with the damage test rewritten to the amended DW4. It was run, then deleted before the benchmarks. The damage test's essential change: `truncated`, `trailing-garbage`, `bad-time`, `missing-time` and `bad-origin` are expected to fail outside the window. `renamed-field`, `unknown-field` and `bad-results-value` are expected to load outside the window. All of them are expected to fail inside the window, and a `notes.json` name is expected to fail.

## Measured cost

`go test ./internal/codehealth -run '^$' -bench HealthHistoryLoadDashboard -benchmem -count=1` on go1.26.2 linux/amd64, Ryzen 7 7800X3D, WSL2, warm cache, single run:

| LoadDashboard | C-959 | C-960 |
|---|---|---|
| single n=1,000 | 34.0 ms, 8.1 MB | 38.3 ms, 8.1 MB |
| multi n=1,000 | 129 ms, 35.4 MB | 136 ms, 35.4 MB |
| heavy n=100 | 76.7 ms, 34.8 MB | 274.8 ms, 34.7 MB (b.N=10, noisy) |
| heavy n=1,000 | 561.6 ms, 168.0 MB | 613.6 ms, 168.0 MB |

Single n=1,000 stays inside DW6's 250 ms budget. Heavy cost grows with snapshot bytes, which matches the amended Outcome, DW6, T-097's amendment and Design.md ("about 0.56 s at 1,000 maximum-size snapshots"). Memory matches C-959 exactly. The heavy n=100 time is a single noisy sample: its allocations are unchanged, and DW6 sets no budget for it.

## Gate, platform and health evidence

- Fresh `make ci` at head 37dffd3 on 2026-10-03 (go1.26.2 linux/amd64): exit 0, with no `--- FAIL` lines. It covered the uncached report-writing test run, six target builds, dist, the npm wrapper and the package dry-run. `sha256sum -c checksums.txt` passes all six archives. The harness was added only after the gate and the health check had finished.
- `git diff --check d2791fa HEAD`, excluding Check records, is clean.
- Native CI: run 36975366853 on exact head 37dffd3, with `ci` and `windows-tests` both successful. I read the run metadata only and dispatched nothing.
- After the gate, the official `./savepoint health check O-032` exited 0 and created the snapshot named in the frontmatter. Code Health does not block clearance. All five instances are optional. The OSV result again reports two unknown-severity groups (golang.org/x/sys GO-2026-5024, golang.org/x/text GO-2026-5970) that need review. This is advisory, unchanged since C-958, and not an Issue.

## Design reconciliation

Design.md line 30 now matches `window.go`, README, CHANGELOG and AGENTS.md: the whole-file read, window-only validation, the non-JSON failure outside the window, "not checked" wording and the cost growth. T-097 records the owner amendment beside the original decision. T-098 keeps its original head-read Technical Evidence and labels it as history. No drift remains.

## Materiality

No Issues remain, so no materiality actions are required.

## Observations (non-blocking)

- `TestLoadDashboardReadsNoSnapshotBodyOutsideTheWindow` proves that bodies outside the window are not validated, not that they are not read. Its name overstates what it checks. This was carried from C-959.
- `window_test.go:216` asserts that the full-history basis does not contain `"not read"`. The current wording is "not checked", so this assertion can no longer fail. The adjacent `HasSuffix(textOlderNotRead)` check still guards the window side. The constant name `textOlderNotRead` also predates the wording change.
- The `LoadWindow` doc comment still says "left unread" for `cut`. That is accurate only in the sense that those files are not decoded.
- Untracked `.coverage.out.*` and `.go-test.json.*` files were at the repo root before this run and were left alone.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches** — `TestLoadWindowStillFailsOnNamesAndUnreadableHeads` covers the non-JSON and bad-time branches of `readHead`. Truncated and trailing-garbage files outside the window are proven only by this Check's harness.
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth** — the head read uses the decoder's own JSON rule.
- [ ] STYLE-08 **Comments explain why** — the `LoadWindow` comment says "left unread" where the code now reads but does not decode (`window.go:30`).
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs**

Style is advisory and did not affect the verdict.

## Owner validation and closure

Every owned Task (T-085 to T-098) is `done`. This Objective Check is CLEAR and current, and every material Issue linked to the Check chain is now resolved `verified`. O-032 is ready for owner acceptance and closure. This Check does not set the Objective's status.
