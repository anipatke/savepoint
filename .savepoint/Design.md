---
type: project-design
status: active
last_audited: E04-templates-and-prompts (2026-04-27)
---

# Savepoint — System Architecture

> Project-level architecture. Audit-kept fresh: every epic's audit step merges its delta into this document.
>
> **Visual identity** lives separately in `.savepoint/visual-identity.md` and is loaded only for TUI/theme/visual tasks.

## 1. Architecture model

- **File-only.** No MCP server in v1. Agents read and edit Markdown + YAML files directly using their native file tools.
- **Agent-agnostic via the Router Pattern:**
  1. `AGENTS.md` at the project root (spec-compliant; readable by all major coding agents).
  2. `AGENTS.md` routes the agent into `.savepoint/router.md` (the state-machine file).
  3. `.savepoint/router.md` conditionally points the agent at the next template prompt based on project state.
  4. Templates contain **embedded HTML-comment instructions** (`<!-- AGENT: ... -->`) the agent follows verbatim.
- **Token-efficiency principle.**
  - Cold session bootstrap: ~5–7K tokens (one-time per conversation).
  - Per-task incremental: <2KB.
  - Audit: 5–15KB.
  - Anything that breaks these bounds violates the wedge.
- **Template asset boundary:** established in epic `E04-templates-and-prompts` (2026-04-27). Project scaffolding and workflow prompts now live as versioned markdown/YAML assets under `templates/`, with small TypeScript helpers in `src/templates/` for manifest lookup, path resolution, loading, and interpolation.

## 2. Directory layout

```
<project-root>/
├── AGENTS.md                       ← uppercase, spec-standard, agent entry point
└── .savepoint/
    ├── PRD.md                      ← project vision (rare changes)
    ├── Design.md                   ← project architecture (this file)
    ├── visual-identity.md          ← design system; loaded conditionally for TUI work
    ├── router.md                   ← state-machine routing
    ├── config.yml                  ← theme, quality_gates, verify_strict
    ├── audit/
    │   └── {E##-epic}/
    │       ├── snapshot.md
    │       └── proposals.md
    └── releases/
        └── v1/
            ├── PRD.md              ← release-scoped PRD
            └── epics/
                └── E##-{epic-name}/
                    ├── Design.md   ← epic delta
                    └── tasks/
                        └── T001-slug.md
```

- `AGENTS.md` is at root (uppercase, cross-vendor spec).
- `Design.md` lives in `.savepoint/` (working doc, not public-facing).
- `visual-identity.md` is conditional — only loaded by router for TUI/theme/visual tasks.
- **Subtasks are inline checklists** inside the task `.md` — never separate files.
- Epic folders are prefixed for navigation (`E01-scaffolding/`). Task files and IDs are prefixed the same way (`T001-package-baseline.md`).
- The Savepoint package itself owns scaffold assets under `templates/project/`, `templates/release/v1/`, and `templates/prompts/`; generated projects receive rendered copies, not hardcoded strings from command handlers.

## 3. Hierarchy semantics

| Level        | Definition                                                                             |
| ------------ | -------------------------------------------------------------------------------------- |
| **Release**  | The thing being built. One PRD per release. v1 = MVP.                                  |
| **Epic**     | A major feature within a release. Has its own Design.md (delta from project Design).   |
| **Task**     | Independently buildable. Objective-led. **Requires implementation plan before build.** |
| **Sub-task** | Inline checklist item — _evidence of the implementation plan_, not standalone work.    |

## 4. Status model & gates

Five statuses, with explicit gates:

| Status        | Meaning                           | Entry gate                                    |
| ------------- | --------------------------------- | --------------------------------------------- |
| `backlog`     | Task exists, no plan              | created                                       |
| `planned`     | Implementation plan written       | plan section non-empty                        |
| `in_progress` | AI building                       | all `depends_on` are `done`                   |
| `review`      | Build done, awaiting verification | all sub-tasks checked off                     |
| `done`        | Verified, locked                  | assertion + (configurable) quality gates pass |

- `blocked` is a **flag**, not a status — `in_progress` + `blocked: "reason"` is valid.
- Status names and transition validation live in `src/domain/status.ts`.
- `done -> review` is allowed so completed work can be reopened for audit-stale rechecks. The warning itself belongs to future CLI/TUI write layers.
- **Verification mode is configurable per project** (`verify_strict: true|false`). Default: `false` (vibe-coder soft mode).

## 5. Dependencies

- Declared in YAML frontmatter: `depends_on: [E##-epic/T###-task-id, ...]` (repo-relative IDs).
- Task dependency IDs are parsed with the same domain parser as task IDs (`src/domain/ids.ts`).
- `src/validation/dependencies.ts` detects duplicate task IDs, missing dependencies, and dependency cycles.
- `src/readers/tasks.ts` reads an epic task directory, validates all task markdown, then runs graph validation only after readable task files are parsed.
- Cross-epic deps allowed but warned (signal that epic boundaries may be wrong).
- TUI shows blocked tasks as visually locked.
- `savepoint doctor` detects cycles.

## 6. CLI surface (4 commands, no extras)

