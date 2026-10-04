---
id: I-137
title: Keep stale selection advice consistent across surfaces
type: defect
status: open
source:
  kind: check
  check: C-965
  actor: {role: checker, session: check-o033-20261004}
  at: '2026-10-04T03:26:27Z'
tasks: [T-103, T-105, T-107]
checks: [C-965]
history:
  - at: '2026-10-04T03:26:27Z'
    actor: {role: checker, session: check-o033-20261004}
    kind: observed
    check: C-965
    note: Initial independent Full Objective Check.
---

# I-137: Keep stale selection advice consistent across surfaces

## Summary

Keep stale selection advice consistent across surfaces before O-033 clearance.

## Evidence

T-103 DW2 requires missing/stale selection to withhold launch advice; T-105 DW2 and T-107 DW3 require surfaces to agree for the same selection/evidence. README also says stale router selection withholds the suggestion.

Use enabled O-001 with independent ready lane Tasks T-002 and T-004; router selects O-001/T-002 and missing Issue I-999. Load once and render resume and T-002 detail. Both see the same Task/index. Expected: same stale-selection constraint/no parallel launch block. Actual: resume withholds via Next.Concurrency, but Task detail offers begin T-002 and copyable Start instructions. A missing router Task T-999 similarly yields no resume advice while Objective detail/plain advertise launches.

`internal/board/v2/detail.go:350`, `lanes.go:79` and `plain.go:112` recompute projection with Enabled only; `internal/data/next.go:376` supplies SelectionDiagnostic. Independent TestO033StaleIssueSameTaskParity and TestO033StaleSelectionParity fail. Existing parity tests use valid router selections; pure selection test does not cover board consumers.

## Proof Needed

For the same selected Objective/Task and stale router evidence, keep withholding/reasons consistent in resume, detail, plain and heading readiness. Replay missing Issue with valid Task and missing Task, normal current selection, feature off, and replan/status-change cases. Preserve ordinary navigation/Next selection and do not add an execution gate. Frozen M3/M6.
