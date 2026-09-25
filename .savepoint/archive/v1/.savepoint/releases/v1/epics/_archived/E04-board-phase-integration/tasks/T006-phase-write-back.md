---
id: E04-board-phase-integration/T006-phase-write-back
status: planned
objective: "Write phase field alongside status in task frontmatter via write-status.ts"
depends_on: [E04-board-phase-integration/T005-phase-gates]
---

# T006: Phase Write-Back

## Acceptance Criteria

- `writeTaskStatus` accepts optional `newPhase?: TaskPhase`.
- When writing back, `phase` field is included in frontmatter if task is in_progress.
- When transitioning to planned or done, `phase` is removed from frontmatter.
- Mtime conflict detection still works.
- Write-status tests pass.

## Implementation Plan

- [ ] Read `write-status.ts` and `test/tui/io/write-status.test.ts`.
- [ ] Add `newPhase?: TaskPhase` parameter to `writeTaskStatus`.
- [ ] Update frontmatter write logic to include/exclude phase based on status.
- [ ] Update write-status tests.
- [ ] Run `npm test`.
