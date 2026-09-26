---
id: R-007
title: Savepoint v2.1 — Code Health
status: planned
---

## Outcome

AI-assisted solo builders can see whether the code their agents produce is staying healthy, improving, or deteriorating through deterministic local measurements, plain-language explanations, historical trends, and a dedicated Savepoint TUI view.

## Why

Savepoint helps users plan, build, and independently verify work, but it does not translate ongoing software-health signals for users unfamiliar with coverage, complexity, duplication, change hotspots, or vulnerability reports. Code Health extends Objective Checks without turning Savepoint into a language-specific analysis platform.

## Success Conditions

- A separate internal Code Health module owns its normalized model, supported provider adapters, execution, classification, snapshot history, and persistence; the TUI and Objective Check use narrow interfaces.
- Initial scaffolding discovers likely capabilities, proposes them for confirmation, records human-readable intentions in `Design.md`, and keeps executable settings in dedicated machine-readable configuration.
- Tests, coverage, complexity, duplication, change hotspots, and known dependency vulnerabilities each have one supported official provider path using project-owned tools or structured reports. Missing capabilities remain explicitly unavailable.
- Full Objective Checks and explicit manual refreshes collect health data. Official trends use Full Objective Check snapshots; manual results remain visible without altering the official baseline.
- Compact versioned snapshots identify repository state and provider provenance, distinguish collection failure from unhealthy code, preserve comparable history, and never store unrestricted scanner output.
- The dedicated TUI view uses deterministic Good, Watch, and Needs Attention states without a composite score or AI assessment, explains every warning, and never starts a scan merely by opening.
- Required capabilities and explicit blocking rules govern Objective Check clearance. Stale, partial, failed, or absent data is never presented as unconditional health.
- Existing projects receive an explicit upgrade and design-reconciliation path; existing Savepoint behavior, local-first operation, Windows support, and the single-binary installation remain intact.

## Boundaries

**In scope:** modular health records and services; confirmed discovery and configuration; safe use of project-owned tools; six supported capabilities; history, classification, and staleness; Full Objective Check integration; a dedicated TUI view; upgrades, tests, and documentation.

**Out of scope:** AI review or health judgment, cloud analysis, language-specific analysis engines inside Savepoint, composite numerical scores, automatic tool installation, background monitoring, daemons, watchers, dead-code analysis as a core capability, automatic Issue creation, web dashboards, enterprise reporting, and a generic plugin ecosystem.

## Confirmed Design Decisions

Confirmed by the owner in chat on 2026-09-26.

- Code Health is a separate top-level internal module compiled into the main binary. It owns its domain and storage; only narrow TUI and Objective Check integrations cross the boundary.
- Savepoint uses project-owned tools and existing structured reports. It neither installs scanners nor implements analysis engines. Agents may scaffold deterministic configuration but never invent measurements or health judgments.
- Discovery occurs during initial scaffolding and proposes providers for confirmation rather than executing guesses. `Design.md` records human intent; exact providers, commands, paths, timeouts, exclusions, and thresholds live in dedicated Code Health configuration.
- Existing projects adopt Code Health through an explicit upgrade and Design reconciliation. Opening a project never writes or scans.
- Collection runs only during Full Objective Checks and explicit manual refreshes. Providers run sequentially in v2.1 with provider-specific default timeouts, overrides, cancellation, and graceful partial failure.
- Capabilities are required or optional. A required provider that cannot produce valid evidence blocks Objective Check clearance. Failing tests and high/critical vulnerabilities block by default; other health warnings block only through confirmed project policy.
- Existing gates own build and test execution. Code Health consumes fresh structured artifacts where possible instead of rerunning work.
- External tools are invoked without a shell through structured executable and argument values. Snapshots store normalized bounded results and sanitised evidence references, not full raw reports or unrestricted output.
- Official trends and baselines use comparable Full Objective Check snapshots. A recent range begins after three comparable snapshots. Only the ten newest manual snapshots are retained, and pruning is an explicit maintenance action rather than a collection side effect.
- Snapshots are versioned machine-readable project records intended for Git. Code Health owns them separately; Checks reference their identities. Commit identity plus a relevant working-tree fingerprint supports honest dirty-tree and stale-data reporting.
- Current values set the base state; a material decline may worsen it by one level, while improvement changes explanation without disguising a poor current value. Confirmed project overrides are allowed except for hard minimum rules.
- Healthy requires all required capabilities to complete with none needing attention. Partial, stale, failed, and absent data receive explicit alternative summaries.
- Generated, vendored, and third-party code is excluded by default through provider-native mechanisms. Scope and configuration fingerprints control comparability.
- Monorepos may configure multiple scoped provider instances per capability. Savepoint does not aggregate values that are not semantically compatible.
- All six capabilities ship with one tested official path selected from a supported catalogue. External tools remain project-owned. Complexity and change-hotspot calculation remain in approved external tools rather than Savepoint.
- Security database behavior belongs to the project scanner. Savepoint records freshness when available and preserves unknown severity instead of guessing.
- Health warnings never create Issues automatically; an independent checker decides whether durable Issue capture is warranted.
- The public concept is Goal, while the current schema retains the `R-###` compatibility identity. A future `G-###` migration is separate from v2.1.
