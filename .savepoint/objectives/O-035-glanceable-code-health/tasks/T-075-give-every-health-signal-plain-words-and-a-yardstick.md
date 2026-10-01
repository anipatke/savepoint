---
id: T-075
title: Give every health signal plain words and a yardstick
objective: O-035
status: done
depends_on: []
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: planning-o035-20261002}
complexity_tier: medium
complexity_reason: Pure data and copy in one new file plus the row builder, but five signals times four labels of wording and value formatting need careful tests.
check_waiver:
    task: T-075
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-01T22:06:52Z"
---

# Give every health signal plain words and a yardstick

## Outcome

Each dashboard row carries plain-language wording a non-expert can read: the question the signal answers, its value in words, the `aim` yardstick, what the result means, one next step, and where the problem is (the file when there is only one, otherwise a count). The board is not changed.

## User Check

None beyond the Full Objective Check; the wording is first seen in the popover.

## Done When

- `DashboardRow` gains `Question`, `Value`, `Aim`, `Meaning`, `NextStep`, and `Where`, all plain strings; the board still receives no snapshot or provider type.
- `Question` per signal: Tests "Do the tests pass?", Coverage "How much code do the tests actually run?", Complexity "How tangled is the hardest code?", Duplication "How much is copy-pasted?", Dependency vulnerabilities "Known security problems in libraries we use".
- `Value` reads as words from the result: for example "all 3,045 pass" or "2 failing", "86%", "hardest function scores 46", "8% copy-pasted", "2 unrated" or "none". An unmeasured row says "not measured".
- `Aim` comes from the signal's Good threshold, using the instance's configured thresholds or else the built-in defaults, so every measured signal has one: coverage "aim for 80% or more"; complexity, duplication, failing tests and vulnerabilities "aim for N or less" (none for zero).
- `Meaning` and `NextStep` come from one table keyed by signal and label (Good, Watch, Needs Attention, Unknown). The Watch boundary appears only when the value is beyond it. `NextStep` is a short instruction a builder can hand to their agent, for example "Ask your agent to split the most tangled function".
- `Where` names the file when the evidence covers exactly one distinct path ("internal/doctor/repairs.go"); with two or more it says the count and points to the full list ("5 files; savepoint health check lists them"); empty when there is no evidence. It never calls a file the top or worst, because stored evidence is in path order (O-035 Decision 2).
- All wording lives in data tables in one file (STYLE-09). The stored classification explanation, thresholds and snapshots are unchanged, and the new wording never says a signal is "treated as blocking".
- Tests cover every signal at each label, a configured threshold overriding the default, unmeasured and not-configured rows, no evidence, one evidence path (including several entries in the same file), and many evidence paths.

## Context Files

`.savepoint/objectives/O-035-glanceable-code-health/Objective.md`; `internal/codehealth/dashboard.go`, `internal/codehealth/dashboard_test.go`, `internal/codehealth/dashboard_copy.go` (new), `internal/codehealth/dashboard_copy_test.go` (new); `internal/codehealth/classification.go` (read only: thresholds); `internal/codehealth/history.go` (read only: number formatting); `internal/codehealth/reader_tests.go`, `internal/codehealth/reader_coverage.go`, `internal/codehealth/reader_lizard.go`, `internal/codehealth/reader_jscpd.go`, `internal/codehealth/reader_osv.go` (read only: detail keys).

## Design References

O-035 Confirmed Design Decisions (aim text, affected files, verdict basis); O-031 Architectural Considerations (read-only dashboard model).

## Guardrails

ARCH-04, DATA-03, STYLE-07, STYLE-09, TEST-01, TEST-02, TEST-04.

## Implementation Plan

1. Read `dashboardRows`, `measuredRow`, `thresholdFor` and the detail-key constants; return REPLAN REQUIRED if instance thresholds are not reachable from `measuredRow`.
2. Add the six fields to `DashboardRow` and the copy tables in `dashboard_copy.go`.
3. Fill them in `measuredRow`, `notConfiguredRow` and `noResultRow`.
4. Add table-driven tests for the cases in Done When.

## Boundaries

