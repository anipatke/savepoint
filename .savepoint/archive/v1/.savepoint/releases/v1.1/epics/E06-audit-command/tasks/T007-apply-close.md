---
id: E06-audit-command/T007-apply-close
status: done
objective: "Agent applies audit proposals and closes epic with status: audited"
depends_on:
    - E06-audit-command/T005-proposals
    - E06-audit-command/T014-tab-indicator
---

# T007: Apply + Close

## Acceptance Criteria

- [ ] Agent reads approved proposals from E##-Audit.md
- [ ] Agent applies each proposal using Replace/Insert/Delete operations
- [ ] Agent marks epic E##-Detail.md with `status: audited`
- [ ] Agent updates Design.md `last_audited` field
- [ ] Agent updates router to advance to next epic
- [ ] Shows summary: X applied, findings closed

## Implementation Plan

- [ ] In audit skill (SKILL.md), add step after proposals approved:
  - Read E##-Audit.md proposals section
  - Apply each proposal (Replace/Insert/Delete) to live files
  - Update E##-Detail.md frontmatter: `status: audited`
  - Update Design.md `last_audited: {release}/{epic}`
  - Update router.md to set state to next epic (or `epic-design` for new epic)
  - Show apply summary output
- [ ] Add apply helper functions to data module:
  - `ApplyProposal(target, old, new) error`
  - `UpdateEpicStatus(epic, status) error`
  - `UpdateLastAudited(release, epic) error`
- [ ] Test apply workflow
- [ ] Run quality gates