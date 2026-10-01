---
id: C-953
scope: {kind: objective, id: O-035}
result: NEEDS WORK
checked_by: {role: checker, session: recheck-o035-20261002-independent}
executed_session: user-request-i107-i112-repair
checked_at: '2026-10-01T22:48:04Z'
health_snapshot: sha256:2370ba359ca13976f423dc36cf89feab4c93bc261d8e76194c85c806ea49a15b
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
    - internal/board/v2/health_realcopy_test.go
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
issues: [I-107, I-112]
supersedes: C-952
---

# C-953: O-035 First Full Recheck

## Closure Map

| Prior Issue | Closure state | Recheck result |
|---|---|---|
| I-107 | Still open | Original lost-aim reproduction passes; repair widens frame beyond the frozen maximum |
| I-108 | Still open, repair proven | Early/restart/thin notes present in selected detail |
| I-109 | Still open, repair proven | Persisted recency matches dashboard despite mtime changes |
| I-110 | Still open, repair proven | Thin/confidence/baseline explanations name their cause |
| I-111 | Still open, repair proven | Provider-free progress on every signal and both modes |
| I-112 | Still open | Instances reachable, but repair scrolls the five-row window |

Result remains NEEDS WORK for the two existing layout Issues. I-108–I-111 are technically reproduced as repaired, but not resolved as verified: the shared issue-capture reference defines verified as proof from a Check recording CLEAR. The later successful Objective Check can resolve them. This session never implemented the reviewed work; it remains independent from the repair executor. All Tasks stay done with their existing owner waivers. No owner acceptance or exception is inferred.

## Scope And Admission Ledger

C-952's numbered scope lock, matrix axes, external-boundary exclusions, workflow inventory, supported-path/materiality boundary and all acceptance requirements remain frozen and are incorporated by reference. This is the first full recheck, not a new initial Check. An admission ledger was written in /tmp before the first recheck probe:

| Item | Prior claim | Exact frozen cell | Admitted outcome |
|---|---|---|---|
| Real copy/aims/frame | I-107 | Popover normal 80x20/24/40, real copy final output; lock item 1 includes T-079 <=72 | Original loss repaired; existing width requirement fails |
| Notes with/without blocks | I-108 | Spark final early/restart/thin; Popover all selections | Pass |
| Reversed saves/mtime/ties | I-109 | Chip differing save order/timestamps/ties and load vs dashboard | Pass |
| Good 1/2/3 checks, partial/stale, decline | I-110 | Copy thin history vs level wording, failure/confidence, final copy | Pass |
| Progress all signals/modes | I-111 | Popover refresh/progress output; no-provider requirement | Pass |
| Multiple first/middle instances/final blocker | I-112 | Chip multiple instances; Popover all selections; lock item 1 includes never-scroll | Reachability repaired, no-scroll requirement fails |
| Remaining original cases | C-952 matrix | Every original Completed Coverage Matrix row | Pass as detailed below |

No new axis, threshold interpretation, dependency layer or public configuration was added. The dimension and no-scroll requirements were explicit in original T-079 and Objective records, included in C-952 lock item 1, and applied to the same supported terminal sizes and duplicate-instance reproduction. They are regression checks of the repair, not newly invented design requirements. Neither governing record has been amended and no owner exception is recorded.

## Findings And Exact Reproductions

**I-107 remains open.** The same realistic five values and default aims now appear, and better/worse direction is complete. However `renderHealthPopover(80)` is 76 cells wide at 80x20, 80x24 and 80x40. T-079 requires at most 72. `healthPopoverMaxWidth = 100` at health_view.go:17 and min(100, termWidth-4) at :105 directly cause this. Independent TestCheckO035FrozenDimensions asserts actual rendered width with Lipgloss; all three fail with width=76. The revised real-copy test checks only fitting within the terminal, so does not enforce the original cap. Restore the cap while retaining required scales, or have the owner explicitly change the design outside this Check.

**I-112 remains open.** Prepend a second Tests instance to the original five rows; press down five times. Before: Tests, Tests, Coverage, Complexity, Duplication. After: Tests, Coverage, Complexity, Duplication, Dependencies, with an up cue. The final signal is accessible, but an earlier instance disappears and the rows shift. `healthWindowTop` at health_view.go:279 and `moveHealth` at health.go:167 implement that movement. O-035 condition 2 says the popover never scrolls; T-079 says no key scrolls. Independent TestCheckO035MultipleInstances checks reachable final signal and unchanged no-scroll contract and fails only the latter. The repair test proves reachability, not conformity to that requirement. A bounded representation must preserve accessible instance information without scrolling unless the owner explicitly revises this design.

