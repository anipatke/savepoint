---
id: I-063
title: The board cannot close or reopen a Goal
type: defect
status: open
source:
  kind: report
  actor: {role: owner, session: user}
  at: '2026-09-25T21:59:32Z'
severity: medium
history:
  - at: '2026-09-25T21:59:32Z'
    actor: {role: owner, session: user}
    kind: observed
    note: >-
      With every R-006 Objective complete, Next reads "Close R-006" but the
      board offers no way to close a Goal. Owner asked for a C key in the
      Goal selector that toggles close and reopen, with a status marker on
      each Goal row.
  - at: '2026-09-25T22:01:48Z'
    actor: {role: executor, session: i063-repair-20260926}
    kind: repair_attempted
    note: >-
      data.WriteReleaseV2 patches only a Goal's status through writeV2Record.
      The board's writeGoalStatusCmd re-reads the index, refuses Goals with a
      legacy completion, reopens a done Goal to in_progress, and otherwise
      closes only when data.ResolveReleaseCompletion allows it, naming the
      blockers when it does not; it never writes the router. C in the Goal
      selector runs it. Rows show releaseStatusMarker (green [✓] done, grey
      [ ] otherwise, in badges.go). The selector hint line now reads
      "enter:select  C:close/reopen  v:detail  esc:cancel" and fits the
      52-cell overlay; help and the Design.md keybindings paragraph describe
      C. goal_close_test.go covers close with byte-exact status-only write
      and unchanged router, reopen, blocked close naming O-001, archived
      refusal, markers, hints, and help. make build && make test-fast
      passed. Not exercised by eye in a live terminal. Issue remains open for
      independent verification.
---

# I-063: The board cannot close or reopen a Goal

## Summary

When every member Objective of a Goal is complete, Next reads
`Close R-###` and resume says "Record Goal R-### as done", but the board has
no action that records it. O-025 T-047 removed the Goal Check and explicitly
left "No board action to close a Goal", so the only path is hand-editing
`status:` in the Goal's `Release.md`. The Goal selector also gives no sign
of which Goals are done.

## Evidence

- `internal/resume/resume.go:88` returns `Close` for `NextReleaseReady`;
  `internal/resume/resume.go:487` says "Record Goal %s as done."
- `internal/board/v2/update.go` `handleReleaseKey` handles only navigation,
  `enter`, `v`/`d`, and cancel keys.
- `internal/board/v2/releases.go` `renderReleaseSelector` draws rows as
  `► R-### — Title` with no status marker.
- `.savepoint/objectives/O-025-remove-the-goal-check/tasks/T-047-*.md`
  Boundaries: "No board action to close a Goal."

## Proof Needed

- Each Goal selector row shows a green `[✓]` when the Goal is done and a grey
  `[ ]` otherwise, legible without colour.
- `C` on the focused Goal closes it (`status: done`) only when
  `data.ResolveReleaseCompletion` allows it, re-resolved from a fresh index
  immediately before the write; a refusal names the blocking Objectives and
  writes nothing.
- `C` on a done Goal reopens it to `status: in_progress`.
- `C` on a Goal carrying a migrated legacy completion is refused and writes
  nothing.
- The write patches only `status` in `Release.md` through the V2 record
  writer with its conflict check, never changes router `release:`, and
  reloads the board.
- The selector hint line, help, and Design.md keybindings describe `C`.
- Tests cover close, reopen, blocked close, legacy refusal, and the markers;
  `make build && make test-fast` pass.
