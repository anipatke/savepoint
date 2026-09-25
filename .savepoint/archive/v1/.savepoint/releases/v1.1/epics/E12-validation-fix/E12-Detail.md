---
type: epic-design
status: audited
---

# E12: Task File Validation & Auto-Fix

## Purpose

Fix task file parsing to provide helpful defaults and clear error messages when validation fails. Currently tasks fail with cryptic errors like `"invalid in_progress phase "": use build, test, or audit"`.

## What this epic adds

- Default `phase: build` when status=in_progress but phase missing
- Default `status: planned` when both status/column missing
- Better error hints with suggested fixes
- Validate on write with helpful messages

## Components

| Module | Purpose |
|--------|---------|
| `internal/data/parser.go` | Default phase to "build" in firstStage() |
| `internal/data/parser.go` | Default status to "planned" in normalizeColumn() |
| `internal/data/lifecycle.go` | Better error messages with hints |
| `internal/data/write.go` | Validate on write |

## Boundaries

**In scope:**
- Default phase for in_progress tasks
- Default status for tasks
- Better error messages
- Write-time validation

**Out of scope:**
- New UI features
- New commands

## Implemented as

- Parser-side defaults live in `internal/data/parser.go`: empty status/column normalizes to `planned`, and parsed `in_progress` tasks without phase/stage default to `build`.
- Lifecycle validation and user-facing hints live in `internal/data/lifecycle.go`.
- Status writes validate lifecycle rules through `internal/data/write.go` and persist the default `phase: build` when callers write `in_progress` with no stage.
- Regression coverage lives in `internal/data/parser_test.go`, `internal/data/lifecycle_test.go`, and `internal/data/write_test.go`.
