---
type: health-check-ceremony
status: active
last_reviewed: 2026-05-15
---

# Savepoint Health Check

## Purpose

Health checks are Savepoint gates that produce compact evidence. They do not
define policy. `GUARDRAILS.md` defines policy, and release audit plans map that
policy to the active work.

Agents should load this file only when building, auditing, verifying, or
preparing a release. Do not load old `.claude` health-check skills.

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
- release guardrails audit plan;
- only the relevant `GUARDRAILS.md` rule IDs or compact sections.

Check:

- acceptance criteria have concrete evidence;
- changed behavior has named tests or an explicit scenario;
- relevant blocker and required guardrails are satisfied or explicitly waived;
- high-risk files have their required negative/failure-path evidence;
- known debt was reduced, preserved, or widened;
- file reality evidence confirms changed/read files were verified against the
  current filesystem and no unexplained phantom files remain;
- `git diff --check` passes;
- targeted tests for touched behavior pass.

WSL / Codex sandbox note:

- if a backend or browser gate hangs/fails in a way that may involve loopback
  networking, run `make sandbox-canary` first;
- if `make sandbox-canary` fails, rerun the affected gate outside the sandbox
  with approval and record both the canary output and unsandboxed result;
- if FastAPI/Starlette `TestClient` hangs after `make sandbox-canary` passes,
  verify with a minimal `TestClient` canary before treating it as product
  evidence;
- if Playwright/Vite cannot bind a local dev server in the sandbox (`EPERM` on
  `127.0.0.1`/`localhost`) and `make sandbox-canary` fails, rerun that browser
  gate outside the sandbox with approval;
- do not fail a health check on sandbox-only infrastructure behavior.

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
- release guardrails audit plan;
- scoped source and test files changed by the epic.

Check:

- every task acceptance criterion has evidence;
- every task has Quick health-check evidence;
- `make build` and `make test` results are recorded;
- mapped guardrail rule IDs in the release audit plan have evidence;
- file reality evidence confirms files named in task logs exist, were
  intentionally deleted, or are explicitly recorded as discarded scratch files;
- unresolved blocker rules are absent;
- required-rule waivers are explicit and owner-approved;
- drift notes are reflected in `.savepoint/Design.md` when needed.

Output:

- `### Guardrails Verification` inside the epic `E##-Audit.md`;
- named tests, commands, staging checks, and waivers.

## Deep Check

Use Deep only for pre-launch, major refactors, security/privacy concerns, or
when the user explicitly asks.

Inputs:

- release PRD;
- release Opus traceability;
- release guardrails audit plan;
- all epic audits;
- staging evidence.

Check:

- every Opus critical/high finding is closed or explicitly accepted by owner;
- no known debt is accidentally widened;
- billing replay, retention, auth boundaries, service-role/RLS posture,
  frontend/API contracts, async LLM runtime, and staging journey evidence are
  coherent across epics;
- release freeze remains intact;
- production launch criteria are satisfied.

Output:

- release-level readiness note or `NEEDS WORK` list;
- no code or task changes unless separately requested.

## Rule Boundary

Health checks may fail work only on rules defined in `GUARDRAILS.md`, unmet
Savepoint acceptance criteria, missing evidence, or explicit release gates.
They must not invent new blocking policy.



# DO NOT DELETE OR AMEND WITHOUT USER INPUT #

# GUARDRAILS.md - Engineering Policy

## 1. Purpose

This document is the authoritative engineering policy for QuizKids.

If a rule can block work, fail a health check, or require remediation, it must
be defined here. Skills, hooks, scripts, and sub-agents may reference these
rules, but must not invent blocking policy.

Savepoint defines when evidence is required. This file defines what must be
protected.

## 2. Severity Model

| Severity | Meaning |
|---|---|
| Blocker | Must be fixed before task handoff, approval, or deploy unless the owner explicitly approves an exception. |
| Required | Must be satisfied before approval unless an explicit waiver is recorded. |
| Guideline | Improve when practical. Does not block progress by itself. |

Blockers normally cover auth/authz failures, billing integrity risks, secret
exposure, child-data privacy violations, retention violations, unsafe schema
changes, and unsafe external-service or job behavior in critical flows.

## 3. Rule Index

### Security

| ID | Severity | Rule |
|---|---|---|
| SEC-01 | Blocker | Backend request bodies must use typed validation schemas. |
| SEC-02 | Blocker | Protected routes must require approved authenticated user context. |
| SEC-03 | Blocker | Secrets, tokens, and credentials must not be hardcoded or committed. Production credentials must come from approved secret storage. |
| SEC-04 | Required | User-facing API identifiers must not expose guessable sequential IDs where that creates enumeration risk. |
| SEC-05 | Blocker | Security-sensitive routes must enforce protection server-side, not client-side only. |
| SEC-06 | Blocker | Auth, quiz generation, and billing endpoints must apply server-side rate limiting. |
| SEC-07 | Blocker | Client-facing error responses must not expose stack traces, internal file paths, raw exception messages, secrets, payment data, or child data. |
| SEC-08 | Required | Production CORS origin allowlists must be explicit. Wildcard origins are not permitted in production. |

### Privacy And Retention

