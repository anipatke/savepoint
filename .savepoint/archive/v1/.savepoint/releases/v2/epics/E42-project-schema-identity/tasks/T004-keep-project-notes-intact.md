---
id: E42-project-schema-identity/T004-keep-project-notes-intact
title: Keep project notes intact during updates
status: done
objective: Update owned V2 record fields without altering unknown frontmatter or authored Markdown content.
depends_on:
    - E42-project-schema-identity/T002-read-v2-work-safely
complexity_tier: medium
complexity_reason: Extends shared parsing and writing while preserving byte-sensitive user-authored record content.
---

# T004: Keep project notes intact during updates

## Problem

Typed V2 records need managed updates, but reconstructing files from typed structs would discard unknown YAML and user-authored Markdown needed by future versions and migration.

## Context Files

- `internal/data/parser.go`
- `internal/data/parser_test.go`
- `internal/data/write.go`
- `internal/data/write_test.go`
- `internal/data/objective_v2.go`
- `internal/data/objective_v2_test.go`
- `internal/data/task_v2.go`
- `internal/data/task_v2_test.go`
- `internal/data/project.go`
- `internal/data/project_test.go`

## Acceptance Criteria

- [x] Loaded V2 records retain the parsed source document needed to patch only fields owned by a managed write.
- [x] A no-op write changes neither file bytes nor modification time.
- [x] A managed-field write preserves unknown YAML keys and values, authored Markdown bodies, and the supported source line-ending form.
- [x] Managed writes validate the resulting V2 record and refuse malformed or unsafe lifecycle values before replacing the file.
- [x] Write failures return path-qualified errors without truncating the existing record.
- [x] Existing V1 status-write preservation and concurrency behavior remain covered and unchanged.

## Implementation Plan

- [x] Define the lossless source-document contract shared by V2 decoders and writers.
- [x] Add YAML-node patch helpers that update only explicit owned fields and preserve unknown mappings and values.
- [x] Implement no-op detection and pre-replacement validation for V2 Objective and Task writes.
- [x] Keep the established safe replacement and V1 status-write contracts intact rather than creating an unscoped filesystem framework.
- [x] Test unknown nested fields, multiline YAML, authored headings/body text, LF/CRLF input, no-op mtimes, validation refusal, and injected write failures.
- [x] Run focused parser/write/data tests and record preservation evidence.

## Context Log

**Files read:** `internal/data/parser.go`, `internal/data/parser_test.go`, `internal/data/write.go`, `internal/data/write_test.go`, `internal/data/objective_v2.go`, `internal/data/objective_v2_test.go`, `internal/data/task_v2.go`, `internal/data/task_v2_test.go`, `internal/data/project.go`, `internal/data/project_test.go`, `internal/data/errors.go`, `internal/data/lifecycle.go` (ValidateTaskLifecycleStateForWrite/IsCanonicalTaskStatus contracts), `AGENTS.md`, `.savepoint/Guardrails.md`, `.savepoint/router.md`, `E42-Detail.md`, `T002-read-v2-work-safely.md` (Context Log for prior decode-boundary decisions), `.savepoint/releases/v2/v2-Design.md` (write/preservation passages only).

**Files edited:**
- `internal/data/parser.go` — added `V2SourceDocument.CRLF` and set it in `ParseV2Document` from the original (pre-normalization) content, so a loaded V2 record retains the line-ending form a later managed write must reproduce.
- `internal/data/write.go` — added the V2 managed-write boundary: `mappingFieldValue` (read a mapping key's current scalar value), `v2FieldPatch`/`patchV2Mapping` (apply only named-key patches and report whether anything actually changed, supporting both set and remove), `cloneYAMLNode` (deep-copy a `yaml.Node` so patching never mutates the caller's loaded record), `writeV2Record` (the shared no-op/validate-before-replace/CRLF-reproducing write path), and the two public writers `WriteObjectiveV2` (patches `status`) and `WriteTaskV2` (patches `status` and `stage`, removing `stage` entirely rather than writing it empty when leaving `in_progress`). Both writers validate the patched content by round-tripping it through `DecodeObjectiveV2`/`DecodeTaskV2` before any file is replaced, reusing the existing strict decoders as the single source of validation truth instead of duplicating lifecycle rules. Also changed `removeMappingField` to return `bool` (whether a key was present and removed) — existing V1 call sites in `WriteTaskStatus`/`WriteDefectStatus` discard the new return value and are unaffected.
- `internal/data/write_test.go` — added `TestWriteObjectiveV2_*` and `TestWriteTaskV2_*` covering: unknown top-level and nested-mapping fields plus authored heading/body text preserved across a status change; a true no-op (write already-current status) changing neither file bytes nor modtime, for both record types; validation refusing an unsupported Objective status and an in_progress Task with no stage, in both cases leaving the file byte-for-byte untouched; `depends_on` records (including a non-default `requires: accepted`) surviving a managed write; stage removed (not blanked) when a Task leaves `in_progress`; CRLF source round-tripping back to CRLF output while an LF source stays LF; and a read-only-file write failure returning a path-qualified error while leaving the original file content intact (skipped under root, where permission bits do not block writes).

**Named cases and results:**
- `TestWriteObjectiveV2_updatesStatusPreservesUnknownFieldsAndBody` — PASS
- `TestWriteObjectiveV2_noOpLeavesBytesAndMtimeUnchanged` — PASS
- `TestWriteObjectiveV2_refusesUnsupportedStatusAndLeavesFileUntouched` — PASS
- `TestWriteObjectiveV2_writeFailureIsPathQualifiedAndLeavesRecordIntact` — PASS
- `TestWriteTaskV2_updatesStatusAndStagePreservesDependencies` — PASS
- `TestWriteTaskV2_removesStageWhenLeavingInProgress` — PASS
- `TestWriteTaskV2_refusesInProgressWithoutStageAndLeavesFileUntouched` — PASS
- `TestWriteTaskV2_noOpLeavesBytesAndMtimeUnchanged` — PASS
- `TestWriteTaskV2_preservesCRLFLineEndings` — PASS
- `TestWriteTaskV2_preservesLFLineEndingsByDefault` — PASS
- Full `go test ./internal/data/...` — PASS (no regressions; all existing `TestWriteTaskStatus_*`, `TestWriteDefectStatus_*`, `TestWriteRouterState_*` cases unchanged and passing)
- `go vet ./...` — clean

**Quality gates:** `make build && make test` — PASS across all packages (`.`, `cmd`, `internal/board`, `internal/buildtool`, `internal/data`, `internal/doctor`, `internal/init`, `internal/styles`).

No `.savepoint/Health-Check.md` in this project; Quick check step skipped per `savepoint-build-task`.

No Drift Notes: all edits are to `parser.go` and `write.go`, the two files E42-Detail.md's Components table already names for the raw-document boundary and preserving writes; no new files or packages, no architecture change. Deliberately did not add V1-style optimistic-concurrency (`expectedMtime`) to `WriteObjectiveV2`/`WriteTaskV2`: the ACs ask that V1's existing concurrency behavior stay covered and unchanged, not that V2 gain the same contract, and dependency-satisfaction/lifecycle-transition gating is explicit E43 scope per E42-Detail.md's Boundaries — adding it here would be speculative scope beyond what this task's ACs call for.
