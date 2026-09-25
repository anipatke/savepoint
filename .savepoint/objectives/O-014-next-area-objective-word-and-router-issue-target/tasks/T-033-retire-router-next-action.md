---
id: T-033
title: Stop the router carrying its own next-step text
objective: O-014
planned_by: {role: planner, session: o014-design-20260924}
status: done
complexity_tier: medium
complexity_reason: "Removes a decoded field while keeping old routers loadable, adds a doctor retired-field report, and must confirm every hand-off the prose carried is expressed by Next."
depends_on: [{task: T-032, requires: clear}]
owner_validation: {required: false}
check_waiver:
    task: T-033
    reason: "Waived - deferred to O-014 Full Objective Check"
    actor:
        role: owner
        session: owner-chat
    recorded_at: "2026-09-24T08:29:30Z"
---

# T-033: Stop the router carrying its own next-step text

## Outcome

The router holds only `state` and its Goal/Objective/Task/Issue selection.
An old router with `next_action` still loads; doctor reports the field once
as retired, and no surface shows it.

## User Check

Run `savepoint doctor` on this repository before the router is cleaned: it
reports `next_action` as retired. After the cleanup the report is gone and
the board and resume are unchanged.

## Done When

- `RouterStateV2.NextAction` is removed; `ReadStateV2` still accepts a
  `next_action` key (no KnownFields failure) and records only that it was
  present, for doctor.
- Doctor reports a retired `next_action` once, with a repair hint to delete
  the line; it never prints the value.
- `WriteRouterStateV2` leaves an existing `next_action` line byte-identical
  (removal is the owner's or skill's edit) and never adds one.
- `templates/project-v2/.savepoint/router.md` and this repository's
  `.savepoint/router.md` no longer carry `next_action`.
- Evidence lists each hand-off `next_action` expressed in this repo's recent
  router history (for example "request a Task Check or record a waiver",
  "select the next Objective") and the `data.Next` rung or resume phrase that
  now states it; any gap gets a resume phrase, not free text.
- Tests: decode with and without the key, doctor report, writer
  preservation, template freshness.
- `git diff --check`, `make build && make test-fast` pass.

## Context Files

`internal/data/router_v2.go`, `internal/data/router_v2_test.go`,
`internal/data/write.go`, `internal/data/write_test.go`,
`internal/doctor/v2_runtime.go`, `internal/doctor/v2_runtime_test.go`,
`internal/resume/resume.go`, `templates/project-v2/.savepoint/router.md`,
`.savepoint/router.md`, `internal/init/template_freshness_test.go`.

## Design References

Design sections 1, 6, and 11.

## Guardrails

DATA-01, DATA-03, DATA-04, TPL-02, TPL-04, TEST-02, TEST-03, TEST-08.

## Implementation Plan

1. Replace the field with a retired-key presence flag.
2. Add the doctor report and its test.
3. Update the template and live router; run template tests.
4. Record the hand-off coverage table; add phrases only for gaps.

## Boundaries

No guidance rewrite (T-034). The legacy V1 router reader used by migrate is
untouched.

## Technical Verification

Focused data/doctor/init tests during iteration; `make build && make
test-fast` at handoff.

## Technical Evidence

Extra reads before implementation:

