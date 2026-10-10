---
id: T-114
title: Drop migration history and use one term per role
objective: O-039
status: done
depends_on: [{task: T-113, requires: clear}]
complexity_tier: low
complexity_reason: Wording edits in the same parity-checked assets, with content tests.
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: plan-o039-20261006}
check_waiver:
    task: T-114
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-10T06:53:04Z"
---

# Drop migration history and use one term per role

## Outcome

Active skills, references, and the managed block no longer narrate V1/V2 history outside the managed block's Required Goal Context. Prose says "owner" for the human decision-maker and "Goal" for the record, and `release:` appears only as the field name.

## User Check

Search the skills and references for "V1", "resurrect", "user may", and "Release" in prose, and find no history narrative or mixed terms.

## Done When

- Each removed history passage is listed in Technical Evidence and is either covered by Required Goal Context or no longer needed by an agent.
- "user" never names the decision-maker where "owner" is meant, and "Release" never names a Goal in prose. Field names and R-### identities are unchanged.
- Live and scaffold copies are byte-identical (TPL-01). Content tests are updated, with one assertion against the old wording returning.
- O-034's scenario walkthroughs are re-walked, with the result recorded per scenario.
- Token weight before and after is recorded.
- `make build && make test-fast` passes.

## Context Files

`templates/project-v2/AGENTS.md`; `AGENTS.md`; the four workflow skills and three references under `agent-skills/` and their copies under `templates/project-v2/agent-skills/`; `internal/init/agent_skills_test.go`; `internal/init/skill_validation_test.go`; `internal/init/template_freshness_test.go`; `.savepoint/objectives/O-034-skills-optimisation/tasks/T-099-review-the-skills-against-what-the-code-health-goal-taught-us.md`.

## Design References

O-039 Confirmed Design; AGENTS.md Terminology.

## Guardrails

TPL-01, TPL-02, TPL-04, FS-02, POL-01, POL-02, TEST-01, TEST-02, TEST-03, TEST-08, STYLE-07, STYLE-10.

## Implementation Plan

1. Record the before token weight.
2. Remove history passages and settle terms file by file, mirroring each change to the scaffold copy.
3. Update content tests.
4. Re-walk the scenarios, record the after weight, and run the gate.

## Boundaries

No rule changes, no field or identity renames, and no Go runtime changes.

## Technical Verification

Focused tests during iteration; handoff gate from AGENTS.md's Verification Policy.

## Technical Evidence

Executor session: T-114 start, 2026-10-10. Not a Check; claims for a fresh `savepoint-check` to verify.

**Per-criterion outcomes**

- History removed: the "does not resurrect a separate defect record, status, or public phase in V2" clause (issue-capture) is dropped; the owner needs nothing from it. "V2" qualifiers ("the V2 index", "V2 record", "V2 Task Markdown", "existing V2 project", "the V2 workflow", "the V2 higher-level integration gate") are removed from the four skills, references and the managed block; the sentences keep their meaning. The V1/migration passages inside the managed block's Required Goal Context are kept, as the Outcome allows. The Legacy Health-Check.md mapping in `commands-and-procedures.md` stays: it is a live procedure for migrated projects, not narrative.
- Terms: "Only the user may set a Task to `status: done`" and "Prompt the user" now say owner in the managed block; "The user's stated intent" (idea Read list) and "a word the user says" (issue-capture) now say owner. Kept on purpose: the Idea `User` section and `## User Check` headings (the product's user, and the task's own check), "user-authored files" (files, not a decision-maker), `release:` field names, `Release.md`, R-### identities, and "release gates" (a kind of gate, not a Goal).
- Live and scaffold copies byte-identical (`diff -r` differs only by live-only `bubbletea-tui-design`); this repository's managed block refreshed with `./savepoint upgrade-assets` (AGENTS.md only).
- Tests: updated V2-wording assertions in `agent_skills_test.go`, `template_freshness_test.go`, `v2_scaffold_test.go`; the "resurrect" assertion became an absence check; new `TestActiveGuidanceKeepsNoHistoryAndOneTermPerRole` fails if the old wording or "user" as decision-maker returns.
- Token weight (managed block + four skills + three references; words / bytes): before 13297 / 89913, after 13271 / 89794. Counted with `wc`.
- `make build && make test-fast` passes (exit 0); `savepoint doctor` ALL CLEAN.

**O-034 scenario re-walk**

1. Next to skill routing: unchanged text; works as before.
2. REPLAN REQUIRED re-entry: unchanged; correct.
3. Waived Task Check to Full Objective Check: only "V2" removed from a gate name; same behaviour.
4. Issue-only repair: only the defect-mapping sentence trimmed; repair and router handoff unchanged.
5. Advisory Code Health report: unchanged; correct.

Limitations: re-walk is by reading the revised text, no live agent run. Term checks cover the named phrases, not every paraphrase.

## Drift Notes

None yet.
