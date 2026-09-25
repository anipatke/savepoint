---
id: E16-pre-prod-refinement/T002-skill-context-budgets
status: done
stage: build
objective: Update bundled skills and scaffolded skill templates to enforce minimal file reads and token-efficient phase workflows
depends_on: []
---

# T002: Add Skill Context Budgets

## Context Files

- `AGENTS.md` — live workflow rules and codebase map
- `agent-skills/savepoint-system-design/SKILL.md` — epic-design phase guidance
- `agent-skills/savepoint-create-task/SKILL.md` — task creation phase guidance
- `agent-skills/savepoint-build-task/SKILL.md` — implementation phase guidance
- `templates/project/AGENTS.md` — scaffolded workflow rules for new projects
- `templates/project/agent-skills/savepoint-system-design/SKILL.md` — scaffolded epic-design skill
- `templates/project/agent-skills/savepoint-create-task/SKILL.md` — scaffolded task creation skill
- `templates/project/agent-skills/savepoint-build-task/SKILL.md` — scaffolded implementation skill
- `internal/init/template_freshness_test.go` — template freshness expectations
- `agent_skills_test.go` — bundled/scaffolded skill parity checks

## Acceptance Criteria

- [ ] Phase skills state that agents must read the fewest files needed for the current phase
- [ ] Epic-design guidance limits source inspection to the implementation boundary and discourages full test-suite reads before implementation
- [ ] Create-task guidance defines a strict "create epic/task only" read budget and prohibits source reads unless explicitly requested
- [ ] Build-task guidance keeps required context-file reads but discourages exploratory reads outside listed context files before implementation starts
- [ ] Scaffolded skill templates match the live skill guidance
- [ ] Scaffolded `AGENTS.md` includes concise context-budget and tool-call discipline rules
- [ ] Existing template freshness or skill parity tests are updated where needed
- [ ] `make build && make test` passes

## Implementation Plan

- [x] Add a concise context-budget section to live `AGENTS.md`
- [x] Update live `savepoint-system-design` guidance with epic-design read limits
- [x] Update live `savepoint-create-task` guidance with create-only read limits
- [x] Update live `savepoint-build-task` guidance with pre-build context discipline
- [x] Mirror the same changes into `templates/project/AGENTS.md`
- [x] Mirror the same changes into scaffolded skill templates
- [x] Update freshness/parity tests if their assertions need to cover the new guidance
- [x] Run `make build && make test`
