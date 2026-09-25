---
id: O-022
title: Require every router and Objective to name a Goal
status: done
depends_on: [O-014]
release: R-006
---

# O-022: Require every router and Objective to name a Goal

## Outcome

In every Savepoint project, the router selects a declared Goal and every live
Objective belongs to one. A project that breaks this rule still opens, but
Next, the board, and doctor name exactly what is missing and how to fix it.
New and migrated projects start valid. There is no project-wide, Goal-less
view of the work.

## Why

On 2026-09-23 the owner directed that Goal (Release-compatible) context is
mandatory. On 2026-09-24 the owner confirmed this applies to every Savepoint
project, not only this repository, and split it out of O-014 because it
changes the shipped product contract that Goals are optional.

What this delivers is one path instead of two. Today every surface carries a
"no Goal" branch: an unscoped board, unassigned Objectives, "a Goal is
optional" guidance, and init/migrate outputs with no Goal. Removing those
branches simplifies the product and unblocks O-020, whose ranking assumes
every Objective sits in exactly one Goal.

Already true before this Objective, and not re-planned here: Next is exactly
the router's selection with no project-wide search (O-014); a malformed or
unknown Objective `release:` already fails the load; an unknown or archived
router Goal already produces a selection diagnostic; the board's Goal
selector already cannot clear the Goal.

## Success Conditions

- A router with a missing, blank, or `none` Goal produces a typed diagnostic.
  Next reads `Choose a Goal`, names the fix, and substitutes no other work.
  Resume, the board, and doctor show the same diagnostic.
- An Objective with no `release:` loads, but the index records it as missing
  a Goal. Doctor lists each such Objective with the exact fix. The board and
  resume flag that such Objectives exist.
- With no valid Goal selected, the board shows no project-wide Tasks or
  Objectives. It shows the diagnostic and how to choose a Goal. With zero
  Goals, it points to doctor.
- `savepoint init` scaffolds R-001, titled after the project, with stub
  sections, and a router that selects it. The Idea skill fills in its Outcome
  with the owner.
- `savepoint migrate` gives every converted Objective a Goal reference. It
  retains the V1 router's live Goal. When that selection is missing or
  unresolvable, it selects the live Goal containing the router's active
  Objective, or the sole live Goal, when that choice is clear. If selected
  active work belongs only to a historical Goal, migrate creates a live
  continuation and moves those active Objectives into it. An unresolved
  release lifecycle decision remains in the preview and blocks Apply.
- AGENTS.md, Design.md, the phase skills, and the scaffold state that a Goal
  is required, replacing "a Goal is optional". Live and scaffold copies stay
  byte-identical. `Choose` joins the owner-reported Next words.

## Architectural Considerations

- Membership stays derived from `Objective.release`; the Goal record gains no
  member list. Storage names (`release:`, `R-###`, `scope.kind: release`) stay
  the compatibility boundary.
- Owner decision (2026-09-25): an existing project missing references loads,
  flags, and guides; it does not fail closed. Missing references are a typed,
  non-fatal index fact. Unknown or malformed references stay fatal.
- `ResolveSelection` owns the router rule, as a new selection diagnostic kind.
  `LoadV2Index` owns the Objective membership rule. Resume, board, and doctor
  render these facts; none re-validates on its own.
- Owner decision (2026-09-25): init creates an R-001 placeholder Goal through
  the existing `{{PROJECT_NAME}}` interpolation. For migration fallback, reuse
  a uniquely identifiable existing live Goal for the selected active work.
  Create a continuation when selected work belongs only to a historical Goal,
  and move that work into it. An unresolved lifecycle decision stays in the
  preview.
- O-020's Boundaries already assume this rule (no unassigned group or
  all-Objectives view).

## Boundaries

**In scope:** router and Objective Goal diagnostics, the `Choose a Goal` Next
line, Goal-scoped board views, init and migrate Goal defaults, doctor repair
guidance, and guidance and Design reconciliation.

**Out of scope:** renaming storage fields; Goal publishing or deployment;
Objective ranking (O-020); changing Next's selection-only rule (O-014);
automatically assigning Goals to existing Objectives; a board action that
edits an Objective's Goal.
