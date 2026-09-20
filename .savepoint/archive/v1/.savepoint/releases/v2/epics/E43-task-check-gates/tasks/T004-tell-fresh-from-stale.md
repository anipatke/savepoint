---
id: E43-task-check-gates/T004-tell-fresh-from-stale
title: Tell fresh clearance from stale
status: done
objective: Resolve a record's clearance state and answer whether each Task dependency is satisfied.
depends_on:
    - E43-task-check-gates/T002-know-when-work-was-checked
complexity_tier: medium
complexity_reason: Deterministic evaluation over loaded records with no new schema or write surface.
---

# T004: Tell fresh clearance from stale

## Problem

E42 built the dependency graph structurally but never answered whether a dependency is satisfied, and nothing turns a Check plus a freshness assessment into a single clearance answer that every control can share.

## Context Files

- `internal/data/gate_v2.go`
- `internal/data/gate_v2_test.go`
- `internal/data/dependency.go`
- `internal/data/dependency_test.go`
- `internal/data/check_v2.go`
- `internal/data/evidence_v2.go`
- `internal/data/project.go`

## Acceptance Criteria

- [x] Clearance resolves to exactly one of `current`, `needs_work`, `stale`, `unknown`, or `missing`, and carries the Check and basis it was derived from.
- [x] Clearance is `current` only when the record's latest Check is CLEAR and a freshness assessment names that same Check with `state: current`.
- [x] A freshness assessment naming any other Check — including a superseded one — resolves to `stale`, not current.
- [x] No freshness assessment resolves to `unknown`; a record with no Check at all resolves to `missing`; a latest Check of NEEDS WORK resolves to `needs_work`.
- [x] Clearance is derived only from recorded evidence: no file scanning, hashing, timestamp comparison, or Git inspection decides it.
- [x] A `requires: clear` dependency is satisfied only when the dependency Task is `done` with `current` clearance.
- [x] A `requires: accepted` dependency additionally requires recorded owner acceptance of that Task's current Check.
- [x] A dependency that is `done` but not currently cleared is unsatisfied, with a reason naming the clearance state.
- [x] Every unsatisfied dependency returns a typed reason naming the target and the unmet requirement, not a bare boolean.
- [x] Evaluation returns decisions rather than errors for recoverable states, and re-resolves evidence from the loaded index on every call.

## Implementation Plan

- [x] Define the clearance state type and its resolution in a new `gate_v2.go`, reading the latest Check and freshness assessment from the index.
- [x] Define the typed reason vocabulary shared by clearance and dependency answers.
- [x] Add V2 dependency satisfaction to `dependency.go` for `clear` and `accepted`, leaving V1 reference matching untouched.
- [x] Test each clearance state, the superseded-freshness case, and acceptance bound to a superseded Check.
- [x] Test satisfied and unsatisfied `clear` and `accepted` dependencies, including done-but-not-cleared and a chain across several Tasks.
- [x] Run the focused `internal/data` suite.

## Context Log

**Files read:** `internal/data/dependency.go`, `internal/data/dependency_test.go`, `internal/data/check_v2.go`, `internal/data/evidence_v2.go`, `internal/data/task_v2.go`, `internal/data/objective_v2.go`, `internal/data/project.go`, `internal/data/errors.go`, `internal/data/check_v2_test.go` (targeted read for fixture style), `internal/data/evidence_v2_test.go` (targeted read for fixture style), `E43-Detail.md`, `T002-know-when-work-was-checked.md` (targeted read to confirm dependency status and prior design decisions), `T005-decide-when-tasks-close.md` (targeted read to confirm T005 owns start/completion/authority/exception, keeping this task scoped to clearance + dependency satisfaction only), `.savepoint/Guardrails.md`, `AGENTS.md`, `.savepoint/router.md`, `agent-skills/savepoint-build-task/SKILL.md`.

