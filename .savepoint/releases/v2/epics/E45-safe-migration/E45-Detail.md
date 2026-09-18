---
type: epic-design
status: audited
---

# E45: Convert active work without losing history

## Purpose

Turn a V1 project into a V2 project once, predictably, with a write-free preview, exact-byte archives of everything it replaces, an accountable destination for every source file, and a recorded operation that can be resumed or explained after an interruption — never a silent partial state and never a fabricated verification result.

This is the V1 delivery epic for V2 product Objective O005. E42, E43, and E44 built the V2 record families in `internal/data`; nothing has ever written one from real project history. E45 is the first epic that produces V2 records from V1 input, and the last one that can still be wrong without a user noticing, because its output becomes the project.

## What this epic adds

- A new `internal/migrate` package owning one-time conversion: raw source inventory, classification, deterministic target planning, reference mapping, backups, an operation journal, recovery, and schema activation.
- A write-free preview that enumerates every planned target record, archive, ID mapping, conflict, and required owner decision without creating staging directories, probes, backups, or manifests.
- A raw source inventory keyed by project-relative path with exact-byte SHA-256, size, and mode, taken from bytes rather than from parsed models, which is also the freshness basis within one `Plan`-then-`Apply` call and across a resume: apply re-reads and refuses when a source no longer matches what that same call's plan (or, on resume, the operation's recorded hashes) recorded. A separately invoked `--apply`, run after an earlier separate `migrate` preview with no state carried between the two commands, always plans fresh from the project's current bytes rather than comparing against that earlier preview — see the note under "The inventory is raw, and it is the freshness basis" below.
- Deterministic conversion of active V1 records into V2 Objectives, Tasks, and Issues, preserving authored bodies verbatim and allocating fresh global `O###`/`T###`/`I###` identities that are never reused.
- A source-qualified reference map at `.savepoint/migrations/v1-to-v2.yml` recording every legacy identity, its destination or archive reference, and the typed legacy prerequisites that cannot become V2 dependencies.
- A blocking ambiguity contract: ambiguous active mappings are named with stable IDs and refuse application until the owner supplies concrete decisions through an explicit input file; they are never healed, guessed, or discarded.
- A byte-preserving archive at `.savepoint/archive/v1/`, verified by hash before any original is removed.
- A recorded, resumable operation under `.savepoint/.migration/<op-id>/` — backup, staging, and journal — where schema activation is the last write and therefore the commit point.
- Refusal of ordinary project writes by `upgrade-assets`, the board's write actions, and doctor's repair guidance while an operation is incomplete, with recovery guidance naming the operation.
- Rerun semantics: resume identical validated work, report a conflict when the user edited sources or installed files, and change nothing on a project already at schema version 2.
- `savepoint migrate [dir]` preview with an explicit `--apply`, dispatched through the existing thin command pattern.

## Components and files

