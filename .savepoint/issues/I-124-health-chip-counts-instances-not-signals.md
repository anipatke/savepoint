---
id: I-124
title: Header health chip counts configured instances, so it shows /6 instead of /5
type: defect
status: open
source:
  kind: report
  actor: {role: owner, session: deep-time-migration-2026-10-03}
  at: '2026-10-03T00:00:00Z'
tasks: [T-076, T-079]
history:
  - at: '2026-10-03T00:00:00Z'
    actor: {role: executor, session: user-request}
    kind: observed
    note: 'Owner saw "♥ Health 1/6" in the board header after migrating the deep-time project. Cause traced in source; not yet repaired. Related to resolved I-112, which fixed the same multi-instance counting for the popover only.'
  - at: '2026-10-03T00:00:00Z'
    actor: {role: executor, session: user-request}
    kind: repair_attempted
    note: 'Reproduced on the real deep-time data: LoadChip gave Good:1 Signals:6 before and Good:1 Signals:5 after. Dashboard.Chip now groups rows by capability and counts a signal Good only when its worst instance is Good. The worst-instance ordering moved to codehealth.WorseInstance and the popover uses it, so chip and popover share one rule. TestDashboardChipCountsSignalsNotInstances covers every-instance-good (5/5) and one-bad-instance (4/5). The old chip test built rows without a Capability, so it now sets them. Full go test ./... passes. Not run: the board rendered on deep-time, and the narrow-width chip tests only via the existing suite.'
  - at: '2026-10-03T00:00:00Z'
    actor: {role: owner, session: user-request}
    kind: reopened
    note: 'After 2.1.2, deep-time still showed "Official check · 4 of 6 need a look". The first repair fixed only the header chip; the dashboard headline (headlineText in dashboard_copy.go) still counted instances. With coverage now configured, deep-time has 6 rows for 5 signals.'
  - at: '2026-10-03T00:00:00Z'
    actor: {role: executor, session: user-request}
    kind: repair_attempted
    note: 'Added signalRows (one row per signal, the worst instance) and used it in both Dashboard.Chip and headlineText. Real deep-time data: headline "4 of 6 need a look" before, "3 of 5 need a look" after; chip stays 2/5. New TestHeadlineCountsSignalsNotInstances; the old headline table test now gives its rows distinct capabilities. Full go test ./... passes. Not run: the rendered board. Other row counts in internal/board were searched; the remaining len(rows) uses are Issues lists, not health.'
  - at: '2026-10-03T00:40:00Z'
    actor: {role: executor, session: user-request}
    kind: rechecked
    note: 'Rechecked on the released 2.1.4 in deep-time through a real pseudo-terminal board: header chip shows ''♥ Health 3/5'', and the Health popover (H) shows ''3 Oct 10:25 · Official check · 2 of 5 need a look'' with five signal rows, the tests row labelled ''Tests failing ×2'' and ''worst of 2''. Chip and headline agree (3 Good, 2 need a look, 5 signals). Before the repair the same project showed /6 and ''4 of 6''. Executor evidence only; it does not claim verified. The Issue stays open for a checker or an explicit owner decision.'
---
# I-124: Header health chip counts configured instances, so it shows /6 instead of /5

## Summary

The board header chip shows `♥ Health n/6` for a project that configures one signal more than once. Savepoint defines five signals and T-076 specifies the chip as `n/5`. The chip counts dashboard rows, and rows are one per configured instance, so each extra instance inflates both the denominator and the numerator. The configuration is valid; the defect is in Savepoint.

Related: I-112 (resolved, verified by C-954) made the popover show one row per signal standing for its worst instance. The chip was not brought in line with that rule. Reopening I-112 was not chosen because the symptom, location and requirement differ (header chip, T-076, versus popover, T-079).

## Evidence

Observed on the deep-time project, whose `.savepoint/health/config.json` has two `tests` instances (`pytest-junit`, `vitest-junit`): 2 tests + 1 coverage placeholder + complexity + duplication + dependency health = 6 rows, shown as `1/6`. Not reproduced by running the board; inferred from source and that config.

`Dashboard.Chip()` (`internal/codehealth/chip.go:50`) sets `Signals: len(d.Rows)` and counts Good per row. `dashboardRows` (`internal/codehealth/dashboard.go:335`) emits one row per configured instance, or one placeholder for a capability with none. The per-signal grouping exists only in the board layer: `HealthOverlay.rows()` and `worseInstance` (`internal/board/v2/health.go:134`), so `codehealth` has no shared grouping to reuse.

## Proof Needed

The chip denominator is the number of signals (5) however many instances are configured. A signal counts as Good only when its worst instance is Good, using the popover's ordering (blocking first, then label severity) so chip and popover cannot disagree. A test with two instances of one signal shows `n/5`, including a case where one instance is Good and the other is not. The existing chip tests keep passing, and the `♥ 3/5` short form still shows at narrow widths. Decide whether to move the worst-instance grouping into `codehealth` so both screens share it.
