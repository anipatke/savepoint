---
id: E04-templates-and-prompts/T002-release-and-prompt-assets
status: done
objective: "Create the release starter and agent workflow prompt templates used by Savepoint-managed projects."
depends_on: ["E04-templates-and-prompts/T001-project-template-assets"]
---

# T002: release-and-prompt-assets

## Implementation Plan

- [x] Create `templates/release/v1/PRD.md` with v1 release frontmatter, epic planning placeholders, and agent-facing guidance.
- [x] Create `templates/prompts/prd.prompt.md` for project PRD creation from user intent.
- [x] Create `templates/prompts/design.prompt.md` for project architecture design from an approved PRD.
- [x] Create `templates/prompts/epic-design.prompt.md` for writing an epic `Design.md` from release scope.
- [x] Create `templates/prompts/task-breakdown.prompt.md` for creating independently buildable task files with `depends_on` and checkbox plans.
- [x] Create `templates/prompts/task-planning.prompt.md` for repairing or late-adding a single planned task.
- [x] Create `templates/prompts/task-building.prompt.md` for executing one selected task, ticking checkboxes, and stopping at review.
- [x] Create `templates/prompts/audit-reconciliation.prompt.md` that instructs agents to produce one `.savepoint/audit/{epic}/proposals.md` bundle with delta-only edits where possible, not separate proposal files.
- [x] Run a focused format check over the new release and prompt markdown files.
