---
id: T-038
objective: O-022
title: Start new projects with a Goal
status: done
complexity_tier: low
complexity_reason: Adds one scaffolded Release record and a router default; reuses existing name interpolation.
depends_on: []
owner_validation:
    required: true
    accepted_check: ""
planned_by: {role: planner, session: o022-goal-required-20260925}
check_waiver:
    task: T-038
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-25T01:47:47Z"
---

# Start new projects with a Goal

## Outcome

`savepoint init` creates a project that already has a Goal. R-001 is titled
after the project, and the router selects it. The Idea skill then fills in
that Goal's Outcome with the owner, instead of treating a Goal as optional.

## User Check

Run `savepoint init` in an empty folder. Confirm
`.savepoint/releases/.../Release.md` exists as R-001, titled with the folder's
project name, and the router reads `release: R-001`. Run `savepoint resume` and
`savepoint doctor` there and confirm neither reports a missing Goal.

## Done When

- The V2 scaffold contains an R-001 Release record with stub Outcome, Why,
  Success Conditions, and Boundaries. Its title uses `{{PROJECT_NAME}}`, and it
  loads as a live Goal, not an archived one.
- The scaffold router has `release: R-001`.
- A freshly initialized project loads strictly with no T-036 diagnostics, and
  its Next line is unchanged apart from the Goal context.
- The Idea skill, live and scaffold byte-identical, tells the planner to fill
  in R-001's sections with the owner rather than create a new Goal. The design
  skill's "Allocate a stable global R-###" step still applies to additional
  Goals.
- Upgrade-assets behaviour for existing projects is unchanged: it does not
  add a Goal to a project's records.
- Tests cover the scaffold contents, interpolation, the router default, strict
  load, and template freshness.

## Context Files

`internal/init/scaffold.go`, `internal/init/v2_scaffold_test.go`,
`internal/init/template_freshness_test.go`, `cmd/init_test.go`, `main.go`,
`templates/project-v2/.savepoint/router.md`,
`agent-skills/savepoint-idea/SKILL.md`,
`templates/project-v2/agent-skills/savepoint-idea/SKILL.md`.
New file: `templates/project-v2/.savepoint/releases/R-001-first-goal/Release.md`.

## Design References

Design sections 2, 6, and 10; O-022 Architectural Considerations.

## Guardrails

TPL-01, TPL-02, TPL-03, TPL-04, FS-01, TEST-01, TEST-02, TEST-05.

## Implementation Plan

1. Add the R-001 scaffold record and the router default. Confirm the
   `all:templates/project-v2/.savepoint` embed picks up the new directory.
2. Update the Idea skill (live and scaffold).
3. Add scaffold and strict-load tests.
4. Run focused tests while iterating, then `make build && make test-fast`.

## Boundaries

No changes to existing projects' records, and no Goal-selection UI.

## Technical Verification

Focused tests during iteration; `make build && make test-fast` for handoff.

## Technical Evidence

Execution started at `stage: build`. Router selection already matched O-022/T-038; Task dependencies are empty and O-022 is in progress. Router `release: R-006` is unchanged.

Extra reads beyond Context Files:
- `.savepoint/Design.md`, sections 2, 6, and 10, to check the planned Goal record against the current storage and CLI architecture.
- `.savepoint/releases/R-006-v2/Release.md`, to use the existing V2 Goal record's serialized shape for the fresh-project placeholder.
- `internal/init/agent_skills_test.go`, around `TestV2SkillsTeachOptionalGoalWorkflow`, because the fast gate exposed assertions that still require the Idea skill to describe Goals as optional; updated those assertions to the new required-Goal behavior.
- `main_test.go`, around `TestMainInitWritesNoV1OnlyPathOrSkill`, because the full Task gate exposed an assertion that forbids all `.savepoint/releases/` content; narrow it to V1-only artifacts while requiring the new R-001 Goal.
- `main_resume_matrix_test.go`, around the fresh-init case in `TestResumeMatrix_everyRungReachedExactlyOnce`, because the fast gate exposed expectations for the old missing-Goal diagnostic and its repair instruction.

