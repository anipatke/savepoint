---
id: I-036
title: Goal detail closed by reload reports "RELEASE ... no longer exists"
type: defect
status: resolved
source:
  kind: check
  check: C-911
  actor: {role: checker, session: o013-full-check-20260923}
  at: '2026-09-23T09:45:00Z'
tasks: [T-007]
checks: [C-911]
severity: low
resolution:
  disposition: accepted
  actor: {role: owner, session: user}
  at: '2026-09-23T09:47:00Z'
  reason: >-
    Owner directed closure of the repair without a fresh-session recheck. The repair was made by the same session that wrote C-911, so this is owner acceptance, not a technical verified closure.
history:
  - at: '2026-09-23T09:45:00Z'
    actor: {role: checker, session: o013-full-check-20260923}
    kind: observed
    check: C-911
    note: >-
      The reload path that closes a detail overlay whose record vanished
      formats its status with the raw DetailKind, which is "RELEASE" for a
      Goal detail.
  - at: '2026-09-23T09:47:00Z'
    actor: {role: executor, session: o013-full-check-20260923}
    kind: repair_attempted
    note: >-
      At the owner's direction, added detailKindLabel in internal/board/v2/detail_view.go (shared with detailHeading) and used it for the removed-detail reload status in update.go, which now reads "GOAL R-### no longer exists; detail closed." Added TestReloadClosingARemovedGoalDetailNamesItAGoal. make build and make test-full pass (2026-09-23, go1.26.2 linux/amd64).
---

# I-036: Goal detail closed by reload reports "RELEASE ... no longer exists"

## Violated Requirement

T-007 Done When: "Selector, selected-context header, detail overlay, help,
and status/error messages use the Goal vocabulary", and "Reload ... cases
remain understandable". O-013 Success Condition 1 lists status messages.

## Scenario

1. Create a V2 project with Goals R-001, R-002, and an empty R-003 that no
   Objective references.
2. Open the board, press `g`, move to R-003, and press `v` to open its detail.
3. Delete `releases/R-003-empty/` and let the board reload.

Expected: a status such as "Goal R-003 no longer exists; detail closed."
Actual: `RELEASE R-003 no longer exists; detail closed.` This was reproduced
by a temporary probe test in the Check session and then deleted.

Evidence: `internal/board/v2/update.go:587` formats
`snapshot.DetailKind` directly. `DetailRelease` is `"RELEASE"`
(`internal/board/v2/detail.go:26`). T-007 routed the heading through
`detailHeading` (`detail_view.go:61`), but this status path was missed.
`TestReleaseReloadPreservesFocusAndDiagnosesRemovedSelection` covers only the
removed-selection message, not a removed open Goal detail.

## Proof Needed

- The removed-detail status for a Goal detail uses Goal wording.
- A focused test covers the removed open Goal detail on reload.
