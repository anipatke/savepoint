---
id: T-051
title: Give health evidence one precise vocabulary
objective: O-027
status: planned
depends_on: []
owner_validation: {required: false}
planned_by: {role: planner, session: planning-o027-20260929}
---

# Give health evidence one precise vocabulary

## Outcome

`internal/codehealth` defines the versioned, stack-agnostic configuration and snapshot contracts that every later Code Health component can validate and interpret without turning missing, failed, stale, or incomparable evidence into a healthy value.

## User Check

No separate owner validation is required. An optional Task Check may inspect the schema invariants; if skipped, the owner records the required waiver when completing the Task.

## Done When

- The fixed capability set is tests, coverage, complexity, duplication, and dependency vulnerabilities; unknown capabilities and duplicate configured instances receive named validation errors.
- Collection outcome and freshness are orthogonal typed values, including every owner-confirmed state, and only an available or partial result may carry measured values.
- Versioned configuration, snapshot, capability-result, provider provenance, repository identity, bounded evidence-reference, origin, threshold, and summary types have explicit validation rules.
- Canonical serialization produces stable snapshot identities and series identities; series identity includes provider, schema, measurement definition, effective configuration, exclusions, and scope, but excludes the observed commit.
- Validation rejects absolute or sensitive evidence references, unbounded details, incompatible values/units, malformed timestamps and identities, and impossible origin/retention combinations.
- Fixtures and table-driven tests cover valid records plus missing, partial, unsupported, unavailable, failed, timed-out, cancelled, stale, unknown-freshness, malformed, duplicate, and incompatible examples.

## Context Files

`.savepoint/objectives/O-027-code-health-model-and-history/Objective.md`; `.savepoint/objectives/O-026-code-health-provider-feasibility/Objective.md`; `.savepoint/Design.md`; `.savepoint/Guardrails.md`; `go.mod`; `internal/codehealth/model.go`; `internal/codehealth/model_test.go`; `internal/codehealth/testdata/valid-config-v1.json`; `internal/codehealth/testdata/valid-snapshot-v1.json`.

## Design References

Design sections 1, 2, 3, 4, 11, and 13; O-026 Confirmed Code Health Module Boundary and series compatibility decisions.

## Guardrails

DATA-02, DATA-03, ARCH-03, ARCH-04, CFG-01, CFG-02, DEP-01, TEST-01, TEST-02, TEST-04, TEST-08.

## Implementation Plan

1. Add the dedicated `internal/codehealth` package with closed capability, outcome, freshness, origin, and classification vocabularies.
2. Define versioned configuration and snapshot records with bounded normalized detail and sanitized repository-relative evidence references.
3. Define canonical content and hashing rules for snapshot and comparison-series identities, keeping observation identity separate from compatibility identity.
4. Add validation that rejects impossible state/value combinations, unsafe paths, duplicate instances, invalid bounds, and unsupported schema versions with named errors.
5. Add deterministic JSON fixtures and parameterized tests covering every state and boundary, including mutation-after-validation through a fresh validation call.

## Boundaries

No provider execution, filesystem persistence, Git subprocesses, Check wiring, TUI rendering, automatic Issues, generic provider API, or composite score.

## Technical Verification

Run focused `internal/codehealth` model tests during iteration, then `make build && make test-fast`. A later Check follows `agent-skills/references/check-method.md`.

## Technical Evidence

Pending execution: record schema cases, identity vectors, commands, changed files, and limitations.

## Drift Notes

Return to planning if canonical identity requires storing source content, secrets, absolute host paths, or a sixth capability.
