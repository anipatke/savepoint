---
id: O-042
title: Make the workflow review a note on the Goal
status: done
depends_on: [O-039, O-040, O-043, O-044]
release: G-002
priority: low
rank: 6
---

# O-042: Make the workflow review a note on the Goal

## Outcome

The end-of-Goal workflow review is a short note the planner writes in the Goal's `Release.md`, not an Objective. A review that finds nothing costs no Objective, Task or Check. A change it does find becomes ordinary work: an Objective, or an Issue.

## Why

`savepoint-design` required a final retrospective Objective for every Goal. An Objective with no Tasks never closes (`resolveSelectedObjective` in `internal/data/next.go` routes it to Plan), and every Objective needs a Full Objective Check, so even a "no change" review forced a paperwork Task and a full Check. The owner rejected that on 2026-10-10 as against the project's ethos.

## Success Conditions

- `savepoint-design` no longer asks for a retrospective Objective. It tells the planner to write a `## Workflow Review` section in the Goal's `Release.md`, including "no change, because…", and to route any real change to an Objective or Issue.
- The review keeps its scope: workflow skills, shared references, AGENTS.md routing, scaffold documents, against the Goal's REPLAN REQUIRED Tasks, NEEDS WORK Checks, Issues and lessons carried in. In a package-receiving project it still tunes only project-owned files and records skill suggestions as Issues.
- The scaffold copy stays byte-identical (TPL-01), content tests pin the new wording and forbid the old "retrospective Objective" wording, and `upgrade-assets` delivers it.
- G-002's own review is recorded in `.savepoint/releases/G-002-skills-optimisation/Release.md` (written by the planner, 2026-10-10).

## Architectural Considerations

Guidance-only change. No runtime, field, state or command change: the zero-Task routing stays as it is, and the review note is free Markdown that the Goal loader already tolerates.

## Confirmed Decisions — 2026-10-10

Owner confirmed through the Objective Decision Interview:

1. The retrospective becomes a note on the Goal, not an Objective, Task or Check. Rejected: dropping reviews entirely (repeat problems go unseen); allowing a zero-Task Objective to close (runtime change, still a full Check).
2. O-042 is repurposed to carry this rule change as its one real Task, rather than deleted.

## Boundaries

**In scope:** the Goal Workflow Retrospective section of `savepoint-design` and its scaffold copy, the content and upgrade tests that pin it, and G-002's review note.

**Out of scope:** runtime changes to Objective closing, new fields, states or commands, a Goal-level Check, and editing past retrospective Objectives such as O-038.

## Readiness and Verification

Dependencies O-039, O-040, O-043 and O-044 are done. Verification: `make build && make test-fast` at Task handoff, then the mandatory Full Objective Check with fresh `make test-full` under `agent-skills/references/check-method.md`. No platform-sensitive code is touched, so no native Windows evidence is needed beyond the normal CI run.
