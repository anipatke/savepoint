---
id: T-049
title: Verify the six health inputs
objective: O-026
status: done
depends_on: []
owner_validation:
  required: true
  accepted_check: C-935
  accepted_by: {role: owner, session: user}
last_check: C-935
planned_by: {role: planner, session: planning-2026-09-26-o026}
---

# Verify the six health inputs

## Outcome

The Objective records an evidence-backed catalogue that chooses one official project-owned tool or structured report path for each of tests, coverage, complexity, duplication, change hotspots, and dependency vulnerabilities, with explicit limits for Go, JavaScript, TypeScript, and Python.

## User Check

Review the proposed official paths, unsupported stack and capability combinations, installation ownership, and any runtime or licence implications. Confirm the catalogue before the interface Task uses it. An optional Task Check may independently inspect the cited evidence; if skipped, record the owner waiver in Task evidence.

## Done When

- Each capability has at least two plausible candidates compared where alternatives exist, and one official path selected with an explicit rationale and unsupported cases.
- The comparison covers stack support, deterministic machine output and schema stability, licence, platform and runtime requirements, offline behavior, bounded performance and report size, and project-owned installation.
- Candidate output contracts name fields needed for normalized values, scopes, failure states, freshness, provenance, and sanitized references; an available representative report or documented sample is mapped to each selected path.
- The catalogue states how generated, vendored, and third-party code is excluded and identifies fixture cases for valid, absent, malformed, partial, and unsupported output.
- Duplication never creates a mandatory Savepoint runtime; complexity and hotspots remain external; no scanner installation or production adapter is added.
- Owner confirmation of the selected catalogue and the configured Task gate result are recorded as evidence.

## Context Files

`.savepoint/objectives/O-026-code-health-provider-feasibility/Objective.md`; `.savepoint/releases/R-007-v2-1-code-health/Release.md`; `.savepoint/Guardrails.md`; `go.mod`.

## Design References

Design sections 1, 3, 12, and 13; R-007 Success Conditions and Confirmed Design Decisions.

## Guardrails

ARCH-03, ARCH-04, CFG-02, CFG-03, DEP-01, TEST-01, TEST-02, TEST-08.

## Implementation Plan

1. Verify current primary tool documentation, licence terms, supported language and platform claims, report formats, and runtime requirements for the four confirmed stacks.
2. Compare alternatives per capability in a compact decision table, noting unsupported cases rather than treating missing evidence as success.
3. Inspect representative machine output or authoritative schemas without installing tools; record exact fields, output size and performance bounds, and limitations of any unverified claims.
4. Choose the official paths and record the catalogue, source links, rationale, fixture strategy, and owner decision in the Objective.

## Boundaries

No scanner installation, production code, analysis algorithm, dashboard, Check integration, custom-provider API, or generic plugin framework.

## Technical Verification

Cross-check selected report fields against cited primary schemas or representative reports and include a malformed or unsupported case. Record the configured Task gate result. A later Check follows `agent-skills/references/check-method.md`.

## Technical Evidence

Build research completed on 2026-09-27. The scoped Context Files were read. Additional reads were `.savepoint/config.yml` (lint, typecheck, and test gates are `null`; `block_on_failure: true`) and `.savepoint/Design.md` sections 1, 3, 12, and 13, as required by the Task's Design References. The proposed six-capability catalogue, alternative comparisons, field mapping, support gaps, exclusions, fixture cases, source URLs, and resource-cap proposal are recorded in O-026's “Proposed Provider Catalogue (T049, 2026-09-27)”.

Primary documentation and license/runtime details were reviewed without installing or executing providers. Representative documented samples were mapped for Vitest JUnit, pytest JUnit, Go coverage profiles, Lizard output, jscpd JSON, and OSV-Scanner JSON. The git-hotspots alpha exposes JSON mode and file-level history evidence, but its schema stability and non-Linux platform support remain unverified. No real-repository performance measurements or local provider fixtures were run; the 120-second, 32-MiB, and 10,000-row ceilings are proposals only. Owner confirmation is pending, especially for the git-hotspots alpha/platform gap, Code Maat's GPL-3.0/Java alternative, and the proposed caps.

Criterion evidence:

1. Candidate comparison and selection: proposed paths and alternatives for all six capabilities are recorded in O-026; change-hotspot provider remains an owner decision due alpha/platform constraints.
2. Stack, output, license, platform/runtime, offline, performance, report size, and installation: provider-specific support and constraints are in the catalogue. Proposed 120-second, 32-MiB, and 10,000-row caps are clearly marked unmeasured policy proposals.
3. Normalized fields and representative mapping: documented Go/Vitest/pytest test and coverage reports, Lizard CSV, jscpd JSON, and OSV-Scanner JSON are mapped to normalized fields. No provider binaries were installed or executed.
4. Exclusions and output fixtures: catalogue defines project-owned exclusions and valid, absent, malformed, partial, unsupported, stale, and excluded-source fixture cases.
5. Architecture boundaries: no scanner was installed; no production code or adapter was added; tests/coverage are consumed as existing artifacts; complexity and hotspots remain external; jscpd v5 requires no Savepoint runtime.
6. Owner decision and gate: owner accepted Check C-935 in conversation on 2026-09-27. `make build && make test-fast` passed at 2026-09-27 07:03 Australia/Sydney (Go 1.26.2, linux/amd64; exit code 0).

Check C-935 was recorded CLEAR on 2026-09-27. Owner acceptance is recorded. Only the owner may mark the Task done.

`./savepoint resume` confirms Task T-049 is ready for owner closure. The Full Objective Check remains mandatory before O-026 can close.

Owner decision (2026-09-29): the owner explicitly approved all five in-scope provider/report paths and the proposed module boundary after C-937 identified their pending status as I-084. O-026 now records the paths as the selected official catalogue. The historical hotspot research remains excluded from V2.1. This metadata-only repair changes no code, tests, fixtures, dependencies, or gate definitions; it reuses the successful `make test-full` run recorded by C-937 at 2026-09-29T10:49:37Z.

## Drift Notes

Record any discovered constraint that changes the Objective's scope before the interface Task proceeds.
