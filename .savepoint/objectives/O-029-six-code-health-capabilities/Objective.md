---
id: O-029
title: Measure all five Code Health signals
status: in_progress
depends_on: [O-026, O-027, O-028]
release: R-007
priority: high
rank: 2
---

# O-029: Measure all five Code Health signals

## Outcome

Savepoint reads the reports of the approved provider paths for tests, coverage, complexity, duplication, and dependency vulnerabilities and turns each into a trustworthy normalized result through the reader seam O-028 left in `Collect`.

## Why

The v2.1 promise requires every in-scope measure to have a real supported implementation while preserving honest unavailability and keeping analysis outside Savepoint. Until the readers exist, every configured instance is reported as `unsupported`.

## Success Conditions

- Tests read Go test JSON, Vitest JUnit, and pytest JUnit. The value is the failed-test count; totals, skips, and errors are supporting details; failing tests are named as affected items. A Go package that fails to build counts as one failure with a note naming it. No structured report is `absent`, never a pass.
- Coverage reads Go cover profiles, Vitest V8 `coverage-final.json`, and coverage.py JSON. The headline value is statement coverage on every stack; branch and function counts are supporting details where the report has them. The least-covered files are the affected items.
- Complexity reads Lizard CSV. The value is the highest function cyclomatic complexity; the most complex functions are the affected items. Unparsed files are reported as partial, never as complexity zero.
- Duplication reads jscpd JSON. The value is the duplicated-lines share; the largest clones are the affected items. jscpd's own health score, complexity, and dead-code output are ignored.
- Dependency vulnerabilities read OSV-Scanner JSON. The value is the total; critical, high, medium, low, and unknown counts are recorded. Only high, critical, and unknown severity block (O-027 rule). The scanner database freshness and partial scans (manifests without resolved versions) stay explicit.
- Every reader returns repository-relative paths only, bounded details and affected items, and provider/version provenance where the report exposes it, and preserves the instance scope. Each reader has fixtures for a valid report, an empty one, a malformed one, a partial one, an unsupported variant, and a mixed-language project.
- The production readers are registered in one place. An end-to-end collection over fixture reports proves that one provider's failure, timeout, or malformed report never removes or alters another instance's valid result and never becomes a bad-code classification.

## Architectural Considerations

Readers live in `internal/codehealth` behind the existing `Reader` interface and never leak provider schemas into persistence, Check, or TUI packages. External tools own the analysis algorithms and vulnerability data. The snapshot schema stays at version 1; the only model change is admitting the medium, low, and unknown severity detail keys beside high and critical. There is still no collection command: O-030 (Check) and O-031 (TUI refresh) call `Collect` with the registered readers.

## Boundaries

**In scope:** one tested reader per approved provider path, normalization, details, affected items, fixtures, the shared reader groundwork they need, registration, and bounded integration tests.

**Out of scope:** change hotspots and dead code (de-scoped 2026-09-27), installation, running tests or tools in tests, language analysis inside Savepoint, full application security, linting, formatting, architecture analysis, generic adapters, a collection command, and network-dependent tests.

## Confirmed Design Decisions

The owner explicitly required one official provider per signal, project-owned external tools, and explicit unknown vulnerability severity on 2026-09-26. The 2026-09-27 de-scoping of change hotspots reduced this Objective from six signals to five.

On 2026-10-01, the owner confirmed the O-029 detail design:

- **Tests without a structured report** are `absent`; there is no fallback to gate pass/fail.
- **Coverage headline** is statement coverage on every stack, since Go profiles only measure statements; branch and function counts are supporting detail.
- **Vulnerability severity** records all five buckets (critical, high, medium, low, unknown). Only high, critical, and unknown block; medium and low are shown, not blocking.
- **Go build failures** count as failing tests, one per package that failed to build, with a note naming the package.

Technical decisions recorded with that confirmation: OSV severity buckets come from each vulnerability group's `max_severity` CVSS score using the CVSS qualitative ranges (9.0+ critical, 7.0+ high, 4.0+ medium, above 0 low; no score is unknown); `ReportInput` gains the project root so readers can turn absolute and Go import paths into repository-relative ones; affected items use the existing `EvidenceRef` (path, line, note) within the existing bounds, so the snapshot schema stays at version 1.
