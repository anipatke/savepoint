---
id: I-080
title: Packaged skills hard-code Savepoint's own make targets and Go paths
type: defect
status: in_progress
source:
  kind: report
  actor: {role: owner, session: user}
  at: '2026-09-26T05:18:49Z'
severity: medium
history:
  - at: '2026-09-26T05:18:49Z'
    actor: {role: owner, session: user}
    kind: observed
    note: >-
      Raised by an independent review of the packaged Savepoint skills; claim verified against the code in a follow-up review session before capture.
  - at: '2026-09-26T05:25:10Z'
    actor: {role: executor, session: skill-review-fixes}
    kind: repair_attempted
    note: >-
      Packaged skills, check-method, commands-and-procedures, and the scaffold AGENTS.md now name gates through config.yml quality_gates or the project's AGENTS.md; the Context Files example uses language-neutral paths; Savepoint's own make gates remain only in the repo AGENTS.md. Packaged template copies re-synced; make build and make
      test-fast passed.
---

# I-080: Packaged skills hard-code Savepoint's own make targets and Go paths

## Summary

The skills shipped by `savepoint init` name this repository's Make targets
as the handoff and Check gates, and use this repository's Go files as the
worked Context Files example. A downstream project in another language gets
gates it cannot run. It also contradicts the skills' own rule that
`config.yml` `quality_gates` is the single place project commands live.

## Evidence

- `agent-skills/` and `templates/project-v2/agent-skills/` are byte-identical
  apart from the repo-local `bubbletea-tui-design`.
- `savepoint-design/SKILL.md:74,232`, `savepoint-task/SKILL.md:91-94`,
  `savepoint-check/SKILL.md:52`, `references/check-method.md:268` name
  `make test-focused`, `make build`, `make test-fast`, `make test-full`,
  `make ci`, or `make test`.
- `savepoint-design/SKILL.md:207` lists `cmd/board.go`,
  `internal/data/project.go`, `internal/resume/resume.go` and others as the
  Context Files example.
- `references/commands-and-procedures.md:14-20` names `quality_gates` as the
  single home for project commands.

## Proof Needed

- Packaged skills name gates through `.savepoint/config.yml`
  `quality_gates`, not Make targets, and use a language-neutral path example.
- Savepoint's own Make conventions live in repo-only guidance, not in
  packaged text.
- A freshly initialized project's skills contain no Savepoint-repo commands
  or paths; packaged copies stay in parity.

## Repair Attempt Evidence

- No `make` target, `cmd/`, or `internal/` path remains under `templates/project-v2/`.
- The scaffold AGENTS.md Build section is now a place for the project to list its own gates.
- `TestSavepointTaskSkillEvidenceRequirement` now asserts `configured build and test gates`.
- Packaged copies are byte-identical to `agent-skills/`; `make build` and `make test-fast` passed.
