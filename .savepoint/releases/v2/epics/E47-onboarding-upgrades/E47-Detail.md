---
type: epic-design
status: planned
---

# E47: Start and refresh coherent V2 projects

## Purpose

Make `savepoint init` produce a V2 project, and make `savepoint upgrade-assets` refresh the assets that match the project it is pointed at — so a new project starts on the four-role workflow, an existing codebase can adopt it without its files being rewritten, and a V1 project that has not migrated yet is neither broken nor silently converted.

This is the V1 delivery epic for V2 product Objective O007. E45 made conversion possible and E46 wrote the four skills, but nothing yet hands either to a user: `init` still scaffolds a V1 project, and `upgrade-assets` still installs V1 instructions over every project it finds. E47 is where the V2 workflow stops being documentation and becomes what a project actually receives.

## What this epic adds

- A V2 project scaffold — Idea, Design, Guardrails, config at `schema_version: 2`, a router that starts at `state: idea`, an empty `objectives/`, and a V2 agent guide whose routing table is live rather than marked inactive.
- `savepoint init` producing that scaffold by default, with V1 scaffolding retired in the same change. There is no `--v1` escape hatch.
- A bootstrap prompt that accepts one rough sentence and routes it to `savepoint-idea`, with no prepared requirements document demanded first.
- Guidance for adopting Savepoint into an existing codebase: reconstruct Design from targeted reads of the code, ask the owner for intent, and change no file the user authored.
- `upgrade-assets` selecting its asset stream from the project's own `schema_version`: V2 assets for a V2 project, today's V1 assets for a project that has not migrated, and a named diagnostic for a version it cannot interpret.
- Retirement of the nine V1 skills from a project that has become V2, archived before deletion so local edits survive.

## Components and files

| Module | Purpose |
|--------|---------|
| `templates/project-v2/` (new) | The V2 scaffold tree: `.savepoint/Idea.md`, `Design.md`, `Guardrails.md`, `config.yml`, `router.md`, `objectives/.gitkeep`, `AGENTS.md`, and `agent-skills/` holding the four V2 skills plus the three shared references. |
| `templates/project/` | Unchanged as the legacy asset stream for projects still at schema version 1, minus the four V2 skills and three V2 references that move to the new tree. Deleted by E50. |
| `templates/prompts/magic-prompt.prompt.md` | Rewritten for V2 bootstrap: one rough sentence in, `savepoint-idea` out. It stays the only prompt template; `init` no longer creates V1 projects, so no second prompt is needed. |
| `main.go` | Embed directives for the new tree, and the tree selection passed into `Scaffold` and `UpgradeProjectAssets`. |
| `internal/init/scaffold.go` | Unchanged walk semantics over whichever tree it is given. The caller chooses the tree; scaffolding does not detect versions. |
| `internal/init/upgrade.go` | Schema-version dispatch at the front of the upgrade, tree-specific asset policy behind it, and the V1 retirement step for a V2 project. |
| `internal/init/retire_v1_skills.go` (new) | Archive-then-delete of the nine retired V1 skills, generalizing the pattern `migrate_audit_skill.go` already established for the split audit skill. |
| `internal/init/manifest.go` | Provenance entries for retired skills are dropped when their files are, so the manifest never records a file the project no longer has. |
| `internal/data/config.go` | `ReadSchemaVersion` is the version gate upgrade reads. Unchanged; upgrade never writes the field. |
| `internal/init/*_test.go` | Per-tree parity contracts, the fresh-V2-init lifecycle, schema dispatch, retirement, and the existing ownership matrix re-proven against the V2 tree. |
| `AGENTS.md`, `.savepoint/Guardrails.md` | The live guide's V2 routing section becomes "active for V2 projects, not for this repository until E50"; TPL-01 is amended to name both shipped trees. |

## Architectural delta

**Two shipped trees, chosen by the project's own declared version.** Today one template tree serves every project. V2 cannot share it: `config.yml`, `router.md`, and `AGENTS.md` exist at the same paths in both lifecycles with irreconcilable content. So `templates/project-v2/` is added beside `templates/project/`, and the selection between them is made once, at the top of each command, from `data.ReadSchemaVersion`. Nothing below that point knows which lifecycle it is serving — `Scaffold` and the upgrade walk take a tree and apply the same ownership rules to it.

This is why the version gate is the project's `schema_version` and nothing else. Not the package version, not the upgrade manifest's version, not whether V2 files happen to be present. A project declares what it is; asset installation reads that declaration and never changes it. `migrate` remains the only operation that writes `schema_version`.

**`init` flips; `upgrade-assets` does not.** New projects are V2 from this epic on, and V1 scaffolding is gone in the same change — a half-flipped default, where `init` can still produce either, is a second lifecycle to maintain for no user. But an existing V1 project that runs `upgrade-assets` keeps receiving exactly the V1 assets it receives today, byte for byte, with one added informational line naming `savepoint migrate` as the route to V2. Upgrading assets has never been allowed to change what a project is, and giving a V1 project V2 instructions would do precisely that: its router has no `idea` state and its records are epics and tasks, so the four skills would route it nowhere.

