# Savepoint — Agent Orchestration Review (Pass 1)

Date: 2026-10-11 · Repo state: branch `2.3.0-skills-optimisation` at `280e715` (v2.3.0) · Type: architecture review, no code changed

**Summary.** Savepoint already has every piece the orchestration loop needs except the one that starts an agent: a deterministic `Next` projection, a file-watching Bubble Tea board whose every write re-resolves its gate from disk, a shell-free subprocess runner with process-group kill on Unix and `taskkill` on Windows, and a parsed-but-unused `agent_launcher` config block. The recommended architecture is **B, a detached `savepoint run` process that owns one agent and writes its status and events to `.savepoint/runs/`, with the TUI as a pure watcher of those files.** It costs little more than option A because "files are state" forces the file writer anyway, and it gives reconnection after the TUI closes for free. All three CLIs (Claude Code 2.1.296, Codex 0.162.1, OpenCode 1.18.34) are installed here and all have a headless mode, model selection, and session resume; Claude Code is the only one whose event schema, permission-denial semantics, cancellation behaviour, and usage reporting are all documented, so Spike 01 hard-codes it. Spike 01 proves: press a key on a Task → Claude Code runs headless in the project directory with permissions the owner pre-approved → real tool events stream into the board → the process ends or is cancelled → Savepoint records what it observed (not what the agent claimed) → the owner decides the next step from the board without opening another terminal.

Evidence labels: **confirmed** (seen in code or in `--help` output run on this machine), **verified** (official docs URL cited), **assumed / unverified** (stated as such; nothing downstream relies on it without a flag).

**Opt-in rule (owner decision, applies to every section below).** Orchestration is off unless the owner turns it on, and the manual loop stays the primary path forever. Concretely: the feature is enabled only by `agent_launcher.enabled: true` in `config.yml` (the existing, disabled-by-default flag); with it off, the board shows no run key, no Runs panel, and no change to any screen, and `savepoint resume` output is byte-identical to today. With it on, the `Next` line is still printed and shown exactly as it is now, so the owner can always copy it into any terminal, any agent, any machine, and those hand-started sessions behave exactly as today. Savepoint never requires, nudges toward, or defaults to starting the agent itself; the `a` key is an addition beside the paste workflow, not a replacement for it. Nothing in the run records feeds a gate, so a project that never uses a run is indistinguishable from one that does.

---

## 1. Current-state assessment

### How Savepoint works today