| ID | Severity | Rule |
|---|---|---|
| PRIV-01 | Blocker | Child-authored quiz text and sensitive derived text must be scheduled for deletion within 48 hours of quiz completion and deleted within 60 hours. (Amended per owner-signed waiver 2026-07-04 — see `docs/decisions.md`; disposition §4 item 1. Worst case = `purge_at` interval +48h plus cron period 12h.) |
| PRIV-02 | Blocker | Logs must not contain secrets, payment details, or child-authored text. |
| PRIV-03 | Blocker | Sensitive child data must not be stored in ad hoc files, local caches, or uncontrolled intermediate storage. |
| PRIV-04 | Required | New child-data flows must state what is stored, where it is stored, and when it is deleted. |
| PRIV-05 | Blocker | Child-facing quiz fetch payloads must not expose correct answers, marking keys, explanations, or grading-only metadata unless the active PRD explicitly allows it. |

### Auth And Admin Boundaries

| ID | Severity | Rule |
|---|---|---|
| AUTH-01 | Blocker | Admin access must be enforced server-side and never inferred from client state alone. |
| AUTH-02 | Required | Admin routes and admin-only actions must be explicit in code and tests. |
| AUTH-03 | Blocker | Parent, child, and admin access paths must not leak capabilities across roles. |
| AUTH-04 | Blocker | If service-role Supabase access is used in a request path, the task must document and test the application-layer ownership check. Do not claim RLS protects that path. |

### Billing And External Services

| ID | Severity | Rule |
|---|---|---|
| EXT-01 | Blocker | Critical external calls must define timeout and failure behavior. |
| EXT-02 | Required | Retries must be deliberate and used only where duplicate side effects are safe. |
| EXT-03 | Blocker | Billing, email, webhook, and LLM integrations must not leak sensitive payloads in errors or logs. |
| EXT-04 | Blocker | Critical flows must fail in a controlled and user-safe way when external services fail. |
| BILL-01 | Blocker | Billing and webhook side effects must be idempotent and duplicate-safe. |
| LLM-01 | Required | Child/user-authored text sent to LLMs must be treated as untrusted data, delimited from instructions, and validated against a bounded output schema. |
| LLM-02 | Required | User-triggered LLM generation must have quota, timeout, retry, and cost controls appropriate to the route. |

### Jobs And Events

| ID | Severity | Rule |
|---|---|---|
| JOB-01 | Blocker | Background jobs must be safe to retry. |
| JOB-02 | Blocker | Duplicate webhook, cron, or queue events must not create duplicate side effects. |
| JOB-03 | Required | Long-running work must not depend on in-memory process state for correctness. |
| JOB-04 | Required | State transitions in asynchronous flows must be explicit and defensible. |

### Architecture

| ID | Severity | Rule |
|---|---|---|
| ARCH-01 | Required | Route handlers should stay thin. Use services when behavior spans repositories, external APIs, jobs, or non-trivial business rules. |
| ARCH-02 | Required | Database and external clients must be injected, not hidden as module-level application globals. |
| ARCH-03 | Required | Application correctness must not depend on local file state or in-memory session state. |

### Database And Migrations

| ID | Severity | Rule |
|---|---|---|
| DATA-01 | Blocker | Schema changes must go through approved migrations. |
| DATA-02 | Blocker | Destructive schema changes require explicit owner approval. |
| DATA-03 | Blocker | Data migrations and backfills must include a safety plan appropriate to risk. |
| DATA-04 | Blocker | Multi-step writes for one user action must be atomic: all succeed, or no partial user-visible state remains. |
| DATA-05 | Blocker | RLS must be enabled on application tables unless an explicit documented exemption/service-only reason is approved. |
| DATA-06 | Required | Migrations that create or change application tables must state RLS posture in the migration header and include policy SQL or an explicit exemption. |
| DATA-07 | Required | Recoverable user-facing records should use soft delete semantics where historical traceability matters. |
| DATA-08 | Required | Changes affecting auth, billing, child data, or retention must consider rollback or containment. |

### API Contracts And Frontend Safety

| ID | Severity | Rule |
|---|---|---|
| API-01 | Required | Client-facing response shape changes must update the declared contract or equivalent source of truth. |
| API-02 | Required | Backend/frontend contract changes must have explicit validation through tests or equivalent evidence. |
| API-03 | Blocker | Breaking contract changes must be explicit, never accidental. |
| FE-01 | Blocker | Sensitive authorization, billing, and privacy protections must not be enforced client-side only. |
| FE-02 | Blocker | Frontend code must not expose secrets, privileged tokens, or sensitive internal configuration in the client bundle. |
| FE-03 | Required | Protected navigation and session-sensitive flows must fail safely when auth state changes or expires. |

### Configuration And Dependencies

| ID | Severity | Rule |
|---|---|---|
| CFG-01 | Blocker | Required runtime configuration must fail clearly when missing or invalid. |
| CFG-02 | Required | Security-sensitive behavior must not vary silently between environments. |
| CFG-03 | Required | Development-only shortcuts must be explicit and must not leak into production paths. |
| DEP-01 | Required | New dependencies affecting auth, billing, privacy, persistence, jobs, or external services require explicit justification. |
| DEP-02 | Guideline | Other new dependencies should be justified against existing project patterns. |

### Testing And Evidence

