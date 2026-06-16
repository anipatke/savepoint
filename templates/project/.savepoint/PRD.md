---
type: project-prd
status: active
---

# {{PROJECT_NAME}} — Product Vision

> Fill each section with concrete, falsifiable claims. Vague PRDs produce vague agents.

## What it is

State the product in three sentences: (1) the category it sits in, (2) the user action that proves it works, (3) the pipeline or lifecycle it commits to.

> A public OSS CLI + Bubble Tea TUI that scaffolds an opinionated AI-driven development workflow. The user runs `npx <tool> init` in an empty directory, points any AI agent at the project, and the embedded prompt templates carry the agent through **PRD → Design → Epics → Tasks → Build → Audit** with hard gates at each transition.

## Why

List the 2–3 concrete failure modes this product removes. For each, name the mechanism that fixes it in one sentence.

> 1. **Inconsistency.** No repeatable process from intent to MVP. Fixed by a file-based state machine any agent can follow.
> 2. **Token bloat.** Monolithic backlogs burn context. Fixed by a per-task read budget enforced in the agent guide.
> 3. **Documentation drift.** `Design.md` goes stale. Fixed by an audit gate before the next epic can start.

## Target user

Name the primary persona in a bold tag, define them in one sentence, then name the persona you are explicitly **not** for.

> **Vibe coders** — builders with minimal-to-moderate dev experience on minimal AI plans who want AI agents to drive most of the implementation.
> Not: experienced engineers with their own systems. (They can still use it; they're not the audience.)

## Headline differentiator

What is the one feature a competitor cannot copy in a weekend? State it as a bold tag and justify in one sentence. If a competitor already has it, pick a different feature.

> **The Audit Loop.** When the last task in an epic moves to `done`, the next epic cannot start until `Design.md`, `AGENTS.md`, and the epic's own design have been reconciled with the actual code that was built. No existing markdown-first task tool has this gate.

## Success metrics

List 3–4 measurable V1 outcomes. Each metric must be falsifiable — a number, a time bound, or a test the agent can run.

> - **Token budget:** agents complete a task reading <2KB of context per task. Audit budget bounded to ~5–15KB.
> - **Documentation accuracy:** zero drift — `AGENTS.md` always correctly maps the current codebase, enforced by gate.
> - **Agent reach:** works with any agent that can read markdown and edit files. No MCP required, no per-agent adapters.
> - **Time-to-first-PR:** a vibe coder can go from `npx <tool> init` to a merged epic in one weekend.

## Constraints

Hard limits only: architecture, platform, or stack refusals. No preferences, no nice-to-haves. If you find yourself hedging, cut the line.

> - File-only architecture for v1. No MCP server.
> - Agent-agnostic via the [Router Pattern](Design.md). No per-agent forks.
> - Recommended planning model is top-tier (Opus / Gemini Pro / GPT-5.5 equivalent). Lighter models work for execution but planning fidelity drops.

## Out of scope (forever or for now)

What are you explicitly **not** building? Use the strongest phrasing you can — "No X. No Y. No Z." — because this list is what keeps the AI focused.

> - No telemetry. Ever.
> - No multi-user collaboration or cloud sync.
> - No mouse / drag-and-drop in the TUI.
> - No per-language adapter code (recommend tools, don't ship them).
