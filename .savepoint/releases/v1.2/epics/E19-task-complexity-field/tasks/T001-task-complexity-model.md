---
id: E19-task-complexity-field/T001-task-complexity-model
status: planned
objective: Add complexity fields to task parsing, validation, and persistence
depends_on: []
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

- [ ] `data.Task` can hold a complexity tier and a short complexity reason
- [ ] Task parsing accepts valid complexity frontmatter values for existing and new task files
- [ ] Task parsing preserves backward compatibility when complexity is absent
- [ ] Invalid complexity tiers or malformed reasons are rejected by the validation path
- [ ] Task write helpers round-trip complexity metadata without dropping it
- [ ] Tests cover valid parsing, missing complexity compatibility, invalid tier handling, and write round-trips
- [ ] `make build && make test` passes for the updated data and validation paths

## Implementation Plan

- [ ] Add complexity fields to the task model and task frontmatter struct
- [ ] Extend parse and write helpers so complexity survives round-trips
- [ ] Add validation for allowed tiers and the short reason rule
- [ ] Update doctor checks so malformed complexity files are reported
- [ ] Add focused parser, write, and validation tests
- [ ] Run `make build && make test`

## Context Log

- Files read:
- Estimated input tokens:
- Notes:
