---
id: T-096
title: Explain Code Health and its release limits
objective: O-032
status: planned
depends_on: [{task: T-085, requires: clear}, {task: T-086, requires: clear}, {task: T-087, requires: clear}, {task: T-088, requires: clear}, {task: T-089, requires: clear}, {task: T-090, requires: clear}, {task: T-091, requires: clear}, {task: T-092, requires: clear}, {task: T-093, requires: clear}, {task: T-094, requires: clear}, {task: T-095, requires: clear}]
complexity_tier: medium
complexity_reason: Public guidance must reconcile all completed features and honest diagnostic/release evidence.
owner_validation: {required: true}
planned_by: {role: planner, session: planning-o032-20261002-owner-confirmed}
---

# Explain Code Health and its release limits

## Outcome

Builders have a clear Code Health walkthrough and v2.1 release notes reflecting verified behavior, provider prerequisites and remaining limitations.

## User Check

Follow the textual walkthrough and review the release limitations as a builder would.

## Done When

- Document all five signals and nine supported report formats, supported/unsupported stacks, setup/check/report, H/R history/refresh and agent handoff, official/manual, labels/thresholds/trends/staleness and Check boundaries.
- Explain local-first/no telemetry and project-owned tools honestly, including possible provider network use and inherited environment. State partial/unavailable evidence and why Code Health is no guarantee.
- Describe no automatic pruning, existing internal maintenance boundary and deferred owner cleanup command; document intentional template/history duplication and agreed maintained-code comparisons.
- Include textual TUI walkthrough, packaging/provider rationale, measured performance evidence, health investigation/fix outcomes and separately named follow-ups. Do not claim unresolved scanner failures are vulnerability-free or all indicators are green.
- Reconcile project Design to final implemented reality and confirmed tools without duplicating exact commands outside config or adding future backlog architecture. Owner reviews copy and walkthrough; mandatory Full Objective Check follows all owner-completed Tasks.

## Context Files

`.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-085-find-why-the-vulnerability-scan-fails.md` (dependency evidence); `.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-086-measure-duplication-in-maintained-code.md` (dependency evidence); `.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-087-simplify-diagnostic-repair-rules.md` (dependency evidence); `.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-088-separate-test-stream-processing-responsibilities.md` (dependency evidence); `.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-089-simplify-board-reload-state-restoration.md` (dependency evidence); `.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-090-protect-saved-health-evidence-under-failures.md` (dependency evidence); `.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-091-verify-provider-processes-stop-safely.md` (dependency evidence); `.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-092-preserve-projects-when-adopting-code-health.md` (dependency evidence); `.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-093-exercise-all-health-readers-together.md` (dependency evidence); `.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-094-measure-the-cost-of-health-history.md` (dependency evidence); `.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-095-verify-the-v2-1-binary-and-package-boundary.md` (dependency evidence); `README.md`; `.savepoint/Design.md`; `.savepoint/health/config.json`; `AGENTS.md`; `CHANGELOG.md` (new file).

## Design References

O-032 Confirmed Release-Hardening Design and Owner-Requested Health Analysis; Design sections 1, 6, 7, 9, 12 and 13 as applicable.

## Guardrails

FS-01, FS-04, FS-05, FS-06, ARCH-01, ARCH-02, ARCH-03, CFG-01, CFG-02, TEST-01, TEST-02, TEST-03, TEST-04, TEST-05, TEST-08, TEST-09, STYLE-01, STYLE-02, STYLE-07, TPL-01, TPL-02, TPL-03, TPL-04.

## Implementation Plan

1. Reconcile Task evidence and known follow-up decisions into an accurate feature/release narrative.
2. Add README walkthrough and a v2.1 changelog section, preserving existing content if CHANGELOG exists.
3. Update Design only for implemented changes; retain command ownership boundaries and config separation.
4. Validate command examples against existing tests/help evidence and submit owner-readable copy for review.

## Boundaries

Documentation and implemented Design reconciliation only; no source repair, health setup/config changes, release tagging or publication.

## Technical Verification

`make build && make test-fast` for ordinary handoff; focused `make test-focused TEST=... PKGS=...` only during iteration. Record named per-criterion evidence and limitations. The later independent Check applies `agent-skills/references/check-method.md`; the mandatory Full Objective Check runs its fresh full gate and official Code Health collection. This executor never writes a Check or runs savepoint health setup/report/check.

## Technical Evidence

Pending execution: named per-criterion cases, command/time/toolchain/results, files read/changed, decision deliverable where applicable and limitations.

## Drift Notes

Record implemented responsibility or interface changes for planner reconciliation before the mandatory Full Objective Check. If evidence requires a material unknown repair, return REPLAN REQUIRED rather than inventing another scope.