| Area | Statement | Status | Where |
|---|---|---|---|
| Entry | `main.go` dispatches `init`, `board`, `doctor`, `upgrade-assets`, `migrate`, `resume`, `health`, `create-task`. `resume` loads the V2 index and router, calls `data.ResolveNext`, renders via `internal/resume`. | confirmed | `main.go:106-122,218-270`, `cmd/resume.go` |
| State files | All workflow state is Markdown with YAML frontmatter under `.savepoint/`: `router.md` (state, release, objective, task, issue), `objectives/O-###-slug/Objective.md`, `objectives/.../tasks/T-###-slug.md`, `checks/C-###-slug.md`, `issues/`, `releases/`, `config.yml`. | confirmed | `.savepoint/router.md`, `internal/data/task_v2.go:55-68`, sample `checks/C-977-o042-objective-check.md` |
| Task lifecycle | `status`: `planned` → `in_progress` → `done`; `stage` (required while in progress): `build` → `test` → `audit`. Advance/retreat are pure functions. | confirmed | `internal/data/lifecycle.go:95-142`, `AGENTS.md` Terminology |
| Router state | `state` ∈ {idea, design, task, check}; selection keys `release`, `objective`, `task`, `issue`. | confirmed | `internal/data/router_v2.go:13-58` |
| Next projection | `ResolveNext` returns one of 11 `NextKind` rungs (replan, dependency, execute, check_needed, owner_validation_required, assess_decision, objective_integration, objective_ready, release_ready, select_task, nothing_selected). `internal/resume` turns the rung into a verb line: Start/Build/Test/Check/Assess/Plan/Replan/Fix/Choose/Accept/Close/Blocked. The verb is what AGENTS.md tells an agent to act on. | confirmed | `internal/data/next.go:242-282,353`, `internal/resume/resume.go:50-136` |
| Gates | Task start/advance/completion and Objective/Goal closure are decided by resolvers in `internal/data` (`gate_v2.go`, `objective_gate_v2.go`, `release_gate_v2.go`, `dependency.go`). Evidence types: `Freshness`, `OwnerValidation`, `Exception`, `Replan`, `CheckWaiver`. | confirmed | `internal/data/evidence_v2.go:12-90` |
| Board | Bubble Tea v1.3.10 program (`tea.WithAltScreen`). `Init` runs a load command; an fsnotify watcher on `objectives`, `checks`, `issues`, `releases`, `router.md`, `config.yml` with a 100 ms debounce sends `v2FileChangeMsg`, which reloads. Non-TTY falls back to deterministic plain text. | confirmed | `internal/board/v2/run.go`, `watch.go:14-20,74-77`, `plain.go` |
| Board writes | Every owner key (space advance, backspace retreat, acceptance, exception closure, Goal toggle) is a `tea.Cmd` that re-reads the index and re-resolves the gate immediately before writing, so a stale screen can never write a transition its gate would refuse. Writes are field patches to frontmatter with mtime conflict detection. | confirmed | `internal/board/v2/io.go:46-50,276-285`, `internal/data/write.go:14,70-80` |
| Board keys in use | `↑↓←→ hjkl`, `enter`, `v`, `d`, `i`, `I`, `H`, `o`, `g`, `?`, `q`, `space`, `backspace`, `1-4`, `J/K`, `x`. **`a` is free.** | confirmed | `internal/board/v2/update.go` case list, `help.go` |
| Subprocess runner | `codehealth.ExecRunner` runs an argv (no shell) with `CI=1 NO_COLOR=1`, bounded stdout/stderr buffers, `WaitDelay`, and on cancel kills the whole process group (Unix `Setpgid` + `SIGKILL -pgid`) or the tree (`taskkill /T /F`). | confirmed | `internal/codehealth/runner.go:52-99`, `runner_unix.go`, `runner_windows.go` |
| Gate runner | `doctor` runs `quality_gates` commands with a 60 s default timeout via `splitCommand`. | confirmed | `internal/doctor/gates.go:22-60` |
| Launcher config | `config.yml` already documents an opt-in, disabled-by-default `agent_launcher` block: `builder`/`auditor` as `{command, args[]}` with `{{prompt}}` and `{{project_root}}` placeholders, plus a `terminal` adapter. It is parsed and validated and **used nowhere else in the binary**. | confirmed | `.savepoint/config.yml`, `internal/data/launcher_config.go`, grep of `AgentLauncher` (only `config.go`) |
| Claude Code integration | `savepoint init` installs `.claude/settings.json` with a `SessionStart` hook that injects the `Next` line and a `PreToolUse` guard that blocks agent edits setting `status: done`, router edits in a worktree lane, and owner-only `savepoint` commands. | confirmed | `templates/project-v2/.claude/settings.json`, `.claude/hooks/guard.js`, `session-start.js` |
| Worktree lanes | Documented policy: in a lane, agents do not touch `router.md` and do not create records. Enforced by the guard hook, not by Go code. | confirmed (policy), assumed (no Go enforcement found by grep) | `AGENTS.md` Worktree Lanes, `guard.js` |
| Size | 39,321 lines of non-test Go, 63,546 lines of tests. Dependencies: `bubbletea`, `fsnotify`, `lipgloss`, `yaml.v3`. No `bubbles` components are used directly. | confirmed | `go.mod`, `wc` |
| Today's loop | The human copies the `Next` line into an agent terminal, watches it there, returns to the board to advance/accept, opens a *fresh* session for Check. Nothing in the binary starts a process other than gates, health tools, `npm install`, and clipboard helpers. | confirmed | grep of `exec.Command` |

**Observation that matters for design:** the board is already an event-sourced reader of files. Adding a run is "add a watched directory and a card type", not "add a process manager to the TUI".

---

## 2. Provider capability matrix

All three binaries were run on this machine (WSL2, Linux 6.18). Sources: H = `--help` output run here; D = official docs fetched 2026-10-11.

