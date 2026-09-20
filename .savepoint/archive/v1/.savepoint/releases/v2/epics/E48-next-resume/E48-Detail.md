---
type: epic-design
status: audited
---

# E48: Recover current work and the next action

## Purpose

Give a V2 project one derived, authoritative answer to "what am I doing, what is it waiting on, and what happens next" — and a read-only command that speaks it — so reopening a project does not require reading records by hand or trusting a hand-written router line.

This is the V1 delivery epic for V2 product Objective O008. E42 through E44 built the records and every gate decision over them: `ResolveClearance`, `ResolveTaskStart`, `ResolveTaskAdvance`, `ResolveTaskCompletion`, `ResolveObjectiveCompletion`, `ResolveObjectiveDependency`, and the consistency inspectors. What does not exist is the layer above them — the one that decides *which* record a user is looking at and *which* of those decisions matters right now. Today that judgment lives nowhere, so every surface would have to invent it. E48 builds it once, and resume is its first consumer.

## What this epic adds

- A shared Next projection in `internal/data`: one ordered precedence ladder over the already-resolved gate decisions, returning one selected Objective/Task, the reason it was selected, and the single next action.
- A V2 router reader for the `state: idea|design|task|check` anchor with `objective`, optional `task`, and a human `next_action` — read strictly, and treated as a selection hint that never contributes to a completion decision.
- Named diagnostics for a selection that cannot be honored: a missing, malformed, or archived `task`/`objective` is reported as itself and paired with the action that *is* available. A similarly numbered record is never substituted for it.
- `savepoint resume`: a read-only narrative rendering of that projection — selected work, its implementation/technical/owner state, recorded evidence and its freshness basis, relevant Issues, and Next.
- Honest evidence wording throughout: `missing`, `unknown`, and `stale` clearance are each said in their own words, and no surface claims to have verified anything it only read.
- Onboarding rather than an error for a project with no Tasks yet: a fresh V2 project resumes into Idea/Design planning.
- Practical health in doctor: malformed data, missing evidence, failing configured gates, and "semantic review still needed" are reported as four different things, and an advisory Issue backlog stops being a reason to call a project unhealthy.

## Components and files

| Module | Purpose |
|--------|---------|
| `internal/data/next.go` (new) | The Next projection: `ResolveNext` over a `V2Index`, a router hint, and injected migration state. Owns the precedence ladder and selection resolution; owns no gate rule of its own. |
| `internal/data/router_v2.go` (new) | Strict V2 router state decoding from the existing `## Current state` anchor. V1's `RouterState` is untouched and keeps its own vocabulary until E50. |
| `internal/data/gate_v2.go`, `objective_gate_v2.go`, `evidence_v2.go` | Read, not changed. Next calls these resolvers; it must not re-derive clearance, dependency, or acceptance anywhere. |
| `internal/data/project.go` | `LoadProject`/`LoadV2Index` stay the only loader. Next takes an index it is handed; it performs no discovery and no IO. |
| `internal/resume/resume.go` (new) | Deterministic rendering of one resolved projection to an `io.Writer`. No IO beyond the writer, no subprocess, no TTY requirement. |
| `internal/resume/evidence.go` (new) | The evidence and freshness phrasing shared by resume's sections, kept out of the layout code so the wording has one source. |
| `cmd/resume.go` (new) | Thin argument parsing, `--help`, optional project directory, and runner injection, matching `cmd/board.go`. No domain parsing. |
| `main.go` | `resume` dispatch, and the wiring that supplies migration state from `internal/migrate` to the projection. |
| `internal/doctor/report.go`, `checks.go` | Health reporting separated into the four categories above; advisory backlog demoted from "unhealthy". |
| `internal/migrate/operation.go` | Read, not changed. `PendingOperation` is the existing incomplete-migration signal the wiring passes in. |
| `AGENTS.md` | Codebase Map rows for `internal/resume/` and for `internal/data`'s new projection responsibility (ARCH-04). |

