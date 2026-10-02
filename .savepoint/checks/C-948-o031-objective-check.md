---
id: C-948
scope: {kind: objective, id: O-031}
result: NEEDS WORK
checked_by: {role: checker, session: check-o031-20261002-independent}
executed_session: prior-o031-executor-sessions
checked_at: '2026-10-01T20:47:19Z'
reviewed:
  base_commit: 01c3067
  head_commit: 9db1988
  files:
    - .gitignore
    - AGENTS.md
    - Makefile
    - internal/board/v2/detail.go
    - internal/board/v2/detail_test.go
    - internal/board/v2/detail_view.go
    - internal/board/v2/health.go
    - internal/board/v2/health_test.go
    - internal/board/v2/health_view.go
    - internal/board/v2/help.go
    - internal/board/v2/io.go
    - internal/board/v2/load.go
    - internal/board/v2/model.go
    - internal/board/v2/update.go
    - internal/board/v2/view.go
    - internal/buildtool/main.go
    - internal/buildtool/main_test.go
    - internal/codehealth/collect.go
    - internal/codehealth/collect_test.go
    - internal/codehealth/dashboard.go
    - internal/codehealth/dashboard_freshness.go
    - internal/codehealth/dashboard_freshness_test.go
    - internal/codehealth/dashboard_test.go
    - internal/codehealth/errors.go
    - internal/healthcheck/healthcheck.go
    - internal/healthcheck/healthcheck_test.go
    - internal/styles/styles.go
  dependencies:
    - internal/codehealth/storage.go
    - internal/codehealth/repository.go
    - internal/codehealth/classification.go
    - internal/codehealth/history.go
    - internal/codehealth/runner.go
    - internal/codehealth/runner_windows.go
issues: [I-104, I-105, I-106]
supersedes: null
---

# C-948: O-031 Full Objective Check

## Result

NEEDS WORK: signal details/history are inaccessible when they overflow, refresh Help swallows Ctrl+C, and an interrupted test subprocess replaces the previous complete JSON report with an incomplete stream. This session is independent from all executor conversations. No implementation, Task status, router, or Design edits were made. Scope includes the uncommitted T-074 implementation and owner-authored record/router changes present at the start.

Every owned Task T-070–T-074 is done and carries a task-specific owner waiver with actor, reason, and timestamp. Waivers are not CLEAR; all five Tasks were reviewed here. No prior Check or matching Issue was linked to O-031; duplicate search included existing scroll, health, cancellation, and report Issues before allocating I-104–I-106. Scoped evidence files exist; temporary Check probes were deliberately removed and are reproduced below.

## Frozen Scope Lock

# O-031 Full Check scope lock

1. Requirements: all O-031 success conditions and Done When criteria of T-070–T-074; ARCH-02/03/04, DATA-03, FS-05/06, CFG-02/03, TEST-01/02/04/08/09. Full gate make test-full, git diff --check, official health collection. STYLE rules advisory. Owner terminal validation is separate from technical clearance.
2. Base 01c3067, head 9db1988 plus current uncommitted T-074 changes. Changed dashboard/freshness, collection cancellation, healthcheck error mapping, board health model/view/commands/help/styles, compact Objective Check lookup/detail, buildtool reports and Makefile/gitignore. Public surfaces: LoadDashboard, SnapshotLabels, DashboardFreshness, Collect, healthcheck.Run, Model.Update/View (H/R/Enter/v/arrows/Esc/q/ctrl+c/?), loadProject, buildtool test/focused-test.
3. Dependencies: Store config/history load/save, Assess/Overall/Trend persisted policy, Git Observe/Relate, ExecRunner process cancellation, Bubble Tea async commands, go test JSON/profile subprocess. No unrelated provider normalization or lifecycle internals. No implementation repairs.
4. Matrix rows: D dashboard (unconfigured/first-run/measured; five capabilities; multiple/missing instances; required/optional failure, partial/absent/unsupported; official/manual, ten-history bound, compatible/reset trends); F freshness (same/behind/ahead/diverged, dirty and Git failure/cancellation); U board (opening/closing/exclusive/help, overview/details/history at tall/short/minimum-width sizes, repeated arrows and Unicode text; load/error/reopen and delayed messages); R refresh (configuration guard, sequential progress, repeated R, success/reload, failure, Esc cancellation/retry, q and ctrl+c including help); C collection (nil/progress, cancellation before/mid/after final instance, save preservation and real process stop); S summary (official/missing/no field/missing or damaged directory, Task/non-TTY unchanged); B reports (pass/failure/repeat/focused, interrupted/error stream publication, atomic replacement, output/exit unchanged, relative paths/platform). Representations are real saved config/snapshots, injected board models, rendered ANSI/plain output, raw test JSON/profile. Numeric nonfinite/mixed types are N/A for new APIs: existing strict validators supply typed models. Remote server/redirect/secrets N/A: local-only work.
5. Materiality boundary: supported project data and documented keyboard/gate paths; reproducible unmet acceptance or named rule introduced/touched here. External boundary cells: local Git/tool/go executables, success/nonzero/unavailable/timeout/cancel/malformed reports, cleanup/no orphan and saved state; no network redirects. Windows execution evidence must be classified separately from cross-build.

