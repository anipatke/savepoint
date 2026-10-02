---
id: C-950
scope: {kind: objective, id: O-031}
result: NEEDS WORK
checked_by: {role: checker, session: native-o031-recheck-20261002}
executed_session: user-request-repair-claude-session-01B2T7Qen3b1EsFCeVPDf92A
checked_at: '2026-10-01T21:45:21Z'
health_snapshot: sha256:7fb0a766dc4a48fc1b793c107c5acd43829627f1cac37a4f6297d8ae20045ff9
reviewed:
  base_commit: 9db1988
  head_commit: 42c8d32dd74e3a6340985b4bb20763315ca03ba7
  files:
    - internal/buildtool/main.go
    - internal/buildtool/main_test.go
    - internal/codehealth/runner_windows.go
    - internal/codehealth/runner_test.go
    - internal/codehealth/collect_test.go
    - internal/board/v2/health_test.go
    - Makefile
  dependencies:
    - go.mod
    - go.sum
    - .github/workflows/ci.yml
    - .savepoint/health/config.json
issues: [I-106]
supersedes: C-949
---

# C-950: O-031 targeted native Windows recheck

## Closure map

| Original Issue | Current state | Evidence |
|---|---|---|
| I-104 | Still open; repair proven in C-949 | No code change since the passing frozen viewport/text matrix; full native board suite passes. No verified closure without a CLEAR Check. |
| I-105 | Still open; repair proven in C-949 | No code change since the passing frozen interrupt probe; full native board suite passes. No verified closure without a CLEAR Check. |
| I-106 | Still open; reproduced natively on Windows | Native TerminateProcess replaces prior complete JSON report with a start-only stream. |

NEEDS WORK for the existing I-106 interruption-preservation defect on Windows. The previously missing native runtime evidence has now been obtained; it reveals an actual in-scope failure rather than merely remaining unavailable. No new Issue ID was allocated. This checker session did not execute or repair implementation.

## Frozen scope and admission

This is the targeted follow-up to full recheck C-949, using C-948's immutable lock. The user explicitly requested the existing native PowerShell route. Before probes, the C-949 admission ledger had named native C2/B3 and the native interruption equivalent required by I-106. Additional admission detail for this run:

| Item | Prior finding or claim | Exact frozen cell | Allowed result |
|---|---|---|---|
| Full native Windows suite in isolated current-source copy | C-949 missing native evidence | C2/B3 | prove runtime evidence or explicit failure |
| Native child interruption after one start event, old report sentinels | I-106 native equivalent required | B2 on original Windows subprocess boundary | prove preservation or retain I-106 |
| Code/input identity and reuse of local full gate | C-949 tested HEAD and metadata-only continuation | Scope-lock release gate / B1/B3 | reuse only with all gate inputs unchanged |
| Official configured health collection | Required Objective Check support | Scope-lock full gate + collection | record official snapshot and verdict |

No new text class, exit-code range, provider policy or dependency layer was admitted. Windows TerminateProcess is the direct native analogue of the initial SIGTERM reproduction: abnormal child termination after partial stream. The ordinary Windows full suite does not test that cell because its Unix-script interruption test skips there.

## Results and exact evidence

| Frozen cells | Result |
|---|---|
| D1–D3/F1/U1–U3/R1–R2/C1/S1/B1/E1–E3 | C-949's current independent full-matrix and race proof retained: source/tests/fixtures/dependencies/gates are unchanged. Native package suites also passed where applicable. E1/E2 N/A classifications retain their exact original reasons. |
| C2 | Native Windows codehealth suite passes, including collection cancellation and real process-stop tests. Some filesystem tests skip due to local symlink/permission prerequisites; no broader filesystem perimeter was introduced. |
| B3 | Native complete Go suite passes with exit 0 on go1.26.2 windows/amd64, in addition to C-949 Linux full gate and six cross-builds. This is local native evidence, not a claim that GitHub CI ran this HEAD. |
| B2 | FAIL on Windows: I-106 persists. Current completed pass/failure/focused/report-cleanup tests pass natively, but the required interrupted child loses the previous complete JSON report. |

All prior matrix classifications remain explicit; the targeted native cells are now exercised. No remaining cell was dropped after this failure was found.

### Native full suite

PowerShell 5.1.26100.9444; Windows Go `C:\Program Files\Go\bin\go.exe`, Git `C:\Program Files\Git\cmd\git.exe`; Windows-local TEMP `C:\Users\User\AppData\Local\Temp`. Sandbox WSL interop initially failed with socket failed 1; approved escalation enabled native execution.

A ZIP of current tracked repository files was extracted into `C:\Users\User\AppData\Local\Temp\savepoint-o031-01414fb4e1f347e4b875853cfd84e8ca`. This avoids UNC working-directory and Git ownership issues and writes no source into the shared working tree. Every Go file plus go.mod, go.sum and Makefile was hashed before the copy and compared afterward: Source manifest differences: []. All tracked files came from the current working tree, not a stale commit archive.

