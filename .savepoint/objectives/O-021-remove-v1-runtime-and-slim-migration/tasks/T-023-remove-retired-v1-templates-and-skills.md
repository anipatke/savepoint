---
id: T-023
title: Remove retired V1 templates and skills
objective: O-021
planned_by: {role: planner, session: o021-design-20260923}
status: planned
complexity_tier: medium
complexity_reason: "The files themselves are unused, but many init contract tests pin V1 and V2 template trees together, so the tests must be narrowed to V2 without losing the V2 guarantees they also carry."
depends_on: [{task: T-022, requires: clear}]
owner_validation:
    required: false
    accepted_check: ""
---

# T-023: Remove retired V1 templates and skills

## Outcome

The V1 scaffold trees, this repository's retired V1 skills, and the pre-V2
audit-skill upgrade shim are gone. `init` and `upgrade-assets` behave as
they do today for V2 projects, and `upgrade-assets` still refuses a V1
project with its existing migrate guidance.

## User Check

None beyond the gate. Optionally run `savepoint init` in an empty temp
directory and `savepoint upgrade-assets --dry-run` on this repository and
confirm unchanged output.

## Done When

- `templates/project/` and `templates/release/v1/` are deleted.
- `UpgradeProjectAssets` no longer takes a `v1Templates` argument; `main.go`
  and callers are updated. V1-project refusal wording is unchanged.
- The nine retired V1 skill folders in this repository's `agent-skills/`
  (`savepoint-audit-epic`, `savepoint-audit-register`,
  `savepoint-audit-task`, `savepoint-build-task`, `savepoint-create-defect`,
  `savepoint-create-plan`, `savepoint-create-task`, `savepoint-draft-prd`,
  `savepoint-system-design`) and `agent-skills/references/audit-method.md`
  are deleted. `agent-skills/bubbletea-tui-design/` and the four V2 skills
  and three V2 references remain.
- `internal/init/migrate_audit_skill.go` and its test are deleted, and
  `upgradeProjectAssets` no longer calls it. `retire_v1_skills.go` and its
  test remain and still retire the nine V1 skills on a migrated project.
- Tests that pinned V1 trees (for example V1 skill parity in
  `template_freshness_test.go`, V1 sections of `agent_skills_test.go`,
  `lifecycle_test.go`'s `shippedTemplates`, and V1 entries in
  `integration_test.go` and `manifest_test.go`) are removed or narrowed to
  V2. Every V2 guarantee they carried keeps a test.
- A repository search finds no live reference to the deleted paths outside
  `.savepoint/archive/`, immutable Checks, Issues, and closed Task evidence.
- `git diff --check`, `make build && make test-fast` pass.

## Context Files

`templates/project/AGENTS.md`, `templates/release/v1/PRD.md`,
`main.go`, `main_test.go`, `internal/init/upgrade.go`,
`internal/init/upgrade_test.go`, `internal/init/upgrade_schema_test.go`,
`internal/init/migrate_audit_skill.go`,
`internal/init/migrate_audit_skill_test.go`,
`internal/init/retire_v1_skills.go`,
`internal/init/retire_v1_skills_test.go`,
`internal/init/template_freshness_test.go`,
`internal/init/agent_skills_test.go`, `internal/init/lifecycle_test.go`,
`internal/init/integration_test.go`, `internal/init/manifest_test.go`,
`internal/init/scaffold_test.go`, `internal/init/stale_reference_test.go`,
`agent_skills_test.go`.

## Design References

Design sections 2 and 5 (init and upgrade-assets).

## Guardrails

FS-01, FS-02, TPL-01, TPL-02, TPL-03, TPL-04, TEST-03, TEST-06, TEST-08.

## Implementation Plan

1. Delete the template trees and retired skills; build and list failing
   tests.
2. Remove the `v1Templates` parameter and the audit-skill shim call.
3. For each failing test, delete V1-only assertions and keep V2 ones.
4. Search for stale references to the deleted paths and fix live ones.
5. Run `make build && make test-fast`.

## Boundaries

No change to V2 templates, V2 skills, the upgrade provenance manifest
format, or `.savepoint/archive/v1/`. Do not rewrite historical records
that mention deleted paths.

## Technical Verification

Focused `internal/init` and root tests during iteration;
`make build && make test-fast` at handoff (TEST-03 evidence that V2 user
content is still preserved by upgrade).

## Technical Evidence

Pending execution.

## Drift Notes

None expected.