| Module | Purpose |
|--------|---------|
| `internal/migrate/inventory.go` (new) | Walk project-owned `.savepoint/` content plus the managed guide and skills, recording project-relative normalized paths, exact-byte SHA-256, size, and mode; reject traversal, symlink escape, and case-collision before anything else runs. |
| `internal/migrate/classify.go` (new) | Assign each inventoried file a source role using the same role vocabulary the E41 fixture manifests record, so classification is checked against frozen evidence rather than against itself. |
| `internal/migrate/plan.go` (new) | Build the immutable conversion plan: target records, archive entries, ID allocations, reference map, conflicts, and ambiguities. Pure function of inventory plus decisions plus an injected clock and operation-ID source. Writes nothing. |
| `internal/migrate/convert.go` (new) | Render planned V2 Objective, Task, and Issue file content from V1 sources, relocating authored text under V2 headings verbatim and never paraphrasing it. |
| `internal/migrate/decisions.go` (new) | Stable ambiguity IDs, the owner decision input file schema, and the rule that an unresolved blocking ambiguity refuses application. |
| `internal/migrate/manifest.go` (new) | The `.savepoint/migrations/v1-to-v2.yml` model: source hashes, legacy-to-global ID map, archive references, typed legacy prerequisites, recorded owner decisions. Create-only. |
| `internal/migrate/operation.go` (new) | Operation journal under `.savepoint/.migration/<op-id>/`: backup, staging, per-path planned/installed hashes and step state, incomplete-operation detection, and the recovery report other packages read. |
| `internal/migrate/apply.go` (new) | Ordered publish: verify backup, create additive records, replace in-place files, archive-then-remove, activate schema last. Each step journalled before and after. |
| `internal/migrate/replace_windows.go`, `replace_unix.go` (new) | The platform replacement primitive behind the operation, selected by the bounded experiment named under Open decisions. No truncating-copy fallback for protected replacements. |
| `cmd/migrate.go` (new) | Argument parsing and dispatch only: optional directory, `--apply`, `--dry-run` (explicit synonym of the default), `--decisions FILE`, `--recover`, help, named errors, exit status. |
| `main.go` | One dispatch case for `migrate`, matching the existing command wiring. |
| `internal/data/config.go` | Unchanged decoding. `ReadSchemaVersion` is the version gate migration reads; activation writes the same field. |
| `internal/data/discover.go`, `project.go` | Unchanged. V2 discovery is already confined to `objectives/`, `checks/`, and `issues/`, so `archive/`, `migrations/`, and `.migration/` are invisible to `LoadV2Index` — a property E45 tests rather than adds. |
| `internal/init/write.go`, `upgrade.go` | Guard the asset-refresh write path against an incomplete operation. Migration reuses low-level ownership helpers only where their contracts already fit; `AtomicWrite`'s truncating-copy fallback is not one of them. |
| `internal/doctor/checks.go` | Report an incomplete migration operation as a named diagnostic with recovery guidance, and refuse repair suggestions that would write while it is pending. |
| `internal/board/io.go` | Refuse the existing status and router write commands while an operation is incomplete, at the write boundary only. |
| `internal/migrate/*_test.go`, `internal/data/migration_*_test.go` | Golden conversion of both E41 fixtures, write-free preview, hash revalidation, backup verification, injected failure at every journal step, resume convergence, second-run no-op, conflict reporting, ambiguity blocking, and post-migration `LoadV2Index` cleanliness. |
| `AGENTS.md` | One Codebase Map row for `internal/migrate`, as ARCH-04 requires. |

## Architectural delta

Migration is the first Savepoint operation that changes many files at once, and the first that can destroy user history. Everything below exists to make that operation honest rather than to make it clever.

**Migration is not asset upgrade.** `internal/init` owns provenance-aware refresh of shipped assets; it must never change `schema_version`, and `migrate` must never silently refresh assets. They meet at exactly two points: migration inventories and backs up the managed guide and skills because it can disturb them, and `upgrade-assets` refuses to run while a migration operation is incomplete. E45 archives the pre-migration copies of `AGENTS.md` and `agent-skills/` and otherwise leaves them alone — the V2 routing block and the four public skills are E46's and E47's work.

**Preview is a pure function; apply is a journalled machine.** `Plan` takes the inventory, the decision file, an injected clock, and an injected operation-ID source, and returns an immutable plan. It touches no path other than reads. This is what makes FS-03 provable by snapshotting bytes and mtimes around a preview, and what makes golden conversion tests of the two frozen E41 fixtures possible at all: with the clock and op-ID injected, the same fixture must render byte-identical target content every run. Apply consumes a plan and owns every write.

