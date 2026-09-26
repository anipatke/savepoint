---
id: O-027
title: Define trustworthy health records and history
status: planned
depends_on: [O-026]
release: R-007
priority: critical
rank: 2
---

# O-027: Define trustworthy health records and history

## Outcome

Code Health owns a versioned stack-agnostic model and compact immutable history that distinguishes measurement, absence, failure, staleness, comparability, trend, and overall health without false precision.

## Why

Providers and UI need stable semantics. Missing measurements must never become zero, passing, or healthy, and incomparable results must never form a false trend.

## Success Conditions

- The model represents scoped results for all six signals and distinguishes available, unavailable, not configured, failed, timed out, and cancelled outcomes where applicable.
- Versioned machine-readable snapshots live in a dedicated `.savepoint` health area with origin, provider/version, configuration and exclusion fingerprints, repository state, bounded detail, and sanitised evidence references.
- Repository identity handles commits, dirty tracked and relevant untracked inputs, ancestor-based “commits behind” wording, divergence, and repositories without a useful commit.
- Deterministic classification implements current-value thresholds, one-level worsening for material decline, hard blockers, partial-data summaries, and a recent range after three comparable official snapshots.
- History never compares incompatible providers, measurement definitions, versions, configuration, or scope; old series remain visible.
- Full Objective Check snapshots remain permanent. Only the ten newest manual snapshots are retained, with pruning performed explicitly outside collection.

## Architectural Considerations

Code Health owns these records. Check records reference snapshot identities rather than embedding the schema, and other packages consume a narrow read model.

## Boundaries

**In scope:** schema, validation, persistence, provenance, repository identity, retention, classification, trends, baselines, summaries, and fixtures.

**Out of scope:** execution, discovery, scanners, Check wiring, TUI rendering, automatic Issues, and composite scores.

## Confirmed Design Decisions

The owner confirmed versioned records intended for Git, separate Check references, dirty-tree fingerprints, official-only baselines, ten manual snapshots, comparability resets, configurable guidance with hard minimums, and no numerical score on 2026-09-26.
