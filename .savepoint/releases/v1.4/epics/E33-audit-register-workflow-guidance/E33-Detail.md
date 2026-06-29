---
type: epic-design
status: planned
---

# E33: Audit Register Workflow Guidance

## Purpose

Teach agents how to use the Audit Register consistently: open an audit run, reconcile findings, map findings to work items, and preserve stable IDs across repeated audits.

## What this epic adds

- New `savepoint-audit-register` skill guidance.
- Updates to existing audit guidance so register-backed audits are preferred when `.savepoint/audit/` exists.
- Generated project instructions that explain the Audit Register without duplicating phase prompts.
- User-facing documentation for the read-only TUI section and markdown source of truth.

## Components and files

| Module | Purpose |
|--------|---------|
| `agent-skills/savepoint-audit-register/SKILL.md` | Canonical workflow for register-backed audit work |
| `agent-skills/savepoint-audit/SKILL.md` | Hand off to the register workflow when appropriate |
| `AGENTS.md` | Repo guidance for Savepoint contributors |
| `templates/project/AGENTS.md` | Generated project guidance |
| `README.md` | User-facing workflow overview |

## Architectural delta

No runtime architecture changes. This epic adds the human/agent operating procedure that makes the file-backed register reliable.

## Boundaries

**In scope:**
- Skill guidance
- Generated agent instructions
- README-level user workflow
- Manual reconciliation and proof rules

**Out of scope:**
- Editing actual findings for this repository
- Replacing release-level audit files
- Creating a CLI command for reconciliation
- Changing router state names

## Quality gates

- Guidance uses canonical lifecycle terminology.
- Guidance keeps PRD, design, task, build, and audit responsibilities separated.
- `make test` passes after documentation/template changes.

## Open decisions

None.
