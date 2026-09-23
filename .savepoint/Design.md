---
type: project-design
status: active
last_audited: v2/E51-first-class-releases
---

# Savepoint — System Architecture

> Project-level architecture. Check-kept fresh: each Objective and Goal verification merges its reconciled delta into this document.

> **Current V2 architecture:** E50's owner-approved migration activated schema 2 in this repository. `internal/data` owns the identity-keyed V2 index, Task/Objective/Issue/Check gates, and first-class Goal completion through Release-compatible records; board, doctor, resume, and migration use those boundaries without a second readiness policy.
>
> **Historical V1 evidence:** Former V1 project records remain byte-preserved under `.savepoint/archive/v1/`. The V1 readers retained in `internal/data` are used by `internal/migrate`; they are not part of ordinary runtime commands.
>
> **Visual identity** lives separately in `.savepoint/visual-identity.md` and is loaded only for TUI/theme/visual tasks.

## 1. Architecture model

- **File-only.** No MCP server. Agents read and edit Markdown + YAML files directly using their native file tools.
- **Agent routing:** AGENTS.md → `.savepoint/router.md` → phase skills. See AGENTS.md Workflow section.
- **Bundled Agent Skills:** The active V2 workflow uses `savepoint-idea`, `savepoint-design`, `savepoint-task`, and `savepoint-check`, with three non-triggerable shared references. Retired V1 workflow files are preserved only as historical evidence.
- **Token-efficiency principle.**
  - Cold session bootstrap: ~5–7K tokens (one-time per conversation).
  - Per-task incremental: <2KB.
  - Audit: 5–15KB.
  - Anything that breaks these bounds violates the wedge.
