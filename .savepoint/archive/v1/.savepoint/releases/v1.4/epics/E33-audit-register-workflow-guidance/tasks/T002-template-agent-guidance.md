---
id: E33-audit-register-workflow-guidance/T002-template-agent-guidance
status: done
objective: Update generated project guidance for the audit-register workflow.
depends_on:
    - E33-audit-register-workflow-guidance/T001-audit-register-skill
complexity_tier: low
complexity_reason: The task is focused on generated markdown guidance.
---

# T002: Template Agent Guidance

## Problem

Generated projects need concise instructions that route agents to the audit-register workflow without bloating every phase prompt.

## Context Files

- `templates/project/AGENTS.md`
- `templates/project/.savepoint/router.md`
- `templates/project/.savepoint/audit/prompt.md`
- `agent-skills/savepoint-audit-register/SKILL.md`

## Acceptance Criteria

- [x] Generated `AGENTS.md` explains when to use the Audit Register.
- [x] The instructions preserve existing router state and task lifecycle terminology.
- [x] Audit Register guidance does not duplicate the full skill body.
- [x] The generated router template remains unchanged unless audit-register routing is explicitly needed.
- [x] Context-budget guidance remains compatible with the register workflow.

## Implementation Plan

- [x] Add a short Audit Register section to the generated agent guide.
- [x] Point agents to `.savepoint/audit/prompt.md` and the register skill.
- [x] Keep state names, task statuses, and defect lifecycle language canonical.
- [x] Verify the generated router template does not gain speculative states.

## Context Log

Read: templates/project/AGENTS.md, templates/project/.savepoint/router.md,
templates/project/.savepoint/audit/prompt.md, agent-skills/savepoint-audit-register/SKILL.md
(from T001).

Replaced the single register bullet in the generated guide's Audit section with a dedicated
`## Audit Register` section: file layout in one line, when to activate the
`savepoint-audit-register` skill, prompt-first read, stable `F###` reconciliation, proof
before `verified`, user-only waived/owner_decision, and the read-only board `A` overlay.
Four bullets plus one intro line — no duplication of the skill body. Router template
untouched (verified via git diff: no changes; existing `audit-pending` routing reaches the
register skill through the savepoint-audit hand-off, so no new states). Section is
compatible with the context budget: audit files are read only when audit work starts.
Quality gates: `make build && make test` pass (run after T003; see T003 context log).
