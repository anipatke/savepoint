---
id: C-979
scope: {kind: objective, id: O-045}
result: CLEAR
checked_by: {role: checker, session: check-o045-recheck-20261011}
executed_session: separate-Claude-executor-ced16e7
checked_at: '2026-10-10T21:44:08Z'
health_snapshot: sha256:2abaf34606605ba05e24c8e38b5e53b0c6cec8db9d4d0f1ed06f167e463b765d
reviewed:
  base_commit: 87ee065
  head_commit: ced16e7b20ad42932483c4e5d84228773f3d0d70
  files:
    - .claude/hooks/guard.js
    - AGENTS.md
    - cmd/run.go
    - cmd/run_test.go
    - internal/data/launcher_config.go
    - internal/data/launcher_config_test.go
    - internal/doctor/lost_runs.go
    - internal/doctor/lost_runs_test.go
    - internal/doctor/v2_runtime.go
    - internal/init/claude_guard_test.go
    - internal/run/argv.go
    - internal/run/claude.go
    - internal/run/claude_test.go
    - internal/run/events.go
    - internal/run/fixes_test.go
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
    - main.go
    - templates/project-v2/.claude/hooks/guard.js
  dependencies: []
issues: []
supersedes: C-978
---

# C-979: O-045 Full Objective Re-check

## Closure Map

| Prior Issue | Result | Proof |
|---|---|---|
| I-149 | Verified | Real helper runner exits; agent remains alive; stale saved record now returns RunnerLost=true and LostRuns names it. runner_pid is distinct from agent pid. |
| I-150 | Verified | The three original unknown/malformed shapes now yield unparsed; replay retains exact bytes and counts three. Known housekeeping/plain text/stream deltas remain ignored. |
| I-151 | Verified | Original synthetic URL credential absent from stored command. Independent child echo proves original URL still delivered as argv; raw log preserves it under its separate contract. |
| I-152 | Verified | Heartbeat write failure stops agent and sets Err; final/summary write failures set Err. Launch checks Err and exits failed with an explicit persistence diagnostic. |
| I-153 | Verified | Native Windows full-suite job succeeded for ced16e7 in CI run 38087209247. |

CLEAR: technical clearance. Owner T-130 real-run acceptance remains pending. No Task or Objective status was changed. This checker did not implement the remediation; it arrived as a separate executor commit between review turns. C-978 is unchanged and superseded by this immutable run.

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


All steps completed against C-978's frozen lock. No axis, supported path, dependency layer, or interpretation was added. The child-echo probe is the original I-151 Proof Needed cell (execution unchanged), not scope expansion. No third remediation cycle was started.

## Frozen Scope and Admission Ledger

C-978's numbered scope lock, public matrix, external matrix and workflow inventory remain the definition of complete. The original lock is reproduced unchanged:

Initial Full Check scope lock (2026-10-11)
1. O-045 SC1–7, T-127–130 Done When, named Guardrails. Local build/full tests and native Windows gate. Owner real-run acceptance stays separate.
2. internal/run/**, cmd/run*, main run dispatch, launcher_config*, doctor lost-run wiring, guard tests and AGENTS map. Entry points: ParseRunArgs/RunRun, Launch, Start/Wait, BuildArgv, ParseClaudeLine, RunnerLost/LostRuns, doctor projection.
3. Config/schema/index/Next resolution; version subprocess; record creation, append, heartbeat, atomic replace, final summary; cancel/timeout/process-tree kill; saved-record reader. No lifecycle writes; raw log preserved, derived text sanitised.
4. Matrix rows: CLI refusal/startable verbs/version boundaries; argv placeholders/defaults; provider normal/missing/unknown/malformed/duplicate data; running to each terminal state; liveness fresh/exact30s/over30s crossed live/dead agent and runner; record IO initial/heartbeat/final failure; cancel idle/tree/stubborn and timeout; replay/result-free/nonzero; redirected output and safe diagnostics. Linux exercised; native Windows/macOS/WSL evidence classified separately. Text width, colour/layout, redirects/network, transactions and board launch N/A (not implemented here).
5. Issues require supported path and named promise/guardrail violated with credible consequence. Read-only check, no implementation repair. Separate owner validation from technical clearance. Existing Task waivers not exceptions to criteria.
Workflow: guards (fatal/no writes) -> version (fatal/no writes) -> mkdir/initial atomic record/open events (fatal) -> child start (failed record on failure) -> PID atomic write (fatal/kill) -> append/adapter + heartbeat (record accuracy required) -> wait/cancel/kill -> flush/close -> summary/final atomic record -> visible completion. Oracle: filesystem content, independently known OS process PID, literal expected text, induced filesystem errors.


Before re-check probes, each admitted blocking cell was matched to C-978:

| Re-check item | Prior Issue / claim | Exact frozen cell | Allowed result |
|---|---|---|---|
| Runner identity and stale record | I-149 runner_pid repair | Liveness finite matrix; SC6; real dead-runner/live-agent sequence | Verify or retain I-149 |
| Unknown subtype/null/numeric message | I-150 parser repair | Provider parser, provider-shapes; SC2 | Verify or retain I-150 |
| Stored argv redaction and unchanged execution | I-151 cleanArgv repair | Argv factory / terminal-safe text / raw-derived representations; SC5 | Verify or retain I-151 |
| Heartbeat/final/summary IO failures | I-152 Err + Launch handling | Heartbeat and final IO failures; workflow operations 6 and 9 | Verify or retain I-152 |
| Windows suite | I-153 covering CI | Platform/output, SC7, external timeout/cancel | Verify or retain I-153 |
| Original unaffected cells | Regression and integration | Every remaining original public/external/workflow row | Reconfirm original classification |

## Matrix Completion and Acceptance Coverage

| Original matrix cells | Re-check classification / evidence |
|---|---|
| CLI parser / RunRun; Launch guards; approved verbs | Passed full suite cmd/run_test.go and launch_test.go; shared resolver and unchanged startable-verb table inspected |
| Version boundary; argv factory | Passed original independent 2.1.258/259/260/2.2.0/3.0.0/missing-text cells and TestBuildArgv |
| Provider parser; provider fold | Passed original ParseClaudeLine/replay tests, new fixes_test cases and original overlay shape probes |
| Raw / derived representations; terminal-safe text | Passed fixture byte expectations and sanitisation cases; saved command probe no password; echo subprocess proves original argv delivered |
| Run states; cancel complete sequences | Passed full Linux and native Windows suites, all TestStop and context-cancel cases; race check passes |
| Liveness finite matrix | Passed original 54-cell predicate matrix, runner-specific state/time tests, actual exited helper with live agent, read-only LostRuns |
| File safety / repetition | Passed original Start refusal and Launch no-project-record-write cases; exclusive run-directory creation unchanged |
| Heartbeat and final IO failures | Passed original induced record-directory and summary-directory reproductions through new error API; heartbeat child dies; final/summary Err present |
| Platform/output | Linux suite Passed; Windows full native suite Passed on exact reviewed commit; Darwin cross-build Passed. Physical macOS/WSL execution remains outside available environment as disclosed in C-978; no new native-macOS/WSL CI gate inferred |
| Limits / mutable inputs | Original tail/timeout tests Passed; cap/dedup/copy paths unchanged except command sanitisation; existing source classifications unchanged |

External-boundary cells retain C-978 classification: configured/runtime target and discovery order, startup/unavailable, success/non-success, timeout/cancel, malformed, retry/cleanup/partial side effects, and stored secret safety now pass their required cells. Network/redirect/reuse remains N/A for the local subprocess boundary. No extra network or toolchain edge introduced.

Workflow inventory was retraced in order: guards/version before writes; mkdir/initial atomic record/open append log before Start; PID rewrite and failure kill; append-before-adapter; heartbeat; stop/Wait; flush/close/heartbeat stop; summary/final record; Launch reporting and read-only doctor. Initial/PID failures still return errors. Heartbeat failure now stores the first error and kills the agent tree. Summary/final writes retain their first error, and Launch returns ExitFailed with a named persistence diagnostic after Wait. The process state may truthfully be finished even when persistence failed; Err plus command failure now distinguishes that outcome from saved success, satisfying the original requested error-reporting correction.

Every required cell was reconciled before verdict. Original unaffected cells were rerun through the full suite and reviewed source, not replaced by the new regression tests alone. No remaining in-scope Issue or material unverified criterion remains. The adversarial pass repeated the same parser, stored-argv, liveness and post-start-persistence bypasses. No runner record feeds lifecycle, no project planning record is written, and no board launch behavior is claimed.

| Acceptance | Final classification |
|---|---|
| O-045 SC1 | Proven: off refusal no writes; initial and heartbeat records; surfaced persistence failures |
| SC2 | Proven: mapped fixtures; raw retention and unparsed correction |
| SC3 | Proven against required Linux/native-Windows platform gates: stop/tree/timeout tests pass; shared Unix implementation cross-builds for Darwin |
| SC4 | Proven: nonzero/no-result failed record and sanitised stderr tests |
| SC5 | Proven: exact prompt/root and shell-free execution; clean stored command with original argv delivered |
| SC6 | Proven: guard copies block run; runner_pid lost-run projection read-only |
| SC7 | Proven: fresh full local gate and native Windows CI at reviewed commit |
| T-127 DW1–6 | Proven by original Start tests, record/heartbeat probes, gate and waiver evidence |
| T-128 DW1–6 | Proven by original fixtures/parser/BuildArgv plus I-150/I-151/I-152 reproductions |
| T-129 DW1–5 | Proven by original Stop matrix, new runner-PID tests, native Windows and full gate |
| T-130 DW1–7 technical portion | Proven thin command, refusal/Next/version/context-cancel, guard and doctor tests, full/native gates. Required owner validation remains separate and pending |

## Gates and Native Evidence

- Fresh `make test-full`: exit 0, go1.26.2 linux/amd64, 2026-10-10T21:42Z. Full Go suite with reports and Linux/Darwin/Windows builds.
- Fresh `make build`: exit 0; `git diff --check`: exit 0.
- `go test -race -count=1 ./internal/run ./cmd ./internal/doctor`: exit 0.
- Original independent overlay harness rerun (adapted only to observe Err and runner_pid): exit 0. Full output below.
- [Native Windows full suite](https://github.com/anipatke/savepoint/actions/runs/38087209247/job/114316104586): success, 2026-10-10T21:20:05Z–21:22:52Z, head ced16e7b20ad42932483c4e5d84228773f3d0d70. Workflow uses windows-latest and `go run ./internal/buildtool test -json -count=1 ./...`. Successful Linux `make ci` job at the same head also confirmed. Read `gh run view` JSON and native job log; this is actual runtime evidence, not cross-compilation.
- Configured quality_gates remain unset, so no invented lint/typecheck command was added.
- Official health collection after full gate: sha256:2abaf34606605ba05e24c8e38b5e53b0c6cec8db9d4d0f1ed06f167e463b765d, does not block clearance. Optional vulnerability collection failed again; that does not establish a vulnerability. Other configured optional measurements collected without blockers.
- Reviewed files and claimed fixtures exist; new `internal/run/fixes_test.go` exists. No phantom artifact. No implementation change during Check.

## Prior Owner Decisions and Closure

All four Tasks remain done with recorded optional-Check waivers; those do not exempt any requirement here. No scoped owner acceptance or criterion exception has been added since C-978, so no carried_forward entry is required. The approved absent-builder/default behavior is unchanged. T-130 still declares owner_validation.required with no accepted_check. The owner must validate the real terminal Claude run and record acceptance before Objective closure; CLEAR does not record that decision or close the Objective.

Design reconciliation is unchanged: package responsibility and file-first/lifecycle boundaries agree; T-130 foreground run and deferred O-046 board detachment remain explicit; full Design reconciliation documentation is assigned to O-047. No new acceptance perimeter was invented for deferred work.

## Issues, Materiality and Observations

No materiality actions remain for this Check. I-149–I-153 are verified by this CLEAR record and receive append-only closure history. No new Issues opened. Observations carried from C-978: physical macOS/WSL/real-terminal execution not performed here; optional vulnerability collection incomplete; duplicated sanitiser/platform code remains advisory. None introduces a new blocking cell on this re-check.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches**
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries** — originally reported persistence errors now retained and checked by Launch.
- [ ] STYLE-07 **One source of truth** — documented sanitiser/platform duplication and repeated timeout defaults remain unchanged, advisory.
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs**

## Reviewed Fingerprints

- `.claude/hooks/guard.js` sha256 `ae8dde236eb5e58dbfe1e9d18bdd7e9d29395cafdae4942c40041f0a993f181d`
- `AGENTS.md` sha256 `b80d3f9244cca3e56a62e2134e1d5381c512c9cf608485ff5bdee8f69aa60e5c`
- `cmd/run.go` sha256 `a7ab3962b93a4146b6828462ce7c2453fd9847fb9c312d917973778e1986b9f9`
- `cmd/run_test.go` sha256 `f1d13f39442177d2b6922f8db065d7740112278a7ed39b3d656be9d99273ef8f`
- `internal/data/launcher_config.go` sha256 `09db55d338d5f25e99646f345c583a83b210a7b35ecff4caa82c082a00fd0eec`
- `internal/data/launcher_config_test.go` sha256 `633107638305d0781aa98765c0b0c7ad6825dc9964d0ab1e2e757b03c4661feb`
- `internal/doctor/lost_runs.go` sha256 `9448d72706162a485e97e0a23b75f7e1505f92aa37fc29fa5d26f97bf91bb0ce`
- `internal/doctor/lost_runs_test.go` sha256 `a77144cab49519aa468a3e2be541dec2730e5eb7efbe3c45ceab335873ad9e34`
- `internal/doctor/v2_runtime.go` sha256 `96027b1556f18c745af7a588107071d2aac4528766ecf029bc6d6fd81e711fb2`
- `internal/init/claude_guard_test.go` sha256 `8c342d2ce329573309b9ef252f86e0d4c31703dff716988807b06e5b08ea88c2`
- `internal/run/argv.go` sha256 `2cbdc3f8ec0beaabececc5ef59580331222dffe786bdc3e72cced5a631612402`
- `internal/run/claude.go` sha256 `c5dc4696ab98debbafc7c263a5c6179dbc81ed179fafb9c419c810e1e65fcba2`
- `internal/run/claude_test.go` sha256 `3797930bd0cfdb431a1b68d42a4080fe3af2151807c0c3cbbe92709eea8ea811`
- `internal/run/events.go` sha256 `e0572d36162a37b3cba5dd99d8f29df07bfe0f82e32880584a901522c6fb15df`
- `internal/run/fixes_test.go` sha256 `28f18c9036a56da9fdcac99e6c4ff8427fb45bedb7700e7239f5bdc61d00d883`
- `internal/run/launch.go` sha256 `2f45de704c0fe44c01eaa49b215a4b9d797416085c8f4edd1effb4d2fcff70fd`
- `internal/run/launch_test.go` sha256 `be3f1d371a7dfe116bda3743e8532ff064911ad3e219bc0194060a537bc6a58c`
- `internal/run/lost.go` sha256 `547e019bfdddb41107fc67778eff465b34e98c73712a362f4b96c9500b81cf17`
- `internal/run/proc_unix.go` sha256 `f8c1f37fc0698a36fb753126917694d7460125f5e8345b12f7522ede50575bad`
- `internal/run/proc_windows.go` sha256 `6d56f406c2cf7b04ecfd2b5c6780228f9e5938badfb2f21df0e64a7dbf1133cc`
- `internal/run/record.go` sha256 `16e81a7ecdc350f3a59283ceb02ff16b06d50938d6aebb163c590f72c9df8dd7`
- `internal/run/run.go` sha256 `d6121fbb66279d70c7e96c4ea22cd0b9ea70e9c46c4776abcccf98fd1a440646`
- `internal/run/run_test.go` sha256 `b79238dc44d6b00b51933e05fc60ccfbb05296eb4b2a59b38c9d80444f418037`
- `internal/run/sanitize.go` sha256 `afde49e30d9f5073a523cf0363db07f332afe55e005d8311b2ae2f5b17fa0930`
- `internal/run/stop.go` sha256 `656bb0545ae94382f96f2a47d6a3ec40bd8ea214cd98e0bd8ea246cdf9cab9f6`
- `internal/run/stop_test.go` sha256 `a7dc2e2cdbd37b2d030372cc0e0908e1cc0d639a46b388afdbd5932681fc13c0`
- `internal/run/testdata/claude-2.1.296/edit-and-denial.jsonl` sha256 `b5ca1539a4236aafbf0a1491e7ee812db1ce59337a19065e26e1d5abfd5d0a15`
- `internal/run/testdata/claude-2.1.296/failing-command.jsonl` sha256 `a7d644423bc71fa01abb7b3aa8985315d2909be392a63f6193bfaba5a6840c1c`
- `main.go` sha256 `4b7d84bf3b48dac3b1e045d305dc98af9696ace98836d988510cb29fa6063f5b`
- `templates/project-v2/.claude/hooks/guard.js` sha256 `ae8dde236eb5e58dbfe1e9d18bdd7e9d29395cafdae4942c40041f0a993f181d`

## Independent Probe Output

The old synthetic PID-only record still returns false by the deliberate legacy fallback: it contains no runner_pid and cannot establish runner identity. The actual helper process case below is the original supported-path reproduction and now reports lost correctly. Final/summary process state stays finished, while Err is nonnil and Launch exits failed; this is explicit process-versus-persistence evidence, not hidden saved success.

```text
=== RUN   TestCheckO045Independent
=== RUN   TestCheckO045Independent/lost-runner-live-agent
    run_test.go:234: runner died, agent PID alive: RunnerLost=false (expected true for dead runner)
=== RUN   TestCheckO045Independent/unknown-system
    run_test.go:237: unknown system subtype events=[{Kind:unparsed ToolUseID: Tool: Text: ExitCode:<nil> Failed:false SessionID: Model: Usage:map[] TotalCostUSD:<nil> Denials:[]}] (expected unparsed)
=== RUN   TestCheckO045Independent/stored-command-credentials
    run_test.go:242: run.yml contains URL password=false
=== RUN   TestCheckO045Independent/final-write-failure
    run_test.go:249: final-write Err=write /tmp/TestCheckO045Independentfinal-write-failure2016383796/001/.savepoint/runs/20261010T214237Z-T-001/run.yml: rename /tmp/TestCheckO045Independentfinal-write-failure2016383796/001/.savepoint/runs/20261010T214237Z-T-001/.run-3596578508 /tmp/TestCheckO045Independentfinal-write-failure2016383796/001/.savepoint/runs/20261010T214237Z-T-001/run.yml: file exists
    run_test.go:249: Wait state=finished disk record is directory=true staterr=<nil>
--- PASS: TestCheckO045Independent (0.62s)
    --- PASS: TestCheckO045Independent/lost-runner-live-agent (0.00s)
    --- PASS: TestCheckO045Independent/unknown-system (0.00s)
    --- PASS: TestCheckO045Independent/stored-command-credentials (0.01s)
    --- PASS: TestCheckO045Independent/final-write-failure (0.61s)
=== RUN   TestCheckO045Matrix
=== RUN   TestCheckO045Matrix/liveness-state-time
=== RUN   TestCheckO045Matrix/provider-shapes
    run_test.go:268: shape  => []
    run_test.go:268: shape [] => [{Kind:unparsed ToolUseID: Tool: Text: ExitCode:<nil> Failed:false SessionID: Model: Usage:map[] TotalCostUSD:<nil> Denials:[]}]
    run_test.go:268: shape 7 => [{Kind:unparsed ToolUseID: Tool: Text: ExitCode:<nil> Failed:false SessionID: Model: Usage:map[] TotalCostUSD:<nil> Denials:[]}]
    run_test.go:268: shape null => [{Kind:unparsed ToolUseID: Tool: Text: ExitCode:<nil> Failed:false SessionID: Model: Usage:map[] TotalCostUSD:<nil> Denials:[]}]
    run_test.go:268: shape {"type":7} => [{Kind:unparsed ToolUseID: Tool: Text: ExitCode:<nil> Failed:false SessionID: Model: Usage:map[] TotalCostUSD:<nil> Denials:[]}]
    run_test.go:268: shape {"type":"assistant","message":null} => [{Kind:unparsed ToolUseID: Tool: Text: ExitCode:<nil> Failed:false SessionID: Model: Usage:map[] TotalCostUSD:<nil> Denials:[]}]
    run_test.go:268: shape {"type":"assistant","message":{"content":7}} => [{Kind:unparsed ToolUseID: Tool: Text: ExitCode:<nil> Failed:false SessionID: Model: Usage:map[] TotalCostUSD:<nil> Denials:[]}]
    run_test.go:268: shape {"type":"result","usage":null} => [{Kind:final ToolUseID: Tool: Text: ExitCode:<nil> Failed:false SessionID: Model: Usage:map[] TotalCostUSD:<nil> Denials:[]}]
    run_test.go:268: shape {"type":"result","usage":[]} => [{Kind:final ToolUseID: Tool: Text: ExitCode:<nil> Failed:false SessionID: Model: Usage:map[] TotalCostUSD:<nil> Denials:[]}]
    run_test.go:268: shape {"type":"system","subtype":"future_unknown"} => [{Kind:unparsed ToolUseID: Tool: Text: ExitCode:<nil> Failed:false SessionID: Model: Usage:map[] TotalCostUSD:<nil> Denials:[]}]
=== RUN   TestCheckO045Matrix/version-boundaries
=== RUN   TestCheckO045Matrix/heartbeat-write-failure
    run_test.go:282: heartbeat disk failure with agent alive=false
    run_test.go:283: final state=failed error=write /tmp/TestCheckO045Matrixheartbeat-write-failure874367803/001/.savepoint/runs/20261010T214238Z-T-001/run.yml: rename /tmp/TestCheckO045Matrixheartbeat-write-failure874367803/001/.savepoint/runs/20261010T214238Z-T-001/.run-3668095251 /tmp/TestCheckO045Matrixheartbeat-write-failure874367803/001/.savepoint/runs/20261010T214238Z-T-001/run.yml: file exists
=== RUN   TestCheckO045Matrix/summary-write-failure
    run_test.go:288: unwritable summary final state=finished error=write /tmp/TestCheckO045Matrixsummary-write-failure1266361332/001/.savepoint/runs/20261010T214238Z-T-001/summary.md: rename /tmp/TestCheckO045Matrixsummary-write-failure1266361332/001/.savepoint/runs/20261010T214238Z-T-001/.run-3198255211 /tmp/TestCheckO045Matrixsummary-write-failure1266361332/001/.savepoint/runs/20261010T214238Z-T-001/summary.md: file exists
--- PASS: TestCheckO045Matrix (0.72s)
    --- PASS: TestCheckO045Matrix/liveness-state-time (0.00s)
    --- PASS: TestCheckO045Matrix/provider-shapes (0.00s)
    --- PASS: TestCheckO045Matrix/version-boundaries (0.00s)
    --- PASS: TestCheckO045Matrix/heartbeat-write-failure (0.11s)
    --- PASS: TestCheckO045Matrix/summary-write-failure (0.61s)
=== RUN   TestCheckO045RealLostRunner
    run_test.go:300: runnerPID=139 alive=false agentPID=145 alive=true stale-record RunnerLost=true LostRuns=[20261010T214138Z-T-001]
--- PASS: TestCheckO045RealLostRunner (0.01s)
=== RUN   TestCheckO045ArgvExecutionUnchanged
    run_test.go:310: child received original URL; saved command redacted
--- PASS: TestCheckO045ArgvExecutionUnchanged (0.01s)
PASS
ok  	github.com/opencode/savepoint/internal/run	1.359s
```

## Harness Adaptation

C-978 contains the original harness. This re-check reused it as a Go overlay on current run_test.go, adding Err logging at final/summary/heartbeat failures. Helper Start now naturally writes runner_pid. For I-151's unchanged-execution proof, a helper mode prints its os.Args[1:] to stdout and exits; the probe checks raw log equals the original synthetic credential URL plus newline while run.yml lacks the password. Repo code and tests were not modified for this probe.
