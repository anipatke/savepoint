# Audit Prompt

Canonical, reusable prompt for a register-backed Savepoint audit. Frozen here as
migration source evidence; kept close to the shipped template shape so a converter can
recognize an authored, in-use prompt rather than an edited one.

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

## Changelog

- **v1 (initial)** — Reconcile against the register, require stable IDs, per-finding
  fields, proof for `verified`, and examined/unexamined coverage accounting.
