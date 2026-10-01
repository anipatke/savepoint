---
id: T-070
title: Prepare what the Code Health screen shows
objective: O-031
status: planned
depends_on: []
owner_validation: {required: false}
planned_by: {role: planner, session: planning-o031-20261002}
complexity_tier: medium
complexity_reason: New read-only projection over existing classification, trend, and repository-relation logic; no new policy.
---

# Prepare what the Code Health screen shows

## Outcome

Code Health can hand the board one plain, read-only description of current health — per-signal labels, reasons, trends, provenance, history, and whether the code has moved on — so the board never touches snapshot or provider types and never decides a label itself.

## User Check

None beyond the Full Objective Check.

## Done When

- `codehealth.LoadDashboard(root)` reads only the saved configuration and snapshots (no subprocess, no tool run) and returns one of three states: not configured, configured with no snapshot (first run), or measured.
- A measured dashboard describes the newest snapshot of either origin and says which origin it is. It has one row per configured instance plus a not-configured row per missing capability, for the five signals in canonical order. Each row carries the Good / Watch / Needs Attention / Unknown label, the explanation, trend words, the comparison basis, the outcome in plain words (failed, timed out, unavailable, absent, partial, unsupported), required or optional, provider and tool version, scope, collection time, and up to the stored evidence references.
- Labels and explanations come from the persisted summary and `Assess`/`Trend`; the dashboard adds no new classification rule. The overall label is `Overall`'s; there is no number. Good appears only when every required result is complete and good; partial and failed optional results stay visible as rows.
- The history list gives the most recent snapshots (bounded) with time, origin, and overall label, newest first.
- `codehealth.DashboardFreshness(ctx, root, git, snapshot)` uses `ObserveRepository` and `Relate` to return a plain state: matches the code, code changed since (with commit count when known), from another branch (diverged/ahead), or unknown, plus whether the working tree differs from the measured inputs. Git failure yields unknown, never an error the screen must handle.
- No wording claims proof of correctness, safety, maintainability, or production readiness (a test asserts against a banned-phrase list over all produced text).
- Tests use temporary projects and fake Git runners and cover: not configured, first run, all five signals measured, partial support, required and optional failures, manual newest over official, reset and comparable trends, same commit, behind, diverged, dirty, and Git unavailable.

## Context Files

`.savepoint/objectives/O-031-code-health-tui/Objective.md`; `internal/codehealth/classification.go`, `internal/codehealth/history.go`, `internal/codehealth/storage.go`, `internal/codehealth/snapshot.go`, `internal/codehealth/model.go`, `internal/codehealth/repository.go`, `internal/codehealth/gate.go` (wording conventions only); new `internal/codehealth/dashboard.go`, `internal/codehealth/dashboard_test.go`, `internal/codehealth/dashboard_freshness.go`, `internal/codehealth/dashboard_freshness_test.go`.

## Design References

Design sections 1 (Code Health readers), 8, and 11; O-031 Architectural Considerations and Confirmed Design Decisions.

## Guardrails

ARCH-02, ARCH-03, ARCH-04, DATA-03, STYLE-07, STYLE-09, TEST-01, TEST-02, TEST-04.

## Implementation Plan

1. Confirm `Assess`, `Overall`, `Trend`, `LoadSnapshots`, `ObserveRepository`, and `Relate` match these paths; return REPLAN REQUIRED if the persisted summary cannot reproduce labels without re-running collection.
2. Define the dashboard types with plain strings and the existing closed vocabularies only; keep `Snapshot` and `CapabilityResult` out of the exported surface.
3. Build `LoadDashboard` from config plus snapshots, rebuilding each row's history the way `Collect` does so trends match.
4. Build `DashboardFreshness` over the repository relation.
5. Keep all user-facing phrases in one table (STYLE-09) and test them against the banned-claims list.

## Boundaries

No change to classification thresholds, `Collect`, the snapshot schema, `Evaluate`, or storage. No board code.

## Technical Verification

Focused `make test-focused TEST=./internal/codehealth/...` while iterating; `make build && make test-fast` for handoff. Later evaluated per `agent-skills/references/check-method.md`.

## Technical Evidence

Pending execution.

## Drift Notes

Adds a dashboard responsibility to `internal/codehealth`; the AGENTS.md Codebase Map entry is updated in this Task (ARCH-04).
