---
id: C-946
scope: {kind: objective, id: O-030}
result: NEEDS WORK
checked_by: {role: checker, session: check-o030-20261001-independent}
executed_session: prior-o030-executor-sessions
checked_at: '2026-10-01T20:06:17Z'
reviewed:
  base_commit: 6247050
  head_commit: 99375ad
  files:
    - internal/codehealth/config.go
    - internal/codehealth/errors.go
    - internal/codehealth/gate.go
    - internal/codehealth/gate_test.go
    - internal/data/check_v2.go
    - internal/data/check_v2_test.go
    - internal/doctor/health_snapshot_refs.go
    - internal/doctor/health_snapshot_refs_test.go
    - internal/doctor/repairs.go
    - internal/doctor/v2_runtime.go
    - internal/healthcheck/healthcheck.go
    - internal/healthcheck/healthcheck_test.go
    - cmd/health.go
    - cmd/health_test.go
    - main.go
    - main_health_test.go
    - AGENTS.md
    - templates/project-v2/AGENTS.md
    - agent-skills/savepoint-check/SKILL.md
    - agent-skills/references/check-method.md
    - templates/project-v2/agent-skills/savepoint-check/SKILL.md
    - templates/project-v2/agent-skills/references/check-method.md
    - .savepoint/Design.md
  dependencies: []
issues: [I-102, I-103]
supersedes: null
---

# C-946: O-030 Full Objective Check

## Result and independence

NEEDS WORK. All four owned Tasks reviewed, including their owner Task-check waivers. This conversation did not implement this work; executor provenance above is an aggregate label for the prior T-066–T-069 sessions, whose precise session IDs are absent from Task logs. No implementation, Task status, router, Design or prior Check changes were made. Source baseline 6247050 to 99375ad plus the existing uncommitted T-069 guidance/Design changes.

## Frozen Scope Lock

# O-030 initial scope lock
1. Acceptance: all O-030 success conditions and T-066–T-069 Done When criteria. Guardrails FS-03/04/05/06, DATA-01/03/04/05, CFG-01/02/03, ARCH-01/02/03/04, TPL-01/02/04, TEST-01/02/04/06/07/08/09, POL-01/02.
2. Baseline 6247050; head 99375ad plus T-069 working changes. Public surfaces: config decode/validate/digest, Evaluate/Blocks/Render, Check decoder, doctor snapshot reference pass, health args/dispatch, healthcheck.Run, main signal wiring, live/template checker instructions. Changed files are git diff 6247050 file inventory, excluding unrelated records.
3. Runtime: ResolveTarget/schema/index -> config -> Collect(history/repository -> existing report readers or executed tools -> Assess -> immutable SaveSnapshot) -> LoadSnapshots/Evaluate -> stdout. No automatic Check/Issue/status writes; checker creates Check separately. Existing collector internals only to establish orchestration and report reuse, not reader normalization reassessment.
4. Matrix rows: blocking config (allowed/disallowed/omitted/digest); verdict capability x required/optional x available/partial/failure outcomes x fresh/stale x Good/Watch/NeedsAttention and default/opt-in flags; default test/vulnerability numeric zero/positive; unknown severity; manual/official origins; missing result/not configured; deterministic sorting/rendering. Check reference absent/null/empty/blank/newline/128/129/wrong type; doctor official/manual/missing/corrupt/no reference and read-only. CLI missing/malformed/valid ID, directory/default/missing/nonproject, extra args/flags/help; workflow config absent/valid/invalid, Objective known/unknown, report pass/fail/absent/stale/partial, cancellation, storage failure; official persisted identity/reload/output and no test subprocess. Guidance mirrors and Design reconciliation.
5. Boundary: supported configured project paths and validated official snapshots, standard plain redirected CLI output. Direct invalid Snapshot model inputs outside Evaluate's documented validated-input precondition are N/A. Interactive cursor/color/Unicode width N/A (plain text only); network redirects/servers N/A (local CLI). Runtime external boundary: process availability, success/nonzero, cancellation/timeouts, malformed/absent reports are collector tests plus policy outcome matrix. No per-commit refresh or manual TUI review.

| Order | Operation/effect | Failure/final state | Independent oracle |
|---|---|---|---|
|1|parse/resolve/schema/index/objective|fatal before tools/writes|named error and zero runner calls|
|2|load config|absent report only; invalid fatal before collection|tree unchanged|
|3|history/repository discovery|fatal before snapshot|no saved snapshot|
|4|existing report read or configured tool|instance outcome; later instances continue|runner call inventory, expected measurements|
|5|classify/encode/save immutable snapshot|write failure fatal; existing content preserved|Store round trip and origin|
|6|reload/evaluate|errors report saved identity; no rollback|snapshot identity vs printed identity|
|7|print verdict|visible result supports checker, never clearance/status|literal expected statements and unchanged planning records|


