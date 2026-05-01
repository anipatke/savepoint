---
id: E01-archive-and-reset/T004-update-router
status: planned
objective: "Rewrite router.md with 3-state model and point at E01"
depends_on: [E01-archive-and-reset/T003-create-epic-stubs]
---

# T004: Update Router

## Acceptance Criteria

- `router.md` defines 3 states: `planning`, `building`, `reviewing`.
- `router.md` points at epic `E01-archive-and-reset` as the active epic.
- `router.md` contains a `## Current state` yaml block with valid frontmatter.
- No references to old 6-state model remain.

## Implementation Plan

- [ ] Read existing router.md.
- [ ] Rewrite state machine section with 3 states.
- [ ] Update Current state yaml block: state = planning, release = v1, epic = E01-archive-and-reset.
- [ ] Verify yaml is parseable.
