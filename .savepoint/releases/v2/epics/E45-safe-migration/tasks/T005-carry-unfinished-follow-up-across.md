---
id: E45-safe-migration/T005-carry-unfinished-follow-up-across
title: Carry unfinished follow-up across
status: planned
objective: Convert unresolved V1 defects and audit findings into V2 Issues by disposition, archiving the rest without implying a Check.
depends_on:
    - E45-safe-migration/T003-plan-the-conversion-before-touching-anything
complexity_tier: high
complexity_reason: Maps ten finding states and defect lifecycle onto E44 Issue obligations without fabricating proof.
---

# T005: Carry unfinished follow-up across

## Problem

V1 spreads follow-up across three shapes: release defects, ten audit finding states, and a separately maintained register. V2 has one Issue record, and E44 gave it obligations that a careless conversion will violate immediately.

The specific traps: a `fixed` finding is not a verified one — it still awaits independent proof, and converting it to `resolved` would launder an unverified repair into a closed record. A `verified` finding's proof was a V1 audit, not a V2 Check, so it must archive rather than arrive with a fabricated `resolution.check`. A `duplicate` finding can only resolve as `duplicate` if its canonical counterpart also becomes an Issue, because E44 requires `duplicate_of` to name an existing record. And a `deferred` or `owner_decision` finding is still open — E44 made deferral a dated history entry precisely so it could not become a fourth lifecycle state.

Dates are the quiet one. An Issue's seeded history entry needs a timestamp, and inventing one would put a fictional date into an append-only record.

## Context Files

- `internal/migrate/convert_issues.go`
- `internal/migrate/convert_issues_test.go`
- `internal/migrate/convert.go`
- `internal/migrate/plan.go`
- `internal/data/issue_v2.go`
- `internal/data/defect.go`
- `internal/data/audit_finding.go`
- `internal/data/audit_run.go`
- `internal/data/audit_register.go`
- `internal/data/testdata/migration/v1-history/manifest.yml`
- `.savepoint/releases/v2/v2-Design.md`

## Acceptance Criteria

- [ ] Rendered Issue content decodes through `DecodeIssueV2` with no diagnostic for every converted record in the `v1-history` fixture.
- [ ] An unresolved V1 defect becomes an Issue with `type: defect` and status `open` or `in_progress` mapped from its V1 state; a resolved defect is archived and produces no Issue.
- [ ] Findings in `open`, `triaged`, or `mapped` become open Issues preserving their reasons, proof requirements, and relations.
- [ ] A `deferred` finding becomes an open Issue carrying a dated `deferred` history entry and no resolution; an `owner_decision` finding that resolved nothing becomes an open Issue carrying a dated `owner_decision` entry.
- [ ] An `in_progress` or `fixed` finding becomes an `in_progress` Issue with a `repair_attempted` history entry; neither is ever converted to `resolved`, and a test proves `fixed` still awaits verification.
- [ ] A `verified` finding is archived with its original proof and status and produces no Issue and no V2 Check.
- [ ] A `waived` finding is archived with its disposition; when active converted work references it, a typed historical decision entry is recorded in the manifest rather than an Issue being invented to hold it.
- [ ] A `duplicate` finding resolves as `duplicate` with `duplicate_of` naming the converted canonical Issue only when that canonical also converts; when the canonical is archive-only, the duplicate is archive-only too, and a test covers both directions using the fixture's `F003-duplicate`.
- [ ] Every converted Issue's `source` records `kind: migration` with its source-qualified legacy key, and never borrows a Check's authority.
- [ ] A seeded history entry uses the date the V1 record recorded; when the source recorded none, it uses the migration time and carries an explicit note saying the original date was absent.
- [ ] Converted Issues link to converted repair Tasks through `tasks` when the V1 record named one, and to nothing when it did not; no link is inferred from narrative similarity.
- [ ] The V1 audit register, runs, and prompt are archived intact and no register is recreated; Issue listings in V2 are derived, as E44 established.
- [ ] A finding whose narrative names no stable identity and cannot be deterministically mapped produces a blocking ambiguity rather than a guessed Issue.

## Implementation Plan

- [ ] Add `convert_issues.go` with the disposition mapping table expressed as data, not as a chain of conditionals.
- [ ] Implement defect conversion over the existing `internal/data/defect.go` reader, and finding conversion over `audit_finding.go`, treating both as read-only migration inputs.
- [ ] Implement the duplicate rule with its dependency on whether the canonical converts, resolving it after allocation so both outcomes are decidable.
- [ ] Implement `source` construction and the dated-history seeding rule, including the absent-date note.
- [ ] Route unmappable narrative findings to blocking ambiguities rather than to Issues.
- [ ] Add round-trip tests decoding every rendered Issue through `DecodeIssueV2`.
- [ ] Test each of the ten finding states, both defect states, both duplicate directions, the archived-with-reference waiver path, `source` provenance, and both history-date paths.
- [ ] Run the focused `internal/migrate` and `internal/data` suites, then `make build && make test`.

## Context Log

Pending.