| Command                | Purpose                                                                           |
| ---------------------- | --------------------------------------------------------------------------------- |
| `savepoint init`       | Scaffold `.savepoint/`, print magic prompt to stdout + clipboard                  |
| `savepoint board`      | Launch TUI; auto-falls-back to plain table on non-TTY                             |
| `savepoint audit`      | Run audit pipeline (`--skip --reason`, `--epic`)                                  |
| `savepoint doctor`     | Integrity check + ad-hoc quality-gate run + Layer-2 prompt for AI semantic review |
| `--version` / `--help` | Standard global flags                                                             |

- Bare `savepoint` prints help.
- As of E03, the binary has a real command shell around stub handlers:
  - `src/cli.ts` owns process globals and delegates to `runCli()`.
  - `src/cli/args.ts` parses the fixed command and global flag surface.
  - `src/cli/run.ts` dispatches help, version, command stubs, unknown commands, and unknown flags.
  - `src/cli/help.ts` generates deterministic top-level and command-level help.
  - `src/cli/environment.ts` detects TTY, color, and platform capability from injected inputs.
  - `src/commands/*.ts` contains deterministic not-yet-implemented stubs.
- Implemented global flags: `--help`, `-h`, `--version`, `-v`.
- Implemented command help flags: `--help`, `-h` after `init`, `board`, `audit`, or `doctor`.
- Command behavior remains future work. Until later epics fill it in, command stubs return `EXIT_NOT_IMPLEMENTED`.
- **Explicitly rejected:** `task new`, `epic new`, `release new`, `plan`, `next`, `status`, `task done`. All are file edits or TUI actions.
- **Agents must not run `savepoint` commands.** Stated in AGENTS.md.

**Names:** npm package `savepoint`; binary `savepoint`. No `vk` alias.

## 7. Audit pipeline (6 steps)

```
0. Quality Gates  — CLI runs configured commands. Halts on failure if block_on_failure: true.
1. Snapshot       — CLI writes file tree (gitignore-respecting) + changed-files list. NO code contents.
2. Diff Brief     — State flips to audit_pending: {E##-epic}. Magic prompt printed to user.
3. Reconcile      — Agent reads epic Design + snapshot + scoped code. Writes one proposal bundle to
                    .savepoint/audit/{E##-epic}/proposals.md.
4. Review         — TUI shows side-by-side per-proposal diff. Approve / reject / edit each.
5. Commit         — Approved proposals overwrite live files. Epic gets status: audited. Next epic unlocks.
```

- **First epic = E01-scaffolding by convention.** First audit establishes the baseline. ✓ _E01-scaffolding audited 2026-04-27._
- `audit_pending` is a **hard gate**: next epic's tasks cannot enter `in_progress` until prior epic is `audited`.
- **High-divergence guard:** if a proposal changes >50% of the live file, TUI requires extra confirmation (threshold tunable in `config.yml`).
- **Skip allowed** via `savepoint audit --skip --reason "..."`. Logged to `.savepoint/audit-log.md`. Permanent `⚠ skipped` badge in TUI.
- **Proposal bundles** use delta-shaped edits by default: `Insert After`, `Replace`, or `Delete` blocks anchored to exact text. Agents avoid rewriting whole sections unless a section is genuinely being replaced.
- **Quality review** is a section inside the proposal bundle. A standalone `quality-review.md` is only needed when the review UI explicitly asks for separate files.
- **Snapshot availability is an audit precondition.** The router should enter `audit-pending` only after `.savepoint/audit/{E##-epic}/snapshot.md` exists. While the audit CLI is still a stub, the closeout task should create one manual snapshot instead of making the next agent probe for a missing file at audit startup.
- **Codebase Map** is generated mechanically in `AGENTS.md` between markers from changed-module metadata. Agents provide or refine short purpose text for new modules; they do not hand-rewrite the whole map unless the generator is unavailable.

### Quality gates (audit step 0)

```yaml
# .savepoint/config.yml
quality_gates:
  lint: "<command>" # null to disable
  typecheck: "<command>"
  test: "<command>"
  block_on_failure: true
```

Three layers:

- **Layer 1 (mechanical):** user's chosen linter (we recommend, don't ship). TS: `eslint` + `dependency-cruiser`. Python: `radon` + `pylint`. Go: `gocyclo` + `staticcheck`. **Cross-language fallback:** `lizard`.
- **Layer 2 (AI semantic review):** baked into the audit reconcile prompt. Outputs a quality-review section in the proposal bundle. **Advisory, not blocking.**
- **Layer 3:** `savepoint doctor` runs Layer 1 + prints Layer 2 prompt for ad-hoc use.

## 8. TUI

**Theming:** Atari-Noir is the default theme. **For full design tokens, palette, and rendering rules, see `.savepoint/visual-identity.md`** (loaded conditionally for TUI tasks). Live values in `config.yml` `theme:` section.

Acknowledged terminal limits: fonts, scanlines, glows, letter-spacing, mouse-driven motion don't translate. Lean on color discipline + box-drawing geometry + uppercase headings.

