---
id: I-101
title: Unparsed complexity files lack verification
type: verification
status: resolved
source:
  kind: check
  check: C-944
  actor: {role: checker, session: o029-full-check-20261001}
  at: '2026-10-01T09:35:25Z'
tasks: [T-062]
checks: [C-944]
guardrail_ids: [TEST-01, TEST-02]
severity: medium
history:
  - at: '2026-10-01T09:35:25Z'
    actor: {role: checker, session: o029-full-check-20261001}
    kind: observed
    check: C-944
    note: Initial Full O-029 Check; see C-944 frozen matrix and embedded independent harness.
  - at: '2026-10-01T10:30:00Z'
    actor: {role: executor, session: o029-issue-repair-20261001}
    kind: deferred
    note: "Not repaired in code. Lizard's headerless CSV lists only function rows, so it cannot tell 'no functions' from 'file not parsed'; any fix needs either a different provider signal or a change to O-029's unparsed-files requirement. Per the Issue's own guidance this is routed to the owner/planner for a decision rather than reinterpreting an unparsed file as an outside-root row."
  - at: '2026-10-01T11:00:00Z'
    actor: {role: owner, session: o029-issue-repair-20261001}
    kind: owner_decision
    note: Owner accepted the Lizard limitation. The no-functions reason now states that the CSV cannot show unparsed files, and the O-029 complexity requirement was reworded to match.
resolution:
  disposition: accepted
  actor: {role: owner, session: o029-issue-repair-20261001}
  at: '2026-10-01T11:00:00Z'
  reason: Owner accepted the limitation that Lizard's CSV cannot report unparsed files.
---

# I-101: Unparsed complexity files lack verification

## Summary

O-029 complexity success condition: unparsed files must be reported as partial, never complexity zero. The per-reader partial/unsupported-variant fixture requirement is also explicit. Frozen M3 unparsed/unsupported and no-functions cells. Classification is Unverified, not a claimed reproduction of a real Lizard parser failure.

## Evidence

reader_lizard.go:49-75 decides partial only from rows whose file path is outside the project. No-functions CSV returns measured zero. decodeLizard at :105-140 reads only function rows; it receives no parsed/unparsed-file inventory. TestLizardReaderValue tests empty and outside-root reports, and TestLizardReaderRejectsMalformedReports tests bad columns/numbers; neither establishes what happens when Lizard cannot parse an in-scope file. T-062 evidence expressly says no real Lizard run was made and lists only those fixtures.

The reviewer verified Lizard's official CSV producer is the supported boundary (https://raw.githubusercontent.com/terryyin/lizard/master/lizard.py). Available reader/test evidence cannot distinguish no functions from no successfully parsed files. A CSV outside-root row is not proof of an unparsed in-scope file. The claimed requirement remains unverified; no analyzer algorithm is requested inside Savepoint.

## Proof Needed

Provide a documented provider/report signal and fixture for an unparsed in-scope file, with an independent scenario proving partial/no false zero through Collect. If headerless CSV cannot carry that information, route the limitation to planning/owner decision to reconcile the existing requirement; do not silently reinterpret unparsed as an outside-root function row.

Repair under this Issue by default and preserve every completed Task status. This record grants no clearance; a later independent Check verifies remediation.