No board rendering, no sign-off wording, no sparkline, no change to classification text or stored snapshots.

## Technical Verification

Focused tests during iteration (`make test-focused`); `make build && make test-fast` for ordinary handoff. Later evaluated per `agent-skills/references/check-method.md`.

## Technical Evidence

REPLAN REQUIRED (partial work preserved, status/stage unchanged).

- Gap: `identity.go` `canonical()` sorts `Evidence` by path/line/note before storage, so `evidence[0]` of a loaded result is the alphabetically first path, not the worst. Confirmed on the real snapshot in `.savepoint/health/snapshots/`: complexity evidence begins `internal/board/v2/objectives_test.go` "has complexity 29" while 36 is the worst. A `Where` built from it would name an arbitrary file as the top affected file. Reader notes are free text, so recovering the rank by parsing them is not sound.
- Options for the planner: (a) keep rank in stored evidence (changes canonical form or adds a rank field, which conflicts with this Task's "snapshots unchanged" boundary and affects snapshot IDs); (b) drop the "top" claim and have `Where` name the first path plus count, worded neutrally; (c) derive rank per signal from the note where the reader format is fixed.
- Done so far: `internal/codehealth/dashboard_copy.go` (new: question, value, aim, meaning/next-step tables, Where helper) and the six `DashboardRow` fields filled in the three row builders; `dashboard_copy_test.go` (new). All new tests pass except `TestDashboardRowsCarryPlainWords`, whose Where expectation assumes rank order.
- Extra reads: none beyond Context Files, plus `internal/codehealth/identity.go`, `reader_paths.go`, `storage.go`, `snapshot.go` and the real snapshot file, to find the cause of the failing Where expectation.
- Not run: `make build && make test-fast` (handoff gate), because the plan is blocked.

Replan 2026-10-02 (planner, owner chose option d): `Where` names a file only when there is one distinct path, otherwise a count plus `savepoint health check`. No storage change. Resume from the preserved work: update `whereText`, its comment and the Where expectations, then run the handoff gate.

## Drift Notes

New module or architecture delta beyond the documented Codebase Map, reconciled through the planner before Check.

## Handoff Evidence (2026-10-02, after replan)

Gate: `make build && make test-fast` — exit 0, no failures.

Per criterion:
- Row fields: `DashboardRow` has `Question`, `Value`, `Aim`, `Meaning`, `NextStep`, `Where` (plain strings), filled in `measuredRow`, `notConfiguredRow`, `noResultRow`; board untouched. Tests: `TestDashboardRowsCarryPlainWords`, `TestDashboardInstanceMissingFromSnapshotHasWords`.
- Questions: `TestQuestionsAreFixedPerSignal`.
- Value in words, incl. "not measured": `TestValueReadsAsWords`, `TestUnmeasuredResultSaysNotMeasured`.
- Aim from defaults or configured Good: `TestAimComesFromDefaultsOrConfiguredGood`, plus the configured override in `TestDashboardRowsCarryPlainWords`.
- Meaning/NextStep for every signal and label; Watch line only beyond it: `TestMeaningAndNextStepExistForEverySignalAndLabel`, `TestWatchBoundaryAppearsOnlyBeyondIt`.
- Where per replan (one path / count + command / empty / same file at several lines): `TestWhereNamesOneFileOrCountsMany`.
- Wording in tables in one file; no "treated as blocking", no banned claims: `TestPlainWordingNeverRepeatsBlockingClaimsOrBannedClaims`. Classification text, thresholds, snapshots unchanged.
- Not-configured and no-result rows: `TestNotConfiguredRowSaysHowToSetItUp`.

Files changed: `internal/codehealth/dashboard.go`, `dashboard_copy.go` (new), `dashboard_copy_test.go` (new).
Extra reads (beyond Context Files): `identity.go`, `reader_paths.go`, `storage.go`, `snapshot.go`, a real snapshot (for the replan).
Limitations: the Meaning for Needs Attention can come from a worsening trend while the value is within Watch; the text stays neutral then. Where is not checked against a real `savepoint health check` run. No Task Check requested; owner decides on waiver or Check.
