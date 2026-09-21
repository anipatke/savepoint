---
id: I016
title: E50 cutover rollback snapshot is no longer available
type: verification
status: resolved
source:
  kind: check
  check: C905
  actor: {role: checker, session: e50-objective-check-20260921}
  at: '2026-09-21T08:06:07Z'
tasks: [T001, T002]
checks: [C905]
severity: high
resolution:
  disposition: accepted
  actor: {role: owner, session: owner-accept-20260921}
  at: '2026-09-21T19:15:00Z'
  reason: >-
    Owner accepts the residual risk rather than gating cutover on a fresh
    Check re-verifying the recreated backup: low residual risk, not worth
    another verification loop. Not a technical CLEAR — the repair
    (repository.bundle + worktree.tar + SHA256SUMS at
    /home/user/code/savepoint-cutover-backups/0e7c9ed/) is not independently
    re-proven by this decision.
history:
  - at: '2026-09-21T08:06:07Z'
    actor: {role: checker, session: e50-objective-check-20260921}
    kind: observed
    note: The recorded rollback snapshot and SHA256SUMS are absent at the only named location.
  - at: '2026-09-21T18:47:00Z'
    actor: {role: executor, session: remediation-20260921}
    kind: repair_attempted
    note: >-
      Root cause: the original backup was written to `/tmp/savepoint-cutover-backups/0e7c9ed/`,
      which is ephemeral and deviates from the archived runbook's own documented location
      (`../savepoint-cutover-backups/<sha>/`, outside the project root). Recreated the
      snapshot at the runbook-correct durable location, `/home/user/code/savepoint-cutover-backups/0e7c9ed/`:
      `repository.bundle` (`git bundle create --all`, `git bundle verify` passed, contains
      commit 0e7c9ed and a complete history), `worktree.tar` (`git archive 0e7c9ed`, a
      deterministic reconstruction of that commit's tree), `NOTE.md` (documents the
      reconstruction and why the original worktree/index patches can't be recovered), and
      `SHA256SUMS` covering all three (`sha256sum -c` passed). Recorded in T001's Context
      Log, which supersedes the stale `/tmp` reference there. Not closing this Issue —
      a fresh Check should independently verify the new location and checksums.
  - at: '2026-09-21T19:15:00Z'
    actor: {role: owner, session: owner-accept-20260921}
    kind: owner_decision
    note: >-
      Owner accepted the risk instead of commissioning a fresh Check to
      re-verify the recreated backup and checksums. See resolution.reason.
---

# I016: E50 cutover rollback snapshot is no longer available

## Summary

T002's migration-accountability criterion requires the verified backup to remain available. T001 records that backup only at `/tmp/savepoint-cutover-backups/0e7c9ed/`. The Full Objective Check found neither that directory nor its `SHA256SUMS` file.

## Evidence

- Violated criterion: T002 acceptance criterion at line 75.
- Recorded location: T001 Context Log at line 81.
- Reproduction: `test -d /tmp/savepoint-cutover-backups/0e7c9ed` and `test -f /tmp/savepoint-cutover-backups/0e7c9ed/SHA256SUMS` both returned false on 2026-09-21.
- Expected: the verified rollback snapshot and checksum inventory remain available through cutover acceptance.
- Actual: the sole named rollback snapshot is absent.
- Missing evidence: no durable repository or owner-controlled location is recorded as a replacement.

## Proof Needed

Record a durable, available rollback snapshot for the reviewed cutover revision, verify its checksum inventory independently, and update the handoff evidence to name that location. A later Check must verify the files rather than accepting a path claim.
