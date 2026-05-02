---
id: E05-tasking-permissions/T002-update-router-md
status: planned
objective: "Add agent router-update guidance to router.md task-building state and update current state for E05"
depends_on: ["E05-tasking-permissions/T001-update-agents-md"]
---

# T002: Update Router.md

## Acceptance Criteria

- Router.md `task-building` state section includes guidance: "When starting work, set task `status: in_progress` and press `m` in TUI to update router task"
- Current state in router.md updated to point to E05/T001
- Router.md continues to describe `audit-pending` transition correctly (user confirms, router updates)

## Implementation Plan

- [ ] Edit `.savepoint/router.md` `task-building` state block to add agent router-update instruction
- [ ] Update router.md `## Current state` yaml block to point to E05 tasking-permissions epic
- [ ] Run `make build && make test` to verify

## Context Log

Files read:
- .savepoint/router.md
- .savepoint/releases/v1.1/epics/E05-tasking-permissions/E05-Detail.md

Estimated input tokens: 2500

Notes:
