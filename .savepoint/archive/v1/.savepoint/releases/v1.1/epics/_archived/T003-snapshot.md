---
id: E06-audit-command/T003-snapshot
status: planned
objective: "Generate file tree and changed-files snapshot"
depends_on: ["E06-audit-command/T002-quality-gates"]
---

# T003: Snapshot Generation

## Acceptance Criteria

- Generates `.savepoint/audit/{release}/{epic}/snapshot.md`
- Includes file tree (gitignore-respecting)
- Includes list of changed files (from git)
- Does NOT include file contents (no code in snapshot)
- Snapshot path matches expected format

## Implementation Plan

- [ ] Add `internal/audit/snapshot.go`
- [ ] Implement `GenerateSnapshot(epic, release) error`
- [ ] Walk directory tree respecting .gitignore
- [ ] Run `git diff --name-only` for changed files
- [ ] Write snapshot.md with tree + changed-files sections
- [ ] Create audit directory if not exists
- [ ] Test snapshot output format
- [ ] Run `make build && make test`