- `.savepoint/objectives/O-014-next-area-objective-word-and-router-issue-target/tasks/T-032-advance-the-router-on-task-closure.md` — confirmed the predecessor's recorded `done` status and owner Task-check waiver; the runtime resolver still decides whether `requires: clear` is met.
- `.savepoint/Guardrails.md` — read the Task's named guardrails and applicable style rules before code changes.
- `go doc` for `LoadV2Index`, `ResolveTaskStart`, `ResolveTaskDependencyV2`, `ResolveObjectiveDependency`, `GateDecision`, `DependencyDecision`, `V2Index`, and `ObjectiveV2` — identified the runtime APIs and result types needed for a start preflight.
- Targeted `rg --files .savepoint` lookup for O-014 and T-033 — located the specifically requested Task and its owning Objective.
- `git status --short` — identified existing uncommitted T-030/T-031/T-032 and O-014 changes to preserve while editing overlapping files.
- A targeted `git diff` of the already-modified T-033 context files — distinguish the existing T-030/T-031/T-032 changes from T-033 edits before modifying shared files.
- A temporary in-repository Go helper loaded `.savepoint` and invoked the named runtime gate resolvers — required because the router selection was stale, and the repository forbids running Savepoint CLI commands.
- Targeted `git log -p` for `.savepoint/router.md` — recover recent committed `next_action` hand-offs so each can be mapped to existing Next/resume behavior.
- Targeted read of `internal/data/next.go` — identify the typed Next rung for each recovered hand-off; the Task requires evidence for any hand-off whose instruction is otherwise unexpressed.
- Targeted repository-wide Go search for `NextAction` and `next_action` — find every remaining consumer before removing the decoded value from `RouterStateV2`.
- Targeted search in `internal/resume/resume_test.go` — inspect the current assertions for the `NextCheckNeeded` hand-off before changing its shared ActionPhrase.
- Targeted read of `main_resume_matrix_test.go` — the top-level gate found an assertion still expecting the old `NextCheckNeeded` phrase; update this cross-surface expectation to the replacement owner choice.

Preflight:

- `go run tmp_task033_preflight.go` — `ResolveTaskStart(T-033)` allowed; T-032's `requires: clear` dependency was satisfied by its recorded owner Task-check waiver; O-014's O-012 Objective dependency was satisfied. The temporary helper was removed immediately after the run.

Handoff coverage:

| Recent router hand-off intent | Next/resume behavior that carries it |
| --- | --- |
| Select the next Objective after an Objective closes | `NextNothingSelected` renders `Nothing selected`; `ActionPhrase` says to select an Objective and how. |
| Start or continue the selected Task, then run its Task gate | `NextExecute` identifies the selected Task and `ActionPhrase` says to start it or advance its stage; the Task's Technical Verification and `savepoint-task` workflow retain the exact gate command. |
| At Task audit, request the optional Task Check or record an owner waiver; still run the mandatory Full Objective Check | `NextCheckNeeded` now states the optional Check/waiver choice and the mandatory Objective Check. The previous generic “Record a fresh Check” phrase left this choice implicit. |
| After all Tasks, record the mandatory Objective integration Check and then close it as owner | `NextObjectiveIntegration`/`NextObjectiveReady` render the Check rung and the existing integration/owner phrases. |
| Run a fresh full gate before the mandatory Goal Check | `NextReleaseCheckNeeded` asks for a fresh Goal Check; the Verification Policy and `savepoint-check` procedure require `make test-full`. |

Implementation evidence:

- `RouterStateV2` now exposes only `HasRetiredNextAction`; `ReadStateV2` recognizes the legacy key under `KnownFields(true)`, including an empty value, and does not expose its contents. `TestReadStateV2_decodesSelectedTask`, `TestReadStateV2_recordsEmptyRetiredNextActionPresence`, and `TestReadStateV2_decodesIssueAloneWithoutObjective` cover present, empty, and absent keys.
- Doctor emits one pending-review finding for the key, with a delete-line repair hint and no value. `TestRunV2ChecksReportsRetiredNextActionWithoutItsValue` checks one warning, value redaction, and disappearance after cleanup.
- `WriteRouterStateV2` still changes only selection fields. `TestWriteRouterStateV2_preservesRetiredNextActionLineByteForByte` and `TestWriteRouterStateV2_doesNotAddRetiredNextAction` cover preservation and non-creation.
- Both router files omit the key and say the router records selection; `TestRouterFilesOmitRetiredNextActionKey` covers both files.
- The `NextCheckNeeded` ActionPhrase now gives the owner the optional Task Check/waiver choice and keeps the mandatory Full Objective Check explicit. `TestRender_checkNeeded` and `TestActionPhraseSelectionGuidance` cover the phrase. `NextNothingSelected` already says how to select an Objective.
- The targeted production-code search found no V2 consumer that renders the retired value. V1 parsing and migration preservation remain untouched, per the Task boundary.

