---
id: T-039
objective: O-022
title: Give migrated projects a live Goal
status: done
complexity_tier: medium
complexity_reason: Changes migration output; golden fixtures and the router conversion must agree.
depends_on: []
owner_validation:
    required: true
    accepted_check: ""
planned_by: {role: planner, session: o022-goal-required-20260925}
check_waiver:
    task: T-039
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-25T02:09:31Z"
---

# Give migrated projects a live Goal

## Outcome

A V1 project converted by `savepoint migrate` has a Goal on every Objective and
a router that selects a live Goal. It retains a live Goal selected by the V1
router. If that selection is missing or unresolvable, it selects the live Goal
containing the router's active Objective, or the sole live Goal, when that
choice is clear. If selected active work belongs only to a historical Goal,
migrate creates a continuation and moves those active Objectives into it.
Unresolved release lifecycle decisions remain visible in the preview and
block Apply.

## User Check

Migrate a V1 fixture whose releases are all complete. Confirm migrate adds
one new live Goal, the router selects it, and resume and doctor report no
missing Goal. Migrate a fixture whose router names an in-progress release and
confirm that release's Goal is selected and no extra Goal is created. Migrate
fixtures with missing and unresolvable router releases where one existing live
Goal contains the active Objective; confirm that Goal is selected and no empty
continuation is generated. Migrate a fixture whose selected active Objective
belongs to a historical Goal; confirm the continuation contains that Objective
and the router resolves to it.

## Done When

- Every converted Objective carries a `release:` reference. If its V1 release
  has no mapped Goal and no unresolved lifecycle decision explains why, the
  preview fails with a named error instead of producing an Objective without
  a Goal. If the release lifecycle is unresolved, the preview shows that
  decision and does not plan the Objective until it is resolved.
- When the V1 router's release resolves to a live V2 Goal, the router selects
  it, as today.
- For a missing or unresolvable router selection, migration selects the live
  Goal containing the router's active Objective, or the sole live Goal, when
  that choice is clear; it does not generate an empty continuation in those
  cases.
- If selected active work belongs only to a historical Goal, migration plans
  a live continuation with the next free `R-###`, a title stating it continues
  work after migration, and stub sections. The router selects it, the active
  Objectives move into it, and the preview explains the decision. The old
  "the router now selects no Release" outcome is removed.
- A release with an unresolved lifecycle status keeps its decision in the
  preview and blocks Apply instead of failing with a missing-Goal error.
- The migrated project loads strictly with no T-036 router diagnostic.
- Golden and end-to-end tests cover all three cases.

## Context Files

`internal/migrate/convert_docs.go`, `internal/migrate/convert_docs_test.go`,
`internal/migrate/convert_releases.go`,
`internal/migrate/convert_releases_test.go`, `internal/migrate/plan.go`,
`internal/migrate/plan_test.go`, `internal/migrate/convert.go`,
`internal/migrate/end_to_end_test.go`, `internal/migrate/decisions.go`.

## Design References

Design sections 10 and 11; O-022 Architectural Considerations.

## Guardrails

FS-01, FS-03, DATA-01, DATA-02, DATA-04, TEST-01, TEST-02, TEST-03, TEST-05.

## Implementation Plan

1. Resolve the selected existing live Goal from the router and active work;
   plan a continuation only when the selected active work belongs to a
   historical Goal.
2. Plan the selection once so preview, Objective conversion, and router
   conversion share the same decision.
3. Enforce a Goal reference on every converted Objective.
4. Update the golden fixtures and add the three router cases.
5. Run focused tests while iterating, then `make test-full` (migration work).

## Boundaries

No changes to V1 Release archiving or legacy completion evidence, and no
changes to already-migrated V2 projects.

## Technical Verification

Focused tests during iteration; `make test-full` for handoff.

## Technical Evidence

Execution started at `stage: build`. Router selection already matched O-022/T-039; Task dependencies are empty and O-022 is in progress. Router `release: R-006` is unchanged.

