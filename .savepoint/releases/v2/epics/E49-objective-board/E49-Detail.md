---
type: epic-design
status: audited
---

# E49: Operate V2 through a clear terminal board

## Purpose

Give a V2 project a board a person can act from: Objectives, Task outcomes, Check evidence, owner waits, and one unified Issues view, in the language of the work rather than the language of the lifecycle — and let the keys on that board obey the gates instead of re-deciding them.

This is the V1 delivery epic for V2 product Objective O009. Two facts make it urgent rather than cosmetic. First, E47 made the V2 scaffold the `init` default, and `internal/board` still traverses `releases/epics/tasks`: a project created by `savepoint init` today has no `releases/` directory, so `savepoint board` fails at `data.Discover.ListReleases` with `releases directory not found`. The default new project cannot open its own board. Second, E48 built the shared interpretation — `data.ResolveNext` — and proved it through exactly one consumer. Section 11 of the release design requires board and resume to share it. Until a second consumer calls it, "shared" is an intention rather than a property.

## What this epic adds

- Schema dispatch at the board's entry point: a V2 project runs a V2 board; a V1 project runs today's board unchanged.
- An Objective sidebar and selection in place of the release picker and epic panel, with each Objective's status and its own integration clearance visible.
- Three columns — Planned, In Progress, Done — carrying the Task's human title, with implementation stage, clearance state, owner wait, dependency wait, replan, and done-by-exception shown as badges rather than as new columns.
- A prominent Next area rendered from `data.ResolveNext`: the same call, over the same input, that `savepoint resume` makes.
- One Issues overlay over `index.Issues`, filtered by type, with Defect as one of those types rather than a separate collection and a separate overlay.
- Task and Objective detail carrying the Outcome, the recorded evidence and its freshness basis, and the Check history for that target in recorded order.
- Gate-aware actions: a key exists only where a gate decision grants owner authority, and a refused action names the unmet requirement and whose authority it needs.
- V2 file watching over `objectives/`, `checks/`, `issues/`, `router.md`, and `config.yml`, with reload through explicit commands.
- Plain, deterministic non-TTY output for a V2 project, leading with Next.
- A readable load-diagnostic screen: a project that `LoadV2Index` refuses is reported as invalid project data, naming the file and the problem, with no partial board and no fallback.

## Components and files

| Module | Purpose |
|--------|---------|
| `internal/board/board.go` | The dispatch seam: resolve the project root once, call `data.LoadProject`, then run the V1 board or the V2 board. The only V1 board file this epic edits. |
| `internal/board/v2/` (new package) | The whole V2 board. Model, reducer, layout, rendering, overlays, async I/O, and plain output, with no V1 board type reachable from it. |
| `internal/board/v2/model.go` | V2 model state: the loaded `*data.V2Index`, the decoded `*data.RouterStateV2`, injected migration state, the resolved `data.Next`, and the focus/selection/scroll state of each surface. |
| `internal/board/v2/load.go` | Startup and reload `tea.Cmd`s: `data.LoadProject`, `ReadStateV2`, `migrate.PendingOperation`, `ResolveNext`, returning one typed message. No load work in `Update` (ARCH-02). |
| `internal/board/v2/update.go` | The reducer. Key dispatch by active surface: overlay, then sidebar, then columns. |
| `internal/board/v2/view.go` | Layout assembly and terminal-size decisions: sidebar, Next area, three columns, status bar. |
| `internal/board/v2/card.go`, `column.go` | Column and card rendering; the card's label is `TaskV2.Title`. |
| `internal/board/v2/badges.go` | The single mapping from typed `data` values — `ClearanceState`, `GateBlockKind`, `ProgressStage`, `FreshnessState`, exception and replan presence — to glyph, label, and accent. One file, so no second copy of the vocabulary exists (STYLE-09). |
| `internal/board/v2/next_panel.go` | The Next area, rendered from one `data.Next` value. |
| `internal/board/v2/objectives.go` | Objective sidebar: list, selection, per-Objective status and clearance, and Task filtering by ownership from `index.ObjectiveTasks`. |
| `internal/board/v2/detail.go`, `checks.go` | Task and Objective detail; Check history from `index.ScopeChecks` in recorded order, with `index.LatestCheck` marked. |
| `internal/board/v2/issues.go` | The Issues overlay: listing, type filter, and Issue detail including append-only history. |
| `internal/board/v2/actions.go` | Which key exists for the focused record, derived from the gate decision, and the refusal wording when none does. |
| `internal/board/v2/io.go` | Write commands: router selection and owner acceptance, both behind the pending-migration guard and both mtime-guarded. |
| `internal/board/v2/watch.go` | The V2 watch set and the debounced reload trigger. |
| `internal/board/v2/plain.go` | Non-TTY rendering. |
| `internal/data/write.go` | `WriteRouterStateV2`: the V2 router selection writer this epic needs and `internal/data` must own (DATA-02, STYLE-07). |
| `internal/data/next.go` | A nil-input guard on `ResolveNext`, carried forward from the E48 audit's non-blocking observation now that the second consumer exists. |
| `internal/data/gate_v2.go`, `objective_gate_v2.go`, `evidence_v2.go`, `project.go` | Read, not changed. Every readiness statement the board makes is a value obtained from these. |
| `cmd/board.go` | Adds `--objective`; argument parsing only, no schema knowledge (ARCH-01). |
| `internal/styles/styles.go` | Badge and clearance style tokens, in the existing palette. |
| `AGENTS.md` | Codebase Map rows for `internal/board/v2/` and the revised `internal/board/` responsibility (ARCH-04). |

