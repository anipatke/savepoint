---
id: E06-audit-command/T004-ai-reconcile
status: planned
objective: "Generate AI prompt for semantic review, parse findings"
depends_on: ["E06-audit-command/T003-snapshot"]
---

# T004: AI Reconcile (Semantic Review)

## Acceptance Criteria

- Generates a prompt file for human to give to their AI
- Prompt includes: snapshot, epic E##-Detail.md, changed files list, AGENTS.md Code Style rules
- Prompt instructs AI to check against 10 Code Style rules
- After human runs AI, parses findings for "Must Fix Before Close" section
- Fails audit if "Must Fix Before Close" has any items

## Implementation Plan

- [ ] Add `internal/audit/reconcile.go`
- [ ] Implement `GeneratePrompt(epic, release) (string, error)`
- [ ] Embed 10 Code Style rules in prompt template
- [ ] Include snapshot.md content in prompt
- [ ] Include epic E##-Detail.md in prompt
- [ ] Include changed files list in prompt
- [ ] Implement `ParseFindings(output) (hasIssues bool, err)`
- [ ] Check for "Must Fix Before Close" header and non-empty content
- [ ] Test prompt generation and finding parsing
- [ ] Run `make build && make test`