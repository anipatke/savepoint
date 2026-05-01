---
type: epic-design
status: planned
---

# Epic E01: Archive and Reset

## Purpose

Archive all v1 epics, reset the project to a clean slate, and establish the v1 release structure for the simplified savepoint. This epic is pure filesystem and metadata housekeeping — no code changes.

## What this epic adds

- `_archived/` directory containing all 11 original v1 epics (E01–E11).
- A rewritten `PRD.md` scoped to the simplified savepoint: board-only CLI, phase model, no audit pipeline.
- Five new epic directories under `releases/v1/epics/` with Design.md stubs.
- A rewritten `router.md` pointing at the first new epic.

## Definition of Done

- All 11 original epic directories live under `releases/v1/epics/_archived/`.
- `releases/v1/PRD.md` describes the simplified scope.
- Five new epic directories exist: E01-archive-and-reset, E02-domain-phase-model, E03-cli-simplify, E04-board-phase-integration, E05-project-cleanup.
- `router.md` uses the 3-state model (`planning`, `building`, `reviewing`) and points at E01.
- No references to old epics remain in active files.

## Components and files

| Path | Purpose |
|------|---------|
| `releases/v1/epics/_archived/` | Archive of all original v1 epics |
| `releases/v1/PRD.md` | Rewritten release scope for simplified savepoint |
| `releases/v1/epics/E01-archive-and-reset/Design.md` | This file |
| `releases/v1/epics/E01-archive-and-reset/tasks/` | Archive + reset tasks |
| `releases/v1/epics/E02-domain-phase-model/Design.md` | Phase model epic stub |
| `releases/v1/epics/E03-cli-simplify/Design.md` | CLI simplification epic stub |
| `releases/v1/epics/E04-board-phase-integration/Design.md` | Board phase epic stub |
| `releases/v1/epics/E05-project-cleanup/Design.md` | Cleanup epic stub |
| `router.md` | 3-state router, points at E01 |
