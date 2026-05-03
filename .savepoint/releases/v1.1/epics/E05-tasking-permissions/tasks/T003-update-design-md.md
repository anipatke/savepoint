---
id: E05-tasking-permissions/T003-update-design-md
status: done
objective: "Revise Design.md Section 4 to document agent vs user task status capabilities"
depends_on: ["E05-tasking-permissions/T001-update-agents-md"]
---

# T003: Update Design.md

## Acceptance Criteria

- Design.md Section 4 (`## 4. Status model & gates`) documents:
  - Agent can only set `status: in_progress`
  - Only user can set `status: done` or retreat
  - Router update is explicit via TUI `m` key, not automatic on navigation
- `last_audited` field updated (or not — this is documenting current design, not auditing)
- No other sections modified

## Implementation Plan

- [x] Edit `.savepoint/Design.md` Section 4 to add new permission row or note
- [x] Update `last_audited` in frontmatter
- [x] Run `make build && make test` to verify

## Context Log

Files read:
- .savepoint/Design.md
- .savepoint/releases/v1.1/epics/E05-tasking-permissions/E05-Detail.md
- .savepoint/releases/v1.1/epics/E05-tasking-permissions/tasks/T001-update-agents-md.md
- .savepoint/releases/v1.1/epics/E05-tasking-permissions/tasks/T002-update-router-md.md
- .savepoint/releases/v1.1/epics/E05-tasking-permissions/tasks/T003-update-design-md.md
- .savepoint/router.md
- agent-skills/savepoint-build-task/SKILL.md
- Makefile

Files edited:
- .savepoint/Design.md
- .savepoint/router.md
- .savepoint/releases/v1.1/epics/E05-tasking-permissions/tasks/T003-update-design-md.md

Estimated input tokens: 8500

Notes:
- `make build && make test` could not run literally because `make` is unavailable in this environment.
- Equivalent gates passed: `go run ./internal/buildtool -version "" build` and `go test ./...`.
- Task remains `status: in_progress`; only the user may set `status: done`.
