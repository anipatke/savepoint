---
type: epic-validation
status: recorded
epic: E50
recorded_at: 2026-09-20
---

# E50 Validation Evidence

This record covers the two disposable migration trials and three bounded agent
scenarios required before the E50 cutover handoff. It is an implementation
validation record, not the independent E50 audit. The evidence distinguishes
what a test observed from what the observations support; it makes no claim
about universal agent reliability.

## Evidence discipline

- Every migration trial is built by the test's `t.TempDir()` copy helper. The
  frozen fixture and the working repository are inputs only; neither is a
  mutable test fixture.
- Per-file SHA-256 and modification-time snapshots are captured before and
  after preview, apply, and retry. The migration tests compare both values and
  fail on a removed, created, changed, or retimestamped source file.
- `session` and `model` are recorded as `not supplied` for the agent runs in
  this record. Session labels embedded in V2 fixture records (for example
  `build-001`, `sess-1`, `sess-2`, and `owner-1`) are fixture provenance, not
  claims that a particular model ran the scenario.
- The repository workflow prohibits invoking the Savepoint CLI from an agent
  task. The commands below call the exported Go command path and the V2
  consumers used by the packaged commands; this is stated explicitly instead
  of presenting a shell CLI transcript that was not run.

## Input snapshots

The following hashes and metadata were captured before this validation file was
created. The aggregate hashes are an input anchor for the disposable copies;
the tests' per-file snapshots remain the integrity proof.

| trial input | source shape | source hash evidence | source mtime evidence |
| --- | --- | --- | --- |
| `v1-basic/project` (tiny trial) | 9 files, 64 KiB | tree aggregate `a75df68354f13b14c124a250ecd4ead747491a9e6f2673d7b701ec34f8772144`; manifest `66af04449769255e3376404a8e4edcbaed40a39efdda137aea664f5bd8b66dba` | `manifest.yml` mtime `2026-09-14 19:34:20.538601718 +1000` |
| repository working tree (existing-codebase trial) | 996 files at capture, `.git` excluded by the copy helper | source commit `ae56fcbf9cc3eddb70f6d7e4675e174f6858f95a`; tree aggregate `79a9f8e61f54b2d353768f110fce9d0b0b98d9a770620e0db2345cf501d3752c`; the repository fixture's release manifest is `967be14aead81deb3a60979b258252f3331f492a7eb0cd211d05a946d675c8e3` | `v1-history/manifest.yml` mtime `2026-09-14 20:09:24.998937001 +1000` (the migration fixture hash used by the release-accountability checks) |

Commands used for those anchors:

```text
sha256sum internal/data/testdata/migration/v1-basic/manifest.yml
sha256sum internal/data/testdata/migration/v1-history/manifest.yml
find internal/data/testdata/migration/v1-basic/project -type f -print0 | sort -z | xargs -0 sha256sum | sha256sum
find . -path './.git' -prune -o -type f -print0 | sort -z | xargs -0 sha256sum | sha256sum
git rev-parse HEAD
stat -c '%n %Y %y' internal/data/testdata/migration/v1-basic/manifest.yml internal/data/testdata/migration/v1-history/manifest.yml
```

## Trial TR-01 — tiny project

**Scope.** `v1-basic/project` is copied to a disposable directory by
`copyFixture`; the source contains nine files and is the smallest frozen V1
input. The trial exercises preview, apply, recovery/conflict probes, retry,
doctor, board/plain, resume, and archive/reference inspection through the
same isolated-copy conventions.

**Migration command and observed result.**

```text
go test ./internal/migrate -run 'TestEndToEnd_(bothFixturesMigrateThroughTheCommandPath|frozenFixturesAreNeverMutated|temporaryRepositoryCopyMigratesWithReleaseAccountability|everyReferenceResolves|archivedFilesMatchTheFixtureManifestHashes|secondFullRunChangesNothing|previewWritesNothing)|TestApply_(sourceChangedSincePreview_conflictNamesPath|resumesAfterInterruption_convergesToSameFinalState|resumeAfterUserEditedInstalledFile_reportsConflict)|TestRunCommand_(previewAndRecoveryReportAreReadOnlyOnReadable0555Project|recoverApplyAcrossInvocationsUsesRecordedOperation|recoverApplyAcrossInvocationsRejectsEditedSource)' -count=1
ok  github.com/opencode/savepoint/internal/migrate  109.913s
```

