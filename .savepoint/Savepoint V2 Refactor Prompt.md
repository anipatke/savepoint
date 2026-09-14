You are working in the `anipatke/savepoint` repository.

Your immediate job is to assess the repository and produce a design proposal for **Savepoint V2** as a product simplification and workflow-quality release. Implementation follows a separately reviewed design and scoped task planning; this prompt does not authorize immediate implementation.

## Confirmed product decisions

These decisions were agreed with the owner during review. Treat them as requirements, not questions to reopen without concrete conflicting evidence. Exact schemas and implementation mechanisms remain for the design proposal.

- First-release success combines a simpler user workflow with a more reliable Task → Check handoff. Include safe migration, coherent skills and board behavior, and `resume`. Defer automatic evidence capture, hooks, and push enforcement.
- Every Task belongs to an Objective, including in tiny projects. Establish overall Design and Objective order first; write detailed Task plans one Objective at a time.
- Every Task needs a verifiable outcome. Technical and research Tasks are valid; research ends with evidence and a decision.
- The planner assesses Design and Guardrails readiness for the next Objective and brings product decisions and meaningful trade-offs to the owner.
- Tasks requiring owner validation need explicit owner acceptance. Other Tasks may complete after independent technical clearance. The planner specifies the validation requirement and steps.
- Dependent Tasks may proceed after technical clearance unless owner approval is an explicit prerequisite.
- Independent Check requires a fresh session, using the same or a different model. Independence and provenance are recorded locally; mandatory model integration is out of scope.
- Clearance applies to a recorded reviewed scope and revision or working-tree snapshot, including uncommitted changes. Relevant code or dependency changes make it stale. Agents assess and record freshness initially; automatic detection is deferred.
- Checkers write Check results and Issues, but do not repair implementation or Design. Repairs return to the executor, or the planner for design decisions, and affected work is rechecked.
- An Objective receives its integration Check after its Tasks are complete and before dependent Objectives begin. Independent Objectives may continue.
- Task IDs are globally unique and independent of Objective, title, and path. Issues use one global `I###` sequence, with Defect as a displayed type. Migration preserves old-reference mappings.
- Only material Issues affecting the Task block technical clearance, with reasons grounded in acceptance criteria, applicable blocking Guardrails, or missing required evidence. Advisory or unrelated Issues do not block.
- An independent checker verifies and closes repaired Issues. Required owner validation and acceptance of an unresolved problem remain owner decisions.
- In-scope repairs stay within the unfinished Task. Other repairs, including defects found after completion, get a linked Task under an Objective.
- `REPLAN REQUIRED` preserves partial work and pauses dependent execution. The executor records changes, verification, and the decision needed; the planner revises the plan before execution resumes.
- Executors may make targeted supporting reads and record them. Broad discovery or missing architectural decisions require replanning. Context is explicit and measured, without a universal byte cap.
- Migration converts active work and unresolved Issues; completed history remains intact in an archive with reference mappings. Do not invent retrospective V2 clearance for historical work.
- Savepoint's own controls enforce completion prerequisites and warn on advisory concerns. Explicit owner exceptions are recorded and displayed distinctly from clean clearance. Direct markdown edits cannot be prevented; validation must diagnose inconsistent records.
- Verification commands live in config, Task-specific checks in Tasks, and shared methodology in skills. Substantial reusable project procedures may remain in an optional referenced document.
- Validate with software tests and a small repeatable real-agent evaluation covering execution, replanning, and independent detection of a seeded defect.
- Use V1 to manage this repository's V2 build. Exercise V2 on fixtures and trial projects; migrate this repository after migration and the core workflow pass validation. Do not maintain parallel authoritative backlogs.

These are future V2 product rules. Existing repository lifecycle and agent-ownership rules remain in force while V1 manages the build.

The core product direction is:

# Savepoint V2
## Idea → Design → Task → Check

The goal is to make Savepoint substantially easier to understand and use without weakening the rigor underneath.

Savepoint should feel simple on the surface while preserving strong planning, scoped execution, guardrails, independent checking, defect handling, durable issue tracking, and project continuity.

Use these principles throughout:

> **Simple surface, rigorous engine.**

> **Plan deeply. Execute cheaply. Check independently.**

> **Every abstraction must earn its place by reducing ambiguity for the user, planner, executor, or checker.**

If a proposed file, state, object, lifecycle, helper document, or concept does not clearly reduce ambiguity for one of those actors, challenge whether it belongs in V2.

---

# 1. PRODUCT CONTEXT

Savepoint is currently:
- a Go CLI
- a Bubble Tea TUI
- distributed through npm
- local-first
- file-first
- agent-agnostic
- markdown/YAML driven
- designed to keep AI-assisted software projects coherent over time

The current system has accumulated too many user-facing concepts and helper documents:
- PRD
- Design
- Guardrails
- Health Check
- Releases
- Epics
- Tasks
- Audits
- Router
- Defects
- Audit Findings
- Audit Register
- multiple skills

Many of these concepts are useful internally, but the user-facing mental model is too complex.

V2 should simplify that model while preserving the valuable mechanics.

The primary danger in V2 is not missing capability.

The danger is recreating V1 complexity using friendlier terminology.

Actively resist that.

---

# 2. TARGET USER

The target user is an **AI-assisted solo builder**.

They:
- use coding agents such as Codex, Claude Code, Gemini, Cursor, or similar tools
- are comfortable with Git and the terminal
- may have limited or moderate traditional software-development experience
- may not be capable of meaningfully reviewing implementation code
- can inspect the resulting product behaviour and decide whether it works

This has an important design consequence:

> **The user validates outcomes. Savepoint verifies implementation.**

Do not assume the user can approve code quality by reading diffs.

Savepoint should act as a control layer between the user and coding agents.

The user should mainly need to answer questions such as:

- Does this behave the way I asked?
- Does the screen look right?
- Can I complete the workflow?
- Did Savepoint preserve my existing project state?
- Is this outcome acceptable?

They should not be required to answer:

- Is this abstraction well-factored?
- Is this parser architecture correct?
- Did this agent introduce an unsafe state transition?
- Is this implementation technically maintainable?

Those are primarily responsibilities of planning and checking models.

---

# 3. MODEL ROLE SEPARATION

This is a core Savepoint design principle and must not be lost in V2.

Savepoint should support three capability roles.

## Planner

Expected capability:
- frontier / strongest available reasoning model

Example only:
- Astra

Responsibilities:
- refine the Idea
- produce the Design
- refine Guardrails once Design is sufficiently mature
- define Objectives
- break Objectives into Tasks
- decide Task boundaries
- identify dependencies
- make architectural decisions
- identify exact context files
- write detailed implementation plans
- define user validation steps
- define technical verification requirements
- identify where uncertainty remains
- create research/spike work where the implementation approach is not yet known

The planner should do the expensive reasoning up front.

## Executor

Expected capability:
- cheaper, efficient model

Example only:
- Luna

