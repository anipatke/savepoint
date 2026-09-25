---
id: O-024
title: Move finished Objectives below open work
status: planned
depends_on: [O-020]
release: R-006
priority: critical
rank: 1
---

# O-024: Move finished Objectives below open work

## Outcome

The Goal's open Objectives fill the top of the sidebar under their priority
headings. Every finished Objective sits below them under one `DONE` heading, in
ID order. The piped board output uses the same order and headings.

## Why

Most of a mature Goal's Objectives are finished. In R-006, 10 of 13 are done.
O-020 kept finished Objectives in their ranked position, so they are
interleaved with the work the owner still needs to focus on. With the
"Planned / In Progress / Done" row line gone, a finished row is also hard to
tell apart from an open one.

## Success Conditions

- The sidebar shows open Objectives (`planned` or `in_progress`) first, under
  the `CRITICAL` / `HIGH` / `MEDIUM` / `LOW` headings, ordered as O-020 defined
  (priority, rank, then ID). A single `DONE` heading follows, listing every
  `done` Objective in ascending ID order. No `TO DO` heading is added. Empty
  groups, including an empty `DONE`, show no heading.
- The piped `savepoint board` output lists the same rows under the same
  headings in the same order.
- `1`–`4`, `K`, `J`, `shift+↑`, and `shift+↓` do nothing on a done row. On an
  open row, `K`/`J` move only among the open rows of its priority group, and a
  priority change places the row at the bottom of the destination group's open
  rows.
- Reorder writes never rewrite a done Objective's file. Doctor's duplicate-rank
  warning considers open Objectives only, so a done Objective's leftover rank
  never produces a warning.
- Reopening a done Objective returns it to its priority group at its recorded
  rank. A collision there is the existing duplicate-rank warning, healed by
  the next move.
- The cursor, selection, Check and wait badges, `▸`/`●` markers, Next, and
  the router behave as they do after O-020. A router selection that names a
  done Objective shows its `●` in the `DONE` section.
- Tests cover all four priority groups plus `DONE`, a Goal with no done
  Objectives, a Goal with only done Objectives, the no-op keys on done rows,
  `K`/`J` skipping done rows in the same group, no write to done files,
  duplicate-rank detection ignoring done rows, and TUI and non-TTY agreement.
  `git diff --check` and `make build && make test-fast` pass.

## Architectural Considerations

- `internal/data` owns the order. `OrderedObjectiveIDsForGoal` returns open
  Objectives in O-020 order, followed by done Objectives in ID order. The
  board and plain output render that list without sorting.
- The board decides each row's heading from the recorded status (`DONE` for
  `done`) and otherwise from priority. It never infers status from position.
- Lifecycle, dependencies, gates, `ResolveNext`, and Goal completion are
  unchanged.

## Boundaries

**In scope:** the data order function, duplicate-rank facts, sidebar and
plain headings, reorder-key handling of done rows, tests, and the Design §8
wording at reconciliation.

**Out of scope:** collapsing or hiding the `DONE` section, a `TO DO` heading,
priority sub-headings inside `DONE`, changes to Task columns, and new keys.

## Confirmed Design Decisions

Confirmed by the owner on 2026-09-25 in chat after reviewing the proposal.

- **Layout:** no `TO DO` heading. Open rows keep the priority headings, and one
  `DONE` heading follows at the bottom.
- **Done order:** ascending ID only, with no priority sub-headings.
- **Keys:** order keys do nothing on done rows, and `K`/`J` swap only with open
  neighbours.
- **Placement:** this is a follow-up Objective. O-020 closed unchanged. This
  Objective supersedes O-020's rule that a completed Objective keeps its ranked
  position.
- **Planner choice, flagged for owner review:** done Objectives are left out of
  reorder writes and duplicate-rank detection. This replaces the earlier
  proposal to keep them at the end of the write list. Both approaches avoid
  false duplicate warnings. This one never rewrites a finished record.
