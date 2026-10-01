---
id: O-031
title: Explain Code Health in the TUI
status: planned
depends_on: [O-027, O-028, O-029, O-030]
release: R-007
priority: medium
rank: 1
---

# O-031: Explain Code Health in the TUI

## Outcome

The existing TUI contains a dedicated Code Health view giving non-expert builders an immediate, honest answer about current health, trends, staleness, incomplete measurement, and why any signal needs attention.

## Why

Users should understand engineering signals without learning provider jargon or mistaking missing, failed, or stale measurement for healthy code.

## Success Conditions

- A dedicated top-level view follows existing visual, navigation, and overlay conventions, and introduces the board's first progress and cancellation display. Objective Check shows only a compact summary and route to it.
- Opening the dashboard is instant, reads persisted data, and never scans. First run explains the feature and offers explicit refresh.
- `[R] Refresh` invokes the Code Health service directly for confirmed providers, displays sequential progress, supports cancellation, and persists a manual result without AI.
- The overview uses plain-language Good, Watch, Needs Attention, partial, unknown, and stale states with no overall number.
- Details explain classification, trend, affected areas, provider outcome and provenance, component scope, and comparison basis.
- Healthy is reserved for complete successful required data. Partial support and optional failures remain visible, and no wording claims proof of correctness, safety, maintainability, or production readiness.
- Tests cover first run, the five signals, partial support, failures, dirty/stale/divergent states, comparable and reset trends, refresh progress, and cancellation.

## Architectural Considerations

The board consumes a read-only dashboard model and sends explicit collection commands. Rendering performs no filesystem or subprocess work, and provider/storage types stay inside Code Health.

## Boundaries

**In scope:** navigation, overview, details, history, first run, stale/partial/failure states, manual refresh, and compact Check presentation.

**Out of scope:** web UI, background refresh, watchers, timers, cloud reporting, AI summaries, automatic setup, and automatic Issues.

## Confirmed Design Decisions

The owner chose a dedicated main-board view, direct deterministic refresh, no implicit scanning, complete-data Healthy semantics, and narrow UI integration on 2026-09-26.

On 2026-10-02 the owner confirmed the screen's shape in chat:

- **Five signals.** Tests, coverage, complexity, duplication, and dependency vulnerabilities; change hotspots were de-scoped in O-029 on 2026-09-27.
- **Navigation.** `H` on the board opens a full-screen Code Health view; Esc returns to the board where the owner was, like the Issues view. Enter on a signal opens its details; a short history lists earlier snapshots with origin and overall label.
- **Opening.** The view appears from saved configuration and snapshots only. Without configuration it explains the feature and points to `savepoint health setup`; with configuration and no snapshot it offers `R`. A quick Git comparison then runs as an explicit board command and marks results whose code has moved on or come from another branch. Opening never runs health tools.
- **Refresh.** `R` runs `Collect` with origin `manual` over the confirmed configuration and shows sequential progress (for example "2 of 5: coverage"). Savepoint does not run tests: test and coverage instances re-read existing reports, and a report older than the code shows as stale with "run your tests, then refresh".
- **Cancel.** Esc during refresh cancels it and saves nothing; the view keeps the last saved result.
- **Check summary.** An Objective detail whose Check names a `health_snapshot` shows one Health line with that snapshot's overall label and a pointer to `H`.
