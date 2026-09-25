---
id: C-924
scope: {kind: objective, id: O-022}
result: CLEAR
checked_by: {role: checker, session: o022-independent-recheck-20260925}
executed_session: o022-repair-followup-20260925
checked_at: '2026-09-25T04:53:45Z'
reviewed:
  base_commit: 35ce82dffbcfa02a88121f99496d668c67d69db9
  files:
    - .savepoint/Design.md
    - .savepoint/objectives/O-022-require-goal-context/Objective.md
    - .savepoint/objectives/O-022-require-goal-context/tasks/T-036-report-a-missing-router-goal.md
    - .savepoint/objectives/O-022-require-goal-context/tasks/T-037-keep-the-board-inside-one-goal.md
    - .savepoint/objectives/O-022-require-goal-context/tasks/T-038-create-a-goal-on-init.md
    - .savepoint/objectives/O-022-require-goal-context/tasks/T-039-give-migrated-projects-a-live-goal.md
    - .savepoint/objectives/O-022-require-goal-context/tasks/T-040-say-a-goal-is-required-everywhere.md
    - internal/data/project.go
    - internal/doctor/checks.go
    - internal/doctor/v2_runtime.go
    - internal/board/v2/next_panel.go
    - internal/migrate/plan.go
    - internal/migrate/preview.go
    - internal/migrate/convert_docs.go
    - internal/migrate/convert_docs_test.go
    - internal/migrate/end_to_end_test.go
    - internal/doctor/v2_runtime_test.go
    - internal/doctor/checks_test.go
    - internal/board/v2/releases_test.go
    - internal/init/v2_scaffold_test.go
    - AGENTS.md
    - README.md
    - templates/project-v2/AGENTS.md
    - agent-skills/savepoint-check/SKILL.md
    - agent-skills/savepoint-design/SKILL.md
    - agent-skills/savepoint-idea/SKILL.md
    - agent-skills/savepoint-task/SKILL.md
    - templates/project-v2/agent-skills/savepoint-check/SKILL.md
    - templates/project-v2/agent-skills/savepoint-design/SKILL.md
    - templates/project-v2/agent-skills/savepoint-idea/SKILL.md
    - templates/project-v2/agent-skills/savepoint-task/SKILL.md
  dependencies: []
issues: []
supersedes: C-923
---

# C-924: O-022 Full Objective Recheck

## Recheck Convergence and Closure Map

C-923 remains immutable with result NEEDS WORK. This independent Check
rechecks its four findings against repaired code, regression coverage, a fresh
full gate, and the reconciled Objective, Design, and guidance. Each is now
proven repaired and is closed as verified below.

| Issue | C-923 finding | C-924 evidence | Disposition |
| --- | --- | --- | --- |
| I-050 | With no live Goals and a stale router selection, the board sent the owner to doctor while doctor reported all clean. | The board and doctor share V2Index.HasLiveGoal; unknown and archived selections with no live Goal produce doctor guidance to create one. Regressions cover both cases and the board pointer. | Verified |
| I-051 | Fresh init and generated migration Goals with no Objectives failed doctor before the owner could add work. | An empty in-progress Goal is accepted by doctor; the completion gate still prevents completing an empty Goal. Fresh scaffold and migration checks exercise doctor. | Verified |
| I-052 | Missing/unresolvable router fallback selected an empty continuation and hid active work; archived source work was not selected. | Per the owner’s choice, fallback reuses the live Goal associated with the router’s active Objective, or the sole live Goal when clear. If selected active work belongs only to a historical Goal, migration creates a continuation and moves that Objective into it. End-to-end tests resolve Next and verify Goal-scoped work. | Verified |
| I-053 | An audited or unrecognized release status with active work failed preview with a misleading no-resolvable-Goal error instead of showing the lifecycle decision. | Preview reports the unresolved lifecycle decision and blocks Apply. Audited and unrecognized active-work cases were checked; supplying the decision restores planning. | Verified |

