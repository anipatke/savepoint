![Savepoint Banner](assets/banner.png)

# Savepoint

> **Hard gates for AI-driven development.**
>
> Local files. Tight context. No telemetry.

Savepoint gives AI coding agents an engineering process they can actually
follow. It turns a fuzzy idea into a durable plan, limits each implementation
step to the files it needs, and asks for independent evidence before the work
advances.

The result is a small, local-first control plane for shipping software with
agents — stored in your repository, inspectable by your team, and usable with
the tools you already have.

[Website](https://www.getsavepoint.dev/) ·
[GitHub](https://github.com/anipatke/savepoint) ·
[npm](https://www.npmjs.com/package/savepoint)

## Why Savepoint?

AI-assisted development tends to fail at the boundaries: context gets too
large, scope gets vague, architecture drifts, and “tests passed” becomes a
substitute for someone checking the actual outcome.

Savepoint makes those boundaries explicit:

- **Plan before implementation.** Intent, architecture, and constraints are
  written down before an agent starts changing code.
- **Keep execution bounded.** Each Task has one observable outcome and a
  strict list of Context Files.
- **Keep policy durable.** Guardrails live beside the project instead of being
  hidden in a prompt or a chat transcript.
- **Check at the right level.** A Task Check is optional and may be explicitly
  waived by the owner. The mandatory Full Objective Check verifies every owned
  Task and its integration; a Full Goal Check is mandatory for every Goal.
- **Keep ownership clear.** Agents implement and prove their work; people
  decide whether the outcome is what they wanted.
- **Every project has a Goal.** The router selects a live Goal, and every
  Objective belongs to exactly one Goal. A mandatory Full Goal Check and owner
  acceptance are required before marking it done. Goals do not publish,
  deploy, tag, or generate changelogs.

Existing V2 projects keep the compatibility storage names: stable `R-###`
identities under `.savepoint/releases/` as `Release.md`, Objective `release:`
references, router `release:` selections, and Check `scope.kind: release`. The
V2 board presents these records as Goals. `savepoint init` creates and selects
R-001, titled after the project. `savepoint migrate` retains the V1 router's
live Goal and, on an unresolved selection, reuses a uniquely identifiable
existing live Goal for the active work when possible. If selected work belongs
only to a historical Goal, migration creates a continuation Goal and moves
that Objective into it. An unresolved release lifecycle decision stays in the
preview. If an existing project is missing its router Goal,
Next says `Choose a Goal`; doctor explains how to select or create one. A
missing Objective `release:` remains loadable, but doctor names the Objective
and exact line to add.

Savepoint does not replace Git, your test runner, or human judgment. It gives
those things a shared workflow.

## The workflow

```text
IDEA  ─────►  DESIGN  ─────►  TASK  ─────►  OBJECTIVE CHECK  ─────►  GOAL CHECK
  intent       architecture    bounded      mandatory Full           mandatory Full
  & outcome    & guardrails    execution     integration              cross-Objective
                                      ╰─ optional Quick Task Check
                                         or explicit owner waiver
```

### Idea

Capture what you are building, who it is for, why it matters, and what is out
of scope. A rough sentence is enough to begin.

### Design

Turn intent into architecture: components, interfaces, boundaries, decisions,
and durable engineering guardrails. Design describes the system that exists;
it is not a second backlog.

### Task

Break the next outcome into a small execution packet. A Task records its
objective, dependencies, acceptance criteria, implementation plan, and the
exact files the agent may need to read.

If the plan is materially wrong, the agent stops with `REPLAN REQUIRED` rather
than quietly inventing a new architecture.

### Check

An independent checker may run an optional Quick Task Check when the owner
requests it. If the owner skips that local review, the Task evidence records an
explicit waiver; the waiver is not technical `CLEAR` and does not waive any
acceptance criterion or guardrail. Before an Objective closes, a mandatory
Full Objective Check verifies every owned Task, cross-Task integration, and
Design reconciliation. A mandatory Full Goal Check for every Goal verifies
cross-Objective integration before the owner accepts that exact Check.

## Quick start

Install Savepoint into the repository you want to work on:

```bash
npx savepoint init
```

Then open the board and inspect the next action:

```bash
npx savepoint board
npx savepoint resume
```

Use the V2 objective filter when you want to focus the board, and use the
command help when you need the complete option contract:

```bash
npx savepoint board --objective O-001
npx savepoint --help
npx savepoint board --help
```

Ask your coding agent to read the generated `AGENTS.md`. That file routes the
agent to the current Savepoint state, the matching workflow skill, and the
bounded context for the active Task.

When you want a deterministic project check:

```bash
npx savepoint doctor
```

Commit the generated `.savepoint/` files and `AGENTS.md` with your project.
They are the project memory that lets a new agent, a new session, or a human
teammate pick up where the last one stopped.

## The command line

| Command | What it does |
| --- | --- |
| `savepoint init [dir]` | Scaffolds Savepoint's project files and agent guidance. |
| `savepoint create-task --objective O-### --draft <path> [dir]` | Creates a Task from an ID-free draft, assigns its project-wide ID, and validates the full V2 index. |
| `savepoint board [--objective O-###]` | Opens the keyboard-driven V2 board, optionally focused on one Objective. |
| `savepoint resume [dir]` | Prints the current state and the next recorded action. |
| `savepoint doctor` | Runs deterministic project diagnostics and configured quality gates. |
| `savepoint migrate [dir]` | Converts a legacy Savepoint project to V2. Preview is the default. |
| `savepoint upgrade-assets [dir]` | Refreshes shipped skills and templates in an existing project. |

Run `savepoint --help` for the command list or `savepoint <command> --help` for
command-specific options. The package is also available through `npx` for
projects that do not need a global installation.

Planners create new Tasks with `savepoint create-task`; they provide an
Objective and a complete Task draft without an `id`, and the command assigns
the next project-wide ID and path. It strict-loads the full V2 index before it
reports success. Agents may use this command only for new Tasks. After
creating or renaming another identity-bearing record, use `savepoint resume`
to require a strict load of the full index.

## The terminal board

`savepoint board` is a fast, keyboard-driven view of the work recorded in your
project:

- **Next** shows the single action selected by the project state.
- **Objectives** group related outcomes without forcing a large hierarchy on
  small projects.
- **Task columns** show planned, in-progress, and done work, including the
  current build/test/audit stage.
- **Detail views** expose acceptance criteria, dependencies, evidence, and
  issues without leaving the terminal.
- **Goal context** is required: the board shows the router-selected Goal and
  only its Objectives and Tasks. Press `g` to switch Goals (`r` remains an
  undisplayed compatibility alias); membership comes from each Objective's
  `release: R-###` field. With no valid Goal selected, the board shows
  `Choose a Goal` and no project-wide work.
- **Router priority** lets you focus the next task without rewriting the
  history of the project.

The board is a view over the files. It is not a second database and does not
silently invent state that is missing from the repository.

## What lives in the repository

Savepoint uses Markdown and YAML as its source of truth:

```text
.savepoint/
├── config.yml             # Project settings and quality gates
├── router.md              # Current workflow state and next action
├── Idea.md                # Intent, user, scope, and success criteria
├── Design.md              # Architecture and verified technical state
├── Guardrails.md          # Durable engineering policy
├── objectives/            # Outcomes and their bounded Tasks
│   └── O-001-example/
│       ├── Objective.md
│       └── tasks/
├── releases/               # Compatibility storage for required V2 Goals (R-###)
│   └── R-001-example/Release.md
├── checks/                # Independent verification evidence
└── issues/                # Durable follow-up and discovered problems

AGENTS.md                  # The agent's entrypoint
agent-skills/              # The workflow instructions it activates
```

The files are intentionally ordinary. You can read them in an editor, review
them in a pull request, diff them with Git, or recover from them without a
Savepoint server.

The board calls these delivery contexts Goals. Existing record filenames and
fields remain Release-compatible (`Release.md`, `R-###`, and `release:`).

## Safe migration

Existing V1 projects can be moved to V2 with a preview-first workflow:

```bash
npx savepoint migrate
npx savepoint migrate --apply
```

The preview reports planned records, identity mappings, archived source,
conflicts, and decisions that still need an owner. Nothing is written unless
`--apply` is used. Apply requires a Git work tree, and every planned write or
removal path must be free of modified, untracked, or ignored files. Commit or
stash changes at those paths and move ignored files away from them before
retrying. Legacy Release PRDs
become accountable `R-###` records that the V2 board presents as Goals, with
exact archive mappings; historical completion is displayed as historical
evidence, never as a fabricated current Check. The migration writes converted
files directly, removes archived source files, records the source hashes and
mappings in `.savepoint/migrations/`, and
sets `schema_version: 2` last, creating `config.yml` at that final step when
the V1 project did not have one. If an apply error occurs after writes, its
output lists the written paths and prints Git commands scoped to those paths
to restore tracked inputs and remove generated outputs.

For existing projects, refresh only the shipped Savepoint assets with:

```bash
npx savepoint upgrade-assets --dry-run
npx savepoint upgrade-assets
```

`migrate` and `upgrade-assets` are deliberately separate. Use `migrate` for a
legacy V1 project; after it is converted to V2, use `upgrade-assets` to refresh
the managed skills and templates. `upgrade-assets` does not perform a
migration.

User-authored project files remain outside the managed asset region. Repeated
updates are designed to be safe and reviewable.

## Built for agent work

Savepoint is agent-agnostic. It works with any coding agent that can read
files, including Claude Code, Cursor, Codex, Gemini, and Aider.

The important integration point is `AGENTS.md`: it tells the agent what to
read, which workflow applies, what the current objective is, and where the
active Task sets its context boundary. The agent does not need a plugin,
account, or hosted workspace to participate.

## Principles

- **Local first:** project state stays in your Git filesystem.
- **Small context:** agents read the minimum relevant context before acting.
- **One source of truth:** the router, records, checks, and issues are plain
  files, not mirrored in a hidden service.
- **Human outcomes, machine evidence:** people judge the result; tools verify
  the implementation.
- **Safe updates:** managed scaffolding can be refreshed without silently
  overwriting user-authored content.
- **No telemetry:** Savepoint has no server, database, authentication, billing,
  or external network dependency in the core workflow.

## Development

Savepoint is a Go CLI and Bubble Tea terminal UI:

```bash
make build
make test
```

The npm package wraps platform-specific binaries so most users can run the
tool with `npx savepoint`.

Distribution validation is local and does not publish or tag a release:

```bash
make dist
make package-check
```

## License

[MIT](LICENSE) © anipatke
