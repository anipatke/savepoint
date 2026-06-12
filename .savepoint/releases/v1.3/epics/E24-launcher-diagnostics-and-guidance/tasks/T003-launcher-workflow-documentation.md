---
id: E24-launcher-diagnostics-and-guidance/T003-launcher-workflow-documentation
status: planned
objective: Document launcher setup, scope guarantees, role behavior, and platform limitations.
depends_on:
  - E23-board-launch-actions/T004-epic-audit-action-and-help
  - E24-launcher-diagnostics-and-guidance/T002-scaffold-launcher-config
complexity_tier: medium
complexity_reason: Guidance must align live docs, scaffold instructions, and existing skill-owned workflow rules.
---

# T003: Launcher Workflow Documentation

## Problem

Users and launched agents need precise guidance on opt-in setup, action eligibility, bounded scope, and what Savepoint does not enforce or automate.

## Context Files

- `README.md`
- `AGENTS.md`
- `templates/project/AGENTS.md`
- `agent-skills/savepoint-build-task/SKILL.md`
- `agent-skills/savepoint-audit/SKILL.md`
- `templates/project/agent-skills/savepoint-build-task/SKILL.md`
- `templates/project/agent-skills/savepoint-audit/SKILL.md`
- `.savepoint/PRD.md`
- `.savepoint/Design.md`

## Acceptance Criteria

- [ ] README documents opt-in configuration, builder/auditor roles, keys, terminal behavior, and troubleshooting.
- [ ] Guidance states that scope is workflow-enforced rather than an OS filesystem sandbox.
- [ ] Build guidance requires the selected task or defect as the entrypoint and forbids unrelated backlog work.
- [ ] Audit guidance preserves the fresh-session epic rule and explains optional item audits.
- [ ] Manual workflows remain fully documented and no phase instructions are duplicated unnecessarily.
- [ ] Project vision and design are updated only for the approved launcher capability and architecture delta.

## Implementation Plan

- [ ] Add a concise launcher section and configuration example to README.
- [ ] Add minimal launch-entrypoint language to live and scaffolded agent guidance.
- [ ] Update build/audit skills only where launched-session behavior needs an explicit rule.
- [ ] Record the optional process-launch boundary in project PRD and Design after implementation is verified.
- [ ] Review terminology for Build, Audit, task status, defect status, and stage consistency.

## Context Log

Pending.
