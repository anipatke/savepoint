---
id: E47-onboarding-upgrades/T002-start-a-new-project-on-the-v2-workflow
title: Start a new project on the V2 workflow
status: done
objective: Make savepoint init scaffold the V2 tree by default, retire V1 scaffolding, and prove a fresh project loads as a valid empty V2 project.
depends_on:
    - E47-onboarding-upgrades/T001-ship-the-v2-project-assets-from-their-own-tree
complexity_tier: medium
complexity_reason: Flips a user-facing default and invalidates several existing scaffold tests that assert V1 output.
---

# T002: Start a new project on the V2 workflow

## Problem

`savepoint init` creates a V1 project: a release skeleton, a PRD, a Concept, a Health-Check, an audit register, and a router at `state: pre-implementation`. Everything E42 through E46 built is invisible to anyone who runs it.

This task points `init` at the tree T001 authored, and retires V1 scaffolding in the same change. There is no `--v1` flag: a command that can still produce either lifecycle is a second scaffold to keep correct, verify, and document, for a user who does not exist — anyone with a V1 project already has one, and `migrate` is how it becomes V2.

The flip has to be proven at the level that matters, which is not "the files were written". A fresh project must load through `data.LoadProject` as a V2 project with an empty index and no diagnostic. If the scaffold cannot be read back by the loader that owns V2 records, it is not a V2 project, whatever its file names say.

Init into a directory that already holds source files is the same operation and must stay one: the scaffold adds its own files and touches nothing else. An existing agent guide keeps everything outside the managed block, exactly as it does today.

Several existing scaffold tests assert V1 output — the release skeleton, the audit register assets, the Guardrails-and-Health-Check pair. Those assertions describe the legacy tree, not the command, so they move to the legacy tree or retire with it; the behavior they protect is re-proven against `templates/project-v2/` where it still applies.

## Context Files

- `main.go`
- `main_test.go`
- `cmd/init.go`
- `cmd/init_test.go`
- `internal/init/scaffold.go`
- `internal/init/scaffold_test.go`
- `internal/init/validate.go`
- `internal/init/integration_test.go`
- `internal/init/manifest.go`
- `internal/init/manifest_test.go`
- `internal/init/lifecycle_test.go`
- `internal/data/project.go`
- `internal/data/config.go`
- `AGENTS.md`
- `.savepoint/releases/v2/v2-Design.md`

## Acceptance Criteria

- [x] `savepoint init` in an empty directory writes the `templates/project-v2/` tree and writes no file from `templates/project/`.
- [x] `data.ReadSchemaVersion` on the scaffolded `config.yml` returns `SchemaVersionV2`.
- [x] `data.LoadProject` on the scaffolded `.savepoint/` returns a project at `SchemaVersionV2` with a non-nil V2 index holding no Objectives, Tasks, Checks, or Issues, and no error.
- [x] The scaffolded project contains no `releases/`, `epics/`, `PRD.md`, `Concept.md`, `Health-Check.md`, or `audit/` path.
- [x] `.savepoint/.upgrade-manifest.yml` is written with a provenance entry for each of the four V2 skills and for no V1 skill.
- [x] Init into a directory holding existing source files leaves every pre-existing file byte-identical, verified by content hash before and after.
- [x] Init into a directory with an existing `AGENTS.md` preserves all content outside the managed block, and the casing variant behavior in `TestIntegration_ExistingAgentGuideCasingVariant` still holds.
- [x] `init` exposes no flag that scaffolds a V1 project, and `ParseInitArgs` rejects an unknown flag as it does today.
- [x] `ValidateTarget`'s refusals are unchanged: missing target, non-directory, unwritable, existing `.savepoint/` without `--force`, and an existing `agent-skills/` directory without `--force`.
- [x] Scaffold tests that assert V1-only output — release skeleton, audit register assets, Health-Check delivery — either target `templates/project/` explicitly or are removed with a recorded reason; no test silently passes by no longer being reached.
- [x] The `AGENTS.md` Codebase Map row for `internal/init/` describes what the package now does, satisfying ARCH-04.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Point `initRunner` in `main.go` at the `templates/project-v2` sub-filesystem; leave `Scaffold` itself unchanged so the tree stays the caller's choice.
- [x] Walk `internal/init/scaffold_test.go` and `internal/init/integration_test.go` case by case, deciding for each whether it describes the command or the legacy tree, and re-point or retire it accordingly with the reason recorded in the Context Log.
- [x] Add a fresh-init test that reads the scaffold back through `data.ReadSchemaVersion` and `data.LoadProject` and asserts an empty, diagnostic-free V2 index.
- [x] Add an absence test for the V1-only paths, driven by a list, so a re-added template file fails loudly.
- [x] Add the existing-codebase test: hash every file in a populated temporary directory, run init, and compare hashes.
- [x] Confirm the manifest records the four V2 skills and nothing else.
- [x] Update the `internal/init/` Codebase Map row in `AGENTS.md` and mirror it into `templates/project/AGENTS.md`.
- [x] Run `go test ./internal/init/... ./cmd/... .`, then `make build && make test`.

