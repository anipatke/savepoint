---
id: I009
title: Concurrent Check and Issue creation has no reload-and-retry contract
type: defect
status: open
source:
    kind: migration
    actor:
        role: owner
        session: .savepoint/releases/v2/defects/D010-concurrent-check-issue-id-allocation.md
    at: "2026-09-20T05:03:29Z"
severity: high
---
## Migrated from V1

Relocated verbatim from the V1 defect body at `.savepoint/releases/v2/defects/D010-concurrent-check-issue-id-allocation.md` (release `v2`).

## V1 Body (verbatim)

# D010: Concurrent Check and Issue creation has no reload-and-retry contract

## Symptom

`CreateCheckV2()` and `CreateIssueV2()` allocate one past the highest ID in a
caller-provided index. Their create-only hard-link publication correctly
prevents overwriting an existing record, but two agents can read the same
index, choose the same next ID, and leave the losing creator with only a
collision error. The API has no tested, clean reload/retry path for that race.

## Expected Behavior

Concurrent creation must preserve create-only/no-overwrite guarantees and give
the losing writer a typed, actionable collision result that tells the caller to
reload the index and retry allocation.

## Reproduction

1. Load one index snapshot containing the same highest Check or Issue ID for
   two concurrent creators.
2. Start `CreateCheckV2()` or `CreateIssueV2()` from both callers with that
   snapshot.
3. Observe that both choose the same next ID; one publishes and the other
   receives a collision error without an integration-level retry contract.

## Impact

Parallel agents can fail with an ugly or ambiguous write error, and callers
that do not know to reload can lose a requested Check or Issue creation even
though the underlying record allocation is safely create-only.

## Fix Plan

Confirmed approach: define a typed collision error for `CreateCheckV2` and
`CreateIssueV2` that tells the caller to reload the index and retry allocation,
preserving the existing create-only no-overwrite guarantee. Then add concurrent
tests that assert exactly one file wins, no bytes are overwritten, and the
losing caller can reload and succeed with the next ID.

## Acceptance Criteria

- [ ] Concurrent Check creation has deterministic collision and retry
      behavior.
- [ ] Concurrent Issue creation has the same behavior.
- [ ] The losing writer receives an actionable typed diagnostic and no partial
      record is left behind.
- [ ] Tests prove existing records remain byte-identical.

## Resolution Notes

Pending.
