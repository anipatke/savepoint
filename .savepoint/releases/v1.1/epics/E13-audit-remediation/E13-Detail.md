---
type: epic-design
status: audited
---

# E13: Codebase Audit Remediation

## Purpose

Address findings from the consolidated codebase audit (Opus 4.6 + GLM 5.1). The audit identified 2 critical, 8 high, 11 medium, and 13 low issues across code health, architecture, correctness, and maintainability. This epic remediates the actionable items in priority order.

## What this epic adds

- Fix cycle detection bug producing inaccurate error paths
- Remove dead code and stdlib reimplementations
- Centralize duplicated logic (line normalization, frontmatter extraction, shared utilities)
- Fix AtomicWrite cross-device fallback and accent default fill
- Remove committed binaries and add gitignore/linter config
- Decompose `update.go` monolith into focused handlers
- Extract synchronous I/O from Update() into tea.Cmd async pattern
- Add test coverage for untested packages

## Components

| Module | Purpose |
|--------|---------|
| `internal/doctor/checks.go` | Fix cycle detection path reconstruction |
| `internal/doctor/repairs.go` | Replace stdlib reimplementations, improve repair matching |
| `internal/doctor/gates.go` | Add quality gate timeout |
| `internal/doctor/report.go` | Remove dead CheckResult type |
| `internal/data/write.go` | Extract SplitFrontmatterBody, normalizeLineEndings |
| `internal/data/parser.go` | Use shared normalization |
| `internal/data/config.go` | Fix accent default fill |
| `internal/init/write.go` | Fix AtomicWrite cross-device fallback |
| `internal/board/update.go` | Decompose into per-overlay handlers, extract I/O |
| `internal/board/model.go` | Move I/O methods to tea.Cmd |
| `internal/board/card.go` | Consolidate shortID |
| `internal/board/view.go` | Consolidate shortRouterID |
| `internal/board/epic_panel.go` | Extract stripFrontmatter, consolidate epicIndex |
| `internal/board/release.go` | Consolidate releaseIndex |
| `internal/board/detail.go` | Move WrapText/SplitLongWord to utility |
| `internal/board/column.go` | Remove dead taskLabel, move colOverhead |
| `internal/board/layout.go` | Co-locate layout constants |
| `internal/buildtool/main.go` | Replace trimSpace with stdlib |
| `.gitignore` | Add binary exclusions |
| `.golangci.yml` | Add linter configuration |

## Implemented As

- `internal/data` now owns shared line-ending normalization and frontmatter/body splitting for parser/write paths.
- `internal/doctor` now uses stack-based dependency cycle reconstruction, typed repair sentinels, and timed quality-gate command execution.
- `internal/init` uses a copy-based AtomicWrite fallback after failed rename attempts.
- `internal/board` now has shared utility helpers and dispatches board filesystem reads/writes through Bubble Tea command messages; `update.go` was decomposed but did not meet the original 30% line-count reduction target.
- `internal/buildtool` and `internal/styles` gained focused package tests.
- Repository hygiene now ignores build outputs and archives, removes tracked binaries, and includes `.golangci.yml`; lint execution still depends on `golangci-lint` being installed locally.

## Boundaries

**In scope:**
- All Phase 1 (safe cleanup) items from audit
- Critical bug fix (cycle detection)
- Structural improvements (Update decomposition, utility consolidation)
- I/O extraction from Update() (architectural fix)
- Test coverage gaps (buildtool, styles)

**Out of scope:**
- New UI features
- New commands
- Interface extraction for data types (H4) — deferred to future epic
- ColumnType/TaskStatus unification (L1) — low priority, deferred
- Markdown parser for epic_panel (M7) — deferred
- Windows build targets (M10) — separate epic
- CI configuration — separate concern
