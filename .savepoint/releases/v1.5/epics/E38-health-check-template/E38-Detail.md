---
type: epic-design
status: planned
---

# E38: Health-Check.md Template

## Purpose

Add a new Health-Check.md template to the shipped project template bundle, providing a canonical file for health check modes (Quick/Full/Deep), check procedures, and evidence output templates. Content is adapted from a customised Savepoint project: environment-specific sandbox notes and migration references are removed.

## What this epic adds

- `templates/project/.savepoint/Health-Check.md` — genericized health check template with Quick/Full/Deep modes, check procedures, output templates, and rule boundary.
- Updated `templates/project/.savepoint/Design.md` directory listing to include `Health-Check.md`.

## Components and files

| Module | Purpose |
|--------|---------|
| `templates/project/.savepoint/Health-Check.md` | The new scaffold template — frontmatter `type: health-check`, modes table (Quick/Full/Deep), check inputs/procedures, output templates, rule boundary |
| `templates/project/.savepoint/Design.md` | Directory listing row addition |

## Architectural delta

Template file only — embedded at `main.go:18`, scaffolded verbatim at `internal/init/scaffold.go:32-51`, skipped by `upgrade-assets` at `internal/init/upgrade.go:118-121`. No Go code changes required. `internal/init/integration_test.go:63-76` enumerates an explicit expected-file list; a `Health-Check.md` existence assertion must be added there.

## Cleanup required

- Frontmatter `type: health-check-ceremony` → `type: health-check`
- Frontmatter `last_reviewed: 2026-05-15` → `last_audited: never`
- Remove line: "Do not load old `.claude` health-check skills."
- Remove WSL/Codex sandbox note (environment-specific, ~13 lines)
- `GUARDRAILS.md` references → `.savepoint/Guardrails.md` throughout
- "release guardrails audit plan" (all occurrences) → softened optional reference: "the release's guardrails mapping, if your project maintains one — otherwise the relevant `.savepoint/Guardrails.md` rule IDs directly"
- Deep Check inputs: remove "release Opus traceability" (project-specific tooling reference)
- Deep Check checks: "every Opus critical/high finding…" → generic "every critical/high audit finding is closed or explicitly accepted by owner"; the "billing replay, retention, auth boundaries, service-role/RLS posture, frontend/API contracts, async LLM runtime, and staging journey evidence" list → placeholder "your project's critical cross-cutting concerns (e.g. billing, retention, auth boundaries) are coherent across epics"
- Keep all other content: Purpose, Modes table, Quick/Full/Deep check structure, output templates, Rule Boundary

## Boundaries

**In scope:**

- The Health-Check.md template file with genericized content
- Minimal Design.md directory listing addition
- Integration test existence assertion (`internal/init/integration_test.go:63-76`)

**Out of scope:**

- Go code changes beyond a file-existence assertion
- Guardrails.md template (separate epic E37)
- Overlap with `internal/doctor/` health-check functionality (this is a user-facing template, not a Go design document)
- Cross-reference updates in AGENTS.md, audit-register, or build-task (handled by E35)

## Quality gates

- `make build && make test` passes.
- Fresh scaffold includes `Health-Check.md` at the expected path with `type: health-check` frontmatter.
- `templates/project/.savepoint/Design.md` directory listing includes `Health-Check.md`.
- No `.claude` reference or WSL/Codex sandbox note remains.
- No project-specific Deep Check content remains (Opus traceability, billing/RLS/LLM-runtime concern list).
- No root `GUARDRAILS.md` path or mandatory "release guardrails audit plan" reference remains.

## Open decisions

None.
