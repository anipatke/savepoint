---
id: E18-template-skill-optimisation/T002-prune-prompts
status: planned
objective: Remove redundant phase prompt templates and keep only the init bootstrap prompt
depends_on: [E18-template-skill-optimisation/T001-canonical-guides]
---

# T002: Prompt Pruning

## Problem

The prompt directory repeats the same workflow ideas already owned by skills and AGENTS, which creates maintenance overhead without adding new meaning. The init path only needs one prompt, so the rest of the prompt surface is duplication.

## Context Files

- `templates/prompts/magic-prompt.prompt.md`
- `templates/prompts/prd.prompt.md`
- `templates/prompts/design.prompt.md`
- `templates/prompts/epic-design.prompt.md`
- `templates/prompts/task-breakdown.prompt.md`
- `templates/prompts/task-building.prompt.md`
- `templates/prompts/task-planning.prompt.md`
- `templates/prompts/audit-reconciliation.prompt.md`
- `main.go`
- `internal/init/prompt.go`
- `internal/init/prompt_test.go`
- `internal/init/integration_test.go`

## Acceptance Criteria

- [ ] Only `magic-prompt.prompt.md` remains in `templates/prompts`
- [ ] The init path still renders the bootstrap prompt successfully
- [ ] The scaffold and integration tests no longer expect deleted prompt assets
- [ ] No runtime code path depends on the removed prompt templates
- [ ] `make build && make test` passes after the prompt cleanup

## Implementation Plan

- [ ] Delete the obsolete prompt template files
- [ ] Keep the magic prompt wired into init and verify it still renders correctly
- [ ] Update prompt and integration tests to reflect the reduced prompt surface
- [ ] Run `make build && make test`

## Context Log

- Files read:
- Estimated input tokens:
- Notes:
