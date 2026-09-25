---
id: E05-phase-transitions/T004-write-router
status: done
objective: "Write router.md Current state when epic or release changes"
depends_on: [E05-phase-transitions/T003-write-task]
---

# T004: Write Router

## Acceptance Criteria

- Selecting a new epic updates `router.md` epic field.
- Selecting a new release updates `router.md` release field.
- Mtime check before write.
- Next action string preserved.

## Implementation Plan

- [x] Add `WriteRouterState(root, state)` to `internal/data/write.go`.
- [x] Read router.md, replace yaml block.
- [x] Check mtime, write back.
- [x] Wire into epic/release selection in Update loop.
- [x] Write tests.
- [x] Run `go test`.

## Context Log

- **Files read:** internal/data/router.go, internal/data/write.go, internal/data/write_test.go, internal/data/parser.go, internal/data/task.go, internal/data/discover.go, internal/data/config.go, internal/data/errors.go, internal/board/model.go, internal/board/update.go, internal/board/board.go, internal/board/update_test.go, internal/board/model_test.go, internal/board/transitions.go, internal/board/transitions_test.go
- **Estimated input tokens:** ~7,200
- **Quality gates:** `go build ./...` ok, `go vet ./...` ok, `go test ./...` — 17/17 new tests pass (4 WriteRouterState tests + 13 existing), all packages pass, board tests pass, data tests pass
- **Notes:** Drift: `Model.Root` field added to board, not yet in Codebase Map. Next audit should reconcile.