Responsibilities:
- execute one Task at a time
- follow the Task implementation plan
- stay within scoped context
- respect Guardrails
- avoid reopening architectural decisions
- collect technical evidence
- report drift
- stop if the plan is materially incomplete, contradictory, or wrong

The executor should not need to rediscover the architecture or recreate the plan.

## Checker

Expected capability:
- strong independent reasoning model

Example only:
- Sol

Responsibilities:
- independently inspect the completed Task
- verify technical correctness
- verify Guardrails
- test important boundaries
- check implementation against the Task plan and Design
- identify meaningful drift
- distinguish proven behaviour from owner validation still required
- create or update Issues where follow-up is genuinely required
- return `CLEAR` or `NEEDS WORK`

The checker must not simply trust the executor's completion claim.

Use this principle:

> **Strong model plans. Cheap model executes. Independent model checks.**

Savepoint must remain model-agnostic.

Do not hard-code Astra, Luna, or Sol into the product except as examples in documentation.

The workflow should define capability roles:

```text
planner
executor
checker
```

Do not make automated model routing a mandatory V2 requirement.

The role separation matters more than automatic invocation.

---

# 4. CORE USER-FACING WORKFLOW

The primary V2 workflow becomes:

```text
IDEA → DESIGN → TASK → CHECK
                 ▲       │
                 └───────┘
```

Every Task belongs to an Objective:

```text
Objective
├── Task
├── Task
└── Task
```

Objectives are planning structure, not another mandatory workflow phase the user has to consciously operate.

The user should be able to understand Savepoint in under 60 seconds as:

> Understand the idea.  
> Design the system.  
> Complete one Task.  
> Check the result.  
> Repeat.

Objectives prevent a flat, unstructured Task backlog and define the integration Check boundary. A very small project may use one Objective.

---

# 5. PRODUCT VOCABULARY

Use the following user-facing terminology.

## Idea

What are we building, for whom, and why?

## Design

How should the system work?

## Objective

What meaningful outcome are we trying to achieve?

## Task

What is the smallest discrete, verifiable outcome an agent can deliver, with owner validation where applicable?

## Check

Did the Task actually work, and did the implementation remain technically sound?

## Issue

What requires follow-up, remediation, owner decision, verification, or explicit acceptance?

## Defect

A specific type of Issue where observed behaviour is broken, incorrect, or regressed.

Avoid making the user learn concepts such as:
- phase
- router state
- audit-pending
- task-building
- epic-task-breakdown
- audit register mechanics
- finding reconciliation terminology

Those may remain internal implementation details where useful.

---

# 6. DEFAULT PROJECT DOCUMENT MODEL

The V2 project structure should move toward something like:

```text
.savepoint/
├── Idea.md
├── Design.md
├── Guardrails.md
├── router.md
├── objectives/
│   └── O001-example/
│       ├── Objective.md
│       └── tasks/
│           ├── T001-example.md
│           └── T002-example.md
├── checks/
└── issues/
```

This is directional, not mandatory.

Inspect the current repository and propose the smallest coherent schema that supports the V2 model.

Do not preserve complexity merely because it already exists.

Do not remove useful structure merely for aesthetic simplicity.

The filesystem may remain somewhat richer than the user-facing mental model.

The TUI and agent guidance should act as the simplification layer.

---

# 7. Idea.md

`Idea.md` replaces the user-facing PRD concept.

It owns:
- what is being built
- who it is for
- why it matters
- the core experience
- high-level scope
- explicit out-of-scope items
- success criteria where useful

It should not contain technical architecture.

The user should be able to begin with a rough sentence such as:

> A tool that lets kids explore the solar system in 3D.

The planning model should refine that into a useful Idea.

Do not require the user to arrive with a complete PRD.

Refining the idea is part of the product value.

---

# 8. Design.md

`Design.md` owns the current technical shape of the system.

It should include:
- architecture
- technical approach
- major components
- key interfaces
- system boundaries
- data flow
- important implementation decisions
- important technology choices
- codebase map where useful
- current known technical state

Design should be concise enough that planner and checker models can read it regularly.

It should not duplicate Guardrails.

Design should answer:

> **How is this system intended to work?**

The planner should use Design as the technical foundation for Objective and Task decomposition.

---

# 9. Guardrails.md

Keep `Guardrails.md`.

It is valuable and should remain a first-class project artifact.

Guardrails should answer:

> **What must the agent not casually break?**

It should contain durable constraints such as:
- security rules
- destructive-change constraints
- source-of-truth rules
- architectural invariants
- data-handling requirements
- dependency restrictions
- operational requirements
- reliability rules
- style rules where genuinely useful

Guardrails should not become:
- another Design document
- generic best-practice boilerplate
- a giant engineering constitution
- a dumping ground for every preference

Prefer roughly 10–20 durable rules grouped sensibly unless the project genuinely requires more.

Use stable rule IDs where useful so Tasks and Checks can reference them.

---

# 10. GUARDRAILS REFINEMENT TIMING

This is important.

`Guardrails.md` should not be treated as fully settled at project initialization.

The expected sequence is:

```text
Idea
  ↓
Design
  ↓
Refine Guardrails
  ↓
Objectives
  ↓
Tasks
```

The planner should first progress `Design.md` until the architecture is sufficiently clear to support meaningful planning.

Only then should `Guardrails.md` be refined and right-sized.

This matters because useful Guardrails depend on understanding:
- architecture
- data model
- runtime environment
- dependencies
- trust boundaries
- persistence model
- risk profile
- delivery constraints

Writing Guardrails too early risks producing:
- generic policy noise
- duplicated Design content
- irrelevant rules
- constraints disconnected from the actual system

By the time Objective and Task planning begins, Guardrails should:
- contain only durable rules that materially constrain implementation
- reflect the actual Design
- avoid duplicating architecture
- remove irrelevant or generic rules
- use clear rule IDs where useful
- be concise enough for executor and checker models to read routinely

Do not create a detailed Objective/Task plan until:
1. Design is sufficiently mature
2. Guardrails have been reviewed against the Design

Treat this as a planning maturity gate.

The planner assesses readiness for the next Objective. Owner review covers product decisions and meaningful trade-offs rather than requiring technical approval of every design detail. Establish the overall Design and Objective order first, then fully detail one Objective's Tasks at a time. Reassess assumptions as earlier work completes.

---

# 11. OBJECTIVES

Replace the current heavy `Release → Epic → Task` hierarchy with a simpler Objective model where practical.

The desired default becomes:

```text
Objective → Task
```

An Objective is:

> A meaningful user or product outcome that groups related Tasks.

Example:

```text
O003 — Improve project recovery

T001 — Show the current Task after reopening Savepoint
T002 — Show the last successful Check
T003 — Show the next recommended action
```

Objectives should be created by the planner.

The planner decides:
- what Objectives are required
- their order
- dependencies
- how many Tasks are required
- Task boundaries

An Objective should contain:
- outcome
- why it matters
- success conditions
- dependencies
- associated Tasks
- major architectural considerations where needed

