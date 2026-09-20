---
id: E48-next-resume/T002-resolve-the-selection-or-say-why-not
title: Resolve the selection or say why not
status: done
objective: Turn the router's hint into an exactly matched Objective/Task selection, or a named diagnostic that never substitutes a different record.
depends_on:
    - E48-next-resume/T001-read-the-projects-own-current-state
complexity_tier: medium
complexity_reason: Small surface but the substitution and mismatch rules carry most of the epic's trust risk.
---

# T002: Resolve the selection or say why not

## Problem

The router names what the user was last working on. Those names are hints written by a human or an agent, and the records they point at move, get renamed, get converted, and get finished. Resolution is therefore a step with two honest outcomes and one forbidden one.

The forbidden one is substitution. If the router names `T014` and the index holds `T014a`, `T104`, or one Task whose title looks similar, the answer is not "you probably meant this" — it is that `T014` does not resolve. Picking a near match silently reattributes a user's evidence, their owner acceptance, and their next action to a record they never selected. This is the single rule in this task that everything else exists to protect.

The second failure is quieter: a `task` and an `objective` that both resolve but disagree, where the named Task's `objective` field is not the named Objective. Ownership in `V2Index.ObjectiveTasks` comes from the Task record, never from the router, so the Task record wins and the disagreement is reported. It usually means the router line is stale, and stale-but-plausible is worse than absent.

An unresolved selection is not the end of the answer. The user still needs to know what to do, so resolution returns the diagnostic *and* leaves the project in a state where T003's ladder can still produce a next action from the records themselves. A project whose router points at nothing is still a project with work in it.

Archived records are out of scope for lookup. The live index is the only thing consulted; a record that was migrated into an archive is simply not in it, and the diagnostic says the ID is not among the project's live records rather than guessing where it went.

## Context Files

- `internal/data/router_v2.go`
- `internal/data/router_v2_test.go`
- `internal/data/project.go`
- `internal/data/project_test.go`
- `internal/data/task_v2.go`
- `internal/data/objective_v2.go`
- `internal/data/errors.go`
- `.savepoint/releases/v2/epics/E48-next-resume/E48-Detail.md`

## Acceptance Criteria

- [x] A selection resolver in `internal/data/next.go` takes a decoded V2 router state and a `V2Index` and returns either an exactly matched Objective and optional Task, or a typed selection diagnostic.
- [x] Matching is by exact global ID only; no prefix, suffix, numeric-proximity, title, or path-based fallback exists anywhere in the resolver.
- [x] A `task` ID absent from `index.Tasks` returns a typed "selection not found" diagnostic naming the ID, and the same holds for an absent `objective` ID.
- [x] A resolvable `task` whose own `objective` field differs from the router's `objective` returns a typed "selection mismatch" diagnostic naming both IDs, and the Task record's ownership is treated as authoritative.
- [x] A router with an `objective` and no `task` resolves to that Objective with no Task selected, and is not a diagnostic.
- [x] A router in `idea` or `design` with no `objective` and no `task` resolves cleanly to no selection.
- [x] Every diagnostic is a typed value the caller can branch on, not a formatted string, and carries the IDs it read.
- [x] Resolution reads only the index and the decoded router state: no filesystem access, no discovery, no write.
- [x] A project whose only Task matches the router ID exactly still resolves through the exact-match path, proven by a fixture holding both `T014` and `T140`.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Create `internal/data/next.go` with the selection result and diagnostic types.
- [x] Implement exact-ID lookup against `index.Tasks` and `index.Objectives`, with no fallback branch to add later.
- [x] Add the ownership cross-check against the Task record's `objective` field.
- [x] Return a resolved-selection value that distinguishes "no selection" from "selection failed", so T003 can tell an empty project from a broken router line.
- [x] Add fixture projects covering: exact match, absent task, absent objective, mismatch, objective-only, no selection, and the near-miss ID pair.
- [x] Assert no filesystem access in the resolver by constructing the index in memory in at least one test.
- [x] Run `go test ./internal/data/...`, then `make build && make test`.

## Context Log

**Files read:** `.savepoint/router.md`, `E48-Detail.md`, this task file, `internal/data/router_v2.go`, `internal/data/router_v2_test.go`, `internal/data/project.go`, `internal/data/project_test.go` (LoadV2Index/fixture conventions), `internal/data/task_v2.go`, `internal/data/objective_v2.go`, `internal/data/errors.go`, `internal/data/gate_v2.go` and `internal/data/objective_gate_v2_test.go` (to match the existing "typed decision struct returned by value, not an error" convention `ResolveClearance`/`ResolveObjectiveCompletion` already use), `internal/data/dependency_test.go` (reused `newV2TestIndex()` helper), completed `T001` task file (to confirm `ReadStateV2`'s invariant that a `task` selection is never decoded without an `objective`), `.savepoint/Guardrails.md` (DATA-02, DATA-03, ARCH-01/03/04, TEST-01/02/04/08).

**Files edited:**
- `internal/data/next.go` (new) — `SelectionRecordKind` (`objective`/`task`), `SelectionDiagnosticKind` (`SelectionNotFound`, `SelectionMismatch`), `SelectionDiagnostic`, `Selection`, and `ResolveSelection(index *V2Index, router *RouterStateV2) (Selection, *SelectionDiagnostic)`. Matching is two plain map lookups (`index.Objectives[router.Objective]`, `index.Tasks[router.Task]`); there is no fallback branch of any kind, so exact-ID-only matching is structural rather than a rule to remember. An empty `router.Objective` returns `Selection{}, nil` (no selection) — relying on `ReadStateV2`'s existing invariant that a `Task` is never decoded without an `Objective`, so this function does not re-validate that pairing. A `SelectionDiagnostic` is returned alongside a zero-value `Selection` in both failure cases, so a caller (T003) tells "no selection" (nil diagnostic, zero Selection) from "selection failed" (non-nil diagnostic, zero Selection) by the diagnostic pointer alone, per the task's own distinction. Followed the existing `Clearance`/`GateDecision` convention (`internal/data/gate_v2.go`) of a typed value returned directly rather than wrapped in `error`, since these are resolvable outcomes, not Go error conditions (DATA-02/DATA-03 — named diagnostics, no parallel vocabulary).
- `internal/data/next_test.go` (new) — one test per acceptance criterion: exact match, the `T014`/`T140` near-miss pair proving no numeric-proximity fallback exists, absent task, absent objective, mismatch (naming `RouterObjective`, `Task`, and `TaskObjective`), objective-only, no-selection, and an in-memory-only `V2Index` literal (no `LoadV2Index`, no disk) proving the resolver reads only its two arguments.

**Quality gates:** `go build ./...` (clean), `gofmt -l internal/data/next.go internal/data/next_test.go` (clean), `go vet ./internal/data/...` (clean), `go test ./internal/data/... -run 'TestResolveSelection' -v` (8/8 pass), `go test ./internal/data/...` (pass), `make build && make test` (pass, all packages). No `.savepoint/Health-Check.md` in this project, so the Quick check step is skipped per the build-task skill.

No drift: no new package, only a new file inside `internal/data` at the exact location `E48-Detail.md`'s Components and files table already names for the projection's selection step.
