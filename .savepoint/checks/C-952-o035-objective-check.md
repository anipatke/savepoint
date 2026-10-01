---
id: C-952
scope: {kind: objective, id: O-035}
result: NEEDS WORK
checked_by: {role: checker, session: check-o035-20261002-independent}
executed_session: t075-t081-executor-sessions-not-recorded
checked_at: '2026-10-01T22:30:12Z'
health_snapshot: sha256:003ca03d393d927dd5e8e1aa4e069c49eb1a1cc38cc65d299d4cac1e73b48508
reviewed:
  head_commit: 42c8d32dd74e3a6340985b4bb20763315ca03ba7
  files:
    - internal/codehealth/dashboard.go
    - internal/codehealth/dashboard_test.go
    - internal/codehealth/dashboard_copy.go
    - internal/codehealth/dashboard_copy_test.go
    - internal/codehealth/dashboard_signoff_test.go
    - internal/codehealth/chip.go
    - internal/codehealth/chip_test.go
    - internal/codehealth/spark.go
    - internal/codehealth/spark_test.go
    - internal/board/v2/health.go
    - internal/board/v2/health_view.go
    - internal/board/v2/health_test.go
    - internal/board/v2/health_popover_test.go
    - internal/board/v2/health_history_test.go
    - internal/board/v2/view.go
    - internal/board/v2/update.go
    - internal/board/v2/help.go
    - internal/board/v2/load.go
    - internal/board/v2/detail_test.go
  dependencies:
    - internal/codehealth/storage.go
    - internal/codehealth/classification.go
    - internal/codehealth/history.go
    - internal/codehealth/gate.go
    - internal/codehealth/dashboard_freshness.go
    - internal/board/v2/io.go
    - internal/styles/styles.go
    - go.mod
    - go.sum
    - Makefile
issues: [I-107, I-108, I-109, I-110, I-111, I-112]
supersedes: null
---

# C-952: O-035 Full Objective Check

## Verdict

NEEDS WORK. All seven owned Tasks T-075–T-081 are done and carry explicit owner Task-check waivers. This conversation did not implement them and is independent from their executor sessions; those sessions have no precise identities in their logs. Owner terminal validation for T-076/T-079/T-081 remains separate and pending. No completion, acceptance, router change or implementation repair is recorded here.

The full gate and gate-backed sign-off parity pass. The display does not yet meet the Objective: normal values remove yardsticks, spark notes disappear, the chip can select a different snapshot, explanations misstate thin-history readings, progress exposes providers, and multiple instances hide a capability.

## Frozen Scope Lock

1. Requirements: all eleven Objective success conditions and T-075..T-081 Done When; ARCH-02/03/04, DATA-03, TEST-01/02/04/05/08/09. Advisory STYLE-01..10. No classifier, threshold, retention or collection changes are permitted. Owner terminal acceptance remains pending.
2. Surfaces: LoadChip, Store.LatestSnapshot, Dashboard.Chip, LoadDashboard/dashboardRows, copy helpers, spark helpers, board loadProject/header/View, Health reducer/open/close/key dispatch/help/hints. Files are the seven Tasks' named source/test files plus directly relied-on storage/history/gate/classifier and board IO/geometry. Includes existing working-tree O-031 repairs only as dependencies.
3. Workflow: saved config -> snapshot listing/decode -> rows/copy/series/gate/history -> explicit board load -> render; H -> load -> freshness -> select/history -> R -> progress -> completion/reload or cancellation -> Esc restore / quit. Rendering does no IO. Collection persistence is unchanged; cancellation must preserve saved results.
4. Matrix (every row crosses applicable columns):
| Row | normal/boundary | missing/malformed/failure | state/sequence | output/representations |
| Copy five signals | four labels; Good/Watch thresholds below/at/above; default/override; 0/1/many paths | absent/unavailable/partial/stale | thin history versus level wording | dashboard -> rendered real-length strings |
| Gate wording | failing tests/high/unrated/opt-in | required missing/failed/stale; invalid Evaluate config | official/manual/no official | same disposition as Evaluate; displayed wording |
| Spark | 0..12 points; spread below/at/above material size; rising/falling/flat | partial/manual/incompatible/nonfinite excluded | restarted and manual newest | eight blocks; early/restart/thin notes through renderer |
| Chip | three states/four labels; 5 signals/multiple instances | missing/corrupt storage | differing save order/timestamps/ties; load versus dashboard update | 50..100 header widths; normal/ANSI/no colour |
| Popover | 80x20/24/40; all selections; history 0/1/10/11+ | loading/load error/unknown/partial | H/select/h/Esc; removed keys; cursor restore; refresh/progress/cancel/retry/help quit | real copy/yardsticks/spark notes; ASCII/wide/combining/emoji/ANSI truncation |
5. Material boundary: supported saved config/snapshots and interactive board at documented sizes; admit reproducible lost required information, false explanations, sign-off disagreement, read/write or gate failure. No unrelated dependency internals. Non-TTY has no interactive popover (N/A); network, redirects, authentication, billing N/A: no such boundary. External refresh runner remains unchanged and is tested by existing collection/board regression tests for success, failure, timeout/cancel, cleanup and secret-safe failure. Immutable record persistence for this review is separately strict-loaded.

