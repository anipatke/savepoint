---
id: E30-audit-register-templates/T002-register-and-finding-templates
status: planned
objective: Add scaffolded register and finding-folder guidance for durable audit findings.
depends_on:
  - E30-audit-register-templates/T001-audit-prompt-template
complexity_tier: low
complexity_reason: The task adds file templates and field guidance without parsing behavior.
---

# T002: Register and Finding Templates

## Problem

There is no durable markdown location for current audit findings, stable finding IDs, or convergence summaries.

## Context Files

- `templates/project/.savepoint/audit/register.md`
- `templates/project/.savepoint/audit/findings/README.md`
- `templates/project/.savepoint/audit/prompt.md`

## Acceptance Criteria

- [ ] The register template includes columns or sections for ID, title, status, severity, confidence, linked work item, first seen, last seen, and proof.
- [ ] The findings folder guidance defines the `F###-slug.md` naming convention.
- [ ] Finding guidance lists required frontmatter and body sections.
- [ ] Status values match the v1.4 PRD lifecycle.
- [ ] The register template includes a convergence summary area for net-new, reopened, verified, deferred, and coverage-gap counts.

## Implementation Plan

- [ ] Add the repo-wide register template.
- [ ] Add findings folder guidance with a complete finding record shape.
- [ ] Document stable ID and duplicate handling rules.
- [ ] Document proof requirements for `verified` findings.
- [ ] Cross-reference the prompt template's reconciliation requirements.

## Context Log

Pending.
