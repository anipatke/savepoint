---
id: C-951
scope: {kind: objective, id: O-031}
result: CLEAR
checked_by: {role: checker, session: recheck-o031-final-repair-20261002-independent}
executed_session: owner-reported-i106-stream-completeness-repair
checked_at: '2026-10-01T21:56:44Z'
health_snapshot: sha256:1e6fe5430070f9c19ac24f3d977bf60f0c5bf4645b65b8dc9107d3f8dd11e9f7
reviewed:
  base_commit: 42c8d32dd74e3a6340985b4bb20763315ca03ba7
  head_commit: 42c8d32dd74e3a6340985b4bb20763315ca03ba7
  files:
    - internal/buildtool/main.go
    - internal/buildtool/main_test.go
    - internal/board/v2/health.go
    - internal/board/v2/health_view.go
    - internal/board/v2/health_test.go
    - internal/board/v2/update.go
    - internal/codehealth/dashboard.go
    - internal/codehealth/dashboard_test.go
    - internal/codehealth/dashboard_freshness.go
    - internal/codehealth/dashboard_freshness_test.go
    - internal/codehealth/collect.go
    - internal/codehealth/collect_test.go
    - internal/board/v2/detail.go
    - internal/board/v2/detail_test.go
    - internal/board/v2/load.go
    - Makefile
    - .gitignore
  dependencies:
    - go.mod
    - go.sum
    - internal/codehealth/runner_windows.go
    - internal/codehealth/repository.go
    - .github/workflows/ci.yml
    - .savepoint/health/config.json
issues: []
supersedes: C-950
---

# C-951: O-031 Full Objective Repair Recheck

## Closure map

| Prior Issue | Disposition | Current proof |
|---|---|---|
| I-104 | Closed, verified by C-951 | Original viewport fixture and 72-cell text/viewport matrix pass; every lower detail field, twelve references and final history are reachable at 80x20/24/40, with bounds, resize and Esc restoration. |
| I-105 | Closed, verified by C-951 | Original Ctrl+C-through-Help cancel/quit probe and existing active-refresh cancellation tests pass, including the native board suite. |
| I-106 | Closed, verified by C-951 | Exact C-950 native TerminateProcess reproduction now preserves both previous complete reports, with nonzero exit. Portable early-exit regression and complete-stream tests pass natively; ordinary completed pass/fail reports still publish. |

CLEAR. Every frozen acceptance/matrix cell now has current proof or its original explicit N/A classification. No new Issue was found. This checker session is independent of the user's repair execution; the checker only reviewed and ran temporary probes, then wrote Check/Issue evidence. The user reported the new repair in chat; no executor repair_attempted entry was manufactured on their behalf.

## Scope and admission

C-948's scope lock, numbered items, finite cells and eight-step workflow inventory are inherited unchanged. C-949's corrected reachability oracle is retained: references must be visible somewhere while scrolling, not permanently on the final screen; fitting content need not have positive Offset. The old q-from-Help assertion remains the explicitly nonblocking convention observation. C-950's native TerminateProcess probe is restored unchanged.

The owner reported a bounded repair in main.go/main_test.go and requested continuation of the pending Windows verification through this conversation. This is an owner-triggered new recheck, not an autonomous expansion of the prior convergence cycle. Native helper/completeness behavior is evaluated only in original B1/B2/C2/B3 cells; no new parser class or remote dependency perimeter is introduced. All Tasks T-070–T-074 retain owner-completed done status and their recorded Task-check waivers. Optional Task waivers are not technical clearance, and this Full Check covers all owned work.

Reviewed source includes the uncommitted user repair on HEAD 42c8d32:

- internal/buildtool/main.go SHA256 4a93cdedb2e32a548a7e146810a6bf1bcc84d28301b76ae1dd2ffb8b2758fd65
- internal/buildtool/main_test.go SHA256 34bd4ffb17a723947bf0b9f4d9cdf8c6d6984776fcc89a3581f52d9f719ad7ae

These identify the actual working-tree contents, rather than pretending the repair is committed. All tracked files were hashed when copied for native execution and compared afterward: All tracked input differences since Windows snapshot: []. All scoped evidence files exist. Temporary checker test files were explicitly removed before the host full gate; the native original reproduction used an external Go overlay only. No implementation, router, Task/Objective status, Design, health configuration or unrelated O-035 work was edited.

Admission ledger, written before probes:

# Owner-requested repair recheck admission