Do not turn Objective into another PRD.

Objectives are **mandatory**. A tiny project may have one Objective. Every Objective defines an integration Check after its Tasks complete and before dependent Objectives begin; unrelated Objectives need not wait.

---

# 12. RELEASES

Releases should no longer be a mandatory structural layer for normal work.

Where useful, retain Release as optional packaging/version metadata:

```yaml
release: v2.0
```

Do not force every project through:

```text
Release → Objective → Task
```

unless there is a genuine need.

A weekend project should not need enterprise ceremony.

---

# 13. TASK MODEL

Task remains the atomic implementation unit.

Do not rename Task to Build.

Task is clearer because:
- Build is a verb
- Task is a discrete work item
- the user can understand “current Task”
- the planner can define Tasks
- the executor can execute Tasks
- the checker can Check Tasks

Define a Task as:

> **The smallest discrete, verifiable outcome a planner can specify and an executor can deliver in one focused run, with owner validation where applicable.**

A Task should not be merely a coding to-do.

It is a **bounded execution packet**.

---

# 14. TASK SIZE IS A HARD PLANNING RULE

Task size should not be subjective hand-waving.

A Task should normally:
- be executable by one agent in one focused run
- have bounded context
- avoid requiring broad repository rediscovery
- avoid multiple unresolved architectural decisions
- produce one observable outcome
- allow technical outcomes verified through evidence and research outcomes ending in a decision
- have independently testable acceptance criteria
- have a clear handoff to Check

If a Task requires:
- major architectural invention
- broad open-ended codebase exploration
- many loosely related changes
- several independent user outcomes
- repeated replanning during execution

then it is probably too large.

Split it.

Do not optimise for the fewest Tasks.

Optimise for reliable execution.

---

# 15. TASKS MUST DESCRIBE OUTCOMES

The Task title and outcome should be written in observable product language where possible.

Poor Task:

> Refactor router parsing into a normalized adapter layer.

Better Task:

> Existing Savepoint projects reopen correctly after migration.

The underlying implementation may absolutely require a parser refactor.

That technical detail belongs in the implementation plan.

Use this rule:

> **Task titles describe outcomes. Implementation plans describe code changes.**

User-visible behavior is preferred where applicable, not mandatory for every Task. Internal refactors need demonstrated behavior preservation and their intended technical improvement. Research Tasks need evidence and a decision, not invented product behavior.

---

# 16. TASK AS A HANDOFF CONTRACT

A Task serves three audiences.

## User

The user needs:
- a clear outcome
- simple validation steps
- confidence about what “done” means

## Executor

The executor needs:
- exact scope
- context files
- dependencies
- relevant Design references
- Guardrail references
- a detailed implementation plan
- explicit boundaries
- technical verification requirements

## Checker

The checker needs:
- the promised outcome
- user validation criteria
- the planned implementation approach
- technical expectations
- relevant Guardrails
- expected evidence

Therefore a Task must be detailed enough that a cheaper execution model does not need to redo the planner's work.

---

# 17. TASK PLANNING QUALITY

Task quality is one of the most important parts of Savepoint.

The planner must:
- decompose Objectives carefully
- keep Tasks small enough to execute reliably
- identify dependencies explicitly
- identify exact context files
- include detailed implementation instructions
- avoid hidden architectural decisions
- identify relevant Design sections
- identify relevant Guardrail IDs
- define user-visible validation
- define technical verification
- define explicit out-of-scope boundaries
- flag genuine uncertainty rather than hiding it

If the planner does not know how something should be implemented, it should create a spike/research Task rather than pretending the plan is complete.

Do not optimise Task files purely for brevity.

Optimise for:

> **Minimum execution ambiguity with bounded context.**

Use explicit, focused context and measure actual usage without a universal byte cap. Replace the current universal <2KB incremental-context target when updating V2 guidance. Targeted supporting reads are permitted and recorded; repeated broad discovery indicates an incomplete plan.

---

# 18. TASK TEMPLATE

Design toward a Task structure like this:

```markdown
---
id: T004
status: planned
objective: O003
depends_on: []
planned_by: planner
owner_validation_required: true
---

# T004: Resume unfinished work

## Outcome

When I reopen Savepoint, I can see what I was working on and what to do next.

## User Check

1. Start a Task.
2. Close Savepoint before it is complete.
3. Reopen the project.
4. Run `savepoint resume`.
5. Confirm the current Task and next action are correct.

## Done When

- [ ] Current Task is correct.
- [ ] Current status is correct.
- [ ] Next action is understandable.
- [ ] Existing project state is preserved.

## Context Files

- `cmd/resume.go`
- `internal/data/project.go`
- `internal/data/router.go`

## Design References

- `Design.md` — Project State
- `Design.md` — Resume Flow

## Guardrails

- `STATE-01`
- `DATA-03`

## Implementation Plan

1. Add a project-state read function that resolves the active Objective and Task from router state.
2. Reuse the canonical Task parser rather than reading Task frontmatter directly from the command layer.
3. Add `resume` command wiring using the existing CLI registration pattern.
4. Implement non-TTY output first using the resolved state object.
5. Add TTY presentation using the same underlying state.
6. Keep `resume` read-only.
7. Add regression tests for:
   - active Task
   - completed Task awaiting Check
   - missing Task
   - malformed router state
8. Do not modify Objective creation or Check lifecycle in this Task.

## Boundaries

Do not:
- alter Objective planning
- change Check semantics
- redesign router persistence
- add automatic model routing

## Technical Verification

- [ ] Unit tests pass.
- [ ] Existing board tests remain green.
- [ ] `savepoint resume` performs no filesystem writes.
- [ ] Missing or malformed state returns a clear error.
- [ ] Existing project state remains unchanged.

## Technical Evidence

Pending.

## Drift Notes

None.
```

The exact schema should be adjusted after inspecting the current repository.

Task IDs are globally unique and stable across Objective moves, title changes, and path changes. The owner-validation field above is illustrative; the final schema must also represent dependency approval prerequisites, technical clearance, owner acceptance, and clearance freshness without duplicating state.

Do not remove implementation detail in the name of simplicity.

The Task must remain a strong handoff from planner to executor.

---

# 19. TASK PROVENANCE

Preserve lightweight provenance.

A Task should be able to record, where useful:
- what role/model planned it
- what role/model executed it
- what role/model checked it
- timestamps or run identifiers where useful

Do not turn this into telemetry.

This is local project provenance.

Its purpose is debugging and accountability.

For example:

```yaml
planned_by: planner
executed_by: executor
checked_by: checker
```

Vendor/model names may optionally be recorded if the user or agent provides them.

Do not make specific model names mandatory.

---

# 20. EXECUTOR BEHAVIOUR

The executor should:
- read the Task and its scoped dependencies, allowing recorded targeted supporting reads
- follow the implementation plan
- avoid redesigning the system
- avoid expanding scope
- satisfy acceptance criteria
- run required technical checks
- collect evidence
- record meaningful drift
- stop if the plan is materially incomplete or contradictory

