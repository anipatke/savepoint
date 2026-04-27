---
type: project-prd
status: active
---

# Savepoint — Product Vision

## What it is

A public OSS CLI + Ink-based TUI that scaffolds an opinionated AI-driven development workflow. The user runs `npx savepoint init` in an empty directory, points any AI agent (Claude / Cursor / Cline / Gemini / Aider / Codex) at the project, and the embedded prompt templates carry the agent through:

> **PRD → Design → Epics → Tasks → Build → Audit**

…with hard gates at each transition.

## Why

Three failure modes plague AI-driven development today:

1. **Inconsistency.** No repeatable process from high-level intent to working MVP.
2. **Token bloat.** Monolithic backlogs and MCP overhead burn context for everyone, but especially users on minimal AI plans.
3. **Documentation drift.** `Design.md` and agent instructions go stale after the first iteration; nobody updates them.

Savepoint addresses all three with a single mechanism: **a file-based state machine that any agent can follow, where every epic completion forces a documentation audit before the next epic can start.**

## Target user

**Vibe coders** — builders with minimal-to-moderate development experience, on minimal AI plans, who want AI agents to drive most of the implementation while a structured workflow keeps the project coherent.

Not: experienced engineers who already have their own systems. (They can still use it; they're not the audience.)

## Headline differentiator

**The Audit Loop.** When the last task in an epic moves to `done`, the next epic cannot start until `Design.md`, `AGENTS.md`, and the epic's own design have been reconciled with the actual code that was built. No existing markdown-first task tool has this gate.

Token-efficient hierarchy and markdown-first storage are table stakes. The audit loop is the marketing-first feature.

## Success metrics

- **Token usage:** AI agents complete tasks reading <2KB of context per task. Audit budget bounded to ~5–15KB.
- **Documentation accuracy:** zero drift — `AGENTS.md` always correctly maps the current codebase, enforced by gate.
- **Agent reach:** works with any agent that can read markdown and edit files (no MCP required, no per-agent adapters).
- **Time-to-first-PR:** a vibe coder can go from `npx savepoint init` to a merged epic in one weekend.

## Constraints

- File-only architecture for v1. No MCP server.
- Agent-agnostic via the [Router Pattern](Design.md). No per-agent forks.
- Recommended planning model is top-tier (Opus / Gemini Pro / GPT-5.5 equivalent). Lighter models work for execution but planning fidelity drops.

## Out of scope (forever or for now)

- Telemetry. Ever.
- Multi-user collaboration / cloud sync.
- Mouse / drag-and-drop in the TUI.
- Per-language adapter code (we recommend tools, don't ship them).
