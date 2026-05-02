---
id: E01-tui-optimisation/T001-next-activity-header
status: done
objective: "Show Next Activity indicator in the TUI header bar, right-aligned, displaying the current router state"
depends_on: []
---

# T001: Next Activity Header Display

## Acceptance Criteria

- The header bar displays "Next Activity:" followed by a formatted state string on the right side
- When router state is `task-building`, displays "Build T{NNN} (E{NN}) v{N}" format (e.g., "Build T010 (E06) v1")
- When router state is `audit-pending`, displays "Audit E{NN}" format (e.g., "Audit E06")
- For `epic-design` state, displays "Design E{NN}" format
- For `epic-task-breakdown` state, displays "Plan E{NN}" format
- For `pre-implementation` state, displays "Planning v{N}" format
- The Next Activity text (excluding "Next Activity:" label) is capped at 20 characters maximum, with ellipsis truncation if exceeded
- The display gracefully degrades at narrow terminal widths (header text remains intact, Next Activity truncates but doesn't break layout)
- All existing tests pass

## Implementation Plan

- [x] Add `HeaderRight` style to `internal/styles/styles.go` (dim foreground, right-aligned)
- [x] Create `FormatNextActivity(state *data.RouterState) string` helper in `internal/board/view.go` that:
  - Parses router state and formats compact string based on state type
  - Truncates to max 20 chars using `xansi.Truncate` with ellipsis
- [x] Update `view.go` `View()` method — modify header rendering:
  - Split header into left (icon + "SAVEPOINT") and right sections
  - Use `lipgloss.PlaceHorizontal` or manual string construction to position Next Activity on right
  - Maintain existing `HeaderFrame` width and padding
- [x] Update `main.go` or board initialization to read router state and populate `Model` with parsed state
- [x] Add test in `internal/board/view_test.go` for Next Activity rendering at various widths
- [x] Add test for `FormatNextActivity` covering all router states and truncation behavior
- [x] Run quality gates — `make build && make test` unavailable in this shell; equivalent `go build ./...` and `go test ./...` pass

## Context Log

Files read:
- `internal/board/view.go`
- `internal/board/model.go`
- `internal/styles/styles.go`
- `internal/data/router.go`
- `internal/board/view_test.go`
- `.savepoint/router.md`
- `.savepoint/releases/v1.1/epics/E03-header-activity/E03-Detail.md`

Estimated input tokens: 2800

Notes:
- Uses existing `data.RouterState` model — no new data parsing needed
- Router state already available via `m.RouterTask` in Model — may need to add `*data.RouterState` field or parse separately
- Truncation uses same `xansi` package already imported in view.go
- Max 20 char constraint keeps header clean at all terminal widths
- Audit must-fix applied: stale implementation checklist items were ticked after verifying the implemented source and tests.
- Quality gates verified during audit: `go build ./...` passed; `go test ./...` passed.
