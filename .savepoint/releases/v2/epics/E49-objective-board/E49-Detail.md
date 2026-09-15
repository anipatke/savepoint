---
type: epic-design
status: planned
---

# E49: Operate V2 through a clear terminal board

## Purpose

Users see Objectives, Task outcomes, Checks, owner waits, and unified Issues without internal workflow jargon.

V1 delivery epic for V2 product Objective O009. The current V1 lifecycle governs this work; this is the single authoritative backlog record.

## What this epic adds

Objective sidebar/selection; three columns with badges; prominent Next; Issues/Defect filtering; Check/history details; gate-aware actions; V2 watchers and plain output. Task cards, details, and plain output use the required human-facing Task title.

## Components and files

internal/board/model.go; interfaces.go; io.go; watch.go; update.go; view.go; detail.go; epic_panel.go; audit_overlay.go; defect_overlay.go; plain.go; cmd/board.go; internal/styles/styles.go.

Paths not yet present are proposed targets to confirm during this epic's design. Task Context Files must name exact files when the detailed plan is written.

## Architectural delta

Replace release/epic-specific model and duplicate overlays with V2 presentation; I/O stays in explicit commands and gates in data. V2 presentation reads Task `title` as the display label and keeps the detailed objective/Outcome as supporting detail. A missing V2 title is shown as invalid project data from E42, never hidden by falling back to technical objective text; any V1 fallback remains confined to transitional V1 rendering.

Reference: `.savepoint/releases/v2/v2-Design.md`.

## Boundaries

**In scope:** the outcome and additions above.

**Out of scope:** No extra Kanban columns for lifecycle flags, mouse/cloud features, or arbitrary UI redesign beyond Atari-Noir.

## Dependencies

- E47-onboarding-upgrades
- E48-next-resume

## Quality gates

Narrow/short terminals, Unicode width, stable borders/focus, no accidental wrapping, color fallbacks, no-TTY output, edits/reloads/conflicts, and blocked actions verified. Card, detail, and plain-output tests prove the human-facing title is displayed, the detailed objective remains secondary, and missing V2 titles surface as diagnostics rather than fallback labels.

Implementation handoff requires `make build && make test` and named outcome evidence. Epic closeout requires a fresh independent V1 audit; this planning session is not that audit. Apply relevant current Guardrails, keeping STYLE advisory.

## Open decisions

Read visual identity and Bubble Tea guidance during scoped epic design; do not treat this outline as a layout spec.
