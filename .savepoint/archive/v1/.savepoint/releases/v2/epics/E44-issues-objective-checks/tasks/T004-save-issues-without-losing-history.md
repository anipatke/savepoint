---
id: E44-issues-objective-checks/T004-save-issues-without-losing-history
title: Save issues without losing history
status: done
objective: Create Issue records once, patch their managed fields safely, and append history entries without rewriting the past.
depends_on:
    - E44-issues-objective-checks/T001-record-durable-follow-up
    - E44-issues-objective-checks/T003-prove-an-issue-is-closed
complexity_tier: medium
complexity_reason: Extends the existing preserving write path with append-only list semantics and a new create path.
---

# T004: Save issues without losing history

## Problem

An Issue is mutable where a Check is immutable, and its history is the part that must never be rewritten. The current write path patches scalars and nested blocks, but nothing can append to a list while refusing to shorten, reorder, or edit what is already recorded. Objective evidence also has no writer.

## Context Files

- `internal/data/write.go`
- `internal/data/write_test.go`
- `internal/data/issue_v2.go`
- `internal/data/issue_v2_test.go`
- `internal/data/project.go`
- `internal/data/errors.go`
- `internal/data/objective_v2.go`
- `internal/data/evidence_v2.go`

## Acceptance Criteria

- [x] `CreateIssueV2` allocates the next unused `I###` over active records, writes create-only to `.savepoint/issues/{id}-{slug}.md`, validates the marshalled content decodes as an Issue before creating the file, and refuses an existing path rather than overwriting it.
- [x] A rejected Issue creation leaves no file behind.
- [x] Managed Issue writes patch only `status`, `resolution`, `duplicate_of`, `tasks`, `checks`, and `severity`, preserving every other YAML key, unknown fields, and the authored Markdown body unchanged.
- [x] A history write appends entries only: a write whose entry list is shorter than, reorders, or edits any already-recorded entry is refused with a named diagnostic and leaves the file untouched.
- [x] Appending a history entry preserves every earlier entry's fields byte-for-byte in the rewritten frontmatter.
- [x] `WriteObjectiveEvidenceV2` patches an Objective's `last_check`, `freshness`, `owner_validation`, `exception`, and `replan` through the same patch builders `WriteTaskEvidenceV2` uses, with no second evidence contract.
- [x] Every write validates the patched content decodes before any file is replaced, so a malformed value is refused and the file is left untouched.
- [x] Writing values a record already has is a no-op: the file's bytes and modification time are unchanged.
- [x] Preserving writes retain the supported line-ending form, matching existing V2 write behavior.

## Implementation Plan

- [x] Add `NewIssueV2`, `CreateIssueV2`, and `nextV2IssueID` to `write.go`, reusing the create-only file sequence and next-unused allocation pattern `CreateCheckV2` established.
- [x] Add `WriteIssueV2` for managed field patches through `writeV2Record`, validating with `DecodeIssueV2`.
- [x] Add an append-only history patch that compares the incoming entry list against the recorded one and refuses any non-append change before writing.
- [x] Add the named append-only violation sentinel to `errors.go`.
- [x] Generalize the evidence patch builder so `WriteObjectiveEvidenceV2` reuses it, validating with `DecodeObjectiveV2`.
- [x] Test create-only collision, validation failure leaving no file, managed patches preserving unknown fields and bodies, append-only acceptance and every refusal form, Objective evidence patches, no-op stability, and line-ending preservation.
- [x] Run the focused `internal/data` suite.

## Context Log

**Read:** `.savepoint/router.md`, `E44-Detail.md`, this task, `T001` and `T003` (for schema/obligations already settled), `.savepoint/Guardrails.md` STYLE rules, and the listed context files in `internal/data`.

**Edited:** `internal/data/write.go`, `internal/data/write_test.go`, `internal/data/errors.go`, `internal/data/issue_v2.go`.

**Schema/design decisions made during build:**

