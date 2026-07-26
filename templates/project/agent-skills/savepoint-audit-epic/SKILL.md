---
name: savepoint-audit-epic
description: Performs an independent Savepoint audit or re-audit of one completed epic during audit-pending, returning a Product Owner-readable CLEAR or NEEDS WORK decision, running the Full health check, and writing the single required E##-Audit.md handoff file.
---

# Savepoint Skill: Audit Epic

## Purpose

Audit every completed task in one epic with fresh eyes, verify guardrails and
drift, and write one concise, Product Owner-readable epic audit handoff.

## Trigger

Use this skill when router `state` is `audit-pending`, or when the user
explicitly asks for an audit or re-audit of a completed epic. For a read-only
review of one in-progress task, use `savepoint-audit-task` instead.

## Read

- `.savepoint/router.md`
- Active epic detail and all task files for the epic
- `.savepoint/Design.md`
- `AGENTS.md`
- The release's guardrails mapping, if your project maintains one — otherwise
  the relevant `.savepoint/Guardrails.md` rule IDs directly
- `.savepoint/Health-Check.md`
- `agent-skills/savepoint-audit-register/SKILL.md` when `.savepoint/audit/` exists
- Complete scoped source and test files changed by the epic
- Current `git diff` and file reality for scoped files
- `agent-skills/references/audit-method.md` in full

Treat `.savepoint/Guardrails.md` and `.savepoint/Health-Check.md` as optional:
apply them when the project has them, and skip the related step when they are
absent rather than reporting a finding.

## Workflow

1. Confirm the router is `audit-pending` for the requested epic, or that the
   user explicitly asked to audit this completed epic.
2. Stop if you are the same agent session that built the epic. This audit
   requires a session independent from the builder. If this session built the
   epic, state that limitation and do not call the review independent unless the
   user explicitly asks to continue.
3. When `.savepoint/audit/` exists, follow the `savepoint-audit-register` skill
   alongside this one: it owns register reconciliation, stable `F###` IDs, run
   records, and proof rules. Reconcile against the register before recording new
   findings.
4. Apply the shared audit method to every completed task in the epic, every
   acceptance criterion, and every applicable guardrail. Build and execute its
   mandatory coverage matrix before deciding the verdict; a prose checklist or a
   handful of manual probes is not matrix evidence. For multi-step or
   side-effecting work, enforce the shared Workflow And Side-Effect Audit Lock
   before returning a verdict.
5. Freeze and record the initial epic audit scope lock. On re-audit, reuse it
   without adding axes, dependency layers, acceptance interpretations, or
   previously unrecorded values. Build the shared method's admission ledger
   before running re-audit probes.
6. Verify file reality: every file named in task context logs exists, was
   intentionally deleted, or is explicitly recorded as discarded scratch work.
   Treat an unexplained phantom file as a finding.
7. Review every task `## Drift Notes` entry and reconcile material drift with
   `.savepoint/Design.md`.
8. Apply the Full health check and the release guardrails in scope.
9. Write exactly one
   `.savepoint/releases/{release}/epics/{E##-slug}/E##-Audit.md`.
10. Use the Product Owner structure below. Keep the full audit internally
    thorough, but compress the scope lock, coverage matrix, workflow inventory,
    guardrail result, and gates in the artifact instead of printing every cell
    or repeating evidence.
11. Put mechanical replacement blocks only under `## Proposed Changes`.
12. Enforce the shared convergence limit, then stop and ask the user to review
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
8. **Code Style Review:** retain the required checklist. Write one checkbox per
   `STYLE` rule found in `.savepoint/Guardrails.md`, labelled with the rule ID
   and its own wording, in the order the file lists them. Add short notes only
   where a box is unchecked or the judgment is not obvious. When the project has
   no `.savepoint/Guardrails.md`, or that file defines no `STYLE` rules, keep the
   section and state "Code style is not defined for this project." instead of a
   checklist. `STYLE` rules are Guideline severity: an unchecked box is an
   observation, never a blocker or a `NEEDS WORK` cause on its own.
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

CLEAR or NEEDS WORK in plain language, plus the repository handoff result.

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

- [ ] {STYLE rule ID} {rule wording from `.savepoint/Guardrails.md`}
- [ ] Repeat once per STYLE rule, in file order

(When no `STYLE` rules are defined, replace the list with: Code style is not
defined for this project.)

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
- Do not audit an epic you built in this session.
- Do not use a health-check mode other than Full.
- Do not hardcode code-style rule labels. Source the `## Code Style Review`
  checklist from the `STYLE` rules in `.savepoint/Guardrails.md`, and never fail
  the epic on one.
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
