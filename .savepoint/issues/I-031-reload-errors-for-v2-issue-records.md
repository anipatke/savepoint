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
  - at: '2026-09-23T22:17:20Z'
    actor: {role: owner, session: user-report}
    kind: observed
    note: >-
      Owner reported another RELOAD issue link diagnostic: issue I-040 names
      missing check C-914. Exact diagnostic is recorded in the Evidence section.
  - at: '2026-09-25T10:22:13Z'
    actor: {role: owner, session: user-report}
    kind: observed
    note: >-
      Owner reported that a RELOAD rejected Task T-046 because
      owner_validation.accepted_by is missing. The reported diagnostic is
      recorded in the Evidence section.
  - at: '2026-09-25T20:29:12Z'
    actor: {role: planner, session: i031-analysis-20260926}
    kind: observed
    note: >-
      Analysis only, no code changed. The reports have three independent
      causes: the board reloading mid-way through non-atomic agent edits, a
      stale dist/npm binary behind npx, and a retired free-text router field.
      Every reported record loads cleanly at its committed state. Causes and
      recommended fixes are recorded in the Analysis section.
---

# I-031: Frequent V2 board load errors for project records

## Summary

The owner reports frequent V2 board load errors during startup and reload,
including diagnostics for Issue, Task, and Objective records and the router. These
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
- `RELOAD v2 issue link names a record that does not exist: issues/I-040-data-loadproject-unreachable-v1-dispatch.md: issue I-040 names missing check C-914`
- `RELOAD v2 record is missing a required field: objectives/O-015-owner-advance-issues/tasks/T-046-use-space-and-backspace-in-the-issues-panel.md: task T-046 missing required field owner_validation.accepted_by`

## Analysis

Recorded 2026-09-26 (analysis only; no code changed). Every reported record
loads cleanly at its committed state; a full load of this project takes about
13 ms and succeeds. The ten evidence lines have three independent causes.

### A. Reload reads a half-finished multi-step edit

The board reloads 100 ms after the last filesystem event
(`internal/board/v2/watch.go`), and one invalid record fails the whole load
(`LoadV2Index`, `internal/data/project.go`). The board's own writes are atomic
(temp file plus rename, `internal/data/write.go`); agent and editor writes are
not, and one logical change often spans several edits and several files.

- I030 named T018 after the owner removed T018 and T019 from O016, until the
  Issue was updated.
- I-040 names C-914, and I029 names C908, where both records first appear in
  one commit (`ab2ef37`, `cd7ef36`): the Issue was written before its Check.
- T-046: `accepted_check` without `accepted_by` is rejected
  (`decodeOwnerValidationV2`, `internal/data/evidence_v2.go`); an edit that
  adds them one at a time passes through that state. The committed file has
  both.
- I029 "no frontmatter found" and I-038 "no closing frontmatter delimiter
  found": the file was empty or partly written when read. Which writer produced
  it is not determined from git history.

Latent contributor, not observed: `projectLoadedMsg` carries no sequence, so an
older failed load finishing after a newer good one would leave the RELOAD line
up until the next change. With a 13 ms load against a 100 ms debounce this is
unlikely.

### B. `npx savepoint` runs a stale local binary (O-001)

Inside this repo `npx savepoint` resolves to `bin/savepoint.js`, which launches
`dist/npm/<platform>/savepoint`. That binary was built 2026-09-22 19:00, before
O-018 hyphenated identities (`bb03eaf`), and still contains the old rule text
"must match O plus" with no "must match O- plus". It rejects valid `O-001`. The
npm advisory notice is unrelated npm output.

### C. Free-text router field broke YAML (router line 5)

The router's line 5 was an unquoted prose `next_action:` value; any `: ` in it
yields "mapping values are not allowed". Committed router `a4e3db6` has this
shape. `next_action` was removed in `f159fc8` (2026-09-24), so the cause is
already retired.

### Recommended fixes

1. Board: on a failed reload keep the last good board, retry once after about
   300-500 ms, and show RELOAD only if the failure persists; add a load
   sequence so a stale result cannot overwrite a newer one. Covers cause A.
2. npx: rebuild `dist/npm` as part of `make build` so the local npm launcher
   cannot lag the source. Covers cause B.
3. Router: no action.

Deferred: loading a dangling Issue link as a per-Issue warning instead of a
project-wide failure. Reconsider only if link errors persist after fix 1.

## Proof Needed

Determine why these reported cases produce board load errors, resolve the
failures, and show that board startup and repeated reloads no longer produce
these diagnostics. Include regression evidence for all reported cases;
append additional observations to this Issue as the owner supplies them.
