---
type: health-check
status: active
last_audited: never
---

# Health Check

## Purpose

Health checks are Savepoint gates that produce compact evidence. They do not define policy. `.savepoint/Guardrails.md` defines policy. What connects that policy to the active work is the release's guardrails mapping, if your project maintains one — otherwise the relevant `.savepoint/Guardrails.md` rule IDs directly.

Agents should load this file only when building, auditing, verifying, or preparing a release.

## Modes

| Mode | When | Output |
|---|---|---|
| Quick | Before task handoff | Short task `Context Log` evidence block |
| Full | Before epic audit closeout | Epic audit `Guardrails Verification` evidence |
| Deep | Before production launch or after major concern | Release-level risk review summary |

## Quick Check

Use Quick for every task handoff.

Inputs:

- active task file;
- changed files;
- the release's guardrails mapping, if your project maintains one — otherwise the relevant `.savepoint/Guardrails.md` rule IDs directly (load only the rule IDs or compact sections that apply).

Check:

- acceptance criteria have concrete evidence;
- changed behavior has named tests or an explicit scenario;
- relevant blocker and required guardrails are satisfied or explicitly waived;
- high-risk files have their required negative/failure-path evidence;
- known debt was reduced, preserved, or widened;
- file reality evidence confirms changed/read files were verified against the current filesystem and no unexplained phantom files remain;
- `git diff --check` passes;
- targeted tests for touched behavior pass.

Output template for task `## Context Log`:

```md
Health Check: Quick
- Guardrails rule IDs:
- Acceptance evidence:
- File reality evidence:
- Tests/commands:
- Known debt:
- Waivers:
```

## Full Check

Use Full before an epic is considered audit-ready and during the epic audit.

Inputs:

- all tasks in the epic;
- each task context log;
- the release's guardrails mapping, if your project maintains one — otherwise the relevant `.savepoint/Guardrails.md` rule IDs directly;
- scoped source and test files changed by the epic.

Check:

- every task acceptance criterion has evidence;
- every task has Quick health-check evidence;
- `make build` and `make test` results are recorded;
- mapped guardrail rule IDs have evidence;
- file reality evidence confirms files named in task logs exist, were intentionally deleted, or are explicitly recorded as discarded scratch files;
- unresolved blocker rules are absent;
- required-rule waivers are explicit and owner-approved;
- drift notes are reflected in `.savepoint/Design.md` when needed.

Output:

- `### Guardrails Verification` inside the epic `E##-Audit.md`;
- named tests, commands, staging checks, and waivers.

## Deep Check

Use Deep only for pre-launch, major refactors, security/privacy concerns, or when the user explicitly asks.

Inputs:

- release PRD;
- the release's guardrails mapping, if your project maintains one — otherwise the relevant `.savepoint/Guardrails.md` rule IDs directly;
- all epic audits;
- staging evidence.

Check:

- every critical/high audit finding is closed or explicitly accepted by owner;
- no known debt is accidentally widened;
- your project's critical cross-cutting concerns (e.g. billing, retention, auth boundaries) are coherent across epics;
- release freeze remains intact;
- production launch criteria are satisfied.

Output:

- release-level readiness note or `NEEDS WORK` list;
- no code or task changes unless separately requested.

## Rule Boundary

Health checks may fail work only on rules defined in `.savepoint/Guardrails.md`, unmet Savepoint acceptance criteria, missing evidence, or explicit release gates. They must not invent new blocking policy.
