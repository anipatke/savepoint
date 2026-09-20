---
id: E01-go-setup/T004-makefile
status: done
objective: "Create Makefile for build, test, run"
depends_on: [E01-go-setup/T003-directory-structure]
---

# T004: Makefile

## Acceptance Criteria

- `Makefile` exists with `build`, `test`, `run`, `clean` targets.
- `make build` produces a binary.
- `make test` runs `go test ./...`.

## Implementation Plan

- [x] Create `Makefile`.
- [x] Add `build` target: `go build -o savepoint main.go`.
- [x] Add `test` target: `go test ./...`.
- [x] Add `run` target: `go run main.go`.
- [x] Add `clean` target: `rm -f savepoint`.
- [x] Verify all targets work.
