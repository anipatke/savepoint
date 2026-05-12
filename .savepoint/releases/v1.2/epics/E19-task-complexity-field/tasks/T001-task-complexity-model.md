---
id: E19-task-complexity-field/T001-task-complexity-model
status: done
objective: Add complexity fields to task parsing, validation, and persistence
depends_on: []
complexity_tier: medium
complexity_reason: "Touches task model, parser, lifecycle validation, doctor checks, and write tests"
---

# T001: Task Complexity Model

## Problem

Task planning needs a typed complexity signal before the board and skill updates can depend on it. The data boundary has to accept the new fields, preserve them through writes, and reject malformed values consistently.

## Context Files

- `internal/data/task.go`
- `internal/data/parser.go`
- `internal/data/write.go`
- `internal/data/parser_test.go`
- `internal/data/write_test.go`
- `internal/data/lifecycle.go`
- `internal/doctor/checks.go`
- `internal/doctor/checks_test.go`

## Acceptance Criteria

- [x] `data.Task` can hold a complexity tier and a short complexity reason
- [x] Task parsing accepts valid complexity frontmatter values for existing and new task files
- [x] Task parsing preserves backward compatibility when complexity is absent
- [x] Invalid complexity tiers or malformed reasons are rejected by the validation path
- [x] Task write helpers round-trip complexity metadata without dropping it
- [x] Tests cover valid parsing, missing complexity compatibility, invalid tier handling, and write round-trips
- [x] `make build && make test` passes for the updated data and validation paths

## Implementation Plan

- [x] Add complexity fields to the task model and task frontmatter struct
- [x] Extend parse and write helpers so complexity survives round-trips
- [x] Add validation for allowed tiers and the short reason rule
- [x] Update doctor checks so malformed complexity files are reported
- [x] Add focused parser, write, and validation tests
- [x] Run `make build && make test`

## Context Log

- Files read: task.go, parser.go, write.go, parser_test.go, write_test.go, lifecycle.go, checks.go, checks_test.go
- Estimated input tokens: ~6000
- Notes: Pre-existing failure in TestBundledSavepointSkillsHaveDiscoveryFrontmatter (savepoint-audit SKILL.md) unrelated to this task — confirmed present on base branch.
