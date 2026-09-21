---
id: O013
title: Rename Release context to Goals
status: planned
depends_on: [O012]
release: R006
---

# O013: Rename Release context to Goals

## Outcome

The optional delivery context that groups Objectives is presented as a Goal (or
Goals) throughout the V2 board and active project guidance, while existing V2
projects continue to load and behave exactly as before.

## Why

“Release” describes a delivery boundary, but the product experience uses this
concept to give related Objectives a navigable outcome context. “Goals” is the
clearer term for that user-facing grouping and should be consistent across the
board, plain output, workflow guidance, and repository documentation.

## Success Conditions

- V2 user-facing UI consistently calls the context Goal/Goals in selectors,
  headers, detail views, help, status messages, and non-TTY output.
- `g` is the canonical V2 shortcut for the Goal selector; any retained `r`
  alias is compatibility-only and is not the displayed vocabulary.
- Existing `R###` records, `.savepoint/releases/` paths, Objective `release:`
  references, and router `release:` selections remain readable without a data
  migration.
- Objective membership, filtering, completion gates, cutover composition, and
  historical migration behavior are unchanged.
- Active architecture, repository guidance, V2 skills, and shipped V2
  templates describe Goals while documenting the compatibility boundary where
  internal `Release` names remain.
- Archived V1 material, immutable Checks, Issues, and historical Release
  evidence remain untouched.
- Focused UI/data/documentation checks, `git diff --check`, `make build`, and
  `make test` pass before handoff.

## Architectural Considerations

- This is a terminology and presentation change, not a persisted-schema
  migration. `internal/data` remains the owner of the existing `R###`,
  `release:`, membership, and completion contracts.
- `internal/board/v2` remains a read/write presentation surface over the same
  identity-keyed index and resolvers; it must not create a second Goal
  membership map or duplicate completion policy.
- The V2 board's canonical public shortcut becomes `g`; preserving `r` as an
  alias is allowed only to avoid breaking existing muscle memory and must not
  reintroduce Release wording in help or documentation.
- Canonical active skills and their `templates/project-v2/` copies must remain
  byte-aligned after documentation changes.
- The current O012 work remains the predecessor Objective. This Objective is
  planned under R006 but is not activated by changing the current router.

## Boundaries

**In scope:**

- V2 Goal terminology and selector shortcut presentation.
- Compatibility regression coverage for existing Release-backed V2 projects.
- Active architecture, workflow, README, repository-guide, and V2 template
  documentation.

**Out of scope:**

- Renaming persisted `R###` identities to `G###`.
- Moving `.savepoint/releases/` to `.savepoint/goals/` or changing `release:`
  frontmatter and router keys.
- Rewriting migration archives, historical V1 guidance, immutable Checks, or
  Issues.
- Changes to Objective/Goal completion policy, Check semantics, or Release
  membership behavior.