| ID | Severity | Rule |
|---|---|---|
| TEST-01 | Required | Every changed behavior must have named outcome evidence: a test, integration check, staging check, or explicit scenario validation. |
| TEST-02 | Required | New or changed behavior must start with a failing automated test or explicit failing scenario unless a narrow exception is recorded. |
| TEST-03 | Required | New or changed backend behavior must include happy-path evidence and one relevant failure-path case. |
| TEST-04 | Blocker | Auth, ownership, role-boundary, and protected-resource changes must include negative authorization evidence. |
| TEST-05 | Blocker | Persistence, external-service, billing, job, webhook, and multi-step-write changes must include duplicate-safe or controlled-failure evidence. |
| TEST-06 | Required | "Existing tests cover it" is acceptable only when the exact test file and test case names are recorded. |
| TEST-07 | Required | Coverage percentage alone does not satisfy evidence for changed behavior. |
| TEST-08 | Required | Bug fixes must include a regression test or explicit failing scenario that proves the bug. |
| TEST-09 | Required | Unit tests must not depend on real network calls or uncontrolled external state. |
| TEST-10 | Blocker | Production deploys require named First Family Journey end-to-end evidence. |

### Observability And Runtime Safety

| ID | Severity | Rule |
|---|---|---|
| OBS-01 | Required | Important flows must emit structured logs with enough context to diagnose failures without leaking sensitive data. |
| OBS-02 | Required | Critical failures must be diagnosable from logs or monitoring without requiring local reproduction. |
| OBS-03 | Required | Purge failures and billing errors must produce observable signals that can trigger human response. |
| RUN-01 | Required | Blocking I/O must not run on the event loop in backend request handling. |
| RUN-02 | Required | User-facing data loads should avoid obvious N+1 query patterns. |
| RUN-03 | Required | User-triggered LLM work must not block backend request paths for launch-critical flows unless explicitly approved in the active PRD/design. |

### Operational Quality And Policy Boundaries

| ID | Severity | Rule |
|---|---|---|
| OPS-01 | Required | Exceptions must not be silently swallowed. |
| OPS-02 | Required | Runtime application code should use structured logging rather than `print()`. |
| OPS-03 | Required | Touched backend code should carry accurate type hints on public functions and important domain models. |
| POL-01 | Blocker | No skill, hook, script, or sub-agent may fail compliance on a rule that is not defined in this document. |
| POL-02 | Required | Operator guides and skill docs may reference guardrails, but must not redefine them. |

## 4. Savepoint Enforcement

Savepoint health checks define how these policies are applied:

- Quick: task handoff evidence.
- Full: epic audit evidence.
- Deep: release readiness evidence.

Release audit plans map active epics to the rule IDs they must verify. Health
checks may fail work only on rules defined here, unmet Savepoint acceptance
criteria, missing evidence, or explicit release gates.

Required waivers must be explicit and documented. Blocker exceptions require
direct owner approval.


---
name: savepoint-audit-epic
description: Performs an independent Savepoint audit or re-audit of one completed epic during audit-pending, returning a Product Owner-readable CLEAR or NEEDS WORK decision, running the Full health check, and writing the single required E##-Audit.md handoff file.
---

# Savepoint Skill: Audit Epic

## Purpose

Audit every completed task in one epic with fresh eyes, verify guardrails and
drift, and write one concise, Product Owner-readable epic audit handoff.

## Read

- `.savepoint/router.md`
- Active epic detail and all task files for the epic
- `.savepoint/Design.md`
- `AGENTS.md`
- Release guardrails audit plan and relevant `GUARDRAILS.md` rules
- `.savepoint/Health-Check.md`
- Complete scoped source and test files changed by the epic
- Current `git diff` and file reality for scoped files
- `../shared/savepoint-audit-method.md` in full

## Workflow

1. Confirm the router is `audit-pending` for the requested epic.
2. Prefer a fresh session. If this session built the epic, state that limitation
   and do not call the review independent unless the user explicitly asks to
   continue.
3. Apply the shared audit method to every completed task, acceptance criterion,
   and applicable guardrail. Build and execute its mandatory coverage matrix
   before deciding the verdict; a prose checklist or a handful of manual probes
   is not matrix evidence. For multi-step or side-effecting work, enforce the
   shared Workflow And Side-Effect Audit Lock before returning a verdict.
4. Freeze and record the initial epic audit scope lock. On re-audit, reuse it
   without adding axes, dependency layers, acceptance interpretations, or
   previously unrecorded values. Build the shared method's admission ledger
   before running re-audit probes.
5. Verify every file named in task context logs exists, was intentionally
   deleted, or is explicitly recorded as discarded scratch work. Treat an
   unexplained phantom file as a finding.
6. Review every task `## Drift Notes` entry and reconcile material drift with
   `.savepoint/Design.md`.
7. Apply the Full health check and the release guardrails audit plan.
8. Write exactly one
   `.savepoint/releases/{release}/epics/{E##-slug}/E##-Audit.md`.
9. Use the Product Owner structure below. Keep the full audit internally
   thorough, but compress the scope lock, coverage matrix, workflow inventory,
   guardrail result, and gates in the artifact instead of printing every cell or
   repeating evidence.
