---
id: E39-code-style-ownership/T001-point-guidance-at-style-rules
status: done
objective: Replace the restated code-style rules in AGENTS.md, build-task, and the audit skeleton with references to the STYLE rules in .savepoint/Guardrails.md.
depends_on: []
complexity_tier: medium
complexity_reason: Spans three file pairs plus the repo guide, and must keep the audit checklist renderable while removing the hardcoded rules.
---

# T001: Point Guidance at the STYLE Rules

## Problem

The ten code-style rules are defined in the managed AGENTS.md block and restated as a hardcoded checklist in the audit skeleton. The managed block is overwritten on every upgrade, so project tailoring is lost, and the restatement violates the Guardrails template's own POL-01 and POL-02 rules. No builder skill references code style at all, so the agent writing the code never sees the rules.

## Context Files

- `AGENTS.md`
- `templates/project/AGENTS.md`
- `agent-skills/savepoint-build-task/SKILL.md`
- `templates/project/agent-skills/savepoint-build-task/SKILL.md`
- `agent-skills/savepoint-audit-epic/SKILL.md`
- `templates/project/agent-skills/savepoint-audit-epic/SKILL.md`
- `templates/project/.savepoint/Guardrails.md`
- `templates/project/.savepoint/Design.md`

## Acceptance Criteria

- [x] `templates/project/AGENTS.md` `## Code Style` contains a one-line pointer to the STYLE rules in `.savepoint/Guardrails.md` and no longer lists the ten rules.
- [x] The repo's own `AGENTS.md` `## Code Style` carries the same pointer.
- [x] `.savepoint/Guardrails.md` exists in this repository, adapted from the E37 template, and contains the `STYLE-01..10` category.
- [x] `savepoint-build-task` (live and template) reads the STYLE rules during build and skips gracefully when `.savepoint/Guardrails.md` is absent.
- [x] `savepoint-audit-epic` (live and template) `## Code Style Review` keeps its checkbox format but instructs one checkbox per STYLE rule found in `.savepoint/Guardrails.md`, with no hardcoded rule labels.
- [x] When `.savepoint/Guardrails.md` is absent, the audit skeleton degrades to a stated "code style not defined for this project" note rather than an empty or invented checklist.
- [x] Code style remains Guideline severity and advisory; no skill treats a STYLE rule as a blocker.
- [x] `templates/project/.savepoint/Design.md` audit-pipeline notes still describe `## Code Style Review` as a required audit-file section.
- [x] No live catalog, skill, template, or guide restates the ten rules.
- [x] Canonical and generated skill copies remain byte-identical (`agent_skills_test.go:28-54`).
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Create `.savepoint/Guardrails.md` for this repository from the E37 template, keeping the STYLE category and trimming rules that do not apply to a Go CLI.
- [x] Replace the `## Code Style` body in both AGENTS.md copies with the pointer line.
- [x] Add the STYLE-rule read to build-task (both copies) alongside the existing Health-Check Quick reference.
- [x] Rewrite the audit skeleton's `## Code Style Review` to source checkboxes from STYLE rule IDs, including the absent-file fallback.
- [x] Verify canonical/template parity and that no file restates the rules.
- [x] Run `make build && make test`.

## Context Log

**Files read:** `.savepoint/router.md`, `E39-Detail.md`, this task file, `AGENTS.md`, `templates/project/AGENTS.md`, `agent-skills/savepoint-build-task/SKILL.md`, `agent-skills/savepoint-audit-epic/SKILL.md` (both template copies verified byte-identical before editing), `templates/project/.savepoint/Guardrails.md`, `templates/project/.savepoint/Design.md`, `.savepoint/Design.md`, `agent_skills_test.go`, `examples.md` (targeted verification read).

**Files created:**

- `.savepoint/Guardrails.md` — Savepoint's own policy file, adapted from the E37 template. Kept the severity model, `STYLE-01..10` verbatim, and `POL-01/02`. Replaced the template's web/SaaS categories (auth, billing, privacy retention, jobs, migrations, frontend, LLM) with categories that match a local single-user Go CLI: filesystem safety (`FS-01..06`), planning-data integrity (`DATA-01..05`), templates and shipped assets (`TPL-01..04`), architecture, configuration/dependencies, testing (`TEST-01..08`), and release/distribution (`REL-01..03`).

**Files edited:**

- `AGENTS.md` and `templates/project/AGENTS.md` — `## Code Style` body replaced with a pointer to the `STYLE` rules in `.savepoint/Guardrails.md`, stating Guideline severity and the absent-file fallback.
- `agent-skills/savepoint-build-task/SKILL.md` (+ template copy) — added the `STYLE` rules to `## Read`, folded them into workflow step 5 with the skip-when-absent clause, and added two rules: STYLE is advisory and never blocks handoff, and do not restate the rules.
- `agent-skills/savepoint-audit-epic/SKILL.md` (+ template copy) — `## Code Style Review` output instruction now says one checkbox per `STYLE` rule found in `.savepoint/Guardrails.md`, labelled by rule ID in file order, with the "Code style is not defined for this project." fallback and an explicit statement that an unchecked box never causes `NEEDS WORK`. The skeleton's ten hardcoded labels were replaced with the sourced-checkbox pattern, and a new rule forbids hardcoding rule labels.
- `.savepoint/Design.md:127` — the audit-pipeline note said `## Code Style Review` "contains the 10 AGENTS.md code style checks", which is no longer true; it now points at the `STYLE` rules in `.savepoint/Guardrails.md`. See `## Drift Notes`.

**Verification:**

- Canonical/template parity: `diff` clean for both edited skills before and after; `TestScaffoldedSavepointSkillsMatchBundledSkills` passes.
- Restatement sweep: `grep -rln "One job per file" . --exclude-dir=.git --exclude-dir=.savepoint` leaves only `examples.md`, `internal/board/epic_panel_test.go`, and `project-audit/audit_report_fable_5.md`. None is a live catalog, skill, template, or guide: `examples.md` is the user-owned source reference carrying its own "DO NOT DELETE OR AMEND WITHOUT USER INPUT" banner, the test file uses the labels as render fixtures, and `project-audit/` is a historical report. Historical `E##-Audit.md` records under `.savepoint/releases/` are out of scope per the epic boundaries.
- `templates/project/.savepoint/Design.md:118-119` still names `## Code Style Review` as a required user-facing audit-file section — unchanged, no edit needed.

**Quality gates:** `make build && make test` — pass. All nine packages ok; `internal/init` 1.194s, `internal/board` cached.

**Health check:** skipped — this repository has no `.savepoint/Health-Check.md`.

## Drift Notes

- `.savepoint/Guardrails.md` is new in this repository. It is project policy data, not a code module, so the AGENTS.md Codebase Map is unchanged. T002 adds the upgrade path that delivers this file (and `Health-Check.md`) to existing projects.
- `.savepoint/Design.md:127` was corrected because the audit pipeline's Code Style Review source moved from the AGENTS.md managed block to `.savepoint/Guardrails.md`. This records where the section's content now comes from; the audit pipeline itself, its section list, and the advisory (non-blocking) Layer 2 posture are unchanged.