## Independence and Scope

This checker session is independent from executor session
o022-repair-followup-20260925. The work remains uncommitted on branch v2,
based on 35ce82d; the working tree also contains O-019 changes. I reviewed
the repaired O-022 behavior and final contract/test edits without changing
implementation or Task completion.

All five Tasks T-036 through T-040 remain done with owner Task-check waivers.
The waivers are not technical clearance. Each Task, the Objective outcome,
cross-Task integration, and Design reconciliation were reviewed for this
mandatory Full Objective Check.

The scope lock is inherited unchanged from C-923: O-022 Success Conditions
1–6 and Outcome; T-036 through T-040 Done When; guardrails DATA-01..04,
ARCH-01/03, TPL-01..04, FS-01/03, and TEST-01..05/08; entry points
data.ResolveSelection/ResolveNext, data.LoadV2Index, resume, doctor, the TUI
and plain board, init, and migration planning/conversion. The locked axes
remain router Goal {absent, blank, none, unknown, archived, live} × other
selections {none, Objective, Objective+Task, Issue only}; Goal inventory
{zero, archived-only, live}; Objective release {present, missing, one,
several, malformed/unknown}; surfaces {resume, board TUI, board plain,
doctor}; origin {fresh init, migrate with live/missing/unresolvable/archived
router Goal, release lifecycle ambiguity}; and guidance {live, scaffold}.
No additional scope was introduced.

The owner selected the I-052 behavior: reuse an identifiable existing live
Goal for active work; create and select a continuation only when the selected
active work belongs to a historical Goal, moving that work into it. An
unresolved release lifecycle choice remains in preview and blocks Apply. O-022
Success Condition 5, T-039, Design §10, README, AGENTS, and live/scaffold
guidance now record that same decision. T-040 Drift Notes records this
reconciliation.

## Coverage Matrix

