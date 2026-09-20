---
type: epic-design
status: audited
---

# E44: Converge follow-up and verify integrated outcomes

## Purpose

Give every durable follow-up one record with a stable identity, a proof requirement, and an honest resolution, and decide Objective completion from an independent integration Check rather than from its Tasks being individually done.

This is the V1 delivery epic for V2 product Objective O004. The current V1 lifecycle remains authoritative while it is built: E44 makes the V2 decision canonical in `internal/data`, it does not move any live workflow, board surface, or skill onto it.

## What this epic adds

- Mutable Issue records under `.savepoint/issues/`, one file per durable follow-up, with a global `I###` identity that survives repair, recheck, and reopening.
- A descriptive Issue `type` — `defect`, `drift`, `guardrail`, `verification`, `other` — that classifies follow-up without ever deciding a gate by itself.
- A resolution contract that distinguishes `verified` (proof Check required), `accepted` (explicit owner decision, not a repair), and `duplicate` (canonical Issue named, not proof of anything), so closure cannot be claimed by setting a status.
- An append-only history of observations, repair attempts, rechecks, deferrals, reopenings, and owner decisions, where deferral is a dated entry on an open Issue rather than a fourth lifecycle state.
- Resolvable links in both directions: a Check's `issues` references stop being identity-only strings and resolve against real Issue records, and an Issue names the Tasks carrying its repair.
- Objective integration gating: an Objective closes only when its owned Tasks are complete **and** its own Objective-scoped Check clearance is current, or an explicit recorded exception applies.
- Objective dependency readiness, so a Task cannot start while its owning Objective waits on a dependency Objective that has no current integration clearance — and so an unrelated Objective is never blocked by one.
- Derived Issue listings and counts computed from the index, replacing the V1 audit register's separately maintained summaries.
- Named diagnostics for malformed Issues, dangling links, contradictory resolutions, and Objectives whose recorded status disagrees with their integration evidence, all surfaced through doctor.

## Components and files

| Module | Purpose |
|--------|---------|
| `internal/data/issue_v2.go` (new) | Define and strictly decode the Issue record — identity, type, status, source, links, severity, resolution, `duplicate_of`, history — and own the derived listings and counts consumers read instead of a register. |
| `internal/data/objective_gate_v2.go` (new) | Evaluate Objective integration completion, Objective dependency readiness, and Objective/Issue consistency, keeping Task gating in `gate_v2.go` untouched. |
| `internal/data/discover.go` | Add confined `.savepoint/issues/` discovery through the existing `v2PathConfiner`, matching the filename/ID and duplicate-identity rules `DiscoverV2Checks` already uses. |
| `internal/data/project.go` | Extend `V2Index` with Issues keyed by `I###`, the Task→Issue and Check→Issue link maps, Issue reference validation at load, and next-unused `I###` allocation. |
| `internal/data/check_v2.go` | Unchanged schema. Its `issues` field keeps decoding as `I###` references; E44 resolves them at index time rather than adding a second Check contract. |
| `internal/data/gate_v2.go` | Extend `ResolveTaskStart` with owning-Objective dependency readiness and add the matching blocker kind; clearance and completion rules are otherwise unchanged. |
| `internal/data/write.go` | Create Issue records, patch their mutable managed fields, append history entries, and patch Objective evidence — all through the existing preserving write path. |
| `internal/data/errors.go` | Add the named Issue and Objective-gate diagnostics doctor reports. |
| `internal/doctor/checks.go`, `internal/doctor/repairs.go` | Report the new diagnostics by stable name with manual repair guidance, and report Issue posture from derived counts. No new parsing, no policy, no automatic Issue creation. |
| `internal/data/*_test.go`, `internal/doctor/*_test.go` | Cover Issue decoding, identity and links, resolution dispositions, append-only history, reopening, Objective integration and dependency gates, consistency diagnostics, preserving writes, and doctor integration on temporary projects. |

## Architectural delta

E42 established structure and E43 established evaluation of a single Task. E44 adds the two things neither could express: a follow-up record that outlives any one Check, and a decision about an Objective as a whole.

`LoadV2Index` gains a fourth record family. Issues are discovered from `.savepoint/issues/{I###-slug}.md` — flat, like `checks/` — keyed by global `I###` ID, using the same confinement, filename-prefix, and duplicate-identity rules. Structural problems keep failing closed at load; evaluation still never returns an error for a recoverable state.

