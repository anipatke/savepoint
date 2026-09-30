---
id: O-027
title: Define trustworthy health records and history
status: done
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

- The model represents scoped results for the five V2.1 capabilities—tests, coverage, complexity, duplication, and dependency vulnerabilities—and distinguishes available, partial, absent, unsupported, unavailable, not configured, failed, timed out, and cancelled outcomes where applicable.
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

On 2026-09-29, the owner confirmed the reconciled five-capability design. Collection outcome and freshness are separate dimensions: freshness is `fresh`, `stale`, or `unknown`, so missing or failed evidence cannot masquerade as stale-but-valid data. Configuration and immutable snapshots live under `.savepoint/health/`; stable snapshot identities derive from canonical snapshot content. Repository fingerprints cover configured scope and relevant tracked and untracked inputs without storing sensitive file content. History compares only compatible provider, schema, measurement-definition, configuration, and scope identities. Full Objective Check snapshots remain permanent; explicit maintenance retains only the ten newest manual snapshots. Classification produces deterministic Good, Watch, and Needs Attention summaries without a composite score.

On 2026-10-01, resolving the T-054 replan, the owner confirmed three classification rules:

- **Built-in default thresholds**, overridable per project: coverage Good ≥ 80%, Watch ≥ 60%; highest function complexity Good ≤ 10, Watch ≤ 20; duplication Good ≤ 3%, Watch ≤ 5%. Beyond Watch is Needs Attention.
- **Material decline** is measured against the median of the last three comparable official snapshots: coverage down 5 or more points, duplication up 2 or more points, complexity up 5 or more, or any increase in failing tests or vulnerabilities. Smaller movement is not material.
- **Vulnerability severity**: dependency-vulnerability results carry `high` and `critical` count details alongside the total. A high or critical count above zero is a hard blocker (Needs Attention). Vulnerabilities whose severity is unknown (total above zero with the severity counts missing) also block.
