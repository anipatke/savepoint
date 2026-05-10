---
id: E18-template-skill-optimisation/T003-template-surface-tests
status: planned
objective: Add regression coverage that locks in the simplified template surface and canonical workflow wording
depends_on: [E18-template-skill-optimisation/T001-canonical-guides, E18-template-skill-optimisation/T002-prune-prompts]
---

# T003: Template Surface Tests

## Problem

The cleanup only holds if the test suite fails when stale template or prompt wording returns. The repo needs focused assertions that protect the slimmer surface and the canonical terminology.

## Context Files

- `internal/init/template_freshness_test.go`
- `internal/init/integration_test.go`
- `internal/init/prompt_test.go`
- `internal/init/upgrade_test.go`
- `main.go`

## Acceptance Criteria

- [ ] Tests verify that the simplified AGENTS and skill wording stays in sync
- [ ] Tests verify that the prompt set only includes the bootstrap prompt
- [ ] Tests verify that init and upgrade paths still behave correctly after prompt pruning
- [ ] Tests catch stale workflow terms such as prompt-based phase instructions or inconsistent status/stage wording
- [ ] `make build && make test` passes with the new regression coverage

## Implementation Plan

- [ ] Update the freshness and integration tests to assert the canonical wording and reduced prompt surface
- [ ] Strengthen prompt tests to cover the remaining bootstrap-only prompt path
- [ ] Update any helper fixtures that still expect the removed prompt files
- [ ] Run `make build && make test`

## Context Log

- Files read:
- Estimated input tokens:
- Notes:
