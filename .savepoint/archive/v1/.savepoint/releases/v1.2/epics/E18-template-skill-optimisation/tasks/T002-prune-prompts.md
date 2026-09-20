---
id: E18-template-skill-optimisation/T002-prune-prompts
status: done
objective: Remove redundant phase prompt templates and keep only the init bootstrap prompt
depends_on: [E18-template-skill-optimisation/T001-canonical-guides]
complexity_tier: low
complexity_reason: "Deletes obsolete prompt templates and updates tests to match the reduced prompt surface"
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

- [x] Only `magic-prompt.prompt.md` remains in `templates/prompts`
- [x] The init path still renders the bootstrap prompt successfully
- [x] The scaffold and integration tests no longer expect deleted prompt assets
- [x] No runtime code path depends on the removed prompt templates
- [x] `make build && make test` passes after the prompt cleanup

## Implementation Plan

- [x] Delete the obsolete prompt template files
- [x] Keep the magic prompt wired into init and verify it still renders correctly
- [x] Update prompt and integration tests to reflect the reduced prompt surface
- [x] Run `make build && make test`

## Context Log

- Files read: `.savepoint/router.md`; `.savepoint/releases/v1.2/epics/E18-template-skill-optimisation/E18-Detail.md`; this task file; `templates/prompts/magic-prompt.prompt.md`; `templates/prompts/prd.prompt.md`; `templates/prompts/design.prompt.md`; `templates/prompts/epic-design.prompt.md`; `templates/prompts/task-breakdown.prompt.md`; `templates/prompts/task-building.prompt.md`; `templates/prompts/task-planning.prompt.md`; `templates/prompts/audit-reconciliation.prompt.md`; `main.go`; `internal/init/prompt.go`; `internal/init/prompt_test.go`; `internal/init/integration_test.go`.
- Files edited: deleted obsolete prompt templates except `magic-prompt.prompt.md`; updated `internal/init/prompt_test.go`; updated this task file.
- Estimated input tokens: ~14,000.
- Notes: `templates/prompts` now contains only `magic-prompt.prompt.md`. Targeted search found no runtime references to removed prompt filenames. `go test ./internal/init` passed. `make build` passed. `make test` passed.
