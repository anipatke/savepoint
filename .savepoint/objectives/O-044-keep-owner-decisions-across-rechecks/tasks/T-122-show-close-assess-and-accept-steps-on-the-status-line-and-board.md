---
id: T-122
title: Show close, assess and accept steps on the status line and board
objective: O-044
status: in_progress
stage: audit
depends_on: [{task: T-121, requires: clear}]
owner_validation: {required: true}
planned_by: {role: planner, session: planning-o044-20261010}
---

# Show close, assess and accept steps on the status line and board

## Outcome

`savepoint resume` and the board name the right next step from the gate result: `Close` with "ready to close by exception" when proof and still-applicable owner decisions satisfy completion; `Assess` when a checker still needs to confirm an older decision for the latest Check; `Accept` for only the decision a material change affected, naming the change. Technical clearance is shown separately and still reads NEEDS WORK when evidence is waived. Accepting on the board adds an owner renewal entry and keeps the originating Check.

## User Check

On a temporary copy of the TheShed-shaped fixture, open the board and run `savepoint resume` at three points: with carry entries (expect `Close O-…` and "ready to close by exception", clearance NEEDS WORK), without them (expect `Assess O-…` naming C-006, not `Check`), and with one decision marked changed (expect `Accept` naming that decision and the change). On a Task needing renewal, press `a` and confirm the file keeps `accepted_check` and gains an owner `carried_forward` entry.

## Done When

1. Next projection: unassessed blockers map to a new assessment rung whose verb is `Assess`; materially changed blockers map to the owner rung (`Accept`); allowed-by-exception maps to `Close`. The board and resume read the same projection and the same wording (STYLE-07, STYLE-09).
2. Evidence lines: "Technical clearance" and "Completion" are separate lines; by-exception completion reads "Ready to close by exception" and names the originating Check and the Check it was carried to, with the assessor and reason. Assess and Accept lines name the decision, the Checks involved, and the material change or uncovered requirement IDs. Authored text is sanitised (ARCH-05).
3. The action phrase for `Assess` says to confirm whether the recorded decision still applies, and that it is not a new Check.
4. Board accept (`a`) is offered for current-clearance acceptance blocks and for changed or unassessed acceptance; it appends an owner `carried_forward` entry at the latest Check when an originating Check exists, sets `accepted_check` only when none does, and never removes earlier entries (FS-01). Closing by exception (`x`) works for carried exceptions. The board never closes anything on its own.
5. Detail view shows each decision's originating Check, scope or requirements, and every carry entry.
6. Resume and board tests cover the three routes, the renewal write, and a failure path (renewal refused when the gate no longer offers it) (TEST-01, TEST-02, TEST-03); exact-text tests are updated, not deleted. `make build && make test-fast` passes (TEST-08).
7. The owner reviews the resume output and board screens for the three routes and records acceptance. Per-criterion evidence follows AGENTS.md's Verification Policy.

## Context Files

`internal/data/next.go`, `internal/data/next_test.go`, `internal/resume/resume.go`, `internal/resume/evidence.go`, `internal/resume/resume_test.go`, `internal/board/v2/actions.go`, `internal/board/v2/actions_test.go`, `internal/board/v2/io.go`, `internal/board/v2/badges.go`, `internal/board/v2/detail_view.go`, `internal/board/v2/detail_test.go`, `internal/board/v2/next_panel.go`, `internal/board/v2/next_panel_test.go`, `main_resume_matrix_test.go`, `main_board_next_parity_test.go`, `.savepoint/objectives/O-044-keep-owner-decisions-across-rechecks/Objective.md`.

## Design References

O-044 Confirmed Design (unassessed decisions; routing). I-141 Planned Fix, Routing.

## Guardrails

ARCH-02, ARCH-05, FS-01, STYLE-07, STYLE-09, TEST-01, TEST-02, TEST-03, TEST-08.

## Implementation Plan

1. Confirm T-121's blockers and the exposed carry entry exist; return REPLAN REQUIRED if not.
2. Add the assessment rung to `next.go` and map the new blocker kinds in task and Objective routing.
3. Add the `Assess` verb, evidence lines and action phrases in `internal/resume`, keeping copy in data.
4. Change the board accept writer to append a renewal entry; update actions, detail and Next panel.
5. Update and add resume, board and parity tests.

## Boundaries

No gate rule changes (T-121 owns them) and no skill or guide changes (a later O-044 Task owns them). No new board keys beyond the existing `a`, `x` and Space. Do not change resume's overall output layout beyond the lines named here.

## Technical Verification

Focused `make test-focused TEST=...` during iteration; `make build && make test-fast` for handoff. Owner validation of the three routes is recorded before closure. The Full Objective Check needs fresh `make test-full` under `agent-skills/references/check-method.md`.

## Technical Evidence

Gate: `make build && make test-fast` passed (2026-10-10).

1. `NextAssessDecision` rung (verb `Assess`) added in `internal/data/next.go`; changed decisions map to the owner rung (`Accept`); by-exception completion stays `Close`. Task and Objective routing share `decisionRung`.
2. `internal/resume/evidence.go` adds `CompletionLines`, `CarryPhrase`, `DecisionBlockerLines`: clearance and completion are separate lines; text passes `cleanText`.
3. `ActionPhrase` for Assess says to confirm the decision still applies and that it is not a new Check.
4. Board `a` is offered for current-clearance acceptance and for unassessed/changed acceptance; it appends an owner `carried_forward` entry and keeps `accepted_check` (`io.go`, `actions.go`, `badges.go`). `x` is unchanged.
5. Detail view lists scope and every carry entry for acceptance and exception.
6. Tests: `TestResolveNext_unassessedAndChangedAcceptanceRoute`, `TestRender_assessAndAcceptRoutes`, `TestRender_closeByExceptionNamesCarry`, `TestOwnerAcceptanceRenewalAppendsCarryAndKeepsOrigin` (includes the refusal); two exact-text tests updated.

Not done: criterion 7 (owner review of the three routes on the TheShed-shaped fixture) awaits the owner. No extra reads outside Context Files other than `internal/resume/concurrency.go` (for `cleanText`) and `fixture_test.go`.

## Drift Notes

None yet.
