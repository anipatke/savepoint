---
id: T-110
title: Deliver and verify Advanced Options
objective: O-037
status: planned
depends_on: [{task: T-109, requires: clear}]
owner_validation: {required: true}
planned_by: {role: planner, session: codex-options-first-2026-10-03}
complexity_tier: medium
complexity_reason: Preference adoption must preserve authored configuration and work independently of lane implementation.
planned_reads:
  - internal/data/feature_preferences.go
  - internal/board/v2/options.go
  - internal/init/upgrade.go
  - internal/init/upgrade_test.go
  - internal/init/v2_scaffold_test.go
  - templates/project-v2/.savepoint/config.yml
  - README.md
  - CHANGELOG.md
  - .github/workflows/ci.yml
planned_writes:
  - templates/project-v2/.savepoint/config.yml
  - internal/init/upgrade_test.go
  - internal/init/v2_scaffold_test.go
  - internal/board/v2/options_integration_test.go
  - README.md
  - CHANGELOG.md
  - .savepoint/Design.md
---

# Deliver and verify Advanced Options

## Outcome

Fresh and existing projects receive the off-by-default preference and documented Advanced Options screen, with standalone integration evidence before lane implementation begins.

## User Check

Open Advanced Options, save the preference, restart, and switch it off. Verify existing Code Health remains unchanged. At this delivery stage, the preference is saved for optional lane advice delivered by O-033; no groupings are promised yet.

## Done When

- Fresh init has features.parallel_planning false; old projects with no key remain off and existing explicit choices survive upgrades.
- Fixture-based upgrades preserve owner configuration, comments, unrelated keys and edited assets; no advice records or worktrees are created.
- Integration exercises save/restart/off/on/external edit/stale or failed write/close focus using temporary projects, with identical existing lifecycle and Code Health behavior.
- README/CHANGELOG explain the single option, manual config editing and optional later lane advice; reconcile Design to implemented settings only.
- Record per-criterion evidence, fresh full gate and owner validation. An optional Task Check needs a fresh checker or explicit owner waiver; mandatory Full Objective Check remains independent.

## Context Files

- `internal/data/feature_preferences.go`
- `internal/board/v2/options.go`
- `internal/init/upgrade.go`
- `internal/init/upgrade_test.go`
- `internal/init/v2_scaffold_test.go`
- `templates/project-v2/.savepoint/config.yml`
- `README.md`
- `CHANGELOG.md`
- `.github/workflows/ci.yml`

feature_preferences.go and options.go are intentionally created by preceding Tasks. options_integration_test.go is an explicitly new write target.

## Design References

O-037 Confirmed Design — 2026-10-03; project Design sections 1, 8, 9 and 13.

## Guardrails

FS-01, FS-02, FS-04, FS-05, DATA-01, ARCH-02, CFG-02, CFG-03, TPL-02, TPL-04, TEST-01, TEST-02, TEST-03, TEST-04, TEST-08, TEST-09, STYLE-01..10.

## Implementation Plan

1. Confirm the settings API/screen work without a lane projection; return REPLAN REQUIRED for a material contract gap.
2. Add the off scaffold default and temporary-fixture upgrade preservation coverage, without running human-only Savepoint commands.
3. Add standalone settings integration cases, including no lane data at all and unchanged health surfaces.
4. Write adoption guidance/release notes and reconcile implemented Design, then run the full gate and record owner walkthrough.

## Boundaries

No lane engine, record model, grouping or session instructions; no Code Health switch or provider execution. No production upgrade rewrite unless required by demonstrated preservation failure and replanned. Actual worktree rules and ordinary evidence/owner authority remain applicable. Suggested manifests are not execution gates.

## Technical Verification

Fresh make test-full. Native windows-tests CI evidence is produced by repository CI and supplied by the owner for the mandatory independent Full Objective Check under agent-skills/references/check-method.md. No executor writes a Check or runs health check.

## Technical Evidence

Pending execution: named cases/results, actual read/write files and extra-read reasons, command/time/toolchain, owner validation and limitations.

## Drift Notes

Record material acceptance/interface deviations for planner review.

