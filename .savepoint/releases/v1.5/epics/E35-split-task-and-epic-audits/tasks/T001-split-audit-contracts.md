---
id: E35-split-task-and-epic-audits/T001-split-audit-contracts
status: done
objective: Replace the generic audit skill with distinct task and epic contracts backed by one shared audit method with enriched rigor.
depends_on: []
complexity_tier: high
complexity_reason: Replaces a cross-workflow contract while preserving audit authority and closeout semantics, and adds scope locks, coverage matrices, and convergence limits.
---

# T001: Split Audit Contracts

## Problem

The generic audit skill conflates read-only review of one in-progress task with artifact-producing audit and closeout of a completed epic, and lacks rigorous audit machinery (scope locks, coverage matrices, materiality tables).

## Context Files

- `agent-skills/savepoint-audit/SKILL.md`
- `agent-skills/savepoint-audit-register/SKILL.md`
- `agent-skills/savepoint-build-task/SKILL.md`
- `agent-skills/savepoint-audit-task/SKILL.md`
- `agent-skills/savepoint-audit-epic/SKILL.md`
- `agent-skills/references/audit-method.md`
- `templates/project/agent-skills/savepoint-audit/SKILL.md`
- `templates/project/agent-skills/savepoint-audit-register/SKILL.md`
- `templates/project/agent-skills/savepoint-build-task/SKILL.md`
- `templates/project/agent-skills/savepoint-audit-task/SKILL.md`
- `templates/project/agent-skills/savepoint-audit-epic/SKILL.md`
- `templates/project/agent-skills/references/audit-method.md`

## Acceptance Criteria

### Structural split (existing)
- [x] `savepoint-audit-task` triggers only for an explicit audit or re-audit of one in-progress task, keeps router state `task-building`, performs no writes, uses the Quick health check, and returns only `CLEAR` or `NEEDS WORK`.
- [x] Task-audit guidance explicitly forbids `E##-Audit.md` creation and task/router lifecycle changes.
- [x] `savepoint-audit-epic` triggers for `audit-pending` or an explicit completed-epic audit, requires a session independent from the builder, audits every completed task, and uses the Full health check.
- [x] Epic-audit guidance reviews file reality, drift, and release guardrails; writes exactly one `E##-Audit.md`; and preserves proposal approval, apply, and closeout rules.
- [x] Both skills use one shared reference requiring general-invariant AC reasoning; normal, boundary, malformed, failure, and bypass scenarios; public-entry-point and alternate-representation checks; and applicable behavior matrices.
- [x] The shared method requires full re-audit after remediation, treats passing tests as supporting evidence, classifies every criterion as Proven/Finding/Unverified, and returns `NEEDS WORK` for findings or material unverified criteria.
- [x] Consolidated findings include the violated rule, smallest reproduction, expected and actual behavior, exact file evidence, and missing test evidence.
- [x] `savepoint-build-task` explicitly defers audit and re-audit requests to `savepoint-audit-task`.
- [x] The live and generated generic skill folders are removed, and all changed canonical skills match their generated copies.

### Enriched audit rigor (new)
- [x] The shared method requires a frozen scope lock (numbered, listing AC, guardrails, gates, changed files, entry points, dependencies, matrix axes, materiality boundary) before adversarial probes; re-audit reuses it immutably without adding axes.
- [x] The shared method requires a mandatory coverage matrix with named axes: public surfaces, input shape, state, environment/output, boundaries, sequences, representations, and text classes; a finite external-boundary matrix for server/dependency code; prose checklist is not matrix evidence.
- [x] The shared method includes a Workflow And Side-Effect Audit Lock for multi-step or side-effecting work: per-operation failure-timing table, independent oracle, and matrix completion lock before verdict.
- [x] The shared method defines a default convergence limit: initial audit, one full re-audit, one targeted remediation, then stop.
- [x] The shared method requires an admission ledger for re-audit: every blocking result maps to an exact frozen matrix cell; credible-blocker exception for secret exposure, cross-tenant access, destructive data loss, or sensitive-data harm.
- [x] Both skills require file reality evidence: every file named in task logs exists, was intentionally deleted, or is recorded as discarded scratch work.
- [x] Epic audit produces a materiality table (Finding × Likelihood × Impact × Materiality × Recommendation) for every finding.
- [x] Epic audit includes a named Guardrails Verification subsection with Full health check.
- [x] Epic audit produces a repository handoff result: CLEAR TO COMMIT/PUSH or NOT READY TO COMMIT/PUSH.
- [x] Epic audit Main Findings defaults to 500-900 words; task audit defaults to 350-600 words.
- [x] Both skills return a Final Response Output: compact chat summary with verdict, materiality table (if NEEDS WORK), gate result, and audit file link.
- [x] Both skills reference `.savepoint/Guardrails.md` and `.savepoint/Health-Check.md` gracefully (skip when absent).
- [x] `savepoint-build-task` references Health-Check.md Quick mode at task handoff (skip when absent).

### Source-material adaptation
- [x] Source-material paths are adapted: `../shared/savepoint-audit-method.md` becomes the real relative path to `agent-skills/references/audit-method.md` in both skills' Read sections, and root `GUARDRAILS.md` references become `.savepoint/Guardrails.md`.
- [x] "release guardrails audit plan" references are softened to: "the release's guardrails mapping, if your project maintains one — otherwise the relevant `.savepoint/Guardrails.md` rule IDs directly."

## Implementation Plan

