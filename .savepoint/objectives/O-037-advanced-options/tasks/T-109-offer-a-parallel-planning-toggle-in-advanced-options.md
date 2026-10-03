---
id: T-109
title: "Offer a parallel planning toggle in Advanced Options"
objective: O-037
status: planned
depends_on: [{task: T-108, requires: clear}]
owner_validation: {required: true}
planned_by: {role: planner, session: codex-o033-advisory-replan-2026-10-03}
complexity_tier: medium
complexity_reason: "A settings overlay must save explicitly and retain reload/focus behavior."
planned_reads:
  - "internal/data/feature_preferences.go"
  - "internal/data/config.go"
  - "internal/board/v2/model.go"
  - "internal/board/v2/update.go"
  - "internal/board/v2/view.go"
  - "internal/board/v2/help.go"
  - "internal/board/v2/footer_test.go"
  - "internal/board/v2/io.go"
  - "internal/board/v2/load.go"
  - "internal/board/v2/load_test.go"
  - "internal/board/v2/watch.go"
  - ".savepoint/visual-identity.md"
planned_writes:
  - "internal/board/v2/options.go"
  - "internal/board/v2/options_test.go"
  - "internal/board/v2/model.go"
  - "internal/board/v2/update.go"
  - "internal/board/v2/view.go"
  - "internal/board/v2/help.go"
  - "internal/board/v2/footer_test.go"
  - "internal/board/v2/io.go"
  - "internal/board/v2/load.go"
  - "internal/board/v2/load_test.go"
  - "internal/board/v2/watch.go"
---

# Offer a parallel planning toggle in Advanced Options

## Outcome

Owners can switch optional parallel advice on/off from Advanced Options with saved-state feedback and no effect on core Code Health.

## User Check

Open Advanced Options, toggle Parallel planning, restart the board, then toggle it off; confirm saved plans remain and Code Health still appears as before.

## Done When

- A documented Settings action opens a keyboard-accessible Advanced Options screen containing exactly the real Parallel planning option and a short explanation that suggestions can be ignored.
- Board load/reload reads the saved preference; toggles persist through explicit Bubble Tea commands. Show saved state only after success and actionable messages on stale/conflicting or failed writes.
- Settings survive board restart and refresh after external edits to config.yml. Esc/close returns focus to the original surface; overlays and terminal sizing follow existing board patterns.
- No Code Health toggle or behavior changes, provider execution, destructive cleanup, shell calls or rendering IO. Turning this option off preserves all lane metadata and evidence.
- Record per-criterion evidence and a fresh full gate; optional Task Check waiver is recorded only on explicit owner instruction. Mandatory independent Full Objective Check remains required.

## Context Files

- `internal/data/feature_preferences.go`
- `internal/data/config.go`
- `internal/board/v2/model.go`
- `internal/board/v2/update.go`
- `internal/board/v2/view.go`
- `internal/board/v2/help.go`
- `internal/board/v2/footer_test.go`
- `internal/board/v2/io.go`
- `internal/board/v2/load.go`
- `internal/board/v2/load_test.go`
- `internal/board/v2/watch.go`
- `.savepoint/visual-identity.md`

internal/data/feature_preferences.go is intentionally created by the preference prerequisite. Other existing Context Files must exist; a material interface mismatch returns REPLAN REQUIRED.

## Design References

O-037 Confirmed Design — 2026-10-03: preference storage and Advanced Options delivery before lane implementation. Project Design sections 1, 8 and 9 describe existing ownership, rendering and concurrency boundaries.

## Guardrails

FS-01, FS-04, FS-05, FS-06, DATA-01, DATA-02, DATA-03, ARCH-02, ARCH-03, CFG-01, CFG-02, CFG-03, TEST-01, TEST-02, TEST-03, TEST-04, TEST-08, TEST-09, STYLE-01..10.

## Implementation Plan

1. Confirm the prerequisite preference interface; no lane projection dependency; add new options.go view/model support and options_test.go with existing overlay conventions.
2. Choose one unused discoverable Settings key and document it in help/footer. Dispatch toggle saves as explicit commands followed by reload.
3. Watch preference-file changes through the existing refresh boundary and preserve board selection/focus while reloading.
4. Verify default/off/on, save/restart, external edit, stale or unwritable save, keyboard close/focus and unchanged Code Health surfaces.

## Boundaries

Implement the scoped preference/screen only. Anticipated manifests and lane choice remain advice, not execution allowlists or gates. Retain ordinary extra-read logging and return REPLAN REQUIRED only for an actual material plan/acceptance gap. Code Health, completion authority, existing dependency gates and actual-worktree rules stay unchanged. No agent settings toggling for the owner's real project during execution; use temporary fixtures. No Check or owner waiver is invented.

## Technical Verification

Fresh make test-full for configuration replacement/path and reload-watch behavior. Focused tests are iteration only. Native windows-tests CI evidence is produced by repository CI and supplied by the owner before the independent Full Objective Check under agent-skills/references/check-method.md. Owner validation is required for the settings screen.

## Technical Evidence

Pending execution: record actual reads/writes, logged extra reads, per-criterion cases/results, command/time/toolchain, full gate and limitations. No self-clearance.

## Drift Notes

Record material implementation/acceptance deltas for planner reconciliation. Ignoring a suggested lane or differing from an anticipated manifest alone requires no replan.


Prerequisite-created setting interface: `internal/data/feature_preferences.go` is intentionally created by T-108; its absence before that dependency completes is expected.
