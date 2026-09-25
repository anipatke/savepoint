---
id: I-054
title: Design TUI section does not describe Objective priority groups or reorder keys
type: drift
status: resolved
source:
  kind: check
  check: C-926
  actor: {role: checker, session: o020-objective-check-20260925}
  at: '2026-09-25T06:55:00Z'
tasks: [T-042, T-043]
checks: [C-926, C-927]
severity: low
resolution:
  disposition: verified
  check: C-927
  actor: {role: checker, session: o020-objective-recheck-20260925}
  at: '2026-09-25T07:29:07Z'
  reason: Design §8 now documents the implemented Objective groups and plain list, order keys, and persistence/recovery behavior; the current sources agree and the Full Objective gate passed.
history:
  - at: '2026-09-25T06:55:00Z'
    actor: {role: checker, session: o020-objective-check-20260925}
    kind: observed
    check: C-926
    note: Design §8 was not reconciled with O-020's grouped sidebar, plain Objective list, order keys, or multi-record order write.
  - at: '2026-09-25T07:29:07Z'
    actor: {role: checker, session: o020-objective-recheck-20260925}
    kind: rechecked
    check: C-927
    note: Design §8 correction matched the implementation and tests; the Full Objective recheck passed.
---

# I-054: Design TUI section does not describe Objective priority groups or reorder keys

## Summary

O-020 changed the board, but `.savepoint/Design.md` §8 (TUI) still describes
the pre-O-020 board. Only §1 gained the new "Objective planning order" line.
T-042's and T-043's Drift Notes both said §8 would be updated "at Objective
reconciliation". That reconciliation did not happen before the Full Objective
Check, which covers reconciliation against Design.

## Evidence

- `.savepoint/Design.md` §8 **Layout** (line ~179) does not mention the
  `CRITICAL`/`HIGH`/`MEDIUM`/`LOW` sidebar headings, the removed
  Planned / In Progress / Done row line, or the grouped Objective list that
  plain output now prints after `Selected:` (`internal/board/v2/plain.go`
  lines 56–73).
- §8 **Keybindings** (line ~206) does not list `1`–`4` (set priority) or
  `K`/`J`/`shift+↑`/`shift+↓` (move within group), which
  `internal/board/v2/update.go` `sidebarObjectiveOrderChange` dispatches and
  `internal/board/v2/help.go` lists.
- §8 **Board persistence and refresh** (line ~204) describes only Task status
  writes. It does not describe the Objective order write
  (`data.WriteObjectiveGroupOrderV2`): freshness is checked for every record
  before any write, each file is replaced atomically, and an interrupted
  write heals on the next move.
- T-042 Drift Notes: "Design section 8 gains the grouped sidebar and plain
  Objective list at Objective reconciliation." T-043 Drift Notes: "Design
  section 8 keybindings gain the new keys at Objective reconciliation."

## Proof Needed

Planner updates Design §8 Layout, Keybindings, and Board persistence to
describe the implemented behavior. The update changes no code. A recheck
confirms §8 matches `objectives.go`, `plain.go`, `update.go`, `help.go`, and
`write.go`.
