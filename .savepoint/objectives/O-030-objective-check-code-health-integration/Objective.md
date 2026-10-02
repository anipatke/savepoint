---
id: O-030
title: Measure health during Full Objective Checks
status: done
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

On 2026-10-01 the owner confirmed the integration shape in chat:

- **Trigger.** A narrow agent-runnable command, `savepoint health check O-### [dir]`, is the checker's only collection entry point. It is permitted to agents only during a Full Objective Check, the same kind of narrow exception as `create-task`, and it only ever writes an `official` snapshot. Manual snapshots come only from the TUI refresh (O-031).
- **Artifact reuse.** The checker runs the project's full gate first; test and coverage instances then read the reports that gate wrote rather than running tests again. A required instance whose report is absent or stale cannot produce valid evidence and blocks.
- **Default blockers.** Any failing test, and any known high or critical vulnerability, block clearance whether the instance is required or optional.
- **Opt-in blockers.** A per-instance `blocking: true` in the health configuration makes a Needs Attention result for coverage, complexity, or duplication block; Watch is reported only. The setting is not accepted on tests or vulnerabilities, which already block by default.
- **Unknown-severity vulnerabilities** are reported as needing review and do not block; the checker decides whether an Issue is warranted.
- **Required versus optional.** A required instance that fails, times out, is unavailable, absent, unsupported, or stale blocks. An optional one is reported without blocking. Partial evidence is reported as incomplete and blocks only through a default or opt-in rule.
- **Distinct states.** The verdict keeps collection failure, incomplete coverage, stale artifacts, unhealthy measurement, Check findings, and clearance as separate statements. It never creates Issues.
- **Reference.** A Check records only the snapshot identity in an optional `health_snapshot` field. `internal/data` checks its shape without importing Code Health; `savepoint doctor` reports a reference to a missing or non-official snapshot. A project without health configuration states that Code Health is not configured, and the Check proceeds.
