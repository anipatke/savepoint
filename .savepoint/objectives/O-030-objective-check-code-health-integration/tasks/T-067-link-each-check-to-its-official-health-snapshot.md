---
id: T-067
title: Link each Check to its official health snapshot
objective: O-030
status: planned
depends_on: []
owner_validation: {required: false}
planned_by: {role: planner, session: planning-o030-20261001}
complexity_tier: medium
complexity_reason: Strict decoder field plus a cross-module doctor diagnostic that must not couple data to codehealth.
---

# Link each Check to its official health snapshot

## Outcome

A Check record can name the official health snapshot it relied on, and `savepoint doctor` tells the owner when that name points at nothing or at a manual snapshot. The Check and health records stay separately owned.

## User Check

None beyond the Full Objective Check.

## Done When

- `DecodeCheckV2` accepts an optional `health_snapshot` string and exposes it on `CheckV2`. An empty, multi-line, or over-long value is a named malformed-Check diagnostic. Existing Checks without the field load unchanged, and the field never affects clearance resolution.
- `internal/data` does not import `internal/codehealth` (a test or `go list` check proves it).
- Doctor reports a Check whose `health_snapshot` names no stored snapshot, and one that names a `manual` snapshot, with distinct named diagnostics, the Check path, and a plain repair hint. A valid official reference produces no diagnostic. A missing health directory is not a finding for Checks without the field.
- Doctor stays read-only, with tests that use temporary projects.

## Context Files

`.savepoint/objectives/O-030-objective-check-code-health-integration/Objective.md`; `internal/data/check_v2.go`, `internal/data/check_v2_test.go`; `internal/doctor/checks.go`, `internal/doctor/report.go`, `internal/doctor/v2_runtime.go`; `internal/codehealth/storage.go`, `internal/codehealth/model.go`; new `internal/doctor/health_snapshot_refs.go`, `internal/doctor/health_snapshot_refs_test.go`.

## Design References

Design sections 1 (V2 evidence and identity boundary), 7, and 11; O-030 Architectural Considerations.

## Guardrails

DATA-01, DATA-03, DATA-04, ARCH-04, FS-03, STYLE-07, TEST-01, TEST-02, TEST-04.

## Implementation Plan

1. Confirm the Check decoder and the doctor problem pipeline match these paths; return REPLAN REQUIRED if doctor has no V2 Check iteration point.
2. Add the frontmatter field, shape validation, and decoder tests.
3. Add the doctor pass: load snapshots once through `codehealth.Store`, index them by ID, and report missing or manual references.
4. Wire it where the V2 doctor problems are assembled.

## Boundaries

No change to `ResolveTaskCompletion`, the Objective gate, Check immutability, or the snapshot schema. No collection and no wording in skills.

## Technical Verification

Focused tests while iterating; `make build && make test-fast` at handoff.

## Technical Evidence

Pending execution.

## Drift Notes

Doctor now reads Code Health storage; the Codebase Map line for `internal/doctor` is updated in the skill-and-docs Task.