Native command: `go run ./internal/buildtool test -reports -json -count=1 ./...`. Result WINDOWS_GATE_EXIT=0. Full suite ran on go1.26.2 windows/amd64 and wrote both conventional reports in the isolated copy. Slowest packages: codehealth 28.715s, board/v2 20.062s, migrate 18.534s, init 14.947s, root 13.057s, data 11.819s, buildtool 6.421s. It reported 32 skips; notably TestRunGoTestWithReportsKeepsCompleteReportsWhenInterrupted skips Windows because its fixture is a shell script. The independent equivalent below covers the actual Windows interruption cell.

### Native interruption reproduction

Compiled using `GOOS=windows GOARCH=amd64 go test -overlay /tmp/o031-windows-overlay.json -c -o /tmp/o031-buildtool-windows.exe ./internal/buildtool`, exit 0. The overlay adds only a temporary Windows-only checker test; it never edits implementation. Executed natively through PowerShell with `-test.run=^TestO031NativeWindowsInterruptedReports$ -test.v -test.count=1`.

```text
=== RUN   TestO031NativeWindowsInterruptedReports
    o031_native_windows_scratch_test.go:16: go-test.json replaced after native TerminateProcess: "{\"Action\":\"start\",\"Package\":\"fixture\"}\n"; read=<nil> run=exit status 1
--- FAIL: TestO031NativeWindowsInterruptedReports (0.30s)
FAIL
```

The helper exe prints a single start event, calls TerminateProcess with code 1, and never emits a package terminal event. Both previous destination reports are complete sentinels. `goTestRunCompleted` at main.go:131-136 accepts nonnegative ExitCode; Windows termination has such a code, so main.go:113-117 renames the partial JSON. Atomic rename is byte-safe but still loses complete evidence. The returned error remains exit status 1. Coverage sentinel is preserved because the helper produces no profile. This independently violates T-074 criterion 2 through the originally promised native path.

## Gates, health and unchanged inputs

C-949's fresh `make test-full` result is reused only for this metadata-only/native-evidence continuation. Original command: `make test-full`, Go 1.26.2 linux/amd64, completed 2026-10-01 around 21:19 UTC, recorded in C-949 checked at 2026-10-01T21:22:45Z; exit 0, full uncached host tests and six cross-builds. `git diff --name-only 42c8d32 -- internal go.mod go.sum Makefile templates agent-skills .github/workflows/ci.yml` is empty; no implementation/test/fixture/dependency/gate-definition changes occurred. Temporary probes use an external Go overlay only. No untracked code file is introduced. The native manifest also confirms no Go/gate input changed during execution. Only Check/Issue metadata and permitted health artifacts are written here. C-949's race suite and full frozen text matrix therefore remain current; no earlier failed run is reused.

After the successful native full suite and retained current Linux full gate, `./savepoint health check O-031` exits 0 and saves official snapshot sha256:7fb0a766dc4a48fc1b793c107c5acd43829627f1cac37a4f6297d8ae20045ff9. Code Health does not block clearance. Optional tests/coverage reports are stale due to metadata changes since the root gate; optional dependency-vulnerability collection failed; complexity and duplication have no blocking finding. This records incomplete/stale measurement honestly, without treating it as bad code or automatically opening an Issue. The health verdict cannot excuse I-106.

`git diff --check` passes. Final `savepoint resume` strict-loads the new Check and updated Issue. No Task/Objective status, router, health configuration, implementation or Design was edited. Owner's task-ids/O-035 work is retained untouched.

## Acceptance, policy and materiality

All local criteria proven in C-949 stay proven. Native Windows full-suite/process behavior is now observed rather than inferred from cross-builds. T-074 Done When 2 remains Issue for B2; ordinary report pass/failure and cleanup criteria pass. Required complete-data/explanation/refresh/cancel behavior has the same C-949 evidence. The full gate is green, yet independent acceptance review is NEEDS WORK.

| Issue | Likelihood | Impact | Materiality | Recommendation |
|---|---|---|---|---|
| I-106 | Low: interrupted Windows test child | Medium: previous complete generated evidence lost, no user source loss | Medium | Fix platform-independent completion detection; preserve incomplete reports while still publishing completed failures; then recheck the frozen interruption cells |

I-104/I-105 are not repeated as current defect failures. No out-of-scope observation affects clearance. Native suite skips beyond the original subprocess boundary are recorded, not converted into a new review perimeter. Design reconciliation observations and T-072 owner terminal validation remain as recorded in C-949.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [ ] STYLE-03 **Test branches** — main_test.go skips interrupted-report behavior on Windows; independent native reproduction fails.
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth**
- [x] STYLE-08 **Comments explain why**
- [ ] STYLE-09 **Content lives in data** — unchanged C-949 health_view.go inline labels; advisory only.
- [x] STYLE-10 **Small diffs**

## Handoff and convergence

This completes the targeted native evidence pass after C-949. Stop autonomous review/remediation cycling here and hand the reproduced I-106 repair to the executor/owner. Do not broaden the frozen scope. Completed T-074 stays done; repair directly under I-106. A future owner-requested independent recheck writes a new immutable Check. Only a CLEAR Check may verify-close the Issues; no owner acceptance is inferred. The Objective remains open.

## Exact native harness

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
