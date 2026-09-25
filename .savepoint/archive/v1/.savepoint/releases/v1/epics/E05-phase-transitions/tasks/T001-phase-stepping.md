---
id: E05-phase-transitions/T001-phase-stepping
status: done
objective: "Space advances, Backspace retreats through phases"
depends_on: [E04-board-components/T006-help-overlay]
---

# T001: Phase Stepping

## Acceptance Criteria

- `Space` on planned task moves to in_progress/build.
- `Space` on in_progress advances phase: build → test → audit → done.
- `Backspace` retreats: done → audit → test → build → planned.
- Illegal transitions are blocked (governed by gates).

## Implementation Plan

- [x] Create `internal/board/transitions.go`.
- [x] Implement `Advance(task)` and `Retreat(task)`.
- [x] Wire Space/Backspace into Update loop.
- [x] Write tests.
- [x] Run `go test`.

## Context Log

- Files read: AGENTS.md, internal/data/task.go, internal/board/update.go
- Quality gates: `go test ./...` passed.
- Notes: Implemented Advance and Retreat logic. Mapped space and backspace keys to trigger transitions. Tests written and passing.
