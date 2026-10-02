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
  Task and its integration.
- **Keep ownership clear.** Agents implement and prove their work; people
  decide whether the outcome is what they wanted.
- **Every project has a Goal.** The router selects a live Goal, and every
  Objective belongs to exactly one Goal. A Goal is complete when every member
  Objective is complete. Goals do not publish, deploy, tag, or generate
  changelogs.

Goals are stored under `.savepoint/releases/` as `Release.md`, and Objectives
and the router select them through a `release:` field. New Goals use `G-###`
identities; existing projects keep their stable `R-###` identities. The V2
board presents these records as Goals. `savepoint init` creates and selects
G-001, titled after the project. `savepoint migrate` retains the V1 router's
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
IDEA  ─────►  DESIGN  ─────►  TASK  ─────►  OBJECTIVE CHECK
  intent       architecture    bounded      mandatory Full
  & outcome    & guardrails    execution     integration
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
Design reconciliation.

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
| `savepoint health setup [dir] [--apply]` | Previews suggested Code Health tools; `--apply` writes the health config. See [Code Health](#code-health). |
| `savepoint health check O-### [dir]` | Collects an official Code Health snapshot during a Full Objective Check. |
| `savepoint health report [dir]` | Rewrites the Code Health report from the newest saved snapshot. |
| `savepoint upgrade-assets [dir] [--dry-run] [--force]` | Refreshes shipped skills and templates in an existing project. `--dry-run` lists what would change; a skill you edited is kept and the new text is offered beside it unless you pass `--force`. |

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
  `release:` field. With no valid Goal selected, the board shows
  `Choose a Goal` and no project-wide work.
- **Router priority** lets you focus the next task without rewriting the
  history of the project.

The board is a view over the files. It is not a second database and does not
silently invent state that is missing from the repository.

## Code Health

Code Health answers one question: is the code your agents produce staying
healthy, improving or getting worse? It reads reports from tools your project
already runs, saves what they said, and explains it in plain words on the board
and in an agent-readable report. It is supporting evidence. It is not a
guarantee that the code is correct, secure or maintainable (see
[Limits](#what-code-health-does-not-tell-you)).

### The five signals

| Signal | What it asks | Good / watch limits |
| --- | --- | --- |
| Tests | Do the tests pass? | Any failure needs attention |
| Coverage | How much of the code do tests run? | Good at 80% or more, watch from 60% |
| Complexity | How tangled are the largest functions? | Good at 10 or less, watch up to 20 (cyclomatic complexity) |
| Duplication | How much code is copied? | Good at 3% or less, watch up to 5% |
| Dependency vulnerabilities | Do dependencies have known advisories? | Any finding is at least watch; high or critical needs attention |

Every signal carries one label: **Good**, **Watch**, **Needs attention** or
**Unknown**. Unknown means nothing was measured, which is never treated as
healthy. Good also needs three comparable official checks, so a new project
starts at Watch. Thresholds are defaults; your project's
`.savepoint/health/config.json` is where they and the exact tool commands live.

### Supported report formats

Savepoint reads nine report formats. It runs no analysis of its own.

| Signal | Formats |
| --- | --- |
| Tests | `go test -json`, Vitest JUnit XML, pytest JUnit XML |
| Coverage | Go cover profile, Vitest V8 coverage, coverage.py JSON |
| Complexity | lizard CSV |
| Duplication | jscpd JSON |
| Dependency vulnerabilities | osv-scanner JSON |

Go, JavaScript/TypeScript (Vitest) and Python (pytest) are the stacks with
test and coverage readers; lizard, jscpd and osv-scanner cover many languages,
and a project can mix stacks, with one instance per stack. A stack with no
matching format is unsupported: that signal shows as not configured or
unavailable, and the other signals still work. There are no plugins and no
custom signals.

### Prerequisites

Savepoint ships as one binary and bundles or installs no providers. You install
the tools you want (for example `lizard`, `jscpd` and `osv-scanner`) and your
test runner writes its own report. A missing tool is reported as unavailable
with an instruction to install it yourself; it is never installed for you.

### Set up, check and report

```bash
savepoint health setup            # preview suggested tools and exclusions
savepoint health setup --apply    # write .savepoint/health/config.json
savepoint health check O-###      # official check, run in a Full Objective Check
savepoint health report           # rewrite .savepoint/health/report.md
```

- `setup` is for you. It proposes tools by looking at your project, skips
  Savepoint's own records and the usual dependency and generated paths, and
  skips `templates/**` and `**/testdata/**` for duplication, since those are
  copies by design. It writes nothing until `--apply`, and never removes an
  entry you already confirmed. `savepoint init` ends with the same preview.
- `check` is run by the independent checker during a Full Objective Check,
  after the full gate has written its reports. It saves one **official**
  snapshot and prints a verdict. The Check records the snapshot ID. Task Checks
  never collect.
- `report` is for you. It rewrites the derived report from the newest saved
  snapshot without running a tool.

### Official and manual evidence

An **official** snapshot comes from `savepoint health check` and is permanent.
A **manual** snapshot comes from refreshing in the board (`R`) and is for your
own curiosity: it can change what the dashboard shows, but it never counts for
sign-off and is never part of a trend. A Check may only cite an official
snapshot.

A blocking verdict (failing tests, high or critical vulnerabilities, a failed or
stale required instance, or a rule you opted in as blocking) prevents `CLEAR`.
Everything else is advisory; the checker may turn a warning into an Issue.

### Trends, history and staleness

- Trends compare only comparable official checks: same tool version, scope and
  configuration. Changing the scope or tool version restarts that series and the
  dashboard says so. Movement smaller than a per-signal noise margin is not
  reported as better or worse.
- History view (`h`) shows the last ten checks.
- Evidence is **stale** when the code has moved on since it was measured. The
  dashboard says whether the code matches, is some commits ahead or is on
  another branch, using read-only local Git queries. A required stale instance
  blocks an official check.
- The dashboard reads a bounded window of recent history (the ten newest
  official checks and anything saved after the oldest of them), so opening
  Code Health stays quick as history grows. When older history was not checked,
  the explanation says so. For a series present in every checked snapshot the
  trend, labels and sign-off match what full history gives. A signal added or
  removed within the older history can show a shorter trend (or "No trend
  yet") than full history would.
- Savepoint never prunes history automatically. Official snapshots are
  permanent. Internal maintenance exists that can keep the ten newest manual
  snapshots, but there is no owner-facing cleanup command in this release
  (see the v2.1 limitations in [CHANGELOG.md](CHANGELOG.md)).

### A walkthrough in the board

1. Open `savepoint board`. The header shows `Health: not set up`,
   `Health: no check yet`, or a heart with how many signals are Good, such as `2/5`.
2. Press `H` to open Code Health. Five rows, one per signal, each with its
   label, the number, an aim, a short plain-language meaning and a trend line.
3. Move with `↑`/`↓`. The selected signal shows its question and number, what
   it means, whether it counts for sign-off, where the evidence came from and
   a next step.
4. Press `h` to see the last ten checks; `h` or `esc` goes back.
5. Press `R` to refresh now. This runs your configured tools and saves a manual
   snapshot; `esc` cancels. The screen says a manual refresh is not an official
   check.
6. Press `esc` to leave. After any saved snapshot Savepoint also rewrites
   `.savepoint/health/report.md`.

### Hand a finding to an agent

`.savepoint/health/report.md` is written for an agent. It lists each signal
that is not Good, blocking ones first, with the question, label, number, aim,
meaning, next step, trend and affected files, followed by the instruction to investigate,
propose a fix, apply it and re-run the official check. Complexity has an aim
(10 or less) and a watch line (20 or less); the shipped skills tell an agent to
bring a signal back to the watch line, not chase the aim, and to report what
remains. Give the file to your
agent, or point it at the signal you care about. The agent does not run
`setup` or `report`; it re-runs `savepoint health check O-###` only inside a
Full Objective Check.

### Privacy and local-first behaviour

Code Health is local-first. Savepoint sends no telemetry, makes no network
requests of its own, and installs nothing. Reports and snapshots stay in
`.savepoint/health/`.

The tools you configure are yours, and they run only when you or an
authorised workflow runs them. They run as a plain argument list (no shell) in
your project directory and inherit Savepoint's environment. A provider may use
the network: for example `osv-scanner` sends dependency names and versions to
`api.osv.dev`, and fails (so the signal is Unknown, not clean) when it cannot
reach it. Savepoint does not sandbox providers or make them work offline.
Output is bounded: the tail of an oversized error message is kept so the real
failure stays visible.

### What Code Health does not tell you

- It measures only what your tools report, within the scope and exclusions you
  confirmed. Passing tests, good coverage, low complexity and no known
  advisories do not prove the code is correct, safe or well designed.
- **Partial** or **unavailable** evidence is shown as exactly that. A failed,
  timed-out, missing or malformed report is "not measured", never "zero
  problems" and never "bad code". A partial report may carry a value, but it
  covers only part of the code, so it is shown as incomplete and never counts
  as Good. A failed vulnerability scan does not mean the
  project is free of vulnerabilities.
- Numbers from different scopes or tool versions are not comparable. Narrowing
  a scope can lower a number without improving the code.
- Duplication is measured over maintained code. Savepoint's own records,
  archived history, template mirrors (which must match the shipped skills byte
  for byte) and test fixtures are excluded on purpose.

## What lives in the repository

Savepoint uses Markdown and YAML as its source of truth:

```text
.savepoint/
├── config.yml             # Project settings and quality gates
├── router.md              # Current workflow state and selection
├── Idea.md                # Intent, user, scope, and success criteria
├── Design.md              # Architecture and verified technical state
├── Guardrails.md          # Durable engineering policy
├── health/                # Code Health config, snapshots, and report
├── objectives/            # Outcomes and their bounded Tasks
│   └── O-001-example/
│       ├── Objective.md
│       └── tasks/
├── releases/              # Storage for required V2 Goals (G-###, or R-### if migrated)
│   └── G-001-example/Release.md
├── checks/                # Independent verification evidence
└── issues/                # Durable follow-up and discovered problems

AGENTS.md                  # The agent's entrypoint
agent-skills/              # The workflow instructions it activates
```

The files are intentionally ordinary. You can read them in an editor, review
them in a pull request, diff them with Git, or recover from them without a
Savepoint server.

The board calls these delivery contexts Goals. Record filenames and fields
stay Release-compatible (`Release.md` and `release:`).

## Safe migration

Existing V1 projects can be moved to V2 with a preview-first workflow:

```bash
npx savepoint migrate
npx savepoint migrate --apply
```

The preview is a short summary: the Goals, Objectives, Tasks, and Issues it
will create, how many V1 files it will archive, any conflicts, and any
decisions that still need an owner. Nothing is written unless `--apply` is
used. When a V1 record has a status V2 cannot map on its own, the preview
lists the choices and the exact entries to put in a decisions file:

```bash
npx savepoint migrate --decisions decisions.yml
npx savepoint migrate --decisions decisions.yml --apply
```

Add `--verbose` to list every planned record, identity mapping, archived
file, and advisory note. Apply requires a Git work tree, and every planned write or
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

User-authored project files remain outside the managed asset region. A skill
you have edited is never overwritten without `--force`, and repeated updates
are designed to be safe and reviewable.

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