## Architectural delta

**One projection, computed above the gates and below every surface.** The gate resolvers from E43/E44 each answer a closed question about one named record: may this Task start, is this Objective's dependency satisfied, how fresh is this clearance. None of them answers "which record". That gap is the whole epic. `ResolveNext` sits directly above them: it selects, it orders, and it explains. It computes no gate rule itself — every statement it makes about readiness is a `GateDecision`, `Clearance`, or `ObjectiveDependencyDecision` it obtained from the existing resolver and carries forward intact. A second place that decides clearance is the failure mode this design exists to prevent (STYLE-07, DATA-02), so the projection holds resolver results as values rather than re-reading the fields they were derived from.

**The precedence ladder is total and ordered, and the order is the contract.** Incomplete migration or an invalid target → recorded replan → unsatisfied dependency or owner prerequisite → execute/verify the selected Task → fresh Check needed → required owner validation → Objective integration Check → next ready Task or Objective → plan the next Objective. Every reachable project state lands on exactly one rung, including the empty project, and the ladder returns a value rather than an error for all of them. Two rungs deserve their reasons stated: migration comes first because a project midway through conversion has records that mean something different than they appear to, and owner validation sits below fresh Check because asking the owner to accept work that has no current technical clearance inverts the authority model E43 established.

**Migration state is injected, not imported.** `internal/migrate` already imports `internal/data`, so `internal/data/next.go` calling `migrate.PendingOperation` would close an import cycle. The projection therefore takes migration state as an input value on its input struct, and `main.go` fills it from `migrate.PendingOperation` at the same point it resolves the project root. This keeps `internal/data` free of a dependency on the conversion tool, and it keeps the first rung of the ladder testable from a plain value rather than from a staged operation directory.

**Router selection is a hint with a resolution step, never a second source of truth.** The V2 router keeps its `## Current state` anchor, and `state`/`objective`/`task` are read strictly — an unknown `state` value is a named diagnostic, not a healed default (DATA-03). What the projection then does with the selection is the load-bearing part: it resolves the named IDs against the index, and it honors them only on an exact identity match. A `task` that names nothing, names an archived record, or disagrees with the `objective` it sits under produces a named selection diagnostic *and* the next action that is actually available — both, because a user whose router line rotted still needs to know what to do. Substituting a similarly numbered Task for a missing one is the specific behavior this rule forbids. `next_action` is human prose: it is displayed, it is never parsed, and it never decides anything.

**Resume reads; it does not act.** No writes, no evidence generation, no subprocess, no board process, no network. Its no-write guarantee covers the failure paths too — a malformed router, a missing project, an unresolvable selection, and a writer that fails mid-render all leave the project byte-identical and mtime-identical. It reports recorded evidence dates, the freshness basis, and what the owner is waiting on; it never restates a recorded Check as a verification it performed, and "no evidence recorded" is printed as that rather than as a pass or a failure. This is what makes resume safe to run reflexively on a project mid-flight, which is the only time anyone runs it.

**Resume serves V2 projects and says so for V1 ones.** `LoadProject` dispatches on `schema_version`, and building a second projection over V1 releases/epics would create the parallel interpretation this epic exists to avoid — for a lifecycle E50 removes. Pointed at a V1 project, resume names the schema it found and the route to V2, writes nothing, and exits nonzero. The live Savepoint repository is V1 until E50, so its own maintainers meet that message; that is the accurate state of a mid-release cutover, not a gap.

**The board is the next consumer, not this one.** Section 11 requires resume and board to share the interpretation, and the sharing is structural: the projection is a value in `internal/data` with no presentation in it, so the board can consume the same call. Adopting it in the TUI is E49's work, along with the badges and Objective sidebar that render it. E48's obligation is that nothing in the projection's shape or location makes that adoption require a second interpretation — no rendering, no terminal awareness, no resume-specific fields.

