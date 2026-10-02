---
name: savepoint-idea
description: Guides Savepoint idea intake when router state is idea, turning a rough sentence into .savepoint/Idea.md through owner questions and handing off to design without resolving architecture or product uncertainty itself.
---

# Savepoint Skill: Idea

## Purpose

Turn a rough, unstructured idea into `.savepoint/Idea.md` through a short back-and-forth with the owner. Own intent, users, core experience, scope, exclusions, and observable success; do not design a solution.

Every Savepoint project has at least one live Goal selected by the router, and every live Objective names exactly one Goal through `release:`. `savepoint init` creates G-001, titled after the project; use the owner's input during Idea intake to fill its Outcome, Why, Success Conditions, and Boundaries. Projects converted by `savepoint migrate` may carry R-### Goal IDs; treat them like G-### Goals. If Next says `Choose a Goal`, or `savepoint doctor` reports a missing Goal, report it to the owner. A Goal does not own Tasks or publish, deploy, tag, or generate changelogs.

## Trigger

Use this skill when router `state` is `idea`. Legacy input requires `savepoint migrate` before the V2 workflow.

## Next

Start from the `Next` line as AGENTS.md's Workflow describes; if `savepoint` is unavailable, follow AGENTS.md rather than guessing. AGENTS.md's Router Selection section says who changes the router.

## Read

- `.savepoint/router.md`
- `.savepoint/Idea.md`
- The Goal record selected by the router, when it resolves to a live Goal
- The user's stated intent for this idea
- Targeted existing-project evidence, when the idea builds on what's already there — read to understand what exists, not to design a solution

Read nothing else. Design.md, Guardrails.md, Objective or Task files, and untargeted source code are out of scope for this skill.

## Workflow

1. Read the router, any existing `.savepoint/Idea.md`, and its selected Goal record when that selection resolves.
2. Accept a single rough sentence as a valid starting input. Do not require a prepared requirements document, research document, or a completed template before the conversation starts.
3. Use the router-selected Goal as the project's planning context. On a fresh project, this is the G-001 placeholder created by `savepoint init`; tell the owner it already exists and do not create another Goal for the same initial outcome.
4. Use the Idea Interview below to fill each Idea section: Intent, User, Core Experience, Scope, Out of Scope, Success Criteria. For the fresh G-001 placeholder, also gather the owner's wording for Outcome, Why, Success Conditions, and Boundaries.
5. When grounding the idea against an existing project, read only the targeted evidence needed for that grounding; do not let it turn into designing a solution. Preserve existing Goal content; only fill the fresh scaffold placeholder from owner-provided answers.
6. Write `.savepoint/Idea.md` using the Idea artifact template below, and fill the fresh G-001 placeholder's sections with the owner's answers while preserving its identity and project-name title.
7. When the Idea is ready, set router `state: design` and hand off to `savepoint-design`. This skill does not detail any Objective itself; that belongs to `savepoint-design`. Do not write a free-text next action; follow AGENTS.md's Router Selection section for any selection changes.

## Idea Interview

Ask one focused question at a time and wait for the answer before asking the next. Offer concrete choices when useful, using an available structured question tool or plain text when none is available. Briefly acknowledge each decision and follow up on uncertainty that affects another Idea section; do not ask again about settled choices unless new evidence or changed intent requires it.

Use targeted evidence allowed by Read to answer factual questions yourself. Ask the owner directly about intent and material product uncertainty; never resolve those choices by inference. Keep architecture and implementation decisions for `savepoint-design`. Stop when the Idea sections and, when applicable, fresh Goal placeholder have enough owner-provided detail for the design handoff; summarise the settled intent, scope, exclusions, and success criteria before writing them.

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

- Write `.savepoint/Idea.md`, owner-provided sections in the fresh G-001 placeholder, and the routing handoff. Do not create another Goal or write Design, Guardrails, Objectives, Tasks, Checks, Issues, or production code.
- Do not design architecture or name components/interfaces.
- Do not detail any Objective; hand off to `savepoint-design` for that.
- Use `state` only for router phase, task `status` only for task lifecycle, and `stage` only when an item is `in_progress`.
