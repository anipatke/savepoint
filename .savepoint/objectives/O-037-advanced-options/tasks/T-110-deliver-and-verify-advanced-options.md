---
id: T-110
title: Deliver and verify Advanced Options
objective: O-037
status: done
depends_on: [{task: T-109, requires: clear}]
owner_validation:
  required: true
  accepted_check: ""
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
check_waiver:
  task: T-110
  reason: Owner completed this Task via the board without requesting a Task Check.
  actor:
    role: owner
    session: board-owner
  recorded_at: "2026-10-03T22:03:31Z"
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

Commands: `make build && make test-full` on 2026-10-04, go toolchain per `go version`, exit 0 (fresh, after the last code change; docs-only edits followed none). Native windows-tests CI evidence is not produced locally; the owner supplies it for the Full Objective Check.

Per criterion:
- Fresh init off; old projects stay off; explicit choices survive: `TestV2ScaffoldConfigShipsParallelPlanningOff` (template now carries `features.parallel_planning: false`); `TestUpgradePreservesOwnerConfigAndFeatureChoices` (no key, explicit on, explicit off with CRLF, quoted value) keeps `config.yml` byte-identical; absent key reads off in `TestParallelPlanningDefaultsOffWithoutConfigFile`/`TestOptionsOpensWithParallelPlanningOffAndExplainsIt` (T-108/T-109).
- Upgrades preserve config, comments, unrelated keys and edited assets; no advice or worktrees: same upgrade test, with a manifest-tracked edited task skill kept and `.worktrees`, `worktrees`, `.savepoint/lanes`, `.savepoint/advice` absent. Upgrade production code was not changed.
- Integration: `TestOptionsJourneyOnTemporaryProject` (scaffold config, save, restart, off, on, external edit, stale refusal then retry, close; only `config.yml` written, nothing created) and `TestOptionsLeaveLifecycleAndBoardIdenticalWhenOn` (same board, Code Health chip and Task-advance record files with the option on). Failed write and close focus are covered by `TestOptionsUnwritableConfigReportsAndSavesNothing` and `TestOptionsCloseReturnsFocusToTheOriginSurface` from T-109; no lane data exists in any of them.
- README (Advanced Options section), CHANGELOG (Unreleased) and Design (`config.yml` layout line, `o` keybinding, Advanced Options paragraph in section 8) describe only the implemented option.
- Gate and owner validation: full gate recorded above. Owner walkthrough (open, save, restart, switch off, Code Health unchanged) and any Task Check or waiver are pending; none is recorded here.

Files changed: templates/project-v2/.savepoint/config.yml, internal/init/upgrade_test.go, internal/init/v2_scaffold_test.go, internal/board/v2/options_integration_test.go (new), README.md, CHANGELOG.md, .savepoint/Design.md. Extra reads: internal/init/upgrade.go (planned), internal/board/v2/options_test.go, objectives_test.go, columns_view_test.go, template_freshness_test.go (test helpers), AGENTS.md and agent-skills/savepoint-task/SKILL.md (not extra reads).

Limitations: no real `savepoint init` run (human-only); scaffold checked through the template file. Windows CI not run locally. The plan's `.savepoint/Design.md` heading names "sections 1, 8, 9 and 13"; only section 2's layout line and section 8 were edited.

## Drift Notes

Record material acceptance/interface deviations for planner review.

