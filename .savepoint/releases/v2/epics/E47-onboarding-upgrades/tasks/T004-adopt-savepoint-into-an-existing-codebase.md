---
id: E47-onboarding-upgrades/T004-adopt-savepoint-into-an-existing-codebase
title: Adopt Savepoint into an existing codebase
status: done
objective: Ship guidance that reconstructs Design from an existing codebase through targeted reads and owner questions, without inventing intent or touching authored files.
depends_on:
    - E47-onboarding-upgrades/T002-start-a-new-project-on-the-v2-workflow
complexity_tier: low
complexity_reason: Authored scaffold guidance plus its assertions; the skills and code paths it describes are unchanged.
---

# T004: Adopt Savepoint into an existing codebase

## Problem

The common case is not an empty directory. It is a codebase that already exists, where someone runs `savepoint init` and the scaffold hands the agent an empty `Design.md` describing a system that is sitting right there in the repository.

Nothing currently tells that agent what to do. `savepoint-design` describes maintaining Design for a project whose Design already says something; `savepoint-idea` may read "targeted existing-project evidence" but owns intent, not architecture. Between them is the case neither covers: the code is the only source of truth about what exists, and the owner is the only source of truth about what it is for.

The failure this guidance prevents is an agent that fills Design.md by inference — writing a confident Codebase Map from filenames, inventing an Intent the owner never stated, and producing a document that reads as reconstructed fact while being partly fiction. The rule is the same one the rest of V2 follows: read what you are told to read, record what you actually found, and ask the owner for what only they know.

This lives in scaffold content, not in the skills. The four skills are byte-parity-locked E46 contracts and are not reopened here. The V2 agent guide, the V2 `Design.md` template's Current Technical State section, and the router's opening `next_action` are what a fresh agent in an existing codebase reads, and that is where the guidance goes.

It must also stay honest about scope. No scan of the whole repository, no model service, no background analysis — reconstruction is an agent reading files under the same targeted-read discipline every other V2 role follows. And no file the user authored is changed: adoption adds Savepoint's own files and nothing else, which is T002's behavior, stated here so the guidance and the code agree.

## Context Files

- `templates/project-v2/AGENTS.md`
- `templates/project-v2/.savepoint/Design.md`
- `templates/project-v2/.savepoint/Idea.md`
- `templates/project-v2/.savepoint/router.md`
- `agent-skills/savepoint-idea/SKILL.md`
- `agent-skills/savepoint-design/SKILL.md`
- `internal/init/template_freshness_test.go`
- `internal/init/integration_test.go`
- `.savepoint/releases/v2/v2-Design.md`

## Acceptance Criteria

- [x] `templates/project-v2/AGENTS.md` carries an existing-codebase adoption section stating that Design is reconstructed from the code through targeted reads, that intent comes from the owner, and that no authored file is modified by adoption.
- [x] That section names which record receives which finding: what exists goes to `Design.md`'s Components/Codebase Map and Current Technical State; what it is for goes to `Idea.md` through the owner.
- [x] It states explicitly that reconstruction performs no whole-repository scan, no automatic analysis, and no model service call, and that unread areas are recorded as unknown rather than inferred.
- [x] `templates/project-v2/.savepoint/Design.md`'s Current Technical State section instructs the reader to record only what was verified by reading, and to name what has not been examined.
- [x] The guidance degrades when optional files are absent: it never requires a Concept, a Health-Check, a procedures file, or a release document, and absence is stated as normal rather than as a gap (TPL-03).
- [x] The four V2 skills and three shared references are byte-identical to their state before this task in both trees.
- [x] A test asserts the adoption section's presence and its three load-bearing statements — targeted reads, owner-supplied intent, no authored file modified.
- [x] An integration test initialises into a populated temporary directory and proves every pre-existing file is byte-identical afterwards, including one file whose name collides with nothing Savepoint writes.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Draft the adoption section: when it applies, what to read, where findings go, what to ask the owner, and what reconstruction explicitly does not do.
- [x] Add it to `templates/project-v2/AGENTS.md` and write the matching instruction into the Design template's Current Technical State section.
- [x] Align the router's opening `next_action` so the empty-directory and existing-codebase starts read as one route with two branches, not two workflows.
- [x] Re-read the section against `savepoint-idea` and `savepoint-design` and remove anything that restates or contradicts their owned rules.
- [x] Add the content assertions and the byte-parity check over the untouched skills.
- [x] Add the populated-directory integration test.
- [x] Run `go test ./internal/init/...`, then `make build && make test`.

## Context Log

### Per-criterion evidence

