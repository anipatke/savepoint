---
id: T-112
title: Keep only repo-specific rules in our AGENTS.md
objective: O-039
status: done
depends_on: []
complexity_tier: low
complexity_reason: One project-owned file edited above the managed block, plus a parity check.
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: plan-o039-20261006}
check_waiver:
    task: T-112
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-10T06:45:04Z"
---

# Keep only repo-specific rules in our AGENTS.md

## Outcome

This repository's AGENTS.md has no section above `<!-- SAVEPOINT:BEGIN -->` that repeats the managed block. What remains there is repository-specific: make gate commands, the Codebase Map, runtime gate code pointers, the build-from-source rule for `savepoint`, and Reporting to the Owner.

## User Check

Open AGENTS.md and confirm that Workflow, Skill Activation, Router Selection, Verification Policy, Required Goal Context, Terminology, Issue Capture, and Worktree Lanes each appear once, inside the managed block.

## Done When

- Every section above the managed block is either repository-specific or removed. Each removed sentence that differs from the managed block is listed in Technical Evidence with where it now lives, or why it was dropped. The `health report` human-only rule is listed for the Task "Give each shared rule one home" to move into the template.
- The managed block's bytes are unchanged (FS-02).
- `savepoint resume` and `savepoint doctor` still run cleanly.
- Token weight of AGENTS.md before and after is recorded.
- `make build && make test-fast` passes.

## Context Files

`AGENTS.md`; `templates/project-v2/AGENTS.md`; `internal/init/agents.go`; `internal/init/agents_test.go`.

## Design References

O-039 Confirmed Design; Design section 1 (V2 agent workflow assets, upgrade-assets).

## Guardrails

TPL-01, TPL-02, TPL-04, FS-02, POL-01, POL-02, TEST-01, TEST-02, TEST-03, TEST-08, STYLE-07, STYLE-10.

## Implementation Plan

1. Record AGENTS.md token weight.
2. Diff each section above the block against its managed counterpart; keep repository-specific text and remove the rest.
3. Note every dropped difference and where it now lives.
4. Run resume, doctor, and the gate; record the after weight.

## Boundaries

No edits inside the managed block or to the template; no rule changes.

## Technical Verification

Focused tests during iteration; handoff gate from AGENTS.md's Verification Policy.

## Technical Evidence

Executor session: T-112 start, 2026-10-10. Not a Check; claims for a fresh `savepoint-check` to verify.

**Per-criterion outcomes**

- Sections above the block are repository-specific or removed: AGENTS.md above `<!-- SAVEPOINT:BEGIN -->` now holds Build (make gates), Codebase Map, Runtime Gates (`ResolveTaskDependencyV2`, `CheckWaiver`), Building `savepoint` (build-from-source rule) and Reporting to the Owner. Removed: Workflow, Skill Activation, Router Selection, Verification Policy, Required Goal Context, Terminology, Issue Capture, Worktree Lanes, Code Style, Context Budget, CLI Rules.
- Managed block bytes unchanged (FS-02): sha256 of the old lines 180-337 and of the new block from the BEGIN marker to the end are both `56b9af9d6fc3d73efb010c33060f469c1487a6106467e60990954260f1b64cd1`.
- `./savepoint resume` runs and returns `Next action: Advance Task T-112 to its next stage.`; `./savepoint doctor` reports ALL CLEAN (exit 0).
- Token weight (AGENTS.md; words / bytes): before 5125 / 34086, after 3171 / 21263 (-38% words, -38% bytes). Counted with `wc`, not a model tokenizer.
- `make build && make test-fast` passes (exit 0) after the test change below.

**Removed sentences that differed from the managed block, and where each now lives**

- Worktree Lanes: "When `features.parallel_planning` is on, `savepoint-design` may suggest lanes and read/write manifests; they are advisory, owners and agents may ignore them, ..." Lives in `templates/project-v2/AGENTS.md` only. The live managed block still has the older sentence ("`savepoint-design` shapes Tasks for lanes where practical ..."), so the live block lags the template; refreshing it is an upgrade-assets matter, not done here (FS-02).
- Verification Policy: the `make` gate commands and "CI runs the full gate with `make ci`" / "A Full Objective Check requires current successful `make test-full` evidence" stay, in Build. "`ResolveTaskDependencyV2` ... decides" stays, in Runtime Gates. "V2 equivalent of the epic-level integration gate" was dropped; the block says "V2 higher-level integration gate". The block adds a "Gate commands are project-owned" bullet the repository copy lacked.
- Code Style: "never block on their own" dropped; the block says "never block a Check by themselves".
- CLI Rules: the `savepoint health report [dir]` human-only rule (also mentioned in the repository half of the `health setup` sentence) was dropped. **Hand-off to the Task "Give each shared rule one home" (T-113): move this rule into `templates/project-v2/AGENTS.md`.** Until then it is in neither place; it remains in the Codebase Map row for `internal/healthcheck/`. The build-from-source rule was kept.
- Existing Codebase Adoption exists only in the managed block, so nothing was lost there.
- All other sections (Workflow, Skill Activation, Router Selection, Required Goal Context, Terminology, Issue Capture, Context Budget) read as repeats of the block; I compared them by eye, not with a line-exact tool, so a Check should re-diff them.

**Files read:** `AGENTS.md`, `agent-skills/savepoint-task/SKILL.md`, `.savepoint/router.md`, this Task, `Objective.md` of O-039, `internal/init/template_freshness_test.go` (lines 485-530). Context Files `templates/project-v2/AGENTS.md` (grep and sed only), `internal/init/agents.go` and `agents_test.go` were not opened.

**Files changed:** `AGENTS.md`; this Task (status, stage); `Objective.md` (`planned` -> `in_progress`); `internal/init/template_freshness_test.go`.

**Extra read / out-of-list edit:** `internal/init/template_freshness_test.go` is not a Context File. `make test-fast` failed in `TestGuidanceKeepsParallelPlanningAdvisory` because it required the advisory sentence in the live AGENTS.md, where only the removed section had it. I limited that assertion to the template copy, which still holds it. No rule changed.

**Limitations:** token weight is word and byte count, not a tokenizer. `make test-full` was not run (not migration or platform sensitive). No Task Check requested or waived.

## Drift Notes

The live managed block lags `templates/project-v2/AGENTS.md` (Worktree Lanes sentence). See Technical Evidence.
