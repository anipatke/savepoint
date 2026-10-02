---
id: T-054
title: Explain health trends without false precision
objective: O-027
status: done
depends_on: [{task: T-051, requires: clear}]
owner_validation:
    required: true
    accepted_check: ""
planned_by: {role: planner, session: planning-o027-20260929}
check_waiver:
    task: T-054
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-30T21:27:22Z"
---

# Explain health trends without false precision

## Outcome

Pure Code Health evaluation turns current compatible evidence into deterministic Good, Watch, or Needs Attention results and honest recent-history explanations without composite scores or comparisons across incompatible series.

## User Check

Review a compact scenario table covering healthy values, threshold crossings, hard blockers, material decline, improvement, missing/partial/stale/unknown evidence, fewer than three official observations, incompatible history, and manual-only history. Confirm every summary is understandable and avoids a numerical overall score.

## Done When

- Current-value thresholds determine the base classification, with confirmed hard blockers for failing tests and high/critical vulnerabilities and configurable guidance that cannot weaken hard minimums.
- A material decline may worsen the base classification by at most one level; improvement changes the explanation without upgrading a poor current value.
- Only comparable official Full Check snapshots contribute to baselines, recent ranges, and trends; manual snapshots and incompatible series remain visible but never alter official classification.
- Recent-range and trend wording begins only after at least three comparable official observations and handles ties, zero values, missing values, and non-finite input without false precision.
- Partial, absent, unsupported, unavailable, not configured, failed, timed-out, cancelled, stale, and unknown-freshness evidence yield explicit deterministic summaries and never silently become Good.
- Pure parameterized tests cover every capability, classification transition, threshold boundary, compatibility reset, history length, failure state, and hard-blocker override with an independent expected-result table.

## Context Files

`.savepoint/objectives/O-027-code-health-model-and-history/Objective.md`; `.savepoint/releases/R-007-v2-1-code-health/Release.md`; `.savepoint/Design.md`; `.savepoint/Guardrails.md`; `internal/codehealth/model.go`; `internal/codehealth/classification.go`; `internal/codehealth/classification_test.go`.

## Design References

Design sections 1, 3, 4, and 13; R-007 Confirmed Design Decisions.

## Guardrails

DATA-03, ARCH-03, ARCH-04, CFG-01, TEST-01, TEST-02, TEST-04, TEST-07, TEST-08.

## Implementation Plan

1. Encode per-capability current-value threshold evaluation and non-overridable hard minimums as pure typed rules, using the owner-confirmed rules in the Objective (2026-10-01): built-in default thresholds when a project configures none; failing tests above zero, and `high`/`critical` vulnerability detail counts above zero, block. A vulnerability total above zero with missing severity counts also blocks.
2. Select only compatible official observations and calculate bounded recent ranges and trend direction after three points.
3. Apply the one-level material-decline rule (the confirmed per-capability minimum movement against the median of the last three comparable official snapshots) and the improvement explanation, without generating a composite score.
4. Produce bounded structured explanations for incomplete, failed, stale, unknown, incompatible, and insufficient-history states.
5. Add exhaustive boundary and transition matrices using independently authored expected classifications and explanations.

## Boundaries

No provider collection, snapshot persistence, repository inspection, user-authored expression language, AI judgment, TUI rendering, Check clearance policy, automatic Issues, or composite numerical score.

## Technical Verification

Run focused pure classification tests during iteration, then `make build && make test-fast`. A later Check follows `agent-skills/references/check-method.md`.

## Technical Evidence

Files changed: `internal/codehealth/classification.go` (thresholds, hard blockers, `Assess`, `Overall`), `internal/codehealth/history.go` (series selection, baseline, trend, wording), `internal/codehealth/classification_test.go`, and the `internal/codehealth` row of the `AGENTS.md` Codebase Map. Files read: the Context Files plus `internal/codehealth/{snapshot,primitives,config,identity,errors}.go` and `AGENTS.md`/`savepoint-task` SKILL (extra reads: needed the snapshot/config/series-identity types the plan builds on). Context Files `classification.go` and `classification_test.go` did not exist; they are this Task's new files, not a plan gap. Renamed my `Observation` type to `HistoryEntry` and `plural` to `pluralize` to avoid collisions with `repository.go`.

Owner decision (2026-10-01, after the User Check table): Good is capped at Watch until three comparable official observations exist (this result counts when official and complete); covered by `TestGoodIsCappedAtWatchUntilThreeComparableChecks`. Re-ran `make build && make test-fast` (exit 0).

Design choices beyond the Objective's confirmed rules (owner may revisit): tests and vulnerabilities use default thresholds Good 0 / Watch 0 and Good 0 / Watch unbounded, with hard blockers overriding; a decline baseline needs three earlier comparable official points (the trend needs three including a current official, complete result); partial, stale, and unknown-freshness evidence caps Good at Watch but never softens Watch or Needs Attention; every non-measured outcome is Unknown; `Overall` ranks Needs Attention > Unknown > Watch > Good, so any Unknown capability keeps Overall from Good; manual current results are judged against the official baseline but never join history; partial or incompatible-series or non-finite history points are excluded from baselines and trends.

Per-criterion evidence (`go test ./internal/codehealth`):
1. Thresholds and hard blockers: `TestAssessCurrentValueThresholds` (boundary vectors for coverage 80/79.9/60/59.9, complexity 10/11/20/21, duplication 3/3.1/5/5.1, tests 0/1, vulnerabilities total/high/critical/missing severity, configured guidance, and guidance unable to weaken failing-test, high/critical, and unknown-severity blockers).
2. One-level decline and improvement: `TestAssessMaterialDeclineAndImprovement` (exact minimum movements, Good→Watch, Watch→Needs Attention, floor, one level only, median of last three, no baseline under three, improvement never upgrades) and `TestAssessTestsAndHardBlockersIgnoreDecline`.
3. Official-only comparable history: `TestAssessHistoryFiltering` (manual, provider/version/schema/definition/config/scope/exclusion resets, other capability, failed/partial/missing/non-finite history, unknown origin, manual or partial current).
4. Three-observation range and trend, ties, zeros, non-finite: `TestAssessTrendWording`, `TestAssessTrendExactWording`.
5. Incomplete evidence never Good: `TestAssessUnmeasuredAndIncompleteEvidence` (all nine outcomes plus stale, unknown freshness, NaN/Inf/negative/missing values) and `TestEverySummaryFitsTheSnapshotContract` (every capability × outcome × freshness is bounded and never Good without available, fresh evidence); `TestOverallHasNoScoreAndNeverHidesProblems`.
6. Independent expected tables: all expectations are hand-written literals from the Objective's confirmed rules, not computed from the code under test. Also `TestAssessIsDeterministicAndDoesNotMutateInput`, `TestBoundTextKeepsUTF8AndLimit`.

Commands: `go vet ./internal/codehealth`, `go test ./internal/codehealth` (pass), `make build` (ok), `make test-fast` (exit 0).

Limitations: no `make test-full` (not required for ordinary handoff). Explanations are capped at 200 bytes, so very long wording would drop trailing notes; current worst case fits. Nothing yet calls `Assess`; wiring to snapshots and Checks belongs to later Tasks. The User Check scenario table has not been presented to the owner yet.

Replan history: on 2026-10-01 the executor returned REPLAN REQUIRED because default thresholds, the material-decline rule, and vulnerability severity inputs were undefined. The owner confirmed all three the same day; they are recorded in the Objective's Confirmed Design Decisions.

## Drift Notes

Return to planning if a capability lacks an owner-confirmed threshold or hard-blocking rule needed to classify it deterministically.
