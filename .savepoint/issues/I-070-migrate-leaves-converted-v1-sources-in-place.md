---
id: I-070
title: Migrate leaves the V1 sources of converted Tasks, defects, and findings in place
type: defect
status: open
source:
  kind: report
  actor: {role: owner, session: user}
  at: '2026-09-26T03:34:20Z'
severity: high
history:
  - at: '2026-09-26T03:34:20Z'
    actor: {role: owner, session: user}
    kind: observed
    note: >-
      After migrating galaxy with 2.0.5, 31 V1 source files (21 findings, 3
      defects, 7 tasks) that had been converted to V2 Issues and Tasks were
      still at their V1 paths and had no copy under .savepoint/archive/v1/.
      The executor moved them into the archive by hand with git mv.
---

# I-070: Migrate leaves the V1 sources of converted Tasks, defects, and findings in place

## Summary

`buildApplyBatch` archives and removes only `plan.Archives`. Release PRDs
and epic details that convert are also archived, but a V1 task, defect, or
finding that converts to a V2 Task or Issue is neither archived nor removed.
The project keeps a live V1 tree beside the V2 records, contrary to Design:
"Former V1 project records remain byte-preserved under
`.savepoint/archive/v1/`". The manifest records their hashes, so nothing is
lost, but the leftovers are duplicates an agent may read or edit.

## Evidence

- `internal/migrate/apply.go` `buildApplyBatch`: copy and remove are built
  only from `plan.Archives`.
- galaxy after apply: 31 converted sources remained under
  `.savepoint/audit/findings/`, `.savepoint/releases/v1/defects/`, and
  `.savepoint/releases/v1/epics/E10-natural-rendering-exploration/tasks/`;
  none had an archive copy.

## Proof Needed

- Every V1 source consumed by a converted target is archived byte-for-byte
  under `.savepoint/archive/v1/` and removed from its V1 path, and the preview
  counts it among archived files.
- After apply no V1 record remains outside the archive.
- Recovery and the manifest still map each converted record to its archived
  source.
- End-to-end tests over the fixtures assert no V1 file remains at its
  source path; `make build && make test-fast` and `make test-full` pass.
