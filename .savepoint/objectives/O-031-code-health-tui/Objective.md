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

- A dedicated top-level view follows existing visual, navigation, overlay, progress, and cancellation conventions. Objective Check shows only a compact summary and route to it.
- Opening the dashboard is instant, reads persisted data, and never scans. First run explains the feature and offers explicit refresh.
- `[R] Refresh` invokes the Code Health service directly for confirmed providers, displays sequential progress, supports cancellation, and persists a manual result without AI.
- The overview uses plain-language Good, Watch, Needs Attention, partial, unknown, and stale states with no overall number.
- Details explain classification, trend, affected areas, provider outcome and provenance, component scope, and comparison basis.
- Healthy is reserved for complete successful required data. Partial support and optional failures remain visible, and no wording claims proof of correctness, safety, maintainability, or production readiness.
- Tests cover first run, six signals, partial support, failures, dirty/stale/divergent states, comparable and reset trends, refresh progress, and cancellation.

## Architectural Considerations

The board consumes a read-only dashboard model and sends explicit collection commands. Rendering performs no filesystem or subprocess work, and provider/storage types stay inside Code Health.

## Boundaries

**In scope:** navigation, overview, details, history, first run, stale/partial/failure states, manual refresh, and compact Check presentation.

**Out of scope:** web UI, background refresh, watchers, timers, cloud reporting, AI summaries, automatic setup, and automatic Issues.

## Confirmed Design Decisions

The owner chose a dedicated main-board view, direct deterministic refresh, no implicit scanning, complete-data Healthy semantics, and narrow UI integration on 2026-09-26.
