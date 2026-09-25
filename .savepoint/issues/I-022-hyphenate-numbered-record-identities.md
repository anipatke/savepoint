---
id: I-022
title: Hyphenate every numbered V2 record identity
type: other
status: resolved
source:
  kind: report
  actor: {role: owner, session: user}
  at: '2026-09-22T09:48:33Z'
severity: high
resolution:
  disposition: accepted
  actor: {role: owner, session: user}
  at: '2026-09-22T00:05:00Z'
  reason: >-
    Superseded by O-018 (hyphenate-numbered-record-identities), which now owns
    the design and repair. Owner override: closing this Issue as redundant
    rather than leaving it open in parallel with the Objective that carries
    it.
history:
  - at: '2026-09-22T09:48:33Z'
    actor: {role: owner, session: user}
    kind: observed
    note: Numbered identities with the digits run straight onto the kind letter are harder to scan than O-015 and T-012; adopt the hyphenated form for every active numbered record type and migrate this project.
  - at: '2026-09-22T00:00:00Z'
    actor: {role: planner, session: o018-design-20260922}
    kind: deferred
    note: >-
      Repair spans the identity grammar and every generator/consumer in
      internal/data, a new interruption-safe V2->V2 migration tool, fresh-init
      templates, skills, and this repository's entire active V2 graph — too
      broad for an inline repair. Carried into O-018
      (hyphenate-numbered-record-identities), which plans the shared grammar,
      the migration tool, and the repository migration itself (including
      renaming this Issue to I-022 as the migration's own final step). This
      Issue stays open until O-018's Full Objective Check verifies the repair.
  - at: '2026-09-22T00:05:00Z'
    actor: {role: owner, session: user}
    kind: owner_decision
    note: >-
      Owner override: close this Issue now (status resolved, disposition
      accepted) rather than leaving it open in parallel with O-018, which
      already owns the design and repair. Not a technical closure — O-018's
      own Full Objective Check still proves the repair.
---

# I-022: Hyphenate every numbered V2 record identity

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

This bootstrap Issue is recorded as `I-022` because that is the only form the
current schema accepts. The migration must rename it to `I-022` along with the
rest of the active Issue set.

## Evidence

- `internal/data/release_v2.go`, `objective_v2.go`, `task_v2.go`,
  `check_v2.go`, and `issue_v2.go` currently validate unhyphenated identities
  and report requirements such as "T plus at least three digits."
- Evidence and relationship fields repeat those formats for `release`,
  `objective`, `depends_on`, `last_check`, freshness, owner acceptance,
  exceptions, supersession, Issue links, and duplicate links.
- Active project paths and records currently use forms such as `O-012`, `T-009`,
  `C-906`, and `I-022`.
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
