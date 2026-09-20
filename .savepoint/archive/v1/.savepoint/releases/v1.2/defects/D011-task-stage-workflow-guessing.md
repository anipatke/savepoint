---
id: v1.2/D011-task-stage-workflow-guessing
release: v1.2
status: resolved
severity: high
title: "Task stage workflow contract is not enforced consistently"
---

# D011: Task Stage Workflow Contract Is Not Enforced Consistently

## Symptom

While testing Savepoint, agents still update task `status` and `stage` inconsistently instead of following the documented workflow. The behavior feels guessed rather than mechanically guided.

## Expected Behavior

Task lifecycle handling should have one canonical contract:

- Tasks use `status: planned`, `status: in_progress`, or `status: done`.
- Tasks with `status: in_progress` must include `stage: build`, `stage: test`, or `stage: audit`.
- Tasks outside `in_progress` must not include `stage`.
- Agents should receive, see, and be validated against the same field names and lifecycle rules.

## Reproduction

Investigation found conflicting task lifecycle signals:

- `AGENTS.md`, scaffolded `templates/project/AGENTS.md`, and `agent-skills/savepoint-build-task/SKILL.md` instruct agents to write `stage: build`.
- `internal/data/write.go` writes task progress back as `phase: <value>` and removes `stage`.
- `internal/data/parser.go` accepts both `phase` and `stage`, but prefers `phase` over `stage`.
- Task lifecycle validation defaults missing in-progress stage to build instead of treating omission as a parse error.
- Task parser/lifecycle error messages still tell users to add `phase: build`.
- Several task and board tests assert `phase:` output, while template freshness tests assert that current guidance says `stage:`.

## Impact

Agents are exposed to contradictory examples and weak validation. The docs say `stage`, the board/data writer emits `phase`, and missing stage is silently repaired to build. This makes the workflow dependent on agent interpretation instead of a strict product invariant.

## Fix Plan

1. Choose `stage` as the canonical task frontmatter field because it matches current AGENTS/skill guidance and defect lifecycle terminology.
2. Update task parsing to accept legacy `phase` only as backward-compatible input, but prefer canonical `stage` when both are present and produce a repair warning or doctor problem for legacy `phase`.
3. Update task writing to emit `stage` for `in_progress` tasks and remove legacy `phase`; continue removing `stage` when status moves to `planned` or `done`.
4. Tighten task lifecycle validation so parsed task files with `status: in_progress` and neither `stage` nor legacy `phase` fail validation instead of defaulting to build.
5. Update all task lifecycle error messages from `phase` wording to `stage` wording.
6. Add doctor checks for task files that report missing in-progress stage, invalid stage, stage outside in-progress, and legacy `phase` usage with a concrete repair suggestion.
7. Update board/data tests that still assert `phase:` to assert `stage:` and add compatibility tests for old `phase:` files.
8. Add template/skill freshness tests that reject `phase:` examples in task workflow docs and reject code/test messages that instruct agents to write task `phase`.
9. Add a compact workflow preflight to the build skill instructions: before implementation, verify the active task frontmatter has canonical `status`/`stage`; after handoff, do not mark `done`.

## Acceptance Criteria

- [x] Task writer persists `stage: build|test|audit` for in-progress tasks and does not persist `phase`.
- [x] Task parser still reads legacy `phase` files but canonical `stage` wins when both fields exist.
- [x] Missing task stage for `status: in_progress` is reported as invalid by parser or doctor instead of silently defaulting.
- [x] Doctor reports actionable problems for task lifecycle field errors.
- [x] Agent-facing docs, scaffolded docs, skills, parser messages, and tests use one task lifecycle vocabulary.
- [x] `make build && make test` passes.

## Resolution Notes

Resolved. Task lifecycle handling now uses `stage` as the canonical in-progress frontmatter field. Legacy task `phase` remains readable for compatibility, but task writes remove it and doctor reports it as a repairable problem. Missing in-progress stage and invalid values such as `done` now produce explicit stage diagnostics. Verification: `go test ./internal/data ./internal/doctor ./internal/board`, `make build`, and `make test` passed.
