---
id: I-133
title: Explain future advice in Advanced Options
type: defect
status: resolved
source:
  kind: check
  check: C-963
  actor: {role: checker, session: check-o037-20261004}
  at: '2026-10-04T02:31:20Z'
tasks: [T-109, T-110]
checks: [C-963, C-964]
guardrail_ids: [TPL-02]
resolution:
  disposition: verified
  check: C-964
  actor: {role: checker, session: recheck-o037-20261004}
  at: '2026-10-04T02:43:55Z'
  reason: 'Rendered screen states only choice saved now, no suggestions yet, advice later.'
history:
  - at: '2026-10-04T02:31:20Z'
    actor: {role: checker, session: check-o037-20261004}
    kind: observed
    check: C-963
    note: Initial independent Full Objective Check.
  - at: '2026-10-04T02:43:55Z'
    actor: {role: checker, session: recheck-o037-20261004}
    kind: rechecked
    check: C-964
    note: 'CLEAR. Rendered screen states only choice saved now, no suggestions yet, advice later.'
---

# I-133: Explain future advice in Advanced Options

## Summary

Explain future advice in Advanced Options before O-037 clearance.

## Evidence

O-037 SC6 explicitly requires public wording to explain that this delivery stores a preference and lane suggestions arrive later.

Open the board and press o. The screen says “Shows optional suggestions for which Tasks could run side by side in separate worktrees.” No screen text explains that no suggestions exist yet. `internal/board/v2/options.go:21` is the visible copy; renderOptions uses it at line 218. README/CHANGELOG correctly explain later delivery, but that does not correct the screen.

Independent TestO037IndependentSequencingCopy prints the supported screen and fails. Existing explanation test checks only ignorable/nonblocking wording.

## Proof Needed

Explain directly on the screen that only the preference is saved now and optional advice arrives later. Preserve the ignorable/no-Code-Health/no-lifecycle explanation. Verify actual rendered copy, README, CHANGELOG and Design agree in frozen M4/M5 cells.
