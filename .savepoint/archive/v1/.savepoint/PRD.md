---
type: project-prd
status: active
---

# Savepoint — Product Vision

## V2 product direction

The next release serves AI-assisted solo builders through **Idea → Design → Task → Check**: a simpler visible workflow and a more reliable execution handoff. Strong planning prepares one Objective's Tasks at a time; efficient executors deliver bounded, verifiable outcomes; fresh checker sessions verify technical integrity. Every Task belongs to an Objective, including tiny projects.

Technical clearance and owner acceptance are separate. Owner validation is required where the planner identifies meaningful behavior or decisions for the owner to assess; technical Tasks may complete after independent clearance. Checks preserve material findings as Issues, with Defect retained as a type. Stable identity and checked evidence support project continuity through the board and read-only `resume`.

The first V2 release includes safe migration of active work, archived historical records, coherent agent guidance, lifecycle gates, onboarding, board changes, and resume. Automatic evidence capture, hooks, and push enforcement are deferred. Context is focused and measured rather than constrained by a universal byte limit. Software tests and small repeatable agent evaluations provide evidence of workflow quality; they do not guarantee correctness across all agents.

The confirmed scope and ordered delivery backlog are in [v2-PRD.md](releases/v2/v2-PRD.md); proposed technical contracts are in [v2-Design.md](releases/v2/v2-Design.md). The existing V1 workflow governs this repository until the tested cutover. Future V2 ownership rules do not change current task-completion authority.

## Existing V1 baseline

The following describes the existing product and its historical targets. V2 supersedes the hierarchy, mandatory owner acceptance for technical Tasks, and fixed context-budget targets through the explicit migration and workflow cutover above.

## What it is

A public OSS CLI + Bubble Tea TUI that scaffolds an opinionated AI-driven development workflow. The user runs `npx savepoint init` in an empty directory, points any AI agent (Claude / Cursor / Cline / Gemini / Aider / Codex) at the project, and the embedded prompt templates carry the agent through:

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
