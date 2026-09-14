---
id: E41-migration-source-fixtures/T002-history-source-fixture
status: done
objective: Preserve scoped identity collisions and unresolved audit history as interpretable V1 migration input.
depends_on:
    - E41-migration-source-fixtures/T001-basic-source-fixture
complexity_tier: medium
complexity_reason: Extends frozen source evidence to scoped references and audit dispositions without implementing conversion.
---

# T002: Preserve history and reused IDs before migration

## Problem

The simple fixture cannot prove how migration will distinguish reused Task/Defect IDs or retain fixed-but-unverified findings, owner dispositions, duplicate links, and immutable audit runs. A second representative source project must capture these cases before a converter flattens the hierarchy.

## Outcome

The migration source corpus includes a multi-release project with explicit scoped expectations and intact authored audit/defect history. Tests distinguish raw states from loader defaults and show why IDs require source-qualified mapping.

## User Check

Technical evidence is sufficient for this outcome; inspect the preservation and identity table. Current V1 still requires the owner to mark completion.

## Context Files

- `.savepoint/releases/v2/v2-Design.md` — sections 1, 3, 6, 9, and 12 only
- `.savepoint/Guardrails.md` — DATA-01..04, TEST-01..04, TEST-08 and applicable STYLE
- `internal/data/audit.go`
- `internal/data/audit_finding.go`
- `internal/data/audit_run.go`
- `internal/data/audit_test.go`
- `internal/data/audit_finding_test.go`
- `internal/data/audit_run_test.go`
- `internal/data/parser.go`
- `internal/data/dependency.go`
- `templates/project/.savepoint/audit/findings/README.md`
- `templates/project/.savepoint/audit/register.md`
- `templates/project/.savepoint/Health-Check.md` — project procedure shape only
- `.savepoint/releases/v1.2/defects/D008-filename-style-dependency-resolution.md`
- `.savepoint/releases/v1.5/epics/E40-upgrade-safety/E40-Audit.md` — authored proof/history shape only
- `internal/data/migration_source_test.go` — helpers from T001
- `internal/data/testdata/migration/README.md` — extend provenance/coverage
- `internal/data/migration_history_test.go` — new
- `internal/data/testdata/migration/v1-history/manifest.yml` — new
- `internal/data/testdata/migration/v1-history/project/.savepoint/config.yml` — new
- `internal/data/testdata/migration/v1-history/project/.savepoint/router.md` — new
- `internal/data/testdata/migration/v1-history/project/.savepoint/PRD.md` — new
- `internal/data/testdata/migration/v1-history/project/.savepoint/Design.md` — new
- `internal/data/testdata/migration/v1-history/project/.savepoint/Health-Check.md` — new, short custom procedure
- `internal/data/testdata/migration/v1-history/project/.savepoint/releases/v1/v1-PRD.md` — new
- `internal/data/testdata/migration/v1-history/project/.savepoint/releases/v1/epics/E01-example/E01-Detail.md` — new
- `internal/data/testdata/migration/v1-history/project/.savepoint/releases/v1/epics/E01-example/tasks/T001-shared.md` — new
- `internal/data/testdata/migration/v1-history/project/.savepoint/releases/v1/defects/D001-shared.md` — new
- `internal/data/testdata/migration/v1-history/project/.savepoint/releases/v1.1/v1.1-PRD.md` — new
- `internal/data/testdata/migration/v1-history/project/.savepoint/releases/v1.1/epics/E01-example/E01-Detail.md` — new
- `internal/data/testdata/migration/v1-history/project/.savepoint/releases/v1.1/epics/E01-example/tasks/T001-shared.md` — new
- `internal/data/testdata/migration/v1-history/project/.savepoint/releases/v1.1/epics/E01-example/tasks/T002-follow-up.md` — new
- `internal/data/testdata/migration/v1-history/project/.savepoint/releases/v1.1/epics/E01-example/E01-Audit.md` — new
- `internal/data/testdata/migration/v1-history/project/.savepoint/releases/v1.1/defects/D001-shared.md` — new
- `internal/data/testdata/migration/v1-history/project/.savepoint/audit/prompt.md` — new
- `internal/data/testdata/migration/v1-history/project/.savepoint/audit/register.md` — new
- `internal/data/testdata/migration/v1-history/project/.savepoint/audit/findings/F001-awaiting-proof.md` — new
- `internal/data/testdata/migration/v1-history/project/.savepoint/audit/findings/F002-owner-waiver.md` — new
- `internal/data/testdata/migration/v1-history/project/.savepoint/audit/findings/F003-duplicate.md` — new
- `internal/data/testdata/migration/v1-history/project/.savepoint/audit/runs/2026-07-01-example.md` — new

