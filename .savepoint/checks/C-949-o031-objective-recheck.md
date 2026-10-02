---
id: C-949
scope: {kind: objective, id: O-031}
result: NEEDS WORK
checked_by: {role: checker, session: recheck-o031-20261002-independent}
executed_session: user-request-repair-claude-session-01B2T7Qen3b1EsFCeVPDf92A
checked_at: '2026-10-01T21:22:45Z'
health_snapshot: sha256:ddf789716891ed8a0047ed112a55d277eac2ee544afd1e160c6170807a5d96b9
reviewed:
  base_commit: 9db1988
  head_commit: 42c8d32dd74e3a6340985b4bb20763315ca03ba7
  files:
    - internal/board/v2/health.go
    - internal/board/v2/health_view.go
    - internal/board/v2/health_test.go
    - internal/board/v2/update.go
    - internal/buildtool/main.go
    - internal/buildtool/main_test.go
    - Makefile
    - .gitignore
    - internal/codehealth/dashboard.go
    - internal/codehealth/dashboard_test.go
    - internal/codehealth/dashboard_freshness.go
    - internal/codehealth/dashboard_freshness_test.go
    - internal/codehealth/collect.go
    - internal/codehealth/collect_test.go
    - internal/board/v2/io.go
    - internal/board/v2/detail.go
    - internal/board/v2/detail_test.go
    - internal/board/v2/load.go
  dependencies:
    - internal/codehealth/repository.go
    - internal/codehealth/runner.go
    - internal/codehealth/runner_windows.go
    - .savepoint/health/config.json
issues: []
supersedes: C-948
---

# C-949: O-031 Full Objective Recheck

## Closure map

| Prior Issue | Disposition this run | Evidence |
|---|---|---|
| I-104 | Still open; repair proven in frozen Linux viewport cells | Every field/evidence reference and final history are reachable at 80x20, 80x24, 80x40; repeated up/down and resize preserve bounds and navigation. |
| I-105 | Still open; repair proven in frozen interrupt cells | Ctrl+C through Help calls the active cancel function and returns a quit command; direct quit, Esc cancellation and no-save regression tests pass. |
| I-106 | Unverified platform component; original Linux reproduction repaired | Interrupted child leaves prior reports unchanged and cleans temporary files; completed pass/fail reports still publish. No native Windows result exists for reviewed HEAD. |

The result remains NEEDS WORK solely because C-948's native Windows runtime/full-suite evidence cells C2/B3 are still Unverified. No original code defect was reproduced and no new defect was admitted. The checker cannot close Issues as verified with a NEEDS WORK Check: the workflow requires proof from a CLEAR Check. Their statuses and append-only histories remain unchanged rather than claiming verified closure prematurely.

This is the first full recheck after C-948. It is independent of both executor conversations and the repair session; this conversation only reviewed the work and wrote Check artifacts. A later targeted recheck should verify the missing native Windows evidence against this same frozen lock rather than reopen broad review.

## Scope and admission ledger

C-948's five numbered scope-lock items, finite matrix cells D1–D3/F1/U1–U3/R1–R2/C1–C2/S1/B1–B3/E1–E3 and eight workflow operations are immutable and inherited without expansion. Repairs 064eeca and 06271a1 change scrolling, interrupt dispatch and interrupted report publication. 42c8d32 commits the prior Check and owner completion records. The intervening jscpd reader/discovery change 710d4dc is outside the frozen provider-normalization perimeter: no new normalization case was used to block this recheck. It participates in the full gate and configured collection only.

All five owned Tasks T-070–T-074 were reread; they remain done with owner Task-check waivers. No Task, Objective, router, Design, user health configuration, or unrelated O-035 record was edited. The project is not a worktree lane. Existing uncommitted task-ids.yml, health configuration/history, and O-035 are owner work. The official collection adds its own permitted snapshot/report artifacts.

The following ledger was written before recheck probes:

# Recheck admission ledger — C-948 frozen cells

|Probe|Prior finding or claim|Exact frozen cell|Allowed result|
|---|---|---|---|
|Original persisted outcome matrix and Git harness|D/F proof retained|D1,D2,D3,F1|prove or in-cell failure|
|Original production board workflow|integration proof retained|U1,R1,C1|prove or in-cell failure|
|Original viewport and all-field union/resize/up-down tests|I-104 repair|U2 at 80x20/24/40 and wrapping/multiple evidence|prove or retain I-104|
|Original Ctrl+C from Help, direct quit/Esc/retry|I-105 repair|R2 and R1|prove or retain I-105; q Help convention observation only|
|Original interrupted report, pass/fail/repeat/focused/cleanup|I-106 repair|B2 and B1|prove or retain I-106|
|Saved labels/missing records and unchanged plain output|T-073 proof retained|S1,E2|prove or in-cell failure|
|NO_COLOR/dumb and eight finite text classes at frozen sizes|C-948 missing evidence|E3,U3|prove or explicit Unverified; no new case classes|
|Linux full gate and available native Windows CI evidence|C-948 missing evidence|C2,B3|prove or explicit Unverified|
|Official health check following gate|mandatory Check collection|scope lock 1 / Code Health gate|record snapshot, blocking verdict prevents CLEAR|

