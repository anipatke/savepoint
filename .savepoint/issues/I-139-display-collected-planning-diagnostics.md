---
id: I-139
title: Display collected planning diagnostics
type: defect
status: open
source:
  kind: check
  check: C-965
  actor: {role: checker, session: check-o033-20261004}
  at: '2026-10-04T03:26:27Z'
tasks: [T-102, T-105, T-107]
checks: [C-965]
history:
  - at: '2026-10-04T03:26:27Z'
    actor: {role: checker, session: check-o033-20261004}
    kind: observed
    check: C-965
    note: Initial independent Full Objective Check.
---

# I-139: Display collected planning diagnostics

## Summary

Display collected planning diagnostics before O-033 clearance.

## Evidence

O-033 Confirmed Design Optional plan records requires visible planning diagnostics; SC8 requires missing/stale advice to state its limits. T-102 DW1 requires named advisory diagnostics.

On an enabled temporary project, T-004 carries planned_writes: ['bad/*.go']. Index validation produces a diagnostic naming the record/file/field/glob. Render Objective/Task advice, resume and plain board. Expected: the actionable nonfatal diagnostic is visible. Actual: text only says "T-004 has unusable planning metadata; see the planning diagnostics". No user-facing surface displays those collected diagnostics; rg of board/resume/main/doctor finds no PlanDiagnostics renderer. Owner cannot see the offending field or cause, and an Objective containing only malformed lane data can produce no advice body at all.

`internal/resume/concurrency.go:26–46` and objectiveAdvice ignore c.Diagnostics. Independent TestO033MalformedDiagnosticVisible fails; existing data tests verify the in-memory diagnostic, while integration only checks generic degradation text.

## Proof Needed

Render sanitized named diagnostics in the shared formatter for the relevant scope, including an otherwise empty/withheld advice body, without blocking execution or demanding metadata repair. Replay malformed glob Task, duplicate/missing-title Objective lane, stale explanation, valid/off/old records, and surface parity. Frozen M1/M6; feature off still omits optional advice.
