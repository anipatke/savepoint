---
id: E17-defect-workflow-tui/T007-create-defect-skill
status: done
objective: Add a create-defect skill and scaffold guidance for capturing release-level defects before repair work starts
depends_on: [E17-defect-workflow-tui/T006-doctor-docs-and-templates]
complexity_tier: medium
complexity_reason: "Adds a new workflow skill plus scaffold copy and guidance across template and live docs"
---

# T007: Create Defect Skill

## Problem

The defect lane now exists, but agents need explicit guidance for the moment when a user reports a bug and wants it recorded. That conversation is different from planned epic/task work and different from repairing an already-selected defect.

## Context Files

- `agent-skills/savepoint-create-task/SKILL.md` - existing create-phase skill structure
- `agent-skills/savepoint-build-task/SKILL.md` - existing implementation skill structure
- `templates/project/agent-skills/savepoint-create-task/SKILL.md` - scaffolded create-phase skill structure
- `templates/project/AGENTS.md` - scaffolded agent methodology guidance
- `templates/project/.savepoint/router.md` - scaffolded router guidance
- `AGENTS.md` - live methodology and Codebase Map guidance

## Acceptance Criteria

- [x] A live `savepoint-create-defect` skill exists under `agent-skills/`
- [x] Scaffolded projects receive the same `savepoint-create-defect` skill under `templates/project/agent-skills/`
- [x] Live AGENTS.md distinguishes defect capture from normal task planning/building
- [x] Scaffolded AGENTS.md distinguishes defect capture from normal task planning/building
- [x] Scaffolded router guidance explains manual defect capture and optional transition to `defect-building`
- [x] The E17 task ledger records this late addition

## Implementation Plan

- [x] Create live `agent-skills/savepoint-create-defect/SKILL.md`
- [x] Create scaffold template `templates/project/agent-skills/savepoint-create-defect/SKILL.md`
- [x] Update live AGENTS.md defect workflow guidance
- [x] Update scaffolded AGENTS.md defect workflow guidance
- [x] Update scaffolded router.md with manual defect-capture guidance
- [x] Record this task under E17 for audit traceability

## Context Log

- Files read: `agent-skills/savepoint-create-task/SKILL.md`, `agent-skills/savepoint-build-task/SKILL.md`, `templates/project/agent-skills/savepoint-create-task/SKILL.md`, `templates/project/AGENTS.md`
- Files edited: `AGENTS.md`, `templates/project/AGENTS.md`, `templates/project/.savepoint/router.md`
- Files created: `agent-skills/savepoint-create-defect/SKILL.md`, `templates/project/agent-skills/savepoint-create-defect/SKILL.md`, this task file
- Quality gates: not run; markdown/template-only change
