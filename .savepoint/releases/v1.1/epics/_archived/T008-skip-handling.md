---
id: E06-audit-command/T008-skip-handling
status: planned
objective: "Implement audit skip with logging"
depends_on: ["E06-audit-command/T001-cli-entrypoint"]
---

# T008: Skip Handling

## Acceptance Criteria

- `--skip --reason "..."` flags are required together
- Missing reason shows error
- Skipped audit logged to `.savepoint/audit-log.md`
- Log entry includes: timestamp, epic, release, reason, skipped-by
- Skipped epic still advances router (no changes to docs)

## Implementation Plan

- [ ] Add `internal/audit/log.go`
- [ ] Implement `LogSkip(epic, release, reason) error`
- [ ] Parse existing audit-log.md or create new
- [ ] Append new entry in format:
    ```markdown
    ## {timestamp}

    - **Epic:** {epic}
    - **Release:** {release}
    - **Reason:** {user reason}
    - **Status:** ⚠ skipped
    ```
- [ ] Update router to advance past skipped epic
- [ ] Test skip flow with valid/invalid args
- [ ] Run `make build && make test`