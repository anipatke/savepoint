---
type: guardrails
status: active
last_audited: never
---

# Guardrails — Engineering Policy

> The rules below are a worked example covering the risks most projects share. **Replace these rules with your project's own engineering policy** — delete what does not apply, and add the rules your domain actually needs.

## Purpose

This document is the authoritative engineering policy for {{PROJECT_NAME}}.

If a rule can block work, fail a health check, or require remediation, it must be defined here. Skills, hooks, scripts, and sub-agents may reference these rules, but must not invent blocking policy.

Savepoint defines when evidence is required. This file defines what must be protected.

## Severity Model

| Severity | Meaning |
|---|---|
| Blocker | Must be fixed before task handoff, approval, or deploy unless the owner explicitly approves an exception. |
| Required | Must be satisfied before approval unless an explicit waiver is recorded. |
| Guideline | Improve when practical. Does not block progress by itself. |

Blockers normally cover auth/authz failures, billing integrity risks, secret exposure, sensitive-data privacy violations, retention violations, unsafe schema changes, and unsafe external-service or job behavior in critical flows.

## Rule Index

### Security

| ID | Severity | Rule |
|---|---|---|
| SEC-01 | Blocker | Backend request bodies must use typed validation schemas. |
| SEC-02 | Blocker | Protected routes must require approved authenticated user context. |
| SEC-03 | Blocker | Secrets, tokens, and credentials must not be hardcoded or committed. Production credentials must come from approved secret storage. |
| SEC-04 | Required | User-facing API identifiers must not expose guessable sequential IDs where that creates enumeration risk. |
| SEC-05 | Blocker | Security-sensitive routes must enforce protection server-side, not client-side only. |
| SEC-06 | Blocker | Auth, resource-intensive generation, and billing endpoints must apply server-side rate limiting. |
| SEC-07 | Blocker | Client-facing error responses must not expose stack traces, internal file paths, raw exception messages, secrets, or payment data. |
| SEC-08 | Required | Production CORS origin allowlists must be explicit. Wildcard origins are not permitted in production. |

### Privacy And Retention

| ID | Severity | Rule |
|---|---|---|
| PRIV-01 | Blocker | Sensitive user-authored text and sensitive derived text must have a defined retention window and a deletion schedule that meets it. |
| PRIV-02 | Blocker | Logs must not contain secrets, payment details, or sensitive user-authored text. |
| PRIV-03 | Blocker | Sensitive user data must not be stored in ad hoc files, local caches, or uncontrolled intermediate storage. |
| PRIV-04 | Required | New sensitive-data flows must state what is stored, where it is stored, and when it is deleted. |
| PRIV-05 | Blocker | Client-facing payloads must not expose privileged fields, grading or scoring keys, or internal-only metadata unless the active PRD explicitly allows it. |

### Auth And Admin Boundaries

| ID | Severity | Rule |
|---|---|---|
| AUTH-01 | Blocker | Admin access must be enforced server-side and never inferred from client state alone. |
| AUTH-02 | Required | Admin routes and admin-only actions must be explicit in code and tests. |
| AUTH-03 | Blocker | Distinct user roles must not leak capabilities across role boundaries. |
| AUTH-04 | Blocker | If privileged database or service-role access bypassing row-level security is used in a request path, the task must document and test the application-layer ownership check. Do not claim row-level security protects that path. |

### Billing And External Services

| ID | Severity | Rule |
|---|---|---|
| EXT-01 | Blocker | Critical external calls must define timeout and failure behavior. |
| EXT-02 | Required | Retries must be deliberate and used only where duplicate side effects are safe. |
| EXT-03 | Blocker | Billing, email, webhook, and LLM integrations must not leak sensitive payloads in errors or logs. |
| EXT-04 | Blocker | Critical flows must fail in a controlled and user-safe way when external services fail. |
| BILL-01 | Blocker | Billing and webhook side effects must be idempotent and duplicate-safe. |
| LLM-01 | Required | User-authored text sent to LLMs must be treated as untrusted data, delimited from instructions, and validated against a bounded output schema. |
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
| DATA-05 | Blocker | Row-level security must be enabled on application tables unless an explicit documented exemption/service-only reason is approved. |
| DATA-06 | Required | Migrations that create or change application tables must state row-level-security posture in the migration header and include policy SQL or an explicit exemption. |
| DATA-07 | Required | Recoverable user-facing records should use soft delete semantics where historical traceability matters. |
| DATA-08 | Required | Changes affecting auth, billing, sensitive user data, or retention must consider rollback or containment. |

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
| TEST-10 | Blocker | Production deploys require named end-to-end evidence for the project's critical user journey. |

### Observability And Runtime Safety

| ID | Severity | Rule |
|---|---|---|
| OBS-01 | Required | Important flows must emit structured logs with enough context to diagnose failures without leaking sensitive data. |
| OBS-02 | Required | Critical failures must be diagnosable from logs or monitoring without requiring local reproduction. |
| OBS-03 | Required | Retention/deletion job failures and billing errors must produce observable signals that can trigger human response. |
| RUN-01 | Required | Blocking I/O must not run on the event loop in backend request handling. |
| RUN-02 | Required | User-facing data loads should avoid obvious N+1 query patterns. |
| RUN-03 | Required | User-triggered LLM work must not block backend request paths for launch-critical flows unless explicitly approved in the active PRD/design. |

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

Savepoint health checks define how these policies are applied:

- Quick: task handoff evidence.
- Full: epic audit evidence.
- Deep: release readiness evidence.

Release audit plans map active epics to the rule IDs they must verify. Health checks may fail work only on rules defined here, unmet Savepoint acceptance criteria, missing evidence, or explicit release gates.

Required waivers must be explicit and documented. Blocker exceptions require direct owner approval.
