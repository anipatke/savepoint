---
id: C-934
scope: {kind: task, id: T-022}
result: CLEAR
checked_by: {role: checker, session: i037-check-20260926}
executed_session: i031-i037-repair-20260926
checked_at: '2026-09-25T21:30:30Z'
reviewed:
  head_commit: 59aa3fb3bd5516960538fa85fc210eef180ecc70
  files:
    - internal/doctor/v2_runtime.go
    - internal/doctor/checks_test.go
    - .savepoint/issues/I-037-doctor-names-unsupported-schema-version-malformed.md
  dependencies: []
issues: []
supersedes: null
---

# C-934: I-037 schema-version diagnostic re-check

## Result

**CLEAR.** I-037's repair is independently verified. `RunV2Checks` reports
well-formed unsupported integer versions with
`[schema-version-unsupported]` and the unsupported-version repair, while
non-integer values retain `[schema-version-malformed]` and its distinct
repair. No material blocker remains inside the frozen scope.

## Closure Map

| Issue | Result | Evidence |
| --- | --- | --- |
| I-037 | Closed | Regression test and independent boundary probe both distinguish unsupported integer values from malformed values and verify the matching repair text. |

## Scope Lock

1. Prove the I-037 requirement and TEST-01, TEST-02, TEST-05, TEST-06, and
   TEST-08: explicit integer schema versions other than 2 are unsupported;
   syntactically non-integer versions are malformed; each diagnostic carries
   its corresponding repair.
2. Review the changed production and regression-test files and the public
   `doctor.RunV2Checks` entry point. CLI preflight output is observed but is
   not the named-diagnostic surface fixed by I-037.
3. The relied-on behavior is `data.ReadSchemaVersion` returning errors that
   preserve `data.ErrUnsupportedSchemaVersion` for `errors.Is`, followed by
   `V2ProblemRepair` mapping the selected name.
4. Matrix axes are input class (supported, unsupported integer, malformed),
   boundary/adjacent values, diagnostic name, and repair. Filesystem writes,
   state transitions, external services, and side effects are not applicable:
   this path only reads configuration and constructs a report.
5. A blocking result requires a supported `RunV2Checks` call within those
   cells to misclassify the input or select the wrong repair.

## Admission Ledger

| Re-check item | Prior claim | Frozen cell | Allowed result |
| --- | --- | --- | --- |
| Explicit 99 and 1 | I-037 repair | Unsupported integer / diagnostic and repair | Close or remain open |
| Text `nope` | I-037 repair | Malformed / diagnostic and repair | Close or remain open |
| Integers 0 and 3 | Adjacent unsupported values | Unsupported integer / diagnostic and repair | Close or remain open |
| Quoted `"2"` | Adjacent wrong-type value | Malformed / diagnostic and repair | Close or remain open |

## Coverage Matrix

| Input | Class | Expected diagnostic | Expected repair | Result |
| --- | --- | --- | --- | --- |
| `2` | Supported | No schema-version problem | N/A | Proven by existing runtime behavior; unchanged branch |
| `99` | Unsupported integer | `schema-version-unsupported` | Unsupported repair | Proven by regression test |
| `1` | Unsupported legacy integer | `schema-version-unsupported` | Unsupported repair | Proven by regression test |
| `0` | Lower adjacent integer | `schema-version-unsupported` | Unsupported repair | Proven by independent probe |
| `3` | Upper adjacent integer | `schema-version-unsupported` | Unsupported repair | Proven by independent probe |
| `nope` | Malformed scalar | `schema-version-malformed` | Malformed repair | Proven by regression test |
| `"2"` | Wrong YAML type | `schema-version-malformed` | Malformed repair | Proven by independent probe |

Empty and missing values are not additional I-037 blocking cells: the Issue
requires preserving the malformed class while correcting explicit unsupported
integers, and the existing missing-version behavior was not changed by the
repair. External-boundary and workflow side-effect matrices are not applicable.

## Evidence

- Source inspection: `RunV2Checks` defaults to
  `schema-version-malformed`, changes the name only when
  `errors.Is(err, data.ErrUnsupportedSchemaVersion)`, and selects repair text
  through `V2ProblemRepair(name)`.
- Regression command:
  `go test ./internal/doctor -run 'TestRunV2Checks_SchemaVersionNamesUnsupportedAndMalformed|TestI037IndependentBoundaryProbe' -count=1`
  passed. The second test was a checker-only temporary harness covering 0, 3,
  and quoted `"2"`; it was removed after the run.
- Independent CLI observation: isolated projects with explicit 99 and 1
  exited nonzero as unsupported, while `nope` exited nonzero as malformed.
  The CLI preflight intentionally reports the class before `RunV2Checks` and
  therefore does not render the internal bracketed diagnostic name or repair.
- Required Quick gate: `git diff --check && make build && make test-fast`
  passed on 2026-09-25T21:30Z with Go package tests green.
- File reality: all three reviewed paths exist; the temporary checker harness
  was intentionally deleted; `git status --short` showed only the router
  transition made for this Check before the immutable record was written.

## Acceptance Classification

| Requirement | Classification | Evidence |
| --- | --- | --- |
| Unsupported integer versions use the unsupported diagnostic and repair | Proven | 99 and 1 regression cells; 0 and 3 independent cells |
| Malformed versions retain the malformed diagnostic and repair | Proven | `nope` regression cell; quoted `"2"` independent cell |
| Required build and fast-test gate passes | Proven | `git diff --check`, `make build`, and `make test-fast` exit 0 |

## Materiality

No materiality actions are required. No Issues remain within the frozen scope.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches**
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth**
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs**

## Owner Validation

None required. This Check verifies the repair and closes I-037 technically;
it does not alter T-022 or O-021 lifecycle status.
