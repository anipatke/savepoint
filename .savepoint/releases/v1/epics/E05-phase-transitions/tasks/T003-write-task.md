---
id: E05-phase-transitions/T003-write-task
status: done
objective: "Write updated status/phase back to task markdown files"
depends_on: [E05-phase-transitions/T002-gates]
---

# T003: Write Task

## Acceptance Criteria

- Task `.md` file is updated with new `status` and `phase`.
- Mtime is checked before write; conflict reported if changed.
- Phase field removed when status is planned or done.
- Backup not required (git handles history).

## Implementation Plan

- [x] Create `internal/data/write.go`.
- [x] Implement `WriteTaskStatus(path, task, expectedMtime)`.
- [x] Read file, parse frontmatter, update fields.
- [x] Check mtime, write back.
- [x] Write tests.
- [x] Run `go test`.

## Context Log

- **Files read:** internal/data/task.go, internal/data/parser.go, internal/data/errors.go, internal/data/write.go, internal/data/write_test.go, internal/board/transitions.go, internal/board/update.go, internal/board/model.go, internal/data/router.go, internal/data/config.go
- **Estimated input tokens:** ~4,500
- **Quality gates:** `go build ./...` ok, `go vet ./...` ok, `go test ./...` — 7/7 new tests pass, all existing tests pass
- **Notes:** Drift: `internal/data/write.go` and `internal/data/write_test.go` added, mapped under E02-data-readers umbrella in Codebase Map but belong to E05-phase-transitions — next audit should reconcile.
