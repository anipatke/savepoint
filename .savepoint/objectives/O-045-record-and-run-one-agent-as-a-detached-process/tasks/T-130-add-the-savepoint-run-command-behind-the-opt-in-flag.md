---
id: T-130
title: Add the savepoint run command behind the opt-in flag
objective: O-045
status: done
depends_on: [{task: T-129, requires: clear}]
owner_validation:
    required: true
    accepted_check: ""
planned_by: {role: planner, session: planning-g003-20261011}
check_waiver:
    task: T-130
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-10T21:07:00Z"
---

# Add the savepoint run command behind the opt-in flag

## Outcome

An owner can type `savepoint run --task T-### [dir]` in a terminal and get a real Claude Code run for that Task, recorded under `.savepoint/runs/`, but only when `agent_launcher.enabled` is true. With the flag off or absent the command exits 2 with one plain line and writes nothing. Agents still cannot run it, and `doctor` names any run whose runner was lost.

## User Check

Owner validation required. In this repository with `agent_launcher.enabled: true` set locally, run the command on a planned Task from a terminal, watch `events.jsonl` grow, confirm the Task's own record changes appear on the board as today, then set the flag back to false and confirm the command refuses and `savepoint resume` and the board are unchanged. This is the first real run; the owner's judgement of "did it do what I expected and nothing more" is the acceptance.

## Done When

