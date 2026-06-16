---
id: E31-epic-audited-shortcut/T001-mark-audited-keybinding
status: planned
objective: Add an Epic-view keyboard shortcut that writes status audited via UpdateEpicStatus and refreshes the glyph
depends_on: []
complexity_tier: medium
complexity_reason: "Wires a new keybinding to an mtime-guarded write command plus in-memory refresh and hint text, across update/io/model"
---

# T001: Mark Audited Keybinding

## Problem

`data.UpdateEpicStatus(path, status)` (`internal/data/write.go:31`) can write an
epic's status, but no TUI path calls it. The epic-detail overlay handler
`handleEpicDetailOverlay` (`internal/board/update.go:376`) only handles `1`/`2`
tabs, nav, and `esc`. So a user who has reviewed an epic cannot mark it `audited`
from the board — they must hand-edit frontmatter. We need a shortcut that writes
`audited` through the existing helper and refreshes the in-memory `EpicStatus`
map so the sidebar glyph updates immediately.

## Context Files

- `internal/board/update.go`
- `internal/board/io.go`
- `internal/board/model.go`
- `internal/board/help.go`
- `internal/board/epic_panel.go`
- `internal/board/update_test.go`
- `internal/board/io_test.go`
- `internal/data/write.go`

## Acceptance Criteria

- [ ] Pressing `a` in the epic-detail overlay (Epic view) on a focused epic writes
      `status: audited` to that epic's `E##-Detail.md` via
      `data.UpdateEpicStatus`.
- [ ] The write is mtime-guarded the same way task status writes are; on
      conflict/error a clear `StatusMessage` is shown and no partial state is
      left.
- [ ] After a successful write, the in-memory `EpicStatus` map is updated so the
      epic's glyph reflects `audited` without requiring `ctrl+r`.
- [ ] If the epic is already `audited`, the key is a no-op with an
      "already audited" status message (no spurious write).
- [ ] Help text and the audit-tab hint mention the shortcut (`a:mark audited`).
- [ ] Tests cover: key writes `status: audited` to the file; already-audited
      no-op; the `EpicStatus` map refresh; error/conflict messaging.
- [ ] `make build && make test` passes.

## Implementation Plan

- [ ] In `internal/board/io.go`, add `writeEpicStatusCmd(epicID, path, status,
      expectedMtime)` returning a `tea.Cmd` that calls `data.UpdateEpicStatus`
      and emits a result msg, mirroring `writeDefectStatusCmd` / the task status
      write command (including mtime-conflict handling).
- [ ] Add a result message type (e.g. `epicStatusWrittenMsg`) and handle it in
      `Update`: on success update `m.EpicStatus[epicID]` and set a confirmation
      `StatusMessage`; on error set an error `StatusMessage`.
- [ ] In `handleEpicDetailOverlay` (`internal/board/update.go:376`), add
      `case "a":` — resolve the focused epic's ID, detail path, and mtime; if
      already `audited`, set the no-op message; otherwise return
      `writeEpicStatusCmd(..., data canonical "audited", mtime)`. Reuse the
      `epicDetailEpic()` / path-building logic already used by the `"2"` audit-tab
      branch.
- [ ] Confirm the epic detail path resolution matches the audit-tab branch
      (`releases/{release}/epics/{slug}/{short}-Detail.md`) and that the epic's
      mtime is available on the model (add it to the loaded epic state if not).
- [ ] In `internal/board/help.go`, add a help row for `a`; in
      `internal/board/epic_panel.go`, extend the audit-tab hint to mention it.
- [ ] Add tests in `update_test.go` (keybinding behaviour, already-audited no-op,
      map refresh) and `io_test.go` (file gets `status: audited`).
- [ ] Run `make build && make test`.

## Notes

- Use the canonical `audited` literal/constant introduced by E30 if merged first;
  otherwise the string `"audited"` matches today's `statusGlyph` and audit skill.
- Scope is setting `audited` only — not router advancement and not writing the
  `E##-Audit.md` file (the audit skill owns those). Default gating: allow from any
  non-`audited` state per the epic's Open decisions.
