---
id: E18-template-skill-optimisation/T003-template-surface-tests
status: done
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

- [x] Tests verify that the simplified AGENTS and skill wording stays in sync
- [x] Tests verify that the prompt set only includes the bootstrap prompt
- [x] Tests verify that init and upgrade paths still behave correctly after prompt pruning
- [x] Tests catch stale workflow terms such as prompt-based phase instructions or inconsistent status/stage wording
- [x] `make build && make test` passes with the new regression coverage

## Implementation Plan

- [x] Update the freshness and integration tests to assert the canonical wording and reduced prompt surface
- [x] Strengthen prompt tests to cover the remaining bootstrap-only prompt path
- [x] Update any helper fixtures that still expect the removed prompt files
- [x] Run `make build && make test`

## Context Log

- Files read: `.savepoint/router.md`; `.savepoint/releases/v1.2/epics/E18-template-skill-optimisation/E18-Detail.md`; `agent-skills/savepoint-build-task/SKILL.md`; `internal/init/template_freshness_test.go`; `internal/init/integration_test.go`; `internal/init/prompt_test.go`; `internal/init/upgrade_test.go`; `main.go`
- Files edited: `internal/init/template_freshness_test.go`; `internal/init/integration_test.go`; `internal/init/prompt_test.go`; `internal/init/upgrade_test.go`; `.savepoint/releases/v1.2/epics/E18-template-skill-optimisation/tasks/T003-template-surface-tests.md`
- Estimated input tokens: ~22k
- Quality gates: `go test ./internal/init` passed; `make build` passed; `make test` passed
- Notes: TUI router priority key `p` could not be pressed from the non-interactive shell. No drift notes required.
