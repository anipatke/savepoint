---
id: C-941
scope: {kind: objective, id: O-028}
result: NEEDS WORK
checked_by: {role: checker, session: o028-check-20261001}
executed_session: o028-executors-20261001
checked_at: '2026-10-01T08:38:18Z'
reviewed:
  base_commit: 64eb046
  head_commit: 919544c0da6e53f73e2b61b01fcb7e24098b72fa
  files:
    - .savepoint/Design.md
    - AGENTS.md
    - agent-skills/savepoint-design/SKILL.md
    - cmd/health.go
    - cmd/health_test.go
    - internal/codehealth/classification.go
    - internal/codehealth/collect.go
    - internal/codehealth/collect_test.go
    - internal/codehealth/config.go
    - internal/codehealth/discovery.go
    - internal/codehealth/discovery_catalogue.go
    - internal/codehealth/discovery_test.go
    - internal/codehealth/identity.go
    - internal/codehealth/instances_test.go
    - internal/codehealth/model.go
    - internal/codehealth/repository.go
    - internal/codehealth/runner.go
    - internal/codehealth/runner_alive_unix_test.go
    - internal/codehealth/runner_alive_windows_test.go
    - internal/codehealth/runner_test.go
    - internal/codehealth/runner_unix.go
    - internal/codehealth/runner_windows.go
    - internal/codehealth/setup.go
    - internal/codehealth/setup_test.go
    - internal/codehealth/snapshot.go
    - main.go
    - main_health_test.go
    - templates/project-v2/AGENTS.md
    - templates/project-v2/agent-skills/savepoint-design/SKILL.md
  dependencies: []
issues: [I-092, I-093, I-094, I-095, I-096]
supersedes: null
---

# C-941: O-028 Full Objective Check

## Result and authority

**NEEDS WORK.** Five in-scope Issues are recorded together below. This conversation is independent from the executors: it did not build O-028. Executor provenance is the aggregate O-028 task execution on 2026-10-01, including the T-057 executor explicitly recorded in that Task. The pasted Next selected this Objective directly; router state is still `task` with O-028 and no Task. The runtime Next projection, rather than that state label, owns the requested action. No router or Task status was changed.

Every owned Task (T-055–T-058) is done and carries the owner's board-recorded optional Task-check waiver naming Task, reason, actor and time. Those satisfy TEST-09 but do not clear these implementation/integration Issues. There are no prior O-028 Checks or matching prior Issues; search preceded allocation. C-940 covers O-027 dependency history only and is not superseded. T-058 declares owner validation and its recorded scratch-project walkthrough remains outstanding; this Check records no owner acceptance.

## Acceptance coverage

Criteria are numbered in their authored order; no interpretation was changed to match implementation.

| Scope | Proven | Issue / remaining evidence |
| --- | --- | --- |
| O-028 Success Conditions | 1 bounded confirmation-ready suggestions; 2 exact config vs plain Design intention boundary; 3 explicit setup/reconcile and no startup scanning; 5 distinct outcomes at service seam; 6 report-only reuse, no tests run; 7 exclusions/fingerprint; 8 separate scoped instance values | 4 safe execution: Windows descendants survive (I-095). End-to-end promised discovery/configuration/execution outcome fails I-092. Rename continuity and supported capacity fail I-093/I-094. |
| T-055 Done When 1–5 | 1 named distinct scoped instances and duplicate validation; 3 timeouts/defaults/bounds; 4 unchanged version/strict decode; 2 name serialization and unnamed identity vectors | 2 rename continuity and 5 collected history separation/continuity: I-093. Independent Assess filtering passes; Collect's upstream filtering does not. |
| T-056 Done When 1–8 | Deterministic valid individual proposals; supported stacks/settings/known reports; explicit gaps; argument vectors and exclusion translation; report-only gate instructions; catalogue exclusions; bounded depth/files/bytes; read-only and symlink handling | Cross-Task report-output contract is I-092 (all three executed proposal forms). Discovery-only tests pass. |
| T-057 Done When 1–7 | Sequential order, outcomes, report-only freshness/no execution, reader seam/invalid-result containment, per-instance required flags, preservation/no pruning/origin retention; snapshot assembly for ordinary configurations | 1/6 saved snapshot for valid many-instance config: I-094. 2 process children: I-095. 5/6 setup's reports never reach reader: I-092. History assembly: I-093. |
| T-058 Done When 1–6 | Read-only preview; apply config-only/repeat; confirmed field preservation/reconcile/missing-tool and input attention; init preview/failure warning wiring; untouched board/resume/doctor/upgrade paths; design skill/template equality and human-only CLI documentation | Required architecture-map reconciliation remains false (I-096). Owner User Check is unverified owner validation, distinct from these technical criteria. |

