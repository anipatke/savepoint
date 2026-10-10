---
id: T-129
title: Stop a run cleanly by cancel file, timeout, or lost runner
objective: O-045
status: done
depends_on: [{task: T-128, requires: clear}]
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: planning-g003-20261011}
check_waiver:
    task: T-129
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-10T20:59:52Z"
---

# Stop a run cleanly by cancel file, timeout, or lost runner

## Outcome

A run can be stopped from outside the runner by creating a `cancel` file in its directory, ends itself on a configurable timeout, and can be recognised as lost when the runner process died without writing a final state. On every supported platform the whole process tree ends, nothing survives, and the record says `cancelled`, `timed_out`, or (as an inference by a reader) `runner_lost`.

## User Check

Not required. Technical Task: the Full Objective Check verifies the platform evidence named below.

## Done When

1. A `cancel` file appearing in the run directory is noticed within 1 s (poll or fsnotify); on Unix the runner sends SIGINT to the agent, waits up to 10 s, then uses the process-group kill; on Windows it uses the tree kill directly. State becomes `cancelled` within 15 s and no agent or grandchild process survives (O-045 condition 3, CFG-02).
2. A timeout from `agent_launcher.timeout` (default 30 min; parsed like `quality_gates.gate_timeout`) stops the run the same way with state `timed_out`.
3. A reader-side helper reports a run as `runner_lost` when `state` is `running`, `updated_at` is older than 30 s, and `pid` is not alive; PID liveness is implemented for Unix and Windows and the helper never writes to `run.yml` (DATA-04 spirit: the inference is reportable, not healed).
4. Cancel and timeout record `ended_at` and leave `events.jsonl` intact; a cancel while the agent is between turns and a cancel mid-command both end cleanly.
5. Tests cover cancel, timeout, and lost-runner on Unix and on Windows build tags, using the fake provider with a child that spawns a grandchild which must not survive (TEST-01, TEST-02, TEST-04). `make build && make test-fast` pass (TEST-08). Per-criterion evidence follows AGENTS.md's Verification Policy, including an explicit owner waiver if the optional Task Check is skipped.

## Context Files

`internal/run/` (as landed by T-127 and T-128), `internal/codehealth/runner_unix.go`, `internal/codehealth/runner_windows.go`, `internal/codehealth/runner_alive_unix_test.go`, `internal/codehealth/runner_alive_windows_test.go`, `internal/doctor/gates.go` (timeout parsing only), `.savepoint/objectives/O-045-record-and-run-one-agent-as-a-detached-process/Objective.md`, `docs/reviews/agent-orchestration-pass-1.md` (section 5, failure cases).

## Design References

Design §11 (failure modes), §13 (testing, Windows CI). O-045 Confirmed Design (cancel and liveness). Claude Code SIGINT/SIGTERM behaviour as cited in the review.

## Guardrails

FS-05, ARCH-03, ARCH-04, CFG-02, CFG-03, TEST-01, TEST-02, TEST-04, TEST-05, TEST-08.

## Implementation Plan

1. Confirm the process setup from T-127 exposes a cancel hook on both platforms; return REPLAN REQUIRED if Windows tree kill is not reachable from the runner.
2. Add the cancel-file watcher and the SIGINT-then-kill sequence on Unix; tree kill on Windows.
3. Add the timeout from config with the default.
4. Add PID liveness (`kill -0` equivalent on Unix; `OpenProcess` on Windows) and the `runner_lost` helper.
5. Tests with a grandchild-spawning fake provider on both platforms; push the branch so the native Windows CI job runs.

## Boundaries

No CLI subcommand, no board change, no doctor output yet (next Task), no change to `codehealth` behaviour.

## Technical Verification

Focused `make test-focused TEST=...` during iteration. Platform-sensitive: handoff uses a fresh `make test-full`, and the Full Objective Check needs the native Windows CI job's result for this branch, produced by the executor pushing the branch before handoff.

## Technical Evidence

### Per-criterion outcomes