No new provider-reader perimeter: 710d4dc reader repair is outside the C-948 lock. No new issue admitted without an exact frozen cell. Original q-help assertion was already a nonblocking observation.


## Frozen matrix results

| Exact original cell | Current evidence | Classification |
|---|---|---|
| D1 | Full suite: TestLoadDashboardNotConfigured, FirstRun, NeverRunsAnythingOrWrites, ReportsDamagedStorage; independent production board first-run load | Proven |
| D2 | Original TestO031IndependentSavedMatrix rerun: all 14 available/partial/failed/timed-out/unavailable/absent/unsupported x required/optional cells; five rows and manual provenance; current multi-instance/missing-instance tests | Proven |
| D3 | Current persisted-summary, partial, newest-manual, comparable/reset/manual-excluded trend, bounded/newest-first history and banned-claims tests plus original independent saved matrix | Proven |
| F1 | Original TestO031IndependentRealGitFreshness passes same commit, dirty edit and one committed change; full scripted same/behind/ahead/diverged/shallow/unavailable/cancel tests pass | Proven |
| U1 | Original TestO031IndependentProductionWorkflow passes real saved config/report → H first run → R progress → exactly one manual snapshot → dashboard reload and real Git freshness. Full load-error, late-freshness, exclusivity, cursor-restoration and initial-screen tests pass | Proven |
| U2 | Original viewport fixture at 80x20/24/40 plus field-visibility union, twelve evidence rows, repeated up/down, final history and resize tests; TestHealthScrollReachesAllDetailAndHistoryContent passes | Proven; I-104 reproduction repaired |
| U3 | Existing five-signal, labels/reasons/trends, non-Good rows, Help and minimum-width tests; explicit E3 text matrix at frozen sizes below | Proven |
| R1 | Full and race refresh progress/order/second-R/success/reload/error/Esc/no-save/direct-quit tests; original production refresh workflow rerun | Proven |
| R2 | Original TestO031IndependentQuitFromRefreshHelp/ctrl+c passes; current TestHealthCtrlCFromHelpDuringRefreshCancelsBeforeQuitting passes | Proven; I-105 reproduction repaired |
| C1 | Full and race callback-order and before-first/mid/after-last no-save tests; full healthcheck TestRun_cancelledContextSavesNothing | Proven |
| C2 | Current executed-outcome/malformed/absent/partial/isolation/report-only tests, real Linux process-stop and cleanup; race suite passes. Native Windows collection/process-stop execution not obtained | Linux Proven; native Windows Unverified |
| S1 | Current stored-official/missing/no-field/no-directory snapshot label and Objective/Task detail tests; unchanged summary/load source traces; full non-TTY tests | Proven |
| B1 | Current pass/fail/focused/no-temp-files/timing/exit tests; fresh make test-full writes/replaces both conventional root reports; Makefile and gitignore inspected | Proven on Linux |
| B2 | Original TestO031IndependentInterruptedReports passes. Current TestRunGoTestWithReportsKeepsCompleteReportsWhenInterrupted additionally preserves both sentinel reports and leaves no temp files | Proven for original Linux interruption; native equivalent remains with B3 |
| B3 | Fresh full host suite and all six Linux/Darwin/Windows cross-build targets pass; GitHub exact-HEAD query returns [] | Linux/cross-build Proven; native Windows Unverified |
| E1 | No added numeric validator, remote/server/browser/redirect surface; typed projection APIs remain behind existing strict config/snapshot readers | N/A, original scope reason retained |
| E2 | Original redirected/non-TTY surface remains unchanged; full board plain-output tests pass | Proven unchanged; refresh UI N/A in redirected mode |
| E3 | Independent 72-cell text/view matrix: ASCII, control/tab/newline/CR, combining marks, variation selector, wide characters, emoji modifier, regional flag, joined emoji x normal/NO_COLOR/TERM=dumb x 80x20/24/40; every intermediate detail fits dimensions, all fields/references and bounded history reachable | Proven for frozen text classes and layout. Uses established Lip Gloss cell-width oracle; no new width algorithm or input class |

All initial matrix rows are classified. No new axis, dependency layer, public value or acceptance meaning was added. The finite local boundary behavior and eight-operation workflow lock were rechecked by the full regression suite, original independent fixtures, source tracing and repeated report publication. No network redirect/server cells were invented. No fresh out-of-scope finding affects the verdict.

## Reproduction oracle correction