The executor should not silently make new architectural decisions.

If blocked, it should not force a binary success/failure outcome.

Introduce an explicit execution escape hatch:

# REPLAN REQUIRED

The executor should return `REPLAN REQUIRED` when:
- the implementation plan is materially wrong
- an assumption in Design is invalid
- required context cannot be obtained through targeted supporting reads within the Task's boundaries
- the Task cannot be completed within its stated boundaries
- execution reveals an architectural decision the planner must make
- the Task is substantially larger than planned

Example:

```text
REPLAN REQUIRED

The Task assumes router state exposes the current Objective directly,
but the current data model does not contain that relationship.

Planner decision required:
Either extend router state or derive Objective ownership from the Task path.
```

This is preferable to executor improvisation.

The planner should then update Design, Guardrails, Objective, or Task as appropriate before execution resumes.

Preserve partial changes when returning `REPLAN REQUIRED`; do not automatically revert them. Record changed files, verification results, and the decision needed. Pause this Task's execution and work dependent on its unresolved outcome until the planner resolves the gap. Routine supporting reads do not alone trigger replanning.

---

# 21. TASK LIFECYCLE

Do not collapse all completion concepts into one state.

Distinguish between:

## Implementation complete

The executor believes the implementation work is finished.

## Check clear

The independent checker has verified technical integrity.

## Owner accepted

The user has validated any required visible outcome.

The planner declares whether owner validation is required and supplies the steps. Technical Tasks without that requirement may complete after independent clearance. For Tasks requiring it, clearance and acceptance remain separate: `CLEAR` means technical clearance, and pending owner validation remains visible.

Dependent Tasks may start after current technical clearance unless owner approval is an explicit dependency prerequisite. The proposal must define legal transitions for failed Checks, stale clearance, replanning, and owner exceptions. Do not impose owner acceptance on technical Tasks that do not require it.

These are not the same thing.

Design the lifecycle so the product can represent that distinction without becoming overly complicated.

For example, a Task may conceptually move through:

```text
PLANNED
  ↓
IN PROGRESS
  ↓
IMPLEMENTED
  ↓
CHECKED
  ↓
ACCEPTED (when owner validation is required)
```

However, do not adopt these exact states blindly.

Inspect the current lifecycle and propose the smallest representation that preserves these distinctions.

Avoid lifecycle explosion.

The key requirement is:

> Executor completion must not automatically equal final Task acceptance.

---

# 22. CHECK MODEL

Replace the user-facing Audit concept with:

# Check

There should be one public concept: Check.

Internally, Check may have multiple depths:

## Task Check

Focused verification after a single Task.

## Objective Check

Deeper reconciliation across all Tasks in an Objective.

Run after the Objective's Tasks complete and before dependent Objectives begin. Verify combined behavior, cross-task interactions, and Design accuracy. Independent Objectives may continue.

Do not make users learn multiple audit taxonomies.

---

# 23. CHECK MUST REMAIN RIGOROUS

Do not weaken the existing audit discipline just because the product terminology changes from Audit to Check.

The UX may become friendlier.

The underlying behaviour should remain:
- independent
- skeptical
- evidence-driven
- boundary-aware
- capable of finding implementation drift
- capable of challenging executor claims

The checker should actively try to disprove completion where appropriate.

Do not turn Check into:

> Tests passed, looks fine.

Independence requires a fresh session, not a different vendor or model. Record the session/run provenance locally without claiming that file-only tooling can prove independence.

Record the reviewed scope and revision or working-tree snapshot, including uncommitted content. A commit SHA alone is insufficient for a dirty tree. Changes to reviewed code or relevant dependencies make clearance stale; unrelated changes do not automatically invalidate it. Agents assess freshness initially. `resume` must distinguish recorded, stale, and unknown evidence without pretending to detect all changes automatically.

Preserve the adversarial spirit of the existing audit mechanism.

---

# 24. CHECK HAS TWO RESPONSIBILITIES

A Check should separate:

## User Outcome

Can the user see that the requested behaviour works?

Example:

```text
USER OUTCOME

✓ Resume shows the correct current Task
✓ Next action is understandable
```

If this cannot be automatically verified:

```text
OWNER CHECK REQUIRED

Please verify:
Open the settings panel on mobile and confirm the Save button
is visible without horizontal scrolling.
```

## Technical Integrity

Did the implementation remain technically sound?

Example:

```text
SYSTEM

✓ Tests pass
✓ Guardrails pass
✓ Task plan materially followed
✓ No unexpected files changed
✓ Design remains accurate
```

The final result should remain simple:

```text
CHECK T004

USER OUTCOME
✓ Current Task shown correctly
✓ Next step is understandable

SYSTEM
✓ Tests passed
✓ Guardrails passed
✓ No material design drift

RESULT
CLEAR
```

or:

```text
RESULT
NEEDS WORK
```

---

# 25. CHECK PHILOSOPHY

The user is not expected to perform deep code review.

Therefore:

> **The user validates the outcome.**

> **The checker verifies the implementation.**

The checker should:
- inspect implementation reality
- verify acceptance criteria
- run or inspect technical verification
- check relevant Guardrails
- review material deviation from the implementation plan
- identify meaningful drift against Design and specify the reconciliation needed
- identify missing owner validation
- create Issues where action is genuinely required
- avoid blindly trusting executor evidence

A deviation from the Task implementation plan is not automatically a failure.

However, material deviations should be:
- explained
- justified
- reflected in Design if they materially alter architecture

The checker records findings and Issues but does not repair code or Design. Route remediation to the executor or planner and recheck affected work. An Issue blocks clearance only when it materially affects the Task's acceptance criteria, applicable blocking Guardrails, or required evidence. Explain the reason; advisory improvements and unrelated Issues do not block.

---

# 26. ISSUES

Introduce **Issue** as the single umbrella concept for durable follow-up items.

An Issue is:

> Something discovered during development or checking that requires remediation, owner decision, follow-up verification, or explicit acceptance.

This replaces the need for separate user-facing concepts for:
- Audit Findings
- Audit Register Findings
- miscellaneous unresolved Check findings

The user should not need to understand a separate “finding” object model.

Use:

> **Checks discover Issues.**

Not every observation should become an Issue.

Create an Issue only when durable follow-up is genuinely required.

Examples:
- broken behaviour
- design drift requiring reconciliation
- Guardrail violation
- missing verification
- unresolved owner decision
- material risk requiring acceptance

Minor observations should remain inside the Check result as non-blocking notes.

Do not create Issues for trivia.

Avoid recreating Jira in a terminal.

---

# 27. DEFECTS AND ISSUES

Retain **Defect** as a clear first-class concept.

A Defect is:

> An Issue where observed product or system behaviour is incorrect, broken, or regressed.

Do not eliminate the word Defect merely to simplify taxonomy.

Conceptually:

