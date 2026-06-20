---
id: E33-release-docs-view/T003-release-docs-renderer
status: done
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

- [x] Epic detail still renders as the default overlay content.
- [x] Release Docs renders PRD and Design selectors using existing board styles.
- [x] The selected document body renders in a scrollable viewport.
- [x] Missing documents render a clear read-only empty state.
- [x] Document content wraps within the available body width after borders and
      padding are accounted for.
- [x] Rendering does not introduce a new panel/card nesting style inconsistent
      with the board.

## Implementation Plan

- [x] Extend the Epic detail tab/header renderer to include Release Docs.
- [x] Add a Release Docs body renderer in `internal/board/epic_panel.go`.
- [x] Wrap document bodies with a line-preserving wrapper. NOTE: `WrapText`
      (`internal/board/util.go`) is **not** ANSI-aware (it measures with
      `len([]rune)`) and collapses whitespace/newlines via `strings.Fields`, so
      it flattens Markdown paragraph breaks, indentation, and code blocks — do
      not use it as-is for raw doc bodies. Wrap per source line (preserving blank
      lines), and if styled content needs measuring, use the ANSI-aware
      `xansi`/`lipgloss.Width` approach already used in `view.go`. Extract a
      shared helper rather than ad hoc string slicing.
- [x] Apply height-aware viewport slicing consistent with task and epic detail
      overlays.
- [x] Add missing-doc and empty-doc render branches.
- [x] Add focused rendering tests for selector labels, selected state, empty
      state, and width-bounded output.

## Context Log

- Read: `internal/board/epic_panel.go`, `view.go`, `update.go`, `detail.go`,
  `util.go`, `layout.go`, `model.go`, `theme.go`, `epic_panel_test.go`, and
  `internal/data/release_doc.go` (T001) plus the T002 state/key handling to
  match the existing Detail/Audit overlay render pattern.
- `internal/board/epic_panel.go`: generalized `renderTabIndicator` to three
  tabs (Detail/Audit/Docs) via a new `tabLabel` helper; updated the Detail and
  Audit footers to advertise `3:Docs`. Added `RenderEpicReleaseDocs` (title,
  tab strip, PRD/Design selector, scrollable body, footer, viewport via
  `visibleDetailLines`), `renderReleaseDocSelector`, `releaseDocBody`
  (no-docs/missing/empty/normal branches), and the line-preserving wrapper
  `renderReleaseDocBody` → `wrapDocLine` → `styledWrap`/`leadingWhitespace`.
  The wrapper strips frontmatter, styles `#`/`##`/`###` headings, preserves
  blank lines and leading indentation, and wraps each source line to the body
  width (no `strings.Fields` flattening across lines).
- `internal/board/view.go`: `OverlayEpicDetail` now switches on `EpicDetailTab`
  with a `case 2` dispatching to `RenderEpicReleaseDocs`, using the per-doc
  scroll offset.
- `internal/board/update.go`: added `selectedReleaseDocOffset()` so the view
  reads the selected doc's stored offset.
- Tests (`epic_panel_test.go`): added header/selector/tab/selected-body/missing/
  empty/no-docs/footer cases for `RenderEpicReleaseDocs`, width-bound +
  blank-line + indentation + frontmatter cases for `renderReleaseDocBody`, and
  `TestView_releaseDocsTabRendered`.
- Wrapping gate: D017 (epic-view line wrapping) was already fixed on this
  branch; the Release Docs body uses an independent line-preserving wrapper
  (not the paragraph-reflow path), and the width-bound test confirms no line
  exceeds the body width including long path-like tokens.
- Quality gates: `make build && make test` pass; `go vet ./internal/board/`
  clean.
