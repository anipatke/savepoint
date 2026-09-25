---
id: I-056
title: Scaffold AGENTS.md does not say the owner can resolve and reopen Issues from the board
type: drift
status: resolved
source:
  kind: check
  check: C-930
  actor: {role: checker, session: o015-objective-check-20260925}
  at: '2026-09-25T09:30:30Z'
tasks: [T-046]
checks: [C-930, C-931, C-932]
guardrail_ids: [TPL-02]
severity: low
resolution:
  disposition: verified
  check: C-932
  actor: {role: checker, session: o015-surgical-recheck-20260925}
  at: '2026-09-25T10:14:41Z'
  reason: The four scaffold templates/project-v2/AGENTS.md passages match the repository AGENTS.md wording (owner Space/Backspace authority, Issue Capture resolution roles, the task-status stop rule, and the Check section's Issues panel wording), and the Full Objective gate (make build && make test-full) passed.
history:
  - at: '2026-09-25T09:30:30Z'
    actor: {role: checker, session: o015-objective-check-20260925}
    kind: observed
    check: C-930
    note: The repository AGENTS.md gained the board Space/Backspace Issue wording, but the scaffolded templates/project-v2/AGENTS.md shipped to new projects did not.
  - at: '2026-09-25T09:40:04Z'
    actor: {role: executor, session: o015-issue-repair-20260925}
    kind: repair_attempted
    note: Updated the four scaffold AGENTS.md passages to document Space resolution as accepted and Backspace reopening, matching the repository guide. make build and make test-fast passed; Issue remains open for independent verification.
  - at: '2026-09-25T09:44:54Z'
    actor: {role: executor, session: o015-issue-repair-20260925}
    kind: repair_attempted
    note: Restored the directly-instructed owner-decision rule, rewrote the Issue Capture bullet to remove the checker-only contradiction, and aligned wording to the Issues panel. make build and make test-fast passed; Issue remains open for independent verification.
  - at: '2026-09-25T09:46:30Z'
    actor: {role: executor, session: o015-issue-repair-20260925}
    kind: repair_attempted
    note: Matched all four scaffold passages to the repository wording, including the directly-instructed owner-decision rule, the checker/owner/planner Issue Capture sentence, and Issues panel terminology. make build and make test-fast passed; Issue remains open for independent verification.
  - at: '2026-09-25T09:58:00Z'
    actor: {role: checker, session: o015-objective-recheck-20260925}
    kind: rechecked
    check: C-931
    note: The four scaffold AGENTS.md passages now match the repository AGENTS.md, and make test-full passed. The repair is proven, but C-931 is NEEDS WORK because of I-057, so this Issue stays open until a CLEAR Check can close it as verified.
  - at: '2026-09-25T10:14:41Z'
    actor: {role: checker, session: o015-surgical-recheck-20260925}
    kind: rechecked
    check: C-932
    note: Re-confirmed the four scaffold templates/project-v2/AGENTS.md passages against the current repository AGENTS.md by exact-text diff; still identical. C-932 is CLEAR, so this Issue is resolved as verified.
---

# I-056: Scaffold AGENTS.md does not say the owner can resolve and reopen Issues from the board

## Summary

O-015's last Success Condition says "AGENTS.md, the active skills/references,
and their V2 scaffold copies say the owner may resolve (as `accepted`) and
reopen Issues from the board." The repository `AGENTS.md` and every
`agent-skills/` file and its scaffold copy were updated. The scaffold
`templates/project-v2/AGENTS.md`, which `savepoint init` ships to new
projects, was not. T-046's Done When narrowed the scaffold copies to
`templates/project-v2/agent-skills/`, so the Task evidence did not catch it.

## Evidence

- `AGENTS.md:104`, `:112`, `:124` name Space (resolve as `accepted`) and
  Backspace (reopen) in the board's Issues panel.
- `templates/project-v2/AGENTS.md:101`, `:109`, `:115`, `:120` still carry the
  pre-O-015 wording: the owner may close an Issue as `accepted` "through an
  explicit decision", with no mention of the board or of reopening.
  `templates/project-v2/AGENTS.md:109` also says "only `savepoint-check`
  verifies the proof and closes it", which the board's Space key now
  contradicts in shipped guidance (TPL-02).
- `git diff --stat -- templates/project-v2/AGENTS.md` is empty.

## Repair Attempt

Updated the four Issue-authority passages in `templates/project-v2/AGENTS.md`
to document owner resolution as `accepted` with Space and reopening with
Backspace, matching the repository's `AGENTS.md`. The wording was reviewed
directly, and `make build && make test-fast` passed. This Issue remains open
for independent verification; no resolution disposition is recorded.

### Follow-up Repair

Restored the scaffold's rule that an agent may record an owner decision only
when directly instructed. Rewrote the Issue Capture bullet to match the
repository's checker/owner/planner wording without saying only a checker
verifies and closes Issues. The four passages now use “Issues panel” wording
matching `AGENTS.md`. `make build && make test-fast` passed again. The Issue
remains open for independent verification.

### Final Wording Alignment

Copied the repository's exact wording for the four relevant guidance
passages: owner action and direct-instruction authority, Issue Capture
resolution roles, the task status stop rule, and the Check section's Issues
panel wording. `make build && make test-fast` passed after this correction.
The Issue remains open for independent verification.

## Proof Needed

Update the four scaffold `AGENTS.md` passages to match the repository
`AGENTS.md` wording about owner resolution and reopening from the board. No
code change. A recheck compares the passages and reruns the gate.
