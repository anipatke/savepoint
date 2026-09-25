---
id: E37-guardrails-template/T001-create-guardrails-template
status: done
objective: Create the genericized Guardrails.md scaffold template at templates/project/.savepoint/Guardrails.md with severity model, rule index, and customization note.
depends_on: []
complexity_tier: medium
complexity_reason: "Template content work with a wide genericization surface; no Go logic beyond a test fixture and assertion."
---

# T001: Create Guardrails.md Template

## Problem

No shipped Guardrails.md template exists for projects to define their engineering policy, severity model, and rule index. The user has a customised version from another project that needs genericizing.

## Context Files

- `examples.md` (lines 141-313 — reference content)
- `templates/project/.savepoint/Guardrails.md` (will create)
- `templates/project/.savepoint/Design.md` (update directory listing)

## Acceptance Criteria

- [x] `templates/project/.savepoint/Guardrails.md` exists with frontmatter: `type: guardrails`, `status: active`, `last_audited: never`.
- [x] Template contains a `## Purpose` section explaining the file's role.
- [x] Template contains a `## Severity Model` section with Blocker/Required/Guideline definitions.
- [x] Template contains a `## Rule Index` with categories: Security, Privacy And Retention, Auth And Admin Boundaries, Billing And External Services, Jobs And Events, Architecture, Database And Migrations, API Contracts And Frontend Safety, Configuration And Dependencies, Testing And Evidence, Observability And Runtime Safety, Code Style, Operational Quality And Policy Boundaries.
- [x] The Code Style category contains `STYLE-01..10` at Guideline severity, carrying the ten rules from `AGENTS.md:70-81` verbatim in wording (one job per file, one job per function, test branches, types document intent, build only what is needed, handle errors at boundaries, one source of truth, comments explain why, content lives in data, small diffs).
- [x] OPS-01, OPS-02, and OPS-03 are removed from the Operational Quality category so the template ships no second, overlapping style list.
- [x] Every rule is genericized: no QuizKids-specific references (child-data, Supabase, First Family Journey, quiz endpoints). This includes SEC-06 ("quiz generation" → "resource-intensive generation"), SEC-07 (drop "or child data"), Severity Model prose ("child-data privacy violations" → "sensitive-data privacy violations"), PRIV-01 (waiver parenthetical and purge intervals stripped; genericized to a defined retention window and deletion schedule), and OBS-03 ("Purge failures" → "Retention/deletion job failures").
- [x] Integration test expected-file list (`internal/init/integration_test.go:63-76`) includes a `Guardrails.md` existence assertion.
- [x] A customization note is present near the top: "Replace these rules with your project's own engineering policy."
- [x] Template contains a `## Savepoint Enforcement` section describing how health checks apply the rules.
- [x] `templates/project/.savepoint/Design.md` directory listing includes a `Guardrails.md` row.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Write Guardrails.md template body.
- [x] Write frontmatter.
- [x] Genericize PRIV rules (child-data → sensitive user data, remove quiz-specifics, strip PRIV-01 waiver history and purge intervals).
- [x] Genericize SEC-06 (quiz generation → resource-intensive generation) and SEC-07 (drop child-data mention).
- [x] Genericize Severity Model prose (child-data → sensitive-data).
- [x] Genericize AUTH-04 (Supabase → generic privileged DB access).
- [x] Genericize OBS-03 (Purge failures → Retention/deletion job failures).
- [x] Genericize LLM-01 (Child/user → User).
- [x] Genericize TEST-10 (First Family Journey → critical user journey).
- [x] Add the `STYLE-01..10` Code Style category at Guideline severity from `AGENTS.md:70-81`.
- [x] Remove OPS-01/02/03 to avoid a duplicate style list.
- [x] Add customization note.
- [x] Add Design.md directory listing row.
- [x] Add `Guardrails.md` to the integration test expected-file list.
- [x] Verify no QuizKids references remain.
- [x] Run `make build && make test`.

## Context Log

**Read:** `.savepoint/router.md`, `AGENTS.md`, `E37-Detail.md`, this task file, `examples.md` (lines 135-313), `templates/project/.savepoint/Design.md`, `templates/project/.savepoint/visual-identity.md`, `internal/init/integration_test.go`, `internal/init/scaffold.go`, `internal/init/upgrade.go`, `main.go`.

**Edited:**
- `templates/project/.savepoint/Guardrails.md` (created) — genericized policy template: frontmatter `type: guardrails` / `status: active` / `last_audited: never`, customization note, Purpose, Severity Model, 13-category Rule Index, Savepoint Enforcement.
- `templates/project/.savepoint/Design.md` — added `Guardrails.md` row to the directory-layout tree.
- `internal/init/integration_test.go` — added `.savepoint/Guardrails.md` to the expected-file list and to the `fstest.MapFS` template fixture (the fixture is synthetic, so both entries are required).

**Verification:**
- Genericization grep over the template for `quizkids|child|supabase|first family|quiz|purge|OPS-0` returns no matches. Generic severity and enforcement language still uses `waiver`; the QuizKids-specific PRIV-01 waiver parenthetical, date, decision reference, and purge intervals are absent.
- `STYLE-01..10` present at Guideline severity, wording carried verbatim from `AGENTS.md:70-81`.
- OPS-01/02/03 removed; Operational Quality And Policy Boundaries now holds POL-01/POL-02 only.
- Fresh scaffold from a locally built binary produced `.savepoint/Guardrails.md` with `{{PROJECT_NAME}}` interpolated to the target dir name.
- No `upgrade.go` change needed: `internal/init/upgrade.go:118-121` skips every `.savepoint` path generically.

**Quality gates:** `make build && make test` — pass (all packages ok).