The original viewport harness has an overly strong incidental assertion: after 60 Down keys it expects a reference at the start of a long evidence note to remain on the final screen, and always expects Offset > 0 even if the whole body fits. The repaired UI legitimately scrolls that earlier reference offscreen and uses Offset=0 when content fits. Restoring the fixture unchanged produced offset=20 at height 20 and offset=0 at height 40, despite reachable content. This is a Check-harness oracle error, not a remaining product defect.

The same frozen fixture and sizes were rerun with the acceptance-aligned oracle: accumulate screens while scrolling and require the reference to appear somewhere. Then the explicit twelve-reference/all-field matrix verified both start and end content, upper/lower navigation bounds and resize. This corrects how reachability is measured; it does not change the frozen scope or waive a criterion. The original q-from-Help assertion remains the explicitly nonblocking C-948 observation and was not admitted as an Issue.

## Gate, platform and health evidence

- Toolchain: go1.26.2 linux/amd64; reviewed HEAD 42c8d32dd74e3a6340985b4bb20763315ca03ba7. Fresh run on 2026-10-01 UTC / 2026-10-02 Australia/Sydney.
- `make build`: exit 0.
- Fresh `make test-full`: exit 0 after all temporary Check test files were removed. Full uncached Go suite and all declared cross-platform builds passed; the gate rewrote go-test.json and coverage.out. No previous full result was reused.
- `go test -race ./internal/board/v2 ./internal/codehealth ./internal/buildtool -run 'Health|Dashboard|Collect|RunGoTest|ExecRunnerStops' -count=1`: exit 0 for all three packages.
- Original persisted/dashboard/Git and interrupted-report harnesses: exit 0. Original production board workflow and Ctrl+C-from-Help subtest: exit 0. Acceptance-aligned original viewport fixture: exit 0.
- `go test ./internal/board/v2 -run TestO031RecheckTextAndViewportMatrix -count=1`: exit 0, 72 finite cases, 19.537 seconds.
- `git diff --check`: exit 0 before record writing; final validation repeats it.
- Read-only `gh run list --limit 5 --json databaseId,headSha,status,conclusion,workflowName`: required sandbox network escalation, then succeeded. Latest listed CI run 36296646883 is for 9beffa5b01aa4e92425498154f5055a8ca642776, not reviewed HEAD. `gh run list --commit 42c8d32dd74e3a6340985b4bb20763315ca03ba7 --limit 10 --json databaseId,headSha,status,conclusion,workflowName`: exit 0, output []. No native Windows runtime/CI proof for this work was found. No CI workflow was triggered and no push occurred.
- After the successful full gate, `./savepoint health check O-031`: exit 0; created official snapshot sha256:ddf789716891ed8a0047ed112a55d277eac2ee544afd1e160c6170807a5d96b9. Printed verdict: Code Health does not block clearance. Tests, coverage, complexity and duplication had no blocking finding. Optional dependency vulnerability collection failed; it is incomplete measurement, not bad code, and does not automatically create an Issue. The owner's new configuration was read as-is; setup was not run.

CFG-03 and T-074 criterion 6 explicitly require native Windows full-suite evidence. C-948 also left T-071's native process-stop component unverified. Cross-compilation proves builds only. These exact C2/B3 evidence gaps prevent CLEAR despite passing local repairs, so this Check does not manufacture clearance or infer a platform waiver.

## Acceptance, guardrails and Design reconciliation

T-070 criteria 1–7 remain Proven; T-073 criteria 1–5 remain Proven. T-072 technical criteria now have proof including former U2/R2 failures and the formerly unverified text-class matrix. T-071 Linux callback/cancellation/diagnostic/process-stop criteria pass, with native Windows runtime evidence Unverified. T-074 criteria 1–5 pass on Linux including interrupted report preservation; criterion 6 native Windows CI remains Unverified. Every O-031 outcome has supporting implementation evidence, but the required platform gate has not been established.

ARCH-02 is maintained: render/reducer calculate layout and issue commands; persisted reads and subprocess work stay in commands. ARCH-03 project-root conversion is unchanged. DATA-03 diagnostics and error retention pass. FS-05 uses filepath for reports; temporary-file rename and interrupted sentinel preservation prove no replacement with the incomplete Linux stream. TEST-01/02/04/08/09 have current local outcome/failure/temp-project/full-gate/owner-waiver evidence. CFG-02/03 Windows execution is the explicit missing evidence, not silently treated as Linux behavior.

Design remains unchanged from C-948. The same nonblocking reconciliation observations apply: section 8 does not yet describe H/dashboard/refresh/cancel/compact Health, and section 13 omits report artifacts. No Design remediation occurred in this Check. No new architecture or classification decision is introduced by the repairs.

## Materiality and remaining action

