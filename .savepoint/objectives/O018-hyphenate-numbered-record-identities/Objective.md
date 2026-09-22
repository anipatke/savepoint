---
id: O018
title: Hyphenate every numbered V2 record identity
status: planned
release: R006
---

# O018: Hyphenate every numbered V2 record identity

## Outcome

Every active V2 record identity — Release, Objective, Task, Check, and Issue
— uses the hyphenated form `R-###`/`O-###`/`T-###`/`C-###`/`I-###` everywhere
it appears: validation, generation, filenames, directories, cross-references,
templates, skills, and this repository's own live records. No mixed or
unhyphenated active identity is left behind.

## Why

I022: bare numbered identities such as `O015` and `T012` are harder to scan
than `O-015` and `T-012`. The owner asked for the hyphenated form on every
active numbered record type and explicitly required this repository's own
active graph to be migrated as part of the work — not a display-only alias.

## Success Conditions

- One canonical identity grammar (kind + hyphen + digits) shared by all five
  record kinds, replacing the five divergent regexes in `internal/data`
  (`release_v2.go`, `objective_v2.go`, `task_v2.go`, `check_v2.go` ×2).
- Every ID-bearing surface uses it: Check/Issue generation in `write.go`, the
  V1→V2 `idAllocator` in `migrate/plan.go`, filenames and directory names for
  Objectives/Tasks/Checks/Issues/Releases, every cross-reference field
  (`depends_on`, `release`, `objective`, `tasks`, `checks`, `duplicate_of`,
  `last_check`, router.md selections, evidence/acceptance/exception blocks),
  CLI filters, and board/doctor/resume rendering.
- A previewable, interruption-safe migration tool: strict-loads the index
  first (fails closed on anything already broken), computes the full rename
  map, refuses up front on any collision or partial-migration state, rewrites
  references and renames files/directories atomically, and is safely
  re-runnable if interrupted.
- Fresh-init and upgrade-assets templates emit only hyphenated identities;
  canonical/scaffold skill copies stay byte-identical per existing
  convention.
- This repository's complete active V2 graph is migrated: every Release,
  Objective, Task, Check, and Issue (including this Objective's own
  originating Issue, `I022` → `I-022`, renamed last so nothing mid-migration
  references its own stale filename), `router.md`, and active prose in
  `AGENTS.md`/skills. Archived V1 material and immutable Check contents are
  left untouched.
- Old unhyphenated identity strings are not accepted anywhere once migration
  lands — this is a one-shot atomic cutover, not a gradual rollout, so no
  dual-format compatibility path is built. (Flag to the owner if a compat
  window is actually wanted; default assumption is a hard cutover.)
- A repository-wide targeted search shows no unintended active unhyphenated
  identity remains; `git diff --check`, `make build`, `make test`, and the
  mandatory Full Objective Check (plus R006's Release Check) pass.
- Negative-path tests cover malformed, mixed, duplicate, and
  partially-migrated identities.

## Architectural Considerations

- `internal/data` stays the sole owner of the identity grammar and
  validation; board, doctor, resume, and the CLI consume the resolved type
  rather than parsing identities locally.
- Build the migration as a narrow, purpose-built V2→V2 tool rather than
  routing it through `internal/migrate`'s V1→V2 schema converter, but reuse
  that package's proven shape (preview → guarded apply → interruption
  recovery → collision rejection) and its atomic platform file-replace
  primitive (`replace.go`, `replace_unix.go`, `replace_windows.go`) instead
  of duplicating that machinery.
- Objective and Task identities currently have no code-level generator — they
  are minted by hand during `savepoint-design`/`savepoint-task` — so this
  Objective also touches skill prose/examples (`O###`/`T###` in `SKILL.md`
  templates), not only Go code.
- O013 (active, in progress) keeps `R###` as the persisted Release identity
  and only relabels it "Goal" in the UI, which is compatible with
  hyphenating the persisted form underneath it. The two Objectives touch
  overlapping board/resume rendering code, so sequence this Objective's Tasks
  after O013 lands rather than interleaving them.

## Boundaries

**In scope:**

- Shared identity grammar/validation and every internal generator/consumer
  of it.
- The V2→V2 migration tool (preview, apply, interruption recovery, collision
  rejection).
- Fresh-init/upgrade-assets template updates.
- Migrating this repository's active V2 graph and reconciling active prose
  (`AGENTS.md`, skills, `Design.md`/`Guardrails.md` where they reference
  identities).

**Out of scope:**

- Any change to archived V1 material, `.savepoint/archive/`, or immutable
  Check contents.
- Any product/UI behavior change beyond identity formatting — no relation to
  O013's Release→Goals wording work beyond not colliding with it.
- Building general dual-format (hyphen/non-hyphen) compatibility support.

## Originating Issue

I022 — Hyphenate every numbered V2 record identity.
