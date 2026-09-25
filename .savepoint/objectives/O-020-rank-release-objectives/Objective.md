---
id: O-020
title: Prioritize and rank Objectives within a Release
status: done
depends_on: [O-022]
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
- With a Goal selected and the sidebar focused, keys `1`–`4` set the focused Objective's priority directly (Critical, High, Medium, Low). The row moves to the bottom of its new group immediately, preserving the other rows' relative order. `K`/`J` (and `shift+up`/`shift+down` where the terminal reports them) move it up or down within its current group. No manual rank-number entry is needed.
- The cursor and selection follow the same Objective after either action, so repeated shortcuts remain predictable. Changes persist across restart and normal watcher reload; an external valid edit appears on reload.
- Board writes preserve all other authored Objective content. Each file replacement is atomic. A write that is interrupted partway through a group renumber leaves a deterministic order (ties fall back to ID), doctor names the duplicate rank, and the next move in that group heals it. A conflicting external edit to any record in the write set produces a diagnostic before anything is written instead of silently overwriting it.
- A completed Objective stays in its owner-set position; the next-candidate calculation skips it without reshuffling the list. A blocked Objective stays in position with its blocker visible. The list never infers “In Progress” from position.
- Moving an Objective to another Release removes the old Release position and gives it a defined initial position in the destination group. Every Objective belongs to a Release (O-022), so every Objective has a position in exactly one Release.
- Interactive and non-TTY Release views use the same order and labels. Ranking is always scoped to the selected Release; there is no all-Objectives view or unassigned group (O-022).
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

**Out of scope:** Renumbering Objective IDs; using priority or rank as a new lifecycle state or completion rule; changing Task card columns; rewriting archived V1 or immutable Check records; making Release membership optional again.

## Confirmed Design Decisions

Confirmed by the owner on 2026-09-25 during planning (O-013, O-018, and O-022 are done).

- **Dependencies:** the owner dropped O-013 from `depends_on` on 2026-09-25. O-013 closed by owner exception (C-911 NEEDS WORK), which the dependency gate does not count as clearance. The Goal vocabulary O-020 needs has already shipped, so O-022 is the only dependency.

- **Storage:** two optional Objective frontmatter fields, `priority: critical|high|medium|low` (missing reads as `medium`) and `rank: <positive integer>` (missing sorts after ranked rows in its group). An invalid value is a strict decode error, like every other malformed V2 field. Display order is priority group, then rank ascending, then ID.
- **Writes are self-healing, not transactional:** one data writer sets `priority` and `rank` on a list of Objectives in their intended group order (1..n), writing only records whose values change, after checking every record in the set is fresh. A move swaps two neighbours and renumbers the group; a priority change appends the row to the destination group and renumbers that group. The old group may keep a gap, which is harmless. Duplicate ranks in one Goal and group are a doctor warning, not a load error.
- **Keys (sidebar focused):** `1`–`4` set priority; `K`/`J` and `shift+up`/`shift+down` move within the group.
- **Display:** the sidebar shows `CRITICAL` / `HIGH` / `MEDIUM` / `LOW` headings only for non-empty groups, with rows beneath; the separate Planned / In Progress / Done line is removed. Non-TTY output gains the same grouped Objective list.
- **Next is unchanged:** ranking is display and planning order only. Next remains exactly the router's selection, closing an Objective still clears it, and the owner presses `p` on the row to work on. Blocked rows keep their wait badges in place.
- **Cross-Goal moves:** there is no board action to move an Objective between Goals. When `release:` is hand-edited, the record keeps its `priority` and `rank` and is placed in the destination group by them; a rank collision there is the same doctor warning, healed by the next move.
