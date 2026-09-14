---
type: epic-design
status: planned
---

# E48: Recover current work and the next action

## Purpose

Resume and board can use one truthful account of selected work, evidence, and next action.

V1 delivery epic for V2 product Objective O008. The current V1 lifecycle governs this work; this is the single authoritative backlog record.

## What this epic adds

Shared Next projection; read-only resume command/renderer; missing/stale/unknown evidence wording; selection and dependency explanations; practical health.

## Components and files

new internal/data/next.go; internal/data/project.go; new internal/resume/; cmd/resume.go; main.go; internal/doctor/report.go; internal/board/model.go.

Paths not yet present are proposed targets to confirm during this epic's design. Task Context Files must name exact files when the detailed plan is written.

## Architectural delta

One derived interpretation feeds multiple presentation surfaces; router selection is a hint and never a duplicate completion source.

Reference: `.savepoint/releases/v2/v2-Design.md`.

## Boundaries

**In scope:** the outcome and additions above.

**Out of scope:** No subprocess verification from resume, hidden selection writes, or separate health policy engine.

## Dependencies

- E44-issues-objective-checks

## Quality gates

No-write resume including errors; exact selection; no-task onboarding; replan/Check/owner/integration precedence; no claims of automatic evidence verification.

Implementation handoff requires `make build && make test` and named outcome evidence. Epic closeout requires a fresh independent V1 audit; this planning session is not that audit. Apply relevant current Guardrails, keeping STYLE advisory.

## Open decisions

Use section 11 and Task example as a contract direction; confirm preceding APIs before detailed Tasks.

