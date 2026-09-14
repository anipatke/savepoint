---
type: epic-design
status: planned
---

# E43: Make Task completion trustworthy

## Purpose

Completion reflects independent technical clearance and conditional owner acceptance.

V1 delivery epic for V2 product Objective O003. The current V1 lifecycle governs this work; this is the single authoritative backlog record.

## What this epic adds

Check records; freshness assessments; actor authority; Task dependency gates; replanning; exception evidence; diagnostics shared by all controls.

## Components and files

internal/data/lifecycle.go; dependency.go; write.go; new Check/evaluation files; internal/board/transitions.go; internal/doctor/checks.go.

Paths not yet present are proposed targets to confirm during this epic's design. Task Context Files must name exact files when the detailed plan is written.

## Architectural delta

Replace status-only eligibility with canonical evidence evaluation. Keep three Task statuses and retain V1 authority during transitional operation.

Reference: `.savepoint/releases/v2/v2-Design.md`.

## Boundaries

**In scope:** the outcome and additions above.

**Out of scope:** No automatic freshness detection, model authentication, push hooks, or checker repair.

## Dependencies

- E42-project-schema-identity

## Quality gates

Clean/failed/missing/stale/unknown Checks, owner-required vs technical Tasks, accepted dependencies, scope changes, and invalid external edits have explicit results.

Implementation handoff requires `make build && make test` and named outcome evidence. Epic closeout requires a fresh independent V1 audit; this planning session is not that audit. Apply relevant current Guardrails, keeping STYLE advisory.

## Open decisions

Use design sections 4–5; resolve API names and exact YAML with existing implementation after E42.

