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

## Trigger

Use this skill only when the user explicitly asks to audit or re-audit one
in-progress task. Router `state` stays `task-building`; this is a
request-qualified skill override, not a new workflow state. For a completed
epic, use `savepoint-audit-epic` instead.

## Read

- `.savepoint/router.md`
- Active epic detail and task files
- `AGENTS.md`
- `.savepoint/Design.md` when architecture or drift is in scope
- The release's guardrails mapping, if your project maintains one — otherwise
  the relevant `.savepoint/Guardrails.md` rule IDs directly
- `.savepoint/Health-Check.md`
- Complete scoped source and test files changed by the task
- Current `git diff` and file reality for scoped files
- `agent-skills/references/audit-method.md` in full

Treat `.savepoint/Guardrails.md` and `.savepoint/Health-Check.md` as optional:
apply them when the project has them, and skip the related step when they are
absent rather than reporting a finding.

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
   without adding axes or dependency layers, and build the shared method's
   admission ledger before running re-audit probes.
5. Verify file reality: every file named in the task context log exists, was
   intentionally deleted, or is explicitly recorded as discarded scratch work.
   Treat an unexplained phantom file as a finding.
6. Apply the Quick health check.
7. Enforce the shared convergence limit and return the task-audit output below.

## Output

Write the result for a Product Owner in plain, everyday language. Keep the
audit work thorough, but do not make the user read the internal audit machinery
to understand the decision. Default to roughly 350–600 words unless there are
more than five findings or the user asks for full technical detail.

Return exactly one result value: `CLEAR` or `NEEDS WORK`.

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
5. **Audit evidence:** compress the scope lock, coverage matrix, file reality,
   guardrail result, focused/full gates, and not-applicable cells into a short
   evidence section. Name the requirements, surfaces, dependencies, axes,
   supported-path boundary, failed or unverified cells, relevant guardrail
   result, and command outcomes. Do not print the full matrix cell by cell.
6. **Non-blocking observations:** include only when useful and label them
   clearly. Omit the section when there are none.
7. Confirm that no code or planning files were changed.

For `CLEAR`, use the same structure but keep it shorter: explain what was
proven, retain a compact evidence section, and state that no materiality actions
are required.

## Final Response Output

The chat response is the whole deliverable for a task audit. Return a compact
summary in this order:

1. State the verdict and finding count.
2. For `NEEDS WORK`, reproduce the materiality table with these columns:

   | Finding | Likelihood | Impact | Materiality | Recommendation |
   |---|---|---|---|---|

   Preserve the ratings and recommendation wording. Do not add a separate
   findings list.
3. State the gate result.
4. Link to the task file under audit. Do not link or create an audit artifact
   file; a task audit has none.

For `CLEAR`, state that no materiality actions are required instead of showing
an empty table.

## Rules

- Do not write product code during audit.
- Do not write any file. Do not edit the task, router, epic detail, design
  records, or other planning files unless the user separately requests that
  record or state change.
- Do not create an epic audit artifact. `E##-Audit.md` belongs to
  `savepoint-audit-epic` only.
- Do not change task `status`, task `stage`, or router `state`.
- Do not apply proposed fixes during audit.
- Do not use a health-check mode other than Quick.
- Do not return a result value other than `CLEAR` or `NEEDS WORK`.
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
- Follow `## Final Response Output` for the user-facing handoff.
- Use `state` only for router phase, task `status` only for task lifecycle, and
  `stage` only when an item is `in_progress`.
