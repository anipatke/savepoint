---
id: E21-document-template-optimisation/T003-design-md-template-optimisation
status: done
objective: Update the Design.md template so its sections and prompts match the current directory layout, architecture model, and audit workflow.
complexity_tier: low
complexity_reason: Single template file; section headings and prompts need updating, no code changes.
---

# T003: Design.md Template Optimisation

## Problem

The `templates/project/.savepoint/Design.md` template has nine numbered sections, but their headings and prompts have not been updated since the directory layout and audit workflow evolved. A scaffolded project receives a template whose section order and language may not match what the `savepoint-system-design` skill produces or what the audit loop expects.

## Context Files

- `templates/project/.savepoint/Design.md`
- `.savepoint/Design.md`
- `templates/project/agent-skills/savepoint-system-design/SKILL.md`
- `internal/init/template_freshness_test.go`

## Acceptance Criteria

- [x] Section headings in the template match the structure the `savepoint-system-design` skill is expected to produce for a new project.
- [x] Section prompts are concrete and match the specificity of the live `.savepoint/Design.md`.
- [x] The audit workflow section (§7) accurately describes the skill-driven audit loop, not any CLI pipeline.
- [x] Directory layout section (§2) matches the canonical directory structure in the live Design.md.
- [x] Template freshness test passes.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Read `templates/project/.savepoint/Design.md` and the live `.savepoint/Design.md` side by side.
- [x] Read `templates/project/agent-skills/savepoint-system-design/SKILL.md` to understand what the skill writes into this file.
- [x] Update section headings and prompts to align with the current architecture model and audit workflow.
- [x] Run `go test ./internal/init` to check freshness coverage.
- [x] Run `make build && make test`.

## Context Log

### Divergences found

1. **Section count.** The template had 9 sections; the live Design.md has 14. Sections 8–14 had been collapsed, renamed, or dropped:
   - Template §8 = "Testing strategy" → live §8 = "TUI", live §13 = "Testing".
   - Live has additional §9 "Concurrency", §10 "Release versioning (PRDs)", §11 "Failure modes", §12 "Distribution & build", §14 "Package versioning" — none of these existed in the template.
2. **Audit workflow (§7) was directionally right but underspecified.** The current text said "Savepoint audit is agent-led and skill-driven" but lacked the numbered 0–5 workflow steps, the `E##-Audit.md` user-facing/admin section split, the apply/close rule, the "no CLI pipeline" explicit prohibition, and the three-layer model.
3. **Directory layout (§2) was a single placeholder.** The template just said `<!-- Expected file structure. -->` — no tree, no placement rules, no AGENTS.md casing note, no inline-subtask rule.
4. **Hierarchy semantics (§3) was a single placeholder.** The table shape (Level / Definition) was missing; the live uses a 5-row table with bold terms and ownership lines.
5. **Status model & gates (§4) was a single placeholder.** No table, no `blocked` flag note, no `internal/data` ownership bullet, no defect-lifecycle callout.
6. **CLI surface (§6) was a single placeholder.** No command table, no rejected-commands bullet, no naming-convention line.
7. **Architecture model (§1) was a single placeholder.** No format constraint, no token-efficiency principle anchor.
8. **Section 9 (template) was a single placeholder for "Release versioning"** — the live splits this into §10 (release versioning) and §14 (package versioning) with different audiences.

### Changes applied

- `templates/project/.savepoint/Design.md` rewritten end-to-end. New structure mirrors the live Design.md exactly: 14 sections in the same order, same headings, same role per section.
- Each section opens with a pointed prompt (the "ask") and is followed by a blockquote example modeled on the live Design.md's voice. The author's job is to delete the example and write their own.
- §1 Architecture model — 4-bullet rule, token-efficiency principle anchored.
- §2 Directory layout — fenced tree copied from live Design.md, plus 2 placement rules.
- §3 Hierarchy semantics — 5-row table with bold terms and ownership lines.
- §4 Status model & gates — 3-row status table with explicit "Who may set it" column, plus 5 ownership bullets covering `blocked`, `internal/data`, agent/user boundary, defect lifecycle, and verification mode pointer.
- §5 Dependencies — 3 bullets covering declaration, doctor checks, cross-epic warnings.
- §6 CLI surface — 5-row command table, explicit "rejected" bullet, naming convention.
- §7 Agent audit workflow — 0–5 numbered steps, `E##-Audit.md` section rules, apply/close rule, explicit "no CLI pipeline" line, three-layer model.
- §8 TUI — theming, render fallbacks, layout, keybindings (skipped in projects without a TUI).
- §9 Concurrency — mtime optimistic-concurrency mechanism, no lockfile.
- §10 Release versioning (PRDs) — integer scheme, doctor warn behavior.
- §11 Failure modes — 8-row diagnosable-failures table.
- §12 Distribution & build — license, runtime, build, cross-platform, artifacts, smoke, no telemetry.
- §13 Testing — 6-row layer table, 70% coverage target, behavior coverage priority.
- §14 Package versioning — 4-milestone semver list with pre-1.0 breakability note.
- Top-of-file italic line added: "Fill each section with concrete, falsifiable claims. Generic architecture docs produce generic agents."

### Verification

- `go test -count=1 ./internal/init/...` passes. No canonical-string assertions on Design.md content; integration and scaffold tests verify `{{PROJECT_NAME}}` interpolation only, which the template still uses.
- `make build && make test` passes across all packages.