## Acceptance Coverage

| Scope | Invariant and evidence | Classification |
|---|---|---|
|T-066.1|Blocking flag accepted only on coverage/complexity/duplication; absent flag compatibility/digest unchanged: TestBlockingFlagValidation, TestBlockingFlagDoesNotChangeDigest, config diff review|Proven|
|T-066.2|Required/optional failure and all default/opt-in dispositions: TestEvaluateRules plus independent 156-cell matrix. Partial evidence is not always labelled incomplete|Issue I-102|
|T-066.3|Manual snapshot rejected: TestEvaluateRefusesManualSnapshot; explicit origin guard|Proven|
|T-066.4|Report avoids claiming Objective clearance, but composite evidence statements disappear|Issue I-102|
|T-066.5–6|Deterministic sorting/golden renderer and table policy tests pass; they assert single-kind output and miss combinations|Proven for determinism; I-102 for semantic completeness|
|T-067.1|Absent/empty/blank/newline and 128/129 byte boundary tests; independent 1/127/128/129, CR, sequence/map/null probes. Field not used in clearance resolver|Proven|
|T-067.2|TestDataDoesNotImportCodeHealth and source imports|Proven|
|T-067.3–4|Official/manual/missing/unreadable/no-field tests; independent combined supersession chain reports manual and missing only, tree byte-identical|Proven|
|T-068.1–3|cmd argument/help tests and main help test; strict schema/index/Objective before config/tools; absent config live command and fixture|Proven|
|T-068.4|Official snapshot, ID, new/stored wording, verdict, blockers exit-success established; failing output returns error after persistence|Issue I-103; output semantics I-102|
|T-068.5|Origin is fixed Official, no manual flag/API field; permanent official round trips|Proven|
|T-068.6|signal.NotifyContext main wiring, cancelled-context test, existing runner deadline/cancel tests|Proven; no real Ctrl-C delivery performed|
|T-068.7|Independent pass/fail/absent/stale/partial flow with real temp Git and test+coverage reports, only fake Lizard executes; invalid config runs nothing|Proven for orchestration; I-102 for partial output|
|T-068.8|One-line main wrapper delegates to internal/healthcheck, updated help|Proven|
|T-069.1–5|Full-only gate-before-collection, optional field, no-config/manual rules and independent Issue judgment; live/template mirror tests, command/map/Design reconciliation|Proven for integration instructions; I-103 contradicts exit wording|
|O-030 success conditions 1–4,6–7|Full-only authorized workflow, artifact reuse, required/optional and blocking rules, no automatic Issues/completion, manual exclusion|Proven|
|O-030 success condition 5|Composite incomplete/stale/unhealthy states remain distinguishable|Issue I-102|

## Matrix Completion and Adversarial Pass

| Row | Cells and result |
|---|---|
|Config|Allowed 3 capabilities, rejected 2 defaults, omitted flag and digest: pass. Invalid JSON fails before tools: pass. Config duplicate/wrong-type validation covered by configured full gate tests.|
|Disposition|Independent 156 cells: six non-measured outcomes × required/optional (12), plus 3 opt-in capabilities × required/optional × blocking on/off × available/partial × fresh/stale × Good/Watch/NeedsAttention (144), all dispositions pass. Default tests 0/positive, severity critical/high/medium/unknown, missing/not-configured, official/manual: named gate tests pass.|
|Composite output|Fresh/stale partial failing tests: I-102. Source tracing identifies same early returns for severe vulnerabilities/opt-in results and stale+partial suppression. Boundary zero vs positive disposition preserved; no unsafe clearance reproduced.|
|Reference parser|Empty/blank/newline, 128/129, absent, valid named ref: existing tests pass; independent CR and 1/127/128/129 pass. Sequences/maps reject with named diagnostic; null absent. Scalar true/42 YAML-coerce into strings: observation only (shape validation deliberately does not require SHA ID).|
|Doctor|Single official/manual/missing/corrupt/no-ref tests pass; combined chain independently passes and stays read-only. Initial scratch combination omitted supersedes and was rejected before the reference pass; corrected fixture obeys immutable Check lineage.|
|Args and dispatch|Missing/bad/extra args, unknown/manual flag, valid/default/explicit directory and help pass named cmd/main tests. ResolveTarget missing/nonproject errors traced before writes; cwd independent absolute root joins.|
|Workflow|Config absent/invalid, Objective unknown, pre-save corrupt history, cancellation pass existing and independent tests. Independent report pass/fail/absent/stale/partial saves one official snapshot, coverage read and only Lizard executes; partial output I-102. Post-save sink failure I-103. Existing/immutable snapshot handling through Store tests and Created branch; no rewriting/pruning introduced.|
|External boundary|Actual tool target equals configured spec, runner dispatch order and no test executable verified. Unavailable/nonzero/malformed/timeout/cancel outcomes covered by full collector/runner tests and independent disposition cells. Local file read, no HTTP redirects or network server. Secondary failure at stdout independently reproduced.|
|Guidance|Both canonical changed files match their scaffold copies; full mirror/template tests pass. No Quick Task collection or automatic Issue creation. ARCH ownership/map reconciled.|
|N/A|Terminal cursor/color/width, mutable validated object ownership, transaction rollback, network redirects/auth/tenant boundaries: no such scoped behavior. Numeric nonfinite and schema roundtrip validation belong to existing validated snapshot/config contract and are covered by full gate; Evaluate explicitly trusts a validated snapshot.|