Workflow side-effect inventory:
| Order | Operation | Effect | Failure/final state | Cleanup/oracle |
| 1 | Config/snapshot read | none | diagnostic (dashboard), quiet chip fallback per Task | independently ordered stored CreatedAt; file bytes |
| 2 | pure projection/gate | none | unavailable sign-off on Evaluate error | numeric thresholds, gate dispositions, expected copy |
| 3 | board H/load/freshness | overlay state only | last result retained on load error | injected counted funcs, stale snapshot ID suppression |
| 4 | render/select/history | cursor/overlay state only | clipping must retain required content | terminal cell width and literal semantic content |
| 5 | R collect/progress | explicit runner/manual save only on success | cancel/failure preserves old result | existing real-process cancellation tests and fake save counters |
| 6 | completion/reload/close/quit | chip/result update or cancellation | no retry during active refresh | command trace, previous cursor and saved bytes |

## Completed Coverage Matrix

The row names and axes above were fixed before the independent probes. No recheck or supersession is involved. Cells below classify all applicable axes; N/A reasons are part of the lock. Existing deterministic parameterized tests supplement independent scenarios, not the other way around.

| Frozen row | Normal and boundaries | Failure and malformed | Sequence/state | Final output | Acceptance |
|---|---|---|---|---|---|
| Copy | Default/configured Good/Watch and all signals/labels pass the copy tests; 0/1/many paths pass | missing/unmeasured and not configured pass; confidence-driven Watch misstates measured result | fresh first official 90% reproduces I-110 | real values reach row correctly; fit loses aims I-107 | Issue |
| Gate | normal, failing tests, unknown/high severity, opt-in match Evaluate | required failed/stale/absent and invalid config handling pass | official/manual/no official pass | overall and each row agree with gate | Proven |
| Spark | 0..12 and five capability movement thresholds below/at/above pass independent matrix; rise/fall tested | partial/manual/incompatible excluded by existing row tests; NaN/Inf excluded by independent scenario | restart series selection passes | early/restart notes lost I-108; long direction clipped I-107 | Issue |
| Chip | 3 states/4 colours, 50..100 widths and load/update tests pass | absent/broken storage fallback per explicit Task passes | mtime versus CreatedAt reversed: I-109; equal mtime uses ID instead of measurement ordering by source | header and popover disagree in reproduction | Issue |
| Popover | 80x20/24/40, every selection, long explanation, history 0/1/10/14 pass; extra instances I-112 | loading/load error/unknown/partial pass | open/select/h/Esc/removed keys/cursor restore pass; refresh order/cancel/retry/quit/help regressions pass | ASCII, combining marks, wide characters, variation selectors, modifiers, flags, joined emoji, ANSI width pass independent test; real required copy I-107, provider progress I-111 | Issue |

