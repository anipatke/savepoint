---
id: T-115
title: Load the Savepoint guide in Claude Code automatically
objective: O-043
status: done
depends_on: []
complexity_tier: medium
complexity_reason: New managed-region file merged into a possibly existing user CLAUDE.md through init and upgrade, reusing the AGENTS.md merge pattern.
owner_validation:
    required: true
    accepted_check: ""
planned_by: {role: planner, session: plan-o043-20261010}
check_waiver:
    task: T-115
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-10T07:27:09Z"
---

# Load the Savepoint guide in Claude Code automatically

## Outcome

`savepoint init` and `savepoint upgrade-assets` give the project a `CLAUDE.md` whose Savepoint-managed region imports the agent guide with `@AGENTS.md`. Claude Code then loads the guide at session start with no pasted prompt and no extra tool call. A user's own `CLAUDE.md` text outside the region is never changed.

## User Check

Run `savepoint init` in an empty scratch Git repo, open Claude Code there, and ask "what should I do next?". Claude should mention Savepoint and its Next step without being told to read anything. Repeat in a scratch repo that already has a `CLAUDE.md` with your own line in it, and confirm that line is still there.

## Done When

- A fresh init writes `CLAUDE.md` with a managed region containing the `@AGENTS.md` import and one plain sentence saying Savepoint manages that region.
- Init and upgrade into a project with an existing `CLAUDE.md` add or refresh only the managed region. Bytes outside it are unchanged (FS-01, FS-02, TEST-03).
- If the existing `CLAUDE.md` already imports `@AGENTS.md` outside the region, no second import is added.
- A second upgrade on unchanged input reports unchanged (FS-04). Dry run writes nothing (FS-03).
- This repository's `CLAUDE.md` is regenerated the same way, replacing the "See `AGENTS.md`" line.
- Session-start token weight (guide plus `CLAUDE.md`) is recorded before and after, for this repo and a fresh project.
- `make build && make test-fast` passes.

## Context Files

`main.go`; `internal/init/scaffold.go`; `internal/init/scaffold_test.go`; `internal/init/agents.go`; `internal/init/agents_test.go`; `internal/init/upgrade.go`; `internal/init/upgrade_test.go`; `internal/init/manifest.go`; `templates/project-v2/AGENTS.md`; `CLAUDE.md`.

## Design References

O-043 Architectural Considerations and Confirmed Design; Design section 1 (V2 agent workflow assets, upgrade-assets) and section 12.

## Guardrails

FS-01, FS-02, FS-03, FS-04, FS-05, FS-06, TPL-02, TPL-04, ARCH-04, CFG-03, TEST-01, TEST-02, TEST-03, TEST-04, TEST-08.

## Implementation Plan

1. Record current session-start token weight for this repo and a fresh init project.
2. Add `templates/project-v2/CLAUDE.md` with the managed region.
3. Generalise the AGENTS.md managed-region merge so `CLAUDE.md` uses the same code path, or add a sibling with shared helpers. Keep the Codebase Map accurate (ARCH-04).
4. Wire it into init and upgrade-assets, including dry run and the existing-import check.
5. Test fresh, existing-file, already-imported, re-run, and dry-run cases, plus one unwritable-target failure.
6. Regenerate this repo's `CLAUDE.md`; record the after weights.

## Boundaries

No change to AGENTS.md content, the skills, or `.claude/`. No support for other agents' files (`GEMINI.md`).

## Technical Verification

Focused tests during iteration; handoff gate from AGENTS.md's Verification Policy. Native Windows CI must pass for the merge path (CFG-03).

## Technical Evidence

- Criteria: fresh init writes CLAUDE.md with the managed `@AGENTS.md` block (verified by scratch `savepoint init`, and `TestScaffold_writesClaudeGuideWithImport`). Existing CLAUDE.md keeps its bytes outside the block on init and upgrade (`TestScaffold_keepsUserClaudeText`, `TestUpgradeProjectAssets_claudeGuideLifecycle`, `..._RefreshesOnlyBlock`). A hand-written `@AGENTS.md` import adds no second one (`TestScaffold_noSecondImport`, `..._AlreadyImportedUnchanged`). Second upgrade reports unchanged and dry run writes nothing (lifecycle test; scratch `upgrade-assets` printed `unchanged CLAUDE.md`). This repo's CLAUDE.md regenerated from the template. One failure case: CLAUDE.md unreadable -> `failed` entry plus error.
- Design choice: unlike AGENTS.md, an unmarked CLAUDE.md is merged (block appended), not a conflict, because the user's bytes stay untouched; `--force` has no extra effect for it.
- Session-start weight (bytes; tokens roughly bytes/4): before, this repo CLAUDE.md 37 B and AGENTS.md not auto-loaded (~0 guide tokens); after, 168 B + AGENTS.md 21,539 B = ~21.7 KB (~5.4k tokens). Fresh project: before ~0; after 168 B + 15,016 B = 15,184 B (~3.8k tokens). Estimate only, no tokenizer run.
- Commands: `make build`, `go test ./internal/init/`, `make test-fast` (passed). Not run: `make test-full`, native Windows CI (CFG-03), a live Claude Code session.
- Files changed: templates/project-v2/CLAUDE.md, internal/init/agents.go, scaffold.go, upgrade.go, claude_guide_test.go, CLAUDE.md, AGENTS.md (Codebase Map row). Files read: only Context Files plus upgrade.go helpers already listed. Extra reads: internal/init/upgrade_test.go helpers (to reuse test fixtures).
- Limitation: the import is `@AGENTS.md`; a project whose guide has a different casing (agents.md) on a case-sensitive disk would not resolve it.

## Drift Notes

None yet.
