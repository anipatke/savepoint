---
id: E45-safe-migration/T004-carry-active-work-across
title: Carry active work across
status: planned
objective: Render V2 Objective and Task records from active V1 epics and tasks without paraphrasing prose or fabricating evidence.
depends_on:
    - E45-safe-migration/T003-plan-the-conversion-before-touching-anything
complexity_tier: high
complexity_reason: Cross-schema rendering of two record families with lifecycle mapping and strict no-fabrication rules.
---

# T004: Carry active work across

## Problem

This is the point where V1 planning becomes V2 records, and where a well-meaning converter would quietly lie. Three temptations have to be designed out rather than resisted.

The first is healing. V1's loader turns an unrecognized status into `planned` so the board keeps working; a one-time conversion that does the same silently decides what the user meant about work they have not finished.

The second is fabricated evidence. E43 makes absent freshness `unknown`, which blocks completion — the correct outcome for work whose V1 record never carried a Check. Writing a `last_check` or a `freshness` block to make a migrated project look tidy would defeat every gate the last three epics built.

The third is paraphrase. Authored plans, acceptance criteria, and evidence notes are the user's words. They move under V2 headings; they do not get rewritten, reflowed, or summarized.

## Context Files

- `internal/migrate/convert.go`
- `internal/migrate/convert_test.go`
- `internal/migrate/plan.go`
- `internal/data/objective_v2.go`
- `internal/data/task_v2.go`
- `internal/data/evidence_v2.go`
- `internal/data/lifecycle.go`
- `internal/data/parser.go`
- `internal/data/release_doc.go`
- `internal/data/testdata/migration/v1-basic/manifest.yml`
- `internal/data/testdata/migration/v1-history/manifest.yml`
- `.savepoint/releases/v2/v2-Design.md`

## Acceptance Criteria

- [ ] Rendered Objective content decodes through `DecodeObjectiveV2` and rendered Task content through `DecodeTaskV2` with no diagnostic, for every active record in both fixtures.
- [ ] An active V1 epic becomes one Objective whose `title` comes from the epic title, whose `status` is `planned` or `in_progress`, whose `depends_on` maps declared epic dependencies to allocated `O###` values, and whose optional `release` preserves the source release as filtering metadata.
- [ ] An epic whose Tasks are all done but which was never audited converts to an `in_progress` Objective with no clearance, so it must still earn an integration Check rather than arriving complete.
- [ ] An active V1 Task becomes a Task owned by exactly one `O###`, with `depends_on` entries naming only converted active Tasks, and with the `{task, requires}` shape E43 decodes.
- [ ] No converted record carries `last_check`, `freshness`, `owner_validation.accepted_check`, or an exception unless the V1 source recorded the corresponding fact; a test asserts a converted fixture Task has unknown clearance and is therefore not completable.
- [ ] A `status: in_progress` Task carries a stage mapped from the source's raw `stage` or legacy `phase`, applying the recorded `implementation` → `build` alias and nothing beyond it.
- [ ] An unrecognized raw status or stage on an active record produces a blocking ambiguity rather than a healed value, and a test proves the V1 loader's healing is not reachable from this path.
- [ ] Authored body content is relocated verbatim under V2 headings, byte for byte including the fixtures' CRLF line endings and multiline authored text, with the V1 section each block came from named in the rendered output.
- [ ] Unknown frontmatter fields present on a V1 source are preserved on the converted record or recorded in the manifest; neither silently disappears.
- [ ] A converted Task whose V1 dependency was an archived completed Task carries the authored legacy-prerequisite line naming the archive path and the original completion evidence, and carries no dependency on a nonexistent Task.
- [ ] Rendering is deterministic: the same source and injected clock produce byte-identical output across runs.

## Implementation Plan

- [ ] Add `convert.go` with the Objective and Task renderers, taking planned targets from T003 and returning file content rather than writing it.
- [ ] Implement frontmatter construction against the `ObjectiveV2` and `TaskV2` field contracts, deliberately omitting every evidence field the source did not record.
- [ ] Implement the lifecycle mapping table — including the unaudited-epic rule and the legacy stage alias — and route every unrecognized value to a blocking ambiguity.
- [ ] Implement verbatim body relocation with source-section attribution, preserving original line endings rather than normalizing them.
- [ ] Render the authored legacy-prerequisite lines from the manifest entries T003 produces.
- [ ] Add round-trip tests that decode every rendered record through the real V2 decoders and assert no diagnostic.
- [ ] Test the no-fabrication rules, the unaudited-epic outcome, unrecognized lifecycle values, CRLF and unknown-field preservation, and byte-identical repeat rendering.
- [ ] Run the focused `internal/migrate` and `internal/data` suites, then `make build && make test`.

## Context Log

Pending.