- **Go data-reader boundary:** `internal/data` owns Savepoint file parsing and discovery for the Go implementation: Objective/Task/Check/Issue/Release models (the internal Release model backs the public Goal), markdown YAML extraction, V2 router state parsing, config/theme defaults, identity-keyed record discovery, lifecycle validation/defaulting, write-time status validation, and boundary error sentinels. Its `CheckRuntimeSchema` front door rejects schema 1 with the named migrate-required diagnostic, directing the owner to preview with `savepoint migrate --dry-run` and apply with `savepoint migrate --apply` without importing conversion code.
- **V2 evidence and identity boundary:** `internal/data` strictly loads identity-keyed Objective, Task, Check, Issue, and Release-compatible Goal records through confined paths, preserves authored record content on managed writes, resolves numeric Check history and freshness (a CLEAR Check signed by a checker is current on its own; an optional freshness assessment naming that Check can only mark it stale or unknown), and owns canonical dependency, lifecycle, authority, acceptance, exception, and replan decisions. `internal/doctor` reports the same structural and evidence diagnostics without rewriting project files.
- **V2 follow-up and integration boundary:** `internal/data` owns the mutable Issue family and bidirectional Check/Issue/Task links, while Objective and Goal completion remain derived from owned work, independent integration evidence, material Issue posture, and exact owner acceptance. The Goal is backed by the existing Release record and completion resolver. `internal/doctor`, `internal/board/v2`, and `internal/resume` report those decisions without duplicating policy. Issue resolution carries four dispositions — `verified` (Check-proven repair), `accepted` (owner risk decision), `duplicate` (points at a canonical Issue), and `escalated` (points at the Objective, `escalated_to: O-###`, the repair was promoted into) — each with its own proof obligation enforced by `internal/data`. Owner acceptance may close an Issue directly as `accepted` without a Check and without claiming technical `CLEAR`. Escalation is the Issue closure the planning workflow performs directly rather than `savepoint-check`: `savepoint-design` retires an Issue with disposition `escalated` at the moment it promotes that Issue's repair into a new Objective, since that Objective's own mandatory Full Objective Check and owner acceptance become the proof.
- **V2 agent workflow assets:** The four active skills and three non-triggerable shared references are byte-identical between the live and V2 scaffold trees. Task `planned_by` provenance and Check `executed_session` provenance are strict typed fields; the decoder rejects a Check that claims the executor and checker were the same session.
- **Configured build gate:** established in E46. `quality_gates.build` is decoded by `internal/data` and executed by `internal/doctor` after typecheck and before test, using the existing timeout and result-reporting path.
- **V2 Goal boundary (Release-compatible storage):** A Goal is optional context grouping related Objectives under a navigable outcome; it does not own Tasks or publish, deploy, tag, or generate changelogs. Existing Goals remain stable `R-###` records under `.savepoint/releases/<slug>/Release.md`; Objective and router references retain `release: R-###`, and Goal Checks retain `scope.kind: release`. Membership comes from Objective records, and completion uses the existing Release-scoped Check resolver, material Issues, and exact owner acceptance. Migration maps every V1 Release PRD to a live record plus a byte-preserved archive; historical completion is typed evidence, never a fabricated current Check.
- **Template assets** live under `templates/project-v2/` for the active workflow; the V1 scaffold tree and its retired skills are gone (O-021), and `savepoint migrate` converts a legacy project onto this same tree.
- **Init command** (`savepoint init`) validates targets and scaffolds `templates/project-v2/`, including the four V2 skills, Idea/Design/Guardrails/router files, and schema version 2. Existing user content is preserved through the managed-guide boundary.
- **Upgrade-assets command** (`savepoint upgrade-assets [dir] [--dry-run] [--force]`) refreshes package-owned V2 skills and shared references with provenance; migration history and project records remain untouched.
- **Board command** (`savepoint board`, and bare `savepoint`) loads the V2 index and router once, resolves the shared `data.Next`, and renders the Objective/Task/Check/Issue surface in TUI or deterministic non-TTY form. Goal context is optional and derived from Objective membership; V1 filters are refused for V2 projects.
- **Doctor command** (`savepoint doctor`) runs read-only V2 structure, lifecycle, dependency, Issue, evidence, quality-gate, and canonical Goal-readiness diagnostics through the existing Release resolver, with named repair guidance and no automatic writes.
- `internal/board/v2`, `internal/doctor`, and `internal/resume` consume the existing first-class Release records and canonical Release completion resolver for the public Goal context. The board's `g` key is the canonical Goal selector; `r` remains an undisplayed compatibility alias. Goal detail, the plain board, and resume expose that context.
- **Audit remediation baseline** (v1.1 E13) centralizes frontmatter/body splitting and line-ending normalization in `internal/data`, uses typed sentinel errors for doctor repair suggestions, applies a configurable `quality_gates.gate_timeout`, removes tracked build artifacts from source control, adds `.golangci.yml`, and moves board filesystem reads/writes behind Bubble Tea command messages while preserving direct file I/O inside command helpers.
- **Structural improvement baseline** (v1.1 E14) groups board `Model` fields into focused embedded state structs, defines consumer-side board/doctor data-access interfaces, routes doctor orphan discovery through `Discover.ListRootDirs`, renders audit-tab hidden sections via exact heading matches, improves quality-gate shell tokenization for quoted and escaped arguments, removes the separate `TaskStatus` enum in favor of `ColumnType`, and adds `internal/testutil` for shared Go test fixtures.
- **Hardening baseline** (v1.1 E15) adds board render/layout benchmarks, data frontmatter fuzz targets, debug logging via CLI `--debug` or `SAVEPOINT_DEBUG`, abbreviation-aware task checklist sentence splitting, root test package isolation, documented audit-tab hidden-section allowlisting, repo-local CI, `make ci`, distribution SHA256 checksums, and Windows amd64/arm64 build outputs.
- **Independent Check workflow** is skill-driven, not a CLI pipeline. A fresh `savepoint-check` session writes an immutable `C-###` record; an individual Task Check is optional and owner-waivable, while the Full Objective Check is mandatory and the Full Goal Check is mandatory whenever a Goal exists. Goal Checks retain serialized `scope.kind: release`; the executor cannot clear its own work or infer Issue acceptance and may only record the owner's explicit `accepted` decision.
- **Board-recorded Task-check waivers**: pressing Space on the V2 board to complete a Task at stage check that has no recorded Check at all (`ClearanceMissing`) is itself the explicit owner action TEST-09 requires — only a human at the interactive keyboard reaches that key, never an agent. The board auto-records the `check_waiver` block (task, reason, `actor: {role: owner, session: board-owner}`, time) in the same write that sets the Task done, rather than requiring the owner to hand-write that fact first. This never applies when a Check was actually recorded and found a problem (`needs_work`, `stale`, `unverified`): that result stands, and completion stays refused.

## 2. Directory layout

