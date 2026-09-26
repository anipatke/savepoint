---
id: I-073
title: New Goals still get R-### Release identities
type: drift
status: resolved
source:
  kind: report
  actor: {role: owner, session: user}
  at: '2026-09-26T03:43:13Z'
severity: low
history:
  - at: '2026-09-26T03:43:13Z'
    actor: {role: owner, session: user}
    kind: observed
    note: >-
      An agent creating the v2.1 Goal had to allocate R-007. O-013 renamed
      the public vocabulary to Goals but kept R-### identities out of scope,
      so every new Goal still reads as a Release. Owner direction: fix it
      forward-looking only; existing R-### records are not migrated.
  - at: '2026-09-26T04:04:50Z'
    actor: {role: executor, session: codex-savepoint-task}
    kind: repair_attempted
    note: >-
      Implemented G-### identities for new Goals across V2 parsing, Objective
      references, router reads and writes, Goal-scoped Checks, init scaffolding,
      and migration-generated continuation Goals. V1 Release records still
      migrate to R-###; existing R-### paths and references remain supported.
      `make build`, `make test-fast`, and a subsequent standalone
      `make test-full` passed; the full gate also built Linux, Darwin, and
      Windows targets. `./savepoint resume` loaded this repository's existing
      R-007 Goal. See Repair Attempt Evidence for changed areas, test names,
      and the one unavailable external-project check.
  - at: '2026-09-26T04:41:01Z'
    actor: {role: owner, session: user}
    kind: owner_decision
    note: >-
      Owner instructed the executor to mark this Issue resolved after the
      repair merged to v2 with make build, make test-fast, and make test-full
      passing on the merged branch. No independent Check was run and no
      technical CLEAR is implied.
resolution:
  disposition: accepted
  actor: {role: owner, session: user}
  at: '2026-09-26T04:41:01Z'
  reason: Owner accepted the repair committed to v2 in c9a5e50.
---

# I-073: New Goals still get R-### Release identities

## Summary

O-013 renamed Release to Goal in the board, output, and guidance but
deliberately kept the persisted `R-###` identity (its Out of scope list:
"Renaming persisted `R###` identities to `G###`"). The result is that every
Goal created today, by `savepoint init`, by `savepoint migrate`, or by an
agent, still gets an `R-###` ID, which contradicts the Goal vocabulary.

Owner direction: fix forward only. New Goals get `G-###` identities.
Existing `R-###` Goals keep their IDs, paths, and references unchanged and
still load. No data migration.

## Evidence

- `.savepoint/objectives/O-013-rename-release-to-goals/Objective.md`,
  Out of scope: G### rename and the `releases/` path excluded.
- The identity checks accept only `R`: `internal/data/release_v2.go:81`,
  `internal/data/router_v2.go:80`, `internal/data/project.go:183`,
  `internal/data/write.go:1000`, `internal/data/check_v2.go:229`.
- The V2 scaffold ships `templates/project-v2/.savepoint/releases/R-001-first-goal`.
- Guidance (AGENTS.md "Required Goal Context", `savepoint-design`) tells
  agents that a new Goal uses the first unused `R-###`.

## Proof Needed

- A Goal record with a `G-###` ID loads, and Objective `release:` and router
  `release:` references to it resolve, alongside existing `R-###` Goals.
- `savepoint init` creates `G-001`; a Goal that `savepoint migrate` creates
  gets a `G-###` ID. Migrate keeps existing V1 Releases as `R-###`.
- Guidance and templates tell agents to give new Goals the next `G-###`.
- Existing `R-###` projects, including this repository and galaxy, load and
  behave exactly as before, with no file changes.
- Tests cover mixed `R-###`/`G-###` projects; `make build && make test-fast` pass.

## Repair Attempt Evidence

- V2 Goal identity validation now accepts `R-###` and `G-###`. Objective
  `release:` references, router `release:` selections and writes, and
  release-scoped Checks accept either prefix. Discovery keeps using the
  Release-compatible `.savepoint/releases/<slug>/Release.md` storage path.
- `savepoint init`'s template now creates and selects `G-001`. Migration still
  assigns `R-###` IDs to converted V1 Releases; an automatically generated
  continuation Goal uses the next unused `G-###` identity.
- Active guidance and scaffold copies describe `G-###` for new Goals and
  preserve existing `R-###` records. Canonical and scaffolded skills compare
  byte-identically.
- Regression coverage includes `TestDecodeReleaseV2_acceptsNewGoalIdentity`,
  `TestLoadV2Index_releasesDeriveObjectiveMembership`,
  `TestLoadV2Index_routerReleaseSelectionResolvesExistingRecords`,
  `TestReadStateV2_decodesSelectedGoalWithNewIdentity`,
  `TestWriteRouterStateV2_setsGGoalContext`,
  `TestDecodeCheckV2_acceptsNewGoalScopeIdentity`,
  `TestV2ScaffoldCreatesProjectGoalWithInterpolatedName`,
  `TestPlan_routerSelectsLiveGoalOrPlansContinuationGoal`,
  `TestConvertRelease_activeAndHistoricalForms`, and
  `TestV2DiagnosticName_invalidGoalIdentity`.
- Verification: `make build` passed; the final `make test-fast` passed; the
  standalone `make test-full` passed and built Linux, Darwin, and Windows
  targets. The first fast run exposed a missing closing YAML fence in the new
  router fixture; it was corrected before the passing run. The first chained
  full run reported a root-package setup failure; `go test . -count=1` passed,
  followed by a passing standalone `make test-full`. `./savepoint resume`
  loaded this repository with its existing R-007 Goal. `git diff --check` and
  canonical/scaffold skill parity checks passed.
- Relevant reads covered `.savepoint/Guardrails.md`, the active task and
  issue-capture workflow, V2 identity/reference parsers, migration planning,
  init templates, doctor repairs, board Goal selection, and their regression
  tests. These paths were read to keep identity validation, allocation,
  diagnostics, scaffolding, and shipped instructions aligned.
- Galaxy was not available at a known path in this workspace, so its R-only
  project was not loaded or edited. No Galaxy files were changed.
