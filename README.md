# SAVEPOINT

> **HARD GATES. ZERO DRIFT. ALL VIBE.**

Savepoint is a file-based state machine and TUI designed to carry AI agents (Claude, Cursor, Aider, Gemini) through a rigorous development lifecycle. It’s the wedge that stops agents from hallucinating your architecture into a corner.

---

### THE MISSION

Most AI-driven development falls apart for three reasons:

1. **The Drift:** Documentation and code stop talking to each other after day one.
2. **The Bloat:** Monolithic backlogs burn your token budget and confuse the agent.
3. **The Chaos:** No repeatable process from a "vibe" to a working MVP.

Savepoint fixes this by turning your project into a series of **Hard Gates**.

### THE LOOP

`PRD` → `DESIGN` → `EPICS` → `TASKS` → `BUILD` → **`[ AUDIT ]`**

The **Audit Loop** is the heart of the machine. When an epic is finished, the next one stays locked until your docs (Design, PRD, and Agent guides) are reconciled with the actual code that was built.

**No stale designs. No documentation debt. Just clean, verified progress.**

### WHY I'M SHARING

I'm building Savepoint to dogfood the very workflow it enables. This entire repository is being built _by_ agents, _through_ Savepoint’s own state machine.

I’m sharing the journey because:

- **Workflow > Model:** The real power of AI isn't just the LLM; it's the structure you give it.
- **Radical Transparency:** Showing how to build a complex tool from scratch using nothing but "Vibe Coding" and a strict process.
- **The Wedge:** Proving that token-efficient, documentation-first development is the only way to build at scale with AI.

### THE STACK (ATARI-NOIR)

- **CLI/TUI:** Built with Ink for a cinematic, technical experience.
- **FILE-FIRST:** Your filesystem is the database. No MCP server, no proprietary cloud.
- **AGENT-AGNOSTIC:** Works with any agent that can read Markdown and edit files.

---

### COMMANDS (0.1.0-MVP)

| Command            | Action                                      |
| :----------------- | :------------------------------------------ |
| `savepoint init`   | Scaffold the loop and get the magic prompt. |
| `savepoint board`  | Launch the Atari-Noir Kanban TUI.           |
| `savepoint audit`  | Sync the map with the territory.            |
| `savepoint doctor` | Check the integrity of the machine.         |

**License:** MIT  
**Status:** Recursive Construction (v1 MVP in progress)