Workflow inventory:
|Order|Operation|Effect|Failure owner/final state|Oracle|
|---|---|---|---|---|
|1|H/load config and snapshots|overlay, read only|status line, previous data retained|fake call counts plus Store fixtures|
|2|Git comparison|freshness message, read only|unknown; no scan|scripted Git and real temporary Git repository|
|3|R/config and Collect|cancellable context, sequential callbacks/tool effects|failed tool is per-instance; fatal error retains dashboard|event sequence and snapshot listing|
|4|classify/save immutable snapshot|one manual snapshot on success|cancel before save leaves store unchanged|before/after listings, official reader tests|
|5|reload dashboard/freshness|replace displayed result|load error visible|reducer and persisted result|
|6|Esc/q/ctrl+c|cancel/close/quit|no saved cancelled run, process stop|context and real runner alive tests|
|7|go test JSON/profile temporary files|temporary reports|test failure still reports; interruption must retain complete reports|temporary projects, sentinel old files|
|8|close/rename reports|atomic report publication and summary|IO failure clear, no half report|parse raw JSON/profile and directory contents|

Initial lock written before independent probes. No prior linked Check/Issue found.


## Coverage Matrix and acceptance classification

The rows below apply the finite cells in the numbered lock. “Proven” covers cells verified by independent probes, current outcome tests, and code trace; it never means that the whole Objective passed. Issue cells retain their exact frozen identity for recheck. Unverified platform evidence is explicit rather than silently treated as Linux evidence.