10. Put mechanical replacement blocks only under `## Proposed Changes`.
11. Enforce the shared convergence limit, then stop and ask the user to review
    the audit before applying proposals.

## Audit File Output

Write `E##-Audit.md` in plain, everyday language. A Product Owner should
understand the delivery decision without reading the internal audit machinery.
Default `## Main Findings` to roughly 500–900 words unless there are more than
five findings or the user asks for full technical detail. Mechanical proposed
changes and compact command evidence do not count toward that guide.

Use this order:

1. **Verdict:** `CLEAR` or `NEEDS WORK`, followed by two or three sentences
   explaining whether the epic is ready, how many issues remain, and whether
   any owner-run evidence or waiver is still needed. Also state one repository
   handoff result:
   - `CLEAR TO COMMIT/PUSH` when all repository findings and gates are clear and
     only post-push owner evidence remains; or
   - `NOT READY TO COMMIT/PUSH` when a repository finding remains.
2. **What needs attention:** findings ordered by impact. Give each finding a
   short plain-language title, then:
   - why it matters;
   - what should happen next; and
   - one compact `Evidence:` line with the smallest reproduction and exact file
     reference.
   Keep implementation detail out of the main explanation unless it changes the
   product or delivery decision.
3. **Materiality summary:** keep this table for every unresolved finding:

   | Finding | Likelihood | Impact | Materiality | Recommendation |
   |---|---|---|---|---|

4. **What is proven / not proven:** group acceptance across the epic into
   concise plain-language lists. Name every unverified criterion, owner waiver,
   or missing release gate without repeating full task wording.
5. **Audit evidence:** compress the frozen scope lock, coverage matrix,
   side-effecting workflows, supported-path boundary, not-applicable cells,
   file reality, drift review, and focused/full gates. Name the requirements,
   public surfaces, direct dependencies, axes, failed or unverified cells,
   command outcomes, and evidence limitations. Do not print the full matrix or
   operation inventory cell by cell.
6. **Guardrails Verification:** retain this named subsection for the Full health
   check. List mapped rule IDs, blocker/waiver status, and the concrete evidence
   once; refer to earlier evidence instead of copying it.
7. **Non-blocking observations:** include only when useful and label them
   clearly. Omit the section when there are none.
8. **Code Style Review:** retain the required checklist, with short notes only
   where a box is unchecked or the judgment is not obvious.
9. **Proposed Changes:** include only proportionate changes for unresolved
   findings. Put exact mechanical replacement blocks here and nowhere else.

For `CLEAR`, use the same structure but keep it shorter: explain what was
proven, retain compact audit and guardrail evidence, state that no materiality
actions are required, and omit `## Proposed Changes` when there is nothing to
apply.

On re-audit, put a compact closure map of the prior findings before current
findings. Admit a new blocking item only when the shared admission ledger maps
it to an exact frozen matrix cell or it meets the named credible-blocker
exception. Put all other newly noticed cases under
`### Non-Blocking Observations`; they do not affect the verdict, finding count,
repository handoff, or required remediation.

Use this skeleton:

````markdown
---
type: audit-findings
audited: {date}
---

# Audit Findings: E## Epic Title

## Main Findings

### Verdict

CLEAR or NEEDS WORK in plain language.

### What Needs Attention

Finding title, consequence, next step, and one compact Evidence line.

### Materiality Summary

| Finding | Likelihood | Impact | Materiality | Recommendation |
|---|---|---|---|---|
| Finding 1 | Low / Medium / High | Low / Medium / High | Low / Medium / High | Proportionate disposition and scope |

### What Is Proven / Not Proven

Concise grouped acceptance result.

### Audit Evidence

- Scope lock:
- Coverage and workflow result:
- File reality and drift:
- Gates:

### Guardrails Verification

- Rule IDs checked:
- Health check mode: Full
- Evidence:
- File reality evidence:
- Waivers or unresolved findings:

### Non-Blocking Observations

Include only when useful.

## Code Style Review

- [ ] One job per file
- [ ] One job per function
- [ ] Test branches
- [ ] Types document intent
- [ ] Build only what is needed
- [ ] Handle errors at boundaries
- [ ] One source of truth
- [ ] Comments explain WHY
- [ ] Content in data files
- [ ] Small diffs

## Proposed Changes

Only unresolved, proportionate changes. Mechanical blocks live here only.

### Target File
path/to/file

### Replace
```md
Existing text to replace.
```

### With
```md
Replacement text.
```
````

## Final Response Output

After writing the audit artifact, return a compact chat summary in this order:

1. State the verdict, finding count, and repository handoff result.
2. For `NEEDS WORK`, reproduce the audit artifact's materiality summary table
   with these columns:

   | Finding | Likelihood | Impact | Materiality | Recommendation |
   |---|---|---|---|---|

   Preserve the ratings and recommendation wording from the audit artifact. Do
   not add a separate findings list.
3. State the gate result.
4. Link to the audit file.

For `CLEAR`, state that no materiality actions are required instead of showing
an empty table.

Do not repeat the audit's narrative, evidence, or proposed changes in the final
response.

## Apply And Close

Only after the user says to apply the audit:

1. Apply approved `## Proposed Changes` blocks.
2. Update `E##-Audit.md` visible sections to describe the applied outcome.
3. Mark the epic audited and advance the router.
4. Report: "Updated audit findings."