- A newly created Issue always starts `status: open` — `NewIssueV2` carries no resolution or history fields; those are recorded later through `WriteIssueV2` and `WriteIssueHistoryV2`, matching the design's "Open -> in_progress when repair starts" flow.
- `CreateIssueV2` writes `.savepoint/issues/{id}-{slug}.md`, matching the AC's filename form and the `{ID}-{slug}.md` convention Tasks and Epics already use. Added `issueV2Slug` (lowercase, hyphen-separated, no external deps) since no slug helper existed in `internal/data`. Discovery only requires a filename to start with its declared ID, so an empty slug falls back to the bare ID.
- `createV2CheckFile` was generalized into `createV2RecordFile(path, content, existsErr)`, parameterizing the collision sentinel; `createV2CheckFile` is now a thin wrapper passing `ErrV2CheckImmutable`, so `CreateIssueV2` reuses the exact same create-only temp-file-then-hard-link sequence. Added `ErrV2IssueAlreadyExists` (deliberately not "immutable" — an Issue is mutable after creation; only the create-only step refuses the collision) rather than reusing `ErrV2CheckImmutable`, which is worded specifically for Checks.
- Managed Issue writes (`status`, `resolution`, `duplicate_of`, `tasks`, `checks`, `severity`) all remove their key when the in-memory value is empty/nil rather than writing an empty scalar or `[]`, matching the "absent optional field stays absent" convention evidence blocks already use. Added `issueV2Frontmatter` `omitempty` tags for the same fields so `CreateIssueV2`'s direct marshal agrees with `WriteIssueV2`'s patch semantics — without this, a freshly created Issue with no tasks/checks/severity would round-trip through `WriteIssueV2` as a non-no-op (the patch would remove a key `CreateIssueV2` had written empty). Verified with `TestCreateIssueV2ThenWriteIssueV2_roundTripsAsNoOp`.
- `taskEvidencePatches` renamed to `evidencePatches` (it was already Task-agnostic — it only touches the shared `Evidence` type) so `WriteObjectiveEvidenceV2` calls the identical builder `WriteTaskEvidenceV2` uses; no second evidence contract was written.
- History append-only enforcement (`validateAppendOnlyIssueHistory`) requires the incoming list to be at least as long as the recorded one and to reproduce every recorded entry's fields (`at`, `actor`, `kind`, `note`, `check`) in the same order; any shorter, reordered, or edited entry returns `ErrV2IssueHistoryNotAppendOnly` before any write happens.

**Acceptance verification.** All AC and plan items above map to a passing test in `write_test.go`:

- Create-only: `TestCreateIssueV2_writesNewFileAllocatesFirstID`, `_allocatesNextIDOverPopulatedIndex`, `_refusesExistingPathAndLeavesItUntouched`, `_rejectsMalformedRecordLeavesNoFileBehind`.
- Managed patches: `TestWriteIssueV2_updatesManagedFieldsPreservesUnknownFieldsAndBody`, `_clearingResolutionOnReopenRemovesKey`, `_refusesUnsupportedStatusAndLeavesFileUntouched`, `_preservesCRLFLineEndings`.
- No-op stability: `TestWriteIssueV2_noOpLeavesBytesAndMtimeUnchanged`, `TestCreateIssueV2ThenWriteIssueV2_roundTripsAsNoOp`.
- History append-only: `TestWriteIssueHistoryV2_appendsEntryPreservesEarlierEntries`, `_refusesShorterListLeavesFileUntouched`, `_refusesEditedEntryLeavesFileUntouched`, `_refusesReorderedEntriesLeavesFileUntouched`, `_noOpLeavesBytesAndMtimeUnchanged`.
- Objective evidence reuse: `TestWriteObjectiveEvidenceV2_setsSubBlocksPreservesUnknownFieldsAndBody`, `_noOpLeavesBytesAndMtimeUnchanged`.

**Quality gates.** `go test ./internal/data/...` ok. `make build && make test` ok — all packages pass. `go vet ./...` clean. `gofmt -l internal/data/*.go` reports only `router.go` and `task_test.go`, both pre-existing and untouched by this task (per T001/T003 notes). `.savepoint/Health-Check.md` is absent from this project, so the Quick check step was skipped.

**Not done here, by design.** Wiring Issue creation/writes into any skill, doctor repair, or board surface is out of scope per the epic boundaries (E46/E48/E49 own consumers). `guardrail_ids` is not a managed field and is left untouched by `WriteIssueV2`, matching the AC's exact field list.
