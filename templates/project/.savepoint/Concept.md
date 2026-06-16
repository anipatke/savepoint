---
type: project-concept
status: active
---

# {{PROJECT_NAME}} — Concept

> Use this file **before** the PRD. Concept.md is where raw ideas land; PRD.md is where you commit. Fill in what you know, leave the rest as open questions, and only promote to PRD when the answers are concrete enough to defend.

## When to use this

- **Concept.md:** you have a feeling, a sketch, a half-formed use case. The hardest questions (scope, constraints, metrics) are still open. You're trying to decide if this is worth a PRD.
- **PRD.md:** you can already answer "what is it," "for whom," and "what's out of scope" in one sentence each. You're committing to a direction, not exploring one.

If you find yourself writing success metrics or hard constraints here, you have already outgrown this file — promote to PRD.md.

## Core impulse

What is the raw idea? Write it the way you'd describe it to a friend on a walk, not the way you'd describe it in a deck. One paragraph, no jargon, no positioning.

> A CLI that reads your bug report, opens a fix branch, writes the failing test, and opens the PR with the test passing — all before you've finished typing the original ticket.

## Target feeling

How should the user feel **while using it**? Pick one verb-led feeling (relieved, smug, in flow, surprised, calm) and one sentence of context. Avoid utility language ("fast," "easy," "intuitive") — those describe capability, not feeling.

> **Relieved.** The bug report is gone, the fix is in flight, and I didn't have to context-switch into ticket hygiene to get there.

## The problem in one sentence

State the problem this product removes in exactly one sentence. If you need a second sentence, the problem is two problems — pick one.

> Bug reports stall in triage because the report-to-fix loop is bottlenecked on the human's attention, not the human's skill.

## Who this is NOT for

Name the persona you are explicitly **not** serving. One bold tag, one sentence of definition. If the persona is "everyone," pick a sharper one — "everyone" means the idea is not sharp yet.

> **Not: enterprise teams with formal SDLC tooling.** They already have Jira + Jenkins + on-call rotation; this is friction removal for the rest of us.

## Open questions

List 2–3 questions that must be answered before you can write a defensible PRD. Each question is a test: if you can't answer it after one focused hour, the idea is not ready for PRD.

> 1. **What is the minimum input the CLI needs to act?** A bug report alone? A repo URL? An issue tracker token? Pick the laziest input that still works.
> 2. **Where does the human approve the PR?** Auto-merge on green CI? A reviewer ping? A `?` confirmation prompt in the CLI? Each answer changes the trust model.
> 3. **What does this cost to run per fix?** Token cost + CI minutes + LLM latency. If a single fix costs more than a human doing it manually, the product is a demo, not a tool.

## Promoting to PRD.md

When all three open questions have answers, copy this file's sections into PRD.md and rewrite them in committed form:

- "Core impulse" → PRD's "What it is" (three sentences, not a paragraph).
- "Target feeling" → PRD's "Why" (failure modes that this feeling fixes).
- "The problem in one sentence" → PRD's "Why" lead.
- "Who this is NOT for" → PRD's "Target user" (with a positive persona added).
- "Open questions" → PRD's "Constraints" and "Out of scope."

Delete the open questions and the promotion section once the PRD exists.