All frozen rows classified; no required cell silently omitted. Workflow inventory is in the scope lock. Probes exercise alternate policy entry and command persistence, not just copied unit examples. Scratch tests were removed after their source was preserved in the appendix; no production fix was made.

## Gates and Code Health

- Fresh `make test-full` passed on Go 1.26.2 linux/amd64 at approximately 2026-10-01T20:00Z: full noncached Go suite and linux/darwin/windows cross-builds. Native Windows CI was not run locally; existing suite has no new platform skip introduced by O-030.
- `make build` passed; `git diff --check` passed.
- `./savepoint health check O-030` after the full gate returned 0: **Code Health not configured**. No snapshot saved; `health_snapshot` omitted. No configuration was created.
- Independent probes intentionally fail on I-102/I-103; these failures are separate from the successful production full gate. Source/test/gate inputs were restored after scratch probes. The full run started before probe files were created and completed successfully; no production input changed after it.
- Decoder/doctor combined probes passed after the Go cache sandbox refusal was rerun with approved access.

## Findings and Materiality

|Issue|Likelihood|Impact|Materiality|Recommendation|
|---|---|---|---|---|
|I-102|Medium: partial or stale measured reports are supported and expected|Medium: checker loses missing-scope/staleness/review information; blockers themselves still block|Medium|Fix the narrow verdict/output composition before accepting O-030|
|I-103|Low: broken output sink after successful collection|Low: official record survives, but exit status contradicts saved collection and identity is unavailable on stdout|Low|Fix post-save success/error reporting alongside this command repair|

No waiver inferred. These findings violate explicit acceptance even though neither proves false clearance or data loss.

## Observations and Owner Validation

T-069's required User Check asks the owner to read the updated workflow and CLI exception. Its board completion waiver is present; this Check does not record the owner's reading/acceptance or create a Task Check. Do not infer owner validation from technical review.

Boolean/numeric health reference scalars coerce to strings under YAML. Doctor then reports a missing reference. This is outside the explicit empty/multiline/length malformed-reference criteria, so no separate finding.

Prior O-028 notes about executed reports outside the managed report directory remain outside this Check's reader/collector implementation scope; O-030 did not change those paths. No dependency Issues were reopened.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches** — relevant branches covered; outcome composition gaps are tracked as acceptance findings.
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries** — post-save disposition is tracked separately in I-103.
- [x] STYLE-07 **One source of truth**
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs** — new package and policy files are bounded to the Objective.

## Handoff

Direct repair under I-102 and I-103, then a new independent Check superseding C-946 using this frozen lock. Every completed Task remains done; Objective status and router selection remain owner-controlled. No acceptance, closure or clearance was granted.

## Independent Harness Appendix

The following temporary test files were run in their named packages and then removed. They reuse existing test fixture helpers but add new cross-product and end-to-end scenarios. To repeat, place each code block in its package as `o030_check_probe_test.go`, run the named command above, and discard the scratch file afterward. They are Check probes, not implementation changes.

### internal/codehealth

