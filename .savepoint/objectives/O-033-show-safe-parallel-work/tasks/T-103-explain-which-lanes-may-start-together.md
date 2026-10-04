---
id: T-103
title: "Explain which lanes may start together"
objective: O-033
status: done
depends_on: [{task: T-102, requires: clear}]
owner_validation:
  required: false
  accepted_check: ""
planned_by: {role: planner, session: codex-o033-planning-2026-10-03}
complexity_tier: high
complexity_reason: "Conservative compatibility combines lifecycle gates, transitive dependencies, reviewed scopes and recorded active work."
lane: core
planned_reads:
  - "internal/data/feature_preferences.go"
  - "internal/data/config.go"
  - "internal/data/concurrency_plan_v2.go"
  - "internal/data/task_v2.go"
  - "internal/data/objective_v2.go"
  - "internal/data/project.go"
  - "internal/data/dependency.go"
  - "internal/data/gate_v2.go"
  - "internal/data/objective_gate_v2.go"
  - "internal/data/next.go"
  - "internal/data/next_test.go"
planned_writes:
  - "internal/data/concurrency_v2.go"
  - "internal/data/concurrency_v2_test.go"
  - "internal/data/next.go"
  - "internal/data/next_test.go"
check_waiver:
  task: T-103
  reason: Owner completed this Task via the board without requesting a Task Check.
  actor:
    role: owner
    session: board-owner
  recorded_at: "2026-10-04T02:56:57Z"
---

# Explain which lanes may start together

## Outcome

One pure deterministic data projection reports persisted lanes, start-eligible Tasks, pairwise compatible opportunities and concrete sequential reasons within the selected Objective.

## User Check

Compare temporary projects with independent lanes, dependencies, shared inputs, active work and missing scopes; every withheld opportunity states why.

Owner validation is not required for this technical Task; owner Task completion remains required.

## Done When

- Read the saved parallel-planning preference from the predecessor interface. Off/absent removes advisory suggestions while preserving metadata; on adds advice without changing any Next selection, start, advance or completion gate. Unusable optional metadata never blocks ordinary work.
- Reuse ResolveTaskStart and Objective dependency gates; missing/stale selections and recorded replan prevent a launch suggestion without changing Next or its selected record.
- Recommend one Task at a time per suggested lane; the owner can ignore that recommendation. Do not recommend together work with a dependency path, write/write conflict or unexplained write/read overlap; ordinary execution remains available under its existing gates. Shared reads alone are allowed; valid narrow independence explanations affect only their bound read/write overlaps.
- Conservatively account for remaining lane scopes and recorded active work. A missing manifest or unresolved path alias leaves safety unknown; no filesystem/worktree monitoring, claimed state or inference that a done Task has merged.
- Produce deterministic pairwise-compatible groups and reasons; do not merely collect pairwise edges into a group that contains an incompatible pair. Empty or singleton opportunities do not claim parallelism.
- Projection includes stable lane membership for every lifecycle column, namespaced by Objective, and can be consumed by board detail and resume without duplicated safety logic. Explain future sequential constraints without presenting blocked work as start-ready.
- Per-criterion evidence and the configured handoff gate are recorded. Optional Task Checks require a fresh session; skipping one requires an explicit owner waiver naming this Task, reason, actor and time. All work remains subject to the mandatory Full Objective Check.

## Context Files

- `internal/data/feature_preferences.go`
- `internal/data/config.go`
- `internal/data/concurrency_plan_v2.go`
- `internal/data/task_v2.go`
- `internal/data/objective_v2.go`
- `internal/data/project.go`
- `internal/data/dependency.go`
- `internal/data/gate_v2.go`
- `internal/data/objective_gate_v2.go`
- `internal/data/next.go`
- `internal/data/next_test.go`

Prerequisite-created context: `internal/data/concurrency_plan_v2.go`. These files are intentionally created by preceding Tasks; verify them after dependencies are met. Other Context Files must exist. A material interface mismatch returns REPLAN REQUIRED.

