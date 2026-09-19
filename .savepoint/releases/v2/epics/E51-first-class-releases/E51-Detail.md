---
type: epic-design
status: planned
---

# E51: Make releases first-class in V2

## Purpose

Preserve Release as an optional but real planning and delivery boundary in V2. A small project can still work directly as Objective → Task, while a project that declares a Release can navigate its promised outcomes, preserve its source brief through migration, and prove whether the whole delivery is ready without reviving the V1 Release → Epic → Task hierarchy.

This closes a product gap in the original V2 simplification. Optional `release` text on an Objective can filter a list, but it cannot carry a stable identity, a release promise, integration evidence, owner acceptance, or an accountable migration destination for a V1 release PRD.

## What this epic adds

- An optional first-class V2 Release record with global `R###` identity, title, lifecycle, outcome, success conditions, release-level Check evidence, and explicit owner acceptance.
- Objective membership expressed once, by an Objective's optional `release: R###` reference. Release membership is derived; Releases do not own a second Objective list and Tasks remain owned only by Objectives.
- A release-completion decision that requires every member Objective to be complete, a current CLEAR Release Check, no unexcepted blocker, and owner acceptance of that Check. `done` means the recorded release outcome is accepted, not published or deployed.
- Release-aware project indexing, router selection, `data.Next`, board navigation, `resume`, doctor diagnostics, file watching, and deterministic plain output, all consuming the same canonical interpretation.
- Preservation of the existing board interaction: `r` opens the Release selector, starts on the current Release, supports keyboard navigation and cancel, and switches the visible delivery context on Enter.
- Migration of every V1 release into a V2 Release record, with release PRD content given a live destination, original bytes archived, and legacy identity/proof retained without fabricating a new V2 CLEAR result.
- Optional Release templates and agent-workflow guidance. Fresh projects do not receive a synthetic Release merely to satisfy structure.
- A release-readiness gate that E50 can consume before cutover, distinct from package publishing and deployment.

## Components and files

| Module | Purpose |
|--------|---------|
| `internal/data/release_v2.go` | Typed Release record, `R###` identity, lifecycle/evidence validation, source-preserving managed writes, and release-completion decisions. |
| `internal/data/project.go` | Discover and index Releases; derive `ReleaseObjectives`, Release Check history, and Objective-to-Release links without duplicating ownership. |
| `internal/data/check_v2.go`, `evidence_v2.go` | Extend Check scope and freshness resolution to `release`; retain one Check schema and one evidence vocabulary. |
| `internal/data/next.go` | Add selected Release context and release-integration/owner-acceptance outcomes to the shared projection without changing Task or Objective gate rules. |
| `internal/data/router.go`, `write.go` | Preserve an optional `release: R###` selection with the same unknown-field and mtime safety as Objective/Task selection. |
| `internal/migrate/` | Convert V1 releases and release PRDs, assign deterministic `R###` identities, preserve exact source bytes, map legacy references, and resume safely after interruption. |
| `internal/board/v2/` | Optional Release selector, filter, detail, readiness display, watch coverage, and release-aware Next/plain output while keeping Objectives as the work surface. |
| `internal/resume/` | Render Release context, integration evidence, and owner wait from `data.Next`; perform no independent release inference. |
| `internal/doctor/` | Diagnose duplicate/missing Release identity, dangling Objective references, invalid lifecycle/evidence, and migration/reference inconsistencies. |
| `agent-skills/savepoint-idea/`, `agent-skills/savepoint-design/`, `agent-skills/savepoint-check/` | Decide when a Release is useful, plan within it, and perform an independent Release Check without adding a fifth public phase. |
| `templates/project-v2/` | Ship matching optional Release record and workflow templates; do not create a default Release record for a fresh project. |
| `.savepoint/releases/v2/epics/E50-release-validation-cutover/` | Consume the recorded Release readiness and migration evidence as a cutover prerequisite. |
| `README.md`, `.savepoint/Design.md`, `.savepoint/releases/v2/v2-Design.md`, `AGENTS.md` | Reconcile the public workflow, implemented architecture, release-design delta, and codebase map when the epic is delivered. |

## Architectural delta

**Release is optional structure, not optional text.** The original V2 design removed mandatory releases so a small project would not need packaging ceremony. That product constraint remains. When no Release exists, Objective → Task, routing, Next, board, and resume behave exactly as they do after E49. When a Release is declared, however, it is a typed record under `.savepoint/releases/R###-slug/Release.md`; a title or path change does not change its identity. The first unused `R###` is allocated with the same no-reuse and duplicate-diagnostic rules as other global IDs.