## Completed frozen matrix

The scope lock below was written before the first independent probe and was not widened afterward. P = passed, I = Issue, U = unverified platform capability, N/A = explicitly not applicable. Every row has a classification; Issue discovery did not stop remaining checks. Existing regression tests support coverage; the embedded independent harness additionally varies exact limits and spans integration boundaries.

| Cell | Classification and evidence |
| --- | --- |
| M1 | P: ConfigInstances, ConfigNameAndRequiredRoundTrip, EffectiveTimeout, EveryExecutedProviderHasADefaultTimeout, TimeoutBounds; independent -1/0/1/59/60/120/599/600/601 and wrong-type/unknown/missing JSON. Existing model tests cover nonfinite thresholds. |
| M2 | P: named snapshot roundtrip, empty name compatibility, order permutation, identity vectors, changed scope/exclusion and separate Assess values. I: integrated rename → saved history, I-093. |
| M3 | P: Go/Vitest/V8/Python/unsupported and polyglot discovery; PATH gaps, report presence; independent two-component setup configuration validates. |
| M4 | P: DiscoverSearchDepthIsBounded, ReadLimitIsAGap, escaping-symlink/malformed/cancel/deterministic/read-only tests; independent size cap-1/cap/cap+1 and excluded node_modules component. Directory walk cap traced; bounded enumeration stops at the cap. |
| M5 | I: independent delivery probe for Lizard, jscpd, OSV file outputs all reads stdout (I-092). P: collector's manually configured stdout/placeholder report, absent report and malformed reader paths; discovery/collection isolated tests. Cleanup passes for supported placeholder form; setup form leaves persistent project reports, included in I-092. |
| M6 | P: ExecRunnerResults for success/nonzero/unavailable/not runnable/env/direct argv; limited buffers accept oversized output without unbounded growth; independent bounded file report at cap-1/cap/cap+1. sanitizeLine bounds/UTF-8/control handling traced. Shell metacharacter test is weak (helper ignores extras), but actual exec.CommandContext proves direct argv with no shell. |
| M7 | P: Unix deadline/cancel/pre-cancel/process-group child death. I: native Windows child liveness after cancellation (I-095). Windows code path/test skip reviewed; native Windows probe reproduces, so this cell is no longer U. Native Windows complete suite was not run here; full gate's cross-builds passed and no claim of native-suite clearance is made. |
| M8 | P: every distinct executed outcome, OSV findings exit, cancelled remaining instances, required flags, unsupported never executes, reader malformed/wrong unit/partial and intact neighboring success. |
| M9 | P: report-only fresh/stale/absent/nonregular/escaping-link/oversize, no Runner calls; independent exact-equal input/report mtime is fresh. Observation failure is visibly unknown freshness by source trace; cancellation maps to cancelled. |
| M10 | P: history failure before execution, invalid config/origin refusal, ordinary immutable Store writes/conflicts/repeat, origin retention and collection never pruning. I: valid 32-instance config executes 32 tools then refuses 36 results (I-094). Store's existing failure/partial-write tests support relied-on persistence; arbitrary OS fault timing during every syscall was not injected. |
| M11 | P: preview bytes/mtimes, apply/reapply, edited fields kept when adding tools, untouched handformatted bytes with no new proposals, missing tool/scope retained, invalid config/conflict/non-project clear failure. Independent polyglot preview+apply+repeat also passes. |
| M12 | P: cmd parser and main health/init integration tests, explicit setup existing-project path, unchanged board/resume/doctor/upgrade source paths; canonical/scaffold design skill byte comparison. I: guide says no collection (I-096). Design CLI table matches the shipped command. |

Network redirect/server/retry axes are N/A because no network client is implemented here; OSV disclosure is verified. Interactivity/color/Unicode width are N/A for the scoped plain-text service; bounded control/UTF-8 stderr is in M6. Production report normalization belongs to O-029, so fake readers are an intended seam, not missing work. Tool installation and collection CLI/TUI are explicitly outside scope.

