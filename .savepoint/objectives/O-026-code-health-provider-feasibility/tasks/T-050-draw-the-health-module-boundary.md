---
id: T-050
title: Draw the health module boundary
objective: O-026
status: planned
depends_on: [{task: T-049, requires: clear}]
owner_validation: {required: true}
planned_by: {role: planner, session: planning-2026-09-26-o026}
---

# Draw the health module boundary

## Outcome

The Objective specifies narrow discovery, execution, normalization, snapshot, and dashboard contracts for a separate Code Health module, plus compatibility rules that determine when a provider or configuration change starts a new comparison series.

## User Check

Review the interface proposal against the confirmed provider catalogue and verify that unsupported and failed capabilities remain visible. Confirm the module boundary and comparability decisions. An optional Task Check may inspect the contract; if skipped, record the owner waiver in Task evidence.

## Done When

- The contract names the Code Health module's owned inputs, outputs, storage boundary, and error states without implementing a generic plugin framework.
- Discovery yields proposals for owner confirmation; execution accepts structured executable and argument values, timeouts, cancellation, and project-owned configuration without shell interpretation or automatic installation.
- Normalization preserves capability, scope, provider/version, measurement definition, exclusions, provenance, absence and failure states, and bounded sanitized evidence references.
- Snapshot and dashboard consumers receive a narrow read model; Objective Checks reference snapshot identities, while collection remains limited to explicit refresh and Full Objective Checks.
- Compatibility rules cover provider identity and version, report definition, scope, exclusions, and relevant configuration, with old series retained and incomparable points never trended together.
- The contract identifies which decisions O-027 needs for its versioned schema and tests, with the owner confirmation and configured Task gate result recorded.

## Context Files

`.savepoint/objectives/O-026-code-health-provider-feasibility/Objective.md`; `.savepoint/objectives/O-027-code-health-model-and-history/Objective.md`; `.savepoint/Design.md`; `.savepoint/Guardrails.md`; `internal/data/check_v2.go`; `internal/data/objective_gate_v2.go`; `internal/board/v2/load.go`; `internal/doctor/v2_runtime.go`.

## Design References

Design sections 1, 3, 8, 11, and 13; R-007 Confirmed Design Decisions; O-027 Architectural Considerations.

## Guardrails

FS-01, DATA-01, DATA-02, DATA-03, ARCH-01, ARCH-02, ARCH-03, ARCH-04, CFG-02, CFG-03, DEP-01, TEST-01, TEST-02, TEST-08.

## Implementation Plan

1. Read T-049's confirmed catalogue and inspect the targeted current Check, board, and doctor boundaries.
2. Define the smallest named requests and read results for discovery, execution, normalization, snapshots, and dashboard consumption, including error and unsupported outcomes.
3. Define stable comparison identity and explicit series reset conditions using the selected provider output contracts.
4. Record the proposed interfaces, compatibility table, O-027 handoff, and owner confirmation in the Objective.

## Boundaries

No production adapters, health schema implementation, scanner execution, dashboard, Check wiring, custom-provider API, or generic plugin framework.

## Technical Verification

Trace each selected report field through the proposed normalized result and read model; walk one valid, one failed or unavailable, and one incompatible-series example. Record the configured Task gate result. A later Check follows `agent-skills/references/check-method.md`.

## Technical Evidence

Pending execution: interface table, provider-to-field trace, series reset examples, current-source basis, owner decision, gate command and result, and limitations.

## Drift Notes

Record any boundary that cannot be honored by the selected provider path and return it for planning before O-027 is detailed.
