---
id: T-067
title: Link each Check to its official health snapshot
objective: O-030
status: done
depends_on: []
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: planning-o030-20261001}
complexity_tier: medium
complexity_reason: Strict decoder field plus a cross-module doctor diagnostic that must not couple data to codehealth.
check_waiver:
    task: T-067
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-01T19:53:55Z"
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

Per-criterion outcomes:

- Decoder field: `DecodeCheckV2` accepts optional `health_snapshot`, exposed as `CheckV2.HealthSnapshot`. Empty, blank, multi-line, and over-long (>128 bytes) values return `ErrV2CheckMalformed` naming `health_snapshot`; absent loads unchanged and the field is not read by clearance resolution. Proven by `TestDecodeCheckV2_healthSnapshot` (and the existing decoder tests still pass).
- No data→codehealth coupling: `TestDataDoesNotImportCodeHealth` runs `go list -deps` on `internal/data`.
- Doctor diagnostics: `[health-snapshot-missing]` and `[health-snapshot-manual]` carry the Check path and a repair hint; an official reference is clean; a Check without the field never loads the health directory, so a missing directory is not a finding. Proven by `TestHealthSnapshotRefs_*` in `internal/doctor/health_snapshot_refs_test.go`.
- Read-only: `TestHealthSnapshotRefs_doctorStaysReadOnly` compares the full project tree before and after `RunV2Checks`; all tests use temporary projects.

Commands run: `go test ./internal/data`, `go test ./internal/doctor`, `make build`, `make test-fast` (all passed).

Files changed: `internal/data/check_v2.go`, `internal/data/check_v2_test.go`, `internal/doctor/health_snapshot_refs.go`, `internal/doctor/health_snapshot_refs_test.go`, `internal/doctor/repairs.go`, `internal/doctor/v2_runtime.go`.

Extra reads and edits beyond Context Files (recorded for transparency):

- `internal/doctor/repairs.go` was edited to add repair text for the three new diagnostic names (the plan's `V2ProblemRepair` home); `internal/doctor/checks_test.go` and `report_test.go` were read for test helpers; `internal/codehealth/snapshot.go` and `identity.go` were read to build valid test snapshots; `internal/data/project.go` was grepped for the index's Check map.

Limitations and design notes:

- A third diagnostic, `health-snapshot-unreadable`, is reported once when stored snapshots fail to load while some Check references one; the Task named only two. It is categorised Malformed Data; the missing and manual ones are Missing Evidence.
- `health_snapshot: null` decodes as absent; only an explicit empty string is malformed.
- No Task Check requested and no owner waiver recorded.

## Drift Notes

Doctor now reads Code Health storage; the Codebase Map line for `internal/doctor` is updated in the skill-and-docs Task.
