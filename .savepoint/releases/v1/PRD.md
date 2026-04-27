---
version: 1
name: "MVP"
status: in_progress
---

# Release v1 — MVP

## What ships in v0.1.0

The first public release. Establishes the workflow end-to-end with the audit loop minus the AI semantic-review layer.

### In scope

- `savepoint init` — scaffolds `.savepoint/`, prints magic prompt to stdout + clipboard, optional one-shot dev-deps install.
- `savepoint board` — Ink-based 5-column Kanban TUI with detail pane and audit-review mode.
- `savepoint audit` — full 6-step pipeline (quality gates → snapshot → diff brief → reconcile → review → commit).
- `savepoint doctor` — integrity checks + ad-hoc quality-gate run.
- File-only architecture, agent-agnostic via Router Pattern.
- Atari-Noir default theme, fully overridable.
- TypeScript / ESM / Node 20.10+, single-package repo, `tsup` build.

### Out of v0.1.0 scope (deferred to later releases)

- AI semantic review in audit reconcile (→ v0.2.0)
- Broader language quality-gate presets beyond TS (→ v0.2.0)
- File watching in TUI (→ v0.3.0)
- Search in TUI (→ v0.3.0)
- MCP server (→ v1.0.0)

## Proposed epic breakdown

These are the proposed epics for v1. Order matters: each builds on the prior. The first epic is scaffolding (per the savepoint convention). Confirm or revise before any epic gets a `Design.md`.

| # | Epic name | Purpose |
|---|---|---|
| 01 | `scaffolding` | TS project init, lint/typecheck/test setup, `tsup` build, package.json, repo skeleton. **First audit establishes the codebase baseline.** |
| 02 | `data-model` | YAML frontmatter parsing, task/epic/release readers, status state machine, dependency resolution, cycle detection. Pure logic, no I/O beyond fs reads. |
| 03 | `cli-foundation` | Argument parsing, command dispatch, `--help` / `--version`, exit codes, non-TTY detection. The 5-command shell with stub bodies. |
| 04 | `init-command` | `savepoint init` — scaffold writer, magic-prompt printer, clipboard integration, one-shot dev-deps install, language detection. |
| 05 | `tui-board` | Ink board: 5-column Kanban, detail pane, navigation, status transitions with gate enforcement, theme system, render fallbacks. |
| 06 | `audit-pipeline` | 6-step audit pipeline (without AI semantic review): snapshot generation, proposal scaffolding, TUI review mode, commit step, audit-log. |
| 07 | `doctor-command` | Integrity checks (corrupt YAML, broken deps, cycles, orphans, duplicates) + ad-hoc quality-gate runner. |
| 08 | `templates-and-prompts` | The shipped template files (`PRD.prompt.md`, `Design.prompt.md`, etc.) with embedded HTML-comment instructions. The Router Pattern made real. |
| 09 | `polish-and-release` | README, license, `npm publish` setup, manual e2e matrix run with Claude / Cursor / Gemini / Aider, version 0.1.0 cut. |

## Success criteria for v0.1.0

- A new user can `npx savepoint init`, point Claude / Cursor / Gemini / Aider at the project, and reach a working epic-1 audit without intervention.
- All five CLI commands work end-to-end on macOS / Linux / Windows-Terminal.
- Vitest suite passes; manual e2e matrix passes against ≥3 of 4 agents.
- `npm publish` produces a clean install via `npx savepoint init`.

## Risks tracked at this release level

- **Agent compliance variance.** Lighter models may ignore embedded prompts. Mitigation: capability disclaimer + manual e2e matrix.
- **Terminal-rendering edge cases.** Atari-Noir doesn't translate fully; `NO_COLOR` and 16-color terminals exist. Mitigation: explicit fallbacks tested in epic 05.
- **Self-bootstrap correctness.** This repo uses savepoint conventions before savepoint exists. Risk: we drift from our own design while building it. Mitigation: every epic's audit must reconcile with this `Design.md`.
