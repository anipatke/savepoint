---
id: E21-document-template-optimisation/T002-prd-template-optimisation
status: planned
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

- [ ] Every section prompt is a concrete question or example, not a generic instruction.
- [ ] The template models the specificity visible in the live `.savepoint/PRD.md` (one-sentence examples, pointed constraints, explicit out-of-scope items).
- [ ] No section is removed; additions or rewrites only sharpen the prompt text.
- [ ] Template freshness test passes.
- [ ] `make build && make test` passes.

## Implementation Plan

- [ ] Read the live `.savepoint/PRD.md` to extract the level of specificity and voice that actually worked.
- [ ] Rewrite each section prompt in `templates/project/.savepoint/PRD.md` to be a pointed question or a one-line example that models the expected output.
- [ ] Run `go test ./internal/init` to check freshness coverage.
- [ ] Run `make build && make test`.