```
<project-root>/
├── AGENTS.md                       ← active V2 routing and policy entry point
├── agent-skills/                   ← four V2 skills plus shared references
└── .savepoint/
    ├── Idea.md                     ← product intent
    ├── Design.md                   ← current architecture (this file)
    ├── Guardrails.md               ← durable engineering policy
    ├── router.md                   ← V2 state and next action
    ├── config.yml                  ← schema_version, theme, quality gates
    ├── releases/                   ← compatibility storage for optional V2 Goals (R-###)
    │   └── R-###-slug/Release.md
    ├── objectives/                 ← Objective records and owned Tasks
    │   └── O-###-slug/
    │       ├── Objective.md
    │       └── tasks/T-###-slug.md
    ├── checks/                     ← immutable independent evidence
    ├── issues/                     ← durable follow-up records
    ├── archive/v1/                 ← byte-preserved historical source
    └── migrations/v1-to-v2.yml    ← source hashes and conversion mappings
```

Goal records are optional first-class V2 contexts, stored as Release records:
Objective membership is derived from each Objective's existing `release: R-###`
field, and Goal completion does not publish, deploy, tag, or generate changelogs.

AGENTS.md at root is the active cross-vendor guide. Design.md in `.savepoint/`
is the working architecture record. `visual-identity.md` is conditional and
loaded only for TUI/theme/visual tasks. V2 Tasks are owned by Objectives;
legacy E##/T## paths remain only under the preserved V1 archive. Scaffold
assets live under `templates/project-v2/`; generated projects receive rendered
copies, not hardcoded strings.

## 3. Hierarchy semantics

| Level        | Definition                                                                             |
| ------------ | -------------------------------------------------------------------------------------- |
| **Goal**     | Optional `R-###` context grouping Objectives under an outcome; stored in Release-compatible records. |
| **Objective**| A durable outcome with an optional Goal reference (serialized as `release:`) and owned Tasks. |
| **Task**     | Independently buildable work owned by exactly one Objective; requires an implementation plan before build. |
| **Check**    | Immutable independent evidence: optional Quick evidence for a requested Task Check, mandatory Full integration evidence for an Objective, and mandatory Full cross-Objective evidence for a Goal (`scope.kind: release`). |
| **Issue**    | Durable follow-up for a defect, drift, guardrail gap, or verification problem; its lifecycle is separate from Task status. |
| **Sub-task** | Inline checklist item — _evidence of the implementation plan_, not standalone work.    |

## 4. Status model & gates

Three statuses, with explicit gates and ownership boundaries:

| Status        | Meaning                    | Entry gate                                                      | Who may set it                         |
| ------------- | -------------------------- | --------------------------------------------------------------- | -------------------------------------- |
| `planned`     | Ready to build             | plan section non-empty                                          | User or planning workflow              |
| `in_progress` | AI building                | all `depends_on` are `done`                                     | Agent when starting implementation     |
| `done`        | Complete for current scope | all implementation items checked; verification per project mode | User only                              |

- `blocked` is a **flag**, not a status — `in_progress` + `blocked: "reason"` is valid.
- `internal/data` is the single owner of task lifecycle rules: canonical statuses, canonical stages, parse compatibility for legacy `phase`, write validation, and transition helpers must flow through that package. `Task.Column` and `Task.Stage` are the canonical in-memory lifecycle fields; no denormalized status mirror exists on the Task struct.
- Canonical task files and workflow guidance must use `status` plus `stage` only while `status: in_progress`; legacy `phase` is accepted only as read compatibility and should be reported as drift by diagnostics/templates.
- Agents may only advance a task into `in_progress`; they must not set `done` or retreat a task to an earlier status.
- Only the user may set a task to `done` or retreat it from `done` to `in_progress` when follow-up work is required.
- Router updates are explicit TUI actions: after setting a task to `in_progress`, the agent prompts the user to press `p` in the board to mark the focused task as router priority. Navigation alone must not change router task priority.
- Verification mode: see `config.yml`. Every Task still records implementation evidence and configured quality-gate results; an optional Task Check may be skipped only with an explicit owner waiver, which is not technical `CLEAR`. The Full Objective Check remains mandatory as the V2 epic-level integration gate and includes every owned Task, including waived Tasks; a Goal Check remains mandatory whenever a Goal exists.
- This verification contract's runtime enforcement lives in `internal/data`: `evidence_v2.go` decodes an explicit `check_waiver` Evidence sub-block (Task-only, owner-attributed), and `gate_v2.go`'s `ResolveTaskCompletion` grants completion `AllowedByWaiver` only when no Task Check was ever requested — never once a Check exists, and never for Objective or Goal completion, which stay mandatory and unaffected.

