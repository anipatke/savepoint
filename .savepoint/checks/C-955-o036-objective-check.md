---
id: C-955
scope: {kind: objective, id: O-036}
result: NEEDS WORK
checked_by: {role: checker, session: check-o036-20261002-independent}
executed_session: o036-task-execution-prior-session
checked_at: '2026-10-02T03:01:24Z'
health_snapshot: sha256:07a6e5b0ac90d85a0096219f0d65fa472993b93aeed663889452df6598879ae1
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
  dependencies: []
issues: [I-113, I-114]
supersedes: null
---

# C-955: O-036 Full Objective Check

NEEDS WORK. This is a fresh independent review session; it did not execute T-082, T-083 or T-084. The supplied Check line selects O-036 directly. All three Tasks are owner-completed with named, dated owner waivers. No Task or Objective lifecycle or router selection was changed. Review includes uncommitted code over the named HEAD; the earlier O-035 chip/value/time edits are read as integration dependencies, not a reopened O-035 review. No prior O-036 Check or linked Issue existed at entry. Issue search found no duplicate of either symptom.

## Findings and Materiality

| Issue | Likelihood | Impact | Materiality | Recommendation |
|---|---|---|---|---|
| I-113 | Medium: report path/permissions can fail independently of snapshot storage | Medium: silent stale report can misdirect the agent after a successful refresh | Medium | Fix the warning bridge, preserving successful snapshot collection |
| I-114 | High: every red row with a missing report uses this wording | Low: ambiguous instruction; command and path otherwise work | Low | Restore the required literal in the same narrow repair |

## Acceptance Coverage

Numbering follows Objective Success Conditions and each Task's Done When order.

| Requirements | Classification | Evidence |
|---|---|---|
| O-036 SC1, T-083 DW1 save/rewrite portion | Proven | Collect saves then calls RefreshReport; TestCollectRewritesReportForTheSavedSnapshot (manual), TestRun_rewritesReportForTheSavedSnapshot (official); independent actual refresh persists one manual snapshot |
| O-036 SC2, T-083 DW2–4 | Proven | TestRunReport_rewritesDeletedReportWithoutNewSnapshot compares complete bytes and unchanged snapshot count; nothingToReportWritesNothing, unwritableReportIsAnError, namesTheRoutersObjective; parser rejection/help/dispatch tests. RunReport uses only filesystem dashboard/render/write, no tool runner |
| O-036 SC3, T-082 DW1–2 | Proven | Independent renderer matrix calls RenderReport directly; fixed brief names O-036. Renderer clones rows before sorting and performs no IO |
| O-036 SC4–6, T-082 DW3–6 | Proven | 160 independent signal × label × text-class cells, independent sort oracle, persisted newest manual scenario, TestReportShowsEverySignalWithDashboardWording, EvidenceOneManyAndNone, SeveralInstancesGetOneSectionEach, ManualNewestSnapshot, StatesOfficialAndDate, Refusals, KeepsInternalsOut. Row next step is the underlying signal action; popover-only pointer is separately described by SC7 |
| O-036 SC7, T-084 DW1 | Issue I-114 for absent-report literal; remainder Proven | Persisted red fixture with absent/present report; WatchAndUnknownKeepTheirNextStep, RedSignalPointsAtTheReport. Present-report pointer exact; absent-report required wording differs |
| O-036 SC8, T-084 DW3 | Proven | TestWhereNamesOneFileOrCountsMany; whereText preserves single path and emits only count for many |
| O-036 SC9, T-082 DW7 ignore portion | Proven | WriteReport creates local ignore once; owner-edited ignore byte-preserved independently; CreatesIgnoreFileOnce, NeverOverwritesOwnersIgnoreFile. Official collection created local ignore and Git status omits generated report. Explicit owner-ignore preservation takes precedence over adding a missing pattern to an existing authored ignore |
| O-036 SC10, T-083 DW1 warning portion | Issue I-113 | Official warning verified by TestRun_reportWriteFailureIsOnlyAWarning; actual popover refresh drops the warning. Snapshot retention independently proven |
| O-036 SC10 atomic portion, T-082 DW7–8 | Proven within supported report workflows | Same-directory writeTemp write/sync/close/chmod then rename; cleanup on error; TestWriteReportReplacesOlderReportAtomically, CreatesFilesAndIsRepeatable; independent directory/symlink/owner-ignore preservation matrix. No in-place writes |
| O-036 SC11, T-084 DW2,4–5 | Proven except assertions encode I-114 and omit I-113 | noteReport performs one existence query at dashboard load; pure rendering calls PopoverNextStep. Existing frame matrix 80x20/24/40 shows full present/missing strings without new rows; no render IO. Focused tests and full suite pass |

