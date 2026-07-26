---
id: E37-guardrails-template/T001-create-guardrails-template
status: planned
objective: Create the genericized Guardrails.md scaffold template at templates/project/.savepoint/Guardrails.md with severity model, rule index, and customization note.
depends_on: []
complexity_tier: medium
---

# T001: Create Guardrails.md Template

## Problem

No shipped Guardrails.md template exists for projects to define their engineering policy, severity model, and rule index. The user has a customised version from another project that needs genericizing.

## Context Files

- `examples.md` (lines 141-313 — reference content)
- `templates/project/.savepoint/Guardrails.md` (will create)
- `templates/project/.savepoint/Design.md` (update directory listing)

## Acceptance Criteria

- [ ] `templates/project/.savepoint/Guardrails.md` exists with frontmatter: `type: guardrails`, `status: active`, `last_audited: never`.
- [ ] Template contains a `## Purpose` section explaining the file's role.
- [ ] Template contains a `## Severity Model` section with Blocker/Required/Guideline definitions.
- [ ] Template contains a `## Rule Index` with categories: Security, Privacy And Retention, Auth And Admin Boundaries, Billing And External Services, Jobs And Events, Architecture, Database And Migrations, API Contracts And Frontend Safety, Configuration And Dependencies, Testing And Evidence, Observability And Runtime Safety, Code Style, Operational Quality And Policy Boundaries.
- [ ] The Code Style category contains `STYLE-01..10` at Guideline severity, carrying the ten rules from `AGENTS.md:70-81` verbatim in wording (one job per file, one job per function, test branches, types document intent, build only what is needed, handle errors at boundaries, one source of truth, comments explain why, content lives in data, small diffs).
- [ ] OPS-01, OPS-02, and OPS-03 are removed from the Operational Quality category so the template ships no second, overlapping style list.
- [ ] Every rule is genericized: no QuizKids-specific references (child-data, Supabase, First Family Journey, quiz endpoints). This includes SEC-06 ("quiz generation" → "resource-intensive generation"), SEC-07 (drop "or child data"), Severity Model prose ("child-data privacy violations" → "sensitive-data privacy violations"), PRIV-01 (waiver parenthetical and purge intervals stripped; genericized to a defined retention window and deletion schedule), and OBS-03 ("Purge failures" → "Retention/deletion job failures").
- [ ] Integration test expected-file list (`internal/init/integration_test.go:63-76`) includes a `Guardrails.md` existence assertion.
- [ ] A customization note is present near the top: "Replace these rules with your project's own engineering policy."
- [ ] Template contains a `## Savepoint Enforcement` section describing how health checks apply the rules.
- [ ] `templates/project/.savepoint/Design.md` directory listing includes a `Guardrails.md` row.
- [ ] `make build && make test` passes.

## Implementation Plan

- [ ] Write Guardrails.md template body.
- [ ] Write frontmatter.
- [ ] Genericize PRIV rules (child-data → sensitive user data, remove quiz-specifics, strip PRIV-01 waiver history and purge intervals).
- [ ] Genericize SEC-06 (quiz generation → resource-intensive generation) and SEC-07 (drop child-data mention).
- [ ] Genericize Severity Model prose (child-data → sensitive-data).
- [ ] Genericize AUTH-04 (Supabase → generic privileged DB access).
- [ ] Genericize OBS-03 (Purge failures → Retention/deletion job failures).
- [ ] Genericize LLM-01 (Child/user → User).
- [ ] Genericize TEST-10 (First Family Journey → critical user journey).
- [ ] Add the `STYLE-01..10` Code Style category at Guideline severity from `AGENTS.md:70-81`.
- [ ] Remove OPS-01/02/03 to avoid a duplicate style list.
- [ ] Add customization note.
- [ ] Add Design.md directory listing row.
- [ ] Add `Guardrails.md` to the integration test expected-file list.
- [ ] Verify no QuizKids references remain.
- [ ] Run `make build && make test`.

## Context Log

Pending.
