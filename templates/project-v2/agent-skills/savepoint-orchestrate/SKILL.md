---
name: savepoint-orchestrate
description: Coordinate an already planned Savepoint Objective through sequential builder and independent checker agents. Use only when the owner explicitly requests orchestration; a Next line alone never selects this mode.
---

# Savepoint Skill: Orchestrate

## Purpose and Scope

Reduce handoffs for one planned Objective. The coordinator delegates implementation and independent Checks, tracks progress, and presents owner decisions together. It does not implement or check the work itself.

This version is sequential: one builder at a time in the selected checkout. Planning, parallel lanes, automatic merges, publishing and pushing are outside this mode. Existing AGENTS.md and the active Savepoint skills remain authoritative. Orchestration adds no router state, lifecycle gate or exception.

## Trigger

Use only when the owner explicitly requests orchestration of a planned Objective. A Next line alone never invokes this skill.

## Invocation

The owner invokes `orchestrate O-015` or `/savepoint-orchestrate O-015` where supported. Optional settings are `builder=<model>`, `checker=<model>` and `stop-after=T-052`. Use the host's available default models when none are named; report unavailable choices rather than silently substituting them.

`stop-after` hands back after the named Task's Check and any authorized commit. If no Objective is named, show the router selection and ask the owner to select the scope before starting.

Reviewing or editing this skill does not invoke orchestration. Invocation authorizes delegation for this workflow, but does not authorize owner acceptance, completion, exceptions, pushes or unrelated cleanup.

## Roles

- **Coordinator:** selects eligible work under AGENTS.md, briefs agents, reports progress and makes authorized commits. It does not write implementation or Check verdicts.
- **Builder:** follows `savepoint-task`; writes scoped code, tests and Task evidence. It never commits, pushes, marks work done, records owner decisions or writes Checks.
- **Checker:** a fresh agent following `savepoint-check`; writes immutable Checks and permitted Issue/evaluation records. It never repairs implementation.
- **Owner:** decides completion, required validation, waivers and exceptions.

A different agent name or model does not prove independence. Start the checker without the builder's conversation history, provide the scope and raw artifacts, and record distinct executor/checker session identities. If the host cannot provide an independent checking session, stop at handoff for an external Check; do not issue self-clearance.

## Start

1. Read AGENTS.md, the router, the selected Objective and its Tasks, and configured gates. Run read-only `savepoint resume` unless the owner supplied the Next selection. Verify live Goal context and the Objective's selection under AGENTS.md.
2. Confirm that the Tasks are planned and the checkout is suitable. Follow Worktree Lane restrictions when applicable; this mode creates no lanes. On `Plan`, `Replan`, `Choose a Goal`, stale selection or `REPLAN REQUIRED`, hand back with the reason.
3. Show a short Task progress list, including required owner validation. Identify unrelated or pre-existing changes so they stay outside commits.
4. Use dependency metadata to suggest ordering, but use Savepoint's runtime result to determine eligibility before each start. Do not infer dependency clearance from this skill.

## Workflow

1. **Build.** Brief one builder with the selected Next line, Task path, project root, relevant earlier observations and the instruction to follow `savepoint-task`, including its Read section, stage changes, configured gates and evidence requirements. Return the final stage, changed paths, evidence, drift and unverified items. Preserve existing evidence; use separate output paths when the project's tools require them.
2. **Check.** At `audit`, delegate a requested Quick Task Check to a fresh checker. Explicit orchestration requests these optional Checks; they are not a new project requirement and are not justified merely by `requires: clear`. The checker loads `savepoint-check` and the shared method in full and follows them. Retain the source, command and output of independent probes in Check evidence, even when scratch files are removed.
3. **Remediate.** On `NEEDS WORK`, resume the builder with the recorded Issues. Delegate rechecking independently under the original frozen scope and the shared convergence limit. Do not repair in the coordinator or checker. Stop for a planning decision or when that limit is reached.
4. **Owner decisions.** On `CLEAR`, report technical readiness. Batch completion decisions where practical; required User Checks still need the owner's actual acceptance. Do not mark a Task done without an explicit owner instruction covering it. An owner may name several Tasks in one instruction. Record only the decision actually given, using the canonical skill's fields and citing actor, time and scope. Never infer acceptance from silence or from a Check result.
5. **Commit.** When the owner has authorized commits, the coordinator stages an explicit allowlist of intended files, reviews the staged diff and commits with the Task/Check references. Do not use `git add -A`, include unrelated work, or alter attribution requirements. Make a separate rename/move commit only when the Task requires it. Otherwise leave the work ready to commit.
6. **Continue.** Run `savepoint resume` after record edits and before selecting the next work. A Task awaiting owner validation may allow a dependant to start only when the runtime reports its dependency met. Carry relevant observations into the next builder's brief; do not convert observations into invented requirements.

Keep one or two plain sentences of progress at each handoff. Put technical evidence in the project records.

## Owner Validation

Offer the Task's User Check with concrete steps and what passing means. Use screenshots, temporary comparison worktrees or live servers only when they help validate that project's criteria; no browser or frontend tooling is assumed.

Record acceptance only from the owner's explicit words, against the relevant Check using `savepoint-check` Closure Rules. Required owner validation and Task completion are separate decisions unless the owner explicitly gives both.

Track resources created for validation. Stop only those servers/processes and remove only owned temporary resources after checking for unsaved work. Preserve dirty worktrees and report them rather than removing them automatically.

## Pause and Resume

On a stop request, stop scheduling work and interrupt active agents or bring them to a safe point as requested. Do not check, commit or clean up work the owner asked to preserve.

Save a short handoff: Objective, Tasks and stages, Check/Issue references, current agents, owned resources, explicit owner decisions and next action. On resume, re-read current records and run `savepoint resume`; use that handoff as context, not as a substitute for current state. Resume a builder when possible, and give each Check an independent session. Never prune or delete worktrees merely because they appear stale.

## Objective Handoff

When all Tasks are done, delegate the mandatory Full Objective Check to a fresh checker. It applies the shared method, current full gate, Code Health collection and Design reconciliation. Task Checks never replace this gate.

A Full Check while Tasks await owner decisions is preliminary: it may establish technical clearance, but cannot make unfinished Tasks or the Objective complete. State remaining decisions explicitly.

Report what was built, the latest Check, unresolved Issues and owner decisions still needed. Leave Objective closure to the owner. Report measured time/token usage only when the host exposes it; mark unavailable measurements rather than inventing totals.