The named tests provide these observations for the tiny copy (and separate
fresh temporary copies for each negative-path probe):

| workflow rung | observed evidence |
| --- | --- |
| preview | `TestEndToEnd_previewWritesNothing` runs the default, `--dry-run`, and `--apply --dry-run` forms and compares the full snapshot; `TestRunCommand_previewAndRecoveryReportAreReadOnlyOnReadable0555Project` also proves a readable project is not probed as writable. |
| apply | `TestEndToEnd_bothFixturesMigrateThroughTheCommandPath` reaches schema V2 through `cmd.RunMigrate` and `migrate.RunCommand`; the fixture copy, not the source, is changed. |
| interruption/recovery | `TestApply_resumesAfterInterruption_convergesToSameFinalState`, `TestApply_interruptedBeforeActivation_stillLoadsAsV1`, `TestApply_interruptedAfterActivation_loadsAsV2`, and `TestApply_recoversAtEveryPublishBoundaryWithoutOverwritingUserEdits` cover the durable journal and every publish boundary. |
| source/destination conflict | `TestApply_sourceChangedSincePreview_conflictNamesPath`, `TestApply_resumeAfterUserEditedInstalledFile_reportsConflict`, and `TestRunCommand_recoverApplyAcrossInvocationsRejectsEditedSource` refuse edited bytes and preserve the edit. |
| unchanged retry | `TestEndToEnd_secondFullRunChangesNothing` and `TestApply_secondFullRun_isNoOp` compare bytes, mtimes, IDs, and the migration manifest after a second apply. |
| archive/reference inspection | `TestEndToEnd_archivedFilesMatchTheFixtureManifestHashes` compares every archive byte to the fixture manifest; `TestEndToEnd_everyReferenceResolves` and `TestEndToEnd_everyFixtureFileHasAnAccountableDestination` check live/archive destinations and links. |

V2 consumer probes were run over their own disposable V2 projects, as the
consumer packages require, and are part of this trial's surface checklist:

```text
go test ./internal/data -run 'TestE43_EpicScenario|TestE44_EpicScenario|TestResolveTaskStart_blockedByReplan|TestResolveTaskAdvance_replanBlocksAuditToCompletion|TestResolveNext_replanOutranksExecution|TestResolveNext_replanOutranksCheckNeeded|TestGateDecisions_unaffectedByOpenIssuesOfEveryType|TestResolveTaskCompletion_needsWorkCheckReferencingIssueBlocksThroughClearanceAlone' -count=1
ok  github.com/opencode/savepoint/internal/data  0.085s

go test ./internal/doctor -run 'TestCheckProject_v2ValidNoProblems|TestCheckProject_OpenAdvisoryIssueDoesNotFailLoad|TestCheckProject_IssueVerifiedProofSuperseded|TestCheckProject_ReadOnly|TestCheckProject_ConsistencyReadOnly|TestCheckReleaseReadiness_reportsMissingUnknownStaleNeedsWorkAndAcceptance|TestDiagnosticReport_FullRunWritesNothing' -count=1
ok  github.com/opencode/savepoint/internal/doctor  0.049s

go test ./internal/board/v2 -run 'TestRunWithoutTTYLeadsWithNextAndReportsCounts|TestRunWithoutTTYIncludesTitlesBadgesAndIssueSummary|TestRunWithoutTTYRefusedIndexFailsAndWritesNoBoard|TestRunWithoutTTYIsDeterministic|TestTaskDetailRendersAReplanByItsRecordedReason|TestDetailListsTheIssuesLinkedToTheRecord|TestIssuesNavigationAndFilteringAreReadOnly' -count=1
ok  github.com/opencode/savepoint/internal/board/v2  0.041s

go test . -run 'TestResumeMatrix_everyRungReachedExactlyOnce|TestResumeMatrix_runResumeThroughTheRealCommandAlsoWritesNothingTwice' -count=1
ok  github.com/opencode/savepoint  0.158s
```