Frozen C-948 scope; current C-950 repro. No new perimeter.
|Item|Claim/Issue|Exact frozen cells|Allowed result|
|---|---|---|---|
|Original native TerminateProcess helper|I-106 user repair|B2 native Windows incomplete child and sentinels|prove or retain I-106|
|Completed passing/failing/build-fail reports and focused/temp cleanup|I-106 no regression|B1|prove or in-cell failure|
|Current completeness helper and early-exit tests|goTestStreamComplete repair|B2 + B1|prove existing incomplete/terminal-package rules; no new parser classes|
|Original D/F/UI/R independent fixtures, 72-cell text/viewport matrix and current regression suite|C-949 repair proof retained|D1-D3,F1,U1-U3,R1-R2,C1-C2,S1,E1-E3|prove original cells only|
|Fresh host full gate and native full suite on current uncommitted source|required gates|B3,C2|pass or unverified/failure|
|Official collection after gate|Full Check evidence|scope-lock 1|record saved snapshot/verdict|


## Frozen coverage matrix

| Exact frozen cell | Current proof | Classification |
|---|---|---|
| D1 | Current LoadDashboard unconfigured/first-run/damaged-history/no-tool/no-write tests and original production opening fixture pass | Proven |
| D2 | Original independent saved 14-cell matrix (seven outcomes x required/optional), five rows/manual provenance; full missing/multiple-instance regression tests pass | Proven |
| D3 | Persisted classification/overall, partial/optional failures, newest manual/official, comparable/reset/manual-excluded trends, bounded newest-first history, banned-claim tests; saved fixture rerun | Proven |
| F1 | Original real-Git same/dirty/one-commit fixture rerun and scripted same/behind/ahead/diverged/shallow/unavailable/cancel tests pass | Proven |
| U1 | Original production config/report → H/first run → R/progress → one manual snapshot → reload/freshness fixture passes; full/race initial-screen/load-error/exclusivity/cursor/late-freshness tests pass | Proven |
| U2 | Original viewport reproduction with C-949 oracle plus twelve-reference and all-field reachability, repeated arrows, final history, resize/bounds at 80x20/24/40 pass | Proven, I-104 closed |
| U3 | Five-signal labels/reasons/trends, partial/optional/required failure visibility, no overall number, glyph/word semantics, Help and minimum width; same 72-cell text/layout matrix passes | Proven |
| R1 | Progress order/count, second R ignored, successful reload/freshness, collection error retained result, Esc cancellation/no-save/retry, direct q/ctrl+c tests plus production workflow | Proven |
| R2 | Original Ctrl+C from refresh Help probe and native/full/race regression pass | Proven, I-105 closed |
| C1 | Current nil/callback compatibility, before-first/mid/after-last cancellation/no-write and named healthcheck cancellation tests pass | Proven |
| C2 | Full native Windows and Linux codehealth suites: outcome isolation, report-only reuse, malformed/absent/partial results, real process cancellation/cleanup | Proven for frozen native process cells; unrelated platform filesystem prerequisites retain existing skips |
| S1 | Stored official label, missing/no-field/no-directory/unreadable lookup, Objective-only Health line, unchanged Task/Next/badges/plain-table tests pass | Proven |
| B1 | Pass/fail/focused/no-temp-files/summary/exit tests pass on Linux and Windows; both fresh full gates write and replace conventional reports. New stream completeness table covers package terminal results without treating test terminals as package terminals | Proven |
| B2 | Original Unix SIGTERM probe passes; exact Windows TerminateProcess probe passes and preserves both sentinel files; native portable early-exit test and helper table explicitly have pass events | Proven, I-106 closed |
| B3 | Fresh make build && make test-full: complete Linux suite and six cross-builds. Fresh native full Go suite using current working-tree files: exit 0. Evidence is local native execution, not a fabricated GitHub CI run | Proven native runtime and full gate |
| E1 | No new numeric validator or remote/server/browser/redirect entry point; existing typed API/config boundaries remain | N/A, original reason retained |
| E2 | Full non-TTY board parity tests pass; health UI remains interactive-only and compact line detail-only | Proven unchanged; refresh UI N/A in redirected sink |
| E3 | Original independent 72 cells: eight text classes x normal/NO_COLOR/TERM=dumb x frozen 80x20/24/40; all fields/references/history reachable and every intermediate screen within bounds | Proven |

E3 includes ASCII, control/tab/newline/CR, combining marks, variation selector, wide characters, emoji modifiers, regional flag and joined emoji. It reuses the established Lip Gloss width oracle, adds no new width policy and has no independent claim about terminal fonts. Full original matrix rows were rerun through the original fixtures and current full/race suites, rather than stopping when I-106's exact probe passed.