## Context Log

**Files read (Context Files list):** `main.go`, `main_test.go`, `cmd/init.go`, `cmd/init_test.go`, `internal/init/scaffold.go`, `internal/init/scaffold_test.go`, `internal/init/validate.go`, `internal/init/integration_test.go`, `internal/init/manifest.go`, `internal/init/manifest_test.go`, `internal/init/lifecycle_test.go`, `internal/data/project.go`, `internal/data/config.go`, `AGENTS.md`, `.savepoint/Guardrails.md` (STYLE rules).

**Extra reads beyond Context Files (logged per savepoint-build-task):**
- `.savepoint/router.md`, `.savepoint/releases/v2/epics/E47-onboarding-upgrades/E47-Detail.md`, T001's task file — to confirm T001 shipped `templates/project-v2/` and that the repo's uncommitted working tree already carried that work.
- `internal/init/v2_scaffold_test.go` (T001-authored) — to see the existing direct-`Scaffold()`-against-`templates/project-v2` proof and avoid duplicating it; this task's new tests instead exercise the real `init` command end to end.
- `templates/project-v2/**` file listing, `templates/prompts/magic-prompt.prompt.md`, `internal/init/prompt.go`, `internal/init/agents.go` (grep only, to confirm `FindAgentGuide`/`MergeAgentGuide`/`AtomicWrite` call sites), `internal/init/clipboard.go` — to confirm the prompt template stays shared (T003's scope, unchanged here) and that `CopyToClipboard` fails soft in a headless test environment.
- `.savepoint/releases/v2/epics/E47-onboarding-upgrades/tasks/T003-begin-from-one-rough-sentence.md` — to confirm the magic-prompt rewrite is explicitly T003's scope, not this task's.
- `templates/project/AGENTS.md` — to check whether its `## Codebase Map` table has a row to mirror; it is an empty placeholder table (header only) shipped for the user's own project, so there is no `internal/init/` row to update there. The Implementation Plan's "mirror it into `templates/project/AGENTS.md`" does not apply: that file's Codebase Map has no content to mirror, and its content is confirmed by `TestV2ScaffoldSavepointFileSet`/parity tests to be about the user's project, not Savepoint's own package layout. Only the live root `AGENTS.md` row was updated.

**Files edited:**
- `main.go` — `initRunner` now builds its scaffold sub-filesystem from `projectTemplatesV2`/`templates/project-v2` instead of `projectTemplates`/`templates/project`. The prompt sub-filesystem (`templates/prompts`) is untouched, matching the epic's "one prompt template" decision. `Scaffold` itself was not touched.
- `internal/init/scaffold_test.go` — removed `TestScaffold_createsReleaseSkeleton` and `TestScaffold_createsAuditRegisterAssets`. Reason recorded here per the Implementation Plan: both used a synthetic `fstest.MapFS` fixture (not a real shipped tree) whose path names were V1-flavored (`releases/v1/epics`, `v1-PRD.md`, `.savepoint/audit/...`) to exercise nothing beyond `Scaffold`'s generic nested-directory-creation and multi-file-write mechanics — mechanics already covered by `TestScaffold_createsParentDirs` (deep nested dirs) and `TestScaffold_createsDirectories` (dir + file). `TestScaffold_installsSplitAuditSkillsAndSharedMethod` and `TestScaffold_installsGuardrailsAndHealthCheck` (the third test the Problem section names, "the Guardrails-and-Health-Check pair") already target `templates/project/` explicitly via `scaffoldFromRealTemplates()` and needed no change. `internal/init/integration_test.go`'s `runInitPipeline`-based tests also use a synthetic, tree-agnostic fixture to test pipeline sequencing (validate → scaffold → render prompt) and agent-guide merge/casing behavior independent of which real tree ships; left unchanged as they describe the command's mechanism, not the legacy tree's content.
- `main_test.go` — added `internal/data` and `internal/init` (aliased `savepointinit`, matching `main.go`) imports, and five new tests that run the real compiled `init` command via the existing `runMainForTest` helper (so they embed the same `//go:embed` trees production does): `TestMainInitScaffoldsV2ProjectWithEmptyValidIndex` (schema version + `data.LoadProject` empty-index proof), `TestMainInitWritesNoV1OnlyPathOrSkill` (absence list for V1-only paths and all nine V1 skills, driven by `v1OnlyPaths`/`v1OnlySkills`), `TestMainInitManifestRecordsExactlyTheFourV2Skills`, and `TestMainInitOnExistingCodebasePreservesEveryFileByteIdentical` (seeds `package.json`, `README.md`, `src/index.js`, `.gitignore`, hashes before/after `init`). Added a small `fileHash` test helper (sha256/hex already imported for the migrate tests).
- `AGENTS.md` — the `internal/init/` Codebase Map row now states that `Scaffold` writes from a caller-selected template tree and names both defaults (`init` → `templates/project-v2`, `upgrade-assets` → `templates/project`), satisfying ARCH-04 without restating STYLE rules.

**Files not edited despite being named in the Implementation Plan, with reason:** `templates/project/AGENTS.md` — see "Extra reads" above; its Codebase Map is an empty placeholder, so there is no row to mirror.

**Commands run:**
- `go build ./...` — pass, after the `main.go` change.
- `go vet ./internal/init/...`, `go vet ./...` — clean.
- `gofmt -l main.go main_test.go internal/init/scaffold_test.go` — clean.
- `go test ./... -run 'TestMainInit' -v` — all 5 new tests pass.
- `go test ./internal/init/... ./cmd/... .` — pass.
- `make build && make test` — pass, all packages green.
- Manual smoke test: built a throwaway binary, ran `init` against an empty temp dir, confirmed the on-disk tree is exactly the V2 scaffold (four skills, `Idea.md`/`Design.md`/`Guardrails.md`/`config.yml`/`router.md`/`objectives/.gitkeep`/`AGENTS.md`/`.upgrade-manifest.yml`, no V1 paths), then removed the scratch directory and binary.

**Limitations:** This repository's own `.savepoint/` stays on the V1 lifecycle per the epic's boundaries — no live record here was touched beyond the `AGENTS.md` Codebase Map row and this task file. `templates/prompts/magic-prompt.prompt.md` and the V2 router's opening `next_action` still name V1-era routing prose; rewriting them for the V2 bootstrap is T003's scope and was left untouched here. No `.savepoint/Health-Check.md` exists in this repository, so the Quick health-check step in `savepoint-build-task`'s workflow was skipped per its own rule (absence is not a finding).
