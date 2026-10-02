---
name: savepoint-design
description: Maintains Savepoint Design and Guardrails and the current Objective, detailing Tasks only for the next ready Objective, when router state is design or an executor returns REPLAN REQUIRED.
---

# Savepoint Skill: Design

## Purpose

Own the technical shape needed by the next unit of work: `Design.md` describes implemented reality, `Guardrails.md` holds durable constraints, and one Objective at a time carries detailed Tasks. Later Objectives stay named outcomes with Boundaries.

Do not write production code. Route product choices to the owner; technical readiness does not require the owner to review code.

## Trigger

Use this skill when the `Next` line starts with `Plan` or `Replan`, or router `state` is `design`. A pasted `Next` line is the selection. When an executor returns `REPLAN REQUIRED`, that routes back into this skill's workflow rather than into a separate phase.

## Next

Start from the `Next` line as AGENTS.md's Workflow describes; if `savepoint` is unavailable, follow AGENTS.md rather than guessing. AGENTS.md's Router Selection section says who changes the router. New Task creation uses the narrow `savepoint create-task` exception in the Task Creation section below.

## Read

- `.savepoint/router.md`
- `.savepoint/Idea.md`, when it exists, for the original intent
- `.savepoint/Design.md`
- `.savepoint/Guardrails.md`
- The current Objective file
- Targeted implementation evidence needed to settle readiness for the next Objective — read to verify interfaces, ownership, and constraints, not to design speculatively

Read nothing else. Do not detail Tasks for any Objective beyond the next one, and do not open source files outside the evidence a readiness decision actually needs.

## Workflow

