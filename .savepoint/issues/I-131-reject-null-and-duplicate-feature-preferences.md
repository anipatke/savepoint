---
id: I-131
title: Reject null and duplicate feature preferences
type: defect
status: open
source:
  kind: check
  check: C-963
  actor: {role: checker, session: check-o037-20261004}
  at: '2026-10-04T02:31:20Z'
tasks: [T-108]
checks: [C-963]
guardrail_ids: [DATA-03, CFG-01]
history:
  - at: '2026-10-04T02:31:20Z'
    actor: {role: checker, session: check-o037-20261004}
    kind: observed
    check: C-963
    note: Initial independent Full Objective Check.
---

# I-131: Reject null and duplicate feature preferences

## Summary

Reject null and duplicate feature preferences before O-037 clearance.

## Evidence

T-108 Done When 1 and 3 promise boolean decoding and malformed-configuration refusal. O-037 SC1 and SC3 require one saved boolean and safe diagnostics.

`features: {parallel_planning: null}` and an empty `parallel_planning:` decode successfully as false. A block containing `parallel_planning: false` followed by `parallel_planning: true` succeeds as true. Expected: named malformed-preference/duplicate diagnostic. Actual: silent default or last-value-wins. `internal/data/feature_preferences.go:40–49` bypasses mapping duplicate validation and accepts null through Decode(bool).

Independent `TestO037IndependentDecode` in the C-963 overlay reproduces all three failures. Existing malformed tests cover strings/lists/maps but omit null and duplicate keys.

## Proof Needed

Reject null/empty and duplicate parallel_planning keys through ConfigReader.Read and direct UnmarshalYAML. Test false/true, absent key, unknown unrelated keys, top-level duplicates, non-boolean values and malformed-write refusal; original bytes must remain untouched. Re-run frozen C-963 M1 and M2 cells.
