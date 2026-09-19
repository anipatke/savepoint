---
type: epic-design
status: planned
---

# E50: Release evaluated V2 and retire transitional runtime

## Purpose

V2 is validated on realistic projects, packaged for supported platforms, and ready for explicit maintainer migration.

V1 delivery epic for V2 product Objective O010. The current V1 lifecycle governs this work; this is the single authoritative backlog record.

## What this epic adds

Agent evaluation scenarios; end-to-end migration trials; packaging/docs/help; legacy live-reader removal; archive/reference access; own-project cutover runbook.

## Components and files

internal/data/ transitional dispatch; internal/board/; internal/doctor/; internal/init/; internal/migrate/; internal/buildtool/main.go; bin/savepoint.js; package.json; Makefile; .github/workflows/ci.yml; .github/workflows/publish.yml; README.md; .savepoint/Design.md; .savepoint/Guardrails.md; AGENTS.md.

Paths not yet present are proposed targets to confirm during this epic's design. Task Context Files must name exact files when the detailed plan is written.

## Architectural delta

Final runtime becomes V2-only; keep V1 parsing solely for explicit migration and historical fixtures. Reconcile current Design/Guardrails and canonical/shipped guidance at the proven cutover.

Reference: `.savepoint/releases/v2/v2-Design.md`.

## Boundaries

**In scope:** the outcome and additions above.

**Out of scope:** No automatic publication, universal model-quality claims, benchmark program, hooks, or telemetry.

## Dependencies

- E49-objective-board
- E51-first-class-releases

## Quality gates

Three agent scenarios recorded; two trial projects; read-only/backup/recovery evidence; six platform artifacts and checksums; retired V1 live paths; documentation matches shipped behavior.

Implementation handoff requires `make build && make test` and named outcome evidence. Epic closeout requires a fresh independent V1 audit; this planning session is not that audit. Apply relevant current Guardrails, keeping STYLE advisory.

## Open decisions

Do not migrate this repository until core workflow, release semantics, and
recovery pass; task completion and publishing remain owner-controlled actions.
