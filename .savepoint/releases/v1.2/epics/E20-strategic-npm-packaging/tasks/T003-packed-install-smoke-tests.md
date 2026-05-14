---
id: E20-strategic-npm-packaging/T003-packed-install-smoke-tests
status: done
objective: Add smoke tests that install packed npm artifacts and execute the Savepoint CLI
depends_on: [T001-platform-package-architecture, T002-binary-format-validation]
complexity_tier: high
complexity_reason: "Coordinates npm packing, temporary installs, CLI execution, and CI-friendly cleanup"
---

# T003: Packed Install Smoke Tests

## Problem

Source-tree tests do not prove that `npm pack`, package contents, optional dependencies, and `npx savepoint` work after installation. The release process needs a smoke test that exercises the same artifacts users install.

## Context Files

- `package.json`
- `bin/savepoint.js`
- `internal/buildtool/main.go`
- `internal/buildtool/main_test.go`
- `.github/workflows/publish.yml`

## Acceptance Criteria

- [x] Smoke test packs the root package and relevant platform package artifacts
- [x] Smoke test installs packed artifacts into a temporary project
- [x] Installed CLI runs `savepoint --version` through npm/npx
- [x] Installed CLI runs `savepoint upgrade-assets --dry-run` against a fixture project or safe temporary project
- [x] Temporary files are isolated from the repository working tree
- [x] Smoke test can run locally and in CI without requiring publish credentials

## Implementation Plan

- [x] Add a buildtool or npm script command for packed-install smoke testing
- [x] Create temporary package/install directories during the smoke test
- [x] Install packed tarballs using local file paths
- [x] Execute `savepoint --version` through the installed npm bin
- [x] Add a dry-run `upgrade-assets` smoke check against a controlled project fixture
- [x] Document the smoke command in the publish workflow

## Context Log

- Added `pack-smoke` subcommand in `internal/buildtool/pack_smoke.go`: builds platform packages via `buildNPM`, `npm pack`s root + host platform package into a temp dir, installs both tarballs into a private smoke project, then exercises `savepoint --version`, `savepoint init <fixture>`, and `savepoint upgrade-assets <fixture> --dry-run`. Temp dir cleaned via `defer os.RemoveAll` so nothing leaks into the working tree.
- Wired dispatch in `internal/buildtool/main.go`, added `make pack-smoke` target and `npm run pack-smoke` script, and inserted a `Packed install smoke test` step in `.github/workflows/publish.yml` before `npm publish` so CI exercises the same artifacts users install. No publish credentials required.
- Added `internal/buildtool/pack_smoke_test.go` covering `parseNpmPackFilename` (happy path, whitespace, empty, invalid JSON, empty array, missing filename), `installedBinaryPath` (windows vs unix), `npmExecutable` (host-conditional), `hostTarget`, and the smoke project manifest invariants (private, trailing newline).
- Quality gates: `go build ./...` clean. `go test ./internal/buildtool` passes (all new + existing). `node --test test/resolve-platform.test.js` 13/13. Repo-wide `go test ./...` shows one unrelated pre-existing failure (`TestBundledSavepointSkillsHaveDiscoveryFrontmatter` trips on CRLF in `agent-skills/savepoint-audit/SKILL.md`) — confirmed present on `git stash`, outside this task's scope.