## Workflow and adversarial evidence

Opening remains read-only config/history projection followed by an explicit bounded Git comparison; rendering does no IO. Manual refresh invokes the configured collector through a cancellable context, delivers sequential callback messages, saves one snapshot only on success and reloads persisted output. Before/mid/after-final cancellation cases preserve snapshot listings, real native runners stop the process tree, and error paths retain previous view. Full current tests cover configured target, discovery ordering, unavailable/nonzero/timeout/malformed/absent/partial outcomes and cleanup without adding a server/network perimeter.

The changed report path now independently establishes both process completion and package-stream completion before publishing. Every package start needs a package pass/fail/skip, while test-level pass/fail cannot close its package. A failed empty stream is not complete. The final file rename remains byte-atomic. Completed passing and failing Go runs still publish both reports; absent coverage from a build failure retains the previous coverage profile as already specified in T-074 evidence. The new reader performs only bounded local JSON scanning after closing the writer; no public CLI/schema or dependency changes were added.

Most importantly, C-950's native helper still uses actual Windows TerminateProcess with exit code 1 after one package start, not merely the implementation's own completion calculation. It now returns a nonzero error with both old report contents byte-identical. This checks the same failure and sentinel oracle that previously failed. The portable new regression is supporting evidence, and its native pass event was independently read from the saved Windows JSON report.

## Commands and results

- Linux Go 1.26.2 linux/amd64; Windows Go 1.26.2 windows/amd64; PowerShell 5.1.26100.9444. Executed 2026-10-01 UTC / 2026-10-02 Australia/Sydney against current uncommitted user repair.
- `go test ./internal/codehealth ./internal/buildtool -run 'TestO031Independent|TestRunGoTest|TestGoTestStreamComplete' -count=1`: exit 0; codehealth 0.253s, buildtool 1.041s.
- Original board production/viewport/Ctrl+C probes plus `TestO031RecheckTextAndViewportMatrix`: exit 0, 19.441s. Their exact source and C-949 oracle are preserved in the prior immutable record, and the current restoration was saved at /tmp/o031-fixed-recheck-harness.md before deletion.
- Fresh `make build && make test-full`: exit 0 after removing temporary Check test files. Uncached full host suite and Linux/Darwin/Windows builds pass, and go-test.json/coverage.out are regenerated. No earlier full result was reused for this code change.
- `go test -race ./internal/board/v2 ./internal/codehealth ./internal/buildtool -run 'Health|Dashboard|Collect|RunGoTest|GoTestStreamComplete|ExecRunnerStops' -count=1`: exit 0; board 2.463s, codehealth 7.965s, buildtool 2.067s.
- Current tracked source was zipped, hashed and extracted into Windows-local temporary path `C:\Users\User\AppData\Local\Temp\savepoint-o031-78584bfb9f6e46f09a51212b322d9d48`. Native `go run ./internal/buildtool test -reports -json -count=1 ./...`: WINDOWS_GATE_EXIT=0. Codehealth 25.934s, board 18.313s, buildtool 5.842s. The existing Unix shell fixture skips on Windows; the new portable early-exit regression runs there. No new Windows skip was added by this repair.
- Independent reading of the saved native go-test.json explicitly printed `pass TestRunGoTestWithReportsKeepsCompleteReportsWhenChildExitsEarly` and `pass TestGoTestStreamComplete`.
- Exact C-950 overlay: `GOOS=windows GOARCH=amd64 go test -overlay /tmp/o031-windows-overlay.json -c -o /tmp/o031-buildtool-windows-fixed.exe ./internal/buildtool`: exit 0 after required cache-access escalation. The initial sandbox cache write was refused and was not treated as test evidence.
- Native `o031-buildtool-windows-fixed.exe -test.run=^TestO031NativeWindowsInterruptedReports$ -test.v -test.count=1`: PASS, 0.31s. Same TerminateProcess/sentinel helper as C-950; no implementation edit by the checker.
- After the successful fresh full gate, `./savepoint health check O-031`: exit 0; creates official sha256:1e6fe5430070f9c19ac24f3d977bf60f0c5bf4645b65b8dc9107d3f8dd11e9f7. Code Health does not block clearance. Tests/coverage/complexity/duplication have no blocking finding. Optional dependency vulnerability collection failed; this remains incomplete measurement, not proof of bad code, and creates no automatic Issue.
- `git diff --check`: clean; final strict-loading and diff check follow writing the record and Issue verification metadata.

