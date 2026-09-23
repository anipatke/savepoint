---
id: O014
title: Give the Next area an Objective word and let the router target Issues
status: planned
depends_on: [O012]
release: R006
---

# O014: Give the Next area an Objective word and let the router target Issues

## Outcome

The board's Next area states a lifecycle word for whichever record `data.Next`
is pointing at — a Task or an Objective alike — instead of only a Task, and
the router's context selection can name a specific Issue the same way it
already names a Release, Objective, or Task.

## Why

Design.md's own Next-area contract (Section 8) says the panel shows "its own
Task or Objective — its lifecycle word ... its identity, and its title." The
actual renderer (`internal/board/v2/next_panel.go`) only builds a word for the
Task branch; an Objective at the ready-for-Check rung prints just its ID and
title, silently breaking that documented parity. This was noticed live while
handing O012 off to its mandatory Full Objective Check, where the Next area
read "O012 — Make Task cards easier to scan" with nothing telling the owner
whether it meant "still building" or "ready for Check."

Separately, the router's selection type (`RouterSelectionV2`) has fields for
Release, Objective, and Task, but nothing for an Issue — there is no way to
point the board or `savepoint resume` at one specific Issue the way the other
three record kinds can already be targeted.

Third, the Next area ignores the router's selected Objective when that
Objective has no Tasks yet. On 2026-09-23, after O016 closed, the router
selected O018 in R006, which had no Tasks. Next still showed "Planned T006 —
Preserve existing project records" from O013. The selected Objective only
wins in two cases: the router also names an unfinished `task`, or the
Objective has Tasks and is waiting on its integration Check. Otherwise the
ready search picks the lowest-numbered startable Task in the whole project
(or Release). It considers a Task-less Objective for planning only after
every ready Task. Task IDs are global, so planning O018's Tasks will not fix
this: their numbers will still sort after T006. The owner added this
behavior to O014's scope on 2026-09-23 instead of opening a separate Issue.

## Success Conditions

- `nextLines`/`renderNext` show one lifecycle word for `next.Objective` the
  same way they already do for `next.Task`: `Planned` (not started), `Build`
  (in progress, an owned Task still needs work), `Check` (every owned Task
  done, ready for the mandatory Full Objective Check), or `Done` — derived
  from the Objective's own status and completion resolution, never borrowing
  a Task-only stage word (`Test`, `Audit`).
- The non-TTY plain rendering states the same word, so piping the board and
  reading it still agree, matching the existing Task-word parity guarantee.
- The router's `## Current state` block can name a single Issue (`issue:
  I###`), decoded and validated with the same strictness and the same "none"
  sentinel convention `release`/`objective`/`task` already use; an existing
  router with no `issue:` key keeps parsing exactly as before (backward
  compatible).
- Selecting an Issue through the router is readable by the board/resume
  without inventing a second selection mechanism alongside
  `data.ResolveSelection`.
- When the router selects an Objective and names no Task, Next reflects
  that Objective ahead of ready work elsewhere in the project or Release. If
  the Objective has no Tasks, Next asks for it to be planned. If it has a
  startable Task, Next offers that Task. The only exceptions are the rungs
  that already outrank Objective selection, such as a pending migration or
  a selected unfinished Task. Board, non-TTY, and `savepoint resume` output
  agree because they share the one `data.Next` projection.
- When the router selects no Objective, Next falls back to the existing
  project-wide (or Release-scoped) search in ID order, and the result is
  unchanged.
- Apart from the selected-Objective rule above, existing
  Release/Objective/Task selection behavior, the `ResolveNext` ladder's
  precedence, and non-TTY parity are unchanged for any project that never
  selects an Issue.
- Active documentation (Design.md Section 8's Next-area contract) is
  reconciled to state the implemented word vocabulary precisely.
- Focused board/data tests, `git diff --check`, `make build`, and `make test`
  pass before handoff.

## Architectural Considerations

- `internal/data` remains the sole owner of Objective status/completion
  resolution and Issue lifecycle; this Objective adds presentation (the
  Next-area word) and a router selection field, not a new gate or a second
  status vocabulary.
- `internal/board/v2/badges.go` and `next_panel.go` remain the single
  translation from typed `data` values to glyph/label/accent/word; no second
  word-mapping for Objectives elsewhere in the package.
- Whether a selected Issue can itself become the resolved `data.Next` (a new
  rung) or stays a pure view-context the way Release selection narrows scope
  without being "the" next action is an open interface question — settle it
  during readiness/Task planning for this Objective, not here.
- The selected-Objective rule is a `ResolveNext` precedence change owned by
  `internal/data/next.go`. The board and resume only render its result. The
  exact rung placement is an owner product decision to confirm before Task
  detailing: whether a selected Objective's integration Check, its own ready
  Tasks, and its planning request each outrank project-wide ready work, and
  how a selected Objective with unmet dependencies or only blocked Tasks
  reports. Keep the rule consistent with O020, which requires the canonical
  Next resolver to keep respecting explicit router selection.
- Keep `router.md`'s YAML shape backward compatible: an existing V2 router
  with no `issue:` key must keep decoding exactly as it does today.

## Boundaries

**In scope:**

- The Objective lifecycle word in the Next area (`internal/board/v2/next_panel.go`,
  and `internal/data/next.go` if the word needs a typed source rather than
  being derived ad hoc in the board).
- Router selection support for a single Issue (`internal/data/router_v2.go`,
  `internal/data/write.go`'s `RouterSelectionV2`): decode, validate, write,
  and the "none" sentinel, plus whatever minimal board surface is needed to
  set it.
- Next honors the router-selected Objective ahead of unrelated ready work
  (`internal/data/next.go` ladder, plus the matching resume/board phrasing
  if a new rung or Next kind is introduced).
- Design.md Section 8 reconciliation for the exact word vocabulary and the
  selected-Objective precedence.

**Out of scope:**

- Changing Issue lifecycle states, severity, or type vocabulary.
- Changing the existing Task lifecycle words (`Build`/`Test`/`Check`/`Planned`/`Done`).
- Any Objective, Task, Release, or Issue gate or completion-policy change
  (the Next precedence change above decides what to show, not what may
  start or complete).
- Renaming Release to Goals (O013's scope) or any other unrelated Next-area
  or router surface not named above.
