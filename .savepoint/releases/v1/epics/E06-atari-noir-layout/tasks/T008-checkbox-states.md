---
id: E06-atari-noir-layout/T008-checkbox-states
status: planned
objective: "Parse and render checked/unchecked state for Implementation Plan checklist items"
depends_on: []
---

# T008: Parse and Render Checkbox States in Task Detail

## Acceptance Criteria

- `Task.Checklist` stores the checked state of each item (not just the text)
- `- [x]` markers in the markdown source are parsed as "done" items
- `- [ ]` markers in the markdown source are parsed as "undone" items
- Legacy `- ` (no checkbox prefix) items are treated as undone
- The task detail overlay renders done items with a `☑` (checked box) glyph and undone items with `□` (empty box)
- Done items render in the `TagDone` (green) style; undone items render in `CardMeta` (dim) style
- All existing `make test` tests pass, with updates to any tests that construct `Checklist` values

## Implementation Plan

- [ ] Edit `internal/data/task.go` — add `CheckItem` struct with `Text string` and `Done bool` fields; change `Checklist []string` to `Checklist []CheckItem`.
- [ ] Edit `internal/data/parser.go` — update `extractChecklistSection()` to detect `- [x] ` vs `- [ ] ` vs `- ` prefixes and set `Done` accordingly; strip prefix from `Text`.
- [ ] Update all usages of `Checklist` across the codebase (likely only `board/detail.go` and tests).
- [ ] Edit `internal/board/detail.go` — in the Implementation Plan section, render done items with `☑ ` + green `TagDone` styling and undone items with `□ ` + dim `CardMeta` styling.
- [ ] Run `make build && make test` to verify no regressions.

## Context Log

Files read:
- `internal/data/task.go`
- `internal/data/parser.go`
- `internal/board/detail.go`
- `internal/styles/styles.go`

Estimated input tokens: 800

Notes:
