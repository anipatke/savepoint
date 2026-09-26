---
id: I-067
title: Migrate reports some owner decisions resolved but never applies them
type: defect
status: resolved
source:
  kind: report
  actor: {role: owner, session: user}
  at: '2026-09-26T00:15:19Z'
severity: blocker
history:
  - at: '2026-09-26T00:15:19Z'
    actor: {role: owner, session: user}
    kind: observed
    note: >-
      Migrate preview of the owner's galaxy project with a decisions file
      showed v2 epics E10 and E11 (status "deferred") as "resolved -> planned
      (from decisions.yml)", yet planned no Objective for either and did not
      archive their detail files.
  - at: '2026-09-26T00:19:26Z'
    actor: {role: executor, session: i067-repair-20260926}
    kind: repair_attempted
    note: >-
      internal/migrate: PlannedTarget gains DecidedStatus, persisted with the
      plan. planEpic and task planning read an unrecognized_lifecycle
      decision (decidedLifecycle) and plan the record with the decided
      status; ConvertObjective, ConvertTask, and ConvertIssue use
      DecidedStatus instead of re-resolving the V1 value. An unresolved
      narrative duplicate finding is archived (archive) or converted as an
      open Issue (treat_as_original). A duplicate task identity archives the
      second declaration (keep_first) or archives the first and plans the
      second under the first's identity (keep_second). The finding
      unrecognized_lifecycle site is unreachable because the V1 finding
      parser already heals an unknown status to open, so it is unchanged.
      decisions_applied_test.go covers epic, task (planned and done), both
      narrative choices, and both identity choices, and asserts every
      resolved blocking ambiguity leaves its source planned or archived; the
      epic and task tests fail with the decision lookup disabled, with the
      galaxy symptom. The galaxy preview now plans v2 E10 and E11 as O-006
      and O-007. make build, make test-fast, and make test-full passed; one
      earlier test-full run failed TestV2WatcherDebouncesRapidWrites under
      load, which passed 20 of 20 in isolation and on the rerun. Issue
      remains open for verification or owner acceptance.
  - at: '2026-09-26T03:34:53Z'
    actor: {role: owner, session: user}
    kind: owner_decision
    note: >-
      Owner instructed the executor to mark this Issue resolved. The fix
      shipped in savepoint 2.0.5 on npm (tag v2.0.5), and in real use the galaxy migration planned the decided epics E10 and E11 as O-006 and O-007.
      No technical CLEAR is implied.
resolution:
  disposition: accepted
  actor: {role: owner, session: user}
  at: '2026-09-26T03:34:53Z'
  reason: Owner accepted the fix as deployed in savepoint 2.0.5.
---

# I-067: Migrate reports some owner decisions resolved but never applies them

## Summary

`addAmbiguityDiscriminated` marks an ambiguity resolved whenever the
decisions file names it, and a resolved blocking ambiguity no longer blocks
apply. But several planner sites record the ambiguity and then skip the
record without ever reading the decision:

- unrecognised epic status (`planEpic`): the epic, its Objective, and its
  Tasks are neither planned nor archived;
- unrecognised task status: the task is neither planned nor archived;
- unrecognised finding status: the finding is neither converted nor archived;
- a duplicate finding whose `duplicate_of` names no finding
  (`treat_as_original` / `archive`): the finding is skipped;
- a duplicate task identity (`keep_first` / `keep_second`): the second
  file is always skipped, whatever was chosen.

Apply then proceeds, and those V1 files stay at their old paths, unconverted
and unarchived, with no warning. Nothing is deleted, but the records vanish
from the V2 project. Release lifecycle, release completion, and duplicate
release source decisions are applied correctly; missing dependency
(`drop_dependency`) matches what the planner already does.

## Evidence

- `internal/migrate/decisions.go:215` sets `Resolved` from the decisions map
  alone.
- `internal/migrate/plan.go:965`, `:1060`, `:1345`, `:1308`, and `:1068`
  record the ambiguity and `return`/`continue` without reading
  `b.decisions`.
- galaxy preview with `decisions.yml`: E10/E11 detail files appear in neither
  Planned records nor Archives.

## Proof Needed

- An unrecognised epic, task, or finding status with a decision is planned
  exactly as if its source carried the decided status.
- An unresolved duplicate finding with `treat_as_original` converts to an
  Issue; with `archive` it is archived.
- A duplicate task identity decision is either applied or refused, never
  reported resolved and ignored.
- Every blocking ambiguity reported resolved has its decision applied; a
  test fails if a kind is marked resolved without a consumer.
- Tests cover each case; `make build && make test-fast` and
  `make test-full` pass.