| Capability | Claude Code 2.1.296 | Codex CLI 0.162.1 | OpenCode 1.18.34 |
|---|---|---|---|
| Headless mode | **verified** `-p/--print` (H; D [1]) | **verified** `codex exec [PROMPT]`, stdin prompt (H) | **verified** `opencode run [message]` (H; D [4]) |
| Structured event stream | **verified** `--output-format stream-json`: `system/init`, `assistant`/`user` messages carrying `tool_use`/`tool_result` blocks, `system/permission_denied`, `system/api_retry`, final `result` (D [1]) | **verified (partial)** `--json` JSONL: `thread.started`, `turn.started`, `item.*` (`agent_message`, `command_execution` named), `turn.completed`, `turn.failed`, `error`. Exact file-change item name **unverified** (D [2]) | **unverified**: `--format json` "raw JSON events" exists (H; D [4]) but no schema is documented |
| Model selection | **verified** `--model <alias\|name>`, `--effort` (H) | **verified** `-m/--model` (H) | **verified** `-m provider/model`, `--variant` (H) |
| Session resume (finished session) | **verified** `--resume <id>`, `--continue`, `--session-id <uuid>`, `--fork-session`; session id is in the `result` message (H; D [1]) | **verified** `codex exec resume <id\|--last> [prompt]`, `codex exec fork` (H; D [2]) | **verified** `-c/--continue`, `-s/--session <id>`, `--fork` (H; D [4]) |
| Continue a *live* session | **verified** `--bg` background sessions with `claude attach/logs/stop`; `-p` rejects `--bg` (H; D [1]). Multi-turn input into one running `-p` process via `--input-format stream-json` (H) | **unverified**: `codex queue` "Queue a message for an existing session" and `codex agents` daemon exist (H) but semantics untested | **verified (mechanism)**: `opencode serve` + `run --attach <url>` / `opencode attach` (H; D [4]) |
| Cancellation | **verified** SIGINT ends the turn; SIGTERM exits 143, kills Bash process trees, runs `SessionEnd` hooks, leaves the turn unfinished and resumable (D [1]) | **unverified** (not covered in docs) | **unverified** |
| Permissions when headless | **verified** `--permission-mode {acceptEdits,auto,bypassPermissions,manual,dontAsk,plan}`, `--allowedTools`, `--disallowedTools`, `--permission-prompts none` (≥2.1.259), `--permission-prompt-tool`. A denied action is reported to the model, the run continues, `permission_denied` events and `result.permission_denials` list them (H; D [1]) | **verified** `--sandbox read-only\|workspace-write\|danger-full-access`, `--approve-for-me`, `--dangerously-bypass-approvals-and-sandbox`; `exec` defaults to read-only (D [2], [3]). `exec` has **no** `-a` flag (H, grep). Top-level `-a never`: "Execution failures are immediately returned to the model" (H). What the model sees when the sandbox blocks a command in `exec`: **unverified** | **verified** config `permission` with `allow/ask/deny` and glob patterns (D [5]); `--auto` approves anything not explicitly denied (H). Behaviour of an `ask` rule inside `run` without `--auto`: **unverified** |
| Usage / token reporting | **verified** `json`/`stream-json` result carries `usage` and `total_cost_usd` (client-side estimate); subscription limits **not exposed** (D [1]) | **verified** `turn.completed.usage`: `input_tokens`, `cached_input_tokens`, `output_tokens`, `reasoning_output_tokens`. No cost field (D [2]) | **verified (aggregate only)** `opencode stats` shows token usage and cost per project/days (H; D [4]). Per-run figures in JSON events **unverified** |
| Uses existing CLI login | **verified** yes by default; `--bare` switches to API-key only (D [1]) | **verified** "reuses saved CLI authentication by default"; `CODEX_API_KEY` optional (D [2]) | **verified (mechanism)** credentials via `opencode providers`/`auth` (H); that `run` uses them is assumed |

Sources: [1] https://code.claude.com/docs/en/headless · [2] https://learn.chatgpt.com/docs/non-interactive-mode (redirect of developers.openai.com/codex/noninteractive) · [3] https://learn.chatgpt.com/docs/agent-approvals-security · [4] https://opencode.ai/docs/cli/ · [5] https://opencode.ai/docs/permissions/

### The permission question

A headless agent cannot stop and ask. Each provider resolves this differently, and it decides what "blocked" means in the UI:

- **Claude Code** turns a would-be prompt into a *denial the model is told about* and keeps going (`--permission-prompts none`), or stops the process waiting for a host that never answers (without the flag and with a permission tool). For Savepoint the first is the only sane mode: the run finishes, the denials are in the `result`, and the board shows "finished with 2 denied actions", not a hung process. Pre-approval is a tool allow-list (`Bash(make *)`, `Edit`), a mode (`acceptEdits`), or the `auto` classifier. `bypassPermissions` exists but violates non-negotiable 5 as a default.
- **Codex** does not ask at all in `exec`; the sandbox is the permission. `workspace-write` is the realistic mode (read-only cannot build). Network is off in `workspace-write` unless configured. A blocked command is, per the `-a never` text, returned to the model as a failure; whether `exec --json` surfaces that as a distinguishable item is **unverified**, so a Codex "blocked" state would initially have to be inferred from `turn.failed`/`error` or the final message. On Linux/WSL2 the sandbox is `bwrap`+`seccomp`; WSL1 is unsupported; Windows uses MXC.
- **OpenCode** has a static allow/ask/deny policy. Without `--auto`, what `run` does on `ask` is undocumented; this must be tested before OpenCode is a supported provider. `--auto` is "dangerous!" by its own help text.

