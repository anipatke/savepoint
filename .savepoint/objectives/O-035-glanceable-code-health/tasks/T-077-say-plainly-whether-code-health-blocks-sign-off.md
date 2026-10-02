---
id: T-077
title: Say plainly whether Code Health blocks sign-off
objective: O-035
status: done
depends_on: [{task: T-075, requires: clear}]
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: planning-o035-20261002}
complexity_tier: low
complexity_reason: One lookup of the existing gate result per row plus three dashboard-level sentences and a date helper.
check_waiver:
    task: T-077
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-01T22:12:52Z"
---

# Say plainly whether Code Health blocks sign-off

## Outcome

The dashboard says in plain words whether the result blocks sign-off, using the same gate as `savepoint health check`, plus a one-sentence headline and friendly dates, so the screen and the command never disagree.

## User Check

None beyond the Full Objective Check.

## Done When

- `Dashboard` gains `Headline` ("2 of 5 need a look", "All 5 look fine", or "Not enough to judge yet"), `SignOff`, and `MeasuredText`; `DashboardHistoryEntry` gains `WhenText`. `DashboardRow` gains `SignOff`.
- When the newest snapshot is official, `Evaluate` runs on it with the saved configuration. Each row's `SignOff` is "Blocks sign-off" or "Advisory only" according to its disposition, and `Dashboard.SignOff` is "Blocks sign-off" or "Doesn't block sign-off".
- When the newest snapshot is a manual refresh, rows carry no `SignOff` and `Dashboard.SignOff` says a manual refresh does not affect sign-off. With no official snapshot it says there is no official check yet.
- An `Evaluate` error leaves `SignOff` as "Sign-off status unavailable" and loads the rest of the dashboard; it is never a panic or a silent default (DATA-03).
- A test asserts the dashboard's sign-off statements equal what `Evaluate` returns for the same snapshot, including the unknown-severity vulnerability case, which the gate reports only.
- `MeasuredText` and `WhenText` use one fixed UTC layout such as "1 Oct 20:53 UTC", from the stored time, not the machine's clock or zone, so tests are deterministic.
- Tests cover blocks, advisory, manual newest, no official, evaluate error, headline counts, and both date forms.

## Context Files

`.savepoint/objectives/O-035-glanceable-code-health/Objective.md`; `internal/codehealth/dashboard.go`, `internal/codehealth/dashboard_test.go`, `internal/codehealth/dashboard_copy.go`, `internal/codehealth/dashboard_copy_test.go`; `internal/codehealth/gate.go` (read only: `Evaluate`, `Verdict`).

## Design References

O-035 Confirmed Design Decisions (sign-off sentence from `Evaluate`); O-030 official snapshot gate.

## Guardrails

ARCH-03, DATA-03, STYLE-07, STYLE-09, TEST-01, TEST-02, TEST-04.

## Implementation Plan

1. Confirm the previous Task's fields exist and `Evaluate` is callable from `LoadDashboard` with the loaded config; return REPLAN REQUIRED otherwise.
2. Locate the newest official snapshot, evaluate it, and map each `ResultVerdict` to its row by instance key.
3. Add the headline, sign-off and date wording to the copy file and fill the new fields.
4. Add the tests in Done When.

## Boundaries

No change to `Evaluate`, classification, snapshot creation or `health check` output; no board rendering.

## Technical Verification

Focused tests during iteration (`make test-focused`); `make build && make test-fast` for ordinary handoff. Later evaluated per `agent-skills/references/check-method.md`.

## Technical Evidence

Per-criterion outcomes:
- Fields: `Dashboard.Headline/SignOff/MeasuredText`, `DashboardRow.SignOff`, `DashboardHistoryEntry.WhenText` added. Met.
- Official newest: `Evaluate` runs with the saved config; rows say "Blocks sign-off"/"Advisory only", dashboard "Blocks sign-off"/"Doesn't block sign-off". Met (`TestDashboardSignOffMatchesEvaluate`, `TestDashboardSignOffBlocksAndAdvisoryRows`).
- Manual newest / no official: rows empty, dashboard sentences as specified. Met (`TestDashboardSignOffManualNewest`, `TestDashboardSignOffNoOfficial`).
- Evaluate error: "Sign-off status unavailable" on dashboard and rows, rest loads, no panic. Met (`TestDashboardSignOffEvaluateErrorKeepsDashboard`, via `applySignOff` with an invalid config).
- Parity with `Evaluate` incl. unknown-severity vulnerability (reported only). Met (same test, per-case).
- Dates in fixed UTC from stored time ("1 Oct 20:53 UTC"). Met (`TestDashboardDatesAreFixedUTC`).
- Headline counts. Met (`TestDashboardHeadline`).

Commands: `go test ./internal/codehealth`, `make build`, `make test-fast` — all passed.

Files changed: `internal/codehealth/dashboard.go`, `internal/codehealth/dashboard_copy.go`, new `internal/codehealth/dashboard_signoff_test.go`. Files read: Context Files plus `chip.go`, `primitives.go`, `history.go`, `classification.go`, `snapshot.go` (extra reads for helper/field names). `gate.go` unchanged.

Limitations: headline counts any non-Good row (including Unknown) as needing a look; "Not enough to judge yet" applies only when every row is Unknown. A not-configured placeholder row carries no SignOff. Evaluate-error path is tested through `applySignOff` directly because stored configs are validated on load. No board rendering yet (T-076/later Tasks).

## Drift Notes

New module or architecture delta beyond the documented Codebase Map, reconciled through the planner before Check.
