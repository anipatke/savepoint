---
id: R001
title: Simplified Board
status: done
legacy_completion:
    source_path: .savepoint/releases/v1/v1-PRD.md
    archive_path: .savepoint/archive/v1/.savepoint/releases/v1/v1-PRD.md
    sha256: 0c93d7cbf11a157d30d6143b61938cfd14ca59aac5ad0a70c9b4bc266f84ae99
---
## Outcome

The V1 release promise is preserved in the Legacy Source section below. Source: `.savepoint/releases/v1/v1-PRD.md`.

## Why

This Release carries the V1 delivery boundary forward with a stable V2 identity.

## Success Conditions

- Converted Objectives sourced from this V1 release reference this Release identity.
- Historical completion, when present, remains typed archive evidence rather than a V2 Check.

## Boundaries

Release membership is derived from Objective records; Tasks remain owned by Objectives.

## Legacy Source (verbatim)

Source path: `.savepoint/releases/v1/v1-PRD.md`

SHA-256: `0c93d7cbf11a157d30d6143b61938cfd14ca59aac5ad0a70c9b4bc266f84ae99`

Original frontmatter:

```yaml
version: 1
name: "Simplified Board"
status: done
```

````markdown


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
````
