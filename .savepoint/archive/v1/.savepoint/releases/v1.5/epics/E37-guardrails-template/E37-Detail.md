---
type: epic-design
status: audited
---

# E37: Guardrails.md Template

## Purpose

Add a new Guardrails.md template to the shipped project template bundle, providing a canonical file for project-level engineering policy, severity model, and rule index. Content is adapted from a customised Savepoint project: QuizKids-specific rules are genericized (child-data → sensitive data, Supabase → generic privileged access, First Family Journey → critical user journey), and a customization note is added.

## What this epic adds

- `templates/project/.savepoint/Guardrails.md` — genericized engineering policy template with severity model, rule index, and Savepoint enforcement section.
- Updated `templates/project/.savepoint/Design.md` directory listing to include `Guardrails.md`.

## Components and files

| Module | Purpose |
|--------|---------|
| `templates/project/.savepoint/Guardrails.md` | The new scaffold template — frontmatter `type: guardrails`, purpose, severity model (Blocker/Required/Guideline), rule index with genericized categories including `STYLE-01..10` at Guideline severity, Savepoint enforcement |
| `templates/project/.savepoint/Design.md` | Directory listing row addition |

## Architectural delta

Template file only — embedded at `main.go:18`, scaffolded verbatim at `internal/init/scaffold.go:32-51`, skipped by `upgrade-assets` at `internal/init/upgrade.go:118-121`. No Go code changes required. `internal/init/integration_test.go:63-76` enumerates an explicit expected-file list; a `Guardrails.md` existence assertion must be added there.

## Genericization required

- "QuizKids" → `{{PROJECT_NAME}}` (interpolated at init by `internal/init/scaffold.go:64`)
- PRIV rules: "child data" → "sensitive user data"; remove quiz-specific purge schedules and answer-exposure rules
- PRIV-01: strip the owner-signed waiver parenthetical (`docs/decisions.md` reference, 2026-07-04 date, purge intervals); genericize to "sensitive user-authored text must have a defined retention window and deletion schedule"
- SEC-06: "Auth, quiz generation, and billing endpoints" → "Auth, resource-intensive generation, and billing endpoints"
- SEC-07: drop "or child data" from the exposure list ("sensitive user data" covers it)
- Severity Model prose: "child-data privacy violations" → "sensitive-data privacy violations"
- AUTH-04: "service-role Supabase access" → "privileged database or service-role access bypassing row-level security"
- OBS-03: "Purge failures and billing errors" → "Retention/deletion job failures and billing errors"
- LLM-01: "Child/user-authored" → "User-authored"
- TEST-10: "First Family Journey" → "critical user journey"
- All other rules (SEC-01-05, SEC-08, AUTH-01-03, ARCH, DATA, API, FE, CFG, DEP, TEST-01-09, OBS-01-02, RUN, OPS, POL, JOB, EXT, BILL, LLM-02) are already generic
- Add customization note: "Replace these rules with your project's own engineering policy"

## Code-style rules

The template also becomes the home for code style, which today lives in the managed AGENTS.md block and is overwritten on every upgrade. E39 points all guidance at these rules; E37 authors them so the duplication is never shipped.

- Add a `### Code Style` category with `STYLE-01..10` at Guideline severity, carrying the ten existing rules from `AGENTS.md:70-81` verbatim in wording: one job per file, one job per function, test branches, types document intent, build only what is needed, handle errors at boundaries, one source of truth, comments explain why, content lives in data, small diffs.
- Guideline severity matches the existing framing in `templates/project/.savepoint/Design.md:122` that code-style review is advisory, not blocking.
- Remove OPS-01, OPS-02, and OPS-03 from the Operational Quality category. They restate "handle errors at boundaries" and "types document intent" in language specific to Python (`print()`, type hints) and would ship as a second, overlapping style list.

## Boundaries

**In scope:**

- The Guardrails.md template file with genericized content
- Minimal Design.md directory listing addition
- Integration test existence assertion (`internal/init/integration_test.go:63-76`)

**Out of scope:**

- Go code changes beyond a file-existence assertion
- Health-Check.md template (separate epic E38)
- Runtime enforcement of guardrails content
- Cross-reference updates in AGENTS.md, audit-register, or build-task (handled by E35)

## Quality gates

- `make build && make test` passes.
- Fresh scaffold includes `Guardrails.md` at the expected path.
- `templates/project/.savepoint/Design.md` directory listing includes `Guardrails.md`.
- No QuizKids references remain in the template body.
- Customization note is present in the template.
- The `STYLE-01..10` category is present at Guideline severity and matches the ten rules in `AGENTS.md:70-81`.
- OPS-01, OPS-02, and OPS-03 are absent, so no second overlapping style list ships.

## Open decisions

None.
