---
id: E33-release-docs-view/T003-release-docs-renderer
status: planned
objective: Render the Release Docs subview inside the Epic detail overlay
depends_on:
  - E33-release-docs-view/T002-epic-doc-subview-state
complexity_tier: medium
complexity_reason: "Extends a styled TUI overlay with selectors, empty states, wrapping, and viewport slicing."
---

> **Wrapping prerequisite:** this task's wrapping acceptance criteria overlap
> with the open defect `v1.2/D017-epic-view-line-wrapping` (same render path).
> Resolve D017 first, or fold its wrapping fix into this task — otherwise the
> Release Docs body inherits the same awkward-wrapping bug. See E33-Detail.md
> "Dependencies and risks".

# T003: Release Docs Renderer

## Problem

The Epic detail overlay needs a compact read-only document view that feels
native to the board and handles long Markdown content inside constrained
terminal widths.

## Context Files

- `internal/board/epic_panel.go`
- `internal/board/layout.go`
- `internal/board/util.go`
- `internal/board/theme.go`
- `internal/board/epic_panel_test.go`

## Acceptance Criteria

- [ ] Epic detail still renders as the default overlay content.
- [ ] Release Docs renders PRD and Design selectors using existing board styles.
- [ ] The selected document body renders in a scrollable viewport.
- [ ] Missing documents render a clear read-only empty state.
- [ ] Document content wraps within the available body width after borders and
      padding are accounted for.
- [ ] Rendering does not introduce a new panel/card nesting style inconsistent
      with the board.

## Implementation Plan

- [ ] Extend the Epic detail tab/header renderer to include Release Docs.
- [ ] Add a Release Docs body renderer in `internal/board/epic_panel.go`.
- [ ] Wrap document bodies with a line-preserving wrapper. NOTE: `WrapText`
      (`internal/board/util.go`) is **not** ANSI-aware (it measures with
      `len([]rune)`) and collapses whitespace/newlines via `strings.Fields`, so
      it flattens Markdown paragraph breaks, indentation, and code blocks — do
      not use it as-is for raw doc bodies. Wrap per source line (preserving blank
      lines), and if styled content needs measuring, use the ANSI-aware
      `xansi`/`lipgloss.Width` approach already used in `view.go`. Extract a
      shared helper rather than ad hoc string slicing.
- [ ] Apply height-aware viewport slicing consistent with task and epic detail
      overlays.
- [ ] Add missing-doc and empty-doc render branches.
- [ ] Add focused rendering tests for selector labels, selected state, empty
      state, and width-bounded output.

## Context Log

Pending.
