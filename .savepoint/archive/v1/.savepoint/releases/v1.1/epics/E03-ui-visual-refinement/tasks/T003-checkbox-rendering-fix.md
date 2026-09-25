---
id: E03-ui-visual-refinement/T003-checkbox-rendering-fix
status: done
objective: "Fix checkbox rendering to place checkboxes at sentence start, not at line break positions"
depends_on:
  - E03-ui-visual-refinement/T001-border-resize-fix
---

# T003: Fix Checkbox Rendering — Place at Sentence Start

## Acceptance Criteria

- Checkboxes (`[ ]` / `[x]`) render at the beginning of each task description sentence, not at arbitrary line-break positions
- If a task description is hard-wrapped in markdown at column 80, the checkbox appears only once per sentence
- Multi-sentence tasks show one checkbox per sentence, visually aligned at column 0 of each sentence
- Already-checked tasks show `[x]` at sentence start consistently
- No duplicate checkboxes for the same sentence
- Single-sentence tasks (most common case) render exactly one checkbox at the start, unchanged from current behavior

## Implementation Plan

- [x] Read `internal/data/task.go` — locate the markdown parsing / TUI rendering path for task descriptions
- [x] Read `internal/board/view.go` — locate where task items are rendered in the board columns
- [x] Identify the parsing layer — this is likely where markdown task body text is split into the displayed description string before being passed to rendering
- [x] Edit the identified file(s) — change the sentence-splitting logic to split on sentence boundaries (`. `, `!\n`, `?\n`) rather than on hard line breaks
- [x] In the rendering layer, emit one checkbox per sentence start, aligned to column 0 of each sentence block
- [x] Add tests in the relevant test file for sentence-boundary detection, single-sentence, multi-sentence, and hard-wrapped cases
- [x] Run `make build && make test` to verify no regressions

## Context Log

Files read:
- `internal/data/task.go`
- `internal/data/parser.go`
- `internal/data/parser_test.go`
- `internal/board/view.go`
- `internal/board/detail.go`
- `internal/board/detail_test.go`
- `agent-skills/savepoint-build-task/SKILL.md`
- `agent-skills/ink-tui-design/SKILL.md`
- `.savepoint/visual-identity.md`

Files edited:
- `internal/data/parser.go`
- `internal/data/parser_test.go`
- `internal/board/detail.go`
- `internal/board/detail_test.go`

Estimated input tokens: 8,500

Notes:
- Parser now joins hard-wrapped continuation lines into the current checklist item.
- Detail rendering now splits checklist item text into semantic sentences and emits `[ ]` / `[x]` only at sentence starts; wrapped continuation lines are indented under the sentence text.
- Focused tests: `go test ./internal/data ./internal/board` passed.
- Required literal quality gate: `make build` could not run because `make` is not installed in this PowerShell environment.
- Equivalent gates: `go run ./internal/buildtool build` passed; `go test ./...` passed.
