---
id: C-956
scope: {kind: objective, id: O-036}
result: CLEAR
checked_by: {role: checker, session: recheck-o036-20261002-independent}
executed_session: o036-issue-repair-separate-session
checked_at: '2026-10-02T03:12:42Z'
health_snapshot: sha256:bcc7043ac9e301037c48d37851b740eaec8e79a9e6e2f8de32ef0f39f0cacbb6
reviewed:
  head_commit: d3edcc4c90af1e53cc03cb2ffca84aafba737334
  files:
    - internal/codehealth/report.go
    - internal/codehealth/report_test.go
    - internal/codehealth/collect.go
    - internal/codehealth/collect_test.go
    - internal/codehealth/dashboard.go
    - internal/codehealth/dashboard_copy.go
    - internal/codehealth/dashboard_copy_test.go
    - internal/codehealth/dashboard_test.go
    - internal/codehealth/storage.go
    - internal/codehealth/storage_test.go
    - internal/codehealth/snapshot.go
    - internal/healthcheck/healthcheck.go
    - internal/healthcheck/healthcheck_test.go
    - cmd/health.go
    - cmd/health_test.go
    - main.go
    - main_health_test.go
    - internal/board/v2/io.go
    - internal/board/v2/health.go
    - internal/board/v2/health_view.go
    - internal/board/v2/health_realcopy_test.go
    - internal/board/v2/health_test.go
    - AGENTS.md
    - .savepoint/Design.md
    - internal/board/v2/health_refresh_test.go
    - internal/board/v2/health_popover_test.go
  dependencies: []
issues: []
supersedes: C-955
---

# C-956: O-036 Full Objective Recheck

## Closure Map

| Prior Issue | Result | Independent proof |
|---|---|---|
| I-113 | Closed as verified by C-956 | Actual healthRefreshCmd invokes production refreshHealth in a temporary Git project with successful and blocked report writes; warning is separate from fatal error, snapshot remains saved, and warning survives model reload at 80x20/24/40 |
| I-114 | Closed as verified by C-956 | Persisted red complexity dashboard with no report returns the exact investigate-it instruction; present/missing prompt survives all original frame sizes with absent/thin/restarted history-note inputs |

CLEAR. This conversation performed the initial Check and recheck only; it did not implement the Objective or these repairs. The user requested a recheck of O-036. Every original C-955 scope-lock row was rerun without expanding the blocking perimeter. C-955 remains immutable. All three owner-completed Tasks and their explicit waivers remain unchanged. No owner acceptance or Objective closure is recorded.

## Scope and Admission Ledger

The immutable authority is C-955's Frozen Initial Scope Lock, coverage matrix and seven-operation side-effect inventory. This recheck uses their original requirements, classifications of unsupported inputs, and materiality boundary. Changed repair surfaces are refreshHealth/healthRefreshCmd, HealthFuncs/completion/model warning, status rendering, next-step wrapping, and the missing-report copy and its tests. The report, collection, official command and snapshot behavior were revalidated as cross-Task integration.

# C-956 admission ledger
| Item | Prior issue/claim | Exact C-955 frozen cell | Allowed result |
|---|---|---|---|
| renderer 160 cells + sort/mutation | unchanged behavior | Renderer five capabilities × four labels, Unicode/control evidence, stable order, no mutation | Proven or in-scope regression |
| writer preservation/repeat/error matrix | unchanged behavior | Writer directory/symlink/owner-ignore/same input/temp cleanup | Proven or in-scope regression; oversize remains observation |
| persisted missing/present prompt and manual newest | I-114 | Loaded projection report missing/present, all-label next step, newest manual | Proven or I-114 remains open |
| actual refresh success/blocked report + warning bridge/reload | I-113 | Popover refresh actual refresh success/report failure, snapshot retained, warning displayed, reloaded pointer | Proven or I-113 remains open |
| pointer full frame across sizes/history notes | I-114 remediation wrapping | Popover render present/missing × 80x20/24/40; unchanged other rows; normal/no-color | Proven or regression inside original size cell |
| report/check/parser and collection boundary tests | prior proven criteria | Report command, Official check, External boundary rows | Proven or frozen-cell regression |
| full gate and official health | required gate | Scope lock requirement 1 full gate; official check row | Proven or gate blocker |

## Completed Coverage Matrix