Issues use `open`, `in_progress`, and `resolved`; `stage` is required only while an Issue is `in_progress`. A user-reported defect maps to `type: defect` on an Issue and does not create a separate router state.

Task files may include `complexity_tier` (`low`, `medium`, `high`, or `spike`) and `complexity_reason` as a short planning signal. The pair is validated together, preserved through task status writes, displayed on task cards/details, and required by the create-task planning skill for newly planned tasks.

## 5. Dependencies

- Declared in YAML frontmatter. Task dependencies use `T-###` references and Objective dependencies use `O-###`; the board transition gate and doctor diagnostics resolve them through `internal/data.ResolveDependency`.
- Doctor dependency checks detect duplicate task IDs, missing dependencies, and dependency cycles.
- Cross-Objective dependencies are explicit integration prerequisites and are evaluated by the canonical Objective gate. A Task-check waiver satisfies a Task dependency that requires `clear` — the waiver is the owner's own completion decision, standing in for "clear" there — but never one that requires `accepted`, since there is no Check for the owner to have accepted.

## 6. CLI surface

| Command                | Purpose                                                                           |
| ---------------------- | --------------------------------------------------------------------------------- |
| `savepoint init`       | Scaffold `.savepoint/`, merge the managed agent guide block, print magic prompt to stdout + clipboard |
| `savepoint board`      | Launch TUI; auto-falls-back to plain table on non-TTY                             |
| `savepoint doctor`     | Integrity check + ad-hoc quality-gate run + Layer-2 prompt for AI semantic review |
| `savepoint migrate [dir]` | Preview-first V1-to-V2 conversion; apply checks planned paths in Git before writing |
| `savepoint resume [dir]` | Print the shared V2 Next projection without writing project files |
| `savepoint upgrade-assets [dir] [--dry-run] [--force]` | Refresh package-owned agent skills and the managed agent-guide block without touching project state |
| `--version` / `--help` | Standard global flags                                                             |

- Bare `savepoint` prints help.
- Source modules: see AGENTS.md Codebase Map.
- **Explicitly rejected:** `task new`, `epic new`, `release new`, `plan`, `next`, `status`, `task done`. All are file edits or TUI actions.

`migrate --apply` requires a Git work tree. Every planned write or removal
path must be free of modified, untracked, or ignored files; apply then writes
the converted files directly and sets `schema_version: 2` last. The conversion
manifest records source hashes and mappings, and an error after partial writes
reports the touched paths with Git commands to undo them.

**Names:** npm package `savepoint`; binary `savepoint`. No `vk` alias.
## 7. Independent Check workflow

```
0. Quality Gates       — Executor runs configured build/test gates before handoff.
1. Optional Task Check — A fresh checker runs Quick evidence only when the owner requests it; an explicit owner waiver may skip it.
2. Full Objective Check — A fresh checker must verify every owned Task, integration, and Design reconciliation before Objective closure; every material Issue linked to the current Check must be resolved, including explicit owner acceptance recorded as an Issue resolution.
3. Goal Check          — When a Goal exists, a fresh checker must verify cross-Objective integration before owner acceptance (`scope.kind: release` in stored Check records).
4. Repair              — `NEEDS WORK` records Issues; a Task Check returns the executor to `stage: build` within that Task, while an Objective/Goal Check routes repair to new or newly selected work linked to the Objective without retreating a completed Task; a fresh re-check supersedes the prior Check.
5. Clear               — `CLEAR` is evidence, not automatic ownership; a Task waiver is not `CLEAR`, and only the user closes a Task or accepts an Objective/Goal outcome.
```

- Check records are immutable at `.savepoint/checks/C-###-slug.md`; a re-check writes a new record naming the one it supersedes.
- `savepoint-check` alone writes Check records and closes proven Issues as `verified`; the owner may close an Issue as `accepted` after visual inspection with reason, actor, and time. The executor may record only the owner's explicit decision and does not grant technical clearance.
- Requested Quick Task Checks and mandatory Full Objective/Goal Checks apply the shared non-triggerable method in `agent-skills/references/check-method.md`.

Three layers:

- **Layer 1 (mechanical):** user's chosen linter. Recommended: eslint+dependency-cruiser (TS), radon+pylint (Python), gocyclo+staticcheck (Go). Cross-language fallback: `lizard`. Quality gate config: see `.savepoint/config.yml`.
- **Layer 2 (semantic review):** supplied by the independent Check method. It records evidence and materiality in immutable Check/Issue records; style observations are advisory, not blocking.
- **Layer 3:** `savepoint doctor` runs Layer 1 + prints Layer 2 prompt for ad-hoc use.

## 8. TUI

**Theming:** Atari-Noir is the default theme. **For full design tokens, palette, and rendering rules, see `.savepoint/visual-identity.md`** (loaded conditionally for TUI tasks). Live values in `config.yml` `theme:` section.

Acknowledged terminal limits: fonts, scanlines, glows, letter-spacing, mouse-driven motion don't translate. Lean on color discipline + box-drawing geometry + uppercase headings.

**Render fallbacks:** 256-color → 16-color hard-coded → `NO_COLOR=1` monochrome with glyphs → non-TTY plain table.

**Layout:** the V2 board uses an Objective sidebar, three Task columns (`planned`, `in_progress`, `done`), optional Goal selection, focused detail overlays, static Atari-Noir surfaces, and a deterministic non-TTY plain table. The selected Objective and Task are filtered from the identity-keyed index. The Next area is a one-line glance at the shared `data.Next` projection's own Task or Objective — its lifecycle word (Build/Test/Check, Planned, or Done; never "Audit"), its identity, and its title — and nothing more; the TUI panel and the deterministic non-TTY plain table render that same one line, so piping the board and looking at it still agree with each other. The rung label, per-criterion evidence, owner-wait/exception/dependency wording, and the action sentence remain `savepoint resume`'s narrative (`internal/resume`) and the record's own detail overlay, both unchanged, both still read from the same resolvers — board and resume are intentionally no longer required to render identical wording for that fuller detail; the board is a glance, resume is the report.

**Task-card review outcome (O-012):** each non-planned Task card shows at most one review-outcome badge — `[ ] CHECK`, `[✓] CHECK`, `[!] NEEDS WORK`, `[!] REVIEW` (stale/unknown clearance and a checker-authority gate failure fold into this one label), `[✓] WAIVED`, or `[✓] OWNER ACCEPTED` — at fixed precedence (an owner-accepted exception first, an owner Check waiver second, resolved clearance otherwise), so a card never states two competing outcomes. Completion is the Done column's own fact now: the retired `✓ DONE`, `⚠ DONE`, `BY WAIVER`, and `BY EXCEPTION` badges no longer appear. An open Task's card states the outcome only when it is itself actionable (`NEEDS WORK` or the `REVIEW` fold) or the Task carries an owner's own waiver/exception; `ClearanceCurrent` and `ClearanceMissing` read as Done-column vocabulary instead. Waiver and owner-accepted risk keep the same green accent a current Check gets but stay distinct labels — neither is independent technical clearance, and neither replaces the mandatory Objective or Goal Check; `internal/data` still owns every underlying clearance, waiver, and exception decision (Section 1). This is presentation only.

**Visual guardrail:** the terminal board intentionally uses one black background for Background, Surface, and Surface 2. Do not restore subtly different dark panel fills; depth should come from spacing, dividers, glyphs, and focused Atari Orange borders.

**Terminal color policy:** the board must use a deterministic Lipgloss color profile and one canonical terminal black across truecolor, ANSI256, and ANSI fallbacks. In 256-color mode, Background, Surface, and Surface 2 must map to the same actual black value, not adjacent dark-gray values. Full-screen/root surfaces may paint that one black background for consistency; nested task cards, task text, glyphs, tags, metadata, and router-priority labels should remain foreground-only unless a component explicitly owns a filled visual region. This prevents padded text from creating gray bars in terminals such as PowerShell, Windows Terminal, VS Code terminal, and Warp.

**Border policy:** focus must not change geometry or introduce terminal-specific broken border rendering. Use one consistent box-border family across columns, cards, and overlays. If rounded borders render as dash bars or broken segments in Warp, prefer the single-line border style already allowed by `.savepoint/visual-identity.md`; do not mix rounded and single-line borders as an ad-hoc per-component workaround.

**Board persistence and refresh:** Task status writes use canonical V2 frontmatter with mtime conflict checks. Startup and every reload check the project schema with `data.CheckRuntimeSchema`, load the identity-keyed records with `data.LoadV2Index`, decode router state, and resolve `data.ResolveNext`; rendering performs no filesystem work. The board watches the V2 Objective/Task/Release-compatible Goal paths, preserves selection across reloads, and never infers completion or Goal readiness from sidebar position.

