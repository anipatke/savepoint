---
id: I-105
title: Code Health refresh Help swallows Ctrl+C cancellation
type: defect
status: open
source:
  kind: check
  check: C-948
  actor: {role: checker, session: check-o031-20261002-independent}
  at: '2026-10-01T20:47:19Z'
tasks: [T-072]
checks: [C-948]
guardrail_ids: [TEST-01, TEST-02]
history:
  - at: '2026-10-01T20:47:19Z'
    actor: {role: checker, session: check-o031-20261002-independent}
    kind: observed
    check: C-948
    note: Reproduced inside the frozen initial Full Objective Check scope of O-031.
  - at: '2026-10-02T07:30:00Z'
    actor: {role: executor, session: user-request}
    kind: repair_attempted
    note: 'Ctrl+C now bypasses the Help handler in handleKey so it cancels the active refresh and quits. Added TestHealthCtrlCFromHelpDuringRefreshCancelsBeforeQuitting.'
---

# I-105: Code Health refresh Help swallows Ctrl+C cancellation

## Summary

Code Health refresh Help swallows Ctrl+C cancellation. Repair directly under this Issue; keep the completed Task status unchanged.

## Evidence

T-072 Done When 6 requires Ctrl+C during refresh to cancel before quitting.

Start a refresh with R, open Help with ?, then press Ctrl+C. The update returns no command and never invokes the collection cancel function. The running tool continues. `internal/board/v2/update.go:88-92` returns from the Help handler before the cancellation/quit branch at 101-106.

Independent `TestO031IndependentQuitFromRefreshHelp/ctrl+c` (embedded in C-948) supplies an active refresh with an observed cancel function, opens Help through Update, and sends Ctrl+C. Actual cancelled=false, quit=false; expected cancel followed by quit. The direct-refresh test `TestHealthQuitDuringRefreshCancelsBeforeQuitting` never opens Help. The q variant also closes Help without cancelling; that existing close-help convention is recorded separately as an observation, not part of this Issue's required repair.

## Proof Needed

Ctrl+C from refresh Help must cancel the same active context and return the quit command. Re-run direct-refresh quit and Help close/back navigation; show no cancelled snapshot is saved.
