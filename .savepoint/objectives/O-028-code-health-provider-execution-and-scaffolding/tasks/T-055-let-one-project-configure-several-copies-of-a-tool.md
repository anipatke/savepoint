---
id: T-055
title: Let one project configure several copies of a tool
objective: O-028
status: planned
depends_on: []
owner_validation: {required: false}
planned_by: {role: planner, session: planning-o028-20261001}
---

# Let one project configure several copies of a tool

## Outcome

The Code Health configuration can hold several named, scoped instances of one capability and provider, mark each instance required or optional, and resolve a provider-specific default timeout, and every snapshot result says which instance produced it, without changing existing comparison series.

## User Check

None beyond the Full Objective Check; this is a schema extension with no owner-visible surface yet.

## Done When

- `CapabilityConfig` gains an optional `name` (short token) and `required` (bool, default false). The instance key becomes capability + provider + name. Two instances with the same key are still refused with `ErrDuplicateInstance`; two instances sharing capability and provider must each have a non-empty `name` and non-identical `scope`, or validation names the problem.
- `CapabilityResult` and `CapabilitySummary` carry the instance `name` (omitted when empty), so existing snapshot identity vectors and series identities are unchanged for unnamed instances. Renaming an instance does not start a new series; scope, exclusions, and the config digest still do.
- The catalogue gives every provider a default timeout (Lizard and jscpd 60s, OSV-Scanner 120s; report-only providers for tests and coverage have none because nothing is executed). An `EffectiveTimeout()` helper returns the override when set, else the default. `MaxTimeoutSecond` rises to 600; out-of-range values are refused.
- Config and snapshot `version` stay 1. Unknown fields are still rejected.
- Classification and history filtering treat each instance separately: two instances never share a trend or combine values.

## Context Files

`.savepoint/objectives/O-028-code-health-provider-execution-and-scaffolding/Objective.md`; `internal/codehealth/model.go`; `internal/codehealth/config.go`; `internal/codehealth/snapshot.go`; `internal/codehealth/identity.go`; `internal/codehealth/classification.go`; `internal/codehealth/model_test.go`; `internal/codehealth/identity_test.go`; `internal/codehealth/classification_test.go`.

## Design References

Design sections 1 and 12; O-026 "Confirmed Code Health Module Boundary"; O-028 Confirmed Design Decisions (2026-10-01).

## Guardrails

DATA-03, ARCH-04, CFG-01, STYLE-07, TEST-01, TEST-02, TEST-06, TEST-08.

## Implementation Plan

1. Add the provider default-timeout table beside `providerCapability` in `model.go`, with a lookup that reports "not executed" for report-only providers. Raise `MaxTimeoutSecond` to 600.
2. Extend `CapabilityConfig` with `Name` and `Required`; validate `Name` as a token; change `instanceKey` to include it; add the shared-provider rule (non-empty distinct names, non-identical scope).
3. Add `Name` (omitempty) to `CapabilityResult` and `CapabilitySummary`, validated as a token; keep it out of `SeriesID` and the config `Digest`.
4. Confirm the snapshot ID vector test still passes unchanged; add a vector with a named instance.
5. Make sure `Assess` history filtering keys on series identity per instance, and `Overall` never merges values across instances.
6. Tests: duplicate named/unnamed instances, same scope refused, name round-trip, timeout default and override and bounds, series unchanged by rename, two instances kept apart in history.

## Boundaries

No discovery, execution, storage layout, CLI, or report parsing. No schema version bump. No change to classification thresholds.

## Technical Verification

Focused `make test-focused TEST=./internal/codehealth/...` while iterating; `make build && make test-fast` at handoff.

## Technical Evidence

Pending execution.

## Drift Notes

None expected; the AGENTS.md Codebase Map entry for `internal/codehealth` already covers configuration records.
