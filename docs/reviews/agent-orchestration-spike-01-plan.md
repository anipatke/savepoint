# Agent Orchestration Spike 01 — Goal and Objective plan

Date: 2026-10-11 · Source: `docs/reviews/agent-orchestration-pass-1.md` (approved direction: option B, Claude Code first, opt-in) · Status: draft for `savepoint-design` to turn into records

This file is written in the shape of Savepoint's own Goal, Objective, and Task records so the design step can create them without reinterpretation. IDs are deliberately absent: `savepoint-design` allocates the Goal and Objective IDs, and `savepoint create-task` allocates Task IDs. Nothing here is a record yet.

**How to turn this into records.** Close G-002 on the board (its Next is already `Close`). Then, in a fresh agent session, say: *"Plan from docs/reviews/agent-orchestration-spike-01-plan.md: create the Goal, select it, and create its three Objectives with their Tasks."* The design skill will ask you to confirm the Goal choice; the router's `release:` changes only on that explicit answer.

---

## Goal — Agent Orchestration Spike 01

```yaml
title: Agent Orchestration Spike 01
status: planned
```

### Outcome

With `agent_launcher.enabled: true`, an owner can press one key on a startable Task and Savepoint runs Claude Code headless in the project directory, shows the agent's real file, command, and permission events on the board as they happen, survives the board being closed, and records the outcome in files that say what Savepoint observed, separately from what the agent reported. With the flag off (the default), nothing on any screen or in `savepoint resume` changes. The paste-the-Next-line workflow in any terminal remains the primary, fully supported way to work.

### Why

Today the human is the router: copy the Next line, switch to an agent terminal, watch there, switch back to update the board, open another session for the Check. The review found that Savepoint already has the next-action projection, the file-watching board, a shell-free subprocess runner, and a parsed-but-unused `agent_launcher` config. The missing piece is a run: a detached process that starts the agent and writes its status and events under `.savepoint/`. Claude Code 2.1.296 is the only installed provider whose event format, permission-denial behaviour, cancellation, and usage reporting are documented, so the spike hard-codes it.

### Success Conditions

- Non-negotiables from the review hold in every Objective: local-first; run records are files from the first commit; agent output is labelled as agent-reported and never feeds a gate; Savepoint never writes a Task status on an agent's behalf; the owner approves permissions in advance and sees every denial; manual operation is untouched; no AI orchestrator; the prompt is the Next line and nothing more; existing Go and Bubble Tea code is extended, not rewritten.
- A real Claude Code run started from the board on WSL2 and on native Windows produces at least one edit event and one command event, finishes with exit code, session ID, usage, and denials recorded, and can be cancelled with no surviving child process.
- Closing and reopening the board during a run loses nothing.
- Per-run usage, as reported by the provider, is not higher than a hand-started session given the same Next line.
- With the flag off, board output, key map, help, and `savepoint resume` are byte-identical to v2.3.0 on the same project, proven by a golden test.

### Boundaries

**In scope:** a new `internal/run` package and `savepoint run` subcommand; a Claude Code `stream-json` adapter; run records under `.savepoint/runs/`; the board's opt-in run key, Runs overlay, activity view, and cancel; `doctor` reporting of lost runners; Design.md, README, config template, and AGENTS.md wording for the opt-in; validation evidence on two platforms.

**Out of scope (later spikes):** Check runs and BUILD → CHECK handoff; interactive takeover (`tea.ExecProcess`); Codex and OpenCode adapters; per-stage provider or model assignment; parallel runs and worktree lanes; any permission-editing UI; a daemon; PTY or transcript views; cost estimation; changes to gates, `ResolveNext`, or the workflow skills.

---

## Objective 1 — Run records and the detached runner

```yaml
title: Record and run one agent as a detached process
depends_on: []
priority: high
```

### Outcome

`savepoint run --task T-### [dir]` starts one Claude Code session for the named Task in the project root, owns the process, and writes `.savepoint/runs/<timestamp>-<task>/run.yml`, `events.jsonl`, and `summary.md`. The record states what the runner observed: state, exit code, session ID, denials, usage, heartbeat, and the agent's final message under `agent_reported`. Cancel by file, timeout, and process-tree kill work on Linux, macOS, WSL2, and Windows. Nothing in the record is read by any lifecycle gate.

### Why

"Files are state" and "agent output is evidence, not truth" are met only if the runner, not the board, owns the process and writes the record. Making the runner a subcommand rather than board code is what lets a run outlive the board (review §3, option B).

### Success Conditions

