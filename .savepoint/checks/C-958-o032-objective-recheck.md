---
id: C-958
scope: {kind: objective, id: O-032}
result: NEEDS WORK
checked_by: {role: checker, session: recheck-o032-20261002-independent}
executed_session: repair-o032-20261002
checked_at: '2026-10-02T06:19:18Z'
health_snapshot: sha256:900a8a18079af14871a7153483d4b439821e96a616b89c02ab46d54d65785916
reviewed:
  base_commit: 67171b02f11faeb73e0356a9f895c4be3acad231
  head_commit: d00543d5252fc5557aeba067995427236da2823f
  files:
    - .savepoint/Design.md
    - .savepoint/health/config.json
    - AGENTS.md
    - CHANGELOG.md
    - README.md
    - internal/board/v2/health_bench_test.go
    - internal/board/v2/update.go
    - internal/buildtool/main.go
    - internal/buildtool/main_test.go
    - internal/codehealth/dashboard.go
    - internal/codehealth/discovery.go
    - internal/codehealth/discovery_catalogue.go
    - internal/codehealth/discovery_test.go
    - internal/codehealth/history.go
    - internal/codehealth/history_bench_test.go
    - internal/codehealth/readers_polyglot_test.go
    - internal/codehealth/runner.go
    - internal/codehealth/runner_safety_test.go
    - internal/codehealth/runner_test.go
    - internal/codehealth/storage_failure_test.go
    - internal/codehealth/window.go
    - internal/codehealth/window_test.go
    - internal/data/write_test.go
    - internal/doctor/checks.go
    - internal/doctor/repairs.go
    - internal/doctor/repairs_test.go
    - internal/init/agent_skills_test.go
    - internal/init/skill_validation_test.go
    - internal/init/upgrade.go
    - internal/init/upgrade_test.go
    - main_health_test.go
    - .github/workflows/ci.yml
  dependencies: [go.mod, go.sum, package.json]
issues: [I-121]
supersedes: C-957
---

# C-958: O-032 Full Objective Recheck

## Result and authority

NEEDS WORK. Six original findings are technically proven repaired or evidenced; I-121 remains materially open. This checker conversation performed no implementation, and is independent from repair-o032-20261002. This is the first complete recheck after C-957 remediation. No scope expansion, new Issue, Task retreat, owner acceptance, router change or implementation repair. All fourteen owner-done Tasks retain their status. Prior C-957 is immutable.

## Frozen scope and admission

C-957's full initial scope lock, matrix, Task criteria, supported representations, Guardrails and Design reconciliation are incorporated unchanged. Its reviewed dirty tree was committed with repairs in c5f240a; head-only diffs against 136fe69 therefore also include already-reviewed original work. Reviewed repair paths are storage_failure_test.go, window.go/window_test.go, runner.go/collect.go/runner_safety_test.go, buildtool/main.go/main_test.go, README/CHANGELOG and ci.yml. Scope admission was written before the adjacent probes:

# C-957 frozen recheck admission ledger

Recheck independent of repair-o032-20261002. Same checker conversation contains only prior review, no execution. Recheck supersedes C-957, never edits it. No scope amendment or new axis.

| Item | Original Issue/claim | Exact frozen cell | Allowed result |
|---|---|---|---|
| Verify current native job, full Linux distribution job and native replacement test | I-115/I-116; run 36971878006 d00543d | R native Windows current tree; S complete read/replace/platform refusal | Pass or original Issue remains unverified/open |
| Original 36-case official/manual count matrix | I-117 windowStart repair | W official count 0–11 x leading/trailing/interleaved manuals | Pass or I-117 remains open |
| Original repeated created_at and origin before normal header; repetitions after normal fields already named in I-118 proof | I-118 parseHead repair | W valid duplicate headers, order/full-decoder/newest identity | Pass or I-118 remains open; no new ordering layer |
| Original normal/compact/reordered/512-prefix/wrong-type/trailing/symlink representations; inside/outside damage and no-write | I-118 repair adjacency and original W pass | W exact recorded representation and failure cells | Pass or original acceptance in-scope failure |
| Original sparse-series fixture, copy matches shorter trend and measured partial semantics | I-119 README/CHANGELOG repair | W sparse/changing series; F history/limits/partial report wording | Pass or I-119 remains open; approved sparse trend behavior not a code Issue |
| Original passing child/broken output, failing child/broken output, independent totals; malformed output named original I-120 | I-120 first-error/join repair | B diagnostic writer failure, primary child failure; valid/malformed/mixed events | Pass or I-120 remains open |
| Original synthetic URL credential normal stderr through sanitizers and Collect; normal and oversized stderr through actual ExecRunner/Collect as named I-121 proof | I-121 redaction repair | E original secret-safe failure output, 64KiB stderr -1/exact/+1 overflow, control/UTF8/bounded reason; I-121 proof explicitly names oversized stderr | Pass or I-121 remains open; no new secret class (same URL userinfo) |
| Original D oracle, L oracle, S retention, P 45-input readers, A preservation, F benchmarks | Original C-957 passing rows | Exact C-957 D/L/S/P/A/F completed cells and same scratch harness | Pass or matching in-scope failure |
| Fresh make ci/full gate and official health after gate | C-957 release gate and health evidence | R full gate/checksums/version/package; P official-only integration | Pass, health blocks, or unverified (no reuse after code change) |