**The inventory is raw, and it is the freshness basis.** Hashes come from file bytes, never from a re-marshalled parse, so a source the loader would have healed still hashes as what the user actually wrote. Apply re-reads and re-hashes every inventoried source before its first write; any difference from the hashes it was handed aborts with a named conflict naming the paths. This is a real, byte-level check, not a timestamp comparison — but it is a check against *the plan `Apply` was actually called with*, not against some earlier command invocation the process has no memory of. Within a single `savepoint migrate --apply` call, `Plan` and `Apply` run back to back, so this closes the gap between "what the plan says" and "what the project now contains" for that one call. On a resume, the check is against the operation's own recorded source hashes from when it was created, which is the same guarantee. What it deliberately does **not** cover: a user who runs `savepoint migrate` (preview only) in one command, walks away, and later runs `savepoint migrate --apply` in a separate command — that second invocation re-plans from the project's current bytes with no reference to the earlier preview at all, so a source edited in between is migrated silently, with no conflict raised. Closing that gap would mean persisting the previewed plan's hashes somewhere between two separate command invocations, which this epic does not build; the guarantee this epic makes is scoped to a single invocation's `Plan`-then-`Apply` call and to resume, not to two independently invoked commands.

**Classification reuses the E41 role vocabulary.** The frozen fixture manifests already name each source file's role — `config`, `router`, `product-prd`, `architecture`, `health-check`, `release-prd`, `epic-detail`, `epic-audit`, `task`, `defect`, `audit-prompt`, `audit-register`, `finding`, `audit-run` — together with raw status, raw phase, scoped ID, and expected classification. E45 classifies against that vocabulary and asserts agreement with the manifests, which turns E41's fixtures from passive samples into the compatibility contract they were frozen to be. A file matching no role is inventoried, archived, and reported as unclassified; it is never dropped and never guessed at.

**Identity is allocated, not translated.** A V1 `T001` is not a V2 `T001`. The legacy key is source-qualified — release, epic, path, and original ID — exactly as `v1-history`'s recurring `E01-example/T001-shared` across `v1` and `v1.1` proves it must be. Allocation walks sources in a fixed order and assigns the next unused global ID, so conversion is deterministic and reruns produce identical IDs. Identities are never reused: archived-only records reserve nothing, because they receive no V2 ID at all.

**Completed work is archived, not converted.** Done Tasks and audited epics become archive content plus a mapping entry. This creates the one place where a V2 schema cannot express a V1 fact: an active Task may depend on a completed Task that has no V2 identity, and `TaskV2.DependsOn` requires a `T###` that exists in the index, so writing that dependency would make the migrated project fail to load. **Resolution: the converted Task's `depends_on` contains only converted active Tasks; a dependency on an archived completed Task becomes a typed `legacy_prerequisite` entry in `v1-to-v2.yml`, keyed by the new `T###`, naming the archive path and the original recorded completion or waiver evidence, plus an authored line in the Task body.** No V2 record family changes, no loader changes, no fabricated `last_check` and no fabricated CLEAR. E48 and E49 may display these as legacy prerequisites when they adopt V2; E45's obligation is only that the fact is recorded, resolvable, and visibly not a Check result.

**Nothing invents evidence, configuration, or prose.** Converted records carry no `last_check`, no `freshness`, and no `owner_validation.accepted_check` unless the V1 source recorded the corresponding fact; absent freshness is `unknown` by absence, which E43's rules already treat as blocking. `quality_gates` are carried across from V1 `config.yml` verbatim; `Health-Check.md` is archived and preserved as an optional procedure file with its candidate commands *listed in the preview for the owner*, because parsing authored prose into executable configuration is exactly the guess this epic exists to refuse. Authored bodies are relocated under V2 headings with their source section named, never rewritten.

**Lifecycle values are mapped or refused, never healed.** V1's loader heals unknown status to `planned` and aliases `stage: implementation` to `build`; migration may apply the recorded alias but must refuse an unrecognized raw status on an *active* record as a blocking ambiguity. Healing is how V1 keeps working; it is not how a one-time conversion decides what a user meant.

