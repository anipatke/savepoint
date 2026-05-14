---
id: E20-strategic-npm-packaging/T004-ci-publish-hardening
status: done
objective: Harden GitHub Actions publish flow for multi-package npm releases
depends_on: [T001-platform-package-architecture, T002-binary-format-validation, T003-packed-install-smoke-tests]
complexity_tier: high
complexity_reason: "Changes release automation, publish ordering, and cross-platform verification gates"
---

# T004: CI Publish Hardening

## Problem

The current publish workflow runs tests and publishes one npm package from Ubuntu. Strategic packaging needs CI to build all platform artifacts, run packed smoke tests on supported operating systems, and publish platform packages before the root wrapper package.

## Context Files

- `.github/workflows/publish.yml`
- `package.json`
- `internal/buildtool/main.go`
- `internal/buildtool/main_test.go`
- `README.md`

## Acceptance Criteria

- [x] Publish workflow builds all npm package artifacts before publishing
- [x] Workflow runs packed-install smoke tests on Windows, Linux, and macOS where feasible
- [x] Workflow publishes platform packages before the root `savepoint` package
- [x] Workflow prevents partial root publish when required platform packages fail validation
- [x] README install and upgrade guidance reflects the strategic package shape
- [x] Release verification notes document how to recover or retry a failed publish safely

## Implementation Plan

- [x] Split publish workflow into build/verify and publish stages
- [x] Add OS matrix smoke testing for packed npm artifacts
- [x] Add deterministic package publish ordering
- [x] Keep npm auth scoped to publish steps only
- [x] Update README install and local upgrade guidance if user commands change
- [x] Run the release quality gates and record any residual constraints

## Context Log

- Read: `.savepoint/router.md`, `E20-Detail.md`, `.github/workflows/publish.yml`, `package.json`, `internal/buildtool/main.go`, `internal/buildtool/main_test.go`, `internal/buildtool/npm.go`, `internal/buildtool/pack_smoke.go`, `Makefile`, `README.md`, `bin/savepoint.js`.
- Edited `.github/workflows/publish.yml`: split into `verify` (matrix ubuntu/windows/macos running `make ci` + `make pack-smoke`) and `publish` (needs verify; rebuilds artifacts, asserts each `dist/npm/<os-arch>/package.json` exists, publishes platform packages in deterministic sorted order, then publishes root). `NODE_AUTH_TOKEN` scoped to publish steps only; `verify` never sees the secret.
- Edited `README.md`: rewrote Updating section into "Installing & Updating" describing the root + optional platform package shape, npm's `os`/`cpu` resolution, and recovery on missing platform binary.
- Added `.savepoint/releases/v1.2/epics/E20-strategic-npm-packaging/RELEASE-VERIFICATION.md`: publish order, pre-publish checklist, partial-publish recovery (bump-and-republish preferred over `npm unpublish`), failed-verify triage, auth hygiene, post-publish manual smoke matrix.
- Quality gates: `make build` ok. `make test` fails on pre-existing unrelated `TestBundledSavepointSkillsHaveDiscoveryFrontmatter` (`agent-skills/savepoint-audit/SKILL.md missing YAML frontmatter`) — not introduced by this task, and outside its scope. All `internal/...` packages green.
