---
id: E40-upgrade-safety/T002-skill-conflict-policy
status: done
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

- [x] `UpgradeAction` gains `ActionConflict = "conflict"`, and `Format()` reports conflicts in the summary counts and per-path list.
- [x] Skill file missing on disk → written, hash recorded, reported `updated`.
- [x] Skill file byte-identical to the incoming template → hash recorded, reported `unchanged`, file untouched.
- [x] Skill file differs and its hash matches the manifest (unmodified, merely outdated) → overwritten, hash recorded, reported `updated`.
- [x] Skill file differs and its hash does **not** match the manifest (user-customized) → the user's file is left byte-identical, the incoming version is written to `<path>.new`, and the result is reported `conflict`.
- [x] Skill file differs and has **no** manifest entry (project upgrading from a pre-manifest version) → existing content saved to `<path>.bak`, then overwritten, hash recorded, reported `updated`.
- [x] `--force` overrides the conflict case only: the existing content is saved to `<path>.bak`, the file is overwritten, the hash is recorded, and the result is reported `updated`. `force` is read in `UpgradeProjectAssets`, not just accepted.
- [x] Dry run reports identical actions for every case above and writes nothing — no `.new`, no `.bak`, no manifest, no target file.
- [x] A second consecutive upgrade with no intervening edits is fully idempotent: all `unchanged`, and no `.new` or `.bak` files created.
- [x] Re-running upgrade against an unresolved conflict re-reports `conflict` and overwrites the stale `<path>.new`, rather than accumulating variants.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Add `ActionConflict` and extend `Format()` counts and output.
- [x] Extract the skill branch of the walk into a named function taking the manifest, `dryRun`, and `force` — the current inline block is already the longest in the walk.
- [x] Implement the decision table above, ordering checks so the cheap byte-equality test short-circuits before hashing.
- [x] Read `force` in `UpgradeProjectAssets` and thread it to the skill path.
- [x] Extend `upgrade_test.go` with one test per AC row, including the dry-run and idempotency cases.

## Context Log

**Read:** `.savepoint/router.md`, `AGENTS.md`, `agent-skills/savepoint-build-task/SKILL.md`, `.savepoint/Guardrails.md` (STYLE), `E40-Detail.md`, `T001-upgrade-manifest.md`, this task, `internal/init/upgrade.go`, `internal/init/manifest.go`, `internal/init/write.go`, `internal/init/upgrade_test.go`, `internal/init/manifest_test.go`, `internal/testutil/fs.go`, `main.go` (upgrade wiring).

**Edited:** `internal/init/upgrade.go`, `internal/init/upgrade_test.go`.

**Design notes:**

- The skill branch moved *above* the dry-run block and became `upgradeSkillAsset`, which takes `dryRun` and gates every write internally. Dry-run parity is therefore structural — one decision path, not two that must be kept in step. What remained in the dry-run block is AGENTS.md only.
- `UpgradeEntry` gained a `Note` so a `.bak` or `.new` sidecar is named in the report. The pre-manifest migration reports `updated` like any other replacement, so without the note the backup would be invisible; the epic requires it be stated.
- Conflict policy is scoped to manifest-covered paths (`agent-skills/*/SKILL.md`). Shared references under `agent-skills/references/` stay package-owned and refresh directly, per the epic lifecycle matrix — routing them through the untracked branch would emit a `.bak` on every upgrade that changed them.
- Decision order is byte-equality, then manifest scope, then provenance, so the cheap comparison short-circuits before any hashing.
- A conflict leaves the manifest entry at the old hash: the file on disk is the user's, not ours, so recording it would falsely claim provenance and silently convert the next upgrade into an overwrite.
- `.new` is written to a fixed path, so a re-run replaces the stale sidecar instead of accumulating `.new.new` variants.
- `main.go:106` already passed `opts.Force`; `UpgradeProjectAssets` now threads it to the skill path, so `--force` stops being inert with no CLI change.

**Acceptance evidence** (`internal/init/upgrade_test.go`):

| AC | Test |
|---|---|
| `ActionConflict` in counts and per-path list | `TestUpgradeReport_formatConflict` |
| Missing → written, hashed, `updated` | `TestUpgradeSkill_missingFileIsWritten` |
| Identical → `unchanged`, untouched | `TestUpgradeSkill_identicalFileIsUnchanged` |
| Tracked + outdated → replaced, `updated` | `TestUpgradeSkill_trackedAndOutdatedIsReplaced` |
| Customized → kept, `.new` written, `conflict` | `TestUpgradeSkill_customizedFileConflicts` |
| No manifest entry → `.bak`, replaced, `updated` | `TestUpgradeSkill_untrackedFileIsBackedUpAndReplaced` |
| `--force` overrides conflict only | `TestUpgradeSkill_forceOverridesConflict`, `TestUpgradeSkill_forceLeavesUnrelatedCasesAlone` |
| Dry run: same action, writes nothing | `TestUpgradeSkill_dryRunMatchesRealRunAndWritesNothing` (6 cases, compared against a real run of each) |
| Second run idempotent | `TestUpgradeSkill_secondRunIsIdempotent` |
| Unresolved conflict re-reports, no variants | `TestUpgradeSkill_unresolvedConflictRefreshesTheSameSidecar` |

Pre-existing skill tests (`updatesAgentSkills`, `skillIdempotent`, `dryRunDoesNotWrite`, `multipleSkills`, `createsMissingSkillFile`) pass unchanged, confirming the untracked-project path keeps prior behaviour plus a backup.

**Quality gates:** `make build && make test` — pass, all packages ok.

**Health check:** skipped; this project has no `.savepoint/Health-Check.md`.

## Notes

The pre-manifest migration path is the one place this epic accepts replacing possibly-customized content. It is one-time and recoverable via the `.bak`, and avoids emitting a `.new` for every skill on the first upgrade after this ships. The upgrade report must make the backups visible. If that trade proves wrong in practice, switching to conflict-on-unknown is a single branch change.
