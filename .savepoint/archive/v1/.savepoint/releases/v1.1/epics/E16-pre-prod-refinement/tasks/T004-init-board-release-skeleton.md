---
id: E16-pre-prod-refinement/T004-init-board-release-skeleton
status: done
objective: Ensure `savepoint init` creates the release skeleton expected by `savepoint board` and the scaffolded router
depends_on: []
---

# T004: Repair Init/Board Release Skeleton Contract

## Problem

`npx savepoint init` currently scaffolds `.savepoint/config.yml`, `.savepoint/router.md`, `.savepoint/PRD.md`, `.savepoint/Design.md`, and related files, but it does not create `.savepoint/releases`. Running `npx savepoint board` immediately after init fails because board discovery treats the missing releases directory as fatal.

The scaffolded router also instructs agents to read `.savepoint/releases/v1/v1-PRD.md`, but init does not create that file or its parent release directory.

## Context Files

- `internal/init/scaffold.go` — scaffold walk/write behavior and release number interpolation
- `internal/init/scaffold_test.go` — scaffold unit coverage for created files and directories
- `internal/init/integration_test.go` — init pipeline coverage for complete generated projects
- `templates/project/.savepoint/router.md` — scaffolded router contract that references release-scoped PRD
- `templates/project/.savepoint/PRD.md` — root PRD template that may seed or align with release PRD content
- `internal/data/discover.go` — release discovery behavior used by board
- `internal/board/board.go` — board startup path and plain-output behavior
- `internal/board/board_test.go` — board coverage for generated or empty Savepoint projects

## Acceptance Criteria

- [x] Running `savepoint init` creates `.savepoint/releases/v1/epics` in the target project
- [x] Running `savepoint init` creates the release PRD path referenced by the scaffolded router, `.savepoint/releases/v1/v1-PRD.md`
- [x] The release PRD content is generated from existing templates or a single source of truth, without duplicating conflicting project-vision copy in code
- [x] Running `savepoint board` immediately after init no longer fails with `releases directory not found`
- [x] Board behavior for existing projects with populated releases remains unchanged
- [x] Tests cover the init scaffold release skeleton and the board startup path for a freshly initialized, zero-epic project
- [x] `make build && make test` passes

## Implementation Plan

- [x] Reproduce the current contract mismatch in a focused test: init output lacks `.savepoint/releases/v1/epics` and board discovery fails
- [x] Add file-backed scaffold templates or scaffold helper logic to create `.savepoint/releases/v1/epics`
- [x] Add `.savepoint/releases/v1/v1-PRD.md` generation aligned with the scaffolded router
- [x] Add or adjust init integration tests to assert the full release skeleton exists after init
- [x] Add or adjust board tests so a freshly initialized zero-epic project renders or exits cleanly instead of failing
- [x] Keep missing/corrupt release errors meaningful for genuinely invalid existing projects
- [x] Run `make build && make test`

## Context Log

Files read:
- `internal/init/scaffold.go`, `internal/init/scaffold_test.go`, `internal/init/integration_test.go`
- `templates/project/.savepoint/router.md`, `templates/project/.savepoint/PRD.md`
- `internal/data/discover.go`, `internal/board/board.go`, `internal/board/board_test.go`
- `main.go`, `internal/testutil/fixture.go`, `internal/testutil/fs.go`

Notes:
- Root cause: templates had no `releases/` subtree; embed directive excluded dotfiles so `.gitkeep` needed `all:` prefix.
- `discover.ListReleases` kept fatal — error is still meaningful for corrupt existing projects.
- `board.loadBoardData` naturally handles zero-epic release (empty epics slice).
