---
id: E38-health-check-template/T001-create-health-check-template
status: done
objective: Create the genericized Health-Check.md scaffold template at templates/project/.savepoint/Health-Check.md with Quick/Full/Deep modes, check procedures, and output templates.
depends_on: []
complexity_tier: low
complexity_reason: Single template file plus a directory-listing row and one integration-test assertion; no Go behaviour changes.
---

# T001: Create Health-Check.md Template

## Problem

No shipped Health-Check.md template exists for projects to define health check modes, check procedures, and evidence templates. The user has a customised version from another project that needs genericizing.

## Context Files

- `examples.md` (lines 1-135 — reference content)
- `templates/project/.savepoint/Health-Check.md` (will create)
- `templates/project/.savepoint/Design.md` (update directory listing)

## Acceptance Criteria

- [x] `templates/project/.savepoint/Health-Check.md` exists with frontmatter: `type: health-check`, `status: active`, `last_audited: never`.
- [x] Template contains a `## Purpose` section explaining health checks as gates that produce compact evidence.
- [x] Template contains a `## Modes` table with Quick, Full, and Deep rows (When/Output).
- [x] Template contains a `## Quick Check` section with inputs, check list, and output template.
- [x] Template contains a `## Full Check` section with inputs, check list, and output format.
- [x] Template contains a `## Deep Check` section with inputs, check list, and output format.
- [x] Template contains a `## Rule Boundary` section defining that health checks only fail work on rules defined in Guardrails.md.
- [x] No `.claude` migration reference remains.
- [x] No WSL/Codex sandbox note remains.
- [x] All `GUARDRAILS.md` references are updated to `.savepoint/Guardrails.md`.
- [x] Every "release guardrails audit plan" occurrence is softened to: "the release's guardrails mapping, if your project maintains one — otherwise the relevant `.savepoint/Guardrails.md` rule IDs directly."
- [x] Deep Check contains no project-specific content: "release Opus traceability" is removed, "every Opus critical/high finding" becomes "every critical/high audit finding," and the billing/retention/RLS/LLM-runtime concern list becomes a generic "critical cross-cutting concerns (e.g. billing, retention, auth boundaries)" placeholder.
- [x] `templates/project/.savepoint/Design.md` directory listing includes a `Health-Check.md` row.
- [x] Integration test expected-file list (`internal/init/integration_test.go:63-76`) includes a `Health-Check.md` existence assertion.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Write Health-Check.md template body.
- [x] Write frontmatter with corrected type and last_audited.
- [x] Remove `.claude` migration reference.
- [x] Remove WSL/Codex sandbox note.
- [x] Update `GUARDRAILS.md` paths to `.savepoint/Guardrails.md`.
- [x] Soften "release guardrails audit plan" references to the optional-mapping wording.
- [x] Genericize Deep Check (remove Opus references, replace the project-specific concern list with a placeholder).
- [x] Add Design.md directory listing row.
- [x] Add `Health-Check.md` to the integration test expected-file list.
- [x] Verify no `.claude`, WSL/Codex, Opus, or root-`GUARDRAILS.md` references remain.
- [x] Run `make build && make test`.

## Context Log

Files read:

- `examples.md` (health-check reference content, lines 1-135)
- `templates/project/.savepoint/Guardrails.md` (frontmatter/style reference)
- `templates/project/.savepoint/Design.md`
- `internal/init/integration_test.go`

Files edited:

- `templates/project/.savepoint/Health-Check.md` (created)
- `templates/project/.savepoint/Design.md` (directory listing row)
- `internal/init/integration_test.go` (template fixture entry + expected-file assertion)

Evidence:

- Frontmatter is `type: health-check`, `status: active`, `last_audited: never`.
- Sections present: Purpose, Modes (Quick/Full/Deep), Quick Check, Full Check, Deep Check, Rule Boundary.
- `grep -n '\.claude\|WSL\|Codex\|Opus\|GUARDRAILS' templates/project/.savepoint/Health-Check.md` → no matches (exit 1).
- Quick Check's two overlapping guardrails inputs merged into one softened bullet to avoid a duplicated reference; the "load only the rule IDs or compact sections that apply" qualifier is preserved inline.
- Real scaffold check: `savepoint init <tmp>` produced `.savepoint/Health-Check.md` with the expected frontmatter.

Quality gates:

- `make build && make test` → all packages ok (`internal/init` 0.793s).
- `gofmt -l internal/init/` flags only pre-existing `clipboard.go` and `scaffold_test.go`; the edited `integration_test.go` is clean.