Focused iteration:

- The first `go test ./internal/data ./internal/doctor ./internal/init ./internal/resume -count=1` run exposed empty-value presence handling and two fixture/assertion mistakes; those were corrected.
- `go test ./internal/doctor -run TestRunV2ChecksReportsRetiredNextActionWithoutItsValue -count=1` exposed the fixture's existing legacy key and a case-sensitive repair assertion; both were corrected.
- Final `go test ./internal/data ./internal/doctor ./internal/init ./internal/resume -count=1` — all four packages passed.
- The first `make build && make test-fast` passed the build but found a top-level resume-matrix expectation for the old Task Check phrase. `go test .` reproduced it in `TestResumeMatrix_everyRungReachedExactlyOnce/task_at_audit_with_no_recorded_check_needs_one`; `main_resume_matrix_test.go` now expects the replacement phrase.

Handoff gates:

- `go test ./internal/data ./internal/doctor ./internal/init ./internal/resume -count=1` — passed.
- `make build && make test-fast` — passed after updating the top-level resume-matrix expectation.
- `git diff --check` — passed; a targeted check found no `next_action:` key in the live or scaffold router.
- Cross-surface tests in the full fast gate passed, including `TestBoardNextAndResumeReportTheSameAnswer`, `TestBuiltBoardAndResumeReportTheSameAnswer`, and `TestResumeMatrix_everyRungReachedExactlyOnce`.

Files read:

- Routing/workflow records: `.savepoint/router.md`, `agent-skills/savepoint-task/SKILL.md`, the O-014 Objective, this T-033 record, and predecessor T-032.
- Task Context Files: `internal/data/router_v2.go`, `internal/data/router_v2_test.go`, `internal/data/write.go`, `internal/data/write_test.go`, `internal/doctor/v2_runtime.go`, `internal/doctor/v2_runtime_test.go`, `internal/resume/resume.go`, `templates/project-v2/.savepoint/router.md`, `.savepoint/router.md`, and `internal/init/template_freshness_test.go`.
- Extra reads and reasons are listed above; follow-up verification reads included `internal/resume/resume_test.go` and `main_resume_matrix_test.go` for the shared NextCheckNeeded phrase.

Files changed for T-033:

- `.savepoint/objectives/O-014-next-area-objective-word-and-router-issue-target/tasks/T-033-retire-router-next-action.md`, `.savepoint/router.md`, and `templates/project-v2/.savepoint/router.md`.
- `internal/data/router_v2.go`, `internal/data/router_v2_test.go`, `internal/data/write.go`, `internal/data/write_test.go`, `internal/doctor/v2_runtime.go`, `internal/doctor/v2_runtime_test.go`, `internal/init/template_freshness_test.go`, `internal/resume/resume.go`, `internal/resume/resume_test.go`, and `main_resume_matrix_test.go`.
- The worktree had existing T-030/T-031/T-032 changes in several shared files. Those changes were preserved; only the T-033 additions described above were made to the overlapping files.

Limitations and handoff:

- The live `savepoint doctor` User Check was not run because the active AGENTS.md rule says never to run Savepoint CLI commands. `TestRunV2ChecksReportsRetiredNextActionWithoutItsValue` exercises the runtime doctor path on a temporary V2 project before and after removing the key.
- No Full Objective Check was run; O-014 still requires current `make test-full` evidence before it can close. T-032's owner waiver satisfies T-033's `requires: clear` dependency but is not technical CLEAR. The owner explicitly waived T-033's optional Task Check to defer review to O-014's mandatory Full Objective Check; this waiver is not technical CLEAR and does not replace that Objective Check.
- `go run tmp_task033_validate.go` — as a post-waiver validation, `LoadV2Index(".savepoint")` succeeded with T-033 recorded as `in_progress`/`audit`; the temporary helper was removed immediately after the run.

## Drift Notes

None expected.