**Files edited:**
- `internal/data/gate_v2.go` (new) — `ClearanceState` (current/needs_work/stale/unknown/missing) and `Clearance` (State, Check, Freshness); `ResolveClearance(index, targetID)` re-reads `index.LatestCheck`/`index.Checks` and the target's evidence on every call (Task or Objective, via `evidenceFreshness`), deriving clearance purely from recorded Checks and freshness — no file scanning, hashing, timestamps, or Git.
- `internal/data/dependency.go` — added `DependencyBlockKind` (not_done/not_cleared/not_accepted), `DependencyBlock` (Target, Requires, Kind, Clearance — reusing `ClearanceState` so a done-but-not-cleared block names the same vocabulary `ResolveClearance` reports), `DependencyDecision` (Satisfied, Block), and `ResolveTaskDependencyV2(index, dep)` implementing `clear`/`accepted` satisfaction; V1 `ResolveDependency` untouched.
- `internal/data/gate_v2_test.go` (new) — `mustCheck` fixture helper; clearance coverage: missing (no Check), needs_work, unknown (CLEAR with no freshness), current (freshness names latest CLEAR Check as current), stale via a superseded-Check reference, stale via freshness naming the latest Check but not asserting `current`, and clearance resolving for an Objective target (not just Task) to prove the shared evidence path.
- `internal/data/dependency_test.go` — extended `newV2TestIndex` to initialize `Checks`, `ScopeChecks`, `LatestCheck`; added dependency satisfaction coverage: `clear` satisfied, not-done unsatisfied, done-but-not-cleared table over needs_work/unknown/stale/missing (each asserting the named `Clearance` on the block), `accepted` satisfied when owner accepted the current Check, `accepted` unsatisfied with no acceptance, `accepted` unsatisfied when acceptance is bound to a Check a newer one supersedes, and a three-Task chain (T003→T002→T001) proving each dependency re-resolves independently against the live index rather than a cached judgement.

**Design decisions not fully spelled out in the task/epic:**
- Given `Freshness.State` can independently be `current`, `stale`, or `unknown` regardless of which Check it names, "current" clearance requires both `freshness.Check == latest Check ID` *and* `freshness.State == current`; any other combination (different Check, or same Check with a non-current state) resolves to `stale` rather than `unknown`, satisfying the AC that "no freshness assessment resolves to unknown" once a freshness block is actually present — `unknown` is reserved for the no-freshness-at-all case.
- `ResolveClearance` looks up evidence generically across `index.Tasks` then `index.Objectives` (`evidenceFreshness`) since both share the identical `Evidence`/`Freshness` shape and E43-Detail.md's Components table names `gate_v2.go` as the "one canonical gate API... every V2 consumer will call" — this keeps E44's later Objective integration from needing a second clearance function, while E43's own boundary (Task scope only) is respected because this task only wires `ResolveTaskDependencyV2` to consume it for Tasks.
- `DependencyBlock.Clearance` is left as the zero value (`""`) for `DependencyBlockNotDone`, since "not done" is not a clearance question at all; it is populated for both `not_cleared` and `not_accepted` (the latter reports the dependency's own current clearance, which is `current` by construction since acceptance is only checked after clearance passes) so a caller can always inspect it without a type switch on `Kind` first.
- Left `gate_v2.go` free of any Task/Objective start/completion decision, actor-authority, exception, or replan logic — those are T005's explicit scope per its Problem/AC/Implementation Plan, which this task's `## Context Files` and epic Components table also reserve for a later task.

**Named cases and results:**
- `TestResolveClearance_missingWhenNoCheck`, `_needsWorkFromLatestCheck`, `_unknownWhenNoFreshness`, `_currentWhenFreshnessNamesLatestCheckAsCurrent`, `_staleWhenFreshnessNamesADifferentCheck`, `_staleWhenFreshnessNamesLatestButNotCurrentState`, `_resolvesObjectiveTargetsToo` — PASS
- `TestResolveTaskDependencyV2_clearSatisfiedWhenDoneAndCurrent`, `_notDoneUnsatisfied`, `_doneButNotClearedUnsatisfied` (table: needs_work/unknown/stale/missing), `_acceptedSatisfiedWhenOwnerAcceptedCurrentCheck`, `_acceptedUnsatisfiedWithNoAcceptance`, `_acceptedUnsatisfiedWhenAcceptanceBoundToSupersededCheck`, `_chainAcrossSeveralTasks` — PASS
- Full `go test ./internal/data/...` — PASS, no regressions in E42/T001–T003 suites
- `go vet ./internal/data/...` — clean
- `gofmt -l` on all edited/added files — clean

**Quality gates:** `make build && make test` — PASS across all packages (`.`, `cmd`, `internal/board`, `internal/buildtool`, `internal/data`, `internal/doctor`, `internal/init`, `internal/styles`).

No `.savepoint/Health-Check.md` in this project; Quick check step skipped per `savepoint-build-task`.

No Drift Notes: `gate_v2.go` is the exact new file E43-Detail.md's Components table names for clearance/decision evaluation, and `dependency.go` is the file this task's own Implementation Plan directs V2 dependency satisfaction into.
