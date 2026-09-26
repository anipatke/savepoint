---
id: I-068
title: The migrate preview is too long and noisy to review
type: other
status: open
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
