---
id: v2/D003-ci-workflow-missing-token-permissions
release: v2
status: in_progress
stage: build
severity: medium
title: "CI workflow grants GITHUB_TOKEN default permissions instead of least privilege"
---

# D003: CI workflow grants GITHUB_TOKEN default permissions instead of least privilege

## Symptom

`.github/workflows/ci.yml` declares no `permissions` block, so the `ci` job's
`GITHUB_TOKEN` inherits the repository or organization default scope. GitHub
code scanning reports this as open alert #1
(`actions/missing-workflow-permissions`, CWE-275), created 2026-05-12 and still
open.

## Expected Behavior

The workflow declares explicit, least-privilege token permissions. The `ci` job
only checks out code, installs Go, and runs `make ci` (`test build`), so
`contents: read` is sufficient. `publish.yml` already declares its own
permissions (`contents: write`, `id-token: write`), so only `ci.yml` is affected.

## Reproduction

1. Open https://github.com/anipatke/savepoint/security/code-scanning/1 and
   confirm the alert is open against `.github/workflows/ci.yml`.
2. Inspect `.github/workflows/ci.yml`: no `permissions` key appears at workflow
   root or job level; the `ci` job spans lines 9-20.

## Impact

If any action or step in the CI job is compromised, the token can carry write
scope it never needs, leaving the repository open to unauthorized content
changes during otherwise read-only CI runs.

## Fix Plan

Add a workflow-root `permissions: contents: read` block to
`.github/workflows/ci.yml`. Land the change on `master` (directly or through the
current `v2` line) so code scanning closes alert #1 on its next analysis.

## Acceptance Criteria

- [x] `.github/workflows/ci.yml` declares `permissions: contents: read` at workflow root.
- [ ] The fix reaches `master` and code scanning alert #1 closes on the next analysis.
- [ ] Focused and repository quality gates pass.

## Resolution Notes

Fix applied: `.github/workflows/ci.yml` now declares workflow-root
`permissions: contents: read`, parsed and verified. Remains `in_progress`
because the change is not yet on `master` (pushes to `v2` do not trigger
`ci.yml`), so alert #1 cannot close yet, and because repository-wide
`make test` is currently red from an unrelated in-flight E46 provenance
repair (`planned_by`/`executed_session` decoding landed in `internal/data`,
`internal/doctor`, and `internal/migrate`; doctor fixtures and migrate
goldens are not yet updated). D003's workflow-only change cannot affect Go
tests. Resolve once the fix reaches `master` and the concurrent repair
restores the gates.
