---
id: T-104
title: "Keep lane headings as Tasks move columns"
objective: O-033
status: planned
depends_on: [{task: T-103, requires: clear}]
owner_validation: {required: true}
planned_by: {role: planner, session: codex-o033-planning-2026-10-03}
complexity_tier: high
complexity_reason: "Grouping must preserve card navigation, counts, reload selection and measured terminal geometry."
lane: board
planned_reads:
  - "internal/board/v2/card.go"
  - "internal/board/v2/card_test.go"
  - "internal/board/v2/column.go"
  - "internal/board/v2/column_test.go"
  - "internal/board/v2/view.go"
  - "internal/board/v2/view_test.go"
  - "internal/board/v2/columns_view_test.go"
  - "internal/board/v2/model.go"
  - "internal/board/v2/update.go"
  - "internal/board/v2/load.go"
  - "internal/board/v2/load_test.go"
  - "internal/board/v2/width.go"
  - ".savepoint/visual-identity.md"
  - "internal/data/concurrency_v2.go"
planned_writes:
  - "internal/board/v2/card.go"
  - "internal/board/v2/card_test.go"
  - "internal/board/v2/column.go"
  - "internal/board/v2/column_test.go"
  - "internal/board/v2/view.go"
  - "internal/board/v2/view_test.go"
  - "internal/board/v2/columns_view_test.go"
  - "internal/board/v2/model.go"
  - "internal/board/v2/update.go"
  - "internal/board/v2/load.go"
  - "internal/board/v2/load_test.go"
---

# Keep lane headings as Tasks move columns

## Outcome

The existing three status columns group Task cards under stable Lane / Proposed worktree headings with concise canonical readiness wording.

## User Check

Move sample Tasks from Planned to In Progress to Done; verify the same heading persists wherever that lane has cards. Navigate and scroll grouped columns in compact and wide terminals.

Owner validation is required after technical evidence; do not infer it from tests.

## Done When

- With the option off/default, render the existing ungrouped board even if saved lanes exist. Turning it on/off preserves task state and metadata. Ignoring the displayed grouping never refuses a status action, and Code Health remains unchanged.
- Lane headings follow saved membership across all three columns; omit empty headings. Legacy no-lane projects keep their existing layout, and mixed projects expose ungrouped sequential work.
- Headings are not selectable Tasks and never affect counts or status writes. Preserve focused Task by identity when grouping, moving status or reloading changes its index.
- Stable lane order and deterministic Task order agree across columns. Goal-wide views namespace headings by Objective and do not suggest cross-Objective parallelism.
- Readiness indications consume the canonical data projection; rendering performs no IO or safety decisions.
- Focused card stays visible with heading height included in window budgeting, including a window beginning mid-lane, long titles, compact view and narrow supported widths. Existing badges, palette and focus geometry remain consistent.
- Per-criterion evidence and the configured handoff gate are recorded. Optional Task Checks require a fresh session; skipping one requires an explicit owner waiver naming this Task, reason, actor and time. All work remains subject to the mandatory Full Objective Check.

## Context Files

- `internal/board/v2/card.go`
- `internal/board/v2/card_test.go`
- `internal/board/v2/column.go`
- `internal/board/v2/column_test.go`
- `internal/board/v2/view.go`
- `internal/board/v2/view_test.go`
- `internal/board/v2/columns_view_test.go`
- `internal/board/v2/model.go`
- `internal/board/v2/update.go`
- `internal/board/v2/load.go`
- `internal/board/v2/load_test.go`
- `internal/board/v2/width.go`
- `.savepoint/visual-identity.md`
- `internal/data/concurrency_v2.go`

Prerequisite-created context: `internal/data/concurrency_v2.go`. These files are intentionally created by preceding Tasks; verify them after dependencies are met. Other Context Files must exist. A material interface mismatch returns REPLAN REQUIRED.

## Design References

O-033 Confirmed Design — Revised 2026-10-03, including Persisted plan and ownership, Canonical safety projection, Board/resume/handoff, and Guidance and verification. Project Design sections 1, 4, 5, 7, 8 and 9 provide implemented boundaries; Objective deltas govern the planned feature.

## Guardrails

FS-01, FS-02, FS-04, FS-05, DATA-01, DATA-02, DATA-03, DATA-05, TPL-01, TPL-02, TPL-04, ARCH-01, ARCH-02, ARCH-03, CFG-02, CFG-03, TEST-01, TEST-02, TEST-03, TEST-04, TEST-08, TEST-09, STYLE-01..10.

## Implementation Plan

1. Resolve lane/readiness values outside rendering using the predecessor's canonical projection; group cards without creating a new lifecycle collection.
2. Render lane headings as non-card rows, accounting for their wrapped height and repeating the context heading when scrolling into a lane.
3. Adjust identity-based selection restoration only where grouping needs it; keep all filesystem operations behind existing load/action commands.
4. Add outcome regressions for status moves, empty/mixed lanes, duplicate lane names in different Objectives, card counts, keyboard focus and scrolling. Log any necessary extra read before accessing it.

## Boundaries

This Task implements its scoped plan and records execution evidence and authorised lifecycle updates. planned_writes describes anticipated work, not an enforcement allowlist. Necessary targeted reads must be logged before access. A difference from advisory manifests may be noted but is not itself a blocker or reason for REPLAN REQUIRED; existing material acceptance/architecture gaps still use the ordinary replan policy. Do not update another Task, allocate identity records in a lane, automate worktree/branch creation, push, merge, change Goal selection, add lifecycle states, or write a Check. In a lane skip router writes and commit locally. Set the Objective in_progress only when starting its first Task; later lanes must not rewrite the shared Objective simply to restate that status. Design reconciliation belongs to the final integration Task on main, not parallel lanes.

## Technical Verification

make build && make test-fast; record supported terminal-size scenarios. Native Windows CI remains mandatory at the Full Objective Check.

Use focused make test-focused TEST=... only for iteration. Record named happy-path and failure cases for each Done When criterion, actual commands/time/toolchain/results, reads and changes, and limitations. Independent verification follows agent-skills/references/check-method.md; no executor session performs its own Check.

## Technical Evidence

Pending execution. Record per-criterion evidence, extra-read reasons, actual scope versus manifests, gate result, owner validation when required, and any explicit owner Task-check waiver. Planning evidence is not technical clearance.

## Drift Notes

Pending execution. Record material acceptance/architecture deviations under normal policy. Lane choice or anticipated-scope deviation alone is advisory and never a required replan.


## Advisory Feature Contract — Revised 2026-10-03

Parallel planning defaults off via features.parallel_planning. Code Health remains core and is unchanged. Lane metadata and read/write manifests are suggestions: the owner and all skills may ignore them even when enabled. No lane rule feeds existing start, advance, completion or Check gates, and no manifest deviation alone requires replan. Actual dependencies, acceptance criteria, Guardrails and rules for a real worktree remain authoritative. This revised Objective contract supersedes any earlier mandatory lane wording.
