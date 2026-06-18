---
type: epic-design
status: audited
---

# E32: Header Release Indicator

## Purpose

The board header shows the SAVEPOINT wordmark and an open-defect count, but never
tells the user which release they are looking at. With multiple releases on disk
(v1, v1.1, v1.2, v1.3) and `R` switching between them, there is no persistent
visual indicator of the current release. This epic adds the selected release to
the header's top-right cluster, next to the defect count.

## What this epic adds

- The current release (`m.SelectedRelease`) rendered in the top-right of the
  header, alongside the open-defect count.
- A header that shows the release even when the open-defect count is zero (today
  the right cluster only appears when there are open defects).

## Components and files

| Module | Purpose |
|--------|---------|
| `internal/board/view.go` | `renderHeader` right-cluster layout + release label |
| `internal/board/view_test.go` | Header rendering tests (release shown, alignment, narrow width) |
| `internal/board/styles` (`internal/styles`) | Reuse/extend header styles if a distinct release style is wanted |

## Architectural delta

Before: `renderHeader` (`internal/board/view.go:110`) builds `left` (icon +
wordmark) and an optional `right` (`⚠ N open`) only when `openDefectCount() > 0`,
then gap-fills between them. `m.SelectedRelease` is available but unused in the
header. After: the right cluster always renders and leads with the release label
(e.g. `v1.2 │ ⚠ 3 open`, or just `v1.2` at zero defects), keeping the existing
right-alignment gap math and narrow-width fallback.

## Boundaries

**In scope:**
- Header rendering of the release label and the always-on right cluster, plus
  tests.

**Out of scope:**
- The release switcher overlay (`R`) and release discovery — unchanged.
- Per-release theming/colour beyond an optional header style.
- Footer phase line (`Build v1.2 …`) — already shows release context separately.

## Quality gates

- The header shows the selected release at all defect counts, including zero.
- Right-cluster alignment is preserved and degrades gracefully at narrow widths
  (no overflow, no negative gap).
- `make build && make test` passes.

## Open decisions

- Same-line (`v1.2 │ ⚠ 3 open`) vs stacked second header line under the defect
  count. Default: same line, release first; fall back to release-only when no
  open defects. Confirm in task breakdown.
