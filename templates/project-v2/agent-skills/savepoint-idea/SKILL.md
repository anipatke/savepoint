---
name: savepoint-idea
description: Guides Savepoint idea intake when router state is idea, turning a rough sentence into .savepoint/Idea.md through owner questions and handing off to design without resolving architecture or product uncertainty itself.
---

# Savepoint Skill: Idea

## Purpose

Turn a rough, unstructured idea into `.savepoint/Idea.md` through a short back-and-forth with the owner. This skill owns intent and boundary only: what is being built, for whom, the core experience, what's in and out of scope, and how success is judged. It does not design a solution.

Goal is optional planning context, not a required phase. Use it only when the owner wants related Objectives grouped under a navigable outcome. Existing V2 Goals remain stored as `R-###` Release records with `release:` references; in the V2 board, `g` is the canonical selector and `r` remains an undisplayed compatibility alias. A Goal does not own Tasks or publish, deploy, tag, or generate changelogs.

## Trigger

Use this skill when router `state` is `idea` in a V2 project. Legacy input is handled by the explicit `savepoint migrate` workflow before this V2 state is available; this skill does not route or maintain a legacy lifecycle.

## Next

Run the read-only `savepoint resume` command and act on its `Next` line; a `Next` line pasted by the owner is also an explicit selection. This is the only Savepoint CLI command agents may run. If the binary is unavailable, follow AGENTS.md: read `.savepoint/router.md`, report the missing tool, and do not guess the next step. For an owner Task closure, use the board's router advance; see AGENTS.md's Router Selection section for the other selection owners and the `release:` rule.

## Read

- `.savepoint/router.md`
- `.savepoint/Idea.md`
- The user's stated intent for this idea
- Targeted existing-project evidence, when the idea builds on what's already there — read to understand what exists, not to design a solution

Read nothing else. Design.md, Guardrails.md, Objective or Task files, and untargeted source code are out of scope for this skill.

## Workflow

1. Read the router and any existing `.savepoint/Idea.md`.
2. Accept a single rough sentence as a valid starting input. Do not require a prepared requirements document, research document, or a completed template before the conversation starts.
3. Ask whether the owner wants a Goal to group related Objectives. Treat it as optional: carry that context forward only when the owner says it is useful; otherwise continue with Objective → Task. Do not infer a Goal from project size or make it a prerequisite for a small project.
4. Ask focused questions to fill each Idea section: Intent, User, Core Experience, Scope, Out of Scope, Success Criteria.
5. When a question turns on material product uncertainty — a choice only the owner can make, not one inferable from context — ask the owner directly. Do not resolve product choices by inference.
6. When grounding the idea against an existing project, read only the targeted evidence needed for that grounding; do not let it turn into designing a solution.
7. Write `.savepoint/Idea.md` using the Idea artifact template below.
8. When the Idea is ready, set router `state: design` and hand off to `savepoint-design`. This skill does not detail any Objective itself; that belongs to `savepoint-design`. Do not write a free-text next action; follow AGENTS.md's Router Selection section for any selection changes.

## Idea Artifact Template

Write `.savepoint/Idea.md` with this structure:

```markdown
---
type: idea
status: planned
---

# Idea: <Short Name>

## Intent

What this idea is and why it matters.

## User

Who this is for.

## Core Experience

The essential experience this idea delivers.

## Scope

What's included.

## Out of Scope

What's explicitly excluded.

## Success Criteria

How to tell this idea succeeded.
```

## Rules

- Write only `.savepoint/Idea.md` and the routing handoff. Do not write Design, Guardrails, Objectives, Tasks, Checks, Issues, or production code.
- Do not design architecture or name components/interfaces.
- Do not detail any Objective; hand off to `savepoint-design` for that.
- Ask the owner about material product uncertainty instead of deciding it by inference.
- Use `state` only for router phase, task `status` only for task lifecycle, and `stage` only when an item is `in_progress`.
