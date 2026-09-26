---
id: I-071
title: Migrate keeps the V1 router's retired next_action line
type: defect
status: resolved
source:
  kind: report
  actor: {role: owner, session: user}
  at: '2026-09-26T03:34:20Z'
severity: low
history:
  - at: '2026-09-26T03:34:20Z'
    actor: {role: owner, session: user}
    kind: observed
    note: >-
      galaxy's converted router.md kept the V1 line "next_action: Run the E06
      epic audit (savepoint-audit-epic)...", and doctor then warned
      [router-next-action-retired] with the repair "Delete the next_action
      line".
  - at: '2026-09-26T04:37:46Z'
    actor: {role: executor, session: codex}
    kind: repair_attempted
    note: >-
      Removed next_action from the V2 router render model. Converted routers
      retain their mapped state, Goal, Objective, Task, and surrounding prose;
      retired V1 selection keys are omitted. Tests confirm doctor reports no
      router-next-action-retired warning on migrated fixtures. Reviewed and
      regenerated all five migration goldens. make build, make test-fast, and
      the final make test-full passed.
  - at: '2026-09-26T04:41:01Z'
    actor: {role: owner, session: user}
    kind: owner_decision
    note: >-
      Owner instructed the executor to mark this Issue resolved after the
      repair merged to v2 with make build, make test-fast, and make test-full
      passing on the merged branch. No independent Check was run and no
      technical CLEAR is implied.
resolution:
  disposition: accepted
  actor: {role: owner, session: user}
  at: '2026-09-26T04:41:01Z'
  reason: Owner accepted the repair merged to v2 in 5d3075f.
---

# I-071: Migrate keeps the V1 router's retired next_action line

## Summary

The migrated `router.md` carries the V1 `next_action` field forward. V2
retired it, so every freshly migrated project that had one starts with a
doctor warning. Its text also names V1 skills that no longer exist.

## Evidence

- galaxy `.savepoint/router.md` after apply, line 9:
  `next_action: Run the E06 epic audit (savepoint-audit-epic)...`.
- doctor: `[router-next-action-retired] router.md has a retired next_action
  field`.

## Proof Needed

- The router document conversion drops `next_action`, and any other key V2
  retired, while preserving every live key and the router body.
- doctor reports no router warning on a freshly migrated fixture whose V1
  router had `next_action`.
- Tests cover it; `make build && make test-fast` pass.
