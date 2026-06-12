---
id: E22-agent-launcher-foundation/T002-launch-request-and-prompts
status: planned
objective: Define launch requests and deterministic scoped prompts for every supported action.
depends_on:
  - E22-agent-launcher-foundation/T001-launcher-config-contract
complexity_tier: medium
complexity_reason: Five action types share prompt rules but require distinct lifecycle and scope instructions.
---

# T002: Launch Request and Prompts

## Problem

The launcher needs one structured input contract and predictable instructions that keep each new agent session focused on the selected task, defect, or epic.

## Context Files

- `internal/launcher/request.go`
- `internal/launcher/request_test.go`
- `internal/launcher/prompt.go`
- `internal/launcher/prompt_test.go`
- `internal/data/task.go`
- `internal/data/defect.go`
- `internal/data/router.go`

## Acceptance Criteria

- [ ] Typed requests cover task Build/Audit, defect Build/Audit, and epic Audit.
- [ ] Every request contains the project root, release, item identifier, item file, action, and agent role.
- [ ] Task and defect prompts name the exact item file and require its declared context boundaries.
- [ ] Epic prompts require the existing fresh-session audit workflow and active audit gate.
- [ ] All prompts prohibit automatic `done`, `resolved`, or `audited` transitions.
- [ ] Prompt output is deterministic and safe to pass as one structured process argument.

## Implementation Plan

- [ ] Add request and action types with validation for required identifiers and paths.
- [ ] Implement separate small prompt builders for build, item audit, and epic audit behavior.
- [ ] Include the root `AGENTS.md` entrypoint and selected artifact path in every prompt.
- [ ] Encode scope language without duplicating full phase-skill instructions.
- [ ] Add exact-output or focused fragment tests for all supported action types and invalid requests.

## Context Log

Pending.