```text
Issue
├── Defect
├── Drift
├── Guardrail
├── Verification
└── Other
```

In implementation, prefer a common Issue structure with a small `type` field.

Example:

```yaml
id: I014
type: defect
status: open
source: check
task: T004
title: Resume displays the wrong active Task
```

Suggested Issue types:

```text
defect
drift
guardrail
verification
other
```

Keep the type list small.

---

# 28. ISSUE LIFECYCLE

Prefer a simple lifecycle such as:

```text
open → in_progress → resolved
```

Support an explicit owner disposition for accepting an unresolved problem, with a representation such as:

```text
accepted
```

where the owner explicitly chooses not to remediate an Issue.

Avoid recreating the current complex finding lifecycle unless specific behaviour genuinely requires it.

Assess whether concepts such as:
- deferred
- waived
- owner_decision
- duplicate
- fixed
- verified

are all still necessary.

Prefer collapsing states where possible.

However, do not lose this principle:

> An Issue is not resolved merely because an agent claims it is fixed.

Resolution should require evidence or a subsequent Check where appropriate.

An independent checker verifies and closes repaired Issues. Owner-validation items still require owner acceptance. Only the owner may accept an unresolved problem, and that disposition must not imply a verified repair.

Repair an in-scope finding within its unfinished Task. Otherwise create a linked repair Task under an Objective, including for defects found after a Task has completed. Issues are durable follow-up records, not a second execution workflow.

---

# 29. ISSUE SOURCE

An Issue may originate from:
- user-reported Defect
- Task Check
- Objective Check
- `doctor`
- migration
- another explicit validation mechanism

Record source where useful.

Do not make Issue source into another workflow hierarchy.

---

# 30. CHECK → ISSUE RELATIONSHIP

A Check should remain an evaluation event, not a permanent backlog object.

A Check can:
- return `CLEAR`
- return `NEEDS WORK`
- create or update Issues
- identify owner validation requirements
- close or verify previously known Issues where evidence supports it

Conceptually:

```text
Task
  ↓
Check
  ├── CLEAR
  └── NEEDS WORK
         ↓
       Issues
```

This allows Check output to remain focused while Issues hold durable follow-up state.

---

# 31. ISSUE DEDUPLICATION

Preserve the useful intent behind the current Audit Register:

> The same problem should not be rediscovered as a brand-new finding every time a Check runs.

The Issue system should:
- preserve stable Issue IDs
- update an existing Issue when the same problem is observed again
- avoid duplicate Issues
- retain useful history where practical
- record evidence of resolution

Do not expose “Audit Register” as another major user-facing concept unless there is a compelling reason.

The Issues collection itself should ideally serve as the durable register.

---

# 32. SKILLS MODEL

Refactor toward four primary public skills:

```text
savepoint-idea
savepoint-design
savepoint-task
savepoint-check
```

Specialist/internal skills or references may remain where necessary.

Issue and Defect mechanics should preferably support these core skills rather than creating many additional top-level skills.

If a dedicated internal helper is needed for Issue reconciliation, keep it internal.

---

# 33. savepoint-idea

Expected capability:
- planner/frontier model

Owns:
- refining rough project intent
- writing `Idea.md`
- clarifying scope
- defining target user
- defining core outcome
- defining out-of-scope boundaries

Must not:
- design architecture
- create implementation Tasks
- write production code

---

# 34. savepoint-design

Expected capability:
- planner/frontier model

Owns:
- architecture
- technical design
- Design.md
- refining Guardrails once Design is mature
- Objective creation
- Task decomposition
- detailed Task planning

This skill should progress through:

```text
Idea
→ Design
→ Refine Guardrails
→ Objectives
→ Tasks
```

Do not produce detailed Tasks before Design and Guardrails are sufficiently mature.

---

# 35. savepoint-task

Expected capability:
- cheaper execution model

Owns:
- executing one already-planned Task
- following the implementation plan
- staying within scoped context
- respecting Guardrails
- satisfying technical verification
- collecting evidence
- recording drift
- returning `REPLAN REQUIRED` where appropriate
- stopping for Check

It should not:
- redesign the architecture
- redefine the Objective
- recreate the Task plan
- silently expand scope

---

# 36. savepoint-check

Expected capability:
- strong independent checker

Owns:
- independently verifying Task or Objective completion
- technical verification
- Guardrail verification
- reporting Design drift and routing reconciliation to the planner or executor
- acceptance criteria verification
- identifying owner validation requirements
- creating/updating Issues where meaningful follow-up is required
- returning `CLEAR` or `NEEDS WORK`

Preserve strong audit discipline where useful:
- scoped verification
- independent review
- materiality
- coverage reasoning
- verification of drift reconciliation after planner/executor updates
- proof before closure

Move heavy audit mechanics into shared references where possible.

Use:

> **Simple skills, rigorous methods.**

Avoid exposing internal terms such as:
- matrix cell
- scope lock
- admission ledger
- side-effect lock
- convergence limit

unless needed for implementation.

---

# 37. PLAN / EPIC / RELEASE SIMPLIFICATION

Review the existing:
- Release
- Epic
- Task
- Release PRD
- Task breakdown
- Audit hierarchy

The desired simplification is broadly:

```text
Release → Epic → Task
```

becoming:

```text
Objective → Task
```

with Release retained only where useful as optional version metadata.

Planning remains important.

What changes is the taxonomy, not the rigor.

The frontier planner still owns detailed decomposition.

Do not remove detailed planning.

---

# 38. HEALTH-CHECK DOCUMENT

Remove `Health-Check.md` as a mandatory default document. Put verification commands in config, Task-specific checks in Tasks, and shared checking methodology in skills/references. Guardrails remain the source of policy, not a collection of command instructions.

Preserve substantial reusable project-specific procedures in an optional referenced document when needed. During migration, account for all custom instructions before retiring a legacy file. This repository currently has no `.savepoint/Health-Check.md`; handle absence gracefully.

Do not reduce verification quality.

If `Health-Check.md` is removed, preserve its useful mechanics elsewhere.

---

# 39. VISUAL IDENTITY

`.savepoint/visual-identity.md` should not be a mandatory core project concept.

Treat it as optional context for UI-heavy projects.

Do not include it in the default V2 mental model.

Preserve Savepoint's own Atari-Noir identity.

---

# 40. HOOKS / AUTOMATION

Automatic evidence capture, hooks, and push enforcement are deferred beyond the first V2 release. Initially agents record evidence and assess freshness. Existing configured technical checks remain available; deferring capture does not remove verification requirements.

Savepoint's own controls must enforce completion prerequisites, including current technical clearance and owner acceptance when required. Advisory concerns produce warnings. Provide an explicit owner-exception path with recorded reasons and a display distinct from clean clearance. Diagnose invalid states introduced through direct markdown edits; do not claim to prevent external edits.

The following describes later automation opportunities, not first-release acceptance criteria.

Useful evidence includes:
- changed files
- current branch
- current commit SHA
- test result
- build result
- lint result
- typecheck result
- files touched outside Task scope
- material implementation-plan divergence
- design drift signals