New files are short curated records, not entire copies of unrelated projects. T001's existing fixture helpers bound test setup; do not rediscover the repository.

## Acceptance Criteria

- [x] Two releases reuse the same epic and short Task ID; one Task is done and the other in progress. A dependency in v1.1 resolves within v1.1, never to the done v1 Task. Two D001 records remain distinct by release.
- [x] The fixture includes fixed-but-unverified F001 with unmet proof, explicitly waived F002 with owner reason, F003 duplicate_of F001, an immutable run, an epic audit, and a custom Health-Check procedure. Authored history and unknown metadata remain raw source evidence.
- [x] A source-qualified manifest records exact bytes/hashes, scoped IDs, original disposition, proof requirements, explicit duplicate link, and active-vs-history expectation without claiming V2 conversion or technical clearance.
- [x] `ParseRawFindingFile` and `LoadAuditRegisterSet` retain canonical finding distinctions. A temporary unknown-status variant demonstrates raw preservation and the current loader's open default; diagnostics are observed without rewriting input.
- [x] Named tests verify raw fixed does not mean verified, waiver is an owner disposition rather than proof, duplicate link points to F001, and historical run content/commit is retained.
- [x] Malformed run/frontmatter and unresolved reference variants in temporary copies report the current reader/resolver behavior with named context; frozen project bytes stay unchanged.
- [x] Focused `MigrationSource`/`MigrationHistory` tests and `make build && make test` pass; evidence distinguishes characterization from future migration verification.

## Implementation Plan

- [x] Create the exact v1-history files listed above, adapting the named fixture/template/source shapes. Keep provenance explicit and original fixture timestamps fixed; do not use runtime timestamps.
- [x] Set v1 T001 done, v1.1 T001 in_progress with canonical stage build, and v1.1 T002 planned with a same-epic dependency on T001-shared. Set one D001 resolved and the other open, with release-qualified frontmatter IDs.
- [x] Populate F001 status fixed with nonempty Proof Needed but no verified proof; F002 waived with explicit owner reason; F003 duplicate pointing to F001. Use valid audit run metadata and short meaningful history sections.
- [x] Add a reviewed inventory manifest using T001's schema and helpers. Include source paths in every identity expectation, ensuring duplicate short IDs cannot collapse into a dictionary entry.
- [x] Add `TestMigrationHistoryScopedReferences`, `TestMigrationHistoryDispositions`, `TestMigrationHistoryRawAndNormalized`, and `TestMigrationHistoryFailures`. Use the existing parser/raw parser, audit loader, and dependency resolver; assert exact states/targets and unchanged input bytes.
- [x] Exercise the remaining finding-status vocabulary through a bounded table of temporary copies when needed; do not add ten near-identical frozen files. Preserve unknown metadata in raw YAML assertions rather than claiming typed structs round-trip unknown fields.
- [x] Reuse test-local helpers; avoid changing production loaders or introducing migration mapping code. Record unsupported/ambiguous source behavior for E45 instead of inventing semantic repairs.
- [x] Run focused/full verification and complete Context Log. Hand the completed epic to a fresh V1 audit session after the owner completes its Tasks.

## Boundaries

No Issue conversion, archive creation in the live project, ID reassignment, writer implementation, product CLI execution, or template refactor. Existing findings are examples of input evidence; this Task does not close, waive, or reclassify any live record.

## Technical Verification

Run `go test ./internal/data -run 'MigrationSource|MigrationHistory' -count=1`, then `make build && make test`. Match each criterion to named assertions and state the remaining migration/recovery verification belongs to E45.

## Context Log