## Architectural delta

**The V2 board is a sibling package, not a second personality inside the current one.** `internal/board` is fifty files and twelve thousand lines of release/epic/defect/audit-register presentation that E50 deletes. Adding a V2 board beside it in the same package forces every new symbol to carry a `v2` prefix to avoid colliding with the V1 one it replaces — `v2Model`, `v2View`, `renderV2Card` — and leaves nothing preventing a V2 render path from reaching a `data.Task`, a `data.Defect`, or an `AuditRegisterSet` that has no V2 meaning. A package boundary is the only mechanism Go offers that makes that unreachable rather than merely discouraged, and it is worth having precisely because "two parallel interpretations of the same project" is the failure this release keeps guarding against. So `internal/board/v2` is its own package, it uses the natural names E50 keeps, and at cutover E50 deletes the V1 files and promotes it. `internal/styles` is shared by both; nothing else is.

**Dispatch happens once, at the entry point, on the same signal every other command uses.** `board.Run` resolves the root, calls `data.LoadProject`, and branches on `Project.SchemaVersion` — the same dispatch `resume`, `doctor`, and `upgrade-assets` already perform. The V2 board therefore never sees a V1 project and never needs a fallback, and the V1 board is edited in exactly one file. `--release` and `--epic` are V1 filters and `--objective` is the V2 filter; a filter that does not apply to the resolved schema is a named error naming both, not a silently ignored flag (CFG-01).

**A load diagnostic is a screen, not a crash, and never a fallback.** `LoadV2Index` fails closed on duplicate identity, unsafe paths, missing ownership, dangling references, and — relevant here — a Task with no `title`, which `DecodeTaskV2` rejects outright. This is what makes the epic's title requirement structural rather than a rendering rule: no titleless Task can reach a card, because no project containing one loads. The board's obligation is the other half — to show that diagnostic as readable text naming the file and the problem, to render no partial board, to exit nonzero in plain mode, and under no circumstance to retry through V1 discovery. Falling back would turn a named data defect into a board that quietly displays something else (DATA-03, FS-06).

**Everything on screen is a record field or a resolver value; the board derives no readiness of its own.** Columns come from `TaskV2.Status`. Badges are a presentation mapping over `ResolveClearance`, `ResolveTaskStart/Advance/Completion`, `ResolveObjectiveCompletion`, `ResolveObjectiveDependency`, and the recorded evidence block — the board chooses a glyph for `ClearanceStale`, it does not decide that clearance is stale. Objective membership comes from `index.ObjectiveTasks`, built from each Task's own `objective` field. Check history comes from `index.ScopeChecks`. Issues come from `index.Issues` and the link maps. The badge vocabulary lives in one file so a second copy cannot drift, and it introduces no lifecycle word of its own (DATA-02). This is the same constraint E48 placed on the projection, applied one layer up.

**Next is the same call, not the same idea.** The Next area is `data.ResolveNext` over a `NextInput` built from the loaded index, the decoded V2 router, and `migrate.PendingOperation` — byte-for-byte the input `main.go` builds for `resume`. The board formats that value; it does not consult the index to second-guess it, and it does not read the router's `next_action` prose to decide anything. The proof obligation this creates is specific and testable: for a fixture project, the board's Next area and `savepoint resume` must report the same rung, the same selected record, and the same action. That test is what makes section 11's "same interpretation" a property of the code rather than a claim in a document.