Changed-file and command-result capture is mechanical evidence. Scope relevance, material plan divergence, and Design drift require semantic assessment; command success alone does not prove them.

Prefer high-level configuration such as:

```yaml
checks:
  technical:
    - make build
    - make test

enforcement:
  require_check_before_push: false
```

Avoid a giant hook lifecycle taxonomy.

Use this principle:

> **Users express intent; Savepoint handles ceremony.**

Future hooks should default to warnings. Keep optional hook enforcement separate from the completion prerequisites enforced by Savepoint's own workflow controls.

---

# 41. ONBOARDING

V2 onboarding should be a first-class experience.

## New project

`savepoint init`

The user should be able to begin with a rough idea.

Example:

```text
▣ SAVEPOINT

What are you building?

> A tool that lets kids explore the solar system in 3D.

✓ Savepoint initialized

Next:
Use your planning model to shape the Idea.
```

Then the planner should lead:

```text
Idea
→ Design
→ Refine Guardrails
→ Objectives
→ Tasks
```

The user should not have to perform technical decomposition themselves.

---

# 42. EXISTING PROJECT ONBOARDING

If Savepoint is initialized in an existing codebase, the planner should be able to reconstruct:
- Idea
- Design
- relevant Guardrails
- current Objectives
- current Tasks
- current Issues where appropriate
- sensible next Task

Do not require the user to manually recreate planning documents.

---

# 43. RESUME

Include in the first V2 release:

```bash
savepoint resume
```

This should answer:

> What was I doing, what changed, and what should I do next?

Example:

```text
WELCOME BACK

Objective
O003 — Improve project recovery

Current Task
T014 — Add repository import

Status
IMPLEMENTATION COMPLETE · CHECK REQUIRED

Issues
1 open issue

Evidence
Tests passing
4 files changed

NEXT
Run Check
```

Persistent project continuity without relying on AI memory is a key Savepoint benefit.

`resume` reports recorded state and evidence without running verification automatically. Show the evidence's reviewed scope/revision and recorded freshness; distinguish missing, stale, or unknown evidence from current clearance. Do not present historical test success as proof of the current working tree.

---

# 44. BOARD / TUI

The TUI should act as the simplification layer over the richer filesystem model.

The board should prioritize:

> **What do I do next?**

Surface prominently:
- current Objective
- current Task
- Task outcome
- current Task state
- next action
- Check status
- owner validation required
- open Issues
- project drift/health where useful

Keep the Kanban-style board if useful, but do not let Kanban become the whole product metaphor.

Example:

```text
NEXT

TASK T004 — Add planet selection

Outcome:
Click a planet and open its information panel.

User Check:
Click a planet and confirm the correct panel appears.

Status:
READY
```

or:

```text
NEXT

CHECK T004

Implementation complete.
Technical verification passed.
Owner validation required.
```

Review the current separate Defect and Audit Finding experiences.

Prefer a unified **Issues** view that can filter by type, while still displaying Defects clearly as Defects.

Example:

```text
ISSUES

I014  DEFECT        Resume shows wrong Task
I015  DRIFT         Design missing new parser boundary
I016  VERIFICATION  Mobile layout owner check required
```

Preserve:
- Atari-Noir identity
- narrow-terminal support
- readable hierarchy
- non-TTY fallback

Treat accidental line wrapping as a bug.

---

# 45. PROJECT HEALTH

Keep project health brutally practical.

It should answer:

> Is this project coherent enough to continue?

Prefer surfacing existing `doctor` evidence rather than inventing a second policy system.

Example:

```text
PROJECT HEALTH   HEALTHY

Idea        ✓
Design      ✓
Guardrails  ✓
Task        ✓
Check       ✓
Issues      0
```

or:

```text
PROJECT HEALTH   NEEDS ATTENTION

2 open Issues
1 unresolved Defect

Next:
Review Issues.
```

Do not let project health become another complicated framework.

`doctor` should remain practical validation, not another planning methodology.

---

# 46. DEFECTS

Retain Defects.

Do not collapse the concept away.

Defect remains useful because it means something specific:

> Observed behaviour is wrong.

A Defect should use the common Issue infrastructure where practical.

For example:

```yaml
type: defect
```

A user should still be able to say:

> Create a defect for this.

Use the single global `I###` Issue sequence, including Defects. Display Defect clearly as a type. Identity stays stable when classification changes. Preserve mappings from legacy `D###` and finding IDs, including their original release/epic scope, during migration.

---

# 47. MIGRATION

There is currently effectively one real user, so do not overengineer permanent V1/V2 compatibility.

However, do not accidentally break the maintainer's existing Savepoint boards.

Convert active work and unresolved Issues into V2. Preserve completed history intact in an archive with reference mappings rather than converting every historical record into a live V2 object. Carry forward dependency evidence without inventing retrospective V2 Checks.

Expected conceptual migration for active records:

```text
PRD.md             → Idea.md
Epic               → Objective
Task               → Task
Audit              → preserved evidence/reference; new verification uses Check
Audit Finding      → Issue
Audit Register     → Issues collection/history
Defect             → Issue with type: defect, while retaining Defect terminology
Release            → optional release metadata
Health-Check.md    → Check methodology/config where possible
```

Preserve existing:
- Task implementation plans
- Task evidence
- Design
- Guardrails
- audit history
- Defects
- meaningful findings
- user-authored content

Do not simplify old Tasks into shallow outcome-only records.

Do not discard audit findings simply because the Audit Register is being simplified.

Translate meaningful unresolved findings into Issues.

A one-time explicit command such as:

```bash
savepoint migrate
```

is acceptable.

Migration must support:
- dry-run or preview
- backups
- deterministic output
- idempotency
- safe failure
- validation after migration

Also require safe interruption and recovery, deterministic globally unique ID assignment, old-reference mappings, preserved cross-record links and dependency meaning, and preservation of unknown fields and user-authored content. Historical owner decisions and audit outcomes must remain available without being reinterpreted as V2 clearance. Define how completed prerequisites are referenced from active Tasks through archived records.

Keep one or two realistic V1 projects as regression fixtures.

---

# 48. COMPATIBILITY EXPECTATION

Do not spend weeks maintaining permanent dual-mode infrastructure.

Prefer:

1. safely migrate existing V1 projects
2. validate them
3. move the codebase forward to V2
4. retain only the legacy parsing required for migration and reasonable resilience

The goal is migration safety, not indefinite compatibility complexity.

Use the existing V1 workflow to manage this repository's V2 implementation. Validate migration and the core Task → Check loop with fixtures and trial projects first. Migrate this repository only after those pass. Do not create parallel authoritative V1/V2 planning records or change current lifecycle ownership rules prematurely.

---

# 49. WHAT TO PRESERVE

Do not lose these V1 strengths:

- local-first
- file-first
- no telemetry
- agent-agnostic
- markdown/YAML source of truth
- scoped reads
- scoped context
- explicit handoffs
- detailed Task decomposition
- detailed implementation plans
- clear acceptance criteria
- small executable Tasks
- planner/executor separation
- independent checking
- Design drift detection
- Guardrails
- doctor/validation
- TUI
- Defect tracking
- stable durable issue/finding identity where useful
- strong audit/check rigor
- upgrade safety
- no hidden cloud dependency
- user-owned files are never silently overwritten

Especially preserve:

> **Expensive reasoning happens during planning so cheaper execution remains reliable.**

---

# 50. WHAT TO REDUCE

Actively reduce:
- duplicated concepts
- helper-doc sprawl
- mandatory Release/Epic hierarchy
- user-facing router jargon
- user-facing lifecycle jargon
- duplicate Design/Guardrail content
- generic policy files
- repeated instructions across skills
- giant SKILL.md files where shared methodology can be referenced
- execution-model replanning
- hidden architectural decisions inside Tasks
- Tasks phrased only as code changes
- separate Audit Finding and Audit Register concepts where Issues can replace them
- excessive Issue types
- trivial Issues
- excessive backward compatibility
- unnecessary ceremony
- new abstractions that do not demonstrably reduce ambiguity

---

# 51. IMPLEMENTATION APPROACH

Do not start making broad changes immediately.

First perform a scoped repository assessment. Use the list below as a coverage checklist, reading relevant implementation and representative records in bounded batches. Record evidence and uncertainties; do not load every historical artifact into one context.

Cover the following areas, treating optional or absent files accordingly:
- `.savepoint/PRD.md`
- `.savepoint/Design.md`
- `.savepoint/Guardrails.md`
- `.savepoint/Health-Check.md`
- `.savepoint/router.md`
- `.savepoint/config.yml`
- `.savepoint/visual-identity.md`
- current releases
- current epics
- current Tasks
- current audit files
- current Audit Register
- current finding files
- current Defects
- `AGENTS.md`
- all Savepoint skills
- shared audit references
- board/TUI code
- parser/data model
- doctor
- init
- upgrade-assets
- CLI structure
- tests
- build tooling
- npm wrapper/package
- migration-sensitive code

Then produce a V2 design proposal before implementation. Apply the confirmed decisions above; distinguish evidence-backed implementation trade-offs from product choices already settled.

Do not make speculative code changes before the new model is coherent.

---

# 52. FIRST DELIVERABLE

Before making broad code changes, produce a concise but thorough V2 design proposal containing:

## A. Current-state diagnosis

Explain:
- what is currently too complex
- what is duplicated
- what should remain
- what should become internal
- where V1 mechanics are valuable despite awkward terminology

## B. Proposed V2 file tree

Show the exact recommended project structure.

## C. Proposed lifecycle

Define:
- user-facing lifecycle
- router states
- Task states
- distinction between implementation complete, technically clear, and owner accepted
- Check states
- Objective lifecycle if needed
- Issue lifecycle
- `REPLAN REQUIRED` handling

## D. Model-role architecture

Define:
- planner responsibilities
- executor responsibilities
- checker responsibilities

Explain how the system supports strong planner / cheaper executor / strong checker separation while remaining model-agnostic.

## E. Idea model

Show the proposed `Idea.md` structure.

## F. Design model

Show the proposed `Design.md` structure.

Explain how Design becomes sufficiently mature before planning begins.

## G. Guardrails model

Show:
- proposed Guardrail categories
- refinement timing
- relationship to Design
- rule ID strategy
- expected size
- when it is read

## H. Objective model

Show:
- Objective schema
- why Objectives are mandatory, including a single Objective for tiny projects
- integration Check timing and dependency gates
- how they replace Epic
- how to prevent flat Task backlogs

## I. Task schema

Provide a complete example `T001.md` showing:
- verifiable outcome, user-visible where applicable
- User Check
- Done When
- exact Context Files
- Design References
- Guardrail references
- detailed Implementation Plan
- explicit boundaries
- Technical Verification
- Technical Evidence
- Drift Notes
- provenance fields
- globally stable Task ID independent of Objective/path
- owner-validation requirement and dependency approval prerequisites
- clearance scope, freshness, and relevant reviewed revision/snapshot

The example must be sufficiently detailed that a cheaper execution model could implement it without reopening architectural decisions.

## J. Task sizing model

Define concrete planning rules for deciding whether a Task is small enough.

Explain:
- what constitutes a focused execution unit
- signals that a Task should be split
- how research/spike Tasks work
- how `REPLAN REQUIRED` feeds back into planning

## K. Check model

Explain:
- Task Check
- Objective Check
- user validation
- automated technical verification
- independent checker responsibilities
- adversarial verification
- `CLEAR`
- `NEEDS WORK`
- identifying Design drift and routing reconciliation to the planner or executor

## L. Issue model

Define:
- what constitutes an Issue
- when a Check creates an Issue
- when an observation should remain non-blocking instead
- Issue schema
- Issue types
- lifecycle
- stable IDs
- duplicate handling
- resolution evidence
- owner acceptance

Explicitly assess how much of the current Audit Register machinery should survive.

Prefer a single durable Issues collection over separate Findings + Register concepts unless repository evidence demonstrates a clear reason not to.

## M. Defect model

Explain:
- why Defect remains a useful concept
- how Defects relate to Issues
- unified `I###` identities and legacy reference mappings
- how Defect repair flows through Task → Check

## N. Skill model

Describe:
- `savepoint-idea`
- `savepoint-design`
- `savepoint-task`
- `savepoint-check`

For each define:
- intended model capability
- reads
- writes
- responsibilities
- forbidden actions
- escalation behaviour

## O. Health Check migration

Explain what happens to `Health-Check.md`.

## P. Release/Epic simplification

Explain exactly what happens to:
- releases
- release PRDs
- epics
- task breakdown

## Q. Audit Finding / Audit Register simplification

Explain exactly how existing:
- audit findings
- stable finding IDs
- finding history
- run history
- waivers
- owner decisions
- duplicate detection
- verified/fixed states

map into the proposed Issues model.

Do not retain complexity automatically.

Do not throw away useful convergence/history mechanics automatically either.

## R. Migration strategy

Explain how existing projects are safely converted.

## S. Board/TUI changes

Explain:
- what remains
- what changes
- how the TUI becomes the simplification layer
- how “Next” becomes more prominent
- how Objectives and Tasks are shown
- how owner validation is shown
- how Issues and Defects are surfaced
- how implementation-complete vs checked vs accepted states are represented without clutter

## T. Resume design

Define `savepoint resume`.

## U. Hooks/evidence

Define:
- first-release agent-recorded evidence and freshness assessment
- completion prerequisites enforced by Savepoint controls and recorded owner exceptions
- advisory concerns that do not block
- deferred automation opportunities, clearly outside first-release scope

## V. Provenance

Define the minimum useful local provenance for:
- planner
- executor
- checker

Do not introduce telemetry.

## W. Files/concepts to delete, merge, rename, retain