**Render fallbacks:** 256-color → 16-color hard-coded → `NO_COLOR=1` monochrome with glyphs → non-TTY plain table.

**Layout:** single screen — header (release/epic context) + 5-column Kanban + detail pane.

**Keybindings:** arrow + vim. `space`/`enter`/`e`/`r`/`R`/`E`/`A`/`q`. **No file watching MVP** (manual `r` refresh).

**Out of MVP:** file watching, drag-and-drop, multi-select, search, inline editing.

## 9. Concurrency

- **mtime-based optimistic concurrency.** TUI reads mtime, re-stats before write, prompts "Reload? [Y/n]" on conflict.
- Agents edit freely; the TUI defers.
- **No lockfile.**

## 10. Release versioning (PRDs)

- Sequential integer (`v1`, `v2`). Optional `name` in YAML.
- `savepoint doctor` warns when creating `v2` while `v1` has un-audited epics.

## 11. Failure modes

All routed through `savepoint doctor`. Doctor diagnoses and proposes; never auto-destructive.

| Failure                                      | Behavior                                                    |
| -------------------------------------------- | ----------------------------------------------------------- |
| Corrupt YAML                                 | Doctor flags file:line. TUI marks `⚠ corrupt`, refuses ops. |
| Missing dep                                  | Doctor flags. TUI shows `⚠ broken dep`.                     |
| Dependency cycle                             | Doctor refuses to start either side; prints cycle path.     |
| Duplicate task ID                            | Doctor flags.                                               |
| Audit proposals without `audit_pending` flag | Doctor offers cleanup or restore.                           |
| Task in nonexistent epic                     | Doctor moves to `.savepoint/orphans/`.                      |
| Missing `config.yml`                         | All commands except `init` refuse.                          |
| Unknown CLI flag                             | Show help, exit 1.                                          |

## 12. Distribution & build

- **License:** MIT.
- **Install:** primary `npx savepoint init`, persistent `npm i -g savepoint` → `savepoint`.
- **Runtime:** Node 20.10+ LTS, ESM-only, no native deps. macOS / Linux / Windows-Terminal.
- **Repo:** single package. TypeScript strict. `tsup` build → `dist/`. Bin `dist/cli.js` shebanged.
- **No telemetry.** Ever.
- **Baseline scaffold:** established in epic `E01-scaffolding` (2026-04-27). Package name `savepoint`, version `0.1.0`. Build, typecheck, lint, format, and test gates all pass.
- **Read-only data model:** established in epic `E02-data-model` (2026-04-27). `js-yaml` is the only runtime dependency added for structured YAML/frontmatter parsing.
- **CLI foundation:** established in epic `E03-cli-foundation` (2026-04-27). The binary now has a testable parser, help/version handling, command dispatch, terminal capability detection, and deterministic stubs for `init`, `board`, `audit`, and `doctor`.
- **Template and prompt assets:** established in epic `E04-templates-and-prompts` (2026-04-27). Default project, release, config, visual identity, and workflow prompt content now lives under `templates/`, with `src/templates/` exposing typed manifest, path, load, and render helpers.

## 13. Testing

| Layer                                                    | Tool                             | Coverage                                                               |
| -------------------------------------------------------- | -------------------------------- | ---------------------------------------------------------------------- |
| Unit: file ops, YAML, frontmatter, snapshot gen          | `vitest`                         | High                                                                   |
| Unit: state transitions, dep resolution, cycle detection | `vitest`                         | High                                                                   |
| Integration: CLI commands in temp dirs                   | `vitest` + `tmp`                 | Medium                                                                 |
| TUI reducers (state, isolated from rendering)            | `vitest` + `ink-testing-library` | Medium                                                                 |
| TUI rendering (snapshot tests)                           | —                                | **None.** Brittle.                                                     |
| End-to-end with real AI agents                           | Manual matrix                    | Pre-release: `[Claude, Cursor, Gemini, Aider]` × `[init, plan, audit]` |

~70% line coverage target; behavior coverage prioritized.

E02 adds focused unit coverage for ID parsing, status validation, task/release/epic/router/config frontmatter validation, markdown read-boundary failures, project root discovery, dependency graph errors, and epic task-set reading. As of the E02 audit, `npm test` reports 16 passing test files and 176 passing tests when run outside the sandbox.

E03 adds focused CLI coverage for bare invocation, global help/version flags, command-level help, command dispatch, unknown top-level commands and flags, unknown command-level flags, environment detection, and command stubs. As of E03 closeout, the focused CLI runner suite reports 35 passing tests.

E04 adds focused template coverage for required project/release/prompt assets, router states, embedded agent markers, audit prompt shape, template path resolution, missing-template boundary errors, and placeholder interpolation. As of the E04 audit snapshot, `npm test` reports 26 passing test files and 299 passing tests.

## 14. Package versioning

- `0.1.0` — first public release: scaffolding, status model, CLI, basic TUI, audit (no AI semantic review).
- `0.2.0` — AI semantic review + broader quality-gate language presets.
- `0.3.0` — file watching, search.
- `1.0.0` — MCP server + production stability.

Strict semver. Pre-1.0 minors may break.
