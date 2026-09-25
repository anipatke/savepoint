---
id: C-923
scope: {kind: objective, id: O-022}
result: NEEDS WORK
checked_by: {role: checker, session: o022-objective-check-20260925}
executed_session: unrecorded-external-executor
checked_at: '2026-09-25T03:27:22Z'
reviewed:
  base_commit: 35ce82dffbcfa02a88121f99496d668c67d69db9
  files:
    - internal/data/next.go
    - internal/data/project.go
    - internal/resume/resume.go
    - internal/resume/evidence.go
    - internal/doctor/v2_runtime.go
    - internal/board/v2/card.go
    - internal/board/v2/next_panel.go
    - internal/board/v2/objectives.go
    - internal/board/v2/plain.go
    - internal/board/v2/releases.go
    - internal/board/v2/update.go
    - internal/board/v2/view.go
    - internal/migrate/plan.go
    - internal/migrate/convert.go
    - internal/migrate/convert_docs.go
    - internal/migrate/convert_releases.go
    - templates/project-v2/.savepoint/releases/R-001-first-goal/Release.md
    - templates/project-v2/.savepoint/router.md
    - AGENTS.md
    - .savepoint/Design.md
  dependencies: []
issues: [I-050, I-051, I-052, I-053]
supersedes: null
---

# C-923: O-022 Full Objective Check

## Independence and Scope

This session did not build O-022. The work is uncommitted on branch `v2` on
top of `35ce82d`, mixed with O-019's uncommitted changes. The review covers the
O-022 diff in the working tree at the time of this Check. All five Tasks
(T-036..T-040) are `done` with owner Task-check waivers. Those waivers were
inspected and are not treated as technical CLEAR. Each Task was reviewed
directly.

### Scope lock

1. Criteria: O-022 Success Conditions 1–6 and Outcome; Done When of T-036..T-040;
   guardrails DATA-01..04, ARCH-01/03, TPL-01..04, FS-01/03, TEST-01..05/08.
2. Entry points: `data.ResolveSelection`/`ResolveNext`, `data.LoadV2Index`,
   `resume.NextLine`/`NextVerb`/`Render`, `doctor.RunV2Checks`, the V2 board
   (TUI model and plain renderer), `savepoint init`, and `migrate.Plan` with
   router/Objective/Release conversion.
3. Relied-on behavior: router `none`/blank normalization, Release liveness
   (`LegacyCompletion`), existing selection diagnostics, doctor's
   Release-readiness checks, and migrate's release lifecycle ambiguities.
4. Matrix axes: router Goal {absent key, blank, `none`, unknown, archived, live}
   × other selections {none, Objective, Objective+Task, Issue only}; Goals
   {zero, archived-only, live}; Objective `release:` {present, missing, one,
   several, malformed/unknown}; surfaces {resume, board TUI model, board plain,
   doctor}; project origin {fresh init, migrate: live/missing/unresolvable/archived
   router release, release lifecycle ambiguity}; guidance {live, scaffold}.
5. Materiality: an Issue must be reproducible through a supported command on a
   supported project shape and must break an O-022 or Task criterion, the
   Objective Outcome, or a guardrail.

## Coverage Matrix

| Row | Probe | Result |
| --- | --- | --- |
| Router Goal `none` / blank / key absent | Scratch copy of this repo, real binary | `Choose a Goal — press g on the board` for all three. **Pass** |
| Missing Goal + Objective/Task, Issue only | Same; router objective O-022, then issue I-049 only | Choose a Goal in both cases; no work substituted. **Pass** |
| Missing Goal: doctor | Same | `[router-goal-missing]` with `Choose a Goal with g on the board.` **Pass** |
| Missing Goal: board plain, `--objective O-022` | Same | 0 cards, `Selected: no Goal selected`, diagnostic shown. **Pass**; project-wide counts remain (observation) |
| Objectives without `release:` (two) | Removed from O-019, O-022 | Loads. Resume lists `O-019, O-022`. Board notice says 2 Objectives. Doctor names each file and the `release: R-###` line. **Pass** |
| Malformed/unknown Objective reference | Existing `TestLoadV2Index_invalidGoalReferencesRemainFatal` | Stays fatal. **Pass** |
| Zero Goals, router `none` | Fresh init with releases removed | Doctor: `Create a Goal first, then choose it with g`. **Pass** |
| Zero Goals, router names deleted R-001 | Same, router untouched | Board: `No Goals exist; run savepoint doctor.` Doctor: ALL CLEAN. **Issue I-050** |
| Archived-only Goals, router on archived Goal | Copied archived Release into fresh project | No doctor Goal guidance and no board pointer. **Issue I-050** |
| Goal selector cannot clear | `TestGoalSelectorCannotClearTheCurrentGoal` | **Pass** |
| Goal-scoped label, TUI and plain | Code trace + `TestGoalScopedSelectionWithObjectiveFilter`; plain probe | `ALL OBJECTIVES IN GOAL R-###` / `all Objectives in Goal R-006`. **Pass** |
| Fresh init: R-001, title, router | Real `savepoint init` in `my-proj` | R-001 `title: my-proj`, `in_progress`, four stub sections, router `release: R-001`. Resume shows no missing-Goal diagnostic. **Pass** |
| Fresh init: doctor | Same | `v2-release-no-objectives`, PROBLEMS FOUND exit 1. **Issue I-051** |
| Migrate: live router release kept | `TestPlan_routerSelectsLiveGoalOrPlansContinuationGoal/live_selection_is_kept` | **Pass** |
| Migrate: missing router release | Materialized `v1-router-missing` golden, real binary | Router selects empty R-002, while live R-001 holds the active O-001/T-001. Board shows 0 cards. **Issue I-052**; doctor also fails on R-002 (**I-051**) |
| Migrate: archived router release | Materialized `v1-router-archived` golden | Router O-001/T-001 is not honored (Goal mismatch), `Nothing selected`. **Issue I-052** |
| Migrate: unresolvable router release | Golden + `TestEndToEnd_routerFallbackCasesMigrateToLiveGoal` | Continuation Goal created and selected; same data shape as the missing case. Covered by **I-052** |
| Migrate: epic in release with no PRD | `TestPlan_activeEpicWithoutResolvableReleaseFailsWithNamedError` | Named error. **Pass** |
| Migrate: audited/unknown release status + active epic | Temporary probe test (removed) | Plan fails with "no resolvable V2 Goal"; ambiguity never shown. With the decision supplied, the plan is appliable. **Issue I-053** |
| Guidance: live and scaffold | `cmp` of four skills and three references; AGENTS Required Goal Context section md5 and the `Choose` routing line | Byte-identical. **Pass** |
| Guidance: optional wording regression | `TestProjectGuidanceRequiresGoalContext` | **Pass** |
| Guidance matches behavior | AGENTS.md:79 / Design §1 doctor promise vs zero-Goal probe | Doctor does not say to create one when the router names a missing or archived Goal. **I-050** (TPL-02) |

