---
name: bubbletea-tui-design
description: Bubble Tea and Lip Gloss TUI design and implementation guide for Savepoint. Use when creating or modifying the Go terminal UI in internal/board or internal/styles, including Bubble Tea models, Update/View behavior, tea.Cmd I/O, keyboard handling, responsive terminal layout, Lip Gloss styling, non-TTY rendering, and Go tests for terminal UI behavior.
---

# Bubble Tea TUI Design

Use this skill for Savepoint's Go terminal UI mechanics. `.savepoint/visual-identity.md` owns the Atari-Noir visual direction, palette, and product design constraints; this skill owns how to implement that direction with Bubble Tea, Lip Gloss, and Go tests.

## Current Stack

- TUI framework: `github.com/charmbracelet/bubbletea`
- Styling: `github.com/charmbracelet/lipgloss`
- Board package: `internal/board`
- Shared styles: `internal/styles`
- Project data source: markdown and YAML frontmatter read through `internal/data`

Do not introduce Ink, React, TypeScript TUI components, or `ink-testing-library` for current board work.

## Quick Workflow

1. Read the active task and its context files first.
2. Read `.savepoint/visual-identity.md` only when the task touches rendering, layout, theme, glyphs, or design-system behavior.
3. Locate the smallest board surface involved:
   - model/state: `internal/board/model.go`
   - event handling: `internal/board/update.go`
   - asynchronous commands/messages: `internal/board/io.go` and adjacent focused files
   - rendering: `internal/board/view.go`, `card.go`, overlays, detail, release, or epic panel files
   - styles: `internal/styles`
4. Keep `Update` as an event reducer. Push filesystem work into `tea.Cmd` helpers that return typed messages.
5. Add or update focused Go tests for the changed branch, render path, and edge case.

## Package Boundaries

| Area | Responsibility |
| --- | --- |
| `internal/board/model.go` | Bubble Tea model state and selected/focused UI state. |
| `internal/board/update.go` | Message dispatch, keyboard handling, and state transitions. |
| `internal/board/io.go` | `tea.Cmd` helpers and typed messages for async file operations. |
| `internal/board/view.go` | Top-level layout assembly and terminal-size decisions. |
| `internal/board/*_overlay.go`, `detail.go`, `epic_panel.go`, `card.go` | Focused render surfaces and interaction-specific helpers. |
| `internal/styles` | Palette, Lip Gloss styles, fallback colors, and shared visual tokens. |
| `internal/data` | Parsing, validation, lifecycle rules, and write helpers. Board code should not duplicate these rules. |

## Bubble Tea Rules

- Keep model fields explicit. Do not hide workflow state in package globals.
- Treat `Update(msg tea.Msg)` as a deterministic reducer over typed messages.
- Return commands for I/O, timers, reloads, and writes. Avoid synchronous filesystem reads or writes inside key branches.
- Prefer typed message structs over stringly status plumbing.
- Dispatch keyboard events by active mode/surface first: modal or overlay, then sidebar/detail, then board columns.
- Make repeated key presses safe. Actions such as refresh, priority write, close overlay, and failed transition handling must be idempotent.
- Preserve router semantics. Browsing around the UI must not silently rewrite `.savepoint/router.md`; use explicit actions such as the priority hotkey.
- Keep quit handling simple and predictable. `ctrl+c` and `q` should exit only when no higher-priority modal behavior owns the key.

## Lip Gloss Layout Rules

- Terminal geometry is fixed-width cell math. Account for borders, padding, margins, and separators before rendering content.
- Focus must not change geometry. Change color, glyphs, or style, not border thickness, padding, or component width.
- Use one border family consistently across related surfaces. Prefer the single-line border style approved by visual identity unless a task explicitly changes it.
- Truncate long text before it can wrap into adjacent surfaces. Treat accidental wrapping as a bug.
- Derive widths from the current terminal size and clamp for narrow screens.
- Use `internal/styles` for colors and shared styles. Do not create one-off hex literals in board rendering code unless the style package is being extended.
- When adding colors, provide true-color, 256-color, and ANSI fallbacks if the surrounding style code uses `lipgloss.CompleteColor`.

## Rendering Rules

- Keep render helpers small and named after the surface they produce.
- Separate data selection from string rendering when the branch logic grows.
- Reinforce color state with text or glyphs so the board remains readable in limited color terminals.
- Avoid wide emoji for structural UI. Prefer stable-width symbols already used by the board.
- Keep copy short. Terminal UI should scan quickly.
- Non-TTY output should remain deterministic and free of control sequences.

## Testing Rules

- Test `Update` by sending Bubble Tea messages and asserting model state plus returned command behavior when practical.
- Test command helpers with temp directories and explicit file fixtures.
- Test rendering with string assertions for stable labels, glyphs, truncation, ordering, and absence of accidental wrapping.
- Cover narrow-width and wide-width layout paths when layout math changes.
- Cover non-TTY fallback behavior when command output changes.
- Prefer targeted package tests during iteration, then run broader gates when the task is complete:

```bash
go test ./internal/board ./internal/styles
go test ./...
```

## Common Traps

- Re-implementing lifecycle rules in board code instead of using `internal/data`.
- Updating router priority as a side effect of navigation.
- Performing file writes directly inside `Update`.
- Letting focus styles change component dimensions.
- Adding a new glyph or color without considering monochrome or low-color terminals.
- Relying on full-screen golden snapshots when a focused string or state assertion would be more stable.
