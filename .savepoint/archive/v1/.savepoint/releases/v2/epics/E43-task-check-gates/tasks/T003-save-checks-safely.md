---
id: E43-task-check-gates/T003-save-checks-safely
title: Save new checks without losing evidence
status: done
objective: Create Check records once and patch Task evidence fields without altering unknown fields or authored content.
depends_on:
    - E43-task-check-gates/T001-record-check-results
    - E43-task-check-gates/T002-know-when-work-was-checked
complexity_tier: medium
complexity_reason: Extends the existing preserving write path to nested fields and adds create-only records.
---

# T003: Save new checks without losing evidence

## Problem

A Check must be written once and never rewritten, and evidence lives in nested frontmatter that the current managed write path can only patch as scalars. Without both, a correction would overwrite history or a write would flatten authored content.

## Context Files

- `internal/data/write.go`
- `internal/data/write_test.go`
- `internal/data/check_v2.go`
- `internal/data/evidence_v2.go`
- `internal/data/task_v2.go`
- `internal/data/project.go`
- `internal/data/errors.go`

## Acceptance Criteria

- [x] Creating a Check writes a new file under `.savepoint/checks/` and refuses an existing path with a named diagnostic; no code path rewrites or deletes an existing Check.
- [x] A new Check ID is the next unused ID across active records; allocation never reuses an ID already present.
- [x] Written Check content decodes as a valid Check before the file is created, and a rejected record leaves no file behind.
- [x] Task evidence writes patch only the named nested fields, preserving every other YAML key and value, the authored Markdown body, and the file's line-ending form.
- [x] Removing evidence — clearing a replan flag — removes the key rather than writing an empty value.
- [x] Writing values a record already has changes no bytes and no modified time.
- [x] Evidence writes keep E42's stale-source refusal: a source changed since load is rejected without touching the user's edit.
- [x] Lifecycle writes through the existing Task write path continue to behave exactly as they do today.

## Implementation Plan

- [x] Extend the managed patch representation in `write.go` to set and remove nested mapping values, keeping the existing clone-validate-replace boundary.
- [x] Add the Task evidence write entry point, patching only `last_check`, `freshness`, `owner_validation`, `exception`, and `replan`.
- [x] Add create-only Check creation with next-ID allocation over the index, validating the marshalled content before the file exists.
- [x] Add the immutability diagnostic to `errors.go`.
- [x] Test unknown-field and body preservation, CRLF preservation, no-op writes, key removal, stale-source refusal, refusal to overwrite a Check, and ID allocation against a populated project.
- [x] Run the focused `internal/data` suite.

## Context Log

