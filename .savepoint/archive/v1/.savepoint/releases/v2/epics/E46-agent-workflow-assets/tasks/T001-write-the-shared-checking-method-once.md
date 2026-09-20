---
id: E46-agent-workflow-assets/T001-write-the-shared-checking-method-once
title: Write the shared checking method once
status: done
objective: Add the non-triggerable V2 check method reference and its shipped copy, carrying the audit method's rigor without register bookkeeping.
depends_on: []
complexity_tier: medium
complexity_reason: Distilling 390 lines of audit method into a smaller V2 reference is judgement work, but it touches two mirrored files.
---

# T001: Write the shared checking method once

## Problem

The four V2 skills are supposed to be thin contracts. The rigor that makes a Check worth anything — freezing scope before probing, turning acceptance criteria into invariants, proving coverage instead of asserting it, the adversarial pass, materiality, bounded re-audit convergence — cannot live in four copies, and it cannot live in the public skill surface without making that surface enormous.

`agent-skills/references/audit-method.md` already holds that method for V1, but it is written around machinery V2 deletes: the audit register, `F###` finding IDs, immutable audit runs, `.savepoint/audit/` bookkeeping, and epic-shaped scope. V2 replaces all of that with Check records and Issues, which are derived and listed rather than separately maintained.

This task produces the V2 shared method as its own reference so the later skills can point at it instead of restating it. The V1 method stays exactly where it is and exactly as it is: the V1 audit skills still run this repository until cutover.

## Context Files

- `agent-skills/references/audit-method.md`
- `agent-skills/references/check-method.md`
- `templates/project/agent-skills/references/check-method.md`
- `internal/init/skill_validation_test.go`
- `internal/init/agent_skills_test.go`
- `internal/init/template_freshness_test.go`
- `internal/init/upgrade.go`
- `internal/init/manifest.go`
- `.savepoint/releases/v2/v2-Design.md`

## Acceptance Criteria

- [x] `agent-skills/references/check-method.md` exists and `templates/project/agent-skills/references/check-method.md` is byte-identical to it.
- [x] Its frontmatter is `type: check-method-reference` and `triggerable: false`, and carries no `name` key, so nothing can discover or trigger it as a skill.
- [x] It carries forward, in V2 vocabulary: scope freezing before the first probe, acceptance-to-invariant conversion, the coverage matrix requirement, the workflow and side-effect lock, the adversarial pass, materiality, and bounded re-audit convergence.
- [x] It states the depth difference explicitly: a Task Check is focused on one Task's outcome and evidence; an Objective Check adds integration across the Objective's Tasks and Design reconciliation.
- [x] It defines the Quick and Full evidence modes in one place and states that absent optional project files (`Guardrails.md`, a project verification procedure) are skipped steps, never findings.
- [x] It contains no register vocabulary: no `.savepoint/audit/`, no `register`, no `F###`, no `audit run`, and no finding-state list.
- [x] It states that it is loaded by `savepoint-check` and referenced by `savepoint-task` and `savepoint-design`, and that it never triggers on its own.
- [x] `agent-skills/references/audit-method.md` and its shipped copy are unchanged.
- [x] `internal/init/agent_skills_test.go` asserts the frontmatter contract, the presence of each named method section, the absence of register vocabulary, and live/template byte parity, in both trees.
- [x] The existing skill-count parity test in `internal/init/template_freshness_test.go` still passes unchanged, because references live outside `savepoint-*` directories.

## Implementation Plan

- [x] Read `audit-method.md` in full and list which sections are method (carry forward), which are V1 register machinery (drop), and which need V2 wording (Check, Issue, Objective, Task).
- [x] Write `agent-skills/references/check-method.md` with the reference frontmatter and the carried-forward sections in V2 vocabulary.
- [x] Add the Task Check versus Objective Check depth section and the Quick/Full evidence-mode section.
- [x] Copy the file verbatim to `templates/project/agent-skills/references/check-method.md`.
- [x] Create `internal/init/agent_skills_test.go` with the reference contract test, reusing `frontmatterField`, `sectionBody`, and `skillRoots` from `skill_validation_test.go`.
- [x] Record in `## Drift Notes` that the new V2 test file borrows those three helpers from the V1 `skill_validation_test.go`, so E50 moves them before deleting that file rather than discovering the coupling at cutover.
- [x] Confirm `internal/init/upgrade.go` and `manifest.go` already treat `agent-skills/references/` as package-owned so the new file ships and refreshes; record the finding rather than changing ownership rules.
- [x] Run `go test ./internal/init/...`, then `make build && make test`.

## Context Log

- Read: `agent-skills/references/audit-method.md` (full, unchanged), `.savepoint/releases/v2/v2-Design.md` sections 4-8 (Task lifecycle, Check records, Objective/Issue contracts, skill/verification model), `internal/init/skill_validation_test.go` (reused `frontmatterField`, `sectionBody`, `skillRoots`), `internal/init/template_freshness_test.go` (confirmed skill-count parity test scopes only `savepoint-*` dirs), `internal/init/upgrade.go` and `manifest.go` (confirmed `isPackageSkillAsset`/`isManifestPath` already treat all of `agent-skills/references/` as package-owned by prefix, not by an explicit `check-method.md` allowlist entry — no code change needed).
- Wrote: `agent-skills/references/check-method.md` (new, V2 shared check method), `templates/project/agent-skills/references/check-method.md` (verbatim copy), `internal/init/agent_skills_test.go` (new: frontmatter contract, required-section presence, register-vocabulary absence, live/template byte parity).
- Quality gates: `go test ./internal/init/...` → ok. `make build && make test` → ok (all packages pass).
- Finding (non-blocking, recorded per Implementation Plan): `isPackageSkillAsset` and `isManifestPath` in `internal/init/{upgrade,manifest}.go` already key off the `agent-skills/references/` prefix generically, so `check-method.md` ships and refreshes on upgrade with no ownership-rule change required.

## Drift Notes

`internal/init/agent_skills_test.go` is a new V2 test file that reuses `frontmatterField`, `sectionBody`, and `skillRoots` from the V1 `internal/init/skill_validation_test.go` (same package, called directly, no new export). This is intentional coupling per the task's Implementation Plan: E50 must move those three helpers to a shared location before it deletes `skill_validation_test.go` during V1 cutover, or this file's build breaks.
