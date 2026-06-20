---
id: v1.2/D017-epic-view-line-wrapping
release: v1.2
status: resolved
severity: medium
title: "Epic view text wraps at awkward positions"
---

# D017: Epic View Text Wraps at Awkward Positions

## Symptom

In the Epic view, a line is breaking in weird places instead of wrapping at
natural word or content boundaries.

## Expected Behavior

Epic view lines should wrap predictably within the available panel width,
preserving readable word boundaries and avoiding visually fragmented text.

## Reproduction

1. Open the board.
2. Navigate to the Epic view.
3. Inspect long release, epic, task, or detail lines.
4. Observe one or more lines breaking at awkward positions.

## Impact

The Epic view becomes harder to scan and can make supporting release context or
task detail text look broken even when the underlying content is correct.

## Fix Plan

Investigate the Epic view detail rendering and width calculations, especially
where styled text, separators, or wrapped content are measured. Apply wrapping
at the rendered panel content width and add coverage for long lines that should
wrap cleanly.

## Acceptance Criteria

- [x] Epic view text wraps at readable word boundaries within the available
      panel width.
- [x] Styled content width calculations do not cause premature or awkward line
      breaks.
- [x] Regression coverage exercises the line that previously wrapped badly.

## Resolution Notes

**Root cause (primary).** The Epic detail body (`epicDetailBody`/`epicAuditBody`,
`internal/board/epic_panel.go`) rendered Markdown **line by line**, wrapping each
source line independently with `WrapText`. The epic detail files are themselves
hard-wrapped at ~76 columns, while the overlay content width is narrower
(`overlayW - 4`). So every source line that was a few columns too wide shed its
last word or two onto a ragged orphan line — the Purpose/description paragraph
rendered as full line, then "it does", full line, then "that", full line, then
"project PRD", etc. That is the awkward wrapping the report describes. Re-wrapping
per source line, not the per-token split, was the real cause; this was confirmed
by rendering the live `E33-Detail.md` through `RenderEpicDetail` at width 80.

**Fix (primary).** Added a `paragraphFlusher` helper that reflows consecutive
prose source lines into a single paragraph before wrapping, flushing on blank
lines, headings, table rows, and list items (each list item is its own reflowed
block so adjacent bullets stay separate). `epicDetailBody` and `epicAuditBody`
now buffer prose through it. The Purpose paragraph now fills each line fully with
only a natural short final line.

**Root cause (secondary).** Even once reflowed, a single token longer than the
width was hard-cut by rune count by `SplitLongWord` at an arbitrary column —
long file paths and task/defect IDs (`internal/board/epic_panel_test.go`,
`v1.2/D017-epic-view-line-wrapping`) snapped mid-token (`…epic_panel_test`/`.go`).

**Fix (secondary).** Rewrote `SplitLongWord` to break an over-long token after
structural separators (`/`, `\`, `-`, `_`, `.`, `:`) when possible, packing
separator-terminated segments greedily and falling back to a rune cut only for an
unbroken run. Output now reads `…epic_panel_`/`test.go` and `…line-`/`wrapping`.

Both fixes keep every emitted line ≤ width (no styled-width regression); they
change only *where* breaks land. The wrap path is shared, so the task/epic detail
overlays and task cards (`detail.go`, `card.go`) benefit too.

**Coverage.**
- `internal/board/epic_panel_test.go`: `TestEpicDetailBody_reflowsHardWrappedParagraph`
  feeds a paragraph pre-wrapped near 70 cols and asserts the greedy-fill
  invariant (no interior line could pull the next line's first word) plus the
  width bound — the line that previously wrapped badly.
- `internal/board/util_test.go`: the two long-token regression cases plus
  word-boundary, width-bound, hard-cut, and separator-packing cases.

**Follow-on (sibling render paths).** Auditing the other detail overlays for the
same class of bug found the opposite failure in `RenderDefectDetail`
(`internal/board/defect_detail.go`) and the `Description` field of `RenderDetail`
(`internal/board/detail.go`): each passed a whole multi-line section to a single
`WrapText`, whose `strings.Fields` collapses every newline — so numbered steps
and `- [x]` checklist items rendered as one run-on paragraph. Extracted
`renderSectionBody` (reflows prose, but keeps blank lines, bullets, `N.` ordered
items, and `### ` sub-headings as breaks) and routed both overlays through it.
The epic body's orphan re-wrap could not occur there (they reflowed whole
sections), and the long-token cut was already fixed centrally; this closes the
structure-flattening variant before T003 builds on the shared path. Coverage:
`TestRenderDefectDetail_preservesListStructure` and
`TestRenderSectionBody_reflowsAndKeepsItems`.

**Note.** The separate `WrapText` weaknesses called out for E33 T003 — not
ANSI-aware and `strings.Fields` collapsing whitespace/newlines — remain by
design; current callers pass unstyled text and *want* prose reflow. T003's *raw
Markdown* body renderer (code blocks, indentation) should still introduce its own
line-preserving wrapper rather than reuse `WrapText` directly.