## Acceptance and guardrails

T-070 criteria 1–7, T-071 criteria 1–5 including native runner evidence, T-072 technical criteria 1–9, T-073 criteria 1–5 and T-074 criteria 1–6 are Proven within the original acceptance interpretation. The native full suite supplies the Windows runtime evidence that cross-compilation could not; the existing CI workflow runs its native Windows job, while this Check records local execution honestly rather than asserting an unobserved CI run. All O-031 technical success conditions are supported. Owner visual validation is a separate acceptance decision below.

ARCH-02/03/04 boundaries remain intact; malformed saved data still has DATA-03 diagnostics. FS-05 paths use filepath, and incomplete report preservation is proven independently on both platforms. CFG-02/03 platform subprocess semantics now have native proof and a portable regression. TEST-01/02/04/08/09 have named happy/failure/temp-project/fresh-full-gate/owner-waiver evidence. No classification rule, schema or provider definition was rewritten to match the outcome. No source remediation was performed in this Check.

No materiality actions are required: all original Issues have proof of repair and there are no new material findings.

## Design reconciliation and observations

The original Design observations are unchanged: section 8 needs planner reconciliation for H/dashboard/refresh/cancel/compact Health, and section 13 should name generated reports. These were nonblocking in the frozen scope; they are not newly promoted into blocking requirements. Architecture still matches the saved projection/direct deterministic collection boundary. This checker did not modify Design.

The existing Help-q close convention and first-run wording observation remain nonblocking. Optional vulnerability collection failure does not withhold CLEAR or prove safe dependencies. Colours/fonts and the owner terminal walkthrough are not inferred from model-level/native suite tests.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches** — interrupted/early-exit and completed-report paths have current Linux/native Windows evidence.
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth**
- [x] STYLE-08 **Comments explain why**
- [ ] STYLE-09 **Content lives in data** — unchanged inline detail labels/history/evidence copy in health_view.go; original advisory observation only.
- [x] STYLE-10 **Small diffs**

## Owner validation and closure

All owned Tasks are already owner-completed. This current independent Full Objective CLEAR Check and verified Issue resolutions provide technical evidence for the owner's Objective decision. T-072 declares the real-terminal User Check; no owner acceptance of this current Check was provided in this conversation and none is recorded on the owner's behalf. The owner should confirm the terminal walkthrough and then accept/close O-031 through the board. The checker does not set Task or Objective status to done, change router/Goal selection or publish anything.

C-948/C-949/C-950 remain immutable. C-951 supersedes C-950; Issues I-104/I-105/I-106 append this verification and point their verified resolution to C-951.

## Native reproduction source

The following is unchanged from the C-950 failing reproduction, rerun successfully against the repaired source through an external overlay:

```go
//go:build windows
package main
import("bytes";"fmt";"os";"path/filepath";"syscall";"testing")
func init(){if os.Getenv("O031_WINDOWS_TERMINATE_HELPER")=="1"{fmt.Println(`{"Action":"start","Package":"fixture"}`);h,_:=syscall.GetCurrentProcess();syscall.TerminateProcess(h,1);os.Exit(1)}}
func TestO031NativeWindowsInterruptedReports(t *testing.T){
 dir:=t.TempDir();bin:=t.TempDir()
 exe,err:=os.Executable();if err!=nil{t.Fatal(err)}
 data,err:=os.ReadFile(exe);if err!=nil{t.Fatal(err)}
 if err=os.WriteFile(filepath.Join(bin,"go.exe"),data,0700);err!=nil{t.Fatal(err)}
 t.Setenv("PATH",bin+string(os.PathListSeparator)+os.Getenv("PATH"));t.Setenv("O031_WINDOWS_TERMINATE_HELPER","1")
 prior:=[]byte("previous complete report\n")
 for _,name:=range []string{goTestReportName,goCoverReportName}{if err=os.WriteFile(filepath.Join(dir,name),prior,0600);err!=nil{t.Fatal(err)}}
 var out bytes.Buffer
 err=runGoTestWithReports(dir,[]string{"-json","./..."},&out)
 if err==nil{t.Fatal("terminated process should fail")}
 for _,name:=range []string{goTestReportName,goCoverReportName}{got,e:=os.ReadFile(filepath.Join(dir,name));if e!=nil||!bytes.Equal(got,prior){t.Errorf("%s replaced after native TerminateProcess: %q; read=%v run=%v",name,got,e,err)}}
}
```
