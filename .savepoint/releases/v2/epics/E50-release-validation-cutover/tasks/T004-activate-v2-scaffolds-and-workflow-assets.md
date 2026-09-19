---
id: E50-release-validation-cutover/T004-activate-v2-scaffolds-and-workflow-assets
status: planned
objective: Make V2 scaffolds, routing, skills, and upgrade behavior the only shipped live workflow.
depends_on:
    - E50-release-validation-cutover/T002-make-live-command-routing-v2-only
complexity_tier: high
complexity_reason: Coordinates embedded templates, upgrades, canonical skills, and managed guidance.
---

# T004: Activate V2 scaffolds and workflow assets

## Problem

V2 templates and skills exist but transitional V1 activation remains in the repository and upgrade paths. Cutover must make the four V2 phases authoritative without overwriting customized project assets or stranding legacy projects before migration.

## Context Files

- `.savepoint/releases/v2/epics/E50-release-validation-cutover/E50-Detail.md`
- `main.go`
- `internal/init/scaffold.go`
- `internal/init/upgrade.go`
- `internal/init/retire_v1_skills.go`
- `internal/init/v2_scaffold_test.go`
- `internal/init/upgrade_schema_test.go`
- `internal/init/retire_v1_skills_test.go`
- `internal/init/provenance_contract_test.go`
- `internal/init/template_freshness_test.go`
- `templates/project-v2/AGENTS.md`
- `templates/project-v2/agent-skills/savepoint-idea/SKILL.md`
- `templates/project-v2/agent-skills/savepoint-design/SKILL.md`
- `templates/project-v2/agent-skills/savepoint-task/SKILL.md`
- `templates/project-v2/agent-skills/savepoint-check/SKILL.md`
- `templates/project-v2/agent-skills/references/check-method.md`
- `templates/project-v2/agent-skills/references/issue-capture.md`
- `templates/project-v2/agent-skills/references/commands-and-procedures.md`
- `agent-skills/savepoint-idea/SKILL.md`
- `agent-skills/savepoint-design/SKILL.md`
- `agent-skills/savepoint-task/SKILL.md`
- `agent-skills/savepoint-check/SKILL.md`

## Acceptance Criteria

- [ ] Fresh init always writes a V2 project and the four-phase V2 routing table; no supported option creates a new V1 project.
- [ ] Canonical and shipped V2 skills and shared references are byte-identical and contain no active V1 phase guidance.
- [ ] Upgrade on V2 preserves customized assets, provenance sidecars, unknown files, and managed-guide conflict behavior.
- [ ] Upgrade on V1 refuses normal asset mutation with explicit migration guidance, except any proven migration-support asset contract.
- [ ] A successfully migrated project retires the nine V1 skills without deleting user-authored content or byte-preserved archives.
- [ ] Init, upgrade, lifecycle, stale-reference, and template-freshness tests cover the final contract.

## Implementation Plan

- [ ] Remove transitional scaffold selection and make the V2 tree the single init default.
- [ ] Reconcile the managed AGENTS block and canonical/shipped V2 skill bytes.
- [ ] Gate upgrade behavior by final schema and preserve provenance/conflict guarantees.
- [ ] Retire V1 workflow assets only after successful migration and verified archive mapping.
- [ ] Replace transitional tests with final init/upgrade/retirement matrices.
- [ ] Run focused init tests and the required full gates; record evidence.

## Context Log

Pending.