**Integrity observation.** `TestEndToEnd_frozenFixturesAreNeverMutated`
records a before/after map of every source file's SHA-256 and `os.Stat`
modification time around both preview and apply. It reports no change. The
trial therefore observes no live-project mutation; it does not infer that an
arbitrary caller's filesystem is immutable.

## Trial TR-02 — existing-codebase copy

**Scope.** `TestEndToEnd_temporaryRepositoryCopyMigratesWithReleaseAccountability`
copies the current working tree to `t.TempDir()`, excluding only `.git`. It
previews, applies, loads the V2 index, checks Release identities and archive
entries, then applies again. The source working tree is never the migration
root.

**Observed result.** The migration portion of the command above passed. Its
explicit assertions show that preview leaves the repository-copy snapshot
unchanged, the apply has a first-class Release and matching live/archive
manifest entries, and the second apply leaves bytes, mtimes, IDs, and the
manifest unchanged. The same migration command also ran the tiny fixture
subtests, so both trial inputs share the recovery, conflict, archive, and
reference evidence listed in TR-01.

**Surface result and boundary.** The data, doctor, board/plain, and resume
commands above passed over fresh disposable V2 projects. They verify the
packaged consumers' read-only behavior and shared `data.Next` interpretation;
they are not presented as a claim that one long-lived repository copy was
mutated through every UI surface.

## Agent scenario SCN-01 — planned execution and Check handoff

| field | record |
| --- | --- |
| session/model | session and model not supplied; fixture provenance uses `planning-fixture`, `build-001`, `sess-1`, and `owner-1` |
| scope | Execute the explicit Task context for the E43 fixture, then hand the resulting scope to an independent checker. |
| reads | Router/task context and the V2 design/PRD are the declared reads; the targeted extra read is `internal/data/e2e_v2_test.go`, which defines the isolated file-backed scenario. No unrelated repository scan is part of the record. |
| plan and evidence | `TestE43_EpicScenario` starts two file-backed Tasks. It records a current checker Check for the technical Task, an owner wait for the owner-validated Task, and a superseding rerun that makes the old freshness stale. Unknown frontmatter and authored Markdown remain intact. |
| Check handoff | The fixture writes checker evidence with `sess-1` and owner acceptance with `owner-1`; the executor's `build-001` is a separate executed-session value. The completion resolver, not the executor, decides whether the Task may close. |
| outcome | PASS for the bounded fixture protocol: explicit scope, targeted extra read, recorded evidence, independent checker provenance, and preserved authored content. |
| limitation | This is a fixture-backed scenario record, not a claim that a particular model session completed every possible Task without error. |

Supporting commands: `go test ./internal/data -run TestE43_EpicScenario -count=1` and the `internal/data` gate tests in TR-01.

## Agent scenario SCN-02 — materially invalid plan and replan

| field | record |
| --- | --- |
| session/model | session and model not supplied; planner provenance in the fixtures is `planner-fixture`, `planner-1`, or `sess-2` |
| scope | A Task is already `in_progress` at `stage: build` when its approach no longer fits. The executor must stop, return `REPLAN REQUIRED`, and leave partial work available to the revised plan. |
| detection | `TestResolveTaskStart_blockedByReplan` and `TestResolveTaskAdvance_replanBlocksAuditToCompletion` return the typed replan blocker. `TestRender_replan` renders “Resolve the recorded replan before resuming build, test, or audit.” The board detail test renders the recorded reason and planner provenance. |
| preserved work | `TestWriteTaskEvidenceV2_setsAllSubBlocksPreservesUnknownFieldsAndBody` keeps `in_progress/build`, dependencies, unknown fields, and authored body while recording the replan. |
| revised plan/resume | `TestWriteTaskEvidenceV2_removesReplanKeyRatherThanEmptyValue` clears the flag only after the revised plan; it removes the key, preserves other evidence, and leaves no empty-value sentinel. The gate remains blocking while the flag is present. |
| outcome | PASS for the bounded replan protocol: partial work is retained, the required action is explicit, and resumption is not silently inferred before the replan is cleared. |
| limitation | The invalidity is deliberately seeded in a disposable fixture; no universal claim is made about detecting every real-world planning gap. |

Supporting commands: `go test ./internal/data -run 'TestResolveTaskStart_blockedByReplan|TestResolveTaskAdvance_replanBlocksAuditToCompletion|TestResolveNext_replanOutranksExecution|TestResolveNext_replanOutranksCheckNeeded|TestWriteTaskEvidenceV2_setsAllSubBlocksPreservesUnknownFieldsAndBody|TestWriteTaskEvidenceV2_removesReplanKeyRatherThanEmptyValue' -count=1` (`ok github.com/opencode/savepoint/internal/data 0.024s`) and `go test ./internal/resume -run TestRender_replan -count=1` (`ok github.com/opencode/savepoint/internal/resume 0.002s`).

## Agent scenario SCN-03 — fresh checker, seeded defect, and repair Check

| field | record |
| --- | --- |
| session/model | session and model not supplied; checker fixture sessions are `sess-1` for the observation and `sess-2` for the later recheck |
| scope | Seed a material defect, let a fresh checker record a durable Issue, ensure advisory posture is not a false completion blocker, then verify the repair with a later Check. |
| finding | `TestE44_EpicScenario` creates the durable `Latency regression` defect from an Objective Check, links it to the Task/Check graph, and detects when a later Objective Check supersedes the proof. `TestWriteIssueV2_updatesManagedFieldsPreservesUnknownFieldsAndBody` records a verified resolution with the reason “Recheck confirmed the fix.” and preserves authored fields/body/history. |
| advisory boundary | `TestCheckProject_OpenAdvisoryIssueDoesNotFailLoad` and `TestDiagnosticReport_IssuePostureAdvisoryOnlyOnV1Project` demonstrate that an advisory Issue is reported without being promoted to a completion blocker. V2 material Issue and proof problems remain named diagnostics. |
| later Check | `TestCheckProject_IssueVerifiedProofSuperseded`, `TestInspectIssueConsistency_verifiedProofSuperseded`, and the E44 scenario report the durable Issue again when its proof Check is no longer latest; the repair path therefore requires a fresh proof rather than silently closing. |
| outcome | PASS for the seeded, bounded protocol: a checker records an Issue, advisory noise is not treated as a universal blocker, and later evidence can verify or reopen the Issue through the canonical Check link. |
| limitation | The checker identity is fixture provenance rather than a live second model session; the record proves the data contract and read-only diagnostics, not a benchmark of checker quality. |

Supporting commands:

```text
go test ./internal/data -run 'TestE44_EpicScenario|TestDecodeIssueV2_valid|TestWriteIssueV2_updatesManagedFieldsPreservesUnknownFieldsAndBody|TestInspectIssueConsistency_verifiedProofSuperseded|TestLoadV2Index_issueVerifiedProofSupersededStillSatisfiesAtLoad' -count=1
ok  github.com/opencode/savepoint/internal/data  0.041s

go test ./internal/doctor -run 'TestCheckProject_OpenAdvisoryIssueDoesNotFailLoad|TestCheckProject_IssueVerifiedProofSuperseded|TestCheckProject_IssueResolutionMissingProof|TestDiagnosticReport_IssuePostureAdvisoryOnlyOnV1Project' -count=1
ok  github.com/opencode/savepoint/internal/doctor  0.015s
```

## Limitations and closeout checks

- The trials use exported command/package paths and test fixtures rather than
  a shell invocation of `savepoint`; this is required by the repository agent
  rules and keeps the evidence read-only with respect to the live project.
- Recovery, board/plain, doctor, and resume probes intentionally use separate
  disposable roots. Their passing results are observations about each package's
  contract, not a claim that one trial root was reused across incompatible
  lifecycle phases.
- No session/model identifiers were supplied, no telemetry was collected, and
  no reliability or benchmark claim is made.
- A full `make build && make test` gate is still required at handoff; this
  record's focused commands are the named evidence that changed during T006.