**A Release owns a promise and evidence, not the work tree.** Release body sections are `Outcome`, `Why`, `Success Conditions`, and `Boundaries`. Its lifecycle is `planned|in_progress|done`, using the same lifecycle words as Objective and Task, but with release-specific gates. Objectives optionally name one Release through `release: R###`; Tasks continue to name one Objective and never name a Release. `ReleaseObjectives` is derived from the Objective records, so moving an Objective changes one authoritative field and no membership list can drift.

**Release completion is an integration and owner decision.** `planned` and `in_progress` describe delivery progress. A current Release can become `done` only when it has at least one member Objective, every member Objective satisfies its existing completion decision, the latest Release-scoped Check is current and CLEAR for the recorded Release outcome, linked material Issues are resolved or explicitly excepted, and the owner has accepted that exact Check. A migrated settled Release whose member Objectives are archived may instead preserve its completed V1 disposition only through a typed `legacy_completion` archive reference; that is historical proof, not current CLEAR evidence or a bypass for live members. The owner acceptance is mandatory at the current boundary because declaring the delivery promise met is a product decision. It does not mean published, deployed, tagged, or commercially released; those remain outside Savepoint's lifecycle.

**Checks gain a third scope, not a parallel release audit system.** `Check.scope.kind` becomes `task|objective|release`. Release Checks use the existing immutable record, supersession, reviewed-scope, Issue, freshness, and independence contracts. `internal/data` adds Release equivalents of the Objective integration/evidence resolvers instead of teaching the board, resume, doctor, or skills separate readiness rules. A failed Release Check creates or reuses ordinary Issues and routes repair to bounded Tasks in the owning Objectives.

**Selection is contextual and cannot hide global work.** The router may carry optional `release: R###`. A selected Objective must either belong to that Release or produce a named mismatch diagnostic; no similarly numbered or same-titled record is substituted. With a Release selected, Next chooses only its eligible Objectives and, after their completion, reports Release Check or owner-acceptance work. With no Release selected, the existing global precedence remains and each Objective is independently actionable. Board filtering changes visibility, not dependencies or gate outcomes.

**The board stays Objective-led and keeps the Release switch.** The existing `r` interaction remains the way to open the Release selector. It opens over the board with the current Release focused; `j`/`k` or arrow keys move, Enter selects and closes, and Esc or `q` cancels without changing context. Selecting a Release refreshes the visible Objectives and Tasks and persists the router's Release context without changing lifecycle state. The selector is a delivery-context filter alongside the Objective sidebar, not a new nested board. Release detail shows its outcome, member Objective progress, latest/superseded Checks, relevant Issues, freshness basis, and owner acceptance. Cards remain Tasks grouped only by status. Projects without Releases keep the shortcut inert or show the established `(none)` state without requiring a synthetic Release. File watching adds `releases/`; reload still goes through the single V2 load command, and plain output remains deterministic and leads with the shared Next result.

**Migration maps the boundary before schema activation.** Each V1 release is assigned a deterministic `R###` from the frozen source inventory, and every converted Objective points to it. The release PRD supplies the live Release outcome/source documentation; epic details continue to supply Objective scope. Content that cannot be mapped structurally remains in the live Release body under a preserved legacy-source section and in the exact-byte archive. The migration manifest maps source release path/name to `R###` and records the archived source hash.

Settled historical releases may enter `done` only with a typed legacy-completion reference to archived evidence; that reference is displayed as historical proof and never as a new CLEAR Check. Active or ambiguous releases enter `in_progress` with unknown Release clearance, and ambiguity is a preview-blocking owner decision. Unresolved defects/findings remain Issues and therefore cannot disappear behind a migrated release status. Preview, apply, retry, conflict, backup, and schema-activation behavior remain the recoverable migration protocol proven by E51 T009.

**E50 consumes readiness; E51 does not publish.** A migrated copy of this repository must expose its V2 Release record, member Objectives, historical mappings, open Issues, and current Release readiness before E50 may cut over the live repository. E50 uses the canonical Release decision plus its own distribution evidence. E51 adds no version-tagging, changelog generation, artifact upload, deployment, or package-manager behavior.

