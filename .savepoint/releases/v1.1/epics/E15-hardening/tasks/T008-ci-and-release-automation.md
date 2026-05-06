---
id: E15-hardening/T008-ci-and-release-automation
status: done
objective: Add repo-local CI coverage for build/test checks and keep the buildtool-driven release flow reproducible
depends_on: []
---

# T008: Add CI and Release Automation Guardrails

## Context Files

- `Makefile` - local entrypoint for build, test, and dist commands
- `internal/buildtool/main.go` - build orchestration and packaging logic
- `.savepoint/releases/v1.1/epics/E15-hardening/E15-Detail.md` - epic scope and component map

## Acceptance Criteria

- [x] `.github/workflows/ci.yml` exists and runs on push and pull_request events
- [x] The CI workflow runs `go test ./...` and `make build` at minimum
- [x] The CI workflow uses the repo-local Makefile targets rather than duplicating buildtool logic inline
- [x] `Makefile` exposes a `ci` target that runs the local verification steps used by CI
- [x] Generated binaries and archives remain excluded from version control via `.gitignore`

## Implementation Plan

- [x] Add a GitHub Actions workflow for CI in `.github/workflows/ci.yml`
- [x] Add a `ci` Makefile target that wraps the local verification commands
- [x] Reuse existing `build`, `test`, and `smoke-test` targets instead of duplicating orchestration
- [x] Confirm artifact ignore rules still cover local binaries and archives
- [x] Run `make build && make test`

## Context Log

Files read:
- Makefile
- internal/buildtool/main.go
- .gitignore
- .savepoint/releases/v1.1/epics/E15-hardening/E15-Detail.md
- main.go
- go.mod

Estimated input tokens:
- 

Notes:
- 
