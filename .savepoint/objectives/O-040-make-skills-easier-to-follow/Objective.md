---
id: O-040
title: Make the skills easier to follow
status: planned
depends_on: [O-039]
release: G-002
priority: medium
rank: 2
---

# O-040: Make the skills easier to follow

## Outcome

An agent reaches every rule it must apply in one step from its skill or from the always-loaded guide. It can see the full scope of a long reference from its first lines, and it tracks a Full Objective Check against a copyable progress checklist. After writing a record, it validates and fixes the record, and the skill descriptions name the `Next` words that select each skill.

## Why

Some reference chains are two or three hops deep: Check → `check-method.md` → AGENTS.md, and Task → AGENTS.md Worktree Lanes → `issue-capture.md`. `check-method.md` is 362 lines with 20 headings and no contents list. The Full Check's required steps are easy to skip without a checklist.

## Success Conditions

- Detailed when this Objective is next.

## Boundaries

**In scope:** one-hop references, contents lists for references over 100 lines, a Full Check progress checklist, a write → `savepoint resume` → fix loop after record edits, and trigger words in skill descriptions.

**Out of scope:** new commands or validators, changed rule content, and work already done by O-039.