## Completed Matrix Classification

The frozen scope and inventory are included below verbatim. Every row was completed; both findings were collected before verdict.

| Frozen row | Classification and evidence |
|---|---|
| Renderer | Passed all 160 independent capability/label/text cells, stable ordering and non-mutation; named persisted tests cover evidence counts, multi-instance, manual/official, dates, no snapshot/config and banned internals. Invalid closed vocabularies/numeric snapshot states NA: validated storage inputs |
| Loaded projection | Passed newest manual refresh, report absent/present, exact present pointer and preserved underlying next step; I-114 absent pointer text; current-config missing row represented as Unknown via dashboardRows/noResultRow; saved evidence copied in stored order |
| Writer | Passed regular replacement, same-input byte/mtime stability, owner ignore preservation, blocked report directory, symlink refusal/target unchanged and temp cleanup. Root/unwritable permissions tests skipped as root; directory obstruction independently exercises write failure. Oversized arbitrary direct text observation below; unsupported report content outside a loaded bounded snapshot is not a supported owner workflow |
| Report command | Passed normal explicit root, selected/no Objective, absent snapshot/config, blocked write, no new snapshot and no tool; parser default-dir/help/flags/extra arguments tests pass. Pure report has no terminal-specific logic |
| Official check | Passed save/rewrite, error warning while snapshot and verdict stand, official live collection after full gate |
| Popover refresh | Passed actual production refresh snapshot retention on report failure; I-113 warning unavailable. Source traces Collect as common save path and reload completion; existing healthRefresh tests cover progress, success, failure, cancellation, repeat/no overlap |
| Popover render | Passed required frame sizes and present/missing pointer visibility, Where count/single and other labels unchanged; realcopy suite exercises normal/no-color layouts. Existing text-width tests plus independent Unicode evidence cells cover relevant text classes; report uses no truncation |
| External boundary | Existing full-suite Collect/ExecRunner tests cover configured target, ordered progress, unavailable and successful/non-success tools, malformed reports, timeout/cancel and cleanup. Health report introduces no external call. Network/redirect/connection-refusal NA to this new filesystem-only report path |

## Workflow and Side Effects

Parsing, root/schema/config/history failures terminate before new report/snapshot effects. Existing tool failures remain per-instance outcomes; cancellation returns before snapshot save. Snapshot save precedes report reload/render/ignore/temp/rename. Official report failures go to stderr without changing the verdict; manual refresh failure is lost (I-113). Explicit report command failure remains fatal to that command. Ignore link prevents overwrite; temporary files are removed; unchanged report returns without rename. Final report is exposed only by completed rename. Direct overlay probes compare saved snapshot count, bytes and untouched symlink target; existing tests compare identical-write mtime. No rendering accesses report bytes or performs filesystem work.

Atomic write/sync/close/chmod and rename failure ownership was traced through shared writeTemp and WriteReport. Directory/symlink failures are independently exercised before replacement; existing store tests and report replacement tests cover temp cleanup and complete existing bytes. A physical disk failure at every individual syscall was not injected; this is a test technique limitation, not an unverified promised behavior. Snapshot collection and report output remain separately owned effects.

## Gates and Code Health

- Fresh `make test-full`: PASS on 2026-10-02 in this Check session. It ran `go run ./internal/buildtool test -reports -json -count=1 ./...`, then build-linux/build-darwin/build-windows. Full test reports and coverage regenerated. Toolchain: Go 1.26.2 linux/amd64. Native Windows tests were not run locally; unchanged CI is the native platform gate.
- `make build`: PASS. `git diff --check`: PASS before and after collection.
- `go test -overlay /tmp/o036-overlay.json ./internal/codehealth ./internal/board/v2 -run '^TestO036Independent' -v`: PASS, recording the two actual discrepancies rather than asserting the implementation is correct.
- `go test ./cmd ./internal/codehealth ./internal/healthcheck ./internal/board/v2 . -run 'Report|Where|MeaningAndNext|HealthPopoverNext|HealthRefresh|HealthArgs|HealthDispatch|Help' -count=1 -v`: PASS. First sandbox run could not write a Go build-cache entry; approved rerun exited 0. This is an environment failure, not a test finding.
- After the full gate, `./savepoint health check O-036`: exit 0, official snapshot sha256:07a6e5b0ac90d85a0096219f0d65fa472993b93aeed663889452df6598879ae1 created. Verdict: Code Health does not block clearance. Optional dependency-vulnerability instance failed collection; tests, coverage, complexity and duplication were reported without blocking findings. Optional scanner failure is not proof of vulnerable code and does not automatically create an Issue.
- Strict load via `./savepoint resume` at entry and after new identity-bearing records. File reality: all 24 reviewed paths above exist; no unexplained phantom file. Dependencies unchanged; no source, test, fixture or gate was repaired during review.

