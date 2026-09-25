---
id: O-025
title: Remove the Goal Check
status: planned
release: R-006
priority: critical
rank: 1
---

# O-025: Remove the Goal Check

## Outcome

Savepoint has two kinds of Check: the mandatory Full Objective Check and the
optional Task Check. A Goal is complete when every member Objective is
complete. Resume, the board, doctor, Design, Guardrails, AGENTS.md, the
skills, and the new-project scaffold no longer ask for or describe a Goal
Check.

## Why

The owner decided on 2026-09-25 that there is no Goal Check (I-060). A Full
Goal Check repeats a full review of every Objective, each of which already
passed its own mandatory Objective Check, and it is too expensive to run.
R-006 has all 13 Objectives done, yet resume says `Record a fresh Goal Check
for R-006.` and doctor exits 1 with `v2-release-clearance-missing`.

## Success Conditions

- `ResolveReleaseCompletion` allows completion when the Goal has member
  Objectives and every one passes `ResolveObjectiveCompletion`. It no longer
  reads Goal-scoped clearance, Goal-level Issues linked to a Goal Check, a
  Goal exception, or Goal owner acceptance. A migrated Goal's legacy
  completion keeps working as it does today.
- When every member Objective is complete, Next for the selected Goal is the
  existing Goal-ready action. The Goal-Check-needed and Goal-owner-acceptance
  Next kinds and their resume phrases are gone.
- Doctor reports no Goal clearance, Goal Check, or Goal owner-acceptance
  problem. It still reports a `done` Goal whose member Objectives are
  incomplete, and a Goal with no Objectives, as today.
- The board's Goal detail no longer shows Goal Check clearance or an owner
  acceptance section for Goal completion.
- Check records with `scope.kind: release` still load. No completion
  decision reads them, and no guidance tells anyone to write one.
- Design, Guardrails (including TEST-08 and TEST-09), AGENTS.md,
  `.savepoint/router.md`, README, every skill and shared reference, and
  their `templates/project-v2` copies describe only the Objective Check and
  the optional Task Check. The template freshness and agent-skill tests pass.
- On this repository, `savepoint resume` no longer asks for a Goal Check for
  R-006, and `savepoint doctor` reports no Goal Check problem.
- `make build && make test-fast` pass at each Task handoff. The Full
  Objective Check has current `make test-full` evidence.

## Architectural Considerations

- `internal/data` stays the single owner of the completion decision
  (ARCH-01, DATA-02). Resume, the board, and doctor only render what
  `ResolveReleaseCompletion` and `ResolveNext` return; none of them
  re-derives Goal completion.
- `CheckScopeRelease` stays a valid decoded scope for compatibility.
  Storage names (`R-###`, `release:`, `scope.kind: release`) do not change.
- Goal `status: done` is still the owner's own record change. This Objective
  adds no board action to close a Goal.
- Past Check, Issue, Objective, and Task records that mention a Goal Check
  are history and are not edited.

## Boundaries

**In scope:** the Goal completion gate and its tests; Next kinds and the
resume copy; the board Goal detail; doctor's Goal diagnostics; Design,
Guardrails, AGENTS.md, router, README, skills, shared references, and their
scaffold copies; migration goldens only where their expected text changes.

**Out of scope:** renaming Release storage, a board action to close a Goal,
changing the Objective or Task Check rules, rejecting existing
`scope.kind: release` records, and editing historical records.

## Confirmed Design Decisions

Confirmed by the owner on 2026-09-25 in chat.

- **No Goal Check:** "No such thing as goal check. Mandatory Obj check and
  optional task check only."
- **Goal completion:** complete when all member Objectives are complete. No
  Goal-level Check and no Goal-level owner acceptance.
- **Existing Goal-scoped Checks:** keep loading, ignored by every decision.
- **Task split:** one runtime Task, then one guidance Task.
