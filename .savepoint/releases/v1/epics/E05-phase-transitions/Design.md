---
type: epic-design
status: planned
---

# Epic E05: Phase Transitions

## Purpose

Enable task status and phase mutations from the board. Space advances, backspace retreats. Gates enforce phase order. Writes update task frontmatter and router state.

## What this epic adds

- `Space` advances task through phases: planned → build → test → audit → done.
- `Backspace` retreats: done → audit → test → build → planned.
- Gate enforcement: cannot skip phases; cannot reach done without audit.
- Task frontmatter write-back: update `status` and `phase` fields in the `.md` file.
- Optimistic mtime conflict detection before writes.
- Router write-back: update `router.md` Current state when epic/release changes.

## Definition of Done

- Space/backspace correctly walk the phase chain.
- Gate blocks illegal transitions with a visible message.
- Task `.md` files are updated with new status/phase.
- Mtime conflicts are detected and reported (no silent overwrites).
- Router state updates when switching epic or release.
- All transition tests pass.
- `go test ./...` passes.

## Components and files

| Path | Purpose |
|------|---------|
| `internal/board/transitions.go` | Phase transition logic, gates |
| `internal/data/write.go` | Task frontmatter write-back |
| `internal/data/write_test.go` | Write tests |
| `internal/board/transitions_test.go` | Gate and transition tests |
