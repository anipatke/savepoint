---
type: audit-findings
audited: 2026-09-18
applied: 2026-09-18
---

# Audit Findings: E45 Convert active work without losing history

## Main Findings

### Verdict

**CLEAR — fixes applied, epic marked audited.** The original pass found four
findings, one a Blocker: on a real project, `savepoint migrate --apply` could
silently destroy a user's file and report success. All four are now resolved —
findings 1 and 2 with code changes plus new regression tests, findings 3 and 4
by narrowing the documentation to the guarantee the code actually makes, per
the owner's-call option this report itself offered. `go build ./...`, `go vet
./...`, and `make build && make test` are clean against the applied tree.
Everything else the epic promised was already genuinely built and proven in
the original pass — the preview really writes nothing, the conversion really
is deterministic, nothing fabricates evidence, the archives really are
byte-exact, and the three other write paths really do stand down. The
Non-Blocking Observations below were not required for this verdict and remain
open as recorded.

### What Needs Attention (as originally found — see resolution note on each)

**1. RESOLVED. Migration writes its archive outside Savepoint's own directory, and
overwrites whatever is already there.**

*Resolution:* `archivePathFor` (`internal/migrate/plan.go`) now targets
`.savepoint/archive/v1/...`, `Inventory` excludes `.savepoint/archive/` from
its own walk (`archiveDirName`, `internal/migrate/inventory.go`), and the
golden fixtures and end-to-end/apply assertions were regenerated and updated
to the new path. Covered by `TestPlan_archivePathLivesInsideSavepoint` and
`TestInventory_excludesArchiveDir`.

The architecture record and the epic both place the byte-preserved archive at
`.savepoint/archive/v1/`, alongside `objectives/`, `issues/`, and `migrations/`.
The shipped code writes it to `archive/v1/` at the *project root* instead — a new
top-level directory in the user's repository. Because that location sits outside
the tree migration inventories, a file already living there is invisible to the
whole safety apparatus: it is not hashed, not reported as a conflict, not named
in the preview, and not backed up. Apply simply overwrites it and exits 0. The
user's original bytes are then gone from the project entirely — there is no
archive copy, no backup, and the operation directory that might have held one is
deleted on success. This is the single thing this epic exists to prevent, and it
violates the Blocker rule on overwriting user-authored content.

The two problems are one problem: putting the archive back inside `.savepoint/`
brings it under inventory, confinement, and conflict detection, which is why the
design put it there. Fixing the path alone is not enough — see finding 2 for the
collision behaviour that then takes over — so both need to land together.

*Evidence:* place any file at `<project>/archive/v1/.savepoint/PRD.md`, run
`savepoint migrate <project> --apply`; it exits 0, the preview never mentions the
collision, and the file is replaced with the V1 PRD with no copy left anywhere in
the project. Path built at `internal/migrate/plan.go:1085` (`archivePathFor`);
intended location at `.savepoint/releases/v2/v2-Design.md:46` and
`E45-Detail.md:22`.

**2. RESOLVED. A file sitting at a path migration plans to create bricks the migration,
permanently.**

*Resolution:* `Plan` now runs `checkDestinationCollisions` (`internal/migrate/plan.go`)
against every create-kind destination (converted records, `Idea.md`, archive
entries) and reports a named `ConflictDestinationExists` at plan time — before
any operation directory is created — rather than failing apply with an
internal invariant message and no way out. Covered by
`TestPlan_preExistingIdeaDestination_isNamedConflict`,
`TestPlan_preExistingArchiveDestination_isNamedConflict`,
`TestPlan_replacedRouterDocument_isNeverACollision` (proving the replaced
`router.md` path is deliberately exempt), and
`TestApply_preExistingCreateDestination_refusesCleanly`.

If any file already exists exactly where the plan intends to create one — the
most likely real case being a `.savepoint/Idea.md` written by a user who read the
V2 documentation and started early — apply stops with an internal invariant
message (`... has action create, which has no backup step`) rather than a named
diagnostic. It fails safe in the sense that nothing is overwritten, but it leaves
an operation directory behind, and there is no way out: `--apply`, `--recover`,
and `--recover --apply` all reproduce the same failure indefinitely. Meanwhile
the guards this epic added are doing their job — `upgrade-assets`, the board's
status and router writes, and doctor's repairs all refuse while that operation
exists. One stray file therefore takes the whole tool offline with no documented
escape short of manually deleting `.savepoint/.migration/`. The underlying cause
is that the plan happily assigns one path two fates at once (archive it *and*
create over it); the recovery report shows the path listed twice.

*Evidence:* `printf 'my idea' > <project>/.savepoint/Idea.md` on a copy of the
`v1-basic` fixture, then `savepoint migrate <project> --apply` → exit 1 with the
message above; every subsequent run repeats it. Contradicts T009 AC7, T011 AC4,
and guardrail FS-06.

**3. RESOLVED (by narrowing the claim, the owner's-call option this report offered).
"Your preview is stale" is documented but never actually happens.**

*Resolution:* the epic did not gain cross-invocation preview staleness
detection — that would mean persisting a previewed plan's hashes between two
separate `migrate` command invocations, a real feature this epic never built.
Instead, `E45-Detail.md`'s freshness-basis prose and its determinism/freshness
quality gates were rewritten to state precisely what the code guarantees: a
real, byte-level revalidation within one `Plan`-then-`Apply` call and across a
resume, and an explicit statement that two separately invoked commands (a
`migrate` preview followed later by a separate `migrate --apply`) are not
covered. The claim no longer overstates the implementation.

The epic states that apply re-hashes every source and refuses when anything
changed since the preview, and names this as a quality gate. That protection is
real inside the package — `Apply` does revalidate against the plan it is handed,
and `internal/migrate/apply_test.go` proves it — and it is real when resuming an
interrupted operation. It does not exist for the command a user actually types,
because `savepoint migrate` keeps nothing between runs: the second invocation
re-plans from scratch, so the hashes always match themselves. A user who
previews, goes away, and comes back to apply after something changed gets a
migration of the new state with no warning that the plan they approved no longer
describes the project.

*Evidence:* preview the `v1-basic` fixture, append a line to its `router.md`,
then `--apply` → completes successfully with no conflict. Contradicts
`E45-Detail.md:18`, `:59`, the "Freshness tests" quality gate, and T009 AC2.

**4. RESOLVED. Two acceptance statements describe the frozen fixtures incorrectly.**

T012's Context Log flags both of these itself, and I verified both independently:
the notes are accurate and the substitutions are sound. They are recorded here
only because the *source* documents still carry the wrong text, and a future
re-audit reading them cold would measure the epic against criteria it cannot
meet.

First, the epic's own determinism gate and T012 AC7 require `v1-history`'s
recurring `E01-example/T001-shared` to become "two distinct Tasks with distinct
global IDs". It cannot, and the epic is right that it cannot: the `v1` copy is
`status: done` under a `status: done` epic, and the epic's fixed rule archives
completed work without giving it a V2 identity. The substance is properly proven
in two pieces — one test shows both recurrences get their own recorded
destination on the frozen fixture, the other shows two distinct global IDs when
both recurrences are active, on a temporary copy. The fixtures were not touched.
Second, T012 AC12 says `v1-basic` has "no `AGENTS.md`"; it does (role
`managed-guide`), and `v1-history` is the fixture without one. The test was
correctly written against each fixture's own `absent_by_design` list, which
between them covers all three absences the criterion meant to name.

*Evidence:* `internal/data/testdata/migration/v1-history/.../v1/epics/E01-example/tasks/T001-shared.md`
(`status: done`) and `.../E01-Detail.md` (`status: done`);
`internal/data/testdata/migration/v1-basic/manifest.yml:22` (`project/AGENTS.md`)
versus `.../v1-history/manifest.yml:135` (`absent_by_design: project/AGENTS.md`).

### Materiality Summary

| Finding | Likelihood | Impact | Materiality | Recommendation | Outcome |
|---|---|---|---|---|---|
| 1. Archive written outside `.savepoint/`, overwrites user files unrecoverably | Medium — needs a name collision under a root `archive/`, but migration runs on real repositories that often have one | High — irreversible loss of user-authored content, the exact harm the epic exists to prevent (FS-01 Blocker) | High | Fix now, together with finding 2 | **Fixed** — archive relocated under `.savepoint/`, inventoried, tested |
| 2. Pre-existing file at a planned create path bricks migration with no exit | Medium — a hand-started `.savepoint/Idea.md` is a plausible user state | Medium — no data loss, but migration and three other write paths are stuck until manual cleanup, with an internal error message | Medium | Fix now: detect the collision at plan time as a named conflict | **Fixed** — named `ConflictDestinationExists` at plan time, tested |
| 3. Stale-preview refusal absent at the command level | Medium — any edit in the window between two runs | Medium — a safety property the epic advertises does not hold for the shipped command; wrong state migrated without warning | Medium | Fix now or correct the claim; owner's call which | **Resolved by narrowing the claim** — docs now scope the guarantee to one invocation and resume |
| 4. Two acceptance statements factually wrong about the fixtures | High — already present in the files | Low — disclosed and substantively proven; documentation only | Low | Apply the text corrections below | **Fixed** — text corrected in `E45-Detail.md` and `T012` |

### What Is Proven / Not Proven

**Proven, with independent evidence beyond the task tests.** Preview is genuinely
write-free — I snapshotted bytes and modification times across the default,
`--dry-run`, and `--apply --dry-run` forms over both fixtures and nothing moved.
Both fixtures migrate end to end and load through `LoadV2Index` clean. Nothing
fabricates evidence: no Check records, no `last_check`, no `freshness`, no
accepted validation anywhere in either migrated project. Identity is genuinely
source-qualified and the manifest resolves every legacy fact, including the typed
legacy prerequisite onto an archived completed Task. A second run reports
"already declares schema_version: 2" and changes no content and no modification
time. Missing, non-Savepoint, and unwritable targets each produce their own named
error; an unknown flag prints usage. Two incomplete operations produce a named
diagnostic, and `upgrade-assets` and `doctor` both refuse and name the operation.
The Windows replacement decision in T001 is recorded with observed evidence, the
alternatives are rejected on observation rather than prediction, and no
dependency was added.

**Not proven.** No fixture carries a `verified` finding, so that branch of the
archival rule is asserted but never exercised (T012 records this). All
behavioural evidence ran on Linux; Windows evidence exists only for the
replacement primitive itself, from T001's recorded NTFS run, which cannot be
re-executed from this session. Findings 1 and 2's fixes are proven by the new unit tests named above and, for
finding 2, independently reproduced against the compiled binary using this
report's own repro command (`printf 'my idea' > <project>/.savepoint/Idea.md`
on a copy of `v1-basic`, then `--apply`): it now exits 1 with the named
`destination_exists_conflict` diagnostic, creates no `.migration/` operation
directory, and the preview's Archives section confirms finding 1's
`.savepoint/archive/v1/...` destination in the same run.

### Audit Evidence

- **Scope lock:** E45-Detail's twelve quality gates and Boundaries, the twelve
  task files' acceptance criteria, guardrails FS-01/03/04/05/06,
  DATA-01/03, ARCH-01/04, CFG-02, TEST-03/04, and the `STYLE` set; the 62 files
  changed since `117518b` (the commit before T001 — note the brief named
  `1cdc869` as the pre-epic base, which is actually T010; the real base is
  `b2e9023`/`117518b`). Public surfaces: the `migrate` command and its five flag
  forms, and `internal/migrate`'s exported `Plan`, `Apply`, `ReplaceFile`,
  `PendingOperation`, `ReadDecisionsFile`, `FormatPreview`, `RunCommand`.
- **Coverage and workflow result:** the publish order was traced in
  `internal/migrate/apply.go` against its stated order and matches — backups,
  additive creates, archive, manifest, in-place replaces, removals, then schema
  activation last, with every step idempotent over recorded journal state. The
  independent oracle throughout was the compiled binary driven against
  disposable copies of the frozen fixtures, never the package's own test
  helpers. Cells classified as findings: the archive destination, the
  create-path collision, and cross-invocation freshness. Not-applicable:
  concurrency and multi-user cells (single-user local tool), and network cells
  (no network boundary exists).
- **File reality and drift:** every file named in all twelve Context Logs exists,
  except T001's cross-compiled Windows probe binary and its `C:` scratch
  directory, both explicitly recorded as removed with reproduction commands —
  acceptable. The one Drift Note of consequence (T012 extracting the command
  body to `internal/migrate/command.go`) is justified, behaviour-preserving, and
  was reconciled into the AGENTS.md Codebase Map. The archive location is
  undeclared drift from `v2-Design.md` that no task recorded — that is finding 1.
- **Gates:** `go build ./...`, `go vet ./...`, `go test ./... -count=1`, and
  `make build && make test` all pass. `git diff --check` clean. `gofmt -l`
  reports the same five pre-existing files before and after the epic
  (`internal/board/release.go`, `internal/data/router.go`,
  `internal/data/task_test.go`, `internal/init/clipboard.go`,
  `internal/init/manifest_test.go`) — verified by running it against a checkout
  of the pre-epic commit, so E45 introduced no new formatting drift.

### Guardrails Verification

- **Rule IDs checked:** FS-01, FS-03, FS-04, FS-05, FS-06, DATA-01, DATA-03,
  ARCH-01, ARCH-04, CFG-02, DEP-01, TEST-02, TEST-03, TEST-04, TEST-08, REL-01,
  TPL-01, STYLE-01..10.
- **Health check mode:** Full. This project has no `.savepoint/Health-Check.md`
  and no `.savepoint/audit/` register, so those steps are skipped rather than
  reported, per AGENTS.md.
- **Evidence:** FS-01 **fails** — finding 1. FS-06 **fails** — finding 2 leaves a
  recorded operation with no clear path out and an unnamed internal error.
  FS-03, FS-04, FS-05, DATA-01, DATA-03, ARCH-01 (`cmd/migrate.go` is 79 lines
  importing only `context`, `fmt`, `io`), ARCH-04 (Codebase Map row present and
  current), DEP-01 (`go.mod`/`go.sum` unchanged), TEST-04, TEST-08, TPL-01 all
  pass; see Audit Evidence above rather than repeated here. CFG-02 is satisfied
  for the replacement primitive by T001's recorded Windows/Linux contract and
  its platform-tagged tests; it could not be re-run from this Linux session.
  REL-01 rests on T001's recorded `GOOS=windows go build ./...`.
- **File reality evidence:** as recorded under Audit Evidence.
- **Waivers or unresolved findings:** no waivers requested. Findings 1–4 open.

### Non-Blocking Observations

- **Converted active work is left behind at its V1 path.** After migrating
  `v1-history`, the project holds both the new `objectives/.../T001-shared.md`
  and the untouched `.savepoint/releases/v1.1/.../T001-shared.md`, likewise for
  defects and findings — while their epic details, audit register, and prompt
  *were* archived and removed. The result is a half-gutted V1 tree containing
  duplicate copies of exactly the active work, which can silently diverge if
  edited. This matches `v2-Design.md`'s mapping table, which promises "original
  bytes archived" only for release PRDs and epic details, so it is not a rule
  violation — but nothing documents it and the preview never says a source will
  be left in place. Worth an explicit owner decision before V2 ships.
- **`WriteManifestCreateOnly` in `internal/migrate/manifest.go` is unreachable
  production code.** Only its own test calls it. The real create-only guarantee
  comes from the plan-time `ConflictExistingManifest` check, which I confirmed
  works; the manifest itself is published through the ordinary journal path. The
  test therefore proves a property of a function the product never uses.
- **`upgrade-assets` surfaces the migration refusal as a per-file `failed` row**
  inside a normal report rather than a single up-front refusal, so the output
  reads "Failed: 1 / Skipped: 2" before the real reason. It does name the
  operation, which is what T010 required.
- **The duplicated build episode left no code drift.** Commit `425b8b3`, labelled
  "Merge E45 T012", in fact merged a second independent build of T011 and is
  content-identical to its parent `f5d7bcc`. README, the T011 task file, and all
  sources are single-copy; only the commit subject is misleading.
- **T001's Drift Notes asked this audit to rule on an out-of-scope issue.** The
  Windows experiment established that `internal/data`'s `replaceV2File`
  (`os.Rename`) fails against a destination another process holds open and
  discards its attributes and ACL, and that `internal/init`'s `AtomicWrite`
  falls back to a truncating copy in precisely that case. Both are real and both
  are outside E45's boundary. Recommend capturing them as a release-level defect
  rather than reopening this epic.

## Code Style Review

- [x] STYLE-01 **One job per file** — split files when responsibilities mix. The
      package splits cleanly across inventory, classify, plan, convert,
      convert_docs, convert_issues, decisions, manifest, operation, apply,
      replace, preview, and command. `plan.go` at 1233 lines is by far the
      largest and carries planning, identity allocation, ambiguity detection,
      and document planning; coherent, but the natural next split.
- [x] STYLE-02 **One job per function** — small, named, testable units.
- [x] STYLE-03 **Test branches** — cover meaningful conditionals and edge cases.
      Both collision branches behind findings 1 and 2 now have tests
      (`TestPlan_preExisting*DestinationIsNamedConflict`,
      `TestPlan_replacedRouterDocument_isNeverACollision`,
      `TestApply_preExistingCreateDestination_refusesCleanly`,
      `TestInventory_excludesArchiveDir`).
- [x] STYLE-04 **Types document intent** — prefer explicit types over comments.
      `EntryAction`, `StepState`, `Role`, `AmbiguityKind`, `DocumentKind`, and
      `LegacyKey` all carry meaning in the type rather than in prose.
- [ ] STYLE-05 **Build only what is needed** — no speculative abstractions.
      `WriteManifestCreateOnly` is exported, tested, and called by nothing.
- [x] STYLE-06 **Handle errors at boundaries** — validate inputs, APIs, IO, and
      external data. The create-path collision (finding 2) now surfaces as a
      named `ConflictDestinationExists` at plan time instead of an internal
      invariant message.
- [x] STYLE-07 **One source of truth** — no duplicated rules, constants, state,
      or config. Notably well handled: T012 extracted the command body rather
      than let the test re-implement its branching, and `migrationStateDir` and
      `archivePathFor` are each defined once.
- [x] STYLE-08 **Comments explain why** — not what the code already says. The doc
      comments on `activateSchema`, `publishWrite`, and `revalidateSources`
      explain the safety reasoning rather than restating the code.
- [x] STYLE-09 **Content lives in data** — keep copy/config out of logic. The
      finding-disposition rules, the role vocabulary, and the ambiguity choice
      lists are all tables.
- [x] STYLE-10 **Small diffs** — minimal, reviewable, behaviour-preserving
      changes. Each task's diff stayed within its declared surface.

`STYLE` rules are Guideline severity throughout and never block on their own.
STYLE-05 remains unchecked as a standing, non-blocking observation
(`WriteManifestCreateOnly` is still unreachable production code — see
Non-Blocking Observations); it is unrelated to findings 1–4 and was not part
of this apply.

## Proposed Changes — Applied 2026-09-18

All proposed changes below landed as described, plus the additional work each
turned out to need in practice: finding 2's plan-time conflict (per the
option this report offered) rather than a `--recover`-abandon path, and finding
3's documentation narrowing (the other option this report offered, chosen over
persisting previewed hashes between separate command invocations). The
`archivePathFor` move to `.savepoint/archive/v1/` also required — beyond what
the block below shows — excluding `.savepoint/archive/` from `Inventory`
(`internal/migrate/inventory.go`), regenerating both golden fixtures via
`SAVEPOINT_REGENERATE_MIGRATION_GOLDEN=1`, and updating the archive-path
assertions in `end_to_end_test.go` and `apply_test.go`, exactly as flagged
below. `go build ./...`, `go vet ./...`, and `make build && make test` are
clean against the applied tree; finding 2's fix was additionally verified
against the compiled binary (see "What Is Proven / Not Proven" above).

### Target File
internal/migrate/plan.go

### Replace
```go
// archivePathFor is the deterministic archive/v1/ destination for a
// project-relative V1 source path: the same relative layout, moved under
// archive/v1/, so an archived path is always recoverable by inspection alone.
func archivePathFor(sourcePath string) string {
	return filepath.ToSlash(filepath.Join("archive", "v1", sourcePath))
}
```

### With
```go
// archivePathFor is the deterministic .savepoint/archive/v1/ destination for a
// project-relative V1 source path: the same relative layout, moved under
// .savepoint/archive/v1/, so an archived path is always recoverable by
// inspection alone. It lives inside .savepoint/ — as v2-Design.md section 2
// specifies — so that the destination is inventoried, path-confined, and
// collision-checked like every other path the operation writes. A destination
// outside .savepoint/ is invisible to the inventory and therefore to every
// conflict and backup guarantee this operation makes.
func archivePathFor(sourcePath string) string {
	return filepath.ToSlash(filepath.Join(".savepoint", "archive", "v1", sourcePath))
}
```

This one-line move also requires regenerating
`internal/migrate/testdata/golden/v1-basic.yml` and
`v1-history.yml` through the documented
`SAVEPOINT_REGENERATE_MIGRATION_GOLDEN=1` procedure, and updating the archive
paths asserted in `internal/migrate/end_to_end_test.go` and `apply_test.go`. It
must land together with finding 2's collision guard, because moving the archive
inside `.savepoint/` brings a pre-existing archive file under the inventory and
into exactly the two-fates-for-one-path state that currently bricks the
migration.

### Target File
.savepoint/releases/v2/epics/E45-safe-migration/E45-Detail.md

### Replace
```md
- Determinism tests prove that with an injected clock and operation ID, each fixture converts to byte-identical target content and identical allocated IDs across repeated runs, and that ID allocation is source-qualified — `v1-history`'s recurring `E01-example/T001-shared` across two releases yields two distinct global Tasks.
```

### With
```md
- Determinism tests prove that with an injected clock and operation ID, each fixture converts to byte-identical target content and identical allocated IDs across repeated runs, and that ID allocation is source-qualified — `v1-history`'s recurring `E01-example/T001-shared` across two releases is recorded as two separate destinations, one per release, that never collapse onto a single record. On the frozen fixture the `v1` recurrence is `done` under a `done` epic, so it archives rather than converting, per this epic's "completed work is archived" rule; the two-distinct-global-IDs property is proven on a temporary copy with both recurrences active.
```

### Target File
.savepoint/releases/v2/epics/E45-safe-migration/tasks/T012-migrate-a-whole-project-end-to-end.md

### Replace
```md
- [x] `v1-basic`, which has no `audit/`, no `Health-Check.md`, and no `AGENTS.md`, migrates with no finding for any of those absences. *(Covered for both fixtures' declared absences; see Deviation 2 — `v1-basic` does carry an `AGENTS.md`, `v1-history` is the fixture without one.)*
```

### With
```md
- [x] Each fixture migrates with no finding for the optional artifacts its own `manifest.yml` declares `absent_by_design`: `v1-basic` has no `audit/` and no `Health-Check.md`, and `v1-history` has no `AGENTS.md`.
```
