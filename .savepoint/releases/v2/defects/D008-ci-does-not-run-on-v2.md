---
id: v2/D008-ci-does-not-run-on-v2
release: v2
status: open
severity: high
title: "CI workflow does not trigger on the active v2 branch"
---

# D008: CI workflow does not trigger on the active v2 branch

## Symptom

`.github/workflows/ci.yml` triggers `push` and `pull_request` only for
`master`, while the active development branch is `v2`. Direct V2 pushes and
pull requests can therefore bypass the repository's automatic `make ci`
quality gate.

## Expected Behavior

While V2 is under active development, the CI workflow must run for direct
pushes to `v2` and for pull requests targeting `v2`, in addition to the
repository's existing protected branch behavior.

## Reproduction

1. Inspect the `push.branches` and `pull_request.branches` entries in
   `.github/workflows/ci.yml`.
2. Push a V2-only change or open a pull request targeting `v2`.
3. Observe that this workflow has no matching trigger because both filters
   contain only `master`.

## Impact

V2 changes can land or accumulate without the build and test gate running on
the branch where the rewrite is being developed.

## Fix Plan

Confirmed approach: add `v2` to the CI push and pull-request branch filters,
then verify the workflow triggers and passes on the active branch without
weakening the existing permissions or `make ci` steps. Sequencing: the
repository's `make test` is currently red from the in-flight E46 provenance
repair, so land that repair before or with this change; otherwise `v2` CI
arrives red.

## Acceptance Criteria

- [ ] A push to `v2` triggers the CI workflow.
- [ ] A pull request targeting `v2` triggers the CI workflow.
- [ ] Existing `master` triggers and least-privilege permissions remain intact.

## Resolution Notes

Pending.
