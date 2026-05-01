---
id: E08-board-workflow-cleanup/T001-acceptance-criteria-model
status: done
objective: Add task-document acceptance criteria parsing with compatibility for existing tasks that do not declare criteria.
depends_on: []
---

# T001: acceptance-criteria-model

## Acceptance Criteria

- Task documents with a `## Acceptance Criteria` section expose normalized criteria to callers.
- Task documents without acceptance criteria remain readable and carry an explicit compatibility signal instead of failing schema validation.
- Existing frontmatter validation remains unchanged, so status writes still preserve unknown frontmatter and task bodies.
- Tests cover present, missing, empty, and checkbox/bullet acceptance criteria formats.

## Implementation Plan

- [x] Add a focused Markdown section parser for task body sections, keeping frontmatter validation in `src/domain/task.ts` strict and unchanged.
- [x] Extend the task document model returned to callers with acceptance criteria data and a missing-criteria compatibility marker.
- [x] Update `test/domain/task.test.ts` to cover acceptance criteria extraction, missing-section compatibility, empty-section compatibility, and mixed bullet/checkbox criteria.
- [x] Update the relevant status-write tests to prove task body and criteria text survive status changes unchanged.
- [x] Run focused task-domain and status-write tests before handing off.
- [x] Update this task's context log with implementation files read before moving it to review.

## Context Log

- Files read: `.savepoint/router.md`; `AGENTS.md`; `.savepoint/releases/v1/epics/E08-board-workflow-cleanup/Design.md`; `src/domain/task.ts`; `test/domain/task.test.ts`; `test/tui/io/write-status.test.ts`.
- Estimated input tokens: ~7,500
- Notes: Implementation pass. Added `AcceptanceCriteria` interface, `parseAcceptanceCriteria`, and `parseTaskDocument` to `src/domain/task.ts`. Tests: 9 new in task.test.ts, 1 new in write-status.test.ts. 35/35 pass.