```go
package codehealth
import("testing";"fmt";"strings")
func TestO030IndependentMatrix(t *testing.T){
 outcomes:=[]Outcome{OutcomeFailed,OutcomeTimedOut,OutcomeUnavailable,OutcomeAbsent,OutcomeUnsupported,OutcomeCancelled}
 cells:=0
 for _,required:=range []bool{false,true}{for _,o:=range outcomes{
 cc:=CapabilityConfig{Capability:CapabilityCoverage,Provider:ProviderGoCoverProfile,Required:required}
 v,e:=Evaluate(gateSnapshot(ClassificationUnknown,gateResult(cc.Capability,cc.Provider,o,FreshnessUnknown,0)),gateConfig(cc));if e!=nil||v.Blocks()!=required{t.Fatalf("failure cell %v %s: %+v %v",required,o,v,e)};cells++
 }}
 for _,cap:=range []Capability{CapabilityCoverage,CapabilityComplexity,CapabilityDuplication}{
 p:=map[Capability]ProviderKey{CapabilityCoverage:ProviderGoCoverProfile,CapabilityComplexity:ProviderLizardCSV,CapabilityDuplication:ProviderJscpdJSON}[cap]
 for _,required:=range []bool{false,true}{for _,blocking:=range []bool{false,true}{for _,partial:=range []bool{false,true}{for _,fresh:=range []Freshness{FreshnessFresh,FreshnessStale}{for _,class:=range []Classification{ClassificationGood,ClassificationWatch,ClassificationNeedsAttention}{
 o:=OutcomeAvailable;if partial{o=OutcomePartial};cc:=CapabilityConfig{Capability:cap,Provider:p,Required:required,Blocking:blocking}
 v,e:=Evaluate(gateSnapshot(class,gateResult(cap,p,o,fresh,20)),gateConfig(cc));want:=(fresh==FreshnessStale&&required)||(fresh==FreshnessFresh&&blocking&&class==ClassificationNeedsAttention)
 if e!=nil||v.Blocks()!=want{t.Fatalf("policy cell: %+v %v want %v",v,e,want)};cells++
 }}}}}
 }
 fmt.Printf("O030 independent policy cells: %d passed\n",cells)
 for _,f:=range []Freshness{FreshnessFresh,FreshnessStale}{
 r:=gateResult(CapabilityTests,ProviderGoTestJSON,OutcomePartial,f,1);r.Reason="package example/missing has no terminal event"
 v,e:=Evaluate(gateSnapshot(ClassificationNeedsAttention,r),gateConfig(CapabilityConfig{Capability:CapabilityTests,Provider:ProviderGoTestJSON,Required:true}));if e!=nil||!v.Blocks(){t.Fatal(v,e)}
 out:=v.Render();fmt.Printf("partial failing tests (%s): %s",f,out)
 if !strings.Contains(out,"incomplete"){t.Errorf("lost incomplete evidence state and reason: %s",out)}
 if f==FreshnessStale&&!strings.Contains(out,"stale"){t.Errorf("lost stale evidence state: %s",out)}
 }
}

```

### internal/healthcheck

