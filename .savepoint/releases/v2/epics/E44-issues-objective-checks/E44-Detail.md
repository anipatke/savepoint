---
type: epic-design
status: planned
---

# E44: Converge follow-up and verify integrated outcomes

## Purpose

Material follow-up retains identity, verified closure, and an integration gate across related Tasks.

V1 delivery epic for V2 product Objective O004. The current V1 lifecycle governs this work; this is the single authoritative backlog record.

## What this epic adds

Issue schema and deduplication references; history/proof/dispositions; linked repairs; Objective integration and dependency readiness.

## Components and files

internal/data/audit_finding.go; audit_run.go; audit_register.go; audit_backlinks.go; defect.go; dependency.go; internal/doctor/checks.go.

Paths not yet present are proposed targets to confirm during this epic's design. Task Context Files must name exact files when the detailed plan is written.

## Architectural delta

New V2 Issues replace parallel defect/finding stores; Check events replace mutable register bookkeeping. Legacy parsers remain migration inputs.

Reference: `.savepoint/releases/v2/v2-Design.md`.

## Boundaries

**In scope:** the outcome and additions above.

**Out of scope:** No generic issue tracker, automatic semantic deduplication, or mandatory Issue for every owner check.

## Dependencies

- E43-task-check-gates

## Quality gates

Repeated observations reuse identity; fixed is not resolved without proof; accepted/duplicate differ from verified; unrelated Objectives can continue.

Implementation handoff requires `make build && make test` and named outcome evidence. Epic closeout requires a fresh independent V1 audit; this planning session is not that audit. Apply relevant current Guardrails, keeping STYLE advisory.

## Open decisions

Detailed Task breakdown follows E43; type list and lifecycle direction are settled in design section 6.

