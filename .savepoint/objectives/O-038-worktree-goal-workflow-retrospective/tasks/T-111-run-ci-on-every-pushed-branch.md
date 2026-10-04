---
id: T-111
title: Run CI on every pushed branch
objective: O-038
status: planned
depends_on: []
owner_validation: {required: true}
planned_by: {role: planner, session: plan-o038-2026-10-04}
---

# Run CI on every pushed branch

## Outcome

Pushing any working branch runs the full CI workflow, including the native `windows-tests` job, without anyone adding the branch name to the workflow first. Tag pushes and the pull-request trigger behave as before.

## User Check

Push the working branch and open its GitHub Actions run. Confirm that `ci` and `windows-tests` both ran for the push event on the pushed commit, and that the workflow file no longer lists individual push branches.

## Done When

1. In `.github/workflows/ci.yml`, the `push` trigger uses `branches: ['**']` in place of the hard-coded list, so every branch push runs CI and tag pushes do not. The `pull_request` trigger, the jobs and the permissions are unchanged.
2. No other workflow, Go source, template or skill file changes.
3. Evidence records the reviewed trigger diff, a successful `make build && make test-fast` result with time and toolchain, and, once the owner has pushed, the CI run URL, revision and the results of `ci` and `windows-tests`. Before the push, record that run as an outstanding owner action rather than claiming it.
4. Per-criterion evidence follows AGENTS.md's Verification Policy. If the optional Task Check is skipped, record an explicit owner waiver naming this Task, the reason, the actor and the time. The mandatory Full Objective Check still applies.

## Context Files

`.github/workflows/ci.yml`, `.github/workflows/publish.yml`, `.savepoint/objectives/O-038-worktree-goal-workflow-retrospective/Objective.md`.

## Design References

O-038 Confirmed Retrospective Decisions — 2026-10-04, decision 1.

## Guardrails

CFG-03, TEST-01, TEST-08, STYLE-10.

## Implementation Plan

1. Confirm `publish.yml` alone owns tag-triggered work, so excluding tags from the CI push trigger loses nothing. If it does not, return REPLAN REQUIRED.
2. Replace the push branch list with `branches: ['**']`.
3. Run `make build && make test-fast`.
4. Record the evidence, and leave the post-push CI run as an owner action.

## Boundaries

Do not change jobs, runners, Go versions, the `pull_request` trigger, the publish workflow, or any skill, template or Go code. Do not add CI orchestration or required-check settings. Agents do not push; the owner pushes.

## Technical Verification

Use `make build && make test-fast` for handoff. Native platform evidence is the GitHub Actions `windows-tests` result on the pushed working branch, produced by repository CI after the owner pushes and supplied before the Full Objective Check. The Full Objective Check needs fresh `make test-full` under `agent-skills/references/check-method.md`.

## Technical Evidence

Pending execution.

## Drift Notes

None yet.