Input duplicate/wrong/mixed types and mutation after validation: stored snapshots/configs are validated by existing storage/model tests; these changes introduce no new parser or write format. Multiple distinct configured instances are a supported direct model and are exercised explicitly (I-112). Mutable dashboard projections use copied slices where stored data is transformed; no change to serialization. Nonfinite comparisons are independently rejected. Lifecycle backward/terminal revival is N/A for a view with no persisted lifecycle transitions; retry, overlapping refresh guard, late freshness message rejection, history return and cancellation are tested by the named reducer regressions. Control bytes use the existing ANSI-aware truncator; this work defines no new control-byte acceptance guarantee. Non-TTY, dumb/no-cursor and redirected sinks have no interactive popover; health view colour semantics use glyphs as well as styles. No unsupported direct invalid cursor state is admitted as a blocker.

## Workflow And External Boundary Evidence

Saved data is read without tools by LoadChip/LoadDashboard; rendering uses only model strings and state. Config/snapshot read errors are visible in the dashboard and retain last result; chip fallback is explicitly required by T-076. Evaluate errors retain dashboard data and say unavailable. H dispatch, explicit freshness, and snapshot-ID filtering match ARCH-02/03/04. Refresh uses the unchanged manual Collect boundary, with configured project root and DefaultReaders, serialized progress, context cancellation and one completion message. Supported failures, refusal/unavailable tools, malformed reports, timeout/cancellation and cleanup are covered by the full collection/reader/runner suite; remote requests/redirects are N/A because the renderer has no network/client boundary. No collection, storage-write or runner remediation was attempted.

Board regression evidence: TestHealthRefreshReportsProgressInOrderAndReloadsOnCompletion; TestHealthEscCancelsARefreshAndKeepsThePreviousResult; TestHealthRefreshErrorIsAStatusLine; TestHealthQuitDuringRefreshCancelsBeforeQuitting; TestHealthCtrlCFromHelpDuringRefreshCancelsBeforeQuitting; TestHealthStaleFreshnessForAnOlderSnapshotIsIgnored; TestHealthLoadErrorKeepsThePreviousView; TestHealthEscRestoresTheBoardCursor; TestHeaderChipRefreshesAfterRefreshResult; TestHealthHistoryReturnPaths. These assert semantic result/cursor state and counted explicit command calls. Pure renders do not enter IO. Cancellation before saving and report preservation also passed the fresh full suite, including the existing I-105/I-106 dependencies already reviewed in C-951.

## Acceptance Coverage

| Objective success condition (in listed order) | Classification | Evidence |
|---|---|---|
| 1 Header chip, saved reads, one row | Issue | I-109; other chip states/layout pass |
| 2 Fixed popover, Esc restore | Proven | documented-size tests and independent real-copy dimensions |
| 3 Names/questions in one place | Proven | questions/copy tables; accessibility with multiple instances is I-112 |
| 4 Value, mark and aim on each row | Issue | I-107, I-112 |
| 5 Accurate meaning, sign-off, next step | Issue | I-110; gate parity passes |
| 6 Official spark and history qualification | Issue | I-108; data selection/scaling passes |
| 7 Copyable Objective command | Proven | selected/sidebar, router and placeholder tests |
| 8 Ten newest checks, manual excluded/dim | Proven | dashboard history and popover tests; source uses CardMeta for manual |
| 9 No providers/hashes/raw times | Issue | I-111; measured/history omission otherwise passes |
| 10 Refresh/cancel/stale/first-run | Proven | full reducer/state tests; progress provider leak remains condition 9 |
| 11 Tests cover stated behaviours | Issue | existing tests pass but use short synthetic copy and omit rendered notes/multi-instance/cause fidelity |

