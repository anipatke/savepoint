---
id: E47-onboarding-upgrades/T007-prove-the-whole-onboarding-loop-protects-user-files
title: Prove the whole onboarding loop protects user files
status: planned
objective: Assert the epic's ownership contract end to end across both trees — edited skills, marked and unmarked guides, absent optional files, repeat runs, and write-free previews.
depends_on:
    - E47-onboarding-upgrades/T003-begin-from-one-rough-sentence
    - E47-onboarding-upgrades/T004-adopt-savepoint-into-an-existing-codebase
    - E47-onboarding-upgrades/T006-retire-the-old-instructions-without-losing-edits
complexity_tier: medium
complexity_reason: A cross-cutting matrix over two trees whose findings may expose real gaps in the preceding tasks.
---

# T007: Prove the whole onboarding loop protects user files

## Problem

Each preceding task proves its own piece. What none of them proves is the property the epic actually promises: that across the whole loop — scaffold a project, adopt it into a codebase, upgrade it, retire what it replaced, upgrade it again — no byte a user wrote is lost without a recoverable copy and an entry in the report naming where it went.

That property is only visible as a set. The ownership rules that make it true are spread across three mechanisms: manifest provenance decides whether a skill was customized, the marker pair decides whether an agent guide is managed, and the install-if-missing allowlist decides whether a file is ever touched at all. Each is tested in isolation today, and each is now exercised against a second tree it has never seen. A rule that silently stops applying to the V2 tree fails no existing test.

The gaps this is looking for are specific. A V2 skill edited by the user and then conflicted — does it keep the user's bytes and get its incoming copy beside it, or does the untracked-path branch back it up and replace it because the manifest carries V1 entries from before migration? An unmarked `AGENTS.md` in a migrated project — does it still conflict rather than being adopted? A retirement archive and a conflict sidecar landing in the same run — does the report name both? A project with no Health-Check, no Concept, and no procedures file — does anything in the new guidance treat their absence as a problem?

Where this finds a real gap, the fix belongs in the mechanism, not in the assertion. A matrix adjusted until it passes proves only that it was adjusted.

## Context Files

- `internal/init/lifecycle_test.go`
- `internal/init/upgrade.go`
- `internal/init/upgrade_test.go`
- `internal/init/upgrade_failure_test.go`
- `internal/init/integration_test.go`
- `internal/init/manifest.go`
- `internal/init/agents.go`
- `internal/init/retire_v1_skills.go`
- `internal/init/scaffold.go`
- `templates/project-v2/AGENTS.md`
- `agent-skills/references/commands-and-procedures.md`
- `.savepoint/Guardrails.md`

## Acceptance Criteria

- [ ] A lifecycle matrix runs the full loop — init, upgrade, retire, upgrade again — against a fresh V2 project, a migrated V2 project, and a project still at schema version 1, asserting the resulting actions and bytes for each.
- [ ] An edited V2 skill conflicts: the user's bytes are unchanged, the incoming version is written to the `.new` sidecar, and the report names the path and the sidecar.
- [ ] A V2 skill whose manifest entry is absent — the state a migrated project is in — is backed up to `.bak` before replacement, and the backup is byte-identical to what the user had.
- [ ] `--force` adopts a conflicted V2 skill and a conflicted unmarked guide, in each case leaving a recoverable backup named in the report.
- [ ] A marked `AGENTS.md` refreshes only between the markers; every byte outside them is unchanged, including content below the closing marker.
- [ ] An unmarked `AGENTS.md` conflicts, stays byte-identical, and receives the merged result as a `.new` sidecar; a half-marked guide behaves the same way.
- [ ] An agent guide under a casing variant keeps its on-disk name for both the write and every sidecar.
- [ ] A project with no `Health-Check.md`, no Concept, and no procedures file completes every step with no entry, warning, or diagnostic about their absence (TPL-03).
- [ ] Every step of the loop is idempotent: a second run changes no file content, no mtime, and no manifest entry, on all three project kinds.
- [ ] A dry run of every step reports the same actions as the real run and writes nothing, verified by snapshotting bytes and mtimes across the whole project directory.
- [ ] A run that produces both a retirement archive and a conflict sidecar reports both, with each path and its recovery location named.
- [ ] Any gap the matrix exposes is fixed in the mechanism that owns the rule, with the change and its reason recorded in the Context Log; no assertion is weakened to make the matrix pass.
- [ ] `make build && make test` passes.

## Implementation Plan

- [ ] Build the three project fixtures as helpers: fresh V2, migrated V2 with a V1-era manifest and V1 skills present, and legacy V1.
- [ ] Write the matrix as a table of project kind, step, and expected actions, so a missing combination is visible as a missing row.
- [ ] Add a whole-directory snapshot helper recording path, content hash, and mtime, and assert against it for every dry run and every second run.
- [ ] Run the matrix and triage each failure as either a wrong expectation or a real gap before changing anything.
- [ ] Fix real gaps in `upgrade.go`, `manifest.go`, `agents.go`, or `retire_v1_skills.go`, recording each in the Context Log.
- [ ] Add the combined retirement-and-conflict reporting case.
- [ ] Add the absent-optional-files case across both trees.
- [ ] Run `go test ./internal/init/...`, then `make build && make test`.

## Context Log

Pending.