**Files read:** `internal/data/write.go`, `internal/data/write_test.go`, `internal/data/check_v2.go`, `internal/data/evidence_v2.go`, `internal/data/task_v2.go`, `internal/data/project.go`, `internal/data/errors.go`, `internal/data/parser.go` (targeted read for `V2SourceDocument`/`ParseV2Document`), `internal/data/discover.go` (targeted read for `v2ChecksDirName`, Check file naming convention, path-confinement rules), `internal/data/objective_v2.go` (targeted read to confirm `Evidence` is shared with Objective but out of this task's scope), `.savepoint/releases/v2/v2-Design.md` sections 3-5 (targeted read: confirmed `checks/C001.md` has no slug in its filename, and "allocate the next unused ID from active records" phrasing), `AGENTS.md`, `.savepoint/router.md`, `agent-skills/savepoint-build-task/SKILL.md`, `E43-Detail.md`, `T001-record-check-results.md`, `T002-know-when-work-was-checked.md` (targeted reads to confirm dependency status and prior design decisions, e.g. RFC 3339 timestamp convention, four-role actor vocabulary), `.savepoint/Guardrails.md` (STYLE rules).

**Files edited:**
- `internal/data/write.go` — extended `v2FieldPatch` with a `Node *yaml.Node` field for nested mapping patches (scalar `Value` behavior unchanged); added `setMappingNode`, `mappingFieldNode`, `nodesEqualByEncoding` (style-normalized marshal comparison so a parsed node's original quote/flow style never reads as a content change), `styleNormalizedV2Node`/`clearV2NodeStyle`, and `encodeV2Node`; updated `patchV2Mapping` to route `Node`-bearing patches through the new helpers while leaving the existing scalar path untouched. Added `WriteTaskEvidenceV2` plus five small per-field patch builders (`taskEvidencePatches`, `lastCheckV2Patch`, `freshnessV2Patch`, `ownerValidationV2Patch`, `exceptionV2Patch`, `replanV2Patch`) that convert `*Evidence`'s typed sub-blocks into `v2FieldPatch` values (nil sub-block → `Remove: true`) and reuse `writeV2Record`/`DecodeTaskV2` exactly as `WriteTaskV2` does, so the clone-validate-replace and stale-source-refusal boundary is shared, not reimplemented. Added `NewCheckV2` (caller-supplied Check content minus ID/Source), `CreateCheckV2` (allocates via `nextV2CheckID`, marshals through the existing `checkV2Frontmatter`/`reviewedFrontmatter` types from `check_v2.go`, validates with `DecodeCheckV2` before any file exists, then writes via `createV2CheckFile`), `nextV2CheckID` (one past the highest numeric ID in `index.Checks`, `%03d`-padded), and `createV2CheckFile` (temp file written+synced, then `os.Link`'d into place so an existing path is never truncated/overwritten — `os.Rename` was rejected here specifically because it silently overwrites, which would violate Check immutability). Added `strconv` import.
- `internal/data/errors.go` — added `ErrV2CheckImmutable` for the Check-already-exists diagnostic.
- `internal/data/write_test.go` — added `TestWriteTaskEvidenceV2_{setsAllSubBlocksPreservesUnknownFieldsAndBody,removesReplanKeyRatherThanEmptyValue,noOpLeavesBytesAndMtimeUnchanged,refusesStaleSourceWithoutOverwritingUserEdit,preservesCRLFLineEndings,rejectsMalformedEvidenceLeavesFileUntouched}` and `TestCreateCheckV2_{writesNewFileAllocatesFirstID,allocatesNextIDOverPopulatedIndex,refusesExistingPathAndLeavesItUntouched,rejectsMalformedRecordLeavesNoFileBehind}`.

**Named cases and results:**
- `TestWriteTaskEvidenceV2_*` (6) — PASS
- `TestCreateCheckV2_*` (4) — PASS
- Full `go test ./internal/data/...` — PASS (no regressions in E42/T001/T002 suites, including existing `TestWriteTaskV2_*`/`TestWriteObjectiveV2_*` proving lifecycle writes are unaffected)
- `go vet ./...` — clean
- `gofmt -l` on edited files — clean

**Quality gates:** `make build && make test` — PASS across all packages (`.`, `cmd`, `internal/board`, `internal/buildtool`, `internal/data`, `internal/doctor`, `internal/init`, `internal/styles`).

No `.savepoint/Health-Check.md` in this project; Quick check step skipped per `savepoint-build-task`.

**Design decisions not fully spelled out in the task/epic:**
- No-op detection for nested patches compares marshalled YAML after clearing each node's `Style` field recursively, not raw bytes: a node parsed from disk keeps its author's original quote/flow style (e.g. `'2026-09-14T00:00:00Z'`), while a freshly `Encode`d replacement node has no style set and yaml.v3 defaults to double-quoting a string whose plain form would otherwise resolve to a different implicit type. Comparing raw `yaml.Marshal` output on the two would report every untouched sub-block as "changed" on its first write-again; style-normalizing both sides first (while leaving each node's explicit `!!str` tag alone, so semantic-preserving auto-quoting still applies identically on both sides) makes the comparison reflect content, matching AC "Writing values a record already has changes no bytes and no modified time."
- `createV2CheckFile` fully writes and `Sync`s a temp file in the target directory, then finalizes with `os.Link` rather than `os.Rename`: `Rename` silently replaces an existing destination on POSIX, which is exactly the immutability guarantee a Check write must not have even a narrow race window around. `Link` fails with `EEXIST` if the destination already exists, giving true create-only semantics; the temp file is removed afterward in both the success and failure path via a single deferred cleanup.
- Check filenames have no slug, matching `v2-Design.md`'s `checks/C001.md` example (unlike Objective/Task directories, which do carry a slug) — `NewCheckV2` and `CreateCheckV2` write to `checks/{id}.md` accordingly, with no filename parameter exposed to the caller.
- `NewCheckV2.Body` is documented and used as the exact bytes appended after the frontmatter closing delimiter, matching `V2SourceDocument.Body`'s existing convention (including its own leading newline), rather than CreateCheckV2 inferring or normalizing a heading/blank-line separator — keeps the one body-assembly convention in one place instead of introducing a second, slightly different one.
- ID allocation (`nextV2CheckID`) reads only `index.Checks`, i.e. active records currently loaded — migration reservations are explicitly out of scope for E43 per `v2-Design.md` ("Allocate the next unused ID from active records and migration reservations") and the epic's Boundaries section, which assigns migration to E45.

No Drift Notes: `write.go`, `write_test.go`, and `errors.go` are exactly the files E43-Detail.md's Components table already names for this write path (`internal/data/write.go | Create Check records create-only, and patch nested Task evidence fields through the existing preserving write path`; `internal/data/errors.go | Add the named Check/evidence diagnostics that doctor reports`).