**Read:** `.savepoint/router.md`, `AGENTS.md`, `agent-skills/savepoint-build-task/SKILL.md`, `E41-Detail.md`, this task, `.savepoint/Guardrails.md`, `.savepoint/releases/v2/v2-Design.md` sections 1/3/6/9/12, `internal/data/audit.go`, `audit_finding.go`, `audit_run.go`, `audit_register.go`, `audit_test.go`, `audit_finding_test.go`, `audit_run_test.go`, `parser.go`, `dependency.go`, `lifecycle.go`, `task.go`, `defect.go`, `discover.go`, `templates/project/.savepoint/audit/findings/README.md`, `.../register.md`, `.../prompt.md`, `.../runs/README.md`, `templates/project/.savepoint/Health-Check.md`, `.savepoint/releases/v1.2/defects/D008-filename-style-dependency-resolution.md`, `.savepoint/releases/v1.5/epics/E40-upgrade-safety/E40-Audit.md`, `internal/data/migration_source_test.go` (T001's fixture helpers, reused unmodified), `internal/data/testdata/migration/README.md`, and the full `v1-basic/` fixture tree as a style/schema template.

**Created:** `internal/data/testdata/migration/v1-history/manifest.yml` and, under `project/.savepoint/`, `config.yml`, `router.md`, `PRD.md`, `Design.md`, `Health-Check.md` (short custom Quick-only procedure); `releases/v1/{v1-PRD.md, epics/E01-example/E01-Detail.md, epics/E01-example/tasks/T001-shared.md (done), defects/D001-shared.md (resolved)}`; `releases/v1.1/{v1.1-PRD.md, epics/E01-example/E01-Detail.md, epics/E01-example/tasks/T001-shared.md (in_progress, stage build), epics/E01-example/tasks/T002-follow-up.md (planned, depends_on: [T001-shared]), epics/E01-example/E01-Audit.md, defects/D001-shared.md (open)}`; `audit/{prompt.md, register.md, findings/F001-awaiting-proof.md, findings/F002-owner-waiver.md, findings/F003-duplicate.md, runs/2026-07-01-example.md}`. Extended `internal/data/testdata/migration/README.md` with the `v1-history/` provenance row. Added `internal/data/migration_history_test.go` (new).

**Design notes:**

- Both releases' `T001-shared.md` carry the identical frontmatter `id: E01-example/T001-shared`; only the explicit `release` field (and, in tests, the release/epic context discovery would normally supply — mirroring T001) distinguishes them. `E01-example/T002-follow-up` (v1.1) depends on the filename-style short reference `T001-shared`; `ResolveDependency`'s `isShortTaskRef`/`taskShortID` path strips the `-shared` suffix to `T001` and then requires `sameRelease` + `sameEpic`, so a candidate list containing *both* releases' `T001-shared` still lands on v1.1's in-progress record, never v1's done one — `TestMigrationHistoryScopedReferences` asserts `TaskStatus == in_progress` explicitly, not just a nonempty match, to make this discriminating.
- The two `D001-shared` defects use release-qualified frontmatter `id` (`v1/D001-shared`, `v1.1/D001-shared`) with matching `release` fields, so they never collide despite the shared short filename — the D008 defect fix this mirrors was about *task* short-ID resolution, not defect IDs, so defect distinctness here is purely a frontmatter-identity convention, not a resolver behavior.
- `F001` is `status: fixed` with a nonempty `proof_needed` and no `verified_proof` — `TestMigrationHistoryDispositions` asserts `VerifiedProof == ""` explicitly so "raw fixed" is never conflated with "verified". `F002` is `status: waived` with a nonempty `waiver_reason` (an owner disposition, not a technical fix — no `## Proof` content is claimed). `F003` is `status: duplicate` with `duplicate_of: F001`.
- `F001`'s frontmatter carries an unknown top-level field `reviewer_note` with no corresponding `AuditFinding` struct field. `TestMigrationHistoryRawAndNormalized` unmarshals the raw frontmatter into a `map[string]any` (via `SplitFrontmatterBody` + `yaml.Unmarshal`, the same primitives `ParseRawFindingFile` uses) and asserts the field survives there, since the typed model has nowhere to put it and correctly drops it.
- The "temporary unknown-status variant" required by AC4 replaces F001's `status: fixed` with `status: escalated` (not a canonical `FindingStatus`) in a `t.TempDir()` copy only. `ParseRawFindingFile` preserves the literal `"escalated"` string (a `FindingStatus` is just a named string type, so an unrecognized value round-trips through YAML without error); `ParseFindingFile` heals it to `FindingOpen` via `NormalizeFindingForLoad`; `DiagnoseFinding` reports `FindingInvalidStatusCode` from the raw value. No production loader changed.
- `LoadAuditRegisterSet` is exercised directly against the fixture's `project/.savepoint` root (not a temp copy) and asserts all three findings and the one run load distinctly with their canonical statuses intact, plus `Prompt.Available`/`Version` and `Register.Available`/`HasSummary`.
- `manifest.yml` reuses the exact `migrationManifest`/`migrationManifestFile` schema from `migration_source_test.go` unmodified (no production or shared-test-helper code changed). Fields beyond `path`/`sha256`/`role` (`scoped_id`, `raw_status`, `expected_stage_after_parse`, `depends_on`, `expected_classification`) are hand-authored review annotations only — as in `v1-basic`, no test reads them programmatically; `expected_classification` here uses `history`/`active` directly to name the AC's "active-vs-history expectation" in the reviewable document. `scoped_id` for the shared-short-ID entries is release-prefixed (`v1/E01-example/T001-shared` vs `v1.1/E01-example/T001-shared`) specifically so a plain `T001-shared` key could never collapse the two entries in a lookup dictionary, per the implementation plan's explicit instruction.
- Considered exercising the full ten-value finding-status vocabulary via a temporary-copy table (as `audit_finding_test.go`'s existing `TestParseFindingFile_AllStatusesValid` already does at the production-parser level). Judged not needed here: that table already exists and passes; this task's characterization need is specifically the *invalid/unknown*-status heal-and-diagnose path in a migration context, which the single `escalated` variant demonstrates without duplicating existing coverage.
- Failure-variant tests (malformed finding frontmatter, malformed run frontmatter, broken `depends_on` reference) mutate copies only in `t.TempDir()`; `TestMigrationHistoryFailures` finishes with `assertFixtureBytesMatchManifest` re-verifying the frozen bytes, and `assertHistoryInventoryMatchesManifest` (called from `TestMigrationHistoryDispositions`, mirroring `TestMigrationSourceBasicInventory`) diffs the on-disk file set against the manifest in both directions and confirms `project/AGENTS.md` is absent by design — this V1 project variant never adopted a managed agent guide.
- No unsupported/ambiguous source behavior was found that needed deferral to E45; the existing parser/resolver/audit-loader behavior fully accounts for every scoped/disposition case this fixture models.

**Acceptance evidence** (`internal/data/migration_history_test.go`):

| AC | Test |
|---|---|
| Two releases share epic/short Task ID; v1 done vs v1.1 in-progress; v1.1 dependency resolves within v1.1; two D001 distinct by release | `TestMigrationHistoryScopedReferences` |
| F001 fixed/unmet proof, F002 waived/owner reason, F003 duplicate_of F001, immutable run, epic audit, custom Health-Check, authored history/unknown metadata | `TestMigrationHistoryDispositions`, `TestMigrationHistoryRawAndNormalized` |
| Manifest records bytes/hashes, scoped IDs, disposition, proof, duplicate link, active-vs-history | `manifest.yml`; verified by `assertHistoryInventoryMatchesManifest` (called from `TestMigrationHistoryDispositions`) |
| `ParseRawFindingFile`/`LoadAuditRegisterSet` retain distinctions; unknown-status variant raw-preserved + healed to open; diagnostics observed | `TestMigrationHistoryRawAndNormalized` |
| Named tests: raw fixed ≠ verified, waiver is owner disposition, duplicate → F001, run content/commit retained | `TestMigrationHistoryDispositions` |
| Malformed run/frontmatter, unresolved reference → existing named error/missing result; frozen bytes unchanged | `TestMigrationHistoryFailures` (three subtests + final byte check) |
| Focused tests + `make build && make test` pass | see Quality gates below |

**Quality gates:** `go test ./internal/data -run 'MigrationSource|MigrationHistory' -v -count=1` — pass (4/4 new named tests incl. 3/3 subtests, plus T001's 3/3 existing tests incl. 2/2 subtests, all uncached). `make build && make test` — pass, all packages ok. `gofmt -l internal/data/migration_history_test.go` — no output. `go vet ./...` — clean. `git diff --check` — clean.

**Health check:** skipped; this repository (the live Savepoint project, not the frozen fixture) has no `.savepoint/Health-Check.md`.

**Limits:** This task adds read-only characterization fixtures and tests only; no production loader, writer, or migration-mapping code changed. The remaining migration/recovery verification (actual V1→V2 conversion, ID reassignment, recoverable writes) is out of scope here and belongs to E45, per the task Boundaries.

## Drift Notes

None. `internal/data/testdata/migration/v1-history/` is new test-fixture data with its
test coverage added to the existing `internal/data` package via
`migration_history_test.go` — no new production module, and the AGENTS.md Codebase Map's
`internal/data` entry already covers this package's scope. Any new production
responsibility would require replanning rather than expanding this fixture task.
