---
id: E35-split-task-and-epic-audits/T002-generated-routing-and-guidance
status: done
objective: Route generated projects to the correct split audit skill, remove generic audit terminology from live guidance, and update cross-references across router, AGENTS, Design, audit-register, build-task, and README.
depends_on:
  - E35-split-task-and-epic-audits/T001-split-audit-contracts
complexity_tier: medium
complexity_reason: Router and catalog changes must align request intent without adding a new workflow state; cross-reference updates span five file pairs.
---

# T002: Generated Routing and Guidance

## Problem

Generated router, contributor, and user guidance still presents audit as one ambiguous phase skill and does not describe the task-audit override. Cross-references to the old generic skill remain in AGENTS.md, Design.md, audit-register, build-task, and README.

## Context Files

- `AGENTS.md`
- `templates/project/AGENTS.md`
- `templates/project/.savepoint/router.md`
- `templates/project/.savepoint/Design.md`
- `README.md`
- `agent-skills/savepoint-audit-task/SKILL.md`
- `agent-skills/savepoint-audit-epic/SKILL.md`
- `agent-skills/savepoint-build-task/SKILL.md`
- `agent-skills/savepoint-audit-register/SKILL.md`
- `templates/project/agent-skills/savepoint-audit-register/SKILL.md`
- `.savepoint/Design.md`

## Acceptance Criteria

- [x] The live and generated skill tables map `audit-pending` to `savepoint-audit-epic`.
- [x] Generated routing maps `task-building` plus an explicit task audit or re-audit request to `savepoint-audit-task` as a request-qualified override, not a router state.
- [x] Router template `## Manual Overrides` includes a task-audit override line next to the existing epic-audit override.
- [x] AGENTS.md (live and template) Audit section describes the split skills and references `.savepoint/Guardrails.md` and `.savepoint/Health-Check.md` when present.
- [x] `templates/project/.savepoint/Design.md` audit pipeline reference updates `savepoint-audit` to `savepoint-audit-epic`.
- [x] `agent-skills/savepoint-audit-register/SKILL.md` (live and template) lines 14, 39 reference `savepoint-audit-epic` instead of `savepoint-audit`.
- [x] `agent-skills/savepoint-build-task/SKILL.md` (live and template) references Health-Check.md Quick mode at task handoff (graceful when absent).
- [x] `.savepoint/Design.md` (repo's own) updates bundled skills list (line 17) and audit pipeline reference (line 130) to `savepoint-audit-epic`.
- [x] README documents the two audit intents, shared method, enriched rigor (scope locks, matrices, materiality), and existing-project compatibility behavior.
- [x] No live catalog, router template, generated guide, or user documentation exposes `savepoint-audit` as a triggerable skill or alias.
- [x] Historical release records are left unchanged.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Replace the generic audit row in live and generated skill activation tables.
- [x] Add the request-qualified task-audit override to the generated router without adding state.
- [x] Update AGENTS.md Audit section to name split skills and add Guardrails/Health-Check discoverability notes.
- [x] Update Design.md template audit reference to `savepoint-audit-epic`.
- [x] Update audit-register skill (both copies) lines 14, 39.
- [x] Update build-task skill (both copies) for audit-deferral and Health-Check Quick mode.
- [x] Update repo's own Design.md bundled skills list (line 17) and audit reference (line 130).
- [x] Update README skill catalog and audit workflow documentation.
- [x] Review live guidance for ambiguous generic references while excluding historical release artifacts.

## Context Log

**Read:** all `## Context Files`, plus `README.md` audit/skills/upgrade sections and
`internal/init/template_freshness_test.go` (to keep asserted strings intact).

**Edited:**

- `AGENTS.md` and `templates/project/AGENTS.md` — Skill Activation maps
  `audit-pending` to `savepoint-audit-epic`; a note under the table states that an
  explicit task audit/re-audit request selects `savepoint-audit-task` while
  `state` stays `task-building` (request-qualified override, not a state). The
  `## Audit` section now describes both skills, names the shared method at
  `agent-skills/references/audit-method.md`, and documents `.savepoint/Guardrails.md`
  (policy) and `.savepoint/Health-Check.md` (Quick at task handoff and task audit,
  Full at epic audit) as optional — absence is not a finding. The existing audit
  file path and apply/close lines are unchanged.
- `templates/project/.savepoint/router.md` — table row updated; added a note that
  `savepoint-audit-task` deliberately has no row; `## Manual Overrides` now carries
  an epic-audit line and a new task-audit line stating `state: task-building` is
  left unchanged and the review writes no audit file.
- `templates/project/.savepoint/Design.md` and `.savepoint/Design.md` — audit
  pipeline reference split between `savepoint-audit-epic` and `savepoint-audit-task`,
  both pointing at the shared method.
- `.savepoint/Design.md` — bundled skills list replaces `savepoint-audit` with
  `savepoint-audit-task` and `savepoint-audit-epic`, and records that the shared
  reference is a reference rather than a triggerable skill.
- `README.md` — new `## Audit` section with an intent/trigger/writes/health-check/
  result table for the two skills, the shared method, and the enriched rigor
  (frozen scope lock, mandatory coverage matrix, workflow and side-effect lock,
  convergence limit with admission ledger and credible-blocker exception,
  materiality table), plus graceful Guardrails/Health-Check behavior. Agent Skills
  list updated and a note added for `agent-skills/references/`. `## Updating
  Existing Projects` documents that upgrade installs the split skills and shared
  reference, archives legacy `savepoint-audit` content under a non-triggerable
  `.savepoint/migrations/` archive before removing the active folder, and creates
  no archive for projects that never had it.

**Already satisfied by T001 (verified, not re-edited):** AC for
`savepoint-audit-register` lines 14/39 (both copies name `savepoint-audit-epic`)
and AC for `savepoint-build-task` Health-Check Quick mode at handoff (both copies).

**Verification:**

- `make build && make test` — pass (all packages; `internal/init` 0.933s).
- Stale-name sweep over `*.go`/`*.md`/`*.yml` excluding `.savepoint/releases/`:
  no remaining `savepoint-audit` reference in any catalog, router template,
  generated guide, or user documentation. Remaining hits are (a) historical
  release records, left unchanged as required; (b) `internal/init/upgrade_test.go`
  and `integration_test.go` in-memory fixtures, which T003 owns; (c) `examples.md`
  — see Drift Notes.
- `*.js`/`*.json`/`*.sh` sweep: no references.

**Health Check:** skipped — this repo has no `.savepoint/Health-Check.md`; the file
exists only as `templates/project/.savepoint/Health-Check.md` for generated projects.

## Drift Notes

- `examples.md` at the repo root is transient source material (the imported
  health-check ceremony, `GUARDRAILS.md`, both draft skills, and the draft shared
  method). It still contains `../shared/savepoint-audit-method.md` and root
  `GUARDRAILS.md` paths that were deliberately adapted away in T001. It is not a
  catalog, router, guide, or user documentation, so it does not violate this
  task's AC — but a naive substring check would flag it. T004 must either exclude
  `examples.md` alongside `.savepoint/releases/`, or the file should be deleted
  once E35 finishes consuming it. Recorded here so the decision is explicit rather
  than accidental.
- `README.md` `## Updating Existing Projects` now describes the
  `.savepoint/migrations/` archive behavior that T003 implements. If T003 is
  rescoped or its archive path changes, that README paragraph must change with it.