T-075 data fields/formatting/aim/Where are present, but Meaning correctness has I-110. T-076 loading/header fallback/update work, but newest selection has I-109. T-077 sign-off parity and UTC dates are proven. T-078 series/filter/scaling and strings are proven; integration loses notes. T-079 layout/state/keys work, but integrated content has I-107/I-108/I-112. T-080 states/cancel work; provider omission has I-111. T-081 history bounds/order/back paths are proven. Task waivers satisfy TEST-09 and do not replace this Full Check.

Design reconciliation: the saved-data projection, pure renderer, explicit refresh, existing classifiers/thresholds and storage identities remain within Design and the Codebase Map. New files extend those existing responsibilities; no new package, provider, storage schema or dependency was added. The specified level-based label rule is not a license to attribute every Watch label to a bad measured number: I-110 asks for truthful explanations without modifying classification. No Design edit is made by this checker.

## Commands And Health

- `make test-full`: exit 0, fresh run at approximately 2026-10-01 22:23–22:24 UTC, Go go1.26.2 linux/amd64. All packages passed with uncached reports; Linux/Darwin/Windows build targets passed. Native Windows testing is CI policy; this run does not claim new native Windows evidence.
- `make build`: exit 0; current source binary used for official collection.
- `go test ./internal/codehealth ./internal/board/v2 -count=1 -run 'Test(Dashboard|LoadChip|LatestSnapshot|Spark|Health|HeaderChip)'`: exit 0, both packages pass.
- `go test -overlay /tmp/o035-probes/overlay.json ./internal/board/v2 ./internal/codehealth -run TestCheckO035 -v`: expected exit 1, six independently reproduced acceptance failures. Independent Unicode and finite spark matrix pass. GateCases initially used an invalid blocking flag for Tests; corrected harness limits opt-in to Coverage. Rerunning GateCases passes required failed/stale/absent, opt-in, high/unrated and partial cases. That harness mistake is not a product finding.
- `git diff --check`: exit 0.
- `./savepoint health check O-035`: exit 0, official snapshot recorded above, verdict **does not block clearance**. Tests/coverage/complexity/duplication report no blocking findings. Optional osv-scanner collection failed: missing measurement, not a security finding. No automatic Issue is created for it.

The official collection is supporting evidence; it does not erase reproducible UI findings. No code/tests/fixtures/dependencies/gates changed during this review. Independent Go probe files were external `/tmp` scratch supplied only by an overlay; no test source was installed in the repo. Scope/harness/log scratch may be discarded; durable reproductions and line evidence are in I-107–I-112. Every scoped evidence source exists. Historical T-075 partial scratch is explicitly reconciled by its replan/handoff and current files.

## Issues And Materiality

| Issue | Likelihood | Impact | Materiality | Recommendation |
|---|---|---|---|---|
| I-107 Lost yardsticks/direction | High: normal values | High for this Objective: removes the scale | High | Fix layout with real copy |
| I-108 Lost history notes | High: any 3/4-point or restarted spark | Medium: misleading confidence/continuity | Medium | Combine with I-107 renderer repair |
| I-109 Different latest snapshot | Low: restored/copied/out-of-order files | Medium: inconsistent health at a glance | Medium | Fix saved recency semantics |
| I-110 False meaning | High: normal first checks | Medium: tells owner to repair a good measurement | Medium | Fix cause-aware plain copy |
| I-111 Provider progress leak | High: every refresh | Low: jargon, no secret or privacy harm | Low | Narrow progress copy repair |
| I-112 Hidden fifth capability | Medium: supported multiple instances | High: may hide the blocking signal | High | Preserve five-capability visibility |

All six remain Issues because they violate explicit scoped requirements; low materiality does not waive I-111. Existing Issues were searched by symptom/location/linked work before allocating IDs. I-104 concerned scrolling of the removed full-screen detail/history and is already verified; this popover's fixed-row loss and new yardstick/note requirements have distinct reproductions. None of I-107–I-112 repairs implementation here. Default remediation is direct repair under the Issues, preserving every done Task; later independent recheck uses this frozen lock.