## Issues and materiality

| Issue | Likelihood | Impact | Materiality | Recommendation |
| --- | --- | --- | --- | --- |
| I-092 | High: every generated executed-tool proposal | High: actual tool reports ignored, core setup-to-run outcome unusable | High | Fix report-delivery contract now across discovery and collection. |
| I-093 | Medium: instance rename after history exists | Medium: trend/baseline silently reset | Medium | Fix history selection and integrated rename regression. |
| I-094 | Low: many-instance monorepo near supported capacity | Medium: tools all execute but no snapshot saved | Medium | Align capacity and fail invalid configurations before running. |
| I-095 | Medium: cancellation of a tool with children on Windows | Medium: orphan analysis continues after cancellation | Medium | Implement native descendant cleanup and valid native test. |
| I-096 | High: any agent following the map | Low: inaccurate architecture guidance | Low | Reconcile in the same narrow repair. |

Each Issue names requirements, supported reproduction, exact source location, and inadequate existing tests. No exception was inferred from the executor's disclosed limitations. Owner-only completion and Task statuses are preserved. Remediation belongs to executor/planner; this checker repaired no implementation or Design.

## Evidence and gates

- Fresh `make test-full`, 2026-10-01 approximately 08:03–08:04 UTC, Go 1.26.2 linux/amd64: exit 0; uncached complete package tests and linux/darwin/windows builds. No reuse of executor results. The review started at clean head 919544c; before and after scratch probes, code/tests/fixtures/dependencies/gate definitions are unchanged. Scratch harness was removed from source before handoff. The fresh gate ran before the scratch harness was added, so it exercised the exact final implementation.
- Focused named schema/discovery/runner/collect/setup package tests, `go test ./internal/codehealth -run 'Test(ConfigInstances|ConfigNameAndRequiredRoundTrip|EffectiveTimeout|EveryExecutedProviderHasADefaultTimeout|TimeoutBounds|SnapshotNamedInstances|SnapshotOmitsEmptyName|SnapshotIDIgnoresInstanceOrder|RenameKeepsSeriesButScopeStartsOne|NameStaysOutOfConfigDigest|InstancesAreAssessedSeparately|Discover|LizardExclude|ExecRunner|Collect|PreviewWritesNothing|ApplyThenReapplyIsUnchanged|ReconciliationKeepsOwnerEdits|MissingToolIsReportedNotRemoved|MissingScopeIsReported|ApplyWithNothingNewKeepsHandEditedFileBytes|ApplyRefusesNonProjectWithoutWriting|UnusableExistingConfigFailsClearly|ConflictingProposalIsNotSavedAndIsExplained|UnsupportedProposalsAreListedNotAdded)' -count=1`: PASS, 2.533s.
- Independent matrix: `go test ./internal/codehealth -run '^TestO028IndependentMatrix$' -v -count=1`: expected requirement failures in M5 (3 providers), M2, M10. Follow-up with boundary/setup additions: M1, M3/M4/M11, M4 bytes and M6/M9 all PASS. Pattern `(M1|M3|M4|M6)` also selected M10, whose expected failure kept exit 1; no passing-gate claim is made for that command.
- Native Windows/amd64 overlay probe: build `GOOS=windows GOARCH=amd64 go test -overlay /tmp/o028-windows-overlay.json -c -o /tmp/o028-codehealth.exe ./internal/codehealth`; invoke from PowerShell with quoted `-test.run=^TestO028WindowsChildren$`, `-test.v`, `-test.count=1`, TEMP/TMP set to `\\wsl.localhost\Ubuntu\tmp`. FAIL: context cancelled, surviving child PID 14132, WaitForSingleObject returns 258; harness cleanup kills the child. First invocation had a PowerShell dotted-flag quoting error and was corrected; the reported finding is from the corrected run. Sandbox blocked WSL interop and Go cache access; approved escalated commands enabled these probes.
- `cmp agent-skills/savepoint-design/SKILL.md templates/project-v2/agent-skills/savepoint-design/SKILL.md`: PASS. Bundled/scaffold equality also covered by repository full gate.
- `git diff --check`: PASS. go.mod/go.sum/Makefile unchanged from baseline. File reality: changed files and evidence-named test files exist; cmd/init_test.go was intentionally not changed; discovery fixtures use temporary generated trees rather than an on-disk testdata directory. Temporary checker harnesses are intentionally discarded source probes, embedded below for reproducibility.
- `make build` and `go vet ./internal/codehealth ./cmd`: completion recorded in final validation below.

## Design reconciliation and observations

Design section 6 names setup accurately and preserves exact configuration ownership. O-028's confirmed decisions retain the substantive Code Health design; Design has no comprehensive Code Health architecture section, the same nonblocking planning follow-up noted in C-940. The explicit incorrect map sentence is different: it violates ARCH-04/TPL-02 and is I-096. Do not edit Design as remediation in this Check.

Low-impact testing observations: the init-warning test proves successful init and preservation but cannot capture successful stderr; the no-shell test ignores extra args, so source tracing supplies that proof. Neither is a separate Issue. Setup preserves configured field values when adding entries and preserves complete file bytes when unchanged; JSON reformatting on an explicit changed apply is not treated as data loss. T-058's owner walkthrough remains outstanding.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches** — substantial isolated failure coverage; integration gaps are recorded above.
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth** — capacity mismatch is an explicit functional Issue.
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [ ] STYLE-10 **Small diffs** — Objective adds roughly 3,900 lines across 35 files; source modules are cohesive, but review scope is large. Advisory only.

## Frozen initial scope lock

# O-028 initial Full Check scope lock

1. Requirements: Objective's eight Success Conditions, T-055..T-058 Done When; named task guardrails plus TEST-08/09. Four Task waivers are owner evidence, never technical clearance. T-058 owner walkthrough remains owner validation.
2. Baseline 64eb046, head 919544c. Changed schema/config/identity/classification, discovery/catalogue, runner/platform helpers, collect/repository freshness, setup, cmd/health, main CLI/init wiring and their tests; design skill/template and managed guides. Public surfaces: Config.Validate/DecodeConfig/EffectiveTimeout, snapshot identity/validation, Assess/Overall, Discover, Plan/PlanProject/ProjectProbe/Apply/Preview/Applied, ExecRunner.Run, Collect, health argument dispatch and init integration. No real provider normalization (O-029), collection CLI (O-030), TUI refresh (O-031), installs or unrelated lifecycle remediation.
3. Relied-on boundaries: repository observation and Store history/config/snapshot reads and writes; PATH lookup; direct subprocess target/argv/root/env; provider report delivery; sequential scheduling and cancellation including children. No network calls by Savepoint; OSV external service semantics outside this reader-seam Check, disclosure in scope.
4. Frozen matrix rows:
   M1 config instance names/duplicates/scopes and decode shape (normal, empty, missing, wrong/unknown, duplicate, timeout -1/0/1/59/60/120/599/600/601, nonfinite thresholds);
   M2 snapshot/result name and identities, unnamed compatibility, order, serialization, rename/history integration, changed scope/exclusions and independent instances;
   M3 discovery Go/Vitest/V8/Python/polyglot/unsupported, report present/missing and PATH present/missing;
   M4 discovery depth 2/3/4, bytes limit-1/limit/limit+1, files limit, hidden/vendor/generated exclusions, escaping symlink, malformed manifests, cancellation, deterministic/no writes;
   M5 discovery-to-Plan-to-Apply-to-Collect for all three executed providers, configured target vs actual argv/report source, missing report, malformed report, report cleanup;
   M6 direct runner success/nonzero/unavailable/not runnable, literal shell metacharacters, root/env, stdout below/at/above 32MiB and stderr control/wide text bounds;
   M7 runner timeout/cancel/pre-cancel/children on Unix and Windows; unavailable native Windows evidence classified explicitly;
   M8 collect config order/sequential/required and no-config, unavailable/failed/findings/timeout/cancel remaining/unsupported/reader invalid/partial;
   M9 report-only fresh/equal/stale/missing/oversize/nonregular/escaping symlink, no tests run, observation failure;
   M10 snapshot history/read failure/save failure, max configured-instance capacity plus unconfigured placeholders, preserve neighbors/no pruning and origin retention;
   M11 setup preview/apply/repeat/reconcile new proposals, preserve edited config fields and handformatted bytes, conflict/corrupt/missing inputs/tool/non-project;
   M12 CLI parser/init success/discovery failure, no discovery in board/resume/doctor/upgrade, explicit existing-project setup, docs/template equality/design reconciliation.
   Text width/layout, interactive/color modes, network redirects/retry/HTTP malformed responses: N/A (plain renderer and local subprocess/report seam; no scoped width or network client). Public output/error channels and control sanitization are applicable in M6/M11/M12. Mutable args/config reader seam traced for trust boundary; callers own trusted adapters, hostile injected adapters N/A.
5. Materiality boundary: supported configured analysis tools/report-only providers and project setup. Admit reproducible requirement violations and unsupported evidence; no unrelated dependency internals or guessed provider normalization requirements.

Workflow side-effect inventory:
| Order | Operation | Effect | Failure/final state | Cleanup/oracle |
| 1 | Parse/schema/config validation | none | fatal before tool/write | typed errors/tree bytes |
| 2 | Discovery/list/bounded read/PATH | none | named gaps or root/cancel error | tree bytes+mtime, fake PATH |
| 3 | Reconcile/render | none | conflicts visible | existing entries/bytes |
| 4 | Apply config | health dirs/temp/replace | fatal; preserve confirmed semantics | Store tests, reapply bytes |
| 5 | Load history/observe repo | git reads | fatal before analysis | tool call count |
| 6 | For each instance temp report/start/wait/read | subprocess and report | per-instance distinct failure; remaining cancelled | PID liveness/temp absence/reader bytes |
| 7 | Classify/encode/save snapshot | immutable snapshot | fatal no returned ID | Store decoding, history counts, max-bound probe |
| 8 | Init preview after scaffold | stdout only | warning after successful init | main integration tests |

External boundary cells: exact configured/actual target (M5/M6); discovery before execution (M3/M5); unavailable (M6/M8); success/non-success (M6/M8); timeout/cancel (M7/M8); malformed/partial (M8/M9); cleanup/partial side effects (M5/M7/M10); sanitized failure output (M6). Retry is explicit repeat (M11); redirects N/A local process. Windows is supported and cannot be removed from M7 by absence of a local runner.

## Independent Linux matrix harness

Restore as a scratch package test or inject with Go overlay; it intentionally fails the three unresolved requirements. Remove after probing.

