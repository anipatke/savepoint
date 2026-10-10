---
id: I-151
title: Run command arguments retain URL credentials in run.yml
type: defect
status: open
source:
  kind: check
  check: C-978
  actor: {role: checker, session: check-o045-20261011}
  at: '2026-10-10T21:10:26Z'
tasks: [T-127, T-128]
checks: [C-978]
guardrail_ids: [TEST-01]
history:
  - at: '2026-10-10T21:10:26Z'
    actor: {role: checker, session: check-o045-20261011}
    kind: observed
    check: C-978
    note: Initial Full Objective Check of O-045.
---

# I-151: Run command arguments retain URL credentials in run.yml

## Summary

Execute the original argv, but strip URL credentials and terminal controls from the saved derived command text. The raw event log remains verbatim by its separate contract. Actual record stores Spec.Argv without sanitization.

Violated requirements: O-045 SC5; T-128 DW5.

## Evidence

- Location: `internal/run/run.go:95-99`.
- Smallest reproduction: Start the fake agent with an additional argv entry https://user:checksecret@example.com. The child succeeds; the stored run.yml contains checksecret (synthetic credential).
- Independent harness and full output are embedded in C-978. Probe PASS means the observation ran successfully, not that the requirement passed.

## Proof Needed

Verify stored command and other derived text remove URL credentials/controls without changing argv delivered to the subprocess. Existing sanitization tests exercise provider fields but omit the top-level command.
