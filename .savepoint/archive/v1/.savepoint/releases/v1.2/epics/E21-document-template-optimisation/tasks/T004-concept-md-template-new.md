---
id: E21-document-template-optimisation/T004-concept-md-template-new
status: done
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

- [x] `templates/project/.savepoint/Concept.md` exists with a `type: project-concept` frontmatter field.
- [x] The template covers: core impulse / raw idea, target feeling (how should the user feel using it), the problem it solves in one sentence, who it is NOT for, and 2–3 open questions to resolve before writing the PRD.
- [x] The template includes a brief header note explaining when to use Concept.md vs. PRD.md.
- [x] The file is self-contained — it does not require the PRD to exist first.
- [x] Template freshness test is updated to include `Concept.md`.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Define the `Concept.md` purpose and section structure (see Acceptance Criteria).
- [x] Write `templates/project/.savepoint/Concept.md` with concrete, opinionated section prompts.
- [x] Update `internal/init/template_freshness_test.go` to assert the new file is present and has the expected `type:` frontmatter value.
- [x] Run `go test ./internal/init` to confirm the new freshness assertion passes.
- [x] Run `make build && make test`.

## Context Log

### Purpose

Concept.md is the document that sits **before** the PRD in the workflow. Vibe coders typically arrive with a feeling or impulse, not a formed product vision, and the jump straight to PRD forces premature commitment. Concept.md gives that early energy a structured place to land — pointed prompts, no commitments — and includes an explicit promotion checklist so the author knows when to move on.

### What was created

- **New file** `templates/project/.savepoint/Concept.md` with frontmatter `type: project-concept`, `status: active`. Self-contained: no references to the PRD's content, no `{{PROJECT_NAME}}` dependency on the PRD existing.
- **Header note** `## When to use this` — two paragraphs that pin down Concept vs PRD boundaries:
  - Concept: feeling, sketch, use case; open questions; deciding if a PRD is worth writing.
  - PRD: answers "what is it," "for whom," "what's out of scope" in one sentence each; committing.
  - Promotion trigger: if the author finds themselves writing success metrics or hard constraints, they have already outgrown this file.
- **Six section prompts**, each with a one-line blockquote example modeled on a fictional "bug-report-to-PR" CLI:
  1. `## Core impulse` — raw idea, paragraph form, no jargon.
  2. `## Target feeling` — one verb-led feeling + one sentence of context; utility language explicitly disallowed.
  3. `## The problem in one sentence` — exactly one sentence; second-sentence test for splitting.
  4. `## Who this is NOT for` — bold tag + sentence; "everyone" treated as a sharpening signal.
  5. `## Open questions` — 2–3 questions, each with a one-hour answerability test.
  6. `## Promoting to PRD.md` — section-by-section mapping from Concept → PRD; deletion instructions for the open questions and the promotion section once the PRD exists.

### Why this shape

- **Self-containment** is enforced by content: no link to PRD, no fields that only make sense after the PRD exists.
- **Concrete-ness** matches the prompt style established in T002 (PRD) and T003 (Design): pointed prompt + one-line blockquote example.
- **Promotion is part of the file**, not a hidden convention. The author can read the file top-to-bottom and know both when to use it and when to leave it.
- **No stale lifecycle terms** — no `phase`, no `stage: implementation`, no legacy status values; the freshness test asserts this.

### Freshness test

Added `TestProjectConceptTemplateExists` in `internal/init/template_freshness_test.go:117`. Asserts:
- File exists at the expected path (the `readTemplate` helper would fatal on missing file).
- Frontmatter contains `type: project-concept`.
- All six section headings exist.
- None of the stale lifecycle terms (`status: todo`, `status: doing`, `status: blocked`, `status: review`, `status: audit`, `phase: build|test|audit|implementation`, ``phase` (build/test/audit)`) appear in the file.

### Verification

- `go test -count=1 ./internal/init/...` passes; new test `TestProjectConceptTemplateExists` is green.
- `make build && make test` passes across all packages.

### Files changed

- **New:** `templates/project/.savepoint/Concept.md`
- **Modified:** `internal/init/template_freshness_test.go` (added `TestProjectConceptTemplateExists`)