## Design Reconciliation and Owner Validation

Design section 6 and AGENTS document the report command as human-only, newest-snapshot rewrite without tool execution, and warning-only save behavior. Report responsibility is in codehealth; healthcheck is the command behavior and cmd/main remain thin dispatch. Classification, thresholds and snapshot identity are unchanged by the report. Implementation's manual warning omission contradicts Design and Objective (I-113). ReportStep is intentionally popover-only so the agent report retains concrete signal actions rather than pointing at itself; this matches the architectural intent of an actionable brief.

All three Tasks carry owner Task-check waivers with IDs, reason, actor and time. T-084 Technical Evidence still says pending execution; independent Check evidence above supplies outcome coverage, but executor evidence should be completed when recording repair. No technical clearance, owner acceptance or Objective completion is granted. Real-terminal owner inspection remains an owner decision.

## Non-blocking Observations

A direct WriteReport call with 1,048,577 bytes succeeds once; a second identical call errors because shared readRecord caps old content at 1 MiB. This test bypasses the normal bounded snapshot/dashboard workflow; no supported saved snapshot reproduction was established. Record this API inconsistency for follow-up rather than require an unrelated size-policy expansion. Existing owner-authored local ignore files without report.md are deliberately preserved per T-082; the report may remain unignored in that explicitly preserved case.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches** — remaining outcome gaps are captured as I-113/I-114 above.
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries** — the material warning omission is captured separately.
- [ ] STYLE-07 **One source of truth** — routerObjective in healthcheck.go and selectedObjective in board/v2/io.go duplicate selected-Objective fallback parsing. Advisory only.
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs**

## Frozen Initial Scope Lock

# O-036 initial Check scope lock (2026-10-02)
1. Requirements: all 11 Objective success conditions; all Done When criteria of T-082/083/084; FS-04/05/06, ARCH-01/02/03/04, CFG-01, TEST-01/02/03/04/08/09, applicable user-file safety, and advisory STYLE-01..10. Full gate make test-full, build and diff check. Only owner closes work.
2. Surfaces: RenderReport/writeReportRow, Store.WriteReport/ensureReportIgnore, LoadDashboard/noteReport/PopoverNextStep, RefreshReport, Collect, RunReport/Run, ParseHealthArgs/RunHealth/main wiring, refreshHealth/healthRefreshCmd/completion and healthSelectedLines/popover. Files: report.go/report_test.go; scoped dashboard/copy and tests; collect and tests; storage.go support; healthcheck and tests; cmd/health and tests; main and health tests; board io/health/view/realcopy tests; AGENTS and Design command documentation. Pre-existing O-035 value-label/time/chip behavior is relied-on projection, not new classification work.
3. Effects: config/history read, provider collection (existing boundary), snapshot save, dashboard reload, render, one-time ignore creation, temp write/sync/close/chmod, rename, warning/error output, board dashboard reload. No tool execution in report rewrite. No planner/implementation changes in Check.
4. Coverage matrix frozen below. Cells include failure after snapshot saved, before ignore/report write, atomic replacement cleanup, retry and same input, latest manual/official, config changes, red/report existence and no-color frame sizing. Direct models supplement valid persisted snapshots. Unsupported arbitrary classifications/numeric models are NA (validated snapshot/config); network/redirect/auth/lifecycle terminal transitions NA (local derived renderer, no server/transaction). Existing provider behavior is classified by current collection tests only; it is unchanged.
5. Materiality: supported owner report/check/popover workflows, faithful newest snapshot content, atomic and repeatable derived writes, visible warning preserving snapshot, actionable red pointer and preservation of owner ignore content. No unrelated provider internals or new classification policy.

