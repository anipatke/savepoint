---
type: project-design
status: active
last_audited: E02-data-readers (2026-05-01)
---

# Savepoint — System Architecture

> Project-level architecture. Audit-kept fresh: every epic's audit step merges its delta into this document.
>
> **Visual identity** lives separately in `.savepoint/visual-identity.md` and is loaded only for TUI/theme/visual tasks.

## 1. Architecture model

- **File-only.** No MCP server in v1. Agents read and edit Markdown + YAML files directly using their native file tools.
- **Agent routing:** AGENTS.md → `.savepoint/router.md` → template prompts. See AGENTS.md Workflow section.
- **Bundled Agent Skills:** Savepoint ships with custom skills (`draft-prd`, `system-design`, `create-plan`, `create-task`, `build-task`, `audit`) to enforce each phase of the state machine.
- **Token-efficiency principle.**
  - Cold session bootstrap: ~5–7K tokens (one-time per conversation).
  - Per-task incremental: <2KB.
  - Audit: 5–15KB.
  - Anything that breaks these bounds violates the wedge.
- **Go data-reader boundary:** established in epic `E02-data-readers` (2026-05-01). `internal/data` owns Savepoint file parsing and discovery for the Go implementation: task frontmatter models, markdown YAML extraction, router state parsing, config theme defaults, release/epic/task directory listing, and boundary error sentinels.
- **Template assets** live under `templates/` with helpers in `src/templates/` (epic E04).
- **Init command** (`savepoint init`) validates, scaffolds, prints prompt, clipboard, optional install (epic E05).
- **Board command** (`savepoint board`) reads project, non-TTY fallback, Ink TUI, transition gates, mtime writes, audit signaling (epic E06).
- **Audit pipeline** (`savepoint audit`) resolves epic, skips, quality gates, snapshots, router transition, proposal review (epic E07).

## 2. Directory layout

```
<project-root>/
├── AGENTS.md                       ← agent entry point
└── .savepoint/
    ├── PRD.md                      ← project vision (rare changes)
    ├── Design.md                   ← project architecture (this file)
    ├── visual-identity.md          ← design system; loaded conditionally for TUI work
    ├── router.md                   ← state-machine routing
    ├── config.yml                  ← theme, quality_gates, verify_strict
    ├── audit/
    │   └── {E##-epic}/
    │       ├── snapshot.md
    │       └── proposals/
    │           └── proposals.md
    └── releases/
        └── v1/
            ├── PRD.md              ← release-scoped PRD
            └── epics/
                └── E##-{epic-name}/
                    ├── Design.md   ← epic delta
                    └── tasks/
                        └── T001-slug.md
```

AGENTS.md at root (uppercase, cross-vendor spec). Design.md in `.savepoint/` (working doc, not public-facing). visual-identity.md conditional — only loaded for TUI/theme/visual tasks. Subtasks are inline checklists inside task `.md` — never separate files. Epic folders and task files use `E##`/`T##` prefix. Scaffold assets live under `templates/`; generated projects receive rendered copies, not hardcoded strings.

## 3. Hierarchy semantics

| Level        | Definition                                                                             |
| ------------ | -------------------------------------------------------------------------------------- |
| **Release**  | The thing being built. One PRD per release. v1 = MVP.                                  |
| **Epic**     | A major feature within a release. Has its own Design.md (delta from project Design).   |
| **Task**     | Independently buildable. Objective-led. **Requires implementation plan before build.** |
| **Sub-task** | Inline checklist item — _evidence of the implementation plan_, not standalone work.    |

## 4. Status model & gates

Three statuses, with explicit gates:

| Status        | Meaning                    | Entry gate                                                      |
| ------------- | -------------------------- | --------------------------------------------------------------- |
| `planned`     | Ready to build             | plan section non-empty                                          |
| `in_progress` | AI building                | all `depends_on` are `done`                                     |
| `done`        | Complete for current scope | all implementation items checked; verification per project mode |

- `blocked` is a **flag**, not a status — `in_progress` + `blocked: "reason"` is valid.
- `done -> in_progress` is allowed so completed work can be reopened when follow-up work is required.
- Verification mode: see `config.yml`.

## 5. Dependencies

- Declared in YAML frontmatter: `depends_on: [E##-epic/T###-task-id, ...]` (repo-relative IDs).
- `src/validation/dependencies.ts` detects duplicate task IDs, missing dependencies, and dependency cycles.
- Cross-epic deps allowed but warned (signal that epic boundaries may be wrong).

## 6. CLI surface (4 commands, no extras)

