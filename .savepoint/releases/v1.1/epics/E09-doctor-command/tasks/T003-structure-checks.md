---
id: E09-doctor-command/T003-structure-checks
status: done
objective: "Validate release/epic/task structure and YAML"
depends_on: ["E09-doctor-command/T002-config-router-validation"]
---

# T003: Structure Checks

## Acceptance Criteria

- Release directories exist and contain release PRD
- Epic directories exist and contain E##-Detail.md
- Task files are valid markdown with YAML frontmatter
- Corrupt YAML/frontmatter detected with file:line info
- Task acceptance criteria present and non-empty

## Implementation Plan

- [x] Add to `internal/doctor/checks.go`
- [x] Implement `CheckStructure(root, epicFilter)` returning `[]Problem`
- [x] Walk release directories, validate release PRD
- [x] Walk epic directories, validate E##-Detail.md
- [x] Walk task files, parse frontmatter
- [x] Validate each task has: id, status, objective, depends_on
- [x] Check acceptance criteria present in task body
- [x] Detect corrupt YAML with file:line location
- [x] Test structure validation on valid and invalid projects
- [x] Run `make build && make test`