| Original frozen row | Result | Evidence |
|---|---|---|
| Renderer | Proven | Original independent 160 capability × label × text-class cells rerun, stable order oracle and no mutation pass. Existing Report tests rerun for all five signals, evidence zero/one/many/line-less, multi-instance, official/manual, no config/snapshot, banned internals and date |
| Loaded projection | Proven | Independent persisted red/no-report fixture asserts required literal; report-present pointer unchanged; newest manual fixture produces manual report; RedSignalPointsAtTheReport and WatchAndUnknownKeepTheirNextStep pass |
| Writer | Proven within original supported workflow boundary | Original independent directory/symlink/owner-ignore/temp-cleanup matrix rerun; same-input/replacement/ignore preservation tests pass. Unbounded arbitrary text observation remains nonblocking exactly as C-955 |
| Report command | Proven | Focused parser/help/dispatch and RunReport tests pass, including missing config/snapshot, rewrite without new snapshot/tool, explicit root, selected/no Objective, failed write |
| Official check | Proven | Report rewrite and warning-only failure tests pass; fresh official collection after successful full gate saves named snapshot and does not block clearance |
| Popover refresh | Proven | Independent command bridge drives real refreshHealth for successful and blocked report path; exactly one manual snapshot saved in either case; ReportErr only for blocked path, Err nil; model reload retains visible warning at all three sizes. Existing success/progress/fatal-failure/cancellation/repeat tests pass |
| Popover render | Proven | 18 independent prompt × history-note × original frame-size cells pass. All five signal aims remain visible, bounding dimensions hold, and the exact missing-report prompt wraps without truncation. Current realcopy/frame/no-color tests pass; Where single/count and other-label wording unchanged |
| External boundary | Proven within original unchanged boundary | Fresh full gate reruns existing collection/runner unavailable/success/non-success/malformed/timeout/cancel/cleanup tests. Report rewrite still invokes no tool. Original network/redirect/connection NA reasons unchanged |

Acceptance reconciliation: O-036 SC1–11 and T-082/T-083/T-084 Done When criteria are Proven within the same interpretations recorded in C-955. SC7/T-084 DW1 is now exact, and SC10/T-083 DW1 now reports refresh write failures. No Issue or material unverified criterion remains; no materiality action is required.

## Repair Evidence and Side Effects

refreshHealth now returns reportErr separately from err; healthRefreshCmd carries both in healthRefreshDoneMsg. applyHealth reloads results after successful collection and retains ReportWarning separately from transient Notice; healthStatusLine displays it. Starting a subsequent refresh clears the old warning. Fatal refresh failures and cancellation still take their original paths and do not masquerade as report warnings.

The independent harness invokes the actual production command bridge with an empty valid configuration in a temporary Git repository, so no tool/network is needed. In the blocked case report.md is a directory: collection still saves a valid manual snapshot, the completion contains a report error and no fatal error, and the model renders the warning after load/freshness completion. In the normal case a manual report is written and no warning appears. Named repair tests TestRefreshHealthReportsBlockedReportAsWarning and TestHealthRefreshReportWarningSurvivesReload corroborate those outcomes.

The missing-report copy now matches T-084 literally. healthSelectedLines wraps Next and merges sign-off/Where to preserve it within five detail lines, with history-note fallback when necessary. Independent frame tests check full concatenated text, all five signal aims and frame bounds for both prompts with empty/thin/restarted notes at 80x20/24/40. Renderer and writer are unchanged; original atomic rename, ignore preservation, snapshot-before-report ordering and warning/error ownership remain valid. No code or tests were repaired by this checker.

The first new frame harness mistakenly compared wrapped text including box borders; that assertion failed despite complete prompt visibility. The oracle was corrected to extract frame interior text and rerun the same cells. No implementation changed. Sandbox Go-cache write errors were followed by approved reruns; they are environment failures, not findings.

## Verification

- Fresh make test-full: exit 0, Go 1.26.2 linux/amd64, 2026-10-02 in this recheck. Entire Go suite ran with -count=1 and reports/coverage, then Linux/macOS/Windows builds passed. Native Windows execution remains CI evidence, as in C-955.
- make build: exit 0. git diff --check: exit 0.
- go test -overlay /tmp/o036-recheck-overlay.json ./internal/codehealth ./internal/board/v2 -run '^TestO036Independent' -count=1 -v: exit 0 after corrected frame oracle and approved cache access; all five independent test functions pass.
- go test ./cmd ./internal/codehealth ./internal/healthcheck ./internal/board/v2 . -run 'Report|Where|MeaningAndNext|HealthPopover|HealthRefresh|RefreshHealth|HealthArgs|HealthDispatch|Help' -count=1: exit 0 for all five packages.
- ./savepoint health check O-036 after full gate: exit 0, official snapshot sha256:bcc7043ac9e301037c48d37851b740eaec8e79a9e6e2f8de32ef0f39f0cacbb6 created. Code Health does not block clearance. Optional vulnerability scanner failed collection; other four signals have no blocking finding. Failed collection is not proof of bad code and creates no automatic Issue.
- File reality: all prior reviewed paths and the new repair test exist. No phantom path. Scope includes the uncommitted working tree over the named HEAD. Final strict loading uses ./savepoint resume after this Check and verified Issue resolutions.

## Design, Guardrails and Owner Handoff

