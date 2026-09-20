---
id: E48-next-resume/T001-read-the-projects-own-current-state
title: Read the project's own current state
status: done
objective: Decode the V2 router's state, objective, task, and next_action strictly, with a named diagnostic for anything it cannot interpret.
depends_on: []
complexity_tier: low
complexity_reason: One new file over an existing anchor parser, with strict decoding and named errors.
---

# T001: Read the project's own current state

## Problem

The V2 router keeps the same `## Current state` fenced-YAML anchor the V1 router uses, but nothing else about it is the same: `state` is `idea|design|task|check`, ownership is an `objective` rather than a `release` plus `epic`, `task` is optional, and `next_action` is prose for a human. `data.RouterState` in `internal/data/router.go` decodes the V1 shape and is still the live repository's own reader, so it cannot be widened to cover both — a single struct carrying both vocabularies is exactly the parallel-vocabulary failure DATA-02 forbids, and it would let a V1 value heal into a V2 field.

So the V2 reader is a separate decode over the same anchor. `ReadState`'s anchor location and `ReplaceStateBlock` are reused as-is; only the shape decoded out of the block is new.

Strict means strict (DATA-03). An unrecognized `state` value is a named diagnostic, not a heal to `idea`. A `task` present without an `objective` is a diagnostic, because a Task selection with no owner cannot be resolved later. An absent `task` is legal — a project in `idea` or `design` has no selected Task — and so is an empty `next_action`, which is prose and never load-bearing. Nothing decoded here decides anything: this task produces a *hint*, and T002 is where that hint is resolved against real records.

## Context Files

- `internal/data/router.go`
- `internal/data/router_test.go`
- `internal/data/errors.go`
- `internal/data/task_v2.go`
- `internal/data/objective_v2.go`
- `templates/project-v2/.savepoint/router.md`
- `.savepoint/releases/v2/v2-Design.md`
- `.savepoint/releases/v2/epics/E48-next-resume/E48-Detail.md`

## Acceptance Criteria

- [x] A new `internal/data/router_v2.go` decodes the `## Current state` block into a V2 router state carrying `state`, `objective`, optional `task`, and `next_action`.
- [x] The four valid states — `idea`, `design`, `task`, `check` — are typed constants in `internal/data`, and no other package declares its own copy.
- [x] An unrecognized or empty `state` returns a named diagnostic identifying the value it read; it never defaults.
- [x] `objective` and `task` are validated against the `O###` and `T###` identity shapes already required by the record decoders, and a malformed ID returns a named diagnostic.
- [x] A `task` present with no `objective` returns a named diagnostic; a `task` absent with an `objective` present decodes cleanly.
- [x] An empty or absent `next_action` decodes cleanly and is carried as an empty string.
- [x] A missing `## Current state` block, malformed YAML, and unknown keys each produce a named diagnostic rather than a panic or a silent default.
- [x] The V1 `RouterState` decode, its fields, and its existing tests are unchanged.
- [x] Decoding performs no filesystem write and no repair of the input content.
- [x] `templates/project-v2/.savepoint/router.md`'s shipped state block decodes cleanly through the new reader, asserted by test.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Read `internal/data/router.go` and confirm which anchor-location helpers are reusable without changing V1 behavior.
- [x] Add the state constants and the V2 router state type in `internal/data/router_v2.go`.
- [x] Decode strictly with `yaml.v3` in known-fields mode, returning named errors from `internal/data/errors.go` for each invalid case.
- [x] Validate the `O###`/`T###` ID shapes using the same expressions the V2 record decoders use, rather than a second pattern.
- [x] Add the shipped-template decode test so the scaffold and the reader cannot drift apart.
- [x] Add failure-path tests: unknown state, malformed IDs, task without objective, missing block, malformed YAML, unknown key.
- [x] Run `go test ./internal/data/...`, then `make build && make test`.

## Context Log

**Files read:** `.savepoint/router.md`, `E48-Detail.md`, this task file, `internal/data/router.go`, `internal/data/router_test.go`, `internal/data/errors.go`, `internal/data/task_v2.go`, `internal/data/objective_v2.go`, `templates/project-v2/.savepoint/router.md`, `internal/migrate/convert_docs.go` (to confirm the V2 anchor shape the migration writer already produces), `internal/doctor/checks.go` (to confirm the existing `"" | "none"` absent-selection convention this task reuses rather than inventing a second one).

**Files edited:**
- `internal/data/router.go` — factored the anchor-finding half of `ReadState` (locating `## Current state` and its fenced ```yaml block) out into a new unexported `extractStateBlock`, which `ReadState` now calls. Behavior-preserving: same error strings, same existing tests pass unchanged. This is the "reused as-is" anchor step the task problem statement calls for, shared by the new V2 reader below rather than duplicated.
- `internal/data/router_v2.go` (new) — `RouterPhaseV2` typed constants (`idea`, `design`, `task`, `check`), `RouterStateV2`, and `(*RouterReader).ReadStateV2`. Decodes the block content returned by `extractStateBlock` with a `yaml.v3` `Decoder` in `KnownFields(true)` mode for strict unknown-key rejection (the V1 reader's `TestRouterReader_ignoresUnknownKeys` intentionally stays tolerant; V2 does not). Validates `state` against the closed vocabulary, `objective`/`task` against the existing `objectiveIDPattern`/`taskIDPatternV2` regexes from `objective_v2.go`/`task_v2.go` (no second pattern), and rejects a `task` selected with no `objective`. `"none"` and `""` both normalize to an absent selection via `normalizeRouterSelectionV2`, matching the same sentinel V1's `release`/`epic` fields already use in `internal/doctor/checks.go` — confirmed necessary because the shipped V2 template (`templates/project-v2/.savepoint/router.md`) ships literal `objective: none`/`task: none`, and that file must decode cleanly per acceptance criteria. Errors reuse the existing generic V2 vocabulary (`ErrV2Malformed`, `ErrV2InvalidLifecycle`, `ErrV2InvalidID`, `ErrV2InvalidOwnership`) rather than adding a parallel router-specific set (STYLE-07/DATA-02).
- `internal/data/router_v2_test.go` (new) — valid decode, `none`-sentinel normalization, task-absent/objective-present, empty `next_action`, unknown/empty state, malformed objective/task ID, task-without-objective, missing block, malformed YAML, unknown key (strict — the mirror-image case of the V1 tolerant test), a regression test pinning that `ReadState` (V1) still ignores unknown keys after the `extractStateBlock` refactor, and the shipped-template decode assertion.

**Quality gates:** `go build ./...`, `gofmt -l` (clean), `go vet ./internal/data/...` (clean), `go test ./internal/data/... -run 'TestReadStateV2|TestRouterReader' -v` (all pass), `go test ./internal/data/...` (pass), `make build && make test` (pass, all packages). No `.savepoint/Health-Check.md` in this project, so the Quick check step is skipped per the build-task skill.

No drift: no new module, only a new file inside `internal/data`, the exact Codebase Map location and package E48-Detail.md already names for the projection's inputs.