```go
package codehealth

import (
 "context"
 "fmt"
 "os"
 "path/filepath"
 "strings"
 "testing"
 "time"
)

func TestO028IndependentMatrix(t *testing.T) {
 t.Run("M1 timeouts and malformed configuration",func(t *testing.T){
  for _,n:=range []int{-1,0,1,59,60,120,599,600,601} {cc:=lizardInstance("a","fake");cc.TimeoutSeconds=n;err:=cfgOf(cc).Validate();if (err==nil)!=(n>=0&&n<=600){t.Errorf("timeout %d error %v",n,err)}}
  for _,raw:=range []string{`{}`,`{"version":1,"capabilities":[{"capability":"complexity","provider":"lizard-csv","required":"true"}]}`,`{"version":1,"unexpected":1}`} {if _,err:=DecodeConfig([]byte(raw));err==nil{t.Errorf("accepted malformed %s",raw)}}
 })
 t.Run("M3 M4 M11 polyglot setup read-only repeat",func(t *testing.T){
  root:=tree(t,map[string]string{".savepoint/router.md":"state: idea\n","a/go.mod":"module a\n","b/package.json":`{"devDependencies":{"vitest":"1","@vitest/coverage-v8":"1"}}`,"b/node_modules/x/package.json":"{}"})
  before:=fileState(t,root);look:=func(s string)(string,error){return s,nil}
  ps,err:=Discover(context.Background(),root,look);if err!=nil{t.Fatal(err)}
  plan:=Plan(nil,ps,ProjectProbe(root,look));if len(plan.New)==0{t.Fatal("no tools")}; _=plan.Preview()
  if fmt.Sprint(before)!=fmt.Sprint(fileState(t,root)){t.Error("discovery/preview wrote")}
  if _,err:=plan.Apply(NewStore(root));err!=nil{t.Fatal(err)};cfg,err:=NewStore(root).LoadConfig();if err!=nil{t.Fatal(err)};if err:=cfg.Validate();err!=nil{t.Fatal(err)}
  again,err:=PlanProject(context.Background(),root,look);if err!=nil{t.Fatal(err)};before=fileState(t,root);if changed,err:=again.Apply(NewStore(root));err!=nil||changed{t.Errorf("repeat changed=%v err=%v",changed,err)};if fmt.Sprint(before)!=fmt.Sprint(fileState(t,root)){t.Error("repeat wrote")}
 })
 t.Run("M4 exact discovery file limits",func(t *testing.T){
  for _,n:=range []int{MaxDiscoveryFileBytes-1,MaxDiscoveryFileBytes,MaxDiscoveryFileBytes+1} {root:=tree(t,map[string]string{"package.json":"{}"+strings.Repeat(" ",n-2)});ps,err:=Discover(context.Background(),root,func(s string)(string,error){return s,nil});if err!=nil{t.Fatal(err)};oversize:=false;for _,p:=range ps{oversize=oversize||p.Gap==GapOversizedFile};if oversize!=(n>MaxDiscoveryFileBytes){t.Errorf("bytes %d oversize=%v",n,oversize)}}
 })
 t.Run("M6 M9 exact report bounds and freshness equality",func(t *testing.T){
  for _,n:=range []int64{MaxReportBytes-1,MaxReportBytes,MaxReportBytes+1}{f,err:=os.CreateTemp(t.TempDir(),"report");if err!=nil{t.Fatal(err)};if err:=f.Truncate(n);err!=nil{t.Fatal(err)};f.Close();data,truncated,_,failure:=readBounded(f.Name());if failure!=nil||truncated!=(n>MaxReportBytes)||int64(len(data))!=min(n,MaxReportBytes){t.Errorf("bytes=%d got=%d truncated=%v failure=%v",n,len(data),truncated,failure)}}
  root:=project(t);write(t,root,"source.go","package main\n");write(t,root,"coverage.out","report");stamp:=testClock();for _,p:=range []string{"source.go","coverage.out"}{if err:=os.Chtimes(filepath.Join(root,p),stamp,stamp);err!=nil{t.Fatal(err)}}
  cc:=CapabilityConfig{Capability:CapabilityCoverage,Provider:ProviderGoCoverProfile,Report:"coverage.out",Scope:[]string{"source.go"}}
  got:=collect(t,root,cfgOf(cc),Readers{cc.Provider:okReader(80,UnitPercent)},&fakeTools{t:t});if got.Results[0].Result.Freshness!=FreshnessFresh{t.Errorf("equal mtimes not fresh: %+v",got.Results[0])}
 })
 t.Run("M5 discovered executed reports reach reader", func(t *testing.T) {
  for _, provider := range []ProviderKey{ProviderLizardCSV, ProviderJscpdJSON, ProviderOSVScannerJSON} {
   t.Run(string(provider), func(t *testing.T) {
    root := project(t); write(t, root, "go.mod", "module example\n")
    proposals, err := Discover(context.Background(), root, func(s string)(string,error){return s,nil}); if err != nil {t.Fatal(err)}
    var cc CapabilityConfig
    for _, p := range proposals {if p.Config.Provider == provider {cc=p.Config}}
    tools := &fakeTools{t:t, behavior:map[string]func(context.Context,ToolSpec)(ToolResult,error){}}
    tools.behavior[cc.Executable] = func(_ context.Context, spec ToolSpec)(ToolResult,error){
     path := filepath.Join(spec.Dir, filepath.FromSlash(cc.Report)); if err := os.MkdirAll(filepath.Dir(path),0755); err != nil {return ToolResult{},err}; if err := os.WriteFile(path,[]byte("actual-report"),0644); err != nil {return ToolResult{},err}; return ToolResult{Stdout:[]byte("console-status")},nil
    }
    var input string
    readers := Readers{provider:readerFunc(func(_ context.Context,in ReportInput)(Reading,error){ input=string(in.Data); if input!="actual-report" {return Reading{},fmt.Errorf("expected report, got %q",input)}; return goodReading(1,capabilityUnits[cc.Capability]),nil })}
    got := collect(t,root,cfgOf(cc),readers,tools)
    r:=got.Results[0].Result
    t.Logf("provider=%s configured report=%s reader bytes=%q outcome=%s",provider,cc.Report,input,r.Outcome)
    if input!="actual-report" {t.Errorf("configured tool report was ignored")}
   })
  }
 })
 t.Run("M2 rename preserves collected history",func(t *testing.T){
  root:=project(t); cc:=lizardInstance("api","fake","api/**")
  tools:=&fakeTools{t:t,behavior:map[string]func(context.Context,ToolSpec)(ToolResult,error){"fake":stdout("report")}}
  req:=CollectRequest{Root:root,Origin:OriginOfficial,Config:cfgOf(cc),Readers:Readers{cc.Provider:okReader(5,UnitCCN)},Runner:tools}
  for i:=0;i<3;i++ {n:=i;req.Clock=func()time.Time{return testClock().Add(time.Duration(n)*time.Second)};if _,err:=Collect(context.Background(),req);err!=nil{t.Fatal(err)}}
  cc.Name="backend";req.Config=cfgOf(cc); req.Clock=func()time.Time{return testClock().Add(4*time.Second)}
  got,err:=Collect(context.Background(),req);if err!=nil{t.Fatal(err)}; snapshots,err:=NewStore(root).LoadSnapshots();if err!=nil{t.Fatal(err)}
  for _,s:=range snapshots {if s.ID==got.SnapshotID {for _,summary:=range s.Summary.Capabilities {if summary.Capability==CapabilityComplexity {t.Log(summary.Explanation);if !strings.Contains(summary.Explanation,"4 official checks"){t.Error("rename lost previous 3 official observations")}}}}}
 })
 t.Run("M10 largest valid configuration collects",func(t *testing.T){
  root:=project(t);var items []CapabilityConfig;for i:=0;i<MaxResults;i++ {items=append(items,lizardInstance(fmt.Sprintf("n%d",i),"fake",fmt.Sprintf("s%d/**",i)))}
  cfg:=cfgOf(items...);if err:=cfg.Validate();err!=nil{t.Fatal(err)}
  tools:=&fakeTools{t:t,behavior:map[string]func(context.Context,ToolSpec)(ToolResult,error){"fake":stdout("report")}}
  _,err:=Collect(context.Background(),CollectRequest{Root:root,Origin:OriginManual,Config:cfg,Readers:Readers{ProviderLizardCSV:okReader(1,UnitCCN)},Runner:tools,Clock:testClock})
  t.Logf("valid config entries=%d calls=%d error=%v",len(items),len(tools.calls),err);if err!=nil{t.Error("valid maximum configuration failed after tools ran")}
 })
}

```