| Slice | Probe or evidence | Result |
| --- | --- | --- |
| Router Goal absent, blank, or none; no work selection | C-923 scratch-project scenarios and current ResolveSelection/ResolveNext tests; current full suite. Resume, board, and doctor retain the same Choose a Goal diagnostic and do not substitute work. | Pass |
| Router Goal absent with Objective, Objective+Task, or Issue-only selection | C-923 scratch-project scenarios and current selection tests. The diagnostic remains selected and no unrelated work is substituted. | Pass |
| No live Goal while router names a deleted, unknown, or archived Goal | TestRunV2ChecksReportsNoLiveGoalForUnknownAndArchivedSelections, TestBoardWithZeroGoalsPointsToDoctor, and TestBoardWithOnlyArchivedGoalsPointsToDoctor; shared HasLiveGoal inspected in board and doctor. Doctor says to create a Goal and the board pointer reaches that repair. | Pass; closes I-050 |
| Objective missing release | C-923 scratch-project scenario plus current data/resume/board/doctor regressions. The Objective remains loadable; resume and board flag it; doctor identifies the Objective and exact release: R-### fix. | Pass |
| Objective release malformed or unknown | TestLoadV2Index_invalidGoalReferencesRemainFatal; remains a named load error rather than a missing-reference default. | Pass |
| Goal-scoped board with no valid selection, with and without Objective filter | Board TUI/plain tests and TestGoalScopedSelectionWithObjectiveFilter; no project-wide Objective or Task fallback. Selected-Goal labels remain Goal-scoped, and TestGoalSelectorCannotClearTheCurrentGoal pins the selector rule. | Pass |
| Zero Goals and archived-only Goals | Shared predicate, I-050 regressions, and C-923 scenario for the none router case. Empty project points from board to doctor; doctor gives create-Goal repair for unknown or archived router selection. | Pass |
| Fresh init / scaffold Goal | TestV2ScaffoldLoadsCleanThroughRuntimeSchemaAndIndex verifies R-001, strict loading, and no doctor problem for an empty in-progress Goal. Existing init scenario verifies R-001 is titled after the project and selected by the router. | Pass |
| Empty Goal lifecycle | TestCheckReleaseReadiness_ignoresEmptyActiveGoalAndKeepsCanonicalFindings; doctor accepts an empty in-progress Goal at Idea stage while the done/completion gate still rejects an empty Goal. | Pass; closes I-051 |
| Migration with a live router Goal | TestPlan_routerSelectsLiveGoalOrPlansContinuationGoal/live_selection_is_kept; retains the selected live Goal. | Pass |
| Migration with missing or unresolvable router Goal and active Objective | TestConvertRouter_fallbackSelectionPointsAtLiveGoal, TestEndToEnd_routerFallbackCasesMigrateToLiveGoal, and golden reproducibility checks. Selects the live Goal containing the active Objective (R-001 in these fixtures), creates no empty continuation, strictly loads, resolves Next, and exposes active Objective/Task in that Goal. | Pass; closes I-052 |
| Migration with archived router Goal and selected active work | Converter/end-to-end tests cover the historical source case: create the next continuation Goal (R-002), move the selected active Objective into it, and resolve router/Next without a Goal mismatch. | Pass; closes I-052 |
| Migration with completed historical releases and no live Goal | TestPlan_routerSelectsLiveGoalOrPlansContinuationGoal and end-to-end migration coverage verify creation and selection of a continuation Goal. | Pass |
| Audited or unrecognized release lifecycle with active work | TestPlan_auditedReleaseWithActiveEpicShowsDecisionInPreview, TestPlan_ambiguousReleaseDispositionBlocksCutover, and independent temporary overlay probe for audited and unrecognized statuses. Preview names the owner decision, Apply is blocked without writes, and supplying in_progress restores an appliable plan. No misleading no-resolvable-Goal error remains. | Pass; closes I-053 |
| Live/scaffold guidance and adoption test | Each of the four active skills and the shared-reference directory compared with its scaffold copy; all relevant pairs are byte-identical. TestV2AgentsGuideCarriesExistingCodebaseAdoptionSection pins reconciled migration wording. git diff --check passed. | Pass |

### Acceptance Classification

- **T-036:** all Done When criteria are proven. The typed missing-Goal fact is
  shared across surfaces, missing Objective references remain non-fatal facts,
  and doctor repairs zero-live-Goal cases including unknown and archived
  router selections.
- **T-037:** all Done When criteria are proven. The board remains within the
  selected Goal, displays no project-wide work without a valid selection,
  points zero-Goal cases to doctor, and reports missing-Goal Objectives.
- **T-038:** all Done When criteria are proven. Init/scaffold creates and
  selects R-001 named after the project; a fresh empty in-progress Goal passes
  doctor, while Goal completion constraints remain enforced.
- **T-039:** all revised Done When criteria are proven. Live selections remain
  selected; missing/unresolvable selections resolve to the identifiable live
  Goal for active work; historical-only selected work is moved into a new
  continuation; lifecycle ambiguity stays in preview and blocks Apply.
- **T-040:** all Done When criteria are proven. Required-Goal guidance matches
  implementation in live and scaffold copies, and the adoption test pins the
  reconciled fallback wording.
- **O-022 Success Conditions 1–6 and Outcome:** proven by the slices above.
  Existing Objective membership gaps remain visible and fixable; malformed
  and unknown references remain errors; new and migrated projects start with
  a usable selected Goal; no Goal-less project-wide view is introduced.
- **Design reconciliation:** Design §§1, 2, 3, 8, and 10 agree with runtime
  and migration behavior. The owner's I-052 choice is reflected consistently
  in T-039 and the guidance.

## Migration Workflow and Side Effects

