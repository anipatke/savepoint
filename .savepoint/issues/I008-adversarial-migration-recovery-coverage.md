---
id: I008
title: Migration recovery lacks the complete adversarial integration matrix
type: defect
status: resolved
source:
    kind: migration
    actor:
        role: owner
        session: .savepoint/releases/v2/defects/D009-adversarial-migration-recovery-coverage.md
    at: "2026-09-20T05:03:29Z"
severity: high
resolution:
    disposition: accepted
    actor: {role: owner, session: user-review-20260922}
    at: "2026-09-22T09:13:35Z"
    reason: >-
        Owner accepts the residual migration-recovery risk because Savepoint has one
        user, this repository's V1-to-V2 cutover is complete, and ordinary operation is
        now V2-only. Building the full interruption and cross-platform matrix is not
        proportionate to the remaining legacy migration exposure.
history:
    - at: "2026-09-22T09:13:35Z"
      actor: {role: owner, session: user-review-20260922}
      kind: owner_decision
      note: >-
          Accepted the residual legacy migration risk without a repair or technical
          CLEAR; the completed single-user cutover does not justify the requested
          adversarial recovery matrix.
---
## Migrated from V1

Relocated verbatim from the V1 defect body at `.savepoint/releases/v2/defects/D009-adversarial-migration-recovery-coverage.md` (release `v2`).

## V1 Body (verbatim)

# D009: Migration recovery lacks the complete adversarial integration matrix

## Symptom

The migration suite already covers important operation lifecycles, an
interruption/resume path, source edits, installed-file tampering, idempotent
second runs, and Windows replacement primitives. It does not yet prove the
whole command-level recovery contract at every write boundary or across the
requested partial-deletion, malformed-project, Windows/WSL, and restart
combinations.

## Expected Behavior

Migration must be demonstrably recoverable and data-preserving when execution
stops at each write boundary and when the source tree or installed outputs are
changed before recovery. The evidence must include malformed V1 projects,
partial deletion, repeated runs, and the supported Windows and WSL
environments.

## Reproduction

Exercise the migration fixture with deterministic interruption at each backup,
stage, install, verify, activation, and cleanup boundary, then restart the
command and vary the state between attempts:

1. Edit source files between preview and recovery.
2. Edit installed outputs during recovery.
3. Delete only part of an output or backup set.
4. Run the operation twice and compare the final tree and manifest.
5. Supply malformed V1 records and verify a named, write-free refusal.
6. Repeat the matrix on Windows/NTFS and on WSL where supported.

## Impact

The migration architecture may be sound, but untested crash and platform
windows can still cause data loss, an unrecoverable journal, or an ambiguous
operator path during the release cutover.

## Fix Plan

Confirmed approach, phased. Phase 1 (Linux): add fault injection at every
publish boundary, invoke recovery through a fresh process or command call, and
assert byte preservation and journal convergence across source edits, installed
output edits, partial deletion, repeated runs, and malformed V1 projects. Keep
the existing focused lifecycle tests as unit evidence and add the missing
end-to-end matrix around them. Phase 2: exercise Windows/NTFS and WSL through
platform jobs, or document the supported boundary with harness evidence.

## Acceptance Criteria

- [ ] Every migration write boundary has an interruption-and-recovery case.
- [ ] Source edits, installed-output edits, partial deletion, repeated runs,
      and malformed V1 projects have explicit assertions.
- [ ] Windows/NTFS and WSL behavior is exercised or its supported boundary is
      explicitly documented with evidence.
- [ ] Recovery either converges byte-for-byte or returns a named refusal that
      preserves user data.

## Resolution Notes

Resolved by explicit owner acceptance. This is a risk waiver, not a repair or
technical `CLEAR`. The decision is based on the completed single-user cutover
and V2-only ordinary operation.
