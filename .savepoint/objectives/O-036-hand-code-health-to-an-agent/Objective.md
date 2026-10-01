---
id: O-036
title: Hand Code Health to an agent
status: planned
depends_on: [O-035]
release: R-007
priority: medium
---

# O-036: Hand Code Health to an agent

## Outcome

A builder can give an AI agent everything Code Health knows with one sentence: "Investigate `.savepoint/health/report.md` and fix what it says." Savepoint keeps that file up to date as a plain-text report of the newest saved snapshot: a short brief for the agent, then the five signals with the ones needing attention first, each with its number, aim, what it means, what to do next, whether it blocks sign-off, and every affected file. When a signal needs attention, the popover's next step tells the builder to ask their agent to investigate that file.

## Why

O-035 made Code Health readable at a glance, but the popover is built to fit a small frame: rows are shortened, only the selected signal is explained, and the list of affected files is not shown anywhere. Today it says "13 files; savepoint health check lists them", and `savepoint health check` does not list them, and nothing tells the builder what to do with a red signal. Copying text out of a terminal also picks up box characters and cut-off words. An agent needs the full picture in a form it can read, and the builder needs one thing to point it at.

## Success Conditions

- A report file exists at `.savepoint/health/report.md` and is rewritten whenever a snapshot is saved, by `savepoint health check` or by a refresh from the popover, so it always describes the newest snapshot.
- `savepoint health report` rewrites the file from the newest saved snapshot without running any tool and without creating a snapshot. With no snapshot, or no Code Health set up, it says so plainly and writes nothing.
- The report opens with a short fixed brief to the agent: investigate each signal that is not Good, blocking ones first, propose then apply fixes, and re-run the check. It names the Objective's re-run command.
- All five signals appear, ordered Needs Attention, Watch, Unknown, Good. Each shows its name and question, number, aim, label, the same meaning and next step the popover shows, the same sign-off wording, and the better/worse/steady word.
- Every affected file is listed under its signal as `path:line` with the evidence note, in the stored order; the report never calls a file the worst unless the data says so. Several instances of one signal each get their own section.
- Provider names, snapshot hashes and raw timestamps stay out of the report, as in the popover. The measured date and whether the snapshot is official or manual are stated.
- For a signal that needs attention (red), the popover's `Next:` line says "Ask your agent to investigate .savepoint/health/report.md". When the report does not exist yet, it says to run `savepoint health report` first. Other labels keep their existing next steps.
- The popover's `Where:` line keeps naming the one affected file, and for several says the count only ("13 files"). The wrong claim that `savepoint health check` lists them is removed.
- The report is derived, so it is not committed: a `.gitignore` inside `.savepoint/health/` ignores `report.md` in new and existing projects, without editing the project's own `.gitignore`.
- Writing it never blocks or fails a check: a write error is reported as a warning and the snapshot and verdict stand. Writes are atomic, so an agent never reads half a file.
- Tests cover the report for each signal and label, ordering, one and many files, several instances, manual newest snapshot, no snapshot, not set up, a failed write, the rewrite after check and after refresh, the red and report-missing `Next:` wording, the `Where:` wording for one and many files, and the ignore rule.

## Architectural Considerations

The report is rendered in `codehealth` from the same dashboard model the popover reads, so the two never disagree on wording; nothing in classification, thresholds, snapshots or storage identity changes. The board stays a renderer of plain strings and reads no report; it only shows the path. Whether the report exists is read once when the dashboard loads, never while rendering. Evidence already carries path, line and note per result, so no new data is collected. The file is a derived artifact: it is never read back as input, and deleting it loses nothing. Rendering stays free of file and subprocess work; only the explicit commands write, using the same atomic temp-and-rename approach the snapshot store uses. Following the existing rule that agents run no other `savepoint` command, `report` is documented as an owner command alongside `health setup` and `health check`.

## Boundaries

**In scope:** the report renderer, writing it after each saved snapshot, the `savepoint health report` command, the red `Next:` wording and `Where:` fix in the popover, the Git ignore rule, and documenting the command in the CLI help and `Design.md`.

**Out of scope:** clipboard copy, a popover key to export, history or trend tables in the report, AI summaries, changing classification, thresholds or what is measured, editing the report from Savepoint, one report per Objective, a web view, and running any tool to produce the report.

## Confirmed Design Decisions

The owner confirmed on 2026-10-02:

1. **Delivery.** One derived file the owner points an agent at, not a clipboard feature or a copy key.
2. **Content.** All five signals, problems first, with a short fixed agent brief at the top. Newest snapshot only; history stays in the popover.
3. **Link from the popover.** When a signal is red, its `Next:` line tells the builder to ask their agent to investigate the report file; no new line is added, and `Where:` only gets its false claim removed.
4. **Not committed.** The report is derived and changes with every check, so a `.gitignore` beside it ignores it.
5. **Fallback.** A snapshot saved before this exists has no report until `savepoint health report` runs; a red signal's `Next:` says to run it first.
