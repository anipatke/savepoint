---
id: O-035
title: Make Code Health glanceable
status: done
depends_on: [O-031]
release: R-007
priority: medium
---

# O-035: Make Code Health glanceable

## Outcome

Code Health becomes a compact part of the board that a non-expert can read at a glance. A header chip always shows overall health. `H` opens a small popover, not a full screen, that lists the five signals in plain words, each with its number beside a yardstick, a sparkline of recent official checks, and one line saying what it means, whether it blocks sign-off, and what to do next. `h` inside the popover shows the last ten checks.

## Why

O-031's full-screen view shows numbers and verdicts with no scale ("complexity 46", "duplication 8%"), uses provider jargon, repeats itself, and never says whether a finding blocks anything. A builder who is not an engineer cannot tell good from bad. The screen also disagrees with `savepoint health check`, which calls the same data advisory.

## Success Conditions

- The board header shows a Health chip, coloured by the worst signal, that costs no extra row and reads from saved data only.
- `H` opens a fixed-height popover that fits an 80x20 terminal and never scrolls. Esc closes it and restores the board cursor.
- Each signal has a plain-language name and question, for example "Complexity: how tangled is the hardest code?", kept in one place.
- Each row shows the measured number, a good/watch/needs-attention mark, and an `aim` yardstick taken from the signal's Good threshold, using the project's configured thresholds or else Code Health's built-in defaults, so every signal has one.
- The selected row's one-line explanation says what the result means, whether it blocks sign-off (derived from the same gate `savepoint health check` uses, so the two never disagree), and one next step.
- A sparkline of up to ten comparable official checks appears from three. It is oldest first, uses eight-level blocks, is flat when the movement is below the signal's existing material-movement size, is coloured by health, and carries `better`, `worse` or `steady`. Under five points it is marked early. Fewer than three shows "not enough history yet", and a restarted series says why.
- The popover shows the command to re-run the official check, ready to copy (`savepoint health check O-###`, naming the Objective in view), so nobody has to remember it.
- `h` shows the last ten checks, newest first, with date, official or manual, and overall label; manual ones are dimmed and never feed a sparkline.
- Provider names, snapshot hashes and raw timestamps do not appear in the default views.
- Refresh (`R`), progress, cancellation, staleness and not-configured/first-run behaviour from O-031 keep working inside the popover.
- Tests cover chip states, popover at 80x20/80x24/80x40, selection, every signal's copy, sparkline rules (2, 3, 5 and 10+ points, flat, restarted, partial), blocks-sign-off wording matching the gate, history, refresh and cancel.

## Architectural Considerations

The board still consumes a read-only dashboard model and sends explicit commands; rendering does no file or subprocess work. The dashboard row gains the plain-language copy, the `aim` text, the sparkline values (from the comparable official series Code Health already selects) and the sign-off wording; none of these change classification, thresholds or storage. The full-screen detail, its scrolling and the affected-files list are replaced, so I-104 stops applying once this lands. Chip data is one cheap read of the latest snapshot, like O-031's Check line.

## Boundaries

**In scope:** header chip, popover, plain-language copy for the five signals, `aim` yardsticks, sparklines, `h` history, one consistent blocks-sign-off statement, and removal of the full-screen detail view.

**Out of scope:** new signals, changes to the classifier or thresholds, editing thresholds in the TUI, browsing affected files, deleting snapshots, changing when official snapshots are created, AI summaries, web UI, background refresh.

## Confirmed Design Decisions

The owner confirmed on 2026-10-02:

1. **Sequencing.** O-031 finishes first: its Check is rechecked after the I-104 to I-106 repairs and the owner accepts it. This Objective builds on that screen and replaces it. O-032 keeps depending on both.
2. **Affected files.** The explanation line names the affected file only when there is exactly one ("Where: internal/doctor/repairs.go"). With several it gives the count and points to `savepoint health check`, which keeps the full list ("Where: 5 files; savepoint health check lists them"). Stored evidence is kept in path order, so the popover never claims a file is the worst. Amended 2026-10-02 from "names only the top affected file" after T-075 found the rank is not kept in saved snapshots.
3. **Verdict basis.** The good/watch/needs-attention label stays level-based, as classified today. Trend appears only as the sparkline. A "no worse than last time" verdict is a later change.
4. **Placement.** A header chip plus a popover on `H`. A one-line strip under NEXT was considered and not chosen.

Details settled while planning:

- The `aim` text comes from each signal's Good threshold (coverage "80 or more"; complexity, duplication, failing tests and vulnerabilities "N or less"). The explanation mentions the Watch boundary only when the signal is beyond it.
- The sign-off sentence comes from `Evaluate` on the newest official snapshot. When the newest snapshot is a manual refresh, the line says a manual refresh does not affect sign-off; with no official check it says so.
- The popover shows its own plain copy, not the stored classification explanation, so stored snapshots and the classifier are untouched. That stored text still says unknown-severity vulnerabilities are "treated as blocking" while the gate only reports them; the popover never repeats it.
- Sparkline text is produced in `codehealth`, so the board stays a renderer of plain strings.
