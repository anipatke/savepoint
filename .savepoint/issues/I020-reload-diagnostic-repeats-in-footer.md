---
id: I020
title: Reload diagnostic repeats in the footer
type: defect
status: open
source:
  kind: report
  actor: {role: owner, session: user}
  at: '2026-09-22T09:45:29Z'
severity: medium
history:
  - at: '2026-09-22T09:45:29Z'
    actor: {role: owner, session: user}
    kind: observed
    note: A failed board reload displays the same long structural diagnostic in both the RELOAD banner and the footer status line; the owner requested that it be removed from the footer.
---

# I020: Reload diagnostic repeats in the footer

## Summary

When a watcher reload fails, the board correctly keeps the last valid state and
shows the full failure in its dedicated `RELOAD` diagnostic banner. The same
failure is also copied into `StatusMessage`, so the footer repeats the long
diagnostic. For path-heavy structural errors such as duplicate global record
IDs, this makes both the top and bottom of the screen noisy without adding
information.

The dedicated reload banner should own the detailed reload diagnostic. The
footer should retain its ordinary action/status role and must not repeat the
same error text.

## Evidence

- `internal/board/v2/update.go:543-544` assigns the same load failure to both
  `ReloadDiagnostic` and `StatusMessage`.
- `internal/board/v2/view.go:156-157` renders `ReloadDiagnostic` in the
  dedicated `RELOAD` banner.
- `internal/board/v2/view.go:354-355` independently renders `StatusMessage` in
  the footer, producing the duplicate presentation.
- The owner observed the duplicate during the T006 ID-collision reload on
  2026-09-22 and requested removal from the footer.

## Proof Needed

- A failed watcher reload keeps the last valid board visible and shows the full
  named diagnostic exactly once in the dedicated `RELOAD` banner.
- The footer does not copy the reload diagnostic and continues to show normal
  action/status content according to the existing footer contract.
- Recovery after the invalid file is corrected clears the reload banner and
  preserves existing reload/focus behavior.
- Focused board reload tests, `git diff --check`, `make build`, and `make test`
  pass.
