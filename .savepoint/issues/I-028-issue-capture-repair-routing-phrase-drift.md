---
id: I-028
title: Shared issue-capture reference drifted from its own doc-consistency test
type: drift
status: resolved
source:
  kind: report
  actor: {role: executor, session: v2-chat}
  at: '2026-09-22T00:00:00Z'
severity: low
resolution:
  disposition: accepted
  actor: {role: owner, session: user}
  at: '2026-09-22T00:00:04Z'
  reason: Owner accepted the recorded repair evidence and passing internal/init tests as sufficient. Owner decision, not a Check-verified closure.
history:
  - at: '2026-09-22T00:00:00Z'
    actor: {role: executor, session: v2-chat}
    kind: observed
    note: While repairing I-019, `make test` failed independently on TestSharedIssueCaptureRoleBoundariesAndRepairRouting (internal/init/agent_skills_test.go:1181-1191, added in ca738a9 on 2026-09-19) — it asserts agent-skills/references/issue-capture.md contains the phrase "becomes a new, bounded Task in an Objective", which the shipped file's Out-Of-Scope Repair section never stated (it said "Escalate to a new Task in an Objective"). Confirmed via `git stash` that this failure predates and is unrelated to the I-019 change.
  - at: '2026-09-22T00:00:00Z'
    actor: {role: executor, session: v2-chat}
    kind: repair_attempted
    note: Reworded the Out-Of-Scope Repair section in agent-skills/references/issue-capture.md (and its byte-identical templates/project-v2 copy, TPL-01) so the escalation sentence reads " Escalate only when the repair itself needs planning — an open Design decision, or work spanning multiple Objectives — where the repair becomes a new, bounded Task in an Objective instead of an inline edit." on one unwrapped line so the exact phrase the test checks for is a contiguous substring. `go test ./internal/init/...` now passes; `go build ./...` and `go vet ./...` pass.
  - at: '2026-09-22T00:00:04Z'
    actor: {role: owner, session: user}
    kind: owner_decision
    note: Owner accepted I-028 as resolved based on the recorded repair and passing test evidence.
---

# I-028: Shared issue-capture reference drifted from its own doc-consistency test

## Summary

`agent-skills/references/issue-capture.md`'s "Out-Of-Scope Repair" section
and its own doc-consistency test in `internal/init/agent_skills_test.go`
(`TestSharedIssueCaptureRoleBoundariesAndRepairRouting`) had drifted apart:
the test required the phrase "becomes a new, bounded Task in an Objective",
but the shipped file said "Escalate to a new Task in an Objective", so
`make test` failed independently of any of that day's actual work.

## Evidence

- `internal/init/agent_skills_test.go:1181-1191` (added in commit `ca738a9`,
  2026-09-19) asserts the phrase against both the live and
  `templates/project-v2` scaffold copies of `issue-capture.md`.
- `git stash` on the `v2` branch at `9ece207` reproduced the same two
  failures with no other changes present, confirming the drift predates and
  is independent of the I-019 repair.

## Proof Needed

- `go test ./internal/init/...` passes `TestSharedIssueCaptureRoleBoundariesAndRepairRouting`
  for both the live and template copies. (Done — see history.)
- `agent-skills/references/issue-capture.md` and its
  `templates/project-v2` copy stay byte-identical per TPL-01. (Done — verified with `diff`.)