These two are not new Issues; dated recheck evidence was appended to their existing records. The remaining four Issues have passing reproductions and evidence appended without manufacturing a verified resolution.

## Completed Frozen Matrix And Acceptance Coverage

| C-952 row | Normal/boundary | Missing/malformed/failure | State/sequence | Output |
|---|---|---|---|---|
| Copy | All signal/label and default/override Good/Watch tests pass; 0/1/many paths pass | absent/unmeasured/not configured and confidence cases pass | independent Tests/Coverage/Complexity/Duplication 1/2/3 observations, coverage partial/stale/decline pass | original real-value aims pass; cause-aware copy passes |
| Gate | independent failing tests/high/unrated/opt-in and default cases pass | required failed/stale/absent and invalid config tests pass | official/manual/no-official tests pass | same per-row/overall Evaluate disposition |
| Spark | independent 0..12 × five capabilities × below/at/above material thresholds pass; rising/falling tests pass | manual/partial/incompatible exclusions and NaN/Inf probes pass | restarted/manual newest tests pass | selected early/restarted/thin now visible with and without blocks |
| Chip | states/colours/counts/header widths pass | missing/damaged storage tests pass | independent original reversed-mtime and regression equal-mtime/stored-time-tie cases pass | header/dashboard identical in reproduction |
| Popover | all selections and 80x20/24/40 height fit; history 0/1/10/14 pass; width cap fails I-107 | loading/load error/unknown/partial pass | H/selection/h/Esc/removed keys/cursor restore and refresh/help cancellation pass; duplicate-instance scrolling fails I-112 | real copy, notes, provider omission and Unicode/ANSI width probes pass |

All nonapplicable cells retain C-952's reasons; all workflow operations retain its inventory. Rendering still does no IO; explicit load/freshness/refresh commands and cancellation paths are unchanged except view selection movement. Current full tests exercise real collection/runner cancellation and report preservation dependencies. No new network/subprocess behavior, configuration format, storage identity, snapshot retention, classification or threshold rule is introduced. Persisted recency lookup reads each bounded record for creation time but fully decodes only the chosen newest; unchanged damaged-older-file and missing-storage regressions pass.

Acceptance conditions 1, 3, 5, 6, 7, 8, 9 and 10 are proven by the corrected reproductions plus original focused tests. Conditions 2 and 11 still have in-scope layout regressions; condition 4's aims are repaired, while the multi-instance representation still conflicts with the no-scroll contract in condition 2. All seven Task outcomes were retained from C-952 and reevaluated through the same criteria: T-075/076 repairs pass; T-077/078/080/081 remain proven; T-079 layout cannot clear. No acceptance criterion was rewritten to fit the repair.

## Gates And Code Health

- Fresh `make test-full`: exit 0, Go go1.26.2 linux/amd64; uncached full suite with reports and Linux/Darwin/Windows build targets pass. This is host full evidence, not a new claim of native Windows execution.
- `make build`: exit 0, rebuild current source CLI.
- `go test ./internal/codehealth ./internal/board/v2 -count=1 -run 'Test(Dashboard|LoadChip|LatestSnapshot|Spark|Health|HeaderChip|Meaning)'`: exit 0 in both packages.
- Independent external Go overlay, `go test -overlay /tmp/o035-probes/overlay.json ./internal/board/v2 ./internal/codehealth -run TestCheckO035 -v`: exit 1 for precisely two admitted layout failures. RealCopy, SparkNotes, ProgressProvider, Unicode, ChipRecency, Meaning, IndependentSparkMatrix, GateCases and CauseAdjacent pass. Helper-specific old assertions were adapted to inspect final selected-detail notes and the shortened Dependencies name; the semantic requirements stayed fixed. No test source was installed in the repository. One attempt reported a read-only Go build cache; approved escalation retried it, and subsequent complete overlay execution produced the stated evidence.
- `git diff --check`: exit 0.
- `./savepoint health check O-035`: exit 0; official snapshot above; **Code Health does not block clearance**. Optional osv-scanner collection failed, with no measurement. This is not evidence of vulnerable code and does not block. Other signals have no blocking findings.
- Final `./savepoint resume` strict-loads C-953 and updated Issue histories. No router, Design, implementation, Task status or Objective status write was made.

The full gate and health result do not waive either layout requirement. Source remains uncommitted atop the recorded head. Only dashboard.go/copy.go/copy_test.go, chip.go/chip_test.go, board health.go/health_view.go/health_test.go changed among C-952's hashed files; health_realcopy_test.go is a new repair test. Other original scoped source/test hashes match C-952. All referenced source exists; overlay/ledger/output in /tmp is declared disposable scratch, with durable reproductions here and in Issues.