**Health becomes a report with categories instead of a verdict.** Doctor today collapses problems into "has problems". A V2 project with three advisory Issues open is working normally, and telling its owner it is unhealthy trains them to ignore the report. So the four categories — malformed data, missing evidence, failing configured gates, and semantic review still needed — are reported separately, and only the first three bear on whether the project is structurally sound. Doctor still creates no Issues and repairs no files.

Reference: `.savepoint/releases/v2/v2-Design.md`, sections 4, 5, 10, 11, and the section 13 Task example.

## Boundaries

**In scope:**

- `ResolveNext`, its input contract, the precedence ladder, and selection resolution with named diagnostics.
- Strict V2 router state reading, including the unknown-state diagnostic.
- The `internal/resume` package, `cmd/resume.go`, and `main.go` dispatch and migration-state wiring.
- Evidence, freshness, owner-wait, and Issue phrasing used by resume.
- The doctor health-category split and the demotion of advisory backlog.
- Codebase Map rows for the new package and responsibility.

**Out of scope:**

- Any change to a gate rule, clearance derivation, dependency semantics, or acceptance authority. E43/E44 own them; Next consumes them unchanged.
- Board rendering, badges, Objective sidebar, or any `internal/board` change. E49 owns them.
- Writing the router, writing evidence, creating Issues, or any other mutation from resume or doctor.
- Running configured quality gates, subprocesses, or model calls to produce evidence for a resume rendering.
- A V1 Next projection or V1 resume rendering beyond the named schema message.
- New Task states, new Check states, or a new Issue type.
- Removing the V1 router reader or the transitional V1 readers. E50 owns that.
- Migrating this repository onto V2. E50 owns that.

## Quality gates

- Every rung of the ladder is reached by a fixture project and returns the expected selection, reason, and next action — pending migration, recorded replan, unsatisfied Task dependency, Objective dependency wait, in-progress execution, needed Check, needs-work Check, stale clearance, unknown clearance, required owner validation, Objective integration Check, next ready Task, and nothing left to do.
- A fresh V2 project with no Objectives and no Tasks resumes into Idea/Design planning with no error and no diagnostic.
- A router naming a missing, archived, or mismatched Task reports the selection diagnostic by name and still reports an available next action; no similarly numbered record is selected in its place.
- A malformed router, an unknown `state` value, and a missing project directory each produce a named diagnostic and a nonzero exit, with no write.
- Resume leaves every project file byte-identical and mtime-identical, verified by snapshot before and after, on success, on every named failure path, and on writer failure mid-render.
- Missing, unknown, and stale evidence each render in distinct words, always naming the Check and the recorded freshness basis when one exists; no rendering asserts verification resume did not perform.
- Resume output is deterministic and complete without a TTY, and readable at a narrow width.
- A V2 project resolves the same selection and next action through the projection call regardless of consumer, with the projection carrying no presentation state.
- Doctor reports malformed data, missing evidence, failing gates, and pending semantic review as four categories, and a project whose only findings are advisory open Issues is not reported as structurally unsound.
- Implementation handoff requires focused `internal/data`, `internal/resume`, `cmd`, and `internal/doctor` tests plus `make build && make test`, with named cases recorded in the Tasks. Epic closeout requires a fresh independent V1 audit; this planning session is not that audit. Current Guardrails apply with STYLE advisory; FS-03, FS-06, DATA-02, DATA-03, ARCH-01, ARCH-03, ARCH-04, CFG-01, TEST-01, TEST-02, TEST-04, and TEST-08 are the load-bearing ones.

## Open decisions

None. The projection's location and layering, the injected migration input, the ladder order, selection-hint semantics, resume's read-only contract and V1 message, the deferral of board adoption to E49, and the doctor health split are settled above. The section 13 Task example is treated as contract direction, not as a plan to copy: its named APIs are confirmed to exist — `V2Index`, `ResolveClearance`, the three Task gate resolvers, and the Objective resolvers — and exact Go identifiers for the new projection and renderer may be refined during task breakdown.
