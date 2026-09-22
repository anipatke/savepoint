---
id: I026
title: Add an owner waiver action for Objective Checks
type: other
status: open
source:
  kind: report
  actor: {role: owner, session: user}
  at: '2026-09-22T10:01:43Z'
severity: high
history:
  - at: '2026-09-22T10:01:43Z'
    actor: {role: owner, session: user}
    kind: observed
    note: The board needs an explicit owner action that can waive an Objective Check instead of requiring the owner to edit evidence manually or remain blocked.
---

# I026: Add an owner waiver action for Objective Checks

## Summary

The owner wants a board action for waiving an Objective Check. Today only
optional Task Checks may be waived; Full Objective Checks are mandatory, and
the existing exception model records owner acceptance of a named Check rather
than absence of the Check itself. The requested action therefore changes the
verification contract and must be designed explicitly rather than borrowing
Task-waiver semantics by accident.

The board should offer one clear, contextual owner action when an Objective is
otherwise ready but its Check requirement is the remaining blocker. The
recorded result must remain visibly an owner waiver—not technical `CLEAR`—and
all downstream Objective and Release behavior must agree on what that waiver
does and does not authorize.

## Evidence

- `.savepoint/Guardrails.md:91` says a Task waiver never replaces the mandatory
  Objective or Release Check.
- `.savepoint/Design.md` and the Check skill repeat mandatory Full Objective
  evidence as the current closure contract.
- `internal/data/evidence_v2.go` models `CheckWaiver` as Task-only evidence;
  Objective gate and completion behavior has no corresponding waiver path.
- The owner explicitly requested a button/key action to waive an Objective
  Check.

## Proof Needed

- Obtain and record the product decision: whether the action creates a true
  Objective-Check waiver, an owner exception, or another explicitly named
  outcome. Define availability before/after a Check exists and after a
  `NEEDS WORK` result.
- Define exact effects on Objective completion, technical clearance,
  dependencies, Release membership/readiness, Issue blockers, freshness,
  rechecks, and later supersession. A waiver must never be rendered as
  independent `CLEAR`.
- Add typed evidence and canonical data-layer resolution; the board must not
  invent gate semantics locally.
- Add one concise contextual key/button and confirmation flow, plus detail,
  Next, resume, doctor, non-TTY, and no-colour wording that clearly attributes
  the decision to the owner.
- Reconcile Guardrails, Design, Check/task skills, AGENTS guidance, templates,
  and closure rules together; do not leave mandatory-Check prose contradicting
  runtime behavior.
- Test allowed and refused states, malformed evidence, dependencies, Release
  completion, stale/superseded evidence, owner attribution, concurrency, and
  persistence; pass all focused and full gates.