| Row | Applicable cells | Independent oracle / execution |
|---|---|---|
| Renderer | five capabilities × four labels; stable tied ordering; no mutation; zero/one/many evidence; zero/positive lines; multi-instance; Unicode/control paths/notes; official/manual; measured/not configured/first run; raw internals absent | overlay deterministic matrix plus named report tests |
| Loaded projection | newest official/manual; current config and missing instances; report missing/present; all labels' next steps; one/many Where; date/origin | persisted fixtures, report/copy/dashboard tests and overlay |
| Writer | missing/new project; report absent/identical/old; ignore absent/owner-edited; report directory/symlink; parent non-directory; temp cleanup; size boundary repeat; paths via explicit root | temp fixture byte/mtime/target preservation and independent overlay |
| Report command | default/explicit root; no config/no snapshot/valid latest; objective selected/none; unknown flag/extra positional/help; no tool/no snapshot; write failure | parser/healthcheck tests and overlay |
| Official check | save then rewrite; failure warning with valid snapshot/verdict; repeated report rewrite | healthcheck and collection tests; official collection after full gate |
| Popover refresh | actual refresh success and report failure; snapshot retained; warning displayed; cancel and tool failure unchanged; reloaded red pointer | actual refresh overlay and existing health tests/source trace |
| Popover render | report present/missing × 80x20/24/40; one/many Where; Watch/Unknown/Good unchanged; normal/no-color | realcopy/copy tests and Unicode sizing overlay |
| External boundary | configured/actual tool, start order, unavailable, success/failure, timeout/cancel, malformed, cleanup, partial side effect, safe errors | existing Collect/ExecRunner tests; no new report subprocess; redirect/connection NA local process |

| Order | Operation | Failure ownership / final state | Oracle |
|---|---|---|---|
| 1 | parse/root/schema/config/history read | fatal, no new artifact | parser/errors + snapshot count |
| 2 | collect/read provider, observe Git | per-instance outcome; cancellation no snapshot | collection tests |
| 3 | save immutable snapshot | fatal on failure; existing bytes preserved | saved snapshot ID/content |
| 4 | reload newest dashboard and render report | post-save warning; snapshot stands | compare independently selected newest/date |
| 5 | ensure ignore via temp/link/cleanup | post-save warning or explicit-report error; owner file untouched | byte equality and temp inventory |
| 6 | read old report; temp/write/sync/close/chmod; rename/cleanup | post-save warning or explicit-report error; old complete report retained on failure | bytes/mtime and absence of temp files |
| 7 | report success path / check verdict / refresh completion | snapshot stands; warning must be visible when report fails | stdout/stderr/board result message |

## Independent Reproduction Harness

The overlay adds temporary test files without modifying implementation or repository tests. Test bodies below retain the reproductions for a later session. Use Go overlay Replace entries for internal/codehealth/o036_check_probe_test.go and internal/board/v2/o036_check_probe_test.go pointing to temporary files containing these bodies.

### o036-codehealth-probe_test.go

