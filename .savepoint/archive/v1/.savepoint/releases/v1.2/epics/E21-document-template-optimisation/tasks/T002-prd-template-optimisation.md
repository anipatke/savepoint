---
id: E21-document-template-optimisation/T002-prd-template-optimisation
status: done
objective: Revise the PRD.md template so its section prompts are concrete, opinionated, and fast to complete for a vibe coder starting from a raw idea.
complexity_tier: low
complexity_reason: Single template file; improvement is editorial, no code changes required.
---

# T002: PRD.md Template Optimisation

## Problem

The current `templates/project/.savepoint/PRD.md` uses generic comment prompts (`<!-- Describe what you are building... -->`). For vibe coders, these are too open-ended — they slow down the transition from idea to committed direction. The template should model the exact level of specificity that makes the PRD useful downstream (for Design.md, agent skills, and the audit loop).

## Context Files

- `templates/project/.savepoint/PRD.md`
- `.savepoint/PRD.md`
- `internal/init/template_freshness_test.go`

## Acceptance Criteria

- [x] Every section prompt is a concrete question or example, not a generic instruction.
- [x] The template models the specificity visible in the live `.savepoint/PRD.md` (one-sentence examples, pointed constraints, explicit out-of-scope items).
- [x] No section is removed; additions or rewrites only sharpen the prompt text.
- [x] Template freshness test passes.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Read the live `.savepoint/PRD.md` to extract the level of specificity and voice that actually worked.
- [x] Rewrite each section prompt in `templates/project/.savepoint/PRD.md` to be a pointed question or a one-line example that models the expected output.
- [x] Run `go test ./internal/init` to check freshness coverage.
- [x] Run `make build && make test`.

## Context Log

### What was wrong

The seven section prompts in `templates/project/.savepoint/PRD.md` were generic "Describe …" instructions. They told the author to be specific but gave no format, no length, and no example of what "specific" actually looked like for Savepoint. Result: a vibe coder staring at `<!-- Describe what you are building -->` would write a paragraph of fluff and ship it as their PRD.

### Voice and format extracted from the live PRD

The live `.savepoint/PRD.md` works because each section commits to a single shape:

- **What it is** — three sentences (category, user action, pipeline).
- **Why** — numbered failure modes, each paired with a one-sentence fix.
- **Target user** — bold persona tag + definition + explicit "Not:" persona.
- **Headline differentiator** — bold feature tag + one-sentence justification, with a "no competitor has this" test.
- **Success metrics** — bulleted, falsifiable (numbers, time bounds, runnable checks).
- **Constraints** — hard limits only, no preferences.
- **Out of scope** — "No X. No Y. No Z." phrasing, because strong language is what keeps the AI focused.

### Changes applied

- `templates/project/.savepoint/PRD.md` rewritten end-to-end. All seven original section headings preserved; the comment-style `<!-- ... -->` prompts are replaced with a directive sentence (the "ask") followed by a blockquote example modeled on the live PRD's voice.
- Added a top-of-file italic line: "Fill each section with concrete, falsifiable claims. Vague PRDs produce vague agents." — sets the bar before the author scrolls.
- Added a closing `>` line in each section so the example is unmistakably an example, not part of the prompt.

### Verification

- `go test -count=1 ./internal/init/...` passes; no canonical-string assertions on PRD content, no scaffold/upgrade tests regressed.
- `make build && make test` passes across all packages.