**Issue conversion follows the finding disposition, and duplicates need a live target.** Unresolved defects and `open`/`triaged`/`mapped` findings become open Issues; `deferred` and `owner_decision` become open Issues carrying a dated history entry of that kind rather than a resolution, because neither closes anything; `in_progress` and `fixed` become in-progress Issues with a `repair_attempted` entry, since fixed still awaits independent verification; `verified` and resolved defects are archived with their original proof and imply no V2 Check. A `duplicate` finding resolves as `duplicate` only when its canonical finding also converts, because E44 requires `duplicate_of` to name an existing Issue; when the canonical is archive-only, the duplicate is archive-only too. Every converted Issue's `source` records migration with its legacy key, and its seeded history entry uses the date the V1 record recorded, falling back to migration time with an explicit note rather than inventing a date.

**Ambiguity blocks, with a name and a place to answer.** Preview assigns each ambiguity a stable ID derived from its source path and kind, so the same project previews the same IDs every time. Blocking ambiguities — unrecognized active lifecycle values, a missing dependency target, duplicate source identities, and unresolved narrative findings that name no stable identity — refuse apply until the owner supplies an explicit decisions file (`--decisions FILE`) mapping each ID to a concrete outcome. Decisions are recorded into the manifest with their provenance. Advisory ambiguities are reported and do not block. Savepoint does not decide whether two narrative findings describe the same problem.

**Recovery is recorded, not claimed atomic.** The operation lives at `.savepoint/.migration/<op-id>/` holding `backup/`, `staging/`, and `operation.yml`. That location is deliberate: it is inside the project so it survives with it and stays writable, and it is outside every path the operation writes to, so a partial apply can never damage its own recovery data. It is dot-prefixed like the existing `.upgrade-manifest.yml`, invisible to V2 discovery, and create-only — an existing `<op-id>` directory is refused rather than overwritten, which is how "never overwrite a prior backup or user sidecar" is enforced mechanically.

Apply proceeds in a fixed, journalled order: revalidate source hashes; write backups of every file the operation will replace or remove and verify each backup copy against the recorded source hash; stage complete target files; create additive records (`objectives/`, `issues/`, `.savepoint/archive/`, `migrations/v1-to-v2.yml`); replace modified in-place files; remove originals only after their archive copies verify; and activate `schema_version: 2` last. Until that final write, the project is still a valid V1 project; after it, a valid V2 one. That single-file activation is the commit point and the closest honest approximation of atomicity — the epic claims recoverable multi-file change and truthful failure reporting, never atomic multi-file replacement. FS-06's no-partial-write expectation applies to preflight; an interruption after a valid apply has begun is a recorded recoverable state, and the distinction is stated in the diagnostics rather than blurred.

**`.savepoint/migrations/` is already occupied.** `internal/init/migrate_audit_skill.go` writes archived legacy skill copies and a `README.md` into that exact directory today. The schema manifest coexists as `v1-to-v2.yml` alongside them: existing contents are inventoried, archived, and preserved untouched, the README is not rewritten, and an existing `v1-to-v2.yml` refuses application as an already-migrated or conflicting project rather than being overwritten.

**Other write paths stand down while an operation is pending.** `migrate.PendingOperation(root)` is a read-only detector returning the incomplete operation and its recovery guidance. `upgrade-assets`, the board's status and router writes in `internal/board/io.go`, and doctor's repair suggestions consult it at their write boundary and refuse with that guidance. This is a guard at three existing write sites, not new presentation and not new policy; the board's V2 presentation belongs to E49.

**Rerun is a first-class path.** A second invocation discovers the single incomplete operation, revalidates sources and installed files, and either resumes from the journalled step or reports a conflict when the user changed either — recovery never overwrites a user's edits automatically. More than one incomplete operation is a named diagnostic, not a heuristic choice. A run against a project already at `schema_version: 2` reports that and changes no file, mtime, or ID; a run against an unknown explicit version fails through the existing named diagnostic.

Reference: `.savepoint/releases/v2/v2-Design.md`, sections 2, 3, 9, and 10.

## Boundaries

**In scope:**

