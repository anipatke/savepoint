---
id: v1.2/D016-in-progress-defect-not-resolvable-via-space
release: v1.2
status: open
severity: medium
title: "Board defect overlay cannot resolve an in_progress defect; space is a dead-end no-op"
---

# D016: In-progress defect not resolvable via space in the board overlay

## Symptom

In the board defect overlay (`d`), pressing `space` on a defect whose
`status: in_progress` does nothing except show the message:

> Defect in progress: resolve after lifecycle stage is closed

There is no way in the TUI to "close the lifecycle stage," so the defect is
stuck at `in_progress` and the user cannot mark it `resolved` from the board.

## Expected Behavior

`space` is the user's explicit "mark resolved" action. The defect-build handoff
intentionally parks a finished defect at `status: in_progress` for user sign-off
(AGENTS.md: "Only the user may set… `resolved`"). So `space` on an `in_progress`
defect should resolve it directly — set `status: resolved`, clear `stage` — the
same write path already used for an `open` defect, preserving the missing-path
guard and the `—`/already-resolved messages.

## Reproduction

1. Have a defect file with `status: in_progress` (e.g. left there by the
   defect-build handoff).
2. Open the board, press `d` for the defect overlay, select that defect.
3. Press `space`.
4. Observe the no-op message; the defect stays `in_progress`. No other key in
   the defect overlay or defect-detail overlay resolves it, and no CLI
   subcommand resolves it (`doctor` only heals the stage, never closes it).

## Impact

- A finished defect cannot be closed from the TUI — the documented
  `open → in_progress → resolved` lifecycle has no working `in_progress →
  resolved` transition for the user.
- The only workaround is hand-editing frontmatter (`status: resolved`, drop
  `stage:`), which the board is supposed to make unnecessary.
- The status message promises a "close the lifecycle stage" affordance that does
  not exist anywhere in the tool, so the guidance is misleading.

## Context Files

- `internal/board/update.go` — `handleDefectOverlay`, `case " "` (the
  `data.DefectInProgress` branch is the dead-end no-op).
- `internal/board/defect_overlay_test.go` —
  `TestUpdate_defectOverlaySpaceInProgressIsNoop` currently locks in the no-op
  and must be rewritten to assert resolution.

File reality note: verify each path with `rg --files` before repair; if a new
file is needed, record why in `## Resolution Notes` first.

## Fix Plan

Chosen behavior: **space resolves an `in_progress` defect directly** (smallest
change, matches "only the user marks resolved").

1. In `handleDefectOverlay` (`internal/board/update.go`), fold
   `data.DefectInProgress` into the resolving branch alongside `""` /
   `data.DefectOpen`, keeping the `defect.Path == ""` guard:

   ```go
   case "", data.DefectOpen, data.DefectInProgress:
       if defect.Path == "" {
           m.StatusMessage = "Defect not updated: missing file path"
           return m, nil
       }
       next := defect
       next.Status = data.DefectResolved
       next.Stage = ""
       return m, writeDefectStatusCmd(next, defect.Mtime)
   ```

   Drop the dedicated `case data.DefectInProgress:` no-op message; keep the
   `data.DefectResolved` "already resolved" and `default` invalid-status cases.
2. Rewrite `TestUpdate_defectOverlaySpaceInProgressIsNoop` to assert that space
   on an `in_progress` defect writes `status: resolved` and clears `stage`
   (mirror `TestUpdate_defectOverlaySpaceResolvesPlannedDefect`). Rename it to
   reflect the new behavior.
3. Sanity-check the overlay/help hints still read correctly (`space:resolve`,
   `space: advance focused task / resolve selected defect`) — no copy change
   expected.

## Acceptance Criteria

- [ ] `space` on an `in_progress` defect writes `status: resolved` and removes
      the `stage` field.
- [ ] The missing-path guard, "already resolved" message, and invalid-status
      default are preserved.
- [ ] A test asserts `in_progress` + `space` → `resolved` (replacing the old
      no-op test).
- [ ] No regression to resolving an `open` defect via space.

## Resolution Notes

Pending.
