---
type: epic-design
status: active
---

# Epic E06: tui-board

## Purpose

Implement `savepoint board`, an Ink-based terminal UI for viewing and moving tasks through the Savepoint workflow, including audit-review mode entry points. The board should make the file-backed state machine usable without replacing direct file editing.

## What this epic adds

- Ink app shell launched by `savepoint board`.
- Five-column Kanban view for `backlog`, `planned`, `in_progress`, `review`, and `done`.
- Detail pane for the selected task.
- Keyboard navigation using arrows and vim-style movement.
- Status transition actions with gate enforcement.
- Manual refresh.
- Theme tokens loaded from `config.yml`.
- Render fallbacks for 24-bit/256-color, 16-color, `NO_COLOR=1`, and non-TTY plain table output.
- Optimistic mtime conflict detection before writes.
- Reducer tests for navigation, selection, and transition behavior.

## Components and files

Expected files introduced or extended by this epic:

| Path                            | Purpose                                            |
| ------------------------------- | -------------------------------------------------- |
| `src/commands/board.ts`         | Command entrypoint and non-TTY fallback.           |
| `src/tui/App.tsx`               | Ink app root.                                      |
| `src/tui/Board.tsx`             | Kanban layout.                                     |
| `src/tui/DetailPane.tsx`        | Selected task details.                             |
| `src/tui/state/*.ts`            | Reducers and derived view state.                   |
| `src/tui/theme/*.ts`            | Theme loading and terminal color fallback mapping. |
| `src/tui/render/plain-table.ts` | Non-TTY table rendering.                           |
| `src/tui/io/*.ts`               | Safe task status write helpers and mtime checks.   |
| `test/tui/**/*.test.ts`         | Reducer and fallback rendering tests.              |

## Architectural delta

Before this epic, Savepoint can read project files and expose CLI commands. After this epic, users can inspect and advance workflow state through a terminal interface.

The TUI uses the domain layer for rules and the filesystem layer for reads/writes. It should not define independent workflow logic.

## Visual and interaction design

The terminal adaptation of Atari-Noir applies here:

- Dark background with warm off-white text where terminal color support allows.
- One accent color per major semantic area.
- Quiet borders and compact uppercase headings.
- Focus shown through accent borders or subtle background tint.
- Colored glyphs should reinforce state but never be the only signal.
- Scanlines and web-style glow are skipped in terminal rendering.

## Boundaries

In scope:

- Board navigation and status transitions.
- Gate enforcement for dependencies and status rules.
- Plain table fallback when not in a TTY.
- Manual `r` refresh.
- Initial audit-review mode entry point if audit proposals exist.

Out of scope:

- File watching.
- Search.
- Mouse interaction.
- Drag-and-drop.
- Full audit proposal diff implementation beyond what `audit-pipeline` needs.
- Snapshot tests for rich TUI rendering.

## Quality gates

- Reducers and transition behavior should be covered by tests.
- Rendering tests should focus on deterministic fallback/plain output, not brittle full-screen snapshots.
- Conflict detection should be tested without relying on real-time delays where possible.

## Design constraints

- Keep workflow state in files; the TUI holds only derived view state.
- Use mtime-based optimistic concurrency.
- Do not create a lockfile.
- Keep terminal controls discoverable but avoid loading screens with long instructions.
