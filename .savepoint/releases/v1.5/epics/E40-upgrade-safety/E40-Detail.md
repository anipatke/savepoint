---
type: epic-design
status: planned
---

# E40: Upgrade Safety and Backward Compatibility

## Purpose

Make `savepoint upgrade-assets` safe to run on a project that has diverged from the shipped templates, and prove with tests that template changes do not break files created by earlier versions.

Three defects exist today, all silent:

1. **Skills are overwritten wholesale.** `upgrade.go:200` replaces every `agent-skills/*/SKILL.md` with the shipped copy, with no merge, no backup, and no warning. It is reported as a routine `updated`. A project that tailors a skill loses that work. This is not a hypothetical case: `AGENTS.md:24` instructs agents to read `agent-skills/{skill}/SKILL.md` directly when the skill tool cannot find a skill, so skill files are user-facing by design, and E39 has skills referencing project-owned Guardrails rules.
2. **AGENTS.md gains a duplicate block.** `replaceManagedBlock` (`agents.go:62-71`) swaps the managed block in place when `<!-- SAVEPOINT:BEGIN -->` is present, but falls through to appending when it is absent. A hand-written or pre-marker `AGENTS.md` therefore ends up holding two sets of workflow instructions, the older one unmarked and never refreshed again. For the file whose only job is instructing agents, contradictory duplicated guidance is the worst available outcome.
3. **`--force` does nothing.** `UpgradeProjectAssets` accepts `force bool` (`upgrade.go:71`) and never reads it. `main.go:106` passes `opts.Force`, so the flag is silently inert — the natural response to "my file did not update" produces no change and no error.

Underneath all three sits a missing guarantee. `upgrade.go:118-121` skips the whole `.savepoint/` subtree, so `router.md` is never migrated and old routers persist indefinitely. Compatibility therefore has to be enforced in the reader, not the template, and nothing currently tests that an old-shape file still loads.

## Governing rule

**Never silently destroy, never silently duplicate.** Replace a file only when it is provably unmodified. Otherwise keep the user's version and offer ours alongside it.

## What this epic adds

- An upgrade manifest at `.savepoint/.upgrade-manifest.yml` recording the SHA-256 of each skill file as Savepoint last wrote it, so a genuine user edit is distinguishable from a merely outdated copy.
- A `conflict` upgrade action: the user's file is kept, the incoming version is written beside it as `<name>.new`, and the report names it.
- A real meaning for `--force`: overwrite a customized file anyway, after saving the previous content as `<name>.bak`.
- A fix for the `AGENTS.md` no-marker path so it reports a conflict instead of appending a second block.
- Backward-compatibility fixtures pinning old file shapes, asserted to parse and upgrade cleanly.

## Behaviour by file

| Path | Owner | Upgrade behaviour after this epic |
|------|-------|-----------------------------------|
| `agent-skills/*/SKILL.md` | Savepoint | Unmodified (matches manifest) → replaced. Customized → kept, new version written as `SKILL.md.new`, reported `conflict`. With `--force` → previous saved as `SKILL.md.bak`, then replaced |
| `AGENTS.md` | Shared: managed block plus user prose | Markers present → block replaced in place, surrounding prose untouched (unchanged behaviour). Markers absent → file left alone, merged result written as `AGENTS.md.new`, reported `conflict`. Missing → written whole |
| `.savepoint/router.md` | User — live instance state | Still skipped, permanently. Compatibility is a reader guarantee: additive fields only, `## Current state` and its ```yaml fence frozen as the parse contract, previously valid `state` values always accepted |
| `.savepoint/Design.md`, `PRD.md`, `Concept.md`, `visual-identity.md`, `config.yml` | User | Still skipped. New `config.yml` keys must carry defaults, since existing files will never gain them |
| `.savepoint/audit/*` | User | Unchanged install-if-missing |
| `.savepoint/Guardrails.md`, `Health-Check.md` | User | Owned by E39's policy-asset allowlist; untouched here |

## Manifest scope

The manifest covers `agent-skills/*/SKILL.md` only. `AGENTS.md` already carries its own ownership signal in the marker pair and needs no hash. No other path is wholesale-owned by Savepoint, so nothing else belongs in it.

## Migration for projects with no manifest

A project upgrading from a pre-manifest version has no recorded hashes. Treating every unrecognised skill as a conflict would emit a `.new` file for each one on that first upgrade, which is noise for the common case of an untouched-but-outdated project.

Instead, when a skill has no manifest entry and its content differs from the incoming template, save the existing content as `SKILL.md.bak`, overwrite, and record the hash. This is a one-time, recoverable path, and the upgrade report must state that backups were written. Every subsequent upgrade uses exact provenance.

This is the one place the epic accepts replacing possibly-customized content. The alternative — conflict-on-unknown — is a one-line change if the noise proves preferable in practice.

## Components and files

| Module | Purpose |
|--------|---------|
| `internal/init/manifest.go` | Manifest model, load, record, and atomic save |
| `internal/init/upgrade.go` | `ActionConflict`, manifest-aware skill path, `force` wired through, `.new` and `.bak` emission |
| `internal/init/agents.go` | No-marker path returns a conflict instead of appending |
| `internal/init/scaffold.go` | Write the manifest at `init` so fresh projects start with exact provenance |
| `internal/init/testdata/legacy/` | Pinned old-shape `router.md`, marker-less `AGENTS.md`, customized `SKILL.md` |
| `internal/init/upgrade_test.go` | Conflict, force, backup, dry-run, and idempotency coverage |
| `internal/data/router_test.go` | Old-shape router parses; unknown fields ignored; absent optional fields default |

## Architectural delta

Upgrade gains a notion of provenance it does not have today. Until now every decision was made by comparing on-disk content to the incoming template, which cannot distinguish "the user changed this" from "this is last version's copy". The manifest supplies that third data point.

`UpgradeAction` gains `conflict`, the first action that reports work deliberately *not* done. `Format()` must surface it prominently, since a conflict needs the user to act.

Router compatibility is stated as an explicit reader contract rather than left implicit: `ReadState` (`internal/data/router.go:28`) already tolerates unknown and missing YAML keys, but it hard-fails without the `## Current state` heading and its ```yaml fence. That tolerance is currently accidental and untested; this epic pins it.

## Boundaries

**In scope:**

- The manifest, and writing it at both `init` and `upgrade-assets`
- Conflict handling and `.new` / `.bak` emission for skills and AGENTS.md
- Wiring `--force`
- Backward-compatibility fixtures and the router format-contract test

**Out of scope:**

- Migrating `router.md` or any other `.savepoint/` file into a new shape
- The policy-asset install allowlist for `Guardrails.md` and `Health-Check.md` — E39 owns it
- Moving code-style rules out of AGENTS.md — E39 owns it
- Three-way or line-level merging of skill files; a partially merged skill is a broken workflow
- Any change to what templates contain

## Dependencies

Builds after E39. Both edit `internal/init/upgrade.go`, and E39's policy-asset allowlist changes the same walk this epic makes manifest-aware.

## Quality gates

- `make build && make test` passes.
