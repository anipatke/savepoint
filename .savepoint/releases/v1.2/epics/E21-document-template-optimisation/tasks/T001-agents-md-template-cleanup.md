---
id: E21-document-template-optimisation/T001-agents-md-template-cleanup
status: planned
objective: Audit and tighten the scaffolded AGENTS.md template so it reflects the E20 lifecycle contract and contains no stale or duplicate guidance.
complexity_tier: low
complexity_reason: Single-file edit with a clear target; alignment check against live AGENTS.md.
---

# T001: AGENTS.md Template Cleanup

## Problem

The scaffolded `templates/project/AGENTS.md` was last meaningfully updated in E20 alongside the live file, but may still carry sections that duplicate phase-skill prompts, reference stale terminology (`phase`, `stage: implementation`), or include guidance that belongs in skill files rather than the entry-point agent guide.

## Context Files

- `templates/project/AGENTS.md`
- `AGENTS.md`
- `internal/init/template_freshness_test.go`

## Acceptance Criteria

- [ ] Scaffolded `AGENTS.md` contains no references to stale lifecycle terms (`phase`, `stage: implementation`, legacy status values).
- [ ] Guidance that duplicates phase-skill prompt instructions is removed; the file routes to skills, it does not re-state them.
- [ ] Live `AGENTS.md` and scaffolded `AGENTS.md` are consistent on all lifecycle and terminology guidance.
- [ ] Template freshness test passes without modification (or is updated only to reflect intentional structural changes).
- [ ] `make build && make test` passes.

## Implementation Plan

- [ ] Read `templates/project/AGENTS.md` and `AGENTS.md` side by side; list every divergence and every section that duplicates skill-level prompt content.
- [ ] Remove or condense duplicate and stale sections from the scaffolded template.
- [ ] Mirror any lifecycle or terminology corrections to the live `AGENTS.md` where applicable.
- [ ] Run `go test ./internal/init` to verify freshness tests.
- [ ] Run `make build && make test`.
