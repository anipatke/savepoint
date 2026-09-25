---
id: E44-issues-objective-checks/T001-record-durable-follow-up
title: Record durable follow-up
status: done
objective: Load Issue records into the project index with strict decoding, confined discovery, and a validated duplicate graph.
depends_on: []
complexity_tier: high
complexity_reason: Adds a fourth record family across decoding, confined discovery, and index graph validation.
---

# T001: Record durable follow-up

## Problem

Follow-up that outlives a single evaluation has nowhere to live. A Check can name `I###` references, but no Issue record exists for them to resolve against, so a repeated observation cannot keep its identity and a duplicate cannot point anywhere.

## Context Files

- `internal/data/issue_v2.go`
- `internal/data/issue_v2_test.go`
- `internal/data/check_v2.go`
- `internal/data/errors.go`
- `internal/data/parser.go`
- `internal/data/evidence_v2.go`
- `internal/data/discover.go`
- `internal/data/discover_test.go`
- `internal/data/project.go`
- `internal/data/project_test.go`
- `internal/data/dependency.go`

## Acceptance Criteria

- [x] An Issue record requires a valid global `I` ID, a non-empty `title`, a `type` of `defect`, `drift`, `guardrail`, `verification`, or `other`, a `status` of `open`, `in_progress`, or `resolved`, and a `source` naming what produced it with its actor and time.
- [x] Issue `status` is its own vocabulary: `planned`, `done`, `in_progress` with a stage, or any other Task lifecycle value is rejected with a named diagnostic and never healed into a valid-looking default.
- [x] `tasks`, `checks`, and `guardrail_ids` decode as lists; `tasks` and `checks` are validated as `T###` and `C###` identity references for shape only, and `guardrail_ids` are opaque policy strings that are never resolved against a Guardrails file.
- [x] Optional `severity`, `resolution`, and `duplicate_of` decode when present; an absent optional field stays absent rather than becoming a permissive default.
- [x] `resolution` decodes its disposition (`verified`, `accepted`, `duplicate`), optional proof `check`, actor, time, and reason as a shape; cross-record proof obligations are out of scope for this task.
- [x] `history` decodes as an ordered list of `{at, actor, kind, note, check}` entries whose `kind` is one of `observed`, `repair_attempted`, `rechecked`, `deferred`, `reopened`, `owner_decision`; a malformed or unrecognized entry is a named diagnostic.
- [x] Issues are discovered only from `.savepoint/issues/`, through the same path confinement as Checks, rejecting traversal, symlink escape, case aliasing, duplicate IDs, and a filename that does not start with its declared ID.
- [x] `V2Index` exposes Issues keyed by `I###`.
- [x] `duplicate_of` resolves to an existing, different Issue; a missing target, a self-reference, or a cycle each return a distinct named diagnostic.
- [x] A project with no `issues/` directory loads with no Issues and no error.
- [x] Malformed Issue records fail the load closed, in deterministic ID order, exactly as E42 and E43 structural diagnostics do.

## Implementation Plan

- [x] Add `IssueV2`, its type/status/disposition/history-kind types, and `DecodeIssueV2` in a new `issue_v2.go`, reusing `ParseV2Document` and the existing anchored ID patterns.
- [x] Add the named Issue sentinels to `errors.go` for malformed records, invalid identity, invalid lifecycle vocabulary, and duplicate-graph conflicts.
- [x] Add `v2IssuesDirName` and confined `.savepoint/issues/` discovery in `discover.go`, reusing `v2PathConfiner` and the existing filename/ID mismatch rule.
- [x] Extend `V2Index` and `LoadV2Index` in `project.go` with the Issues map and duplicate-graph validation, reusing `findV2Cycle`/`describeV2Cycle` from `dependency.go`.
- [x] Test valid records, every malformed field, Task vocabulary rejection, duplicate IDs, path mismatch, unsafe paths, history entry shapes, valid and missing and self and cyclic `duplicate_of`, absent `issues/`, and deterministic diagnostic ordering.
- [x] Run the focused `internal/data` suite.

## Context Log

**Read:** `.savepoint/router.md`, `E44-Detail.md`, this task, sibling tasks T002–T004 (scope boundaries only), `.savepoint/Guardrails.md` STYLE rules, `.savepoint/releases/v2/v2-Design.md` sections 5–6, and the listed context files in `internal/data`.

**Edited:** `internal/data/issue_v2.go` (new), `internal/data/issue_v2_test.go` (new), `internal/data/errors.go`, `internal/data/discover.go`, `internal/data/discover_test.go`, `internal/data/project.go`, `internal/data/project_test.go`, `internal/data/evidence_v2.go`.

