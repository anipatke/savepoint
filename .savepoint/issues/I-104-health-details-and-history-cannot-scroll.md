---
id: I-104
title: Code Health details and history cannot scroll to all content
type: defect
status: resolved
source:
  kind: check
  check: C-948
  actor: {role: checker, session: check-o031-20261002-independent}
  at: '2026-10-01T20:47:19Z'
tasks: [T-072]
checks: [C-948, C-951]
guardrail_ids: [TEST-01, TEST-02]
resolution:
  disposition: verified
  check: C-951
  actor: {role: checker, session: recheck-o031-final-repair-20261002-independent}
  at: '2026-10-01T21:56:44Z'
  reason: Original frozen-scope repair proof passed; C-951 records current independent CLEAR.
history:
  - at: '2026-10-01T20:47:19Z'
    actor: {role: checker, session: check-o031-20261002-independent}
    kind: observed
    check: C-948
    note: Reproduced inside the frozen initial Full Objective Check scope of O-031.
  - at: '2026-10-02T07:30:00Z'
    actor: {role: executor, session: user-request}
    kind: repair_attempted
    note: 'Health overview and detail now scroll freely: detail has no cursor pin (noCursorLine) and Down past the last signal scrolls into history (Model.moveHealth, healthWindow). Added TestHealthScrollReachesAllDetailAndHistoryContent at 80x20, 80x24, 80x40.'
  - at: '2026-10-01T21:56:44Z'
    actor: {role: checker, session: recheck-o031-final-repair-20261002-independent}
    kind: rechecked
    check: C-951
    note: Frozen viewport and text matrix rechecked at all sizes; every detail field, evidence reference and history entry is reachable.
---

# I-104: Code Health details and history cannot scroll to all content

## Summary

Code Health details and history cannot scroll to all content. Repair directly under this Issue; keep the completed Task status unchanged.

## Evidence

T-072 Done When 5 and 9 and O-031 Success Conditions 5 promise complete signal details and reachable history at supported terminal sizes.

At 80x20 or 80x24, open H, Enter on a signal, then press Down repeatedly. Lower detail fields and affected areas never become visible. Even at 80x40 a long evidence list cannot be scrolled. Return to the overview and repeatedly press Down: navigation stops at the last signal and the history below remains unreachable in short terminals.

`internal/board/v2/health.go:163-169` increments detail Offset, but `health_view.go:133-134` supplies cursorLine=0 for detail. `health_view.go:226-234` resets any positive offset to that cursor line, both in `syncHealthScroll` (242-249) and render (252-257). Overview navigation clamps Cursor to rows only and uses cursor visibility instead of providing a route into history.

Independent `TestO031IndependentViewportMatrix` (embedded in C-948) fails at heights 20, 24, and 40: detail offset remains 0; at heights 20/24 Manual refresh history is absent after 60 Down keys. Existing `TestHealthFitsTheBoardsMinimumWidth` checks dimensions only; `TestHealthDetailsExplainTheLabelAndEscReturnsToTheList` uses a large terminal with short data; `TestHealthKeepsTheCursorRowOnScreenInAShortTerminal` checks the final signal, not content beneath it.

## Proof Needed

At 80x20, 80x24, and 80x40, repeated scrolling must reach every rendered detail field and the final evidence reference, and all bounded history entries must be reachable. Retain up/down bounds, row selection, Esc restoration, and resize behavior. Cover overflow by wrapping and by multiple evidence references.
