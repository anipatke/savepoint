---
id: E33-audit-register-workflow-guidance/T002-template-agent-guidance
status: planned
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

- [ ] Generated `AGENTS.md` explains when to use the Audit Register.
- [ ] The instructions preserve existing router state and task lifecycle terminology.
- [ ] Audit Register guidance does not duplicate the full skill body.
- [ ] The generated router template remains unchanged unless audit-register routing is explicitly needed.
- [ ] Context-budget guidance remains compatible with the register workflow.

## Implementation Plan

- [ ] Add a short Audit Register section to the generated agent guide.
- [ ] Point agents to `.savepoint/audit/prompt.md` and the register skill.
- [ ] Keep state names, task statuses, and defect lifecycle language canonical.
- [ ] Verify the generated router template does not gain speculative states.

## Context Log

Pending.
