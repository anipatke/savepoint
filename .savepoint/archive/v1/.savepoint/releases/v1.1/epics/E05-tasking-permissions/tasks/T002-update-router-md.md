---
id: E05-tasking-permissions/T002-update-router-md
status: done
objective: "Add agent router-update guidance to router.md task-building state and update current state for E05"
depends_on: ["E05-tasking-permissions/T001-update-agents-md"]
---

# T002: Update Router.md

## Acceptance Criteria

- Router.md `task-building` state section includes guidance: "When starting work, set task `status: in_progress` and press `m` in TUI to update router task"
- Current state in router.md updated to point to E05/T001
- Router.md continues to describe `audit-pending` transition correctly (user confirms, router updates)

## Implementation Plan

- [x] Edit `.savepoint/router.md` `task-building` state block to add agent router-update instruction
- [x] Update router.md `## Current state` yaml block to point to E05 tasking-permissions epic
- [x] Run `make build && make test` to verify

## Context Log

Files read:
- .savepoint/router.md
- .savepoint/releases/v1.1/epics/E05-tasking-permissions/E05-Detail.md
- .savepoint/releases/v1.1/epics/E05-tasking-permissions/tasks/T001-update-agents-md.md
- .savepoint/releases/v1.1/epics/E05-tasking-permissions/tasks/T002-update-router-md.md
- agent-skills/savepoint-build-task/SKILL.md

Files edited:
- .savepoint/router.md
- .savepoint/releases/v1.1/epics/E05-tasking-permissions/tasks/T002-update-router-md.md

Estimated input tokens: 5200

Notes:
- `make build && make test` could not run literally because `make` is unavailable in this environment.
- Equivalent gates passed: `go run ./internal/buildtool -version "" build` and `go test ./...`.
- Router current state already pointed to E05/T002. I left it on the active task instead of moving it backward to completed dependency T001.