Out-of-scope observations cannot block this recheck. Frozen old-body-damage boundary and accepted empty Lizard CSV I-101 remain unchanged. Entire matrix must finish even if one repair fails. No third autonomous repair cycle; checker repairs no implementation.


## Original finding closure map

| Issue | Technical assessment | Disposition |
|---|---|---|
| I-115 | Proven. Read-only GitHub job metadata/logs show run 36971878006 on exact d00543d5252fc5557aeba067995427236da2823f: windows-tests and ci both success. Windows go1.26.8 windows/amd64 runs `go run ./internal/buildtool test -json -count=1 ./...`; Linux ci runs make ci. | Open pending CLEAR Check proof; current native evidence gap is filled. |
| I-116 | Proven. Replacement test no longer skips Windows wholesale; permitted Windows replacement/read refusals are constrained, complete old/new bytes and final temp cleanup asserted. Native full job passed; replacement test is absent from its skipped-test list. Fresh Linux full and scoped race checks pass. | Open pending CLEAR Check proof. |
| I-117 | Proven. windowStart retains leading manuals when no older official exists. Original official-count 0–11 x leading/trailing/interleaved matrix passes all 36 cells. | Open pending CLEAR Check proof. |
| I-118 | Proven for the original finding. Header parser scans the bounded prefix and detects repeated origin/time, falling back to full decoding. Original repeated-origin/time fixtures agree on newest/history; permanent before/after-header regressions and all original representation cases pass. | Open pending CLEAR Check proof. |
| I-119 | Proven. README/CHANGELOG qualify continuous versus sparse comparable series, and partial evidence may carry an incomplete value but never read Good. Original sparse fixture still demonstrates the approved shorter trend rather than demanding a new history algorithm. | Open pending CLEAR Check proof. |
| I-120 | Proven. Remembered writer errors join scanner/child failures; errors.As retains child ExitError. Original success-child/broken-output and failed-child/broken-output reproductions pass, as do malformed passthrough and independent totals. | Open pending CLEAR Check proof. |
| I-121 | Still open. Normal URL userinfo is removed, but raw stderr tail truncation can remove the scheme before sanitization. Actual ExecRunner followed by Collect persists synthetic credentials in the permanent snapshot. | Material blocker; direct repair and targeted independent recheck required. |

Issue capture requires a CLEAR Check for verified disposition. Consequently no Issue is resolved as verified by this aggregate NEEDS WORK record. Histories record the six successful assessments and the remaining failure without altering original reports.

## Complete coverage and acceptance update

The entire frozen matrix was completed despite I-121's failure. D baseline mapping/sentinel/pairwise oracle and L 192-case reducer oracle pass again. B existing malformed/oversized/read/tee/child/interrupt/report cases pass the fresh full gate; original writer-error reproductions now pass. S schema/type/path/conflict/concurrent/replacement/cleanup/prune cells pass full gate and race checks; independent thirteen-count prune matrix passes. E executable/argv/cwd/timeout/orphan/cancellation/pipe/report bounds and UTF-8/control cases pass full gate, with URL-output cap adjacency failing as recorded below. P nine-reader x five malformed/empty inputs, mixed-stack history/gate/official-reference cells pass; accepted I-101 empty Lizard limitation remains unchanged. W all original count/order/representation probes pass, along with permanent body-damage/tie/full-fallback/no-write tests. A adoption/init/upgrade/preservation/skill parity passes full gate. R fresh host full gate, six cross-builds/distribution/package and independent checksum verification pass; native current-head evidence is now proven. F all nine single/multi/heavy load shapes and all seven fixed-frame rendering shapes execute successfully; documentation matches approved limits. N/A cells retain C-957's reasons.

Every C-957 per-criterion assessment remains supported by fresh gate/scenario evidence, with these replacements: T-088 criterion 2, T-090 criterion 2, T-091 criterion 5, T-095 criterion 4, T-096 criteria 1–2 and T-098 criteria 1–2 are now Proven. T-091 criterion 2 remains unmet by I-121. All other T-085–T-098 criteria retain their original classifications and disclosed limitations. O-032 integration/clearance conditions and secret-safe diagnostics cannot clear until I-121 is repaired; native, bounded history and public guidance outcomes are now supported. Owner copy/walkthrough acceptance and Objective closure remain owner decisions.

## Current gate, platform and health evidence

Fresh `make ci`, 2026-10-02, go1.26.2 linux/amd64, AMD Ryzen 7 7800X3D/WSL2: exit 0. Includes uncached make test-full, six target builds, archive verification, npm wrapper tests and dry-run package inventory. `sha256sum -c checksums.txt` passes all six archives. Version remains 2.0.5 under the recorded owner deferral until merge. `git diff --check` passes. Original overlay command across codehealth/doctor/board/buildtool exits 0; adjacent credential command exits 1 only for above-cap persistence. Scoped replacement/save/prune race command exits 0.

Native evidence: https://github.com/anipatke/savepoint/actions/runs/36971878006 ; windows-tests job 110727351995 and ci job 110727352191 both success on exact reviewed head. CI workflow now includes v2.1 push and PR triggers. No workflow revert, push, merge or dispatch performed by this checker.

After the full gate, official `./savepoint health check O-032` exits 0 and creates the cited snapshot at 2026-10-02T06:16:17Z. Code Health does not block clearance. Optional OSV results report two unknown-severity groups needing review; this is advisory, not a claim that dependencies are safe. No provider installation/configuration/dependency change. Collection snapshot/report writes are the explicit Full Check exception.

Performance single-run measurements are diagnostic, not universal speed gates: single n=100/1000 3.34/19.45 ms; heavy n=100/1000 22.96/40.34 ms. Parser scans more header tokens after the duplicate-header repair; single1000 now 8.17 MB/127104 allocs, heavy1000 26.80 MB/184859. Fixed-frame normal/history/overlay cases 0.099–0.242 ms. All disclosed bounded-load tradeoffs remain, with no claim of zero metadata cost or arbitrary manual-history bound.

## Remaining materiality and remedy

I-121 is the original secret class and exact promised oversized-stderr failure path, admitted before testing. With a synthetic URI and newline filler, stderr at 65535 and 65536 bytes redacts correctly; at 65548 bytes the raw tail drops `https://alic`, leaving `e:o032-secret-for-probe@proxy.example.invalid failure`. Sanitization no longer sees `://`; Collect records the credential in result.Reason and the immutable on-disk snapshot. Useful trailing failure text is retained. No real secret or production snapshot was used.

Likelihood low (provider emits credential-bearing URI at a truncated stderr boundary); impact high (permanent local secret exposure); materiality high under T-091 Done When 2. The executor history acknowledges this precise truncation limit; it is not an owner exception. Repair directly under I-121 without retreating completed Tasks. One targeted remediation/recheck remains under the convergence method; do not reopen cleared cells or invent a broader secret scanner requirement.

## Design reconciliation and nonblocking observations

C-957 architecture reconciliation still holds: local immutable evidence, project-owned executables, read-only rendering, no installation/telemetry/schema/dependency change; bounded dashboard loading does not alter full-history collection/pruning. Approved sparse/history damage limits and I-101 accepted CSV limit remain. Unknown-severity vulnerability advisories need separate triage. Prior temporary root report files were left intact. CI v2.1 triggers are retained as supplied; merge/version decisions remain owner-owned.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [ ] STYLE-03 **Test branches** — oversized URL userinfo truncated across the raw stderr cap remains uncovered by permanent redaction tests (I-121).
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries** — original stream error propagation now proven.
- [x] STYLE-07 **One source of truth**
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs**

Style observations are advisory and do not independently determine the verdict.

## Reproducible evidence appendix

Original harness/oracle sources are preserved in C-957. Re-run its overlay across the four scoped packages with `-run '^TestO032' -count=1 -v`; all now pass. The additional probe below uses an overlay target in internal/codehealth, no repository test edit. It invokes the running test binary as the configured provider, and validates Collect plus persisted snapshot contents. Run `go test -overlay /tmp/o032-recheck-overlay.json ./internal/codehealth -run '^TestO032Recheck' -count=1 -v`.

### o032_recheck_adjacent_test.go

```go
package codehealth

import (
 "context"
 "os"
 "strings"
 "testing"
)

// Same synthetic URI-userinfo failure as I-121, crossed with the explicitly
// named stderr cap and control/whitespace cells. This goes through real
// ExecRunner and the complete Collect persistence path.
func TestO032RecheckCredentialCapMatrix(t *testing.T) {
 const credential="o032-secret-for-probe"
 for _,offset:=range []string{"below","exact","above"} {
  t.Run(offset,func(t *testing.T){
   t.Setenv("O032_CREDENTIAL_HELPER",offset)
   exe,err:=os.Executable();if err!=nil {t.Fatal(err)}
   cfg:=cfgOf(lizardInstance("local",exe,"a/**"))
   cfg.Capabilities[0].Args=[]string{"-test.run=^TestO032RecheckCredentialTool$"}
   root:=project(t)
   c:=collect(t,root,cfg,Readers{ProviderLizardCSV:okReader(1,UnitCCN)},ExecRunner{})
   r:=collectedFor(t,c,CapabilityComplexity,"local").Result
   if r.Outcome!=OutcomeFailed {t.Errorf("expected failed provider, got %s",r.Outcome)}
   if strings.Contains(r.Reason,credential) {t.Errorf("%s stderr cap: synthetic credential persisted in reason %q",offset,r.Reason)}
   snaps,err:=NewStore(root).LoadSnapshots();if err!=nil {t.Fatal(err)}
   for _,snap:=range snaps{for _,result:=range snap.Results{if strings.Contains(result.Reason,credential){t.Errorf("%s stderr cap: immutable snapshot contains synthetic credential",offset)}}}
  })
 }
}

func TestO032RecheckCredentialTool(t *testing.T) {
 mode:=os.Getenv("O032_CREDENTIAL_HELPER")
 if mode=="" {return}
 uri:="https://alice:o032-secret-for-probe@proxy.example.invalid"
 size:=maxStderrBytes
 switch mode {case "below":size--;case "above":size+=12}
 // The retained tail starts inside the URL when the raw cap is exceeded.
 // Sanitization later collapses the whitespace, preserving that exposed userinfo.
 message:=uri+strings.Repeat("\n",size-len(uri)-len(" failure"))+" failure"
 os.Stderr.WriteString(message)
 os.Exit(127)
}

func TestO032RecheckCredentialSanitizerRetainsFailure(t *testing.T){
 const secret="o032-secret-for-probe"
 for _,in:=range []string{
  "proxy failed https://alice:"+secret+"@proxy.example.invalid",
  strings.Repeat("progress\n",1000)+"proxy failed https://alice:"+secret+"@proxy.example.invalid",
 }{
  for _,got:=range []string{sanitizeLine(in),sanitizeTail(in)}{if strings.Contains(got,secret){t.Error("userinfo retained")}}
  if got:=sanitizeTail(in);!strings.Contains(got,"proxy.example.invalid"){t.Error("useful host lost")}
 }
}

var _ = context.Background

```

### o032-recheck-adjacent.log

```text
=== RUN   TestO032RecheckCredentialCapMatrix
=== RUN   TestO032RecheckCredentialCapMatrix/below
=== RUN   TestO032RecheckCredentialCapMatrix/exact
=== RUN   TestO032RecheckCredentialCapMatrix/above
    o032_recheck_adjacent_test.go:25: above stderr cap: synthetic credential persisted in reason "codehealth.test exited with status 127: e:o032-secret-for-probe@proxy.example.invalid failure"
    o032_recheck_adjacent_test.go:27: above stderr cap: immutable snapshot contains synthetic credential
--- FAIL: TestO032RecheckCredentialCapMatrix (0.08s)
    --- PASS: TestO032RecheckCredentialCapMatrix/below (0.03s)
    --- PASS: TestO032RecheckCredentialCapMatrix/exact (0.03s)
    --- FAIL: TestO032RecheckCredentialCapMatrix/above (0.03s)
=== RUN   TestO032RecheckCredentialTool
--- PASS: TestO032RecheckCredentialTool (0.00s)
=== RUN   TestO032RecheckCredentialSanitizerRetainsFailure
--- PASS: TestO032RecheckCredentialSanitizerRetainsFailure (0.00s)
FAIL
FAIL	github.com/opencode/savepoint/internal/codehealth	0.081s
FAIL

```

