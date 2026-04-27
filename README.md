# Savepoint — The AI-Driven Development Workflow

> **Stop the drift. Kill the bloat. Focus on the vibe.**

Savepoint is an opinionated, file-based state machine and TUI designed to carry AI agents (Claude, Cursor, Aider, Gemini, etc.) through a rigorous, documented development lifecycle.

## What is it?

A CLI + Ink-based TUI that scaffolds a repeatable workflow in your project. It turns your development process into a series of hard gates:

**PRD → Design → Epics → Tasks → Build → Audit**

By pointing any AI agent at the `.savepoint/` directory, the embedded prompt templates guide the agent to plan, implement, and verify work without you having to manage the "what's next" or worry about the agent losing context.

## Why?

Three failure modes plague AI-driven development today:

1.  **Inconsistency:** No repeatable process from a "vibe" (high-level intent) to a working MVP.
2.  **Token Bloat:** Monolithic backlogs and high-overhead tools burn context, especially on minimal AI plans.
3.  **Documentation Drift:** Design docs and agent instructions go stale the moment the first line of code is written.

Savepoint solves this with **The Audit Loop**. When an epic is finished, the next one cannot start until your documentation (Design, Agents guide, and PRD) is reconciled with the actual code that was built.

## Why I'm Sharing This

I'm building Savepoint to dogfood the very workflow it enables. This project is developed _using_ Savepoint's own conventions—a recursive, self-bootstrapping "savepoint" for the project itself.

I'm sharing my learnings because:

- **Workflow over Code:** The real lever in AI development isn't just the model; it's the structure you give it.
- **Transparency is the best teacher:** Showing how a project can be built from scratch using nothing but AI agents and a strict state machine proves the "vibe coder" thesis.
- **Community Feedback:** Building in public helps stress-test the "Audit Loop" and ensures the tool stays agent-agnostic and token-efficient.

## Target Audience

**Vibe Coders.** Builders with minimal-to-moderate experience who want AI agents to drive implementation while a structured workflow keeps the project coherent and the documentation fresh.

## Status: v1 (MVP) in Progress

We are currently building the MVP. Follow the progress in the `E##-epics` and the TUI board.

### Key Commands (Coming Soon)

| Command            | Purpose                                              |
| :----------------- | :--------------------------------------------------- |
| `savepoint init`   | Scaffold your project and get the magic prompt.      |
| `savepoint board`  | Launch the Atari-Noir TUI Kanban board.              |
| `savepoint audit`  | Run the 6-step audit pipeline to sync code and docs. |
| `savepoint doctor` | Integrity checks and semantic health reviews.        |

## License

MIT