**The board is the owner's surface, and the gates decide which keys exist.** This is the sharpest break from V1, where any focused card could be advanced or retreated with a keypress and the board checked dependencies itself. In V2, `GateDecision.Actor` names the authority a decision grants: executor for start and advance, checker for a completion with no remaining blockers, owner for a completion allowed by a recorded exception. The person at the board is the owner, so the board offers a write only where the decision names owner authority, plus the one write that is unambiguously the owner's to make — recording acceptance of a current Check through `WriteTaskEvidenceV2` / `WriteObjectiveEvidenceV2`. Starting a Task, advancing a stage, writing a Check, and resolving an Issue are shown, never pressed; the key reports whose session performs them and which requirement is unmet, reading the blockers straight out of the decision. A Task whose completion is allowed with checker authority is therefore displayed as clear to close by the checker rather than as a board action. That is restrictive on purpose: if it proves wrong in use, it is a question about the gate's recorded authority, which E43 owns, and it gets raised as an Issue rather than routed around in a render path.

**Router selection writes the selection and nothing else.** The V2 board's priority key sets `objective` and `task` in the router's `## Current state` anchor through a new `data.WriteRouterStateV2`, mtime-guarded and unknown-field-preserving exactly as `WriteRouterState` is (DATA-01). It does not write `state`, because `state` names which skill owns the conversation and the board is not running one. It does not write `next_action`, because that field is human prose and a board-authored sentence would sit where the projection already derives the real answer. The board displays the derived Next rather than the router's prose, so a `next_action` left behind by an earlier session cannot mislead from the board — which is also why not maintaining it is safe here.

**Writes go through commands, behind the migration guard, and survive a conflict.** Every write is a `tea.Cmd` returning a typed message (ARCH-02), refuses first if `migrate.PendingOperation` reports an incomplete conversion — the same boundary the V1 board, `upgrade-assets`, and `doctor` use — and detects a file changed underneath it through the `V2SourceDocument` freshness check `writeV2Record` already performs, surfacing a refresh-and-retry message rather than a partial overwrite (FS-01, FS-04).

**The watch set follows the V2 file model.** `objectives/` recursively, `checks/`, `issues/`, `router.md`, and `config.yml`. Not `releases/`, not `audit/`, not `defects/` — those directories have no V2 meaning, and watching `archive/` would make historical bytes trigger live reloads. A reload rebuilds the index through the same load command startup uses, so there is one load path rather than a startup path and a refresh path that can disagree.

**Plain output leads with the answer.** Non-TTY rendering prints Next first, then the three columns with title and badges, then a one-line Issues summary by type. V1's `hasAuditProposals` scan over `releases/*/epics/*/E##-Audit.md` has no V2 analogue and is not carried across; the V2 signal a project has outstanding follow-up is its open Issues, which the summary already states. Output is deterministic, ANSI-free without a TTY, and readable at a narrow width.

**This repository keeps the V1 board until E50.** Savepoint's own project is `schema_version: 1`, so the dispatch sends its maintainers to the V1 board throughout this epic. The V2 board is exercised against temporary fixture projects, never against the live project (TEST-04). That is the accurate state of a mid-release cutover, not a coverage gap.

Reference: `.savepoint/releases/v2/v2-Design.md`, sections 4, 5, 6, 10, and 11; `.savepoint/visual-identity.md` and its terminal adaptation appendix; `agent-skills/bubbletea-tui-design/SKILL.md`.

## Boundaries

**In scope:**

- The `internal/board/v2` package and the dispatch seam in `internal/board/board.go`.
- Objective sidebar and selection, three columns with badges, the Next area, Task/Objective detail with Check history, and the Issues overlay with type filtering.
- Gate-aware action availability and refusal wording, owner-acceptance writes, and router selection writes.
- `data.WriteRouterStateV2` and the `ResolveNext` nil-input guard.
- V2 watching, reload, conflict handling, and the pending-migration guard on writes.
- V2 plain output and the load-diagnostic screen.
- `--objective` parsing and the filter/schema mismatch diagnostic.
- Badge and clearance style tokens in `internal/styles`, within the existing palette.
- Codebase Map rows for the new package and the revised `internal/board` responsibility.

**Out of scope:**