Implementation: added the R-001 placeholder to the V2 scaffold, selected it in the scaffold router, and updated both Idea skill copies to fill that placeholder from owner answers. Added scaffold content/interpolation, embedded-init, router, strict-load, resume, guidance parity, and upgrade-assets boundary coverage. `internal/init/scaffold.go` already interpolates `{{PROJECT_NAME}}`, and `main.go` already embeds the full `.savepoint` tree with `all:`; neither required a code change.

Per-criterion evidence:
1. R-001 has stub Outcome, Why, Success Conditions, and Boundaries; its title is interpolated and it loads as a live Goal: `TestV2ScaffoldCreatesProjectGoalWithInterpolatedName`, `TestMainInitScaffoldsV2ProjectWithProjectGoal`, and `TestV2ScaffoldLoadsCleanThroughRuntimeSchemaAndIndex` passed.
2. The scaffold router selects `release: R-001`: `TestV2ScaffoldRouterOpensAtIdeaWithObjectiveField` and `TestMainInitScaffoldsV2ProjectWithProjectGoal` passed.
3. A fresh scaffold passes strict runtime/index loading and resolves resume without a missing-Goal prompt: `TestV2ScaffoldLoadsCleanThroughRuntimeSchemaAndIndex`, `TestV2ScaffoldResumesWithItsProjectGoalSelected`, `TestResumeMatrix_everyRungReachedExactlyOnce/fresh_savepoint_init_scaffold`, and the embedded init integration test passed.
4. Both Idea skill copies direct the owner-assisted completion of the R-001 placeholder and remain byte-identical; the additional-Goal allocation guidance remains covered: `TestIdeaGuidanceFillsTheFreshProjectsGoal`, `TestV2SkillsTeachProjectGoalWorkflow`, and `TestProjectGuidanceTemplatesMirrorLiveGuidance` passed.
5. Upgrade-assets does not add a Goal to an existing project: the no-release-directory assertion in `TestUpgradeDeliversPolicyAssetsFromRealTemplates` passed. `TestV2ScaffoldIntoPopulatedDirectoryPreservesExistingFiles` also passed, preserving the existing-content guard.

Verification:
- `go test ./internal/init -run 'TestV2Scaffold|TestIdeaGuidanceFillsTheFreshProjectsGoal|TestProjectGuidanceTemplatesMirrorLiveGuidance|TestUpgradeDeliversPolicyAssetsFromRealTemplates'` passed during iteration.
- After updating stale optional-Goal expectations, `go test ./internal/init -count=1` passed.
- `go test . -run 'TestMainInitScaffoldsV2ProjectWithProjectGoal|TestMainInitWritesNoV1OnlyPathOrSkill' -count=1` passed.
- `go test . -run 'TestResumeMatrix_everyRungReachedExactlyOnce' -count=1` passed.
- Final `make build && make test-fast` passed. Earlier gate attempts identified and corrected stale optional-Goal, forbidden-release-directory, and missing-Goal resume expectations.
- `./savepoint resume` strict-loaded the current V2 index successfully.
- `git diff --check` passed.

Files read: `.savepoint/router.md`; this Task and `.savepoint/objectives/O-022-require-goal-context/Objective.md`; `.savepoint/Guardrails.md`; `internal/init/scaffold.go`; `internal/init/v2_scaffold_test.go`; `internal/init/template_freshness_test.go`; `cmd/init_test.go`; `main.go`; `templates/project-v2/.savepoint/router.md`; `agent-skills/savepoint-idea/SKILL.md`; `templates/project-v2/agent-skills/savepoint-idea/SKILL.md`; and all five logged extra reads above. Files changed: this Task; `internal/init/v2_scaffold_test.go`; `internal/init/template_freshness_test.go`; `internal/init/agent_skills_test.go`; `main_test.go`; `main_resume_matrix_test.go`; `templates/project-v2/.savepoint/router.md`; new `templates/project-v2/.savepoint/releases/R-001-first-goal/Release.md`; and both Idea skill copies.

Limitations: owner validation is required. The full owner-side `savepoint init` → `savepoint resume` → `savepoint doctor` scenario in an empty folder remains to be run; the embedded init path and resume projection are covered by tests. Agent CLI policy permits invoking only `savepoint resume`, which was run after adding the Goal record.

## Drift Notes

None yet.
