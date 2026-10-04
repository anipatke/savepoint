---
id: T-107
title: "Verify the complete lane planning handoff"
objective: O-033
status: done
depends_on: [{task: T-105, requires: clear}, {task: T-106, requires: clear}]
owner_validation:
  required: true
  accepted_check: ""
planned_by: {role: planner, session: codex-o033-planning-2026-10-03}
complexity_tier: medium
complexity_reason: "Cross-surface integration and public guidance must demonstrate the delivered contract without claiming independent clearance."
lane: integration
planned_reads:
  - "internal/data/feature_preferences.go"
  - "internal/data/concurrency_v2.go"
  - "internal/data/concurrency_plan_v2.go"
  - "internal/data/e2e_v2_test.go"
  - "internal/board/v2/run_test.go"
  - "internal/board/v2/fixture_test.go"
  - "internal/board/v2/plain.go"
  - "internal/resume/concurrency.go"
  - "internal/resume/resume.go"
  - "internal/init/upgrade_test.go"
  - "README.md"
  - "CHANGELOG.md"
  - ".github/workflows/ci.yml"
planned_writes:
  - "internal/data/concurrency_integration_test.go"
  - "internal/board/v2/concurrency_integration_test.go"
  - "README.md"
  - "CHANGELOG.md"
  - ".savepoint/Design.md"
check_waiver:
  task: T-107
  reason: Owner completed this Task via the board without requesting a Task Check.
  actor:
    role: owner
    session: board-owner
  recorded_at: "2026-10-04T03:13:35Z"
---

# Verify the complete lane planning handoff

## Outcome

Integrated temporary-project evidence and public guidance demonstrate stable lanes, conservative parallel suggestions and complete session handoffs, with implemented Design reconciled for the independent Objective Check.

## User Check

Follow the documented walkthrough: inspect two lanes, copy a ready Task instruction, observe a blocked lane and move Tasks across columns while headings remain stable.

Owner validation is required after technical evidence; do not infer it from tests.

## Done When

- Integrate the single-option Advanced Options screen and preference: off/default/on/restart, external edits, malformed advice degradation and save failure. Prove that identical Tasks have identical existing start/advance/completion decisions with advice off/on and when grouping/manifests are ignored; Code Health behavior stays unchanged.
- Temporary-project integration covers old/mixed/new records, sequential same-lane work, satisfied/blocked dependencies, shared writes, shared reads, read/write explanation invalidation and unknown active scope.
- Board/detail/plain/resume agree on selected-Objective opportunities, reasons and instructions, including after replan and status changes; Goal-wide views never imply cross-Objective concurrency.
- README and CHANGELOG describe optional metadata, Lane / Proposed worktree headings, conservative safety limits, manually prepared prerequisites/worktrees and owner merge/check responsibility without promising orchestration or monitoring.
- Reconcile Design to actual shipped storage/projection/render/upgrade boundaries only after implementation, preserving original intent and owner authority. Record architecture drift rather than retrofitting criteria.
- Fresh full gate and named integration cases are recorded per criterion; limitations and native Windows CI result availability are explicit. This executor Task writes no Check and grants no independent CLEAR.
- Per-criterion evidence and the configured handoff gate are recorded. Optional Task Checks require a fresh session; skipping one requires an explicit owner waiver naming this Task, reason, actor and time. All work remains subject to the mandatory Full Objective Check.

## Context Files

- `internal/data/feature_preferences.go`
- `internal/data/concurrency_v2.go`
- `internal/data/concurrency_plan_v2.go`
- `internal/data/e2e_v2_test.go`
- `internal/board/v2/run_test.go`
- `internal/board/v2/fixture_test.go`
- `internal/board/v2/plain.go`
- `internal/resume/concurrency.go`
- `internal/resume/resume.go`
- `internal/init/upgrade_test.go`
- `README.md`
- `CHANGELOG.md`
- `.github/workflows/ci.yml`

Prerequisite-created context: `internal/data/concurrency_v2.go`, `internal/data/concurrency_plan_v2.go`, `internal/resume/concurrency.go`. These files are intentionally created by preceding Tasks; verify them after dependencies are met. Other Context Files must exist. A material interface mismatch returns REPLAN REQUIRED.

Explicitly new write targets: `internal/data/concurrency_integration_test.go`, `internal/board/v2/concurrency_integration_test.go`. Their absence before this Task is intentional, not a missing-context defect.

## Design References

O-033 Confirmed Design — Revised 2026-10-03, including Persisted plan and ownership, Canonical safety projection, Board/resume/handoff, and Guidance and verification. Project Design sections 1, 4, 5, 7, 8 and 9 provide implemented boundaries; Objective deltas govern the planned feature.

## Guardrails

FS-01, FS-02, FS-04, FS-05, DATA-01, DATA-02, DATA-03, DATA-05, TPL-01, TPL-02, TPL-04, ARCH-01, ARCH-02, ARCH-03, CFG-02, CFG-03, TEST-01, TEST-02, TEST-03, TEST-04, TEST-08, TEST-09, STYLE-01..10.

## Implementation Plan