This epic supersedes the optional-string treatment in `.savepoint/releases/v2/v2-Design.md` sections 2, 3, 6, 9, 11, and 14. Those sections now describe the implemented Release contract; the current `.savepoint/Design.md` records the transitional implementation while the live repository remains V1 until E50 audit closeout.

## Boundaries

**In scope:**

- Optional V2 Release records at `.savepoint/releases/R###-slug/Release.md`, stable identity, lifecycle, authored source, and confined discovery.
- Optional `release: R###` ownership on Objectives and derived Release membership/index links.
- Release-scoped immutable Checks, evidence freshness, Issues, completion decision, mandatory owner acceptance, and historical-completion display.
- Optional router Release selection, selection diagnostics, and shared release-aware Next projection.
- Release navigation, filtering, detail, watching, and deterministic TTY/non-TTY presentation without changing Task columns.
- Compatibility with the current `r` Release-selector workflow, including current-selection cursor placement, keyboard navigation, cancel, Enter-to-switch, and router persistence.
- V1 release/release-PRD conversion, exact-byte archive, legacy mapping, ambiguous-state decisions, recovery, and idempotence.
- Doctor, skills, templates, documentation, and E50 cutover prerequisites needed to make the model usable end to end.

**Out of scope:**

- Making Releases mandatory, creating a default Release for a fresh project, or preventing unassigned Objectives.
- Reintroducing Epic, nesting Objective or Task files under a Release, or maintaining a second membership list.
- Changing Task, Objective, dependency, Issue, or existing completion semantics except to compose their results into Release readiness.
- Release dependencies, trains, milestones, roadmaps, calendars, budgets, teams, remote coordination, or cloud state.
- Version calculation, changelog generation, Git tags, artifact signing/upload, deployment, package publishing, or environment promotion.
- Treating a migrated historical disposition as a V2 CLEAR Check.
- Migrating the live repository before E51 is independently audited and E50's cutover gates pass.

## Quality gates

- A project with no Releases loads, plans, executes, checks, boards, and resumes exactly as before E51; no Release record or selector is required or synthesized.
- A project with two Releases and unassigned Objectives derives correct membership from Objective records alone, keeps unrelated work independently actionable, and diagnoses every missing or mismatched reference by file and ID.
- Release IDs remain stable across title/path edits, are never silently reused, and duplicate IDs or path escapes fail closed before partial indexing.
- Release completion is refused unless there is at least one member Objective, every member Objective is complete, the latest Release Check is current and CLEAR, material Issues are resolved/excepted, and the owner accepted that exact Check.
- A superseding Check or material change makes the prior Release acceptance stale; legacy-completed Releases remain visibly historical rather than CLEAR.
- Board, plain output, and `resume` report the same selected Release, readiness state, evidence basis, owner wait, and next action from one `data.Next` projection. Browsing remains byte- and mtime-identical.
- Selecting a missing, archived, or Objective-mismatched Release yields a named diagnostic and never substitutes another record. With no Release selected, the existing global Next precedence is unchanged.
- On the V2 board, `r` opens the Release selector over the board, focuses the current Release, supports arrows/`j`/`k`, cancels with Esc/`q`, and switches Release context with Enter; switching updates visible Objectives/Tasks and router selection without changing lifecycle records.
- Release detail and Check history remain readable at narrow widths, with wide/combining characters, monochrome output, and no TTY; Release edits reload through the canonical load path.
- Migration of fixtures with multiple releases, duplicate legacy Task IDs, active and settled work, release PRDs, unresolved defects/findings, and ambiguous status accounts for every source byte and reference in preview and apply.
- Migration dry-run writes nothing; interruption is recoverable; conflicts preserve user edits; a second unchanged run changes no bytes, mtimes, IDs, mappings, or evidence.
- A migrated copy of this repository loads as V2 and exposes accountable live/archive destinations for every release PRD and release-scoped reference before E50 considers live cutover.
- Focused `internal/data`, `internal/migrate`, `internal/board/v2`, `internal/resume`, `internal/doctor`, template/skill parity, and end-to-end tests pass with `make build && make test`; all mutable fixtures use temporary project copies.

## Dependencies

- E49-objective-board
- E45-safe-migration

## Open decisions

None. `R###`, `.savepoint/releases/R###-slug/Release.md`, optional Objective membership, derived membership, `planned|in_progress|done`, Release-scoped Checks, mandatory owner acceptance for Release completion, and no publishing semantics are fixed by this design.
