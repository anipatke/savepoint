---
id: E50-release-validation-cutover/T004-activate-v2-scaffolds-and-workflow-assets
status: done
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

- [x] Fresh init always writes a V2 project and the four-phase V2 routing table; no supported option creates a new V1 project.
- [x] Canonical and shipped V2 skills and shared references are byte-identical and contain no active V1 phase guidance.
- [x] Upgrade on V2 preserves customized assets, provenance sidecars, unknown files, and managed-guide conflict behavior.
- [x] Upgrade on V1 refuses normal asset mutation with explicit migration guidance, except any proven migration-support asset contract.
- [x] A successfully migrated project retires the nine V1 skills without deleting user-authored content or byte-preserved archives.
- [x] Init, upgrade, lifecycle, stale-reference, and template-freshness tests cover the final contract.

## Implementation Plan

- [x] Remove transitional scaffold selection and make the V2 tree the single init default.
- [x] Reconcile the managed AGENTS block and canonical/shipped V2 skill bytes.
- [x] Gate upgrade behavior by final schema and preserve provenance/conflict guarantees.
- [x] Retire V1 workflow assets only after successful migration and verified archive mapping.
- [x] Replace transitional tests with final init/upgrade/retirement matrices.
- [x] Run focused init tests and the required full gates; record evidence.

## Context Log

- Focused evidence: `go test ./internal/init -run 'TestUpgradeProjectAssets_refusesV1WithoutMutation|TestUpgradeProjectAssets_refusesPendingMigrationOnBothTrees|TestUpgradeAssetsFromTree_preservesFrozenPreE47Fixture|TestV2WorkflowAssetsHaveNoActiveV1Routing|TestLifecycleMatrix_fullLoopAcrossProjectKinds' -count=1` and `go test . -run 'TestMainUpgradeAssetsPrintsPartialWorkOnFailure|TestMainUpgradeAssetsV1ProjectRefusesMutationAndNamesMigrateRoute|TestMainUpgradeAssetsV2ProjectInstallsOnlyV2Skills' -count=1`: PASS.
- Full gate: `make build && make test`: PASS (`go test ./...`; `internal/migrate` 135.302s).
- `git diff --check`: PASS.
- Production `init` and `upgrade-assets` now embed/select only the V2 project tree; prompts have their own embed, and legacy projects receive a no-write migration report. V2 upgrades retain the existing manifest, sidecar, managed-guide, unknown-file, and V1-retirement paths.
- Canonical and `templates/project-v2` V2 skills/references remain byte-identical. The explicit `pre-implementation`/legacy skill route was removed from `savepoint-idea`; stale V1 routing terms are covered by `TestV2WorkflowAssetsHaveNoActiveV1Routing`.
- V1 compatibility remains reachable only through the explicit `upgradeAssetsFromTree` history/migration fixture path; production `UpgradeProjectAssets` validates the target, refuses V1 mutation, and names `savepoint migrate --dry-run`/`--apply`.
- Existing V2 scaffold, provenance, retirement, lifecycle, conflict, archive, and template-freshness matrices pass; the V1 lifecycle case now proves the refusal is byte-preserving.
- Extra targeted reads outside the listed Context Files: `internal/init/lifecycle_matrix_test.go` and `main_test.go` to update final lifecycle/CLI matrices; `internal/init/migrate_audit_skill.go` and `internal/data/config.go` to preserve the archive/schema contracts; `internal/migrate/cutover.go` to match migration guidance; `internal/init/agent_skills_test.go`, `internal/init/scaffold_test.go`, `internal/init/validate.go`, `cmd/init.go`, and `cmd/upgrade-assets.go` to verify routing/target semantics.
- No `.savepoint/Health-Check.md` is present, so the Quick health check is skipped per `AGENTS.md`.