**Consequence for intervention.** A blocked run is never an interactive prompt inside Savepoint. It is a *finished or stalled run whose record names what was refused*, plus two owner choices: re-run with a wider allow-list, or take over interactively (resume the same session in the terminal). That keeps the TUI free of a permission UI, which it has no provider-neutral way to build.

**Session continuity, per provider:** Claude Code: live attach (`--bg`), finished-session resume, and fork are all real. Codex: finished-session resume and fork are real; live queueing is unverified. OpenCode: finished-session resume real; live continuation is real only through its HTTP server. "Restart with reconstructed context" is always available and is what Savepoint's files already enable; it should be the fallback, not the design.

---

## 3. Architecture options

| Criterion | A — TUI child processes | B — Detached `savepoint run` writing files | C — tmux window per agent |
|---|---|---|---|
| Complexity | Low: `exec.Command` + a `tea.Cmd` that reads lines | Medium: A plus detach, PID/heartbeat, cancel file, self re-exec | Low code, high integration: tmux CLI, pane naming, pipe-pane capture |
| Fit with codebase | Good: runner exists; but the TUI becomes a process owner, which it never is today | Best: TUI stays a file reader; `run` is a sibling of `health check` (a subprocess-running subcommand) | Poor: no existing tmux layer; board semantics (gates, cards) unchanged but output is not machine-readable |
| Live progress | Direct pipes, lowest latency | fsnotify on `events.jsonl`; existing 100 ms debounce is fine | Only by attaching; no structured events |
| Reliability | Agent dies with the TUI (or is orphaned on crash) | Run survives TUI exit; a dead runner is detectable by PID + heartbeat | tmux server survives everything; but the run record must still be written by someone |
| Reconnection after TUI closes | None without extra work | Free: reload files | Free for the human; nothing for the board |
| Windows / WSL2 / Linux / macOS | Windows: no process groups, no SIGINT to a child; needs `taskkill`. Paths fine. CLI discovery via `exec.LookPath` plus `node_modules/.bin` (as `savepoint-find.js` does) | Same plus detaching: Unix `Setsid`; Windows `CREATE_NEW_PROCESS_GROUP\|DETACHED_PROCESS`. Cancel via a file the runner watches, so no cross-process signals are needed on Windows | No native Windows tmux; WSL2 only. Breaks the npm-distribution promise on Windows |
| Maintenance | Low | Medium: one more subcommand, one run schema | Medium: tmux versions, shells, terminal quirks |
| Risks | Non-negotiable 2 is met only if the TUI also writes files; TUI closed mid-run is a hard failure case | Two processes to keep honest; stale-runner detection must be right | Scope rule 10; Windows; output is a transcript, which non-negotiable 8 says not to hand off |

**Recommendation: B.** The marginal cost over A is detaching and liveness, and in exchange "TUI closed mid-run" becomes a non-event. A is what you would build if files were optional; they are not. C is rejected as the primary mechanism because of Windows and because tmux output is a transcript, not evidence. tmux remains a *user's* choice for their own interactive sessions; Savepoint never needs to know.

```
┌───────────────────────────── savepoint board (Bubble Tea) ─────────────────────────────┐
│ reads: objectives/ checks/ issues/ releases/ router.md config.yml  + runs/**            │
│ writes: Task/Objective/Goal transitions (existing io.go), runs/<id>/cancel (new)        │
│ starts: `savepoint run --task T-### --provider claude` detached, then forgets the PID    │
└────────────────────────────┬────────────────────────────────────────────────────────────┘
                             │ files only (fsnotify)                     ▲ files only
                             ▼                                            │