Extra reads beyond Context Files:
- `.savepoint/Design.md`, sections 10 and 11, to check migration changes against Goal record storage and migration failure boundaries.
- `internal/migrate/apply.go`, `internal/migrate/preview.go`, and `internal/migrate/manifest.go`, to verify a generated Goal target is previewed and written without fabricating a V1 source mapping.
- `internal/migrate/testdata/golden/v1-basic.yml` and `internal/migrate/testdata/golden/v1-history.yml`, to check existing migration golden expectations before adding fallback-Goal coverage.
- `internal/migrate/convert_issues_test.go`, `internal/migrate/convert_test.go`, and `internal/migrate/decisions_test.go`, to locate hand-built active-epic fixtures that need valid V1 Goal inputs after enforcing Objective Goal references.
- `internal/migrate/testdata/golden/v1-router-missing.yml`, `v1-router-unresolvable.yml`, and `v1-router-archived.yml`, to review converted router and generated-Goal outputs for all fallback cases.
- `main_test.go` around the migrate ambiguity and decisions command fixtures, to provide a resolvable V1 Goal so those tests continue to isolate their intended ambiguity behavior.

Implementation evidence:
- `ConversionPlan.GoalSelection` is planned once after V1 releases and epics. A selected live Goal is retained; for a missing or unresolvable selection, a uniquely identifiable live Goal for the router's active Objective (or the sole live Goal) is reused. If selected active work belongs only to a historical Goal, migration creates an in-progress continuation using the next free R-### and moves those Objectives into it.
- `planEpic` defers Objective planning while its release has an unresolved lifecycle ambiguity; otherwise it refuses to create an Objective without a mapped Goal and names the V1 release and epic. `ConvertObjective` also rejects an empty planned Goal reference.
- The generated Goal uses the V2 Release renderer with stub sections. The converted router uses the same plan decision; the preview explains it. Generated Goals are omitted from the V1 identity map because they have no V1 source.
- Golden and end-to-end coverage includes retained live selection and missing, unresolvable, and archived router selections. The archived case has every source release complete.

Iteration evidence:
- `go test ./internal/migrate` — passed after adding the new fallback goldens and valid Goal sources to hand-built active-epic fixtures.
- `make test-focused TEST=Router|ContinuationGoal` — passed with `GOCACHE=/tmp/savepoint-t039-gocache`.
- The first `make test-full` run found two root CLI test fixtures missing a V1 Goal; their shared fixture now includes a live V1 Release PRD so ambiguity behavior remains the tested outcome.

Criterion evidence:
- Objective Goal reference and named missing-Goal failure: `TestPlan_activeEpicWithoutResolvableReleaseFailsWithNamedError`; the fallback end-to-end test also checks every converted Objective target references a Goal in the strict index.
- Existing live router selection is retained: `TestPlan_routerSelectsLiveGoalOrPlansContinuationGoal/live_selection_is_kept` and `TestConvertRouter_activeSelectionResolvesGlobalIDs`.
- Missing and unresolvable router releases select the existing live Goal containing converted active work; an archived selected Goal creates a continuation that contains the active Objective: `TestPlan_routerSelectsLiveGoalOrPlansContinuationGoal`, `TestConvertRouter_fallbackSelectionPointsAtLiveGoal`, `TestEndToEnd_routerFallbackCasesMigrateToLiveGoal`, and the three `v1-router-*.yml` goldens. The archived case marks the fixture's only source Goal complete.
- Strict V2 load and router Goal selection: `TestEndToEnd_routerFallbackCasesMigrateToLiveGoal` runs `data.LoadV2Index`, parses the converted router, and checks its Goal ID.

Full handoff gate:
- Command: `env GOCACHE=/tmp/savepoint-t039-gocache make test-full`
- Toolchain: Go 1.26.2, linux/amd64.
- Result: PASS on 2026-09-25 at 02:08 UTC; all Go packages passed and the gate built linux, darwin, and windows targets.
- This recorded Task gate predates the C-923 repairs. The migration plan,
  tests, and golden fixtures changed after it, so this evidence does not cover
  the repaired Objective; the fresh Full Objective Check must run
  `make test-full`.

## Drift Notes

None yet.
