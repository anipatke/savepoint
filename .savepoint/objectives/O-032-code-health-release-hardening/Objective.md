---
id: O-032
title: Harden Code Health for the v2.1 release
status: planned
depends_on: [O-026, O-027, O-028, O-029, O-030, O-031]
release: R-007
priority: medium
rank: 2
---

# O-032: Harden Code Health for the v2.1 release

## Outcome

Code Health is documented, migration-safe, privacy-conscious, cross-platform, performant, and independently verified as a modular foundation for future providers without weakening existing Savepoint behavior.

## Why

Persisted records, external processes, scaffolding, and a substantial TUI surface require explicit release evidence beyond feature-level tests.

## Success Conditions

- New-project scaffolding and existing-project upgrades preserve user content and maintain confirmed Design/config separation.
- Schema compatibility, strict loading, upgrade/downgrade diagnostics, concurrent access, partial writes, malformed files, and retention maintenance have regression coverage.
- Execution is tested for path safety, cancellation, timeout, process cleanup, bounded output, sanitisation, report limits, supported platforms, and network-independent normal tests.
- Packaging proves Savepoint remains one binary and neither bundles nor dynamically installs providers. Prerequisites and unsupported stacks are documented.
- Documentation explains purpose, limitations, capabilities, provider availability, states, thresholds, trends, refresh, Objective Checks, staleness, history, privacy, local-first behavior, and why Code Health is not a guarantee.
- The Full Objective Check covers all providers, partial/polyglot fixtures, history and staleness, Check integration, TUI behavior, existing behavior, Design reconciliation, and modularity.
- Release notes include a textual TUI walkthrough, provider rationale, packaging implications, supported and unsupported cases, limitations, and separate follow-up work.

## Architectural Considerations

Hardening verifies the module boundary. Provider or platform incompatibility is represented honestly as unavailable rather than solved with bundled runtimes or analysis inside Savepoint.

## Boundaries

**In scope:** compatibility, migrations, packaging, privacy, security, performance, platforms, docs, regression evidence, Design reconciliation, and the mandatory Full Objective Check.

**Out of scope:** more providers, dead-code analysis, generic plugins, cloud services, CI/CD integrations, background monitoring, and unrelated features.

## Confirmed Design Decisions

The owner confirmed modularity, project-owned tools, local-first operation, single-binary distribution, explicit upgrades, bounded persisted output, explicit retention maintenance, and the seven-Objective structure on 2026-09-26.
