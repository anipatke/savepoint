---
id: E21-document-template-optimisation/T003-design-md-template-optimisation
status: planned
objective: Update the Design.md template so its sections and prompts match the current directory layout, architecture model, and audit workflow.
complexity_tier: low
complexity_reason: Single template file; section headings and prompts need updating, no code changes.
---

# T003: Design.md Template Optimisation

## Problem

The `templates/project/.savepoint/Design.md` template has nine numbered sections, but their headings and prompts have not been updated since the directory layout and audit workflow evolved. A scaffolded project receives a template whose section order and language may not match what the `savepoint-system-design` skill produces or what the audit loop expects.

## Context Files

- `templates/project/.savepoint/Design.md`
- `.savepoint/Design.md`
- `templates/project/agent-skills/savepoint-system-design/SKILL.md`
- `internal/init/template_freshness_test.go`

## Acceptance Criteria

- [ ] Section headings in the template match the structure the `savepoint-system-design` skill is expected to produce for a new project.
- [ ] Section prompts are concrete and match the specificity of the live `.savepoint/Design.md`.
- [ ] The audit workflow section (§7) accurately describes the skill-driven audit loop, not any CLI pipeline.
- [ ] Directory layout section (§2) matches the canonical directory structure in the live Design.md.
- [ ] Template freshness test passes.
- [ ] `make build && make test` passes.

## Implementation Plan

- [ ] Read `templates/project/.savepoint/Design.md` and the live `.savepoint/Design.md` side by side.
- [ ] Read `templates/project/agent-skills/savepoint-system-design/SKILL.md` to understand what the skill writes into this file.
- [ ] Update section headings and prompts to align with the current architecture model and audit workflow.
- [ ] Run `go test ./internal/init` to check freshness coverage.
- [ ] Run `make build && make test`.
