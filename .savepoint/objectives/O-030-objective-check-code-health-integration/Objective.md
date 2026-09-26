---
id: O-030
title: Measure health during Full Objective Checks
status: planned
depends_on: [O-027, O-028, O-029]
release: R-007
priority: high
rank: 3
---

# O-030: Measure health during Full Objective Checks

## Outcome

Every Full Objective Check can collect and reference an official Code Health snapshot without weakening checker independence, duplicating existing gates, or confusing provider failure with unhealthy code.

## Why

Objective Check is the canonical measurement point, but snapshots remain deterministic supporting evidence rather than independent clearance.

## Success Conditions

- Full Objective Checks invoke configured collection and reference the resulting snapshot; Task Checks and ordinary Savepoint activity do not run the suite.
- Fresh build/test artifacts are reused so the same underlying gate does not run twice.
- Required providers must produce valid evidence. Optional provider failures are reported without independently blocking clearance.
- Failing tests and high/critical vulnerabilities block by default. Coverage, complexity, duplication, and hotspot results block only through confirmed project policy.
- Collection failure, incomplete coverage, stale artifacts, unhealthy measurements, Check findings, and clearance remain distinct.
- Warnings never create Issues automatically; the independent checker uses the existing issue-capture workflow when durable follow-up is warranted.
- Manual snapshots have no independent verification standing and cannot be mistaken for official Check evidence.

## Architectural Considerations

Check integration depends only on Code Health's collection result and snapshot identity. Snapshot storage does not import Check types, and Checks do not embed the health schema.

## Boundaries

**In scope:** Full Objective Check collection, artifact reuse, required/optional outcomes, blocking policy, references, evidence wording, and independent-checker tests.

**Out of scope:** Task Check changes, automatic completion, automatic Issues, per-commit collection, and manual-refresh clearance.

## Confirmed Design Decisions

The owner confirmed Full Objective Check and manual refresh as the only collection points, required/optional capability policy, explicit blockers, independent Issue judgment, and separate Check/snapshot ownership on 2026-09-26.
