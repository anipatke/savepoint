---
id: E30-audit-register-templates/T001-audit-prompt-template
status: planned
objective: Add the scaffolded reusable audit prompt template.
depends_on: []
complexity_tier: low
complexity_reason: The change is a bounded markdown template with focused scaffold expectations.
---

# T001: Audit Prompt Template

## Problem

Savepoint audits do not have a reusable prompt that requires reconciliation against prior findings, coverage accounting, or closure proof.

## Context Files

- `templates/project/.savepoint/audit/prompt.md`
- `templates/project/AGENTS.md`
- `agent-skills/savepoint-audit/SKILL.md`

## Acceptance Criteria

- [ ] The prompt template explains that each audit must reconcile against the existing register when present.
- [ ] Required per-finding fields include stable ID handling, severity, confidence, source auditor, location, guardrail IDs, proof needed, and work-item mapping.
- [ ] The prompt requires coverage notes for examined and unexamined surfaces.
- [ ] The prompt includes a short changelog section for refinements over time.
- [ ] Generated project guidance points agents to the prompt before starting a register-backed audit.

## Implementation Plan

- [ ] Add `.savepoint/audit/prompt.md` to the project template.
- [ ] Define the required audit output shape in the prompt body.
- [ ] Add a changelog section with an initial version entry.
- [ ] Update generated `AGENTS.md` guidance to mention the audit prompt when the register workflow is used.
- [ ] Align the existing audit skill wording with the new prompt location.

## Context Log

Pending.
