---
id: T-043
title: Reorder Objectives from the keyboard
objective: O-020
status: done
complexity_tier: medium
complexity_reason: Adds two sidebar key families that compute a new group order and call the data writer through a command, with cursor stability and conflict messages.
depends_on: [{task: T-042, requires: clear}]
owner_validation:
    required: true
    accepted_check: ""
planned_by: {role: planner, session: o020-rank-objectives-20260925}
check_waiver:
    task: T-043
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-25T06:49:19Z"
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

Acceptance evidence:
- `1`–`4` route to Critical, High, Medium, and Low. The tests verify appending after existing destination rows, persisted order after opening a fresh board, and no file change when setting the current priority: `TestSidebarPriorityKeysAppendAndKeepSelection`.
- `K`/`J` and the shift-arrow aliases swap only within their current group; group edges do not cross priorities. The High-group test exercises `K`, `K`, then `J`; a first move ranks legacy rows: `TestSidebarPriorityKeysAppendAndKeepSelection`, `TestSidebarGroupMovementRenumbersLegacyRowsAndClamps`, `TestSidebarGroupMovementDoesNotCrossPriorityGroups`.
- Order writes run through Bubble Tea commands. Reloads keep the focused and selected Objective by ID: `TestSidebarPriorityKeysAppendAndKeepSelection`, `TestSidebarGroupMovementRenumbersLegacyRowsAndClamps`.
- The stale-source test checks the Objective-naming conflict, unchanged external file bytes, and unchanged visible order: `TestSidebarReorderConflictKeepsOrderAndExternalEdit`. The injected partial-write test checks the retained visible order and a subsequent healing move: `TestSidebarOrderHealsAfterInjectedMidWriteFailure`.
- Focused help lists the priority and move keys. Tests compare router bytes, the canonical `resume.NextLine`, Objective status, and Task cards across reorder actions: `TestSidebarHelpListsOrderKeysOnlyWhenFocused`, `TestSidebarPriorityKeysAppendAndKeepSelection`, `TestSidebarGroupMovementRenumbersLegacyRowsAndClamps`.

Verification:
- `gofmt -w internal/board/v2/update.go internal/board/v2/io.go internal/board/v2/help.go internal/board/v2/actions_test.go internal/board/v2/objectives_test.go` — completed.
- `git diff --check` — passed.
- `go test ./internal/board/v2` — passed after correcting a missing test import; the initial sandbox attempt could not access the external Go build cache, so the passing focused run used approved elevated access.
- `make build` — passed.
- `make test-fast` — passed, including the latest board tests.
- `./savepoint resume` — `Check T-043 — Reorder Objectives from the keyboard (O-020)`; owner decision is required for an optional Task Check or explicit waiver.

Context read: `.savepoint/router.md`, this Task, O-020, all files in the Task's Context Files list, and the applicable guardrail rules.
Extra reads logged before access: `.savepoint/Guardrails.md` (named policy rules); `.savepoint/Design.md` sections 8 and 9 (cited keybindings, persistence, and concurrency); `internal/resume/resume.go` lines 35–70 (canonical Next-line formatter for an exact display assertion).
T-043 implementation changes: `internal/board/v2/update.go`, `internal/board/v2/io.go`, `internal/board/v2/help.go`, `internal/board/v2/actions_test.go`, and `internal/board/v2/objectives_test.go`, plus this evidence record. The router was not edited; tests confirm its bytes are unchanged by reorder actions.

Limitations: the scratch-project keyboard scenario was verified through temporary-project tests rather than a manual interactive terminal session. No Task Check or waiver is recorded; this Task is ready for the owner-selected handoff.

## Drift Notes

Design section 8 keybindings gain the new keys at Objective reconciliation.
