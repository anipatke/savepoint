---
id: E14-structural-improvements/T002-data-interfaces
status: done
objective: Define consumer-side interfaces for stateless data types and accept them in board/doctor constructors
depends_on: []
---

# T002: Extract Consumer-Side Interfaces for Data-Access Types

## Context Files

- `internal/data/discover.go` — NewDiscover() returns *Discover
- `internal/data/parser.go` — NewParser() returns *Parser
- `internal/data/config.go` — NewConfigReader() returns *ConfigReader
- `internal/data/router.go` — NewRouterReader() returns *RouterReader
- `internal/board/board.go` — calls data.NewDiscover() directly
- `internal/doctor/checks.go` — calls data.NewDiscover() directly
- `internal/board/model.go` — NewModel() signature

## Acceptance Criteria

- [x] Consumer-side interfaces defined for taskDiscoverer, taskParser, configReader, routerReader
- [x] board.NewModel() accepts interfaces instead of calling constructors directly
- [x] doctor.CheckXxx functions accept interfaces instead of calling constructors directly
- [x] Existing concrete types satisfy the interfaces with no changes to internal/data
- [x] Tests can inject mock implementations
- [x] `go test ./...` passes with no regressions

## Implementation Plan

- [x] Define interfaces in board/types.go or at each consumer site
- [x] Update NewModel signature to accept interfaces with zero-value defaults
- [x] Update CheckOrphans, CheckDependencies, CheckRouter signatures
- [x] Add doc comments documenting the interface contracts
- [x] Run `make build && make test`

## Context Log

- Files read: `internal/data/discover.go`, `internal/data/parser.go`, `internal/data/config.go`, `internal/data/router.go`, `internal/board/board.go`, `internal/doctor/checks.go`, `internal/board/model.go`, related board/doctor tests.
- Files edited: `internal/board/interfaces.go`, `internal/board/model.go`, `internal/board/board.go`, `internal/board/tui.go`, `internal/board/io.go`, `internal/board/update.go`, `internal/board/watch.go`, `internal/board/interfaces_test.go`, `internal/doctor/interfaces.go`, `internal/doctor/checks.go`, `internal/doctor/gates.go`, `internal/doctor/interfaces_test.go`.
- Token estimate: ~18k.
- Quality gates: `go test ./internal/board ./internal/doctor` passed; `make build && make test` passed (`go test ./...` all packages).