## Rules

- Do not write product code during audit.
- Do not edit task files, the router, design records, or other planning files
  during the initial audit. The single `E##-Audit.md` is the only allowed write.
- Do not apply proposals before approval.
- Do not create more than one audit file for an epic.
- Do not treat a checked box, context-log claim, prior `CLEAR`, or passing suite
  as proof without checking the current implementation.
- Do not return a verdict until every mandatory matrix cell is classified as
  passed, finding, unverified, or not-applicable with a reason. Finding one bug
  never ends the matrix pass.
- Do not turn an out-of-scope observation into a finding unless the shared
  credible-blocker exception applies.
- Do not exceed the shared re-audit convergence limit without an explicit user
  request that says whether the scope lock is extended.
- Do not expose internal audit terms such as "invariant", "matrix cell",
  "side-effect lock", or "scope perimeter" without a short plain-language
  explanation. Prefer the product consequence over implementation detail.
- Do not repeat the same evidence under findings, acceptance coverage, audit
  evidence, guardrails, and gates. Put the detail once, then summarize or refer
  back to it.
- Follow `## Final Response Output` for the user-facing handoff.
- Use `state` only for router phase, task `status` only for task lifecycle, and
  `stage` only when an item is `in_progress`.

---
name: savepoint-audit-task
description: Performs a fresh, read-only audit or re-audit of one in-progress Savepoint task when the user explicitly requests task review during task-building, returning CLEAR or NEEDS WORK without writing an epic audit artifact.
---

# Savepoint Skill: Audit Task

## Purpose

Audit one in-progress task independently from implementation. Test acceptance
criteria as general invariants, find bypasses and boundary failures beyond the
existing suite, and return one complete verdict without changing product or
planning files.

## Read

- `.savepoint/router.md`
- Active epic detail and task files
- `AGENTS.md`
- `.savepoint/Design.md` when architecture or drift is in scope
- Release guardrails audit plan and relevant `GUARDRAILS.md` rules
- `.savepoint/Health-Check.md`
- Complete scoped source and test files changed by the task
- Current `git diff` and file reality for scoped files
- `../shared/savepoint-audit-method.md` in full

## Workflow

1. Confirm the request is an explicit audit or re-audit of one task. Keep router
   `state: task-building`; this is a request-qualified skill override, not a new
   workflow state.
2. Prefer a fresh session. If this session built the task, state that limitation
   and do not call the review independent unless the user explicitly asks to
   continue.
3. Apply the shared audit method to every task acceptance criterion and all
   applicable guardrail evidence. Build and execute its mandatory coverage
   matrix before deciding the verdict; a prose checklist or a handful of manual
   probes is not matrix evidence. For multi-step or side-effecting work, enforce
   the shared Workflow And Side-Effect Audit Lock before returning a verdict.
4. Freeze and report the initial audit scope lock. On re-audit, reuse that lock
   without adding axes or dependency layers.
5. Apply the Quick health check.
6. Enforce the shared convergence limit and return the task-audit output below.

## Output

Write the result for a Product Owner in plain, everyday language. Keep the
audit work thorough, but do not make the user read the internal audit machinery
to understand the decision. Default to roughly 350–600 words unless there are
more than five findings or the user asks for full technical detail.

Use this order:

1. **Verdict:** `CLEAR` or `NEEDS WORK`, followed by two or three sentences
   explaining whether the task is ready, how many issues remain, and whether
   any owner-run evidence is still needed.
2. **What needs attention:** findings ordered by impact. Give each finding a
   short plain-language title, then:
   - why it matters;
   - what should happen next; and
   - one compact `Evidence:` line with the smallest reproduction and exact file
     reference.
   Keep implementation detail out of the main explanation unless it changes the
   product or delivery decision.
3. **Materiality table:** keep this table for every finding:

   | Finding | Likelihood | Impact | Materiality | Recommendation |
   |---|---|---|---|---|

4. **What is proven / not proven:** group the acceptance criteria into concise
   plain-language lists. Name every unverified criterion, but do not repeat the
   full acceptance wording when a short faithful summary is enough.
5. **Audit evidence:** compress the scope lock, coverage matrix, guardrail
   result, focused/full gates, and not-applicable cells into a short evidence
   section. Name the requirements, surfaces, dependencies, axes, supported-path
   boundary, failed or unverified cells, relevant guardrail result, and command
   outcomes. Do not print the full matrix cell by cell.
6. **Non-blocking observations:** include only when useful and label them
   clearly. Omit the section when there are none.
7. Confirm that no code or planning files were changed.

For `CLEAR`, use the same structure but keep it shorter: explain what was
proven, retain a compact evidence section, and state that no materiality actions
are required.

## Rules

- Do not write product code during audit.
- Do not edit the task, router, or other planning files unless the user
  separately requests that record or state change.
- Do not create an epic audit artifact.
- Do not apply proposed fixes during audit.
- Do not treat a checked box, context-log claim, prior `CLEAR`, or passing suite
  as proof without checking the current implementation.
- Do not return a verdict until every mandatory matrix cell is classified as
  passed, finding, unverified, or not-applicable with a reason. Finding one bug
  never ends the matrix pass.
