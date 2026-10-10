---
id: G-003
title: Agent Orchestration Spike 01
status: planned
---

## Outcome

With `agent_launcher.enabled: true`, an owner can press one key on a startable Task and Savepoint runs Claude Code headless in the project directory, shows the agent's real file, command, and permission events on the board as they happen, survives the board being closed, and records the outcome in files that say what Savepoint observed, separately from what the agent reported. With the flag off (the default), nothing on any screen or in `savepoint resume` changes. Pasting the Next line into any terminal remains the primary, fully supported way to work.

## Why

Today the human is the router: copy the Next line, switch to an agent terminal, watch there, switch back to update the board, open another session for the Check. The review at `docs/reviews/agent-orchestration-pass-1.md` (2026-10-11) found that Savepoint already has the next-action projection, the file-watching board, a shell-free subprocess runner, and a parsed-but-unused `agent_launcher` config. The missing piece is a run: a detached process that starts the agent and writes its status and events under `.savepoint/`. Claude Code 2.1.296 is the only installed provider whose event format, permission-denial behaviour, cancellation, and usage reporting are documented, so the spike hard-codes it. The plan is at `docs/reviews/agent-orchestration-spike-01-plan.md`.

## Success Conditions

- The review's non-negotiables hold in every member Objective: local-first; run records are files from the first commit; agent output is labelled as agent-reported and never feeds a gate; Savepoint never writes a Task status on an agent's behalf; the owner approves permissions in advance and sees every denial; manual operation is untouched; no AI orchestrator; the prompt is the Next line and nothing more; existing Go and Bubble Tea code is extended, not rewritten.
- A real Claude Code run started from the board on WSL2 and on native Windows produces at least one edit event and one command event, finishes with exit code, session ID, usage, and denials recorded, and can be cancelled with no surviving child process.
- Closing and reopening the board during a run loses nothing.
- Per-run usage, as reported by the provider, is not higher than a hand-started session given the same Next line.
- With the flag off, board output, key map, help, and `savepoint resume` are byte-identical to v2.3.0 on the same project, proven by a golden test.

## Boundaries

**In scope:** a new `internal/run` package and `savepoint run` subcommand; a Claude Code `stream-json` adapter; run records under `.savepoint/runs/`; the board's opt-in run key, Runs overlay, activity view, and cancel; `doctor` reporting of lost runners; Design.md, README, config template, and AGENTS.md wording for the opt-in; validation evidence on two platforms.

**Out of scope (later spikes):** Check runs and BUILD → CHECK handoff; interactive takeover (`tea.ExecProcess`); Codex and OpenCode adapters; per-stage provider or model assignment; parallel runs and worktree lanes; any permission-editing UI; a daemon; PTY or transcript views; cost estimation; changes to gates, `ResolveNext`, or the workflow skills.
