---
id: I-031
title: Frequent V2 reload errors for project records
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
      saying no frontmatter was found, followed by a report that I-029 names a
      missing check C-908. Both exact diagnostics are recorded in the Evidence
      section. Owner expects to add further examples.
  - at: '2026-09-23T03:57:02Z'
    actor: {role: owner, session: user-report}
    kind: observed
    note: >-
      Owner reported another RELOAD error: issue I-030 names missing task T-018.
      Exact diagnostic is recorded in the Evidence section.
  - at: '2026-09-23T06:50:30Z'
    actor: {role: owner, session: user-report}
    kind: observed
    note: >-
      Owner reported a RELOAD error decoding router.md because of a YAML
      mapping-values parse error on line 5. Exact diagnostic is recorded in the
      Evidence section.
---

# I-031: Frequent V2 reload errors for project records

## Summary

The owner reports frequent RELOAD diagnostics while loading V2 project records, including Issue records and the router. These
repeated errors disrupt the experience and need to be resolved before release.
The owner will add further examples as they occur. This record captures the
reported symptoms without asserting a cause.

## Evidence

- `RELOAD parse error for issues/I-029-windows-migration-tests-have-two-testmain-functions.md: no frontmatter found`
- `RELOAD v2 issue link names a record that does not exist: issues/I-029-windows-migration-tests-have-two-testmain-functions.md: issue I-029 names missing check C-908`
- `RELOAD v2 issue link names a record that does not exist: issues/I-030-windows-full-suite-has-platform-failures.md: issue I-030 names missing task T-018&#x20;`
- `RELOAD /home/user/code/savepoint/.savepoint/router.md: v2 record could not be decoded: router state: yaml: line 5: mapping values are not allowed in this context`

## Proof Needed

Determine why these reported cases produce reload errors, resolve the failure,
and show that repeated reloads no longer produce these diagnostics. Include
regression evidence for all reported cases; append additional observations to
this Issue as the owner supplies them.