| Command                | Purpose                                                                           |
| ---------------------- | --------------------------------------------------------------------------------- |
| `savepoint init`       | Scaffold `.savepoint/`, print magic prompt to stdout + clipboard                  |
| `savepoint board`      | Launch TUI; auto-falls-back to plain table on non-TTY                             |
| `savepoint audit`      | Run audit pipeline (`--skip --reason`, `--epic`)                                  |
| `savepoint doctor`     | Integrity check + ad-hoc quality-gate run + Layer-2 prompt for AI semantic review |
| `--version` / `--help` | Standard global flags                                                             |

- Bare `savepoint` prints help.
- Source modules: see AGENTS.md Codebase Map.
- **Explicitly rejected:** `task new`, `epic new`, `release new`, `plan`, `next`, `status`, `task done`. All are file edits or TUI actions.

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

- `audit_pending` is a **hard gate**: next epic's tasks cannot enter `in_progress` until prior epic is `audited`.
- **High-divergence guard:** if a proposal changes >50% of the live file, TUI requires extra confirmation (threshold tunable in `config.yml`).
- **Skip allowed** via `savepoint audit --skip --reason "..."`. Logged to `.savepoint/audit-log.md`. Permanent `⚠ skipped` badge in TUI.
- **Proposal bundles** use delta-shaped edits: `Insert After`, `Replace`, or `Delete` blocks anchored to exact text.
- **Quality review** is a section inside the proposal bundle.
- **Snapshot availability is an audit precondition.** The router should enter `audit-pending` only after `.savepoint/audit/{E##-epic}/snapshot.md` exists.

Three layers:

- **Layer 1 (mechanical):** user's chosen linter. Recommended: eslint+dependency-cruiser (TS), radon+pylint (Python), gocyclo+staticcheck (Go). Cross-language fallback: `lizard`. Quality gate config: see `.savepoint/config.yml`.
- **Layer 2 (AI semantic review):** baked into the audit reconcile prompt. Outputs a quality-review section in the proposal bundle. **Advisory, not blocking.**
- **Layer 3:** `savepoint doctor` runs Layer 1 + prints Layer 2 prompt for ad-hoc use.

## 8. TUI

**Theming:** Atari-Noir is the default theme. **For full design tokens, palette, and rendering rules, see `.savepoint/visual-identity.md`** (loaded conditionally for TUI tasks). Live values in `config.yml` `theme:` section.

Acknowledged terminal limits: fonts, scanlines, glows, letter-spacing, mouse-driven motion don't translate. Lean on color discipline + box-drawing geometry + uppercase headings.

**Render fallbacks:** 256-color → 16-color hard-coded → `NO_COLOR=1` monochrome with glyphs → non-TTY plain table.

**Layout:** single screen with a 5-column Kanban board and detail pane. Non-TTY output uses `src/tui/render/plain-table.ts`.

**Implementation modules:** see AGENTS.md Codebase Map (E06 and E07 epic rows).

**Keybindings:** arrow/vim navigation, enter advances, backspace retreats, r/R refreshes, a/A exits toward audit review when proposals exist, q quits.

## 9. Concurrency

- **mtime-based optimistic concurrency.** TUI status writes compare the expected task-file mtime before parsing and again immediately before a no-op or write; conflicts are reported as non-destructive messages that require manual refresh before retry.
- Agents edit freely; the TUI defers.
- **No lockfile.**

## 10. Release versioning (PRDs)

- Sequential integer (`v1`, `v2`). Optional `name` in YAML.
- `savepoint doctor` warns when creating `v2` while `v1` has un-audited epics.

## 11. Failure modes

All failure modes are diagnosed by `savepoint doctor`. Doctor diagnoses and proposes; never auto-destructive.

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

> Audit note: the live repository is transitioning from the documented TypeScript/Node implementation to a Go module (`github.com/opencode/savepoint`). The architecture document still contains substantial TypeScript-era implementation detail and should be reconciled as Go epics are audited.

- **License:** MIT.
- **Install:** primary `npx savepoint init`, persistent `npm i -g savepoint` → `savepoint`.
- **Runtime:** Node 20.10+ LTS, ESM-only, no native deps. macOS / Linux / Windows-Terminal.
- **Repo:** single package. TypeScript strict. `tsup` build → `dist/`. Bin `dist/cli.js` shebanged.
- **No telemetry.** Ever.

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

## 14. Package versioning

- `0.1.0` — first public release: scaffolding, status model, CLI, basic TUI, audit (no AI semantic review).
- `0.2.0` — AI semantic review + broader quality-gate language presets.
- `0.3.0` — file watching, search.
- `1.0.0` — MCP server + production stability.

Strict semver. Pre-1.0 minors may break.