| Cell | Applicable inputs, states, surfaces and boundaries | Evidence | Classification |
|---|---|---|---|
| D1 | LoadDashboard: no config / first run / damaged config or history; no tool run or write | TestLoadDashboardNotConfigured, TestLoadDashboardFirstRun, TestLoadDashboardNeverRunsAnythingOrWrites, TestLoadDashboardReportsDamagedStorage; independent production board open | Proven |
| D2 | Five signals, canonical order, multiple instances and missing snapshot instance; available/partial/failed/timed out/unavailable/absent/unsupported crossed with required/optional | TestDashboardAllFiveSignalsMeasured, MultipleInstancesOfOneCapabilityEachGetARow, InstanceMissingFromSnapshotIsUnknown, FailuresStayDistinctAndRequiredIsShown; independent 14-cell saved matrix (seven outcomes x required/optional) | Proven |
| D3 | Labels and overall from persisted summary, incomplete rows not Good, newest manual/official, comparison-compatible/reset/manual-excluded trend, provenance/scope/evidence, history 10/newest first, banned claims | TestDashboardLabelsComeFromThePersistedSummary, PartialSupportStaysVisibleAndNotGood, ManualNewestWinsOverOfficial, ComparableTrendAndReset, ManualResultsAreShownNotCounted, HistoryIsBoundedAndNewestFirst, WordingTableHasNoBannedClaims; independent saved matrix checks all generated copy | Proven |
| F1 | DashboardFreshness: same/behind/ahead/diverged, shallow count unavailable, dirty/unborn, unavailable Git and cancellation | TestDashboardFreshness table and TestDashboardFreshnessCancelledIsUnknown; independent real temporary Git repository checks same, dirty edit, and one committed change | Proven |
| U1 | H from board, saved load then explicit freshness, no IO in rendering; unconfigured R ignored, first-run R offered, exclusivity, Esc cursor restore, delayed freshness for older snapshot ignored, load errors retain previous data | TestHealthShowsResultsBeforeFreshnessArrives, StaleFreshnessForAnOlderSnapshotIsIgnored, NotConfiguredExplainsAndRefreshDoesNothing, FirstRunOffersRefresh, IsExclusiveWithTheOtherOverlays, EscRestoresTheBoardCursor, LoadErrorKeepsThePreviousView, LoadErrorOnFirstOpenIsAStatusLineNotACrash; independent production wiring calls real LoadDashboard/Freshness | Proven |
| U2 | Overview/detail/Enter/v/history and repeated arrows at 80x20, 80x24, 80x40; wrapping and long evidence; finite lower/upper navigation boundaries | Independent TestO031IndependentViewportMatrix: detail Offset remains 0 at all sizes; history unreachable at 20/24; current dimension tests still pass | Issue I-104 |
| U3 | Atari-Noir label glyphs/words, no overall number; five signals/partial/failed/optional/not-configured and trend/basis copy; minimum-width size and Help | TestHealthShowsFiveSignalsWithLabelsReasonsAndTrends, NonGoodRowsAreNeverShownAsGood, ComparableAndResetTrendsAreShownAsWorded, HelpListsTheKeys, FitsTheBoardsMinimumWidth; render code trace | Proven; access to overflow is U2 |
| R1 | Refresh configuration guard, sequential callback messages/counts, repeated R ignored, success/reload/freshness, failed refresh retains data, Esc cancellation/no save/retry, direct q/ctrl+c cancel | TestHealthRefreshReportsProgressInOrderAndReloadsOnCompletion, EscCancelsARefreshAndKeepsThePreviousResult, RefreshErrorIsAStatusLine, QuitDuringRefreshCancelsBeforeQuitting; independent production H/first-run/R/progress/manual-save/load/freshness workflow | Proven |
| R2 | Refresh → ? Help → Ctrl+C | Independent TestO031IndependentQuitFromRefreshHelp/ctrl+c: cancelled=false and no quit command | Issue I-105 |
| C1 | Collect: nil callback compatibility, progress before each instance/config order; cancel before first / mid-tool / after final before save; directory unchanged; healthcheck named cancelled diagnostic | TestCollectReportsProgressInOrder, TestCollectCancelledSavesNothing, TestRun_cancelledContextSavesNothing; current source trace and full gate | Proven |
| C2 | Local executed provider success/nonzero/unavailable/timeout/cancel, report-only no execution, malformed/absent/partial report isolation, real Linux child-process kill and cleanup | TestCollectExecutedOutcomesStayDistinct, ReaderProblemsAreContained, ReportOnlyNeverRunsATool, WithRealProcessesAndReportFiles, CancelledRealProcessLeavesNoTemporaryFiles; TestExecRunnerStopsOnDeadlineAndCancel | Proven on Linux; native Windows execution unverified |
| S1 | Compact Objective health line from stored official snapshot; missing snapshot / no field / missing or unreadable storage; no Task health line; badges/Next/non-TTY unchanged | TestObjectiveDetailShowsTheHealthLineForAStoredSnapshot, ReportsAMissingHealthSnapshot, WithoutAHealthSnapshotFieldShowsNoHealthLine, TestTaskDetailNeverShowsAHealthLine, TestMissingHealthDirectoryIsNotABoardError, TestSnapshotLabelsNameEachStoredSnapshotAndTolerateMissingStorage; load/detail diff and full suite | Proven |
| B1 | make test/test-fast reports, pass/fail publication, focused writes neither, report names and ignore entries, summary/exit status, temp cleanup, filepath paths, repeated atomic replace | TestRunGoTestWithReportsWritesReportsOnPass/OnFailure/LeavesNoTempFiles, TestRunGoTestFocusedWritesNoReports, TestRunGoTestCommandReportsTimingAndPropagatesFailure; two fresh full gates replace actual reports successfully | Proven for completed runs |
| B2 | Child interrupted after partial JSON stream, previous complete report sentinel, nonzero exit | Independent TestO031IndependentInterruptedReports returns signal: terminated and observes old JSON replaced by one start event | Issue I-106 |
| B3 | Linux full suite and all six cross-build targets; native Windows full-suite result for this working tree | Fresh make test-full passes on linux/amd64; no native Windows runner/result accessible in this session | Linux/cross-build Proven; native Windows Unverified |
| E1 | No network/server/browser/redirect boundary; numeric nonfinite/wrong type/model mutation beyond typed new APIs | No new such entry points; existing config/snapshot validators guard saved input. These are N/A to new typed projection/reducer behavior, not newly added policy | N/A with scope reason |
| E2 | Redirected/non-TTY board output | New health surface is TUI-only and compact summary is detail-only; existing plain output path untouched and full non-TTY board suite passes | Proven unchanged; refresh UI N/A in non-TTY |
| E3 | NO_COLOR/dumb output and Unicode width classes (ASCII, control, combining, variation selector, wide, modifiers, flags, joined emoji) | Scoped health rendering uses existing fitLine/wrapDetailLine/frameColumn and lipgloss.Width rather than a new width algorithm; current minimum-width test covers wrapped paths and glyphs. Independent viewport probe covers overflow; no separate complete text-class oracle was run | Unverified text-class matrix; overflow Issue U2 remains independently proven |