**Implementation modules:** see AGENTS.md Codebase Map.

**Keybindings:** arrow/vim navigation, enter opens focused task detail, space advances, backspace retreats, `p` marks the focused non-done task as router priority, `g` opens Goal selection (`r` remains an undisplayed compatibility alias), `R` refreshes where supported, `d` opens defects, `D` opens Goal details, `A` opens the Audit Register, `?` opens help, and `q` quits or closes overlays.

## 9. Concurrency

- **mtime-based optimistic concurrency.** TUI status writes compare the expected task-file mtime before parsing and again immediately before a no-op or write; conflicts are reported as non-destructive messages that require manual refresh before retry.
- Agents edit freely; the TUI defers.
- **No lockfile.**

## 10. Goal records and PRD history

- Goals use stable `R-###` identities with a title, outcome, success
  conditions, optional scoped Check evidence, material Issue links, and exact
  owner acceptance. Storage remains under `.savepoint/releases/` as
  `Release.md`; membership is derived from the serialized `Objective.release`
  field, and Goal Checks retain `scope.kind: release`.
- Historical `v1`/`v2` PRDs remain in `.savepoint/archive/v1/` and are not
  active Goal records. Goal completion means an accepted cross-Objective
  outcome, never publication, deployment, tagging, or changelog generation.

## 11. Failure modes

All failure modes are diagnosed by `savepoint doctor`. Doctor diagnoses and proposes; never auto-destructive.

| Failure                                      | Behavior                                                    |
| -------------------------------------------- | ----------------------------------------------------------- |
| Corrupt YAML                                 | Doctor flags file:line. TUI marks `⚠ corrupt`, refuses ops. |
| Missing dep                                  | Doctor flags. TUI shows `⚠ broken dep`.                     |
| Dependency cycle                             | Doctor refuses to start either side; prints cycle path.     |
| Duplicate task ID                            | Doctor flags.                                               |
| Check or Issue reference invalid              | Doctor names the missing target and refuses unsafe actions. |
| Task without an Objective                     | Doctor reports the missing owner and refuses V2 loading.    |
| Missing `config.yml`                         | All commands except `init` refuse.                          |
| Unknown CLI flag                             | Show help, exit 1.                                          |

## 12. Distribution & build

> Distribution note: the live repository is a Go module (`github.com/opencode/savepoint`); package and archive checks use the Go build toolchain.

- **License:** MIT.
- **Runtime:** Go CLI binary. Source builds with `go build`; tests run with `go test ./...`.
- **Local build:** `make build` delegates to `internal/buildtool`, builds `savepoint` or `savepoint.exe`, and injects `main.version` from `VERSION` or the latest git tag.
- **Cross-platform builds:** `make build-all` cross-compiles linux-amd64, linux-arm64, darwin-amd64, darwin-arm64, windows-amd64, and windows-arm64 raw binaries into `dist/{platform}-{arch}/savepoint` or `savepoint.exe` for Windows. `make ci` runs the repo-local verification sequence used by CI.
- **Artifacts:** `make dist` creates versioned `.tar.gz` archives in `dist/` for Linux, Darwin, and Windows targets using Go archive APIs, not shell `tar`, and writes SHA256 hashes to `dist/checksums.txt`.
- **Smoke validation:** `make smoke-test` builds the local binary and runs `--version` as a headless exit-0 check.
- **No telemetry.** Ever.

## 13. Testing

| Layer | Tool | Evidence |
| --- | --- | --- |
| Unit and package behavior | `go test ./...` | parser, lifecycle, gate, rendering, and filesystem branches |
| Repository build | `make build` | Go binary and embedded template wiring |
| Distribution | `make build-all`, `make dist`, `make package-check` | six declared platform/architecture targets and checksums |
| Temporary-project integration | migration, board, doctor, resume tests | schema refusal, V2 loading, identity mapping, and non-TTY parity |

Quality gates are named in `.savepoint/config.yml`; coverage percentage alone
never substitutes for acceptance evidence.

## 14. Package versioning

The Go binary reports the injected build version through `--version`; the npm
wrapper distributes the six platform archives and their checksums. Versioning,
publishing, and deployment remain outside Goal completion.
