---
id: O-014
title: Give the Next area an Objective word and let the router target Issues
status: planned
depends_on: [O-012]
release: R-006
---

# O-014: Give the Next area an Objective word and let the router target Issues

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
handing O-012 off to its mandatory Full Objective Check, where the Next area
read "O-012 — Make Task cards easier to scan" with nothing telling the owner
whether it meant "still building" or "ready for Check."

Separately, the router's selection type (`RouterSelectionV2`) has fields for
Release, Objective, and Task, but nothing for an Issue — there is no way to
point the board or `savepoint resume` at one specific Issue the way the other
three record kinds can already be targeted.

Third, the Next area ignores the router's selected Objective when that
Objective has no Tasks yet. On 2026-09-23, after O-016 closed, the router
selected O-018 in R-006, which had no Tasks. Next still showed "Planned T-006 —
Preserve existing project records" from O-013. The selected Objective only
wins in two cases: the router also names an unfinished `task`, or the
Objective has Tasks and is waiting on its integration Check. Otherwise the
ready search picks the lowest-numbered startable Task in the whole project
(or Release). It considers a Task-less Objective for planning only after
every ready Task. Task IDs are global, so planning O-018's Tasks will not fix
this: their numbers will still sort after T-006. The owner added this
behavior to O-014's scope on 2026-09-23 instead of opening a separate Issue.

On 2026-09-23, the owner directed that Release context is mandatory: the
router must always select a declared Release, and every live Objective must
reference one. A missing, blank, `none`, or unknown router Release is invalid;
an Objective without a valid Release reference is invalid. The previously
unassigned O-016, O-017, and O-019 records were assigned to R-006.

Fourth, the router and the board state two different "next" actions. The
board's Next is computed by `data.ResolveNext` from the records. The router's
`next_action` is hand-written prose that no V2 code reads. It is decoded
into `RouterStateV2.NextAction` but only the unreachable V1 board renders it.
Agents still act on it, because AGENTS.md tells them to start from the
router. On 2026-09-23, after the owner closed T-020 and O-018 (C-910), the
router still selected `objective: O-018, task: T-020` with a `next_action`
asking for O-018's Full Objective Check. The board's Next showed T-006
instead: the selected Task was done, so the Release ladder silently fell
through to the lowest-ID in-progress Task in R-006. The executor agent
followed the stale prose until the owner pointed at the board. No role owns
advancing the router after an owner closure, and no surface says that the
router selects finished records. The owner added this to O-014's scope on
2026-09-23 instead of opening a separate Issue.

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
  I###`), decoded and validated with the same strictness and explicit `none`
  sentinel as optional Objective/Task selection; an existing router with no
  `issue:` key keeps parsing exactly as before (backward compatible). Release
  selection follows its mandatory rule below and does not accept `none`.
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
- The router always has a valid `release: R###` selection naming a declared
  Release. Missing, blank, `none`, and unknown Release selections produce a
  clear diagnostic; they never silently show all Objectives.
- Every live Objective has a valid `release: R###` reference to a declared
  Release, including planned and done Objectives. Missing, blank, `none`, and
  unknown references produce a clear diagnostic.
- When the router selects no Objective, Next falls back to the selected
  Release's search in ID order. There is no unscoped project-wide fallback.
- Apart from mandatory Release context and the selected-Objective rule above,
  selection behavior for valid Release/Objective/Task values, the `ResolveNext`
  ladder's precedence, and non-TTY parity are unchanged for projects that
  never select an Issue.
- There is one next action. The router no longer carries a free-text
  `next_action`; it holds only `state` and the Release/Objective/Task/Issue
  selection. An existing router that still has `next_action` keeps loading,
  and the field is ignored and reported once as retired (for example by
  doctor), never rendered as a competing instruction.
- Agents obtain the next action from the same `data.Next` projection the
  board shows, by running the read-only `savepoint resume`. AGENTS.md, the
  router template, and the phase skills name `savepoint resume` as the one
  permitted CLI command for agents and replace "act on the router's
  next_action" with "act on `savepoint resume`'s Next".
- When the router selects a Task or Objective that is already `done`, the
  board's Next area, non-TTY output, `savepoint resume`, and doctor all state
  that the selection is stale and name the record, alongside whatever Next
  the records support. The resolver never silently substitutes other work
  without that diagnostic.
- Advancing the router after an owner closure has a named owner in the
  skills: the role that records or acts on the closure updates the selection,
  and the stale-selection diagnostic catches any miss.
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
  reports. Keep the rule consistent with O-020, which requires the canonical
  Next resolver to keep respecting explicit router selection.
- Release selection and Objective membership are mandatory under the owner's
  confirmed requirement. Keep membership derived from `Objective.release`;
  the router must name a declared Release, and Next/board views remain scoped
  to it. Existing data missing these references needs an explicit, actionable
  repair path rather than a silent unscoped fallback.
- Keep `router.md`'s YAML shape backward compatible: an existing V2 router
  with no `issue:` key must keep decoding exactly as it does today.
- `data.ResolveNext` is the only source of the next action. The stale-selection
  diagnostic extends `SelectionDiagnostic` (a new kind such as
  `SelectionDone`) rather than a board- or resume-side check, so every
  surface reports it from the same value. `internal/resume` owns its wording.
- Retiring `next_action` removes the only prose channel agents used to
  hand off intent. Settle during Task planning whether any of that intent
  (for example "request a Task Check or record a waiver") is already fully
  expressed by `data.Next`'s rungs, and add a rung or phrase only where it is
  not. Do not reintroduce free text.
- Allowing agents to run `savepoint resume` changes the "never run savepoint
  commands" rule. The exception is `resume` alone, which is read-only and
  performs no writes, so it cannot break the owner-only authority over status
  changes.

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
- Require every router to name a declared Release and every live Objective to
  reference a declared Release; reconcile loading, writing, diagnostics,
  current records, and active guidance with that rule while keeping membership
  derived from `Objective.release`.
- Retiring the router's `next_action` (`internal/data/router_v2.go`, the V2
  router template, this repository's router), a stale-selection diagnostic in
  `internal/data/next.go` rendered by board, resume, and doctor, and
  AGENTS.md, scaffold guidance, and phase-skill updates that make
  `savepoint resume` the agent's source of the next action and name who
  advances the router after an owner closure.
- Design.md Section 8 reconciliation for the exact word vocabulary, the
  selected-Objective precedence, and the single-Next contract.

**Out of scope:**

- Changing Issue lifecycle states, severity, or type vocabulary.
- Changing the existing Task lifecycle words (`Build`/`Test`/`Check`/`Planned`/`Done`).
- Any Objective, Task, Release, or Issue gate or completion-policy change
  (the Next precedence change above decides what to show, not what may
  start or complete).
- Renaming Release to Goals (O-013's scope) or any other unrelated Next-area
  or router surface not named above.
