---
id: O-045
title: Record and run one agent as a detached process
status: in_progress
depends_on: []
release: G-003
priority: high
rank: 1
---

# O-045: Record and run one agent as a detached process

## Outcome

`savepoint run --task T-### [dir]` starts one Claude Code session for the named Task in the project root, owns the process, and writes `.savepoint/runs/<timestamp>-<task>/run.yml`, `events.jsonl`, and `summary.md`. The record states what the runner observed: state, exit code, session ID, denials, usage, heartbeat, and the agent's final message under `agent_reported`. Cancel by file, timeout, and process-tree kill work on Linux, macOS, WSL2, and Windows. The feature is off unless `agent_launcher.enabled` is true. Nothing in the record is read by any lifecycle gate.

## Why

"Files are state" and "agent output is evidence, not truth" are met only if the runner, not the board, owns the process and writes the record. Making the runner a subcommand rather than board code is what lets a run outlive the board (review §3, option B, `docs/reviews/agent-orchestration-pass-1.md`).

## Success Conditions

1. Running the subcommand with the feature flag off or absent exits 2 with a one-line reason and writes nothing; with it on, `run.yml` exists before the agent starts and is updated at least every 5 s while it runs.
2. The Claude adapter maps `system/init`, `assistant` `tool_use`, `user` `tool_result`, `system/permission_denied`, and `result` lines to display events and record fields; unknown lines are kept verbatim on disk and shown as unparsed, never dropped.
3. A `cancel` file in the run directory ends the run within 15 s on every supported platform with state `cancelled` and no surviving agent or grandchild process; a timeout does the same with state `timed_out`.
4. A non-zero exit without a `result` line records state `failed` and a sanitised stderr tail.
5. The argv is built from `agent_launcher.builder` with `{{prompt}}` = the Next line verbatim and `{{project_root}}`; no shell is involved; credentials in URLs are stripped from stored text as the health runner already does.
6. `guard.js` still blocks agents from running `savepoint run`; `doctor` lists runs whose heartbeat is stale and whose PID is dead as `runner_lost`.
7. Configured gates pass, including `make test-full`, with native Windows evidence named.

## Architectural Considerations

Settled in Confirmed Design below:

- Process handling is generalised from `internal/codehealth/runner.go`, `runner_unix.go`, `runner_windows.go` into a new `internal/run` package (ARCH-04); the board must not import `codehealth` for this.
- `events.jsonl` is append-only (`O_APPEND`), one provider line per line; `run.yml` is rewritten atomically (FS-04, FS-05).
- Record shape: `provider`, `command` (argv), `cwd`, `task`, `next_line`, `started_at`, `updated_at`, `pid`, `state` ∈ {running, finished, failed, cancelled, timed_out, runner_lost}, `exit_code`, `session_id`, `denials[]`, `usage` (verbatim provider fields), `agent_reported`.
- Default Claude flags: `-p`, `--output-format stream-json`, `--verbose`, `--permission-mode acceptEdits`, `--allowedTools` derived from `quality_gates` as `Bash(<first word> *)`, `--permission-prompts none`. No `--bare`. Version floor 2.1.259 (the `--permission-prompts` flag).
- Lifecycle vocabulary stays in `internal/data` (DATA-02); run states are a separate vocabulary owned by `internal/run` and never feed `ResolveNext` or any gate.

## Confirmed Design

Owner, 2026-10-11, design confirmation after the review and plan:

- **Opt-in.** Off by default behind `agent_launcher.enabled`; the Next line stays copyable and hand-started sessions in any terminal remain the primary path. The subcommand refuses to run when the flag is off.
- **Process model.** Option B: a detached `savepoint run` owns the agent; the board (O-046) only reads files. In this Objective the subcommand is run by the owner from a terminal; no board change.
- **Provider.** Claude Code only, hard-coded behind the `agent_launcher.builder` command spec so a later provider is a config and adapter change.
- **Cancel and liveness.** Cancel is a file the runner watches, so no cross-process signals are needed on Windows; SIGINT first on Unix, then the existing process-group or `taskkill` tree kill. Liveness is heartbeat plus PID.
- **What the runner never does.** It never writes a Task, Objective, Check, or router record, and never summarises or interprets the agent's output beyond labelling it `agent_reported`.
- **Tasks.** Four sequential Tasks: runner core with a fake provider; Claude event adapter; cancel, timeout and liveness across platforms; the subcommand, flag gate, guard and doctor wiring.

## Boundaries

**In scope:**
- `internal/run` (spec, record writer, process ownership, cancel, timeout, heartbeat, Claude adapter)
- `cmd/run.go` and `main.go` wiring; `doctor` line for lost runners
- Tests with a fake `claude` that replays a recorded `stream-json` fixture; native Windows CI evidence

**Out of scope:**
- Any board change (O-046)
- Any Task, Objective, Check, or router write by the runner
- Any provider other than Claude Code; any permission UI; Check runs; worktrees
