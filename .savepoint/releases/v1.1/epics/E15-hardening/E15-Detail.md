---
type: epic-design
status: planned
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

## Components

| Module | Purpose |
|--------|---------|
| `internal/board/view_test.go` | Add render benchmarks |
| `internal/data/parser_test.go` | Add fuzz targets |
| `cmd/main.go` | Add --debug / SAVEPOINT_DEBUG |
| `internal/buildtool/main.go` | Add checksums, Windows targets |
| `internal/board/detail.go` | Fix abbreviation handling |
| `internal/board/epic_panel.go` | Extract allowlist constant |
| `agent_skills_test.go` | Move to cmd_test package |

## Boundaries

**In scope:**
- All Phase 3 audit findings not covered by E13 (benchmarks, fuzz, debug flag)
- Medium/Low hardening items (M10, L7, L10, L11, L12)

**Out of scope:**
- CI/CD configuration — separate concern
- New UI features or commands
- Structural improvements (separate epic E14)
