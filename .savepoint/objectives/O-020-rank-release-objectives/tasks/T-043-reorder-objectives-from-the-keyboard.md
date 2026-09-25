---
id: T-043
title: Reorder Objectives from the keyboard
objective: O-020
status: planned
complexity_tier: medium
complexity_reason: Adds two sidebar key families that compute a new group order and call the data writer through a command, with cursor stability and conflict messages.
depends_on: [{task: T-042, requires: clear}]
owner_validation: {required: true}
planned_by: {role: planner, session: o020-rank-objectives-20260925}
---

# Reorder Objectives from the keyboard

## Outcome

With the sidebar focused, the owner presses `1`–`4` to set an Objective's
priority and `K`/`J` to move it up or down within its group. The row moves at
once, the cursor stays on it, and the new order survives a restart.

## User Check

In a scratch copy of a project, focus the sidebar and press `2` on a Medium
Objective: it moves to the bottom of High with the cursor on it. Press `K`
twice and `J` once and confirm each step. Press `K` at the top of a group and
confirm nothing changes. Quit and reopen: the order is the same. Edit that
Objective's file in another editor, then press `J` without refreshing, and
confirm the board names the conflict and no file changed. Confirm `?` help
lists the new keys and the Next line never changed.

## Done When

- `1`–`4` on a focused sidebar row append that Objective to the bottom of the
  destination group through the data order writer. Pressing its current
  priority does nothing.
- `K`/`shift+up` and `J`/`shift+down` swap it with its neighbour inside its
  group and renumber the group through the same writer. At the group edge they
  do nothing; they never cross into another group.
- Writes run as Bubble Tea commands, never during rendering. After the reload
  the cursor and selection stay on the same Objective.
- A conflict or write error shows a status message naming the Objective and
  leaves the prior order on screen.
- Help lists the keys when the sidebar is focused. The router, Next, Objective
  status, and Task columns are unchanged by every reorder.
- Tests cover each priority key, no-op on the current priority, move up and
  down, edges, a group with unranked legacy rows being renumbered on first
  move, cursor and selection stability across reload, a stale-file conflict,
  an injected mid-write failure followed by a healing move, unchanged router
  bytes, and help text.

## Context Files

`internal/board/v2/update.go`, `internal/board/v2/actions.go`,
`internal/board/v2/actions_test.go`, `internal/board/v2/io.go`,
`internal/board/v2/objectives.go`, `internal/board/v2/objectives_test.go`,
`internal/board/v2/help.go`, `internal/board/v2/footer_test.go`,
`internal/board/v2/watch_test.go`, `internal/data/write.go`,
`.savepoint/objectives/O-020-rank-release-objectives/Objective.md`.

## Design References

Design sections 8 (keybindings, board persistence) and 9 (concurrency). O-020
Confirmed Design Decisions.

## Guardrails

FS-01, ARCH-02, ARCH-03, DATA-01, TEST-01..04, TEST-08, STYLE-02.

## Implementation Plan

1. Add a pure helper that, given the current rows and a key, returns the
   group's new ordered IDs and target priority, or nothing for a no-op.
2. Add a write command in `io.go` that calls the data order writer and returns
   the existing action message, with a reload.
3. Dispatch the keys in `handleSidebarKey`; rely on the existing cursor
   restore by Objective ID after reload, and fix it if a moved row lands
   elsewhere.
4. Add the keys to help.
5. Write the tests listed above.

## Boundaries

No new data fields, display changes, Goal-move action, or change to Next and
router selection.

## Technical Verification

Focused tests while iterating; `make build && make test-fast` at handoff.

## Technical Evidence

Pending execution.

## Drift Notes

Design section 8 keybindings gain the new keys at Objective reconciliation.