Not applicable: text-width classes (no truncation logic changed; narrow widths
are covered by `TestGoalScopedChromeFitsNarrowWidths`); external network
boundaries (none). Migrate apply side effects are unchanged apart from one
generated target, which goes through the existing apply path and path checks.

## Acceptance Classification

- **T-036**: Done When 1–3, 5, 6 **Proven**. Done When 4 (zero Goals → repair says
  create one) is **Proven** for a blank router, but the unknown/archived router
  paths leave doctor silent → **I-050**.
- **T-037**: Done When 2–5 **Proven**. Done When 1's doctor pointer leads to an
  all-clean doctor → **I-050**.
- **T-038**: Done When 1–6 **Proven** as written. The User Check and the
  Objective Outcome "new projects start valid" fail on doctor → **I-051**.
- **T-039**: Done When 2, 4, 5 **Proven**. Done When 1 **Proven** for no-PRD releases
  but breaks the lifecycle-ambiguity path → **I-053**. Done When 3 holds, but its
  output leaves active work outside the selected Goal → **I-052**, and the new
  Goal fails doctor → **I-051**.
- **T-040**: Done When 1–3, 5, 6 **Proven**. Done When 4 / TPL-02 holds except the
  doctor zero-Goal promise → **I-050**.
- **Design reconciliation**: §1, §2, §3, §8, §10 describe the implemented rule.
  The only mismatch is the doctor zero-Goal promise (I-050).

## Gates

- `make test-full`: **passed**, run fresh on the reviewed working tree at
  2026-09-25T03:22Z, `go1.26.2 linux/amd64`. The full Go suite passed and the
  linux/darwin/windows builds ran. Exit 0.
- `make build`: passed. `git diff --check`: passed.
- Probe harness: real `./savepoint` binary on scratch copies (this repo's
  `.savepoint`, fresh `init` folders, and the three materialized fallback
  migrate goldens). One temporary `internal/migrate` probe test was removed
  after the run.

## Materiality

| Issue | Likelihood | Impact | Materiality | Recommendation |
| --- | --- | --- | --- | --- |
| I-050 doctor silent with no live Goal | Low: needs Goals deleted or archived while the router still names one | Medium: the board's only guidance leads to an all-clean doctor | Medium | Fix now; small doctor change and one shared "no live Goal" definition |
| I-051 fresh/migrated Goal fails doctor | High: every `savepoint init` | Low–Medium: doctor exit 1 on day one; the repair cannot be acted on at Idea state | Medium | Fix now, or record a planner decision that this is expected |
| I-052 migration hides active work | Medium: V1 routers with no release, or a completed release | Medium–High: migrated work is invisible on the board, and in the archived case the selection is lost | Medium | Needs a planner decision on which Goal holds converted active work; then a direct repair |
| I-053 ambiguity hidden by the new error | Low: audited or unrecognized release status with active epics | Medium: migration is blocked by a misleading error | Low | Combine with the I-052 migrate repair |

## Observations (non-blocking)

- Plain board output still prints project-wide totals (`Objectives: 13  Tasks: 39`)
  and the project-wide Issue summary with no Goal or inside one Goal. These are
  counts, not record lists, and predate O-022.
- The doctor problems added by O-022 use absolute file paths
  (`filepath.Join(root, ...)`), while other V2 problems print `.savepoint`-relative
  paths. This follows the older router-done entry.
- The working tree mixes O-019 and O-022 changes (75 files) and is uncommitted.
  Committing the two Objectives separately would make later rechecks easier to
  scope.
- Doctor flags Tasks closed under an owner waiver with
  `v2-done-without-clearance`, including all five O-022 Tasks. This predates
  O-022 and is out of scope here.

## Owner Validation Still Needed

All five Tasks declare `owner_validation.required`. Owner acceptance of a
current CLEAR Objective Check is still needed after remediation.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches**
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [ ] STYLE-07 **One source of truth** — "no Goals" is `len(Index.Releases) == 0` in `internal/board/v2/next_panel.go:39` but "no live Goal" in `internal/doctor/v2_runtime.go` `hasLiveGoal`.
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [ ] STYLE-10 **Small diffs** — O-022 is uncommitted and interleaved with O-019 in one 75-file working tree.
