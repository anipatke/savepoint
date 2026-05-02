# Audit Snapshot: v1.1 E01 TUI Optimisation

Manual snapshot created because the audit CLI snapshot was missing and the router permits one manual snapshot when the CLI is unavailable.

## Scope

- Release: `v1.1`
- Epic: `E01-tui-optimisation`
- Router state at audit start: `audit-pending`
- Next action: `Audit v1.1 E01 TUI optimisation.`

## Epic Files

- `.savepoint/releases/v1.1/epics/E01-tui-optimisation/E01-Detail.md`
- `.savepoint/releases/v1.1/epics/E01-tui-optimisation/tasks/T001-next-activity-header.md`
- `.savepoint/releases/v1.1/epics/E01-tui-optimisation/tasks/T002-rename-epic-design-files.md`
- `.savepoint/releases/v1.1/epics/E01-tui-optimisation/tasks/T003-rename-release-prd.md`
- `.savepoint/releases/v1.1/epics/E01-tui-optimisation/tasks/T004-update-instruction-files.md`
- `.savepoint/releases/v1.1/epics/E01-tui-optimisation/tasks/T005-update-cross-references.md`
- `.savepoint/releases/v1.1/epics/E01-tui-optimisation/tasks/T006-column-and-detail-scrolling.md`
- `.savepoint/releases/v1.1/epics/E01-tui-optimisation/tasks/T007-column-focus-border-stability.md`

## Changed Source And Test Files Reviewed

- `internal/board/layout.go`
- `internal/board/view.go`
- `internal/board/column.go`
- `internal/board/detail.go`
- `internal/board/model.go`
- `internal/board/update.go`
- `internal/board/watch.go`
- `internal/styles/styles.go`
- `internal/board/column_test.go`
- `internal/board/detail_test.go`
- `internal/board/layout_test.go`
- `internal/board/view_test.go`

## Documentation And Workflow Files In Scope

- `.savepoint/Design.md`
- `AGENTS.md`
- `.savepoint/router.md`
- `agent-skills/savepoint-audit/SKILL.md`
- `agent-skills/savepoint-build-task/SKILL.md`
- `agent-skills/savepoint-create-plan/SKILL.md`
- `agent-skills/savepoint-create-task/SKILL.md`
- `agent-skills/savepoint-system-design/SKILL.md`
- `templates/project/AGENTS.md`
- `templates/project/.savepoint/router.md`
- `templates/prompts/task-breakdown.prompt.md`
- `templates/prompts/task-building.prompt.md`
- `templates/prompts/task-planning.prompt.md`

## Implemented Delta Observed

- Board header can show a right-aligned `Next Activity:` label from parsed router state.
- Release PRD naming moved from `PRD.md` to `{release}-PRD.md`.
- Epic design naming moved from per-epic `Design.md` to `E##-Detail.md`.
- TUI columns use height-aware virtual viewport slicing with above/more scroll indicators.
- Detail overlay accepts terminal-height-derived caps and scroll offsets.
- Model state now tracks per-column offsets and detail offsets.
- Column navigation and page keys keep the focused task visible.
- Unfocused columns now use a rounded-border style matching focused column dimensions.
- `internal/styles/styles.go` added `HeaderRight`, `ScrollIndicator`, and `ColumnUnfocused`.

## Drift Notes Found

- `T007-column-focus-border-stability.md`: `internal/styles/styles.go` added exported `ColumnUnfocused` style, not yet in Codebase Map.

## Verification Performed During Audit

- `go build ./...`: passed.
- `go test ./...`: passed.

## Audit Caveats

- The required `make build && make test` gate was not rerun because this environment has used direct Go equivalents for recent task closeout when `make` is unavailable.
- Must-fix follow-up applied: `T001-next-activity-header.md`, `T004-update-instruction-files.md`, and `T005-update-cross-references.md` no longer contain unchecked implementation checklist items.
