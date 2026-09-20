---
id: E32-header-release-indicator/T001-header-release-label
status: done
objective: Render the selected release in the header's top-right cluster, always visible, next to the defect count
depends_on: []
complexity_tier: low
complexity_reason: "Single render function change plus right-align math and width-fallback tests"
---

# T001: Header Release Label

## Problem

`renderHeader` (`internal/board/view.go:110`) only builds a right-hand cluster
(`⚠ N open`) when `openDefectCount() > 0`, and never shows which release is
selected. `m.SelectedRelease` is available but unused in the header, so a user
switching releases with `R` has no persistent indicator of where they are.

## Context Files

- `internal/board/view.go`
- `internal/board/view_test.go`
- `internal/styles/` (header styles)

## Acceptance Criteria

- [x] The header right cluster shows the selected release (e.g. `v1.2`) drawn
      from `m.SelectedRelease`.
- [x] The release is shown at all open-defect counts, including zero (the right
      cluster is no longer gated solely on `count > 0`).
- [x] When there are open defects, the release and the `⚠ N open` count render
      together in the top-right (e.g. `v1.2 │ ⚠ 3 open`) with a single separator.
- [x] Existing right-alignment is preserved: the gap math still right-aligns the
      cluster, and at narrow widths it degrades without overflow or a negative
      repeat count (falls back to left-only, as the current code does when
      `gap <= 0`).
- [x] When `m.SelectedRelease` is empty, the header renders without a stray
      separator.
- [x] Tests cover: release shown with defects, release shown at zero defects,
      empty-release case, and narrow-width fallback.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] In `renderHeader`, build a `right` string unconditionally from
      `m.SelectedRelease` (when non-empty), appending ` │ ⚠ N open` only when
      `openDefectCount() > 0`. Style the release with an existing header style
      (or add a small `HeaderRelease` style in `internal/styles`).
- [x] Keep the existing `inner := w - 2` / `gap := inner - width(left) -
      width(right)` right-align logic; only render the gap-filled line when
      `gap > 0`, otherwise fall back to `left` (current behaviour) so narrow
      terminals never overflow.
- [x] Guard the empty-release and zero-defect combinations so no dangling
      separator is emitted.
- [x] Add/extend `view_test.go` cases: defects+release, zero-defects+release,
      empty release, narrow width. Assert the rendered string contains the
      release label and that width never exceeds `w`.
- [x] Run `make build && make test`.

## Notes

- Reuse `lipgloss.Width` for measurement (already imported in view.go) so styled
  ANSI does not break the alignment math.
- Decision per epic Open decisions: same-line, release first, release-only at
  zero defects. Stacked second line is out of scope unless revisited.

## Context Log

**Read:** `internal/board/view.go` (`renderHeader`, `openDefectCount`),
`internal/board/view_test.go` (existing header tests + `plainTerminal` helper from
`detail_test.go`), `internal/styles/styles.go` (header styles, colour vars).

**Edited:**
- `internal/styles/styles.go` — added `HeaderRelease` style (`clrText` foreground)
  beside `HeaderRight`.
- `internal/board/view.go` — `renderHeader` now builds a `[]string` right cluster:
  the styled release (when `m.SelectedRelease != ""`) followed by `⚠ N open` (when
  `openDefectCount() > 0`), joined by a single ` │ ` separator. The existing
  `inner := w - 2` / `gap` right-align math and the `gap > 0` fallback to
  left-only are unchanged, so narrow widths and the all-empty case still degrade
  to the wordmark alone with no stray separator.

**Tests added** (`view_test.go`):
- `TestView_headerShowsReleaseWithDefects` — asserts `v1.2 │ ⚠  1 open`.
- `TestView_headerShowsReleaseAtZeroDefects` — release shown, no separator
  (asserts against `renderHeader` output, not the full view, so the footer
  phase line's `│` doesn't false-positive).
- `TestView_headerEmptyReleaseNoSeparator` — empty release emits no separator but
  still shows the defect signal.
- `TestView_headerReleaseNarrowWidthNoOverflow` — every header line width ≤ `w`
  at width 24.

**Quality gates:** `make build && make test` — all packages pass.
