---
id: E16-pre-prod-refinement/T001-merge-existing-agents-md
status: done
objective: Insert Savepoint init instructions into existing agent guide files instead of rejecting or overwriting them
depends_on: []
---

# T001: Merge Existing Agent Guide During Init

## Context Files

- `internal/init/validate.go` — target validation and scaffold conflict checks
- `internal/init/scaffold.go` — template walk and rendered file writes
- `internal/init/write.go` — atomic write helper used by scaffold writes
- `internal/init/validate_test.go` — validation coverage for conflicting files
- `internal/init/scaffold_test.go` — scaffold write and interpolation coverage
- `internal/init/integration_test.go` — end-to-end init pipeline coverage
- `main.go` — init runner wiring for validation, scaffold, and magic prompt rendering

## Acceptance Criteria

- [ ] `ValidateTarget(dir, false)` succeeds when the target contains an existing root `AGENTS.md`
- [ ] `savepoint init` preserves existing agent guide content and inserts rendered Savepoint instructions into a managed block
- [ ] Re-running scaffold/init refreshes the managed block without duplicating Savepoint instructions
- [ ] `--force` updates the managed block without deleting user-owned content outside that block
- [ ] Existing casing variants such as `Agents.MD` are detected and reused instead of creating a second agent guide
- [ ] Existing `agent-skills` conflict behavior is unchanged
- [ ] Unit and integration tests cover insert, idempotency, casing, validation, and force behavior
- [ ] `make build && make test` passes

## Implementation Plan

- [x] Remove `AGENTS.md` from the validation conflict list while keeping `agent-skills` protected
- [x] Add a focused `internal/init/agents.go` helper for detecting root agent guide filename variants
- [x] Add managed Savepoint block markers and insertion/replacement logic
- [x] Special-case `AGENTS.md` template writes in `Scaffold` to call the merge helper
- [x] Preserve normal `AtomicWrite` behavior for all non-agent-guide template files
- [x] Update validation tests for existing `AGENTS.md`
- [x] Add scaffold tests for preserving content, inserting the block, idempotent replacement, `--force`, and `Agents.MD`
- [x] Add integration coverage for init in a project with an existing agent guide
- [x] Run `make build && make test`
