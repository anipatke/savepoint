---
id: O-015
title: Let the owner advance Issues from the board
status: done
depends_on: [O-012]
release: R-006
priority: high
rank: 3
---

# O-015: Let the owner advance Issues from the board

## Outcome

In the board's Issues panel, Space moves the selected Issue one column right
(Open → In Progress → Resolved) and Backspace moves it one column left
(Resolved → In Progress → Open), the same keys the Task columns use.

## Why

The Issues panel shows Open, In Progress, and Resolved columns but is
read-only. Owners must hand-edit Issue frontmatter to start, resolve, or
reopen an Issue.

## Success Conditions

- With the Issues panel focused and no detail open, Space advances the
  selected Issue one status and Backspace retreats it one status. Focus
  follows the Issue into its new column after reload.
- Space on a Resolved Issue, Backspace on an Open Issue, an empty column, or
  no selection performs no write.
- Open → In Progress and In Progress → Open change only `status` and append
  one history entry.
- In Progress → Resolved sets `status: resolved` and writes
  `resolution.disposition: accepted` with actor `{role: owner, session:
  board-owner}`, the action time, and the fixed reason `Resolved by the owner
  from the board.` It never names a proof Check.
- Resolved → In Progress sets `status: in_progress`, removes the
  `resolution` block (and `duplicate_of` / `escalated_to` when present), and
  appends a `reopened` history entry whose note names the removed
  disposition.
- Every transition appends exactly one owner-attributed, timestamped history
  entry and never edits or reorders existing entries.
- Writes preserve the Markdown body and unknown frontmatter, validate the
  result before replacing the file, and refuse a record changed on disk.
  Failures show in the board status line; success reloads through the
  existing board load path.
- The Issues panel footer and help list Space and Backspace.
- Focused data and board tests cover all four transitions, the no-op cases,
  append-only history, preserved content, a stale-file conflict, and focus
  following the Issue; `make build && make test-fast` passes at handoff.
- AGENTS.md, the active skills/references, and their V2 scaffold copies say
  the owner may resolve (as `accepted`) and reopen Issues from the board.

## Confirmed Design Decisions

Owner, 2026-09-25 (planning chat):

- Only Space (forward) and Backspace (back) in the Issues panel. No prompt,
  no typed reason, no disposition choice.
- Board resolution always records `accepted` with the fixed reason above;
  the owner hand-edits the file for anything more specific.
- Backspace may reopen any resolved Issue, including a verified, duplicate,
  or escalated one; the reopened history note keeps what it was closed as.
- The earlier `superseded` disposition and accepted/superseded prompt are
  dropped from this Objective.
- History kinds reuse the existing vocabulary: `owner_decision` for start,
  resolve, and back-to-open; `reopened` for Resolved → In Progress.

## Architectural Considerations

- `internal/data` owns the Issue transition, the resolution/history patch, and
  the validate-before-replace write through the existing `writeV2Record`
  boundary. The board asks for a transition by Issue ID; it does not build
  Issue YAML.
- `internal/board/v2` schedules the write as a `tea.Cmd` (ARCH-02), re-reads
  the project before writing, and reuses `actionMsg` status/reload handling.
- Owner resolution is not technical clearance: Check, Task, Objective, and
  Goal gates are unchanged.

## Boundaries

**In scope:**

- Space/Backspace Issue transitions in the V2 Issues panel, the data-layer
  Issue writer, footer/help text, tests, and guidance/template wording.

**Out of scope:**

- New dispositions, reason entry, or any prompt.
- Transitions from the Issue detail overlay, or skipping a column.
- Creating Issues, editing Issue text or links, duplicate/escalation flows.
- Changing Issue status/type/history vocabularies, Task behavior, router
  selection, or gates.
- Rewriting I-001 or any other historical record.
