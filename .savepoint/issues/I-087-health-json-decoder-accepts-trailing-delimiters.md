---
id: I-087
title: Health JSON decoder accepts trailing delimiters
type: defect
status: resolved
source:
  kind: check
  check: C-939
  actor: {role: checker, session: o027-full-check-20261001}
  at: '2026-09-30T21:34:30Z'
tasks: [T-051, T-052]
checks: [C-939, C-940]
guardrail_ids: [DATA-03, TEST-02]
severity: medium
resolution:
  disposition: verified
  check: C-940
  actor: {role: checker, session: o027-recheck-20261001}
  at: '2026-09-30T22:21:36Z'
  reason: Repair reproduced and proven by CLEAR Full Objective re-check C-940.
history:
  - at: '2026-09-30T21:34:30Z'
    actor: {role: checker, session: o027-full-check-20261001}
    kind: observed
    check: C-939
    note: Reproduced during the initial Full O-027 Check; see frozen matrix and embedded harness in C-939.
  - at: '2026-09-30T22:21:36Z'
    actor: {role: checker, session: o027-recheck-20261001}
    kind: rechecked
    check: C-940
    note: C-939 reproduction and independent re-check probes pass; see C-940 closure map.
---

# I-087: Health JSON decoder accepts trailing delimiters

## Summary

T-051 malformed-input criteria 5–6 and T-052 criterion 1 require malformed records to be refused. Frozen M1/M3.

Append }, ], " } garbage", or " ] {}" to either valid JSON fixture. json.Valid reports false; DecodeConfig and DecodeSnapshot return nil error in every case. The unchanged snapshot content hash still matches.

Evidence: internal/codehealth/config.go:159 uses Decoder.More after decoding the root value; More answers array/object membership and returns false on a closing delimiter instead of enforcing EOF. Both Store loaders call these decoders. Embedded harness M1_trailing_delimiters reproduces all eight decoder failures.

Expected: any non-whitespace trailing content is a named malformed-record error. Actual: corrupt records load successfully and can enter retention/history operations. Existing TestDecodeConfigRejectsMalformedInput covers a trailing object but not closing delimiters; snapshot malformed tests omit trailing data.

## Proof Needed

Require EOF after the one root record and test closing delimiters, garbage, a second value, and valid whitespace through both decoders and Store loaders.

Repair directly under this Issue and retain completed Task statuses. A fresh independent Check verifies the repair; this record grants no clearance.

