---
id: I022
title: Hyphenate every numbered V2 record identity
type: other
status: open
source:
  kind: report
  actor: {role: owner, session: user}
  at: '2026-09-22T09:48:33Z'
severity: high
history:
  - at: '2026-09-22T09:48:33Z'
    actor: {role: owner, session: user}
    kind: observed
    note: Numbered identities such as O015 and T012 are harder to scan than O-015 and T-012; adopt the hyphenated form for every active numbered record type and migrate this project.
---

# I022: Hyphenate every numbered V2 record identity

## Summary

All active V2 numbered identities should separate the record kind from its
number with a hyphen. The canonical forms become `R-###`, `O-###`, `T-###`,
`C-###`, and `I-###` for Releases, Objectives, Tasks, Checks, and Issues.

This is a schema-wide identity migration, not a display-only alias. It must
update parsing, validation, references, filenames/directories, generated
templates, board/resume/doctor output, skills and guidance, tests and fixtures,
and every active record in this repository. The migration must not rewrite
immutable Check contents casually or leave mixed identities that the strict V2
index cannot load.

This bootstrap Issue is recorded as `I022` because that is the only form the
current schema accepts. The migration must rename it to `I-022` along with the
rest of the active Issue set.

## Evidence

- `internal/data/release_v2.go`, `objective_v2.go`, `task_v2.go`,
  `check_v2.go`, and `issue_v2.go` currently validate unhyphenated identities
  and report requirements such as "T plus at least three digits."
- Evidence and relationship fields repeat those formats for `release`,
  `objective`, `depends_on`, `last_check`, freshness, owner acceptance,
  exceptions, supersession, Issue links, and duplicate links.
- Active project paths and records currently use forms such as `O012`, `T009`,
  `C906`, and `I022`.
- The owner requested the hyphenated format for readability and explicitly
  required this repository to be scanned and migrated as part of the work.

## Proof Needed

- Define one canonical typed identity grammar for all five active V2 record
  kinds: `R-`/`O-`/`T-`/`C-`/`I-` plus at least three digits, with shared
  formatting and validation rather than divergent regex copies.
- Inventory every identity-bearing field, filename, directory, map key,
  selection/router field, renderer, CLI filter, diagnostic, template, skill,
  fixture, test, and migration path before changing data.
- Provide a previewable, interruption-safe migration for existing projects;
  reject mixed or colliding results without partial renames or broken
  references. Define whether old unhyphenated input is accepted temporarily
  and how canonical writes behave during compatibility.
- Migrate this repository's complete active V2 graph atomically: Releases,
  Objectives, Tasks, Checks, Issues, router selections, dependencies,
  freshness/acceptance/exception links, supersession, Issue history links,
  filenames, directories, and active prose references. Preserve archived V1
  history unless an explicit migration boundary says otherwise.
- Update fresh-init and upgrade assets so newly generated projects use only
  hyphenated identities; keep canonical/scaffold skill copies byte-identical.
- Prove strict loading, board interaction, resume, doctor, migration recovery,
  duplicate detection, non-TTY output, and all relationship/gate resolvers
  across the new format. Include negative tests for malformed, mixed,
  duplicate, and partially migrated identities.
- Run repository-wide targeted searches showing no unintended active
  unhyphenated identities remain, then pass `git diff --check`, `make build`,
  `make test`, and the mandatory Full Objective/Release Checks governing the
  migration.
