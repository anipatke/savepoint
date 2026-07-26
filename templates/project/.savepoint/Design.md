---
type: project-design
status: active
last_audited: never
---

# {{PROJECT_NAME}} — System Architecture

> Project-level architecture. Audit-kept fresh: every epic's audit step merges its delta into this document.
>
> **Visual identity** lives separately in `.savepoint/visual-identity.md` and defines the project's design system (palette, typography, patterns).
>
> Fill each section with concrete, falsifiable claims. Generic architecture docs produce generic agents.

## 1. Architecture model

List 4–6 architecture commitments as bullets. Each bullet must be a rule the system follows, not a preference.

> - **File-only.** No MCP server in v1. Agents read and edit Markdown + YAML files directly using their native file tools.
> - **Agent routing:** agent guide → `.savepoint/router.md` → phase skills. See `AGENTS.md` Workflow section.
> - **Bundled agent skills:** ship custom skills (planning, design, task creation, task build, audit) to enforce the state machine and capture release-level defects.
> - **Token-efficiency principle.** Cold session bootstrap: ~5–7K. Per-task incremental: <2KB. Audit: 5–15KB. Anything that breaks these bounds violates the wedge.

## 2. Directory layout

Show the canonical directory tree as a fenced code block, then list 2–3 placement rules below. Depth: project root → `.savepoint/` → `releases/{release}/` → `epics/E##-{slug}/`.

> ```
> <project-root>/
> ├── AGENTS.md                       ← agent entry point
> └── .savepoint/
>     ├── PRD.md                      ← project vision (rare changes)
>     ├── Design.md                   ← project architecture (this file)
>     ├── visual-identity.md          ← design system (palette, type, patterns)
>     ├── Guardrails.md               ← engineering policy, severity model, rule index
>     ├── Health-Check.md             ← Quick/Full/Deep evidence gates
>     ├── router.md                   ← state-machine routing
>     ├── config.yml                  ← theme, quality_gates, verify_strict
>     └── releases/
>         └── {release}/              ← e.g. v1, v1.1
>             ├── {release}-PRD.md    ← release-scoped PRD
>             ├── defects/
>             │   └── D001-slug.md    ← release-level repair record
>             └── epics/
>                 └── E##-{epic-name}/
>                     ├── E##-Detail.md   ← epic delta
>                     ├── E##-Audit.md    ← audit findings + admin apply proposals
>                     └── tasks/
>                         └── T001-slug.md
> ```
>
> AGENTS.md at root (uppercase, cross-vendor spec). Design.md in `.savepoint/` (working doc). Subtasks are inline checklists inside task `.md` — never separate files. Epic folders and task files use `E##`/`T##` prefix.

## 3. Hierarchy semantics

Markdown table with 5 rows: Release, Epic, Task, Defect, Sub-task. Each row: bold term, one-line definition, one-line ownership statement.

> | Level        | Definition                                                                             |
> | ------------ | -------------------------------------------------------------------------------------- |
> | **Release**  | The thing being built. One PRD per release. v1 = MVP.                                  |
> | **Epic**     | A major feature within a release. Has its own E##-Detail.md (delta from project Design). |
> | **Task**     | Independently buildable. Objective-led. **Requires implementation plan before build.** |
> | **Defect**   | Release-level repair artifact for observed bugs or regressions; separate from planned epic/task scope. |
> | **Sub-task** | Inline checklist item — _evidence of the implementation plan_, not standalone work.    |

## 4. Status model & gates

Status table with 4 columns: Status, Meaning, Entry gate, Who may set it. Then a bullet list of ownership rules: `internal/data` as canonical owner, defect-specific lifecycle, complexity_tier fields, blocked flag.

> | Status        | Meaning                    | Entry gate                                                      | Who may set it                         |
> | ------------- | -------------------------- | --------------------------------------------------------------- | -------------------------------------- |
> | `planned`     | Ready to build             | plan section non-empty                                          | User or planning workflow              |
> | `in_progress` | AI building                | all `depends_on` are `done`                                     | Agent when starting implementation     |
> | `done`        | Complete for current scope | all implementation items checked; verification per project mode | User only                              |
>
> - `blocked` is a **flag**, not a status — `in_progress` + `blocked: "reason"` is valid.
> - `internal/data` is the single owner of task lifecycle rules: canonical statuses, canonical stages, parse compatibility for legacy `phase`, write validation, and transition helpers must flow through that package.
> - Agents may only advance a task into `in_progress`; they must not set `done` or retreat a task to an earlier status.
> - Only the user may set a task to `done` or retreat it from `done` to `in_progress`.
> - Defects use defect-specific lifecycle statuses: `open`, `in_progress`, `resolved`. `stage` is required while a defect is `in_progress`.