## Independent native Windows child probe

Inject this file with Go overlay into internal/codehealth. It uses a kernel process handle as an independent liveness oracle and terminates the child during cleanup.

```go
//go:build windows
package codehealth

import (
 "context"
 "os"
 "path/filepath"
 "strconv"
 "syscall"
 "testing"
 "time"
)

func TestO028WindowsChildren(t *testing.T) {
 helperEnv(t)
 dir:=t.TempDir();pidFile:=filepath.Join(dir,"child.pid")
 ctx,cancel:=context.WithCancel(context.Background());defer cancel()
 done:=make(chan error,1)
 go func(){_,err:= (ExecRunner{}).Run(ctx,helperSpec(dir,"child",pidFile));done<-err}()
 var pid int
 for i:=0;i<200;i++ {raw,err:=os.ReadFile(pidFile);if err==nil {pid,_=strconv.Atoi(string(raw));break};time.Sleep(25*time.Millisecond)}
 if pid==0 {cancel();t.Fatal("child did not start")}
 handle,err:=syscall.OpenProcess(0x00100001,false,uint32(pid));if err!=nil{t.Fatal(err)};defer syscall.CloseHandle(handle)
 defer syscall.TerminateProcess(handle,0)
 cancel();err=<-done
 result,waitErr:=syscall.WaitForSingleObject(handle,1000)
 t.Logf("runner returned %v; child pid=%d wait=%d err=%v",err,pid,result,waitErr)
 if waitErr!=nil||result!=0 {t.Fatal("child survived cancellation (WAIT_TIMEOUT=258)")}
}

```

## Final validation

`make build`, `go vet ./internal/codehealth ./cmd`, and `git diff --check` completed with exit 0. `./savepoint resume` strict-loaded the full V2 index successfully and reported C-941 NEEDS WORK with all five linked Issues. Git status contains only this new Check and I-092–I-096; implementation, tests, fixtures, dependencies, gate definitions, router and Task statuses are unchanged. The Check is final and immutable at handoff.