| Phase | Evidence | Result |
| --- | --- | --- |
| Plan / preview | Audited and unrecognized active-release probes; migration planning tests. | Read-only; unresolved lifecycle decision is shown and plan is not appliable. |
| Apply while unresolved | Regression tests and independent overlay probe. | Refuses before planned file writes; no source removal, manifest, schema activation, or Goal selection is applied. |
| Apply after a resolved decision | Migration apply and end-to-end regression suite. | Uses existing Git/path conflict checks and batch planning; no new direct-write path or external network dependency was introduced. |
| Cross-platform behavior | Fresh make test-full. | Full Go suite and Linux, Darwin, and Windows builds passed. |

## Guardrails and Adversarial Pass

The frozen C-923 guardrail set was applied. DATA-01..04 remain satisfied:
conversion preserves source content according to existing parser and writer
contracts; lifecycle decisions remain in the data/migration owners; malformed
references and unresolved migration statuses report named facts; and defaults
such as the initial Goal are explicit and diagnosable. ARCH-01/03 remain
satisfied: CLI dispatch stays thin, planning/rendering behavior is
deterministic, and no new current-working-directory dependency was found.

TPL-01 passes: all four active skill copies and the shared references match
their scaffold copies. TPL-02 passes: required-Goal and migration fallback
guidance matches behavior. TPL-03/04 show no newly required optional file or
unrouted scaffold asset; the initial Goal is part of fresh init/scaffold and
migration behavior is covered by migration itself. FS-01/03 pass: preview is
read-only and unresolved decisions block before Apply writes. TEST-01..05/08
pass through named happy/failure regressions, temporary test data, migration
no-write checks, and the current full gate.

Adversarial cases considered: a stale router after all Goals are removed;
router selection of an archived Goal; empty in-progress versus completed
Goals; selected work whose source Goal is historical; selected work whose
source Goal is already live; and an active Objective whose release status is
audited or unrecognized. Each now has a coherent diagnostic, Goal assignment,
or preview decision. Missing Objective membership remains recoverable and
exactly directed, while malformed/unknown membership remains fatal. No new
material issue was found.

## Gates and Evidence

- GOCACHE=/tmp/savepoint-o022-check-gocache make test-full: **passed fresh**
  at 2026-09-25T04:51Z on go1.26.2 linux/amd64. The full Go test suite and
  Linux, Darwin, and Windows builds completed with exit 0. This run includes
  final migration expectation and scaffold adoption-test wording changes.
- git diff --check: passed on the reviewed working tree before this record.
- Live/scaffold comparisons: all four skills and the shared-reference
  directory passed diff -qr.
- Independent audited/unrecognized active-release probe:
  GOCACHE=/tmp/savepoint-o022-check-gocache go test -overlay=/tmp/savepoint-o022-check-overlay.json ./internal/migrate -run '^TestO022CheckProbeReleaseLifecycleWithActiveEpic$' -count=1 -v — both lifecycle cases passed; temporary probe source was kept under /tmp only.
- Toolchain: go version go1.26.2 linux/amd64; GOOS=linux, GOARCH=amd64.

## Materiality

The four C-923 issues were reproducible against supported project shapes and
material to their stated criteria. Their repairs now pass independent
scenarios and the current full gate. No materiality action remains and no new
finding met the locked materiality threshold.

## Non-blocking Observations

The observations from C-923 remain outside O-022 scope: plain-board project
totals are global counts; some doctor findings use absolute paths; this
uncommitted working tree mixes O-019 and O-022 changes; and doctor reports
waived Task clearance separately. None changes the O-022 acceptance result.

## Owner Validation Still Needed

All five Tasks retain owner_validation.required. The owner must accept this
current CLEAR Objective Check before closing O-022. Task-check waivers do not
replace this Check. Once O-022 is accepted, Goal R-006 still requires its
mandatory Goal Check and exact owner acceptance.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches**
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth** — board and doctor share V2Index.HasLiveGoal.
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [ ] STYLE-10 **Small diffs** — the reviewed working tree is uncommitted and interleaves O-019 and O-022 changes; advisory only.

