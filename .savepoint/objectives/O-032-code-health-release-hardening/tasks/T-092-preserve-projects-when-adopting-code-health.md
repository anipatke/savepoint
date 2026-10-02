---
id: T-092
title: Preserve projects when adopting Code Health
objective: O-032
status: done
complexity_tier: high
complexity_reason: Scaffold, upgrades and migration boundaries must preserve owner-authored files.
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: planning-o032-20261002-owner-confirmed}
check_waiver:
    task: T-092
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-02T04:22:27Z"
---

# Preserve projects when adopting Code Health

## Outcome

Fresh and existing projects receive truthful Code Health guidance without silent configuration, history or authored-content replacement.

## User Check

Inspect a temporary existing project after a tested asset upgrade: health choices, history and edited Design remain intact.

## Done When

- Temporary fresh/init and existing schema-2 upgrade cases preserve pre-existing health config, snapshots, edited Design and unmanaged content; repeat and dry-run behavior is verified.
- Legacy migration and incompatible project-schema paths remain named and preview-first; asset upgrade never rewrites schema or invents health configuration/snapshots.
- Confirm owner-only setup/report and Full-Check-only health collection guidance matches current runtime in scaffold assets. Preserve canonical/scaffold skill identity via existing tests.
- Exercise partial write or asset conflict with exact retained bytes and reported recovery paths; repair only bounded adoption regressions.
- Record current command/guidance separation, named preservation tests and fresh full gate.

## Context Files

`internal/init/scaffold.go`; `internal/init/upgrade.go`; `internal/init/manifest.go`; `internal/init/v2_scaffold_test.go`; `internal/init/upgrade_schema_test.go`; `internal/init/provenance_contract_test.go`; `internal/init/upgrade_failure_test.go`; `internal/migrate/live_boundary_test.go`; `internal/migrate/end_to_end_test.go`; `templates/project-v2/AGENTS.md`; `templates/project-v2/.savepoint/Design.md`; `main_health_test.go`.

## Design References

O-032 Confirmed Release-Hardening Design and Owner-Requested Health Analysis; Design sections 1, 6, 7, 9, 12 and 13 as applicable.

## Guardrails

FS-01, FS-04, FS-05, FS-06, ARCH-01, ARCH-02, ARCH-03, CFG-01, CFG-02, TEST-01, TEST-02, TEST-03, TEST-04, TEST-05, TEST-08, TEST-09, STYLE-01, STYLE-02, STYLE-07, TPL-01, TPL-02, TPL-03, TPL-04, CFG-03, REL-01, REL-02, REL-03.

## Implementation Plan

1. Trace scaffold/managed asset ownership and existing health entry wiring.
2. Build temporary adoption fixtures carrying authored health/config/Design content.
3. Add missing preservation and schema-boundary cases; reconcile shipped guidance to implemented commands.
4. Validate conflicts/repeat/dry-run and migration full gate.

## Boundaries

No agent-running init/upgrade/migrate CLI, provider auto-install/setup, new schema conversion, authored policy replacement or changes to canonical skills.

## Technical Verification

Fresh `make test-full` for this platform/migration/distribution-sensitive Task; focused `make test-focused TEST=... PKGS=...` only during iteration. Record named per-criterion evidence and limitations. The later independent Check applies `agent-skills/references/check-method.md`; the mandatory Full Objective Check runs its fresh full gate and official Code Health collection. This executor never writes a Check or runs savepoint health setup/report/check.

## Technical Evidence

Command: `make build && make test-full` on linux/amd64, go1.26.2, 2026-10-02 ~14:22 AEST: passed (build, test, and linux/darwin/windows builds).

Per criterion:
1. Fresh/init and schema-2 upgrade preservation: new `TestMainUpgradeAssetsPreservesAdoptedHealthAndAuthoredContent` (main_health_test.go) keeps edited health config, snapshots, reports, Design.md, Health-Check.md, config.yml and an unmanaged file byte-identical; dry-run changes nothing; repeat run changes nothing (hash and mtime). Fresh init writes no health files: existing `TestMainInitEndsWithHealthPreviewAndWritesNoHealthConfig`.
2. Migration/schema boundaries: existing `TestUpgradeProjectAssets_refusesV1WithoutMutation`, `_neverWritesSchemaVersion`, `_malformedSchemaVersionRefusesCleanly`, `_unsupportedSchemaVersionRefusesCleanly` and `internal/migrate` preview-first tests pass unchanged; no health config or snapshots are invented.
3. Guidance: templates/project-v2/AGENTS.md lines 189-191 already state setup is human-only and `health check O-###` runs only in a Full Objective Check; matches runtime wiring in main.go. No template change needed; skill identity tests pass.
4. Conflict recovery: the new test's edited-skill conflict keeps the owner bytes, writes `SKILL.md.new` and reports "incoming written to .new". Repair: a repeat run rewrote an identical `.new` sidecar (mtime churn); `writeSidecar` in internal/init/upgrade.go now skips an identical sidecar. Existing `upgrade_failure_test.go` write-failure cases pass.
5. Separation and full gate recorded above.

Files read: internal/init/{scaffold,manifest,upgrade}.go, main.go, main_health_test.go, main_test.go helpers, upgrade_schema_test.go, templates/project-v2/AGENTS.md (grep). Extra reads beyond Context Files: main.go and main_test.go (wiring and snapshot helpers). Files changed: internal/init/upgrade.go, main_health_test.go, this Task.

Limitations: partial-write failure was covered by existing tests rather than a new one; no Task Check or health command was run; Windows behavior verified only by cross-build.

## Drift Notes

Record implemented responsibility or interface changes for planner reconciliation before the mandatory Full Objective Check. If evidence requires a material unknown repair, return REPLAN REQUIRED rather than inventing another scope.