1. Running the subcommand with the feature flag off exits 2 with a one-line reason and writes nothing; with it on, `run.yml` exists before the agent starts and is updated at least every 5 s while it runs.
2. The Claude adapter maps `system/init`, `assistant` `tool_use`, `user` `tool_result`, `system/permission_denied`, and `result` lines to display events and record fields; unknown lines are kept verbatim on disk and shown as unparsed, never dropped.
3. A `cancel` file in the run directory ends the run within 15 s on every supported platform with state `cancelled` and no surviving agent or grandchild process; a timeout does the same with state `timed_out`.
4. A non-zero exit without a `result` line records state `failed` and a sanitised stderr tail.
5. The argv is built from `agent_launcher.builder` with `{{prompt}}` = the Next line verbatim and `{{project_root}}`; no shell is involved; secrets in URLs are stripped from stored text as the health runner already does.
6. `guard.js` still blocks agents from running `savepoint run`; `doctor` lists runs whose heartbeat is stale and whose PID is dead as `runner_lost`.
7. Configured gates pass, including `make test-full`, with Windows evidence named.

### Architectural Considerations

- Generalise or copy the process handling in `internal/codehealth/runner.go`, `runner_unix.go`, `runner_windows.go` into `internal/run`; the board must not import `codehealth` for this.
- `events.jsonl` is append-only (`O_APPEND`), one provider line per line; `run.yml` is rewritten atomically.
- Record shape per review §5: `provider`, `command`, `cwd`, `task`, `next_line`, `started_at`, `updated_at`, `pid`, `state`, `exit_code`, `session_id`, `denials[]`, `usage`, `agent_reported`.
- Default flags: `-p`, `--output-format stream-json`, `--verbose`, `--permission-mode acceptEdits`, `--allowedTools` derived from `quality_gates` as `Bash(<first word> *)`, `--permission-prompts none`. No `--bare`. Version floor 2.1.259.

### Boundaries

In scope: `internal/run`, `cmd/run.go`, `main.go` wiring, `doctor` line, tests with a fake `claude` that replays a recorded fixture. Out of scope: any board change, any Task write, any provider other than Claude Code.

### Tasks

| Task | Done when |
|---|---|
| Runner core with a fake provider | `internal/run.Start` writes `run.yml`, streams a fake provider's stdout to `events.jsonl`, heartbeats, and records exit state; unit tests cover running, finished, failed. |
| Claude Code stream-json adapter | Fixture from a real session parses to the expected events and result fields; unknown lines survive verbatim; denial and usage fields populate `run.yml`. |
| Cancel, timeout, and liveness across platforms | Cancel file and timeout end the process tree on Unix and Windows; stale heartbeat plus dead PID reads as `runner_lost`; tests run on both CI platforms. |
| `savepoint run` subcommand, flag gate, guard, and doctor | Subcommand wired and documented as owner-or-board only; exits 2 when the flag is off; `guard.js` unchanged and tested; `doctor` reports lost runners; `make test-full` evidence recorded. |

---

## Objective 2 — Start and watch a run from the board (opt-in)

```yaml
title: Start and watch a run from the board, opt-in
depends_on: [{objective: <Objective 1>, requires: clear}]
priority: high
```

### Outcome

When `agent_launcher.enabled` is true, `a` on a Task whose Next is Start, Build, or Test and whose gate allows it starts `savepoint run` detached; `A` opens a Runs overlay listing runs with state, elapsed time, and last event; `enter` opens a live activity view of real events; `c` cancels. A finished run shows the agent's message under an "agent reports" label and no verdict of its own. A run with denials reads "needs you" and names the refused actions. When the flag is false, none of this exists and the board is byte-identical to today.

### Why

This is the step that removes the terminal switch. It must be an addition beside the paste workflow, not a replacement, and the board must stay a reader of files so a run never depends on the board being open (review §3 and the opt-in rule).

### Success Conditions

1. Golden test: with the flag off or absent, rendered board, key map, help overlay, status hints, and plain (non-TTY) output match v2.3.0 for the same project.
2. `a` is offered only when the feature is on, the focused record is a Task, its Next verb is Start/Build/Test, and the same gate `actionsForRecord` uses allows it; a Blocked Task cannot start a run; a second run cannot start while one is `running`.
3. The board spawns the runner via `os.Executable()` detached (`Setsid` on Unix; `DETACHED_PROCESS|CREATE_NEW_PROCESS_GROUP` on Windows) and holds no child handle; closing the board does not end the run; reopening shows the run with its newer events.
4. The watcher includes `runs`; the activity view reads `events.jsonl` incrementally from the last offset, keeps a bounded in-memory window, coalesces bursts through the existing debounce, and stays responsive on a 50 MB event file.
5. The overview and activity view follow the Atari-Noir sketches in review §4: real events only, no percentages; orange running, purple agent-reported, red needs-you, dim queued.
6. `c` writes the cancel file and the row reaches `cancelled` without the board touching the process.
7. Configured gates pass, including `make test-full`.

