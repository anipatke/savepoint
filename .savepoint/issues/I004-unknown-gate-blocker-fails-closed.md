---
id: I004
title: Unknown gate blockers fall back to NextDependency
type: defect
status: resolved
source:
    kind: migration
    actor:
        role: owner
        session: .savepoint/releases/v2/defects/D005-unknown-gate-blocker-fails-closed.md
    at: "2026-09-20T05:03:29Z"
severity: high
resolution:
    disposition: accepted
    actor: {role: owner, session: user-review-20260922}
    at: "2026-09-22T09:11:47Z"
    reason: >-
        Owner accepts the residual risk for the current V2 release because strict
        decoding makes the unknown-blocker fallback effectively unreachable. Revisit
        this decision when adding a blocker kind or next changing the Next resolver.
history:
    - at: "2026-09-22T09:11:47Z"
      actor: {role: owner, session: user-review-20260922}
      kind: owner_decision
      note: >-
          Accepted the current release risk without a repair or technical CLEAR;
          collision with a future blocker kind remains an explicit revisit trigger.
---
## Migrated from V1

Relocated verbatim from the V1 defect body at `.savepoint/releases/v2/defects/D005-unknown-gate-blocker-fails-closed.md` (release `v2`).

## V1 Body (verbatim)

# D005: Unknown gate blockers fall back to NextDependency

## Symptom

`internal/data/next.go:rungForBlockers()` maps known blocker kinds to their
workflow rungs, but its final fallback returns `NextDependency`. A newly added
or otherwise unhandled blocker can therefore produce a plausible dependency
instruction without being a dependency.

## Expected Behavior

Unknown blocker state must fail closed with an explicit diagnostic or typed
unknown-state outcome. It must never be silently represented as
`NextDependency`.

## Reproduction

1. Construct a blocked `GateDecision` containing a `GateBlocker` kind that
   `rungForBlockers()` does not recognize.
2. Resolve the task's `Next` projection.
3. Observe that the projection reports the dependency rung and downstream
   resume/board wording instead of identifying the unhandled blocker.

## Impact

Users and agents can receive an incorrect workflow instruction, and adding a
new gate blocker can silently change behavior until every consumer happens to
be updated.

## Fix Plan

Confirmed approach: add a typed `NextKind` (for example `next_invalid_state`)
for an unrecognized blocker instead of the `NextDependency` fallback. Its
action phrasing names the unhandled blocker kind and tells the caller to report
it. Board badge mapping and resume wording handle the new kind; `ResolveNext`
stays error-free. Add an adversarial test that injects an unhandled blocker and
verifies the fail-closed result at each rendering surface.

## Acceptance Criteria

- [ ] An unknown blocker cannot resolve to `NextDependency`.
- [ ] The resulting diagnostic names the unhandled blocker and remains
      actionable for the caller.
- [ ] Known blocker mappings retain their current precedence and coverage.

## Resolution Notes

Resolved by explicit owner acceptance. This is a risk waiver, not a repair or
technical `CLEAR`. The waiver expires when a new blocker kind is introduced or
the Next resolver is next changed.
