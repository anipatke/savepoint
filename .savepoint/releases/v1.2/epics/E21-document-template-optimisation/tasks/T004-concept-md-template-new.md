---
id: E21-document-template-optimisation/T004-concept-md-template-new
status: planned
objective: Author the Concept.md template — a lightweight pre-PRD ideation document for capturing raw ideas, target feelings, and open questions before committing to a product direction.
complexity_tier: medium
complexity_reason: New document type with no prior art; requires defining purpose, structure, and PRD handoff.
---

# T004: Concept.md Template (New)

## Problem

There is no scaffolded document for the ideation phase that precedes the PRD. Vibe coders often arrive with a feeling or impulse — not a formed product vision — and the jump straight to `PRD.md` forces premature commitment. A `Concept.md` gives that early energy a structured place to land before the harder questions of scope, constraints, and metrics are asked.

## Context Files

- `templates/project/.savepoint/PRD.md`
- `templates/project/.savepoint/Design.md`
- `templates/project/AGENTS.md`
- `internal/init/template_freshness_test.go`

## Acceptance Criteria

- [ ] `templates/project/.savepoint/Concept.md` exists with a `type: project-concept` frontmatter field.
- [ ] The template covers: core impulse / raw idea, target feeling (how should the user feel using it), the problem it solves in one sentence, who it is NOT for, and 2–3 open questions to resolve before writing the PRD.
- [ ] The template includes a brief header note explaining when to use Concept.md vs. PRD.md.
- [ ] The file is self-contained — it does not require the PRD to exist first.
- [ ] Template freshness test is updated to include `Concept.md`.
- [ ] `make build && make test` passes.

## Implementation Plan

- [ ] Define the `Concept.md` purpose and section structure (see Acceptance Criteria).
- [ ] Write `templates/project/.savepoint/Concept.md` with concrete, opinionated section prompts.
- [ ] Update `internal/init/template_freshness_test.go` to assert the new file is present and has the expected `type:` frontmatter value.
- [ ] Run `go test ./internal/init` to confirm the new freshness assertion passes.
- [ ] Run `make build && make test`.
