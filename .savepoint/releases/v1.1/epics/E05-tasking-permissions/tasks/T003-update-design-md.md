---
id: E05-tasking-permissions/T003-update-design-md
status: planned
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

- [ ] Edit `.savepoint/Design.md` Section 4 to add new permission row or note
- [ ] Update `last_audited` in frontmatter
- [ ] Run `make build && make test` to verify

## Context Log

Files read:
- .savepoint/Design.md
- .savepoint/releases/v1.1/epics/E05-tasking-permissions/E05-Detail.md

Estimated input tokens: 2500

Notes:
