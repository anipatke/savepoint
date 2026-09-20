---
id: E46-agent-workflow-assets/T008-carry-project-commands-and-procedures-forward
title: Carry project commands and procedures forward
status: done
objective: Write the guidance that maps a legacy Health-Check into config quality gates and one preserved optional procedure, with no substitute invented when it is absent.
depends_on:
    - E46-agent-workflow-assets/T001-write-the-shared-checking-method-once
complexity_tier: medium
complexity_reason: A short reference, but it must match the existing config keys and the migration behavior E45 already implemented.
---

# T008: Carry project commands and procedures forward

## Problem

A V1 project's `Health-Check.md` holds two different things that V2 stores in two different places. The commands — lint, typecheck, test, build — belong in `config.yml` under the existing `quality_gates` key, where execution and timeouts already work and where doctor already runs them. The universal Quick/Full mechanics belong in the shared check method, written once for every project. What is left is genuinely project-specific prose: a manual verification workflow somebody wrote because their project needs it, which has no home in either and must not be deleted.

That residue becomes an ordinary optional procedure file the project references, not a mandatory new default. And when a project has no `Health-Check.md` at all — this repository does not — nothing is generated to stand in for it. Absence is not a gap to fill.

The failure mode this guidance prevents is a parallel config surface: inventing `checks.technical` alongside `quality_gates` because the new model felt cleaner, and leaving two places where commands live.

## Context Files

- `agent-skills/references/commands-and-procedures.md`
- `templates/project/agent-skills/references/commands-and-procedures.md`
- `agent-skills/references/check-method.md`
- `agent-skills/savepoint-design/SKILL.md`
- `templates/project/agent-skills/savepoint-design/SKILL.md`
- `templates/project/.savepoint/Health-Check.md`
- `templates/project/.savepoint/config.yml`
- `internal/data/config.go`
- `internal/doctor/checks.go`
- `internal/migrate/plan.go`
- `internal/init/agent_skills_test.go`
- `.savepoint/releases/v2/v2-Design.md`

## Acceptance Criteria

- [x] `agent-skills/references/commands-and-procedures.md` exists with `type: commands-reference`, `triggerable: false`, no `name` key, and a byte-identical shipped copy.
- [x] It names `quality_gates` in `config.yml` as the single place project commands live, using the key names `internal/data/config.go` already decodes, and states that no parallel `checks.technical` surface is introduced.
- [x] It states that an explicit build command is supported alongside the existing gate commands, and describes where it runs.
- [x] It maps a legacy `Health-Check.md` in three parts: commands into `quality_gates`, universal Quick/Full mechanics into `agent-skills/references/check-method.md`, and project-specific reusable prose into one preserved optional procedure file.
- [x] It states that the preserved procedure is an ordinary referenced project file — named by Tasks that need it — and never a mandatory scaffold default.
- [x] It states that a project with no `Health-Check.md` needs no generated substitute and that its absence is not a finding.
- [x] It states that a Task naming extra verification does so in its own Technical Verification section rather than by editing shared policy.
- [x] The mapping described matches what `internal/migrate` actually plans for project documents; any divergence found is recorded as a drift note and raised, not silently written into guidance as if true.
- [x] `savepoint-design` points at this reference when reconciling a migrated project's commands and procedures, without restating the mapping.
- [x] `internal/init/agent_skills_test.go` asserts the reference frontmatter contract, the `quality_gates` naming, the absence rule, and that no `checks.technical` key appears anywhere in the V2 assets.

## Implementation Plan

- [x] Read design section 8's verification paragraph, `internal/data/config.go` for the live `quality_gates` shape, and `internal/migrate/plan.go` for what migration already does with `Health-Check.md`.
- [x] Write `agent-skills/references/commands-and-procedures.md` with the three-part mapping, the config contract, and the absence rule.
- [x] Add the pointer from `savepoint-design` to this reference for post-migration reconciliation.
- [x] Mirror both files to the template tree.
- [x] Extend `internal/init/agent_skills_test.go` with the reference contract case and the `checks.technical` absence assertion.
- [x] Record any divergence between this guidance and `internal/migrate`'s actual behavior in `## Drift Notes` rather than adjusting the guidance to match an assumption.
- [x] Run `go test ./internal/init/... ./internal/migrate/...`, then `make build && make test`.

## Context Log

**Files read:** `.savepoint/releases/v2/v2-Design.md` (section 8), `internal/data/config.go`, `internal/doctor/gates.go`, `internal/doctor/gates_test.go`, `internal/migrate/plan.go`, `internal/migrate/convert_docs.go`, `.savepoint/releases/v2/epics/E45-safe-migration/tasks/T006-carry-the-project-documents-across.md`, `.savepoint/releases/v2/epics/E45-safe-migration/E45-Detail.md`, `agent-skills/references/check-method.md`, `agent-skills/references/issue-capture.md`, `agent-skills/savepoint-design/SKILL.md`, `internal/init/agent_skills_test.go`, `templates/project/.savepoint/config.yml`, `templates/project/.savepoint/Health-Check.md`.

**Files edited:**
- `internal/data/config.go` — added `QualityGates.Build *string` (`yaml:"build"`), matching the existing `Lint`/`Typecheck`/`Test` shape.
- `internal/doctor/gates.go` — `RunQualityGates` runs the `build` gate between `typecheck` and `test`, with the same timeout/`block_on_failure` semantics as the other gates.
- `internal/doctor/gates_test.go` — added `TestRunQualityGates_BuildRunsBetweenTypecheckAndTest` and `TestRunQualityGates_BuildOnly`.
- `templates/project/.savepoint/config.yml` — added `build: null` to `quality_gates` alongside `lint`/`typecheck`/`test`.
- `agent-skills/references/commands-and-procedures.md` (new) and its mirror under `templates/project/agent-skills/references/` — the config contract, the three-part Health-Check.md mapping, and the absence rule.
- `agent-skills/savepoint-design/SKILL.md` and its template mirror — added `## Commands And Procedures Reconciliation` pointing at the new reference.
- `internal/init/agent_skills_test.go` — added the frontmatter-contract, config-key, mapping-phrase, live/template-match, and `savepoint-design`-pointer tests for the new reference.

**Quality gates:** `go test ./internal/init/... ./internal/migrate/...` and `make build && make test` all pass.

## Drift Notes

Design section 8 (`v2-Design.md`) says to "add explicit build command support" to `quality_gates`, but before this task neither `internal/data/config.go`'s `QualityGates` struct nor `internal/doctor/gates.go`'s `RunQualityGates` had a `build` field or gate — E45 T006 explicitly carried `quality_gates` across from V1 verbatim and never added one. AC3 required the reference to state build support truthfully, so this task implemented the minimal `Build *string` field and wired it into `RunQualityGates` (ordered after `typecheck`, before `test`) rather than document behavior that did not exist. This was raised to the project owner before implementing rather than silently assumed.

Separately, `internal/migrate`'s `Plan` does not itself perform the three-part Health-Check.md split AC4 describes: it archives `Health-Check.md` byte-for-byte and records fenced-code-block lines as `ArchiveEntry.CandidateCommands` for the preview only — it never writes `config.yml` and never creates a procedure file. The reference states this explicitly (`### What migrate actually does`) so the three-part mapping reads as reconciliation guidance for `savepoint-design` to apply by hand, not as a claim about migrate's automated behavior.
