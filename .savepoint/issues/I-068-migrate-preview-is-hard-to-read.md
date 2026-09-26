---
id: I-068
title: The migrate preview is too long and noisy to review
type: other
status: resolved
source:
  kind: report
  actor: {role: owner, session: user}
  at: '2026-09-26T00:15:19Z'
severity: medium
history:
  - at: '2026-09-26T00:15:19Z'
    actor: {role: owner, session: user}
    kind: observed
    note: >-
      Owner: "dry run output is not user friendly..." after previewing the
      galaxy migration.
  - at: '2026-09-26T00:22:13Z'
    actor: {role: executor, session: i068-repair-20260926}
    kind: repair_attempted
    note: >-
      internal/migrate/preview_summary.go adds FormatSummaryPreview, now the
      default: counts of Goals, Objectives, Tasks, Issues, and Documents; the
      archive count; advisory notes as one count; the router Goal; one line
      per Objective with its source epic, open task count, and any decided
      status; decisions still needed in full with choices and paste-ready
      decisions-file entries; decisions applied; and a ready or blocked
      status line. FormatPreview is unchanged and printed with the new
      --verbose flag (cmd/migrate.go, CommandOptions.Verbose, main.go). The
      manifest and apply are unchanged. README's migration section describes
      the summary, decisions file, and --verbose. The galaxy preview drops
      from 1,409 lines to 36. Tests: preview_summary_test.go (collapsed ready
      plan, determinism, blocked plan with entries, decided plan),
      TestMainMigrateVerboseListsEveryPlannedRecord,
      TestParseMigrateArgs_verbose, and the default preview test updated to
      expect the summary. make build, make test-fast, and make test-full
      passed. Issue remains open for verification or owner acceptance.
  - at: '2026-09-26T03:34:53Z'
    actor: {role: owner, session: user}
    kind: owner_decision
    note: >-
      Owner instructed the executor to mark this Issue resolved. The fix
      shipped in savepoint 2.0.5 on npm (tag v2.0.5), and in real use the galaxy preview printed a 36-line summary.
      No technical CLEAR is implied.
resolution:
  disposition: accepted
  actor: {role: owner, session: user}
  at: '2026-09-26T03:34:53Z'
  reason: Owner accepted the fix as deployed in savepoint 2.0.5.
---

# I-068: The migrate preview is too long and noisy to review

## Summary

For the galaxy project the preview is 1,409 lines. It lists every planned
record, then 746 archive lines, then 587 ambiguity lines, 583 of which are
near-identical advisory `unclassified_file` notes. The owner cannot see at a
glance what the migration will create, what it will archive, and what still
needs a decision.

## Evidence

- galaxy `savepoint migrate --decisions decisions.yml`: 1,409 lines; sections
  Planned records (about 50), Documents, Archives (746), Ambiguities (587).

## Proof Needed

- The preview opens with a short summary: Goals, Objectives, Tasks, and
  Issues to create; files to archive; decisions needed and resolved; and
  whether apply is blocked.
- Decisions still needed are listed in full near the top, with their choices
  and the exact decisions-file line to add.
- Archive and advisory sections are collapsed to counts by default, with a
  `--verbose` flag that prints the full listing.
- The written manifest and apply behaviour are unchanged.
- Tests cover the summary, the collapsed and verbose forms, and determinism;
  `make build && make test-fast` pass.
