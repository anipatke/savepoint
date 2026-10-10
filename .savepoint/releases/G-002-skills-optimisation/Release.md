---
id: G-002
title: Skills Optimisation
status: planned
---

## Outcome

Savepoint's workflow skills, shared references, and agent guide cost less context per session and are easier for an agent to follow correctly. Each rule lives in one place, references are one step from the skill that needs them, and re-walking O-034's written scenarios shows that behaviour is unchanged or intentionally improved.

## Why

A review of the skills against Anthropic's published skill-authoring best practices (2026-10-06) found avoidable cost and risk. This repository's agent guide loads two near-copies of the same policy. Shared rules such as Goal Context and Task-check waivers are restated across skills, references, and the guide. Some references sit two or three hops from the skill that needs them, and the 362-line Check method has no contents list.

## Success Conditions

- Each shared policy rule has one home, and the other files point to that home instead of restating it.
- Every reference a skill must apply is reachable in one step from that skill or from the always-loaded agent guide.
- Long references open with a contents list, and the Full Objective Check has a copyable progress checklist.
- Every skill change in this Goal is re-walked against O-034's written scenarios (T-099 baselines), with a result per scenario.
- Owner authority, checker independence, acceptance evidence, mandatory Objective Checks, and live/scaffold byte parity are unchanged.
- Existing projects receive the changes through `savepoint upgrade-assets`.

## Boundaries

**In scope:** the four workflow skills, the three shared references, the managed AGENTS.md block and this repository's own AGENTS.md, their scaffold copies, content tests, and scenario re-walks; and Claude Code integration files (`CLAUDE.md`, `.claude/skills/`, `.claude/hooks/`, `.claude/settings.json`) installed by `init` and `upgrade-assets` (O-043); and the owner-decision carry-forward across re-checks (O-044): its evidence fields, runtime gates, resume and board routing, and guidance.

**Out of scope:** `bubbletea-tui-design`; new lifecycle states, fields, or commands other than those O-044 needs; weakening verification to save tokens; live multi-model agent runs; automated eval tooling; Go runtime behaviour changes other than the `init` and `upgrade-assets` scaffolding O-043 needs and the decision-lifecycle changes O-044 needs.
