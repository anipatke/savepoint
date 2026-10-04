---
id: O-033
title: Show which Tasks can safely run in parallel worktrees
status: done
depends_on: [O-037]
release: G-001
priority: medium
rank: 2
lanes:
  - key: core
    title: Lane / Proposed worktree — Planning and safety
  - key: board
    title: Lane / Proposed worktree — Board and session instructions
  - key: guidance
    title: Lane / Proposed worktree — Workflow and upgrades
  - key: integration
    title: Lane / Proposed worktree — Integration and documentation
---

# O-033: Show which Tasks can safely run in parallel worktrees

## Outcome

An owner can optionally see how Tasks within the selected Objective could run in parallel worktrees, with suggested lane groupings, understandable reasons and copyable session instructions. The owner can ignore every suggestion and execute Tasks normally.

## Why

Persisted planning advice helps owners coordinate parallel agents without making the default workflow more complicated or treating suggestions as execution policy.

## Success Conditions

- Parallel planning is off by default per project and can be switched on or off from an Advanced Options screen. The setting persists across restarts and is also editable in project configuration.
- Code Health remains a core feature. This Objective adds no Code Health visibility or enablement switch and changes none of its collection, evidence or verification behavior.
- When enabled and available, saved lane headings group Tasks within Planned, In Progress and Done; membership remains stable as Tasks move. Empty headings are hidden.
- Suggested opportunities stay within the selected Objective and use the same pure data projection on board and resume.
- Advice accounts for existing dependency gates, dependency paths, anticipated writes and write/read overlap. Shared reads are permitted; unexplained write/read overlap or uncertain scope withholds a recommendation, never an execution action.
- Lane membership and manifests are advisory. Choosing another order, ignoring a lane, running on main, using a different worktree, or changing anticipated file scope does not itself block any board action, skill, Task handoff or Check.
- Genuine dependencies, acceptance criteria, Guardrails, owner authority and the existing actual-worktree filesystem/identity rules remain applicable independently of this feature.
- Every recommended parallel Task has complete copyable optional session instructions. Missing or stale advice states its limits without requiring repair before normal execution.
- Disabling the feature restores the ordinary ungrouped board and removes optional recommendations/instructions without deleting lane plans, Tasks, configuration unrelated to this preference or recorded evidence.
- Existing projects without a setting or lane metadata remain valid. Settings writes preserve unrelated authored content and refuse stale-file conflicts.

## Architectural Considerations

`internal/data` owns optional persisted planning metadata, project preferences and a pure recommendation projection. The existing dependency and lifecycle gates do not consume lane compatibility as a new blocker. Optional recommendations are separate from Next selection and completion. Board persistence uses explicit commands; rendering performs no IO. Planner/executor/checker guidance may offer and assess advice but cannot turn it into an enforcement contract.

## Boundaries

**In scope:** integration with the preference and Advanced Options screen delivered by O-037; advisory lane/read/write metadata; conservative recommendation rules; grouping, reasons and optional instructions on board/resume; canonical/scaffold guidance; safe init/upgrade delivery; integration evidence and documentation.

**Out of scope:** Code Health toggles or behavior changes, generic feature/plugin frameworks, general configuration editing, automatic worktree/branch creation, monitoring actual lanes, claimed/running/merged states, automatic merge/conflict resolution, enforcing suggested lane membership or manifests, new lifecycle/completion gates, multi-user coordination, cloud services and general agent orchestration.

## Confirmed Design — Revised 2026-10-03

Owner-directed split: O-037 owns preference storage, Advanced Options and their standalone delivery. This Objective waits for O-037 and consumes its implemented interface. The settings description below is the downstream integration contract, not duplicate implementation ownership.

The owner confirmed the original design by asking to detail Tasks, subsequently requested opt-in behavior, and then explicitly clarified that Code Health remains core and all lane/worktree advice must be ignorable even while enabled. This revised contract supersedes the original mandatory scope-drift wording and the proposed Code Health toggle. Detailed planning is authorised; implementation and owner Task completion are not.

### Project preference and Advanced Options

Use a single optional boolean `features.parallel_planning` in `.savepoint/config.yml`, defaulting false when absent. A documented Settings action opens the existing board overlay pattern with one plain-language option: “Parallel planning — show suggested lanes and worktree instructions; you can ignore them.” Use standard keyboard navigation, toggle and close controls. Show save success only after persistence succeeds; stale or unwritable configuration produces an actionable error without falsely changing the displayed saved state. This is one real option in a small screen, not a registry or an empty list of speculative features.

A narrow data-owned writer updates only this preference, preserving unrelated keys, comments and authored content, checks source freshness and safely replaces the file. Board reload reads the persisted value. Resume and workflow skills use the same saved setting; non-TUI owners can edit it directly. No migration or automatic opt-in occurs, and existing owner choices survive upgrades. Code Health keeps its existing board surfaces, setup, providers, reports and verification behavior.

### Optional plan records

The Objective may declare named `lanes`, each with an Objective-local key and readable title. Task `lane`, `planned_reads` and `planned_writes` are optional; membership is derived from the Task reference, not a second list. No lane record claims that a worktree exists or changes task ownership. Exact project-relative paths distinguish shared implementation reads from anticipated writes; Context Files still govern executor reading under the existing skill policy. Manifests describe the plan rather than an execution allowlist. Omitted scopes are unknown; explicit empty lists are reviewed empty scopes. Conservative portable path comparisons cannot certify uncertain aliases.

