---
id: E03-ui-visual-refinement/T002-next-activity-below-header
status: done
objective: "Reposition Next Activity below header as a dedicated line with phase-aligned styling"
depends_on:
    - E03-ui-visual-refinement/T001-border-resize-fix
---

# T002: Reposition Next Activity Below Header with Phase-Aligned Styling

## Acceptance Criteria

- Next Activity renders as a dedicated full-width line immediately below the header bar
- The line is hidden (not rendered) when no activity is active or router state is idle
- The line uses `phase:` prefix mapped to existing footer phase styles:
  - `task-building` → FooterPhaseBuild (orange + bold)
  - `audit-pending` → FooterPhaseAudit (green + bold)
  - `pre-implementation` / `epic-design` / `epic-task-breakdown` → FooterPhasePlan (purple + bold)
- Format: `"PLAN: Build T010 (E06) v1"` — phase prefix in styled tag, activity text follows
- Activity text is the `next_action` value from router state (already populated)
- At narrow widths (< 60 chars), the line truncates gracefully: `"PLAN: Build T0…"`
- Existing header layout is preserved — nothing shifts or wraps incorrectly

## Implementation Plan

- [x] Read `internal/board/view.go` — locate current header rendering and any existing Next Activity right-aligned logic
- [x] Read `internal/styles/styles.go` — verify FooterPhaseBuild/FooterPhaseAudit/FooterPhasePlan are available for reuse, or add them if missing
- [x] Read `internal/data/router.go` — confirm RouterState model exposes `state` and `next_action` fields
- [x] Edit `internal/board/view.go` — add `renderNextActivityLine(state)` that returns a styled string using phase-mapped style
- [x] Edit `internal/board/view.go` — insert the rendered line immediately after header output in the `View()` function
- [x] Edit `internal/board/layout.go` — adjust vertical layout height calculation to account for the new line
- [x] Add tests in `internal/board/view_test.go` — verify phase mapping, truncation, and hidden state
- [x] Run `make build && make test` to verify no regressions

## Context Log

Files read:
- `.savepoint/router.md`
- `agent-skills/savepoint-build-task/SKILL.md`
- `.savepoint/releases/v1.1/epics/E03-ui-visual-refinement/E03-Design.md`
- `.savepoint/releases/v1.1/epics/E03-ui-visual-refinement/tasks/T002-next-activity-below-header.md`
- `.savepoint/visual-identity.md`
- `agent-skills/ink-tui-design/SKILL.md`
- `internal/board/view.go`
- `internal/board/layout.go`
- `internal/styles/styles.go`
- `internal/data/router.go`
- `internal/board/view_test.go`
- `internal/board/layout_test.go`

Files edited:
- `internal/board/view.go`
- `internal/board/layout.go`
- `internal/board/view_test.go`
- `internal/board/layout_test.go`
- `.savepoint/releases/v1.1/epics/E03-ui-visual-refinement/tasks/T002-next-activity-below-header.md`
- `.savepoint/router.md`

Estimated input tokens: 9800

Notes:
- Key design decision from user: placement below header as dedicated line (not right-aligned in header); reuse existing FooterPhase* styles
- Acceptance criteria verified by `go test ./internal/board` covering dedicated line ordering, hidden nil/idle/empty states, phase tag mapping, `next_action` text, narrow truncation, and header preservation.
- Quality gates: `make build` could not run because `make` is not installed in this Windows shell. Equivalent `go run ./internal/buildtool build` passed. Equivalent `go test ./...` passed.
