---
id: E06-atari-noir-layout/T007-detail-card-fixes
status: done
objective: "Add line breaks under Acceptance Criteria and Implementation Plan headings in task detail, and word-wrap card titles instead of truncating"
depends_on: []
---

# T007: Detail Overlay Spacing and Card Word-Wrap

## Acceptance Criteria

- Task detail overlay shows a blank line between `"Acceptance Criteria:"` heading and the first bullet item
- Task detail overlay shows a blank line between `"Implementation Plan:"` heading and the first checklist item
- Card titles word-wrap to multiple lines when too long for the available card width, instead of being truncated with `…`
- Cards grow vertically to accommodate wrapped titles
- Short titles remain single-line (no visual change for titles that fit)
- All existing card visual structure (phase glyph, ID line, focus border) is preserved

## Implementation Plan

- [x] Edit `internal/board/detail.go` — add `""` entry after the `Acceptance Criteria:` heading line and after the `Implementation Plan:` heading line.
- [x] Edit `internal/board/detail.go` — make `wrapText` and `splitLongWord` exported (capitalize) so they can be reused from `card.go`.
- [x] Edit `internal/board/card.go` — replace `truncate(t.Title, inner)` with a multi-line approach that word-wraps the title across available width.
- [x] Edit `internal/board/card.go` — combine wrapped title lines with newlines in the card content (joining with `\n`).
- [x] Run `make build && make test` to verify no regressions.

## Context Log

Files read:
- `internal/board/detail.go`
- `internal/board/card.go`
- `internal/board/column.go`

Estimated input tokens: 600

Notes: Exported WrapText/SplitLongWord from detail.go; updated card_test.go TestRenderCard_titleTruncated → TestRenderCard_titleWraps. Build and tests pass.