1. **Adoption section, three load-bearing statements** — `templates/project-v2/AGENTS.md` `## Existing Codebase Adoption` states reconstruction is "through targeted reads", intent is "recorded in `.savepoint/Idea.md` through `savepoint-idea`, never inferred from source", and "Adoption never modifies a file the user authored". Verified by `TestV2AgentsGuideCarriesExistingCodebaseAdoptionSection`.
2. **Finding destinations** — same section: "What exists goes to `.savepoint/Design.md`: concrete structure to Components/Codebase Map, what was actually verified to `Current Technical State`. What the code is for goes to `.savepoint/Idea.md`, through the owner." Verified by the same test.
3. **No scan/analysis/model-service; unread = unknown** — same section: "never through a whole-repository scan, an automatic analysis pass, or a call to a model service" and "An area not yet read is recorded as unknown; it is never filled in by inference." Verified by the same test.
4. **Design.md Current Technical State instruction** — `templates/project-v2/.savepoint/Design.md` now reads "Record only what was verified by reading, and name what has not yet been examined as unknown rather than leaving it implied." Verified by `TestV2DesignCurrentTechnicalStateInstructsVerifiedOnly`.
5. **Optional-file degradation (TPL-03)** — adoption section names Concept, Health-Check, procedures file, and release document as never required, closing with "Their absence is normal, not a finding." Verified by the same adoption test.
6. **Skill/reference byte-parity untouched** — no file under `agent-skills/` or either tree's `agent-skills/` was edited this task; `TestV2SkillSetIsCompleteWithByteParity` (pre-existing) still passes.
7. **Content test for adoption section** — `TestV2AgentsGuideCarriesExistingCodebaseAdoptionSection` added in `internal/init/v2_scaffold_test.go`.
8. **Populated-directory integration test** — `TestV2ScaffoldIntoPopulatedDirectoryPreservesExistingFiles` added in `internal/init/v2_scaffold_test.go`: scaffolds `templates/project-v2` into a directory pre-populated with `main.go`, `src/app.py`, `README.md`, and `requirements.txt` (the last colliding with nothing the V2 tree writes), and asserts every byte is unchanged while the V2 `.savepoint` files still land.
9. **Build/test** — `make build && make test` passes (see commands below).

### Commands run

- `go test ./internal/init/... -run 'TestV2' -v` — all pass, including the 3 new tests.
- `go test ./internal/init/... -run 'TestV2|TestSavepoint|TestShared|TestProject' -v` — no failures.
- `make build && make test` — all packages pass (`ok` for every listed package).
- `gofmt -l internal/init/v2_scaffold_test.go` — clean after `gofmt -w`.

### Files read (beyond Context Files, with reason)

- `internal/init/v2_scaffold_test.go` — the existing V2-scaffold test file; needed to place new tests consistently and reuse its fixtures (`v2ScaffoldSavepointFiles`).
- `internal/init/agent_skills_test.go` — to see the established pattern for section/phrase assertions (`sectionBody`, `v2SkillRoots`) before writing new tests in the same idiom.
- `internal/init/skill_validation_test.go` — to find where `sectionBody`, `frontmatterField`, `skillRoots`/`v2SkillRoots` are defined, confirming they're package-level helpers reusable from `v2_scaffold_test.go`.
- `internal/init/scaffold.go` — to confirm `Scaffold` never touches a path outside the template tree (governs what the populated-directory test needed to prove).
- `agent_skills_test.go` (root) — quick function-name scan, no content dependency taken.
- `.savepoint/releases/v2/epics/E47-onboarding-upgrades/E47-Detail.md` — epic detail, per router's standard read order.
- `.savepoint/router.md`, `.savepoint/Guardrails.md` — router state/next action and current guardrail rule IDs, per the standard read order.
- `.savepoint/releases/v2/epics/E47-onboarding-upgrades/tasks/T002-*.md` (header only) — confirmed the depended-on Task's `status: done` before starting.

### Files changed

- `templates/project-v2/AGENTS.md` — added `## Existing Codebase Adoption` section.
- `templates/project-v2/.savepoint/Design.md` — extended the `## Current Technical State` section with the verified-only/name-unknowns instruction.
- `templates/project-v2/.savepoint/router.md` — extended the opening `next_action` to name the existing-codebase branch.
- `internal/init/v2_scaffold_test.go` — added 3 tests: `TestV2AgentsGuideCarriesExistingCodebaseAdoptionSection`, `TestV2DesignCurrentTechnicalStateInstructsVerifiedOnly`, `TestV2ScaffoldIntoPopulatedDirectoryPreservesExistingFiles`.
- `.savepoint/router.md`, this Task file — lifecycle bookkeeping.

### Limitations

- `internal/init/template_freshness_test.go` and `internal/init/integration_test.go`, named in this Task's Context Files, were read for context (they cover the V1 tree's template-freshness and integration contracts) but not edited — the new V2-specific tests live in `internal/init/v2_scaffold_test.go` instead, alongside the rest of the V2 scaffold test suite, since that file already owns V2-tree content assertions and this keeps the V1 files' contracts undisturbed.
- No owner conversation occurred to validate the adoption section reads well to a real adopting agent; that judgment is left to Check/owner review per this skill's boundaries.
