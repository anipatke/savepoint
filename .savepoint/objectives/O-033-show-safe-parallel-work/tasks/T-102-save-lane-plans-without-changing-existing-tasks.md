---
id: T-102
title: "Save lane plans without changing existing Tasks"
objective: O-033
status: planned
depends_on: []
owner_validation: {required: false}
planned_by: {role: planner, session: codex-o033-planning-2026-10-03}
complexity_tier: high
complexity_reason: "Optional advisory metadata validation must preserve authored records and avoid blocking ordinary lifecycle behavior."
lane: core
planned_reads:
  - "internal/data/task_v2.go"
  - "internal/data/task_v2_test.go"
  - "internal/data/objective_v2.go"
  - "internal/data/objective_v2_test.go"
  - "internal/data/project.go"
  - "internal/data/project_test.go"
  - "internal/data/write.go"
  - "internal/data/write_splice.go"
  - "internal/data/write_test.go"
  - "internal/data/dependency.go"
planned_writes:
  - "internal/data/task_v2.go"
  - "internal/data/task_v2_test.go"
  - "internal/data/objective_v2.go"
  - "internal/data/objective_v2_test.go"
  - "internal/data/project.go"
  - "internal/data/project_test.go"
  - "internal/data/write_test.go"
  - "internal/data/concurrency_plan_v2.go"
  - "internal/data/concurrency_plan_v2_test.go"
---

# Save lane plans without changing existing Tasks

## Outcome

V2 records retain optional named lanes, Task read/write manifests and narrowly scoped read/write independence explanations, while old records load unchanged and unusable advisory plans name the offending record without invalidating ordinary execution.

## User Check

Load old and newly planned temporary projects; inspect stable lane membership, missing versus explicitly empty scopes, and actionable metadata errors.

Owner validation is not required for this technical Task; owner Task completion remains required.

## Done When

- Malformed/unresolvable advisory lane metadata produces nonfatal diagnostics and no affected recommendation; otherwise valid records still load and existing lifecycle gates remain available. Required record/schema validation retains existing failures.
- Objective lanes have unique Objective-local keys and readable titles; Task lane references resolve only within their owner, with no duplicated membership list or new global identity.
- Optional planned_reads/planned_writes preserve omitted versus explicitly empty lists. Advisory validation of exact portable project-relative paths diagnoses globs, directories, absolute/drive/UNC paths, traversal and ambiguous case aliases; anticipated new files remain valid without filesystem probing.
- Independence explanations identify two owned Tasks, exact read/write overlaps and a non-empty reason. Bind their validity to reviewed manifests so edits cannot silently broaden the exception; missing references, duplicate declarations and wrong overlap types produce advisory diagnostics and suppress affected recommendations without blocking otherwise valid records or lifecycle actions.
- Existing Tasks without metadata strict-load successfully and remain unknown for concurrency. Status and Objective priority/rank writes preserve new blocks, unknown keys and authored bodies byte-for-byte outside their managed fields.
- Expose typed immutable planning values and pure validation helpers for the next Task. No concurrency scheduling, board changes, automatic backfill or new schema version.
- Per-criterion evidence and the configured handoff gate are recorded. Optional Task Checks require a fresh session; skipping one requires an explicit owner waiver naming this Task, reason, actor and time. All work remains subject to the mandatory Full Objective Check.

## Context Files

- `internal/data/task_v2.go`
- `internal/data/task_v2_test.go`
- `internal/data/objective_v2.go`
- `internal/data/objective_v2_test.go`
- `internal/data/project.go`
- `internal/data/project_test.go`
- `internal/data/write.go`
- `internal/data/write_splice.go`
- `internal/data/write_test.go`
- `internal/data/dependency.go`

Explicitly new write targets: `internal/data/concurrency_plan_v2.go`, `internal/data/concurrency_plan_v2_test.go`. Their absence before this Task is intentional, not a missing-context defect.

## Design References

O-033 Confirmed Design — Revised 2026-10-03, including Persisted plan and ownership, Canonical safety projection, Board/resume/handoff, and Guidance and verification. Project Design sections 1, 4, 5, 7, 8 and 9 provide implemented boundaries; Objective deltas govern the planned feature.

## Guardrails

FS-01, FS-02, FS-04, FS-05, DATA-01, DATA-02, DATA-03, DATA-05, TPL-01, TPL-02, TPL-04, ARCH-01, ARCH-02, ARCH-03, CFG-02, CFG-03, TEST-01, TEST-02, TEST-03, TEST-04, TEST-08, TEST-09, STYLE-01..10.

## Implementation Plan

1. Extend typed optional decoders and add focused planning metadata helpers in the explicitly new concurrency_plan_v2.go file.
2. Collect nonblocking owner/lane/explanation diagnostics after records are indexed; keep optional metadata errors separate from fatal record/schema validation. Keep existing dependency/gate vocabulary authoritative.
3. Define reviewed-scope binding for independence explanations and a conservative portable lexical path policy; unknown symlink aliases must never be certified by lexical equality alone.
4. Add temporary-project parser and managed-write regressions for absence, empty scopes, new paths, malformed YAML, duplicate keys, wrong owner and content preservation. Change existing production writers only if a regression proves they fail preservation; log the targeted scope extension before doing so.

## Boundaries

This Task implements its scoped plan and records execution evidence and authorised lifecycle updates. planned_writes describes anticipated work, not an enforcement allowlist. Necessary targeted reads must be logged before access. A difference from advisory manifests may be noted but is not itself a blocker or reason for REPLAN REQUIRED; existing material acceptance/architecture gaps still use the ordinary replan policy. Do not update another Task, allocate identity records in a lane, automate worktree/branch creation, push, merge, change Goal selection, add lifecycle states, or write a Check. In a lane skip router writes and commit locally. Set the Objective in_progress only when starting its first Task; later lanes must not rewrite the shared Objective simply to restate that status. Design reconciliation belongs to the final integration Task on main, not parallel lanes.

## Technical Verification

Fresh make test-full for path-sensitive metadata work; record Go toolchain/time and results. Full Objective Check additionally requires the native windows-tests CI job, produced by repository CI and supplied by the owner.

Use focused make test-focused TEST=... only for iteration. Record named happy-path and failure cases for each Done When criterion, actual commands/time/toolchain/results, reads and changes, and limitations. Independent verification follows agent-skills/references/check-method.md; no executor session performs its own Check.

## Technical Evidence

Pending execution. Record per-criterion evidence, extra-read reasons, actual scope versus manifests, gate result, owner validation when required, and any explicit owner Task-check waiver. Planning evidence is not technical clearance.

## Drift Notes

Pending execution. Record material acceptance/architecture deviations under normal policy. Lane choice or anticipated-scope deviation alone is advisory and never a required replan.


## Advisory Feature Contract — Revised 2026-10-03

Parallel planning defaults off via features.parallel_planning. Code Health remains core and is unchanged. Lane metadata and read/write manifests are suggestions: the owner and all skills may ignore them even when enabled. No lane rule feeds existing start, advance, completion or Check gates, and no manifest deviation alone requires replan. Actual dependencies, acceptance criteria, Guardrails and rules for a real worktree remain authoritative. This revised Objective contract supersedes any earlier mandatory lane wording.
