---
type: guardrails
status: active
last_audited: never
---

# Guardrails — Engineering Policy

## Purpose

This document is the authoritative engineering policy for Savepoint.

If a rule can block work, fail a health check, or require remediation, it must be defined here. Skills, hooks, scripts, and sub-agents may reference these rules, but must not invent blocking policy.

Savepoint defines when evidence is required. This file defines what must be protected.

Savepoint is a local, single-user Go CLI and TUI over a project's markdown planning files. It has no server, no database, no authentication, no billing, and no external network calls. The rule set below is scoped to that reality: the real risks are damaging a user's files, shipping broken cross-platform binaries, and drifting the shipped templates out of sync with the canonical ones.

## Severity Model

| Severity | Meaning |
|---|---|
| Blocker | Must be fixed before task handoff, approval, or release unless the owner explicitly approves an exception. |
| Required | Must be satisfied before approval unless an explicit waiver is recorded. |
| Guideline | Improve when practical. Does not block progress by itself. |

Blockers cover user file loss, silent overwrites of user-authored content, corrupt planning data, broken release binaries, and guidance that contradicts the shipped behavior.

## Rule Index

### User Files And Filesystem Safety

| ID | Severity | Rule |
|---|---|---|
| FS-01 | Blocker | Commands must not overwrite or delete user-authored content without an explicit, documented decision to do so. |
| FS-02 | Blocker | `upgrade-assets` must leave user-edited files byte-identical unless the file is inside a declared Savepoint-managed region. |
| FS-03 | Blocker | Dry-run modes must report the planned actions without writing to the filesystem. |
| FS-04 | Required | Writes that replace an existing file must be safe to re-run: a second run on unchanged input reports unchanged and changes nothing. |
| FS-05 | Required | Paths must be built with `filepath` joins, never string concatenation, so Windows targets stay correct. |
| FS-06 | Required | Commands must fail clearly and without partial writes when the target directory is missing, unwritable, or not a Savepoint project. |

### Planning Data Integrity

| ID | Severity | Rule |
|---|---|---|
| DATA-01 | Blocker | Task, router, defect, and audit-register parsing must preserve unknown frontmatter fields and body content on rewrite. |
| DATA-02 | Blocker | Lifecycle vocabulary is owned by `internal/data`: router `state`, task `status`, and `stage` only when `status: in_progress`. No package may introduce a parallel vocabulary. |
| DATA-03 | Required | Malformed or partially written planning files must produce a named diagnostic, not a panic or a silent default. |
| DATA-04 | Required | Self-healing defaults must be visible: anything the loader repairs must be reportable by `savepoint doctor`. |
| DATA-05 | Required | Agents must never set a task to `done` or retreat it to an earlier status; only the user may. |

### Templates And Shipped Assets

| ID | Severity | Rule |
|---|---|---|
| TPL-01 | Blocker | Canonical `agent-skills/{skill}/SKILL.md` must stay byte-identical to its scaffolded copy in whichever shipped tree carries it — the nine V1 skills and `references/audit-method.md` in `templates/project/agent-skills/`, and the four V2 skills plus their three shared references in `templates/project-v2/agent-skills/`. |
| TPL-02 | Required | Shipped guidance must describe behavior that the current code actually has. |
| TPL-03 | Required | A template that references another `.savepoint/` file must degrade gracefully when that file is absent; absence is not a finding. |
| TPL-04 | Required | New scaffold files must reach existing projects through a declared upgrade path, or the epic must state why they are fresh-init only. |

### Architecture

| ID | Severity | Rule |
|---|---|---|
| ARCH-01 | Required | Command entrypoints stay thin: argument parsing and dispatch only, with behavior in `internal/`. |
| ARCH-02 | Required | Rendering must not perform IO. TUI update paths do filesystem work through explicit commands. |
| ARCH-03 | Required | Correctness must not depend on the process's working directory or on in-memory state that survives a single command. |
| ARCH-04 | Required | Each `internal/` package keeps the single purpose recorded in the AGENTS.md Codebase Map; a new responsibility means a new package or a map update. |

### Configuration And Dependencies

