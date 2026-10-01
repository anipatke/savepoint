---
id: T-078
title: Draw a sparkline of recent official checks
objective: O-035
status: done
depends_on: [{task: T-077, requires: clear}]
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: planning-o035-20261002}
complexity_tier: medium
complexity_reason: Windowing, scaling and several edge rules over history the package already selects; correctness is in the edge cases.
check_waiver:
    task: T-078
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-01T22:14:56Z"
---

# Draw a sparkline of recent official checks

## Outcome

Each measured row carries a small sparkline of its last official checks with a better, worse or steady word and a note when history is thin or restarted.

## User Check

None beyond the Full Objective Check.

## Done When

- `DashboardRow` gains `Spark`, `SparkWord` and `SparkNote`, all plain strings.
- Points are the comparable official values already selected for the row's trend, oldest first, plus the current value when the newest snapshot is official, at most 10. Manual, partial, incompatible and non-finite values never contribute.
- With fewer than 3 points `Spark` is empty and `SparkNote` says "not enough history yet (n of 3)". With 3 or 4 points `SparkNote` is "early". Incompatible earlier points set `SparkNote` to say the series restarted because settings changed.
- Each point maps to one of eight blocks `▁▂▃▄▅▆▇█` scaled to the window's own minimum and maximum; taller always means a bigger number. When the spread is smaller than the signal's existing material-movement size, every point is the middle block, so noise does not look like movement.
- `SparkWord` is "better", "worse" or "steady" from the existing trend direction (improving, declining, steady); it is empty when there is no trend.
- The existing trend wording, thresholds, series selection and stored data are unchanged.
- Tests cover 1, 2, 3, 5 and 12 points, a flat series, a below-material wobble, a rising and a falling series for a higher-is-better and a lower-is-better signal, a manual refresh among the points, a partial result, and a restarted series.

## Context Files

`.savepoint/objectives/O-035-glanceable-code-health/Objective.md`; `internal/codehealth/dashboard.go`, `internal/codehealth/dashboard_test.go`, `internal/codehealth/dashboard_copy.go`, `internal/codehealth/spark.go` (new), `internal/codehealth/spark_test.go` (new); `internal/codehealth/history.go`, `internal/codehealth/classification.go` (read only: series selection, trend, material movement).

## Design References

O-035 Confirmed Design Decisions (sparklines, level-based verdict); O-027 history rules.

## Guardrails

ARCH-04, STYLE-05, STYLE-07, STYLE-09, TEST-01, TEST-02, TEST-04.

## Implementation Plan

1. Confirm `selectSeries`, `buildTrend` and the material-movement table are reachable from `measuredRow`; return REPLAN REQUIRED otherwise.
2. Add a pure function from values and capability to the three strings in `spark.go`.
3. Call it from `measuredRow` with the same series the trend uses.
4. Add the tests in Done When.

## Boundaries

No new storage, no change to when snapshots are made, no colouring or layout (the board does that), no change to the trend or classification.

## Technical Verification

Focused tests during iteration (`make test-focused`); `make build && make test-fast` for ordinary handoff. Later evaluated per `agent-skills/references/check-method.md`.

## Technical Evidence

Executed by executor session; ready for an optional Task Check or the Full Objective Check (not claimed clear).

Per criterion:
- `DashboardRow` fields: `Spark`, `SparkWord`, `SparkNote` plain strings added (`dashboard.go`); set in `measuredRow`. `TestDashboardSparklineThroughTheRow`.
- Points: `sparkValues` takes the trend's series values plus the current value when official and usable, newest 10. `TestSparkValuesKeepsTheNewestTen`, `TestDashboardSparklineIgnoresManualPartialAndKeepsTrend` (manual and partial excluded; non-finite is excluded by `usableValue` in `selectSeries`).
- Thin and early notes: `TestSparklinePointCounts` (1, 2, 3, 4, 5, 10 of 12 points), `TestSparklineEarlyBoundary`, `TestDashboardSparklineThinHistory`. Restart note: `TestSparklineRestartedNote`, `TestDashboardSparklineRestartedSeries`.
- Scaling and flat rule: `TestSparklineTallerIsAlwaysBiggerNumber` (rising and falling, higher- and lower-is-better), `TestSparklineFlatAndBelowMaterialDrawTheMiddleBlock` (flat, below-material wobble, at-material shape).
- `SparkWord` from the existing trend direction, empty without trend: `TestSparklineWordFollowsTrendDirection`, `TestDashboardSparklineThroughTheRow`.
- Trend wording, thresholds, series selection and storage unchanged: `TestDashboardSparklineLeavesTrendWordingAlone` and the existing dashboard tests pass.

Commands: `go test ./internal/codehealth` passed; `make build && make test-fast` passed.

Files read: the Context Files plus `internal/codehealth/identity.go` (SeriesID) and `.savepoint/Guardrails.md` rules. Extra reads: `identity.go` to confirm restart (series identity) semantics; `AGENTS.md` and `agent-skills/savepoint-task/SKILL.md` (skill tool reported unknown skill). Files changed: `internal/codehealth/dashboard.go`, `internal/codehealth/spark.go` (new), `internal/codehealth/spark_test.go` (new), this Task.

Limitations: `sparkValues` repeats the first four lines of `buildTrend`'s value assembly because `history.go` was read-only for this Task. A restarted series note says "restarted because settings changed" for any series-identity change (tool version, scope, definition), per the Task wording. `SparkWord` follows the existing 5-point trend window while the drawing can span 10 points. No board rendering was run; the board does not use these fields yet.

## Drift Notes

New module or architecture delta beyond the documented Codebase Map, reconciled through the planner before Check.
