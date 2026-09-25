---
id: O-015
title: Let the owner advance Issues from the board
status: planned
depends_on: [O-012]
release: R-006
priority: high
rank: 3
---

# O-015: Let the owner advance Issues from the board

## Outcome

The owner can advance a focused Issue through its lifecycle directly from the
Issues panel: Open to In Progress when work starts, then In Progress to
Resolved through an explicit owner decision. Resolution distinguishes accepted
risk from a stale Issue superseded by behavior that already exists, and keeps
both distinct from a checker-verified repair.

## Why

The Issues panel currently presents Open, In Progress, and Resolved as a board
but is read-only. Owners must hand-edit frontmatter to start work or exercise
the already-supported `accepted` resolution disposition. Migrated records such
as I-001 can also describe behavior that the current implementation already
supersedes, but the resolution model has no honest disposition for that case:
`accepted` would falsely call it risk acceptance, `verified` requires repair
proof from a Check, and `duplicate` requires another canonical Issue. A direct
board action and an explicit superseded disposition make the visible lifecycle
usable without falsifying why an Issue closed.

## Success Conditions

- With the Issues panel focused, the displayed forward action advances the
  selected Open Issue to In Progress and keeps focus on that Issue in its new
  column after reload.
- Advancing a selected In Progress Issue to Resolved is an explicit owner
  action that requires choosing the applicable non-verification outcome:
  accepted risk or superseded/stale behavior. It records the matching
  disposition, an owner actor, the action time, and a non-empty reason; it
  never records a proof Check or presents either outcome as verified repair.
- `resolution.disposition: superseded` is a first-class, strictly validated
  Issue resolution for a stale report whose claimed behavior no longer matches
  the supported product. It requires an owner actor and a non-empty reason,
  forbids a proof Check and `duplicate_of`, and renders distinctly from
  `accepted`, `verified`, and `duplicate`.
- I-001 remains the motivating stale-record example and records the owner's
  accepted closure under the schema available before this Objective. O-015
  adds the honest superseded path for future stale reports without rewriting
  I-001's recorded owner decision.
- Each successful transition appends an attributed, timestamped Issue history
  entry without editing or reordering existing history.
- A selected Resolved Issue cannot advance further, and an empty column or
  missing selection performs no write.
- The action is available from the Issue list and is labelled in the panel
  footer/help so its owner authority and effect are understandable without
  colour. Issue detail remains a reading surface.
- Writes preserve authored Markdown and unknown frontmatter, validate the
  resulting Issue before replacement, report conflicts or invalid records in
  the board status area, and reload through the existing board load path.
- Open to In Progress does not create a resolution. In Progress to Resolved
  creates only the explicitly selected owner disposition (`accepted` or
  `superseded`); checker-owned `verified` closure and duplicate disposition
  behavior remain unchanged.
- Focused data and board tests cover both transitions, append-only history,
  attribution, preserved content, terminal states, write failures, and focus
  restoration; `git diff --check`, `make build`, and `make test` pass before
  handoff.
- Active workflow guidance and the V2 scaffold are reconciled so owner-accepted
  and owner-superseded closure from the board are explicit authority paths,
  while executors still cannot close Issues and only a Check can claim verified
  repair.

## Architectural Considerations

- `internal/data` remains the sole owner of Issue lifecycle validation,
  accepted-resolution obligations, append-only history, and safe record
  writes. The board requests a typed transition rather than patching Issue
  YAML itself.
- `internal/board/v2` follows its existing command-message pattern: key
  handling schedules filesystem work, and the update/render path remains free
  of direct IO.
- The action should reuse the board's established forward key and asynchronous
  status/reload behavior where practical, while keeping Issue and Task
  lifecycle vocabularies separate.
- Owner acceptance and supersession are closure dispositions, not technical
  clearance. Existing Check freshness, Task/Objective/Goal gates, and Issue
  linkage remain unchanged.
- The current O-012 work remains the predecessor Objective. This Objective is
  planned under R-006 but is not activated by changing the current router.

## Boundaries

**In scope:**

- A forward owner action in the V2 Issues list for Open to In Progress and In
  Progress to Resolved by accepted-risk or superseded disposition.
- The canonical Issue model, validation, write/read round trip, rendering, and
  documentation for the superseded disposition.
- Canonical transition/write behavior, owner attribution, resolution metadata,
  append-only history, conflict handling, reload/focus behavior, footer/help
  text, tests, and active V2 workflow/template reconciliation.
- Regression coverage using I-001's stale-report scenario, without rewriting
  its already-recorded owner-accepted resolution.

**Out of scope:**

- Calling an owner-accepted or superseded Issue `verified`, fabricating a
  Check, or weakening the proof requirements for
  `resolution.disposition: verified`.
- Resolving duplicates, reopening or retreating Issues, editing Issue text or
  links, creating Issues, or adding free-form resolution editing to the board.
- Changing Issue status/type/severity vocabularies, Task lifecycle behavior,
  router selection, Objective/Goal gates, or the existing Issue detail layout.
- Rewriting historical Checks, Issues, migration archives, or V1 guidance.
