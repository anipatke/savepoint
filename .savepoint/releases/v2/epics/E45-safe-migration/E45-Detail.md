---
type: epic-design
status: planned
---

# E45: Convert active work without losing history

## Purpose

A previewed conversion preserves active work, historical sources, scoped references, and recoverable backups.

V1 delivery epic for V2 product Objective O005. The current V1 lifecycle governs this work; this is the single authoritative backlog record.

## What this epic adds

Read-only plan; raw record mapping; ambiguity handling; archive/manifest; interruption recovery; CLI migrate dispatch; version activation.

## Components and files

new internal/migrate/; cmd/migrate.go; main.go; internal/data/ legacy readers; internal/init/write.go and ownership helpers; E41 frozen fixtures.

Paths not yet present are proposed targets to confirm during this epic's design. Task Context Files must name exact files when the detailed plan is written.

## Architectural delta

Migration is separate from upgrade-assets and uses a recorded multi-file recovery protocol, not a claim of atomic multi-file replacement.

Reference: `.savepoint/releases/v2/v2-Design.md`.

## Boundaries

**In scope:** the outcome and additions above.

**Out of scope:** No live maintainer migration before full workflow readiness, permanent dual-mode runtime, or implicit semantic rewriting of authored histories.

## Dependencies

- E44-issues-objective-checks

## Quality gates

Dry run changes nothing; source changes invalidate preview; backups precede replace; injected failures/retry/user edits are safe; second run is unchanged; Windows and Unix recovery verified.

Implementation handoff requires `make build && make test` and named outcome evidence. Epic closeout requires a fresh independent V1 audit; this planning session is not that audit. Apply relevant current Guardrails, keeping STYLE advisory.

## Open decisions

First resolve Windows replacement/recovery through a bounded experiment with a decision and evidence; create a spike Task if the approach remains uncertain.

