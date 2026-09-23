---
id: I-031
title: Frequent V2 board load errors for project records
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
  - at: '2026-09-23T20:33:41Z'
    actor: {role: owner, session: user-report}
    kind: observed
    note: >-
      Owner reported that npx savepoint board rejected objective O-001 as an
      invalid global ID. An npm advisory endpoint notice appeared with the
      command output. Both lines are recorded in the Evidence section.
  - at: '2026-09-23T20:51:46Z'
    actor: {role: owner, session: user-report}
    kind: observed
    note: >-
      Owner reported a parse error for Issue I-038: no closing frontmatter
      delimiter was found. The board displayed that no board is drawn because
      the project's records did not load.
---

# I-031: Frequent V2 board load errors for project records

## Summary

The owner reports frequent V2 board load errors during startup and reload,
including diagnostics for Issue records, the router, and an Objective. These
errors disrupt the experience and need to be resolved before release.
The owner will add further examples as they occur. This record captures the
reported symptoms without asserting a cause.

## Evidence

- `RELOAD parse error for issues/I-029-windows-migration-tests-have-two-testmain-functions.md: no frontmatter found`
- `RELOAD v2 issue link names a record that does not exist: issues/I-029-windows-migration-tests-have-two-testmain-functions.md: issue I-029 names missing check C-908`
- `RELOAD v2 issue link names a record that does not exist: issues/I-030-windows-full-suite-has-platform-failures.md: issue I-030 names missing task T-018&#x20;`
- `RELOAD /home/user/code/savepoint/.savepoint/router.md: v2 record could not be decoded: router state: yaml: line 5: mapping values are not allowed in this context`
- `npx savepoint board` also printed: `npm notice This endpoint is being retired. Use the bulk advisory endpoint instead. See the following docs for more info: https://api-docs.npmjs.com/#tag/Audit`
- `npx savepoint board` failed with: `board: [invalid_v2] V2 project cannot be loaded safely: v2 record has a malformed or invalid global ID: objectives/O-001-release-validation-cutover/Objective.md: objective id "O-001" must match O plus at least three digits; repair the named record before cutover`
- `parse error for issues/I-038-npx-board-rejects-canonical-objective-id.md: no closing frontmatter delimiter found`
- `No board is drawn: this project's records did not load.`

## Proof Needed

Determine why these reported cases produce board load errors, resolve the
failures, and show that board startup and repeated reloads no longer produce
these diagnostics. Include regression evidence for all reported cases;
append additional observations to this Issue as the owner supplies them.