```go
package codehealth
import("testing";"strings";"os";"path/filepath";"reflect")
func TestO036IndependentRendererMatrix(t *testing.T){
 labels:=[]Classification{ClassificationNeedsAttention,ClassificationWatch,ClassificationUnknown,ClassificationGood}
 texts:=[]string{"ascii.go","控制.go","e\u0301.go","v\ufe0f.go","👍🏽.go","🇦🇺.go","👩‍💻.go","tab\tfile.go"}
 for _,label:=range labels {for _,c:=range Capabilities(){for _,path:=range texts{
 row:=DashboardRow{Capability:c,CapabilityText:CapabilityText(c),Label:label,LabelText:classificationText[label],Question:"question",Value:"7",Aim:"aim",Meaning:"meaning",NextStep:"step",SignOff:"sign-off",SparkWord:"steady",Evidence:[]EvidenceRef{{Path:path,Line:7,Note:"first"},{Path:"second.go",Note:"second"}}}
 d:=Dashboard{State:DashboardMeasured,OriginText:"Official check",MeasuredText:"2 Oct 12:00",Rows:[]DashboardRow{row}}
 before:=append([]DashboardRow(nil),d.Rows...); got,err:=RenderReport(d,"O-036");if err!=nil{t.Fatal(err)}
 for _,want:=range []string{"## "+CapabilityText(c),"Label: "+classificationText[label],"- "+path+":7 first","- second.go second","savepoint health check O-036","Trend: steady"}{if !strings.Contains(got,want){t.Errorf("missing %q",want)}}
 if !reflect.DeepEqual(before,d.Rows){t.Fatal("renderer mutated dashboard")}
 }}}
 // Independent ordering oracle: rank follows the fixed objective sequence.
 d:=Dashboard{State:DashboardMeasured};for _,l:=range []Classification{ClassificationGood,ClassificationUnknown,ClassificationWatch,ClassificationNeedsAttention,ClassificationWatch}{d.Rows=append(d.Rows,DashboardRow{CapabilityText:string(l),Label:l})}
 got,_:=RenderReport(d,"O-036");heads:=headingOrder(got);want:=[]string{"needs_attention","watch","watch","unknown","good"};if !reflect.DeepEqual(heads,want){t.Fatalf("headings %v want %v",heads,want)}
 t.Log("160 label/capability/text cells and independent stable ordering passed")
}
func TestO036IndependentWriterMatrix(t *testing.T){
 for _,shape:=range []string{"directory","symlink","owner-ignore","large-repeat"}{t.Run(shape,func(t *testing.T){st,root:=newProject(t);dir:=healthPath(root);os.MkdirAll(dir,0755)
 switch shape {case "directory":os.Mkdir(filepath.Join(dir,reportFile),0755);if _,err:=st.WriteReport("new");err==nil{t.Fatal("directory accepted")}
 case "symlink": target:=filepath.Join(t.TempDir(),"owner");os.WriteFile(target,[]byte("owner"),0644);if err:=os.Symlink(target,filepath.Join(dir,reportFile));err!=nil{t.Skip(err)};if _,err:=st.WriteReport("new");err==nil{t.Fatal("symlink accepted")};b,_:=os.ReadFile(target);if string(b)!="owner"{t.Fatal("owner content changed")}
 case "owner-ignore":p:=filepath.Join(dir,gitignoreFile);os.WriteFile(p,[]byte("# owner\n"),0644);st.WriteReport("one");st.WriteReport("two");b,_:=os.ReadFile(p);if string(b)!="# owner\n"{t.Fatal("ignore changed")}
 case "large-repeat":text:=strings.Repeat("x",maxRecordBytes+1);if _,err:=st.WriteReport(text);err!=nil{t.Fatal(err)};changed,err:=st.WriteReport(text);t.Logf("oversize second identical write changed=%v error=%v",changed,err)
 }
 entries,_:=os.ReadDir(dir);for _,e:=range entries{if isTempName(e.Name()){t.Errorf("temp left %s",e.Name())}}
 })}
}
func TestO036IndependentPersistedProjection(t *testing.T){
 st,root:=dashProject(t);cfg:=dashConfig(t,st,allCapabilities()...);saveOfficialRuns(t,st,cfg,func(r *CapabilityResult){if r.Capability==CapabilityComplexity{r.Value.Number=46}})
 d:=mustLoad(t,root);row:=rowFor(t,d,CapabilityComplexity);got,_:=RenderReport(d,"O-036");t.Logf("missing-report popover=%q report=%q",row.PopoverNextStep(),row.NextStep)
 if !strings.Contains(got,"- report.out report"){t.Fatal("evidence absent")}
 st.WriteReport(got);d=mustLoad(t,root);if rowFor(t,d,CapabilityComplexity).PopoverNextStep()!="Ask your agent to investigate .savepoint/health/report.md"{t.Fatal("pointer missing")}
 saveDash(t,st,cfg,OriginManual,5,goodResults(cfg,nil),dashRepo(5));_,_,err:=RefreshReport(root,"O-036");if err!=nil{t.Fatal(err)};if !strings.Contains(readReport(t,root),"manual refresh"){t.Fatal("not newest manual")}
}
```

### o036-board-probe_test.go

```go
package v2
import("testing";"context";"os";"path/filepath";"os/exec";"github.com/opencode/savepoint/internal/codehealth")
func TestO036IndependentActualRefreshReportFailure(t *testing.T){
 root:=t.TempDir();sp:=filepath.Join(root,".savepoint");os.MkdirAll(filepath.Join(sp,"health","report.md"),0755)
 if out,err:=exec.Command("git","-C",root,"init").CombinedOutput();err!=nil{t.Fatalf("%v %s",err,out)}
 if _,err:=codehealth.NewStore(root).SaveConfig(codehealth.Config{Version:codehealth.ConfigVersion});err!=nil{t.Fatal(err)}
 err:=refreshHealth(context.Background(),sp,nil);snaps,loadErr:=codehealth.NewStore(root).LoadSnapshots();if loadErr!=nil||len(snaps)!=1{t.Fatalf("snapshot not retained %v %d",loadErr,len(snaps))}
 t.Logf("blocked report path, actual refresh returns %v; retained %d snapshot; completion carries no warning",err,len(snaps));if err!=nil{t.Fatalf("refresh falsely failed: %v",err)}
}
```

