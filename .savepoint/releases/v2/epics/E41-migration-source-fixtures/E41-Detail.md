---
type: epic-design
status: audited
---

# E41: Preserve migration source evidence

## Purpose

Freeze representative V1 projects and prove their raw records, scoped identities, and authored history remain available to migration.

V1 delivery epic for V2 product Objective O001. The current V1 lifecycle governs this work; this is the single authoritative backlog record.

## What this epic adds

Two complementary frozen fixtures cover simple active/completed Task dependencies and multi-release findings/defects with custom content.

## Components and files

internal/data/parser.go; internal/data/dependency.go; internal/data/audit_finding.go; internal/data/audit_run.go; internal/init/testdata/legacy/README.md; new internal/data/testdata/migration/ fixtures and focused characterization tests.

Paths not yet present are proposed targets to confirm during this epic's design. Task Context Files must name exact files when the detailed plan is written.

## Architectural delta

Adds migration input characterization only. Existing readers and writers keep their behavior. Future conversion uses raw source metadata plus these expected interpretations.

Reference: `.savepoint/releases/v2/v2-Design.md`.

## Boundaries

**In scope:** the outcome and additions above.

**Out of scope:** No converter, schema changes, public CLI, template changes, or live project mutation.

## Dependencies

None.

## Quality gates

Parser aliases, raw status fidelity, scoped dependency targets, optional artifact absence/presence, and fixture byte hashes are named evidence.

Implementation handoff requires `make build && make test` and named outcome evidence. Epic closeout requires a fresh independent V1 audit; this planning session is not that audit. Apply relevant current Guardrails, keeping STYLE advisory.

## Open decisions

No unresolved design decision. The two task files detail this epic; all later epics remain at design-outline level.

