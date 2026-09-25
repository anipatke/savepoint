---
id: T-036
objective: O-022
title: Flag work that has no Goal
status: done
complexity_tier: medium
complexity_reason: Adds one selection diagnostic kind and one non-fatal index fact, and renders both in resume and doctor.
depends_on: []
owner_validation:
    required: true
    accepted_check: ""
planned_by: {role: planner, session: o022-goal-required-20260925}
check_waiver:
    task: T-036
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-25T01:03:53Z"
---

# Flag work that has no Goal

## Outcome

A router with no Goal, or an Objective with no Goal, is reported plainly
instead of being silently treated as valid. Next reads `Choose a Goal` and
does no other work. Doctor names each missing reference and how to add it.

## User Check

In a scratch copy of a V2 project, set the router `release:` to `none`, then
run `savepoint resume` and `savepoint doctor`. Confirm the first resume line
reads `Choose a Goal`, with the fix (`g` on the board). Confirm doctor reports
the router problem without offering an unrelated record. Remove `release:`
from one Objective, restore the router, and confirm the project still loads,
resume flags it, and doctor names that Objective and the line to add.

## Done When

- `ResolveSelection` returns a new typed diagnostic (for example
  `SelectionReleaseMissing`) when the router Goal is missing, blank, or
  `none`. `ResolveNext` returns a Next that selects no Objective, Task, or
  Issue work because of it, whatever else the router names.
- `LoadV2Index` records every live Objective without `release:` as a typed,
  sorted, non-fatal fact on the index (for example
  `ObjectivesWithoutGoal`). Unknown or malformed references stay fatal, as
  today.
- `resume.NextLine` prints `Choose a Goal — <fix>` for the new diagnostic;
  `NextVerb` returns `Choose`. Resume's body flags Objectives without a Goal
  by ID.
- Doctor reports the router problem and each Objective without a Goal. Each
  report is a named problem with a concrete repair: choose a Goal with `g` on
  the board, or add `release: R-###` to the named Objective file. With zero
  Goals declared, the repair says to create one.
- AGENTS.md's Next-word list adds `Choose` → report to the owner, who decides.
  The scaffold copy matches byte-for-byte.
- Tests cover missing, blank, `none`, unknown, and archived router Goals;
  a router naming an Objective or Issue but no Goal; one and several
  Objectives without a Goal; and a fully valid project with no new output.

## Context Files

`internal/data/next.go`, `internal/data/next_test.go`,
`internal/data/project.go`, `internal/data/project_test.go`,
`internal/data/router_v2.go`, `internal/data/router_v2_test.go`,
`internal/resume/resume.go`, `internal/resume/resume_test.go`,
`internal/doctor/v2_runtime.go`, `internal/doctor/v2_runtime_test.go`,
`internal/doctor/repairs.go`, `internal/doctor/report.go`,
`AGENTS.md`, `templates/project-v2/AGENTS.md`.

## Design References

Design sections 4, 6, 8, 10, and 11; O-022 Architectural Considerations.

## Guardrails

DATA-01, DATA-02, DATA-03, ARCH-01, ARCH-03, TPL-01, TEST-01, TEST-02,
TEST-04, TEST-08, STYLE-07.

## Implementation Plan

1. Add the router diagnostic kind in `ResolveSelection`, checked before the
   Objective, Task, and Issue lookups. Make `resolveLadder` return no work for it.
2. Add the non-fatal Objectives-without-Goal fact in `indexReleaseObjectives`.
3. Render both in resume (`NextLine`, `NextVerb`, body) and in doctor.
4. Add `Choose` to the AGENTS.md Next-word routing, live and scaffold.
5. Run focused tests while iterating, then `make build && make test-fast`.

## Boundaries

No board changes (next Task), no init or migrate changes, no automatic Goal
assignment, and no inference of a Goal from the selected Objective.

## Technical Verification

Focused tests during iteration; `make build && make test-fast` for handoff.

## Technical Evidence

Execution started 2026-09-25. Implementation and criterion evidence:

- Missing Goal selection is a typed `SelectionReleaseMissing` diagnostic, and
  resolves no Objective, Task, or Issue even when the router names those
  records. `TestResolveSelection_missingGoalPreemptsOtherSelections` covers
  an empty selection, Objective plus Task, Issue, and Objective plus Issue.
- `LoadV2Index` records live unassigned Objectives in sorted ID order and
  keeps malformed or unknown Goal references fatal. Evidence:
  `TestLoadV2Index_recordsObjectivesWithoutGoalInIDOrder`,
  `TestLoadV2Index_validGoalReferenceHasNoMissingGoalFact`, and
  `TestLoadV2Index_invalidGoalReferencesRemainFatal`.
- Resume renders `Choose a Goal — press g on the board`, uses `Choose`, names
  the router issue, and lists unassigned Objective IDs. Valid Next values add
  no missing-Goal text. Evidence: `TestNextLine_missingGoalAndObjectiveFacts`,
  `TestRender_validNextHasNoMissingGoalOutput`, and
  `TestRender_flagsUnassignedObjectivesWithGoalSelected`.
- Doctor reports the router and each missing Objective reference with a
  concrete repair, switches to “Create a Goal” when none exist, and stays
  quiet for a valid project. Evidence:
  `TestRunV2ChecksReportsMissingGoalsAndConcreteRepairs` and
  `TestRunV2ChecksFullyValidGoalProjectHasNoMissingGoalOutput`.
- `Choose` is in the Next-word routing list in `AGENTS.md` and the scaffold
  copy (the routing lines match). Resolver coverage includes absent, blank, and
  `none` Goal values,
  missing/archived references, several unassigned Objectives, and valid
  projects without new output.
- Focused iteration passed:
  `GOCACHE=/tmp/savepoint-t036-go-cache go test ./internal/data ./internal/resume ./internal/doctor`.
  The required handoff gate passed at 2026-09-25 00:49 UTC with Go 1.26.2
  linux/amd64: `GOCACHE=/tmp/savepoint-t036-go-cache make build &&
  GOCACHE=/tmp/savepoint-t036-go-cache make test-fast` (exit 0). The first
  fast-gate run identified downstream fixtures that assumed a Goal was
  optional; those test fixtures now provide valid Goal context, while the
  scaffold test expects the new Choose action. `git diff --check` passed.
- After setting `stage: audit`, `./savepoint resume` returned
  `Check T-036 — Flag work that has no Goal (O-022)` and the owner action to
  request the optional Task Check or record an explicit waiver.

Extra reads beyond Context Files:

- `.savepoint/Guardrails.md` — apply the Task-named rules.
- `internal/data/dependency_test.go` — locate the shared in-memory resolver fixture so existing resolver tests can use valid Goal context.
- `internal/resume/evidence.go` — add wording for the new selection diagnostic, which is rendered there.
- `internal/doctor/report_test.go` — update the shared complete-project fixture so valid-project doctor tests include a Goal.
- `internal/doctor/checks_test.go` — update default doctor fixtures for the required Goal reference and preserve the intentional unassigned Objective case.
- `main_resume_test.go`, `main_resume_matrix_test.go`, `main_board_test.go`, `main_board_next_parity_test.go`, and `internal/board/board_test.go` — inspect and correct cross-surface integration fixtures after the required fast gate exposed Goal-less routers.
- `internal/board/v2/fixture_test.go`, `actions_test.go`, `card_test.go`, `view_test.go`, `objectives_test.go`, `footer_test.go`, `load_test.go`, `next_panel_test.go`, `watch_test.go`, and `releases_test.go` — preserve the intent of existing board tests by supplying valid Goal context to ordinary work fixtures and retaining the unconfigured scaffold case. Changes were confined to tests; no board implementation changed.

## Drift Notes

None yet.
