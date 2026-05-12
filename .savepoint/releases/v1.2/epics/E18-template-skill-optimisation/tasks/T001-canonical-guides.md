---
id: E18-template-skill-optimisation/T001-canonical-guides
status: done
objective: Make skills the single source of truth for workflow guidance and remove duplicate prompt-style instructions from AGENTS docs
depends_on: []
complexity_tier: medium
complexity_reason: "Consistency pass across live and scaffolded docs and skills; many files but no production code changes"
---

# T001: Canonical Guides

## Problem

The live and scaffolded workflow docs repeat the same phase instructions in different places, which makes the contract harder to keep consistent. The team needs one canonical explanation of how Savepoint phases work and where each responsibility lives.

## Context Files

- `AGENTS.md`
- `templates/project/AGENTS.md`
- `templates/project/.savepoint/router.md`
- `agent-skills/savepoint-draft-prd/SKILL.md`
- `agent-skills/savepoint-system-design/SKILL.md`
- `agent-skills/savepoint-create-plan/SKILL.md`
- `agent-skills/savepoint-create-task/SKILL.md`
- `agent-skills/savepoint-build-task/SKILL.md`
- `agent-skills/savepoint-audit/SKILL.md`
- `templates/project/agent-skills/savepoint-draft-prd/SKILL.md`
- `templates/project/agent-skills/savepoint-system-design/SKILL.md`
- `templates/project/agent-skills/savepoint-create-plan/SKILL.md`
- `templates/project/agent-skills/savepoint-create-task/SKILL.md`
- `templates/project/agent-skills/savepoint-build-task/SKILL.md`
- `templates/project/agent-skills/savepoint-audit/SKILL.md`

## Acceptance Criteria

- [x] The live and scaffolded AGENTS docs clearly state that skills are the canonical phase workflow source
- [x] The live and scaffolded AGENTS docs no longer conflict on state, status, or stage terminology
- [x] The root skills and scaffolded skill copies use the same simplified phase instructions
- [x] Prompt-style duplication is removed from the workflow guidance text
- [x] `make build && make test` passes after the doc updates

## Implementation Plan

- [x] Rewrite the live and scaffolded AGENTS guidance to remove redundant phase narration
- [x] Align the root skill files and scaffolded skill copies so they describe the same workflow contract
- [x] Normalize terminology so `state`, `status`, and `stage` are used consistently
- [x] Run `make build && make test`

## Context Log

- Files read: `.savepoint/router.md`, `.savepoint/releases/v1.2/epics/E18-template-skill-optimisation/E18-Detail.md`, `.savepoint/releases/v1.2/epics/E18-template-skill-optimisation/tasks/T001-canonical-guides.md`, all files listed in `## Context Files`, plus `internal/init/template_freshness_test.go` to resolve a failing freshness assertion.
- Files edited: `AGENTS.md`, `templates/project/AGENTS.md`, `templates/project/.savepoint/router.md`, root `agent-skills/savepoint-*/SKILL.md`, scaffolded `templates/project/agent-skills/savepoint-*/SKILL.md`.
- Estimated input tokens: ~24k
- Notes: Verified mirrored root/scaffold skill copies with `Compare-Object`. `make build && make test` passed.
