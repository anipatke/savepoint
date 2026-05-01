---
id: E01-go-setup/T003-directory-structure
status: done
objective: "Set up internal package directories"
depends_on: [E01-go-setup/T002-entrypoint]
---

# T003: Directory Structure

## Acceptance Criteria

- `cmd/`, `internal/board/`, `internal/data/`, `internal/styles/` directories exist.
- `main.go` lives at root and imports `internal/board`.
- Build compiles with packages.

## Implementation Plan

- [x] Create `cmd/` directory.
- [x] Create `internal/board/` directory.
- [x] Create `internal/data/` directory.
- [x] Create `internal/styles/` directory.
- [x] Move entrypoint logic into `internal/board/` package.
- [x] Update `main.go` to delegate to `board.Run()`.
- [x] Verify `go build`.
