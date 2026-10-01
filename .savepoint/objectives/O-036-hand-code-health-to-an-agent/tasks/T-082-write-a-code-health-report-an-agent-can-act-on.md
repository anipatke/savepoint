---
id: T-082
title: Write a Code Health report an agent can act on
objective: O-036
status: planned
owner_validation: {required: true}
planned_by: {role: planner, session: planning-o036-20261002}
---

# Write a Code Health report an agent can act on

## Outcome

Code Health can describe its newest saved snapshot as one plain-text report: a short brief for an agent, then the five signals with the ones needing attention first, each with its number, aim, meaning, next step, sign-off and every affected file.

## User Check

With a saved official snapshot that has a red complexity signal, render the report (the Task's tests print it) and read it as an agent would: the brief says what to do, complexity comes first, and each affected function is listed as `path:line` with its note. No provider names, hashes or raw timestamps appear.

## Done When

- A function in `codehealth` renders the report from a loaded `Dashboard` and returns the text; it reads no file and runs nothing.
- The text opens with a fixed brief: investigate each signal that is not Good, blocking ones first, propose then apply fixes, then re-run the check. It names the re-run command for the Objective given to it.
- All five signals appear ordered Needs Attention, Watch, Unknown, Good, keeping the dashboard's order within a label. Each shows its name, question, label, number, aim, meaning, next step, sign-off wording and better/worse/steady word, exactly as the dashboard row words them.
- Every evidence entry of a signal is listed as `path:line` with its note in stored order; an entry with no line shows the path only. A signal with several instances gets one section per instance. A signal with no evidence says so.
- The measured date and whether the snapshot is official or manual are stated. Provider names, snapshot hashes and raw timestamps are absent.
- With no snapshot or no Code Health set up, the renderer returns a clear refusal instead of a report.
- `Store` gains one method that writes the report atomically to `.savepoint/health/report.md` with a `.gitignore` beside it that ignores `report.md`, creating both when needed; a second identical write changes nothing, and an existing `.gitignore` the owner edited there is left alone.
- Tests cover each label and signal, ordering, one, many and no evidence, several instances, a manual newest snapshot, no snapshot, not set up, banned strings, an unwritable directory, atomic replacement of an older report, and the ignore file being created once and never overwritten.

## Context Files

`.savepoint/objectives/O-036-hand-code-health-to-an-agent/Objective.md`; `internal/codehealth/dashboard.go`, `internal/codehealth/dashboard_copy.go`, `internal/codehealth/snapshot.go` (read only: evidence fields), `internal/codehealth/storage.go`, `internal/codehealth/storage_test.go` (helpers), `internal/codehealth/dashboard_test.go` (helpers); new `internal/codehealth/report.go` and `internal/codehealth/report_test.go`.

## Design References

O-036 Confirmed Design Decisions; O-035 plain-language wording decisions.

## Guardrails

FS-04, FS-05, FS-06, ARCH-02, ARCH-04, STYLE-01, STYLE-07, STYLE-09, TEST-01, TEST-02, TEST-03, TEST-04.

## Implementation Plan

1. Confirm the dashboard fields and `Evidence` exist as described; return REPLAN REQUIRED otherwise.
2. Add the renderer in `report.go` using the dashboard row fields and the existing wording; no copy is duplicated.
3. Add the store write method with the same temp-and-rename approach as snapshots, and the one-time `.gitignore`.
4. Add the tests in Done When.

## Boundaries

No command, no hook into snapshot saving, and no popover change; those are the next Tasks. No new data is collected and nothing in classification, thresholds or snapshots changes.

## Technical Verification

Focused tests during iteration (`make test-focused`); `make build && make test-fast` for ordinary handoff. Later evaluated per `agent-skills/references/check-method.md`.

## Technical Evidence

Pending execution: named cases, results, reviewed source basis, files read/changed, and limitations, recorded after the work lands.

## Drift Notes

New module or architecture delta beyond the documented Codebase Map, reconciled through the planner before Check.
