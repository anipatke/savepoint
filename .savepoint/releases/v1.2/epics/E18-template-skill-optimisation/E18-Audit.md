---
type: audit-findings
audited: 2026-05-11
---

# Audit Findings: E18 Template and Skill Optimisation

## Main Findings

All four tasks delivered their acceptance criteria cleanly. No drift notes filed across T001–T004.

**T001 — Canonical Guides:** AGENTS.md carries the explicit statement that "The phase skill is the canonical workflow source." Phase-by-phase narration removed from AGENTS and delegated to skills. Terminology section renamed from "Task Status" to "Terminology"; scoped qualifiers (`Router state`, `Task status`, `Task stage`) eliminate prior ambiguity. All six root skills and their six scaffolded copies are byte-for-byte identical — verified by diff. Skills reduced by 41–65% in word count through removal of duplicated AGENTS content, not just compression.

**T002 — Prompt Pruning:** Seven redundant phase prompt files deleted. `templates/prompts/` now contains only `magic-prompt.prompt.md`. No runtime code path referenced the deleted filenames. Tests updated before deletion; suite stayed green throughout.

**T003 — Template Surface Tests:** Three new test functions protect the reduced surface:
- `TestProjectGuidanceTemplatesMirrorLiveGuidance` — asserts canonical phrases appear in both live and scaffolded AGENTS, and that every root skill has an identical scaffolded copy with no extra scaffold-only entries.
- `TestProjectTemplatesRejectStaleWorkflowTerms` — rejects stale status values (`todo`, `doing`, `blocked`, `review`, `audit`) and prompt-based phase instruction patterns.
- `TestPromptTemplates_onlyMagicPromptRemains` / `TestPromptTemplates_magicPromptIsBootstrapOnly` — assert exactly one prompt file exists and that it contains no phase-state instructions.

**T004 — Artifact Templates:** Explicit markdown template blocks added to three skills that previously described artifact structure only in prose. `savepoint-audit` now carries the canonical `E##-Audit.md` template with exact frontmatter (`type: audit-findings`, `audited: {date}`), a `- [ ]` checkbox Code Style Review, and `### Target File / ### Replace / ### With` block structure. `savepoint-create-task` now carries the canonical `T###-slug.md` template with exact frontmatter (`id`, `status: planned`, `objective`, `depends_on`) and all required sections in order. `savepoint-system-design` now carries the canonical `E##-Detail.md` template with exact frontmatter (`type: epic-design`, `status: planned`) and all required sections. Root and scaffolded copies confirmed identical via `Compare-Object` after edits. `make build && make test` passed. This closes the format consistency gap seen in E18's own audit file diverging from E17's established structure.

**Design.md updated:** Section 1 routing description corrected from `template prompts` to `phase skills`. `last_audited` advanced to `v1.2/E18-template-skill-optimisation`.

## Code Style Review

- [x] One job per file
- [x] One job per function
- [x] Test branches
- [x] Types document intent
- [x] Build only what is needed
- [x] Handle errors at boundaries
- [x] One source of truth
- [x] Comments explain WHY
- [x] Content in data files
- [x] Small diffs

## Proposed Changes

### Target File
`.savepoint/Design.md`

### Replace
```md
- **Agent routing:** AGENTS.md → `.savepoint/router.md` → template prompts. See AGENTS.md Workflow section.
```

### With
```md
- **Agent routing:** AGENTS.md → `.savepoint/router.md` → phase skills. See AGENTS.md Workflow section.
```
