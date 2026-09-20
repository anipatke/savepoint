---
id: E47-onboarding-upgrades/T003-begin-from-one-rough-sentence
title: Begin from one rough sentence
status: done
objective: Rewrite the bootstrap prompt and the V2 router's opening next_action so a rough sentence is a valid start that routes straight to savepoint-idea.
depends_on:
    - E47-onboarding-upgrades/T002-start-a-new-project-on-the-v2-workflow
complexity_tier: low
complexity_reason: Two authored files and their assertions; no code path changes.
---

# T003: Begin from one rough sentence

## Problem

The prompt `init` prints and copies to the clipboard is the first thing an agent reads in a new project, and it currently routes into the V1 state machine: read the agent guide, read the router, follow its next action — which, in a V1 project, means `pre-implementation` and `savepoint-draft-prd`.

After T002 that prompt lands in a V2 project whose router opens at `state: idea`. The route it names is gone, and the handoff it needs does not exist: `savepoint-idea` accepts a single rough sentence as a valid starting input, and the prompt is the only place a user learns that before they have typed anything.

The Idea skill states the rule for the agent that has already been routed to it. Nothing states it for the user sitting in a fresh project deciding whether they need to go write a requirements document first. That is what this task fixes, in the two files a new project is read from: the bootstrap prompt and the router's opening `next_action`.

The prompt stays bootstrap-only. It names where to go; it does not restate the Idea workflow, the section list, or anything else the skill owns — a prompt that duplicates a skill is a second copy to keep in sync, and `TestIntegration_MagicPromptDoesNotCarryPhaseWorkflow` already exists to prevent exactly that. There is one prompt template, not one per lifecycle, because after T002 `init` produces only V2 projects.

## Context Files

- `templates/prompts/magic-prompt.prompt.md`
- `templates/project-v2/.savepoint/router.md`
- `templates/project-v2/.savepoint/Idea.md`
- `agent-skills/savepoint-idea/SKILL.md`
- `internal/init/prompt.go`
- `internal/init/prompt_test.go`
- `internal/init/integration_test.go`
- `main.go`

## Acceptance Criteria

- [x] `templates/prompts/magic-prompt.prompt.md` tells the reader that one rough sentence is enough to start and that no prepared requirements document is needed.
- [x] It names `savepoint-idea` as where that sentence goes, and names the agent guide and `.savepoint/router.md` as what to read first.
- [x] It contains no V1 routing vocabulary: no `pre-implementation`, no `epic`, no `PRD`, no `savepoint-draft-prd`.
- [x] It restates no part of the Idea workflow — not the section list, not the question sequence — and `TestIntegration_MagicPromptDoesNotCarryPhaseWorkflow` passes against the rewritten text.
- [x] `templates/prompts/` still holds exactly one prompt template, and `TestPromptTemplates_onlyMagicPromptRemains` passes unchanged.
- [x] `{{PROJECT_NAME}}` still interpolates, and `TestRenderMagicPrompt_interpolatesAllVariables` passes.
- [x] The V2 router's opening `next_action` names the same starting point in one plain sentence, consistent with the prompt and with `state: idea`.
- [x] `savepoint init` prints the rendered prompt and attempts the clipboard copy exactly as it does today; the clipboard result branches are unchanged.
- [x] `agent-skills/savepoint-idea/SKILL.md` and its shipped copy are unchanged, byte for byte.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Rewrite `templates/prompts/magic-prompt.prompt.md` for V2 bootstrap: what this project is, what to read, that a rough sentence is a valid start, and where it goes.
- [x] Set the V2 router's opening `next_action` to match, keeping it to one sentence.
- [x] Add prompt-content assertions for the rough-sentence statement, the `savepoint-idea` handoff, and the absence of V1 vocabulary.
- [x] Re-run the existing prompt tests and adjust the new text, not the tests, where they fail.
- [x] Verify end to end: scaffold a temporary project, render the prompt, and confirm the route it names resolves to a skill the scaffold actually ships.
- [x] Run `go test ./internal/init/...`, then `make build && make test`.

## Context Log

**Files read:** `templates/prompts/magic-prompt.prompt.md`, `templates/project-v2/.savepoint/router.md`, `templates/project-v2/.savepoint/Idea.md`, `agent-skills/savepoint-idea/SKILL.md`, `internal/init/prompt.go`, `internal/init/prompt_test.go`, `internal/init/integration_test.go`, `main.go`, `.savepoint/Guardrails.md`.

**Files edited:**
- `templates/prompts/magic-prompt.prompt.md` — added the rough-sentence / no-prepared-document statement and the `savepoint-idea` handoff line; no V1 vocabulary, no Idea workflow detail.
- `templates/project-v2/.savepoint/router.md` — condensed the opening `next_action` to one sentence carrying the same rough-sentence statement.
- `internal/init/prompt_test.go` — extended `TestPromptTemplates_magicPromptIsBootstrapOnly` with assertions for `"rough sentence"`, `"savepoint-idea"`, and absence of `epic`/`PRD`/`savepoint-draft-prd`.

**Verification:**
- `go test ./internal/init/...` — pass.
- `make build && make test` — pass (all packages).
- End-to-end manual check (throwaway test, removed after use): scaffolded a temp project from `templates/project-v2`, rendered the real prompt template through `RenderMagicPrompt`, confirmed the router's `next_action` reads as one sentence naming the rough-sentence start, and confirmed `agent-skills/savepoint-idea/SKILL.md` is present in the scaffolded project.
- `diff agent-skills/savepoint-idea/SKILL.md templates/project-v2/agent-skills/savepoint-idea/SKILL.md` — identical; this task made no edit to either.
- No `.savepoint/Health-Check.md` in this project — Quick check step skipped per `savepoint-build-task`.
