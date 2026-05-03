---
id: E07-init-command/T005-magic-prompt
status: done
objective: "Render and output magic prompt to stdout"
depends_on: ["E07-init-command/T004-atomic-writes"]
---

# T005: Magic Prompt

## Acceptance Criteria

- Reads `templates/prompts/magic-prompt.prompt.md` template
- Renders template with project context
- Outputs rendered prompt to stdout
- Prompt is copy-paste ready for AI agent

## Implementation Plan

- [x] Add `internal/init/prompt.go`
- [x] Implement `RenderMagicPrompt(projectName) (string, error)`
- [x] Load magic-prompt.prompt.md template
- [x] Interpolate project name into prompt
- [x] Output to stdout after scaffold completes
- [x] Test prompt rendering
- [x] Run `make build && make test`