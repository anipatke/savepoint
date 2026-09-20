---
id: E49-objective-board/T001-give-the-board-a-router-writer-it-cannot-misuse
title: Give the board a router writer it cannot misuse
status: done
objective: Add a V2 router selection writer that can only set objective and task, and harden ResolveNext against a partly filled input now that a second consumer exists.
depends_on: []
complexity_tier: medium
complexity_reason: Small surface, but a write path whose whole value is what it refuses to touch.
---

# T001: Give the board a router writer it cannot misuse

## Problem

`data.WriteRouterState` writes the V1 router shape — `release`, `epic`, `state: task-building`, and a fabricated `next_action`. None of that applies to a V2 project, and the V2 board needs to record a selection when the user picks an Objective or a Task. There is no V2 equivalent today.

The interesting part is the restriction, not the write. E49 decided that the board writes `objective` and `task` and nothing else: `state` names which skill owns the conversation and the board is not running one, and `next_action` is human prose where `ResolveNext` already derives the real answer. That restriction is worth enforcing in the writer's own signature rather than trusting every caller to pass the right struct. A function that accepts a whole `RouterStateV2` and writes all of it would make the board's discipline a convention; a function that accepts a selection cannot write the rest by accident.

Everything else about the write is already established practice in this package and must not be re-invented: the `## Current state` YAML anchor is edited in place, unknown fields and the surrounding document body survive byte-for-byte (DATA-01), an mtime guard refuses a file that changed underneath the read, and a write that would change nothing changes nothing — including the file's mtime (FS-04). `WriteRouterState` and `writeV2Record` between them already demonstrate each of these; the job is to compose them, not to add a third way of writing a file.

The second item is small and is carried forward from the E48 audit's non-blocking observations: `ResolveNext` panics on a zero or partly filled `NextInput` — a nil index or nil router. It is unreachable through `resume`, because `LoadProject` always yields a non-nil index and `ReadStateV2` a non-nil router. The audit flagged it specifically because E49's board is the named second consumer of that exact call, and a board has more ways to reach a call with state half-assembled: a load that failed, a reload in flight, a model constructed by a test. A named diagnostic or a defined return is the boundary check (STYLE-06); this changes no ladder semantics and no rung, and must not.

## Context Files

- `internal/data/write.go`
- `internal/data/write_test.go`
- `internal/data/router_v2.go`
- `internal/data/router_v2_test.go`
- `internal/data/router.go`
- `internal/data/next.go`
- `internal/data/next_test.go`
- `internal/data/errors.go`
- `.savepoint/Guardrails.md`

## Acceptance Criteria

- [x] `data.WriteRouterStateV2` exists in `internal/data/write.go` and accepts only the project root, the Objective and Task selection, and the expected mtime — its signature makes writing `state` or `next_action` impossible.
- [x] A write sets `objective` and `task` in the `## Current state` anchor and leaves `state`, `next_action`, key order, and the entire surrounding document byte-identical (DATA-01). **Amended during build:** the original wording also claimed unknown anchor fields are preserved. `ReadStateV2` decodes the anchor with `KnownFields(true)`, so the V2 anchor's key set is closed and an unknown key there is a load diagnostic rather than data to preserve. The writer refuses such a router instead, which the criterion below now covers; document-level preservation outside the anchor is unaffected and is proven.
- [x] An anchor carrying a key `ReadStateV2` would reject is refused with that reader's diagnostic and the file is left untouched, so a selection is never written into a document that would then fail to load.
- [x] Clearing a selection is expressible and writes the `none` form the V2 router reader already accepts, not an empty string or a removed key.
- [x] A selection naming an ID that is not a valid `O###`/`T###` shape is refused with a named diagnostic before any write (DATA-03).
- [x] A Task selected with no Objective is refused with a named diagnostic, matching the ownership rule `ReadStateV2` already enforces on read.
- [x] A write whose result equals the file's current content is a no-op: bytes and mtime are both unchanged (FS-04).
- [x] A file whose mtime does not match the expected value is refused with `ErrMtimeConflict` and left byte-identical, with no partial write.
- [x] A missing, unreadable, or anchor-less `router.md` is refused with a named diagnostic and no write (FS-06, DATA-03).
- [x] The written file round-trips through `ReadStateV2` to the selection that was written.
- [x] `ResolveNext` returns a defined value or a named diagnostic — never a panic — for a `NextInput` with a nil `Index`, a nil `Router`, or both, including with `Migration.Pending` set and unset.
- [x] The nil-input change alters no rung, no selection, and no `Next` field for any input that was already valid; the existing `next_test.go` cases pass unchanged.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Read `WriteRouterState` and `writeV2Record` and decide which existing primitive the V2 writer composes rather than duplicating the replace-and-guard mechanics.
- [x] Add `WriteRouterStateV2` with a selection-shaped parameter, validating ID shape and the Task-requires-Objective rule before touching the file.
- [x] Reuse the existing anchor edit and preservation path so the V1 and V2 router writers share one mechanism (STYLE-07).
- [x] Add the no-op short circuit that leaves mtime untouched, mirroring the existing V2 record writers.
- [x] Add `write_test.go` cases: field preservation, body preservation, clearing to `none`, invalid IDs, Task-without-Objective, no-op idempotence, mtime conflict, missing file, missing anchor, and a `ReadStateV2` round trip.
- [x] Add the nil-input guard to `ResolveNext` and record in its doc comment why the guard exists rather than the caller being trusted.
- [x] Add `next_test.go` cases for nil index, nil router, both nil, and each with `Migration.Pending` set.
- [x] Run `go test ./internal/data/...`, then `make build && make test`.

