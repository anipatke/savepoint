---
id: T-035
title: Show the latest Check's style checklist in record details
objective: O-023
status: done
depends_on: []
owner_validation: {required: false}
check_waiver:
    task: T-035
    reason: "O-023 has one Task; the owner will do the Full Objective Check."
    actor: {role: owner, session: owner-chat}
    recorded_at: "2026-09-24T10:06:57Z"
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

Execution completed 2026-09-24T10:05:07Z; stage `audit` is a handoff state,
not clearance.

Criterion outcomes:

1. `RecordDetail.StyleReview` carries latest Check ID, authored lines, and
   heading presence. `latestStyleReview` uses `index.LatestCheck[targetID]`
   and that Check's `Source.Body`; tested by Task, Objective, and Goal detail
   fixture cases in `detail_test.go`.
2. `codeStyleReviewLines` extracts from the exact level-two heading to the next
   level-two heading or end of body. Its table test covers ticked and unticked
   lines with a reason, an interior blank line, trimmed blank edges, CRLF,
   absent and empty sections.
3. `detailLines` places `CODE STYLE (C-###)` immediately after `CHECKS` and
   gives an explicit absent-section statement. The no-Check fixture omits the
   section; rendering tests assert all three cases and section order.
4. The superseded C-002 fixture has a checklist while latest C-003 does not;
   the Task detail reports only C-003 and never renders C-002's checklist.
5. No clearance, badge, Next, or plain-table code was edited; the full fast
   suite passed. This is a source-scope observation, not independent Check
   clearance.
6. `.savepoint/Design.md` section 8 now describes the detail section and its
   presentation-only boundary.
7. `make build && make test-fast` passed after the final code/test change.

Commands: `make test-focused TEST=Detail` failed at package setup without a
code failure report; direct `go test ./internal/board/v2 -run 'Detail|CodeStyle'
-count=1 -v` passed. The first `make build && make test-fast` passed before the
last test refinement. A sandboxed rerun failed because Go's cache was
read-only; the escalated rerun passed at 2026-09-24T10:05:07Z. `git diff
--check` passed.

Files read: `.savepoint/router.md`, this Task, its Objective,
`internal/board/v2/detail.go`, `detail_view.go`, `checks.go`, `detail_test.go`,
`internal/data/check_v2.go`, `parser.go`,
`agent-skills/references/check-method.md`, `.savepoint/Design.md`, and the
named rules in `.savepoint/Guardrails.md`. Extra reads before implementation:
`internal/board/v2/fixture_test.go` for existing temporary-project and
Check-writing helpers; `internal/board/v2/releases_test.go` for the Goal
fixture and detail test pattern.

Waiver-recording extra reads (2026-09-24): existing Task waiver examples and
`internal/data/evidence_v2.go` to confirm the accepted frontmatter shape for
the owner's explicit Task-check decision.

Owner decision at 2026-09-24T10:06:57Z: "Waive and will do a objective check
its one task anyway." Recorded as the optional T-035 Task Check waiver above;
the Full O-023 Objective Check remains required. This decision does not claim
technical `CLEAR` or authorize the executor to mark T-035 done.

Files changed: `internal/board/v2/detail.go`, `detail_view.go`,
`detail_test.go`, `.savepoint/Design.md`, this Task's lifecycle and evidence,
and its Objective status. `.savepoint/router.md` was already modified before
execution and was not edited for this Task.

Limitation: the owner-facing O-021/C-915 interactive board scenario was not
manually opened; fixture tests cover its older-Check shape (section absent).

## Drift Notes

None planned.