| ID | Severity | Rule |
|---|---|---|
| CFG-01 | Required | Invalid or missing required configuration must fail with a clear, actionable message. |
| CFG-02 | Required | Behavior must not vary silently across platforms; platform differences must be explicit and tested. |
| DEP-01 | Required | New third-party dependencies require explicit justification against existing project patterns. |
| DEP-02 | Guideline | Prefer the standard library and the existing Bubble Tea / Lip Gloss stack over new abstractions. |

### Testing And Evidence

| ID | Severity | Rule |
|---|---|---|
| TEST-01 | Required | Every changed behavior must have named outcome evidence: a test, a command transcript, or an explicit scenario validation. |
| TEST-02 | Required | New or changed behavior must include happy-path evidence and one relevant failure-path case. |
| TEST-03 | Blocker | Changes to file writing, scaffolding, or upgrade paths must include evidence that existing user content is preserved. |
| TEST-04 | Required | Tests must use temporary directories and must not depend on network access or on the developer's real project files. |
| TEST-05 | Required | Bug fixes must include a regression test or an explicit failing scenario that proves the bug. |
| TEST-06 | Required | "Existing tests cover it" is acceptable only when the exact test file and test case names are recorded. |
| TEST-07 | Required | Coverage percentage alone does not satisfy evidence for changed behavior. |
| TEST-08 | Blocker | `make build && make test-fast` must pass before ordinary Task handoff. Migration/platform-sensitive Task handoff requires fresh `make test-full`; Full Objective and Goal Checks require current successful `make test-full` evidence. Focused tests are iteration aids only. A prior full result may be reused only for a metadata-only correction with the original run recorded and code, tests, fixtures, dependencies, and gate definitions proven unchanged. |
| TEST-09 | Required | If an optional Task Check is skipped, Task evidence must carry an explicit owner waiver naming the Task, reason, actor, and time; the waiver never replaces the mandatory Objective or Goal Check. |

### Release And Distribution

| ID | Severity | Rule |
|---|---|---|
| REL-01 | Blocker | Cross-compiled release artifacts must build for every declared target, including Windows. |
| REL-02 | Required | Distribution archives must ship with checksums. |
| REL-03 | Required | `--version` must report the version actually built into the binary. |

### Code Style

| ID | Severity | Rule |
|---|---|---|
| STYLE-01 | Guideline | **One job per file** — split files when responsibilities mix. |
| STYLE-02 | Guideline | **One job per function** — small, named, testable units. |
| STYLE-03 | Guideline | **Test branches** — cover meaningful conditionals and edge cases. |
| STYLE-04 | Guideline | **Types document intent** — prefer explicit types over comments. |
| STYLE-05 | Guideline | **Build only what is needed** — no speculative abstractions. |
| STYLE-06 | Guideline | **Handle errors at boundaries** — validate inputs, APIs, IO, and external data. |
| STYLE-07 | Guideline | **One source of truth** — no duplicated rules, constants, state, or config. |
| STYLE-08 | Guideline | **Comments explain why** — not what the code already says. |
| STYLE-09 | Guideline | **Content lives in data** — keep copy/config out of logic. |
| STYLE-10 | Guideline | **Small diffs** — minimal, reviewable, behaviour-preserving changes. |

### Operational Quality And Policy Boundaries

| ID | Severity | Rule |
|---|---|---|
| POL-01 | Blocker | No skill, hook, script, or sub-agent may fail compliance on a rule that is not defined in this document. |
| POL-02 | Required | Operator guides and skill docs may reference guardrails, but must not redefine them. |

## Savepoint Enforcement

Savepoint checks define how these policies are applied:

- Task handoff: per-criterion implementation evidence and configured quality
  gates. An independent Task Check is optional and may be skipped only with
  the explicit owner waiver required by TEST-09.
- Full Objective Check: mandatory before Objective closure; it reviews every
  owned Task, including waived Tasks, plus integration and Design
  reconciliation.
- Full Goal Check: mandatory whenever a Goal exists; it reviews all
  member Objectives and cross-Objective integration before exact owner
  acceptance. Existing Check records retain `scope.kind: release` and the
  Goal's `R-###` compatibility identity.

Checks may fail work only on rules defined here, unmet Savepoint acceptance
criteria, missing evidence, or explicit release gates.

Required waivers must be explicit and documented. Blocker exceptions require direct owner approval.

Code style is Guideline severity throughout. `STYLE-01..10` inform review and are never a blocker on their own.
