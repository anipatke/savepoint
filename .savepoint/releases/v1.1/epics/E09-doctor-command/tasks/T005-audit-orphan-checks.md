---
id: E09-doctor-command/T005-audit-orphan-checks
status: done
objective: "Validate audit state and detect orphaned tasks"
depends_on: ["E09-doctor-command/T004-dependency-checks"]
---

# T005: Audit + Orphan Checks

## Acceptance Criteria

- Audit proposals without matching audit-pending state detected
- Orphaned tasks detected (tasks in nonexistent epics)
- Orphaned tasks reported with suggestion to move to .savepoint/orphans/
- Audit log shows previous audit results

## Implementation Plan

- [x] Add to `internal/doctor/checks.go`
- [x] Implement `CheckAuditState(root) []Problem`
- [x] Find all audit proposal directories
- [x] Check router state is audit-pending for corresponding epic
- [x] Warn if proposals exist without audit-pending state
- [x] Implement `CheckOrphans(root) []Problem`
- [x] Find tasks whose epic directory doesn't exist
- [x] Report orphaned tasks with move suggestion
- [x] Test audit and orphan detection
- [x] Run `make build && make test`