### Architectural Considerations

- Follow the existing overlay pattern (`Issues`, `Health`, `Options`): a nullable overlay struct on `Model`, mutually exclusive with the others, origin restored on close.
- All reads are `tea.Cmd`s; the reducer only appends; no filesystem access in `Update`.
- Run state shown on the board is an inference from files plus PID liveness; the board never rewrites `run.yml`.

### Boundaries

In scope: `internal/board/v2` (new `runs.go`, `runs_view.go`, watcher root, keys, help, hints), config gate, golden tests. Out of scope: re-run with a wider allow-list, takeover, permission editing, Check runs.

### Tasks

| Task | Done when |
|---|---|
| Feature gate and golden parity | Board reads `agent_launcher.enabled`; with it off, golden fixtures for TUI, help, hints, and plain output match the pre-change render byte for byte. |
| Start a run from a Task | `a` appears only under the gated conditions; pressing it spawns the detached runner and shows a status line; a second start while running is refused with a reason. |
| Runs overlay and live activity view | `A` lists runs; `enter` shows incremental events; large-file and burst tests prove responsiveness; rendering matches the review sketches. |
| Cancel and needs-you presentation | `c` cancels by file; denied actions render as a red "needs you" block naming each refusal; finished runs label the agent's message as agent-reported. |

---

## Objective 3 — Prove it on real runs and document the opt-in

```yaml
title: Validate Spike 01 on real runs and document the opt-in
depends_on: [{objective: <Objective 2>, requires: clear}]
priority: medium
```

### Outcome

Spike 01's acceptance criteria are shown to hold with a real Claude Code session on WSL2 and on native Windows, including cancel, timeout, runner loss, denials, and closing the board mid-run, with the evidence recorded in the Objective. Usage per run is compared with a hand-started session on the same Next line. The README, this repository's and the scaffold's AGENTS.md, the config template comment, and `.savepoint/Design.md` describe the feature as opt-in and name the paste workflow as the default.

### Why

The review's final test is "if any of these fail, the direction is wrong, not the implementation". That judgement needs recorded evidence, not a passing unit suite. The docs change is small but is where the opt-in promise becomes visible to users.

### Success Conditions

1. A table in the Objective's evidence, one row per acceptance criterion 1–9 of review §5, with platform, command, time, and result; every row passes on both platforms or names an Issue.
2. A chaos run (`kill -9` of the runner mid-run) shows `runner_lost` on the board and in `doctor`.
3. Usage comparison recorded: provider-reported input and output tokens for one board-started and one hand-started run of the same Next line, with the difference explained.
4. README and both AGENTS.md copies state that orchestration is off by default and that pasting the Next line into any terminal remains fully supported; the config template's `agent_launcher` comment matches the shipped flags; `Design.md` Codebase Map gains `internal/run`.
5. Existing projects receive the guidance through `savepoint upgrade-assets` with user-edited files preserved.
6. Configured gates pass, including `make test-full`.

### Boundaries

In scope: validation runs, evidence tables, Issues for any failure, documentation and template updates, Design.md. Out of scope: fixing defects found (those become Issues routed through the normal repair path), new features.

### Tasks

| Task | Done when |
|---|---|
| Real-run validation on WSL2 and Windows | Evidence table complete for criteria 1–9 plus the chaos run and usage comparison; failures captured as Issues with repro. |
| Document the opt-in and update Design | README, AGENTS.md (both copies), config template, and Design.md updated; scaffold copies byte-identical to canonical; content tests assert the opt-in wording; `upgrade-assets` delivers them. |

---

## Sequencing and dependencies

Objective 1 → Objective 2 → Objective 3, each `requires: clear`. Objective 1 can start immediately after the Goal is selected. Objective 3's validation Task should be run by the owner or a fresh session, not the executor of Objective 2, so the evidence is independent.

## Spikes after this Goal (not planned here)

1. Check runs and the session-ID proof of independence. 2. Interactive takeover with `tea.ExecProcess`. 3. Codex adapter behind `agent_launcher`; OpenCode after its headless behaviour is verified. 4. Per-stage assignments and an allow-list editor. 5. Parallel runs with worktree lanes.
