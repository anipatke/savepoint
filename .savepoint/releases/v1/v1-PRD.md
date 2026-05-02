---
version: 1
name: "Simplified Board"
status: in_progress
---

# Release v1 — Simplified Board

## Simplification Objective

The original v1 scope (11 epics, 5 CLI commands, full audit pipeline) proved over-engineered. This release strips savepoint back to its essential value: a board for tracking AI-driven development work through phases.

**Stack change:** Ink/TypeScript → Bubble Tea/Go. Full rewrite.

**What stays:**
- `savepoint board` — Bubble Tea TUI for viewing and moving tasks
- The `.savepoint/` data structure: router → release → epic → task
- AGENTS.md as the agent bootstrap file

**What goes:**
- `init`, `audit`, `doctor` commands
- Full 6-step audit pipeline
- Template scaffolding and prompt generation
- Quality gates, divergence thresholds, snapshot generation
- Complex router state machine (6 states → 3 states)
- TypeScript, React, Ink, npm, tsup, vitest

**What replaces the audit pipeline:**
- Tasks in `in_progress` now carry a `phase`: `build → test → audit → done`
- The board itself enforces phase progression
- Space advances phase, backspace retreats
- Must reach `audit` phase before advancing to `done`

## What ships in v0.2.0

- `savepoint board` — the only CLI command (Go binary)
- Phase-aware Kanban board: planned / in_progress (build/test/audit) / done
- Task status transitions with phase gating
- Non-TTY plain table fallback
- Atari-Noir theme, fully overridable via `config.yml`
- Go 1.23+, single binary, `go build`

## Epic breakdown

| #   | Epic name                    | Purpose |
| --- | ---------------------------- | ------- |
| 01  | `E01-go-setup`               | Initialize Go module, project structure, dependencies, build loop |
| 02  | `E02-data-readers`           | Parse markdown/yaml, read .savepoint data (router, config, tasks) |
| 03  | `E03-board-tui-core`         | Bubble Tea model, update loop, view, styles, responsive layout |
| 04  | `E04-board-components`       | Columns, cards with phase glyphs, epic panel, detail overlay, dropdowns, help |
| 05  | `E05-phase-transitions`      | Phase stepping, gates, task frontmatter write-back, router state write-back |

## Success criteria

- `savepoint board` launches and renders the phase-aware Kanban board.
- Space/backspace correctly advances/retreats task phases.
- Cannot advance to `done` without passing through `audit` phase.
- `go build` produces a working binary.
- `go test ./...` passes.
- No dead code, no unused dependencies, no references to deleted commands.

## Risks

- **Agent compliance.** Simpler structure should improve compliance. AGENTS.md must be crystal clear.
- **Phase model clarity.** Users must understand build/test/audit as sub-states of in_progress. Board rendering must make this obvious.
- **Go ecosystem.** New stack requires different testing patterns, dependency management, build tooling.