Explicitly new write targets: `internal/data/concurrency_v2.go`, `internal/data/concurrency_v2_test.go`. Their absence before this Task is intentional, not a missing-context defect.

## Design References

O-033 Confirmed Design — Revised 2026-10-03, including Persisted plan and ownership, Canonical safety projection, Board/resume/handoff, and Guidance and verification. Project Design sections 1, 4, 5, 7, 8 and 9 provide implemented boundaries; Objective deltas govern the planned feature.

## Guardrails

FS-01, FS-02, FS-04, FS-05, DATA-01, DATA-02, DATA-03, DATA-05, TPL-01, TPL-02, TPL-04, ARCH-01, ARCH-02, ARCH-03, CFG-02, CFG-03, TEST-01, TEST-02, TEST-03, TEST-04, TEST-08, TEST-09, STYLE-01..10.

## Implementation Plan

1. Confirm the predecessor's typed manifest/explanation interface and return REPLAN REQUIRED if materially inconsistent.
2. Implement pure lane membership, transitive dependency and compatibility helpers in explicitly new concurrency_v2.go; keep algorithm bounded and deterministic, avoiding speculative optimal scheduling.
3. Attach the selected Objective projection to Next and expose the same projection function for a board's focused Objective without changing routing semantics.
4. Exercise happy paths and counterexamples: indirect dependencies, clear/accepted/waived prerequisites, Objective blocks, same-lane Tasks, future scope conflicts, unknown active scope, stale explanation, case aliases, valid shared reads and a three-lane incompatibility triangle.

## Boundaries

This Task implements its scoped plan and records execution evidence and authorised lifecycle updates. planned_writes describes anticipated work, not an enforcement allowlist. Necessary targeted reads must be logged before access. A difference from advisory manifests may be noted but is not itself a blocker or reason for REPLAN REQUIRED; existing material acceptance/architecture gaps still use the ordinary replan policy. Do not update another Task, allocate identity records in a lane, automate worktree/branch creation, push, merge, change Goal selection, add lifecycle states, or write a Check. In a lane skip router writes and commit locally. Set the Objective in_progress only when starting its first Task; later lanes must not rewrite the shared Objective simply to restate that status. Design reconciliation belongs to the final integration Task on main, not parallel lanes.

## Technical Verification

Fresh make test-full because supported-platform path comparisons form part of safety; native windows-tests CI evidence is produced by repository CI for the Full Check.

Use focused make test-focused TEST=... only for iteration. Record named happy-path and failure cases for each Done When criterion, actual commands/time/toolchain/results, reads and changes, and limitations. Independent verification follows agent-skills/references/check-method.md; no executor session performs its own Check.

## Technical Evidence

Executed 2026-10-04 (go1.26.2 linux/amd64). Predecessor interfaces confirmed: `FeaturePreferences`/`Config.ParallelPlanningEnabled` (T-108) and `TaskPlanV2`, `ObjectivePlanV2`, `V2Index.UnusablePlans/PlanIndependence/PlanDiagnostics`, `ComparePlanPaths` (T-102) matched the plan; no REPLAN REQUIRED.

Implementation: `internal/data/concurrency_v2.go` adds `ResolveConcurrencyV2(index, objectiveID, ConcurrencyOptionsV2{Enabled, Selection})` returning `*ConcurrencyV2` (lanes by column, candidates, pairs with reasons, first-fit groups, notes, diagnostics, withheld). `internal/data/next.go` adds `NextInput.ParallelPlanning` and `Next.Concurrency`, attached for any resolved Objective and consumed by nothing in Next selection.

