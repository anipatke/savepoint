---
id: E38-health-check-template/T001-create-health-check-template
status: planned
objective: Create the genericized Health-Check.md scaffold template at templates/project/.savepoint/Health-Check.md with Quick/Full/Deep modes, check procedures, and output templates.
depends_on: []
complexity_tier: low
---

# T001: Create Health-Check.md Template

## Problem

No shipped Health-Check.md template exists for projects to define health check modes, check procedures, and evidence templates. The user has a customised version from another project that needs genericizing.

## Context Files

- `examples.md` (lines 1-135 — reference content)
- `templates/project/.savepoint/Health-Check.md` (will create)
- `templates/project/.savepoint/Design.md` (update directory listing)

## Acceptance Criteria

- [ ] `templates/project/.savepoint/Health-Check.md` exists with frontmatter: `type: health-check`, `status: active`, `last_audited: never`.
- [ ] Template contains a `## Purpose` section explaining health checks as gates that produce compact evidence.
- [ ] Template contains a `## Modes` table with Quick, Full, and Deep rows (When/Output).
- [ ] Template contains a `## Quick Check` section with inputs, check list, and output template.
- [ ] Template contains a `## Full Check` section with inputs, check list, and output format.
- [ ] Template contains a `## Deep Check` section with inputs, check list, and output format.
- [ ] Template contains a `## Rule Boundary` section defining that health checks only fail work on rules defined in Guardrails.md.
- [ ] No `.claude` migration reference remains.
- [ ] No WSL/Codex sandbox note remains.
- [ ] All `GUARDRAILS.md` references are updated to `.savepoint/Guardrails.md`.
- [ ] Every "release guardrails audit plan" occurrence is softened to: "the release's guardrails mapping, if your project maintains one — otherwise the relevant `.savepoint/Guardrails.md` rule IDs directly."
- [ ] Deep Check contains no project-specific content: "release Opus traceability" is removed, "every Opus critical/high finding" becomes "every critical/high audit finding," and the billing/retention/RLS/LLM-runtime concern list becomes a generic "critical cross-cutting concerns (e.g. billing, retention, auth boundaries)" placeholder.
- [ ] `templates/project/.savepoint/Design.md` directory listing includes a `Health-Check.md` row.
- [ ] Integration test expected-file list (`internal/init/integration_test.go:63-76`) includes a `Health-Check.md` existence assertion.
- [ ] `make build && make test` passes.

## Implementation Plan

- [ ] Write Health-Check.md template body.
- [ ] Write frontmatter with corrected type and last_audited.
- [ ] Remove `.claude` migration reference.
- [ ] Remove WSL/Codex sandbox note.
- [ ] Update `GUARDRAILS.md` paths to `.savepoint/Guardrails.md`.
- [ ] Soften "release guardrails audit plan" references to the optional-mapping wording.
- [ ] Genericize Deep Check (remove Opus references, replace the project-specific concern list with a placeholder).
- [ ] Add Design.md directory listing row.
- [ ] Add `Health-Check.md` to the integration test expected-file list.
- [ ] Verify no `.claude`, WSL/Codex, Opus, or root-`GUARDRAILS.md` references remain.
- [ ] Run `make build && make test`.

## Context Log

Pending.