- Do not turn an out-of-scope observation into a finding unless the shared
  credible-blocker exception applies.
- Do not exceed the shared re-audit convergence limit without an explicit user
  request that says whether the scope lock is extended.
- Do not expose internal audit terms such as "invariant", "matrix cell",
  "side-effect lock", or "scope perimeter" without a short plain-language
  explanation. Prefer the product consequence over the implementation detail.
- Do not repeat the same evidence under findings, acceptance coverage, matrix
  coverage, guardrails, and gates. Put the detail once, then summarize or refer
  back to it.
- Use `state` only for router phase, task `status` only for task lifecycle, and
  `stage` only when an item is `in_progress`.

# Shared Savepoint Audit Method

Read and apply this method completely for both task and epic audits.

## Establish Scope

1. Inspect current files and the diff. Treat context logs, checked boxes, prior
   audit notes, remediation claims, and green tests as claims to verify, not
   proof.
2. Verify file reality for every scoped file before referring to it.
3. Load only the acceptance criteria, guardrail rules, design context, and
   health-check requirements that govern the scoped work.

## Freeze The Audit Scope

Before the first adversarial probe, trace the real code and workflow, then write
a numbered scope lock containing:

1. the acceptance criteria, guardrails, and release gates being tested;
2. the changed files and supported public entry points;
3. the directly relied-on runtime orchestration, external effects, and dependency
   behavior needed for those entry points to keep their promise;
4. the selected matrix axes and cells, including explicit not-applicable cells;
5. the supported-path and materiality boundary used to admit findings.

The scope lock defines what the audit means by `complete`. Do not expand into
unrelated dependency internals, unsupported configurations, or increasingly
remote hypothetical states.

During an initial audit only, new source evidence may correct a factual mistake
in the scope lock. Record the amendment and restart the affected matrix pass.
Once the initial verdict is returned, the lock is immutable for every re-audit.
Do not silently introduce a new axis, dependency layer, acceptance
interpretation, or meaning of "adjacent" during remediation review.

## Turn Acceptance Into Invariants

For every acceptance criterion:

1. Restate it internally as a general rule, not as one example.
2. List the inputs, state transitions, output paths, environment modes, and
   public entry points that can affect the rule.
3. Check the normal case, boundary values, malformed input, failure behavior,
   and at least one bypass path.
4. Run at least one independent scenario that is not merely an existing unit
   test repeated unchanged.
5. Record the expected result, actual result, and concrete evidence.

A regression test proves its example. It does not prove the surrounding
invariant.

## Build The Mandatory Coverage Matrix

Before running focused probes, create a concrete matrix for the scoped public
behavior. List every row and applicable axis; mark a cell not-applicable only
with a reason tied to scope or an acceptance rule. A prose checklist is not a
matrix.

Always include these axes when applicable:

1. **Public surfaces:** every constructor, factory, parser/validator, state
   reducer, pure renderer, output backend, and selection/helper entry point.
2. **Input shape:** normal, empty, missing, duplicate, wrong type, mixed type,
   non-finite, mutable input, and mutation after validation.
3. **State:** every state; every allowed transition; backward, skipped,
   overlapping, terminal-revival, post-failure, and representation-switch
   transitions.
4. **Environment/output:** interactive and redirected sinks crossed with normal,
   no-colour, dumb/no-cursor, failure, warning, timeout, and explicit public
   overrides. Exercise both the recommended factory and direct public backend
   construction.
5. **Boundaries:** the exact limit, immediately below and above it, and the full
   small finite range when practical. Do not test only named sample points.
6. **Sequences:** complete workflows, not isolated frames. Include initialization,
   intermediate milestones, completion, failure, retry/repeat, and final output.
7. **Representations:** serialization round trips, direct models, structured
   events, and rendered output where each exists.
8. **Text classes:** ASCII, control characters, combining marks, variation
   selectors, wide characters, emoji modifiers, regional flags, and joined emoji
   when text width or truncation is in scope. Use an independent oracle when one
   is available; otherwise use structural assertions that do not reuse the
   implementation's calculation.

For progress renderers and state-machine output, the minimum matrix is:

- every public model/event/view entry point;
- TTY and non-TTY crossed with `NO_COLOR`, `TERM=dumb`, default detection, and
  explicit overrides;
- all model states and valid/invalid transition classes;
- every integer width from the hard floor through the comfortable minimum,
  plus immediately below/at/above the maximum;
- counted progress across 0%, each milestone boundary on both sides, 100%, and
  completion with and without a result word;
- timed progress start, repeat, advance, equal value, rewind, timeout boundary,
  and counted/timed representation switches;
- malformed values and containers, including mutation after construction;
- the complete Unicode corpus listed above;
- interactive, redirected, ready, warning, failure, and no-colour final output.

Run the matrix through repeatable parameterized tests or one deterministic audit
harness where possible. Record its rows, cell classifications, and command/output
evidence. Ad-hoc probes may supplement the matrix but cannot replace it.

Once the initial audit starts its adversarial probes, this matrix is the frozen
scope lock. A missing required cell discovered during that initial audit is an
audit-process error: amend the lock explicitly and rerun the affected matrix
before returning a verdict. A re-audit may never add the missing cell
retroactively as a new blocking perimeter.

