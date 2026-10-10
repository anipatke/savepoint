---
id: C-978
scope: {kind: objective, id: O-045}
result: NEEDS WORK
checked_by: {role: checker, session: check-o045-20261011}
executed_session: task-t127-20261011 / T-128-T-130 separate executor sessions
checked_at: '2026-10-10T21:10:26Z'
health_snapshot: sha256:8a2a32e697cf2f67c1b2503244f5c0168e9785bae5ea9f0edf7a56a3b869240a
reviewed:
  base_commit: 87ee065
  head_commit: 87ee065
  files:
    - internal/run/argv.go
    - internal/run/claude.go
    - internal/run/claude_test.go
    - internal/run/events.go
    - internal/run/launch.go
    - internal/run/launch_test.go
    - internal/run/lost.go
    - internal/run/proc_unix.go
    - internal/run/proc_windows.go
    - internal/run/record.go
    - internal/run/run.go
    - internal/run/run_test.go
    - internal/run/sanitize.go
    - internal/run/stop.go
    - internal/run/stop_test.go
    - internal/run/testdata/claude-2.1.296/edit-and-denial.jsonl
    - internal/run/testdata/claude-2.1.296/failing-command.jsonl
    - cmd/run.go
    - cmd/run_test.go
    - main.go
    - internal/data/launcher_config.go
    - internal/data/launcher_config_test.go
    - internal/doctor/v2_runtime.go
    - internal/doctor/lost_runs.go
    - internal/doctor/lost_runs_test.go
    - internal/init/claude_guard_test.go
    - AGENTS.md
    - .claude/hooks/guard.js
    - templates/project-v2/.claude/hooks/guard.js
  dependencies: []
issues: [I-149, I-150, I-151, I-152, I-153]
unmet: [O045-SC1, O045-SC2, O045-SC5, O045-SC6, O045-SC7, T127-DW3, T128-DW2, T128-DW3, T128-DW5, T129-DW3, T129-DW5, T130-DW6, T130-DW7, CFG-03, FS-06]
supersedes: null
---

# C-978: O-045 Full Objective Check

NEEDS WORK. Four implementation defects and one required platform-evidence gap. Every owned Task remains done; this Check does not retreat Tasks, repair code, change Design, or record owner acceptance.

Independent fresh session: this conversation did not implement any reviewed work. Reviewed scope is the uncommitted working tree on 87ee065, not that commit alone. This is the first O-045 Check. No linked prior Checks or Issues existed. Task check_waiver decisions remain scoped to skipping optional Task Checks; none waives criteria or this Full Check. T-130 owner_validation.required is true with no accepted_check. The recorded builder-default owner decision is respected; this Check tests the approved optional builder behavior.

## Full Check Progress Checklist

- [x] Establish Scope
- [x] Freeze The Check Scope
- [x] Turn Acceptance Into Invariants
- [x] Build The Mandatory Coverage Matrix
- [x] Finite External-Boundary Matrix
- [x] Workflow And Side-Effect Check Lock
- [x] Matrix Completion Lock
- [x] Parallel Planning Advice
- [x] Perform The Adversarial Pass
- [x] Re-check After Remediation
- [x] Verify File Reality
- [x] Verify Evidence And Gates
- [x] Collect Code Health Evidence
- [x] Complete The Issues Pass
- [x] Summarize Materiality
- [x] Review Code Style

Checked means completed/classified, not that every result passed. Re-check after remediation is N/A (initial run). Parallel advice is N/A; no lane/worktree promise is in scope. No required cell was omitted to permit clearance.

## Frozen Scope Lock