Task criteria mapping: T-070 criteria 1–7 Proven through D/F. T-071 criteria 1–3 and cancellation test cells Proven; criterion 4/5 Windows runtime component Unverified (Linux process-stop cells Proven). T-072 criteria 1–4, 7–8 Proven; criterion 5 Issue U2, criterion 6 Issue R2; criterion 9 incomplete with U2/R2 and text-class evidence Unverified. T-073 criteria 1–5 Proven through S1. T-074 criteria 1,3–5 Proven; criterion 2 Issue B2; criterion 6 Linux/cross-build Proven and native Windows CI result Unverified. Objective success conditions 1–4 and complete-data/nonclaiming semantics have supporting proof; complete details/history and cancellation contain Issues, and mandatory matrix evidence gaps are disclosed.

## Workflow, external effects, and adversarial pass

Opening reads persisted config/snapshots only, then schedules a bounded local Git command; rendering is pure. Production board integration in the independent harness used a temporary V2 project and Git repository with a real Go JSON report, then drove first-run → R → progress → successful manual persistence → real load → freshness. Store independently reported exactly one manual snapshot. No test command ran during report-only refresh. Persisted classification matches the established policy, and optional failure/partial rows remain explicit; no overall score or readiness claim is added.

Collection cancel checks occur after repository observation and immediately before save. The existing before/mid/after-last tests compare snapshot directory contents; real process tests verify Linux children stop. Provider failures remain individual outcomes. Official cancellation mapping names the error and says no snapshot was saved. Snapshot storage remains the existing immutable boundary. Stale snapshots, dirty files, and comparison resets retain distinctions.

The adversarial pass tried repeated navigation rather than checking only initial layout, Help as an alternate path to an active refresh, and abnormal subprocess termination rather than ordinary test failure. All three reveal violations inside the original lock. Report rename is byte-atomic but currently commits incomplete streams. Source review confirms rendering never reads files or invokes tools, board commands convert .savepoint roots to project roots, no classification threshold or snapshot format changed, and summary decoration leaves gate/badge/plain output decisions alone.

Workflow cleanup, stdout summary/exit codes, config load failures and store errors are covered by full package tests and source trace. Saved-state preservation on cancelled refresh is proven on Linux; report preservation on interrupted child is disproven. Native Windows and the complete explicit Unicode/no-colour matrix remain Unverified, preventing claiming every matrix cell has technical proof. No unsupported hypothetical dependency layer was added. All remaining lock cells were classified after finding the Issues; none was skipped because the first Issue was found.

## Commands and evidence

