---
type: epic-design
status: planned
---

# E14: Phase 2 — Structural Improvements

## Purpose

Address the remaining Phase 2 findings from the consolidated codebase audit (Opus 4.6 + GLM 5.1) that were not covered by E13: model decomposition, consumer-side interfaces, and structural Medium/Low items flagged by both audits.

## What this epic adds

- Model fields grouped into focused sub-structs for maintainability
- Consumer-side interfaces for data-access types enabling test mocking
- Discover-based traversal in CheckOrphans replacing raw os.ReadDir
- Exact heading matching in epic_panel replacing fragile substrings
- Improved shell tokenization with quote/escape support
- Consolidated ColumnType/TaskStatus enumeration removing syncTaskStatus
- Shared testutil package reducing fixture duplication

## Components

| Module | Purpose |
|--------|---------|
| `internal/board/model.go` | Group Model fields into sub-structs |
| `internal/board/board.go` | Remove newProgramModel and loadAllTasks dead code |
| `internal/board/` + `internal/doctor/` | Introduce consumer-side interfaces |
| `internal/data/discover.go` | Add ListRootDirs method |
| `internal/data/task.go` | Unify ColumnType/TaskStatus |
| `internal/board/transitions.go` | Remove syncTaskStatus |
| `internal/board/epic_panel.go` | Replace substring heading matching |
| `internal/doctor/gates.go` | Improve splitCommand tokenizer |
| `internal/doctor/checks.go` | Route CheckOrphans through Discover |
| `internal/testutil/` | New package for shared test fixtures |

## Boundaries

**In scope:**
- All Phase 2 audit findings not covered by E13 (Model grouping, interfaces)
- Medium/Low audit findings that fit the structural theme (M4, M7, M9, L1, L9)
- Dead code removal (M2, L3)

**Out of scope:**
- New UI features or commands
- Audit Phase 3 hardening items (separate epic E15)
