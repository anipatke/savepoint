![Savepoint Banner](assets/banner.png)

# Savepoint

> **Hard gates for AI-driven development. Local files, tight context, no telemetry.**<br>
> Official site: [getsavepoint.dev](https://getsavepoint.dev)

Savepoint is a local-first CLI and Bubble Tea terminal board that keeps AI-assisted software projects inside a disciplined engineering workflow.

It acts as a control layer between you and your coding agents (Claude Code, Cursor, Codex, Gemini, Aider). It gives agents a simple 4-beat rhythm to follow, exact scoped context files to read, and hard gates before work drifts away from the plan:

$$\textbf{Idea} \longrightarrow \textbf{Design} \longrightarrow \textbf{Task} \longrightarrow \textbf{Check}$$

```text
Plan deeply. Execute cheaply. Check independently.
```

**No database. No proprietary cloud. No telemetry. Your Git filesystem is the source of truth.**

---

## Quick Start

```bash
npx savepoint init
npx savepoint board
npx savepoint doctor
```

* `init` scaffolds `.savepoint/`, `AGENTS.md`, and agent skills into your repository.
* `board` opens the Atari-Noir keyboard-driven terminal dashboard.
* `doctor` deterministically checks repository sanity, router state, and quality gates.

After `init`, point your agent at `AGENTS.md` and let it follow the plan.

---

## The Core Philosophy

### 1. The Division of Labor
> **The user validates outcomes. Savepoint verifies implementation.**

You shouldn't have to spend your Sunday reviewing 600-line git diffs of generated code just to know if your app is safe.

* **Your job:** Validate the outcome. *Does the button work? Does the screen look right? Can I complete the workflow?*
* **Savepoint's job:** Verify the code. *Did the agent violate `Guardrails.md`? Did it touch files it wasn't supposed to touch? Did it drift from `Design.md`? Did unit tests pass?*

### 2. The Tri-Model Architecture
Smart models are too expensive to write every line of code. Cheap models are too dumb to design systems. And **no model should ever grade its own homework**.

Savepoint splits AI work into three distinct capability roles:

| Role | Capability | Responsibility |
| :--- | :--- | :--- |
| **Planner** | Frontier reasoning (e.g., Claude 3.7 / Opus / GPT-4.5) | Refines the **Idea**, produces technical **Design**, settles durable **Guardrails**, and decomposes work into small, bounded **Tasks**. Does the expensive thinking up front—once. |
| **Executor** | Fast & budget-friendly (e.g., Haiku / Gemini Flash / GPT-4o-mini) | Executes one **Task** at a time within strictly scoped files. Never improvises architecture. If blocked or if the plan is wrong, raises its hand: `REPLAN REQUIRED`. |
| **Checker** | Independent reasoning (e.g., Sonnet / Codex) | Skeptically tests the completed Task. Challenges executor claims. Verifies technical integrity, tests, and guardrails. Returns `CLEAR` or `NEEDS WORK`. |

---

## The 4-Beat Rhythm

```text
IDEA ───► DESIGN ───► TASK ───► CHECK
                        ▲         │
                        └─────────┘
```

1. **Idea (`Idea.md`):** What are we building, who is it for, and why? Start with a rough sentence; let the planning model refine the scope and explicit out-of-scope boundaries.
2. **Design (`Design.md` & `Guardrails.md`):** Architecture before code. Major components, data flow, boundaries, and 10–20 durable guardrails. Settled before implementation starts.
3. **Task (`tasks/T###-slug.md`):** Bounded execution packets. One discrete, observable outcome with strictly scoped context files (2–3 files max).
4. **Check:** Independent verification. Tests pass? Guardrails intact? No design drift? Produces a simple verdict: `CLEAR` or `NEEDS WORK`.

---

## What Savepoint Creates

Savepoint stores project state directly in Markdown and YAML frontmatter next to your code:

```text
.savepoint/
├── Idea.md             # What we're building & why (replaces PRD)
├── Design.md           # Architecture, components, data flow, codebase map
├── Guardrails.md       # Durable constraints the agent must not break
├── router.md           # Current state machine & active task pointer
├── objectives/         # Objectives group related tasks (optional for small projects)
│   └── O001-example/
│       ├── Objective.md
│       └── tasks/
│           ├── T001-setup.md
│           └── T002-feature.md
├── checks/             # Check evaluation results & verification evidence
└── issues/             # Durable follow-up: defects, drift, and guardrail items
AGENTS.md               # The single entrypoint for your coding agents
agent-skills/           # Workflow instructions for planner, executor, and checker
```

---

## Task as a Bounded Execution Packet

In Savepoint, a Task is not a vague to-do item. It is a **handoff contract**:

```markdown
---
id: O003/T004-resume-work
status: in_progress
stage: build              # build | test | audit
objective: O003-improve-project-recovery
depends_on: []
planned_by: planner
---

# T004: Resume unfinished work

## Outcome
When I reopen Savepoint, I can see what I was working on and what to do next.

## User Check
1. Start a task, interrupt Savepoint, reopen project.
2. Run `savepoint resume`.
3. Confirm current task and next action are shown accurately without disk writes.

## Context Files (Strictly Scoped)
- `cmd/resume.go`
- `internal/data/project.go`

## Guardrails
- `STATE-01` (read-only execution)
- `DATA-03` (use canonical parser)

## Implementation Plan
1. Resolve active task from router state.
2. Implement read-only query helper.
3. Add CLI command wiring.
4. Add regression tests for corrupted router state.

## Boundaries
Do not redesign router persistence or add automatic model routing.

## Technical Verification
- [ ] Unit tests pass (`make test`).
- [ ] Zero filesystem writes during execution.
```

If an executor gets stuck or discovers an architectural ambiguity, it returns:
```text
REPLAN REQUIRED
The Task assumes router state exposes Objective directly, but data model requires derivation.
Planner decision required.
```
This halts execution cleanly instead of letting the agent improvise rogue code.

---

## The Terminal Board (`savepoint board`)

`savepoint board` launches the retro **Atari-Noir** Bubble Tea terminal interface:

* **Header:** Displays release status, active objective, and open defect warnings (`⚠ 1 open`).
* **Next Activity Line:** The exact next step derived from `.savepoint/router.md`.
* **Kanban Columns:** `PLANNED`, `IN PROGRESS` (with `[build]`, `[test]`, `[audit]` stage tags), and `DONE`.
* **Sidebar:** Fast navigation across Epics and Objectives.
* **Overlays:**
  * Press `Enter` on any card to view the **Task Detail** modal.
  * Press `d` to open the **Defects Overlay** for release-level bugs and regressions.
  * Press `A` to inspect the **Audit Register / Issues** overlay.
* **Keyboard-Driven:** Fast vim/arrow key navigation (`p` to set router priority, `q` to quit).

---

## Project Sanity (`savepoint doctor`)

Run `savepoint doctor` to run deterministic sanity checks on your project:

```text
$ savepoint doctor
savepoint doctor report
────────────────────────────────

◆ Config Check
  ✓ config

◆ Router Check
  ✓ router (active: O003/T004)

◆ Project Check
  ✓ no problems

◆ Structure Check
  ✓ no problems

◆ Defect Check
  ✓ 1 open defect tracked (D003)

◆ Quality Gates
  [PASS] build (make build)
  [PASS] test (make test)

result: ALL CLEAN (exit code 0)
```

Real tools measuring real files—not AI opinion.

---

## Migrating a Legacy Project (`savepoint migrate`)

`savepoint migrate [dir]` converts a V1 project into V2 once: fresh Objectives, Tasks, and Issues for active work, a byte-preserved archive of everything it replaces, and a recorded reference map at `.savepoint/migrations/v1-to-v2.yml`.

**Preview is the default and writes nothing.** Run it with no flags (or `--dry-run`, an explicit synonym for the same default) to see every planned record, archive, identity mapping, conflict, and owner decision required, with no file, directory, or backup created:

```bash
npx savepoint migrate
```

Add `--apply` when you are ready to write:

```bash
npx savepoint migrate --apply
```

Passing `--apply` together with `--dry-run` still previews — `--dry-run` always wins.

* `--decisions FILE` supplies concrete answers to any blocking ambiguity the preview names (an unrecognized status, a missing dependency target, a duplicate task id, an unresolved duplicate finding). A plan with unresolved blocking ambiguities refuses to apply and exits nonzero, naming every unresolved ID.
* `--recover` reports an incomplete migration operation and how to finish it; combined with `--apply`, it resumes and completes that operation.
* A missing directory, an unwritable directory, or a directory that is not a Savepoint project each fail with a distinct, named error rather than a partial write.

---

## Defects & Issues

Savepoint distinguishes between planned tasks and discovered problems:
* **Defect:** Observed behavior is wrong, broken, or regressed.
* **Issue:** Umbrella tracking for defects, architectural drift, guardrail violations, or required owner verifications.
* Discovered during development or checks; tracked with stable IDs so the same issue isn't rediscovered every run.

---

## Design Principles

* **Simple surface, rigorous engine:** The terminal board hides complexity; the underlying files enforce discipline.
* **File-first & Local-only:** Markdown and YAML are the database. No cloud lock-in, no telemetry, no tracking.
* **Agent-agnostic:** Works with Claude Code, Cursor, Codex, Gemini, Aider, or any tool that reads files.
* **Token-efficient:** Bounded context packets prevent agents from blowing 100k tokens on chat history.
* **Safe updates:** User-authored files (`Idea.md`, `Design.md`, `Guardrails.md`, Tasks) are **never** silently overwritten.

---

## Development

Build and test the CLI locally:

```bash
make build
make test
```

* **CLI & TUI:** Written in Go with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss).
* **Distribution:** Packaged via npm (`npx savepoint`) wrapping cross-compiled native binaries.
* **Marketing Site:** Ultra-lightweight static site at [getsavepoint.dev](https://getsavepoint.dev).

---

## License

[MIT](LICENSE) © anipatke