No new reproducible code Issue was found. The previous High/Medium/Medium defect priorities are not repeated as active defect failures: their local reproductions pass. The remaining evidence gap has Medium materiality: Windows is explicitly supported and report/cancellation behavior crosses platform process semantics, so a native successful run is needed before technical clearance. Obtain native Windows evidence for the reviewed implementation, then run the permitted targeted recheck using C-948's lock. Do not restart implementation or retreat completed Tasks merely to obtain that evidence.

No materiality action is required for advisory copy/style observations or the optional health-tool collection failure. No new Issue ID was allocated; existing I-106 already names the native-equivalent proof, while original C-948 C2/B3 identify the wider Windows gate cells.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches** — regressions now cover scrolling, Help interrupt and interrupted report preservation on Linux; native evidence remains separately Unverified.
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth**
- [x] STYLE-08 **Comments explain why**
- [ ] STYLE-09 **Content lives in data** — unchanged inline detail labels/history/evidence copy in health_view.go; original advisory observation only.
- [x] STYLE-10 **Small diffs**

## Owner validation and handoff

T-072 still declares owner_validation.required, with no recorded accepted Check in its evidence. The real-terminal User Check remains an owner decision after technical clearance. This run does not record acceptance or close a Task/Objective. Router and R-007 selection are unchanged. Issues remain open pending CLEAR proof under the skill's closure boundary. C-948 is untouched; this immutable record supersedes it.

## Recheck harnesses

Temporary files were restored from C-948, their visibility oracle corrected as described, and an E3/U2 matrix added only for already frozen cells. They were deliberately removed before the final full gate. Their exact source is preserved below.

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
   seen:=before
   for i:=0;i<60;i++ { m=press(t,m,"down");seen+="\n"+screen(m) }
   if !strings.Contains(seen,"FINAL_EVIDENCE_SENTINEL") {
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

### internal/board/v2/o031_recheck_scratch_test.go

```go
package v2

import (
 "fmt"
 "strings"
 "testing"
 tea "github.com/charmbracelet/bubbletea"
 "github.com/charmbracelet/lipgloss"
 "github.com/opencode/savepoint/internal/codehealth"
)

func TestO031RecheckTextAndViewportMatrix(t *testing.T) {
 classes:=[]string{"ASCII text", "tab\ttext\nnext\rline", "e\u0301", "\u2764\ufe0f", "漢字", "👍🏽", "🇦🇺", "👩‍💻"}
 for _,mode:=range []string{"normal","no_color","dumb"} {
  t.Run(mode,func(t *testing.T){
   if mode=="no_color" {t.Setenv("NO_COLOR","1")};if mode=="dumb" {t.Setenv("TERM","dumb")}
   for idx,text:=range classes {for _,height:=range []int{20,24,40} {
    t.Run(fmt.Sprintf("class%d_height%d",idx,height),func(t *testing.T){
     rows:=fiveSignals()
     rows[0].Explanation=strings.Repeat(text+" ",35)
     rows[0].Evidence=nil
     for j:=0;j<12;j++ {rows[0].Evidence=append(rows[0].Evidence,codehealth.EvidenceRef{Path:fmt.Sprintf("source%02d.go",j),Note:text})}
     d:=measuredDashboard(rows)
     for j:=0;j<10;j++ {d.History=append(d.History,codehealth.DashboardHistoryEntry{CreatedAt:fmt.Sprintf("history%02d",j),OriginText:"Manual refresh",OverallText:"Watch"})}
     m:=openHealthScreenAt(t,&fakeHealth{dashboard:d},80,height)
     m=press(t,m,"enter")
     seen:=screen(m)
     for j:=0;j<250;j++ {
      m=press(t,m,"down");view:=screen(m);seen+="\n"+view
      for _,line:=range strings.Split(view,"\n") {if lipgloss.Width(line)>80 {t.Fatalf("width overflow: %q",line)}}
      if strings.Count(view,"\n")+1>height {t.Fatal("height overflow")}
     }
     for _,field:=range []string{"Required: required","Provider:","Scope:","Measured:","Snapshot:","AFFECTED AREAS","source00.go","source11.go"} {if !strings.Contains(seen,field){t.Errorf("unreachable %s",field)}}
     next,_:=m.Update(tea.WindowSizeMsg{Width:80,Height:20});m=next.(Model)
     for j:=0;j<250;j++ {m=press(t,m,"up")};if m.Health.Offset!=0{t.Fatal("detail lower bound")}
     m=press(t,m,"esc")
     for j:=0;j<250;j++ {m=press(t,m,"down")}
     if !strings.Contains(screen(m),"history09"){t.Fatal("final history unreachable")}
     for j:=0;j<250;j++ {m=press(t,m,"up")}
     if m.Health.Cursor!=0 || m.Health.Scrolled{t.Fatal("overview return bound")}
     m=press(t,m,"esc");if m.Health!=nil{t.Fatal("close failed")}
    })
   }}
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