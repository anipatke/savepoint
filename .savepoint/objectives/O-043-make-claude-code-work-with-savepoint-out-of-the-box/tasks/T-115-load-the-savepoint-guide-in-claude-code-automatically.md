---
id: T-115
title: Load the Savepoint guide in Claude Code automatically
objective: O-043
status: planned
depends_on: []
complexity_tier: medium
complexity_reason: New managed-region file merged into a possibly existing user CLAUDE.md through init and upgrade, reusing the AGENTS.md merge pattern.
owner_validation: {required: true}
planned_by: {role: planner, session: plan-o043-20261010}
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

Pending execution.

## Drift Notes

None yet.