Issue frontmatter is `id`, `title`, `type`, `status: open|in_progress|resolved`, `source`, `tasks: [T###]`, `checks: [C###]`, `guardrail_ids`, optional `severity`, optional `resolution`, optional `duplicate_of: I###`, and `history`. Issue status is its own vocabulary and never overlaps Task lifecycle values: `planned`, `done`, and any stage are rejected, not healed. `source` records what produced the Issue — a Check, a direct report, or migration — with its actor and time, so an Issue always names its origin. `guardrail_ids` are opaque policy strings validated for shape only; Savepoint does not parse Guardrails.md to confirm them.

Links resolve in both directions at load, and one direction owns the truth. A Check is immutable, so its `issues` list is authoritative for which Issues that evaluation recorded; an Issue's `checks` list is the mutable mirror written in the same operation that creates the Check. Load validation requires every reference on both sides to exist, and requires a Check naming an Issue to appear in that Issue's `checks`. An unpaired link is a named diagnostic naming both records, never a silent reconciliation. `tasks` names the Tasks carrying repair, including a new bounded Task created for an out-of-scope fix, and its targets must exist.

Resolution is where closure is decided, and it is structural rather than narrative. `resolution` is required when status is `resolved` and rejected otherwise, so returning an Issue to `open` must clear it. It carries a disposition with its own obligations: `verified` requires a proof `check: C###` that exists, recorded CLEAR, and appears in the Issue's `checks`; `accepted` requires an owner actor and a reason and must not name a proof Check, because accepting risk is not repair; `duplicate` requires `duplicate_of` to name an existing, different Issue and must not name a proof Check. `duplicate_of` links are validated for existence, self-reference, and cycles exactly as dependency graphs are. Reopening is the same ID with status back to `open` and a dated history entry — a second Issue for a recurring problem is a mistake the identity rules make visible, not a supported workflow.

History is an append-only list of `{at, actor, kind, note, check}` entries in frontmatter, where `kind` is one of `observed`, `repair_attempted`, `rechecked`, `deferred`, `reopened`, `owner_decision`. The managed write may only append: a write that would shorten, reorder, or edit an existing entry is refused. **This refines `v2-Design.md` section 6, which placed History in the authored body.** History must be machine-appended, ordered, dated, and provably append-only — none of which a prose section can support — so the body keeps Summary, Evidence, and Proof Needed as authored content and does not carry a second history. This deviation is reconciled into the release design as part of this epic, not left as an undocumented divergence.

Issues never gate directly. Material blocking is expressed by a Check recording NEEDS WORK, and that Check's result is what E43's clearance rules already read; an Issue's type or severity contributes nothing to any gate decision. An Issue opened after a CLEAR Check does not retroactively block the Task — it makes a new Check necessary, which agents record by reassessing freshness. This keeps "type is descriptive" a property of the code rather than a convention.

Objective completion is the new decision. `ResolveObjectiveCompletion` requires every Task owned by the Objective to be `done` — ownership read from `index.ObjectiveTasks`, never from a manually maintained membership list — and the Objective's own clearance, resolved by the existing `ResolveClearance` over its `scope.kind: objective` Checks, to be `current`. Missing, needs_work, stale, and unknown integration clearance each block with a distinct reason, owner acceptance applies when the Objective declares it, and a recorded exception naming the Objective's latest Check grants completion by exception under owner authority — allowed-by-exception, never reported as clearance. Cross-Task repair goes back through Tasks; an Objective Check never closes a Task.

Objective dependency readiness is what makes integration clearance mean something downstream. `depends_on: [O###]` is satisfied only when the dependency Objective is `done` with current integration clearance; Task-only clearance inside that Objective is explicitly not enough. `ResolveTaskStart` gains one additional requirement — the owning Objective's dependencies must be satisfied — reported as a distinct blocker kind so consumers can explain that the wait is at the Objective level. Evaluation consults only declared dependencies, so an Objective with no dependency on a blocked one keeps running.

Consistency inspection extends to the new records without rewriting any of them: an Objective that is `done` without current integration clearance, an Objective `done` while an owned Task is not, an Issue resolved as `verified` whose proof Check is no longer CLEAR or has been superseded, and a Check-to-Issue link present on only one side. As in E43, Savepoint names the mismatch and states its limit — it does not authenticate who wrote a record or prevent external edits.

Reference: `.savepoint/releases/v2/v2-Design.md`, sections 5 and 6.

## Boundaries

**In scope:**

- Issue record schema, confined discovery, index, link and duplicate-graph validation, derived listings and counts.
- Resolution dispositions, proof requirements, reopening, and append-only history with preserving and append-only writes.
- Objective integration completion, Objective dependency readiness, and the owning-Objective requirement added to Task start.
- Objective evidence writes reusing the E43 evidence patch path.
- Named diagnostics for malformed Issues, dangling or unpaired links, contradictory resolutions, and Objective/Task status-versus-evidence mismatches, reported through doctor.

