---
id: E47-onboarding-upgrades/T001-ship-the-v2-project-assets-from-their-own-tree
title: Ship the V2 project assets from their own tree
status: done
objective: Add templates/project-v2/ holding every V2 scaffold file and the V2 skills, move the V2 assets out of the V1 tree, and make the canonical-to-shipped parity contract per-tree.
depends_on: []
complexity_tier: medium
complexity_reason: Many authored files plus a rewrite of the parity tests that currently derive pairing from one directory listing.
---

# T001: Ship the V2 project assets from their own tree

## Problem

Every shipped asset lives in one tree, `templates/project/`, and that tree describes a V1 project: a router at `state: pre-implementation`, a config with no `schema_version`, an agent guide whose active routing table is the V1 one, a release skeleton, and an audit register. E46 put the four V2 skills in there too, marked as future state.

V2 cannot share that tree. `config.yml`, `router.md`, and `AGENTS.md` exist at the same paths in both lifecycles with content that cannot be reconciled — one says `state: pre-implementation` and the other `state: idea`; one declares no schema version and the other declares 2. A single tree would have to pick one, and picking either breaks the other lifecycle.

So this task creates the second tree and moves the V2 assets into it, with no command yet pointing at it. Nothing about `init` or `upgrade-assets` changes here; the tree is authored and contract-tested first so the behavior changes that follow have something correct to switch to.

Moving the skills breaks a contract on the way. `TestProjectGuidanceTemplatesMirrorLiveGuidance` discovers the canonical-to-shipped pairing by listing `agent-skills/` and requiring every entry to appear in `templates/project/`. With two trees that pairing is no longer derivable from a listing, and it becomes an explicit table. TPL-01 names the single shipped path in its prose, so the rule's wording is amended to match the code — the rule itself does not change.

## Context Files

- `templates/project/AGENTS.md`
- `templates/project/.savepoint/config.yml`
- `templates/project/.savepoint/router.md`
- `templates/project/.savepoint/Design.md`
- `templates/project/.savepoint/Guardrails.md`
- `templates/project/agent-skills/savepoint-idea/SKILL.md`
- `templates/project/agent-skills/savepoint-design/SKILL.md`
- `templates/project/agent-skills/savepoint-task/SKILL.md`
- `templates/project/agent-skills/savepoint-check/SKILL.md`
- `templates/project/agent-skills/references/check-method.md`
- `templates/project/agent-skills/references/issue-capture.md`
- `templates/project/agent-skills/references/commands-and-procedures.md`
- `agent-skills/savepoint-idea/SKILL.md`
- `agent-skills/savepoint-design/SKILL.md`
- `internal/init/skill_validation_test.go`
- `internal/init/template_freshness_test.go`
- `internal/init/agent_skills_test.go`
- `internal/init/provenance_contract_test.go`
- `internal/init/stale_reference_test.go`
- `main.go`
- `AGENTS.md`
- `.savepoint/Guardrails.md`
- `.savepoint/releases/v2/v2-Design.md`

## Acceptance Criteria

- [x] `templates/project-v2/.savepoint/` contains `Idea.md`, `Design.md`, `Guardrails.md`, `config.yml`, `router.md`, and `objectives/.gitkeep`, and contains no `PRD.md`, `Concept.md`, `Health-Check.md`, `releases/`, or `audit/`.
- [x] `templates/project-v2/.savepoint/config.yml` declares `schema_version: 2` and carries `quality_gates` with `lint`, `typecheck`, `build`, `test`, `block_on_failure`, and `gate_timeout`, plus the existing theme keys. It declares no `audit` key.
- [x] `templates/project-v2/.savepoint/router.md` opens at `state: idea` with an `objective` field, and its routing table maps `idea`/`design`/`task`/`check` to the four V2 skills with no V1 state named.
- [x] `templates/project-v2/.savepoint/Design.md` carries the six V2 sections — Architecture, Components/Codebase Map, Interfaces and Data Flow, Boundaries, Decisions, Current Technical State.
- [x] `templates/project-v2/AGENTS.md` presents the four-state routing table as live, with no "inactive until cutover" wording, and uses V2 vocabulary only: no `epic`, no `PRD`, no `audit register`, no `defect-building`, and no `phase` as a lifecycle term.
- [x] `templates/project-v2/agent-skills/` contains exactly `savepoint-idea`, `savepoint-design`, `savepoint-task`, `savepoint-check`, and `references/check-method.md`, `references/issue-capture.md`, `references/commands-and-procedures.md`.
- [x] Those seven files are byte-identical to their canonical sources under `agent-skills/`, and none of them remains under `templates/project/agent-skills/`.
- [x] `templates/project/agent-skills/` contains exactly the nine V1 skills and `references/audit-method.md`, and `templates/project/` is otherwise unchanged byte for byte.
- [x] A per-tree parity test asserts the pairing from an explicit table in both directions: every canonical skill has a shipped copy in its declared tree, and neither tree carries a skill or reference belonging to the other.
- [x] `TestV2SkillSetIsCompleteWithByteParity`, `TestV2SkillSetHasNoObsoleteVocabulary`, `TestV2SkillRoleMatrixPartitionsSensitiveWrites`, `TestSavepointSkillsHaveValidFrontmatter`, and `TestV2ArtifactTemplatesDecodeThroughTypedContracts` all pass against the V2 assets at their new path.
- [x] `main.go` embeds `templates/project-v2` including its dotfiles, and `go build ./...` resolves every embedded path.
- [x] The live `AGENTS.md` V2 routing section states that the table is active for V2 projects and not for this repository until E50, and `templates/project/AGENTS.md` mirrors that wording.
- [x] `.savepoint/Guardrails.md` TPL-01 names both shipped trees; its severity and intent are unchanged.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Create `templates/project-v2/.savepoint/` and author `config.yml`, `router.md`, `Idea.md`, `Design.md`, `Guardrails.md`, and `objectives/.gitkeep`, using the existing templates' instructional style and `{{PROJECT_NAME}}` interpolation.
- [x] Author `templates/project-v2/AGENTS.md`: V2 workflow, four-state routing table as live guidance, V2 terminology, the shared-reference list, and the build command section.
- [x] Move the four V2 skills and three V2 references from `templates/project/agent-skills/` to `templates/project-v2/agent-skills/`, copying bytes rather than re-authoring.
- [x] Replace the discovery-based pairing in `TestProjectGuidanceTemplatesMirrorLiveGuidance` with an explicit skill-to-tree table, and extend `skillRoots()` in `internal/init/skill_validation_test.go` to cover both shipped trees.
- [x] Add the set-completeness assertions: each tree holds exactly its declared assets, in both directions.
- [x] Add the V2 scaffold content assertions — the file set, `schema_version: 2`, the router's opening state, the Design sections, and the agent guide's vocabulary.
- [x] Add the `templates/project-v2` embed directives to `main.go`, including the `all:` form for `.savepoint`.
- [x] Update the live `AGENTS.md` V2 routing statement and mirror it into `templates/project/AGENTS.md`; amend TPL-01 in `.savepoint/Guardrails.md`.
- [x] Run `go test ./internal/init/...`, then `make build && make test`.

## Context Log

**Files read (Context Files list):** `templates/project/AGENTS.md`, `templates/project/.savepoint/{config.yml,router.md,Design.md,Guardrails.md,PRD.md,Concept.md}`, all four V2 `agent-skills/{savepoint-idea,savepoint-design,savepoint-task,savepoint-check}/SKILL.md`, all three V2 `agent-skills/references/{check-method,issue-capture,commands-and-procedures}.md`, `internal/init/{skill_validation_test.go,template_freshness_test.go,agent_skills_test.go,provenance_contract_test.go,stale_reference_test.go}`, `main.go`, `AGENTS.md`, `.savepoint/Guardrails.md`, `.savepoint/releases/v2/v2-Design.md` (E47-Detail.md itself).