- The `internal/migrate` package, its preview/apply split, and the `migrate` command with thin dispatch.
- Raw inventory, role classification against E41 fixture manifests, deterministic planning, and verbatim body relocation.
- Global ID allocation, the source-qualified reference map, typed legacy prerequisites, and the `v1-to-v2.yml` manifest.
- Ambiguity identification, the owner decision input, and refusal to apply with unresolved blocking ambiguities.
- Byte-verified archival, backup, operation journal, ordered publish, schema activation, recovery, and rerun semantics.
- The platform replacement primitive and its Windows and Unix recovery evidence.
- Incomplete-operation refusal at the three existing write boundaries and the matching doctor diagnostic.

**Out of scope:**

- Changing any V2 record schema, loader, gate, or write path from E42, E43, or E44. If a conversion cannot be expressed, it becomes a manifest entry and a named limitation, not a schema change.
- V2 scaffolding, the V2 default for `init`, the V2 `AGENTS.md` routing block, and the four public skills — E46 and E47 own them. E45 archives the pre-migration copies and rewrites neither.
- Board or resume presentation of migrated records, legacy prerequisites, or archive content. E48 and E49 own it; E45 guarantees only that the data resolves.
- Removing V1 live readers or transitional dispatch. E50 owns cutover; migration input readers and the frozen fixtures survive it.
- Migrating any live maintainer project, including this one. Tests use temporary projects and copies of frozen fixtures only.
- Automatic deduplication of narrative findings, inferred quality gates, reconstructed Idea or Design content, and any rewriting of authored prose.
- A general-purpose migration framework, a rollback command, or version-to-version upgrade beyond V1 to V2.

## Dependencies

- E44-issues-objective-checks supplies the Issue record family, resolution dispositions, append-only history, and Objective gating that migrated follow-up and Objectives must satisfy.
- E41-migration-source-fixtures supplies the two frozen V1 projects, their hand-authored role manifests, and the characterization tests this epic converts against.

## Quality gates

- Preview tests prove a full preview over both frozen fixtures creates no file, directory, probe, backup, or manifest, and leaves every source byte and mtime unchanged, on a success path and on every named failure path.
- Determinism tests prove that with an injected clock and operation ID, each fixture converts to byte-identical target content and identical allocated IDs across repeated runs, and that ID allocation is source-qualified — `v1-history`'s recurring `E01-example/T001-shared` across two releases is recorded as two separate destinations, one per release, that never collapse onto a single record. On the frozen fixture the `v1` recurrence is `done` under a `done` epic, so it archives rather than converting, per this epic's "completed work is archived" rule; the two-distinct-global-IDs property is proven on a temporary copy with both recurrences active.
- Classification tests prove every file in both fixture manifests is assigned the role its manifest records, and that an unrecognized file is inventoried, archived, and reported as unclassified rather than dropped.
- Freshness tests prove apply revalidates raw source hashes and aborts with a named conflict, without writing, when any inventoried source changed within a single `Plan`-then-`Apply` call or since an operation was created (the resume case); this scope is deliberate, not a gap — see "The inventory is raw, and it is the freshness basis" above for what a separately invoked preview and apply does not cover.
- Ambiguity tests prove an unrecognized active lifecycle value, a missing dependency target, a duplicate source identity, and an unresolved narrative finding each produce a stable ambiguity ID, refuse apply, and are recorded with provenance once a decisions file resolves them; and that a decisions file naming an unknown ID is itself an error.
- Backup tests prove every file the operation replaces or removes is backed up and hash-verified before the first replacement, that an existing operation directory refuses rather than overwrites, and that an original is removed only after its archive copy verifies.
- Interruption tests inject a failure at every journalled step, proving each leaves a truthful report, recoverable bytes for every affected path, and no case in which both the original and its backup are lost.
- Resume tests prove a rerun after each injected failure converges to the same final state, that a user edit to a source or an installed file between runs reports a conflict instead of overwriting, and that two incomplete operations produce a named diagnostic.
- Idempotence tests prove a successful second full run changes no file content, mtime, or ID, and that a project already at `schema_version: 2` reports no work and writes nothing.
- Result tests prove the migrated fixture projects load through `LoadV2Index` with no diagnostic, that every `objective`, `depends_on`, `issues`, `checks`, and `duplicate_of` reference resolves, that no converted record carries clearance the source did not record, and that `archive/`, `migrations/`, and `.migration/` remain invisible to discovery.
- Archive tests prove every archived file hashes equal to its source bytes, including the CRLF and unknown-field content the fixtures freeze, and that pre-existing `.savepoint/migrations/` content and its README survive untouched.
- Issue conversion tests cover each V1 finding disposition and defect state mapping named above, including a `duplicate` whose canonical converts and one whose canonical is archive-only.
- Guard tests prove `upgrade-assets`, the board's status and router writes, and doctor's repair suggestions refuse while an operation is incomplete, name the operation, and write nothing.
- Platform tests prove the replacement primitive and recovery behave identically on Windows and Unix, or that any difference is explicit and tested, satisfying CFG-02 and REL-01.
- Command tests prove `migrate` defaults to preview, requires `--apply` to write, accepts `--dry-run` as an explicit synonym, reports named errors with nonzero status for a missing, unwritable, or non-Savepoint target, and keeps `cmd/` free of domain parsing.
- Implementation handoff requires focused `internal/migrate`, `internal/data`, `internal/init`, `internal/doctor`, and `cmd` tests plus `make build && make test`, with named cases and outcomes recorded in the Tasks.
- Epic closeout requires a fresh independent V1 epic audit; this planning session is not that audit. Current Guardrails apply, with STYLE advisory. FS-01, FS-03, FS-04, FS-05, FS-06, DATA-01, DATA-03, ARCH-01, ARCH-04, CFG-02, TEST-03, and TEST-04 are the load-bearing ones here.