### o032-recheck-bench.log

```text
goos: linux
goarch: amd64
pkg: github.com/opencode/savepoint/internal/codehealth
cpu: AMD Ryzen 7 7800X3D 8-Core Processor           
BenchmarkHealthHistoryLoadDashboard/single/n=10-16         	       1	   1613097 ns/op	      5266 bytes/snapshot	     52655 fixture-bytes	  686696 B/op	    5445 allocs/op
BenchmarkHealthHistoryLoadDashboard/single/n=100-16        	       1	   3339151 ns/op	      5257 bytes/snapshot	    525682 fixture-bytes	 1658840 B/op	   18195 allocs/op
BenchmarkHealthHistoryLoadDashboard/single/n=1000-16       	       1	  19448587 ns/op	      5256 bytes/snapshot	   5256194 fixture-bytes	 8171904 B/op	  127104 allocs/op
BenchmarkHealthHistoryLoadDashboard/multi/n=10-16          	       1	   4663294 ns/op	     27929 bytes/snapshot	    279287 fixture-bytes	 3289224 B/op	   23015 allocs/op
BenchmarkHealthHistoryLoadDashboard/multi/n=100-16         	       1	   8594347 ns/op	     27893 bytes/snapshot	   2789275 fixture-bytes	 5562528 B/op	   44150 allocs/op
BenchmarkHealthHistoryLoadDashboard/multi/n=1000-16        	       1	  22811826 ns/op	     27890 bytes/snapshot	  27890063 fixture-bytes	12080520 B/op	  153076 allocs/op
BenchmarkHealthHistoryLoadDashboard/heavy/n=10-16          	       1	  14200308 ns/op	    139429 bytes/snapshot	   1394287 fixture-bytes	13627024 B/op	   44262 allocs/op
BenchmarkHealthHistoryLoadDashboard/heavy/n=100-16         	       1	  22961033 ns/op	    139393 bytes/snapshot	  13939275 fixture-bytes	20813480 B/op	   75968 allocs/op
BenchmarkHealthHistoryLoadDashboard/heavy/n=1000-16        	       1	  40340614 ns/op	    139390 bytes/snapshot	 139390063 fixture-bytes	26803528 B/op	  184859 allocs/op
PASS
ok  	github.com/opencode/savepoint/internal/codehealth	7.935s
goos: linux
goarch: amd64
pkg: github.com/opencode/savepoint/internal/board/v2
cpu: AMD Ryzen 7 7800X3D 8-Core Processor           
BenchmarkHealthRenderPopover/normal/instances=1-16         	       1	    207746 ns/op	  122688 B/op	     503 allocs/op
BenchmarkHealthRenderPopover/normal/instances=4-16         	       1	    200027 ns/op	  123928 B/op	     555 allocs/op
BenchmarkHealthRenderPopover/normal/instances=12-16        	       1	    158007 ns/op	  124024 B/op	     556 allocs/op
BenchmarkHealthRenderPopover/history/entries=10-16         	       1	    116231 ns/op	   96360 B/op	     423 allocs/op
BenchmarkHealthRenderPopover/history/entries=1000-16       	       1	     99432 ns/op	   26040 B/op	     417 allocs/op
BenchmarkHealthRenderOverlay/80x24-16                      	       1	    235518 ns/op	  134144 B/op	     672 allocs/op
BenchmarkHealthRenderOverlay/200x50-16                     	       1	    241676 ns/op	  173232 B/op	     877 allocs/op
PASS
ok  	github.com/opencode/savepoint/internal/board/v2	0.010s

```

### o032-recheck-race.log

```text
ok  	github.com/opencode/savepoint/internal/codehealth	1.781s

```

Gate input manifest: 291 source/test/dependency/gate files, sorted JSON SHA256 `5ad776be2e85b77c6abfe07604f194daec8f491c1e88ccca72af06b8ef6ef3ca`. Fresh gate used; no prior full result reused. Code and gate inputs unchanged during this checker run.