1. **Met on Linux; Windows not yet run natively.** `startStopper` (`internal/run/stop.go`) polls for `cancel` every 200 ms. Unix: SIGINT to the agent, up to `Grace` (default 10 s), then group SIGKILL; the group kill also runs after the agent exits on SIGINT so a grandchild cannot outlive it. Windows: `interruptProcess` returns false, so the tree kill (`taskkill /T /F`) runs directly. `TestStop_cancelFileEndsAnAgentBetweenTurns`, `..MidCommandAndItsGrandchild`, `..IgnoresTheInterrupt` pass: state `cancelled`, agent and grandchild gone, all inside the 15 s budget (they finish in under 1 s).
2. **Met.** `Spec.Timeout` (default `DefaultTimeout` 30 min) stops the run the same way with `timed_out`; `TestStop_timeoutEndsTheRunAndItsGrandchild`. `agent_launcher.timeout` is read by `AgentLauncher.RunTimeout()` in `internal/data`, parsed with `time.ParseDuration` like `gate_timeout`; a non-positive or invalid value gives the 30 min default (`TestAgentLauncherRunTimeout_defaultsUnlessAPositiveDurationIsSet`). Nothing calls it yet: T-130 wires the subcommand.
3. **Met.** `RunnerLost(rec, now)` is true only for `running`, `updated_at` older than 30 s, and a dead `pid`; `processAlive` is `kill(pid, 0)` on Unix and `OpenProcess`/`WaitForSingleObject` on Windows. `TestRunnerLost_needsRunningStaleHeartbeatAndDeadPid` covers each condition; `TestRunnerLost_neverWritesTheRecord` shows `run.yml` is byte-identical afterwards.
4. **Met.** Every final state now records `ended_at` (cancel and timeout included); `assertStopped` checks `events.jsonl` still holds the agent's line exactly. Cancel with no child (idle, between turns) and with a running child (mid-command) both end cleanly.
5. **Met on Linux.** Tests use the fake provider (`idle`, `tree`, `stubborn`, grandchild `sleeper`) and are build-tag free, so the same tests run on Windows; per-platform code is in `proc_unix.go` / `proc_windows.go`. `GOOS=windows go vet ./internal/run ./internal/data` is clean. `-race -count=3` on `internal/run` passed. TEST-01/02/04: fakes only, no real `claude`. Native Windows execution is **not done** (see limitations).

### Commands run

- `make build && make test-fast` at 2026-10-10T20:58Z, go1.26.2 linux/amd64: pass.
- `make test-full` (platform-sensitive handoff), fresh, 2026-10-10T20:59Z: pass (exit 0, includes build-linux/darwin/windows).
- `go test -race -count=3 ./internal/run`: pass. `GOOS=windows go vet ./internal/run ./internal/data`: clean.

### Files

- Read (Context Files): `internal/run/*`, `internal/codehealth/runner_unix.go`, `runner_windows.go`, `runner_alive_unix_test.go`, `runner_alive_windows_test.go`, `internal/doctor/gates.go` (timeout parsing), this Task's Objective. Not read: `docs/reviews/agent-orchestration-pass-1.md`.
- Extra reads, logged: `internal/data/launcher_config.go` and `launcher_config_test.go`, because criterion 2 names `agent_launcher.timeout` and the field did not exist yet; `.github/workflows/ci.yml` to confirm the native Windows job runs `./...`.
- Changed: `internal/run/stop.go` (new), `stop_test.go` (new), `proc_unix.go`, `proc_windows.go`, `record.go` (states, `ended_at`), `run.go` (Spec `Timeout`/`Grace`, stopper wiring), `run_test.go` (fake modes), `internal/data/launcher_config.go` and `launcher_config_test.go` (timeout field, `RunTimeout`).

### Limitations

- **Native Windows result is missing.** The branch has not been pushed (pushing is the owner's call), so the Windows CI job has not run these tests. The Full Objective Check needs that result. The Windows tree kill is only type-checked here.
- `pid` in `run.yml` is the agent's pid, not the runner's. If the runner dies but the agent keeps going, `RunnerLost` stays false because the pid is alive. That follows the Done When wording; see Drift Notes.
- Zombie processes in a container with no init could look alive to `kill(pid, 0)`; the tests poll up to 3 s and did not flake here.

## Drift Notes

Record architecture deltas and reconcile through the planner before Check.

- `ended_at` is written for every final state, not only cancel and timeout.
- `agent_launcher.timeout` was added to `internal/data` (not in the Task's file list).
- Planner question: the heartbeat writer is the runner, but `pid` is the agent's. A runner that dies while its agent lives is not reported as `runner_lost`. Consider recording `runner_pid` before T-130's doctor line relies on this.
