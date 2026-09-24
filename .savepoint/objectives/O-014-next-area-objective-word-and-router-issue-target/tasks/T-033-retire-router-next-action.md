---
id: T-033
title: Stop the router carrying its own next-step text
objective: O-014
planned_by: {role: planner, session: o014-design-20260924}
status: planned
complexity_tier: medium
complexity_reason: "Removes a decoded field while keeping old routers loadable, adds a doctor retired-field report, and must confirm every hand-off the prose carried is expressed by Next."
depends_on: [{task: T-032, requires: clear}]
owner_validation: {required: false}
---

# T-033: Stop the router carrying its own next-step text

## Outcome

The router holds only `state` and its Goal/Objective/Task/Issue selection.
An old router with `next_action` still loads; doctor reports the field once
as retired, and no surface shows it.

## User Check

Run `savepoint doctor` on this repository before the router is cleaned: it
reports `next_action` as retired. After the cleanup the report is gone and
the board and resume are unchanged.

## Done When

- `RouterStateV2.NextAction` is removed; `ReadStateV2` still accepts a
  `next_action` key (no KnownFields failure) and records only that it was
  present, for doctor.
- Doctor reports a retired `next_action` once, with a repair hint to delete
  the line; it never prints the value.
- `WriteRouterStateV2` leaves an existing `next_action` line byte-identical
  (removal is the owner's or skill's edit) and never adds one.
- `templates/project-v2/.savepoint/router.md` and this repository's
  `.savepoint/router.md` no longer carry `next_action`.
- Evidence lists each hand-off `next_action` expressed in this repo's recent
  router history (for example "request a Task Check or record a waiver",
  "select the next Objective") and the `data.Next` rung or resume phrase that
  now states it; any gap gets a resume phrase, not free text.
- Tests: decode with and without the key, doctor report, writer
  preservation, template freshness.
- `git diff --check`, `make build && make test-fast` pass.

## Context Files

`internal/data/router_v2.go`, `internal/data/router_v2_test.go`,
`internal/data/write.go`, `internal/data/write_test.go`,
`internal/doctor/v2_runtime.go`, `internal/doctor/v2_runtime_test.go`,
`internal/resume/resume.go`, `templates/project-v2/.savepoint/router.md`,
`.savepoint/router.md`, `internal/init/template_freshness_test.go`.

## Design References

Design sections 1, 6, and 11.

## Guardrails

DATA-01, DATA-03, DATA-04, TPL-02, TPL-04, TEST-02, TEST-03, TEST-08.

## Implementation Plan

1. Replace the field with a retired-key presence flag.
2. Add the doctor report and its test.
3. Update the template and live router; run template tests.
4. Record the hand-off coverage table; add phrases only for gaps.

## Boundaries

No guidance rewrite (T-034). The legacy V1 router reader used by migrate is
untouched.

## Technical Verification

Focused data/doctor/init tests during iteration; `make build && make
test-fast` at handoff.

## Technical Evidence

Pending execution.

## Drift Notes

None expected.