Be explicit.

For every retained or newly introduced abstraction, explain:

> What ambiguity does this remove?

If there is no convincing answer, remove it.

## X. Delivery plan

Propose the overall delivery order and small outcome boundaries. Detail Tasks only for the next Objective after Design and Guardrails readiness is assessed; do not prewrite the entire implementation backlog.

Do not implement the entire V2 in one giant change.

Use the current V1 Savepoint workflow to structure implementation until V2 migration and core workflow validation pass. Map V2 product Objectives to the current release/epic/task structure for the build, without creating two authoritative backlogs.

---

# 53. RECOMMENDED DELIVERY SEQUENCE

Unless repository reality strongly suggests otherwise:

## Phase 1 — Product model

Lock:
- Idea
- Design
- Guardrails
- Objective
- Task
- Check
- Issue
- Defect
- model roles

## Phase 2 — Guardrails / Design boundary

Clarify:
- what belongs in Design
- what belongs in Guardrails
- when Guardrails are refined
- what constitutes planning readiness

## Phase 3 — Task execution contract

Lock:
- Task size
- implementation-plan expectations
- provenance
- `REPLAN REQUIRED`
- executor boundaries
- implementation-complete vs accepted semantics

## Phase 4 — Issue / Defect model

Simplify:
- Audit Findings
- Audit Register
- Defects
- persistent follow-up state

into the smallest coherent Issue model.

## Phase 5 — Migration safety

Create:
- V1 fixtures
- backup strategy
- migration plan
- regression tests

## Phase 6 — Data model

Refactor:
- parser
- router
- Objective model
- Task model
- Check model
- Issue model

## Phase 7 — Skills

Refactor to:
- savepoint-idea
- savepoint-design
- savepoint-task
- savepoint-check

## Phase 8 — Onboarding

Improve:
- fresh project init
- existing project init
- planner handoff

## Phase 9 — Task / Check loop

Ensure:
- planner-created Task
- detailed implementation plan
- bounded Task size
- scoped executor flow
- technical evidence
- REPLAN REQUIRED path
- owner validation
- independent Check
- Issue creation
- drift reconciliation

## Phase 10 — TUI

Update:
- terminology
- hierarchy
- Objective display
- current Task
- next action
- Check state
- owner acceptance state
- unified Issues view

## Phase 11 — Resume / Health

Add:
- `savepoint resume`
- practical project health surfacing

## Phase 12 — Hooks / Evidence

Deferred beyond the first V2 release:
- automatic evidence capture
- hooks
- optional push enforcement

Agent-recorded evidence, freshness assessment, and enforcement of completion prerequisites belong in the first-release Task / Check loop, not this deferred phase.

## Phase 13 — Docs / Packaging

Update:
- README
- CLI help
- examples
- migration docs
- npm package metadata if needed
- upgrade guidance

---

# 54. TESTING EXPECTATIONS

Do not treat this as a documentation-only refactor.

Separate deterministic software tests from real-agent workflow evaluation. Software tests validate file safety, parsing, state transitions, references, evidence metadata, and rendering; they cannot prove an external agent's reasoning quality, independence, or compliance with instructions.

Run a small repeatable agent evaluation using defined inputs and expected outcomes: execute a planned Task, encounter a materially invalid plan and return `REPLAN REQUIRED`, and independently detect a seeded defect in a fresh checker session. Record outcomes and limitations; a broad multi-model benchmark is not required for the first release.

Across software tests and agent evaluations, provide appropriately matched evidence for:

- fresh V2 init
- Idea creation
- Design creation
- Guardrails refinement
- Objective creation
- mandatory Objective ownership, including a single-Objective tiny project
- Task planning
- Task sizing
- Task implementation-plan completeness
- Task scoped context
- Task execution
- executor refusal to re-plan architecture
- `REPLAN REQUIRED`
- Task implementation-complete state
- independent Check
- owner acceptance
- Task provenance
- Technical Evidence
- Check lifecycle
- adversarial Check behaviour
- Check-created Issues
- non-blocking Check observations
- Issue deduplication
- Issue lifecycle
- Defect creation
- Defect repair through Task → Check
- Guardrail references
- Design drift
- Objective Check
- V1 Audit Finding migration
- V1 Audit Register migration
- V1 Defect migration
- stable Issue identity
- V1 migration
- migration rerun/idempotency
- backup behaviour
- migration failure recovery
- preservation of existing Task plans
- preservation of user files
- board rendering after migration
- doctor after migration
- resume output
- non-TTY fallback
- narrow terminal rendering
- malformed frontmatter
- missing context files
- incomplete Task plans
- user-owned files are never silently overwritten
- conditional owner acceptance and dependency approval prerequisites
- reviewed working-tree content, stale clearance, and unknown evidence freshness
- checker remediation boundaries and verified Issue closure
- completion gates, advisory warnings, and explicit owner exceptions
- partial work preservation during replanning
- globally unique IDs, Objective moves, and legacy reference mappings
- completed-history archival and active dependencies on archived work
- Objective integration gates allowing independent work to continue

Use agent evaluation evidence to assess scoped execution, targeted supporting reads, and context consumption. Software tests may validate context declarations, but must not be presented as proof that an external executor never loads the whole project.

---

# 55. NON-GOALS

Do not add:
- cloud accounts
- authentication
- telemetry
- SaaS backend
- collaboration
- billing
- dashboards
- team management
- enterprise policy server
- MCP server unless essential
- model-vendor lock-in
- mandatory automatic model routing
- permanent complex V1/V2 dual-mode support
- a new database
- a heavyweight generic issue tracker

Keep V2 local and focused.

---

# 56. PRODUCT NORTH STAR

Use these as decision filters:

> **Savepoint V2 gives AI-assisted projects a simple rhythm: understand the idea, design the system, complete one well-planned Task, check the result, repeat.**

> **Plan deeply. Execute cheaply. Check independently.**

> **Tasks describe verifiable outcomes, user-visible where applicable, and contain implementation-grade execution plans.**

> **A Task should be small enough for one focused execution run with bounded context.**

> **The planner makes the hard decisions before execution starts.**

> **The executor follows the plan rather than rediscovering it.**

> **If the plan is wrong, return REPLAN REQUIRED rather than improvising architecture.**

> **Executor completion is not the same as technical clearance or owner acceptance.**

> **The user validates the outcome.**

> **The checker verifies the implementation.**

> **Checks remain rigorous even if the user-facing word is friendlier than Audit.**

> **Checks discover Issues. Defects are Issues where behaviour is broken.**

> **Not every observation deserves to become an Issue.**

> **Guardrails are refined after Design is mature enough to make them meaningful.**

> **The TUI hides complexity rather than exposing every internal object.**

> **Every abstraction must earn its place by reducing ambiguity for the user, planner, executor, or checker.**

> **The complexity stays underneath.**

If a feature, file, phase, lifecycle state, or concept does not strengthen this system, challenge whether it belongs in V2.
