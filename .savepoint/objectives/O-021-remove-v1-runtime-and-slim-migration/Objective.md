---
id: O-021
title: Remove the dead V1 runtime and slim the migration engine
status: planned
release: R-006
---

# O-021: Remove the dead V1 runtime and slim the migration engine

## Outcome

Savepoint's code base contains only what the V2 runtime and a one-time V1-to-V2
conversion actually need. V1 board, doctor, and data code that no command can
reach is deleted with its tests. Everyday commands (`board`, `resume`,
`doctor`, `init`, `upgrade-assets`) no longer depend on the migration package;
they recognise a V1 project with a plain schema-version check and point the
user at `savepoint migrate`. The migration engine keeps its conversion logic
and golden tests but drops the transaction machinery that duplicates what git
already provides.

## Why

An independent review on 2026-09-23 measured about 33.4k lines of production Go
and 55.1k lines of tests. Roughly 37% of the production code is either
migration code or V1 code the shipped binary can no longer reach:

- `deadcode` against the real `savepoint` binary reports **6,350 unreachable
  production lines**: ~3,750 in `internal/board` (the whole V1 board, which
  `runWithFilters` refuses to start for a V1 project), ~1,560 in
  `internal/data` (V1 audit register/run loaders and validators, most of
  `lifecycle.go`, V1 write helpers, plus unused V2 writers such as
  `CreateCheckV2`, `CreateIssueV2`, `WriteIssueV2`, `WriteReleaseV2*`), and
  ~950 in `internal/doctor` (V1 checks). Roughly 10–13k lines of tests exist
  only to keep that code working.
- `internal/migrate` is 5.9k production lines and 8.0k test lines. About 2k
  lines are the actual conversion. The rest is a write-ahead journal
  (backup → stage → install → verify per file), crash recovery from that
  journal, source-hash revalidation, create-only adoption, and a Windows
  file-replace primitive. That is database-grade safety for a one-off
  markdown conversion that almost always runs inside a git repository.
- Migration leaks into the runtime. `board`, `board/v2`, `doctor`, `init`,
  and `main` import `migrate`. Every `resume` and board startup calls
  `PreflightCutover` and `PendingOperation`, and on a V1 project
  `PreflightCutover` builds a full conversion plan only to print "migrate
  first".
- It keeps costing work. Since the cutover, migration code and fixtures have
  taken ~2.5k lines of changes. O-018 alone touched ten migration files and
  both golden outputs. Every V2 schema change has to be carried back into
  the converter.

Removing this reduces build and test time, makes the codebase easier to read
for both agents and the owner, and stops future Objectives paying a migration
tax.

## Success Conditions

- `deadcode .` against the `savepoint` main package reports no unreachable
  functions in `internal/board`, `internal/doctor`, or `internal/data`, apart
  from any the owner explicitly keeps and lists in this Objective.
- The V1 board package code is gone. The `board` command's schema dispatch
  (the V1 filter refusal and the V2 hand-off) lives in the smallest package
  that needs it, with no V1 board types left.
- `internal/board`, `internal/board/v2`, `internal/doctor`, `internal/init`,
  and `main.go`'s `board`/`resume` paths do not import `internal/migrate`.
  A shared schema-version check in `internal/data` (building on
  `data.ReadSchemaVersion`) returns one named diagnostic for a V1 project
  ("schema_version 1: run `savepoint migrate`") and for a project with an
  interrupted migration, if interruption remains a reachable state.
- `savepoint migrate` still converts both golden fixtures (`v1-basic`,
  `v1-history`) to byte-identical output, or the golden files are updated
  in the same change with an explained diff. Preview stays the default.
- `migrate --apply` refuses to run unless the target is inside a git work
  tree with no uncommitted changes under `.savepoint/` (or the owner-chosen
  alternative below). It writes the converted tree directly and states
  how to undo (`git checkout -- .savepoint && git clean -fd .savepoint`).