**Retirement follows migration, not installation.** The nine V1 skills are retired from a project when that project is V2 and asks for an asset refresh — never on the strength of the package version alone. Retirement reuses the preserve-then-delete order `migrate_audit_skill.go` already proves: archive the content under `.savepoint/migrations/`, verify it is on disk, then remove the triggerable copy, and only then drop its manifest entry. A user-edited skill is archived with its edits intact; an unmodified one is archived too, because deciding which local changes were worth keeping is not a judgment this command gets to make. The live repository keeps all nine skills and its V1 routing table: Savepoint is still built under the V1 lifecycle until E50, and shipping a default it does not itself use is the honest state of a mid-release cutover, not a contradiction.

**The canonical-to-shipped parity contract becomes per-tree.** TPL-01 pairs each `agent-skills/{skill}/SKILL.md` with one shipped copy, and today's test discovers that pairing by listing the live skill directory and requiring every entry in `templates/project/`. With two trees the pairing is no longer derivable from a directory listing, so it becomes an explicit table: nine V1 skills and `references/audit-method.md` to `templates/project/`; `savepoint-idea`, `savepoint-design`, `savepoint-task`, `savepoint-check` and the three V2 references to `templates/project-v2/`. Byte parity and set-completeness are asserted in both directions, so a skill added to the live tree and forgotten in its shipped tree still fails. TPL-01's wording is amended to match; the rule does not change, only the path it names.

**Onboarding guidance is scaffold content, not skill content.** The four skills are byte-parity-locked contracts from E46 and are not reopened here. What a fresh agent needs — that a rough sentence is a valid start, and that an existing codebase means reconstructing Design from the code rather than inventing it — lives in the V2 `AGENTS.md`, the V2 `router.md`'s opening `next_action`, the V2 `Design.md` template's Current Technical State section, and the bootstrap prompt. Reconstruction stays an agent reading files it is told to read: no model service, no automatic scan, and no write to any file the user authored.

**Optional things stay optional.** The V2 scaffold ships no Concept, no Health-Check, no release PRD, no audit register, and no procedures file. Guidance that mentions them must degrade when they are absent, because absent is the normal case — TPL-03, and the rule `commands-and-procedures.md` already states for a project with no `Health-Check.md`, which is this one.

**Board, doctor, and resume are not part of this.** A fresh V2 project is proven to load through `data.LoadProject` with no diagnostic; what the board draws for it is E49's, and what doctor and resume report is E48's. The scaffold default flips inside an unreleased line, so no user meets a V2 project before those land.

Reference: `.savepoint/releases/v2/v2-Design.md`, sections 2, 3, 7, 9, and 10.

## Boundaries

**In scope:**

- The `templates/project-v2/` tree and the removal of V2 assets from `templates/project/`.
- `init` defaulting to V2 and V1 scaffolding being retired.
- The rewritten bootstrap prompt and the existing-codebase adoption guidance.
- Schema-version dispatch in `upgrade-assets`, including the unreadable-version diagnostic.
- Retirement of the nine V1 skills from a V2 project, with archival and manifest reconciliation.
- The per-tree parity contract, the TPL-01 amendment, and the live guide's V2 routing statement.

**Out of scope:**

- Any change to the four V2 skills or three shared references. They are E46 contracts and are moved, not edited.
- Removing the V1 template tree, the nine live V1 skills, or transitional V1 readers. E50 owns that.
- Board, doctor, or resume presentation of a V2 project, and any V2 change to `internal/board` or `internal/doctor`.
- Writing `schema_version`, converting records, or any migration behavior. `migrate` owns it.
- Reconstructing Idea or Design content by any automatic process, eager scan, or model service.
- Migrating this repository, and any change to its live `.savepoint/` records beyond the Guardrails amendment and the routing statement named above.
- `templates/release/`, which no command embeds and which E50 resolves with the rest of the release scaffolding.

## Quality gates

- Fresh init produces a project whose `schema_version` reads as 2, that loads through `data.LoadProject` with no diagnostic and an empty index, and that contains no epic, release, PRD, Concept, Health-Check, or audit-register file.
- Init into a directory that already holds source files leaves every one of those files byte-identical and preserves an existing agent guide's content outside the managed block.
- Both shipped trees match their canonical sources byte for byte in both directions, and neither tree contains a skill belonging to the other.
- A V1 project's upgrade produces the same actions and the same bytes as it does before this epic, with `schema_version` unwritten and the migrate route named.
- A V2 project's upgrade installs the V2 assets, retires the nine V1 skills through archive-then-delete, and leaves an archived copy of every retired skill, edited or not.
- An edited V2 skill conflicts and keeps the user's bytes; an unmarked agent guide conflicts and stays byte-identical; a marked guide refreshes only inside the markers; `--force` adopts with a recoverable backup.
- A project with no procedures file, no Health-Check, and no Concept upgrades cleanly and produces no finding about their absence.
- A second upgrade after any of the above changes no file content, no mtime, and no manifest entry; a dry run of each reports the same actions and writes nothing.
- A malformed or unsupported `schema_version` refuses with a named diagnostic before the first write.
- Implementation handoff requires focused `internal/init` and `cmd` tests plus `make build && make test`, with named cases recorded in the Tasks. Epic closeout requires a fresh independent V1 audit; this planning session is not that audit. Current Guardrails apply with STYLE advisory; FS-01, FS-02, FS-03, FS-04, FS-06, DATA-02, TPL-01, TPL-02, TPL-03, TPL-04, CFG-01, TEST-03, and TEST-08 are the load-bearing ones.

## Open decisions

None. The two-tree layout, the version gate, the one-way `init` flip, the retirement order, the per-tree parity table, and the placement of onboarding guidance in scaffold content are settled above. Exact Go identifiers may be refined during task breakdown.
