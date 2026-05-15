---
id: E20-clean-up-lifecycle/T005-lifecycle-docs-and-regression
status: done
objective: Update lifecycle documentation, template guidance, and regression guards to match the shared contract.
depends_on: [T004-board-lifecycle-transitions]
complexity_tier: medium
complexity_reason: Documentation and freshness tests span live and scaffolded workflow surfaces.
---

# T005: Lifecycle Docs and Regression

## Problem

The code can be correct while agents still receive slightly different lifecycle instructions from AGENTS, skills, templates, and Design.md. That is how workflow drift reappears.

## Context Files

- `AGENTS.md`
- `templates/project/AGENTS.md`
- `agent-skills/savepoint-build-task/SKILL.md`
- `templates/project/agent-skills/savepoint-build-task/SKILL.md`
- `.savepoint/Design.md`
- `internal/init/template_freshness_test.go`
- `internal/data/parser_test.go`
- `internal/doctor/checks_test.go`
- `internal/board/transitions_test.go`

## Acceptance Criteria

- [x] Live and scaffolded agent guidance describe the same task lifecycle contract.
- [x] Design.md records `internal/data` as the single owner of task lifecycle rules.
- [x] Template freshness or documentation tests protect against reintroducing `phase` as canonical task guidance.
- [x] Regression tests cover the defect pattern that motivated the epic.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Update live and scaffolded workflow guidance only where it clarifies the shared lifecycle contract.
- [x] Update Design.md with the lifecycle ownership decision and compatibility boundary.
- [x] Add or adjust freshness tests for live/scaffolded lifecycle language.
- [x] Run focused package tests, then `make build` and `make test`.

## Context Log

- Read: `.savepoint/router.md`, `.savepoint/releases/v1.2/epics/E20-clean-up-lifecycle/E20-Detail.md`, this task, and all listed context files.
- Edited: `AGENTS.md`, `templates/project/AGENTS.md`, `agent-skills/savepoint-build-task/SKILL.md`, `templates/project/agent-skills/savepoint-build-task/SKILL.md`, `.savepoint/Design.md`, `internal/init/template_freshness_test.go`, `internal/data/parser_test.go`, `internal/doctor/checks_test.go`, and this task file.
- Router priority: router already pointed at this task; no interactive TUI was available to press `p`.
- Verification: `go test ./internal/init ./internal/data ./internal/doctor ./internal/board` passed.
- Verification: `make build` passed.
- Verification: `make test` passed.
