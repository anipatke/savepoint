---
id: E18-template-skill-optimisation/T004-artifact-templates
status: done
objective: Add explicit artifact template blocks to savepoint-audit, savepoint-create-task, and savepoint-system-design skills so any agent produces consistent E##-Audit.md, T###-slug.md, and E##-Detail.md files across all projects
depends_on: [E18-template-skill-optimisation/T001-canonical-guides]
complexity_tier: medium
complexity_reason: "Adds canonical artifact templates across multiple skills and scaffolded copies"
---

# T004: Artifact Templates

## Problem

Three of the four Savepoint artifact types have no explicit format template in their skill. Only defects (`savepoint-create-defect`) include a canonical markdown template block. Tasks, epics, and audit files are described in prose or by section name only, which allows agents to produce inconsistent frontmatter, title formats, and section structures — as seen in the E18 audit file diverging from E17's established format.

## Context Files

- `agent-skills/savepoint-audit/SKILL.md`
- `agent-skills/savepoint-create-task/SKILL.md`
- `agent-skills/savepoint-system-design/SKILL.md`
- `agent-skills/savepoint-create-defect/SKILL.md`
- `templates/project/agent-skills/savepoint-audit/SKILL.md`
- `templates/project/agent-skills/savepoint-create-task/SKILL.md`
- `templates/project/agent-skills/savepoint-system-design/SKILL.md`
- `.savepoint/releases/v1.2/epics/E17-defect-workflow-tui/E17-Audit.md`

## Acceptance Criteria

- [x] `savepoint-audit` includes an explicit `E##-Audit.md` template block with exact frontmatter (`type: audit-findings`, `audited: {date}`), title format, `- [ ]` checkbox list for Code Style Review, and `### Target File` / `### Replace` / `### With` block structure under `## Proposed Changes`
- [x] `savepoint-create-task` includes an explicit `T###-slug.md` template block with exact frontmatter fields (`id`, `status`, `objective`, `depends_on`) and all required sections
- [x] `savepoint-system-design` includes an explicit `E##-Detail.md` template block with exact frontmatter (`type: epic-design`, `status: planned`) and required sections
- [x] Root skill files and scaffolded copies remain identical after changes
- [x] `make build && make test` passes

## Implementation Plan

- [x] Read `savepoint-create-defect` SKILL.md to confirm the defect template block pattern to follow
- [x] Read `E17-Audit.md` to confirm the canonical audit file format
- [x] Add audit template block to `savepoint-audit` SKILL.md
- [x] Add task template block to `savepoint-create-task` SKILL.md
- [x] Add epic detail template block to `savepoint-system-design` SKILL.md
- [x] Mirror all three changes to `templates/project/agent-skills/`
- [x] Run `make build && make test`

## Context Log

Files read:
- `.savepoint/router.md`
- `.savepoint/releases/v1.2/epics/E18-template-skill-optimisation/E18-Detail.md`
- `.savepoint/releases/v1.2/epics/E18-template-skill-optimisation/tasks/T004-artifact-templates.md`
- `agent-skills/savepoint-build-task/SKILL.md`
- `agent-skills/savepoint-audit/SKILL.md`
- `agent-skills/savepoint-create-task/SKILL.md`
- `agent-skills/savepoint-system-design/SKILL.md`
- `agent-skills/savepoint-create-defect/SKILL.md`
- `templates/project/agent-skills/savepoint-audit/SKILL.md`
- `templates/project/agent-skills/savepoint-create-task/SKILL.md`
- `templates/project/agent-skills/savepoint-system-design/SKILL.md`
- `.savepoint/releases/v1.2/epics/E17-defect-workflow-tui/E17-Audit.md`
- `agent_skills_test.go` (targeted read after quality gate failure)

Files edited:
- `.savepoint/releases/v1.2/epics/E18-template-skill-optimisation/tasks/T004-artifact-templates.md`
- `agent-skills/savepoint-audit/SKILL.md`
- `agent-skills/savepoint-create-task/SKILL.md`
- `agent-skills/savepoint-system-design/SKILL.md`
- `templates/project/agent-skills/savepoint-audit/SKILL.md`
- `templates/project/agent-skills/savepoint-create-task/SKILL.md`
- `templates/project/agent-skills/savepoint-system-design/SKILL.md`
- `agent-skills/savepoint-*/SKILL.md` and `templates/project/agent-skills/savepoint-*/SKILL.md` line endings normalized to LF for the existing discovery-frontmatter test.

Quality gates:
- `make build` passed.
- `make test` passed.

Notes:
- Root and scaffolded copies were compared with `Compare-Object` after edits and had no differences.
- Could not press `p` in the non-interactive TUI from this shell session.
