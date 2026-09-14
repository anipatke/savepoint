---
type: epic-design
status: planned
---

# E47: Start and refresh coherent V2 projects

## Purpose

Fresh and existing codebases receive the V2 workflow while authored files and customized skills remain protected.

V1 delivery epic for V2 product Objective O007. The current V1 lifecycle governs this work; this is the single authoritative backlog record.

## What this epic adds

V2 default scaffold; rough Idea input/planner handoff; existing-codebase reconstruction guidance; version-aware asset upgrades; legacy instruction retirement.

## Components and files

internal/init/scaffold.go; upgrade.go; manifest.go; agents.go; templates/project/.savepoint/; templates/prompts/magic-prompt.prompt.md; cmd/init.go; main.go; lifecycle/template tests.

Paths not yet present are proposed targets to confirm during this epic's design. Task Context Files must name exact files when the detailed plan is written.

## Architectural delta

Switch scaffold/schema as a coordinated release change. Schema migration remains explicit; asset install cannot silently upgrade project state.

Reference: `.savepoint/releases/v2/v2-Design.md`.

## Boundaries

**In scope:** the outcome and additions above.

**Out of scope:** No eager reconstruction by a hidden model service, silent replacement of policy, or migration on npm install.

## Dependencies

- E45-safe-migration
- E46-agent-workflow-assets

## Quality gates

Fresh init, existing codebase, edited skills, marked/unmarked guides, missing optional procedures, repeat upgrade and no-write preview are verified.

Implementation handoff requires `make build && make test` and named outcome evidence. Epic closeout requires a fresh independent V1 audit; this planning session is not that audit. Apply relevant current Guardrails, keeping STYLE advisory.

## Open decisions

Ensure V2 data/migration and workflow contracts are present before changing default init. Live V1 canonical build guidance stays until E50.

