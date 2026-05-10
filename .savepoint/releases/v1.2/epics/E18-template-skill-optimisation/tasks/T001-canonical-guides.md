---
id: E18-template-skill-optimisation/T001-canonical-guides
status: planned
objective: Make skills the single source of truth for workflow guidance and remove duplicate prompt-style instructions from AGENTS docs
depends_on: []
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

- [ ] The live and scaffolded AGENTS docs clearly state that skills are the canonical phase workflow source
- [ ] The live and scaffolded AGENTS docs no longer conflict on state, status, or stage terminology
- [ ] The root skills and scaffolded skill copies use the same simplified phase instructions
- [ ] Prompt-style duplication is removed from the workflow guidance text
- [ ] `make build && make test` passes after the doc updates

## Implementation Plan

- [ ] Rewrite the live and scaffolded AGENTS guidance to remove redundant phase narration
- [ ] Align the root skill files and scaffolded skill copies so they describe the same workflow contract
- [ ] Normalize terminology so `state`, `status`, and `stage` are used consistently
- [ ] Run `make build && make test`

## Context Log

- Files read:
- Estimated input tokens:
- Notes:
