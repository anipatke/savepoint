---
id: I-065
title: Migrate parses the V1 findings README as a finding and aborts
type: defect
status: resolved
source:
  kind: report
  actor: {role: owner, session: user}
  at: '2026-09-25T22:43:41Z'
severity: blocker
history:
  - at: '2026-09-25T22:43:41Z'
    actor: {role: owner, session: user}
    kind: observed
    note: >-
      savepoint migrate --dry-run on an existing V1 project (galaxy) failed
      with "parse finding .savepoint/audit/findings/README.md: no
      frontmatter found".
  - at: '2026-09-25T22:46:10Z'
    actor: {role: executor, session: i065-repair-20260926}
    kind: repair_attempted
    note: >-
      internal/migrate/classify.go: Classify returns RoleUnclassified for
      .savepoint/audit/findings/README.md before the finding rule, so the
      unclassified sweep archives it byte-for-byte with an advisory note and
      planFindings never parses it. classify_test adds the path to the
      unclassified cases; TestPlan_v1ShippedAuditReadmesAreArchivedNotParsed
      copies v1-history, adds both shipped READMEs, and asserts both are
      archived, neither is a target, and F001 still converts. It fails with
      the fix reverted, with the reported error. git diff --check, make
      build, make test-fast, and make test-full passed. The owner's galaxy
      project now previews fully (migrate exits 1 only for four owner
      lifecycle decisions). Issue remains open for independent verification.
  - at: "2026-09-25T22:52:38Z"
    actor:
      role: owner
      session: board-owner
    kind: owner_decision
    note: Moved from open to in_progress by the owner from the board.
  - at: '2026-09-25T23:03:15Z'
    actor: {role: owner, session: user}
    kind: owner_decision
    note: >-
      Owner instructed the executor to mark this Issue resolved once its fix
      was deployed. The fix shipped in savepoint 2.0.2 on npm (tag v2.0.2).
      No technical CLEAR is implied.
resolution:
  disposition: accepted
  actor: {role: owner, session: user}
  at: '2026-09-25T23:03:15Z'
  reason: Owner accepted the fix as deployed in savepoint 2.0.2.
---

# I-065: Migrate parses the V1 findings README as a finding and aborts

## Summary

V1 `savepoint init` shipped `templates/project/.savepoint/audit/findings/README.md`
into every project. V2 migration classifies every
`.savepoint/audit/findings/*.md` as `RoleFinding` and `planFindings` parses
each one, returning an error on the README's missing frontmatter. The whole
migration preview and apply abort, so every V1 project that kept the shipped
README cannot migrate with 2.0.0 or 2.0.1.

`audit/runs/README.md` is also shipped but is safe: audit runs are archived
without parsing.

## Evidence

- `git ls-tree v1.3.1` lists `templates/project/.savepoint/audit/findings/README.md`
  and `templates/project/.savepoint/audit/runs/README.md`.
- `internal/migrate/classify.go:58` matches `^\.savepoint/audit/findings/[^/]+\.md$`.
- `internal/migrate/plan.go:1243` returns `parse finding %s` on any parse error.
- The fixture projects `v1-basic` and `v1-history` carry no findings README,
  so no test exercised it.

## Proof Needed

- A V1 project with the shipped `audit/findings/README.md` and
  `audit/runs/README.md` plans and applies; both READMEs are archived
  byte-for-byte and no finding or Issue is planned from them.
- Real findings still convert as before.
- Tests cover classification and a plan over a fixture copy with both
  READMEs; `make build && make test-fast` pass; migration-sensitive handoff
  also runs `make test-full`.