Initial Full Check scope lock (2026-10-11)
1. O-045 SC1–7, T-127–130 Done When, named Guardrails. Local build/full tests and native Windows gate. Owner real-run acceptance stays separate.
2. internal/run/**, cmd/run*, main run dispatch, launcher_config*, doctor lost-run wiring, guard tests and AGENTS map. Entry points: ParseRunArgs/RunRun, Launch, Start/Wait, BuildArgv, ParseClaudeLine, RunnerLost/LostRuns, doctor projection.
3. Config/schema/index/Next resolution; version subprocess; record creation, append, heartbeat, atomic replace, final summary; cancel/timeout/process-tree kill; saved-record reader. No lifecycle writes; raw log preserved, derived text sanitised.
4. Matrix rows: CLI refusal/startable verbs/version boundaries; argv placeholders/defaults; provider normal/missing/unknown/malformed/duplicate data; running to each terminal state; liveness fresh/exact30s/over30s crossed live/dead agent and runner; record IO initial/heartbeat/final failure; cancel idle/tree/stubborn and timeout; replay/result-free/nonzero; redirected output and safe diagnostics. Linux exercised; native Windows/macOS/WSL evidence classified separately. Text width, colour/layout, redirects/network, transactions and board launch N/A (not implemented here).
5. Issues require supported path and named promise/guardrail violated with credible consequence. Read-only check, no implementation repair. Separate owner validation from technical clearance. Existing Task waivers not exceptions to criteria.
Workflow: guards (fatal/no writes) -> version (fatal/no writes) -> mkdir/initial atomic record/open events (fatal) -> child start (failed record on failure) -> PID atomic write (fatal/kill) -> append/adapter + heartbeat (record accuracy required) -> wait/cancel/kill -> flush/close -> summary/final atomic record -> visible completion. Oracle: filesystem content, independently known OS process PID, literal expected text, induced filesystem errors.

## Acceptance Invariants and Coverage

| Requirement | General invariant | Classification and evidence |
|---|---|---|
| O-045 SC1 | Off writes nothing; enabled recording precedes launch and heartbeat persists within five seconds | Issue I-152; off/refusal and normal heartbeat pass named Launch/Start tests |
| SC2 | Recognised provider lines map; unrecognised/malformed lines are retained and reported unparsed | Issue I-150; fixture replay and all mapped line tests otherwise pass |
| SC3 | Cancel/timeout finish in 15 seconds with no live descendants on every supported platform | Linux Proven by all TestStop cases; Windows Unverified I-153; macOS/WSL runtime unverified locally, shared Unix implementation cross-builds |
| SC4 | Nonzero exit without result is failed with safe stderr | Proven: TestStart_failedRunRecordsExitCodeAndSanitisedStderrTail, failing Launch, no-result replay; no terminal controls remain |
| SC5 | Exact Next/root substitution with no shell; stored text strips URL credentials | Issue I-151; BuildArgv tests prove exact substitution/default/custom/error behavior |
| SC6 | Agent guard blocks run; doctor identifies stale dead runners | Issue I-149; guard copies and doctor no-runs/dead-agent tests pass |
| SC7 | Current full gates and native Windows evidence | Linux full gate Proven; native Windows Unverified I-153 |
| T-127 DW1,2,4,5 | Pre-start record, raw append, exit status, shell-free cwd/env/process hooks | Proven Linux, except stored-command safety in I-151; Windows process behavior covered by I-153 |
| T-127 DW3,6 | Persist heartbeat; named tests and handoff evidence | DW3 Issue I-152; local tests/gates pass, named waivers present |
| T-128 DW1,3,4,5,6 | Map fixture fields, save labelled report, argv validation, sanitisation, replay evidence | DW1/4/6 Proven; DW3 Issue I-152, DW5 Issue I-151 |
| T-128 DW2 | Unknown/malformed counted without loss or runner failure | Issue I-150 |
| T-129 DW1,2,4 | Stop/cancel/timeout timestamp and preserve events | Proven Linux TestStop matrix; Windows Unverified I-153 |
| T-129 DW3,5 | Dead heartbeat writer detected read-only; platform tests | Issue I-149; Windows Unverified I-153 |
| T-130 DW1,2,3,4,5 | Thin parse/refusal, strict selection, version floor, foreground/context cancel, guard | Proven via cmd/run_test, launch_test, guard tests and source tracing; physical terminal Ctrl-C remains owner's scenario |
| T-130 DW6,7 | Lost-run advice and fresh cross-platform gates | Issue I-149; Windows Unverified I-153; owner validation pending separately |

## Mandatory Matrix: Public Surfaces, Shapes, States and Sequences

| Cell | Cases | Classification / oracle |
|---|---|---|
| CLI parser / RunRun | Valid task/root orders, missing task/value, unknown flag, extra dir, help; no dispatch on error | Passed TestParseRunArgs and TestRunRun cases; exact options and exit codes |
| Launch guards | Flag absent/false; invalid project/schema/index; done/missing Task; audit/blocked verbs | Passed named Launch tests plus traced shared strict resolver; prewrite boundaries clear |
| Launch approved verbs | planned Start, in-progress Build, Test through shared Next; no synthetic prompt | Existing Start/Build tested; Test dispatch uses same allowed-verb table and resolver, source verified; no lifecycle writes |
| Version boundary | 2.1.258/259/260, 2.2.0, 3.0.0, empty, unreadable, missing executable | Passed deterministic independent version cases and refusal tests |
| Argv factory | nil/default, missing/blank command, configured args, relative root, placeholder-looking prompt, quotes/$HOME | Passed TestBuildArgv; literal strings are independent oracle; no shell |
| Provider parser | init/read/edit/Bash/other/result/denial/partial; text and blocks; missing/wrong-type/null/malformed/unknown | Known cases Passed TestParseClaudeLine; unknown system and null/numeric message content Issue I-150 |
| Provider fold | tool result correlates command; denied tool keeps no exit; repeat denial deduplicated; usage and final session/report | Passed real fixture with literal expected fields; no-result and failure fixture tests |
| Raw / derived representations | complete lines, malformed raw, final line lacking newline; saved YAML and labelled summary | Normal Passed replay/Start; command credentials Issue I-151; final persistence Issue I-152. Final no-newline gets a delimiter per one-line-per-event contract |
| Terminal-safe text | ASCII, controls/ANSI CSI/OSC, UTF-8 multibyte cutoff, URL userinfo | Provider fields Passed TestCleanText and sanitisation replay; top-level argv Issue I-151. Width/emoji grapheme layout N/A, no width renderer |
| Run states | running→finished/failed/cancelled/timed_out, missing executable, process finishes before timeout | Passed Start/Stop/Launch cases; no persisted state-reducer API allows terminal revival/backward transition |
| Liveness finite matrix | all six states × 29/30/31 seconds × absent/dead/alive PID | Passed independent 54-cell predicate test; actual dead runner/live agent Issue I-149 |
| Cancel complete sequences | idle, mid-command+grandchild, interrupt ignored; timeout+grandchild; context cancel | Passed TestStop and TestLaunch_endingTheContextCancelsTheRun; OS liveness and literal final state/event bytes |
| File safety / repetition | invalid task path, missing project, empty command; collision mkdir exclusive; project records byte-identical | Passed TestStart_refusesWhatItCannotRecordSafely and TestLaunch_neverWritesAProjectRecord; no existing run directory overwritten |
| Heartbeat and final IO failures | record replaced with directory; summary path directory | Independent probes Issue I-152, agent continues and finished reported; source shows ignored errors |
| Platform/output | Linux redirected test sink; Windows/macOS cross-build; native Windows/macOS/WSL | Linux Passed except Issues; Windows Unverified I-153; Unix runtime coverage outside Linux unverified. Colour/TTY-specific rendering N/A for plain run lines |
| Limits / mutable inputs | 500 list cap and truncation flag, byte-tail bound, valid/invalid timeout; argv copied | Source verified cap/dedup/copy paths; TestSanitizeTail and timeout config tests Passed. External concurrent mutation of Spec/Adapter internals N/A (not supported contract) |

## Finite External-Boundary Matrix

| Cell | Classification |
|---|---|
| Configured vs runtime target, discovery order | Passed: actual argv uses configured/default executable; version checks same command; config/index/Next before write |
| Dependency unavailable / startup error | Passed missing agent and missing executable tests; failed run record on start error |
| Success/non-success response | Passed fake exit0/exit3, fixture command exit2, unknown/malformed stream; I-150 classification gap |
| Timeout/cancellation | Passed Linux cancel/timeout/context and descendant probes; native Windows Unverified I-153 |
| Malformed response | Passed top-level unparsed cases; nested shapes Issue I-150 |
| Retry/repeat/cleanup/partial effects | Exclusive run-dir creation prevents overwrite; atomic temp cleanup traced; after-start persistence Issue I-152 |
| Secret-safe stored/output evidence | Provider sanitisation Passed; configured stored argv Issue I-151; raw events intentionally remain verbatim |
| Network refusal/redirect/reuse | N/A: runner invokes one local subprocess, no network client, server reuse or redirect handling |

## Workflow and Side-Effect Inventory

| Order / operation | Effect / failure owner and final state | Evidence / classification |
|---|---|---|
| 1 Parse, target/schema/config/index, Next, argv | Refusal/fatal before writes | Thin cmd tests, strict resolver, Launch refusal tests Passed |
| 2 Version subprocess | Missing/old/unreadable refusal; timeout 15s | Boundary test Passed; source confirms command, no shell, context deadline |
| 3 mkdir, initial atomic YAML, open append log | Run artifacts before agent; create failure fatal, no overwrite | Source traced exclusive mkdir, CreateTemp/write/sync/close/rename; Start guard/start tests Passed |
| 4 Start child, save PID | PID rewrite failure kills/waits/closes; start failure saves failed record | Missing executable test and source traced; no child before initial record |
| 5 Stream append / parser, stderr buffer | Raw log written before observe; append error reaches child Wait; final incomplete line flush | Fixture byte oracle Passed; parser issue I-150; flush/close errors also discarded in same I-152 inventory |
| 6 Heartbeat | Atomic replace every5s; failure should be visible | Normal heartbeat Passed; injected after-start failure Issue I-152 |
| 7 Stop poll / timeout / SIGINT / tree kill / Wait | cancelled or timed_out; wait for stopper before publishing | Linux tree/stubborn/idle cases Passed. taskkill errors discarded; native platform outcome Unverified I-153 |
| 8 Flush/close and stop heartbeat | Stop heartbeat before final snapshot | Source order verified, normal raw bytes Passed; secondary IO errors ignored under I-152 |
| 9 Summary and final YAML replace | Success externally visible only after saved artifacts | Independent failed summary/final-record writes still finished: Issue I-152 |
| 10 Launch output / read-only doctor | Prints state; doctor reads saved records, never heals | Normal tests Passed; lost runner/live agent Issue I-149. Failing stdout after persistence not independently verified; no acceptance promise of output-sink recovery |

No transaction, remote API, rollback or multi-project write exists in this scope. Initial/PID write errors are returned; heartbeat and final errors are not. Every remaining cell was classified after the first finding; no widening beyond the initial lock.

## Adversarial Pass

Direct Start bypasses Launch's version/opt-in deliberately as an internal primitive, but retains path/project/command guards. Provider nested representations bypass unparsed classification (I-150). Saved argv bypasses the provider sanitiser (I-151). Dead runner/live agent bypasses lost detection (I-149). After-start IO failure violates file-first success reporting (I-152). State transitions and cancel/repeat tests otherwise agree with independent OS and literal expected oracles. No run-state vocabulary is imported into a Task/Objective gate. No board code changed. No new dependency added.

## Gates, File Reality and Code Health

- Fresh `make test-full`: exit0, Go test suite with reports, Linux/Darwin/Windows builds, go1.26.2 linux/amd64, 2026-10-10T21:07–21:08Z.
- Fresh `make build`: exit0; `git diff --check`: exit0.
- `go test -race -count=1 ./internal/run ./cmd ./internal/doctor`: exit0.
- Independent overlay `go test -overlay /tmp/o045-overlay.json -count=1 -run 'TestCheckO045' -v ./internal/run`: exit0, observations below. PASS means the probe executed, not acceptance clearance.
- `quality_gates` build/lint/typecheck/test have no configured commands; no invented configured gate. Repo Build requirements supplied build/full verification.
- `savepoint health check O-045` after full gate: official snapshot recorded above; verdict does not block clearance. Tests, coverage, complexity and duplication collected; optional osv-scanner collection failed. This is incomplete optional vulnerability evidence, not proof of vulnerabilities and not an Issue by itself.
- Scoped files/fixtures and all changed files claimed in Tasks exist; no phantom file found. No implementation repair or test addition in the repository. Overlay harness was scratch Check work; embedded below for reproduction.
- Design reconciliation: the new package respects file-first evidence and leaves lifecycle in data; no board/renderer dependency added. Foreground command is the explicit T-130 decision, with board detachment deferred to O-046. Complete Design updates deferred explicitly to O-047; recorded original stale Design/config comment requiring a builder is a nonblocking observation under that agreed boundary, not a new acceptance requirement.

## Issues and Materiality

| Issue | Likelihood | Impact | Materiality | Recommendation |
|---|---|---|---|---|
| I-149 lost-run PID | Medium: runner interrupted/crashes while agent lives | Medium: orphan run concealed | Medium | Fix runner identity/liveness |
| I-150 unparsed classification | Medium: provider format evolution or malformed stream | Low: raw evidence retained, display/count incomplete | Medium relative to explicit audit promise | Narrow parser correction |
| I-151 saved URL credential | Low: custom builder argument contains URL userinfo | High: durable credential exposure | High | Redact derived stored command before clearance |
| I-152 ignored IO errors | Low: disk/permissions/path changes after start | Medium: saved truth absent/stale but success claimed | Medium | Surface persistence errors and prove cleanup |
| I-153 Windows evidence | High: evidence absent now | Medium: platform reliability unproven | Medium | Supply covering native Windows CI |

## Owner Validation Still Needed

T-130's required User Check has no recorded owner acceptance: real terminal Claude run, event growth, expected Task/board effects, off refusal and unchanged manual workflow. Context cancellation tests support the mechanism but do not replace this decision. No acceptance or requirement exception was inferred from marking the Task done. No prior scoped acceptance/exception requires a carried_forward entry.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches** — happy/failure branches tested; missing cases recorded as technical Issues above.
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [ ] STYLE-06 **Handle errors at boundaries** — internal/run/run.go:185,222–223 discards persistence failures.
- [ ] STYLE-07 **One source of truth** — run/health sanitiser and platform hooks are duplicated; timeout defaults repeated in data/run, documented reuse choice.
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data** — provider/default command tables centralised.
- [x] STYLE-10 **Small diffs** — scope spans four planned Tasks, each package/file remains focused.

## Review Fingerprints

These identify the uncommitted reviewed code for a later re-check; head_commit alone does not identify it.

- `internal/run/argv.go` sha256 `2cbdc3f8ec0beaabececc5ef59580331222dffe786bdc3e72cced5a631612402`
- `internal/run/claude.go` sha256 `605a77b1eb298f30171f63418c9a1e5cf9b33acc2c43dcbd6f5d45aa18f62d0b`
- `internal/run/claude_test.go` sha256 `3797930bd0cfdb431a1b68d42a4080fe3af2151807c0c3cbbe92709eea8ea811`
- `internal/run/events.go` sha256 `e0572d36162a37b3cba5dd99d8f29df07bfe0f82e32880584a901522c6fb15df`
- `internal/run/launch.go` sha256 `920ee912f894a4620dcb2a42c1563369c3c5fa87b015cfe761d05652df2a9ee9`
- `internal/run/launch_test.go` sha256 `be3f1d371a7dfe116bda3743e8532ff064911ad3e219bc0194060a537bc6a58c`
- `internal/run/lost.go` sha256 `547e019bfdddb41107fc67778eff465b34e98c73712a362f4b96c9500b81cf17`
- `internal/run/proc_unix.go` sha256 `f8c1f37fc0698a36fb753126917694d7460125f5e8345b12f7522ede50575bad`
- `internal/run/proc_windows.go` sha256 `6d56f406c2cf7b04ecfd2b5c6780228f9e5938badfb2f21df0e64a7dbf1133cc`
- `internal/run/record.go` sha256 `8f2d80e451ef3333cec9c814f5c0af88ec192dc95d59c5970610a22146e65c3d`
- `internal/run/run.go` sha256 `64499a23848eb778bdc39b4b754844461b4db5dbc5797df56f375d1edc6d5f83`
- `internal/run/run_test.go` sha256 `b79238dc44d6b00b51933e05fc60ccfbb05296eb4b2a59b38c9d80444f418037`
- `internal/run/sanitize.go` sha256 `afde49e30d9f5073a523cf0363db07f332afe55e005d8311b2ae2f5b17fa0930`
- `internal/run/stop.go` sha256 `b7886a4118d9e470123dfe9fca410d6c1ed37051548ebe7dbcc5aaf62c7f978a`
- `internal/run/stop_test.go` sha256 `a7dc2e2cdbd37b2d030372cc0e0908e1cc0d639a46b388afdbd5932681fc13c0`
- `internal/run/testdata/claude-2.1.296/edit-and-denial.jsonl` sha256 `b5ca1539a4236aafbf0a1491e7ee812db1ce59337a19065e26e1d5abfd5d0a15`
- `internal/run/testdata/claude-2.1.296/failing-command.jsonl` sha256 `a7d644423bc71fa01abb7b3aa8985315d2909be392a63f6193bfaba5a6840c1c`
- `cmd/run.go` sha256 `a7ab3962b93a4146b6828462ce7c2453fd9847fb9c312d917973778e1986b9f9`
- `cmd/run_test.go` sha256 `f1d13f39442177d2b6922f8db065d7740112278a7ed39b3d656be9d99273ef8f`
- `main.go` sha256 `4b7d84bf3b48dac3b1e045d305dc98af9696ace98836d988510cb29fa6063f5b`
- `internal/data/launcher_config.go` sha256 `09db55d338d5f25e99646f345c583a83b210a7b35ecff4caa82c082a00fd0eec`
- `internal/data/launcher_config_test.go` sha256 `633107638305d0781aa98765c0b0c7ad6825dc9964d0ab1e2e757b03c4661feb`
- `internal/doctor/v2_runtime.go` sha256 `96027b1556f18c745af7a588107071d2aac4528766ecf029bc6d6fd81e711fb2`
- `internal/doctor/lost_runs.go` sha256 `9448d72706162a485e97e0a23b75f7e1505f92aa37fc29fa5d26f97bf91bb0ce`
- `internal/doctor/lost_runs_test.go` sha256 `a77144cab49519aa468a3e2be541dec2730e5eb7efbe3c45ceab335873ad9e34`
- `internal/init/claude_guard_test.go` sha256 `8c342d2ce329573309b9ef252f86e0d4c31703dff716988807b06e5b08ea88c2`
- `AGENTS.md` sha256 `b80d3f9244cca3e56a62e2134e1d5381c512c9cf608485ff5bdee8f69aa60e5c`
- `.claude/hooks/guard.js` sha256 `ae8dde236eb5e58dbfe1e9d18bdd7e9d29395cafdae4942c40041f0a993f181d`
- `templates/project-v2/.claude/hooks/guard.js` sha256 `ae8dde236eb5e58dbfe1e9d18bdd7e9d29395cafdae4942c40041f0a993f181d`

## Independent Probe Output

```text
=== RUN   TestCheckO045Independent
=== RUN   TestCheckO045Independent/lost-runner-live-agent
    run_test.go:232: runner died, agent PID alive: RunnerLost=false (expected true for dead runner)
=== RUN   TestCheckO045Independent/unknown-system
    run_test.go:235: unknown system subtype events=[] (expected unparsed)
=== RUN   TestCheckO045Independent/stored-command-credentials
    run_test.go:240: run.yml contains URL password=true
=== RUN   TestCheckO045Independent/final-write-failure
    run_test.go:247: Wait state=finished disk record is directory=true staterr=<nil>
--- PASS: TestCheckO045Independent (0.63s)
    --- PASS: TestCheckO045Independent/lost-runner-live-agent (0.00s)
    --- PASS: TestCheckO045Independent/unknown-system (0.00s)
    --- PASS: TestCheckO045Independent/stored-command-credentials (0.02s)
    --- PASS: TestCheckO045Independent/final-write-failure (0.61s)
=== RUN   TestCheckO045Matrix
=== RUN   TestCheckO045Matrix/liveness-state-time
=== RUN   TestCheckO045Matrix/provider-shapes
    run_test.go:266: shape  => []
    run_test.go:266: shape [] => [{Kind:unparsed ToolUseID: Tool: Text: ExitCode:<nil> Failed:false SessionID: Model: Usage:map[] TotalCostUSD:<nil> Denials:[]}]
    run_test.go:266: shape 7 => [{Kind:unparsed ToolUseID: Tool: Text: ExitCode:<nil> Failed:false SessionID: Model: Usage:map[] TotalCostUSD:<nil> Denials:[]}]
    run_test.go:266: shape null => [{Kind:unparsed ToolUseID: Tool: Text: ExitCode:<nil> Failed:false SessionID: Model: Usage:map[] TotalCostUSD:<nil> Denials:[]}]
    run_test.go:266: shape {"type":7} => [{Kind:unparsed ToolUseID: Tool: Text: ExitCode:<nil> Failed:false SessionID: Model: Usage:map[] TotalCostUSD:<nil> Denials:[]}]
    run_test.go:266: shape {"type":"assistant","message":null} => []
    run_test.go:266: shape {"type":"assistant","message":{"content":7}} => []
    run_test.go:266: shape {"type":"result","usage":null} => [{Kind:final ToolUseID: Tool: Text: ExitCode:<nil> Failed:false SessionID: Model: Usage:map[] TotalCostUSD:<nil> Denials:[]}]
    run_test.go:266: shape {"type":"result","usage":[]} => [{Kind:final ToolUseID: Tool: Text: ExitCode:<nil> Failed:false SessionID: Model: Usage:map[] TotalCostUSD:<nil> Denials:[]}]
    run_test.go:266: shape {"type":"system","subtype":"future_unknown"} => []
=== RUN   TestCheckO045Matrix/version-boundaries
=== RUN   TestCheckO045Matrix/heartbeat-write-failure
    run_test.go:280: heartbeat disk failure with agent alive=true
    run_test.go:281: final state=finished
=== RUN   TestCheckO045Matrix/summary-write-failure
    run_test.go:286: unwritable summary final state=finished
--- PASS: TestCheckO045Matrix (1.23s)
    --- PASS: TestCheckO045Matrix/liveness-state-time (0.00s)
    --- PASS: TestCheckO045Matrix/provider-shapes (0.00s)
    --- PASS: TestCheckO045Matrix/version-boundaries (0.00s)
    --- PASS: TestCheckO045Matrix/heartbeat-write-failure (0.61s)
    --- PASS: TestCheckO045Matrix/summary-write-failure (0.62s)
=== RUN   TestCheckO045RealLostRunner
    run_test.go:298: runnerPID=145 alive=false agentPID=151 alive=true stale-record RunnerLost=false LostRuns=[]
--- PASS: TestCheckO045RealLostRunner (0.01s)
PASS
ok  	github.com/opencode/savepoint/internal/run	1.875s
```

## Reproducible Check Harness

Take the existing `internal/run/run_test.go`, add the `os/exec` import, insert the helper branch below at the beginning of TestMain, and append the Check functions below using a Go overlay rather than editing implementation. The helper exits after Start so its agent remains; the probe kills that process group on cleanup. All credentials are synthetic.

```go
 if root:=os.Getenv("O045_HELPER_ROOT");root!="" {
  os.Unsetenv("O045_HELPER_ROOT");os.Setenv("FAKE_AGENT","idle")
  _,err:=Start(Spec{Dir:root,Task:"T-001",Provider:"fake",Argv:[]string{os.Args[0]},Now:func()time.Time{return time.Now().Add(-time.Minute)}})
  if err!=nil{os.Exit(8)}
  os.Exit(0)
 }

func TestCheckO045Independent(t *testing.T) {
 t.Run("lost-runner-live-agent",func(t *testing.T){
  rec:=Record{State:StateRunning,UpdatedAt:time.Now().Add(-time.Minute).UTC().Format(time.RFC3339),PID:os.Getpid()}
  t.Logf("runner died, agent PID alive: RunnerLost=%v (expected true for dead runner)",RunnerLost(rec,time.Now()))
 })
 t.Run("unknown-system",func(t *testing.T){
  t.Logf("unknown system subtype events=%+v (expected unparsed)",ParseClaudeLine([]byte(`{"type":"system","subtype":"future_unknown"}`)))
 })
 t.Run("stored-command-credentials",func(t *testing.T){
  spec:=fakeSpec(t,project(t),"ok"); spec.Argv=append(spec.Argv,"https://user:checksecret@example.com")
  r,e:=Start(spec);if e!=nil{t.Fatal(e)}; r.Wait()
  b,_:=os.ReadFile(filepath.Join(r.Dir,"run.yml"));t.Logf("run.yml contains URL password=%v",strings.Contains(string(b),"checksecret"))
 })
 t.Run("final-write-failure",func(t *testing.T){
  spec:=fakeSpec(t,project(t),"slow");spec.Heartbeat=time.Hour
  r,e:=Start(spec);if e!=nil{t.Fatal(e)}
  if e=os.Remove(filepath.Join(r.Dir,"run.yml"));e!=nil{t.Fatal(e)}
  if e=os.Mkdir(filepath.Join(r.Dir,"run.yml"),0755);e!=nil{t.Fatal(e)}
  final:=r.Wait();st,e:=os.Stat(filepath.Join(r.Dir,"run.yml"));t.Logf("Wait state=%s disk record is directory=%v staterr=%v",final.State,st.IsDir(),e)
 })
}

func TestCheckO045Matrix(t *testing.T){
 t.Run("liveness-state-time",func(t *testing.T){
  now:=time.Date(2026,10,11,0,0,0,0,time.UTC);dead:=deadPid(t)
  for _,state:=range []string{StateRunning,StateFinished,StateFailed,StateCancelled,StateTimedOut,StateRunnerLost}{
   for _,age:=range []time.Duration{29*time.Second,30*time.Second,31*time.Second}{
    for _,pid:=range []int{0,dead,os.Getpid()}{
     rec:=Record{State:state,UpdatedAt:now.Add(-age).Format(time.RFC3339),PID:pid}
     want:=state==StateRunning&&age>30*time.Second&&pid==dead
     if got:=RunnerLost(rec,now);got!=want{t.Errorf("%+v got %v want %v",rec,got,want)}
    }
   }
  }
 })
 t.Run("provider-shapes",func(t *testing.T){
  cases:=[]string{"",`[]`,`7`,`null`,`{"type":7}`,`{"type":"assistant","message":null}`,`{"type":"assistant","message":{"content":7}}`,`{"type":"result","usage":null}`,`{"type":"result","usage":[]}`,`{"type":"system","subtype":"future_unknown"}`}
  for _,line:=range cases{t.Logf("shape %s => %+v",line,ParseClaudeLine([]byte(line)))}
 })
 t.Run("version-boundaries",func(t *testing.T){
  for _,v:=range []string{"2.1.258","2.1.259","2.1.260","2.2.0","3.0.0","garbage",""}{
   err:=checkAgentVersion("fake",func(string)(string,error){return v,nil});want:=v=="2.1.258"||v=="garbage"||v==""
   if (err!=nil)!=want{t.Errorf("version %s error %v",v,err)}
  }
 })
 t.Run("heartbeat-write-failure",func(t *testing.T){
  spec:=fakeSpec(t,project(t),"slow");spec.Heartbeat=20*time.Millisecond
  r,e:=Start(spec);if e!=nil{t.Fatal(e)}
  if e=os.Remove(filepath.Join(r.Dir,"run.yml"));e!=nil{t.Fatal(e)}
  if e=os.Mkdir(filepath.Join(r.Dir,"run.yml"),0755);e!=nil{t.Fatal(e)}
  time.Sleep(100*time.Millisecond)
  t.Logf("heartbeat disk failure with agent alive=%v",processAlive(r.rec.PID))
  t.Logf("final state=%s",r.Wait().State)
 })
 t.Run("summary-write-failure",func(t *testing.T){
  r,e:=Start(fakeSpec(t,project(t),"slow"));if e!=nil{t.Fatal(e)}
  if e=os.Mkdir(filepath.Join(r.Dir,"summary.md"),0755);e!=nil{t.Fatal(e)}
  t.Logf("unwritable summary final state=%s",r.Wait().State)
 })
}

func TestCheckO045RealLostRunner(t *testing.T){
 root:=project(t)
 cmd:=exec.Command(os.Args[0]);cmd.Env=append(os.Environ(),"O045_HELPER_ROOT="+root)
 if e:=cmd.Run();e!=nil{t.Fatal(e)}
 entries,e:=os.ReadDir(filepath.Join(root,".savepoint","runs"));if e!=nil||len(entries)!=1{t.Fatalf("entries %v %v",entries,e)}
 rec:=readRecord(t,filepath.Join(root,".savepoint","runs",entries[0].Name()))
 agent,e:=os.FindProcess(rec.PID);if e!=nil{t.Fatal(e)}
 defer killTree(&exec.Cmd{Process:agent})
 t.Logf("runnerPID=%d alive=%v agentPID=%d alive=%v stale-record RunnerLost=%v LostRuns=%v",cmd.Process.Pid,processAlive(cmd.Process.Pid),rec.PID,processAlive(rec.PID),RunnerLost(rec,time.Now()),LostRuns(root,time.Now()))
}
```
