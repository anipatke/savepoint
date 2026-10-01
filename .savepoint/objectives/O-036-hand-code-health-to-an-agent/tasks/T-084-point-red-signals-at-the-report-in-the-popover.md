---
id: T-084
title: Point red signals at the report in the popover
objective: O-036
status: planned
depends_on: [{task: T-082, requires: clear}]
owner_validation: {required: true}
planned_by: {role: planner, session: planning-o036-20261002}
---

# Point red signals at the report in the popover

## Outcome

When a signal needs attention, the popover's next step tells the builder to ask their agent to investigate `.savepoint/health/report.md`, and the `Where:` line stops claiming a command lists the files.

## User Check

Open `H` on a project where complexity is red and the report exists: the selected signal's `Next:` reads "Ask your agent to investigate .savepoint/health/report.md". Delete the file and reopen: it says to run `savepoint health report` first. A Watch or Good signal keeps its own next step. A signal with several files says "13 files" and nothing about a command.

## Done When

- A signal labelled Needs Attention has the next step "Ask your agent to investigate .savepoint/health/report.md" when the report exists, and "Run savepoint health report, then ask your agent to investigate it." when it does not. Watch, Unknown and Good keep the existing next steps.
- The dashboard load notes whether the report file exists, once, with a single file-existence check; rendering does no file work.
- The `Where:` text for several files is the count only ("13 files"); the single-file text is unchanged.
- The next-step line still fits the 80x20, 80x24 and 80x40 popover without moving other rows.
- Tests cover a red signal with and without a report, Watch/Unknown/Good unchanged, one and many files in `Where:`, the dashboard load noting a present and an absent report, and the popover at the three sizes.

## Context Files

`.savepoint/objectives/O-036-hand-code-health-to-an-agent/Objective.md`; `internal/codehealth/dashboard.go`, `internal/codehealth/dashboard_copy.go`, `internal/codehealth/dashboard_copy_test.go`, `internal/codehealth/dashboard_test.go` (helpers), `internal/codehealth/report.go` (read only: path); `internal/board/v2/health_view.go` (read only), `internal/board/v2/health_realcopy_test.go`.

## Design References

O-036 Confirmed Design Decisions; O-035 wording decisions.

## Guardrails

ARCH-02, STYLE-07, STYLE-09, TEST-01, TEST-02, TEST-04.

## Implementation Plan

1. Confirm the report path from the report Task exists; return REPLAN REQUIRED otherwise.
2. Have the dashboard loader record whether the report exists and choose the red next step from it; keep the wording in the copy tables.
3. Change the many-files `Where:` text.
4. Add the tests in Done When.

## Boundaries

No report content, command or snapshot hook; no new popover keys or lines.

## Technical Verification

Focused tests during iteration; `make build && make test-fast` for ordinary handoff. Later evaluated per `agent-skills/references/check-method.md`.

## Technical Evidence

Pending execution: named cases, results, reviewed source basis, files read/changed, and limitations, recorded after the work lands.

## Drift Notes

New module or architecture delta beyond the documented Codebase Map, reconciled through the planner before Check.