## Open decisions

One, deliberately bounded: **the Windows file replacement and interruption-recovery primitive.**

Everything above is fixed and does not depend on the answer — the journal, the ordered publish, the verified backup, and the last-write commit point are the same whichever primitive wins. What is unresolved is which call the operation makes when it replaces a protected file on Windows, and this planning session cannot settle it, because this environment is Linux and the question is one of observed platform behavior, not of reasoning.

Task breakdown must therefore create a bounded research Task as E45's first Task, with a decision deliverable and recorded evidence, before any apply-path Task starts. Its candidates and the evidence each must produce:

1. `os.Rename` — Go maps this to `MoveFileEx` with `MOVEFILE_REPLACE_EXISTING`. Same-volume only, which the same-directory temp file in `internal/data`'s existing `replaceV2File` already guarantees. Evidence needed: behavior when the destination is open by another process or marked read-only, and what the destination's attributes and ACLs look like afterwards.
2. `ReplaceFileW` via `golang.org/x/sys/windows` — preserves destination attributes and ACLs and takes a backup path. `golang.org/x/sys` is already an indirect dependency, so adopting it is a DEP-01 justification rather than a new dependency, and that justification must be recorded with the decision.
3. Retry-with-backoff over either, for transient sharing violations only, with a bounded attempt count and a truthful failure report when it exhausts.

The decision Task must record, for the chosen primitive: observed behavior on an open destination, on a read-only destination, and across an interruption at each journalled step; whether recoverable bytes survive in every case; and the reason the alternatives were rejected. The required outcome is fixed and not open: recoverable bytes and truthful failure reports, never silent loss. `AtomicWrite`'s truncating-copy fallback in `internal/init/write.go` is already excluded for protected replacements and is not a candidate.

No other decision is open. The package boundary, preview/apply split, inventory and classification contracts, ID allocation, the legacy-prerequisite resolution, the operation directory and journal, the publish order and commit point, the ambiguity gate, and the consumer guard boundary above are fixed. Exact Go identifiers may be refined during task breakdown.
