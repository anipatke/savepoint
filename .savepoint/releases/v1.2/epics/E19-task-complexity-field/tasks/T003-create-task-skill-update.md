---
id: E19-task-complexity-field/T003-create-task-skill-update
status: planned
objective: Teach create-task to assess complexity and record it in task plans
depends_on: [E19-task-complexity-field/T001-task-complexity-model]
---

# T003: Create-Task Skill Update

## Problem

The task-planning skill still writes tasks without complexity guidance. The planner needs a shared rubric, a repeatable way to choose a tier, and an instruction to write the tier and reason into each new task.

## Context Files

- `agent-skills/savepoint-create-task/SKILL.md`
- `templates/project/agent-skills/savepoint-create-task/SKILL.md`
- `internal/init/template_freshness_test.go`

## Acceptance Criteria

- [ ] The root and scaffolded `create-task` skill docs include the complexity rubric and the tier-selection rules
- [ ] The skill instructions tell the planner to append complexity metadata to each new task plan
- [ ] The rubric includes the L, M, H, and U definitions plus the highest-tier and U-overrides rules
- [ ] The root and scaffolded skill copies stay textually aligned for the shared workflow contract
- [ ] Any freshness test that guards the skill copy surface covers the updated create-task guidance
- [ ] `make build && make test` passes after the skill updates

## Implementation Plan

- [ ] Update the root create-task skill with the complexity rubric and output requirements
- [ ] Mirror the same guidance into the scaffolded skill copy
- [ ] Extend the template freshness test if needed so the sync surface includes create-task
- [ ] Run `make build && make test`

## Context Log

- Files read:
- Estimated input tokens:
- Notes:
