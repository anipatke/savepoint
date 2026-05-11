---
type: epic-design
status: audited
---

# E18: Template and Skill Optimisation

## Purpose

Simplify the Savepoint scaffold and live guidance so the workflow is explained once, by skills and AGENTS, instead of being repeated across multiple prompt files. This epic keeps the current behaviour but removes redundant surfaces, tightens terminology, and makes the template suite easier for agents to read and maintain.

## What this epic adds

- A canonical workflow contract in the skill files and AGENTS guidance.
- A reduced prompt surface that keeps only the init-time magic prompt.
- Consistent wording for state, status, and stage terminology across scaffolded docs and tests.
- Release planning metadata that records the new epic alongside the existing v1.2 defect work.

## Components and files

| Module | Purpose |
|--------|---------|
| `AGENTS.md` | Live agent guide that should reflect the simplified workflow contract |
| `templates/project/AGENTS.md` | Scaffolded agent guide that must mirror the live guidance |
| `agent-skills/savepoint-*/SKILL.md` | Canonical skill instructions for each Savepoint phase |
| `templates/project/agent-skills/savepoint-*/SKILL.md` | Scaffolded copies of the canonical skill instructions |
| `templates/prompts/` | Prompt assets, which should be reduced to the init bootstrap prompt |
| `internal/init/*` | Prompt rendering, integration, upgrade, and freshness tests that protect the template surface |
| `.savepoint/releases/v1.2/v1.2-PRD.md` | Release scope and epic list for v1.2 |

## Architectural delta

Before this epic, the workflow guidance is split across skills, prompts, and scaffold docs with overlapping instructions and inconsistent terminology. After this epic, skills are the single source of truth for phase behaviour, prompts are limited to the bootstrap prompt, and the release documentation explicitly tracks the slimmer template surface.

## Boundaries

**In scope:**
- Simplifying and aligning live and scaffolded AGENTS guidance
- Keeping skill files authoritative and mirrored in the scaffold
- Removing redundant phase prompt files while preserving the magic prompt
- Updating tests and release metadata to protect the simplified template surface

**Out of scope:**
- Changing board behaviour, router runtime behaviour, or task lifecycle rules
- Adding new prompt-based workflows
- Reworking the defect epic or its release data model

## Quality gates

- The live and scaffolded guidance files agree on the canonical workflow contract
- Only the init prompt remains in the prompt template set
- The template freshness and integration tests cover the reduced surface
- `make build && make test` passes after the cleanup

## Open decisions

None. The epic standard is to make skills canonical and remove redundant prompt surfaces.
