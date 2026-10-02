---
id: C-947
scope: {kind: objective, id: O-030}
result: CLEAR
checked_by: {role: checker, session: recheck-o030-20261001-independent}
executed_session: prior-o030-executor-sessions
checked_at: '2026-10-01T20:15:49Z'
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
issues: []
supersedes: C-946
---

# C-947: O-030 Full Objective Recheck

## Closure Map

| Prior Issue | Result | Proof |
|---|---|---|
| I-102 | Closed as verified | Original partial/fresh and partial/stale failing-test reproductions pass; all 156 disposition cells retain policy and limitation wording; 48 additional combinations inside the original severe/unknown/default/partial/stale cells pass. |
| I-103 | Closed as verified | Original broken-output writer returns nil with one official snapshot saved; regression verifies stderr contains saved identity and secondary error. Main passes os.Stderr. Genuine collection errors still return errors. |

CLEAR. This session reviewed the initial Check but did not implement either repair. Repairs were present as working-tree edits when the owner requested this recheck; exact external repair-session IDs are not recorded. The aggregate execution label names the prior implementation, not this independent checker. All four owned Tasks remain done with their recorded owner waivers. No acceptance or completion status was changed.

## Immutable Scope Lock

The following is copied from C-946 unchanged. Historical commit IDs describe the initial scope; this recheck also reviews the current uncommitted repairs in gate.go/gate_test.go, healthcheck.go/healthcheck_test.go and main.go. O-034 planning is outside the technical scope.

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



## Admission Ledger

| Recheck item | Prior claim | Exact C-946 frozen cell | Allowed result |
|---|---|---|---|
|Original 156 policy combinations|I-102 repair|Disposition: 12 failures + 144 opt-in|Pass or in-scope failure|
|Fresh/stale partial failing tests|I-102|Composite output: fresh/stale partial failing tests|Pass or I-102 remains|
|Severe/unknown and opt-in partial combinations|I-102 Proof Needed|Composite output: severe/opt-in, stale+partial, unknown review|Pass or I-102 remains|
|Original pass/fail/absent/stale/partial workflow|I-102|Workflow: report states + coverage reuse|Pass or in-scope failure|
|Broken stdout after save + stderr identity|I-103|Workflow: post-save sink failure|Pass or I-103 remains|
|Reference sizes/types and doctor mixed chain|Unchanged T-067|Reference parser and Doctor combined chain|Pass or in-scope failure|
|Full suite and mirrors|O-030 integration|Gates and Guidance|Pass or unverified gate|

## Acceptance and Matrix Reconciliation

| Frozen row / acceptance | Recheck evidence | Classification |
|---|---|---|
| Config, T-066.1 | Allowed/disallowed flag and unchanged digest tests, config validation in full suite | Proven |
| Policy and composite output, T-066.2–6 and O-030 distinct states | Original independent 156 cells repeated; each applicable partial reason and stale statement now asserted. Original failing-test reproductions pass. Independent 48 cells cross required/optional × fresh/stale × available/partial × passing tests/failing tests/high/critical/unknown/high+unknown. Disposition, incomplete reason, staleness, and unknown review remain visible. Existing golden/manual/not-configured/missing tests pass. | Proven |
| Reference decoder and ownership, T-067.1–2 | Original 1/127/128/129 and null/CR/sequence/map/scalar probes repeated; empty/blank/newline tests and no-data-to-codehealth test pass. Field stays outside clearance resolution. | Proven |
| Doctor references/read-only, T-067.3–4 | Original combined manual/official/missing supersession chain repeated, with tree equality; all official/missing/manual/corrupt/no-field regression cases pass. | Proven |
| Args, Objective and config guards, T-068.1–3 | Full cmd/main suite, unknown Objective no-run regression, independent invalid-config no-effects scenario, and live no-config command pass. Target/schema/index guards precede collection. | Proven |
| Persistence/output, T-068.4–5 | Original workflow pass/fail/absent/stale/partial matrix repeats; each saves one official snapshot. Broken stdout original reproduction now returns nil; named output-failure regression proves saved ID on stderr. Fixed origin and immutable permanent storage retained. | Proven |
| Cancellation/report reuse/thin main, T-068.6–8 | Full runner/collector cancellation and deadlines, cancelled-context regression; independent workflow records only Lizard, with test and coverage reports consumed without executing test tools. Main wrapper remains dispatch only. | Proven |
| Guidance and Design, T-069.1–5 | Canonical skill/reference copies compare byte-identical; full mirror/template tests pass. Review confirms full gate first, Full-only collection, optional health_snapshot, no-config/manual handling, independent Issues, package map and Design ownership. | Proven |
| Objective success conditions 1–7 | All owned Task outcomes and their cross-Task integration now meet the frozen acceptance, including repaired evidence wording and collection success semantics. No automatic Issues, Check, or status mutation in collection. | Proven |
| Frozen N/A cells | Same reasons as C-946: no interactive cursor/width surface, network service/redirects, transactions, auth or tenant boundary. Invalid direct snapshots outside documented validated-input precondition remain outside scope. | Not applicable |

All original rows were revisited through the unchanged harnesses, the supported workflow, source tracing, named regression tests and fresh full suite. No new blocking axis or dependency layer was introduced. Workflow order remains strict-load/Objective guard → config → history/repository → per-instance reports/tools → classify/save → reload/evaluate → output. Per-instance failures remain isolated. Output failure is now secondary to the successful save, with saved identity on stderr; no rollback or rewriting is introduced.

## Commands and Results

- Original C-946 appendix harnesses recreated temporarily in their four packages. `go test ./internal/codehealth ./internal/healthcheck ./internal/data ./internal/doctor -run 'TestO030|TestEvaluate|TestBlocking|TestRenderGolden|TestHealthSnapshotRefs|TestDecodeCheckV2_healthSnapshot|TestDataDoesNotImportCodeHealth|TestRun_' -count=1 -v` passed. Original workflow output: `output failure: error=<nil> snapshots=1`; partial/fresh tests line retains its incomplete reason.
- Extended existing policy cells to assert limitation wording and ran 48 in-scope composite cells: `go test ./internal/codehealth -run TestO030 -count=1 -v` passed; 156 policy cells and 48 composite cells passed.
- All scratch Go files removed before the full gate. No production/test repair was made by this checker. The extension is preserved below; original workflow/reference/doctor harnesses remain in C-946.
- Fresh `make test-full` passed on Go 1.26.2 linux/amd64, completed 2026-10-01T20:14Z: noncached complete Go suite and linux/darwin/windows cross-builds. Native Windows CI not run locally; no new skips introduced.
- `make build` passed. `git diff --check` and both changed canonical/template `cmp` commands passed.
- After the full gate, `./savepoint health check O-030` returned 0: **Code Health not configured**; nothing collected. No snapshot ID exists, so health_snapshot is omitted. No configuration was created.

## Observations, Materiality and Owner Handoff

No materiality actions remain; both original findings are verified. C-946's nonblocking scalar-coercion and out-of-scope collector observations remain unchanged. Post-save concurrent readback corruption was not reproduced in C-946 and was not admitted as a new blocking perimeter here.

The mandatory Objective Check is CLEAR and every owned Task is done. The owner may complete O-030. T-069's requested review of the updated instructions remains an owner-facing validation; this Check does not claim the owner performed it or record acceptance on their behalf.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches**
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth**
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs**

## Independent Composite Harness

This extends only exact I-102 cells already named in C-946 and I-102 Proof Needed. Recreate in internal/codehealth alongside its existing fixture helpers; scratch only.

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
 if e!=nil||v.Blocks()!=want{t.Fatalf("policy cell: %+v %v want %v",v,e,want)}
 out:=v.Render();if partial&&(!strings.Contains(out,"incomplete")||!strings.Contains(out,"two packages missing")){t.Fatalf("partial reason lost: %s",out)};if fresh==FreshnessStale&&!strings.Contains(out,"stale"){t.Fatalf("stale lost: %s",out)};cells++
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

func TestO030CompositeRecheck(t *testing.T){cells:=0
 for _,required:=range []bool{false,true}{for _,fresh:=range []Freshness{FreshnessFresh,FreshnessStale}{for _,outcome:=range []Outcome{OutcomeAvailable,OutcomePartial}{for _,mode:=range []string{"pass","fail","high","critical","unknown","high+unknown"}{
 c,p,n:=CapabilityTests,ProviderGoTestJSON,0.0;if mode=="fail"{n=1};if mode!="pass"&&mode!="fail"{c,p,n=CapabilityDependencyVulnerability,ProviderOSVScannerJSON,1}
 r:=gateResult(c,p,outcome,fresh,n);r.Reason="missing package marker";switch mode{case "high":r.Details=[]Detail{{Key:DetailHighVulnerabilities,Number:1}};case "critical":r.Details=[]Detail{{Key:DetailCriticalVulnerabilities,Number:1}};case "unknown":r.Details=[]Detail{{Key:DetailUnknownVulnerabilities,Number:1}};case "high+unknown":r.Details=[]Detail{{Key:DetailHighVulnerabilities,Number:1},{Key:DetailUnknownVulnerabilities,Number:1}}}
 v,e:=Evaluate(gateSnapshot(ClassificationNeedsAttention,r),gateConfig(CapabilityConfig{Capability:c,Provider:p,Required:required}));if e!=nil{t.Fatal(e)};want:=mode=="fail"||mode=="high"||mode=="critical"||mode=="high+unknown"||(required&&fresh==FreshnessStale);if v.Blocks()!=want{t.Fatal(mode,v,want)}
 out:=v.Render();if outcome==OutcomePartial&&(!strings.Contains(out,"incomplete")||!strings.Contains(out,"missing package marker")){t.Fatal(out)};if fresh==FreshnessStale&&!strings.Contains(out,"stale"){t.Fatal(out)};if strings.Contains(mode,"unknown")&&!strings.Contains(out,"unknown severity need review"){t.Fatal(out)};cells++
 }}}};fmt.Printf("Composite recheck: %d cells passed\n",cells)
}

```