## Materiality

| Issue | Likelihood | Impact | Materiality | Recommendation |
|---|---|---|---|---|
| I-107 width cap | High at documented sizes | Low: still fits the terminal, but breaks explicit compact width | Low | Narrow layout correction or explicit owner design decision |
| I-112 scrolling | Medium: supported multiple instances | Medium: contradicts fixed non-scrolling view design | Medium | Bounded representation or explicit owner design decision |

Neither is data loss, secret exposure or a platform failure. These remain acceptance Issues even though the original higher-impact missing-information symptoms are repaired. No additional out-of-scope observation is promoted to a blocker. This is recheck one of the convergence policy; a targeted repair and recheck may follow within the same frozen scope.

## Owner Validation

Real-terminal acceptance of the chip, popover and history remains pending for T-076/T-079/T-081. No owner permission question is needed to finish this Check. If the owner prefers wider/scrolling presentation, that needs an explicit design decision, not silent checker acceptance. Tasks remain done and the Objective remains in progress.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches** — missing width/no-scroll assertions are findings above.
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [ ] STYLE-07 **One source of truth** — sparkValues still duplicates buildTrend value assembly; stored timestamp comparison is repeated in chip.go rather than shared compareSnapshots.
- [ ] STYLE-08 **Comments explain why** — healthPopoverWidth still says 72 while its maximum constant is 100 (health_view.go:102).
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs**

## Reviewed Content Hashes

- `internal/codehealth/dashboard.go`: `37e8ce0d860591099090f7985d00429884d68dd5ac3884904dd90634c700c57d`
- `internal/codehealth/dashboard_test.go`: `5e70b56d15ebe12e29a466aeed43204ba83a7c7a92e1a8e07f60618d00153e80`
- `internal/codehealth/dashboard_copy.go`: `7a67d823cbf79b8a7eb7108f11b03bcab81831c206565febdbad2052cda5c7b9`
- `internal/codehealth/dashboard_copy_test.go`: `c1ea069078d3072a6f11c76e302e9d397c1129a161600018578eb7455653ff5f`
- `internal/codehealth/dashboard_signoff_test.go`: `b25f3a8aa8356d1ed3db47661f3c612698374425acf3a665b47324500f5d3155`
- `internal/codehealth/chip.go`: `df3dc55d82b53be4d239ce6aa715dc92f60868fbf5e1689298302c97dd0a7d02`
- `internal/codehealth/chip_test.go`: `3878037128752387de4144bd66bd6615e615e9a07dc39915cd8e1371f70b78c9`
- `internal/codehealth/spark.go`: `83485582282409df37cb16fba31692ce608133d1fb8e5d2c887a2309f7d13213`
- `internal/codehealth/spark_test.go`: `0cb59506e9dc807d251033e57652963531c28bf597df2ea4bf834ee5a78a4227`
- `internal/board/v2/health.go`: `501027cb0c3e001daaf3fb9a66b26c34e4afd657de250648b0bc31a04f472404`
- `internal/board/v2/health_view.go`: `0fdba5df22e38780ccebf8a312f5aadffcae74a5a3eb54ef42d17863e1a8b816`
- `internal/board/v2/health_test.go`: `32a0fbeb9b2cd0d87397f21e3783abe2c58c0154f4ae82cfc610dde88c192c24`
- `internal/board/v2/health_popover_test.go`: `4804f8c2f0f596288d9e39b945647dea0de708c4d28b6384fe55e9222335558c`
- `internal/board/v2/health_history_test.go`: `f31a3d873bb50d7f928ceab872c1bbaaa15ce0402b9bae496b2cc4e492c54aad`
- `internal/board/v2/view.go`: `12e0433f0db1777ab19db9ed9db1d7da6a4f5a5e10aac62d1670c8be0082d0a7`
- `internal/board/v2/update.go`: `fe7d030ac6b58e37bc61f7e4bcb117e6031187ec6bc6591ddfccd4ba358d516d`
- `internal/board/v2/help.go`: `94044d488a53cf637e97ef8c6f3f7ce255f4e53d027fc204658f891511457f46`
- `internal/board/v2/load.go`: `d1b0dcd40ed654821fc1019b9c10c04a99a0f1e22be00896103191b6f6984ad6`
- `internal/board/v2/detail_test.go`: `c87b70ad6a0e0a0d756c14a59a4754448ba412e69a49de1ff8d620a5a6ca148c`
- `internal/board/v2/health_realcopy_test.go`: `ef0d7d803e97b9d38e5b205eb3011398b89c0528bf3f90c3f0238e1daaf95d5e`
