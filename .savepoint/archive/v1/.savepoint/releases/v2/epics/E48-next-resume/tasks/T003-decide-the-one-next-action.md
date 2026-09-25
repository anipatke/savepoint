---
id: E48-next-resume/T003-decide-the-one-next-action
title: Decide the one next action
status: done
objective: Build the shared Next projection as one ordered precedence ladder over the existing gate decisions, returning exactly one next action for any project state.
depends_on:
    - E48-next-resume/T002-resolve-the-selection-or-say-why-not
complexity_tier: high
complexity_reason: The epic's core contract, a total ordered ladder over four resolver families, consumed by every later surface.
---

# T003: Decide the one next action

## Problem

E43 and E44 answered every closed question about a named record. `ResolveTaskStart`, `ResolveTaskAdvance`, `ResolveTaskCompletion`, `ResolveClearance`, `ResolveObjectiveCompletion`, and `ResolveObjectiveDependency` each return a decision with typed blockers. What none of them answers is which record matters right now, and in what order the blockers across records should be heard. Without that layer, every surface — resume, the board, doctor, an agent reading the project — invents its own ordering, and the four answers drift.

The ladder is the contract:

1. Incomplete migration, or a selection that did not resolve
2. A recorded replan flag
3. An unsatisfied dependency or owner prerequisite
4. Execute or verify the selected Task
5. A fresh Check is needed
6. Required owner validation
7. The Objective's integration Check
8. The next ready Task or Objective
9. Plan the next Objective

Two rungs carry reasons worth restating in code comments, because a later reader will be tempted to reorder them. Migration is first because a project halfway through conversion holds records whose meaning is not yet settled — reading them as if they were final is how a user gets told to finish a Task that migration is about to replace. Owner validation sits *below* fresh Check because asking an owner to accept work with no current technical clearance inverts the authority model E43 built: the checker clears, then the owner accepts.

The rule that keeps this layer honest is that it computes nothing. Every readiness claim it makes is a `GateDecision`, `Clearance`, or `ObjectiveDependencyDecision` obtained from the existing resolver and carried forward whole. If the projection ever re-reads `Evidence.Freshness` to decide staleness itself, there are two clearance rules in the codebase and one of them will be wrong (STYLE-07).

Migration state arrives as an input value, not an import. `internal/migrate` already imports `internal/data`, so a call to `migrate.PendingOperation` from here closes a cycle. The projection takes a small migration-state field on its input struct; `main.go` fills it in T005.

The ladder must be total. Every reachable project — including one with no Objectives at all — lands on exactly one rung and gets a value back, never an error. A fresh project resumes into Idea/Design planning because that genuinely is its next action, not because the empty case was special-cased.

## Context Files

- `internal/data/next.go`
- `internal/data/next_test.go`
- `internal/data/gate_v2.go`
- `internal/data/gate_v2_test.go`
- `internal/data/objective_gate_v2.go`
- `internal/data/evidence_v2.go`
- `internal/data/issue_v2.go`
- `internal/data/project.go`
- `internal/data/router_v2.go`
- `internal/migrate/operation.go`
- `.savepoint/releases/v2/v2-Design.md`
- `.savepoint/releases/v2/epics/E48-next-resume/E48-Detail.md`
- `AGENTS.md`

## Acceptance Criteria

- [x] `ResolveNext` takes an input carrying a `V2Index`, a decoded V2 router state, and a migration-state value, and returns one `Next` value for every input; it returns no error for any reachable project state.
- [x] The returned value names the rung reached as a typed kind, the selected Objective and Task when there is one, and the gate decision or clearance the rung was derived from.
- [x] The nine rungs are evaluated in the documented order, and each of the nine is reached by a fixture project asserting its kind, selection, and carried decision.
- [x] Pending migration outranks every other rung, proven by a fixture that would otherwise land on a lower one.
- [x] An unresolved selection from T002 reports its diagnostic and still yields an available next action derived from the records.
- [x] A recorded replan flag outranks dependency, execution, and Check rungs.
- [x] An unsatisfied Task dependency and an unsatisfied Objective dependency are distinguishable in the result, carrying the typed blocks from `ResolveTaskStart` and `ResolveObjectiveDependency` respectively.
- [x] `missing`, `needs_work`, `stale`, and `unknown` clearance each reach a distinguishable rung result carrying the `Clearance` value, with no restatement of freshness logic in `next.go`.
- [x] Required owner validation is reported only after technical clearance is current, and a Task needing both reports the Check rung first.
- [x] An Objective whose Tasks are all complete but whose integration clearance is not current reaches the integration rung carrying `ResolveObjectiveCompletion`'s decision.
- [x] A project with no Objectives and no Tasks returns the plan-next-Objective rung with no diagnostic and no error.
- [x] A completion decision allowed only by a recorded exception is reported as exception-allowed, carrying the `Exception`, and is never presented as clearance.
- [x] `next.go` contains no clearance, freshness, dependency, or acceptance derivation of its own, confirmed by review against the acceptance list and by tests that change only the resolver-visible inputs.
- [x] `internal/data` does not import `internal/migrate`, asserted by a test over the package's imports.
- [x] `ResolveNext` performs no filesystem access and no write.
- [x] The `internal/data/` Codebase Map row in `AGENTS.md` names the Next projection responsibility (ARCH-04).
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Define the migration-state input value first — the smallest shape that answers "is a conversion in flight, and which operation" — so the injection boundary is fixed before the ladder is written.
- [x] Define the `Next` result type and the typed rung kinds, with the carried decision as a field rather than a flattened copy of its fields.
- [x] Implement the ladder as one ordered sequence of small named predicates, one per rung, each delegating to its resolver.
- [x] Write the rung-order comment naming why migration is first and why owner validation follows fresh Check.
- [x] Build fixture helpers that construct a `V2Index` in memory for each rung, so ladder tests do not depend on scaffolded directories.
- [x] Add the outranking tests: migration over everything, replan over dependency and execution, Check over owner validation.
- [x] Add the empty-project and exception-allowed cases.
- [x] Add the import-boundary test asserting `internal/data` does not depend on `internal/migrate`.
- [x] Update the `internal/data/` Codebase Map row.
- [x] Run `go test ./internal/data/...`, then `make build && make test`.