## 5. Dependencies

Bullet list: 3–4 dependency rules. Each rule must include the mechanism (YAML frontmatter, resolver function, doctor check).

> - Declared in YAML frontmatter. Full task IDs (`E##-epic/T###-task-id`) are preferred; same-epic shorthand may use `T###` or the task filename stem. The board transition gate and doctor diagnostics resolve these forms through a canonical resolver.
> - Doctor dependency checks detect duplicate task IDs, missing dependencies, and dependency cycles.
> - Cross-epic deps allowed but warned (signal that epic boundaries may be wrong).

## 6. CLI surface

Command table: Command, Purpose. Then a bullet list of explicitly rejected commands and the naming convention.

> | Command                | Purpose                                                                           |
> | ---------------------- | --------------------------------------------------------------------------------- |
> | `<tool> init`          | Scaffold `.savepoint/`, merge the managed agent guide block, print magic prompt to stdout + clipboard |
> | `<tool> board`         | Launch TUI; auto-falls-back to plain table on non-TTY                             |
> | `<tool> doctor`        | Integrity check + ad-hoc quality-gate run + Layer-2 prompt for AI semantic review |
> | `<tool> upgrade-assets [dir] [--dry-run] [--force]` | Refresh package-owned agent skills and the managed agent-guide block without touching project state |
> | `--version` / `--help` | Standard global flags                                                             |
>
> Bare `<tool>` prints help. **Explicitly rejected:** `task new`, `epic new`, `release new`, `plan`, `next`, `status`, `task done`. All are file edits or TUI actions.

## 7. Agent audit workflow

Numbered workflow steps 0–5 (Quality Gates through Apply + Close). Then 4–5 follow-up bullets covering the audit-file sections, the user-facing vs admin split, the three layers, and the explicit "no CLI pipeline" rule.

> ```
> 0. Quality Gates  — Build agent runs configured build/test gates before audit handoff.
> 1. Audit Pending  — Router enters `audit-pending` for the completed epic.
> 2. Reconcile      — Fresh audit agent reads router, epic detail, task files, Design.md, AGENTS.md, and scoped source/test files.
> 3. Findings       — Agent writes exactly one `{epic}/E##-Audit.md`.
> 4. Review         — User reviews the TUI Epic Detail Audit tab.
> 5. Apply + Close  — After user approval, agent applies proposal blocks, updates the audit file's visible findings, marks the epic audited, updates `last_audited`, and advances router.
> ```
>
> - `audit-pending` is a **hard gate**: next epic's tasks cannot enter `in_progress` until prior epic is `audited` or the user explicitly skips the audit.
> - `E##-Audit.md` has two user-facing sections: `## Main Findings` and `## Code Style Review`. File-specific `### Target File` / `### Replace` / `### With` blocks belong under a separate `## Proposed Changes` admin section so the TUI Audit tab can omit them.
> - Apply/close must rewrite `## Main Findings` and `## Code Style Review` in the same `E##-Audit.md` so the TUI Audit tab shows resolved findings and remaining risks instead of stale pre-apply blockers.
> - There is no `<tool> audit` CLI pipeline in the active design. Epic audit is performed by agents using `agent-skills/savepoint-audit-epic/SKILL.md`; an explicit read-only review of one in-progress task uses `agent-skills/savepoint-audit-task/SKILL.md`. Both apply the shared method in `agent-skills/references/audit-method.md`.
>
> Three layers:
> - **Layer 1 (mechanical):** user's chosen linter. Quality gate config: see `.savepoint/config.yml`.
> - **Layer 2 (AI semantic review):** baked into the audit reconcile prompt. Outputs Main Findings and Code Style Review in the epic-local audit file. **Advisory, not blocking.**
> - **Layer 3:** `<tool> doctor` runs Layer 1 + prints Layer 2 prompt for ad-hoc use.

## 8. TUI

Theming, render fallbacks, layout, terminal color policy, border policy, board persistence, keybindings. Skip this section if the project has no TUI.