Per-criterion evidence (all in `internal/data/concurrency_v2_test.go`):
- Preference off/absent, metadata kept: `offOrAbsentKeepsMembershipOnly` (membership identical on/off, no advice off); `TestResolveNext_concurrencyIsAdviceOnly` (Next identical on/off apart from `Concurrency`). Unusable metadata degrades: `ungroupableTasksAreExplained`.
- Reuses ResolveTaskStart/Objective gates; stale selection and replan withhold: `prerequisitesUseOrdinaryStartGate` (clear, waived, none, needs-work, waiver never satisfies accepted), `objectiveDependencyBlocksEveryTask`, `replanAndStaleSelectionWithholdLaunchAdvice`.
- One Task per lane; no dependency-path/shared-write/unexplained overlap; shared reads allowed; explanations bind to their overlap only: `oneTaskPerLane`, `dependenciesNeverTogether` (indirect path, path through a done Task), `sharedWriteAndOverlapAreWithheld`, `independentLanesStartTogether`, `independenceExplanationsAffectOnlyTheirOverlap`, `staleExplanationWithholdsAndSaysSo` (loaded project, digest current vs stale).
- Remaining lane scopes, active work, unknown/alias safety, no filesystem/worktree input: `remainingLaneScopeConflictsWithhold`, `recordedActiveWork` (disjoint, shared write, unknown, unlaned, same lane), `sharedWriteAndOverlapAreWithheld` (case alias, missing reads/writes). The projection takes only the index and options.
- Pairwise-compatible groups; singleton/empty claim nothing: `incompatibleTriangleNeverGroupsAConflictingPair`, `singletonAndEmptyClaimNothing`, `isDeterministic` (20 runs).
- Stable membership in every column, Objective-namespaced (`Ref()`), consumable by board/resume without duplicated logic; future constraints explained without start-ready claims: `offOrAbsentKeepsMembershipOnly`, `ungroupableTasksAreExplained`, start_blocked notes in the dependency tests.

Commands: `go vet ./internal/data`; `go test ./internal/data -run 'Concurrency|ResolveNext'` (pass); `make build && make test-full` at 2026-10-04T02:54Z (pass, all packages, including windows/darwin cross-builds).

Reads: all listed Context Files except `config.go` fully read via `Config` struct only and `project.go` V2Index definition only. Extra reads: `internal/data/concurrency_plan_v2_test.go` and `internal/data/dependency_test.go`/`gate_v2_test.go` helpers (fixture helpers `writePlanProject`, `newV2TestIndex`, `mustCheck`), `agent-skills/savepoint-task/SKILL.md` (skill not exposed as a tool).

Changes: `internal/data/concurrency_v2.go`, `internal/data/concurrency_v2_test.go` (new), `internal/data/next.go`; this Task file. Actual scope matches `planned_writes` except `next_test.go`, where the Next-attachment test lives in `concurrency_v2_test.go`.

Limitations: `main.go` resume and `internal/board/v2/load.go` still call `ResolveNext` without `ParallelPlanning`, so `Next.Concurrency` carries membership only until the board/resume Tasks pass the saved preference. Any selection diagnostic (including an unknown Issue) and any recorded replan on an unfinished Task of the Objective withholds all launch advice; this is a conservative reading of "missing/stale selections and recorded replan". Grouping is first-fit by lane order, not optimal. No Task Check requested and no owner waiver recorded; the native windows-tests CI evidence is for the Full Check.

## Drift Notes

Pending execution. Record material acceptance/architecture deviations under normal policy. Lane choice or anticipated-scope deviation alone is advisory and never a required replan.


## Advisory Feature Contract — Revised 2026-10-03

Parallel planning defaults off via features.parallel_planning. Code Health remains core and is unchanged. Lane metadata and read/write manifests are suggestions: the owner and all skills may ignore them even when enabled. No lane rule feeds existing start, advance, completion or Check gates, and no manifest deviation alone requires replan. Actual dependencies, acceptance criteria, Guardrails and rules for a real worktree remain authoritative. This revised Objective contract supersedes any earlier mandatory lane wording.

Prerequisite-created setting interface: `internal/data/feature_preferences.go` is intentionally created by T-108; its absence before that dependency completes is expected.
