---
id: O-020
title: Prioritize and rank Objectives within a Release
status: planned
depends_on: [O-013]
release: R-006
---

# O-020: Prioritize and rank Objectives within a Release

## Outcome

The owner can order a Release's Objectives by Critical, High, Medium, or Low priority and by position within each group, without renumbering their stable identities. The V2 board updates immediately and removes redundant Objective lifecycle text from sidebar rows while retaining exact status and Check evidence in detail and Next views.

## Why

Objective IDs identify records; their numbers should not imply execution sequence. The current Release sidebar sorts by ID, so an owner cannot express that a later-created Objective should be worked on first. Each sidebar row also prints “Planned,” “In Progress,” or “Done” beneath its heading even though task progress, Check badges, dependencies, and the Next area provide more useful action context. The owner wants fast priority changes and a stable stack within each priority group, without another status claim.

## Success Conditions

- Each Release member has a priority of Critical, High, Medium, or Low and an order within that priority group. Existing records without these fields load as Medium and retain stable ID order until the owner explicitly reorders them.
- A Release's Objectives display by priority group (Critical, High, Medium, Low), then by owner-set position within each group. ID is only a stable fallback for unranked records; it is never presented as the intended work sequence. Invalid or duplicate explicit positions in the same Release and group produce a named diagnostic.
- Priority and position live on Objective records alongside the existing `release:` reference. Release membership continues to be derived from Objective records; the Release file gains no duplicate member list. Priority and position are planning metadata, not identity, lifecycle state, or completion gates.
- With a Release selected, an owner shortcut cycles the focused Objective's priority. The row moves to its new group immediately, preserving the other rows' relative order. Move Up and Move Down change its position within its current group. No manual rank-number entry is needed.
- The cursor and selection follow the same Objective after either action, so repeated shortcuts remain predictable. Changes persist across restart and normal watcher reload; an external valid edit appears on reload.
- Board writes preserve all other authored Objective content. Reordering avoids transient duplicate positions and has recoverable behavior on an interrupted write. A conflicting external edit produces a diagnostic instead of silently overwriting it.
- A completed Objective stays in its owner-set position; the next-candidate calculation skips it without reshuffling the list. A blocked Objective stays in position with its blocker visible. The list never infers “In Progress” from position.
- Moving an Objective to another Release removes the old Release position and gives it a defined initial position in the destination group. Unassigned Objectives have no Release position and are not silently assigned by prioritization.
- Interactive and non-TTY Release views use the same order and labels. The all-Objectives view makes Release grouping and unassigned placement explicit rather than applying one Release's rank globally.
- The Objective sidebar omits the separate “Planned / In Progress / Done” line beneath the heading. Priority, Check, wait, selection, and focus signals remain legible, including without colour and at narrow widths. Objective detail and the canonical Next projection continue to report exact recorded lifecycle and evidence.
- The first unfinished Objective in displayed order is a priority candidate, not automatically “In Progress.” The canonical Next resolver still respects dependencies, Check and owner waits, and explicit router selection. If the top candidate is blocked, the UI names the wait rather than implying it has started.
- Tests cover legacy defaults, all four priority groups, within-group ordering, ties, cycling, cursor stability, completed and blocked rows, cross-Release moves, interrupted/conflicting writes, narrow terminals, non-TTY output, and unchanged Check/owner decisions. Canonical and scaffold guidance remain aligned; `git diff --check` and `make build && make test` pass before the mandatory Full Objective Check.

## Architectural Considerations

- `internal/data` owns priority and position parsing, validation, and the derived ordered member projection. `internal/board/v2` renders that projection and records owner actions through explicit write commands. No view parses these fields independently.
- Priority groups sort first; position applies only within a group and Release. A priority change inserts the Objective into the destination group while preserving other rows' relative order. Existing unranked records stay readable without a data migration.
- Priority and order change presentation and owner planning order. They do not change Objective status, dependencies, technical clearance, Release completion, or `ResolveNext` precedence without an explicit decision during Task planning.
- The existing router selects current work. A reorder must not silently overwrite that selection; Task planning must define the owner-facing handoff when the visible top candidate and router selection differ.
- O-013 changes Release wording to Goal in the UI. This Objective follows that vocabulary when O-013 is complete while preserving persisted `release:` references.
- O-018 may later hyphenate record identities. Priority and position are independent of ID spelling; implementation should be sequenced to avoid conflicting edits to shared data and board files.

## Boundaries

**In scope:** Release-scoped Objective priority and position, owner keyboard actions, sidebar and non-TTY ordering, Objective status-line review, persistence/reload behavior, guidance, and focused verification.

**Out of scope:** Renumbering Objective IDs; using priority or rank as a new lifecycle state or completion rule; changing Task card columns; rewriting archived V1 or immutable Check records; automatically moving unassigned Objectives into a Release.

## Planning Handoff

This is a planned Objective, not the current router selection. Detail its Tasks when it becomes the next ready Objective; settle the board shortcut choices, multi-record write/recovery approach, and router handoff against the then-current O-013/O-018 implementation before build.
