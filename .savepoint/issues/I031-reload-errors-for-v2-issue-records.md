---
id: I031
title: Repeated reload errors for V2 Issue records
type: defect
status: open
source:
  kind: report
  actor: {role: owner, session: user-report}
  at: '2026-09-23T01:53:02Z'
history:
  - at: '2026-09-23T01:53:02Z'
    actor: {role: owner, session: user-report}
    kind: observed
    note: >-
      Owner reported frequent RELOAD errors for an Issue record: a parse error
      saying no frontmatter was found, followed by a report that I029 names a
      missing check C908. Both exact diagnostics are recorded in the Evidence
      section. Owner expects to add further examples.
  - at: '2026-09-23T03:57:02Z'
    actor: {role: owner, session: user-report}
    kind: observed
    note: >-
      Owner reported another RELOAD error: issue I030 names missing task T018.
      Exact diagnostic is recorded in the Evidence section.
---

# I031: Repeated reload errors for V2 Issue records

## Summary

The owner reports frequent RELOAD diagnostics for V2 Issue records. These
repeated errors disrupt the experience and need to be resolved before release.
The owner will add further examples as they occur. This record captures the
reported symptoms without asserting a cause.

## Evidence

- `RELOAD parse error for issues/I029-windows-migration-tests-have-two-testmain-functions.md: no frontmatter found`
- `RELOAD v2 issue link names a record that does not exist: issues/I029-windows-migration-tests-have-two-testmain-functions.md: issue I029 names missing check C908`
- `RELOAD v2 issue link names a record that does not exist: issues/I030-windows-full-suite-has-platform-failures.md: issue I030 names missing task T018&#x20;`

## Proof Needed

Determine why these reported cases produce reload errors, resolve the failure,
and show that repeated reloads no longer produce these diagnostics. Include
regression evidence for all reported cases; append additional observations to
this Issue as the owner supplies them.