- Go toolchain: go1.26.2 linux/amd64. Base 01c3067, HEAD 9db1988 plus pre-existing uncommitted T-074 changes. Fresh execution on 2026-10-01 UTC / 2026-10-02 Australia/Sydney.
- `make build`: exit 0.
- `make test-full`: exit 0 twice; last fresh run after deliberately removing Check scratch files. It ran uncached full Go tests, produced go-test.json and coverage.out, and built Linux, Darwin and Windows targets. No earlier executor result was reused.
- `git diff --check`: exit 0 before record writing; final check follows record validation.
- `go test -race ./internal/board/v2 ./internal/codehealth ./internal/healthcheck ./internal/buildtool -run 'Health|Dashboard|Collect|RunGoTest|ExecRunnerStops' -count=1`: exit 0 (healthcheck had no matching tests; its cancelled Run case is covered by the full gate).
- `go test ./internal/codehealth -run TestO031Independent -count=1`: exit 0, 14-cell persisted outcome/required matrix and real Git checks. An initial harness partial-result fixture was invalid because it omitted its reason; corrected the fixture, then reran successfully. This is a probe setup error, not a product finding.
- `go test ./internal/board/v2 -run TestO031IndependentProductionWorkflow -count=1`: exit 0.
- `go test ./internal/board/v2 -run TestO031Independent -count=1`: exit 1 on independent viewport and refresh-Help assertions; exact outcomes are in I-104/I-105. q from Help closes it by the existing convention; the Ctrl+C cell is the blocking reproduction.
- `go test ./internal/buildtool -run TestO031Independent -count=1`: exit 1: interrupted tool replaced complete report with only a start event, signal: terminated.
- After the successful full gate, `./savepoint health check O-031`: exit 0, “Code Health is not configured for this project; nothing was collected.” Code Health not configured; no health_snapshot field. Setup is human-only and was not run. No manual result from the temporary integration test counts as official Check evidence.

## Materiality

| Issue | Likelihood | Impact | Materiality | Recommendation |
|---|---|---|---|---|
| I-104 | High: normal short terminals or long evidence | Medium: promised explanation, provenance and history cannot be read | High for this explanatory UI Objective | Fix scrolling directly under the Issue; recheck the frozen U2 sizes/content cases |
| I-105 | Medium: Help is offered during refresh | Medium: normal interrupt does not cancel active work or quit | Medium | Fix interrupt dispatch with a Help-active regression |
| I-106 | Low: test child is interrupted | Medium: last complete generated evidence is lost; no user source loss | Medium | Preserve reports on interrupted/incomplete streams while retaining ordinary failure reports |

## Design reconciliation and observations

Design sections 1/6/7/8/9/11/12/13 support the package/service boundaries and gate order. The implementation uses the existing collector, local Git relation and saved classification. Drift notes correctly identify documentation now needing planner reconciliation: section 8 lacks H, dashboard/refresh/cancellation and the Check health line; section 13 omits the new report artifacts. No Design remediation was made by this checker. There is no new classification or provider/schema decision to ratify.

Nonblocking observations: Help q retains the existing close-help convention even during refresh; I-105 narrowly requires Ctrl+C. The first-run wording “Nothing is changed in your project” could be clearer about saving a manual snapshot. These observations do not add repair conditions. Complete no-colour/Unicode matrix and native Windows execution were not obtained here and are Unverified evidence cells, separately from reproducible defects.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [ ] STYLE-03 **Test branches** — current viewport tests miss actual scrolling, Help interruption, and abnormal subprocess report publication; see the independent probes and I-104–I-106.
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth**
- [x] STYLE-08 **Comments explain why**
- [ ] STYLE-09 **Content lives in data** — health_view.go detail field labels and history/evidence strings remain inline despite a central copy table; advisory only.
- [x] STYLE-10 **Small diffs** — additions follow the five scoped Tasks; no unrelated implementation changes.

## Owner validation and handoff

T-072 declares required owner validation. Its real-terminal User Check was not accepted in the evidence read here, and this Check does not record owner acceptance. After repair and a new independent Full Check, the owner still needs the terminal walkthrough for colours, real refresh/cancel and return navigation. This NEEDS WORK Check cannot authorize Objective closure. Completed Tasks retain done status; route direct repair under I-104–I-106 rather than retreating them. Router selection and Goal remain unchanged. Each later recheck must use this frozen lock and supersede C-948 with a new immutable record.

## Reproduction harnesses

The following were temporary Check-only files, deliberately removed before the final full gate. Restore them at their exact paths only while running the independent probes; they are not implementation changes. Their source is preserved here for repeatable remediation checks.

### internal/board/v2/o031_check_scratch_test.go