**Extra reads/edits beyond Context Files (logged per savepoint-build-task):**
- `internal/data/config.go`, `internal/data/project.go`, `internal/data/discover.go` — read to confirm `schema_version: 2` field name, `QualityGates` struct keys, and that `LoadV2Index` tolerates a missing/empty `objectives/` tree, so the fresh scaffold would load clean.
- Root-level `agent_skills_test.go` (package `main_test`, not under `internal/init/`) — not in the Context Files list, but `make test` failed there (`TestScaffoldedSavepointSkillsMatchBundledSkills` assumed one shipped tree). Made it tree-aware instead of hardcoding a single shipped path; also extended `TestBundledSavepointSkillsHaveDiscoveryFrontmatter` to cover `templates/project-v2/agent-skills`.
- Added `internal/init/v2_scaffold_test.go` (new file) to give the Implementation Plan's "V2 scaffold content assertions" item permanent test coverage instead of one-time manual verification: file-set/forbidden-file checks, `config.yml` schema/quality-gates content, `router.md` opening state, `Design.md`'s six sections, `AGENTS.md` vocabulary, and an end-to-end `Scaffold()` → `data.LoadProject()` proof that a fresh V2 tree loads with no diagnostic and an empty index (schema dispatch and `objectives/` handling are existing `internal/data` behavior; this task did not modify `internal/data`).

**Files created:** `templates/project-v2/.savepoint/{config.yml,router.md,Idea.md,Design.md,Guardrails.md,objectives/.gitkeep}`, `templates/project-v2/AGENTS.md`, `internal/init/v2_scaffold_test.go`.

**Files moved (`git mv`, bytes unchanged):** `templates/project/agent-skills/{savepoint-idea,savepoint-design,savepoint-task,savepoint-check}/SKILL.md` and `templates/project/agent-skills/references/{check-method,issue-capture,commands-and-procedures}.md` → the equivalent paths under `templates/project-v2/agent-skills/`.

**Files edited:** `main.go` (new `projectTemplatesV2` embed var), `AGENTS.md` and `templates/project/AGENTS.md` (V2 routing section reworded — heading `## V2 Routing`, body states the table is active for V2 projects and inactive for this repository until E50), `.savepoint/Guardrails.md` (TPL-01 names both shipped trees), `internal/init/skill_validation_test.go` (added `liveSkillRoot`/`v1TemplateSkillRoot`/`v2TemplateSkillRoot`/`v2SkillRoots`/`allSkillRoots`; the two discovery-based structural tests now use `allSkillRoots`), `internal/init/agent_skills_test.go` (53 V2-scoped test functions switched from `skillRoots()` to `v2SkillRoots()`; 9 hardcoded `templates/project/agent-skills` paths repointed to `templates/project-v2/agent-skills`; routing-heading tests updated for the new heading/wording; the 4 V1-scoped tests left on `skillRoots()` unchanged), `internal/init/provenance_contract_test.go` (`TestV2ArtifactTemplatesDecodeThroughTypedContracts` now uses `v2SkillRoots()`), `internal/init/template_freshness_test.go` (`TestProjectGuidanceTemplatesMirrorLiveGuidance` rewritten around a new `assertSkillTreeParity` helper: explicit V1/V2 name tables, byte parity and set-completeness in both directions per tree, plus a check that the live tree's skill set is exactly the union of the two declared tables), root `agent_skills_test.go` (made tree-aware, see above).

**Commands run:**
- `go build ./...` — pass, embeds resolve.
- `go vet ./...` — clean.
- `gofmt -l` on every file touched — clean (pre-existing formatting debt in `internal/init/clipboard.go` and `internal/init/manifest_test.go` predates this task and was left untouched).
- `go test ./internal/init/...` — pass.
- `make build && make test` — pass, all packages green (`internal/init` 2.09s, rest cached/pass).

**Limitations:** No command wires `templates/project-v2` into `init` or `upgrade-assets` yet — by design; T001 authors and contract-tests the tree only. `TestV2ScaffoldLoadsCleanThroughDataLoadProject` exercises `Scaffold()` directly against `templates/project-v2` to prove the tree is structurally sound ahead of that wiring. This repository's own `.savepoint/` stays on the V1 lifecycle, unchanged except the Guardrails/AGENTS.md edits named above, per the epic's boundaries.