### Finite External-Boundary Matrix

When scoped code relies on a server, subprocess, browser runner, provider, or
other external boundary, classify this finite set during the initial audit:

- configured target and actual runtime target;
- startup, discovery, and reuse ordering;
- connection refusal or unavailable dependency;
- successful and non-success responses;
- redirect behavior;
- timeout and cancellation;
- malformed or unexpected response;
- retry, cleanup, and partial side effects where applicable;
- secret-safe failure output.

Mark a cell not-applicable only with a scope reason. Do not invent additional
network or toolchain edge classes during later re-audits unless they meet the
credible-blocker exception below.

## Workflow And Side-Effect Audit Lock

Apply this lock to any command or workflow with multiple operations, external
calls, persistence, transactions, generated artifacts, cleanup, or structured
progress. Derive the inventory from the actual code path and external effects,
not from task notes, an implementation-owned lifecycle registry, or the test
table that is meant to verify it.

Before testing, record this table for the complete workflow:

| Order | Real operation | Side effect / state change | Failure timing | Failure owner and final state | Cleanup / secondary failure | Independent oracle |
|---|---|---|---|---|---|---|

The inventory must include, when applicable:

- input parsing, guards, setup, constructors, and connection/client creation;
- external calls and each operation that can partly succeed;
- transaction begin, write, commit, rollback, and retry behavior;
- artifact build, encode, open, partial write, flush/sync, replace, invalidation,
  and temporary-file cleanup;
- `finally` blocks, resource close, reporter/output close, and interruption;
- the point where success becomes externally visible, including terminal output,
  exit code, structured event, manifest, cache, or database state.

For every real operation:

1. Decide whether failure is command-fatal, a warning, or secondary cleanup
   information. The decision must follow acceptance criteria or guardrails and
   must never silently replace or hide the primary failure.
2. Exercise failure before, during, and after its side effect when those states
   differ. Include partial writes/updates and a primary failure combined with a
   rollback, close, or cleanup failure.
3. Trace the operations actually entered. Compare the success trace with the
   intended order and each failure trace with its expected prefix and cleanup.
4. Verify semantic state, not only file existence, valid syntax, or matching
   counts. Generated artifacts must describe the current external/database
   state; an old but well-formed artifact may still be false.
5. Use an independent oracle. Two registries, tables, reducers, or calculations
   copied from the implementation may agree while sharing the same omission.
   Implementation-owned declarations may support evidence but cannot define the
   audit scope or prove their own completeness.
6. For security filters, parsers, and redactors, derive a corpus from the full
   accepted input grammar: quoted/unquoted values, whitespace, escapes, encoded
   forms, mixed surrounding text, and configured secret values where applicable.
7. Mutate or bypass the general rule, not only the original reproduction. A
   mutation check is useful only if an adjacent omission would also fail.

Classify every operation row and every applicable failure timing as passed,
finding, unverified, or not-applicable with a scope reason. A workflow cannot
return `CLEAR` while a real operation is missing from the inventory, an output
can disagree with external state, a secondary failure is silently swallowed, or
the oracle is self-referential.

### Matrix Completion Lock

Do not begin the verdict or final findings write-up until:

1. every mandatory cell is classified as passed, finding, unverified, or
   not-applicable with a reason;
2. every prior remediation is reproduced inside the matrix or a named adjacent
   matrix cell;
3. every finding has been followed by completion of all remaining cells; and
4. the acceptance-coverage classification has been reconciled against the
   matrix results; and
5. every applicable multi-step or side-effecting workflow satisfies the
   Workflow And Side-Effect Audit Lock above.

If time, tooling, or environment prevents completion, classify the affected
acceptance criterion as unverified and return `NEEDS WORK`; never silently shrink
the matrix.

## Perform The Adversarial Pass

Ask every applicable question:

- Can validation be bypassed through another constructor, factory, direct
  public API, serialization form, or environment path?
- Can state move backward, skip forward, overlap, revive after completion, or
  continue after failure?
- Can switching modes or representations bypass a monotonicity, ownership,
  authorization, idempotency, or safety check?
- Can redirected, no-colour, non-interactive, failure, retry, or timeout output
  lose required information or expose forbidden information?
- What happens immediately below and above every documented limit?
- Does a Unicode, parser, date/time, pagination, or numeric example cover the
  whole input class, or only the tested sample?
- Can duplicate, empty, missing, non-finite, mixed-type, or unexpected values
  create an impossible model or unhandled exception?
- Are tests checking an independent outcome, or reusing the implementation's
  own calculation and assumptions?

For state machines, structured events, parsers, renderers, auth boundaries,
billing, persistence, jobs, or multi-step writes, build a small behavior matrix
for relevant inputs and transitions. Do not rely on informal spot checks.

## Re-audit After Remediation

Use the immutable scope lock from the initial audit. Re-audit:

1. every original matrix cell;
2. every original reproduction;
3. the remediation's changed code paths;
4. only the adjacent cases already named in the original finding or scope lock;
5. the same focused and full gates.

Apply existing axes to remediation code, but do not add new axes, dependency
layers, or interpretations. A failure inside the frozen lock remains a finding.
A newly noticed issue outside it is an observation for follow-up and does not
change the verdict, unless it is a credible blocker: secret exposure,
cross-tenant or role-boundary access, child-data or privacy harm, destructive
data loss, unsafe billing side effects, or an equivalent `GUARDRAILS.md`
blocker.