1. `cmd/run.go` parses `--task T-###` and an optional directory, nothing else, in the thin style of `cmd/resume.go` (ARCH-01); unknown flags and a missing Task are named errors with exit 2.
2. With `agent_launcher.enabled` false or absent, the command prints one line saying the feature is off and how to turn it on, exits 2, and creates no file or directory (O-045 condition 1, FS-03 spirit).
3. With it on, the command resolves the Task through the strict V2 index, refuses a Task that is `done` or whose `ResolveNext` verb is not Start, Build, or Test (reporting the verb), computes the Next line through `internal/resume` for that Task, and starts the run with the argv from T-128; the Claude version floor 2.1.259 is checked and a lower version is a named error before any write (CFG-01).
4. The command runs the agent in the foreground of the terminal that invoked it (detaching is the board's job in O-046); Ctrl-C in that terminal behaves as a cancel.
5. `guard.js` in both the scaffold and this repository still blocks an agent running `savepoint run`; a test asserts it (TPL-01, TPL-02). The agent-runnable command list in AGENTS.md and Design §6 is unchanged.
6. `savepoint doctor` prints one line per run that T-129's helper reports as `runner_lost`, naming the run directory; projects with no `runs/` directory print nothing new (TPL-03 spirit).
7. `make test-full` passes fresh, with native Windows CI evidence named. Per-criterion evidence follows AGENTS.md's Verification Policy, including an explicit owner waiver if the optional Task Check is skipped; owner validation is recorded after the User Check above.

## Context Files

`main.go`, `cmd/resume.go`, `cmd/health.go`, `internal/run/` (as landed by T-127, T-128, T-129), `internal/data/config.go`, `internal/data/launcher_config.go`, `internal/data/next.go`, `internal/resume/resume.go`, `internal/doctor/doctor.go`, `templates/project-v2/.claude/hooks/guard.js`, `.claude/hooks/guard.js`, `internal/init/claude_guard_test.go`, `.savepoint/objectives/O-045-record-and-run-one-agent-as-a-detached-process/Objective.md`.

## Design References

Design §6 (CLI surface and agent-runnable commands), §11 (failure modes), §13 (testing). O-045 Confirmed Design (opt-in, what the runner never does).

## Guardrails

FS-03, FS-06, ARCH-01, ARCH-03, CFG-01, CFG-02, CFG-03, TPL-01, TPL-02, TPL-03, TEST-01, TEST-02, TEST-08.

## Implementation Plan

1. Confirm T-127–T-129 expose `Start`, the argv builder, and the lost-runner helper; return REPLAN REQUIRED if the record shape differs from O-045's Confirmed Design.
2. Add `cmd/run.go` parsing and the `main.go` dispatch with an injected runner, like `resume`.
3. Flag gate first, then Task resolution and verb check, then version check, then start.
4. Doctor line for lost runners.
5. Guard tests for both `guard.js` copies.
6. Fresh `make test-full`; push for Windows CI.

## Boundaries

No board change; no detaching from the board; no Task, Objective, Check, or router write; no second provider; no change to the agent-runnable command list; no Design.md change (O-047 reconciles it).

## Technical Verification

Focused `make test-focused TEST=...` during iteration. Platform-sensitive and new CLI surface: handoff uses a fresh `make test-full`, and the Full Objective Check needs the native Windows CI job's result for this branch before it starts. Owner validation is recorded after the User Check.

## Technical Evidence

Executor session: T-130 build, 2026-10-11. Claims for a fresh Check to verify.

**Per criterion**

1. `cmd/run.go` parses `--task T-###` plus one optional directory; unknown flag, missing `--task`, missing value, non-Task ID and two directories are named errors, exit 2, and never call the runner. Evidence: `cmd` `TestParseRunArgs`, `TestRunRun_*`; `TestPackage_staysThinNoDomainImports` still passes (ARCH-01).
2. Flag off or absent: `run.Launch` returns exit 2 with one line naming `agent_launcher.enabled: true`, before the Task, version or any write. Evidence: `internal/run` `TestLaunch_offWritesNothingAndSaysHowToTurnItOn` (absent and false; asserts no `runs/` directory).
3. Flag on: Task resolved through the strict V2 index; done Task, unknown Task and a Task whose verb is not Start/Build/Test (audit stage) are refused with the reason, no write. Next line comes from `data.ResolveNext` with the Task selected, rendered by `internal/resume`; argv from `BuildArgv`; version floor 2.1.259 checked first (2.1.258, 1.9.999, unreadable and missing agent all refused before any write). Evidence: `TestLaunch_refusesATaskThatCannotStartBeforeWriting`, `TestLaunch_refusesATaskWhoseNextStepIsNotAStartReportingTheVerb`, `TestLaunch_refusesAnOldOrUnreadableAgentBeforeWriting`, `TestLaunch_buildsTheNextLineAndRootIntoTheConfiguredBuilder`, `TestLaunch_enabledWithoutABuilderUsesTheDefaultClaudeCommand`.
4. Foreground run; ending the context (main wires `signal.NotifyContext` for Ctrl-C) writes the cancel file and the run ends `cancelled`. Evidence: `TestLaunch_endingTheContextCancelsTheRun`; failing agent exits 1 (`TestLaunch_aFailingAgentExitsNonZeroAndRecordsFailed`). A real terminal Ctrl-C was not exercised (limitation).
5. `guard.js` already blocked every `savepoint` command except resume, create-task and health check, so no guard change was needed. Added assertions that `savepoint run`, `./savepoint run` and `npx savepoint run` are blocked (`TestGuard_savepointCommands`) and that `.claude/hooks/guard.js` is byte-identical to the scaffold copy (`TestGuard_repositoryCopyMatchesTheScaffold`). Agent-runnable command list in AGENTS.md and Design unchanged.
6. `doctor` prints a `[runner_lost]` line per lost run via `run.LostRuns` + `RunnerLost`, as advice (`Pending Semantic Review`, does not fail doctor); projects with no or empty `runs/` print nothing new. Evidence: `internal/doctor` `TestRunV2Checks_namesALostRunAsAdviceOnly`, `TestRunV2Checks_aProjectWithoutRunsPrintsNothingNew`; `internal/run` `TestLostRuns_namesOnlyRunsWhoseRunnerIsGone`.
7. `make build && make test-fast` exit 0; fresh `make test-full` exit 0 (Linux, go1.26.2, 2026-10-11); `GOOS=windows go vet ./internal/run ./cmd .` clean. Native Windows CI result for this branch: **not yet available** (needs a push, which was not done); the Full Objective Check must wait for it. Owner validation (User Check) is pending.
Also: `TestLaunch_neverWritesAProjectRecord` shows the run changes nothing under `.savepoint/` outside `runs/`.

**Decision recorded (owner, this session):** `agent_launcher.builder` is now optional when enabled; an absent builder means `DefaultAgentBuilder`. Previously `Validate` required it, which made the T-128 default unreachable. Changed `internal/data/launcher_config.go`; replaced the "enabled without builder" error case and `TestLauncherEnabledMissingBuilderFailsRead` with `TestLauncherEnabledWithoutBuilderReadsAndUsesTheDefault`. A builder that is present still needs a command.

**Commands run:** `go build ./...`, `go vet ./...`, `go test ./cmd ./internal/run ./internal/doctor ./internal/init ./internal/data .`, `make build`, `make test-fast`, `make test-full`, `GOOS=windows go vet ./internal/run ./cmd .`.

**Files read:** the listed Context Files, plus extra reads: `internal/run/*.go` (run.go, stop.go, argv.go, record.go, proc_unix.go; part of T-127-T-129's landed code, named in Context Files), `internal/data/router_v2.go`, `project.go`, `objective_v2.go`, `task_v2.go`, `runtime.go` (types and `ResolveTarget` needed to build the selection), `internal/doctor/checks_test.go`, `report_test.go`, `health_snapshot_refs.go` (fixtures and problem pattern), `.savepoint/Guardrails.md` (rule rows).

**Files changed:** new `cmd/run.go`, `cmd/run_test.go`, `internal/run/launch.go`, `internal/run/launch_test.go`, `internal/run/lost.go`, `internal/doctor/lost_runs.go`, `internal/doctor/lost_runs_test.go`; edited `main.go`, `internal/doctor/v2_runtime.go`, `internal/data/launcher_config.go`, `internal/data/launcher_config_test.go`, `internal/init/claude_guard_test.go`, `AGENTS.md` (Codebase Map rows for `cmd/` and `internal/run/`), this Task.

**Limitations:** no real Claude run or terminal Ctrl-C was exercised; no Windows CI result yet; `AGENTS.md` and `launcher_config*.go` already carried uncommitted edits from earlier Tasks, so the diff there mixes T-128 and T-130.

## Drift Notes

Record architecture deltas and reconcile through the planner before Check.