```go
package v2

import (
 "os"
 "os/exec"
 "path/filepath"
 "strings"
 "testing"
 "github.com/opencode/savepoint/internal/codehealth"
)

func TestO031IndependentProductionWorkflow(t *testing.T) {
 m:=openSizedBoard(t,writeValidProject(t),100,30)
 root:=filepath.Dir(m.Root)
 cmd:=exec.Command("git","init","-q");cmd.Dir=root
 if out,err:=cmd.CombinedOutput();err!=nil{t.Fatalf("git init: %v %s",err,out)}
 store:=codehealth.NewStore(root)
 cfg:=codehealth.Config{Version:codehealth.ConfigVersion,Capabilities:[]codehealth.CapabilityConfig{{Capability:codehealth.CapabilityTests,Provider:codehealth.ProviderGoTestJSON,Report:"go-test.json",Required:true}}}
 if _,err:=store.SaveConfig(cfg);err!=nil{t.Fatal(err)}
 report:="{\"Action\":\"run\",\"Package\":\"fixture\",\"Test\":\"TestExample\"}\n{\"Action\":\"pass\",\"Package\":\"fixture\",\"Test\":\"TestExample\"}\n{\"Action\":\"pass\",\"Package\":\"fixture\"}\n"
 if err:=os.WriteFile(filepath.Join(root,"go-test.json"),[]byte(report),0600);err!=nil{t.Fatal(err)}
 next,load:=m.Update(keyMsg("H"));m=settle(t,next.(Model),load)
 if m.Health.Dashboard.State!=codehealth.DashboardFirstRun{t.Fatal("not first run")}
 m,event:=startRefresh(t,m)
 msg:=event();if _,ok:=msg.(healthProgressMsg);!ok{t.Fatalf("progress: %T",msg)};m=feed(t,m,msg)
 msg=event();done,ok:=msg.(healthRefreshDoneMsg);if !ok||done.Err!=nil{t.Fatalf("done: %+v",msg)}
 next,load=m.Update(msg);m=settle(t,next.(Model),load)
 snaps,err:=store.LoadSnapshots();if err!=nil||len(snaps)!=1{t.Fatalf("snapshots: %v %v",snaps,err)}
 if snaps[0].Origin!=codehealth.OriginManual || m.Health.Dashboard.State!=codehealth.DashboardMeasured || m.Health.Freshness==nil{t.Fatal("manual workflow incomplete")}
}

func TestO031IndependentViewportMatrix(t *testing.T) {
 for _, height := range []int{20,24,40} {
  t.Run(string(rune('A'+height)),func(t *testing.T) {
   rows:=fiveSignals()
   rows[0].Evidence=[]codehealth.EvidenceRef{{Path:"FINAL_EVIDENCE_SENTINEL",Note:strings.Repeat("detail ",60)}}
   m:=openHealthScreenAt(t,&fakeHealth{dashboard:measuredDashboard(rows)},80,height)
   m=press(t,m,"enter")
   before:=screen(m)
   for i:=0;i<60;i++ { m=press(t,m,"down") }
   if m.Health.Offset==0 || screen(m)==before || !strings.Contains(screen(m),"FINAL_EVIDENCE_SENTINEL") {
    t.Errorf("detail cannot reach final evidence: height=%d offset=%d",height,m.Health.Offset)
   }
   m=press(t,m,"esc")
   for i:=0;i<60;i++ { m=press(t,m,"down") }
   if !strings.Contains(screen(m),"Manual refresh") { t.Errorf("history cannot be reached at height=%d offset=%d",height,m.Health.Offset) }
  })
 }
}

func TestO031IndependentQuitFromRefreshHelp(t *testing.T) {
 for _, key:=range []string{"ctrl+c","q"} {
  t.Run(key,func(t *testing.T) {
   m:=openHealthScreen(t,&fakeHealth{dashboard:measuredDashboard(fiveSignals())})
   cancelled:=false
   m.Health.Refresh=&healthRefresh{Cancel:func(){cancelled=true}}
   m=press(t,m,"?")
   _,cmd:=sendKey(t,m,key)
   if !cancelled || cmd==nil { t.Errorf("%s while refresh help open: cancelled=%v quit=%v",key,cancelled,cmd!=nil) }
  })
 }
}
```

### internal/codehealth/o031_check_scratch_test.go

