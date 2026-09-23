---
id: O018
title: Hyphenate every numbered V2 record identity
status: planned
release: R006
---

# O018: Hyphenate every numbered V2 record identity

## Outcome

Every V2 record identity — Release, Objective, Task, Check, and Issue — uses
the hyphenated form `R-###`/`O-###`/`T-###`/`C-###`/`I-###`. This
repository's existing records are renamed, and every identity minted from now
on carries the hyphen.

## Why

I022: bare numbered identities such as `O015` and `T012` are harder to scan
than `O-015` and `T-012`. The owner asked for the hyphenated form on every
numbered record type, applied to this repository's own records.

## Success Conditions

- One shared identity rule in `internal/data` (`^[ROTCI]-[0-9]{3,}$` per
  kind) replaces the five per-kind regexes; unhyphenated IDs are refused.
- Check and Issue generators (`write.go`) and the V1→V2 allocator
  (`migrate/plan.go`) emit hyphenated IDs.
- Skill, template, and fixture examples use the hyphenated form;
  canonical/scaffold skill copies stay byte-identical.
- This repository's active records, directories, filenames, cross-references,
  `router.md`, and active prose (`AGENTS.md`, `Design.md`, `Guardrails.md`)
  are renamed by a one-off script kept in the repo for other projects to run.
- Strict load, board, doctor, and resume succeed on this repository; a
  targeted search finds no unhyphenated active IDs; `make build` and
  `make test-full` pass, followed by the mandatory Full Objective Check.

## Boundaries

**In scope:**

- The shared identity rule, generators, and tests.
- Skill/template/fixture examples.
- A simple rename script and its one run against this repository.

**Out of scope:**

- A product migration command, preview/recovery machinery, progress markers,
  or archived copies of original Checks.
- Accepting both identity forms.
- Any change under `.savepoint/archive/`.
- Any UI or behavior change beyond identity formatting.

## Originating Issue

I022 — Hyphenate every numbered V2 record identity.

## Confirmed Design Decisions

On 2026-09-23, the owner cut the earlier recoverable-migration design back to
a simple rename. Other projects migrate by running the in-repo rename script
rather than a `savepoint` command. Existing Checks are renamed and have their
identity references rewritten in place; no other Check content changes. T019
and T020 land in one commit because the repository does not load between them.
