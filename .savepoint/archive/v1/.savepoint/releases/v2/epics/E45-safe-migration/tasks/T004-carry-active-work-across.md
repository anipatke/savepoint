---
id: E45-safe-migration/T004-carry-active-work-across
title: Carry active work across
status: done
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

- [x] Rendered Objective content decodes through `DecodeObjectiveV2` and rendered Task content through `DecodeTaskV2` with no diagnostic, for every active record in both fixtures.
- [x] An active V1 epic becomes one Objective whose `title` comes from the epic title, whose `status` is `planned` or `in_progress`, whose `depends_on` maps declared epic dependencies to allocated `O###` values, and whose optional `release` preserves the source release as filtering metadata.
- [x] An epic whose Tasks are all done but which was never audited converts to an `in_progress` Objective with no clearance, so it must still earn an integration Check rather than arriving complete.
- [x] An active V1 Task becomes a Task owned by exactly one `O###`, with `depends_on` entries naming only converted active Tasks, and with the `{task, requires}` shape E43 decodes.
- [x] No converted record carries `last_check`, `freshness`, `owner_validation.accepted_check`, or an exception unless the V1 source recorded the corresponding fact; a test asserts a converted fixture Task has unknown clearance and is therefore not completable.
- [x] A `status: in_progress` Task carries a stage mapped from the source's raw `stage` or legacy `phase`, applying the recorded `implementation` → `build` alias and nothing beyond it.
- [x] An unrecognized raw status or stage on an active record produces a blocking ambiguity rather than a healed value, and a test proves the V1 loader's healing is not reachable from this path.
- [x] Authored body content is relocated verbatim under V2 headings, byte for byte including the fixtures' CRLF line endings and multiline authored text, with the V1 section each block came from named in the rendered output.
- [x] Unknown frontmatter fields present on a V1 source are preserved on the converted record or recorded in the manifest; neither silently disappears.
- [x] A converted Task whose V1 dependency was an archived completed Task carries the authored legacy-prerequisite line naming the archive path and the original completion evidence, and carries no dependency on a nonexistent Task.
- [x] Rendering is deterministic: the same source and injected clock produce byte-identical output across runs.

## Implementation Plan

- [x] Add `convert.go` with the Objective and Task renderers, taking planned targets from T003 and returning file content rather than writing it.
- [x] Implement frontmatter construction against the `ObjectiveV2` and `TaskV2` field contracts, deliberately omitting every evidence field the source did not record.
- [x] Implement the lifecycle mapping table — including the unaudited-epic rule and the legacy stage alias — and route every unrecognized value to a blocking ambiguity.
- [x] Implement verbatim body relocation with source-section attribution, preserving original line endings rather than normalizing them.
- [x] Render the authored legacy-prerequisite lines from the manifest entries T003 produces.
- [x] Add round-trip tests that decode every rendered record through the real V2 decoders and assert no diagnostic.
- [x] Test the no-fabrication rules, the unaudited-epic outcome, unrecognized lifecycle values, CRLF and unknown-field preservation, and byte-identical repeat rendering.
- [x] Run the focused `internal/migrate` and `internal/data` suites, then `make build && make test`.

## Context Log

**Files read:** `plan.go`, `manifest.go`, `inventory.go` (internal/migrate); `objective_v2.go`, `task_v2.go`, `evidence_v2.go`, `lifecycle.go`, `parser.go`, `write.go`, `task.go` (internal/data); both frozen fixture projects and manifests under `internal/data/testdata/migration/`.

**Files written:**
- `internal/migrate/convert.go` — `ConvertObjective` and `ConvertTask`, rendering V2 file content from a `PlannedTarget` plus the full `ConversionPlan` (needed for epic-level dependency resolution and legacy-prerequisite lookup). Re-derives status/stage from raw frontmatter through `resolveEpicStatus`/`resolveTaskStatus` (reused from plan.go, same package) plus a new non-healing `resolveTaskStage`, so `ParseTaskFile`'s healing path is never on the call chain. Unknown V1 frontmatter fields are preserved under a `legacy_fields` map (yaml.v3 sorts map keys on marshal, so this stays deterministic). Body is relocated verbatim by wrapping the already-newline-normalized `SplitFrontmatterBody` output under one "V1 Body (verbatim)" heading naming the source path, then the whole rendered file is CRLF-reapplied at the end if the source was CRLF — mirroring the existing `writeV2Record` convention in `write.go` exactly, so round-tripping CRLF sources is byte-exact.
- `internal/migrate/convert_test.go` — round-trip decode over every active target in both frozen fixtures; Objective field mapping (title from H1, status, release, depends_on) including a synthetic project for the unaudited-epic-with-all-tasks-done rule and for epic-level `depends_on` resolution/refusal; Task field mapping including the `phase: implementation` → `build` alias and the release-scoped same-short-ID dependency case in `v1-history`; the legacy-prerequisite authored line; two synthetic-content tests proving unrecognized status/stage refuse via `ErrAmbiguousLifecycle` rather than healing (one directly demonstrates `ParseTaskFile` would have healed the same content, to make the "not reachable from this path" claim concrete); CRLF and unknown nested-field (`legacy_note`, `metadata.reviewer`) preservation against the `v1-basic` CRLF fixture; determinism across repeated calls.

**Quality gates:** `go build ./...` clean. `go test ./...` — all packages pass, including the new `internal/migrate` tests. `make build && make test` — pass. `gofmt -l` clean on both new files.

**Health-Check.md:** absent from this project; Quick check step skipped per the build-task skill (absence is not a finding).

**Drift:** none. No new files/modules outside the epic's documented `internal/migrate` package map, and no architecture changed beyond what E45-Detail.md already describes (T004's renderers, exactly as scoped).
