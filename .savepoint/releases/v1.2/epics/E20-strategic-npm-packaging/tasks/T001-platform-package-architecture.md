---
id: E20-strategic-npm-packaging/T001-platform-package-architecture
status: done
objective: Define and generate the root npm package plus platform-specific optional package layout
depends_on: []
complexity_tier: high
complexity_reason: "Redesigns npm package boundaries, launcher resolution, and generated artifact layout"
---

# T001: Platform Package Architecture

## Problem

The tactical npm package bundles all platform binaries in one root package. A strategic fix needs a root CLI wrapper package and platform-specific optional packages so installs fetch only the native binary needed for the current system.

## Context Files

- `package.json`
- `bin/savepoint.js`
- `internal/buildtool/main.go`
- `internal/buildtool/main_test.go`
- `.github/workflows/publish.yml`

## Acceptance Criteria

- [x] Root `savepoint` package keeps the `savepoint` bin entrypoint as a JS launcher
- [x] Root package declares platform packages as optional dependencies using npm-supported manifest fields
- [x] Platform package manifests are generated or maintained with exact OS/CPU targeting metadata
- [x] Platform package artifact names map unambiguously to Savepoint build targets
- [x] Launcher resolution prefers the installed platform package and reports a clear unsupported-platform error
- [x] Tests cover supported and unsupported platform/architecture resolution paths

## Implementation Plan

- [x] Choose the package naming convention for platform packages
- [x] Update or generate package manifests for each supported platform target
- [x] Refactor the JS launcher so platform resolution is testable
- [x] Wire root optional dependencies to the platform package set
- [x] Add focused tests for launcher resolution behavior
- [x] Update buildtool tests for the strategic npm artifact layout

## Context Log

- Naming convention: unscoped `savepoint-{goos}-{goarch}` (six packages: linux/darwin/windows × amd64/arm64). Maps 1:1 to existing build target slugs and to `dist/npm/{slug}/` directory layout.
- JS launcher split into `lib/resolve-platform.js` (pure platform→package mapping, no I/O) and `lib/locate-binary.js` (injectable `resolver` + `existsSync` for tests). `bin/savepoint.js` is now a thin shim that wires `process.platform`/`process.arch` into the resolver, then spawns.
- Launcher resolution order: `require.resolve('savepoint-{os}-{arch}/bin/savepoint(.exe)')` first (npm-installed optionalDependency), then dev fallback at `dist/npm/{slug}/bin/savepoint(.exe)`. Unsupported platform → exact message `savepoint does not support {platform}/{arch}`; missing package → message naming the expected platform package.
- `internal/buildtool/npm.go` owns manifest generation. `buildNPMManifest` returns `{name, version, os:[npm-os], cpu:[npm-cpu], files:["bin"]}` with `linux→linux`, `darwin→darwin`, `windows→win32`, `amd64→x64`, `arm64→arm64`. Version is read from root `package.json` so all platform packages stay in lockstep with the root release.
- `buildNPM` now writes `dist/npm/{slug}/package.json` and `dist/npm/{slug}/bin/savepoint(.exe)` per target, preserving the existing Windows PE header check.
- Root `package.json`: dropped `dist/npm` from `files` (binaries live in platform packages), added `lib/`, added `optionalDependencies` for all six platform packages pinned to the root version, and switched `npm test` to `node --test test/resolve-platform.test.js`.
- Tests: 13 node test cases cover all six supported `(platform, arch)` pairs, two unsupported pairs, `locateBinary` happy-path/dev-fallback/missing-binary branches, and the Windows `.exe` path; new Go tests cover manifest fields per target, unknown OS/arch errors, manifest serialization, and root-version read.
- Verified end-to-end: `go run ./internal/buildtool build-npm` produces six platform package directories each containing the expected manifest and binary. `go test ./internal/buildtool/` and `node --test` pass. Pre-existing `TestBundledSavepointSkillsHaveDiscoveryFrontmatter` failure in the root package is unrelated to E20.
