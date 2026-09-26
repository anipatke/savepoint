---
id: T-049
title: Verify the six health inputs
objective: O-026
status: planned
depends_on: []
owner_validation: {required: true}
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

Pending execution: source dates and links, candidate matrix, report-field mapping, unsupported cases, owner decision, gate command and result, and limitations.

## Drift Notes

Record any discovered constraint that changes the Objective's scope before the interface Task proceeds.
