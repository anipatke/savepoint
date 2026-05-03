---
id: E06-audit-command/T005-proposals
status: done
objective: Update proposals to write to E##-Audit.md format
depends_on:
    - E06-audit-command/T010-audit-file-migration
---

# T005: Proposals to E##-Audit.md

## Acceptance Criteria

- [ ] Proposals written to `{epic}/E##-Audit.md` instead of `.savepoint/audit/`
- [ ] Includes Main Findings section from proposals.md + snapshot.md
- [ ] Includes Code Style Review checklist section
- [ ] Format: frontmatter + Main Findings + Code Style Review

## Implementation Plan

- [x] Update proposal file format:
  ```
  ---
  type: audit-findings
  audited: YYYY-MM-DD
  ---
  # Audit Findings: E## {Epic Name}

  ## Main Findings
  [content from proposals.md + snapshot.md]

  ## Code Style Review
  - [ ] One job per file
  - [ ] One-sentence functions
  - [ ] Test branches
  - [ ] Types are documentation
  - [ ] Build, don't speculate
  - [ ] Errors at boundaries
  - [ ] One source of truth
  - [ ] Comments explain WHY
  - [ ] Content in data files
  - [ ] Small diffs
  ```
- [x] Update skill to write in this format
- [x] Verify new format works