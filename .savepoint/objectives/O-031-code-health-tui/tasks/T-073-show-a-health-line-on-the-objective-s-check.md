---
id: T-073
title: Show a health line on the Objective's Check
objective: O-031
status: done
depends_on: [{task: T-072, requires: clear}]
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: planning-o031-20261002}
complexity_tier: low
complexity_reason: One read-only line in an existing overlay, sourced from one snapshot lookup.
check_waiver:
    task: T-073
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-01T20:37:01Z"
---

# Show a health line on the Objective's Check

## Outcome

When an Objective's latest Check recorded an official health snapshot, the Objective detail shows one plain Health line with that snapshot's overall label and a pointer to the Code Health screen, and nothing more.

## User Check

None beyond the Full Objective Check.

## Done When

- The Objective detail's CHECKS section shows `Health: <label> (snapshot <short id>) — press H for details` when the latest Check names a `health_snapshot` that is stored.
- A reference to a missing snapshot reads `Health: snapshot not found — run savepoint doctor`; a Check without the field shows no Health line. Task details are unchanged.
- The label comes from a small `codehealth` lookup read during the board's load command, not during rendering; the board imports no snapshot type. A missing or unreadable health directory is not an error for the board.
- The line does not change badges, Next, clearance, or the non-TTY plain table.
- Tests cover stored official snapshot, missing snapshot, no field, and no health directory.

## Context Files

`.savepoint/objectives/O-031-code-health-tui/Objective.md`; `internal/board/v2/detail.go`, `internal/board/v2/detail_view.go`, `internal/board/v2/detail_test.go`, `internal/board/v2/load.go`, `internal/board/v2/load_test.go`; `internal/codehealth/dashboard.go`, `internal/codehealth/dashboard_test.go`; `internal/data/check_v2.go` (read only).

## Design References

Design section 8 (Check code style review, as the pattern for a detail-only line) and section 7; O-031 Confirmed Design Decisions (Check summary).

## Guardrails

ARCH-02, DATA-03, STYLE-07, TEST-01, TEST-02, TEST-04.

## Implementation Plan

1. Confirm `CheckV2.HealthSnapshot` and the T-072 `H` key exist; return REPLAN REQUIRED otherwise.
2. Add a `codehealth` lookup that returns the overall label for each stored snapshot ID.
3. Read it in `loadProject` and carry it on the project state.
4. Add the line in the Objective detail and tests.

## Boundaries

No full health rendering in the detail, no change to Check records or doctor.

## Technical Verification

`make build && make test-fast` for handoff. Later evaluated per `agent-skills/references/check-method.md`.

## Technical Evidence

Commands: `make build && make test-fast` — passed, no failures (2026-10-02).

Per criterion:
- Stored snapshot: Objective detail CHECKS shows `Health: Watch (snapshot <8 hex>) — press H for details` — `TestObjectiveDetailShowsTheHealthLineForAStoredSnapshot`.
- Missing snapshot reads `Health: snapshot not found — run savepoint doctor`; no field shows no line; Task detail unchanged — `TestObjectiveDetailReportsAMissingHealthSnapshot`, `...WithoutAHealthSnapshotFieldShowsNoHealthLine`, `TestTaskDetailNeverShowsAHealthLine`.
- Lookup is `codehealth.SnapshotLabels`, read in `loadProject` only when some Check names a snapshot, carried as `ProjectState.Health`; the board imports no snapshot type. Missing/unreadable health storage yields no labels and no error — `TestMissingHealthDirectoryIsNotABoardError`, `TestSnapshotLabelsNameEachStoredSnapshotAndTolerateMissingStorage`.
- Badges, Next, clearance, and the plain table are untouched (no edits to those paths; existing tests pass).

Files read: Context Files plus `internal/doctor/health_snapshot_refs.go` and its test, `internal/codehealth/storage.go`, `internal/board/v2/update.go` (extra reads, to mirror doctor's store path and update the detail call sites).
Files changed: `internal/codehealth/dashboard.go`, `internal/board/v2/{load,detail,detail_view,update}.go`, plus tests in `dashboard_test.go` and `detail_test.go`.

Limitations: the lookup does not check snapshot origin (doctor owns that); not run in a live TTY; `make test-full` not run.

## Drift Notes

None expected beyond Design section 8 wording reconciled at the Objective Check.
