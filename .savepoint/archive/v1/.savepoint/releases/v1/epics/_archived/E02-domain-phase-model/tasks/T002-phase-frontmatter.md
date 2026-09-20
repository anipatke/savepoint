---
id: E02-domain-phase-model/T002-phase-frontmatter
status: planned
objective: "Add phase field to TaskFrontmatter and update validation in task.ts"
depends_on: [E02-domain-phase-model/T001-phase-types]
---

# T002: Phase Frontmatter

## Acceptance Criteria

- `TaskFrontmatter` interface includes optional `phase?: TaskPhase`.
- `validateTaskFrontmatter` accepts and validates `phase` as optional string.
- Invalid phase values are rejected with clear error messages.
- `phase` is absent when `status` is not `in_progress`.
- All task validation tests pass.

## Implementation Plan

- [ ] Read existing `task.ts` and `test/domain/task.test.ts`.
- [ ] Add `phase?: TaskPhase` to `TaskFrontmatter` interface.
- [ ] Update `validateTaskFrontmatter` to parse optional `phase` field.
- [ ] Validate that phase is one of `build`, `test`, `audit` if present.
- [ ] Update `test/domain/task.test.ts` with phase validation cases.
- [ ] Run `npm test` and verify.