**Schema decided during build.** The epic fixed the Issue field list but not two sub-shapes, so this task settled them:

- `source: {kind: check|report|migration, check: C###, actor: {role, session}, at: RFC3339}`. `check` is required when `kind: check` and rejected otherwise, so a reported or migrated Issue cannot borrow an evaluation's authority.
- `resolution: {disposition, check, actor: {role, session}, at, reason}`, decoded for shape only. Which disposition may name a proof Check, and whether a resolution belongs on a given status, are T003's cross-record obligations.
- `severity` is an opaque non-empty string. Guardrails place severity under policy ownership, so no vocabulary is enforced here and no default is supplied; a blank value is a diagnostic rather than an absent one.
- History `note` and `check` stay optional. A `rechecked` entry naming its Check is meaningful without prose; per-kind obligations (a deferral's reason) belong to T003.
- `IssueV2.Origin` holds the decoded `source:` block so `IssueV2.Source` stays the parsed document, matching Objective, Task, and Check.

**Reuse decisions.** `DiscoverV2Checks` and the new `DiscoverV2Issues` are now both thin calls over one generic `discoverV2FlatDir`, so the confinement, filename-prefix, and duplicate-identity rules exist once instead of twice (STYLE-07). Check discovery diagnostics are byte-identical to before and its existing tests pass unchanged. `decodeEvidenceActor`/`decodeEvidenceTimestamp` became `decodeV2Actor`/`decodeV2Timestamp` taking the calling family's malformed sentinel, so an Issue's bad actor reports as an Issue problem rather than an evidence one; all six call sites are inside `evidence_v2.go` and `check_v2.go` is untouched.

**Acceptance verification.** 36 Issue tests plus subtests, all passing:

- Identity, title, type, status, source — `TestDecodeIssueV2_valid`, `_malformedID` (6 cases), `_missingTitle`, `_type`, `_statusVocabulary`, `_source` (11 cases), `_sourceReportAndMigration`.
- Task vocabulary rejection — `TestDecodeIssueV2_rejectsTaskLifecycleVocabulary` covers `planned`, `done`, legacy `todo`/`complete`, and a stage on both `in_progress` and `open`, asserting no record is returned alongside the diagnostic.
- Lists and opacity — `_identityReferences` (5 cases), `_guardrailIDsAreOpaque` proves an unknown policy ID decodes and is never resolved.
- Optional absence — `_absentOptionalsStayAbsent`, `_blankSeverityRejected`.
- Resolution and history shape — `_resolutionShape` (7 cases), `_resolutionDispositionsDecodeAsShape`, `_historyEntryShapes` (8 cases), `_historyKinds` (all 6), `_historyDiagnosticNamesTheEntry` proves the diagnostic names `history[1]`, and `_valid` asserts recorded order is preserved.
- Confined discovery — `TestDiscoverV2Issues_valid`, `_absentIssuesDir`, `_fileNameMismatch`, `_duplicateID`, `_rejectsTraversalFilename`, `_rejectsSymlinkEscape`, `_rejectsCaseCollision`, `_malformedRecordFailsClosed`.
- Index and duplicate graph — `TestLoadV2Index_issuesValid`, `_issuesAbsentDirLoadsEmpty`, `_issueDuplicateOfResolves`, `_issueDuplicateOfMissingTarget` (names both records), `_issueDuplicateOfSelf`, `_issueDuplicateOfCycle` (three-Issue ring), `_issueDuplicateDiagnosticsAreDistinct` (the three sentinels do not match one another), `_issueDeterministicDiagnosticOrder` (identical error across runs, lowest ID first), `_malformedIssueFailsClosed` (nil index alongside the diagnostic).

**Quality gates.** `go test ./internal/data/` ok. `make build` ok. `make test` ok — all packages pass. `go vet ./...` clean. `gofmt -l` reports no file touched by this task; the five files it lists were already unformatted on `v2` before this change.

**Not done here, by design.** Resolving `tasks`/`checks`/Check `issues` references against the index is T002. Resolution proof obligations, reopen/deferral rules, and the no-Issue-gates proof are T003. Issue writes and history append-only enforcement are T004. `.savepoint/Health-Check.md` is absent from this project, so the Quick check step was skipped.

## Drift Notes

`history` is implemented in Issue frontmatter, while `v2-Design.md` section 6 still lists History as an authored body section. This is the deliberate refinement E44-Detail.md records in its Architectural delta and Open decisions, which state it "must be reconciled there during this epic." No task in the breakdown owns that edit, so `v2-Design.md` section 6 is still unreconciled after T001 and needs an explicit owner before epic closeout.
