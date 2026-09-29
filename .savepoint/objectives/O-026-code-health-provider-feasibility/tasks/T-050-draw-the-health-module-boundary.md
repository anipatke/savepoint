---
id: T-050
title: Draw the health module boundary
objective: O-026
status: done
depends_on: [{task: T-049, requires: clear}]
owner_validation:
  required: true
  accepted_check: C-936
  accepted_by: {role: owner, session: user}
last_check: C-936
planned_by: {role: planner, session: planning-2026-09-26-o026}
---

# Draw the health module boundary

## Outcome

The Objective specifies narrow discovery, execution, normalization, snapshot, and dashboard contracts for a separate Code Health module covering the five V2.1 capabilities—tests, coverage, complexity, duplication, and dependency vulnerabilities—plus compatibility rules that determine when a provider or configuration change starts a new comparison series.

## User Check

Review the interface proposal against T-049's five provisional provider paths and verify that unsupported and failed capabilities remain visible. Change hotspots and dead-code measurement are out of scope for V2.1. Confirm the provider paths, module boundary, and comparability decisions. An optional Task Check may inspect the contract; if skipped, record the owner waiver in Task evidence.

## Done When

- The contract names the Code Health module's owned inputs, outputs, storage boundary, and error states without implementing a generic plugin framework.
- Discovery yields proposals for owner confirmation; execution accepts structured executable and argument values, timeouts, cancellation, and project-owned configuration without shell interpretation or automatic installation.
- Normalization preserves capability, scope, provider/version, measurement definition, exclusions, provenance, absence and failure states, and bounded sanitized evidence references.
- Snapshot and dashboard consumers receive a narrow read model; Objective Checks reference snapshot identities, while collection remains limited to explicit refresh and Full Objective Checks.
- Compatibility rules cover provider identity and version, report definition, scope, exclusions, and relevant configuration, with old series retained and incomparable points never trended together.
- The V2.1 contract contains no hotspot provider, collection, storage, trend, or health-score input.
- The contract identifies which decisions O-027 needs for its versioned schema and tests, with the owner confirmation and configured Task gate result recorded.

## Context Files

`.savepoint/objectives/O-026-code-health-provider-feasibility/Objective.md`; `.savepoint/objectives/O-027-code-health-model-and-history/Objective.md`; `.savepoint/Design.md`; `.savepoint/Guardrails.md`; `internal/data/check_v2.go`; `internal/data/objective_gate_v2.go`; `internal/board/v2/load.go`; `internal/doctor/v2_runtime.go`.

## Design References

Design sections 1, 3, 8, 11, and 13; R-007 Confirmed Design Decisions; O-027 Architectural Considerations.

## Guardrails

FS-01, DATA-01, DATA-02, DATA-03, ARCH-01, ARCH-02, ARCH-03, ARCH-04, CFG-02, CFG-03, DEP-01, TEST-01, TEST-02, TEST-08.

## Implementation Plan

1. Read T-049's five provisional provider paths and inspect the targeted current Check, board, and doctor boundaries. Keep provider-specific assumptions explicit and do not present proposals as owner-approved selections.
2. Define the smallest named requests and read results for discovery, execution, normalization, snapshots, and dashboard consumption, including error and unsupported outcomes.
3. Define stable comparison identity and explicit series reset conditions using the proposed provider output contracts; show where a provider decision could change the contract.
4. Record the proposed interfaces, compatibility table, O-027 handoff, and owner confirmation in the Objective.

## Boundaries

No production adapters, health schema implementation, scanner execution, dashboard, Check wiring, custom-provider API, or generic plugin framework.

## Technical Verification

Trace each selected report field through the proposed normalized result and read model; walk one valid, one failed or unavailable, and one incompatible-series example. Record the configured Task gate result. A later Check follows `agent-skills/references/check-method.md`.

## Technical Evidence

REPLAN REQUIRED before settling the interface contract. `./savepoint resume` reported that T-050's dependency is satisfied. However, the allowed O-026 context labels the T-049 provider catalogue as awaiting owner confirmation and leaves the change-hotspot provider decision open. The T-050 plan specifically assumes a confirmed catalogue, and its User Check requires review against one. Provider selection changes the supported platform and report identity that the normalization and comparison contract must describe.

The planner needs to resolve the ordering: record owner confirmation of the T-049 catalogue before T-050 resumes, or revise T-050's basis to permit a conditional proposal over unconfirmed provider choices. No interface contract was drafted, no production files were changed, and no configured Task gate was run. `status: in_progress` and `stage: build` are unchanged from the start of this attempt.

Files read: T-050, O-026, O-027, Design sections 1, 3, 8, 11, and 13, Guardrails, and targeted symbol/consumer declarations in the four listed Go Context Files. Files changed: T-050 only.

Planner replan update (2026-09-27): the owner directed that hotspots be de-scoped for V2.1. O-026 now excludes hotspot measurement and derived signals from the V2.1 scope, and this Task's outcome and acceptance criteria cover five capabilities. The remaining five proposed provider paths are still marked owner-confirmation pending in O-026. They must be confirmed before provider-specific interfaces can be settled; the replan remains active. No implementation resumed, and Task status/stage remain unchanged.

Owner decision and replan resolution (2026-09-27): the owner confirmed the five-measure V2.1 scope (tests, coverage, complexity, duplication, and dependency vulnerabilities), excluded dead code, and directed work to proceed to T-050. The planner resolved the replan by authorizing a conditional interface proposal using T-049's five provider paths as provisional inputs. Provider selection remains pending owner review; the proposal must identify assumptions and any decisions that would change its contracts. Hotspots and dead code are excluded from this release. No Task lifecycle status was changed.

T050 design proposal (2026-09-27): O-026 now contains a proposed Code Health boundary, storage ownership, discovery/collection/read requests, caller responsibilities, separate result/error states, compatibility identity, and O-027 handoff. The provider paths remain provisional; this proposal does not select them. No production code or data schema was added.

Context reads: T-050, O-026, O-027, `.savepoint/Design.md` sections 1, 3, 8, 11, and 13, `.savepoint/Guardrails.md` rules ARCH-01..04, CFG-02..03, DEP-01, TEST-01..02 and TEST-08, and the four named Go Context Files (`check_v2.go`, `objective_gate_v2.go`, `load.go`, and `v2_runtime.go`). Extra read: `.savepoint/releases/R-007-v2-1-code-health/Release.md`, because T-050's Design References name its Success Conditions and Confirmed Design Decisions and the owner requested that the five-measure scope be enshrined there.

Criterion evidence:

1. Proposed `internal/codehealth` ownership, config/snapshot storage paths, caller boundaries, and error states are recorded in O-026.
2. Discovery, sequential collection, structured executable/argument invocation, timeout/cancellation, and no-auto-install boundaries are specified in O-026.
3. Normalized result provenance, repository state, configuration/scope identity, bounded evidence, and the no-zero-on-missing rule are specified in the T049 catalogue and T050 contract.
4. Official/manual snapshot distinctions, snapshot references from Checks, and board read-only queries are specified; the current `CheckV2` has no snapshot reference, so O-027 must define it. Board rendering remains IO-free; doctor does not collect health data.
5. Series reset keys and vulnerability-database comparability treatment are specified in O-026; unsupported configurations do not silently share a series.
6. O-027's outstanding schema, ID, retention, freshness, trend, and fixture decisions are named. Owner confirmed the five-measure scope, but has not yet confirmed the five proposed provider paths or this interface proposal; final owner validation remains pending.

Files changed for this Task: `.savepoint/objectives/O-026-code-health-provider-feasibility/Objective.md`, `.savepoint/objectives/O-026-code-health-provider-feasibility/tasks/T-050-draw-the-health-module-boundary.md`, and `.savepoint/releases/R-007-v2-1-code-health/Release.md`. `make build && make test-fast` passed on 2026-09-27 09:57 Australia/Sydney with Go 1.26.2 linux/amd64 (exit code 0). Limits remain provisional and unmeasured; provider paths remain proposals; O-027's current text still says six signals pending its later reconciliation. The interface and five provider paths await owner review. The optional Task Check has not been requested or waived; T-050 is ready for Check-stage handling, not completion.

Owner decision (2026-09-29): the owner explicitly approved the five provider/report paths, proposed module boundary, and compatibility rules after C-937 identified their pending status as I-084. O-026 now records the boundary as confirmed. The documented limitations and unmeasured resource caps remain explicit and are not converted into measured claims. This metadata-only repair changes no code, tests, fixtures, dependencies, or gate definitions; it reuses the successful `make test-full` run recorded by C-937 at 2026-09-29T10:49:37Z.

## Drift Notes

Record any boundary that cannot be honored by the selected provider path and return it for planning before O-027 is detailed.
