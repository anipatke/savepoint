---
id: E33-audit-register-workflow-guidance/T003-user-workflow-documentation
status: planned
objective: Document how users review audit prompt, findings, and runs in markdown and the TUI.
depends_on:
  - E32-audit-register-tui-review/T004-help-footer-and-render-regression
  - E33-audit-register-workflow-guidance/T002-template-agent-guidance
complexity_tier: low
complexity_reason: The task adds user-facing documentation for already-planned behavior.
---

# T003: User Workflow Documentation

## Problem

Users need to know where audit-register files live, how to ask agents to use them, and how the read-only TUI review section fits the workflow.

## Context Files

- `README.md`
- `templates/project/AGENTS.md`
- `templates/project/.savepoint/audit/register.md`
- `templates/project/.savepoint/audit/runs/README.md`
- `templates/project/.savepoint/audit/findings/README.md`

## Acceptance Criteria

- [ ] README explains the Audit Register at a high level.
- [ ] Documentation names `A` as the board shortcut for read-only audit review.
- [ ] Documentation explains that markdown files remain the editable source of truth in v1.4.
- [ ] Documentation describes stable finding IDs, run history, and proof-based closure.
- [ ] Documentation keeps dashboards, external trackers, and automated matching out of v1.4.

## Implementation Plan

- [ ] Add a concise Audit Register section to README.
- [ ] Explain the `.savepoint/audit/` file layout.
- [ ] Explain the TUI review path and read-only boundary.
- [ ] Cross-reference generated project guidance where appropriate.

## Context Log

Pending.