1. Read the router, the Idea when present, Design, Guardrails, and the current Objective.
2. Update `Design.md` to describe implemented reality, and `Guardrails.md` to hold durable project constraints; keep Objective deltas as the record of planned change until reconciliation. Describe confirmed Code Health tools (the owner's `.savepoint/health/config.json`) in `Design.md` in plain words; the exact commands stay in that file.
3. Keep exactly one Objective active. Objectives beyond it stay named outcomes with Boundaries — no detailed Tasks.
4. Before detailing an Objective's Tasks, inspect its stated requirements against Idea, Design, Guardrails, dependencies, and targeted evidence. Use the Objective Decision Interview below to resolve material choices that are missing, ambiguous, or contradictory; do not fill product or verification gaps by assumption. Summarize the proposed outcome, success conditions, boundaries, and key technical decisions in plain language, and obtain the owner's explicit design confirmation. Record confirmed decisions in the Objective. Then apply the readiness gate below; do not detail Tasks for an unconfirmed or unready Objective.
5. After confirmation, detail the implementation Tasks as ID-free drafts through `savepoint create-task`, then present the plan for owner review before routing to execution. In that review, name the parallel lanes in plain words, for example "T-101 and T-102 can run side by side in worktrees; T-103 waits for both." A confirmed Objective design is not approval to implement or to mark a Task done.
6. When the implementation approach for a piece of work is unknown, write a bounded research Task with a named decision deliverable instead of a confident plan the executor will discover is fiction.
7. Split a Task that carries multiple unrelated outcomes or an unresolved architectural decision into separate Tasks. Where practical, shape Tasks so independent ones can run side by side in worktrees: no `depends_on` path between them and no overlapping Context Files. Keep work that genuinely shares files in one sequential lane rather than forcing a split.
8. When an executor returns `REPLAN REQUIRED`, treat it as re-entry here: reassess Design, Guardrails, or the Objective as needed, then resume from step 3.
9. When a Task needs a verification approach, name it and reference `agent-skills/references/check-method.md` for how it will later be evaluated; do not restate that method here. When the work touches processes, paths, signals, or file replacement, also name in its Technical Verification the platform evidence the Full Check needs (for example the native Windows CI job) and who produces it, so it exists before the Check starts; the project's own gates still decide.
10. When the next Objective's Tasks are detailed and approved, select that Objective and its first unblocked planned Task in the router, set `state: task`, and hand off to `savepoint-task`. Follow AGENTS.md's Router Selection section; do not write a free-text next action.

## Objective Decision Interview

Before the design confirmation in workflow step 4, identify unresolved decisions or assumptions that could change the current Objective's outcome, success conditions, boundaries, interfaces, data ownership, dependencies, or verification. Challenge them against Idea, Design, Guardrails, and targeted evidence allowed by Read. Answer factual questions from that evidence yourself; ask the owner about material product or verification choices.

Ask one material decision at a time, starting with decisions that constrain later ones, and wait for the answer. Use an available structured question tool, or plain text when none is available. Offer concrete options with their tradeoffs and recommend one when justified. Briefly state each settled decision, then follow any consequential uncertainty it exposes. Do not reopen settled choices without new evidence or changed requirements.

Keep the interview scoped to the current Objective. Stop when its material choices are settled and the Readiness Gate can be applied; unknown technical feasibility or implementation approach may instead need a bounded research Task with a named decision deliverable. Do not ask the owner to establish facts that require research, or detail distant work to resolve every possible branch.

Summarise the settled outcome, success conditions, boundaries, and key technical decisions for the existing explicit design confirmation; this interview adds no approval checkpoint. Record confirmed planned decisions in the Objective, keeping `Design.md` about implemented reality. On `REPLAN REQUIRED`, repeat the interview only for decisions affected by what broke, then follow the same confirmation and readiness requirements.

## Task Creation

The planner never chooses, reserves, or writes a Task ID or destination filename. Prepare the complete V2 Task Markdown as an ID-free draft. Its `objective` field may be omitted or must match the selected Objective. Create each new Task with:

```bash
savepoint create-task --objective O-### --draft <path> [project-dir]
```

The command allocates the next project-wide `T-###`, adds the Task identity and Objective to the draft, derives the filename from the title, and strict-loads the complete V2 index before reporting success. It holds the allocator lock through creation and validation; concurrent planners for different Objectives receive distinct IDs. A failed post-reservation creation may retire an ID, so never retry by selecting that number yourself.

For example, planners working concurrently on O-014 and O-015 each prepare an ID-free draft and invoke `savepoint create-task --objective O-014 --draft draft-o-014.md` or `savepoint create-task --objective O-015 --draft draft-o-015.md`. The command outputs the assigned ID and path for each; neither planner predicts, copies, or reserves a number. Use `savepoint resume` after creating or renaming any other identity-bearing V2 record to require a strict load of the complete index. No other agent-run `savepoint` command is permitted.

## Verification Contract

Apply AGENTS.md's Verification Policy to every implementation: per-criterion evidence and configured gates, an explicit owner waiver when a Task Check is skipped, and mandatory Full Objective integration including waived Tasks. Plan verification and owner validation accordingly; never substitute a Task-only result or waiver for Objective clearance.

## Required Goal Context

Every Savepoint project must have at least one live Goal selected by the router, and every live Objective must name exactly one live Goal through `release:`. Apply AGENTS.md's Required Goal Context for identity compatibility, storage, diagnostics, and migration. `savepoint init` supplies G-001; `savepoint migrate` preserves converted R-### identities. If Next says `Choose a Goal` or `savepoint doctor` reports a missing Goal, report it to the owner; do not infer a selection.

When adding another Goal:

1. Allocate a stable global `G-###` identity from the first unused G number. Keep that identity stable across title or path edits, never silently reuse it, and fail closed on duplicates. Existing R-### identities remain unchanged.
2. Author the Goal sections `Outcome`, `Why`, `Success Conditions`, and `Boundaries`. The outcome describes the grouped Objectives' result, not whether anything has been published or deployed.
3. Give each member Objective the required `release:` compatibility field containing its Goal's R-### or G-### identity. Derive membership from those Objective records; do not maintain a second membership list.
4. Keep Objectives and Tasks in their normal locations and ownership: a Goal does not nest files, own Tasks, or recreate an Objective → Task hierarchy.
5. A Goal is complete when every member Objective is complete; completion is derived, not a publishing action.

### Goal Workflow Retrospective

Once a Goal's other Objectives are planned, add one final workflow-retrospective Objective to that Goal. The planner owns it and records its outcome in that Objective, including a "no change, because…" conclusion when nothing needs to change. It reviews the workflow skills, shared references, AGENTS.md routing guidance, and scaffolded project documents against the Goal's records: REPLAN REQUIRED Tasks, NEEDS WORK Checks, Issues, and lessons carried in.

What the review may change depends on who owns the files. In a project that receives the skills from the package, it tunes the project's own Guardrails, AGENTS.md project rules, and configured gates, and records suggestions for the packaged skills as Issues rather than editing them. It adds no field, state, command, or Goal-owned Task list, and a Goal is still complete when every member Objective is complete.

## Objective Artifact Template

Write the Objective file with this structure:

```markdown
---
id: O-###
title: Objective Title
status: planned|in_progress|done
depends_on: [O-###]
release: G-###
# Optional ordering metadata: priority defaults to medium; omit rank to leave unranked.
# priority: critical|high|medium|low
# rank: 1
last_check: optional-check-id
freshness:  # optional; only to mark the latest Check stale or unknown
  state: current|stale|unknown
  check: C-###
  assessed_by: {role: checker, session: <session>}
  assessed_at: '2026-09-19T00:00:00Z'
  basis: what was compared to reach this state
---

# O-###: Objective Title

## Outcome

The concrete result this Objective delivers.

## Why

Why this Objective matters now.

## Success Conditions

Observable conditions that mean this Objective is done.

## Architectural Considerations

Interfaces, data ownership, and constraints this Objective must respect.

## Boundaries

**In scope:**
- Included work

**Out of scope:**
- Excluded work
```

Every Objective must set `release:` to its Goal's stable R-### or G-### identity; this is the compatibility field, not free-form text.

Task membership is derived from which Tasks name this Objective as their owner. Do not also maintain a second, manually kept list of member Tasks in the Objective body — that list drifts from the Tasks themselves and becomes a second source of truth.

## Task Artifact Template

Write each Task file with this structure, filled as a worked example rather than an empty skeleton:

```markdown
---
title: Resume unfinished work without changing project files
objective: O-008
status: planned
depends_on: [{task: T-013, requires: clear}]
owner_validation: {required: true}
planned_by: {role: planner, session: planning-example}
---

# Resume unfinished work without changing project files

## Outcome

Reopening an existing V2 project shows selected work, recorded Check freshness, and one understandable next action without changing project files.

## User Check

Open a project awaiting a Check; invoke resume and compare selection, owner-wait state, evidence date, and Next with the board. Confirm no evidence collection or project writes.

## Done When

Correct selection and Next; explicit stale/unknown evidence; named malformed/missing selection; unchanged file bytes and mtimes. Apply the Verification Contract for Check evidence; record required owner validation after mandatory integration evidence.

## Context Files

Name exact paths only — no globs, no directory-only entries: `src/commands/status.ts`, `src/commands/status.test.ts`, `src/project/index.ts`, `src/cli.ts`, `src/resume/render.ts`, `src/resume/render.test.ts`.

## Design References

Design sections 4, 5, 10, and 11.

## Guardrails

FS-01, FS-03, DATA-02, DATA-03, ARCH-01, ARCH-03, TEST-01..04, TEST-08.

## Implementation Plan

1. Confirm the preceding project/Next APIs exist; return REPLAN REQUIRED if they do not.
2. Add thin resume argument handling with named errors.
3. Render selection, outcome, implementation/technical/owner states, Issues, evidence basis/date, and Next.
4. Keep presentation deterministic; no subprocess to generate evidence.
5. Wire main dispatch through the injected runner pattern.
6. Verify awaiting Check, awaiting owner, stale/unknown evidence, missing target, malformed router, no Task yet, and writer failure.

## Boundaries

No new Task states, evidence collection, Objective creation, automatic model routing, or duplicate gate logic.

## Technical Verification

Focused tests during iteration; handoff and Full Check gates from AGENTS.md's Verification Policy.

## Technical Evidence

Pending execution: named cases, results, reviewed source basis, files read/changed, and limitations, recorded after the work lands.

## Drift Notes

Record architecture deltas and reconcile through the planner before Check.
```

`title` and `objective` are separate required fields with different jobs. `title` is a short, plain-English phrase written for the task's owner; it must never be the Outcome text, a truncation of it, or a restatement of the technical objective. `objective` is always an `O-###` reference to the owning Objective, never free text, and every Task belongs to exactly one Objective. `owner_validation.required` and `planned_by` are recorded at planning time, not left as placeholders.

Evaluate title readability through owner-facing scenarios; field validation alone does not establish that a title reads clearly.

## Design Template

`Design.md` carries these sections: Architecture, Components/Codebase Map, Interfaces and Data Flow, Boundaries, Decisions, Current Technical State. Design describes implemented reality concisely; it is not a backlog. Planned change lives in the active Objective's deltas until the work lands, at which point Design is reconciled to match.

## Guardrails Rule

`Guardrails.md` holds durable project constraints only, once Design is mature enough to know what they are — not task-specific instructions and not a running commentary on in-flight work. Keep roughly 10–20 substantive rules, using stable category IDs (like `FS-01`, `TEST-03`) so tasks and checks can reference a rule instead of restating its prose. Guardrails, not this skill or any task file, owns severity and exception authority for those rules. Tasks and checks must not duplicate guardrail prose; they cite the rule ID.

## Commands And Procedures Reconciliation

When reconciling a migrated project's commands and procedures, reference
`agent-skills/references/commands-and-procedures.md` for the config contract
and the Health-Check.md mapping; do not restate that mapping here.

## Issue Capture

Enter Issue capture as an entry from this workflow when planning surfaces
drift, a guardrail gap, or follow-up that does not belong inside Design or the
current Objective's Tasks. Capture it as an Issue instead of folding it into
Design or a Task plan; see `agent-skills/references/issue-capture.md` for the
artifact template and rules. This skill may read and reference an Issue; the
owner may resolve it as `accepted` or reopen a resolved one. This skill does not perform those
owner actions. When it promotes an Issue's repair into a new Objective, it
retires that Issue immediately with disposition `escalated` naming the new
Objective, per "Escalation Retires The Issue" in
`agent-skills/references/issue-capture.md`. That is the only Issue closure
this skill performs.

## Readiness Gate

Detail an Objective's Tasks only after the owner confirms its design and all of the following are settled:

- Interfaces and data ownership for the work are settled, not still open questions.
- Constraints on the work are scoped, not open-ended.
- The outcomes of any Task or Objective this one depends on are known.
- There is a verification approach for the work, even if the detailed check itself runs later.

When any of these is not yet true, that gap is the thing to resolve next — either by settling it directly or by writing a bounded research Task that produces the missing decision.

## Rules

Write only Design, Guardrails, the current Objective, the next Objective's detailed Tasks, and the routing handoff; never production code or detailed backlog beyond the next Objective. Apply AGENTS.md's lifecycle terminology, verification policy, and owner authority. Keep Task titles distinct from outcomes and Objective references as the artifact contract requires.
