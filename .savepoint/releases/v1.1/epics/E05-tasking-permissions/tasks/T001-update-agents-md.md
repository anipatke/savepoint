---
id: E05-tasking-permissions/T001-update-agents-md
status: planned
objective: "Update AGENTS.md to remove agent 'set status: done' permission and add 'set status: in_progress' guidance with router update instruction"
depends_on: []
---

# T001: Update AGENTS.md

## Acceptance Criteria

- Task Completion Protocol no longer says "Set the task frontmatter to `status: done`"
- Protocol says "When starting implementation, set `status: in_progress`"
- Protocol says "After setting `in_progress`, press `m` in the TUI to update the router"
- The words "Set the task frontmatter to `status: done`" are removed or replaced
- Audit Handoff Rule still says the agent updates router.md to `state: audit-pending` when all tasks are done (this is router state, not task state)

## Implementation Plan

- [ ] In AGENTS.md `## Task Completion Protocol`: replace "Set the task frontmatter to `status: done`" with appropriate guidance
- [ ] In AGENTS.md `## Task Completion Protocol` step 7: keep router update for `audit-pending` (that's router state, not task status)
- [ ] In AGENTS.md add note under `## Task Status Canon` about agent vs user capabilities
- [ ] Run `make build && make test` to verify

## Context Log

Files read:
- AGENTS.md
- .savepoint/router.md
- .savepoint/releases/v1.1/epics/E05-tasking-permissions/E05-Detail.md

Estimated input tokens: 3000

Notes:
