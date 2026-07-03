---
id: E33-audit-register-workflow-guidance/T001-audit-register-skill
status: done
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

- [x] The new skill defines read order for prompt, register, findings, runs, relevant work items, and audited files.
- [x] The workflow distinguishes run history from current register state.
- [x] Reconciliation requires classifying each finding as existing, new, regression, duplicate, deferred, waived, or owner-decision.
- [x] Closure guidance requires named proof before `verified`.
- [x] The skill preserves Savepoint authority rules: agents do not mark tasks done and do not grant owner waivers.

## Implementation Plan

- [x] Create `savepoint-audit-register` skill frontmatter and workflow.
- [x] Define required reads and write targets.
- [x] Document finding ID preservation and duplicate handling.
- [x] Document proof and owner-decision rules.
- [x] Update existing audit guidance to hand off to the register skill when present.

## Context Log

Read: router.md, E33-Detail.md, E33 task files, AGENTS.md, templates/project/AGENTS.md,
agent-skills/savepoint-audit/SKILL.md, .savepoint/audit templates (prompt.md, register.md,
findings/README.md, runs/README.md), internal/data/audit_finding.go (canonical finding
statuses), internal/init/template_freshness_test.go (skill mirror requirement).

Created `agent-skills/savepoint-audit-register/SKILL.md` with numbered read order
(prompt → register → findings → runs → linked work items → audited files), a "run history
vs register state" section (immutable runs vs mutable register), a seven-way reconciliation
classification (existing/new/regression/duplicate/deferred/waived/owner-decision), stable
`F###` ID rules, fixed-vs-verified proof rules, and authority rules (no task `done`, waived
and owner_decision recorded only on explicit user decision). Mirrored the skill to
`templates/project/agent-skills/savepoint-audit-register/SKILL.md` (freshness test requires
identical copies). Updated `agent-skills/savepoint-audit/SKILL.md` (both copies) to hand
off to the register skill when `.savepoint/audit/` exists, and added the register pointer
to the live AGENTS.md Audit section. Quality gates: `make build && make test` pass (run
after T003; see T003 context log).
