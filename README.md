![Savepoint Banner](assets/banner.png)

# Savepoint

**A little terminal app that helps you keep track of what you're building with AI.**

Local files. Tight context. No telemetry. Free.

[Website](https://www.getsavepoint.dev/) ·
[GitHub](https://github.com/anipatke/savepoint) ·
[npm](https://www.npmjs.com/package/savepoint)

Are you an elite developer with 200 repos, enterprise tooling and strong
opinions about Kubernetes? Cool. Savepoint definitely isn't for you.

Are you a bored dad trying to build something cool on the weekend before
someone asks you to mow the lawn? Now we're talking.

## What it actually does

AI agents wander off. Context balloons, scope gets fuzzy, and "tests passed"
quietly replaces "somebody looked at it." Savepoint puts a few gates in the way:

- **Plan first.** Idea and Design get written down before an agent touches code.
- **Small tasks.** Each Task has one outcome and a short list of files the agent
  is allowed to read.
- **Proof before moving on.** An independent check verifies the work before an
  Objective closes. If the plan turns out to be wrong, the agent stops and says
  `REPLAN REQUIRED` instead of freestyling a new architecture.
- **You decide.** Agents build and prove. You say whether it's what you wanted.

Everything is Markdown and YAML in your repo. Every project has a Goal, Goals
contain Objectives, and Objectives contain Tasks.

```text
IDEA  ──►  DESIGN  ──►  TASK  ──►  CHECK
```

## Quick start

```bash
npx savepoint init      # scaffold into your repo
npx savepoint board     # open the TUI
npx savepoint resume    # what's next?
npx savepoint doctor    # sanity check
```

Then tell your coding agent to read the generated `AGENTS.md`. It points the
agent at the current state, the right workflow skill, and the files it's
allowed to look at. Commit `.savepoint/` and `AGENTS.md` so the next session
(or teammate) picks up where you left off.

Works with anything that can read files: Claude Code, Cursor, Codex, Gemini,
Aider.

## Commands

| Command | What it does |
| --- | --- |
| `savepoint init [dir]` | Scaffolds project files and agent guidance. |
| `savepoint create-task --objective O-### --draft <path> [dir]` | Creates a Task from an ID-free draft. |
| `savepoint board [--objective O-###]` | Opens the keyboard-driven board. |
| `savepoint resume [dir]` | Prints the current state and next action. |
| `savepoint doctor` | Runs project diagnostics and quality gates. |
| `savepoint migrate [dir]` | Converts a V1 project to V2. Preview by default; add `--apply`. |
| `savepoint upgrade-assets [dir] [--dry-run] [--force]` | Refreshes shipped skills and templates without clobbering your edits. |
| `savepoint health setup\|check\|report` | Code Health, see below. |

`create-task` takes an Objective and a complete Task draft without an `id`,
assigns the next project-wide ID and path, and strict-loads the full V2 index
before it reports success. Agents use it only for new Tasks.

Run `savepoint <command> --help` for the details.

## Code Health

Is the code your agents write getting better or worse? Savepoint reads reports
from tools you already run (tests, coverage, complexity, duplication,
dependency vulnerabilities), saves them, and shows five signals on the board
(press `H`). It also writes `.savepoint/health/report.md` for you to hand to an
agent.

It runs no analysis of its own, installs nothing, and is supporting evidence,
not proof the code is correct or secure. Missing or broken reports show as
"Unknown", never as "fine". See `savepoint health setup` to get started.

You bring the tools and the reports. `savepoint health setup` lists what it
found, how to install anything missing (`lizard` from pip, `jscpd` from npm,
`osv-scanner` from its releases), and the command that writes each test or
coverage report, such as `vitest run --reporter=junit --outputFile=junit.xml`.
Run those first, then `savepoint health check`. Add the generated files
(`junit.xml`, `coverage/`) to `.gitignore`.

## What's in your repo

```text
.savepoint/
├── config.yml       # settings and quality gates
├── router.md        # current state and selection
├── Idea.md          # intent and scope
├── Design.md        # architecture
├── Guardrails.md    # durable engineering rules
├── objectives/      # outcomes and their Tasks
├── releases/        # Goals
├── checks/          # verification evidence
├── issues/          # follow-ups
└── health/          # Code Health config and snapshots

AGENTS.md            # the agent's entrypoint
agent-skills/        # the workflow instructions it activates
```

Plain files. Read them in an editor, diff them in Git, recover from them
without a server.

## Upgrading

Coming from V1? `npx savepoint migrate` shows a preview and writes nothing
until you add `--apply`. Already on V2? `npx savepoint upgrade-assets --dry-run`
shows what would change. Skills you've edited are kept.

Already have `savepoint` in `package.json`? Update it (`npm install -D
savepoint@latest`) or remove it and use `npx savepoint@latest`. `npx` runs the
copy in your project first, so an old 1.x install runs instead of this one.
`doctor`, `migrate` and `upgrade-assets` warn you when they find one.

## FAQ

**Why a TUI?** Because it makes me feel like Crash Override from the movie *Hackers*.  
**Will it save you tokens?** No goddamn idea. I don't know how to measure that.  
**Is it secure?** Sounds like a question for someone who's actually read the code.  
**Is it production ready?** What's a production?  
**What's the roadmap?** Depends how loudly you complain.  
**What does it cost?** Nothing financially. Emotionally, unclear.  

It's fully local, with no server, no account and no telemetry. Try it, break
it, tell me what's stupid.

## Development

Go CLI with a Bubble Tea TUI. The npm package wraps the platform binaries.

```bash
make build
make test
make dist && make package-check   # local only, publishes nothing
```

## License

[MIT](LICENSE) © anipatke
