---
id: E21-document-template-optimisation/T005-template-freshness-and-regression
status: planned
objective: Extend template freshness tests to cover all four document templates and add a regression guard against future guidance drift between live and scaffolded AGENTS.md.
depends_on: [T001-agents-md-template-cleanup, T002-prd-template-optimisation, T003-design-md-template-optimisation, T004-concept-md-template-new]
complexity_tier: low
complexity_reason: Test-only task; extends existing freshness test patterns to new and updated templates.
---

# T005: Template Freshness and Regression Guards

## Problem

After T001–T004 land, the freshness test suite should assert that all four document templates remain present, carry the correct `type:` frontmatter value, and do not silently diverge again. Without explicit coverage, the next template-touching epic can introduce drift that only surfaces at runtime.

## Context Files

- `internal/init/template_freshness_test.go`
- `templates/project/AGENTS.md`
- `templates/project/.savepoint/PRD.md`
- `templates/project/.savepoint/Design.md`
- `templates/project/.savepoint/Concept.md`

## Acceptance Criteria

- [ ] Freshness tests assert `type:` frontmatter for `PRD.md` (`project-prd`), `Design.md` (`project-design`), and `Concept.md` (`project-concept`).
- [ ] A freshness test asserts that key lifecycle terminology in the scaffolded `AGENTS.md` (e.g. canonical status values, no `phase` references) matches the live `AGENTS.md`.
- [ ] All existing freshness tests continue to pass.
- [ ] `make build && make test` passes.

## Implementation Plan

- [ ] Read `internal/init/template_freshness_test.go` to understand existing assertion patterns.
- [ ] Add assertions for `PRD.md`, `Design.md`, and `Concept.md` frontmatter type fields if not already present after T002–T004.
- [ ] Add a terminology consistency assertion between live and scaffolded `AGENTS.md` (e.g. both contain `status: planned`, neither contains `phase:`).
- [ ] Run `go test ./internal/init` and verify all assertions pass.
- [ ] Run `make build && make test`.
