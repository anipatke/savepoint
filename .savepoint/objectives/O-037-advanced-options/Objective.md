---
id: O-037
title: Choose optional features in Advanced Options
status: in_progress
depends_on: [O-032]
release: G-001
priority: medium
rank: 1
---

# O-037: Choose optional features in Advanced Options

## Outcome

Owners can open Advanced Options and save a project-specific Parallel planning preference before optional lane/worktree advice is implemented.

## Why

Optional complexity needs an accessible opt-in boundary first. The settings screen and preference storage do not depend on the lane recommendation engine.

## Success Conditions

- One saved boolean features.parallel_planning defaults off when absent and survives restart and upgrade.
- A documented Settings action opens a keyboard-accessible Advanced Options screen with the single Parallel planning option and a short explanation that suggestions are ignorable.
- Settings writes preserve unrelated configuration/comments, refuse stale conflicts, safely replace files, and report failed writes without displaying unsaved choices as saved.
- Existing board selection, keyboard focus, reload behavior and lifecycle gates remain unchanged. External preference edits are reflected on reload.
- Fresh scaffolds provide the off default; upgrades preserve owner choices and authored assets.
- The preference and screen are delivered and independently verified before O-033 adds lane advice. At this stage the switch stores a preference for that later feature, not a claim that lane suggestions already exist; public wording explains that sequencing.
- Code Health remains core and unchanged. No health option, tool setup/collection, snapshot deletion or verification waiver is added.

## Architectural Considerations

internal/data owns the typed project preference and narrow persistent writer. The board consumes saved configuration, writes only through explicit Bubble Tea commands, and renders without IO. Keep one actual feature key and screen, avoiding a generic registry or configuration editor. Publish the interface downstream O-033 uses; no Task/Lane projection is required here.

## Boundaries

**In scope:** preference storage, safe writes, one Advanced Options screen, config reload/watch integration, fresh/upgrade compatibility, tests, documentation and implemented Design reconciliation.

**Out of scope:** lane records/recommendations/groupings/instructions, Code Health changes, general feature frameworks, automatic worktrees/branches/merges, new lifecycle or completion gates.

## Confirmed Design — 2026-10-03

The owner instructed creating and reordering Objectives/Tasks to build Advanced Options first. This confirms splitting the already reviewed preference/screen contract ahead of parallel planning; it does not authorise production implementation or owner completion. The saved key is features.parallel_planning, default false. Actual lane advice is entirely ignorable even when the option is on. Settings require safe, source-conflict-aware preservation, explicit save feedback and a normal reload. Code Health remains unchanged.

## Readiness and Verification

O-032 is done. Targeted prior reads establish ConfigReader/YAML ownership in internal/data/config.go, board overlay patterns in model/update/view/help, and command-based persistence/reload in io/load. The preference writer and options files are explicitly new targets. Work is sequential because the screen relies on the preference API and shares reload behavior with final adoption verification.

Every Task records per-criterion evidence. Settings replacement/watch work and final integration require fresh make test-full. The mandatory independent Full Objective Check applies check-method.md and uses current native windows-tests CI evidence from repository CI, supplied by the owner. Optional Task Checks require explicit owner waivers when skipped. Owner validation covers save/restart/close, on/off and unchanged Code Health. Only implemented reality is reconciled into Design.

## Planning Handoff

Existing settings Tasks retain their allocated identities and move into this Objective. Remove their dependencies on lane implementation; downstream parallel planning depends on this Objective. A separate adoption/integration Task covers scaffold defaults, upgrades, public documentation and Design reconciliation. Router selects this Objective in design pending owner plan review. No production implementation, Check, waiver or Task completion is written.
