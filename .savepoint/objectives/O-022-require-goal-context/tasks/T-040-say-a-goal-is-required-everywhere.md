---
id: T-040
objective: O-022
title: Say a Goal is required everywhere
status: done
complexity_tier: low
complexity_reason: Guidance and Design reconciliation across live and scaffold copies, after the behaviour lands.
depends_on: [{task: T-036, requires: clear}, {task: T-037, requires: clear}, {task: T-038, requires: clear}, {task: T-039, requires: clear}]
owner_validation:
    required: true
    accepted_check: ""
planned_by: {role: planner, session: o022-goal-required-20260925}
check_waiver:
    task: T-040
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-25T03:21:28Z"
---

# Say a Goal is required everywhere

## Outcome

Every piece of guidance an agent or owner reads says a Goal is required and
matches the shipped behaviour. No "a Goal is optional" wording remains in the
active or scaffolded guidance.

## User Check

Read AGENTS.md's Goal section and the Idea and Design skills' Goal
sections. Confirm each says every project has a selected Goal and every
Objective belongs to one. Confirm each says what happens when a Goal is
missing: `Choose a Goal`, doctor's fix, the init default, and migrate's live
Goal.

## Done When

- AGENTS.md "Optional Goals" becomes a required-Goal section. The Goal Check
  sentence no longer says "whenever a Goal exists". The scaffold copy is
  byte-identical.
- The design skill's "Optional Goal Boundary" and its Objective template note
  ("Omit `release` when an Objective is intentionally unassigned") are
  replaced. The same goes for the idea skill's "Goal is optional" paragraph
  and any matching wording in the task and check skills. Live and scaffold
  copies are byte-identical.
- The `.savepoint/router.md` lifecycle bullet and the scaffold router say a
  Goal is required.
- Design.md sections 1, 3, 8, and 10 describe the implemented rule: the typed
  router diagnostic, the non-fatal Objective fact, the Goal-scoped board, and
  the init and migration Goal-selection defaults.
- A test fails if "Goal is optional" or equivalent wording reappears in the
  live or scaffold guidance.
- `git diff --check` and `make build && make test-fast` pass.

## Context Files

`AGENTS.md`, `templates/project-v2/AGENTS.md`, `.savepoint/Design.md`,
`.savepoint/router.md`, `templates/project-v2/.savepoint/router.md`,
`agent-skills/savepoint-idea/SKILL.md`,
`templates/project-v2/agent-skills/savepoint-idea/SKILL.md`,
`agent-skills/savepoint-design/SKILL.md`,
`templates/project-v2/agent-skills/savepoint-design/SKILL.md`,
`agent-skills/savepoint-task/SKILL.md`,
`templates/project-v2/agent-skills/savepoint-task/SKILL.md`,
`agent-skills/savepoint-check/SKILL.md`,
`templates/project-v2/agent-skills/savepoint-check/SKILL.md`,
`internal/init/agent_skills_test.go`,
`internal/init/template_freshness_test.go`, `README.md`.

## Design References

Design sections 1, 3, 8, and 10; O-022 Success Conditions.

## Guardrails

TPL-01, TPL-02, TPL-04, TEST-01, TEST-05.

## Implementation Plan

1. Confirm T-036 to T-039 landed as planned; record any drift for the planner.
2. Rewrite the Goal sections in AGENTS.md and the skills, live and scaffold.
3. Reconcile Design.md and the router lifecycle bullet.
4. Add the regression test against optional-Goal wording.
5. Run `git diff --check` and `make build && make test-fast`.

## Boundaries

Guidance and Design only; no behaviour changes. Do not rename storage
fields.

## Technical Verification

Focused tests during iteration; `make build && make test-fast` for handoff.
The mandatory Full Objective Check then runs `make test-full`.

## Technical Evidence

