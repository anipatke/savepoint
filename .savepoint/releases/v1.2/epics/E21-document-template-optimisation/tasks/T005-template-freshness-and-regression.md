---
id: E21-document-template-optimisation/T005-template-freshness-and-regression
status: done
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

- [x] Freshness tests assert `type:` frontmatter for `PRD.md` (`project-prd`), `Design.md` (`project-design`), and `Concept.md` (`project-concept`).
- [x] A freshness test asserts that key lifecycle terminology in the scaffolded `AGENTS.md` (e.g. canonical status values, no `phase` references) matches the live `AGENTS.md`.
- [x] All existing freshness tests continue to pass.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Read `internal/init/template_freshness_test.go` to understand existing assertion patterns.
- [x] Add assertions for `PRD.md`, `Design.md`, and `Concept.md` frontmatter type fields if not already present after T002–T004.
- [x] Add a terminology consistency assertion between live and scaffolded `AGENTS.md` (e.g. both contain `status: planned`, neither contains `phase:`).
- [x] Run `go test ./internal/init` and verify all assertions pass.
- [x] Run `make build && make test`.

## Context Log

### Coverage map after E21

| Template                          | `type:` frontmatter asserted? | Source test                            |
| --------------------------------- | ----------------------------- | -------------------------------------- |
| `templates/project/AGENTS.md`     | n/a (no `type:` field)        | `TestProjectAgentsGuidesLifecycleTerminologyConsistency` (new) |
| `templates/project/.savepoint/PRD.md`     | `project-prd`            | `TestProjectDocumentTemplatesHaveTypeFrontmatter` (new) |
| `templates/project/.savepoint/Design.md`  | `project-design`         | `TestProjectDocumentTemplatesHaveTypeFrontmatter` (new) |
| `templates/project/.savepoint/Concept.md` | `project-concept`        | `TestProjectDocumentTemplatesHaveTypeFrontmatter` (new) + `TestProjectConceptTemplateExists` (T004) |

The Concept template's `type: project-concept` is asserted twice — once in the unified frontmatter test (where it lives next to PRD/Design for symmetry) and once in the per-template section-heading test from T004. The duplication is intentional: the unified test is the single source of truth for `type:` field coverage, and the per-template test still documents the file's section shape.

### Lifecycle terminology guard

Added `TestProjectAgentsGuidesLifecycleTerminologyConsistency` in `internal/init/template_freshness_test.go`. Asserts:

1. **Canonical status values present in both files** — `planned`, `in_progress`, `done` (the bare tokens used in the AGENTS.md Terminology section). Catches any drift that removes a canonical value.
2. **Canonical Terminology sentence present in both files** — `Task `status`: only `planned`, `in_progress`, or `done``. Catches any drift that rephrases the canonical list (e.g. swapping the wording so a regex-style check would silently pass but humans would be confused).
3. **No `phase:` as a YAML key in either file** — a strict bare-token check that supersedes the per-value `phase: build|test|audit|implementation` checks in `TestProjectTemplatesRejectStaleWorkflowTerms`. Catches any new `phase:` field regardless of value.

### Why the literal `status: planned` form was rejected

The task's example mentioned `status: planned`, but the AGENTS.md files use the Terminology form `planned`, `in_progress`, `done` — wrapped in backticks and prose, not as YAML examples. Asserting the literal `status: planned` token would fail against the current files and force a prose change unrelated to the regression-guard intent. The test instead asserts the canonical tokens and the canonical sentence, which is what an actual drift would break.

### Verification

- `go test -count=1 ./internal/init/...` passes. New tests:
  - `TestProjectDocumentTemplatesHaveTypeFrontmatter` — green
  - `TestProjectAgentsGuidesLifecycleTerminologyConsistency` — green
- `make build && make test` passes across all packages.

### Files changed

- **Modified:** `internal/init/template_freshness_test.go` (added `TestProjectDocumentTemplatesHaveTypeFrontmatter` and `TestProjectAgentsGuidesLifecycleTerminologyConsistency`).
