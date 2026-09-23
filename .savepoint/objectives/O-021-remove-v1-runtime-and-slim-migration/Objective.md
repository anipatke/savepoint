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
  ("schema_version 1: run `savepoint migrate`"). An interrupted apply leaves
  `schema_version: 1` and uncommitted changes, which `migrate --apply` refuses
  until the owner restores the tree with git, so no separate
  interrupted-migration state exists.
- `savepoint migrate` still converts both golden fixtures (`v1-basic`,
  `v1-history`) to byte-identical output, or the golden files are updated
  in the same change with an explained diff. Preview stays the default.
- `migrate --apply` refuses to run unless every path the plan writes or
  removes is inside a git work tree with no uncommitted or untracked
  changes. It writes the converted tree directly and states how to undo
  with git.
- The journal, per-file backup/stage/install/verify steps, journal-based
  recovery, pending-operation detection, the `NextPendingMigration` rung,
  the cutover preflight, and the platform replace primitive
  (`replace*.go`) are removed.
- The V1 template trees, this repository's retired V1 skills, and the
  pre-V2 audit-skill upgrade shim are removed with the tests that pin them.
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
- **Sequencing with active work.** O-013 edited the same V2 board and data
  files, so O-021 starts only after O-013 (done 2026-09-23; Confirmed
  Decision 4).
  O-014, O-019, and O-020 then plan against the smaller code base.
- **Task order.** Deletions come first (T-022 code, T-023 assets), so later
  Tasks never edit files that are about to disappear. The runtime schema gate
  (T-024) lands before the apply rewrite (T-025), because it is what removes
  every runtime caller of the journal. T-026 reconciles guidance and records
  the size reduction.

## Confirmed Decisions

The owner confirmed this design on 2026-09-23:

1. **Keep `savepoint migrate`, slimmed.** Conversion logic and golden
   fixtures stay in the current release line. The transaction machinery goes,
   and `migrate` leaves runtime command paths. Removing `internal/migrate`
   entirely is not planned.
2. **A clean git tree is the apply safety model.** `migrate --apply` refuses
   unless every path the plan writes or removes is inside a git work tree with
   no uncommitted or untracked changes. It then writes directly. Undo is a git
   restore. There is no non-git fallback and no backup copy. Preview
   (`--dry-run`, the default) needs no git.
3. **Extra V1 leftovers are in scope:** the V1 template trees
   (`templates/project/`, `templates/release/v1/`) and `UpgradeProjectAssets`'
   unused `v1Templates` parameter; this repository's nine retired V1 skill
   folders in `agent-skills/` plus `agent-skills/references/audit-method.md`
   and the tests that pin them; and the pre-V2 audit-skill upgrade shim
   (`internal/init/migrate_audit_skill.go`). `retire_v1_skills.go` stays,
   because migrated projects still need it on their first V2 upgrade.
   `agent-skills/bubbletea-tui-design/` stays.
4. **Sequencing:** O-021 runs after O-013 completes and before O-014, O-019,
   and O-020. It stays in R-006.

`upgrade-assets` already refuses a V1 project and points at `migrate`, so no
separate decision on it was needed.

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
entirely (Confirmed Decision 1); a non-git apply path (Confirmed
Decision 2).

## Planning Handoff

Tasks T-022 to T-026 are detailed and owned by this Objective. T-021 is
skipped because superseded O-018 planning prose already used that number for
a Task that was never created. O-013 completed on 2026-09-23 (`a6283d5`,
owner exception on C-911). O-021 carries no `depends_on` on it: the
ordering was only to avoid concurrent edits, and an exception closure does
not satisfy a `clear` Objective dependency. Once the owner
approves this plan, the router moves to `state: task`, `objective: O-021`,
`task: T-022`.