Started T-040 after the runtime reported its dependencies satisfied. The
router selects O-022/T-040 and keeps `release: R-006`; the owning Objective is
`in_progress`. Required workflow reads outside this Task's Context Files:
`.savepoint/router.md`, the owning Objective record, `.savepoint/Guardrails.md`,
and `agent-skills/savepoint-task/SKILL.md`.

Plan-required prior Task records were read before editing. T-036 through T-039
are all `done`, with evidence recording the planned outcomes: typed missing
router-Goal diagnostics and non-fatal missing-Objective facts; a Goal-scoped
board with no project-wide fallback; init's selected R-001; and migration's
live Goal selection or generated continuation Goal. Their records report the
required gates passing and no drift. The recorded evidence names the relevant
tests and implementation results, so no additional implementation source was
read for this confirmation.

Before adjusting the adoption regression surfaced by `make test-fast`, also
recording the required extra read of `internal/init/v2_scaffold_test.go`; its
existing assertion treats an absent Goal as an ordinary optional file.

Per-criterion evidence:

1. `AGENTS.md` now requires a router-selected live Goal and a Goal reference
   on every Objective; both Goal Check statements are unconditional. The
   live and scaffold Goal sections match byte-for-byte. Evidence:
   `TestProjectGuidanceRequiresGoalContext`.
2. Idea, Design, Task, and Check guidance now states the required Goal context,
   missing-Goal repair, init default, and migration behavior. Design requires
   `release: R-###` on every Objective. Evidence:
   `TestV2SkillsTeachProjectGoalWorkflow`,
   `TestProjectGuidanceTemplatesMirrorLiveGuidance`, and
   `TestProjectGuidanceRequiresGoalContext`.
3. The live and scaffold router lifecycle guidance requires a selected Goal
   and Objective membership. `TestProjectGuidanceRequiresGoalContext` checks
   both router files for retired optional-Goal wording.
4. Design sections 1, 3, 8, and 10 now record `SelectionReleaseMissing`, the
   non-fatal `ObjectivesWithoutGoal` fact, Goal-scoped board behavior, and the
   R-001 init and migration Goal-selection defaults. Section 2's directory
   description was reconciled as well.
5. `TestProjectGuidanceRequiresGoalContext` scans the live and scaffold
   guidance plus README and Design for optional or unassigned Goal wording,
   and pins the missing-Goal repair/default instructions. The existing
   adoption test now treats Concept, Health-Check, and procedures as optional
   while requiring the Goal repair guidance:
   `TestV2AgentsGuideCarriesExistingCodebaseAdoptionSection`.
6. Final handoff verification passed: `git diff --check` and
   `make build && make test-fast`. Focused evidence also passed:
   `go test ./internal/init -count=1` and
   `go test ./internal/init -run 'TestProjectGuidanceRequiresGoalContext|TestV2SkillsTeachProjectGoalWorkflow|TestV2AgentsGuideCarriesExistingCodebaseAdoptionSection' -count=1`.
   Toolchain: Go 1.26.2, linux/amd64. Final gate result: PASS on
   2026-09-25 UTC.

Iteration note: earlier gate runs exposed stale assertions for optional Goal
wording and the previous Idea sentence. Those assertions were updated; the
final focused package and handoff gate passed.

Extra reads beyond Context Files: `.savepoint/router.md`, the owning O-022
`Objective.md`, `.savepoint/Guardrails.md`,
`agent-skills/savepoint-task/SKILL.md`, the O-022 Task records T-036 through
T-039 (to confirm their recorded outcomes), and
`internal/init/v2_scaffold_test.go` (to reconcile its existing adoption test).

Owner validation is required. No independent Task Check, waiver, or owner
acceptance is recorded; the Task remains in `stage: audit` for the owner.

## Drift Notes

After C-923 identified I-052, the owner directed migration to reuse an
existing live Goal when its association with active work is clear. T-039,
Design §10, AGENTS, README, and the live/scaffold skills now describe that
fallback and retain a continuation Goal for selected work whose source Goal
is historical. This supersedes the continuation-only fallback described in
the earlier evidence above. C-924 independent recheck is pending.
