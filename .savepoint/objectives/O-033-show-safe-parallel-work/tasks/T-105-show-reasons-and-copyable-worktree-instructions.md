---
id: T-105
title: "Show reasons and copyable worktree instructions"
objective: O-033
status: planned
depends_on: [{task: T-103, requires: clear}, {task: T-104, requires: clear}]
owner_validation: {required: true}
planned_by: {role: planner, session: codex-o033-planning-2026-10-03}
complexity_tier: medium
complexity_reason: "Shared narrative must remain consistent across resume, details and plain board output without introducing clipboard or shell execution."
lane: board
planned_reads:
  - "internal/data/feature_preferences.go"
  - "internal/resume/resume.go"
  - "internal/resume/resume_test.go"
  - "internal/resume/evidence.go"
  - "internal/board/v2/detail.go"
  - "internal/board/v2/detail_view.go"
  - "internal/board/v2/detail_test.go"
  - "internal/board/v2/plain.go"
  - "internal/board/v2/next_panel.go"
  - "internal/board/v2/next_panel_test.go"
  - "internal/board/v2/run_test.go"
  - "internal/data/concurrency_v2.go"
  - "internal/data/next.go"
planned_writes:
  - "internal/resume/concurrency.go"
  - "internal/resume/concurrency_test.go"
  - "internal/resume/resume.go"
  - "internal/resume/resume_test.go"
  - "internal/board/v2/detail.go"
  - "internal/board/v2/detail_view.go"
  - "internal/board/v2/detail_test.go"
  - "internal/board/v2/plain.go"
  - "internal/board/v2/next_panel.go"
  - "internal/board/v2/next_panel_test.go"
  - "internal/board/v2/run_test.go"
---

# Show reasons and copyable worktree instructions

## Outcome

Owners can read the same parallel opportunities and sequential reasons on every surface and copy a complete fresh-session instruction for each eligible Task.

## User Check

Read an Objective's detail, a Task's detail, resume and piped board output; compare opportunities/reasons and paste a complete Task instruction into a text editor.

Owner validation is required after technical evidence; do not infer it from tests.

## Done When

- Off/default removes optional lane reasons and instructions from resume/details/plain board without deleting their records. On presents them as ignorable suggestions; the owner may run on main or another worktree and no lane rule becomes a completion condition.
- Use the canonical projection and one shared plain-text formatter for reasons and prompt text. Resume, Objective/Task detail and non-TTY board agree for the same selection and evidence.
- Every eligible parallel Task instruction includes Objective/Task/lane, current Start line, read/write scope, prerequisites and owner-prepared worktree setup with prerequisite changes available.
- Optional instructions use savepoint-task. If the owner chooses an actual worktree, preserve router/Goal, forbid identity allocation there, require Task evidence/local commit and defer merge/Checks to main under existing worktree rules. Instructions explicitly permit normal execution without a lane/worktree. No guessed actual path, automatic branch/worktree creation or clipboard dependency.
- Blocked or unknown Tasks state constraints and are not advertised with actionable parallel Start instructions. Project record strings cannot inject terminal controls into output.
- NextLine remains unchanged in meaning and continues to name selected work. Goal-wide board output can show namespaced lane membership but no cross-Objective opportunity.
- Per-criterion evidence and the configured handoff gate are recorded. Optional Task Checks require a fresh session; skipping one requires an explicit owner waiver naming this Task, reason, actor and time. All work remains subject to the mandatory Full Objective Check.

## Context Files

- `internal/data/feature_preferences.go`
- `internal/resume/resume.go`
- `internal/resume/resume_test.go`
- `internal/resume/evidence.go`
- `internal/board/v2/detail.go`
- `internal/board/v2/detail_view.go`
- `internal/board/v2/detail_test.go`
- `internal/board/v2/plain.go`
- `internal/board/v2/next_panel.go`
- `internal/board/v2/next_panel_test.go`
- `internal/board/v2/run_test.go`
- `internal/data/concurrency_v2.go`
- `internal/data/next.go`

Prerequisite-created context: `internal/data/concurrency_v2.go`. These files are intentionally created by preceding Tasks; verify them after dependencies are met. Other Context Files must exist. A material interface mismatch returns REPLAN REQUIRED.

Explicitly new write targets: `internal/resume/concurrency.go`, `internal/resume/concurrency_test.go`. Their absence before this Task is intentional, not a missing-context defect.

## Design References

O-033 Confirmed Design — Revised 2026-10-03, including Persisted plan and ownership, Canonical safety projection, Board/resume/handoff, and Guidance and verification. Project Design sections 1, 4, 5, 7, 8 and 9 provide implemented boundaries; Objective deltas govern the planned feature.

## Guardrails

FS-01, FS-02, FS-04, FS-05, DATA-01, DATA-02, DATA-03, DATA-05, TPL-01, TPL-02, TPL-04, ARCH-01, ARCH-02, ARCH-03, CFG-02, CFG-03, TEST-01, TEST-02, TEST-03, TEST-04, TEST-08, TEST-09, STYLE-01..10.

## Implementation Plan

1. Add pure shared reason/instruction renderers in explicitly new internal/resume/concurrency.go, keeping safety exclusively in internal/data.
2. Wire resume narrative, detail projections and plain board sections; use the selected Objective without treating board focus as a router write.
3. Preserve wrapping/scrolling and control sanitisation in existing render boundaries.
4. Validate exact-content parity and writer failure, unavailable manifests, dependency blocks, special characters and full instruction completeness using temporary-project fixtures.

## Boundaries

This Task implements its scoped plan and records execution evidence and authorised lifecycle updates. planned_writes describes anticipated work, not an enforcement allowlist. Necessary targeted reads must be logged before access. A difference from advisory manifests may be noted but is not itself a blocker or reason for REPLAN REQUIRED; existing material acceptance/architecture gaps still use the ordinary replan policy. Do not update another Task, allocate identity records in a lane, automate worktree/branch creation, push, merge, change Goal selection, add lifecycle states, or write a Check. In a lane skip router writes and commit locally. Set the Objective in_progress only when starting its first Task; later lanes must not rewrite the shared Objective simply to restate that status. Design reconciliation belongs to the final integration Task on main, not parallel lanes.

## Technical Verification

make build && make test-fast; record text parity and copy/paste scenario validation.

Use focused make test-focused TEST=... only for iteration. Record named happy-path and failure cases for each Done When criterion, actual commands/time/toolchain/results, reads and changes, and limitations. Independent verification follows agent-skills/references/check-method.md; no executor session performs its own Check.

## Technical Evidence

Pending execution. Record per-criterion evidence, extra-read reasons, actual scope versus manifests, gate result, owner validation when required, and any explicit owner Task-check waiver. Planning evidence is not technical clearance.

## Drift Notes

Pending execution. Record material acceptance/architecture deviations under normal policy. Lane choice or anticipated-scope deviation alone is advisory and never a required replan.


## Advisory Feature Contract — Revised 2026-10-03

Parallel planning defaults off via features.parallel_planning. Code Health remains core and is unchanged. Lane metadata and read/write manifests are suggestions: the owner and all skills may ignore them even when enabled. No lane rule feeds existing start, advance, completion or Check gates, and no manifest deviation alone requires replan. Actual dependencies, acceptance criteria, Guardrails and rules for a real worktree remain authoritative. This revised Objective contract supersedes any earlier mandatory lane wording.