The repaired warning behavior now matches Design section 6, O-036 SC10 and T-083; cmd/main remain thin, report rendering performs no IO, and explicit command/update paths own writing. FS-04/05/06, ARCH-01/02/03/04, CFG-01 and TEST-01/02/03/04/08/09 are satisfied within the original scope. Owner ignore content and symlink target content remain byte-identical under probes. No classification, thresholds or snapshot identity changed.

Task-check waivers continue to name each Task, owner actor, reason and time. T-084's pending Technical Evidence text and prior oversized-direct-text/owner-ignore observations remain metadata/advisory follow-ups, not new blocking requirements in this recheck. No new Issue was opened. All owned Tasks are done; the Full Objective Check is CLEAR and the two linked material findings are verified. O-036 is ready for the owner's closure decision. This Check does not close the Objective or record owner validation.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches**
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [ ] STYLE-07 **One source of truth** — routerObjective in healthcheck.go and selectedObjective in board/v2/io.go still duplicate fallback parsing; unchanged advisory observation from C-955.
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs**

## Independent Repair Harness

The original codehealth harness is C-955's appendix, with the missing-report log replaced by an assertion of the required investigate-it literal. The adapted actual bridge/frame harness below runs as a temporary Go overlay and changes no repository source or test.

```go
package v2
import("testing";"context";"os";"path/filepath";"os/exec";"strings";"github.com/opencode/savepoint/internal/codehealth"; xansi "github.com/charmbracelet/x/ansi";tea "github.com/charmbracelet/bubbletea")
func TestO036IndependentRefreshWarningBridge(t *testing.T){for _,blocked:=range []bool{false,true}{
 root:=t.TempDir();sp:=filepath.Join(root,".savepoint");os.MkdirAll(sp,0755);if blocked{os.MkdirAll(filepath.Join(sp,"health","report.md"),0755)}
 if out,err:=exec.Command("git","-C",root,"init").CombinedOutput();err!=nil{t.Fatalf("%v %s",err,out)}
 if _,err:=codehealth.NewStore(root).SaveConfig(codehealth.Config{Version:codehealth.ConfigVersion});err!=nil{t.Fatal(err)}
 // Drive the real command bridge, not just its refresh callback.
 events:=make(chan tea.Msg,4);result:=healthRefreshCmd(context.Background(),HealthFuncs{Refresh:refreshHealth},sp,events)();if result!=nil{t.Fatal("unexpected command result")};done:= (<-events).(healthRefreshDoneMsg)
 if done.Err!=nil||(done.ReportErr!=nil)!=blocked{t.Fatalf("blocked %v: done %+v",blocked,done)}
 snaps,err:=codehealth.NewStore(root).LoadSnapshots();if err!=nil||len(snaps)!=1{t.Fatalf("snapshot %d err %v",len(snaps),err)}
 if !blocked{b,err:=os.ReadFile(filepath.Join(sp,"health","report.md"));if err!=nil||!strings.Contains(string(b),"manual refresh"){t.Fatal("successful actual refresh did not write newest manual report")}}
 for _,size:=range [][2]int{{80,20},{80,24},{80,40}}{f:= &fakeHealth{dashboard:measuredDashboard(realCopyRows())};m:=openHealthScreenAt(t,f,size[0],size[1]);next,cmd:=m.Update(done);m=settle(t,next.(Model),cmd);view:=xansi.Strip(screen(m));popoverBounds(t,view,size[0],size[1]);if strings.Contains(view,"report could not be updated")!=blocked{t.Fatalf("warning display blocked=%v size=%v: %s",blocked,size,view)};if blocked&&m.Health.ReportWarning==""{t.Fatal("warning lost after reload")}}
 }}
func TestO036IndependentMissingPromptFrameMatrix(t *testing.T){for _,step:=range []string{"Ask your agent to investigate .savepoint/health/report.md","Run savepoint health report, then ask your agent to investigate it."}{for _,note:=range []string{"","not enough history yet (2 of 3)","early restarted because settings changed"}{for _,size:=range [][2]int{{80,20},{80,24},{80,40}}{
 rows:=realCopyRows();rows[0].ReportStep=step;rows[0].SparkNote=note;rows[0].Where="13 files";m:=openHealthScreenAt(t,&fakeHealth{dashboard:measuredDashboard(rows)},size[0],size[1]);view:=xansi.Strip(screen(m));popoverBounds(t,view,size[0],size[1]);var inside []string;for _,line:=range strings.Split(view,"\n"){a,b:=strings.Index(line,"│"),strings.LastIndex(line,"│");if a>=0&&b>a{inside=append(inside,strings.ReplaceAll(line[a+len("│"):b],"│"," "))}};flat:=strings.Join(strings.Fields(strings.Join(inside," "))," ");if !strings.Contains(flat,"Next: "+step){t.Fatalf("prompt lost size=%v note=%q: %s",size,note,view)};for _,row:=range rows{if !strings.Contains(view,row.Aim){t.Fatalf("other signal row moved/lost %q",row.Aim)}}
 }}}}
```