Optional narrowly scoped Task-pair explanations can establish independence of exact read/write overlaps; they never explain away shared writes or real dependencies. Bind explanations to the reviewed inputs so changed inputs invalidate only that recommendation. Missing, malformed, unresolved or stale advisory data produces a visible planning diagnostic and no affected parallel recommendation; it must not prevent valid project loading or existing lifecycle actions. Existing required record/schema errors remain errors. Managed writes preserve optional authored data. No automatic backfill or required lane metadata for ordinary planning.

### Recommendation projection

Compute stable membership separately from current recommended opportunities, within one selected Objective. Reuse existing start/dependency decisions only to describe whether a Task is already eligible under the ordinary workflow; add no lane condition to those gates. Recommend one Task at a time within each proposed lane, respecting true dependencies, disjoint writes and no unexplained write/read overlap. Shared reads are allowed. Consider remaining lane scopes and recorded active work conservatively, and withhold advice when inputs are unknown. Report pairwise-compatible groups deterministically without claiming an optimal schedule or guaranteed independence.

Descriptions use advisory language: “Suggested together”, “Parallelism not established”, or an existing dependency reason. They do not say a shared-write Task is blocked solely because of a lane. No actual-worktree/merge monitoring or inferred ownership is added. Unknown scope, deviation from a manifest or multiple active Tasks in a suggested lane affects advice only.

### Board, resume and optional instructions

With the option enabled, Tasks with membership appear under the same “Lane / Proposed worktree” headings in every status column containing them. Hide empty headings; leave ungrouped Tasks as ordinary work. Stable ordering, card counts, keyboard selection, scrolling and existing badges remain correct. The Goal-wide view namespaces headings by Objective and offers no cross-Objective concurrency suggestions. Turning the option off restores the existing board without altering Task files or completion state. The current Next line always names selected work independently of advice.

Objective/Task details and resume/plain board provide the same reasons and complete text instructions from a shared formatter. Recommended instructions name the Task, Objective, suggested lane, existing Start line, anticipated scopes and prerequisites. They explicitly say the owner may ignore the proposed lane or worktree and follow normal execution. If choosing a real worktree, existing AGENTS worktree rules apply: prepare prerequisites, preserve the router, allocate no identity records there, record evidence in the existing Task, commit locally and leave merges/main-branch Checks to the owner. These rules are conditional on an actual worktree, not on a suggested heading. No automatic shell/branch/worktree action, guessed physical path or required clipboard support.

### Skills, verification and delivery

Planner guidance reads the preference and may offer/persist a parallel plan when useful. It may plan sequentially even while the option is on and must not demand lane data, worktree creation or adherence to recommendations. Executors retain their ordinary scope, extra-read and lifecycle responsibilities; lane deviations or manifest expansion alone do not require REPLAN REQUIRED. They may note differences so advice can be updated. Checkers assess recommendation accuracy and presentation as feature behavior; disregarding advice is not a project finding or clearance failure. Actual unmet acceptance criteria, constraints or dependencies remain governed by existing policy.

Canonical/scaffold guidance remains byte-identical. Init/upgrade behavior delivers new defaults and guidance while preserving authored assets, preferences and existing records. Verify feature-on/off/default behavior, data preservation, malformed advisory input degradation, stale settings writes, advice parity, stable grouping and navigation, and unchanged start/advance/completion/Check decisions when advice is ignored. No test needs network or real worktrees. Config/path work and final integration require fresh `make test-full`; ordinary presentation/guidance uses `make build && make test-fast`. The independent Full Objective Check applies check-method.md and obtains native windows-tests CI evidence produced by repository CI and supplied by the owner. Code Health collection remains the existing Full Check procedure only.

### Readiness evidence and handoff

O-032 is owner-completed. Targeted reads confirm typed record ownership and lifecycle/dependency gates in internal/data, a pure resume renderer, existing board overlays and command-based persistence. ConfigReader currently has no feature preference writer, so settings storage is explicit new work rather than an assumed interface. Existing Tasks were allocated through create-task and remain planned; revise them in place and allocate new settings work only through that command. Existing Context Files have been checked; intentionally new or predecessor-created files are named in each Task.

The prior detailed plan is superseded by this advisory/opt-in revision. No production code, owner waiver or Check is written during planning. Router remains selected on this Objective in design, with execution handoff pending owner plan review. Lanes in this plan are themselves optional suggestions; actual Task dependency requirements remain the source of ordering.

## Revised Planning Handoff Evidence — 2026-10-03

Existing planned Tasks were revised in place for optional, non-enforcing advice; the allocator created separate preference storage and Advanced Options Tasks. No Code Health toggle remains. The core sequence settles metadata, saved preference, then advice. Board settings precede grouped-column changes because they share model/update/view/load files; the workflow/scaffold branch can proceed independently after the advice interface. Instructions follow board grouping, and final integration follows instructions and workflow delivery. These lanes are suggestions only; the owner may execute the same dependency-respecting plan sequentially.

Verification specifically requires identical existing lifecycle/completion decisions with the feature on/off and when advice is ignored, plus graceful malformed-advice degradation and preference persistence/preservation. All Tasks remain planned; no production implementation or automatic feature opt-in occurred. Execution handoff remains pending owner plan review.

## Owner-Directed Ordering — 2026-10-03

Advanced Options is now the preceding Objective. Settings Tasks keep their identities under O-037; parallel Tasks remain here and wait for O-037 through the ordinary Objective dependency gate. Existing detailed parallel Tasks are retained at the owner's explicit request to reorder this plan, and must verify the delivered settings interface before execution. No implementation is started.
