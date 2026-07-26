---
id: E35-split-task-and-epic-audits/T002-generated-routing-and-guidance
status: planned
objective: Route generated projects to the correct split audit skill, remove generic audit terminology from live guidance, and update cross-references across router, AGENTS, Design, audit-register, build-task, and README.
depends_on:
  - E35-split-task-and-epic-audits/T001-split-audit-contracts
complexity_tier: medium
complexity_reason: Router and catalog changes must align request intent without introducing a new workflow state, and cross-reference updates span five file pairs.
---

# T002: Generated Routing and Guidance

## Problem

Generated router, contributor, and user guidance still presents audit as one ambiguous phase skill and does not describe the task-audit override. Cross-references to the old generic skill remain in AGENTS.md, Design.md, audit-register, build-task, and README.

## Context Files

- `AGENTS.md`
- `templates/project/AGENTS.md`
- `templates/project/.savepoint/router.md`
- `templates/project/.savepoint/Design.md`
- `README.md`
- `agent-skills/savepoint-audit-task/SKILL.md`
- `agent-skills/savepoint-audit-epic/SKILL.md`
- `agent-skills/savepoint-build-task/SKILL.md`
- `agent-skills/savepoint-audit-register/SKILL.md`
- `templates/project/agent-skills/savepoint-audit-register/SKILL.md`
- `.savepoint/Design.md`

## Acceptance Criteria

- [ ] The live and generated skill tables map `audit-pending` to `savepoint-audit-epic`.
- [ ] Generated routing maps `task-building` plus an explicit task audit or re-audit request to `savepoint-audit-task` as a request-qualified override, not a router state.
- [ ] Router template `## Manual Overrides` includes a task-audit override line next to the existing epic-audit override.
- [ ] AGENTS.md (live and template) Audit section describes the split skills and references `.savepoint/Guardrails.md` and `.savepoint/Health-Check.md` when present.
- [ ] `templates/project/.savepoint/Design.md` audit pipeline reference updates `savepoint-audit` to `savepoint-audit-epic`.
- [ ] `agent-skills/savepoint-audit-register/SKILL.md` (live and template) lines 14, 39 reference `savepoint-audit-epic` instead of `savepoint-audit`.
- [ ] `agent-skills/savepoint-build-task/SKILL.md` (live and template) references Health-Check.md Quick mode at task handoff (graceful when absent).
- [ ] `.savepoint/Design.md` (repo's own) updates bundled skills list (line 17) and audit pipeline reference (line 130) to `savepoint-audit-epic`.
- [ ] README documents the two audit intents, shared method, enriched rigor (scope locks, matrices, materiality), and existing-project compatibility behavior.
- [ ] No live catalog, router template, generated guide, or user documentation exposes `savepoint-audit` as a triggerable skill or alias.
- [ ] Historical release records are left unchanged.
- [ ] `make build && make test` passes.

## Implementation Plan

- [ ] Replace the generic audit row in live and generated skill activation tables.
- [ ] Add the request-qualified task-audit override to the generated router without adding state.
- [ ] Update AGENTS.md Audit section to name split skills and add Guardrails/Health-Check discoverability notes.
- [ ] Update Design.md template audit reference to `savepoint-audit-epic`.
- [ ] Update audit-register skill (both copies) lines 14, 39.
- [ ] Update build-task skill (both copies) for audit-deferral and Health-Check Quick mode.
- [ ] Update repo's own Design.md bundled skills list (line 17) and audit reference (line 130).
- [ ] Update README skill catalog and audit workflow documentation.
- [ ] Review live guidance for ambiguous generic references while excluding historical release artifacts.

## Context Log

Pending.
