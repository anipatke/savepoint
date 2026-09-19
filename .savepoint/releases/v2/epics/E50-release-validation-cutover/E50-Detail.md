---
type: epic-design
status: planned
---

# E50: Release evaluated V2 and retire transitional runtime

## Purpose

V2 is validated on realistic projects, packaged for supported platforms, and ready for explicit maintainer migration.

V1 delivery epic for V2 product Objective O010. The current V1 lifecycle governs this work; this is the single authoritative backlog record.

## What this epic adds

Agent evaluation scenarios; end-to-end migration trials; packaging/docs/help; legacy live-reader removal; archive/reference access; own-project cutover runbook. E51 now supplies the Release proof this epic must consume: migration is previewed and recoverable on temporary copies, and `internal/data.ResolveReleaseCutover` composes the canonical per-Release completion decisions.

## Components and files

internal/data/ transitional dispatch and Release cutover composition; internal/board/; internal/doctor/; internal/init/; internal/migrate/; internal/buildtool/main.go; bin/savepoint.js; package.json; Makefile; .github/workflows/ci.yml; .github/workflows/publish.yml; README.md; .savepoint/Design.md; .savepoint/Guardrails.md; AGENTS.md.

Paths not yet present are proposed targets to confirm during this epic's design. Task Context Files must name exact files when the detailed plan is written.

## Architectural delta

Final runtime becomes V2-only; keep V1 parsing solely for explicit migration and historical fixtures. Before the cutover write, require an appliable, recoverable migration and a clean V2 load, then consume `data.ResolveReleaseCutover` for every declared Release. That function only delegates to `ResolveReleaseCompletion`; E50 must not add a parallel readiness rule. Reconcile current Design/Guardrails and canonical/shipped guidance at the proven cutover.

Reference: `.savepoint/releases/v2/v2-Design.md`.

## Boundaries

**In scope:** the outcome and additions above.

**Out of scope:** No automatic publication, universal model-quality claims, benchmark program, hooks, or telemetry.

## Dependencies

- E49-objective-board
- E51-first-class-releases

## Quality gates

Three agent scenarios recorded; two trial projects; repository-copy Release accountability; read-only/backup/recovery evidence including every publish boundary and edit-conflict retry; canonical Release cutover decision consumed without duplication; six platform artifacts and checksums; retired V1 live paths; documentation matches shipped behavior.

Implementation handoff requires `make build && make test` and named outcome evidence. Epic closeout requires a fresh independent V1 audit; this planning session is not that audit. Apply relevant current Guardrails, keeping STYLE advisory.

## Open decisions

Do not migrate this repository until core workflow, Release semantics,
migration recovery, and the independent E51 audit pass. Refuse cutover while
a migration plan is ambiguous, a pending operation needs recovery, a V2
Release is technically unclear, has an unexcepted material Issue, or lacks
acceptance of the exact current Release Check. Task completion and publishing
remain owner-controlled actions; this epic does not publish or deploy anything.
