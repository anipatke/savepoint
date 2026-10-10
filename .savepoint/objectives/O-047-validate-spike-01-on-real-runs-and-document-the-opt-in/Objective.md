---
id: O-047
title: Validate Spike 01 on real runs and document the opt-in
status: planned
depends_on: [O-046]
release: G-003
priority: medium
rank: 3
---

# O-047: Validate Spike 01 on real runs and document the opt-in

## Outcome

Spike 01's acceptance criteria are shown to hold with a real Claude Code session on WSL2 and on native Windows, including cancel, timeout, runner loss, denials, and closing the board mid-run, with the evidence recorded in this Objective. Usage per run is compared with a hand-started session on the same Next line. The README, this repository's and the scaffold's AGENTS.md, the config template comment, and `.savepoint/Design.md` describe the feature as opt-in and name the paste workflow as the default.

## Why

The review's final test is "if any of these fail, the direction is wrong, not the implementation". That judgement needs recorded evidence from real runs, not a passing unit suite. The documentation change is small but is where the opt-in promise becomes visible to users.

## Success Conditions

1. A table in the evidence, one row per acceptance criterion 1–9 of review §5 (`docs/reviews/agent-orchestration-pass-1.md`), with platform, command, time, and result; every row passes on both platforms or names an Issue.
2. A chaos run (`kill -9` of the runner mid-run) shows `runner_lost` on the board and in `doctor`.
3. Usage comparison recorded: provider-reported input and output tokens for one board-started and one hand-started run of the same Next line, with the difference explained.
4. README and both AGENTS.md copies state that orchestration is off by default and that pasting the Next line into any terminal remains fully supported; the config template's `agent_launcher` comment matches the shipped flags; `Design.md` Codebase Map gains `internal/run` and §6 gains the subcommand.
5. Existing projects receive the guidance through `savepoint upgrade-assets` with user-edited files preserved (FS-02, TPL-04).
6. Configured gates pass, including `make test-full`.

## Architectural Considerations

- Validation is evidence gathering, not a Check: the validation Task should be run by the owner or a fresh session, not the executor of O-046, so the evidence is independent.
- Defects found become Issues routed through the normal repair path; this Objective does not fix them.
- Tasks are detailed only after O-046 lands and this Objective's design is confirmed.

## Boundaries

**In scope:**
- Validation runs and evidence tables; Issues for any failure
- README, AGENTS.md (both copies), config template, Design.md; scaffold parity and content tests

**Out of scope:**
- Fixing defects found; new features; later spikes
