---
id: E39-code-style-ownership/T001-point-guidance-at-style-rules
status: planned
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

- [ ] `templates/project/AGENTS.md` `## Code Style` contains a one-line pointer to the STYLE rules in `.savepoint/Guardrails.md` and no longer lists the ten rules.
- [ ] The repo's own `AGENTS.md` `## Code Style` carries the same pointer.
- [ ] `.savepoint/Guardrails.md` exists in this repository, adapted from the E37 template, and contains the `STYLE-01..10` category.
- [ ] `savepoint-build-task` (live and template) reads the STYLE rules during build and skips gracefully when `.savepoint/Guardrails.md` is absent.
- [ ] `savepoint-audit-epic` (live and template) `## Code Style Review` keeps its checkbox format but instructs one checkbox per STYLE rule found in `.savepoint/Guardrails.md`, with no hardcoded rule labels.
- [ ] When `.savepoint/Guardrails.md` is absent, the audit skeleton degrades to a stated "code style not defined for this project" note rather than an empty or invented checklist.
- [ ] Code style remains Guideline severity and advisory; no skill treats a STYLE rule as a blocker.
- [ ] `templates/project/.savepoint/Design.md` audit-pipeline notes still describe `## Code Style Review` as a required audit-file section.
- [ ] No live catalog, skill, template, or guide restates the ten rules.
- [ ] Canonical and generated skill copies remain byte-identical (`agent_skills_test.go:28-54`).
- [ ] `make build && make test` passes.

## Implementation Plan

- [ ] Create `.savepoint/Guardrails.md` for this repository from the E37 template, keeping the STYLE category and trimming rules that do not apply to a Go CLI.
- [ ] Replace the `## Code Style` body in both AGENTS.md copies with the pointer line.
- [ ] Add the STYLE-rule read to build-task (both copies) alongside the existing Health-Check Quick reference.
- [ ] Rewrite the audit skeleton's `## Code Style Review` to source checkboxes from STYLE rule IDs, including the absent-file fallback.
- [ ] Verify canonical/template parity and that no file restates the rules.
- [ ] Run `make build && make test`.

## Context Log

Pending.