## Context Log

**Files read:** this task file, `E48-Detail.md`, `AGENTS.md` Codebase Map, completed `T001`/`T002` task files, `internal/data/next.go`, `internal/data/next_test.go`, `internal/data/gate_v2.go`, `internal/data/objective_gate_v2.go`, `internal/data/evidence_v2.go`, `internal/data/issue_v2.go`, `internal/data/project.go`, `internal/data/router_v2.go`, `internal/data/task_v2.go`, `internal/data/objective_v2.go`, `internal/data/lifecycle.go`, `internal/data/dependency.go`, `internal/data/dependency_test.go` (`newV2TestIndex` helper), `internal/data/gate_v2_test.go` and `internal/data/objective_gate_v2_test.go` (`mustCheck`/`mustObjectiveCheck`/`mustCurrentTask`/`mustCurrentObjective` fixture helpers, reused rather than duplicated), `internal/migrate/operation.go` (`PendingOperationReport` shape, to confirm `MigrationState` only needs pending + operation ID), `internal/migrate/replace_test.go` (`TestReplaceFile_doesNotReachAtomicWrite`, the `go/build.ImportDir` import-boundary pattern this task's own boundary test mirrors).

**Files edited:**
- `internal/data/next.go` — appended `MigrationState`, `NextKind` (the nine typed rungs), `Next`, `NextInput`, and `ResolveNext` plus its unexported ladder helpers (`resolveLadder`, `objectiveInView`, `resolveTaskRung`, `taskDecisionRung`, `rungForBlockers`, `resolveObjectiveIntegrationRung`, `resolveReadyRung`, `objectiveDependenciesSatisfied`) to the file T002 started. The ladder computes nothing of its own: rungs two, three, five, and six are read by ranking the `GateBlocker.Kind` values `ResolveTaskStart`/`ResolveTaskAdvance`/`ResolveTaskCompletion` already assigned (`rungForBlockers`), rung four is `decision.Allowed` from the same three resolvers (exception-allowed included, since `GateDecision.AllowedByException`/`.Exception` are carried through unchanged rather than restated), rung seven re-reads `ResolveObjectiveCompletion` and falls through when its blockers are the incomplete-owned-Task kind rather than a clearance/owner-acceptance one, and rungs eight/nine are a project-wide scan that still defers every readiness judgement to `ResolveTaskStart`/`ResolveObjectiveDependency`. `Clearance` on `Next` is populated by a second, independent call to the existing `ResolveClearance`/`ResolveObjectiveCompletion`-adjacent `ResolveClearance` (never by reading `Evidence.Freshness` directly), satisfying the "no restatement of freshness logic" acceptance criterion literally. Selection resolution (T002) runs once per call in `ResolveNext`; its diagnostic is attached to whatever `Next` the ladder produces, so an unresolved selection still yields the project's real next action.
- `internal/data/next_test.go` — added 24 tests: one per rung (all nine reached), the outranking fixtures (migration over a ready Task, replan over both execution and check-needed), Task-dependency vs Objective-dependency distinguishability, all four non-current clearance states plus the check-before-owner-validation ordering, the exception-allowed-is-not-clearance assertion, the objective-integration rung and its fall-through when an owned Task is still incomplete, the ready-rung's deterministic lowest-ID Task pick and its Objective-with-no-Tasks fallback, both the empty-project and nothing-left-to-do routes to plan-next-Objective, the unresolved-selection-still-yields-an-action and no-selection-means-no-diagnostic pair, an in-memory-only fixture proving no filesystem read, and the import-boundary test.
- `AGENTS.md` — extended the `internal/data/` Codebase Map row with the Next projection responsibility (ARCH-04).

**Quality gates:** `go build ./...` (clean), `gofmt -l internal/data/next.go internal/data/next_test.go` (clean), `go vet ./internal/data/...` (clean), `go test ./internal/data/... -run 'TestResolveNext|TestNext_' -v` (24/24 pass), `go test ./internal/data/...` (pass), `make build && make test` (pass, all packages). No `.savepoint/Health-Check.md` in this project, so the Quick check step is skipped per the build-task skill.

No drift: no new package, only an extension of the one file `E48-Detail.md`'s Components and files table already names for the projection (`internal/data/next.go`), plus the `AGENTS.md` row it explicitly requires.
