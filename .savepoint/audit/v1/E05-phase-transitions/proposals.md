# Audit Proposals: E05-phase-transitions

## Target File

`.savepoint/Design.md`

## Insert After

```md
- **Board command** (`savepoint board`) reads project, non-TTY fallback, Ink TUI, transition gates, mtime writes, audit signaling (epic E06).
```

## With

```md
- **Phase transitions** are owned by the Go board/data boundary: `internal/board/transitions.go` defines planned/in_progress/done phase movement and gate checks; `internal/data/write.go` owns mtime-checked task and router frontmatter writes.
```

## Target File

`AGENTS.md`

## Insert After

```md
| `internal/data/errors.go`                     | [E02-data-readers](.savepoint/releases/v1/epics/E02-data-readers/E02-Detail.md)                   | Shared data-reader boundary error sentinels                                                          |
```

## With

```md
| `internal/data/write.go`                      | [E05-phase-transitions](.savepoint/releases/v1/epics/E05-phase-transitions/E05-Detail.md)         | Mtime-checked task status/phase and router current-state write helpers                               |
| `internal/board/transitions.go`               | [E05-phase-transitions](.savepoint/releases/v1/epics/E05-phase-transitions/E05-Detail.md)         | Board phase advance/retreat lifecycle and dependency gate checks                                     |
```

## Target File

`.savepoint/releases/v1/epics/E05-phase-transitions/E05-Detail.md`

## Insert After

```md
## Definition of Done
```

## With

```md
## Implemented As

- `internal/board/transitions.go` implements in-memory phase movement and `CanAdvance` gate checks.
- `internal/board/update.go` maps Space and Backspace to transition helpers and displays gate-block reasons in the status bar.
- `internal/data/write.go` adds mtime-checked write helpers for task frontmatter and router current state.
- `internal/board/model.go` carries `Root` and writes selected release/epic changes back to `router.md`.

## Audit Delta

- The task write helper exists, but board Space/Backspace does not yet call it, so task markdown write-back is not implemented on the live interaction path.
- Router write-back is implemented only for release/epic overlay selection, preserving the existing `next_action` value.
```

## Quality Review

## Must Fix Before Close

- `internal/board/update.go:47` and `internal/board/update.go:64`: Space/Backspace mutate `m.AllTasks` and refresh the grouped view, but never call `data.WriteTaskStatus`. The epic Definition of Done says task `.md` files are updated and mtime conflicts are detected/reported. The current live board path is memory-only, so the central write-back acceptance criteria from T003 are not satisfied outside isolated helper tests.
- `internal/board/update.go:47`: the update path has no way to know a selected task file path or expected task mtime. `data.Task` carries task metadata but not source path/mtime, and the board model does not keep a side table for them. That means mtime conflict detection cannot be enforced for task status transitions as designed.

## Must Fix Before Next Epic

- Add reducer/integration tests for Space and Backspace that verify the board either writes through `WriteTaskStatus` successfully or reports `ErrMtimeConflict` without mutating local state. Current tests cover `Advance`, `Retreat`, and `WriteTaskStatus` independently, but not the user-facing transition workflow.
- Harden `WriteTaskStatus` body preservation. `internal/data/write.go:54` computes the body offset from the trimmed frontmatter returned by `extractFrontmatter`; frontmatter with extra blank lines or trailing whitespace can shift the splice point. Prefer deriving the body boundary from delimiter indexes in the original normalized content.

## Carry Forward

- Keep `WriteRouterState` as the data-layer boundary for router YAML updates, but document that callers must preserve `next_action` by reading existing router state first.
- Consider splitting `internal/data/write.go` if task writes and router writes grow separate validation rules; for now it is still small enough to remain a single file.

## Already Fixed

- Phase order helpers cover planned -> in_progress/build -> test -> audit -> done and the reverse path.
- Dependency gates block audit -> done when dependencies are missing or not done.
- Router release/epic selection writes through the data-layer mtime helper and reports errors via `StatusMessage`.
