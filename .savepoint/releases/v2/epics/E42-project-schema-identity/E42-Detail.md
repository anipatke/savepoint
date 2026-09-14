---
type: epic-design
status: planned
---

# E42: Load V2 work with stable identity

## Purpose

A V2 project can be read and diagnosed with globally stable IDs and preserved authored content.

V1 delivery epic for V2 product Objective O002. The current V1 lifecycle governs this work; this is the single authoritative backlog record.

## What this epic adds

Version dispatch; Objective/Task records; project index; raw/typed parse boundaries; stable references; preserving writes and named diagnostics.

## Components and files

internal/data/config.go; parser.go; discover.go; dependency.go; write.go; lifecycle.go; internal/doctor/checks.go; new internal/data/project.go and schema-specific record files.

Paths not yet present are proposed targets to confirm during this epic's design. Task Context Files must name exact files when the detailed plan is written.

## Architectural delta

Introduce V2 schema and shared project index while isolating temporary V1 support. Record ownership explicitly rather than deriving it from task path.

Reference: `.savepoint/releases/v2/v2-Design.md`.

## Boundaries

**In scope:** the outcome and additions above.

**Out of scope:** No migration apply, evidence automation, or board redesign.

## Dependencies

- E41-migration-source-fixtures

## Quality gates

Malformed versions, duplicates, moved Tasks, unknown fields/body preservation, no unsafe defaults granting completion, and reference cycles are covered.

Implementation handoff requires `make build && make test` and named outcome evidence. Epic closeout requires a fresh independent V1 audit; this planning session is not that audit. Apply relevant current Guardrails, keeping STYLE advisory.

## Open decisions

Refine exact public APIs and preserving writer contract against E41 fixtures before writing Tasks.

