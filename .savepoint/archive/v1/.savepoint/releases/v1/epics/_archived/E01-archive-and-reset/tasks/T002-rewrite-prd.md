---
id: E01-archive-and-reset/T002-rewrite-prd
status: planned
objective: "Rewrite releases/v1/PRD.md with simplification scope"
depends_on: [E01-archive-and-reset/T001-archive-epics]
---

# T002: Rewrite PRD

## Acceptance Criteria

- `releases/v1/PRD.md` defines the simplified savepoint scope.
- Scope includes: board-only CLI, phase model (build/test/audit), no audit pipeline, no init command, no doctor command.
- PRD frontmatter is valid: `version`, `name`, `status`.
- No references to the full 11-epic v1 scope remain in the active PRD.

## Implementation Plan

- [ ] Read existing PRD.md for reference.
- [ ] Write new PRD.md with simplified scope: board command, phase model, minimal CLI.
- [ ] Update frontmatter: status = in_progress, version = 1.
- [ ] List 5 new epics in scope table.
