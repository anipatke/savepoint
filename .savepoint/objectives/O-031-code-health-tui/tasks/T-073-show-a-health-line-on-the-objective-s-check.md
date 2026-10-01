---
id: T-073
title: Show a health line on the Objective's Check
objective: O-031
status: planned
depends_on: [{task: T-072, requires: clear}]
owner_validation: {required: false}
planned_by: {role: planner, session: planning-o031-20261002}
complexity_tier: low
complexity_reason: One read-only line in an existing overlay, sourced from one snapshot lookup.
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

Pending execution.

## Drift Notes

None expected beyond Design section 8 wording reconciled at the Objective Check.
