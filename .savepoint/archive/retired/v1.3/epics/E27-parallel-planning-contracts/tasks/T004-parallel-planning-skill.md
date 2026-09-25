---
id: E27-parallel-planning-contracts/T004-parallel-planning-skill
status: planned
objective: Teach task planning to define conservative parallel groups and exact write ownership.
depends_on:
  - E27-parallel-planning-contracts/T003-parallel-eligibility-validation
complexity_tier: medium
complexity_reason: Workflow guidance and templates must align with the validated metadata contract without bloating default tasks.
---

# T004: Parallel Planning Skill

## Problem

Advanced metadata is unsafe unless planning guidance defines when to use it and when to remain sequential.

## Context Files

- `agent-skills/savepoint-create-task/SKILL.md`
- `templates/project/agent-skills/savepoint-create-task/SKILL.md`
- `templates/project/AGENTS.md`

## Acceptance Criteria

- [ ] The planning skill adds execution metadata only when advanced mode is enabled and the user requests parallel planning.
- [ ] Guidance requires exact file ownership, satisfied dependencies, and contract/schema work before parallel consumers.
- [ ] Ambiguous ownership, shared generated files, or overlapping paths require sequential tasks.
- [ ] Default task examples remain unchanged and do not require execution metadata.
- [ ] Generated project guidance matches the repository skill exactly where managed.

## Implementation Plan

- [ ] Add a concise advanced-mode planning section to the source skill.
- [ ] Define group naming, exact write-scope, dependency, and fallback rules.
- [ ] Add the parallel launch-scope exception to the generated AGENTS workflow without weakening normal router routing.
- [ ] Mirror managed skill changes into the project template.
- [ ] Verify package-owned asset consistency tests.

## Context Log

Pending.
