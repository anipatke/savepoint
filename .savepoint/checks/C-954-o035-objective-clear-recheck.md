---
id: C-954
scope: {kind: objective, id: O-035}
result: CLEAR
checked_by: {role: checker, session: recheck-o035-targeted-20261002-independent}
executed_session: user-request-targeted-i107-i112-repair
checked_at: '2026-10-01T23:35:44Z'
health_snapshot: sha256:2cd13b62c62200ad20dcea0e4d4d9e75491630a670cfe5c4accbc75bf0677b92
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
issues: []
supersedes: C-953
---

# C-954: O-035 Targeted Final Recheck

## Closure Map

| Prior Issue | Closure | Independent proof |
|---|---|---|
| I-107 | Verified by C-954 | Real values retain all default aims and direction; actual frame width is <=72 at 80x20/24/40 |
| I-108 | Verified by C-954 | Early/restarted/thin notes visible with and without blocks |
| I-109 | Verified by C-954 | Reversed file times and persisted-time/identity ties agree with dashboard |
| I-110 | Verified by C-954 | Good measurements at 1/2/3 checks and confidence/baseline causes explained accurately |
| I-111 | Verified by C-954 | All signal progress omits providers in both modes |
| I-112 | Verified by C-954 | Fixed five-capability grouping; blocking instance preserved; row content does not move during selection |

CLEAR. This conversation did not implement either repair and remains independent of the executor. All seven Tasks remain done with their recorded owner waivers. This Check grants technical clearance, not owner acceptance or Objective completion. Real-terminal visual acceptance remains an owner activity.

## Frozen Scope And Admission Ledger

This is the targeted second remediation recheck within the convergence limit. C-952's numbered scope lock, complete coverage matrix, side-effect inventory, external-boundary exclusions and materiality rules are unchanged. C-953's admission ledger was carried forward before testing, with only its repair claims updated: I-107 restores 72; I-112 replaces the scrolling window with a grouped row. No new axis, dependency layer, acceptance interpretation or supported-value perimeter was added.

| Recheck item | Prior Issue | Exact frozen cell | Outcome |
|---|---|---|---|
| Real copy and <=72 frame | I-107 | Popover normal 80x20/24/40, real-copy output, original T-079 maximum | Passed |
| Note rendering | I-108 | Spark final early/restart/thin and all selection cells | Passed |
| Persisted recency | I-109 | Chip save-order/timestamps/ties, load vs dashboard | Passed |
| Thin/confidence/decline explanations | I-110 | Copy thin history vs level wording and partial/stale cells | Passed |
| Provider-free progress | I-111 | Popover refresh/progress output | Passed |
| Grouped first/middle instances and final blocker | I-112 | Multiple instances/all selections and no-scroll contract | Passed |
| All other original matrix cells | C-952 | Original completed matrix, same exclusions | Passed |

## Repair Verification

`healthPopoverMaxWidth` is 72 again. The independent original realistic-copy scenario renders five real default aims, complete direction, measured numbers, and the same 17-line frame at 80x20, 80x24 and 80x40. TestCheckO035FrozenDimensions measures rendered cell width, not the implementation constant, and passes every size. The repository's TestHealthPopoverKeepsEveryYardstickWithRealCopy now asserts the same cap, including its larger-terminal supplemental case.

HealthOverlay.rows groups instances per capability. It selects a blocking instance first, otherwise the worst classification, and labels the multiplicity and represented instance. The data projection remains unchanged; it still describes every configured instance. Independent TestCheckO035MultipleInstances feeds duplicate Tests rows including a blocking instance, confirms all five capability rows, the blocking Meaning/SignOff, and the final dependency signal. Across six Down presses, normalized row content is byte-identical, and the cursor stops at four. This reproduces both the original hidden-capability case and C-953's scroll failure with the repaired representation. The named repository regression additionally covers duplicate first and middle capabilities, mixed labels, worst-instance name, final blocking signal and each documented size. The discarded window and up/down overflow cues are gone.

The earlier four repairs remain proven by rerunning their independent scenarios rather than reusing C-953 alone. Note assertions inspect final selected detail after the implementation moved notes there; Dependency's short display name is recognized. No old helper location or all-instance row count is imposed as a requirement. Actual five-capability visibility, truthful representative details and no scrolling are the frozen invariants.

## Completed Coverage Matrix

| Frozen row | Normal and boundary | Failure/missing | State/sequence | Output |
|---|---|---|---|---|
| Copy | every signal/label, default/override thresholds, Good/Watch below/at/above and 0/1/many paths pass | absent/unmeasured/not configured and partial/stale pass | independent four signals at 1/2/3 observations and coverage confidence/decline cases pass | real-length values/aims and truthful meanings pass |
| Gate | all-good, failing tests, high/unrated, opt-in pass | required failed/stale/absent and invalid config pass | official/manual/no-official pass | per-row and overall disposition match Evaluate |
| Spark | independent 0..12 × five signals × spread below/at/above movement pass; rising/falling/flat regressions pass | manual/partial/incompatible excluded; independent NaN/Inf rejection passes | restart/manual newest pass | early/restart/thin notes and direction visible |
| Chip | three states/four labels/counts, 50..100 header widths pass | missing/corrupt fallback passes | reversed save/mtime, equal mtime and timestamp ties pass | header/dashboard agreement and refresh update pass |
| Popover | 80x20/24/40, <=72 width, 17 height, all selections and history 0/1/10/14 pass | loading/error/unknown/partial pass | H/select/h/Esc/removed keys/cursor restore, progress/cancel/retry/help quit and grouped instances pass | real copy, glyphs, ANSI/Unicode width and provider omission pass |

The original N/A cells remain N/A for the same reasons: no new parser or serialization, lifecycle mutation, network/auth/billing boundary or non-TTY interactive view. Independent Unicode uses ASCII, combining marks, wide characters, variation selectors, emoji modifiers, regional flags, joined emoji and ANSI-formatted text. No control-character requirement was added. Original input/storage-validation and workflow/runner tests pass in the full suite. All mandatory frozen cells are classified; no material unverified criterion remains.

All eleven Objective success conditions are Proven against this completed matrix. T-075 copy and T-076 chip repairs remain accurate; T-077 sign-off, T-078 spark data, T-080 state/refresh/cancel and T-081 history remain proven; T-079 now meets the two remaining width/non-scroll requirements without losing required yardsticks or the final capability. The Full Objective Check covers every waived Task and cross-Task integration. Design reconciliation remains as recorded in C-952: read-only health projection, pure rendering and explicit commands, unchanged classifiers, thresholds, provider/storage identities and retention. The grouped display extends board presentation only; the Codebase Map remains accurate.

## Gates And Official Health

- Fresh `make test-full`: exit 0, go1.26.2 linux/amd64, 2026-10-01 approximately 23:31–23:33 UTC. Uncached full Go suite with reports and Linux/Darwin/Windows build targets passed. This does not claim a new native Windows execution.
- `make build`: exit 0; current source CLI used for collection.
- `go test ./internal/codehealth ./internal/board/v2 -count=1 -run 'Test(Dashboard|LoadChip|LatestSnapshot|Spark|Health|HeaderChip|Meaning)'`: exit 0, both packages passed.
- `go test -overlay /tmp/o035-probes/overlay.json ./internal/board/v2 ./internal/codehealth -run TestCheckO035 -v`: exit 0. RealCopy, SparkNotes, ProgressProvider, Unicode, MultipleInstances, FrozenDimensions, ChipRecency, Meaning, IndependentSparkMatrix, GateCases and CauseAdjacent all passed. External overlay and scratch logs were never installed as repository test source. A first invocation hit read-only Go cache permissions; an approved escalation reran it successfully.
- `git diff --check`: exit 0.
- `./savepoint health check O-035`: exit 0, official snapshot in frontmatter. **Code Health does not block clearance.** Tests, coverage, complexity and duplication have no blocking findings. Optional osv-scanner collection failed without a measurement; that is not a security finding and does not prevent clearance.
- `./savepoint resume` strict-loads this immutable record and verified Issue resolutions.

Workflow and side-effect inventory is unchanged from C-952. Rendering performs no IO. Board regression evidence includes explicit load/freshness ordering and stale message suppression, refresh progress/reload, cancellation retaining last result, load error retention, Help interrupt/quit, cursor restoration and history return. The full collection and runner suite verifies the unchanged subprocess cancellation/report preservation boundary. No implementation, health config, router, Design, Task status or Objective status was edited by this checker. Only official health artifacts, this Check and six permitted Issue verifications were written.