┌──────────────── savepoint run (one process per run, no daemon) ────────┴────────────────┐
│ builds argv from config (agent_launcher) + Next line · exec in project root              │
│ owns the agent child (codehealth-style process group) · enforces timeout · honours cancel│
│ writes runs/<id>/run.yml (status, pid, heartbeat, exit, session_id, denials, usage)      │
│        runs/<id>/events.jsonl (provider lines, append-only, never rewritten)             │
│        runs/<id>/summary.md (agent's final message, labelled "agent-reported")           │
└────────────────────────────┬────────────────────────────────────────────────────────────┘
                             │ stdout pipe (stream-json)
                             ▼
                     claude -p … / codex exec --json … / opencode run --format json …
                             │ edits Task files, Context Log, code (as the skill says today)
                             ▼
                 .savepoint/objectives/**  (watched; board reloads; resume re-resolves)
```

**How events reach the TUI.** The runner appends each provider line to `events.jsonl`. The board's watcher gains `runs` as a watched root; on change it reads from the last byte offset it saw (never re-parses the file) and maps lines to a small `RunEvent{Time, Kind, Text}` for display. The full line is only ever on disk.

**How workflow state stays authoritative.** Nothing in `runs/` is a lifecycle input. `ResolveNext`, the gates, and the board writes read exactly what they read today. A run record says "an agent ran, exited N, reported X, was denied Y"; the Task says what the agent wrote into it; the Check says what an independent session verified. Three sources, never merged. The run card renders agent-reported text under a visibly different label than gate results.

**What would trigger a later change.** (a) A provider-level live-session API worth using (Claude's `--input-format stream-json` or OpenCode's server) would justify a long-lived runner per session, still writing files. (b) If users want to type into a running agent, add `tea.ExecProcess` takeover, not a PTY. (c) A daemon is justified only if runs must outlive logout, which nobody has asked for.

**Live progress without a terminal emulator.** Yes, and better: the structured streams name the file edited, the command run, and the denial, which a PTY scrollback does not. For the "I want to drive it" case, `tea.ExecProcess` (in bubbletea v1.3.10) suspends the board, hands the terminal to `claude --resume <id>`, and returns when it exits, on every platform. A tmux window gives the same on Unix with a second tool and no return signal. Use `ExecProcess`; document tmux as "also fine".

---

## 4. Live progress in the TUI

Palette per `.savepoint/visual-identity.md`: uniform black, warm off-white text, orange = active, green = live/ok, purple = AI, red = blocked. No progress bars, no percentages; milestones and events only. Sketch 1 is a new overlay reached with `A` (runs), alongside `i` (issues) and `H` (health).

**4.1 Multi-agent overview**

```
▣  S A V E P O I N T                      ♥ Health: ok   ◆ 2/3 Objectives  ▪ 9/12 Tasks  ✗ 0/4 Issues
──────────────────────────────────────────────────────────────────────────────────────────────
RUNS                                                                        2 active · 1 waiting on you
──────────────────────────────────────────────────────────────────────────────────────────────
 ● T-127  Build  claude/opus     running   12m 04s   last: edit internal/board/v2/runs_view.go
 ● T-128  Test   codex           running    3m 41s   last: $ make test-fast
 ◐ O-045  Check  claude/opus     finished  08m 10s   agent reports: NEEDS WORK  · result unverified
 ✗ T-126  Build  claude/opus     blocked   02m 55s   2 actions denied · needs you
 ○ T-129  Start  —               queued    —         waits for T-127 (requires: clear)
──────────────────────────────────────────────────────────────────────────────────────────────
                         IDEA │ DESIGN │ TASK │ CHECK
     ↑↓:run  enter:activity  c:cancel  r:re-run with wider allow  t:take over in terminal  esc:close
```

● orange = running, ◐ purple = finished with agent-reported outcome not yet checked, ✗ red = needs you, ○ dim = queued. "Decisions waiting on the user" is a count in the header line and a red row, nothing else.

**4.2 Live activity for one run**

```
RUN 2026-10-11T09:14Z · T-127 · Build · claude 2.1.296 · opus · acceptEdits + Bash(make *)      12m 04s
──────────────────────────────────────────────────────────────────────────────────────────────
09:14:02  ▸ started      cwd /home/user/code/savepoint   session 7f3a…   prompt: "Build T-127 — …"
09:14:09  ◦ read         agent-skills/savepoint-task/SKILL.md
09:14:10  ◦ read         .savepoint/objectives/O-045-…/tasks/T-127-….md
09:15:31  ✎ edit         internal/board/v2/runs_view.go            (+41 −0)
09:16:02  $ command      make build                                 exit 0  4.1s
09:16:50  ✎ edit         internal/board/v2/runs_view_test.go        (+88 −0)
09:17:12  $ command      make test-fast                             exit 1  21.3s
09:17:12  ✗ failure      FAIL TestRunsOverlay_Empty (runs_view_test.go:31)
09:18:40  ✎ edit         internal/board/v2/runs_view.go            (+3 −1)
09:19:05  $ command      make test-fast                             exit 0  20.9s
09:26:01  ▸ turn         usage so far: in 48,210 · out 6,114 · est. $0.41 (provider estimate)
──────────────────────────────────────────────────────────────────────────────────────────────
                                 ↑↓:scroll  end:follow  c:cancel  esc:back
```

Only event kinds that the provider actually emits are shown: read, edit, command (with exit and duration), failure (non-zero command or `turn.failed`), denied, turn/usage, finished. Text and thinking deltas are not rendered and not stored (skip `--include-partial-messages`).

**4.3 A blocked run**

```
RUN 2026-10-11T08:40Z · T-126 · Build · claude · acceptEdits                     blocked · 02m 55s
──────────────────────────────────────────────────────────────────────────────────────────────
08:41:17  ✎ edit         templates/project-v2/AGENTS.md
08:42:03  ✗ denied       Bash: git push origin 2.3.0-skills-optimisation     not in allow-list
08:42:40  ✗ denied       WebFetch: https://example.com/spec                  not in allow-list
08:43:01  ▸ finished     exit 0 · agent reports: "Could not push; changes are committed locally."
──────────────────────────────────────────────────────────────────────────────────────────────
 NEEDS YOU   The agent finished but 2 actions were refused. Nothing was approved on its behalf.
   r  re-run T-126 with these allowed:  [ ] Bash(git push *)   [ ] WebFetch
   t  take over: open this session in the terminal (claude --resume 9c1e…)
   d  dismiss: keep the record, decide on the board
```

**Keeping the TUI responsive.** (1) The runner, not the board, parses provider output; the board reads one file incrementally from a remembered offset. (2) Keep a bounded ring (e.g. 500 display events) per run in memory; older lines stay on disk. (3) Reuse the existing debounce, and coalesce "N new lines" into one message. (4) Cap each rendered line at terminal width and each stored `Text` at a few hundred bytes; the raw line is in `events.jsonl`. (5) Never render while reading: the read is a `tea.Cmd`, the reducer only appends. (6) Partial-message deltas are not requested. Claude's own docs note it waits up to 30 s for a slow consumer; the runner must drain stdout continuously, which `codehealth`'s bounded buffers already do in spirit.

---

## 5. Spike 01 — one Task, one provider, headless, recorded

**Goal.** Prove that the board can start a Claude Code run for the focused Task, show real events as they happen, survive the board being closed, record the outcome in files without trusting the agent's own report, and leave the owner with the decision.

**Scope.**
- Gated by `agent_launcher.enabled: true`. When it is false (the default, and what `init` scaffolds), none of the below exists on screen or in `resume`; the paste-the-Next-line workflow is untouched and remains the documented way to work.
- New subcommand `savepoint run --task T-### [dir]` (hidden from the human-facing CLI rules except as "owner or board only"; the `guard.js` list stays unchanged so agents cannot start agents).
- Provider hard-coded: Claude Code. Command: `claude -p "<Next line>" --output-format stream-json --verbose --permission-mode acceptEdits --allowedTools <gate commands from quality_gates, as Bash(prefix *)> --permission-prompts none --model <from config or omitted>`. Prompt is exactly the `Next` line `resume` prints; the skill files carry the rest, as they do today. No `--bare` (it would drop CLAUDE.md, the Savepoint hooks, and subscription login).
- Run record under `.savepoint/runs/<yyyymmddThhmmssZ>-<task>/`: `run.yml`, `events.jsonl`, `summary.md`. `run.yml` fields: `provider`, `command` (argv), `cwd`, `task`, `next_line`, `started_at`, `updated_at` (heartbeat every ~5 s), `pid`, `state` ∈ {running, finished, failed, cancelled, timed_out, runner_lost}, `exit_code`, `session_id`, `denials[]`, `usage` (verbatim provider fields, labelled estimate), `agent_reported` (last message, verbatim).
- Board: `a` on a focused Task starts a run when `Next` for that Task is Start/Build/Test (reuse `actionsForRecord`'s gate so a Blocked Task cannot be started); `A` opens the runs overlay; `enter` opens the activity view; `c` cancels.
- Watcher: add `runs` to `v2WatchDirectories`.
- `savepoint resume` output unchanged. `doctor` gains one line per `runner_lost` run.

**Explicit exclusions.** Codex and OpenCode; Check runs and BUILD→CHECK handoff; per-stage assignments; parallel runs (the board refuses a second run while one is `running`); worktrees; interactive takeover; permission editing UI (re-run with a wider allow-list is a config edit in this spike); any change to gates, `ResolveNext`, or the skills; any Task status write by Savepoint on the agent's behalf.

**Dependencies.** Claude Code ≥ 2.1.259 on PATH (`--permission-prompts none`); Go 1.26; existing `codehealth` runner generalised or copied into an `internal/run` package (do not import `codehealth` from the board). `agent_launcher` config is read but only `builder.command`/`args` are honoured, with `{{prompt}}`/`{{project_root}}`.

**Technical approach.**
1. `internal/run`: `Spec{Task, Provider, Argv, Dir, Timeout}`, `Start(spec) (id, error)` that writes `run.yml` then execs the agent with the `codehealth` process-group setup, stdout piped line-by-line to `events.jsonl` with `O_APPEND`, stderr tail kept in `run.yml`.
2. Claude adapter: parse `type` = `system`/`assistant`/`user`/`result`; from `assistant` content blocks extract `tool_use` (name, file path or command); from `user` blocks, `tool_result` exit codes where present; `system/permission_denied` → denials; `result` → `session_id`, `usage`, `total_cost_usd`, `permission_denials`, `result` text. Unknown lines are kept on disk and shown as `·` (unparsed).
3. Detach: board runs `os.Executable()` with `run --task …`, `Setsid` on Unix, `DETACHED_PROCESS|CREATE_NEW_PROCESS_GROUP` on Windows, stdio to the run directory. The board never holds the child handle.
4. Cancel: board writes `runs/<id>/cancel`; the runner's loop sees it (poll 1 s or fsnotify), sends SIGINT to the agent, waits 10 s, then uses the existing group kill; state `cancelled`. Windows: skip SIGINT, go straight to the tree kill, state `cancelled`.
5. Liveness: the board marks `running` runs whose heartbeat is older than 30 s *and* whose PID is not alive as `runner_lost`; it never rewrites the runner's file, it renders the inference and `doctor` names it.
6. Timeout: default 30 min, `agent_launcher.timeout` optional; state `timed_out`, agent killed via the same path.

**Acceptance criteria.**
1. Pressing `a` on a startable Task creates `runs/<id>/run.yml` within 1 s and the overview shows the run as `running` without any board restart.
2. Within the run, at least one `edit` and one `command` event from a real Claude Code session render in the activity view before the process exits.
3. Closing the board (`q`) while the run is `running` does not end the run; reopening the board shows it still running with new events.
4. On agent exit, `run.yml.state` is `finished`, `exit_code`, `session_id`, `usage`, and `agent_reported` are filled from the `result` line, and the Task file changes the agent made are visible on the board through the existing watcher.
5. The board shows the agent's final message under an "agent reports" label and shows no verdict of its own; `savepoint resume` still derives Next from Task/Check records only.
6. `c` ends the run within 15 s on Linux/macOS/WSL2 and Windows, with state `cancelled` and no surviving `claude` or child `make` process.
7. A run started with an empty allow-list shows `denied` events and `permission_denials` in `run.yml`; the run still reaches `finished`.
8. `make build && make test-fast` pass; new code has unit tests for the Claude line parser (fixtures from a real session), cancel-file handling, and liveness.
9. With `agent_launcher.enabled` false or absent, the board's rendered output, key map, help screen, and `savepoint resume` output are byte-identical to the current release on the same project; a golden test asserts this.

**Validation method.** One scripted session on each of Linux (WSL2) and Windows native, using this repository's own T-### Task; a fake `claude` shell script that replays a recorded `stream-json` fixture for deterministic tests; a chaos check that `kill -9`s the runner mid-run and confirms `runner_lost` is rendered and reported by `doctor`.

**Failure cases.**

| Case | Expected behaviour |
|---|---|
| Agent crashes (non-zero exit, no `result`) | state `failed`, `exit_code` set, stderr tail in `run.yml`, board row red |
| Owner cancels | SIGINT then group kill; state `cancelled`; Claude's own session remains resumable (documented SIGTERM behaviour) |
| Timeout | as cancel, state `timed_out` |
| TUI closed mid-run | run continues; reopened board reads files; no duplicate run can be started while `running` |
| Permission denied | run continues (Claude semantics); denials recorded and rendered; owner chooses re-run or takeover |
| Runner killed / machine sleeps | heartbeat stale + PID dead → `runner_lost`; `doctor` lists it; owner may delete or keep the record |
| Provider missing / wrong version | `run` exits 2 before writing `run.yml`; board status line names the missing binary or the version floor |

**Key risks.** (1) The stream-json shape is documented but evolves; keep raw lines on disk and tolerate unknown types. (2) `acceptEdits` plus a gate allow-list may still deny something routine (e.g. `git status`); the first real runs will tell you what the default allow-list should contain, which is exactly why denials must be visible. (3) Windows detaching and tree-killing are the most likely platform bugs; test them first, not last. (4) Token cost: the agent reads the same skills as today; the only addition is the `Next` line, so the spike should not raise per-run context. Measure with the `usage` field before and after.

**Spikes to follow (one line each, recommended order).**
1. Check runs: `a` on a `Check O-###` Next line starts a run with the same provider; the run records `session_id` so the Check's `executed_session` independence claim is machine-checkable.
2. Interactive takeover: `t` suspends the board with `tea.ExecProcess` and resumes the run's session in the terminal; on return the board reloads.
3. Second provider (Codex) behind the `agent_launcher.builder/auditor` config, with a Codex event adapter; OpenCode only after its `run` permission and JSON behaviour are verified by a real run.
4. Per-stage assignments (Build vs Check provider/model) and a default allow-list editor in Advanced Options.
5. Parallel runs with worktree lanes: one run per lane, the board refusing overlaps on the same Task.

**On the working hypothesis.** "Persistent sessions and reconnection" should not be a late spike: reconnection falls out of B in Spike 01, and session persistence is a provider feature already available. Move interactive takeover up, because it is the honest answer to a blocked run. Put BUILD→CHECK before configurable assignments: a second provider adds a parser and config; a Check run adds the product's actual value and tests the independence rule under automation. Parallel runs stay last.

---

## 6. Decisions and direct answers

| # | Decision | Recommendation |
|---|---|---|
| 0 | Opt-in | Off by default behind `agent_launcher.enabled`; the Next line stays copyable and hand-started sessions in any terminal remain the primary, fully supported path. Decided by the owner on 2026-10-11 |
| 1 | First provider | Claude Code: the only one with documented event schema, denial semantics, signal behaviour, and usage; Savepoint already ships Claude hooks |
| 2 | Process model | B, detached `savepoint run`; TUI reads files only |
| 3 | Default permission posture | `acceptEdits` + allow-list derived from `quality_gates` + `--permission-prompts none`; `auto` and `bypassPermissions` only by explicit per-project config; never silently |
| 4 | Where run records live | `.savepoint/runs/<timestamp>-<task>/`; no `R-###` identity (collides with Releases) and no task-ID allocation |
| 5 | What the prompt contains | The `Next` line verbatim, nothing else; the skills carry the method; no transcript or Task body is injected |
| 6 | Who changes Task status | Unchanged: the agent through its skill, the owner through the board; Savepoint's runner never writes a Task |
| 7 | Reuse or retire `agent_launcher` | Reuse the struct and config block; honour `builder` in Spike 01, `auditor` in the Check spike; drop the `terminal` adapter in favour of `tea.ExecProcess` later |
| 8 | Launch key and gating | `a` on a Task whose Next is Start/Build/Test and whose gate allows it; one running run per project in Spike 01 |
| 9 | Usage display | Show provider fields verbatim with "estimate" and provider name; show nothing for providers that report nothing; never show subscription limits |
| 10 | Agent-run guard | Keep `run` out of the agent-allowed command list in `guard.js` so agents cannot start agents |

**Direct answers.**

1. **Entirely within the TUI?** Yes for the monitored loop: start, watch, cancel, decide. Not for typing to a running agent, which needs a terminal handover (`tea.ExecProcess`) that returns to the board; that is still one terminal.
2. **Smallest implementation that removes meaningful switching?** Spike 01 as specified: one key starts the Build/Test run for the focused Task, events appear in the board, the outcome is recorded, and the owner's advance/accept keys are already there. The Check spike removes the second switch.
3. **Reusable today?** `ResolveNext` and the resume verb line (the prompt), the gate resolvers (whether a run may start), the fsnotify watcher and debounce, the overlay pattern (`Issues`, `Health`, `Options`), `codehealth.ExecRunner` process handling, the `agent_launcher` config and validation, `doctor` for stale-run reporting, and the Claude hooks already installed by `init`.
4. **Postpone.** Codex and OpenCode adapters, per-stage assignments, parallel runs and worktree management, any permission UI, a daemon, live PTY views, cost estimation, and anything that reads an agent's transcript.
5. **What makes it expensive or brittle?** Parsing three evolving event schemas (mitigate: raw lines on disk, display only what parses); Windows process semantics (mitigate: cancel-by-file, runner is the parent); trying to make permissions interactive across providers (do not); letting run records feed gates (never); injecting Task bodies or transcripts into prompts (costs tokens and duplicates the files the agent reads anyway).
6. **What Spike 01 must prove before committing further.** That a real Claude Code run started from the board produces useful, accurate events; that the board survives closing and reconnects from files; that cancel and timeout kill the whole tree on Windows and WSL2; that denials are visible and the run still completes; and that per-run usage is not higher than a hand-started session with the same `Next` line. If any of these fail, the direction is wrong, not the implementation.