```go
package codehealth

import (
 "context"
 "os"
 "os/exec"
 "path/filepath"
 "testing"
)

func TestO031IndependentSavedMatrix(t *testing.T) {
 for _, outcome:=range []Outcome{OutcomeAvailable,OutcomePartial,OutcomeFailed,OutcomeTimedOut,OutcomeUnavailable,OutcomeAbsent,OutcomeUnsupported} {
  for _, required:=range []bool{false,true} {
   store,root:=dashProject(t)
   cfg:=dashConfig(t,store,allCapabilities()...)
   cfg.Capabilities[0].Required=required
   if _,err:=store.SaveConfig(cfg);err!=nil {t.Fatal(err)}
   rs:=goodResults(cfg,nil)
   if outcome==OutcomePartial {rs[0].Outcome=outcome;rs[0].Reason="Some inputs were not measured."} else if outcome!=OutcomeAvailable {rs[0]=unmeasured(rs[0],outcome)}
   saveDash(t,store,cfg,OriginManual,1,rs,dashRepo(1))
   d:=mustLoad(t,root)
   if len(d.Rows)!=5 || d.Origin!=OriginManual || d.Rows[0].Required!=required || d.Rows[0].Outcome!=outcome {t.Fatalf("wrong result %s required=%v: %+v",outcome,required,d)}
   if outcome!=OutcomeAvailable && d.Rows[0].Label==ClassificationGood {t.Fatalf("incomplete %s labelled Good",outcome)}
   assertNoBannedClaims(t,"independent result",dashboardText(d))
  }
 }
}

func TestO031IndependentRealGitFreshness(t *testing.T) {
 root:=t.TempDir()
 run:=func(args ...string){ t.Helper();cmd:=exec.Command("git",args...);cmd.Dir=root;if out,err:=cmd.CombinedOutput();err!=nil {t.Fatalf("git %v: %v %s",args,err,out)} }
 run("init","-q")
 write:=func(s string){if err:=os.WriteFile(filepath.Join(root,"source.go"),[]byte(s),0600);err!=nil{t.Fatal(err)}}
 write("first")
 run("add","source.go")
 run("-c","user.name=Check","-c","user.email=check@example.invalid","commit","-qm","first")
 obs,err:=ObserveRepository(context.Background(),GitRunner{},root,InputScope{});if err!=nil{t.Fatal(err)}
 if f:=DashboardFreshness(context.Background(),root,nil,obs.Identity);f.State!=CodeMatches || f.WorkingTreeDiffers{t.Fatalf("same: %+v",f)}
 write("second")
 if f:=DashboardFreshness(context.Background(),root,nil,obs.Identity);!f.WorkingTreeDiffers{t.Fatalf("dirty: %+v",f)}
 run("add","source.go")
 run("-c","user.name=Check","-c","user.email=check@example.invalid","commit","-qm","second")
 if f:=DashboardFreshness(context.Background(),root,nil,obs.Identity);f.State!=CodeMovedOn || f.Commits!=1{t.Fatalf("moved: %+v",f)}
}
```

### internal/buildtool/o031_check_scratch_test.go

```go
package main

import (
 "bytes"
 "os"
 "path/filepath"
 "runtime"
 "testing"
)

func TestO031IndependentInterruptedReports(t *testing.T) {
 if runtime.GOOS=="windows" {t.Skip("Unix signal reproduction; Windows requires native equivalent")}
 dir:=t.TempDir()
 bin:=filepath.Join(dir,"bin")
 if err:=os.Mkdir(bin,0700);err!=nil{t.Fatal(err)}
 script:="#!/bin/sh\nprintf '%s\\n' '{\"Action\":\"start\",\"Package\":\"fixture\"}'\nkill -TERM $$\n"
 if err:=os.WriteFile(filepath.Join(bin,"go"),[]byte(script),0700);err!=nil{t.Fatal(err)}
 t.Setenv("PATH",bin+string(os.PathListSeparator)+os.Getenv("PATH"))
 t.Chdir(dir)
 prior:=[]byte("prior complete report\n")
 if err:=os.WriteFile(filepath.Join(dir,goTestReportName),prior,0600);err!=nil{t.Fatal(err)}
 var out bytes.Buffer
 err:=runGoTestWithReports(dir,[]string{"-json","./..."},&out)
 if err==nil {t.Fatal("expected interrupted subprocess error")}
 got,readErr:=os.ReadFile(filepath.Join(dir,goTestReportName));if readErr!=nil{t.Fatal(readErr)}
 if !bytes.Equal(got,prior){t.Fatalf("interrupted tool replaced complete report with %q (error %v)",got,err)}
}
```