Before running a re-audit probe, create an admission ledger with one row per
check:

| Re-audit check | Prior finding or remediation claim | Exact frozen matrix cell | Allowed result |
|---|---|---|---|

Require an exact frozen cell for every blocking result. A broad topic,
general-purpose axis, changed helper, plural configuration name, or newly
noticed supported value is not enough. Applying an existing axis to remediation
code means rechecking only the values, classes, and adjacent cases recorded in
the initial scope lock; it does not authorize filling an omitted initial matrix
cell later.

If a probe has no exact frozen cell:

- do not run it as a blocking probe;
- record a useful result as a non-blocking observation;
- do not include it in the finding count or remediation required for closure;
- promote it only when the credible-blocker exception above applies, naming the
  exact Blocker rule.

Start the re-audit result with a closure map of the prior findings: `closed`,
`still open`, or `unverified`. Do not renumber an observation as a new finding.
If remediation closes the recorded cells but owner-run evidence is still
impossible until commit, push, or deploy, report repository readiness separately
from the overall audit verdict.

Default convergence limit:

1. one initial audit;
2. one full re-audit after remediation;
3. if an in-scope failure remains, one targeted remediation and re-audit;
4. then stop and ask the owner to fix now, approve a permitted waiver, create
   follow-up work, or close with non-blocking observations outside the task.

Do not start a third autonomous remediation cycle or broaden scope to keep an
audit active. The owner may explicitly request further work, but that request
must name whether it extends the scope lock.

At the convergence limit, stop. Do not turn a newly noticed, non-blocking edge
case into another remediation round.

## Verify Evidence And Gates

1. Run focused tests for changed behavior and relevant failure paths.
2. Run direct type or lint checks when the default gate excludes scoped files.
3. Run `git diff --check`, `make build`, and `make test` unless the active
   health-check instructions define a narrower approved gate.
4. Apply the health-check mode required by the invoking skill.
5. Treat passing tests and gates as supporting evidence, never as a substitute
   for acceptance review.

## Complete The Findings Pass

Classify every acceptance criterion before returning a verdict:

- **Proven:** independent evidence supports the general rule.
- **Finding:** a reproducible scenario violates the rule.
- **Unverified:** required evidence could not be obtained.

Any finding or material unverified criterion means `NEEDS WORK`. Return `CLEAR`
only when every criterion is proven, relevant guardrails are satisfied, and the
required gates pass.

Before admitting an item as a finding, require all of:

1. it violates a named acceptance criterion, guardrail, or release gate;
2. it is reproducible through a supported path;
3. it lies inside the frozen scope lock;
4. the task introduced, touched, widened, or explicitly promises the behavior;
5. it has a credible consequence rather than only a theoretical possibility.

If an item fails this test, record it as an observation or omit it. Observations
do not change `CLEAR`/`NEEDS WORK`, do not expand remediation scope, and should
name follow-up work only when useful. The credible-blocker exception in the
re-audit section overrides the frozen perimeter, but the auditor must state
which blocker rule makes the exception valid.

Each finding must include:

- the violated acceptance criterion or guardrail rule;
- the smallest reproducible scenario;
- expected and actual behavior;
- exact file and line evidence; and
- the missing or inadequate test evidence.

Report all findings from the completed pass together. Do not stop after the
first issue. Do not invent requirements: findings may rely only on acceptance
criteria, `GUARDRAILS.md`, active Savepoint evidence gates, or explicit release
gates.

## Summarize Materiality

After the evidence-backed findings, summarize every finding in one compact
materiality table:

| Finding | Likelihood | Impact | Materiality | Recommendation |
|---|---|---|---|---|

Use `Low`, `Medium`, or `High` for likelihood, impact, and materiality, with a
short explanation where the rating is not self-evident.

- **Likelihood:** judge the realistic prerequisites, frequency, reachability,
  and whether normal users or only unusual operator states can trigger it.
- **Impact:** judge the credible outcome and existing containment, not only the
  theoretical worst case.
- **Materiality:** combine likelihood and impact with the task or epic's stated
  purpose and launch boundary.
- **Recommendation:** state the proportionate disposition, such as fix now,
  combine with another narrow fix, defer to named follow-up work, or accept with
  an explicit owner waiver.

Do not copy the finding order or guardrail severity into the materiality rating
without this separate check. Do not inflate a rare, contained developer-workflow
issue into a product-critical risk. Equally, do not use low likelihood to excuse
a credible blocker such as secret exposure, cross-tenant access, destructive
data loss, or child-data harm.

Materiality guides priority and remediation scope; it does not silently waive an
acceptance criterion or guardrail. If the check shows that an item does not
actually violate an in-scope requirement, reclassify it as an observation rather
than leaving it as a finding. If any finding remains, the verdict remains
`NEEDS WORK` unless the owner explicitly approves the waiver allowed by policy.
When there are no findings, state that no materiality actions are required
instead of emitting an empty table.

List observations separately from findings. Do not use finding language,
`NEEDS WORK`, or an in-task fix recommendation for an out-of-scope observation
unless the owner explicitly expands the task.
