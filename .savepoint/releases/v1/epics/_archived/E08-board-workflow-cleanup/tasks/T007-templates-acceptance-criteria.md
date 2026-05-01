---
id: E08-board-workflow-cleanup/T007-templates-acceptance-criteria
status: done
objective: Update workflow templates and prompts so new task files require acceptance criteria before implementation planning.
depends_on:
  - E08-board-workflow-cleanup/T001-acceptance-criteria-model
---

# T007: templates-acceptance-criteria

## Acceptance Criteria

- Live workflow docs and scaffold templates describe acceptance criteria as required task content before `## Implementation Plan`.
- Task-breakdown and task-planning prompts require acceptance criteria and distinguish observable outcomes from implementation checklist steps.
- Task-building prompts tell agents to satisfy acceptance criteria by completing the implementation checklist.
- Router and agent templates preserve token-discipline rules and include the Ink guide conditional read where TUI implementation requires it.
- Template tests prove generated task examples contain `## Acceptance Criteria` before `## Implementation Plan`.
- Historical task files are not bulk edited; compatibility behavior from T001 remains the migration path.

## Implementation Plan

- [x] Update live `.savepoint/router.md` workflow instructions only as needed for acceptance criteria and task detail navigation wording, without advancing router state while E07 remains active.
- [x] Update `templates/project/AGENTS.md` and `templates/project/.savepoint/router.md` so generated projects inherit acceptance-criteria and Ink-guide conditional-read guidance.
- [x] Update `templates/prompts/task-breakdown.prompt.md`, `templates/prompts/task-planning.prompt.md`, and `templates/prompts/task-building.prompt.md` to require acceptance criteria and clarify the implementation-plan relationship.
- [x] Update template tests for prompt shape, generated task examples, and router/agent workflow wording.
- [x] Confirm no historical task-file migration is required for tests or templates beyond fixtures that intentionally represent new tasks.
- [x] Run focused template tests before handing off.
- [x] Update this task's context log with implementation files read before moving it to review.

## Context Log

- Files read: `.savepoint/router.md`; `templates/project/AGENTS.md`; `templates/project/.savepoint/router.md`; `templates/prompts/task-breakdown.prompt.md`; `templates/prompts/task-planning.prompt.md`; `templates/prompts/task-building.prompt.md`; `test/templates/prompt-templates.test.ts`; `test/templates/project-templates.test.ts`; `test/templates/router-template.test.ts`; `test/templates/render-integrity.test.ts`; `test/templates/template-registry.test.ts`; `.savepoint/releases/v1/epics/E08-board-workflow-cleanup/Design.md`; `.savepoint/releases/v1/epics/E08-board-workflow-cleanup/tasks/T007-templates-acceptance-criteria.md`; `.savepoint/releases/v1/epics/E08-board-workflow-cleanup/tasks/T001-acceptance-criteria-model.md`; `src/domain/task.ts`; `.savepoint/metrics/context-bench.md`; `templates/prompts/magic-prompt.prompt.md`; `templates/prompts/audit-reconciliation.prompt.md`; `templates/prompts/epic-design.prompt.md`.
- Estimated input tokens: ~19,360
- Notes: Updated live router, template AGENTS.md/router.md, task prompts, and template tests. No historical task-file migration needed. All 48 template tests pass.
