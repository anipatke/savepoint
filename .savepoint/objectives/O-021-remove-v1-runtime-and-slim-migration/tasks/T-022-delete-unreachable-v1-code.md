---
id: T-022
title: Delete code the shipped binary cannot reach
objective: O-021
planned_by: {role: planner, session: o021-design-20260923}
status: done
complexity_tier: medium
complexity_reason: "Mechanical deletion across four packages, but tests entangle live and dead code, so each test file must be split or removed without losing coverage of anything still reachable."
depends_on: []
owner_validation:
    required: false
    accepted_check: ""
check_waiver:
    task: T-022
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-23T10:16:26Z"
---

# T-022: Delete code the shipped binary cannot reach

## Outcome

Every function that `deadcode` reports as unreachable from the `savepoint`
main package is gone, together with the tests that only exercise it. The
V1 board disappears, leaving `internal/board` as the thin `board` command
dispatcher. User-visible behavior does not change.

## User Check

None needed beyond the gate: `savepoint board`, `resume`, `doctor`, `init`,
`upgrade-assets`, and `migrate --dry-run` behave exactly as before on this
repository and on a V1 fixture copy.

## Done When

- `go run golang.org/x/tools/cmd/deadcode@latest .` reports nothing under
  `internal/board`, `internal/data`, `internal/doctor`, `internal/init`, or
  `cmd`. Anything reported under `internal/migrate` that T-025 will delete
  anyway (`apply.go`, `cutover.go`, `operation.go`, `manifest.go` entries)
  may remain and is listed in evidence.
- `internal/board` contains only `board.go` (dispatch and `Filters`) and
  whatever it genuinely calls. Every V1 board file and V1 board test is
  deleted. `internal/board/v2` reachable behavior is untouched.
- The unreachable V1 functions in `internal/data` (audit register/run
  loaders and validators, V1 lifecycle helpers, V1 write helpers, unused
  parser/discover methods) and the unused V2 writers (`CreateCheckV2`,
  `CreateIssueV2`, `WriteIssueV2`, `WriteReleaseV2`,
  `WriteReleaseLifecycleV2`, `WriteReleaseEvidenceV2`, and their private
  helpers) are deleted. V1 readers that `internal/migrate` still reaches stay.
- Unreachable V1 checks in `internal/doctor` are deleted; `RunV2Checks`
  output on this repository is byte-identical before and after.
- Tests that only covered deleted code are deleted; tests covering both
  keep their live cases. No test is weakened to make deletion pass.
- Evidence records production and test line counts per package before
  and after, and the before/after `deadcode` reports.
- `git diff --check`, `make build && make test-fast` pass.

## Context Files

Deletion set: the `deadcode` report, run fresh at start. Board:
`internal/board/board.go`, `internal/board/debug.go`,
`internal/board/legacy.go`, `internal/board/model.go`,
`internal/board/update.go`, `internal/board/view.go`,
`internal/board/dispatch_test.go`, `internal/board/board_test.go`. Data:
`internal/data/audit.go`, `internal/data/audit_backlinks.go`,
`internal/data/audit_finding.go`, `internal/data/audit_register.go`,
`internal/data/audit_run.go`, `internal/data/audit_validate.go`,
`internal/data/lifecycle.go`, `internal/data/write.go`,
`internal/data/write_test.go`, `internal/data/release_gate_v2.go`,
`internal/data/release_doc.go`, `internal/data/parser.go`,
`internal/data/discover.go`, `internal/data/task.go`. Doctor:
`internal/doctor/checks.go`, `internal/doctor/checks_test.go`,
`internal/doctor/interfaces.go`, `internal/doctor/repairs.go`,
`internal/doctor/report.go`, `internal/doctor/gates.go`. Other:
`internal/init/clipboard.go`, `internal/init/upgrade.go`, `cmd/init.go`,
`internal/board/v2/card.go`, `internal/board/v2/model.go`,
`internal/board/v2/objectives.go`, `internal/board/v2/releases.go`,
`internal/board/v2/update.go`, `internal/board/v2/view.go`. Any other
`internal/board/*.go` file is deleted whole once `board.go` no longer
references it; its `_test.go` sibling goes with it.

## Design References

Design sections 2 (Codebase Map) and 11 (testing).

## Guardrails

ARCH-01, ARCH-04, DATA-01, TEST-01, TEST-06, TEST-08.

## Implementation Plan

1. Record baseline line counts and a fresh `deadcode` report.
2. Delete the V1 board: remove every `internal/board` file except
   `board.go`, then restore only what `board.go` needs. Move any board
   test that exercises the live dispatch into `board_test.go`.
3. Delete unreachable `data`, `doctor`, `init`, `cmd`, and `board/v2`
   functions. After each package, run `go build ./... && go vet ./...` and
   its focused tests; delete test cases that no longer compile because
   their subject is gone.
4. Re-run `deadcode` until clean (or only T-025's migrate entries remain).
5. Run `make build && make test-fast`; record counts and results.

## Boundaries

No behavior change, no refactoring of live code, no renaming of
`internal/board/v2`. Do not touch `internal/migrate` beyond compile fixes,
templates, skills, or documentation (T-023 to T-026).

## Technical Verification

`deadcode` report; focused package tests during iteration;
`make build && make test-fast` at handoff. See
`agent-skills/references/check-method.md`.

## Technical Evidence

Executor session `t022-build-20260923`, 2026-09-23. Toolchain go1.26.2
linux/amd64; `deadcode` and `staticcheck` from `golang.org/x/tools` v0.50.0 /
`honnef.co/go/tools` latest, run via `go run` (auto-switched to go1.26.8).

### Per-criterion outcome

1. **deadcode clean outside migrate — met.** Before (fresh at start, 359
   entries): `cmd` 1, `internal/board` 215, `internal/board/v2` 9,
   `internal/data` 86, `internal/doctor` 41, `internal/init` 2,
   `internal/migrate` 5. After, `go run golang.org/x/tools/cmd/deadcode@latest .`
   reports only the five T-025 migrate entries:
   `apply.go:534 publishWrite`, `cutover.go:172 EvaluateCutover`,
   `manifest.go:188 UnmarshalManifest`, `manifest.go:202 WriteManifestCreateOnly`,
   `operation.go:161 CreateOperation`. None reported under board, data,
   doctor, init, or cmd.
2. **`internal/board` is the dispatcher only — met.** It now holds
   `board.go` (dispatch and `Filters`), `debug.go` (the `debugf`/`SetDebug`
   that `board.go` and `main.go` call; `DebugEnabled` removed as unreachable),
   and their tests `board_test.go` (the former `dispatch_test.go`, with the
   one `writeTask` call inlined as `testutil.WriteTask`) and `debug_test.go`.
   All 26 other V1 board source files and 25 V1 board test files are
   deleted. `internal/board/v2` changed only by removing its 9 unreachable
   helpers; its reachable behavior is unchanged (probe below).
3. **data deletions — met.** All 86 reported `data` functions deleted,
   including `CreateCheckV2`, `CreateIssueV2`, `WriteIssueV2`,
   `WriteIssueHistoryV2`, `WriteReleaseV2`, `WriteReleaseLifecycleV2`,
   `WriteReleaseEvidenceV2`, and their private helpers. Types, constants,
   and variables left with no production reference were then removed:
   `audit.go`, `audit_backlinks.go`, `audit_register.go`, `audit_run.go`,
   `audit_validate.go`, `release_doc.go` (whole files); the finding
   diagnostic types and codes in `audit_finding.go`; the lifecycle
   diagnostic, transition, and legacy complexity-alias vocabulary in
   `lifecycle.go`; the `v2CheckFileOps` seam and `ErrProposalNotFound` in
   `write.go`; `ErrSavepointDirectoryMissing` and
   `ErrV2ReleaseCompletionBlocked`. Kept V1 readers that migrate reaches:
   `Parser.ParseFindingFile`, `ParseRawFindingFile`,
   `NormalizeFindingForLoad`, task/defect/epic parsing, and discovery.
4. **doctor deletions, byte-identical `RunV2Checks` — met.** All 41
   reported doctor functions deleted, plus `interfaces.go`
   (`DoctorDependencies` and its interfaces, unused once the V1 checks were
   gone), `taskDep`, and `auditFrontmatterRepair`. Byte-identity proof: a
   throwaway probe (scratchpad only, never committed) printed
   `doctor.RunV2Checks(root).Format()` followed by `board.Run()` output.
   It was built once from `HEAD` (6aa2784, in a temporary worktree) and once
   from the working tree, and each build ran against this repository's
   current `.savepoint` and a copy of the `v1-basic` fixture.
   `cmp` reported both outputs identical: 218 lines for this repository
   and 68 lines for the V1 fixture copy, where the board refuses with
   `[migration_required]` as before.
5. **Tests — met, none weakened.** Test functions: 2027 → 1362. Deleted
   only tests whose subject was deleted. The following tests reached live
   code through a deleted wrapper, so they were kept and redirected rather
   than deleted:
   - `internal/init`: ~45 upgrade tests called the unreachable
     `upgradeAssetsFromTree` wrapper. It now lives in `upgrade_test.go`
     as a test seam over the live `upgradeProjectAssets`, so every test
     still runs.
   - `internal/doctor`: V2 cases of `CheckProject`, `CheckReleaseReadiness`,
     `IssuePostureReport`, `RunAllChecks(root, "")`, and `RunQualityGates`
     now call the live `RunV2Checks(root).Project/.Releases/.Issues`,
     `RunV2Checks(root)`, and `runQualityGatesV2`. Six tests that then
     failed only because their fixture was a V1 project were deleted
     (`TestCheckProject_v1ProjectNoProblems`,
     `TestDiagnosticReport_CleanProject`, `_FormatAllClean`,
     `_AuditRegisterAbsentStaysClean`, `_AuditRegisterProblemsInPlainOutput`,
     `_IssuePostureAdvisoryOnlyOnV1Project`). `ALL CLEAN` output stays
     covered by `TestDiagnosticReport_AdvisoryIssueBacklogIsStructurallySound`.
     `TestCheckProject_SchemaVersionMalformed` now expects the live path's
     `File` (the config file, not the root).
     `TestCheckProject_SchemaVersionUnsupported` now drives the live
     unsupported branch (no `schema_version` key), because on the live path
     `schema_version: 99` is mislabelled "malformed". That is captured as
     I-037, not fixed here.
   - `internal/data`: `TestE43_EpicScenario` and `TestE44_EpicScenario`
     write their Check and Issue records as files through a new
     `writeScenarioCheck` helper (decoded with `DecodeCheckV2`) instead of
     the deleted writers; every gate-resolver assertion is unchanged.
     `migration_source_test.go` reads raw YAML with the existing
     `rawFrontmatterMap` instead of the deleted `Parser.ParseFrontmatter`.
     `migration_history_test.go` dropped only the assertions on the deleted
     run parser, `DiagnoseFinding`, and `LoadAuditRegisterSet`.
     `TestTaskLifecycleContract_exposesCanonicalValuesAndAliases` dropped
     its `CanonicalTaskStages` and complexity-alias assertions and keeps its
     live assertions.
   - `internal/board/v2`: `TestO900OutcomeSpreadRendersEveryOutcomeAndBlockerOnCards`
     calls the live `groupTaskCardsForRelease(index, "", "O-900")` in place
     of the deleted one-line wrapper.
   - `main_board_next_parity_test.go` and `internal/migrate/end_to_end_test.go`
     (compile fixes only): `doctor.CheckReleaseReadiness(x)` →
     `doctor.RunV2Checks(x).Releases`.
6. **Line counts — recorded.** Production / test lines (`wc -l` per package
   directory, non-recursive):

   | Package | Prod before | Prod after | Test before | Test after |
   | --- | ---: | ---: | ---: | ---: |
   | `.` (main) | 302 | 302 | 2075 | 2075 |
   | `cmd` | 357 | 353 | 806 | 806 |
   | `internal/board` | 4930 | 90 | 7826 | 321 |
   | `internal/board/v2` | 6356 | 6292 | 5869 | 5869 |
   | `internal/data` | 9515 | 7255 | 16750 | 12543 |
   | `internal/doctor` | 2496 | 1374 | 3720 | 1821 |
   | `internal/init` | 1659 | 1638 | 8960 | 8944 |
   | `internal/migrate` | 5934 | 5934 | 7999 | 7999 |
   | others (resume, buildtool, testutil, styles) | 1883 | 1883 | 1248 | 1248 |
   | **Total** | **33432** | **25121** | **55253** | **41626** |

   Net: −8,311 production lines (−24.9%), −13,627 test lines (−24.7%).
7. **Gate — met.** `git diff --check` clean. `make build && make test-fast`
   exit 0 at 2026-09-23T10:11Z, then re-run at 2026-09-23T10:14Z after
   comment- and message-only edits: exit 0, all packages pass. `go vet ./...`
   clean. `staticcheck -checks U1000 ./...` reports only the four findings
   that were already present at HEAD (`data/router.go stateBlockEnd`,
   `init/upgrade_schema_test.go dirDiff`, `styles/styles.go clrSurface`,
   `clrSurfaceDark`), all outside this Task.

### Files read

Budgeted Context Files as listed. Extra reads, logged:

- `AGENTS.md`, `agent-skills/savepoint-task/SKILL.md`,
  `agent-skills/references/issue-capture.md`, the O-021 Objective,
  `.savepoint/Guardrails.md` (named rule rows): the workflow read order.
- `internal/board/debug_test.go`, `internal/board/integration_test.go`
  (function list only): deciding which board tests hold live dispatch.
- `internal/doctor/v2_runtime.go`, `internal/doctor/report_test.go`,
  `gates_test.go`, `repairs_test.go`: finding the live equivalents of the
  deleted doctor entry points.
- `internal/data/e2e_v2_test.go`, `migration_history_test.go`,
  `migration_source_test.go`, `lifecycle_test.go`, `discover_test.go` and
  `project_test.go` (fixture helpers only), `gate_v2_test.go`,
  `release_gate_v2_test.go`, `config.go` (`ReadSchemaVersion`):
  test compile fallout and the I-037 root cause.
- `internal/init/upgrade.go` (wrapper vs. core), `internal/init/*_test.go`
  (compile fallout), `internal/board/v2/card_test.go`, `boundary_test.go`,
  `main_board_next_parity_test.go`, `internal/migrate/end_to_end_test.go`,
  `internal/migrate/live_boundary_test.go`: compile fallout and boundary
  regexes.
- `.savepoint/config.yml`, `internal/data/testdata/migration/*/config.yml`,
  `internal/testutil/fixture.go`: confirming `RunV2Checks` runs no gate
  commands in tests or in the probe.
- `Makefile` and `internal/buildtool`: confirming CI does not enforce
  `go mod tidy`.

### Files changed

- Deleted: 26 V1 board sources, 25 V1 board tests (incl. the old
  `board_test.go`; `dispatch_test.go` renamed to `board_test.go`),
  `internal/data/{audit,audit_backlinks,audit_register,audit_run,audit_validate,release_doc}.go`,
  `internal/data/{audit_backlinks,audit_register,audit_run,audit,audit_validate,release_doc}_test.go`,
  `internal/doctor/interfaces.go`, `internal/doctor/interfaces_test.go`.
- Edited (deletions within the file): `cmd/init.go`; `internal/board/board.go`
  (the `Filters` comment now describes its live role), `debug.go`;
  `internal/board/v2/{card,model,objectives,releases,update,view}.go`;
  `internal/data/{audit_finding,discover,lifecycle,parser,release_gate_v2,task,write}.go`;
  `internal/doctor/{checks,gates,repairs,report,v2_runtime}.go` (comments
  in `repairs.go`, `report.go`, and `v2_runtime.go` renamed to live entry
  points); `internal/init/{clipboard,upgrade}.go`.
- Edited tests: see criterion 5, plus the deletion-only changes in
  `internal/data/{audit_finding,discover,fuzz,gate_v2,lifecycle,parser,release_gate_v2,task,write}_test.go`,
  `internal/doctor/{checks,gates,repairs,report}_test.go`, and
  `internal/init/{clipboard,lifecycle,manifest,migrate_audit_skill,template_freshness,upgrade_failure,upgrade_schema,upgrade}_test.go`.
- New: `.savepoint/issues/I-037-doctor-names-unsupported-schema-version-malformed.md`.
- The Task record itself (lifecycle and this evidence).

### Limitations

- `resume`, `init`, `upgrade-assets`, and `migrate --dry-run` were not run
  as commands (AGENTS.md: agents never run `savepoint`). Their behavior is
  covered only by the passing `make test-fast` suites. Only doctor and board
  output were compared byte for byte. The owner's User Check covers the rest.
- `make test-fast` skips the three slow migrate tests by design. This Task
  touches migrate only in one test call, so `make test-full` was not run;
  the mandatory Full Objective Check requires it.
- `internal/doctor.DiagnosticReport` still has the V1 fields `Structure`,
  `Dependencies`, `AuditState`, `Orphans`, `Defects`, and `AuditRegister`.
  `RunV2Checks` never fills them, but `Format` reads them. Removing them is
  a live-code refactor that could change output, so it is out of scope
  here.
- `go.mod` still lists `github.com/charmbracelet/bubbles`,
  `github.com/sahilm/fuzzy`, and `github.com/atotto/clipboard`, which
  nothing imports any more. `go.mod` was already untidy at HEAD and CI
  does not run tidy. Left for the owner or T-026.
- The boundary tests in `internal/board/v2/boundary_test.go` and
  `internal/migrate/live_boundary_test.go` still forbid
  `data.AuditRegisterSet`. That type no longer exists, so the guard is
  inert but harmless.

### Handoff

`stage: audit` — ready for an optional Task Check or, under an explicit
owner waiver, for the mandatory Full Objective Check. This is not a pass.

## Drift Notes

None expected. If a reported function turns out to be reachable through
reflection, templates, or `go:linkname`, keep it and record why.

Execution: no reflection, template, or `go:linkname` use exists in the
module (grep for `text/template`, `html/template`, `go:linkname`, and
`MethodByName` found none), so no reported function was kept. Beyond the
function list, orphaned types, constants, and variables were also removed
(criterion 3), which the Objective's "nothing else in `data` is V1-only"
success condition covers. One latent live defect surfaced and was captured as
I-037 instead of being fixed.
