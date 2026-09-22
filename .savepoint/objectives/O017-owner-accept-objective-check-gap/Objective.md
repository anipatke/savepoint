---
id: O017
title: Give the owner an explicit Accept action for an Objective's Check gap
status: planned
---

# O017: Give the owner an explicit Accept action for an Objective's Check gap

## Outcome

The owner can accept an Objective forward past its Check gap with one board
action and a confirmation step, in either of two distinct situations — a
Check ran and left issues unresolved, or no Check has ever run — without the
result ever being rendered or resolved as technical `CLEAR`.

## Why

I026 asked for a board action to waive an Objective Check instead of the
owner editing evidence by hand. Design review found the mechanism mostly
already exists: `Exception` already applies to Objectives
(`objective_gate_v2.go:68`) and already covers "a Check ran, found problems,
owner accepts the remainder" without weakening the mandatory Full Objective
Check — it requires a real Check to name. What it lacks is a board action to
record one.

A second case `Exception` cannot express at all is "no Check has ever run."
That is a materially bigger authorization — zero independent verification —
and needs its own validated shape and its own, louder wording so it is never
confused with the first case.

## Success Conditions

- One "Accept" board action on an Objective, gated by scenario:

  | Scenario | Show Accept? | Behaviour when Accept hit |
  | --- | --- | --- |
  | Tasks not all done yet | No | — |
  | No Check has ever run, Tasks done | Yes | Records "Accepted without Check" — no Check ID, owner-typed reason required |
  | Check ran, found open issues, Tasks done | Yes | Records "Accepted, N issues outstanding" — lists the open issues, tied to that Check, owner-typed reason required |
  | Check ran, fully clear | No | — (already resolves `Allowed: true` today, no action needed) |
  | Already accepted, no new Check since | No | — (button replaced by an "Accepted" status) |
  | New Check runs after a prior acceptance | Yes (again) | Old acceptance is void; treated as the "found open issues" row against the new Check |

- Accepting always requires an owner-typed reason; there is no confirm-only
  path.
- The confirmation dialog lists the actual open issues (or states plainly
  that no Check has run) — it is not a bare "are you sure."
- `Exception` gains a required discriminator distinguishing "found problems,
  accepted" from "no Check, accepted"; the "no Check" kind carries no Check
  ID, the other kind must name the Objective's real latest Check exactly as
  today.
- Board, doctor, resume, non-TTY, and no-colour output all label the two
  kinds differently ("ACCEPTED — N ISSUES OUTSTANDING" vs. "ACCEPTED WITHOUT
  CHECK") and never as `CLEAR` or `PASSED`.
- Objective dependents and Release membership see "accepted, not clear," not
  plain clearance — reusing the existing `AllowedByException` distinction
  already wired into Task, Objective, and Release gate resolvers.
- A later Check on the same Objective automatically supersedes a prior
  acceptance (existing supersession behavior via `LatestCheck` matching
  extends to the new discriminator).
- Guardrails, Design, the Check/task skills, AGENTS guidance, and templates
  are reconciled so no mandatory-Check prose contradicts this action.
- Tests cover both kinds, the disabled/hidden states, supersession,
  dependency and Release propagation, and malformed/missing evidence.

## Architectural Considerations

- `internal/data` remains the sole owner of this semantics: the board must
  read a resolved decision, never decide gate behavior locally.
- Prefer extending the existing `Exception` type with a `Kind` field over a
  new sibling type — the decode/validate/render plumbing already exists for
  one kind and should not be duplicated.
- `applicableException`'s current rule (exception's Check must equal
  `index.LatestCheck`) stays the validation path for the "found problems"
  kind; the "no Check" kind instead requires `LatestCheck` to be empty and
  requires no Check field at all — these are different preconditions, not the
  same field left blank.
- This Objective governs the Objective-scope action only. Task-scope already
  has its own, separate mechanism (`CheckWaiver`, Task-only, waives the
  *optional* Task Check before one exists) — it is not touched here. A
  Release-scope Accept button is not requested by the owner and is out of
  scope; the underlying `Exception` wiring for Releases already exists
  in `release_gate_v2.go` if a future Objective wants to add that button.

## Boundaries

**In scope:**

- `Exception.Kind` discriminator and its decode/validation rules for both
  kinds.
- The Objective gate resolver's handling of the "no Check" kind.
- The board Accept action, its confirmation dialog (open-issue list or
  no-Check notice, reason field), and per-scenario visibility.
- Board, doctor, resume, non-TTY, and no-colour rendering for both kinds.
- Dependency and Release propagation of "accepted, not clear."
- Guardrails/Design/skill/template reconciliation.

**Out of scope:**

- A Task- or Release-scope Accept button.
- Any change to the mandatory Full Objective Check or Full Release Check
  requirement itself — this action accepts a gap in that requirement, it
  does not remove the requirement.
- Blocking or throttling Accept based on issue count or severity — the
  dialog surfaces counts for owner judgment only, per prior design
  discussion; it never becomes a second gate.

## Originating Issue

I026 — Add an owner waiver action for Objective Checks.
