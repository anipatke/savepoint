---
id: E50-release-validation-cutover/T005-prove-six-platform-distribution-and-cli-contract
status: planned
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

- [ ] Linux, macOS, and Windows archives are produced for amd64 and arm64 with the expected binary names and executable metadata where applicable.
- [ ] Distribution checksums cover exactly the six archives and fail on a missing, extra, renamed, or changed artifact.
- [ ] Native launch smoke tests run where the host permits; non-native archives receive deterministic structural verification.
- [ ] The npm launcher resolves every supported platform/architecture pair and reports unsupported pairs clearly.
- [ ] `--help`, command-specific help, README examples, and CI use only the final V2 command/filter contract and distinguish migration from asset upgrade.
- [ ] No workflow publishes, tags, deploys, or generates a changelog during validation.

## Implementation Plan

- [ ] Reconcile the build target matrix, archive naming, launcher mapping, and checksum inventory.
- [ ] Add negative distribution tests for missing, extra, corrupt, and unsupported artifacts.
- [ ] Exercise native archives and structurally verify the remaining targets.
- [ ] Update help, README, CI, and publish workflow validation without triggering publication.
- [ ] Run package packing and distribution checks in an isolated output directory.
- [ ] Run the required full gates and record the six-target evidence.

## Context Log

Pending.
