---
type: epic-design
status: audited
---

# E33: Release Docs View

## Purpose

The board can show release, epic, task, audit, and defect context, but it does
not provide an in-board way to review the supporting planning documents that
explain why the work exists. Users must leave the board to read the project PRD
or system design. This epic adds a read-only Release Docs section inside the
Epic detail experience so agents and humans can inspect PRD and Design context
without changing workflow state.

## What this epic adds

- A Release Docs subview reachable from the Epic detail overlay.
- A compact document selector for PRD and Design.
- A scrollable read-only document body renderer with width-aware wrapping.
- Missing-document empty states that do not crash the board.
- Test coverage for data loading, navigation, rendering, and wrapping behavior.

## Components and files

| Module | Purpose |
|--------|---------|
| `internal/data/release_doc.go` | Load known supporting document files from the `.savepoint` root |
| `internal/data/release_doc_test.go` | Data loader coverage for present, missing, and unreadable docs |
| `internal/board/model.go` | Store Epic detail subview, selected release doc, and doc scroll state |
| `internal/board/io.go` | Add async board messages/commands for loading release docs |
| `internal/board/update.go` | Handle subview switching, doc selection, and doc scrolling keys |
| `internal/board/epic_panel.go` | Render Epic detail tabs and Release Docs content |
| `internal/board/epic_panel_test.go` | Rendering and navigation coverage for the new subview |
| `internal/board/help.go` | Add concise help text for Release Docs navigation if needed |

## Architectural delta

Before: Epic detail has a focused Epic/Audit style experience that renders
epic-local detail and audit files. Supporting project documents live on disk and
are only read by agents when a phase explicitly needs them. After: the board
loads a bounded set of supporting documents (`.savepoint/PRD.md` and
`.savepoint/Design.md`) through `internal/data`, caches them in board state, and
renders them as a read-only Epic detail subview. The data boundary remains
file-only and the board stays an event/message reducer by loading file contents
through command messages.

## Boundaries

**In scope:**
- Read-only PRD and Design viewing from the Epic detail overlay.
- Keyboard navigation between Epic detail and Release Docs.
- Keyboard navigation between PRD and Design.
- Scrollable wrapped rendering for document bodies.
- Graceful missing-file handling.

**Out of scope:**
- Editing PRD or Design from the board.
- Search, heading outline navigation, or markdown syntax highlighting.
- Loading arbitrary files outside the known supporting docs list.
- Changing router state, task lifecycle rules, or release switching behavior.
- Closing the existing Epic view wrapping defect unless the repair is explicitly
  verified against `v1.2/D017-epic-view-line-wrapping`.

## Quality gates

- The Epic detail overlay defaults to the existing Epic detail behavior.
- Users can switch to Release Docs without changing release, epic, task, or
  router state.
- PRD and Design labels render, and each document body can be selected and
  scrolled independently.
- Missing PRD or Design files render a clear empty state instead of an error
  overlay or panic.
- Long document lines wrap within the content width without awkward overflow.
- `make build && make test` passes.

## Open decisions

- Exact keybinding for switching Epic detail subviews. **Resolved:** the overlay
  uses numbered tab keys today — `1` Epic, `2` Audit (`handleEpicDetailOverlay`,
  `internal/board/update.go`) — not `tab` cycling. Add `3` for Release Docs to
  match the existing pattern; do not introduce a new `tab`-cycle interaction.
  Within Release Docs, switch PRD/Design with `[`/`]` (or left/right) so the
  doc selector does not collide with the `1`/`2`/`3` subview keys.

## Dependencies and risks

- **D017 (`v1.2/D017-epic-view-line-wrapping`, open) is effectively a
  prerequisite for T003.** It is the existing Epic-view wrapping bug in the same
  render path this epic extends, and T003/T004 require clean wrapping — the same
  symptom. Resolve D017 (or fold its wrapping fix into T003) before relying on
  the wrapping quality gate; otherwise the Release Docs body will inherit the
  defect. Note that `WrapText` (`internal/board/util.go`) is **not** ANSI-aware
  and collapses whitespace via `strings.Fields`, so it is not a drop-in for
  rendering raw Markdown bodies (see T003).