## Owner Validation And Observations

Real-terminal visual acceptance for chip, popover and history is still needed after repair. Existing headless style tests often run without an ANSI colour profile; source style wiring is reviewed, but no actual terminal appearance is claimed. Optional dependency scan collection failure is an observation, not bad code or a blocker. No additional out-of-scope blocker is asserted.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches** — semantic integration gaps are separately recorded above.
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [ ] STYLE-07 **One source of truth** — sparkValues duplicates buildTrend value assembly; latest lookup has a second recency policy (chip.go:97, history.go:83).
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs** — changes remain within existing dashboard and board responsibilities.

## Reviewed Working-Tree Identity

The source is uncommitted atop the recorded head. These hashes identify reviewed bytes rather than treating that head as the implementation revision.

- `internal/codehealth/dashboard.go`: `cd4f8db10445fb2736e9fddc874453dcdabb41f976c2fb57e026f7fcc391a3f2`
- `internal/codehealth/dashboard_test.go`: `5e70b56d15ebe12e29a466aeed43204ba83a7c7a92e1a8e07f60618d00153e80`
- `internal/codehealth/dashboard_copy.go`: `d885d4525a9ff7143a192ec79bc00712c22faeddd4c72dbabc06d75f724d66ff`
- `internal/codehealth/dashboard_copy_test.go`: `3c44dd9d839743ff8c95b844f84275d6e79c6f2848bbc0fc9bb439c12a8ee43c`
- `internal/codehealth/dashboard_signoff_test.go`: `b25f3a8aa8356d1ed3db47661f3c612698374425acf3a665b47324500f5d3155`
- `internal/codehealth/chip.go`: `f2945022bc431caed066db3343a09bb562fefb1981baa335d619dc3a6e32c655`
- `internal/codehealth/chip_test.go`: `d85b07ec5ea5dc689d5b18e890c2996f8a8ad01506641366cf6581adec3de13c`
- `internal/codehealth/spark.go`: `83485582282409df37cb16fba31692ce608133d1fb8e5d2c887a2309f7d13213`
- `internal/codehealth/spark_test.go`: `0cb59506e9dc807d251033e57652963531c28bf597df2ea4bf834ee5a78a4227`
- `internal/board/v2/health.go`: `8f346f7040919ce091a0687d5a812e064482c28886f0f267330284b58f32ab8b`
- `internal/board/v2/health_view.go`: `dae9fdc884e426f673df6884950240185538a95376a34770467885cb24b4eecd`
- `internal/board/v2/health_test.go`: `ad54caa0f0f8daaf12459b2a9e68acd132d16c50db41c27bcbbf0cd5ebee0fb8`
- `internal/board/v2/health_popover_test.go`: `4804f8c2f0f596288d9e39b945647dea0de708c4d28b6384fe55e9222335558c`
- `internal/board/v2/health_history_test.go`: `f31a3d873bb50d7f928ceab872c1bbaaa15ce0402b9bae496b2cc4e492c54aad`
- `internal/board/v2/view.go`: `12e0433f0db1777ab19db9ed9db1d7da6a4f5a5e10aac62d1670c8be0082d0a7`
- `internal/board/v2/update.go`: `fe7d030ac6b58e37bc61f7e4bcb117e6031187ec6bc6591ddfccd4ba358d516d`
- `internal/board/v2/help.go`: `94044d488a53cf637e97ef8c6f3f7ce255f4e53d027fc204658f891511457f46`
- `internal/board/v2/load.go`: `d1b0dcd40ed654821fc1019b9c10c04a99a0f1e22be00896103191b6f6984ad6`
- `internal/board/v2/detail_test.go`: `c87b70ad6a0e0a0d756c14a59a4754448ba412e69a49de1ff8d620a5a6ca148c`
