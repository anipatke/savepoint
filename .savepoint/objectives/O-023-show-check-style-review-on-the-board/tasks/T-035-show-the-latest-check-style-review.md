---
id: T-035
title: Show the latest Check's style checklist in record details
objective: O-023
status: planned
depends_on: []
owner_validation: {required: false}
planned_by: {role: planner, session: o023-planning-20260924}
---

# T-035: Show the latest Check's style checklist in record details

## Outcome

Opening a Task, Objective, or Goal detail on the V2 board shows a `CODE STYLE`
section with the latest Check's `## Code Style Review` lines verbatim, or a
plain statement that the latest Check has none.

## User Check

Open the O-021 detail on the board (`v`). Its latest Check, C-915, predates the
checklist, so `CODE STYLE (C-915)` reads that the Check has no Code Style
Review. A fixture Check with the section shows its ticked and unticked lines
as written.

## Done When

- `RecordDetail` carries the latest Check's style review: the Check ID plus the
  section's lines, or a marker that the section is absent. It is resolved in
  `detail.go` from `index.LatestCheck` and `Check.Source.Body`.
- The section runs from the `## Code Style Review` heading to the next `## `
  heading or end of body; lines keep their `- [x]` / `- [ ]` markers and
  reasons; blank edge lines are trimmed. CRLF bodies give the same result.
- `detailLines` renders `CODE STYLE (C-###)` after `CHECKS`. A Check with no
  section renders `(C-### has no Code Style Review)`. With no Check, the
  section is omitted.
- Only the latest Check is read; a superseded Check's section never appears.
- Clearance, badges, Next, and the plain table are unchanged.
- Design section 8 describes the new section.
- `make build && make test-fast` pass.

## Context Files

`internal/board/v2/detail.go`, `internal/board/v2/detail_view.go`,
`internal/board/v2/checks.go`, `internal/board/v2/detail_test.go`,
`internal/data/check_v2.go`, `internal/data/parser.go`,
`agent-skills/references/check-method.md`, `.savepoint/Design.md`.

## Design References

Design section 7 (Check workflow, Layer 2 advisory style) and section 8 (TUI).

## Guardrails

ARCH-02, DATA-03, TPL-02, TEST-01, TEST-02, TEST-04, TEST-08, DEP-02.

## Implementation Plan

1. Add a pure helper in `internal/board/v2` that returns a body's
   `## Code Style Review` lines and whether the heading was found.
2. Resolve the latest Check's review into `RecordDetail` in the Task,
   Objective, and Goal constructors.
3. Render the `CODE STYLE (C-###)` section after `CHECKS` in `detailLines`.
4. Test: section present (ticked, unticked with reason), section absent, no
   Check, superseded Check with a section while the latest has none, section
   followed by another `##` heading, and a CRLF body.
5. Reconcile Design section 8; run the handoff gate.

## Boundaries

Presentation only. No change to `internal/data` types or decoders, no whole
Check bodies, no Check selection, and no edits to existing Check records.

## Technical Verification

Focused `make test-focused TEST=Detail` while iterating; `make build && make
test-fast` at handoff. Evaluated later by the mandatory O-023 Full Objective
Check per `agent-skills/references/check-method.md`.

## Technical Evidence

Pending execution.

## Drift Notes

None planned.