- [x] Add the shared non-triggerable audit-method reference (scope locks, coverage matrices, side-effect locks, adversarial pass, re-audit convergence, materiality) and mirror it into project templates.
- [x] Adapt source-material paths and terminology: shared-method read path, `.savepoint/Guardrails.md` paths, softened audit-plan wording.
- [x] Add the task-audit skill with strict read-only scope, Quick health check, evidence classification, materiality table, and result contract.
- [x] Add the epic-audit skill with Full health check, full-epic coverage, single artifact, frozen scope lock, coverage matrix, proposal approval, guardrails verification, repository handoff, and closeout workflow.
- [x] Update audit-register integration to name the correct split audit owner without becoming a generic alias.
- [x] Update build-task guidance to hand explicit audit and re-audit requests to task audit and reference Health-Check.md Quick mode.
- [x] Remove both copies of the old generic skill and verify canonical/template parity.

## Context Log

**Read:** `.savepoint/router.md`, `AGENTS.md`, `E35-Detail.md`, T002/T003/T004 (for
contract boundaries only), `examples.md` (source material), the four Context Files
that existed, `internal/init/template_freshness_test.go`, `agent_skills_test.go`,
`internal/init/scaffold.go`, `internal/init/upgrade.go`.

**Added:**

- `agent-skills/references/audit-method.md` — non-triggerable shared method
  (frontmatter `type: audit-method-reference`, `triggerable: false`, no `name:`
  key so no skill discovery). Covers scope establishment, frozen scope lock,
  acceptance-as-invariants, mandatory coverage matrix (8 named axes) plus the
  finite external-boundary matrix, Workflow And Side-Effect Audit Lock with the
  matrix completion lock, adversarial pass, re-audit with admission ledger and
  the initial + one re-audit + one targeted remediation convergence limit, file
  reality, gates, Proven/Finding/Unverified classification, and materiality.
- `agent-skills/savepoint-audit-task/SKILL.md` — explicit-request-only trigger,
  router stays `task-building`, Quick health check, no writes at all, no
  `E##-Audit.md`, result limited to `CLEAR`/`NEEDS WORK`, 350–600 word default,
  Final Response Output.
- `agent-skills/savepoint-audit-epic/SKILL.md` — `audit-pending` or explicit
  completed-epic trigger, independent-session stop rule, every completed task,
  Full health check, file reality, drift reconciliation, audit-register handoff,
  exactly one `E##-Audit.md`, materiality summary, named Guardrails Verification
  subsection, repository handoff (`CLEAR TO COMMIT/PUSH` /
  `NOT READY TO COMMIT/PUSH`), 500–900 word default, proposal approval and
  apply/close rules preserved verbatim in behaviour.
- Template mirrors of all three under `templates/project/agent-skills/`.

**Edited:**

- `agent-skills/savepoint-build-task/SKILL.md` (+ template copy) — Trigger and
  Rules defer explicit audit/re-audit to `savepoint-audit-task` (or
  `savepoint-audit-epic` for a completed epic); workflow step 8 applies the
  `.savepoint/Health-Check.md` Quick check at handoff and skips when absent.
- `agent-skills/savepoint-audit-register/SKILL.md` (+ template copy) — names
  `savepoint-audit-epic` as the handoff owner and states that a read-only
  `savepoint-audit-task` review writes no register records.
- `internal/init/template_freshness_test.go` — audit-skill assertions now target
  `savepoint-audit-epic`; the scaffold count check filters to `savepoint-*`
  directories so the new `references/` folder does not fail it; added canonical/
  template parity for `agent-skills/references/audit-method.md`.

**Removed:** `agent-skills/savepoint-audit/` and
`templates/project/agent-skills/savepoint-audit/` (via `git rm`).

**Source-material adaptation:** `../shared/savepoint-audit-method.md` →
`agent-skills/references/audit-method.md`; root `GUARDRAILS.md` →
`.savepoint/Guardrails.md`; "release guardrails audit plan" → "the release's
guardrails mapping, if your project maintains one — otherwise the relevant
`.savepoint/Guardrails.md` rule IDs directly"; "child-data" softened to
"sensitive-data / privacy" so the method stays project-agnostic.

**Verification:**

- `make build && make test` — pass (all packages; `internal/init` 0.802s).
- Fresh scaffold into a scratch dir with a locally built binary produced
  `agent-skills/references/audit-method.md`, `savepoint-audit-task/SKILL.md`,
  `savepoint-audit-epic/SKILL.md` and no `savepoint-audit/` folder.
- `TestScaffoldedSavepointSkillsMatchBundledSkills` and
  `TestBundledSavepointSkillsHaveDiscoveryFrontmatter` pass over the new folders;
  canonical/template parity confirmed for both split skills, the shared
  reference, build-task, and audit-register.

**Health Check:** skipped — this repo has no `.savepoint/Health-Check.md`; the
file exists only as `templates/project/.savepoint/Health-Check.md` for generated
projects.

**Deliberately out of scope (owned by later E35 tasks):** `AGENTS.md`,
`templates/project/AGENTS.md`, `templates/project/.savepoint/router.md`,
both `Design.md` files, and `README.md` still name the removed generic skill —
T002 owns that routing and documentation sweep. `internal/init/upgrade.go` still
copies only `*/SKILL.md` and has no legacy-skill migration — T003. Broader
scaffold, contract, validator, and stale-reference tests — T004.

## Drift Notes

- New directory `agent-skills/references/` (and its generated mirror) holds
  non-triggerable shared skill references. The `AGENTS.md` Codebase Map row for
  `agent-skills/` still reads "Phase-specific skill guides, including defect
  capture guidance" and the Skill Activation table still maps `audit-pending` to
  the removed `savepoint-audit`. Both are updated by T002, which owns generated
  routing and guidance; recorded here so the gap is not lost if T002 is rescoped.
- `.savepoint/Design.md` line 130 and `templates/project/.savepoint/Design.md`
  line 120 still reference `agent-skills/savepoint-audit/SKILL.md` as the audit
  pipeline. Also T002.
