---
type: epic-design
status: planned
---

# E46: Guide planning, execution, and independent checking

## Purpose

Four public skills carry the agreed workflow with explicit context and authority.

V1 delivery epic for V2 product Objective O006. The current V1 lifecycle governs this work; this is the single authoritative backlog record.

## What this epic adds

Idea/design/task/check skill contracts; matching scaffold copies; shared method; Task and Objective templates; Issue capture guidance; command/procedure migration instructions.

## Components and files

agent-skills/ Savepoint skills; agent-skills/references/audit-method.md; templates/project/agent-skills/; templates/project/AGENTS.md; internal/init/skill_validation_test.go; audit_contract_test.go; agent_skills_test.go.

Paths not yet present are proposed targets to confirm during this epic's design. Task Context Files must name exact files when the detailed plan is written.

## Architectural delta

Add V2 skills under distinct names and retain existing V1 instruction pairs until cutover. Preserve method rigor without register ceremony.

Reference: `.savepoint/releases/v2/v2-Design.md`.

## Boundaries

**In scope:** the outcome and additions above.

**Out of scope:** No automatic model routing, universal context byte cap, or mandatory extra planning documents.

## Dependencies

- E44-issues-objective-checks

## Quality gates

Canonical/shipped parity, role write boundaries, targeted reads/replan, conditional acceptance, one-Objective planning, and fresh-session Check are explicit and contract-tested.

Implementation handoff requires `make build && make test` and named outcome evidence. Epic closeout requires a fresh independent V1 audit; this planning session is not that audit. Apply relevant current Guardrails, keeping STYLE advisory.

## Open decisions

Use design sections 7–8; semantic quality is evaluated with real agents in E50, not claimed from text tests.