- The journal, per-file backup/stage/install/verify steps, journal-based
  recovery, and `replace_windows.go` are removed, unless the owner decides
  to keep any of them (see Owner Decisions).
- The V1 `data` readers that migration still needs (discovery, task/defect/
  epic parsing, audit finding loading) stay and are covered by migration
  tests. Nothing else in `data` is V1-only.
- AGENTS.md's Codebase Map and `Design.md` describe the smaller shape. Each
  map row is one or two sentences. `make build && make test-full` passes,
  and the Objective reports the before/after production and test line
  counts per package.

## Architectural Considerations

- **Delete, don't deprecate.** Unreachable code goes, with its tests. Git
  history is the archive; no `legacy_` files or build tags.
- **One schema gate.** `internal/data` owns "which schema is this project".
  Runtime commands consume that result. Only `savepoint migrate` imports
  `internal/migrate`.
- **Conversion stays pure.** The plan, convert, and manifest code keep their
  current role: read V1 sources, produce V2 bytes and a manifest. Only the
  apply step changes.
- **Golden fixtures are the contract.** They are the proof that slimming did
  not change conversion output.
- **Sequencing with active work.** O-013 (Release → Goals wording) is in
  progress and edits the V2 board and data. Phase 1 (deleting V1 board,
  doctor, and data code) mostly touches different files and can run first.
  Phases 2–3 touch `data`, `board/v2` load, and `doctor` startup and should
  start after O-013 lands. O-020's board writes must not build on anything
  removed here.
- **Suggested Task split** (detailed when this Objective becomes next):
  1. Delete the unreachable V1 board, V1 doctor checks, unused V1 `data`
     functions, and unused V2 writers, with their tests.
  2. Replace `PreflightCutover`/`PendingOperation` in runtime commands with
     the `data` schema gate, and drop `migrate` imports outside `main`'s
     migrate dispatch.
  3. Replace the journal-based apply with a clean-git precondition and
     direct writes. Delete the journal, recovery, and platform replace code.
     Keep the golden and conversion tests.
  4. Reconcile AGENTS.md, Design.md, README "Safe migration", and scaffold
     guidance. Record the line-count deltas.

## Owner Decisions Required Before Task Detailing

1. **V1 support policy.** Keep `savepoint migrate` in the current release
   line, or freeze it at a tagged version ("install vX.Y, migrate, then
   upgrade") and delete `internal/migrate` entirely in a later Objective?
   This Objective assumes it is kept, but slimmed.
2. **Apply safety model.** A clean git tree as a hard precondition (recommended).
   Or: allow non-git projects by writing a single timestamped
   `.savepoint.bak/` copy before converting, which is still far smaller than
   the journal.
3. **V1 `upgrade-assets`.** `upgrade-assets` still dispatches to the V1
   template tree for schema-1 projects, and `init` carries
   `retire_v1_skills.go` and `migrate_audit_skill.go`. Keep these for V1
   users who haven't migrated, or refuse and point them at `migrate`?
4. **Release membership.** Assigned to R-006 alongside the other open V2
   Objectives. Leave it there, or make it unassigned.

## Boundaries

**In scope:** Deleting unreachable V1 and unused V2 code and its tests;
removing `migrate` from runtime command paths; a `data`-owned schema gate;
simplifying migration apply; keeping conversion behaviour proven by golden
fixtures; documentation and Codebase Map reconciliation; before/after size
reporting.

**Out of scope:** Changing V2 record formats, gates, or `ResolveNext`;
V2 board features; rewriting archived V1 records or immutable Checks;
editing `.savepoint/archive/v1/`; broad test-suite trimming beyond tests of
deleted code (a separate Objective can review the 5:1 test ratio in
`internal/init` and wording-contract tests); removing `internal/migrate`
entirely (see Owner Decision 1).

## Planning Handoff

This is a proposed Objective, not the current router selection. The router
is unchanged. Settle the four Owner Decisions above, then detail Tasks against
the code as it stands after O-013 lands.
