---
id: E33-audit-register-workflow-guidance/T001-audit-register-skill
status: planned
objective: Add the canonical skill for register-backed audit work.
depends_on:
  - E30-audit-register-templates/T002-register-and-finding-templates
  - E31-audit-register-data-model/T004-audit-register-validation
complexity_tier: medium
complexity_reason: The skill must encode lifecycle authority and reconciliation rules precisely.
---

# T001: Audit Register Skill

## Problem

Agents need one canonical workflow for opening audit runs, reconciling findings, updating the register, and preserving stable finding IDs.

## Context Files

- `agent-skills/savepoint-audit-register/SKILL.md`
- `agent-skills/savepoint-audit/SKILL.md`
- `AGENTS.md`
- `templates/project/AGENTS.md`

## Acceptance Criteria

- [ ] The new skill defines read order for prompt, register, findings, runs, relevant work items, and audited files.
- [ ] The workflow distinguishes run history from current register state.
- [ ] Reconciliation requires classifying each finding as existing, new, regression, duplicate, deferred, waived, or owner-decision.
- [ ] Closure guidance requires named proof before `verified`.
- [ ] The skill preserves Savepoint authority rules: agents do not mark tasks done and do not grant owner waivers.

## Implementation Plan

- [ ] Create `savepoint-audit-register` skill frontmatter and workflow.
- [ ] Define required reads and write targets.
- [ ] Document finding ID preservation and duplicate handling.
- [ ] Document proof and owner-decision rules.
- [ ] Update existing audit guidance to hand off to the register skill when present.

## Context Log

Pending.