> **Theming:** define a single theme as the default. For full design tokens and rendering rules, see `.savepoint/visual-identity.md`. If your project includes a terminal UI, the TUI adaptation appendix covers terminal-specific translations. Live values live in `config.yml` `theme:` section.
>
> **Render fallbacks:** 256-color → 16-color hard-coded → `NO_COLOR=1` monochrome with glyphs → non-TTY plain table.
>
> **Layout:** single screen with a 3-column task board (`planned`, `in_progress`, `done`), optional epic sidebar on wide terminals, centered overlays, height-aware viewport slicing, full-width dividers, and navigation hints. Active router `next_action` renders as a dedicated full-width line below the header with phase-colored prefix styling.
>
> **Keybindings:** arrow/vim navigation, enter opens focused task detail, space advances, backspace retreats, `p` marks the focused non-done task as router priority, `?` opens help, `q` quits or closes overlays.

## 9. Concurrency

Describe the concurrency model in 3–5 bullets. Cover the optimistic-concurrency mechanism, the agent-vs-TUI roles, and whether a lockfile is used.

> - **mtime-based optimistic concurrency.** Status writes compare the expected file mtime before parsing and again immediately before a write; conflicts are reported as non-destructive messages that require manual refresh before retry.
> - Agents edit freely; the TUI defers.
> - **No lockfile.**

## 10. Release versioning (PRDs)

Naming scheme, doctor behavior on roll-overs, optional fields. 2–3 bullets.

> - Sequential integer (`v1`, `v2`). Optional `name` in YAML.
> - `<tool> doctor` warns when creating `v2` while `v1` has un-audited epics.

## 11. Failure modes

Markdown table with 2 columns: Failure, Behavior. 5–8 rows covering the diagnosable failure modes. All are diagnosed by `<tool> doctor`; doctor diagnoses and proposes, never auto-destructive.

> | Failure                                      | Behavior                                                    |
> | -------------------------------------------- | ----------------------------------------------------------- |
> | Corrupt YAML                                 | Doctor flags file:line. TUI marks `⚠ corrupt`, refuses ops. |
> | Missing dep                                  | Doctor flags. TUI shows `⚠ broken dep`.                     |
> | Dependency cycle                             | Doctor refuses to start either side; prints cycle path.     |
> | Duplicate task ID                            | Doctor flags.                                               |
> | Audit proposals without `audit_pending` flag | Doctor offers cleanup or restore.                           |
> | Task in nonexistent epic                     | Doctor moves to `.savepoint/orphans/`.                      |
> | Missing `config.yml`                         | All commands except `init` refuse.                          |
> | Unknown CLI flag                             | Show help, exit 1.                                          |

## 12. Distribution & build

License, runtime, build commands, distribution targets, artifacts. 4–6 bullets.

> - **License:** MIT.
> - **Local build:** `make build` (or the project's equivalent) produces the primary binary and injects the version from `VERSION` or the latest git tag.
> - **Cross-platform builds:** cross-compile the matrix of OS × arch your users need; raw binaries land in `dist/{platform}-{arch}/`. Repo-local `make ci` runs the verification sequence used by CI.
> - **Artifacts:** versioned archives plus SHA256 checksums; prefer language-native archive APIs over shell `tar`.
> - **Smoke validation:** build the local binary and run `--version` as a headless exit-0 check.
> - **No telemetry.** Ever.

## 13. Testing

Layer table with 3 columns: Layer, Tool, Coverage. 4–6 rows. End with a coverage target.

> | Layer                                                    | Tool                                  | Coverage                              |
> | -------------------------------------------------------- | ------------------------------------- | ------------------------------------- |
> | Unit: file ops, YAML, frontmatter, snapshot gen          | language-native test runner           | High                                  |
> | Unit: state transitions, dep resolution, cycle detection | language-native test runner           | High                                  |
> | Integration: CLI commands in temp dirs                   | language-native test runner + `tmp`   | Medium                                |
> | TUI reducers (state, isolated from rendering)            | TUI message-driven unit tests         | Medium                                |
> | TUI rendering (snapshot tests)                           | —                                     | **None.** Brittle.                    |
> | End-to-end with real AI agents                           | Manual matrix                         | Pre-release: agents × init/plan/audit |
>
> Target ~70% line coverage; behavior coverage prioritized.

## 14. Package versioning

Numbered list of milestones with version and the feature added. Strict semver with a note about pre-1.0 breakability.

> - `0.1.0` — first public release: scaffolding, status model, CLI, basic TUI, audit (no AI semantic review).
> - `0.2.0` — AI semantic review + broader quality-gate language presets.
> - `0.3.0` — file watching, search.
> - `1.0.0` — production stability and the planned integration layer.
>
> Strict semver. Pre-1.0 minors may break.
