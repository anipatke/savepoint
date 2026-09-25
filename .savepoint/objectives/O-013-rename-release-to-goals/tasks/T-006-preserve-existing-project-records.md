---
id: T-006
title: Preserve existing project records
objective: O-013
planned_by: {role: planner, session: goals-terminology-20260921}
status: done
complexity_tier: medium
complexity_reason: "The public rename must be proven against V2 parsing, membership, router selection, completion, cutover, and migration compatibility without changing the stored contract."
depends_on: []
owner_validation:
    required: true
    accepted_check: ""
check_waiver:
    task: T-006
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-23T08:50:51Z"
---

# T-006: Preserve existing project records

## Outcome

The compatibility boundary for the Goal rename is explicit and tested: current
V2 projects continue to use their existing `R###` Release records and
`release:` references while the public presentation vocabulary changes.

## User Check

Review the compatibility statement and representative existing-project
fixtures. Confirm that this Objective changes the user-facing name only and
does not silently introduce a `G###`/`goal:` schema migration. If a full
persisted rename is desired, stop and replan this Objective before
implementation.

## Done When

- Existing `R###` Release records and `.savepoint/releases/` paths still load
  through the V2 index.
- Objective `release:` references and router `release:` selections still
  resolve to the same Objectives and Tasks.
- The canonical Release completion and project cutover resolvers remain the
  only gate decisions; no Goal-specific duplicate is introduced.
- Existing migration conversion and historical-completion behavior remains
  covered by named tests or is explicitly recorded as unchanged evidence.
- No parser, writer, or migration path rewrites stored Release data as part of
  the terminology change.

## Context Files

`internal/data/release_v2.go`, `internal/data/release_v2_test.go`, `internal/data/release_gate_v2.go`, `internal/data/release_gate_v2_test.go`, `internal/data/release_cutover.go`, `internal/data/objective_v2.go`, `internal/data/objective_v2_test.go`, `internal/data/project.go`, `internal/data/project_test.go`, `internal/data/router_v2.go`, `internal/data/router_v2_test.go`, `internal/migrate/convert_releases.go`, `internal/migrate/convert_releases_test.go`, `internal/migrate/convert_test.go`.

## Design References

Design sections 2, 3, 4, 5, 6, 10, and 13.

## Guardrails

DATA-02, DATA-03, FS-01, TPL-02, TEST-01, TEST-02, TEST-04, TEST-06,
TEST-08, STYLE-07, and STYLE-10.

## Implementation Plan

1. Confirm the existing R###, `release:`, router, membership, completion, and
   migration contracts before changing presentation code.
2. Add only the missing regression assertions needed to prove old project data
   remains loadable and gate behavior is unchanged; reuse existing tests when
   they already provide exact evidence.
3. Record the compatibility boundary for the UI and documentation Tasks.
4. Run focused data and migration tests, `git diff --check`, `make build`, and
   `make test`.

## Boundaries

No `G###` identity vocabulary, `goal:` frontmatter, goals directory, migration
rewrite, new completion resolver, or changes to the Release/Objective gate
policy.

## Technical Verification

Focused `internal/data` and `internal/migrate` tests covering existing V2
records, Objective membership, router selection, completion/cutover decisions,
and legacy conversion; `git diff --check`; `make build`; and `make test`.

## Technical Evidence

Start evidence: the owner reports the project's `NEXT` displays “Planned T-006”.
The `ResolveNext` ready rung selects a planned Task only when
`ResolveTaskStart` allows it. O-013's predecessor O-012 is `done`, and this
Task has no Task dependencies; implementation may start.

Extra reads beyond Context Files, and why:

- `.savepoint/objectives/O-013-rename-release-to-goals/Objective.md` — checked
  Objective boundaries and the predecessor relationship.
- `.savepoint/objectives/O-012-simplify-task-card-outcomes/Objective.md` —
  verified the predecessor Objective is done.
- `internal/data/next.go` — reconciled the owner's visible `NEXT` with the
  resolver's ready-task selection after the stale router action still named
  completed T-020.
- `.savepoint/objectives/O-018-hyphenate-numbered-record-identities/Objective.md`
  and its `T-020-rename-this-projects-records.md` — checked the conflicting
  router selection and confirmed T-020 is already done.
- `.savepoint/Design.md` sections 2–6, 10, and 13, and the named rule rows in
  `.savepoint/Guardrails.md` — checked the Task's explicit architecture and
  policy references.
- `internal/data/next_test.go` — needed to verify runtime router Release
  selection resolution; the listed `router_v2_test.go` covers decoding only.
- `.savepoint/router.md` and `agent-skills/savepoint-task/SKILL.md` — required
  route/workflow context; during handoff the shared-worktree router was
  observed selecting O-013/T-006 and was left untouched.

Acceptance evidence:

- Existing Release records and paths load; Objective membership is derived
  from stored `release:` values: `TestDecodeReleaseV2_valid`,
  `TestLoadV2Index_releasesDeriveObjectiveMembership`, and
  `TestLoadV2Index_releaseIdentitySurvivesDirectorySlugChange`.
- Parsed router Release/Objective/Task selections resolve to the same existing
  records, and a cross-Release selection is refused:
  `TestLoadV2Index_routerReleaseSelectionResolvesExistingRecords` (added),
  plus the lower-level failure case
  `TestResolveSelection_releaseMismatchNeverSubstitutesObjective`.
- The canonical Release completion and cutover decisions remain unchanged;
  named coverage includes
  `TestResolveReleaseCompletion_requiresMemberObjectivesAndCurrentAcceptance`,
  `TestResolveReleaseCutoverComposesCanonicalReleaseDecisions`, and
  `TestResolveReleaseCutover_refusesTechnicalIssueAndOwnerStates`. No
  Goal-specific resolver was introduced.
- Migration and historical completion behavior remains covered by
  `TestConvertRelease_activeAndHistoricalForms`,
  `TestConvertRelease_historicalArchiveHashMatchesSource`, and
  `TestPlan_releasePRDIdentityAndObjectiveReferenceAreStable`.
- Parser/writer preservation is covered by `TestDecodeReleaseV2_valid` and
  `TestWriteReleaseEvidenceV2_preservesSourceAndRefusesStaleDocument`; the
  historical converter tests above verify archive hashing and typed historic
  completion. No parser, writer, migration, or gate production code changed.

Files inspected: relevant sections of every file listed under Context Files,
plus the extra reads above. Files changed by T-006: this Task record and
`internal/data/release_v2_test.go`. The shared-worktree `.savepoint/router.md`
change selecting O-013/T-006 was preserved, not edited by this Task.

Verification: PASS `go test ./internal/data -run
'^TestLoadV2Index_routerReleaseSelectionResolvesExistingRecords$' -count=1`;
PASS `go test ./internal/data ./internal/migrate -count=1` (data package
0.595s, migration package 110.749s); PASS `git diff --check`; PASS
`make build`; PASS fresh `make test-full` on 2026-09-23 (Go 1.26.2,
linux/amd64), including all Go packages and Linux, Darwin, and Windows builds.
The full gate's slowest package was `internal/migrate` at 2m2.876s. The first
focused-test attempt was blocked by the sandbox's read-only Go build cache;
the same command passed after approval to run outside the sandbox.

Limitations: no interactive board or CLI was run; this Task verifies the
existing data/parser/gate/migration contract and adds only temporary-project
regression coverage. No optional Task Check or owner waiver has been recorded.

## Drift Notes

Any requirement to rename persisted identities, paths, frontmatter, or router
keys is a material scope change and routes back to design before implementation.
