---
id: I-032
title: README identity examples and format placeholders were still bare
type: drift
status: resolved
source:
  kind: report
  actor: {role: executor, session: o018-t020-rename-20260923}
  at: '2026-09-23T08:12:14Z'
tasks: [T-020]
resolution:
  disposition: accepted
  actor: {role: owner, session: user}
  at: '2026-09-23T09:05:00Z'
  reason: >-
    Owner accepted the README repair without an independent recheck. C-910 found no bare IDs or X### placeholders left in README.md. This is owner acceptance, not a technical verified closure.
history:
  - at: '2026-09-23T08:12:14Z'
    actor: {role: executor, session: o018-t020-rename-20260923}
    kind: observed
    note: >-
      While renaming this repository's own V2 records under T-020, a
      repo-wide search for the bare identity pattern also matched
      README.md, which T-020's Context Files and O-018's Done When do not
      name (only router.md, AGENTS.md, Design.md, and Guardrails.md are in
      scope). README.md's `--objective O001` example, its
      `O001-example`/`R001-example` directory tree sample, and its
      `R###`/`O###` format placeholders predate the hyphenated identity
      requirement T-019 added and no longer match a schema that now rejects
      bare IDs.
  - at: '2026-09-23T08:12:14Z'
    actor: {role: executor, session: o018-t020-rename-20260923}
    kind: repair_attempted
    note: >-
      Fixed directly per the default Out-Of-Scope Repair rule (small,
      self-contained doc correction, no planning needed): updated the
      `--objective` example to `O-001`, the directory-tree sample to
      `O-001-example`/`R-001-example`, and the `R###`/`O###`/`release: R###`
      format placeholders to their hyphenated form in README.md. No other
      file changed. Not independently verified by a Check.
---

# I-032: README identity examples and format placeholders were still bare

## Summary

README.md documents the V2 record identity format for end users with a
handful of examples and format placeholders (`--objective O001`,
`O001-example`, `R001-example`, `R###`, `O###`, `release: R###`). None of
these are this repository's own records, so they were outside T-020's
renamed scope and outside O-018's named active-prose files (`router.md`,
`AGENTS.md`, `Design.md`, `Guardrails.md`). But the identity format itself
changed under this same Objective: `internal/data` now rejects a bare
`[ROTCI][0-9]{3,}` identity, so a new user following the old README examples
literally would create a record their own project immediately fails to load.

## Evidence

Repo-wide search `grep -rnoE '\b[ROTCI][0-9]{3,}\b'` against `README.md`
before repair: `O001` at two call-sites and `R001` at one, plus bare
`R###`/`O###` format placeholders at four more locations (the command-line
table, the board feature list, and the migration paragraph). All were stale
relative to the hyphenated schema enforced by T-019
(`internal/data/identity_v2.go`).

## Proof Needed

A checker (or a later targeted search) confirms `README.md` now contains no
bare `[ROTCI][0-9]{3,}` identity or `R###`/`O###`/`T###`/`C###`/`I###`
placeholder, and that the repair changed no other file's meaning.