1. Confirm delivered interfaces and guidance; return REPLAN REQUIRED for an unresolved material contract.
2. Add explicitly new cross-component integration test files using existing fixture conventions and temporary projects; no real worktree/network.
3. Write the user walkthrough and release note, then reconcile implemented Design sections on lane data, projection, columns and handoff.
4. Run fresh make test-full and prepare evidence for a new independent Full Objective Check under check-method.md. Owner supplies native windows-tests CI result; record outstanding evidence without claiming clearance.

## Boundaries

This Task implements its scoped plan and records execution evidence and authorised lifecycle updates. planned_writes describes anticipated work, not an enforcement allowlist. Necessary targeted reads must be logged before access. A difference from advisory manifests may be noted but is not itself a blocker or reason for REPLAN REQUIRED; existing material acceptance/architecture gaps still use the ordinary replan policy. Do not update another Task, allocate identity records in a lane, automate worktree/branch creation, push, merge, change Goal selection, add lifecycle states, or write a Check. In a lane skip router writes and commit locally. Set the Objective in_progress only when starting its first Task; later lanes must not rewrite the shared Objective simply to restate that status. Design reconciliation belongs to the final integration Task on main, not parallel lanes.

## Technical Verification

Fresh make test-full. Repository CI produces native windows-tests evidence; owner makes it available before the mandatory Full Objective Check. Do not run health check or write a Check during this Task.

Use focused make test-focused TEST=... only for iteration. Record named happy-path and failure cases for each Done When criterion, actual commands/time/toolchain/results, reads and changes, and limitations. Independent verification follows agent-skills/references/check-method.md; no executor session performs its own Check.

## Technical Evidence

Commands: `make test-full` on 2026-10-04 14:12 +11:00, go1.26.2 linux/amd64, wall time 17.5s, all packages passed and linux/darwin/windows cross-builds succeeded. Focused iteration only: `go test ./internal/data -run TestIntegration_`, `go test ./internal/board/v2 -run TestIntegration`.

Per criterion (new files `internal/data/concurrency_integration_test.go`, `internal/board/v2/concurrency_integration_test.go`):
1. Advanced Options and decisions: option off/default/on/restart, external edit and save failure are covered by the existing `TestOptionsJourneyOnTemporaryProject` and `TestOptionsLeaveLifecycleAndBoardIdenticalWhenOn` (not rewritten here). New `TestIntegration_adviceNeverChangesExistingDecisions` loads the same Tasks with and without lanes/manifests and projects with the preference off, on, off; start, advance and completion decisions are identical, including one Task with malformed `planned_writes` (diagnostic present). `TestIntegrationMalformedAdviceDegradesWithoutHidingTheBoard` covers board degradation. Code Health untouched (no code change).
2. Records: `TestIntegration_lanesAcrossOldMixedAndNewRecords` (old/mixed/new), `TestIntegration_sequentialDependencyAndSharedPathCases` (same lane, blocked and satisfied dependency, shared write, shared read, running Task with unknown scope), `TestIntegration_explanationInvalidatesWhenEitherManifestChanges` (writer and reader broadened).
3. Surface agreement: `TestIntegrationSurfacesAgreeAfterStatusChanges` (lane head started then done), `TestIntegrationSurfacesWithholdTogetherAfterAReplan`, `TestIntegrationGoalWideViewNeverImpliesCrossObjectiveConcurrency`; pre-existing `TestParallelAdviceAgreesAcrossResumeDetailsAndPlainBoard`.
4. README gains Planning lanes and Try it sections; CHANGELOG lists metadata, headings, reasons, safety limits and manual worktree/merge responsibility. No orchestration or monitoring is promised.
5. Design.md section 8 gains a shipped "Parallel planning (O-033)" paragraph with the drift noted below.
6. Gate recorded above. Limitations: native windows-tests CI result is not available locally and the owner must supply it before the Full Objective Check; resume focuses on the selected Task, so Objective-wide "In progress" lines appear only on the Objective detail and plain board. No Check written, no CLEAR claimed, health check not run.

Files read beyond Context Files: `internal/data/gate_v2.go`, `internal/data/next.go`, `internal/data/evidence_v2.go`, `internal/board/v2/lanes.go`, `internal/board/v2/parallel_advice_test.go`, `internal/board/v2/lanes_test.go`, `internal/data/concurrency_v2_test.go`, `internal/data/concurrency_plan_v2_test.go` (to reuse helpers and avoid duplicating tests), and the O-033 Objective.md. Actual writes match planned_writes.

Owner validation: required and pending; owner Task-check waiver: none recorded.

## Drift Notes

Architecture drift recorded in Design.md: `Next.Concurrency` carries the projection so resume and board share it, and the Goal-wide plain board prints no advice. No material acceptance deviation or REPLAN REQUIRED.


## Advisory Feature Contract — Revised 2026-10-03

Parallel planning defaults off via features.parallel_planning. Code Health remains core and is unchanged. Lane metadata and read/write manifests are suggestions: the owner and all skills may ignore them even when enabled. No lane rule feeds existing start, advance, completion or Check gates, and no manifest deviation alone requires replan. Actual dependencies, acceptance criteria, Guardrails and rules for a real worktree remain authoritative. This revised Objective contract supersedes any earlier mandatory lane wording.
