---
name: savepoint-draft-prd
description: Guides Savepoint PRD drafting and refinement before implementation, interviewing the user until the product scope is clear enough for design.
---

# Savepoint Skill: Draft PRD

## Purpose

Create or refine `.savepoint/PRD.md` before implementation begins. This skill owns product scope only: what is being built, for whom, why it matters, V1 boundaries, and explicit out-of-scope items.

## Trigger

Use this skill when router `state` is `pre-implementation` and the project PRD is missing, incomplete, or being changed.

## Read

- `.savepoint/router.md`
- `.savepoint/PRD.md`
- User-provided idea, outline, or source notes

## Workflow

1. Read the router and PRD.
2. If the product scope is unclear, ask focused questions before editing.
3. Fill the PRD using the existing structure and remove template-only instructional comments.
4. Keep technical design, epic breakdown, and task plans out of this phase.
5. When the PRD is ready for review, update the router `next_action` for the user to review and approve the PRD before design begins.

## Rules

- Do not write code.
- Do not design architecture.
- Do not create epics or tasks.
- Use `state` only for router phase, task `status` only for task lifecycle, and `stage` only when an item is `in_progress`.