```go
package healthcheck
import("testing";"context";"bytes";"fmt";"strings";"os";"path/filepath";"time";"errors";"github.com/opencode/savepoint/internal/codehealth")
type brokenO030Writer struct{}
func(brokenO030Writer)Write(p []byte)(int,error){return 0,errors.New("output sink failed")}
func TestO030WorkflowMatrix(t *testing.T){
 for _,state:=range []string{"pass","fail","absent","stale","partial"}{t.Run(state,func(t *testing.T){
 root:=newProject(t);report:=fixture(t,"tests/go-pass.jsonl");if state=="fail"||state=="partial"{report=fixture(t,"tests/go-fail.jsonl")};if state=="partial"{report+="{\"Action\":\"start\",\"Package\":\"example.com/m/missing\"}\n"}
 writeFile(t,root,"go.mod","module example.com/m\n\ngo 1.26\n");configure(t,root,report)
 cfg,_:=codehealth.NewStore(root).LoadConfig();cfg.Capabilities[0].Required=true
 cfg.Capabilities=append(cfg.Capabilities,codehealth.CapabilityConfig{Capability:codehealth.CapabilityCoverage,Provider:codehealth.ProviderGoCoverProfile,Report:"coverage.out",Required:true})
 writeFile(t,root,"coverage.out","mode: set\nexample.com/m/main.go:1.1,1.2 1 1\n")
 if _,e:=codehealth.NewStore(root).SaveConfig(cfg);e!=nil{t.Fatal(e)}
 future:=time.Now().Add(time.Minute);os.Chtimes(filepath.Join(root,"go-test.json"),future,future);os.Chtimes(filepath.Join(root,"coverage.out"),future,future)
 if state=="absent"{os.Remove(filepath.Join(root,"go-test.json"))};if state=="stale"{old:=time.Now().Add(-time.Hour);os.Chtimes(filepath.Join(root,"go-test.json"),old,old)}
 runner:=&recordingRunner{csv:[]byte(fixture(t,"complexity/mixed.csv"))};out,e:=run(t,root,"O-001",runner);if e!=nil{t.Fatal(e)};assertOnlyLizardRan(t,runner)
 saved:=snapshots(t,root);if len(saved)!=1||saved[0].Origin!=codehealth.OriginOfficial{t.Fatal(saved)}
 var testResult codehealth.CapabilityResult;for _,r:=range saved[0].Results{if r.Capability==codehealth.CapabilityTests{testResult=r}}
 fmt.Printf("workflow %s: %s/%s; %s",state,testResult.Outcome,testResult.Freshness,out)
 if state=="partial"{if testResult.Outcome!=codehealth.OutcomePartial{t.Fatal("did not reproduce partial")};line:=out[strings.LastIndex(out,"- tests via"):];if !strings.Contains(line,"incomplete"){t.Error("command hides incomplete failing report")}}
 })}
 t.Run("output failure after persistence",func(t *testing.T){root:=newProject(t);configure(t,root,fixture(t,"tests/go-pass.jsonl"));runner:=&recordingRunner{csv:[]byte(fixture(t,"complexity/mixed.csv"))}
 e:=Run(context.Background(),Request{Dir:root,Objective:"O-001",Runner:runner},brokenO030Writer{});saved:=snapshots(t,root);fmt.Printf("output failure: error=%v snapshots=%d\n",e,len(saved));if len(saved)!=1{t.Fatal("snapshot not saved")};if e!=nil{t.Error("returns command failure despite saved snapshot")}})
 t.Run("invalid config no effects",func(t *testing.T){root:=newProject(t);writeFile(t,root,".savepoint/health/config.json","invalid");r:=&recordingRunner{};var out bytes.Buffer;e:=Run(context.Background(),Request{Dir:root,Objective:"O-001",Runner:r},&out);if e==nil||len(r.calls)!=0{t.Fatal(e,r.calls)}})
}

```

### internal/data

```go
package data
import("testing";"strings";"fmt")
func TestO030RefMatrix(t *testing.T){
 head:="---\nid: C-001\nscope: {kind: objective, id: O-030}\nresult: CLEAR\nchecked_by: {role: checker, session: independent}\nexecuted_session: executor\nchecked_at: '2026-10-01T00:00:00Z'\n"
 for _,n:=range []int{1,127,128,129}{c,e:=DecodeCheckV2("probe.md",head+"health_snapshot: "+strings.Repeat("x",n)+"\n---\n");if (e==nil)!=(n<=128){t.Fatal(n,c,e)}}
 for _,field:=range []string{"null","[]","{}","true","42","\"a\\rb\""}{c,e:=DecodeCheckV2("probe.md",head+"health_snapshot: "+field+"\n---\n");fmt.Printf("reference %s: decoded=%v error=%v\n",field,c!=nil,e)}
}

```

### internal/doctor

```go
package doctor
import("testing";"strings";"fmt";"os";"path/filepath";"github.com/opencode/savepoint/internal/codehealth")
func TestO030CombinedRefs(t *testing.T){root:=healthProject(t);official:=saveTestSnapshot(t,root,codehealth.OriginOfficial);manual:=saveTestSnapshot(t,root,codehealth.OriginManual);writeHealthRefCheck(t,root,"C-002",manual);writeHealthRefCheck(t,root,"C-003",official);writeHealthRefCheck(t,root,"C-004","missing-identity");for id,prev:=range map[string]string{"C-003":"C-002","C-004":"C-003"}{p:=filepath.Join(root,"checks",id+".md");b,_:=os.ReadFile(p);os.WriteFile(p,[]byte(strings.Replace(string(b),"health_snapshot:","supersedes: "+prev+"\nhealth_snapshot:",1)),0644)};before:=snapshotTree(t,root);messages:=healthRefMessages(RunV2Checks(root));fmt.Printf("combined references: %v\n",messages);if len(messages)!=2||!strings.Contains(messages[0],"manual")||!strings.Contains(messages[1],"missing"){t.Fatal(messages)};if snapshotTree(t,root)!=before{t.Fatal("doctor wrote")}}

```