## Context Log

**Files read:** this task file, `E49-Detail.md`, `.savepoint/router.md`, the `STYLE` rules in `.savepoint/Guardrails.md`, `internal/data/write.go` (`WriteRouterState`, `writeV2Record`, `replaceV2File`, `checkV2SourceFresh`, `setMappingField`/`mappingFieldValue`/`patchV2Mapping`, `WriteObjectiveV2`/`WriteTaskV2` as the no-op and validate-before-replace precedent), `internal/data/write_test.go` (`TestWriteRouterState_*`, `TestWriteIssueV2_noOpLeavesBytesAndMtimeUnchanged`, `TestWriteIssueV2_preservesCRLFLineEndings` as the shapes to mirror), `internal/data/router.go` (`extractStateBlock`, `ReplaceStateBlock`, `stateBlockStart`), `internal/data/router_v2.go`, `internal/data/router_v2_test.go`, `internal/data/next.go`, `internal/data/next_test.go`, `internal/data/errors.go`, `internal/data/objective_v2.go` and `task_v2.go` (the `objectiveIDPattern`/`taskIDPatternV2` regexes reused for selection validation), `templates/project-v2/.savepoint/router.md` (the real document shape the fixture reproduces).

**Files edited:**

- `internal/data/write.go` — added `RouterSelectionV2`, its `validate`, `WriteRouterStateV2`, `patchRouterSelectionV2`, `routerSelectionValueV2`, and the `routerSelectionNoneV2` constant, placed directly above `WriteRouterState` so both router writers sit together. `WriteRouterStateV2` deliberately does **not** reuse `WriteRouterState`'s mechanism: that function marshals a whole `RouterState` struct over the anchor, which would drop any key the struct does not declare — acceptable for V1, fatal for a writer whose contract is that `state` and `next_action` survive untouched. It instead parses the anchor into a `yaml.Node` tree and patches exactly two keys through the existing `setMappingField`/`mappingFieldValue` helpers, so key order, scalar quoting style, and every unpatched key pass through; the document outside the anchor is preserved by the existing `ReplaceStateBlock`, whose whole purpose is that boundary. Validation runs in three places for three different reasons: `validate()` before the file is opened, so a malformed selection never touches disk; `ReadStateV2` on the composed content before replacement, so a document this package could not read back is refused while the file is still untouched; and a second `Lstat` mtime comparison inside `replaceV2File`'s `finalCheck`, closing the read/rename window the way `writeV2Record` does. CRLF is detected on the original and re-applied after `ReplaceStateBlock` normalizes (CFG-02), matching `writeV2Record`'s `source.CRLF` handling.
- `internal/data/next.go` — `ResolveNext` now reads a nil `Index` or `Router` as an empty project and an empty selection instead of panicking, with the reason recorded at the boundary. The substitution happens after the pending-migration rung, so migration still outranks everything even with no index; `relevantIssues` was switched from `input.Index` to the normalized local, which was the actual remaining nil dereference. Go's nil-map reads make an empty `V2Index` safe for every rung, so the guard adds no branch to the ladder itself.
- `internal/data/write_test.go` — added `routerV2FixtureContent`, `writeRouterV2Fixture`, `readFileString`, and eleven tests: full-document byte preservation on a successful write, clearing to the `none` sentinel, not adding sentinel keys a document lacks, adding keys to record a real selection, a seven-case table of refused malformed selections, no-op bytes-and-mtime, stale mtime, a six-case table of unreadable documents, a missing file, a router that could not be read back, a `ReadStateV2` round trip, and CRLF preservation.
- `internal/data/next_test.go` — added `TestResolveNext_nilInputsReturnAValueRatherThanPanicking` (six cases across nil index, nil router, both, zero value, and each with migration pending) and `TestResolveNext_nilInputMatchesAnEmptyProject`, which pins the guard to the reading it claims rather than to an arbitrary return.

**Acceptance criterion amended during build.** The planned criterion claimed a selection write preserves "every unknown field" in the anchor. It does not, and should not: `ReadStateV2` decodes the anchor with `KnownFields(true)`, so the V2 anchor has a closed four-key set and an unknown key there is a named load diagnostic (DATA-03), not author data. A project carrying one never loads, so writing into it would produce a file that still fails to load. The writer refuses it with the reader's own diagnostic and leaves the file untouched; the criterion was rewritten to state that, and a new criterion and test case cover the refusal. Document-level preservation outside the anchor — prose, headings, and a second fenced YAML block placed after the anchor specifically to catch an over-broad edit — is unchanged and proven byte-for-byte.

**Verification note.** Before implementing, a throwaway probe test confirmed that a `yaml.Node` round trip through `yaml.Marshal` does not re-fold the template's ~180-character quoted `next_action` or rewrite flow-style values, which the byte-preservation contract depends on. The probe was deleted; the `TestWriteRouterStateV2_setsSelectionAndPreservesEveryOtherByte` and `..._roundTripsThroughReadStateV2` cases assert the same property permanently against the real document shape.

**Quality gates:** `gofmt -l` clean over all four changed files; `go vet ./...` clean; `git diff --check` clean; `go test ./internal/data/` pass; `make build && make test` pass across all 11 packages. Pre-existing, untouched by this task: `internal/data/task_test.go` is reported by `gofmt -l` and is not in this task's diff.

No `.savepoint/Health-Check.md` in this project, so the Quick check step is skipped per the build-task skill; its absence is not a finding.

No drift: both changes are in `internal/data`, which `E49-Detail.md`'s Components and files table names for this task, and neither adds a file, a package, or a responsibility beyond the Codebase Map's existing `internal/data` row.