**Out of scope:**

- Migrating V1 defects, findings, audit runs, or the audit register into Issues; E45 owns it. E44 fabricates no Issue for historical work, and `internal/data/defect.go`, `audit_finding.go`, `audit_register.go`, `audit_run.go`, and `audit_backlinks.go` are untouched migration inputs. This narrows the previous E44 stub, which named those V1 files as edit targets before E42 and E43 settled the loader and Check boundaries.
- Automatic or probabilistic deduplication. Matching a new observation to an existing Issue is an agent's explicit search-then-link decision; Savepoint validates the link it is given and never merges records on similarity.
- Requiring an Issue for every Check observation. Advisory findings stay in the Check body; only durable follow-up earns an ID.
- A generic issue tracker: no assignees, no priorities beyond the optional severity, no workflow states beyond open/in_progress/resolved, no cross-project queries.
- Rewiring live consumers. The board's defect and audit-register overlays, router, and Next still read V1 records; E48 and E49 adopt these APIs when their consumers move to V2. E44's obligation is that the decision exists once, in data.
- Automatic freshness detection or repository scanning. An Issue opened after a CLEAR Check does not automatically restate freshness; an agent records that reassessment, exactly as in E43.
- V2 skills and shipped guidance describing Issue capture and the checker's Issue authority; E46 owns them.
- Activating any V2 policy change. Owner-only completion under DATA-05 remains the active V1 rule until cutover in E50.

## Dependencies

- E43-task-check-gates supplies the Check record family, shared evidence decoding, clearance resolution, gate decision and blocker contracts, consistency inspection, and the preserving and create-only write paths this epic extends.

## Quality gates

- Issue decoding tests cover valid records, malformed and duplicate `I###` IDs, a filename that does not match its declared ID, path traversal and symlink escape, every valid `type`, and rejection of Task lifecycle vocabulary in `status` without healing.
- Resolution tests prove `verified` without an existing CLEAR proof Check is refused, `accepted` requires an owner actor and reason and never names a proof Check, `duplicate` requires an existing, non-self `duplicate_of` target, `resolution` present on an open or in_progress Issue is a diagnostic, and a `duplicate_of` cycle fails closed.
- History tests prove entries append in recorded order, a write that shortens, reorders, or edits an existing entry is refused, a deferral is a dated entry on a still-open Issue, and reopening the same ID records a dated entry and clears the resolution.
- Link tests prove references resolve in both directions, a Check naming a missing Issue and an Issue naming a missing Task or Check each fail closed with a named diagnostic, and a Check-to-Issue link present on only one side is reported naming both records.
- Blocking tests prove no Issue type, severity, or open count changes any Task gate decision, and that a NEEDS WORK Check referencing an Issue blocks through E43's existing clearance rules alone.
- Objective completion tests cover an Objective with an unfinished Task, with all Tasks done but missing, needs_work, stale, or unknown integration clearance, with current clearance closing under checker authority, with declared owner validation waiting on acceptance, and with an exception closing by exception rather than reporting clearance.
- Objective dependency tests prove a dependency Objective that is done with current integration clearance satisfies readiness, a dependency whose Tasks are individually cleared but whose integration Check is absent does not, a Task cannot start while its owning Objective waits, and an Objective declaring no dependency on a blocked Objective is unaffected.
- Consistency tests prove a done Objective without current clearance, a done Objective with an unfinished Task, and a `verified` Issue whose proof Check was superseded each produce named diagnostics without rewriting any record.
- Write tests prove Issue creation is create-only over next-unused ID allocation, managed field patches preserve unknown frontmatter, authored bodies, and line-ending form, Objective evidence patches reuse the E43 path, and a no-op write leaves the file untouched.
- Doctor tests prove every new diagnostic reports under a stable name with a manual repair suggestion, Issue counts are derived from the index rather than a stored summary, an open advisory Issue does not make a project unhealthy, and doctor writes no project files.
- Implementation handoff requires focused `internal/data` and `internal/doctor` tests plus `make build && make test`, with the named cases and outcomes recorded in the Tasks.
- Epic closeout requires a fresh independent V1 epic audit; this planning session is not that audit. Current Guardrails apply, with STYLE advisory.

## Open decisions

None. Exact Go identifiers may be refined during task breakdown, but the Issue schema, resolution obligations, append-only history location, link authority direction, Objective integration and dependency rules, and the consumer boundary above are fixed. The history relocation from the body to frontmatter is a deliberate, recorded refinement of `v2-Design.md` section 6 and must be reconciled there during this epic.
