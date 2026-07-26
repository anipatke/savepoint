---
id: E40-upgrade-safety/T002-skill-conflict-policy
status: planned
objective: Stop upgrade silently overwriting customized skill files — keep the user's version, write the incoming one as SKILL.md.new, report a conflict, and give --force real meaning.
depends_on: ["E40-upgrade-safety/T001-upgrade-manifest"]
complexity_tier: medium
complexity_reason: Branching rewrite of the skill upgrade path plus a new action, migration fallback, and dry-run parity.
---

# T002: Skill Conflict Policy

## Problem

`upgrade.go:179-205` replaces every `agent-skills/*/SKILL.md` whose content differs from the shipped copy, with no merge, no backup, and no warning, reported as a routine `updated`. Any project that tailored a skill loses that work.

Skill files are user-facing by design: `AGENTS.md:24` tells agents to read `agent-skills/{skill}/SKILL.md` directly when the skill tool cannot find a skill, and E39 has skills referencing project-owned Guardrails rules.

Separately, `UpgradeProjectAssets` accepts `force bool` (`upgrade.go:71`) and never reads it. `main.go:106` passes `opts.Force`, so `--force` is silently inert.

## Context Files

- `internal/init/upgrade.go`
- `internal/init/upgrade_test.go`
- `internal/init/manifest.go` (from T001)
- `main.go` (`--force` wiring, line 106)

## Acceptance Criteria

- [ ] `UpgradeAction` gains `ActionConflict = "conflict"`, and `Format()` reports conflicts in the summary counts and per-path list.
- [ ] Skill file missing on disk → written, hash recorded, reported `updated`.
- [ ] Skill file byte-identical to the incoming template → hash recorded, reported `unchanged`, file untouched.
- [ ] Skill file differs and its hash matches the manifest (unmodified, merely outdated) → overwritten, hash recorded, reported `updated`.
- [ ] Skill file differs and its hash does **not** match the manifest (user-customized) → the user's file is left byte-identical, the incoming version is written to `<path>.new`, and the result is reported `conflict`.
- [ ] Skill file differs and has **no** manifest entry (project upgrading from a pre-manifest version) → existing content saved to `<path>.bak`, then overwritten, hash recorded, reported `updated`.
- [ ] `--force` overrides the conflict case only: the existing content is saved to `<path>.bak`, the file is overwritten, the hash is recorded, and the result is reported `updated`. `force` is read in `UpgradeProjectAssets`, not just accepted.
- [ ] Dry run reports identical actions for every case above and writes nothing — no `.new`, no `.bak`, no manifest, no target file.
- [ ] A second consecutive upgrade with no intervening edits is fully idempotent: all `unchanged`, and no `.new` or `.bak` files created.
- [ ] Re-running upgrade against an unresolved conflict re-reports `conflict` and overwrites the stale `<path>.new`, rather than accumulating variants.
- [ ] `make build && make test` passes.

## Implementation Plan

- [ ] Add `ActionConflict` and extend `Format()` counts and output.
- [ ] Extract the skill branch of the walk into a named function taking the manifest, `dryRun`, and `force` — the current inline block is already the longest in the walk.
- [ ] Implement the decision table above, ordering checks so the cheap byte-equality test short-circuits before hashing.
- [ ] Read `force` in `UpgradeProjectAssets` and thread it to the skill path.
- [ ] Extend `upgrade_test.go` with one test per AC row, including the dry-run and idempotency cases.

## Notes

The pre-manifest migration path is the one place this epic accepts replacing possibly-customized content. It is one-time and recoverable via the `.bak`, and avoids emitting a `.new` for every skill on the first upgrade after this ships. The upgrade report must make the backups visible. If that trade proves wrong in practice, switching to conflict-on-unknown is a single branch change.