## Materiality And Owner Handoff

No materiality actions are required; all admitted Issues are repaired and verified by this CLEAR Check. Optional scan failure remains a nonblocking observation. No new out-of-scope blocker was found. Owner may accept the Objective after any desired real-terminal visual validation; this record does not perform that acceptance or close any Task/Objective.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches**
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [ ] STYLE-07 **One source of truth** — sparkValues duplicates buildTrend value assembly; chip's timestamp ordering and board's label severity mirror existing package rules (spark.go, chip.go, health.go). Advisory only.
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs**

## Reviewed Working-Tree Identity

The implementation remains uncommitted atop the named head. Since C-953, only dashboard.go, health.go, health_view.go and health_realcopy_test.go differ among its recorded content hashes. All other recorded paths match. Every evidence file exists; scratch from the earlier partial T-075 plan was already explicitly reconciled. Current hashes:

- `internal/codehealth/dashboard.go`: `051cfc392eba47e48880a1b03d388a6a4897e22a9f2b6217ebb6714661fb127d`
- `internal/codehealth/dashboard_test.go`: `5e70b56d15ebe12e29a466aeed43204ba83a7c7a92e1a8e07f60618d00153e80`
- `internal/codehealth/dashboard_copy.go`: `7a67d823cbf79b8a7eb7108f11b03bcab81831c206565febdbad2052cda5c7b9`
- `internal/codehealth/dashboard_copy_test.go`: `c1ea069078d3072a6f11c76e302e9d397c1129a161600018578eb7455653ff5f`
- `internal/codehealth/dashboard_signoff_test.go`: `b25f3a8aa8356d1ed3db47661f3c612698374425acf3a665b47324500f5d3155`
- `internal/codehealth/chip.go`: `df3dc55d82b53be4d239ce6aa715dc92f60868fbf5e1689298302c97dd0a7d02`
- `internal/codehealth/chip_test.go`: `3878037128752387de4144bd66bd6615e615e9a07dc39915cd8e1371f70b78c9`
- `internal/codehealth/spark.go`: `83485582282409df37cb16fba31692ce608133d1fb8e5d2c887a2309f7d13213`
- `internal/codehealth/spark_test.go`: `0cb59506e9dc807d251033e57652963531c28bf597df2ea4bf834ee5a78a4227`
- `internal/board/v2/health.go`: `596389fe455a80d8afc3c1b6f0652e726a0ff6326985a9802d55920b748a74d7`
- `internal/board/v2/health_view.go`: `8fb00b0806890c5f02144ceb5aec1bcd23931dbdf1c19ddc17dbf8b0bb455ef8`
- `internal/board/v2/health_test.go`: `32a0fbeb9b2cd0d87397f21e3783abe2c58c0154f4ae82cfc610dde88c192c24`
- `internal/board/v2/health_popover_test.go`: `4804f8c2f0f596288d9e39b945647dea0de708c4d28b6384fe55e9222335558c`
- `internal/board/v2/health_history_test.go`: `f31a3d873bb50d7f928ceab872c1bbaaa15ce0402b9bae496b2cc4e492c54aad`
- `internal/board/v2/view.go`: `12e0433f0db1777ab19db9ed9db1d7da6a4f5a5e10aac62d1670c8be0082d0a7`
- `internal/board/v2/update.go`: `fe7d030ac6b58e37bc61f7e4bcb117e6031187ec6bc6591ddfccd4ba358d516d`
- `internal/board/v2/help.go`: `94044d488a53cf637e97ef8c6f3f7ce255f4e53d027fc204658f891511457f46`
- `internal/board/v2/load.go`: `d1b0dcd40ed654821fc1019b9c10c04a99a0f1e22be00896103191b6f6984ad6`
- `internal/board/v2/detail_test.go`: `c87b70ad6a0e0a0d756c14a59a4754448ba412e69a49de1ff8d620a5a6ca148c`
- `internal/board/v2/health_realcopy_test.go`: `e74078722e68768aef471202f9605fe0d3ac8aacbd668c7ab3ff4a23a3936dee`
