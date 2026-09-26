---
id: O-029
title: Measure all six Code Health signals
status: planned
depends_on: [O-026, O-027, O-028]
release: R-007
priority: high
rank: 2
---

# O-029: Measure all six Code Health signals

## Outcome

Savepoint normalizes trustworthy results for tests, coverage, complexity, duplication, change hotspots, and known dependency vulnerabilities through the six official provider paths selected by O-026.

## Why

The v2.1 promise requires all six signals to have real supported implementations while preserving honest unavailability and keeping analysis outside Savepoint.

## Success Conditions

- Tests use structured counts when available and truthful gate status otherwise; any failed test needs attention and blocks by default.
- Coverage favours understandable line coverage while preserving supported technical detail without instrumentation.
- Complexity reports provider-defined difficult areas and project-relative change without universal cross-language claims.
- Duplication reports repeated-code share without adding a mandatory runtime to Savepoint.
- Change hotspots come from an approved external tool that owns the Git/complexity analysis.
- Security reports dependency findings as critical, high, medium, low, or unknown; high and critical block by default, while database freshness and scan completeness remain explicit.
- Every adapter provides bounded explanations and affected items, records provenance, preserves component scope, and has success, unavailable, malformed, failure, timeout, and mixed-language fixtures.
- One provider failure never destroys other valid results or becomes a bad-code classification.

## Architectural Considerations

Adapters implement Code Health interfaces and cannot leak provider schemas into persistence, Check, or TUI packages. External tools own analysis algorithms and vulnerability data.

## Boundaries

**In scope:** one tested official adapter path for each signal, normalization, details, fixtures, and bounded integration tests.

**Out of scope:** installation, language analysis inside Savepoint, dead code, full application security, linting, formatting, architecture analysis, generic adapters, and network-dependent normal tests.

## Confirmed Design Decisions

The owner explicitly required all six supported paths, one official provider per signal, project-owned external tools, explicit unknown vulnerability severity, and externally calculated change hotspots on 2026-09-26.
