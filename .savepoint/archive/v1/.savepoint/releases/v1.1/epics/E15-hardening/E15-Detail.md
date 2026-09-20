---
type: epic-design
status: audited
---

# E15: Phase 3 — Hardening

## Purpose

Address the remaining Phase 3 findings from the consolidated codebase audit (Opus 4.6 + GLM 5.1) that were not covered by E13: benchmarks, fuzzing, debug instrumentation, and hardening Medium/Low items.

## What this epic adds

- Benchmark tests for render functions
- Fuzz targets for YAML frontmatter parsing
- --debug flag and SAVEPOINT_DEBUG environment variable
- Distribution checksums in buildtool
- Windows build targets in buildtool
- Abbreviation-aware checklist sentence splitting
- Root-level test relocation and audit allowlist consolidation
- Repo-local CI workflow and `ci` Makefile target for build/test verification

## Components

| Module | Purpose |
|--------|---------|
| `internal/board/view_test.go` | Add render and layout benchmarks |
| `internal/board/card_test.go` | Add card render benchmarks |
| `internal/board/column_test.go` | Add column render benchmarks |
| `internal/data/fuzz_test.go` | Add frontmatter and split-body fuzz targets |
| `main.go` | Add global --debug parsing and SAVEPOINT_DEBUG activation |
| `internal/board/debug.go` | Add board debug logging helper |
| `internal/board/board.go` | Log board initialization under debug mode |
| `internal/board/watch.go` | Log watcher and reload activity under debug mode |
| `internal/board/update.go` | Log update dispatch and preserve focused-card visibility |
| `internal/buildtool/main.go` | Add checksums and Windows targets |
| `internal/board/detail.go` | Fix abbreviation handling |
| `internal/board/epic_panel.go` | Document audit hidden-section allowlist |
| `agent_skills_test.go` | Move to external root test package |
| `Makefile` | Add `ci` target for local verification |
| `.github/workflows/ci.yml` | Add repo CI workflow |

## Implemented As

- Debug activation is in root `main.go` because global flags are parsed before dispatching into `cmd/`.
- Fuzz targets live in `internal/data/fuzz_test.go` rather than `internal/data/parser_test.go`.
- Fuzzing found and fixed `SplitFrontmatterBody` delimiter-offset handling and CR line-ending normalization.
- Root `agent_skills_test.go` remains at the repository root but now uses external package `main_test`, which removes the root package coupling without moving fixture-relative paths.
- CI was included as repo-local guardrail scope for E15 via `T008-ci-and-release-automation`.

## Boundaries

**In scope:**
- All Phase 3 audit findings not covered by E13 (benchmarks, fuzz, debug flag)
- Medium/Low hardening items (M10, L7, L10, L11, L12)

**Out of scope:**
- External deployment hosting beyond repo-local CI/release automation
- New UI features or commands
- Structural improvements (separate epic E14)
