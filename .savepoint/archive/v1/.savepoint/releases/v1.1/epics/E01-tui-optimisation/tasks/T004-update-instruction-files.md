---
id: E01-tui-optimisation/T004-update-instruction-files
status: done
objective: "Update all instruction files to reference new naming convention"
depends_on:
  - E01-tui-optimisation/T002-rename-epic-design-files
  - E01-tui-optimisation/T003-rename-release-prd
---

# T004: Update Instruction Files

## Acceptance Criteria

- `.savepoint/router.md` references `E##-Detail.md` and `v1-PRD.md` (not `Design.md` or `PRD.md`)
- `AGENTS.md` references `E##-Detail.md` pattern
- `agent-skills/savepoint-create-plan/SKILL.md` updated
- `agent-skills/savepoint-system-design/SKILL.md` updated
- `agent-skills/savepoint-build-task/SKILL.md` updated
- `agent-skills/savepoint-audit/SKILL.md` updated
- `agent-skills/savepoint-create-task/SKILL.md` updated
- `templates/project/AGENTS.md` updated
- `templates/prompts/task-building.prompt.md` updated
- `templates/prompts/task-planning.prompt.md` updated
- `templates/prompts/task-breakdown.prompt.md` updated
- Root `.savepoint/Design.md` references unchanged (per user request)

## Implementation Plan

- [x] Update `.savepoint/router.md` — all epic `Design.md` → `E##-Detail.md`, `PRD.md` → `v1-PRD.md`
- [x] Update `AGENTS.md` — epic design file path pattern
- [x] Update `agent-skills/savepoint-create-plan/SKILL.md`
- [x] Update `agent-skills/savepoint-system-design/SKILL.md`
- [x] Update `agent-skills/savepoint-build-task/SKILL.md`
- [x] Update `agent-skills/savepoint-audit/SKILL.md`
- [x] Update `agent-skills/savepoint-create-task/SKILL.md`
- [x] Update `templates/project/AGENTS.md`
- [x] Update `templates/prompts/task-building.prompt.md`
- [x] Update `templates/prompts/task-planning.prompt.md`
- [x] Update `templates/prompts/task-breakdown.prompt.md`

## Context Log

Files read:
- All instruction files listed above

Estimated input tokens: 800

Notes:
- Only change path references, do not alter surrounding content
- Root Design.md references stay as-is
- Audit must-fix applied: stale implementation checklist items were ticked after verifying the referenced instruction-file updates were present in the completed task scope.
