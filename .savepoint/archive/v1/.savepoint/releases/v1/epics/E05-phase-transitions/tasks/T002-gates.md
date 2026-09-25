---
id: E05-phase-transitions/T002-gates
status: done
objective: "Enforce phase order: cannot skip, must reach audit before done"
depends_on: [E05-phase-transitions/T001-phase-stepping]
---

# T002: Gates

## Acceptance Criteria

- Cannot advance build → audit (must pass test).
- Cannot advance audit → done if dependencies not done.
- Retreat transitions are always allowed.
- Blocked transitions show a message in the status bar.

## Implementation Plan

- [x] Add `CanAdvance(task, allTasks)` to transitions.go.
- [x] Add dependency check.
- [x] Add phase adjacency check.
- [x] Return reason string on block.
- [x] Write tests.
- [x] Run `go test`.

## Context Log

- **Files read:** transitions.go, transitions_test.go, update.go, model.go, view.go, task.go, model_test.go, update_test.go
- **Input tokens (estimate):** ~3,200
- **Quality gates:** `go build ./...` ok, `go vet ./...` ok, `go test ./internal/board/...` 107/107 pass
- **Notes:** No drift — all changes in existing files within `internal/board/`.
