---
id: I-150
title: Unknown and malformed provider lines disappear from the unparsed count
type: defect
status: open
source:
  kind: check
  check: C-978
  actor: {role: checker, session: check-o045-20261011}
  at: '2026-10-10T21:10:26Z'
tasks: [T-128]
checks: [C-978]
guardrail_ids: [TEST-01, TEST-02]
history:
  - at: '2026-10-10T21:10:26Z'
    actor: {role: checker, session: check-o045-20261011}
    kind: observed
    check: C-978
    note: Initial Full Objective Check of O-045.
---

# I-150: Unknown and malformed provider lines disappear from the unparsed count

## Summary

Keep raw bytes as today, but classify unknown subtypes and malformed message structures as unparsed. Actual behavior silently treats them as ignored housekeeping/plain text, concealing incomplete interpretation. Known housekeeping and stream_event may remain intentionally ignored.

Violated requirements: O-045 SC2; T-128 DW2.

## Evidence

- Location: `internal/run/claude.go:99-118`.
- Smallest reproduction: ParseClaudeLine on {"type":"system","subtype":"future_unknown"}, {"type":"assistant","message":null}, and {"type":"assistant","message":{"content":7}} each returns no events. The deterministic provider-shapes probe records these outputs.
- Independent harness and full output are embedded in C-978. Probe PASS means the observation ran successfully, not that the requirement passed.

## Proof Needed

Replay these three lines through Start and verify exact raw log plus unparsed increments, alongside valid plain-text messages, known housekeeping, and stream deltas. Existing tests cover unknown top-level types, not these bypass shapes.
