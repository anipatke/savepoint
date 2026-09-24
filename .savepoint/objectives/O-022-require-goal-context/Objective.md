---
id: O-022
title: Require every router and Objective to name a Goal
status: planned
depends_on: [O-014]
release: R-006
---

# O-022: Require every router and Objective to name a Goal

## Outcome

In every Savepoint project, the router always selects a declared Goal and
every live Objective belongs to one. Next and board views are always scoped to
that Goal; there is no unscoped, project-wide fallback.

## Why

On 2026-09-23 the owner directed that Goal (Release-compatible) context is
mandatory: a missing, blank, `none`, or unknown router Goal, or an Objective
without a valid Goal reference, is invalid. On 2026-09-24 the owner confirmed
this applies to every Savepoint project, not only this repository, and split
it out of O-014 because it changes the shipped product contract that Goals
are optional.

## Success Conditions

- The router always has a valid `release: R-###` naming a declared Goal.
  Missing, blank, `none`, and unknown values produce a clear, actionable
  diagnostic in load, board, resume, and doctor; they never silently show all
  Objectives.
- Every live Objective, including planned and done ones, has a valid
  `release: R-###` reference to a declared Goal. Missing, blank, `none`, and
  unknown references produce a clear diagnostic.
- With no Objective selected, Next searches only the selected Goal. The
  unscoped project-wide ladder is removed.
- `savepoint init` scaffolds a default Goal and a router that selects it.
  `savepoint migrate` produces a valid Goal reference for every converted
  Objective and a router selection. Existing V2 projects missing references
  get a named repair path, not a silent default.
- AGENTS.md, Design.md, Idea-facing guidance, the phase skills, and the
  scaffold (live and scaffold byte-identical) state that a Goal is required,
  replacing "a Goal is optional".
- The board's Goal selector cannot clear the Goal to an all-Objectives view.

## Architectural Considerations

- Membership stays derived from `Objective.release`; the Goal record gains no
  member list. Storage names (`release:`, `R-###`, `scope.kind: release`) stay
  the compatibility boundary.
- `ResolveSelection` and `ReadStateV2`/`WriteRouterStateV2` own the router
  rule; the Objective decoder or index owns the membership rule. No surface
  re-validates on its own.
- O-020 currently describes "unassigned Objectives" and an all-Objectives
  view; its Boundaries need reconciling with this rule before it is planned.
- Settle during Task planning which migration and init changes are needed
  and whether a legacy V2 project without references fails closed or loads
  with doctor repair guidance.

## Boundaries

**In scope:** router and Objective Goal validation, removing the unscoped
Next fallback, init/migrate defaults, board Goal-selector changes, repair
diagnostics, and guidance and Design reconciliation.

**Out of scope:** renaming storage fields, Goal publishing or deployment,
Objective ranking (O-020), and O-014's Next-area and router-selection work.
