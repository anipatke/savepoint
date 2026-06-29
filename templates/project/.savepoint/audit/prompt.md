# Audit Prompt

Canonical, reusable prompt for a register-backed Savepoint audit. Read this before
auditing so each run converges on the existing register instead of restarting from a
cold scan.

## Before you start

1. If `.savepoint/audit/register.md` exists, treat it as the current source of truth.
   Reconcile every prior finding before recording anything new — do not duplicate a
   finding that already has a stable `F###` ID.
2. If the register is absent, you are seeding it: open findings start at `open`.
3. Record this run as an immutable file under `.savepoint/audit/runs/` and update the
   mutable `.savepoint/audit/register.md` to reflect the reconciled state.

## Reconciliation rules

For each prior finding, decide exactly one disposition and carry the **same stable ID**:

- Still present → keep the ID, refresh `last_seen`, update status if work advanced.
- Resolved with proof → move toward `fixed` then `verified` (proof required, see below).
- No longer applicable → `deferred`, `owner_decision`, or `waived` with a recorded reason.
- Already tracked elsewhere → `duplicate` pointing at the canonical `F###` ID.

A finding reaches `verified` only with named proof — preferably a passing regression
test, otherwise an explicit manual verification note. Never mark `verified` without it.

## Required fields per finding

Every finding the audit reports or updates must carry:

- **Stable ID** — reuse the existing `F###` when the finding already exists; only mint a
  new ID for a genuinely net-new finding.
- **Title** — short, specific summary.
- **Status** — one of the lifecycle values below.
- **Severity** — `critical`, `high`, `medium`, or `low`.
- **Confidence** — `high`, `medium`, or `low` that the finding is real.
- **Source auditor** — who or what produced the finding (agent name/model or person).
- **Location** — file path(s) and line ranges, or the surface under review.
- **Guardrail IDs** — the code-style or guardrail rules implicated, when any.
- **Proof needed** — what evidence would move the finding to `verified`.
- **Work-item mapping** — linked release, epic, task, or defect ID, when one applies.

### Finding lifecycle statuses

`open`, `triaged`, `mapped`, `in_progress`, `fixed`, `verified`, `deferred`,
`owner_decision`, `waived`, `duplicate`.

## Coverage accounting

Close every run with coverage notes so gaps are visible across runs:

- **Examined** — surfaces, modules, or files you actually reviewed this run.
- **Unexamined** — surfaces deliberately or unavoidably skipped, with the reason.

Unexamined surfaces are coverage gaps and feed the register convergence summary.

## Changelog

Refinements to this prompt are recorded here so audits stay comparable over time.

- **v1 (initial)** — Reconcile against the register, require stable IDs, per-finding
  fields, proof for `verified`, and examined/unexamined coverage accounting.
