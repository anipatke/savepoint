---
id: T-054
title: Explain health trends without false precision
objective: O-027
status: planned
depends_on: [{task: T-051, requires: clear}]
owner_validation: {required: true}
planned_by: {role: planner, session: planning-o027-20260929}
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

1. Encode per-capability current-value threshold evaluation and non-overridable hard minimums as pure typed rules.
2. Select only compatible official observations and calculate bounded recent ranges and trend direction after three points.
3. Apply the one-level material-decline rule and improvement explanation without generating a composite score.
4. Produce bounded structured explanations for incomplete, failed, stale, unknown, incompatible, and insufficient-history states.
5. Add exhaustive boundary and transition matrices using independently authored expected classifications and explanations.

## Boundaries

No provider collection, snapshot persistence, repository inspection, user-authored expression language, AI judgment, TUI rendering, Check clearance policy, automatic Issues, or composite numerical score.

## Technical Verification

Run focused pure classification tests during iteration, then `make build && make test-fast`. A later Check follows `agent-skills/references/check-method.md`.

## Technical Evidence

Pending execution: record matrix coverage, threshold vectors, commands, changed files, and limitations.

## Drift Notes

Return to planning if a capability lacks an owner-confirmed threshold or hard-blocking rule needed to classify it deterministically.