- Any change to a gate rule, clearance derivation, dependency semantics, acceptance authority, or the precedence ladder. E43, E44, and E48 own them; the board consumes them unchanged. The nil-input guard is a defensive boundary check, not a semantic change.
- Authoring a Check or an Issue from the board, starting or advancing a Task from the board, or any other write outside owner authority.
- Changes to `internal/board`'s V1 rendering, overlays, actions, or tests beyond the dispatch seam.
- Removing the V1 board, the V1 router reader, or the transitional V1 readers. E50 owns that.
- Migrating this repository onto V2. E50 owns that.
- New columns for lifecycle flags, mouse support, cloud or team features, or visual redesign beyond the Atari-Noir identity.
- Doctor, resume, migrate, or init behavior changes.
- A V2 audit-register or defect surface; Issues with `type: defect` is the replacement, and historical register content stays in the archive.

## Quality gates

- A project created from `templates/project-v2` opens its board, with no Objectives and no Tasks, showing empty columns and the Idea/Design planning Next — no error and no diagnostic.
- Board Next and `savepoint resume` report the same rung, the same selected Objective/Task, and the same action for the same fixture project, proven across at least the rungs a board user reaches: execute, check-needed, owner-validation, dependency, replan, objective-integration, ready, and plan.
- A router naming a missing, archived, or mismatched record shows the selection diagnostic and still shows an available next action; no similarly numbered record is selected in its place.
- Every Task card shows its `title`; a project containing a Task with no title does not load, and the board shows that load diagnostic by name with the file path, renders no columns, and exits nonzero in plain mode.
- Each badge state renders distinctly and is distinguishable without color: stage build/test/audit, clearance missing/needs_work/stale/unknown/current, owner wait, dependency wait, replan, and done-by-exception, with done-by-exception visibly distinct from an ordinary done and stale completion visibly needing attention.
- An action key exists only where the gate decision grants owner authority; a refused action names the unmet requirement and the authority it needs, taken from the decision's blockers, and writes nothing. Repeated presses are idempotent.
- The Issues overlay lists every Issue, filters by each of the five types, shows Defect as a type rather than a separate collection, and renders an Issue's append-only history in recorded order.
- Check history for a Task and for an Objective renders in recorded order with the latest Check marked and superseded entries visible as superseded.
- A router selection write sets only `objective` and `task`, leaves `state`, `next_action`, unknown fields, and the body byte-identical, and refuses without writing when the file changed underneath it or when a migration is pending.
- An owner-acceptance write records the accepted Check, preserves unknown fields and body, is a no-op on a second unchanged run, and refuses without a partial write on conflict or pending migration.
- Rendering is correct at narrow widths and with wide/combining characters: no accidental wrapping, stable borders, and focus that changes color and glyphs without changing geometry.
- Non-TTY output is deterministic, ANSI-free, and complete; the forced-color and monochrome fallbacks render every badge distinguishably.
- An edit to a watched file reloads through the load command and preserves focus and selection where the records still exist; an edit under `archive/` or `releases/` triggers no reload.
- Implementation handoff requires focused `internal/board/v2`, `internal/data`, and `cmd` tests plus `make build && make test`, with named cases recorded in the Tasks. Fixture projects are temporary directories; the live Savepoint project is never used as a mutable fixture. Epic closeout requires a fresh independent V1 audit; this planning session is not that audit. Current Guardrails apply with STYLE advisory; ARCH-01, ARCH-02, ARCH-04, DATA-01, DATA-02, DATA-03, FS-01, FS-04, FS-05, CFG-01, DEP-02, TEST-01, TEST-02, TEST-04, TEST-06, and TEST-08 are the load-bearing ones.

## Open decisions

None. The sibling-package boundary, the single dispatch point, the load-diagnostic-over-fallback rule, the derive-nothing constraint on rendering, the shared `ResolveNext` call and its parity obligation, owner-only action authority, selection-only router writes, the V2 watch set, and the retirement of the release/epic/defect/audit-register surfaces for V2 are all settled above.

Two decisions are recorded as deliberate rather than deferred. First, a Task whose completion decision is `Allowed` with `Actor: checker` gets no board key: the release design says an owner may accept a checked outcome, while `ResolveTaskCompletion` grants checker authority once no blockers remain, and reconciling those belongs to the gate, not to a render path. E49 shows the state and names the authority; if that is wrong in use it becomes an Issue against E43's model. Second, the board does not maintain `next_action`, so a stale sentence can persist in `router.md` — acceptable because the board displays the derived Next instead, and the field is prose the skills author.

Exact Go identifiers in `internal/board/v2` and the precise glyph assignments may be refined during task breakdown; Task Context Files must name exact paths at that point.
