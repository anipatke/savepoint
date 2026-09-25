---
id: E50-release-validation-cutover/T005-prove-six-platform-distribution-and-cli-contract
status: done
objective: Build and verify all supported archives, checksums, launchers, help text, and public V2 command documentation.
depends_on:
    - E50-release-validation-cutover/T003-confine-v1-compatibility-to-migration-and-history
    - E50-release-validation-cutover/T004-activate-v2-scaffolds-and-workflow-assets
complexity_tier: high
complexity_reason: Cross-platform archives, launchers, CI, checksums, and documentation must agree.
---

# T005: Prove six-platform distribution and CLI contract

## Problem

Source-level cutover is not releasable until the npm launcher, archive layout, checksums, workflows, help, and README all describe and exercise the same V2-only product on six supported targets.

## Context Files

- `.savepoint/releases/v2/epics/E50-release-validation-cutover/E50-Detail.md`
- `internal/buildtool/main.go`
- `internal/buildtool/main_test.go`
- `bin/savepoint.js`
- `package.json`
- `Makefile`
- `.github/workflows/ci.yml`
- `.github/workflows/publish.yml`
- `main.go`
- `main_test.go`
- `README.md`

## Acceptance Criteria

- [x] Linux, macOS, and Windows archives are produced for amd64 and arm64 with the expected binary names and executable metadata where applicable.
- [x] Distribution checksums cover exactly the six archives and fail on a missing, extra, renamed, or changed artifact.
- [x] Native launch smoke tests run where the host permits; non-native archives receive deterministic structural verification.
- [x] The npm launcher resolves every supported platform/architecture pair and reports unsupported pairs clearly.
- [x] `--help`, command-specific help, README examples, and CI use only the final V2 command/filter contract and distinguish migration from asset upgrade.
- [x] No workflow publishes, tags, deploys, or generates a changelog during validation.

## Implementation Plan

- [x] Reconcile the build target matrix, archive naming, launcher mapping, and checksum inventory.
- [x] Add negative distribution tests for missing, extra, corrupt, and unsupported artifacts.
- [x] Exercise native archives and structurally verify the remaining targets.
- [x] Update help, README, CI, and publish workflow validation without triggering publication.
- [x] Run package packing and distribution checks in an isolated output directory.
- [x] Run the required full gates and record the six-target evidence.

## Context Log

Read: router, E50 detail, Guardrails, all listed Context Files, and the final
V2 command parsers in `cmd/` for help/filter reconciliation.

Implemented:

- `internal/buildtool` now owns the six-target archive naming contract, exact
  archive/checksum inventory validation, ELF/Mach-O/PE structural checks, and
  host-native archive `--version` smoke launch; `dist` verifies and
  `verify-dist` rechecks an existing distribution.
- `bin/savepoint.js` exposes a six-pair resolver with an explicit unsupported
  matrix; Node tests cover all supported pairs and the unsupported diagnostic.
- Make/CI/package validation covers Windows in `build-all`, distribution
  verification, npm packing, global V2 help, `--objective`, and the
  migration-versus-asset-upgrade distinction. Publish validation is a
  read-only prerequisite; publication/tagging remain outside that job.

Evidence:

- `GOCACHE=/tmp/savepoint-go-cache GOFLAGS=-p=1 make build` passed.
- `GOCACHE=/tmp/savepoint-go-cache GOFLAGS=-p=1 make test` passed for all
  packages (including the long migration integration suite).
- `GOCACHE=/tmp/savepoint-go-cache VERSION=v1.3.1 make dist` passed with
  `distribution verified: 6 archives`; `verify-dist` and an exact six-file,
  six-line checksum assertion also passed.
- `GOCACHE=/tmp/savepoint-go-cache NPM_CONFIG_CACHE=/tmp/savepoint-npm-cache
  VERSION=v1.3.1 make package-check` passed; npm dry-run listed all six
  `dist/npm/<os>-<arch>` binaries and performed no publish.
- `GOCACHE=/tmp/savepoint-go-cache go test ./internal/buildtool`,
  `go test .`, and `node --test bin/savepoint.test.js` passed. Generated
  binaries/distribution output were removed with `make clean`.
