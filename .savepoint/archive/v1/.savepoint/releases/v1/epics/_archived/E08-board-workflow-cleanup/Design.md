---
type: epic-design
status: active
---

# Epic E08: board-workflow-cleanup

## Purpose

Make the board usable as the primary release workflow surface before `0.1.0` ships by adding epic navigation, readable task presentation, acceptance criteria, and a focused task detail view.

## What this epic adds

- Immersive full-screen TUI that utilizes the terminal's alternate screen buffer, hiding shell history while active.
- Epic selection in a fixed left rail (minimum 25 columns wide) so users can browse and work across release epics without manually editing router state.
- A clear separation between selectable epic items and task cards, with `Tab` to switch focus between the two panes.
- Readable task card labels that prioritize task title/slug/objective over full `epic/task` IDs.
- Required task-level acceptance criteria that define what "done" means.
- A workflow rule that the implementation plan is the build checklist for satisfying acceptance criteria.
- A full-screen task detail modal that shows objective, acceptance criteria, implementation checklist, dependencies, status, and transition affordances.
- Updated templates and planning prompts so new task files include acceptance criteria by default.
- Migration/compatibility behavior for existing task files that do not yet include acceptance criteria.

## Backlog.md & OpenCode inspiration

Backlog.md and OpenCode are useful inspiration for release-blocking product behaviors: task files are human-readable Markdown, tasks carry clear acceptance criteria, and the visual board/detail views make task status and task content inspectable without opening files directly. The TUI acts as a standalone, immersive application. Savepoint should borrow those product lessons while preserving Savepoint's epic-first router and audit workflow.

## Components and files

Expected files introduced or extended by this epic:

| Path                                 | Purpose                                                               |
| ------------------------------------ | --------------------------------------------------------------------- |
| `src/domain/task.ts`                 | Task document model with acceptance criteria parsing and validation.  |
| `src/readers/tasks.ts`               | Epic task-set reader updated for acceptance criteria compatibility.   |
| `src/tui/board-data.ts`              | Cross-epic board data and readable task/epic display metadata.        |
| `src/tui/state/*.ts`                 | Board navigation state for epic selection, task selection, and popup. |
| `src/tui/{App,Board,DetailPane}.tsx` | Board keyboard flow, readable cards, layout, and task detail modal.   |
| `src/tui/render/plain-table.ts`      | Non-TTY rendering for epic/task separation and acceptance criteria.   |
| `src/commands/board.ts`              | Alternate screen buffer management and immersive launch.              |
| `templates/project/**`               | Default workflow docs and task templates with acceptance criteria.    |
| `templates/prompts/*.prompt.md`      | Agent prompts updated to require acceptance criteria before planning. |
| `.savepoint/router.md`               | Workflow instructions for acceptance criteria and detail navigation.  |
| `test/domain/*.test.ts`              | Task schema and compatibility tests.                                  |
| `test/readers/*.test.ts`             | Cross-epic task reading tests.                                        |
| `test/tui/**/*.test.ts(x)`           | Board navigation, card labels, popup detail, and render tests.        |
| `test/templates/*.test.ts`           | Template and prompt integrity tests.                                  |

## Architectural delta

Before this epic, the board is scoped to the active epic and shows task IDs as the primary task label. The board prints inline to standard output. Task files have an objective and implementation plan, but no explicit acceptance criteria contract.

After this epic, the board becomes a release-level immersive navigation surface using the alternate screen buffer. Epics are selectable in a left rail, tasks are readable at a glance, and a full-screen task detail modal is available. Task documents gain acceptance criteria as a first-class section, with the implementation plan explicitly subordinate to those criteria.

The router remains the source of truth for the active build task. Browsing another epic in the board must not silently change the active build task unless the user performs an explicit workflow action (`Enter` on the left rail).

## Resolved design decisions

1. **Immersive TUI & Layout:** The board will enter the terminal's alternate screen buffer (`\x1b[?1049h`), restoring shell history cleanly on exit. Layout consists of a fixed Left Rail (min 25 columns width) for Epics, and a Main Board for tasks. Focus is toggled via `Tab` / `Shift+Tab`.
2. **Epic Interaction:** Epics display status via glyphs and colors: `✓` (Completed/Green), `▶` (Active/Yellow), `○` (Planned/Dim). Browsing previews tasks (read-only). Pressing `Enter` on an Epic activates it in `router.md`.
3. **Task Mutation Constraints:** Tasks can only be advanced (`Space`) or retreated (`Backspace`) if the currently browsed epic matches the Active Router Epic. Otherwise, the action is blocked with a warning.
4. **Task Detail Modal:** A full-screen overlay triggered by `Enter` from the board. It displays task details, and the implementation checklist is strictly read-only for `v0.1.0`. Closed via `Escape`.
5. **Acceptance Criteria Parsing:** The `task.ts` parser will extract the `## Acceptance Criteria` markdown section. If missing, it safely returns an empty array for backward compatibility.

## Boundaries

In scope:

- Release-level epic browsing in an immersive full-screen board.
- Readable task labels and full ID demotion into detail/metadata.
- Required acceptance criteria for newly planned tasks.
- Compatibility for existing tasks missing acceptance criteria.
- Full-screen task detail modal in the TUI and plain fallback output.
- Prompt/template updates that make acceptance criteria part of normal workflow.
- Visual bugs to be fixed - scrolling up and down fragments the board, the task surface is too busy for a four-pane workflow, and the TUI should feel like a standalone application.

Out of scope:

- Web UI.
- Search.
- File watching.
- Drag-and-drop or mouse interactions.
- Editing every historical task file by hand unless tests/templates require fixture updates.
- AI semantic review or audit pipeline behavior beyond prompt/template wording.
- Interactive toggling of markdown implementation checklists within the TUI (deferred to future release).

## Quality gates

- Tests cover task documents with and without acceptance criteria.
- Tests cover cross-epic board data loading and epic switching.
- Tests prove readable labels are shown on task cards while full IDs remain visible in detail.
- Tests cover the task detail modal content, full-screen takeover, and keyboard flow.
- Template tests prove new task scaffolds include acceptance criteria before implementation plans.
- Existing board transition behavior remains covered after keybinding changes.

## Design constraints

- Keep the board fast and bounded; release-level browsing must not load unrelated source code.
- Preserve the file-first workflow. The board reads and writes Markdown-backed state; it does not create a hidden database.
- Make navigation explicit. Browsing an epic and changing the router's active task are different actions.
- Keep task files independently buildable and objective-led.
- Acceptance criteria describe observable outcomes; implementation plans describe the build steps to satisfy them.
- Design the layout for terminal constraints: safe truncation, fixed regions, and consistent with Atari-Noir terminal palette rules.
